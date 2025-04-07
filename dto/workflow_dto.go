package dto

import (
	"api-flow/models"
)

// WorkflowDTO 工作流与节点数据传输对象
type WorkflowWithNodesDTO struct {
	Workflow models.Workflow `json:"workflow"`
	Nodes    []models.Node   `json:"nodes"`
}

// WorkflowDTO 工作流与节点响应数据传输对象
type WorkflowDTO struct {
	ID          uint          `json:"id"`
	Name        string        `json:"name"`
	Status      string        `json:"status"`
	Description string        `json:"description"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
	Nodes       []models.Node `json:"nodes"`
}
