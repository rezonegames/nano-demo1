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
			Profile: util.ConvProfileToProtoProfile(p),
			ResOK:   false,
		},
		cs:        s,
		id:        s.UID(),
		updatedAt: z.GetTime(),
		table:     table,
	}

	return c
}

func (c *NormalClient) Save(msg *proto.UpdateState) {
	state := c.player.State
	player, arena, end := msg.Player, msg.Arena, msg.End
	if player != nil {
		pos := player.Pos
		score := player.Score
		matrix := player.Matrix
		if pos != nil {
			state.Player.Pos = pos
		}
		if score != 0 {
			state.Player.Score = score
		}
		if matrix != nil {
			state.Player.Matrix = matrix
		}
	}
	if arena != nil {
		matrix := arena.Matrix
		if matrix != nil {
			state.Arena.Matrix = matrix
		}
	}
	if end {
		c.player.End = true
	}
	c.player.ResOK = msg.ResOK
	c.updatedAt = z.GetTime()
}

func (c *NormalClient) Update(t time.Time) {
	duration := t.Sub(c.updatedAt)
	duration = 1 * time.Millisecond
	if duration > time.Second {
		state := c.player.State
		s := &proto.UpdateState{
			Fragment: "all",
			Player:   state.Player,
			Arena:    state.Arena,
			PlayerId: c.id,
			End:      c.player.End,
		}
		if s.End {
			return
		}
		//
		// 如果刚上线就掉了， 初始化之
		if s.Arena.Matrix == nil {
			s.Arena.Matrix = fill0()
			reset(s)
		}

		//
		// 判断过了几秒，并且pos.y ++
		nSeconds := int(duration.Seconds())
		for i := 0; i < nSeconds; i += 1 {
			s.Player.Pos.Y++
			//
			// 碰撞了，pos.y--
			if collide(s) {
				s.Player.Pos.Y--
				merge(s)
				reset(s)
				//
				// 如果结束了，就放弃循环
				if s.End {
					break
				}
			}
		}

		//
		// 广播之
		c.table.BroadcastPlayerState(c.id, s)
	}

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
