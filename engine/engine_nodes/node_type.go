package engine_nodes

import (
	"api-flow/engine/core"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/jinzhu/gorm"
)

// NodeType 节点类型模型
type NodeType struct {
	core.BasicModel
	Code        string           `gorm:"size:50;unique;not null" json:"code"`
	Name        string           `gorm:"size:100;not null" json:"name"`
	Description string           `gorm:"size:500" json:"description"`
	Category    string           `gorm:"size:50" json:"category"`
	Input       core.ParamFormat `gorm:"type:text" json:"input"`
	Output      core.ParamFormat `gorm:"type:text" json:"output"`
}

// TableName 指定表名
func (NodeType) TableName() string {
	return "node_types"
}

// 实现 driver.Valuer接口
func (n *NodeType) Value() (driver.Value, error) {
	if n == nil {
		return nil, nil
	}
	bytes, err := json.Marshal(n)
	return string(bytes), err
}

// 实现 sql.Scanner接口
func (n *NodeType) Scan(value interface{}) error {
	if value == nil {
		*n = NodeType{}
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
			// 系统节点类型
			*InputNodeType,
			// 系统自带的任务节点类型
			*ApiNodeType,
			*TextNodeType,
		}

		for _, nt := range nodeTypes {
			db.Create(&nt)
		}
	}

	return nil
}
