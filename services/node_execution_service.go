package services

import (
	"api-flow/engine"
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

// ExecuteNode 执行节点
func (s *NodeExecutionService) ExecuteNode(nodeID uint, inputs map[string]interface{}) (*engine.ExecuteResult, error) {
	// 获取节点信息
	node, err := s.nodeService.GetNodeByID(nodeID)
	if err != nil {
		return nil, err
	}

	// 检查节点状态
	// if node.Status != "active" {
	// 	return nil, errors.New("节点未激活，无法执行")
	// }

	// 执行节点
	return s.nodeEngine.ExecuteNode(node, inputs)
}
