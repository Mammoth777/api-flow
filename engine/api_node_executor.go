package engine

import (
	"api-flow/engine/core"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"
	"time"
)

var apiNodeOutputFormat = core.ParamFormat{
	core.NewParamDefination("response", core.DataTypeObject, "API响应数据"),
}
var apiNodeInputFormat = core.ParamFormat{
	core.NewParamString("url", "API 请求URL", ""),
	core.NewParamOptions("method", "API 请求方法", "GET", []interface{}{"GET", "POST", "PUT", "DELETE", "PATCH"}),
	core.NewParamObject("headers", "请求头", map[string]interface{}{}),
	core.NewParamString("body", "请求体", ""),
	core.NewParamObject("params", "请求参数", map[string]interface{}{}),
	core.NewParamNumber("timeout", "请求超时时间（秒）", 30),
	core.NewParamNumber("retry", "重试次数", 0),
	core.NewParamNumber("retryInterval", "重试间隔时间（秒）", 5),
}
var ApiNodeType = &NodeType{
	Code:        "api",
	Name:        "API节点",
	Description: "发送HTTP请求并处理响应的节点",
	Category:    "Task",
	Input: apiNodeInputFormat,
	Output: apiNodeOutputFormat,
}

// APINodeExecutor API节点执行器
type APINodeExecutor struct {
	client *http.Client
}

// NewAPINodeExecutor 创建API节点执行器实例
func NewAPINodeExecutor() *APINodeExecutor {
	return &APINodeExecutor{
		client: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

func (e *APINodeExecutor) GetOutputFormat() core.ParamFormat {
	return apiNodeOutputFormat
}

// ValidateConfig 验证API节点配置
func (e *APINodeExecutor) ValidateConfig(config ItemConfig) error {
	if config == nil {
		return errors.New("配置不能为空")
	}

	url, ok := config["url"].(string)
	if !ok || url == "" {
		return errors.New("url必须是非空字符串")
	}

	method, ok := config["method"].(string)
	if !ok || method == "" {
		return errors.New("method必须是非空字符串")
	}

	// 验证HTTP方法是否支持
	method = strings.ToUpper(method)
	if method != "GET" && method != "POST" && method != "PUT" && method != "DELETE" && method != "PATCH" {
		return errors.New("不支持的HTTP方法")
	}

	return nil
}

func (e *APINodeExecutor) newFailExecuteResult(msg string) *core.ExecuteResult {
	return &core.ExecuteResult{
		Status: core.ExecuteStatusError,
		Data: nil,
		Error: msg,
	}
}

// Execute 执行API请求
func (e *APINodeExecutor) Execute(node *Node, inputs map[string]interface{}) *core.ExecuteResult {
	config := node.Config

	// 获取URL和方法
	urlTpl, _ := config["url"].(string)
	method, _ := config["method"].(string)
	method = strings.ToUpper(method)

	// 处理URL模板
	url, err := renderTemplate(urlTpl, inputs)
	if err != nil {
		return e.newFailExecuteResult(fmt.Sprintf("渲染URL模板失败: %v", err))
	}

	// 准备请求体
	var reqBody []byte
	if body, ok := config["body"]; ok && body != nil {
		if bodyStr, ok := body.(string); ok {
			// 渲染请求体模板
			renderedBody, err := renderTemplate(bodyStr, inputs)
			if err != nil {
				return e.newFailExecuteResult(fmt.Sprintf("渲染请求体模板失败: %v", err))
			}
			reqBody = []byte(renderedBody)
		} else {
			// 非字符串类型的请求体，直接序列化为JSON
			reqBody, err = json.Marshal(body)
			if err != nil {
				return e.newFailExecuteResult(fmt.Sprintf("序列化请求体失败: %v", err))
			}
		}
	}

	// 创建HTTP请求
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return e.newFailExecuteResult(fmt.Sprintf("创建HTTP请求失败: %v", err))
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	if headers, ok := config["headers"].(map[string]interface{}); ok {
		for key, value := range headers {
			if strValue, ok := value.(string); ok {
				req.Header.Set(key, strValue)
			}
		}
	}

	// 执行请求
	resp, err := e.client.Do(req)
	if err != nil {
		return e.newFailExecuteResult(fmt.Sprintf("执行HTTP请求失败: %v", err))
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return e.newFailExecuteResult(fmt.Sprintf("读取响应失败: %v", err))
	}

	// 解析JSON响应
	var responseData interface{}
	if err := json.Unmarshal(respBody, &responseData); err != nil {
		// 如果不是JSON，直接返回字符串
		responseData = string(respBody)
	}

	var status core.ExecuteStatus
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		status = core.ExecuteStatusSuccess
	} else {
		status = core.ExecuteStatusError
	}

	// 返回执行结果
	return &core.ExecuteResult{
		NodeID: node.ID,
		NodeKey: node.NodeKey,
		Status:  status,
		Data:  map[string]interface{}{
			"response": responseData,
		},
	}
}

// renderTemplate 使用输入数据渲染模板字符串
func renderTemplate(tpl string, data map[string]interface{}) (string, error) {
	t, err := template.New("template").Parse(tpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
