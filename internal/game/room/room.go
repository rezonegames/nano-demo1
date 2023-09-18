package room

import (
	"fmt"
	"github.com/lonng/nano/session"
	"tetris/config"
	"tetris/consts"
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
	BackToTable([]*session.Session)
}

type WaiterEntity interface {
	Ready(s *session.Session) error
	GetId() string
	CheckAndDismiss()
	AfterInit()
}

// NewWaiter waiter的设计是等待的模式，比如王者荣耀，斗地主满人就开等等
func NewWaiter(sList []*session.Session, room RoomEntity) WaiterEntity {
	var w WaiterEntity
	conf := room.GetConfig()
	switch conf.RoomType {
	case consts.QUICK:
		w = NewNormalWaiter(sList, room)
	default:
		panic(fmt.Sprintf("NewWaiter unknown room type %s", conf.RoomType))
	}
	w.AfterInit()
	return w
}

// NewTable 根据桌子类型创建道具桌还是正常桌
func NewTable(room RoomEntity, ss []*session.Session) TableEntity {
	conf := room.GetConfig()
	var t TableEntity
	switch conf.TableType {
	case consts.NORMAL:
		t = NewNormalTable(room, ss)
	default:
		panic(fmt.Sprintf("NewTable unknown type %s", conf.TableType))
	}
	t.AfterInit()
	return t
}

// NewRoom 根据房间类型创建房间
func NewRoom(conf *config.Room) RoomEntity {
	var r RoomEntity
	switch conf.RoomType {
	case consts.QUICK:
		r = NewQuickRoom(conf)
	default:
		panic(fmt.Sprintf("NewRoom unknown type %s", conf.RoomType))
	}
	r.AfterInit()
	return r
}
