package engine

import (
	"api-flow/engine/core"

	"github.com/jinzhu/gorm"
)

type Edge struct {
	core.BasicModelWithUUID
	SourceNodeKey string     `json:"sourceNodeKey"`
	TargetNodeKey string     `json:"targetNodeKey"`
	Config        ItemConfig `gorm:"type:json" json:"config"`
	WorkflowID    uint       `json:"workflowId"`
}

// tableName 指定表名
func (Edge) TableName() string {
	return "edges"
}

// MigrateEdge 创建连线表
func MigrateEdge(db *gorm.DB) error {
	return db.AutoMigrate(&Edge{}).Error
}
