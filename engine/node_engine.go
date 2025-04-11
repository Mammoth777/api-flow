package engine

import (
	"api-flow/engine/core"
	"errors"
	"fmt"
)

// NodeExecutor 节点执行器接口
type NodeExecutor interface {
	Execute(node *Node, inputs map[string]interface{}) *core.ExecuteResult
	ValidateConfig(config ItemConfig) error
	GetOutputFormat() core.ParamFormat
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
	engine.RegisterExecutor(ApiNodeType.Code,  NewAPINodeExecutor())
	engine.RegisterExecutor(TextNodeType.Code, NewTextNodeExecutor())

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
func (e *NodeEngine) ExecuteNode(node *Node, inputs map[string]interface{}) (*core.ExecuteResult, error) {
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
	return executor.Execute(node, inputs), nil
}
