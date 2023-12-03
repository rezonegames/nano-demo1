package room

import (
	"fmt"
	"github.com/lonng/nano"
	"github.com/lonng/nano/session"
	"sync"
	"tetris/config"
	"tetris/internal/game/util"
	"tetris/models"
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

	//
	// 处理桌子的清理
	clear := func() {
		for k, v := range r.tables {
			if v.CanDestory() {
				v.Clear()
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

func (r *QuickRoom) Entity(tableId string) (util.TableEntity, error) {
	t, ok := r.tables[tableId]
	if !ok {
		return nil, z.NilError{Msg: tableId}
	}
	return t, nil
}

func (r *QuickRoom) GetConfig() *config.Room {
	return r.config
}

// WaitReady 准备就绪后，进入确认界面，所有玩家都点击确认以后，才开始游戏
func (r *QuickRoom) WaitReady(sList []*session.Session) {
	r.CreateTable(sList)
}

func (r *QuickRoom) BackToWait(sList []*session.Session) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.queue = append(r.queue, sList...)
	for _, v := range sList {
		log.Debug("BackToQueue %d", v.UID())
		v.Push("onState", &proto.GameStateResp{
			State: proto.GameState_WAIT,
		})
	}
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
	if rs, err := models.GetRoundSession(s.UID()); err == nil {
		if table, err := r.Entity(rs.TableId); err == nil {
			table.Leave(s)
		}
	} else {
		models.RemoveRoundSession(s.UID())
	}
	return s.Response(&proto.LeaveResp{Code: proto.ErrorCode_OK})
}

func (r *QuickRoom) Join(s *session.Session) error {
	if rs, err := models.GetRoundSession(s.UID()); err == nil && rs.RoomId != r.config.RoomId {
		return s.Response(&proto.GameStateResp{
			Code:   proto.ErrorCode_AlreadyInRoom,
			ErrMsg: fmt.Sprintf(`在另一个房间 %s`, rs.RoomId),
		})
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	r.queue = append(r.queue, s)
	r.group.Add(s)
	log.Info("%s addToQueue %d", r.config.RoomId, s.UID())
	return s.Response(&proto.GameStateResp{
		Code:  proto.ErrorCode_OK,
		State: proto.GameState_WAIT,
	})
}

// CreateTable waiter或者其它的方式创建的房间
func (r *QuickRoom) CreateTable(sList []*session.Session) {
	t := NewTable(r, sList)
	r.tables[t.GetTableId()] = t
}
