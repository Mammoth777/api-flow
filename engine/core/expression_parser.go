package core

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func getValueFromMap(m interface{}, properties []string) (interface{}, error) {
	if len(properties) == 0 {
		return m, nil
	}
	switch m := m.(type) {
		case map[string]interface{}:
			val, ok := m[properties[0]]
			if !ok {
				return nil, errors.New("未找到属性值")
			}
			return getValueFromMap(val, properties[1:])
		default:
			if len(properties) > 0 {
				return nil, errors.New("不支持的值类型")
			}
			return m, nil
	}
}

// ExpressionParser 表达式解析器
// 用于解析形如 ${nodeId.property} 的表达式
type ExpressionParser struct {
	results []ExecuteResult
}

// NewExpressionParser 创建表达式解析器
func NewExpressionParser(results []ExecuteResult) *ExpressionParser {
	return &ExpressionParser{
		results: results,
	}
}

func (p *ExpressionParser) Parse(expression string) (interface{}, error) {
	reg := regexp.MustCompile(`\${(.+?)}`)
	matches := reg.FindStringSubmatch(expression)
	if len(matches) == 0 {
		return expression, nil
	} else {
		// 解析表达式
		expr := matches[1]
		value, err := p.Evaluate(expr)
		if err != nil {
			return "", err
		}
		return value, nil
	}
}

// Evaluate 评估表达式并返回实际值
func (p *ExpressionParser) Evaluate(expression string) (interface{}, error) {
	parts := strings.Split(expression, ".")
	if len(parts) < 2 {
		return "", errors.New("表达式格式无效，应为 nodeKey.property")
	}

	nodeKey := parts[0]
	properties := parts[1:]

	// 查找节点执行结果
	for _, result := range p.results {
		if result.NodeKey == nodeKey {
			return p.extractProperty(result.Data, properties)
		}
	}

	return nil, errors.New("未找到节点ID对应的执行结果")
}

// extractProperty 从执行结果中提取属性值
func (p *ExpressionParser) extractProperty(output interface{}, properties []string) (interface{}, error) {
	// 处理嵌套属性，如 data.user.name
	return getValueFromMap(output, properties)
}

// Sprint 将任意类型转换为字符串
func Sprint(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}
