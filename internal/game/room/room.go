package room

import (
	"github.com/lonng/nano/session"
	"tetris/config"
	"tetris/consts"
	"tetris/pkg/z"
	"tetris/proto/proto"
)

type RoomEntity interface {
	run()
	Leave(s *session.Session) error
	Join(s *session.Session) error
	UpdateState(s *session.Session, msg *proto.UpdateState) error
	ResumeTable(s *session.Session) error
}

type TableEntity interface {
	getTableId() string
	run()
	broadcastTableState(s int32)
	updateState(s *session.Session, msg *proto.UpdateState) error
	clear()
	leave(s *session.Session)
	join(s *session.Session, tableId string) error
	canDestory() bool
	getInfo() *proto.TableInfo
}

// NewTable 根据桌子类型创建道具桌还是正常桌
func NewTable(conf *config.Room, ss []*session.Session) (TableEntity, error) {
	switch conf.TableType {
	case consts.NORMAL:
		return newNormalTable(conf, ss), nil
	default:
		return nil, z.NilError{Msg: conf.RoomId}
	}
}

// NewRoom 根据房间类型创建房间
func NewRoom(conf *config.Room) (RoomEntity, error) {
	switch conf.RoomType {
	case consts.QUICK:
		return newQuickRoom(conf), nil
	default:
		return nil, z.NilError{Msg: conf.RoomId}
	}
}
