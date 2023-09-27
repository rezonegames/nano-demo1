package room

import (
	"github.com/google/uuid"
	"github.com/lonng/nano"
	"github.com/lonng/nano/scheduler"
	"github.com/lonng/nano/session"
	"tetris/consts"
	"tetris/internal/game/util"
	"tetris/pkg/log"
	"tetris/pkg/z"
	"tetris/proto/proto"
	"time"
)

type NormalWaiter struct {
	uiid        string
	group       *nano.Group
	readys      map[int64]int64
	room        util.RoomEntity
	table       util.TableEntity
	timeCounter int32
	countDown   int
	stime       *scheduler.Timer
	sList       []*session.Session
	profiles    map[int64]*proto.Profile
}

// NewWaiter返回错误代表有人下线了
func NewNormalWaiter(sList []*session.Session, room util.RoomEntity, table util.TableEntity) *NormalWaiter {
	uiid := uuid.New().String()
	w := &NormalWaiter{
		uiid:      uiid,
		group:     nano.NewGroup(uiid),
		room:      room,
		sList:     make([]*session.Session, 0),
		table:     table,
		readys:    make(map[int64]int64, 0),
		countDown: 30,
		profiles:  make(map[int64]*proto.Profile, 0),
	}
	for _, v := range sList {
		w.group.Add(v)
		w.profiles[v.UID()] = util.ConvProfileToProtoProfile(table.GetClient(v).GetProfile())
	}
	w.sList = sList
	return w
}

func (w *NormalWaiter) AfterInit() {
	w.stime = scheduler.NewCountTimer(time.Second, w.countDown, func() {
		w.timeCounter++
		// 都准备好了或者又离开的，解散waiter
		w.CheckAndDismiss()
	})
}

// Dismiss 解散
func (w *NormalWaiter) CheckAndDismiss() {
	// 中途有离开或者10秒倒计时有玩家没有准备，返回到队列
	//log.Debug("CheckAndDismiss %d %d", len(w.readys), w.group.Count())
	if len(w.readys) == len(w.sList) {
		log.Debug("waiter done %s", w.table.GetTableId())
		w.table.BackToTable()
	} else if (w.timeCounter >= int32(w.countDown) && len(w.readys) < len(w.sList)) ||
		w.group.Count() != len(w.sList) {
		log.Debug("Dismiss waiter!!")
		bList := make([]*session.Session, 0)
		for _, v := range w.sList {
			if _, ok := w.readys[v.UID()]; ok {
				bList = append(bList, v)
			} else {
				w.group.Leave(v)
			}
		}
		w.room.BackToQueue(bList)
		w.group.Broadcast("onState", &proto.GameStateResp{
			State: consts.WAIT,
		})
	} else {
		w.group.Broadcast("onState", &proto.GameStateResp{
			State:     consts.WAITREADY,
			SubState:  consts.WAITREADY_COUNTDOWN,
			CountDown: int32(w.countDown) - w.timeCounter,
			Profiles:  w.profiles,
		})
		return
	}
	w.group.Close()
	w.sList = make([]*session.Session, 0)
	w.stime.Stop()
}

// Ready 准备
func (w *NormalWaiter) Ready(s *session.Session) error {
	uid := s.UID()
	if _, ok := w.readys[uid]; ok {
		return s.Response(&proto.ReadyResp{})
	}
	w.readys[uid] = z.NowUnix()
	return w.group.Broadcast("onState", &proto.GameStateResp{
		State:    consts.WAITREADY,
		SubState: consts.WAITREADY_READYLIST,
		Readys:   w.readys,
	})
	return s.Response(&proto.ReadyResp{})
}

func (w *NormalWaiter) Leave(s *session.Session) error {
	return w.group.Leave(s)
}
