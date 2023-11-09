package room

import (
	"github.com/lonng/nano/session"
	"tetris/internal/game/util"
	"tetris/models"
	"tetris/pkg/log"
	"tetris/proto/proto"
)

type NormalClient struct {
	player  *proto.TableInfo_Player
	cs      *session.Session
	profile *models.Profile
}

// 方块
func NewNormalClient(s *session.Session, teamId int32) *NormalClient {
	p, _ := util.GetProfile(s)
	return &NormalClient{
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
		cs: s,
		//profile: p,
	}
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

func (c *NormalClient) Save(msg *proto.UpdateState) {
	fragment := msg.Fragment
	log.Info("client save %s", fragment)

	player := msg.Player
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
	arena := msg.Arena
	if arena != nil {
		c.player.State.Arena = arena
	}
	c.player.End = msg.End
}

func (c *NormalClient) GetSession() *session.Session {
	return c.cs
}
