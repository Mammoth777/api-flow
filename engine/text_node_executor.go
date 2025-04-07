package engine

import (
	"errors"

	"api-flow/models"
)

// TextNodeExecutor 文本节点执行器
type TextNodeExecutor struct{}

// NewTextNodeExecutor 创建文本节点执行器实例
func NewTextNodeExecutor() *TextNodeExecutor {
	return &TextNodeExecutor{}
}

// ValidateConfig 验证文本节点配置
func (e *TextNodeExecutor) ValidateConfig(config models.NodeConfig) error {
	// if config == nil {
	// 	return errors.New("配置不能为空")
	// }

	// _, ok := config["content"]
	// if !ok {
	// 	return errors.New("content字段是必需的")
	// }

	return nil
}

// Execute 执行文本节点逻辑
func (e *TextNodeExecutor) Execute(node *models.Node, inputs map[string]interface{}) (*ExecuteResult, error) {
	config := node.Config

	// 获取文本内容
	content, ok := inputs["content"]
	if !ok {
		content, ok = config["content"]
		if !ok {
			return nil, errors.New("未找到content字段")
		}
	}

	// 将content转换为字符串
	contentStr, ok := content.(string)
	if !ok {
		return nil, errors.New("content字段必须是字符串类型")
	}

	// 获取内容类型
	contentType, _ := config["content_type"].(string)
	if contentType == "" {
		contentType = "plain_text" // 默认为普通文本
	}

	// 返回执行结果
	return &ExecuteResult{
		Success: true,
		Data: map[string]interface{}{
			"output":      "echo: " + contentStr,
			"content_type": contentType,
		},
	}, nil
}
