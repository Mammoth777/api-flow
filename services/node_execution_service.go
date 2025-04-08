package services

import (
	"api-flow/engine"
	"api-flow/models"
	"fmt"
)

// NodeExecutionService 节点执行服务
type NodeExecutionService struct {
	nodeService *NodeService
	nodeEngine  *engine.NodeEngine
}

// NewNodeExecutionService 创建节点执行服务实例
func NewNodeExecutionService(nodeService *NodeService) *NodeExecutionService {
	return &NodeExecutionService{
		nodeService: nodeService,
		nodeEngine:  engine.NewNodeEngine(),
	}
}

func (s *NodeExecutionService) getRealValue(value interface{}, results []engine.ExecuteResult) (interface{}, error) {
	if expression, ok := value.(string); ok {
		parser := engine.NewExpressionParser(results)
		return parser.Parse(expression)
	} else {
		return value, nil
	}
}

func (s *NodeExecutionService) ExecuteNode(node *models.Node, inputs map[string]interface{}, results []engine.ExecuteResult) (*engine.ExecuteResult, error) {
	config := node.Config
	if config != nil {
		// 覆盖默认输入
		configInputs, ok := config["inputs"].(map[string]interface{})
		if ok {
			for key, expressionOrValue := range configInputs {
				fmt.Println("expressionOrValue:", expressionOrValue)
				realVal, err := s.getRealValue(expressionOrValue, results)
				if err != nil {
					fmt.Println("Error parsing expression:", err)
				} else {
					inputs[key] = realVal
				}
			}
		}
	}
	return s.nodeEngine.ExecuteNode(node, inputs)
}

// ExecuteNode 执行节点
func (s *NodeExecutionService) ExecuteNodeById(nodeID uint, inputs map[string]interface{}, results []engine.ExecuteResult) (*engine.ExecuteResult, error) {
	// 检查节点状态
	// if node.Status != "active" {
	// 	return nil, errors.New("节点未激活，无法执行")
	// }

	// 执行节点
	return s.ExecuteNodeWithoutWorkflow(nodeID, inputs)
}

func (s *NodeExecutionService) ExecuteNodeWithoutWorkflow(nodeID uint, inputs map[string]interface{}) (*engine.ExecuteResult, error) {
	// 获取节点信息
	node, err := s.nodeService.GetNodeByID(nodeID)
	if err != nil {
		return nil, err
	}

	// 执行节点
	return s.nodeEngine.ExecuteNode(node, inputs)
}