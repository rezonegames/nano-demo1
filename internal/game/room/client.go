package room

import (
	"github.com/lonng/nano/session"
	"tetris/pkg/log"
	"tetris/proto/proto"
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
			State: &proto.State{
				Arena: nil,
				Player: &proto.Player{
					Matrix: nil,
					Pos:    nil,
					Score:  0,
				},
			},
			End:   false,
			Score: 0,
		},
		Session: s,
	}
}

func (c *Client) getPlayer() *proto.TableInfo_Player {
	return c.player
}

func (c *Client) getTeamId() int32 {
	return c.player.TeamId
}

func (c *Client) isEnd() bool {
	return c.player.End
}

func (c *Client) save(msg *proto.UpdateState) {
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
