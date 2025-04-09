package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	"api-flow/database"
	"api-flow/dto"
	"api-flow/engine"
	"api-flow/models"
)

// WorkflowService 流程服务
type WorkflowService struct {
	DB                   *gorm.DB
	NodeExecutionService *NodeExecutionService
}

// NewWorkflowService 创建流程服务实例
func NewWorkflowService() *WorkflowService {
	nodeService := NewNodeService()
	return &WorkflowService{
		DB:                   database.DB,
		NodeExecutionService: NewNodeExecutionService(nodeService),
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
func (s *WorkflowService) UpdateWorkflow(id uint, workflow *dto.WorkflowDTO) error {
	// 先检查记录是否存在且未被删除
	existingWorkflow := &models.Workflow{}
	if err := s.DB.Where("id = ? AND deleted_at IS NULL", id).First(existingWorkflow).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return errors.New("流程不存在")
		}
		return err
	}

	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新流程的基本信息
	updateWorkflowBasic := tx.Model(existingWorkflow).Updates(map[string]interface{}{
		"name":        workflow.Name,
		"description": workflow.Description,
		"status":      workflow.Status,
	})
	if updateWorkflowBasic.Error != nil {
		tx.Rollback()
		return updateWorkflowBasic.Error
	}

	// 更新节点
	for _, node := range workflow.Nodes {
		updateNodes := tx.Model(&node).Where("id = ?", node.ID).Update(map[string]interface{}{
			"nodeKey": node.NodeKey,
			"nodeType": node.NodeType,
			"name": node.Name,
			"description": node.Description,
			"config": node.Config,
			"status": node.Status,
		})
		if updateNodes.Error != nil {
			tx.Rollback()
			return updateNodes.Error
		}
	}

	// 更新连线
	for _, edge := range workflow.Edges {
		updateEdges := tx.Model(&edge).Where("id = ?", edge.ID).Update(map[string]interface{}{
			"sourceNodeKey": edge.SourceNodeKey,
			"targetNodeKey": edge.TargetNodeKey,
			"workflowID":    edge.WorkflowID,
		})
		if updateEdges.Error != nil {
			tx.Rollback()
			return updateEdges.Error
		}
	}
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
	
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

// SaveWorkflow 保存工作流及其节点和连线
func (s *WorkflowService) SaveWorkflow(workflowDto *dto.WorkflowDTO) (*dto.WorkflowDTO, error) {
	if err := s.flowValidate(workflowDto); err != nil {
		return nil, err
	}
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
		workflowDto.Nodes[i].WorkflowID = workflow.ID

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

	// 保存所有连线
	for i := range workflowDto.Edges {
		workflowDto.Edges[i].WorkflowID = workflow.ID
		if err := tx.Create(&workflowDto.Edges[i]).Error; err != nil {
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
		ID:          workflow.ID,
		Name:        workflow.Name,
		Description: workflow.Description,
		CreatedAt:   workflow.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   workflow.UpdatedAt.Format("2006-01-02 15:04:05"),
		Nodes:       workflowDto.Nodes,
		Edges:       workflowDto.Edges,
	}

	return response, nil
}

func (s *WorkflowService) flowValidate(workflowDto *dto.WorkflowDTO) error {
	// 检查工作流名称是否为空
	if workflowDto.Name == "" {
		return errors.New("工作流名称不能为空")
	}

	// 检查节点是否存在
	for _, node := range workflowDto.Nodes {
		if node.NodeKey == "" {
			return errors.New("节点键不能为空")
		}
	}

	// 检查连线是否存在
	for _, edge := range workflowDto.Edges {
		if edge.SourceNodeKey == "" || edge.TargetNodeKey == "" {
			return errors.New("连线的源节点和目标节点不能为空")
		}
	}

	// 检查循环依赖
	if err := s.checkCircularDependency(workflowDto.Edges); err != nil {
		return errors.New("工作流存在循环依赖")
	}

	return nil
}

func (s *WorkflowService) checkCircularDependency(edges []models.Edge) error {
	// 构建邻接表
	graph := make(map[string][]string)
	for _, edge := range edges {
		graph[edge.SourceNodeKey] = append(graph[edge.SourceNodeKey], edge.TargetNodeKey)
	}
	// 深度优先搜索
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	var dfs func(node string) bool
	dfs = func(node string) bool {
		if !visited[node] {
			visited[node] = true
			recStack[node] = true

			for _, neighbor := range graph[node] {
				if !visited[neighbor] && dfs(neighbor) {
					return true
				} else if recStack[neighbor] {
					return true
				}
			}
		}
		recStack[node] = false
		return false
	}
	for node := range graph {
		if !visited[node] {
			if dfs(node) {
				return errors.New("工作流存在循环依赖")
			}
		}
	}
	return nil
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
	if err := s.DB.Where("workflow_id = ?", workflowID).Find(&nodes).Error; err != nil {
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

// ExecuteWorkflow 执行工作流
func (s *WorkflowService) ExecuteWorkflow(request *dto.WorkflowExecutionRequest) (*dto.WorkflowExecutionResult, error) {
	startTime := time.Now()

	// 获取工作流信息
	workflow, err := s.GetWorkflowByID(request.WorkflowID)
	if err != nil {
		return nil, fmt.Errorf("获取工作流失败: %v", err)
	}

	// 获取工作流的所有节点
	var nodes []models.Node
	if err := s.DB.Where("workflow_id = ?", request.WorkflowID).Find(&nodes).Error; err != nil {
		return nil, fmt.Errorf("获取工作流节点失败: %v", err)
	}

	// 获取工作流的所有连线
	var edges []models.Edge
	if err := s.DB.Where("workflow_id = ?", request.WorkflowID).Find(&edges).Error; err != nil {
		return nil, fmt.Errorf("获取工作流连线失败: %v", err)
	}

	// 建立节点映射，方便快速查找
	nodeMap := make(map[string]*models.Node)
	for i := range nodes {
		nodeMap[nodes[i].NodeKey] = &nodes[i]
	}

	var success bool = true
	var errorMessage string

	// 执行节点
	nodeResults, err := s.executeNodes(nodes, edges, request.Inputs)
	if err != nil {
		success = false
		errorMessage = fmt.Sprintf("执行节点失败: %v", err)
	} else {
		for _, result := range nodeResults {
			if !result.Success {
				success = false
				errorMessage += fmt.Sprintf("节点 %s 执行失败: %s\n", result.NodeKey, result.Error)
			}
		}
	}

	// 构建工作流执行结果
	executionResult := &dto.WorkflowExecutionResult{
		WorkflowID:    workflow.ID,
		WorkflowName:  workflow.Name,
		Success:       success,
		NodeResults:   nodeResults,
		ErrorMessage:  errorMessage,
		ExecutionTime: time.Since(startTime).String(),
	}

	return executionResult, nil
}

func (s *WorkflowService) executeNodes(nodes []models.Node, edges []models.Edge, inputs map[string]interface{}) ([]engine.ExecuteResult, error) {
	// 建立节点映射，方便快速查找
	nodeMap := make(map[string]*models.Node)
	for i := range nodes {
		nodeMap[nodes[i].NodeKey] = &nodes[i]
	}
	// 获取起始节点列表
	startNodes, err := s.getStartNodeKeyList(nodes, edges)
	if err != nil {
		return nil, err
	}
	toBeExecuted := make(chan string, len(nodes))
	for _, startNodeKey := range startNodes {
		toBeExecuted <- startNodeKey
	}
	var getNextNodeKeys = func(nodeKey string) []string {
		nextNodes := make([]string, 0)
		for _, edge := range edges {
			if edge.SourceNodeKey == nodeKey {
				nextNodes = append(nextNodes, edge.TargetNodeKey)
			}
		}
		return nextNodes
	}
	results := make([]engine.ExecuteResult, 0)
nodeloop:
	for {
		select {
		case nodeKey, ok := <-toBeExecuted:
			if !ok {
				fmt.Println("读取节点异常")
				return nil, errors.New("读取节点异常")
			}
			node := nodeMap[nodeKey]
			result, err := s.NodeExecutionService.ExecuteNode(node, inputs, results)
			if err != nil {
				return nil, fmt.Errorf("节点 %s 执行失败: %v", node.NodeKey, err)
			}
			results = append(results, engine.ExecuteResult{
				NodeID:  node.ID,
				NodeKey: node.NodeKey,
				Success: result.Success,
				Data:    result.Data,
				Error:   result.Error,
			})
			nextKeys := getNextNodeKeys(nodeKey)
			for _, nextKey := range nextKeys {
				if _, ok := nodeMap[nextKey]; ok {
					toBeExecuted <- nextKey
				} else {
					fmt.Printf("节点 %s 不存在\n", nextKey)
				}
			}
		default:
			if len(toBeExecuted) == 0 {
				close(toBeExecuted)
				break nodeloop
			}
		}
	}
	return results, nil
}

// getStartNodeKeyList 获取工作流的起始节点列表
func (s *WorkflowService) getStartNodeKeyList(nodes []models.Node, edges []models.Edge) ([]string, error) {
	if len(nodes) == 0 {
		return nil, errors.New("工作流没有节点")
	}
	startNodeKeyList := make([]string, 0)
	// 使用map来存储每个节点的入度
	inDegree := make(map[string]int)
	for _, node := range nodes {
		inDegree[node.NodeKey] = 0
	}

	// 计算每个节点的入度
	for _, edge := range edges {
		inDegree[edge.TargetNodeKey]++
	}

	// 找到入度为0的节点
	for key, val := range inDegree {
		if val == 0 {
			startNodeKeyList = append(startNodeKeyList, key)
		}
	}

	return startNodeKeyList, nil
}

// PublishWorkflow 发布工作流
func (s *WorkflowService) PublishWorkflow(id uint) error {
	var workflow models.Workflow
	if err := s.DB.Where("id = ?", id).First(&workflow).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return errors.New("工作流不存在")
		}
		return err
	}
	workflow.Status = models.WorkflowPublished
	s.DB.Save(&workflow)
	return nil
}
