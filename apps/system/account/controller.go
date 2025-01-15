package account

import (
	"context"
	"fmt"
	"time"

	"github.com/kochabonline/kcloud/apps/common"
	"github.com/kochabonline/kit/core/crypto/bcrypt"
	"github.com/kochabonline/kit/core/util"
	"github.com/kochabonline/kit/log"
)

var _ Interface = (*Controller)(nil)

type Controller struct {
	repo *Repository
	log  log.Helper
}

func NewController(repo *Repository, log log.Helper) *Controller {
	return &Controller{
		repo: repo,
		log:  log,
	}
}

func (ctrl *Controller) Create(ctx context.Context, req *CreateRequest) error {
	// 先查询账户是否存在
	result, _ := ctrl.repo.FindByUsername(ctx, req.Username)
	if result.Id != 0 {
		return common.ErrorAccountExist
	}

	// 创建账户
	now := time.Now().Unix()
	hashedPassword, err := bcrypt.HashPassword(req.Password)
	if err != nil {
		ctrl.log.Errorw("message", "密码加密失败", "error", err.Error())
		return err
	}

	account := Account{
		Meta: common.Meta{
			CreatedAt: now,
		},
		Username: req.Username,
		Password: hashedPassword,
		Status:   common.StatusNormal,
	}

	if err := ctrl.repo.Create(ctx, account); err != nil {
		ctrl.log.Errorw("message", "账户创建失败", "error", err.Error())
		return err
	}

	return nil
}

func (ctrl *Controller) changeQuery(ctx context.Context, code string) (int64, error) {
	accountId, err := util.CtxValue[int64](ctx, "id")
	if err != nil {
		ctrl.log.Errorw("message", "从上下文获取Id失败", "error", err.Error())
		return 0, err
	}

	// 查询账户是否存在
	result, err := ctrl.repo.FindById(ctx, accountId)
	if err != nil {
		ctrl.log.Errorw("message", "查询账户失败", "error", err.Error())
		return 0, err
	}
	if result.Id == 0 {
		ctrl.log.Errorw("message", "账户不存在", "accountId", accountId)
		return 0, common.ErrorAccountNotExist
	}

	// 查询安全码是否正确
	securityCode, err := ctrl.repo.rediscli.Get(ctx, fmt.Sprintf(common.SecurityCodeKey, accountId)).Result()
	if err != nil {
		ctrl.log.Errorw("message", "查询安全码失败", "error", err.Error())
		return 0, err
	}
	if securityCode != code {
		ctrl.log.Errorw("message", "安全码错误", "src", code, "dst", securityCode)
		return 0, common.ErrorSecurityCode
	}

	return accountId, nil
}

func (ctrl *Controller) ChangePassword(ctx context.Context, req *ChangePasswordRequest) error {
	id, err := ctrl.changeQuery(ctx, req.SecurityCode)
	if err != nil {
		return err
	}

	// 修改密码
	if err := ctrl.repo.ChangePassword(ctx, id, req.Password); err != nil {
		ctrl.log.Errorw("message", "修改密码失败", "error", err.Error())
		return err
	}

	return nil
}

func (ctrl *Controller) ChangeEmail(ctx context.Context, req *ChangeEmailRequest) error {
	id, err := ctrl.changeQuery(ctx, req.SecurityCode)
	if err != nil {
		return err
	}

	// 修改邮箱
	if err := ctrl.repo.ChangeEmail(ctx, id, req.Email); err != nil {
		ctrl.log.Errorw("message", "修改邮箱失败", "error", err.Error())
		return err
	}

	return nil
}

func (ctrl *Controller) ChangeMobile(ctx context.Context, req *ChangeMobileRequest) error {
	id, err := ctrl.changeQuery(ctx, req.SecurityCode)
	if err != nil {
		return err
	}

	// 修改手机号
	if err := ctrl.repo.ChangeMobile(ctx, id, req.Mobile); err != nil {
		ctrl.log.Errorw("message", "修改手机号失败", "error", err.Error())
		return err
	}
	return nil
}

func (ctrl *Controller) FindAll(ctx context.Context, req *FindAllRequest) (*Accounts, error) {
	accounts, err := ctrl.repo.FindAll(ctx, req)
	if err != nil {
		ctrl.log.Errorw("message", "查询账户列表失败", "error", err.Error())
		return nil, err
	}
	return accounts, nil
}

func (ctrl *Controller) FindById(ctx context.Context, id int64) (*Account, error) {
	account, err := ctrl.repo.FindById(ctx, id)
	if err != nil {
		ctrl.log.Errorw("message", "Id查询账户失败", "error", err.Error())
		return &Account{}, err
	}
	return account, nil
}

func (ctrl *Controller) FindByUsername(ctx context.Context, username string) (*Account, error) {
	account, err := ctrl.repo.FindByUsername(ctx, username)
	if err != nil {
		ctrl.log.Errorw("message", "用户名查询账户失败", "error", err.Error())
		return &Account{}, err
	}
	return account, nil
}

func (ctrl *Controller) Detail(ctx context.Context) (*Account, error) {
	accountId, err := util.CtxValue[int64](ctx, "id")
	if err != nil {
		ctrl.log.Errorw("message", "从上下文获取Id失败", "error", err.Error())
		return &Account{}, err
	}

	account, err := ctrl.repo.FindById(ctx, accountId)
	if err != nil {
		ctrl.log.Errorw("message", "Id查询账户失败", "error", err.Error())
		return &Account{}, err
	}

	// 脱敏处理
	account.Mobile = util.MobileDesensitization(account.Mobile)
	account.Email = util.EmailDesensitization(account.Email)

	return account, nil
}

func (ctrl *Controller) Delete(ctx context.Context, req *DeleteRequest) error {
	// 查询账户是否存在
	result, err := ctrl.repo.FindById(ctx, req.Id)
	if err != nil {
		ctrl.log.Errorw("message", "查询账户失败", "error", err.Error())
		return err
	}
	if result.Id == 0 {
		ctrl.log.Errorw("message", "账户不存在", "accountId", req.Id)
		return common.ErrorAccountNotExist
	}

	// 删除账户
	if err := ctrl.repo.Delete(ctx, req.Id); err != nil {
		ctrl.log.Errorw("message", "删除账户失败", "error", err.Error())
		return err
	}

	return nil
}
