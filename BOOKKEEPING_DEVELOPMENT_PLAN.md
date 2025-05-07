# 记账项目开发进展与结构说明

## 一、项目整体结构

本项目为基于 Go 语言的记账系统后端，采用 Gin 框架，分层清晰，便于扩展和维护。主要目录结构如下：

- `api/`         —— 路由处理层，负责 HTTP 请求的接收与响应（如分类、账户、流水、用户相关接口）
- `service/`     —— 业务逻辑层，处理具体的业务操作（如增删改查、数据校验等）
- `model/`       —— 数据模型层，定义数据库表结构及相关类型
- `router/`      —— 路由注册与分组
- `utils/`       —— 工具函数（如 JWT、加密、校验等）
- `config/`      —— 配置相关（如数据库、Redis、JWT、日志等）
- `core/`        —— 框架核心初始化（如 GORM、Viper、Logger 等）
- `global/`      —— 全局变量与对象
- `middleware/`  —— 中间件（如 CORS、鉴权等）
- `docs/`        —— Swagger 文档
- `main.go`      —— 程序入口

## 二、已实现主要功能模块

### 1. 分类管理（Category）
- 主要文件：
  - `api/bookkeeping_category_handler.go`
  - `service/bookkeeping_category_service.go`
  - `model/bookkeeping_category.go`
  - `service/dto/bookkeeping_category_dto.go`
- 功能说明：
  - 分类的增、删、改、查（支持层级结构、扁平结构、类型过滤、父子分类等）
  - 接口文档已通过 Swagger 注释完善

### 2. 账户管理（Account）
- 主要文件：
  - `api/bookkeeping_account_handler.go`
  - `service/bookkeeping_account_service.go`
  - `model/bookkeeping_account.go`
  - `service/dto/bookkeeping_account_dto.go`
- 功能说明：
  - 账户的增、删、改、查
  - 支持多账户管理

### 3. 交易流水管理（Transaction）
- 主要文件：
  - `api/bookkeeping_transaction_handler.go`
  - `service/bookkeeping_transaction_service.go`
  - `model/bookkeeping_transaction.go`
  - `service/dto/bookkeeping_transaction_dto.go`
- 功能说明：
  - 交易流水的增、删、改、查
  - 支持按账户、分类、时间等条件筛选

### 4. 用户与鉴权
- 主要文件：
  - `api/user.go`
  - `service/user_info.go`
  - `model/user_info.go`
  - `service/dto/user_dto.go`
  - `utils/jwt.go`
- 功能说明：
  - 用户注册、登录、信息获取
  - JWT 鉴权中间件

## 三、开发建议与后续计划

1. **完善测试**：建议在 `test/` 目录下补充各模块的单元测试和集成测试。
2. **前端对接**：可基于当前接口文档，开发前端页面（如 Vue/React），实现分类、账户、流水的管理与展示。
3. **优化接口**：根据实际业务需求，细化接口权限、增加批量操作、导入导出等功能。
4. **部署与运维**：完善配置文件，支持多环境部署，增加日志、监控等。
5. **文档维护**：保持接口文档与代码同步，便于团队协作和后续维护。

## 四、推荐开发流程

1. 明确需求，梳理数据结构和接口设计。
2. 先实现 model、dto、service 层，确保业务逻辑正确。
3. 编写 api 层，完善接口文档（Swagger 注释）。
4. 编写/完善测试用例，保证功能稳定。
5. 前后端联调，优化接口和交互体验。
6. 持续集成与部署。

---

如需继续开发，建议先阅读本文件和 `README.md`，结合 `api/`、`service/`、`model/` 目录下的代码，快速了解当前进度和待办事项。