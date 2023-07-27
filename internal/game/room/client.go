package room

import (
	"github.com/lonng/nano/session"
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
			State:  nil,
			End:    false,
			Score:  0,
		},
		Session: s,
	}
}
