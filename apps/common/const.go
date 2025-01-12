package common

// 状态
type Status int

const (
	// 1:正常
	StatusNormal Status = iota + 1
	// 2:禁用
	StatusDisabled
	// 3:删除
	StatusDeleted
)

// 角色
type Role int

const (
	// 1:普通账户
	RoleUser Role = iota + 1
	// 2:管理员
	RoleAdmin
)

// 安全设置安全码Key
var SecurityCodeKey = "security:code:%v"

// 验证header
var AuthorizationHeader = "token"
