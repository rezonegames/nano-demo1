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

func ConvRoomListToProtoRoomList(rl []*config.Room) []*proto.Room {
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

// SetTableId 设置桌子
func SetTableId(s *session.Session, tableId string) {
	s.Set("tableId", tableId)
}

func GetTableId(s *session.Session) (string, error) {
	tid := s.String("tableId")
	if tid == "" {
		return "", z.NilError{}
	}
	return tid, nil
}

func RemoveTableId(s *session.Session) {
	s.Remove("tableId")
}

// SetWaiterId 设置waiterId
func SetWaiterId(s *session.Session, waiterId string) {
	s.Set("waiterId", waiterId)
}

func GetWaiterId(s *session.Session) (string, error) {
	wid := s.String("waiterId")
	if wid == "" {
		return "", z.NilError{}
	}
	return wid, nil
}

func ReomveWaiter(s *session.Session) {
	s.Remove("waiterId")
}
