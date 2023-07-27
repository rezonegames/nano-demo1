package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lonng/nano/benchmark/io"
	"github.com/lonng/nano/serialize/protobuf"
	"net/http"
	"sync"
	"testing"
	"tetris/consts"
	"tetris/pkg/z"
	proto2 "tetris/proto/proto"
	"time"
)

func client(deviceId, rid string, wg sync.WaitGroup) {
	defer wg.Done()
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
	var tableId string
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
		c.Request("g.login", &proto2.LoginToGame{UserId: uid}, func(data interface{}) {
			chLogin <- struct{}{}
			v := proto2.LoginToGameResp{}
			ss.Unmarshal(data.([]byte), &v)
			fmt.Println(deviceId, "login", v.Player.Name, v.Player.UserId, v.Player.Coin)
		})
	}
	<-chLogin
	c.On("pong", func(data interface{}) {})

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
		tableId = tableInfo.TableId
		fmt.Println(deviceId, "onStateChange", v.State)
	})

	c.Request("r.quickstart", &proto2.Join{
		RoomId: rid,
	}, func(data interface{}) {
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

	c.On("onStateUpdate", func(data interface{}) {
		v := proto2.UpdateState{}
		ss.Unmarshal(data.([]byte), &v)
		fmt.Println(z.ToString(v))
	})

	ra := z.RandInt(6, 10)
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
			if state == consts.GAMING {

				m := make([]*proto2.Array, 0)
				m = append(m, &proto2.Array{V: append(make([]uint32, 0), 0, 1, 2, 3, 4, 5), I: 0}, &proto2.Array{V: append(make([]uint32, 0), 0, 1, 2, 3, 4, 5), I: 1})

				c.Request("r.updatestate", &proto2.UpdateState{
					RoomId:   rid,
					TableId:  tableId,
					PlayerId: uid,

					Fragment: "all",
					Player: &proto2.Player{
						Matrix: m,
						Pos: &proto2.Pos{
							X: 1,
							Y: 2,
						},
						Score: 10,
					},

					Arena: m,

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

//func TestMatrix(t *testing.T) {
//	//xx := `[[0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0]]`
//
//	m := make([]*proto2.Array, 0)
//	m = append(m, &proto2.Array{V: append(make([]uint32, 0), 0, 1, 2, 3, 4, 5), I: 0}, &proto2.Array{V: append(make([]uint32, 0), 0, 1, 2, 3, 4, 5), I: 1})
//
//	us := proto2.UpdateState{
//		Player: &proto2.Player{
//			Matrix: m,
//			Pos: &proto2.Pos{
//				X: 1,
//				Y: 2,
//			},
//			Score: 0,
//		},
//	}
//
//	//err := json.Unmarshal([]byte(uss), &us)
//	fmt.Println(z.ToString(us))
//}

func TestIO(t *testing.T) {

	// wait server startup
	wg := sync.WaitGroup{}
	for i := 0; i < 6; i++ {
		wg.Add(1)
		time.Sleep(50 * time.Millisecond)
		go func(index int) {
			client(fmt.Sprintf("test%d", index), "3", wg)
		}(i)
	}

	wg.Wait()

	t.Log("exit")
}
