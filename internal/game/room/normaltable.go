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

// 桌子
type Table struct {
	group         *nano.Group
	tableId       string
	conf          *config.Room
	clients       map[int64]util.ClientEntity
	loseCount     int32
	loseTeams     map[int32]int64
	teamGroupSize int32
	lock          sync.RWMutex
	waiter        util.WaiterEntity
	room          util.RoomEntity
	begin         time.Time
	state         proto.TableState
	end           chan bool
	nextFrameId   int64
	randSeed      int64
	pieceList     []int32
	resCountDown  int32 // 检查资源是否加载成功
	res           map[int64]int32
}

func NewNormalTable(room util.RoomEntity, sList []*session.Session) *Table {
	conf := room.GetConfig()
	now := z.NowUnixMilli()
	tableId := fmt.Sprintf("%s:%d", conf.RoomId, now)
	t := &Table{
		group:         nano.NewGroup(tableId),
		tableId:       tableId,
		clients:       make(map[int64]util.ClientEntity, 0),
		conf:          conf,
		loseCount:     0,
		loseTeams:     make(map[int32]int64, 0),
		teamGroupSize: conf.Pvp / conf.Divide,
		room:          room,
		begin:         z.GetTime(),
		end:           make(chan bool, 0),
		resCountDown:  100,
		randSeed:      now,
		pieceList:     make([]int32, 0),
		res:           make(map[int64]int32, 0),
	}

	for i := 0; i < 500; i++ {
		t.pieceList = append(t.pieceList, z.RandInt32(0, 6))
	}

	var teamId int32
	for i, v := range sList {
		uid := v.UID()
		if i > 0 && int32(i)%t.conf.Divide == 0 {
			teamId++
		}
		c := NewClient(v, teamId, t)
		t.clients[uid] = c
		t.Join(v, tableId)

		log.Info("p2t %d %d", uid, teamId)
	}

	t.waiter = NewWaiter(sList, room, t)

	log.Info("newTable start %s", tableId)
	return t
}

func (t *Table) BackToTable() {
	go func() {
		for t.resCountDown > 0 {
			t.ChangeState(proto.TableState_CHECK_RES)
			b := true
			for id, _ := range t.clients {
				if t.res[id] != 100 {
					log.Info("check res not ok %d", id)
					b = false
				}
			}
			if b || t.resCountDown == 1 {
				t.ChangeState(proto.TableState_GAMING)
				return
			}
			time.Sleep(time.Second)
			t.resCountDown--
		}
	}()
}

func (t *Table) AfterInit() {
	go t.Run()
}

func (t *Table) Run() {
	ticker := time.NewTicker(33 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {

		case <-t.end:
			t.group.Close()
			for _, v := range t.clients {
				models.RemoveRoundSession(v.GetId())
			}
			return

		case <-ticker.C:
			switch t.state {
			case proto.TableState_GAMING:
				msg := &proto.OnFrame{
					FrameId:    t.nextFrameId,
					PlayerList: make([]*proto.OnFrame_Player, 0),
				}

				bSend := false
				if t.nextFrameId == 0 {
					bSend = true
					msg.PieceList = t.pieceList
				} else {
					for _, v := range t.clients {
						if al := v.GetFrame(t.nextFrameId); len(al) > 0 {
							msg.PlayerList = append(msg.PlayerList, &proto.OnFrame_Player{
								UserId:     v.GetId(),
								ActionList: al,
							})
							bSend = true
						}
					}
				}

				t.nextFrameId++
				if bSend {
					t.group.Broadcast("onFrame", msg)
				}
			}
			break
		}
	}
}

func (t *Table) ChangeState(state proto.TableState) {
	t.state = state
	var roomList []*proto.Room
	tableInfo := t.GetInfo()
	switch state {
	case proto.TableState_CHECK_RES:
		tableInfo.Res = &proto.TableInfo_Res{
			Players:   t.res,
			CountDown: t.resCountDown,
		}
		break

	case proto.TableState_SETTLEMENT:
		roomList = util.ConvRoomListToProtoRoomList(config.ServerConfig.Rooms...)
		break
	}
	t.group.Broadcast("onState", &proto.GameStateResp{
		State:     proto.GameState_INGAME,
		TableInfo: tableInfo,
		RoomList:  roomList,
	})
}

func (t *Table) WaiterEntity() util.WaiterEntity {
	return t.waiter
}

func (t *Table) Entity(uid int64) util.ClientEntity {
	return t.clients[uid]
}

func (t *Table) Ready(s *session.Session) error {
	return t.waiter.Ready(s)
}

func (t *Table) LoadRes(s *session.Session, msg *proto.LoadRes) error {
	t.res[s.UID()] = msg.Current
	log.Info("LoadRes down %d", s.UID())
	return nil
}

func (t *Table) Update(s *session.Session, msg *proto.UpdateFrame) error {
	if t.state != proto.TableState_GAMING {
		return nil
	}

	uid := s.UID()
	// 有些操作可能是其它客户端上传的
	who := msg.Action.From
	if who != 0 {
		uid = who
	}

	c := t.Entity(uid)
	c.SaveFrame(t.nextFrameId, msg)

	// 如果已经结束了，判断输赢
	if c.IsEnd() {
		isTeamLose := true
		teamId := c.GetTeamId()
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
			t.ChangeState(proto.TableState_GAMING)
			log.Info("队伍 %d 结束", teamId)
		}

		if t.loseCount >= t.teamGroupSize-1 {
			t.ChangeState(proto.TableState_SETTLEMENT)
			log.Info("本轮结束 %s", t.tableId)
		}
	}
	return nil
}

func (t *Table) Clear() {
	t.end <- true
}

func (t *Table) Leave(s *session.Session) {
	t.waiter.Leave(s)
	t.group.Leave(s)
}

func (t *Table) Join(s *session.Session, tableId string) error {
	err := t.group.Add(s)
	if err != nil {
		return err
	}
	return models.SetRoundSession(s.UID(), &models.RoundSession{
		RoomId:  t.room.GetConfig().RoomId,
		TableId: t.tableId,
	})
}

func (t *Table) CanDestory() bool {
	switch t.state {
	case proto.TableState_SETTLEMENT:
		fallthrough
	case proto.TableState_CANCEL:
		return true
	case proto.TableState_WAITREADY:
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
	rooms := util.ConvRoomListToProtoRoomList(t.room.GetConfig())
	return &proto.TableInfo{
		Players:    players,
		TableId:    t.tableId,
		TableState: t.state,
		LoseTeams:  t.loseTeams,
		Waiter:     t.waiter.GetInfo(),
		Room:       rooms[0],
		RandSeed:   t.randSeed,
	}
}

func (t *Table) GetTableId() string {
	return t.tableId
}

func (t *Table) ResumeTable(s *session.Session, msg *proto.ResumeTable) error {
	t.group.Add(s)
	frameList := make([]*proto.OnFrame, 0)
	for i := msg.FrameId; i < t.nextFrameId-1; i++ {

		frame := &proto.OnFrame{
			FrameId: i,
		}
		// 第0帧
		if i == 0 {
			frame.PieceList = t.pieceList
			frameList = append(frameList, frame)
			continue
		}

		// 其它帧
		pl := make([]*proto.OnFrame_Player, 0)
		for k, v := range t.clients {
			if al := v.GetHistoryFrame(i); len(al) > 0 {
				pl = append(pl, &proto.OnFrame_Player{
					UserId:     k,
					ActionList: al,
				})
			}
		}
		if len(pl) > 0 {
			frame.PlayerList = pl
			frameList = append(frameList, frame)
		}

	}
	tableInfo := t.GetInfo()
	tableInfo.FrameList = frameList
	return s.Response(&proto.GameStateResp{
		Code:      proto.ErrorCode_OK,
		State:     proto.GameState_INGAME,
		TableInfo: tableInfo,
	})
}
