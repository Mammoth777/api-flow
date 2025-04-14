package core

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

type BasicModel struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
}



type BasicModelWithUUID struct {
	ID        string     `gorm:"primary_key;size=36" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
}

// BeforeCreate is a GORM hook that generates a UUID if the ID is not set
func (b *BasicModelWithUUID) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == "" {
		b.ID = time.Now().Format("20060102150405") + "-" + strconv.FormatInt(b.CreatedAt.UnixNano(), 10)
	}
	fmt.Println("BeforeCreate ID:", b.ID)
	return
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
	NodeID  string                 `json:"nodeId"`
	NodeKey string               `json:"nodeKey"`
	Status  ExecuteStatus `json:"status"`
	Data    ExecuteOutput        `json:"data,omitempty"`
	Error   string               `json:"error,omitempty"`
}