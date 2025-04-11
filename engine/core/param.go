package core

import (
	"database/sql/driver"
	"encoding/json"
)

type ParamDataType string

const (
	DataTypeString  ParamDataType = "string"
	DataTypeNumber  ParamDataType = "number"
	DataTypeBoolean ParamDataType = "boolean"
	DataTypeArray   ParamDataType = "array"
	DataTypeOptions ParamDataType = "options"
	DataTypeObject  ParamDataType = "object"
	DataTypeNull    ParamDataType = "null"
	DataTypeAny     ParamDataType = "any"
)

type ParamDefination struct {
	Field       string        `json:"field"`
	Datatype    ParamDataType `json:"type"`
	Description string        `json:"desc"`
	Default     interface{}   `json:"default"`
}

func NewParamDefination(field string, datatype ParamDataType, description string) *ParamDefination {
	return &ParamDefination{
		Field:       field,
		Datatype:    datatype,
		Description: description,
		Default:     nil,
	}
}

func NewParamDefinationWithDefault(field string, datatype ParamDataType, description string, defaultValue interface{}) *ParamDefination {
	return &ParamDefination{
		Field:       field,
		Datatype:    datatype,
		Description: description,
		Default:     defaultValue,
	}
}

// driver.Valuer 接口实现
func (e *ParamDefination) Value() (driver.Value, error) {
	if e == nil {
		return nil, nil
	}
	bytes, err := json.Marshal(e)
	return string(bytes), err
}

// sql.Scanner 接口实现
func (e *ParamDefination) Scan(value interface{}) error {
	if value == nil {
		*e = ParamDefination{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}

	return json.Unmarshal(bytes, e)
}

type ParamFormat []*ParamDefination

// 为ExecuteOutputFormat实现driver.Valuer接口
func (e ParamFormat) Value() (driver.Value, error) {
	if e == nil {
		return nil, nil
	}
	bytes, err := json.Marshal(e)
	return string(bytes), err
}

// 为ExecuteOutputFormat实现sql.Scanner接口
func (e *ParamFormat) Scan(value interface{}) error {
	if value == nil {
		*e = ParamFormat{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}

	return json.Unmarshal(bytes, e)
}

func NewParamString(field, description string, defaultValue string) *ParamDefination {
	return NewParamDefinationWithDefault(field, DataTypeString, description, defaultValue)
}

func NewParamNumber(field, description string, defaultValue float64) *ParamDefination {
	return NewParamDefinationWithDefault(field, DataTypeNumber, description, defaultValue)
}

func NewParamBoolean(field, description string, defaultValue bool) *ParamDefination {
	return NewParamDefinationWithDefault(field, DataTypeBoolean, description, defaultValue)
}
func NewParamArray(field, description string, defaultValue []interface{}) *ParamDefination {
	return NewParamDefinationWithDefault(field, DataTypeArray, description, defaultValue)
}
func NewParamOptions(field, description string, defaultValue []interface{}) *ParamDefination {
	return NewParamDefinationWithDefault(field, DataTypeOptions, description, defaultValue)
}
func NewParamObject(field, description string, defaultValue map[string]interface{}) *ParamDefination {
	return NewParamDefinationWithDefault(field, DataTypeObject, description, defaultValue)
}
func NewParamNull(field, description string) *ParamDefination {
	return NewParamDefination(field, DataTypeNull, description)
}
func NewParamAny(field, description string, defaultValue interface{}) *ParamDefination {
	return NewParamDefinationWithDefault(field, DataTypeAny, description, defaultValue)
}
