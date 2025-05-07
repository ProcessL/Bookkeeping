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
- `cmd/`         —— 应用初始化命令
- `test/`        —— 测试文件目录
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
  - 支持收入和支出两种分类类型
  - 已实现层级分类（parent-child 关系）

### 2. 账户管理（Account）
- 主要文件：
  - `api/bookkeeping_account_handler.go`
  - `service/bookkeeping_account_service.go`
  - `model/bookkeeping_account.go`
  - `service/dto/bookkeeping_account_dto.go`
- 功能说明：
  - 账户的增、删、改、查
  - 支持多种账户类型（现金、储蓄卡、信用卡、支付宝、微信钱包、投资账户等）
  - 支持账户初始余额和当前余额管理
  - 支持设置默认账户

### 3. 交易流水管理（Transaction）
- 主要文件：
  - `api/bookkeeping_transaction_handler.go`
  - `service/bookkeeping_transaction_service.go`
  - `model/bookkeeping_transaction.go`
  - `service/dto/bookkeeping_transaction_dto.go`
- 功能说明：
  - 交易流水的增、删、改、查
  - 支持收入、支出、转账三种交易类型
  - 关联账户和分类
  - 自动更新账户余额（通过 GORM 钩子实现）
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
  - 用户关联交易、账户和分类数据

### 5. 数据统计（Statistics）
- 主要文件：
  - `api/bookkeeping_statistics_handler.go`
  - `service/bookkeeping_statistics_service.go`
  - `service/dto/bookkeeping_statistics_dto.go`
- 功能说明：
  - 收支汇总统计（按日/周/月/年/自定义时间范围）
  - 分类消费占比分析
  - 账户余额概览
  - 月度收支趋势分析

### 6. 预算管理（Budget）
- 主要文件：
  - `api/bookkeeping_budget_handler.go`
  - `service/bookkeeping_budget_service.go`
  - `model/bookkeeping_budget.go`
  - `service/dto/bookkeeping_budget_dto.go`
- 功能说明：
  - 支持总体预算和分类预算
  - 支持周期性预算（周/月/年）
  - 预算进度跟踪和监控
  - 预算提醒（当达到或超过提醒阈值时）
  - 自动计算剩余金额和使用率

## 三、开发建议与后续计划

1. **完善测试**
   - 目前测试覆盖率较低，只有少量测试文件
   - 建议为每个核心模块编写单元测试
   - 增加集成测试和 API 端到端测试
   - 实现测试数据生成器

2. **数据导入导出**
   - 支持 CSV/Excel 格式导入导出
   - 批量交易记录导入
   - 账单数据导出

3. **系统优化**
   - 性能优化，特别是大量数据查询的场景
   - 缓存策略改进
   - 数据库索引优化
   - API 接口限流

4. **前端开发**
   - 基于 Vue/React 开发响应式 Web 界面
   - 移动端适配
   - 数据可视化（图表展示）

5. **文档完善**
   - 更新 API 文档（增加更详细的使用说明）
   - 编写开发者指南
   - 部署文档完善

6. **安全性增强**
   - 敏感数据加密
   - 完善鉴权与授权机制
   - 防止 SQL 注入和 XSS 攻击

## 四、近期开发优先级建议

1. **测试完善**：首先增强测试覆盖率，确保核心功能稳定可靠。
2. **数据导入导出**：作为日常使用的关键功能，便于用户数据迁移和备份。
3. **性能优化**：随着数据量增长，需要确保系统的高效运行。
4. **前端开发**：基于现有 API 开发完整的用户界面。

## 五、技术栈推荐

1. **后端**：继续使用 Go + Gin + GORM
2. **前端**：Vue 3 + Element Plus 或 React + Ant Design
3. **数据库**：MySQL（当前使用）
4. **缓存**：Redis（已在配置中，可启用）
5. **容器化**：Docker + Docker Compose
6. **CI/CD**：GitHub Actions 或 GitLab CI

---

如需继续开发，建议先阅读本文件和 `README.md`，结合 `api/`、`service/`、`model/` 目录下的代码，快速了解当前进度和待办事项。