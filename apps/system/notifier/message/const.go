package message

const (
	// 待发送
	Pending string = "pending"
	// 发送成功
	Success string = "success"
	// 发送失败
	Failure string = "failure"
)

const (
	// 信息
	Info string = "info"
	// 警告
	Warning string = "warning"
	// 严重
	Critical string = "critical"
)

const (
	// 文本
	Text string = "text"
	// Markdown
	Markdown string = "markdown"
)

const (
	DefaultTitle string = "你有一条新消息"
)
