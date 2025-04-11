package engine

import (
	"api-flow/engine/core"
	"time"

	"github.com/jinzhu/gorm"
)

// WorkflowInstance 流程实例
type WorkflowInstance struct {
	core.BasicModel
	WorkflowID   uint               `json:"workflowId"`
	WorkflowName string             `json:"workflowName"`
	Status       core.ExecuteStatus `json:"status"`
	StartTime    time.Time          `json:"startTime"`
	EndTime      time.Time          `json:"endTime"`
	Inputs       string             `json:"inputs" gorm:"type:text"`  // JSON字符串
	Results      string             `json:"results" gorm:"type:text"` // JSON字符串
	ErrorMessage string             `json:"errorMessage"`
	Duration     int64              `json:"duration" gorm:"comment:'执行时间(ms)'"`
}

func MigrateWorkflowInstance(db *gorm.DB) error {
	// 自动迁移数据库结构
	return db.AutoMigrate(&WorkflowInstance{}).Error
}
