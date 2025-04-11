package dto

import (
	"api-flow/engine"
	"api-flow/engine/core"
)

// WorkflowDTO 工作流与节点数据传输对象
type WorkflowWithNodesDTO struct {
	Workflow engine.Workflow `json:"workflow"`
	Nodes    []engine.Node   `json:"nodes"`
}

// WorkflowDTO 工作流与节点响应数据传输对象
type WorkflowDTO struct {
	ID          uint                  `json:"id"`
	Name        string                `json:"name"`
	Status      engine.WorkflowStatus `json:"status"`
	Description string                `json:"description"`
	CreatedAt   string                `json:"created_at"`
	UpdatedAt   string                `json:"updated_at"`
	Nodes       []engine.Node         `json:"nodes"`
	Edges       []engine.Edge         `json:"edges"`
}

// WorkflowExecutionRequest 工作流执行请求
type WorkflowExecutionRequest struct {
	WorkflowID uint                   `json:"workflowId" binding:"required"`
	Sync       bool                   `json:"sync"`
	Inputs     map[string]interface{} `json:"inputs"`
}

// WorkflowExecutionResult 工作流执行结果
type WorkflowExecutionResult struct {
	WorkflowID   uint                   `json:"workflowId"`
	WorkflowName string                 `json:"workflowName"`
	Status       core.ExecuteStatus   `json:"status"`
	NodeResults  []core.ExecuteResult `json:"nodeResults"`
	ErrorMessage string                 `json:"errorMessage,omitempty"`
	Duration     int64                  `json:"duration"`
}
