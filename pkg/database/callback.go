package database

import (
	"github.com/Olive1117/gin-blog/pkg/contextutil"
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
	uid := contextutil.GetCurrentUser(db.Statement.Context)
	// 设置 CreatedBy 和 UpdatedBy
	if field := db.Statement.Schema.LookUpField("CreatedBy"); field != nil {
		db.Statement.SetColumn("CreatedBy", uid)
	}
	if field := db.Statement.Schema.LookUpField("UpdatedBy"); field != nil {
		db.Statement.SetColumn("UpdatedBy", uid)
	}
	// db.Statement.SetColumn("CreatedBy", uid)
	// db.Statement.SetColumn("UpdatedBy", uid)
}

func (op *AuditPlugin) beforeUpdate(db *gorm.DB) {
	uid := contextutil.GetCurrentUser(db.Statement.Context)
	if field := db.Statement.Schema.LookUpField("UpdatedBy"); field != nil {
		db.Statement.SetColumn("UpdatedBy", uid)
	}
	// db.Statement.SetColumn("UpdatedBy", uid)
}

func (op *AuditPlugin) beforeDelete(db *gorm.DB) {
	uid := contextutil.GetCurrentUser(db.Statement.Context)
	// 逻辑删除本质是更新，手动设置字段
	if field := db.Statement.Schema.LookUpField("DeletedBy"); field != nil {
		db.Statement.SetColumn("DeletedBy", uid)
	}
	// db.Statement.SetColumn("DeletedBy", uid)
}
