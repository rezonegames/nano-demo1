package service

import (
	"github.com/lonng/nano/component"
	"github.com/lonng/nano/session"
	"tetris/internal/game/util"
	"tetris/models"
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

func (r *RoomService) Entity(roomId string) util.RoomEntity {
	return r.rooms[roomId]
}

func (r *RoomService) Join(s *session.Session, msg *proto.Join) error {
	return r.Entity(msg.RoomId).Join(s)
}

func (r *RoomService) Leave(s *session.Session, _ *proto.Leave) error {
	rs, err := models.GetRoundSession(s.UID())
	if err != nil {
		return s.Response(&proto.LeaveResp{Code: proto.ErrorCode_OK})
	}
	return r.Entity(rs.RoomId).Leave(s)
}

func (r *RoomService) getTable(s *session.Session) (util.TableEntity, error) {
	rs, err := models.GetRoundSession(s.UID())
	if err != nil {
		return nil, err
	}
	table, err := r.Entity(rs.RoomId).Entity(rs.TableId)
	if err != nil {
		return nil, err
	}
	return table, nil
}

func (r *RoomService) Ready(s *session.Session, _ *proto.Ready) error {
	table, err := r.getTable(s)
	if err != nil {
		return err
	}
	return table.Ready(s)
}

func (r *RoomService) LoadRes(s *session.Session, msg *proto.LoadRes) error {
	table, err := r.getTable(s)
	if err != nil {
		return err
	}
	return table.LoadRes(s, msg)
}

func (r *RoomService) Update(s *session.Session, msg *proto.UpdateFrame) error {
	table, err := r.getTable(s)
	if err != nil {
		return err
	}
	return table.Update(s, msg)
}

func (r *RoomService) ResumeTable(s *session.Session, msg *proto.ResumeTable) error {
	table, err := r.getTable(s)
	if err != nil {
		return err
	}
	return table.ResumeTable(s, msg)
}
