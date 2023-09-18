package room

import (
	"github.com/lonng/nano"
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
	group  *nano.Group
	config *config.Room
	tables map[string]TableEntity
	lock   sync.RWMutex
	queue  []*session.Session
}

func NewQuickRoom(conf *config.Room) *QuickRoom {
	r := &QuickRoom{
		group:  nano.NewGroup(conf.RoomId),
		config: conf,
		tables: make(map[string]TableEntity, 0),
	}
	log.Info(conf.Dump())
	return r
}

func (r *QuickRoom) AfterInit() {
	// 处理玩家断线等情况
	session.Lifetime.OnClosed(func(s *session.Session) {
		r.Leave(s)
	})

	// 处理玩家加入队列
	qs := func() {
		r.lock.Lock()
		defer r.lock.Unlock()
	ctrl2:
		pvp := r.config.Pvp
		if len(r.queue) >= int(pvp) {
			vv := r.queue[:pvp]
			r.queue = r.queue[pvp:]
			r.WaitReady(vv)
			goto ctrl2
		}
	}

	// 处理桌子的清理
	clear := func() {
		for k, v := range r.tables {
			if v.CanDestory() {
				delete(r.tables, k)
				log.Info("clear table %s", k)
			}
		}
	}

	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				qs()
				clear()
			default:

			}
		}
	}()
}

func (r *QuickRoom) GetConfig() *config.Room {
	return r.config
}

// WaitReady 准备就绪后，进入确认界面，所有玩家都点击确认以后，才开始游戏
func (r *QuickRoom) WaitReady(sList []*session.Session) {
	wc := len(sList)
	profiles := make(map[int64]*proto.Profile, 0)
	for i, v := range sList {
		p, err := util.GetProfile(v)
		if err != nil {
			sList = append(sList[:i], sList[i+1:]...)
			r.Leave(v)
			break
		}
		profiles[p.UserId] = util.ConvProfileToProtoProfile(p)
	}
	if len(profiles) < wc {
		r.BackToQueue(sList)
		return
	}
	for _, s := range sList {
		s.Push("onState", &proto.GameStateResp{
			State:    consts.WAITREADY,
			Profiles: profiles,
		})
	}
	r.CreateTable(sList)
}

func (r *QuickRoom) BackToQueue(sList []*session.Session) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.queue = append(r.queue, sList...)
}

func (r *QuickRoom) Ready(s *session.Session) error {
	tid, err := util.GetTableId(s)
	if err != nil {
		return err
	}
	t, err := r.Table(tid)
	if err != nil {
		return err
	}
	return t.Ready(s)
}

func (r *QuickRoom) Table(tableId string) (TableEntity, error) {
	t, ok := r.tables[tableId]
	if !ok {
		return nil, z.NilError{Msg: tableId}
	}
	return t, nil
}

func (r *QuickRoom) Leave(s *session.Session) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	// 从group离开
	r.group.Leave(s)

	// 如果在队列从队列离开
	for i, v := range r.queue {
		if v.ID() == s.ID() {
			r.queue = append(r.queue[:i], r.queue[:i+1]...)
			log.Info("removeFromQueue1 %d", s.UID())
			break
		}
	}

	//	如果在桌子，就从桌子离开
	if tid, err := util.GetTableId(s); err != nil {
		if t, err := r.Table(tid); err == nil {
			t.Leave(s)
		}
	}
	return nil
}

func (r *QuickRoom) Join(s *session.Session) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.queue = append(r.queue, s)
	r.group.Add(s)
	log.Info("addToQueue %d", s.UID())
	return s.Response(&proto.GameStateResp{
		State: consts.WAIT,
	})
}

func (r *QuickRoom) UpdateState(s *session.Session, msg *proto.UpdateState) error {
	t, err := r.Table(msg.TableId)
	if err != nil {
		return err
	}
	return t.UpdateState(s, msg)
}

func (r *QuickRoom) ResumeTable(s *session.Session) error {
	tid, err := util.GetTableId(s)
	if err != nil {
		return err
	}
	t, err := r.Table(tid)
	if err != nil {
		return err
	}
	return s.Response(t.GetInfo())
}

// CreateTable waiter或者其它的方式创建的房间
func (r *QuickRoom) CreateTable(sList []*session.Session) {
	t := NewTable(r, sList)
	r.tables[t.GetTableId()] = t
}
