package util

import (
	"github.com/lonng/nano/session"
	"tetris/config"
	"tetris/models"
	"tetris/pkg/z"
	"tetris/proto/proto"
)

func ConvProfileToProtoProfile(pp *models.Profile) *proto.Profile {
	return &proto.Profile{
		Name:   pp.Name,
		Coin:   pp.Coin,
		UserId: pp.UserId,
	}
}

func ConvRoomListToProtoRoomList(rl ...*config.Room) []*proto.Room {
	mm := make([]*proto.Room, 0)
	for _, v := range rl {
		mm = append(mm, &proto.Room{
			RoomId:  v.RoomId,
			Pvp:     v.Pvp,
			Name:    v.Name,
			MinCoin: v.MinCoin,
		})

	}
	return mm
}

// BindUser todo：如果是集群模式，使用remote方法同步profile数据，并保存到session
func BindUser(s *session.Session, p *models.Profile) error {
	s.Set("profile", p)
	s.Bind(p.UserId)
	return s.Bind(p.UserId)
}

func GetProfile(s *session.Session) (*models.Profile, error) {
	v := s.Value("profile")
	if vv, ok := v.(*models.Profile); ok {
		return vv, nil
	}
	return nil, z.NilError{}
}
