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
	GetConfig() *config.Room
	CreateTable(sList []*session.Session)
	BackToWait(sList []*session.Session)
	Entity(tableId string) (TableEntity, error)
}

type TableEntity interface {
	AfterInit()
	GetTableId() string
	CanDestory() bool
	Clear()
	GetInfo() *proto.TableInfo
	BackToTable()
	WaiterEntity() WaiterEntity
	Entity(uid int64) ClientEntity
	ChangeState(state proto.TableState)
	Ready(s *session.Session) error
	LoadRes(s *session.Session, msg *proto.LoadRes) error
	Update(s *session.Session, msg *proto.UpdateFrame) error
	ResumeTable(s *session.Session, msg *proto.ResumeTable) error
	Leave(s *session.Session) error
	Join(s *session.Session, tableId string) error
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
	GetSession() *session.Session
	SetSession(s *session.Session)
	GetId() int64
	GetTeamId() int32
	IsEnd() bool
	// 帧相关的数据
	SaveFrame(frameId int64, msg *proto.UpdateFrame)
	GetFrame(frameId int64) []*proto.Action
	SetLastFrame(frameId int64)
	GetLastFrame() int64
}
