package room

import (
	"github.com/lonng/nano/session"
	"sync"
	"tetris/config"
	"tetris/consts"
	"tetris/internal/game/util"
	"tetris/pkg/log"
	"tetris/pkg/z"
	"tetris/proto/proto"
	"time"
)

type QuickRoom struct {
	conf   *config.Room
	tables map[string]TableEntity
	lock   sync.RWMutex
	queue  []*session.Session
}

func newQuickRoom(conf *config.Room) *QuickRoom {
	r := &QuickRoom{
		conf:   conf,
		tables: make(map[string]TableEntity, 0),
	}
	go r.run()
	return r
}

func (r *QuickRoom) run() {

	// 处理玩家加入队列
	qs := func() {
		r.lock.Lock()
		defer r.lock.Unlock()
	ctrl2:
		pvp := r.conf.Pvp
		if len(r.queue) >= int(pvp) {
			vv := r.queue[:pvp]
			r.queue = r.queue[pvp:]
			t, _ := NewTable(r.conf, vv)
			r.setTable(t.getTableId(), t)
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

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			qs()
			clear()
		default:

		}
	}
}

func (r *QuickRoom) table(tableId string) (TableEntity, error) {
	t, ok := r.tables[tableId]
	if !ok {
		return nil, z.NilError{Msg: tableId}
	}
	return t, nil
}

func (r *QuickRoom) Leave(s *session.Session) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	for i, v := range r.queue {
		if v.ID() == s.ID() {
			r.queue = append(r.queue[:i], r.queue[:i+1]...)
			log.Info("removeFromQueue1 %d", s.UID())
			break
		}
	}
	return nil
}

func (r *QuickRoom) Join(s *session.Session) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.queue = append(r.queue, s)
	log.Info("addToQueue %d", s.UID())
	return s.Response(&proto.GameStateResp{
		State: consts.WAIT,
	})
}

func (r *QuickRoom) UpdateState(s *session.Session, msg *proto.UpdateState) error {
	t, err := r.table(msg.TableId)
	if err != nil {
		return err
	}
	return t.updateState(s, msg)
}

func (r *QuickRoom) ResumeTable(s *session.Session) error {
	tid := util.GetSessionTableId(s)
	t, err := r.table(tid)
	if err != nil {
		return err
	}
	return s.Response(t.getInfo())
}

func (r *QuickRoom) setTable(tableId string, t TableEntity) {
	r.tables[tableId] = t
}
