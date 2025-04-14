package engine_nodes

import (
	"api-flow/engine/core"
	"api-flow/engine/models"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/jinzhu/gorm"
)

// Node 节点模型
type Node struct {
	core.BasicModelWithUUID
	NodeKey     string         `gorm:"size:255;not null" json:"nodeKey"`
	NodeType    string         `json:"nodeType"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Description string         `gorm:"size:1000" json:"description"`
	Config      core.ItemConfig     `gorm:"type:json" json:"config"`
	Status      string         `gorm:"size:50;default:'active'" json:"status"`
	WorkflowID  uint           `json:"workflowId"`
	Ui          models.Record `gorm:"type:json" json:"ui"`
}

// TableName 指定表名
func (Node) TableName() string {
	return "nodes"
}

// 实现 sql.Value 接口
func (n *Node) Value() (driver.Value, error) {
	if n == nil {
		return nil, nil
	}
	bytes, err := json.Marshal(n)
	return string(bytes), err
}

// 实现 sql.Scanner 接口
func (n *Node) Scan(value interface{}) error {
	if value == nil {
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

	return json.Unmarshal(bytes, n)
}

// MigrateNode 创建节点表
func MigrateNode(db *gorm.DB) error {
	return db.AutoMigrate(&Node{}).Error
}
