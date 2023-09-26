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
	UpdateState(s *session.Session, msg *proto.UpdateState) error
	ResumeTable(s *session.Session) error
	GetConfig() *config.Room
	CreateTable(sList []*session.Session)
	Ready(s *session.Session) error
	BackToQueue(sList []*session.Session)
	Entity(tableId string) (TableEntity, error)
}

type TableEntity interface {
	AfterInit()
	GetTableId() string
	BroadcastTableState(s int32)
	UpdateState(s *session.Session, msg *proto.UpdateState) error
	Clear()
	Leave(s *session.Session)
	Join(s *session.Session, tableId string) error
	Ready(s *session.Session) error
	CanDestory() bool
	GetInfo() *proto.TableInfo
	BackToTable()
	Entity() (WaiterEntity, error)
}

type WaiterEntity interface {
	Ready(s *session.Session) error
	CheckAndDismiss()
	AfterInit()
}
