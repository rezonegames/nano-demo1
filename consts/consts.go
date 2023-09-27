package consts

const (
	NONE = 0
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
	WAITREADY
	CANCEL
	COUNTDOWN
	GAMING
	SETTLEMENT
)

// 等待开始的子状态
const (
	WAITREADY_PROFILE int32 = iota + 1
	WAITREADY_COUNTDOWN
	WAITREADY_READYLIST
)

// 进入游戏倒计时
const (
	COUNTDOWN_BEGIN int32 = iota + 1
)

// 进入游戏
const (
	GAME_BEGIN int32 = iota + 1
)

// 结算
const (
	SETTLEMENT_BEGIN int32 = iota + 1
	SETTLEMENT_END
)

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
