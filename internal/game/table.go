package game

import (
	"fmt"
	"github.com/lonng/nano"
	"github.com/lonng/nano/session"
	"tetris/config"
	"tetris/consts"
	"tetris/pkg/log"
	"tetris/pkg/z"
	"tetris/proto/proto"
	"time"
)

type Client struct {
	player *proto.TableInfo_Player
	*session.Session
}

// 方块
func newClient(s *session.Session, teamId int32) *Client {
	return &Client{
		player: &proto.TableInfo_Player{
			TeamId: teamId,
			State:  nil,
			End:    false,
			Score:  0,
		},
		Session: s,
	}
}

// 桌子
type Table struct {
	group         *nano.Group
	tableId       string
	conf          *config.Room
	clients       map[int64]*Client
	loseCount     int32
	teamGroupSize int32
	state         int32 // 从countdown =》settlement 结束

}

func newTable(conf *config.Room, ss []*session.Session) *Table {
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
	for i := 0; i < 5; i++ {
		t.group.Broadcast("onCountDown", &proto.CountDownResp{
			Counter: int32(5 - i),
		})
		time.Sleep(time.Second)
	}

	t.broadcastTableState(consts.GAMEING)
}

func (t *Table) broadcastTableState(s int32) {
	t.state = s
	t.group.Broadcast("onStateChange", &proto.GameStateResp{
		State:     t.state,
		TableInfo: t.getInfo(),
	})
}

func (t *Table) updateState(s *session.Session, msg *proto.UpdateState) error {
	// Sync message to all members in this room
	uid := s.UID()
	p := t.clients[uid].player
	teamId := p.TeamId
	p.End = msg.End
	log.Info("updateState %d %d %v", uid, teamId, z.ToString(msg))

	if msg.End {
		isTeamLose := true
		for _, v := range t.clients {
			if v.player.TeamId == teamId {
				if !v.player.End {
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

		if t.loseCount == t.teamGroupSize-1 {
			t.broadcastTableState(consts.SETTLEMENT)
			log.Info("本轮结束 %s", t.tableId)
			return nil
		}
	}
	return nil
}

func (t *Table) clear() {
	t.group.Close()
	for _, v := range t.clients {
		v.Set("tid", "")
	}
}

func (t *Table) leave(s *session.Session) {
	t.group.Leave(s)
}

func (t *Table) join(s *session.Session, tableId string) error {
	SetSessionTableId(s, tableId)
	return t.group.Add(s)
}

func (t *Table) canDestory() bool {
	return t.state == consts.SETTLEMENT
}

func (t *Table) getInfo() *proto.TableInfo {
	players := make(map[int64]*proto.TableInfo_Player, 0)
	for k, v := range t.clients {
		players[k] = v.player
	}
	return &proto.TableInfo{
		Players:    players,
		TableId:    t.tableId,
		TableState: t.state,
	}
}
