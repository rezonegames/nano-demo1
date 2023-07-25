package game

import (
	"github.com/lonng/nano/component"
	"github.com/lonng/nano/scheduler"
	"github.com/lonng/nano/session"
	"sync"
	"tetris/config"
	"tetris/consts"
	"tetris/pkg/log"
	"tetris/pkg/z"
	"tetris/proto/proto"
	"time"
)

type RoomService struct {
	*component.Base
	serviceName string
	conf        *config.Room
	tables      map[string]*Table
	lock        sync.RWMutex
	queue       []*session.Session
}

func newRoomService(conf *config.Room) *RoomService {
	return &RoomService{
		serviceName: "r" + conf.RoomId,
		conf:        conf,
		tables:      make(map[string]*Table, 0),
	}
}

func (r *RoomService) AfterInit() {

	// 处理玩家断开连接
	session.Lifetime.OnClosed(func(s *session.Session) {
		r.removeFromQueue(s)
	})

	// 处理玩家加入队列
	qs := func() {
		r.lock.Lock()
		defer r.lock.Unlock()
	ctrl2:
		pvp := r.conf.Pvp
		if len(r.queue) >= int(pvp) {
			vv := r.queue[:pvp]
			r.queue = r.queue[pvp:]
			t := newTable(r.conf, vv)
			r.setTable(t)
			goto ctrl2
		}
	}

	// 处理桌子的清理
	clear := func() {
		for k, v := range r.tables {
			if v.canDestory() {
				delete(r.tables, k)
				log.Info("clear table %s", k)
			}
		}
	}

	scheduler.NewTimer(time.Second, func() {
	ctrl:
		for {
			select {
			default:
				qs()
				clear()
				break ctrl
			}
		}
	})

}

func (r *RoomService) table(tableId string) (*Table, error) {
	t, ok := r.tables[tableId]
	if !ok {
		return nil, z.NilError{Msg: tableId}
	}
	return t, nil
}

func (r *RoomService) setTable(t *Table) {
	r.tables[t.tableId] = t
}

func (r *RoomService) removeFromQueue(s *session.Session) {
	r.lock.Lock()
	defer r.lock.Unlock()
	for i, v := range r.queue {
		if v.ID() == s.ID() {
			r.queue = append(r.queue[:i], r.queue[:i+1]...)
			log.Info("removeFromQueue1 %d", s.UID())
			break
		}
	}
}

func (r *RoomService) addToQueue(s *session.Session) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.queue = append(r.queue, s)
	log.Info("addToQueue %d", s.UID())
}

func (r *RoomService) QuickStart(s *session.Session, _ *proto.Empty) error {
	r.addToQueue(s)
	return s.Response(&proto.GameStateResp{State: consts.WAIT})
}

func (r *RoomService) Cancel(s *session.Session, _ *proto.Empty) error {
	r.removeFromQueue(s)
	return s.Response(&proto.GameStateResp{State: consts.IDLE})
}

func (r *RoomService) UpdateState(s *session.Session, msg *proto.UpdateState) error {
	tid := GetSessionTableId(s)
	t, err := r.table(tid)
	if err != nil {
		return err
	}
	return t.updateState(s, msg)
}

func (r *RoomService) ResumeTable(s *session.Session, _ *proto.Empty) error {
	tid := GetSessionTableId(s)
	t, err := r.table(tid)
	if err != nil {
		return s.Response(&proto.TableInfo{TableId: ""})
	}
	return s.Response(t.getInfo())
}
