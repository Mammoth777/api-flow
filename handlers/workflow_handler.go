package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"api-flow/dto"
	"api-flow/services"
)

// WorkflowHandler 处理流程相关API
type WorkflowHandler struct {
	workflowService *services.WorkflowService
}

// NewWorkflowHandler 创建流程处理器实例
func NewWorkflowHandler(workflowService *services.WorkflowService) *WorkflowHandler {
	return &WorkflowHandler{
		workflowService: workflowService,
	}
}

// Get 获取单个流程
func (h *WorkflowHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID参数"})
		return
	}

	workflow, err := h.workflowService.GetWorkflowWithNodes(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workflow)
}

// List 获取流程列表
func (h *WorkflowHandler) List(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)

	workflows, count, err := h.workflowService.GetAllWorkflows(page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": count,
		"data":  workflows,
	})
}

// Update 更新流程
func (h *WorkflowHandler) update(c *gin.Context, workflow dto.WorkflowDTO) {
	id := workflow.ID
	if err := h.workflowService.UpdateWorkflow(uint(id), &workflow); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "流程更新成功",
		"id":      id,
	})
}

// Delete 删除流程
func (h *WorkflowHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID参数"})
		return
	}

	if err := h.workflowService.DeleteWorkflow(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "流程删除成功"})
}

// Save 保存流程
func (h *WorkflowHandler) Save(c *gin.Context) {
	var workflowDTO dto.WorkflowDTO
	if err := c.ShouldBindJSON(&workflowDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if workflowDTO.ID != 0 {
		h.update(c, workflowDTO)
	} else {
		h.createWithNodes(c, workflowDTO)
	}
}

// CreateWithNodes 创建工作流及其节点
func (h *WorkflowHandler) createWithNodes(c *gin.Context, workflowDTO dto.WorkflowDTO) {
	response, err := h.workflowService.SaveWorkflow(&workflowDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "工作流创建成功",
		"id":      response.ID,
	})
}

// ExecuteWorkflow 执行工作流
func (h *WorkflowHandler) ExecuteWorkflow(c *gin.Context) {
	log.Println("执行工作流")
	var request dto.WorkflowExecutionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var result *dto.WorkflowExecutionResult
	var err error
	if request.Sync {
		// 执行同步工作流
		result, err = h.workflowService.ExecuteWorkflow(&request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// 执行异步工作流
		log.Fatalln("异步执行工作流的逻辑尚未实现")
		c.JSON(http.StatusNotImplemented, gin.H{"error": "异步执行工作流的逻辑尚未实现"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *WorkflowHandler) PublishWorkflow(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID参数"})
		return
	}
	if err := h.workflowService.PublishWorkflow(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "工作流发布成功",
	})
}
