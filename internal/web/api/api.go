package api

import (
	"github.com/gorilla/mux"
	"github.com/lonng/nex"
	"net/http"
	"tetris/config"
	"tetris/consts"
	"tetris/models"
	"tetris/pkg/log"
	"tetris/pkg/z"
	proto2 "tetris/proto/proto"
)

// queryHandler 登录，如果玩家没有注册，自动注册之，只实现一种deviceId登录， todo：host暂时写死
func queryHandler(req *proto2.AccountLoginReq) (*proto2.AccountLoginResp, error) {
	log.Info("%v", req)
	partition := req.Partition
	accountId := req.AccountId

	var userId int64
	switch partition {
	case consts.DEVICEID:
		a, err := models.GetAccount(accountId)
		if err != nil {
			if _, ok := err.(z.NilError); ok {
				a = models.NewAccount(partition, accountId)
				err = models.CreateAccount(a)
				if err != nil {
					return &proto2.AccountLoginResp{Code: proto2.ErrorCode_DBError}, nil
				}
			} else {
				return &proto2.AccountLoginResp{Code: proto2.ErrorCode_DBError}, nil
			}
		}
		userId = a.UserId

	default:
		log.Error("queryHandler no account %d %s", partition, accountId)
		return &proto2.AccountLoginResp{Code: proto2.ErrorCode_UnknownError}, nil
	}

	sc := config.ServerConfig
	return &proto2.AccountLoginResp{
		UserId:    userId,
		Host:      "127.0.0.1",
		Port:      sc.Addr,
		AccountId: accountId,
	}, nil
}

func MakeLoginService() http.Handler {
	router := mux.NewRouter()
	router.Handle("/v1/user/login/query", nex.Handler(queryHandler)).Methods("POST")
	return router
}
