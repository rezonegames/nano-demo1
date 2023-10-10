package room

import (
	"fmt"
	"github.com/lonng/nano"
	"github.com/lonng/nano/session"
	"sync"
	"tetris/config"
	"tetris/internal/game/util"
	"tetris/pkg/log"
	"tetris/pkg/z"
	"tetris/proto/proto"
	"time"
)

type QuickRoom struct {
	group  *nano.Group
	config *config.Room
	tables map[string]util.TableEntity
	lock   sync.RWMutex
	queue  []*session.Session
}

func (r *QuickRoom) Entity(tableId string) (util.TableEntity, error) {
	t, ok := r.tables[tableId]
	if !ok {
		return nil, z.NilError{Msg: tableId}
	}
	return t, nil
}

func NewQuickRoom(conf *config.Room) *QuickRoom {
	r := &QuickRoom{
		group:  nano.NewGroup(conf.RoomId),
		config: conf,
		tables: make(map[string]util.TableEntity, 0),
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
	r.CreateTable(sList)
}

func (r *QuickRoom) BackToQueue(sList []*session.Session) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.queue = append(r.queue, sList...)
	for _, v := range sList {
		log.Debug("BackToQueue %d", v.UID())
	}
}

func (r *QuickRoom) Ready(s *session.Session) error {
	entity, err := util.GetTable(s, r)
	if err != nil {
		return err
	}
	return entity.Ready(s)
}

func (r *QuickRoom) Leave(s *session.Session) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	// 从group离开
	r.group.Leave(s)

	// 如果在队列从队列离开
	for i, v := range r.queue {
		if v.ID() == s.ID() {
			r.queue = append(r.queue[:i], r.queue[i+1:]...)
			log.Info("%s player leave queue %d", r.config.RoomId, s.UID())
			break
		}
	}

	//	如果在桌子，就从桌子离开
	if entity, err := util.GetTable(s, r); err == nil {
		entity.Leave(s)
	}
	util.RemoveRoom(s)
	return s.Response(&proto.LeaveResp{})
}

func (r *QuickRoom) Join(s *session.Session) error {
	if rid, err := util.GetRoomId(s); err == nil && rid != r.config.RoomId {
		return s.Response(&proto.GameStateResp{
			Code:   proto.ErrorCode_AlreadyInRoom,
			ErrMsg: fmt.Sprintf(`在另一个房间 %s`, rid),
		})
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	r.queue = append(r.queue, s)
	r.group.Add(s)
	util.SetRoomId(s, r.config.RoomId)
	log.Info("%s addToQueue %d", r.config.RoomId, s.UID())
	return s.Response(&proto.GameStateResp{
		State: proto.GameState_WAIT,
	})
}

func (r *QuickRoom) UpdateState(s *session.Session, msg *proto.UpdateState) error {
	entity, err := util.GetTable(s, r)
	if err != nil {
		return err
	}
	return entity.UpdateState(s, msg)
}

func (r *QuickRoom) ResumeTable(s *session.Session) error {
	entity, err := util.GetTable(s, r)
	if err != nil {
		return err
	}
	return s.Response(entity.GetInfo())
}

// CreateTable waiter或者其它的方式创建的房间
func (r *QuickRoom) CreateTable(sList []*session.Session) {
	t := NewTable(r, sList)
	r.tables[t.GetTableId()] = t
}
