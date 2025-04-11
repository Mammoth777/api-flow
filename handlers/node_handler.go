package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"api-flow/engine"
	"api-flow/services"
)

// NodeHandler 处理节点相关API
type NodeHandler struct {
	nodeService        *services.NodeService
	nodeExecutionService *services.NodeExecutionService
}

// NewNodeHandler 创建节点处理器实例
func NewNodeHandler(nodeService *services.NodeService, nodeExecutionService *services.NodeExecutionService) *NodeHandler {
	return &NodeHandler{
		nodeService:        nodeService,
		nodeExecutionService: nodeExecutionService,
	}
}

// Get 获取单个节点
func (h *NodeHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID参数"})
		return
	}

	node, err := h.nodeService.GetNodeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, node)
}

// Update 更新节点
func (h *NodeHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID参数"})
		return
	}

	var node engine.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.nodeService.UpdateNode(uint(id), &node); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "节点更新成功"})
}

// Delete 删除节点
func (h *NodeHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID参数"})
		return
	}

	if err := h.nodeService.DeleteNode(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "节点删除成功"})
}

// Execute 执行节点
func (h *NodeHandler) Execute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID参数"})
		return
	}

	// 获取输入参数
	var inputs map[string]interface{}
	if err := c.ShouldBindJSON(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的输入参数"})
		return
	}

	// 执行流程
	result, err := h.nodeExecutionService.ExecuteNodeWithoutWorkflow(uint(id), inputs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetNodeTypes 获取所有节点类型
func (h *NodeHandler) GetNodeTypes(c *gin.Context) {
	nodeTypes, err := h.nodeService.GetAllNodeTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nodeTypes)
}
