package consts

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
	WAITREADY
	CANCEL
	COUNTDOWN
	GAMING
	SETTLEMENT
)

// 等待开始的子状态
const (
	WAITREADY_PROFILE int32 = iota
	WAITREADY_COUNTDOWN
	WAITREADY_READYLIST
)

// 进入游戏倒计时
const (
	COUNTDOWN_BEGIN int32 = iota
)

// 进入游戏

// 房间类型
const (
	QUICK int32 = iota + 1
	MATHCH
)

// 桌子类型
const (
	NORMAL int32 = iota + 1
	HAPPY
)
