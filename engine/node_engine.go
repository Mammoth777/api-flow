package engine

import (
	"errors"
	"fmt"

	"api-flow/models"
)

// ExecuteResult 节点执行结果
type ExecuteResult struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
	StatusCode int         `json:"status_code,omitempty"`
}

// NodeExecutor 节点执行器接口
type NodeExecutor interface {
	Execute(node *models.Node, inputs map[string]interface{}) (*ExecuteResult, error)
	ValidateConfig(config models.NodeConfig) error
}

// NodeEngine 节点执行引擎
type NodeEngine struct {
	executors map[string]NodeExecutor
}

// NewNodeEngine 创建节点执行引擎实例
func NewNodeEngine() *NodeEngine {
	engine := &NodeEngine{
		executors: make(map[string]NodeExecutor),
	}
	
	// 注册默认执行器
	engine.RegisterExecutor(models.NodeTypeAPI, NewAPINodeExecutor())
	engine.RegisterExecutor(models.NodeTypeText, NewTextNodeExecutor())
	
	return engine
}

// RegisterExecutor 注册节点执行器
func (e *NodeEngine) RegisterExecutor(nodeTypeCode string, executor NodeExecutor) {
	e.executors[nodeTypeCode] = executor
}

// GetExecutor 获取节点执行器
func (e *NodeEngine) GetExecutor(nodeTypeCode string) (NodeExecutor, error) {
	executor, exists := e.executors[nodeTypeCode]
	if !exists {
		return nil, fmt.Errorf("未找到节点类型 '%s' 的执行器", nodeTypeCode)
	}
	return executor, nil
}

// ExecuteNode 执行节点
func (e *NodeEngine) ExecuteNode(node *models.Node, inputs map[string]interface{}) (*ExecuteResult, error) {
	if node == nil {
		return nil, errors.New("节点不能为空")
	}
	
	if node.NodeType == "" {
		return nil, errors.New("节点类型未设置")
	}
	
	executor, err := e.GetExecutor(node.NodeType)
	if err != nil {
		return nil, err
	}
	
	// 验证节点配置
	if err := executor.ValidateConfig(node.Config); err != nil {
		return nil, fmt.Errorf("节点配置无效: %v", err)
	}
	
	// 执行节点
	return executor.Execute(node, inputs)
}
