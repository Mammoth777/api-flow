package engine

import "api-flow/engine/core"

var textNodeOutputFormat = core.ParamFormat{
	core.NewParamDefination("output", core.DataTypeString, "echo + input text content"),
}

var textNodeInputFormat = core.ParamFormat{
	core.NewParamDefination("content", core.DataTypeString, "请输入文本内容"),
}

var TextNodeType = &NodeType{
	Code:        "text",
	Name:        "文本节点",
	Description: "直接返回配置的文本内容的节点",
	Category:    "Task",
	Input: textNodeInputFormat,
	Output: textNodeOutputFormat,
}

// TextNodeExecutor 文本节点执行器
type TextNodeExecutor struct{}

// NewTextNodeExecutor 创建文本节点执行器实例
func NewTextNodeExecutor() *TextNodeExecutor {
	return &TextNodeExecutor{}
}

func (e *TextNodeExecutor) GetOutputFormat() core.ParamFormat {
	return textNodeOutputFormat
}

// ValidateConfig 验证文本节点配置
func (e *TextNodeExecutor) ValidateConfig(config ItemConfig) error {
	// if config == nil {
	// 	return errors.New("配置不能为空")
	// }

	// _, ok := config["content"]
	// if !ok {
	// 	return errors.New("content字段是必需的")
	// }

	return nil
}

func (e *TextNodeExecutor) newfailExecuteResult(node *Node, msg string) *core.ExecuteResult {
	return &core.ExecuteResult{
		NodeID: node.ID,
		NodeKey: node.NodeKey,
		Status: core.ExecuteStatusError,
		Data:    nil,
		Error:   msg,
	}
}

// Execute 执行文本节点逻辑
func (e *TextNodeExecutor) Execute(node *Node, inputs map[string]interface{}) *core.ExecuteResult {
	config := node.Config

	// 获取文本内容
	content, ok := inputs["content"]
	if !ok {
		content, ok = config["content"]
		if !ok {
			return e.newfailExecuteResult(node, "未找到content字段")
		}
	}

	// 将content转换为字符串
	contentStr, ok := content.(string)
	if !ok {
		return e.newfailExecuteResult(node, "content字段必须是字符串类型")
	}

	// 获取内容类型
	contentType, _ := config["content_type"].(string)
	if contentType == "" {
		contentType = "plain_text" // 默认为普通文本
	}

	// 返回执行结果
	return &core.ExecuteResult{
		NodeID:  node.ID,
		NodeKey: node.NodeKey,
		Status: core.ExecuteStatusSuccess,
		Data: core.ExecuteOutput{
			"output":      "echo: " + contentStr,
		},
	}
}
