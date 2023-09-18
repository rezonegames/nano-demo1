package room

import (
	"fmt"
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

// 桌子
type Table struct {
	group         *nano.Group
	tableId       string
	conf          *config.Room
	clients       map[int64]*Client
	loseCount     int32
	teamGroupSize int32
	state         int32 // 从countdown =》settlement 结束
	lock          sync.RWMutex
	waiter        WaiterEntity
	room          RoomEntity
	begin         time.Time
}

func (t *Table) Ready(s *session.Session) error {
	return t.waiter.Ready(s)
}

func NewNormalTable(room RoomEntity, sList []*session.Session) *Table {
	conf := room.GetConfig()
	tableId := fmt.Sprintf("%s:%d", conf.RoomId, z.NowUnixMilli())
	t := &Table{
		group:         nano.NewGroup(tableId),
		tableId:       tableId,
		clients:       make(map[int64]*Client, 0),
		state:         consts.WAITREADY,
		conf:          conf,
		loseCount:     0,
		teamGroupSize: conf.Pvp / conf.Divide,
		room:          room,
		begin:         z.GetTime(),
	}

	var teamId int32
	for i, v := range sList {
		uid := v.UID()
		if i > 0 && int32(i)%t.conf.Divide == 0 {
			teamId++
		}
		t.clients[uid] = newClient(v, teamId)
		t.Join(v, t.GetTableId())
		log.Info("p2t %d %d", uid, teamId)
	}

	log.Info("newTable start %s", tableId)
	t.waiter = NewWaiter(sList, room, t)
	return t
}

func (t *Table) BackToTable(sList []*session.Session) {
	t.BroadcastTableState(consts.COUNTDOWN)
	time.Sleep(50 * time.Millisecond)
	for i := 0; i < 3; i++ {
		t.group.Broadcast("onCountDown", &proto.CountDownResp{
			Counter: int32(3 - i),
		})
		time.Sleep(time.Second)
	}
	t.BroadcastTableState(consts.GAMING)
}

func (t *Table) AfterInit() {

}

func (t *Table) BroadcastTableState(s int32) {
	t.state = s
	t.group.Broadcast("onState", &proto.GameStateResp{
		State:     t.state,
		TableInfo: t.GetInfo(),
	})
}

func (t *Table) ClientEnd(teamId int32) {
	t.lock.Lock()
	defer t.lock.Unlock()

	isTeamLose := true
	for _, v := range t.clients {
		if v.getTeamId() == teamId {
			if !v.isEnd() {
				isTeamLose = false
			}
		}
	}
	// onTeamLose
	if isTeamLose {
		t.loseCount++
		t.group.Broadcast("onTeamLose", &proto.TeamLose{TeamId: teamId})
		log.Info("队伍 %d 结束", teamId)
	}

	if t.loseCount >= t.teamGroupSize-1 {
		t.BroadcastTableState(consts.SETTLEMENT)
		log.Info("本轮结束 %s", t.tableId)
	}
}

func (t *Table) GetClient(uid int64) *Client {
	return t.clients[uid]
}

func (t *Table) UpdateState(s *session.Session, msg *proto.UpdateState) error {
	// Sync message to all members in this room
	uid := s.UID()
	c := t.GetClient(uid)
	c.save(msg)
	teamId := c.getTeamId()
	log.Info("updateState %d %d %v", uid, teamId, z.ToString(msg))
	if c.isEnd() {
		t.ClientEnd(teamId)
	}
	return t.group.Broadcast("onStateUpdate", msg)
}

func (t *Table) Clear() {
	t.group.Close()
	for _, v := range t.clients {
		util.SetTableId(v.Session, "")
	}
}

func (t *Table) Leave(s *session.Session) {
	t.group.Leave(s)
}

func (t *Table) Join(s *session.Session, tableId string) error {
	util.SetTableId(s, tableId)
	return t.group.Add(s)
}

func (t *Table) CanDestory() bool {
	switch t.state {
	case consts.SETTLEMENT:
	case consts.CANCEL:
		return true
	case consts.WAITREADY:
		// 1分钟了，还是这个状态，销毁之
		if z.GetTime().Sub(t.begin) > time.Minute {
			return true
		}
		return false
	default:
		return t.group.Count() == 0
	}
	return false
}

func (t *Table) GetInfo() *proto.TableInfo {
	players := make(map[int64]*proto.TableInfo_Player, 0)
	for k, v := range t.clients {
		players[k] = v.getPlayer()
	}
	return &proto.TableInfo{
		Players:    players,
		TableId:    t.tableId,
		TableState: t.state,
	}
}

func (t *Table) GetTableId() string {
	return t.tableId
}
