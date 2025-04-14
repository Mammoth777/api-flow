/**
 * 系统节点： 执行时输入节点
 */

package engine_nodes

import "api-flow/engine/core"

// 输入格式定义为空， 前端根据节点类型， 做特殊处理
var inputNodeInputFormat = core.ParamFormat{}

// 输出格式定义, 直接输出全部输入内容
var inputNodeOutputFormat = core.ParamFormat{}

var InputNodeType = &NodeType{
	Code:        "execInput",
	Name:        "执行时输入",
	Description: "用于定义执行时入参格式的节点",
	Category:    "System",
	Input: inputNodeInputFormat,
	Output: inputNodeOutputFormat,
}

// TextNodeExecutor 文本节点执行器
type InputNodeExecutor struct{}

// NewTextNodeExecutor 创建文本节点执行器实例
func NewInputNodeExecutor() *InputNodeExecutor {
	return &InputNodeExecutor{}
}

func (e *InputNodeExecutor) GetOutputFormat() core.ParamFormat {
	return textNodeOutputFormat
}

// ValidateConfig 验证文本节点配置
func (e *InputNodeExecutor) ValidateConfig(config core.ItemConfig) error {
	return nil
}

// Execute 执行文本节点逻辑
func (e *InputNodeExecutor) Execute(node *Node, inputs map[string]interface{}) *core.ExecuteResult {

	var execParams core.ExecuteOutput

	if inputs == nil {
		execParams = make(map[string]interface{})
	} else {
		execParams = inputs
	}

	// 返回执行结果
	return &core.ExecuteResult{
		NodeID:  node.ID,
		NodeKey: node.NodeKey,
		Status: core.ExecuteStatusSuccess,
		Data: execParams,
	}
}
