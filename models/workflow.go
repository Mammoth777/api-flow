package models

import (
	"github.com/jinzhu/gorm"
)

type WorkflowStatus uint8

const (
	WorkflowDraft WorkflowStatus = iota
	WorkflowPublished
)

// Workflow 流程模型
type Workflow struct {
	BasicModel
	Name        string         `gorm:"size:255;not null" json:"name"`
	Description string         `gorm:"size:1000" json:"description"`
	Status      WorkflowStatus `json:"status"`
}

// TableName 指定表名
func (Workflow) TableName() string {
	return "workflows"
}

// MigrateWorkflow 创建流程表
func MigrateWorkflow(db *gorm.DB) error {
	return db.AutoMigrate(&Workflow{}).Error
}
