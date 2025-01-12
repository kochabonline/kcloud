package channal

// 状态
type Status string

const (
	// 1:正常
	StatusNormal Status = "normal"
	// 2:禁用
	StatusDisabled Status = "disabled"
)

// 供应商类型
const (
	ProviderTypeEmail    string = "email"
	ProviderTypeDingTalk string = "dingtalk"
	ProviderTypeLark     string = "lark"
	ProviderTypeTelegram string = "telegram"
)

// 限速器类型
const (
	LimiterTypeSlidingWindow string = "sliding_window"
	LimiterTypeTokenBucket   string = "token_bucket"
)

// 通道限速器key
var (
	SlidingWindowKey = "chananl:sliding_window:%d"
	TokenBucketKey   = "chananl:token_bucket:%d"
)
