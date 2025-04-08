package dto

import (
	"api-flow/engine"
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
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   string        `json:"updated_at"`
	Nodes       []models.Node `json:"nodes"`
	Edges       []models.Edge `json:"edges"`
}

// WorkflowExecutionRequest 工作流执行请求
type WorkflowExecutionRequest struct {
	WorkflowID uint `json:"workflowId" binding:"required"`
	Sync       bool
	Inputs     map[string]interface{} `json:"inputs"`
}

// WorkflowExecutionResult 工作流执行结果
type WorkflowExecutionResult struct {
	WorkflowID    uint                   `json:"workflow_id"`
	WorkflowName  string                 `json:"workflow_name"`
	Success       bool                   `json:"success"`
	NodeResults   []engine.ExecuteResult `json:"node_results"`
	ErrorMessage  string                 `json:"error_message,omitempty"`
	ExecutionTime string                 `json:"execution_time"`
}
