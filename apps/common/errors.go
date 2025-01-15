package common

import "github.com/kochabonline/kit/errors"

var (
	ErrorInvalidParam            = errors.BadRequest("无效参数")
	ErrorAccountExist            = errors.BadRequest("账号已存在")
	ErrorAccountNotExist         = errors.BadRequest("账号不存在")
	ErrorAccountOrPassword       = errors.BadRequest("账户或密码错误")
	ErrorSecurityCode            = errors.BadRequest("安全码错误")
	ErrorAccountLocked           = errors.BadRequest("账号已被锁定")
	ErrorAccountDisabled         = errors.BadRequest("账号已被禁用")
	ErrorAccountDeleted          = errors.BadRequest("账号已被删除")
	ErrorChannelExists           = errors.BadRequest("通道名称已存在")
	ErrorChannelNotExist         = errors.BadRequest("通道不存在")
	ErrorChannelDetail           = errors.BadRequest("查询通道详情失败")
	ErrorMissAuthorizationHeader = errors.Unauthorized("缺少Authorization头")
	ErrorUnauthorized            = errors.Unauthorized("未授权")
	ErrorInvalidCaptcha          = errors.BadRequest("无效验证码")
	ErrorInvalidRole             = errors.BadRequest("无效的角色")
)
