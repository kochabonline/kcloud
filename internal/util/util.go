package util

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"github.com/kochabonline/kit/core/tools"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 分页查询
func Paginate(page int, size int) (offset int, limit int) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	offset = (page - 1) * size
	limit = size
	return
}

// 从上下文中获取账户id,账户角色
func CtxAccount(ctx context.Context) (accountId int64, role int, err error) {
	if accountId, err = tools.CtxValue[int64](ctx, "id"); err != nil {
		return
	}
	if role, err = tools.CtxValue[int](ctx, "role"); err != nil {
		return
	}
	return
}

// 将bson.D反序列化为目标结构体, 目标结构体必须是指针
func BsonUnmarshal(src, dst any) error {
	doc, ok := src.(primitive.D)
	if !ok {
		return errors.New("src must be bson.D")
	}
	if reflect.ValueOf(dst).Kind() != reflect.Ptr {
		return errors.New("dst must be a pointer")
	}

	data, err := bson.Marshal(doc)
	if err != nil {
		return err
	}
	if err := bson.Unmarshal(data, dst); err != nil {
		return err
	}
	return nil
}

// 手机号脱敏
func MobileDesensitization(mobile string) string {
	length := len(mobile)
	if length < 7 {
		return mobile
	}
	return mobile[:3] + "****" + mobile[length-3:]
}

// 邮箱脱敏
func EmailDesensitization(email string) string {
	index := strings.IndexByte(email, '@')
	if index == -1 || index < 4 {
		return email
	}
	return email[:3] + "****" + email[index:]
}
