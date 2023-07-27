package util

import (
	"github.com/lonng/nano/session"
	"tetris/config"
	"tetris/models"
	"tetris/proto/proto"
)

func ConvPlayerToProtoPlayer(pp *models.Profile) *proto.Profile {
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

func SetSessionTableId(s *session.Session, tableId string) {
	s.Set("tid", tableId)
}

func GetSessionTableId(s *session.Session) string {
	return s.String("tid")
}
