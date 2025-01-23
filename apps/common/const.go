package common

// 状态
type Status int

const (
	// 1:启用
	StatusNormal Status = iota + 1
	// 2:禁用
	StatusDisabled
	// 3:删除
	StatusDeleted
)

// 系统角色
const (
	// 1:普通账户
	RoleNormal string = "normal"
	// 2:管理员
	RoleAdmin string = "admin"
	// 3:超级管理员
	RoleSuperAdmin string = "super_admin"
)

// 安全设置安全码Key
var SecurityCodeKey = "security:code:%v"

// 验证header
var AuthorizationHeader = "Authorization"

// 语言header
var LanguageHeader = "Language"
