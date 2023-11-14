package service

import (
	"github.com/lonng/nano/component"
	"github.com/lonng/nano/session"
	"tetris/internal/game/util"
	"tetris/pkg/z"
	"tetris/proto/proto"
)

type RoomService struct {
	*component.Base
	serviceName string
	rooms       map[string]util.RoomEntity
}

func NewRoomService() *RoomService {
	return &RoomService{
		rooms: make(map[string]util.RoomEntity, 0),
	}
}

func (r *RoomService) AfterInit() {
	// 处理玩家断开连接
	session.Lifetime.OnClosed(func(s *session.Session) {
		for _, v := range r.rooms {
			v.Leave(s)
		}
	})
}

func (r *RoomService) AddRoomEntity(roomId string, entity util.RoomEntity) {
	r.rooms[roomId] = entity
}

func (r *RoomService) entity(roomId string) (util.RoomEntity, error) {
	entity, ok := r.rooms[roomId]
	if !ok {
		return nil, z.NilError{Msg: roomId}
	}
	return entity, nil
}

func (r *RoomService) Join(s *session.Session, msg *proto.Join) error {
	entity, err := r.entity(msg.RoomId)
	if err != nil {
		return err
	}
	return entity.Join(s)
}

func (r *RoomService) Ready(s *session.Session, _ *proto.Ready) error {
	rid, err := util.GetRoomId(s)
	if err != nil {
		return err
	}
	entity, err := r.entity(rid)
	if err != nil {
		return err
	}
	return entity.Ready(s)
}

func (r *RoomService) Cancel(s *session.Session, _ *proto.Cancel) error {
	rid, err := util.GetRoomId(s)
	if err != nil {
		return s.Response(&proto.CancelResp{Code: proto.ErrorCode_OK})
	}
	entity, err := r.entity(rid)
	if err != nil {
		return err
	}
	return entity.Leave(s)
}

func (r *RoomService) UpdateState(s *session.Session, msg *proto.UpdateState) error {
	rid, err := util.GetRoomId(s)
	if err != nil {
		return err
	}
	entity, err := r.entity(rid)
	if err != nil {
		return err
	}
	return entity.UpdateState(s, msg)
}

func (r *RoomService) ResumeTable(s *session.Session, msg *proto.Join) error {
	entity, err := r.entity(msg.RoomId)
	if err != nil {
		return err
	}
	return entity.ResumeTable(s)
}
