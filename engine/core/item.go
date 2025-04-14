package core

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
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