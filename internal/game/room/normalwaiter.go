package room

import (
	"github.com/google/uuid"
	"github.com/lonng/nano"
	"github.com/lonng/nano/scheduler"
	"github.com/lonng/nano/session"
	"tetris/consts"
	"tetris/proto/proto"
	"time"
)

type NormalWaiter struct {
	uiid        string
	group       *nano.Group
	readyList   []int64
	room        RoomEntity
	table       TableEntity
	timeCounter int32
	stime       *scheduler.Timer
	sList       []*session.Session
}

// NewWaiter返回错误代表有人下线了
func NewNormalWaiter(sList []*session.Session, room RoomEntity) *NormalWaiter {
	uiid := uuid.New().String()
	w := &NormalWaiter{
		uiid:  uiid,
		group: nano.NewGroup(uiid),
		room:  room,
		sList: make([]*session.Session, 0),
	}
	for _, v := range sList {
		w.group.Add(v)
	}
	w.sList = sList
	return w
}

func (w *NormalWaiter) GetId() string {
	return w.uiid
}

func (w *NormalWaiter) AfterInit() {
	session.Lifetime.OnClosed(func(s *session.Session) {
		w.CheckAndDismiss()
	})

	w.stime = scheduler.NewCountTimer(time.Second, 10, func() {
		w.timeCounter++
		// 都准备好了或者又离开的，解散waiter
		w.CheckAndDismiss()

	})
}

// Dismiss 解散
func (w *NormalWaiter) CheckAndDismiss() {
	// 中途有离开或者10秒倒计时有玩家没有准备，返回到队列
	if len(w.readyList) == w.group.Count() {
		w.table.BackToTable(w.sList)
	} else if w.timeCounter >= 10 && len(w.readyList) < w.group.Count() {
		bList := make([]*session.Session, 0)
		for _, v := range w.sList {
			bc := false
			for _, v0 := range w.readyList {
				if v0 == v.UID() {
					bc = true
				}
			}
			if bc {
				bList = append(bList, v)
			}
		}
		w.room.BackToQueue(bList)
		w.group.Broadcast("onState", &proto.GameStateResp{
			State: consts.WAIT,
		})
	} else {
		return
	}
	w.group.Close()
	w.sList = make([]*session.Session, 0)
	w.stime.Stop()
}

func (w *NormalWaiter) Ready(s *session.Session) error {
	w.readyList = append(w.readyList, s.UID())
	return w.group.Broadcast("onState", &proto.GameStateResp{
		State:     consts.WAITREADY,
		ReadyList: w.readyList,
	})
}
