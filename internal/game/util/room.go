package util

import (
	"github.com/lonng/nano/session"
	"tetris/config"
	"tetris/proto/proto"
	"time"
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
	UpdateState(s *session.Session, msg *proto.UpdateState) error
	BroadcastPlayerState(uid int64, msg *proto.UpdateState) error
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
	Save(msg *proto.UpdateState)
	GetSession() *session.Session
	GetId() int64
	Update(t time.Time)
}
