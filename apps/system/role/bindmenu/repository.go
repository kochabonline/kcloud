package bindmenu

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

// Edit 编辑绑定关系
func (repo *Repository) Edit(ctx context.Context, roleId int64, menuIds []int64) error {
	// 先查询出已经绑定的菜单
	var roleBindMenus []RoleBindMenu
	repo.db.Where("role_id = ? AND menu_id IN ?", roleId, menuIds).Find(&roleBindMenus)

	// 构造map方便查找
	bindMenuMap := make(map[int64]struct{})
	for _, bindMenu := range roleBindMenus {
		bindMenuMap[bindMenu.MenuId] = struct{}{}
	}

	// 需要新增的菜单
	var add []int64
	// 需要删除的菜单
	var del []int64

	// 遍历新传入的菜单，找出需要新增的菜单
	for _, menuId := range menuIds {
		if _, found := bindMenuMap[menuId]; !found {
			add = append(add, menuId)
		} else {
			// 如果菜单已经存在于绑定菜单中，则从 map 中删除
			delete(bindMenuMap, menuId)
		}
	}

	// 剩下的 bindMenuMap 中的菜单即为需要删除的菜单
	for menuId := range bindMenuMap {
		del = append(del, menuId)
	}

	var query = repo.db.WithContext(ctx)
	// 原生sql构造
	var values []any
	insertSql := "INSERT INTO bind_menu (role_id, menu_id) VALUES "
	for _, mid := range add {
		insertSql += "(?, ?),"
		values = append(values, roleId, mid)
	}
	if len(add) > 0 {
		insertSql = insertSql[:len(insertSql)-1]
		if err := query.Exec(insertSql, values...).Error; err != nil {
			return err
		}
	}

	if len(del) > 0 {
		if err := query.Where("role_id = ? AND menu_id IN ?", roleId, del).Delete(&RoleBindMenu{}).Error; err != nil {
			return err
		}
	}

	return nil
}
