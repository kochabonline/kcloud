package common

type Meta struct {
	// 自增主键
	Id int64 `json:"id" gorm:"primaryKey;autoIncrement;comment:自增主键"`
	// 创建时间
	CreatedAt int64 `json:"createdAt" gorm:"not null;comment:创建时间"`
	// 更新时间
	UpdatedAt int64 `json:"updatedAt" gorm:"comment:更新时间"`
	// 删除时间
	DeletedAt int64 `json:"deletedAt" gorm:"index:idx_deleted_at;comment:删除时间"`
}
