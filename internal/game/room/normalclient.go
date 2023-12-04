package room

import (
	"github.com/lonng/nano/session"
	"sync"
	"tetris/internal/game/util"
	"tetris/pkg/z"
	"tetris/proto/proto"
	"time"
)

type NormalClient struct {
	player           *proto.TableInfo_Player
	cs               *session.Session
	lastUpdateAt     time.Time
	table            util.TableEntity
	frames           map[int64][]*proto.Action
	lock             sync.RWMutex //锁，用于存帧
	lastFrameId      int64        // 当前客户端到了哪帧
	speed            int64
	lastAutoDropTime time.Time
}

func NewNormalClient(s *session.Session, teamId int32, table util.TableEntity) *NormalClient {
	p, _ := util.GetProfile(s)
	c := &NormalClient{
		player: &proto.TableInfo_Player{
			TeamId:  teamId,
			End:     false,
			Profile: util.ConvProfileToProtoProfile(p),
		},
		cs: s,
		//updatedAt: z.GetTime(),
		table:  table,
		frames: make(map[int64][]*proto.Action, 0),
		speed:  3000,
	}

	return c
}

func (c *NormalClient) SetLastFrame(frameId int64) {
	c.lastFrameId = frameId
}

func (c *NormalClient) GetLastFrame() int64 {
	return c.lastFrameId
}

func (c *NormalClient) SaveFrame(frameId int64, msg *proto.UpdateFrame) {
	c.lock.Lock()
	defer c.lock.Unlock()

	action := msg.Action
	// 判断游戏是否结束，打一个标记
	if action.Key == proto.ActionType_END {
		c.player.End = true
	}
	uf, ok := c.frames[frameId]
	if !ok {
		uf = make([]*proto.Action, 0)
	}
	uf = append(uf, action)
	c.frames[frameId] = uf
	c.lastUpdateAt = z.GetTime()
}

func (c *NormalClient) GetFrame(frameId int64) []*proto.Action {
	c.lock.RLock()
	defer c.lock.RUnlock()

	al := make([]*proto.Action, 0)
	if uf, ok := c.frames[frameId]; ok {
		al = uf
	}

	// 自动下降
	now := z.GetTime()
	if now.Sub(c.lastAutoDropTime).Milliseconds() >= c.speed {
		al = append(al, &proto.Action{
			Key:     proto.ActionType_DROP,
			ValList: []int32{1},
		})
		c.frames[frameId] = al
		c.lastAutoDropTime = now
	}

	return al
}

func (c *NormalClient) GetSession() *session.Session {
	return c.cs
}

func (c *NormalClient) SetSession(s *session.Session) {
	c.cs = s
}

func (c *NormalClient) GetId() int64 {
	return c.player.Profile.UserId
}

func (c *NormalClient) GetPlayer() *proto.TableInfo_Player {
	return c.player
}

func (c *NormalClient) GetTeamId() int32 {
	return c.player.TeamId
}

func (c *NormalClient) IsEnd() bool {
	return c.player.End
}
