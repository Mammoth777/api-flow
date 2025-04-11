package engine

import (
	"api-flow/engine/core"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/jinzhu/gorm"
)

// ItemConfig 节点/连线配置JSON存储结构
type ItemConfig map[string]interface{}

// Value 实现driver.Valuer接口
func (c ItemConfig) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}
	bytes, err := json.Marshal(c)
	return string(bytes), err
}

// Scan 实现sql.Scanner接口
func (c *ItemConfig) Scan(value interface{}) error {
	if value == nil {
		*c = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("不支持的类型")
	}

	return json.Unmarshal(bytes, c)
}

// Node 节点模型
type Node struct {
	core.BasicModel
	NodeKey     string     `gorm:"size:255;not null" json:"nodeKey"`
	NodeType    string     `json:"nodeType"`
	Name        string     `gorm:"size:255;not null" json:"name"`
	Description string     `gorm:"size:1000" json:"description"`
	Config      ItemConfig `gorm:"type:json" json:"config"`
	Status      string     `gorm:"size:50;default:'active'" json:"status"`
	WorkflowID  uint       `json:"workflowId"`
}

// TableName 指定表名
func (Node) TableName() string {
	return "nodes"
}

// MigrateNode 创建节点表
func MigrateNode(db *gorm.DB) error {
	return db.AutoMigrate(&Node{}).Error
}
