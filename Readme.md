# Geminate

> Geminate is a high-performance, user-friendly tool for flexible API and SQL orchestration, enabling efficient development and seamless data management.

## 功能列表

### 流程
- [x] 流程管理：对于流程对象的的增删改查
- [x] 工作流存储：支持同时存储工作流及其多个节点
- [x] 节点执行引擎：支持节点的独立执行


### 节点

系统节点
- [ ] 输入节点: 定义流程输入参数，即调用流程时，需要传入的参数类型。
  + 后续的节点， 可以直接使用这些参数。
  + 在生成工作流调用文档时， 也可以根据这个节点的内容，定义调用参数。
- [ ] 输出节点: 定义流程输出结果，即(同步)调用流程后，返回的结果类型。
  + 从前置节点的输出，或流程的输入中， 获取值，作为输出
  + 在生成工作流调用文档时， 也可以根据这个节点的内容，定义结果类型。

普通节点
- [x] 文本节点(echo)
- [x] API节点 -- 完善中

## 项目概述

该项目是一个基于ORM的流程管理系统，支持对流程对象的增删改查操作，并通过配置连接本地MySQL数据库。
系统支持多种节点类型，包括API节点和文本节点，各节点可独立执行。

## 项目结构

```
api-flow/
├── config/           # 配置文件
├── database/         # 数据库连接
├── engine/           # 节点执行引擎
├── models/           # 数据模型
├── services/         # 业务服务
├── handlers/         # API处理器
├── router/           # 路由配置
├── main.go           # 主程序入口
└── Readme.md         # 项目文档
```

## 功能特性

### 流程管理
1. **创建流程**：POST /api/workflows
2. **查询流程**：GET /api/workflows 和 GET /api/workflows/:id
3. **更新流程**：PUT /api/workflows/:id
4. **删除流程**：DELETE /api/workflows/:id

### 节点管理
1. **创建节点**：POST /api/nodes
2. **查询节点**：GET /api/nodes 和 GET /api/nodes/:id
3. **更新节点**：PUT /api/nodes/:id
4. **删除节点**：DELETE /api/nodes/:id
5. **执行节点**：POST /api/nodes/:id/execute
6. **获取所有节点类型**：GET /api/node-types

## 节点类型

### API节点
发送HTTP请求并返回结果的节点。

配置示例：
```json
{
  "url": "https://api.example.com/data",
  "method": "GET",
  "headers": {
    "Authorization": "Bearer token",
    "Content-Type": "application/json"
  },
  "body": "{\"query\": \"{{.query}}\"}"
}
```

执行示例：
```json
{
  "query": "search term"
}
```

### 文本节点
直接返回配置的文本内容。

配置示例：
```json
{
  "content": "这是一段文本内容",
  "content_type": "plain_text"
}
```

## 运行说明

1. 确保已安装Go (1.16+)和MySQL
2. 创建MySQL数据库: `api_flow`
3. 配置 `config/config.yaml` 文件中的数据库连接信息
4. 运行项目: `go run main.go`

## API文档

### 创建流程

**请求**: `POST /api/workflows`

请求体:
```json
{
  "name": "测试流程",
  "description": "这是一个测试流程",
  "status": "active"
}
```

> 注意：创建时间和更新时间由系统自动生成，不需要在请求中提供。

### 创建节点

**请求**: `POST /api/nodes`

API节点请求体示例:
```json
{
  "node_type_id": 1,
  "name": "获取用户数据",
  "description": "从API获取用户数据",
  "workflow_id": 1,
  "status": "active",
  "config": {
    "url": "https://api.example.com/users/{{.user_id}}",
    "method": "GET",
    "headers": {
      "Authorization": "Bearer {{.token}}"
    }
  }
}
```

文本节点请求体示例:
```json
{
  "node_type_id": 2,
  "name": "欢迎信息",
  "description": "显示欢迎信息",
  "workflow_id": 1,
  "status": "active",
  "config": {
    "content": "欢迎使用我们的服务，{{.username}}！",
    "content_type": "plain_text"
  }
}
```

### 执行节点

**请求**: `POST /api/nodes/:id/execute`

请求体示例 (API节点):
```json
{
  "user_id": 123,
  "token": "your-auth-token"
}
```

请求体示例 (文本节点):
```json
{
  "username": "张三"
}
```

### 创建工作流及关联节点

**请求**: `POST /api/workflows/with-nodes`

请求体:
```json
{
  "workflow": {
    "name": "综合处理流程",
    "description": "包含多个节点的流程示例",
    "status": "active"
  },
  "nodes": [
    {
      "node_type_id": 1,
      "name": "API获取数据",
      "description": "从外部API获取数据",
      "status": "active",
      "config": {
        "url": "https://api.example.com/data",
        "method": "GET"
      }
    },
    {
      "node_type_id": 2,
      "name": "处理结果",
      "description": "显示处理结果",
      "status": "active",
      "config": {
        "content": "处理完成",
        "content_type": "plain_text"
      }
    }
  ]
}
```

### 获取工作流及关联节点

**请求**: `GET /api/workflows/:id/with-nodes`

响应示例:
```json
{
  "workflow": {
    "id": 1,
    "name": "综合处理流程",
    "description": "包含多个节点的流程示例",
    "status": "active",
    "created_at": "2023-06-01T10:00:00Z",
    "updated_at": "2023-06-01T10:00:00Z"
  },
  "nodes": [
    {
      "id": 1,
      "node_type_id": 1,
      "node_type": {
        "id": 1,
        "code": "api",
        "name": "API节点",
        "description": "发送HTTP请求并处理响应的节点"
      },
      "name": "API获取数据",
      "description": "从外部API获取数据",
      "status": "active",
      "workflow_id": 1,
      "config": {
        "url": "https://api.example.com/data",
        "method": "GET"
      },
      "created_at": "2023-06-01T10:00:00Z",
      "updated_at": "2023-06-01T10:00:00Z"
    },
    {
      "id": 2,
      "node_type_id": 2,
      "node_type": {
        "id": 2,
        "code": "text_node",
        "name": "文本节点",
        "description": "直接返回配置的文本内容的节点"
      },
      "name": "处理结果",
      "description": "显示处理结果",
      "status": "active",
      "workflow_id": 1,
      "config": {
        "content": "处理完成",
        "content_type": "plain_text"
      },
      "created_at": "2023-06-01T10:00:00Z",
      "updated_at": "2023-06-01T10:00:00Z"
    }
  ]
}
```

## 开发

```bash
docker run --name geminate-mysql -p 3306:3306 -e MYSQL_ALLOW_EMPTY_PASSWORD=yes -d mysql
docker exec -it geminate-mysql mysql -u root -e "CREATE DATABASE api_flow;"
go run main.go
```

```bash
cd fe
pnpm i
pnpm run dev
```