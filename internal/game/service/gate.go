package service

import (
	"github.com/lonng/nano"
	"github.com/lonng/nano/component"
	"github.com/lonng/nano/session"
	"tetris/config"
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

func (g *GateService) online(s *session.Session, uid int64) error {
	if ps, err := g.group.Member(uid); err == nil {
		log.Info("close 老的连接", ps.ID())
		ps.Close()
	}
	log.Info("玩家: %d上线", uid)
	g.group.Add(s)
	return s.Bind(uid)
}

func (g *GateService) Register(s *session.Session, req *proto.RegisterGameReq) error {
	if req.AccountId == "" {
		return s.Response(proto.LoginToGameResp{Code: proto.ErrorCode_AccountIdError})
	}
	pp := models.NewPlayer(req.Name, 100)
	uid, err := models.CreatePlayer(pp)
	if err != nil {
		return err
	}
	err = models.BindAccount(req.AccountId, uid)
	if err != nil {
		return err
	}
	g.online(s, uid)

	return s.Response(&proto.LoginToGameResp{Player: util.ConvPlayerToProtoPlayer(pp)})
}

func (g *GateService) Login(s *session.Session, req *proto.LoginToGame) error {
	uid := req.UserId
	pp, err := models.GetPlayer(uid)
	if err != nil {
		return err
	}
	g.online(s, uid)
	return s.Response(&proto.LoginToGameResp{Player: util.ConvPlayerToProtoPlayer(pp)})
}

func (g *GateService) rooms(s *session.Session, _ interface{}) error {
	return s.Response(&proto.GetRoomListResp{RoomList: util.ConvRoomListToProtoRoomList(config.ServerConfig.Rooms)})
}

func (g *GateService) Ping(s *session.Session, _ *proto.Ping) error {
	return s.Response(&proto.Pong{TimeStamp: z.NowUnixMilli()})
}
