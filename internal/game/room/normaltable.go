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
}

func newNormalTable(conf *config.Room, ss []*session.Session) *Table {
	tableId := fmt.Sprintf("%s:%d", conf.RoomId, z.NowUnixMilli())
	t := &Table{
		group:         nano.NewGroup(tableId),
		tableId:       tableId,
		clients:       make(map[int64]*Client, 0),
		state:         consts.WAIT,
		conf:          conf,
		loseCount:     0,
		teamGroupSize: conf.Pvp / conf.Divide,
	}
	var teamId int32
	for i, v := range ss {
		uid := v.UID()
		if i > 0 && int32(i)%t.conf.Divide == 0 {
			teamId++
		}
		t.clients[uid] = newClient(v, teamId)
		t.join(v, tableId)
		log.Info("p2t %d %d", uid, teamId)
	}
	log.Info("newTable start %s", tableId)
	go t.run()
	return t
}

func (t *Table) run() {
	t.broadcastTableState(consts.COUNTDOWN)
	time.Sleep(50 * time.Millisecond)
	for i := 0; i < 3; i++ {
		t.group.Broadcast("onCountDown", &proto.CountDownResp{
			Counter: int32(3 - i),
		})
		time.Sleep(time.Second)
	}
	t.broadcastTableState(consts.GAMING)
}

func (t *Table) broadcastTableState(s int32) {
	t.state = s
	t.group.Broadcast("onStateChange", &proto.GameStateResp{
		State:     t.state,
		TableInfo: t.getInfo(),
	})
}

func (t *Table) clientEnd(teamId int32) {
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
		t.broadcastTableState(consts.SETTLEMENT)
		log.Info("本轮结束 %s", t.tableId)
	}
}

func (t *Table) getClient(uid int64) *Client {
	return t.clients[uid]
}

func (t *Table) updateState(s *session.Session, msg *proto.UpdateState) error {
	// Sync message to all members in this room
	uid := s.UID()
	c := t.getClient(uid)
	c.save(msg)
	teamId := c.getTeamId()
	log.Info("updateState %d %d %v", uid, teamId, z.ToString(msg))
	if c.isEnd() {
		t.clientEnd(teamId)
	}
	return t.group.Broadcast("onStateUpdate", msg)
}

func (t *Table) clear() {
	t.group.Close()
	for _, v := range t.clients {
		util.SetSessionTableId(v.Session, "")
	}
}

func (t *Table) leave(s *session.Session) {
	t.group.Leave(s)
}

func (t *Table) join(s *session.Session, tableId string) error {
	util.SetSessionTableId(s, tableId)
	return t.group.Add(s)
}

func (t *Table) canDestory() bool {
	return t.state == consts.SETTLEMENT || t.group.Count() == 0
}

func (t *Table) getInfo() *proto.TableInfo {
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

func (t *Table) getTableId() string {
	return t.tableId
}
