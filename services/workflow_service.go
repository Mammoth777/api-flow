package services

import (
	"errors"

	"github.com/jinzhu/gorm"

	"api-flow/database"
	"api-flow/dto"
	"api-flow/models"
)

// WorkflowService 流程服务
type WorkflowService struct {
	DB *gorm.DB
}

// NewWorkflowService 创建流程服务实例
func NewWorkflowService() *WorkflowService {
	return &WorkflowService{
		DB: database.DB,
	}
}

// CreateWorkflow 创建新流程
func (s *WorkflowService) CreateWorkflow(workflow *models.Workflow) error {
	return s.DB.Create(workflow).Error
}

// GetWorkflowByID 通过ID获取流程
func (s *WorkflowService) GetWorkflowByID(id uint) (*models.Workflow, error) {
	var workflow models.Workflow
	// 使用Unscoped()不会自动加WHERE deleted_at IS NULL条件，我们手动处理
	if err := s.DB.Where("id = ? AND deleted_at IS NULL", id).First(&workflow).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New("流程不存在")
		}
		return nil, err
	}
	return &workflow, nil
}

// GetAllWorkflows 获取所有流程
func (s *WorkflowService) GetAllWorkflows(page, size int) ([]models.Workflow, int, error) {
	var workflows []models.Workflow
	var count int

	// 获取总记录数，只计算未删除的记录
	s.DB.Model(&models.Workflow{}).Where("deleted_at IS NULL").Count(&count)

	// 分页查询，只查询未删除的记录
	query := s.DB.Where("deleted_at IS NULL")
	if err := query.Offset((page - 1) * size).Limit(size).Find(&workflows).Error; err != nil {
		return nil, 0, err
	}

	return workflows, count, nil
}

// UpdateWorkflow 更新流程
func (s *WorkflowService) UpdateWorkflow(id uint, workflow *models.Workflow) error {
	// 先检查记录是否存在且未被删除
	existingWorkflow := &models.Workflow{}
	if err := s.DB.Where("id = ? AND deleted_at IS NULL", id).First(existingWorkflow).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return errors.New("流程不存在")
		}
		return err
	}

	// 只更新指定的字段，避免更新ID和时间戳字段
	return s.DB.Model(existingWorkflow).Updates(map[string]interface{}{
		"name":        workflow.Name,
		"description": workflow.Description,
		"status":      workflow.Status,
	}).Error
}

// DeleteWorkflow 删除流程 (软删除)
func (s *WorkflowService) DeleteWorkflow(id uint) error {
	// 检查记录是否存在
	var workflow models.Workflow
	if err := s.DB.Where("id = ? AND deleted_at IS NULL", id).First(&workflow).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return errors.New("流程不存在")
		}
		return err
	}

	// 使用Delete进行软删除，GORM会自动设置DeletedAt字段
	return s.DB.Delete(&workflow).Error
}

// SaveWorkflowWithNodes 保存工作流及其节点
func (s *WorkflowService) SaveWorkflowWithNodes(workflowDto *dto.WorkflowDTO) (*dto.WorkflowDTO, error) {
	// 开启事务
	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	workflow := models.Workflow{
		Name:        workflowDto.Name,
		Description: workflowDto.Description,
		Status:      workflowDto.Status,
	}

	// 保存工作流
	if err := tx.Create(&workflow).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 为每个节点设置工作流ID
	for i := range workflowDto.Nodes {
		workflowDto.Nodes[i].WorkflowID = workflowDto.ID

		// 验证节点类型是否存在
		var nodeType models.NodeType
		if err := tx.Where("code = ?", workflowDto.Nodes[i].NodeType).First(&nodeType).Error; err != nil {
			tx.Rollback()
			if gorm.IsRecordNotFoundError(err) {
				return nil, errors.New("节点类型不存在")
			}
			return nil, err
		}
	}

	// 保存所有节点
	for i := range workflowDto.Nodes {
		if err := tx.Create(&workflowDto.Nodes[i]).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 构建响应
	response := &dto.WorkflowDTO{
		ID:          workflowDto.ID,
		Name:        workflowDto.Name,
		Description: workflowDto.Description,
		CreatedAt:   workflow.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   workflow.UpdatedAt.Format("2006-01-02 15:04:05"),
		Nodes:       workflowDto.Nodes,
	}

	return response, nil
}

// GetWorkflowWithNodes 获取工作流及其关联节点
func (s *WorkflowService) GetWorkflowWithNodes(workflowID uint) (*dto.WorkflowDTO, error) {
	// 获取工作流
	workflow, err := s.GetWorkflowByID(workflowID)
	if err != nil {
		return nil, err
	}

	// 获取关联的节点
	var nodes []models.Node
	if err := s.DB.Where("workflow_id = ?", workflowID).Preload("NodeType").Find(&nodes).Error; err != nil {
		return nil, err
	}

	// 构建响应
	response := &dto.WorkflowDTO{
		ID:          workflow.ID,
		Name:        workflow.Name,
		Description: workflow.Description,
		CreatedAt:   workflow.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   workflow.UpdatedAt.Format("2006-01-02 15:04:05"),
		Nodes:       nodes,
	}

	return response, nil
}
