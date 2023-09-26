package test

import (
	"bytes"
	"fmt"
	"github.com/lonng/nano/benchmark/wsio"
	"github.com/lonng/nano/serialize/protobuf"
	"io/ioutil"
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

	url := "http://127.0.0.1:8000/v1/login"
	aa := &proto2.AccountLoginReq{
		Partition: consts.DEVICEID,
		AccountId: deviceId,
	}
	input, _ := ss.Marshal(aa)
	req, err := http.NewRequest("POST", url, bytes.NewReader(input))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-protobuf")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Sprintf("Error http client do %+v", err)
		return
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Sprintf("Error reading response body: %v", err)
		return
	}
	respMsg := &proto2.AccountLoginResp{}
	if err := ss.Unmarshal(respData, respMsg); err != nil {
		fmt.Sprintf("Error unmarshaling response: %v", err)
		return
	}

	// 长链接
	c := wsio.NewConnector()
	chReady := make(chan struct{})
	c.OnConnected(func() {
		chReady <- struct{}{}
	})

	path := respMsg.Addr
	println(path, respMsg.UserId)
	if err := c.Start(path, "/nano"); err != nil {
		panic(err)
	}
	<-chReady

	chLogin := make(chan struct{})
	chEnd := make(chan interface{}, 0)

	var teamId int32
	uid := respMsg.UserId
	state := consts.IDLE
	if uid == 0 {
		c.Request("g.register", &proto2.RegisterGameReq{Name: deviceId, AccountId: aa.AccountId}, func(data interface{}) {
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

	// 倒计时
	c.On("onCountDown", func(data interface{}) {
		v := proto2.CountDownResp{}
		ss.Unmarshal(data.([]byte), &v)
		fmt.Println(deviceId, "onCountDown", v.Counter)
	})

	// 状态变化
	c.On("onState", func(data interface{}) {
		v := proto2.GameStateResp{}
		ss.Unmarshal(data.([]byte), &v)
		state = v.State
		tableInfo := v.TableInfo

		if v.TableInfo != nil {
			for k, v := range tableInfo.Players {
				if k == uid {
					teamId = v.TeamId
				}
			}
			//tableId = tableInfo.TableId
		}
		fmt.Println(deviceId, "onState", v.State)
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

	for /*i := 0; i < 1; i++*/ {

		select {
		case <-chEnd:
			fmt.Println("游戏结束了", uid)
			return
		case <-ticker.C:
			switch state {
			case consts.IDLE:
				c.Request("r.join", &proto2.Join{
					RoomId: rid,
				}, func(data interface{}) {
					v := proto2.GameStateResp{}
					ss.Unmarshal(data.([]byte), &v)
					state = v.State
					fmt.Println(deviceId, "join", state)
				})
			case consts.WAITREADY:
				c.Notify("r.ready", &proto2.Ready{})

			case consts.GAMING:
				m := make([]*proto2.Array, 0)
				m = append(m, &proto2.Array{V: append(make([]uint32, 0), 0, 1, 2, 3, 4, 5), I: 0}, &proto2.Array{V: append(make([]uint32, 0), 0, 1, 2, 3, 4, 5), I: 1})

				c.Request("r.updatestate", &proto2.UpdateState{
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
			case consts.SETTLEMENT:
				state = consts.IDLE
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

func TestGame(t *testing.T) {

	// wait server startup
	wg := sync.WaitGroup{}
	for i := 0; i < 1; i++ {
		wg.Add(1)
		time.Sleep(50 * time.Millisecond)
		go func(index int) {
			client(fmt.Sprintf("robot%d", index), "1", wg)
		}(i)
	}

	wg.Wait()

	t.Log("exit")
}

func TestWeb(t *testing.T) {
	al := proto2.AccountLoginReq{
		Partition: 0,
		AccountId: "11111",
	}

	s := protobuf.NewSerializer()
	data, err := s.Marshal(&al)
	if err != nil {
		fmt.Println("创建失败", err)
	}
	c := &http.Client{}

	url := "http://127.0.0.1:8000/v1/login"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("请求失败", err)
	}
	req.Header.Set("Content-Type", "application/x-protobuf")
	resp, err := c.Do(req)
	if err != nil {
		fmt.Println("发送请求失败", err)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Sprintf("Error reading response body: %v", err)
	}
	respMsg := &proto2.AccountLoginResp{}
	if err := s.Unmarshal(respData, respMsg); err != nil {
		fmt.Sprintf("Error unmarshaling response: %v", err)
	}
	defer resp.Body.Close()
}
