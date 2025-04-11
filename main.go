package main

import (
	"fmt"
	"log"

	"api-flow/config"
	"api-flow/database"
	"api-flow/models"
	"api-flow/router"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库连接
	if err := database.Initialize(cfg.GetDSN()); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer database.Close()

	// 数据库迁移
	if err := models.MigrateWorkflow(database.DB); err != nil {
		log.Fatalf("工作流表迁移失败: %v", err)
	}
	
	// 迁移节点类型表并初始化基础类型
	if err := models.MigrateNodeType(database.DB); err != nil {
		log.Fatalf("节点类型表迁移失败: %v", err)
	}
	
	// 迁移节点表
	if err := models.MigrateNode(database.DB); err != nil {
		log.Fatalf("节点表迁移失败: %v", err)
	}

	// 连线表
	if err = models.MigrateEdge(database.DB); err != nil {
		log.Fatalf("连线表迁移失败: %v", err)
	}

	if err = models.MigrateWorkflowInstance(database.DB); err != nil {
		log.Fatalf("流程实例表迁移失败: %v", err)
	}

	// 设置路由
	r := router.SetupRouter()

	// 启动服务器
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("服务启动在: http://localhost%s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
