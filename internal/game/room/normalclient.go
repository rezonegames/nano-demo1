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
	player       *proto.TableInfo_Player
	cs           *session.Session
	lastUpdateAt time.Time
	table        util.TableEntity
	frames       map[int64]*proto.UpdateFrame
	lock         sync.RWMutex //锁，用于存帧
}

func NewNormalClient(s *session.Session, teamId int32, table util.TableEntity) *NormalClient {
	p, _ := util.GetProfile(s)
	c := &NormalClient{
		player: &proto.TableInfo_Player{
			TeamId:    teamId,
			FrameList: make([]*proto.TableInfo_Frame, 0),
			End:       false,
			Profile:   util.ConvProfileToProtoProfile(p),
			ResOK:     false,
		},
		cs: s,
		//updatedAt: z.GetTime(),
		table:  table,
		frames: make(map[int64]*proto.UpdateFrame, 0),
	}

	return c
}

func (c *NormalClient) SaveFrame(frameId int64, msg *proto.UpdateFrame) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// 判断游戏是否结束，打一个标记
	for _, a := range msg.ActionList {
		if a.Key == proto.ActionType_END {
			c.player.End = true
			break
		}
	}
	if uf, ok := c.frames[frameId]; ok {
		uf.ActionList = append(uf.ActionList, msg.ActionList...)
	} else {
		uf.ActionList = msg.ActionList
	}
	c.lastUpdateAt = z.GetTime()
}

func (c *NormalClient) GetFrame(frameId int64) []*proto.Action {
	c.lock.RLock()
	defer c.lock.RUnlock()

	al := make([]*proto.Action, 0)
	uf, ok := c.frames[frameId]
	if ok {
		al = uf.ActionList
	} else {

		// 没结束，并且玩家这一帧没有操作，服务器主动加一帧
		if !c.player.End {
			if z.GetTime().Sub(c.lastUpdateAt) >= time.Second {
				al = append(al, &proto.Action{
					Key: proto.ActionType_DROP,
					Val: 1,
				})
				c.frames[frameId] = &proto.UpdateFrame{ActionList: al}
				c.lastUpdateAt = z.GetTime()
			}
		}
	}
	return al
}

func (c *NormalClient) GetSession() *session.Session {
	return c.cs
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
