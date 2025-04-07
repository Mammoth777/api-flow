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

// NodeExecutionResult 节点执行结果
type NodeExecutionResult struct {
	NodeID	 uint                  `json:"nodeId"`
	NodeKey  string                `json:"nodeKey"`
	Result   *engine.ExecuteResult `json:"result"`
}

// WorkflowExecutionResult 工作流执行结果
type WorkflowExecutionResult struct {
	WorkflowID    uint                  `json:"workflow_id"`
	WorkflowName  string                `json:"workflow_name"`
	Success       bool                  `json:"success"`
	NodeResults   []NodeExecutionResult `json:"node_results"`
	ErrorMessage  string                `json:"error_message,omitempty"`
	ExecutionTime string                `json:"execution_time"`
}
