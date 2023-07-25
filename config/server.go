package config

type Server struct {
	AppName         string        `yaml:"appName"`
	IsDebug         bool          `yaml:"isDebug"`
	TokenExpireTime int64         `yaml:"tokenExpireTime"`
	AppSecret       string        `yaml:"appSecret"`
	TokenSecret     string        `yaml:"tokenSecret"`
	Addr            string        `yaml:"addr"`
	ServerPort      string        `yaml:"serverPort"`
	EnableChecksum  bool          `yaml:"enableChecksum"`
	Redis           *Redis        `yaml:"redis"`
	Mongo           *Mongo        `yaml:"mongo"`
	GmIps           []interface{} `yaml:"gmIps"`
	Version         string        `yaml:"version"`
	Rooms           []*Room       `yaml:"rooms"`
}

type Redis struct {
	Host string `yaml:"host"`
	Db   int    `yaml:"db"`
}

type Mongo struct {
	Uri string `yaml:"uri"`
	Db  string `yaml:"db"`
}

type Room struct {
	RoomId  string `yaml:"roomId"`
	Pvp     int32  `yaml:"pvp"`
	Divide  int32  `yaml:"divide"`
	Name    string `yaml:"name"`
	MinCoin int32  `yaml:"minCoin"`
}
