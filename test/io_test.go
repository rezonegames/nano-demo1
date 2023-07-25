package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lonng/nano/benchmark/io"
	"github.com/lonng/nano/serialize/protobuf"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"tetris/consts"
	"tetris/pkg/z"
	proto2 "tetris/proto/proto"
	"time"
)

func client(deviceId, rid string) {
	ss := protobuf.NewSerializer()

	url := "http://127.0.0.1:8000/v1/user/login/query"
	aa := &proto2.AccountLoginReq{
		Partition: 0,
		AccountId: deviceId,
	}
	input, _ := json.Marshal(aa)
	req, err := http.NewRequest("POST", url, bytes.NewReader(input))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	var loginResp proto2.AccountLoginResp
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	if err != nil {
		panic(err)
	}

	// 长链接
	c := io.NewConnector()
	chReady := make(chan struct{})
	c.OnConnected(func() {
		chReady <- struct{}{}
	})

	path := loginResp.Host + loginResp.Port
	println(path, loginResp.UserId, loginResp.AccountId)
	if err := c.Start(path); err != nil {
		panic(err)
	}
	<-chReady

	chLogin := make(chan struct{})
	chEnd := make(chan interface{}, 0)

	var teamId int32
	uid := loginResp.UserId
	state := consts.IDLE
	if uid == 0 {
		c.Request("g.register", &proto2.RegisterGameReq{Name: deviceId, AccountId: loginResp.AccountId}, func(data interface{}) {
			chLogin <- struct{}{}

			v := proto2.LoginToGameResp{}
			ss.Unmarshal(data.([]byte), &v)
			ss.Unmarshal(data.([]byte), &v)
			fmt.Println(deviceId, "register", v.Player.Name, v.Player.UserId, v.Player.Coin)
		})
	} else {
		c.Request("g.login", &proto2.LoginToGameReq{UserId: uid}, func(data interface{}) {
			chLogin <- struct{}{}
			v := proto2.LoginToGameResp{}
			ss.Unmarshal(data.([]byte), &v)
			fmt.Println(deviceId, "login", v.Player.Name, v.Player.UserId, v.Player.Coin)
		})
	}
	<-chLogin
	//c.On("pong", func(data interface{}) {})

	// 倒计时
	c.On("onCountDown", func(data interface{}) {
		v := proto2.CountDownResp{}
		ss.Unmarshal(data.([]byte), &v)
		fmt.Println(deviceId, "onCountDown", v.Counter)
	})

	// 状态变化
	c.On("onStateChange", func(data interface{}) {
		v := proto2.GameStateResp{}
		ss.Unmarshal(data.([]byte), &v)
		state = v.State
		tableInfo := v.TableInfo

		for k, v := range tableInfo.Players {
			if k == uid {
				teamId = v.TeamId
			}
		}
		fmt.Println(deviceId, "onStateChange", v.State)
	})

	c.Request(rid+".quickstart", &proto2.Empty{}, func(data interface{}) {
		v := proto2.GameStateResp{}
		ss.Unmarshal(data.([]byte), &v)
		fmt.Println(deviceId, "quickstart", v.State)
	})

	c.On("onTeamLose", func(data interface{}) {
		v := proto2.TeamLose{}
		ss.Unmarshal(data.([]byte), &v)

		if v.TeamId == teamId {
			chEnd <- struct{}{}
		}
	})

	ra := z.RandInt(6, 20)
	ticker := time.NewTicker(time.Duration(ra) * time.Second)
	defer ticker.Stop()
	ticker1 := time.NewTicker(time.Second * 5)
	defer ticker1.Stop()
	for /*i := 0; i < 1; i++*/ {

		select {
		case <-chEnd:
			fmt.Println("游戏结束了", uid)
			return

		case <-ticker1.C:
			c.Request("g.ping", &proto2.Ping{}, func(data interface{}) {
				v := proto2.Pong{}
				ss.Unmarshal(data.([]byte), &v)
				//fmt.Println("pong", v.TimeStamp)
			})
		case <-ticker.C:
			if state == consts.GAMEING {
				c.Request(rid+".updatestate", &proto2.UpdateState{
					End: true,
					//Data: nil,
				}, func(data interface{}) {
					fmt.Println(deviceId, "syncmessage")
				})
			}
		default:

		}
	}
}

func TestIO(t *testing.T) {

	// wait server startup
	for i := 0; i < 6; i++ {
		time.Sleep(50 * time.Millisecond)
		go func(index int) {
			client(fmt.Sprintf("test%d", index), "r3")
		}(i)
	}

	log.SetFlags(log.LstdFlags | log.Llongfile)

	sg := make(chan os.Signal)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)

	<-sg

	t.Log("exit")
}
