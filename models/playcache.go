package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type RoundSession struct {
	RoomId  string `json:"rid"`
	TableId string `json:"tid"`
}

func GetUserRoundSession(uid int64) (*RoundSession, error) {
	v, err := rclient.Get(fmt.Sprintf("rs:%d", uid)).Result()
	if err != nil {
		return nil, err
	}
	var rs RoundSession
	err = json.Unmarshal([]byte(v), &rs)
	if err != nil {
		return nil, err
	}
	return &rs, nil
}

func SetUserRoundSession(uid int64, rs *RoundSession) error {
	v, err := json.Marshal(rs)
	if err != nil {
		return err
	}
	return rclient.Set(fmt.Sprintf("rs:%d", uid), v, 24*time.Hour).Err()
}
