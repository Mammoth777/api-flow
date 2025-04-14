package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Record map[string]interface{}

// 实现 sql.Valuer 接口
func (r Record) Value() (driver.Value, error) {
	if r == nil {
		return nil, nil
	}
	bytes, err := json.Marshal(r)
	return string(bytes), err
}

// 实现 sql.Scanner 接口
func (r *Record) Scan(value interface{}) error {
	if value == nil {
		*r = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("unsupported type")
	}

	return json.Unmarshal(bytes, r)
}