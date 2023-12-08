package room

import (
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
	group   *nano.Group
	config  *config.Room
	tables  map[string]util.TableEntity
	lock    sync.RWMutex
	waitMap map[int64]int64
}

func NewQuickRoom(conf *config.Room) *QuickRoom {
	r := &QuickRoom{
		group:   nano.NewGroup(conf.RoomId),
		config:  conf,
		tables:  make(map[string]util.TableEntity, 0),
		waitMap: make(map[int64]int64, 0),
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
		pvp := int(r.config.Pvp)
		if len(r.waitMap) >= pvp {
			sList := make([]*session.Session, 0)
			for k, _ := range r.waitMap {
				s, err := r.group.Member(k)
				if err == nil {
					sList = append(sList, s)
				}
				delete(r.waitMap, k)
				if len(sList) >= pvp {
					break
				}
			}
			r.WaitReady(sList)
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
	for _, v := range sList {
		r.waitMap[v.UID()] = z.NowUnix()
		if err := models.SetRoomId(v.UID(), r.config.RoomId); err != nil {
			v.Push("onState", &proto.GameStateResp{
				State: proto.GameState_IDLE,
			})
			continue
		}
		log.Debug("BackToQueue %d", v.UID())
		v.Push("onState", &proto.GameStateResp{
			State: proto.GameState_WAIT,
		})
	}
}

func (r *QuickRoom) Leave(s *session.Session) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.group.Leave(s)
	delete(r.waitMap, s.UID())
	if rs, err := models.GetRoundSession(s.UID()); err == nil {
		if table, err := r.Entity(rs.TableId); err == nil {
			return table.Leave(s)
		}
	}
	log.Info("%d leave %s", s.UID(), r.config.RoomId)
	models.RemoveRoundSession(s.UID())
	return nil
}

func (r *QuickRoom) Join(s *session.Session) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if err := models.SetRoomId(s.UID(), r.config.RoomId); err != nil {
		return err
	}
	log.Info("%d addToQueue %s", s.UID(), r.config.RoomId)
	r.waitMap[s.UID()] = z.NowUnixMilli()
	return r.group.Add(s)
}

// CreateTable waiter或者其它的方式创建的房间
func (r *QuickRoom) CreateTable(sList []*session.Session) {
	t := NewTable(r, sList)
	r.tables[t.GetTableId()] = t
}
