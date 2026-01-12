package database

import (
	"context"

	"github.com/Olive1117/gin-blog/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuditPlugin struct{}

type userIDKey struct{}

var kUserID = userIDKey{}

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
	id, ok := db.Statement.Context.Value(kUserID).(uint)
	if !ok {
		logger.L.Error("获取用户ID失败")
		return
	}
	logger.L.Debug("执行创建插件", zap.Uint("当前用户", id))
	// 设置 CreatedBy 和 UpdatedBy
	if field := db.Statement.Schema.LookUpField("CreatedBy"); field != nil {
		db.Statement.SetColumn("CreatedBy", id)
	}
	if field := db.Statement.Schema.LookUpField("UpdatedBy"); field != nil {
		db.Statement.SetColumn("UpdatedBy", id)
	}
}

func (op *AuditPlugin) beforeUpdate(db *gorm.DB) {
	id, ok := db.Statement.Context.Value(kUserID).(uint)
	if !ok {
		logger.L.Error("获取用户ID失败")
		return
	}
	logger.L.Debug("执行更新插件", zap.Uint("id", id))
	if field := db.Statement.Schema.LookUpField("UpdatedBy"); field != nil {
		db.Statement.SetColumn("UpdatedBy", id)
	}
}

func (op *AuditPlugin) beforeDelete(db *gorm.DB) {
	id, ok := db.Statement.Context.Value(kUserID).(uint)
	if !ok {
		logger.L.Error("获取用户ID失败")
		return
	}
	logger.L.Debug("执行删除插件", zap.Uint("id", id))
	// 逻辑删除本质是更新，手动设置字段
	if field := db.Statement.Schema.LookUpField("DeletedBy"); field != nil {
		logger.L.Debug("删除！")
		db.Statement.SetColumn("DeletedBy", id)
	}
}

func SetUserID(ctx context.Context, id uint) context.Context {
	return context.WithValue(ctx, kUserID, id)
}
