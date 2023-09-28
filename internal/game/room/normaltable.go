package room

import (
	"fmt"
	"github.com/lonng/nano"
	"github.com/lonng/nano/scheduler"
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
	clients       map[int64]util.ClientEntity
	loseCount     int32
	loseTeams     map[int32]int64
	teamGroupSize int32
	state         int32 // 从countdown =》settlement 结束
	lock          sync.RWMutex
	waiter        util.WaiterEntity
	room          util.RoomEntity
	begin         time.Time
	countDown     int32
}

func (t *Table) Entity() (util.WaiterEntity, error) {
	return t.waiter, nil
}

func (t *Table) Ready(s *session.Session) error {
	return t.waiter.Ready(s)
}

func NewNormalTable(room util.RoomEntity, sList []*session.Session) *Table {
	conf := room.GetConfig()
	tableId := fmt.Sprintf("%s:%d", conf.RoomId, z.NowUnixMilli())
	t := &Table{
		group:         nano.NewGroup(tableId),
		tableId:       tableId,
		clients:       make(map[int64]util.ClientEntity, 0),
		state:         consts.WAITREADY,
		conf:          conf,
		loseCount:     0,
		loseTeams:     make(map[int32]int64, 0),
		teamGroupSize: conf.Pvp / conf.Divide,
		room:          room,
		begin:         z.GetTime(),
		countDown:     3,
	}

	profiles := make(map[int64]*proto.Profile, 0)
	var teamId int32
	for i, v := range sList {
		uid := v.UID()
		if i > 0 && int32(i)%t.conf.Divide == 0 {
			teamId++
		}
		c := NewClient(v, teamId)
		t.clients[uid] = c
		t.Join(v, t.GetTableId())

		p := util.ConvProfileToProtoProfile(c.GetProfile())
		p.TeamId = teamId
		profiles[p.UserId] = p

		log.Info("p2t %d %d", uid, teamId)
	}
	for _, s := range sList {
		s.Push("onState", &proto.GameStateResp{
			State:    consts.WAITREADY,
			SubState: consts.WAITREADY_PROFILE,
			Profiles: profiles,
		})
	}

	log.Info("newTable start %s", tableId)
	t.waiter = NewWaiter(sList, room, t)
	return t
}

func (t *Table) BackToTable() {
	t.BroadcastTableState(consts.COUNTDOWN, consts.COUNTDOWN_BEGIN)
	time.Sleep(500 * time.Millisecond)

	for t.countDown > 0 {
		t.countDown--
		t.BroadcastTableState(consts.COUNTDOWN, consts.NONE)
		time.Sleep(time.Second)
	}

	t.BroadcastTableState(consts.GAMING, consts.GAME_BEGIN)
}

func (t *Table) AfterInit() {

}

func (t *Table) BroadcastTableState(state int32, subState int32) {
	t.state = state
	var roomList []*proto.Room
	var tableInfo *proto.TableInfo
	switch t.state {
	case consts.GAMING:
		tableInfo = t.GetInfo()
		roomList = util.ConvRoomListToProtoRoomList(config.ServerConfig.Rooms)
		break
	case consts.SETTLEMENT:
		roomList = util.ConvRoomListToProtoRoomList(config.ServerConfig.Rooms)
		break

	}
	t.group.Broadcast("onState", &proto.GameStateResp{
		State:     t.state,
		SubState:  subState,
		TableInfo: tableInfo,
		CountDown: t.countDown,
		RoomList:  roomList,
	})
}

func (t *Table) ClientEnd(teamId int32) {
	t.lock.Lock()
	defer t.lock.Unlock()

	isTeamLose := true
	for _, v := range t.clients {
		if v.GetTeamId() == teamId {
			if !v.IsEnd() {
				isTeamLose = false
			}
		}
	}
	// onTeamLose
	if isTeamLose {
		t.loseCount++
		t.loseTeams[teamId] = z.NowUnix()
		t.BroadcastTableState(consts.GAMING, consts.GAME_LOSE)
		log.Info("队伍 %d 结束", teamId)
	}

	if t.loseCount >= t.teamGroupSize-1 {
		t.BroadcastTableState(consts.SETTLEMENT, consts.SETTLEMENT_BEGIN)
		log.Info("本轮结束 %s", t.tableId)
	}
	scheduler.NewAfterTimer(3*time.Second, func() {
		t.BroadcastTableState(consts.SETTLEMENT, consts.SETTLEMENT_END)
	})
}

func (t *Table) GetClient(s *session.Session) util.ClientEntity {
	return t.clients[s.UID()]
}

func (t *Table) UpdateState(s *session.Session, msg *proto.UpdateState) error {
	// Sync message to all members in this room
	uid := s.UID()
	c := t.GetClient(s)
	c.Save(msg)
	teamId := c.GetTeamId()
	log.Info("updateState %d %d %v", uid, teamId, z.ToString(msg))
	if c.IsEnd() {
		t.ClientEnd(teamId)
	}
	return t.group.Broadcast("onStateUpdate", msg)
}

func (t *Table) Clear() {
	t.group.Close()
	for _, v := range t.clients {
		util.RemoveTable(v.GetSession())
	}
}

func (t *Table) Leave(s *session.Session) {
	t.waiter.Leave(s)
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
		players[k] = v.GetPlayer()
	}
	return &proto.TableInfo{
		Players:    players,
		TableId:    t.tableId,
		TableState: t.state,
		LoseTeams:  t.loseTeams,
	}
}

func (t *Table) GetTableId() string {
	return t.tableId
}
