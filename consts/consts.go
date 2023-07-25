package consts

// 配置文件名字
const (
	ROOM_CONFIG = "rooms"
)

const (
	PLAYER_KEY = "player"
)

// 第三方登录渠道
const (
	DEVICEID = iota
	WX
	FB
	GIT
)

// 游戏状态
const (
	IDLE int32 = iota
	WAIT
	COUNTDOWN
	GAMEING
	SETTLEMENT
)
