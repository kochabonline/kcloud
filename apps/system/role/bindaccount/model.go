package bindaccount

type RoleBindAccount struct {
	Id        int64 `json:"id" gorm:"primaryKey;autoIncrement;comment:主键"`
	AccountId int64 `json:"account_id" gorm:"type:bigint;not null;comment:账户id"`
	RoleId    int64 `json:"role_id" gorm:"type:bigint;not null;comment:角色id"`
}

func (RoleBindAccount) TableName() string {
	return "role_bind_account"
}
