package room

import (
	"github.com/lonng/nano/session"
	"tetris/internal/game/util"
	"tetris/models"
	"tetris/pkg/z"
	"tetris/proto/proto"
	"time"
)

type NormalClient struct {
	player    *proto.TableInfo_Player
	cs        *session.Session
	profile   *models.Profile
	updatedAt time.Time
	robot     util.RobotEntity
	id        int64
	table     util.TableEntity
}

func NewNormalClient(s *session.Session, teamId int32, table util.TableEntity) *NormalClient {
	p, _ := util.GetProfile(s)
	c := &NormalClient{
		player: &proto.TableInfo_Player{
			TeamId: teamId,
			State: &proto.State{
				Arena: &proto.Arena{Matrix: nil},
				Player: &proto.Player{
					Matrix: nil,
					Pos:    nil,
					Score:  0,
				},
			},
			End:     false,
			Score:   0,
			Profile: util.ConvProfileToProtoProfile(p),
		},
		cs:        s,
		id:        s.UID(),
		updatedAt: z.GetTime(),
		table:     table,
	}
	c.robot = NewRobot(c, table)
	return c
}

func (c *NormalClient) Save(msg *proto.UpdateState) {
	//fragment := msg.Fragment
	//log.Info("client save %s", fragment)

	player, arena := msg.Player, msg.Arena

	//
	// piece区域
	if player != nil {
		pos := player.Pos
		score := player.Score
		matrix := player.Matrix

		cp := c.player.State.Player
		if pos != nil {
			cp.Pos = pos
		}
		if score != 0 {
			cp.Score = score
		}
		if matrix != nil {
			cp.Matrix = matrix
		}
	}

	//
	// 面板
	if arena != nil {
		c.player.State.Arena = arena
	}

	//
	// 结束了
	if msg.End {
		c.player.End = true
	}

	//
	// 记录最后更新时间
	c.updatedAt = z.GetTime()
}

func (c *NormalClient) Update(t time.Time) {
	c.robot.Update(c.updatedAt, t)
}

func (c *NormalClient) GetSession() *session.Session {
	return c.cs
}

func (c *NormalClient) GetId() int64 {
	return c.id
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
