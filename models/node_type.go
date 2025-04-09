package models

import (
	"github.com/jinzhu/gorm"
)

// NodeType 节点类型模型
type NodeType struct {
	BasicModel
	Code        string `gorm:"size:50;unique;not null" json:"code"`
	Name        string `gorm:"size:100;not null" json:"name"`
	Description string `gorm:"size:500" json:"description"`
	Category    string `gorm:"size:50" json:"category"`
}

// TableName 指定表名
func (NodeType) TableName() string {
	return "node_types"
}

// 预定义节点类型代码
const (
	NodeTypeAPI  = "api"
	NodeTypeText = "text"
)

// MigrateNodeType 创建节点类型表
func MigrateNodeType(db *gorm.DB) error {
	err := db.AutoMigrate(&NodeType{}).Error
	if err != nil {
		return err
	}

	// 初始化基础节点类型
	var count int
	db.Model(&NodeType{}).Count(&count)
	if count == 0 {
		nodeTypes := []NodeType{
			{
				Code:        NodeTypeAPI,
				Name:        "API节点",
				Description: "发送HTTP请求并处理响应的节点",
				Category:    "Task",
			},
			{
				Code:        NodeTypeText,
				Name:        "文本节点",
				Description: "直接返回配置的文本内容的节点",
				Category:    "Task",
			},
		}

		for _, nt := range nodeTypes {
			db.Create(&nt)
		}
	}

	return nil
}
