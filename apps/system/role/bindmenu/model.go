package bindmenu

type RoleBindMenu struct {
	RoleId int64 `json:"role_id" gorm:"type:bigint;not null;comment:角色Id"`
	MenuId int64 `json:"menu_id" gorm:"type:bigint;not null;comment:菜单Id"`
}

func (RoleBindMenu) TableName() string {
	return "role_bind_menu"
}
