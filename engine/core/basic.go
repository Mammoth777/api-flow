package core

import "time"

type BasicModel struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
}

type ExecuteStatus uint8

const (
	ExecuteStatusReady ExecuteStatus = iota
	ExecuteStatusRunning
	ExecuteStatusSuccess
	ExecuteStatusError
)

// ExecuteResult 节点执行结果
type ExecuteResult struct {
	NodeID  uint                 `json:"nodeId"`
	NodeKey string               `json:"nodeKey"`
	Status  ExecuteStatus `json:"status"`
	Data    ExecuteOutput        `json:"data,omitempty"`
	Error   string               `json:"error,omitempty"`
}