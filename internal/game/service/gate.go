package service

import (
	"github.com/lonng/nano"
	"github.com/lonng/nano/component"
	"github.com/lonng/nano/session"
	"tetris/internal/game/util"
	"tetris/models"
	"tetris/pkg/log"
	"tetris/pkg/z"
	"tetris/proto/proto"
)

type GateService struct {
	*component.Base
	group *nano.Group
}

func NewGateService() *GateService {
	return &GateService{group: nano.NewGroup("all-users")}
}

func (g *GateService) AfterInit() {
	// Fixed: 玩家WIFI切换到4G网络不断开, 重连时，将UID设置为illegalSessionUid
	session.Lifetime.OnClosed(func(s *session.Session) {
		if s.UID() > 0 {
			if err := g.offline(s); err != nil {
				log.Error("玩家退出: UID=%d, Error=%s", s.UID, err.Error())
			}
		}
	})
}

func (g *GateService) offline(s *session.Session) error {
	return g.group.Leave(s)
}

func (g *GateService) online(s *session.Session, uid int64) (*models.Profile, error) {

	if ps, err := g.group.Member(uid); err == nil {
		log.Info("close 老的连接 %d", ps.UID())
		g.group.Leave(ps)
		ps.Close()
	}

	log.Info("玩家: %d上线", uid)
	p, err := models.GetPlayer(uid)
	if err != nil {
		return nil, err
	}
	util.BindUser(s, p)
	return p, g.group.Add(s)
}

func (g *GateService) Register(s *session.Session, req *proto.RegisterGameReq) error {
	if req.AccountId == "" {
		return s.Response(proto.LoginToGameResp{Code: proto.ErrorCode_AccountIdError})
	}
	p := models.NewPlayer(req.Name, 100)
	uid, err := models.CreatePlayer(p)
	if err != nil {
		return err
	}
	err = models.BindAccount(req.AccountId, uid)
	if err != nil {
		return err
	}
	_, err = g.online(s, uid)
	if err != nil {
		return err
	}
	return s.Response(&proto.LoginToGameResp{
		Code:     proto.ErrorCode_OK,
		Profile:  util.ConvProfileToProtoProfile(p),
		RoomList: util.GetRoomList(),
	})
}

func (g *GateService) Login(s *session.Session, req *proto.LoginToGame) error {
	uid := req.UserId
	p, err := g.online(s, uid)
	if err != nil {
		return err
	}
	// 返回所在的房间号和桌子号
	var roomId, tableId string
	if rs, err := models.GetRoundSession(uid); err == nil {
		roomId = rs.RoomId
		tableId = rs.TableId
	}

	return s.Response(&proto.LoginToGameResp{
		Code:     proto.ErrorCode_OK,
		Profile:  util.ConvProfileToProtoProfile(p),
		RoomList: util.GetRoomList(),
		RoomId:   roomId,
		TableId:  tableId,
	})
}

func (g *GateService) Rooms(s *session.Session, _ interface{}) error {
	return s.Response(&proto.GetRoomListResp{RoomList: util.GetRoomList()})
}

func (g *GateService) Ping(s *session.Session, _ *proto.Ping) error {
	return s.Response(&proto.Pong{Ts: z.NowUnixMilli()})
}
