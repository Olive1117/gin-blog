package database

import (
	"context"
	"reflect"

	"github.com/Olive1117/gin-blog/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuditPlugin struct{}

func (op *AuditPlugin) Name() string {
	return "audit_plugin"
}

func (op *AuditPlugin) Initialize(db *gorm.DB) error {
	var err error
	// 注册创建钩子
	err = db.Callback().Create().Before("gorm:create").Register("audit:before_create", op.beforeCreate)
	// 注册更新钩子
	err = db.Callback().Update().Before("gorm:update").Register("audit:before_update", op.beforeUpdate)
	// 注册删除钩子（针对软删除）
	err = db.Callback().Delete().Before("gorm:delete").Register("audit:before_delete", op.beforeDelete)
	if err != nil {
		logger.L.Error("grom callback", zap.Error(err))
		return err
	}
	return nil
}

func (op *AuditPlugin) beforeCreate(db *gorm.DB) {
	if db.Statement.Schema == nil {
		return
	}
	id, ok := db.Statement.Context.Value(kUserID).(int64)
	if !ok {
		logger.L.Error("获取用户ID失败")
		return
	}
	logger.L.Debug("执行创建插件", zap.Int64("当前用户", id))
	// 设置 CreatedBy 和 UpdatedBy
	createdBy := db.Statement.Schema.LookUpField("CreatedBy")
	updatedBy := db.Statement.Schema.LookUpField("UpdatedBy")
	switch db.Statement.ReflectValue.Kind() {
	case reflect.Slice, reflect.Array:
		// 批量创建：必须循环处理每一行
		for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
			rv := reflect.Indirect(db.Statement.ReflectValue.Index(i))
			if createdBy != nil {
				createdBy.Set(db.Statement.Context, rv, id)
			}
			if updatedBy != nil {
				updatedBy.Set(db.Statement.Context, rv, id)
			}
		}
	case reflect.Struct:
		// 单条创建
		rv := db.Statement.ReflectValue
		if createdBy != nil {
			createdBy.Set(db.Statement.Context, rv, id)
		}
		if updatedBy != nil {
			updatedBy.Set(db.Statement.Context, rv, id)
		}
	}
}

func (op *AuditPlugin) beforeUpdate(db *gorm.DB) {
	if db.Statement.Schema == nil {
		return
	}
	id, ok := db.Statement.Context.Value(kUserID).(int64)
	if !ok {
		logger.L.Error("获取用户ID失败")
		return
	}
	logger.L.Debug("执行更新插件", zap.Int64("id", id))
	updatedBy := db.Statement.Schema.LookUpField("UpdatedBy")
	switch db.Statement.ReflectValue.Kind() {
	case reflect.Slice, reflect.Array:
		// 批量创建：必须循环处理每一行
		for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
			rv := reflect.Indirect(db.Statement.ReflectValue.Index(i))
			if updatedBy != nil {
				updatedBy.Set(db.Statement.Context, rv, id)
			}
		}
	case reflect.Struct:
		// 单条创建
		rv := db.Statement.ReflectValue
		if updatedBy != nil {
			updatedBy.Set(db.Statement.Context, rv, id)
		}
	}
}

func (op *AuditPlugin) beforeDelete(db *gorm.DB) {
	id, ok := db.Statement.Context.Value(kUserID).(int64)
	if !ok {
		logger.L.Error("获取用户ID失败")
		return
	}
	logger.L.Debug("执行删除插件", zap.Int64("id", id))
	// 逻辑删除本质是更新，手动设置字段
	if field := db.Statement.Schema.LookUpField("DeletedBy"); field != nil {
		logger.L.Debug("删除！")
		db.Statement.SetColumn("DeletedBy", id)
	}
}

type userIDKey struct{}

var kUserID = userIDKey{}

func SetUserID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, kUserID, id)
}
func GetUserID(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(kUserID).(int64)
	if ok {
		return id, true
	}
	return 0, false
}
