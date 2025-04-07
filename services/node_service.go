package services

import (
	"errors"

	"github.com/jinzhu/gorm"

	"api-flow/database"
	"api-flow/models"
)

// NodeService 节点服务
type NodeService struct {
	DB *gorm.DB
}

// NewNodeService 创建节点服务实例
func NewNodeService() *NodeService {
	return &NodeService{
		DB: database.DB,
	}
}

// GetNodeByID 通过ID获取节点
func (s *NodeService) GetNodeByID(id uint) (*models.Node, error) {
	var node models.Node
	if err := s.DB.First(&node, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New("节点不存在")
		}
		return nil, err
	}
	return &node, nil
}

// UpdateNode 更新节点
func (s *NodeService) UpdateNode(id uint, node *models.Node) error {
	existingNode, err := s.GetNodeByID(id)
	if err != nil {
		return err
	}

	// 只更新允许的字段
	updates := map[string]interface{}{
		"name":        node.Name,
		"description": node.Description,
		"config":      node.Config,
		"status":      node.Status,
	}

	// 如果节点类型有更改，验证新节点类型是否存在
	if node.NodeType != "" && node.NodeType != existingNode.NodeType {
		var nodeType models.NodeType
		if err := s.DB.Where("code = ?", node.NodeType).First(&nodeType).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				return errors.New("节点类型不存在")
			}
			return err
		}
		updates["node_type"] = node.NodeType
	}

	// 如果工作流有更改，验证新工作流是否存在
	if node.WorkflowID != 0 && node.WorkflowID != existingNode.WorkflowID {
		var workflow models.Workflow
		if err := s.DB.First(&workflow, node.WorkflowID).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				return errors.New("工作流不存在")
			}
			return err
		}
		updates["workflow_id"] = node.WorkflowID
	}

	return s.DB.Model(existingNode).Updates(updates).Error
}

// DeleteNode 删除节点
func (s *NodeService) DeleteNode(id uint) error {
	// 先检查节点是否存在
	if _, err := s.GetNodeByID(id); err != nil {
		return err
	}

	return s.DB.Delete(&models.Node{}, id).Error
}

// GetNodesByWorkflowID 获取特定工作流的所有节点
func (s *NodeService) GetNodesByWorkflowID(workflowID uint) ([]models.Node, error) {
	var nodes []models.Node
	if err := s.DB.Where("workflow_id = ?", workflowID).Preload("NodeType").Find(&nodes).Error; err != nil {
		return nil, err
	}
	return nodes, nil
}

// GetNodeTypeByCode 通过代码获取节点类型
func (s *NodeService) GetNodeTypeByCode(code string) (*models.NodeType, error) {
	var nodeType models.NodeType
	if err := s.DB.Where("code = ?", code).First(&nodeType).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New("节点类型不存在")
		}
		return nil, err
	}
	return &nodeType, nil
}

// GetAllNodeTypes 获取所有节点类型
func (s *NodeService) GetAllNodeTypes() ([]models.NodeType, error) {
	var nodeTypes []models.NodeType
	if err := s.DB.Find(&nodeTypes).Error; err != nil {
		return nil, err
	}
	return nodeTypes, nil
}
