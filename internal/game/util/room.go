package util

import (
	"github.com/lonng/nano/session"
	"tetris/config"
	"tetris/proto/proto"
)

type RoomEntity interface {
	AfterInit()
	Leave(s *session.Session) error
	Join(s *session.Session) error
	Update(s *session.Session, msg *proto.UpdateFrame) error
	ResumeTable(s *session.Session) error
	GetConfig() *config.Room
	CreateTable(sList []*session.Session)
	Ready(s *session.Session) error
	BackToQueue(sList []*session.Session)
	Entity(tableId string) (TableEntity, error)
	LoadRes(s *session.Session, msg *proto.LoadRes) error
}

type TableEntity interface {
	AfterInit()
	GetTableId() string
	Update(s *session.Session, msg *proto.UpdateFrame) error
	Leave(s *session.Session)
	Join(s *session.Session, tableId string) error
	Ready(s *session.Session) error
	CanDestory() bool
	Clear()
	GetInfo() *proto.TableInfo
	BackToTable()
	Entity() (WaiterEntity, error)
	GetClient(uid int64) ClientEntity
	ChangeState(state proto.TableState)
	LoadRes(s *session.Session, msg *proto.LoadRes) error
}

type WaiterEntity interface {
	Ready(s *session.Session) error
	Leave(s *session.Session) error
	CheckAndDismiss()
	GetInfo() *proto.TableInfo_Waiter
	AfterInit()
}

type ClientEntity interface {
	GetPlayer() *proto.TableInfo_Player
	GetTeamId() int32
	IsEnd() bool
	SaveFrame(frameId int64, msg *proto.UpdateFrame)
	GetFrame(frameId int64) []*proto.Action
	GetSession() *session.Session
	GetId() int64
}
