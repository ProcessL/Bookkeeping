# 记账系统后端

这是一个基于Go语言的记账系统后端项目，采用Gin框架，提供完整的记账功能API。

## 功能特性

### 用户管理
- 用户注册和登录
- JWT身份验证
- 用户信息管理

### 分类管理
- 支持收入和支出两种分类类型
- 多级分类结构（父子分类）
- 完整的CRUD操作

### 账户管理
- 支持多种账户类型（现金、储蓄卡、信用卡、支付宝、微信钱包等）
- 账户余额自动计算
- 默认账户设置

### 交易记录管理
- 支持收入、支出和转账三种交易类型
- 关联分类和账户
- 自动更新账户余额
- 丰富的查询筛选条件

### 数据统计
- 收支汇总统计（日/周/月/年）
- 分类消费占比分析
- 账户余额概览
- 月度收支趋势分析

### 预算管理
- 支持总体预算和分类预算
- 支持周/月/年周期设置
- 预算进度跟踪
- 预算超支提醒
- 自动计算剩余金额和使用率

## 项目结构

```
├── api/            // 路由处理层，负责HTTP请求的接收与响应
├── service/        // 业务逻辑层，处理具体的业务操作
├── model/          // 数据模型层，定义数据库表结构及相关类型
├── router/         // 路由注册与分组
├── utils/          // 工具函数
├── config/         // 配置相关
├── core/           // 框架核心初始化
├── global/         // 全局变量与对象
├── middleware/     // 中间件
├── docs/           // Swagger文档
├── cmd/            // 应用初始化命令
├── test/           // 测试文件
└── main.go         // 程序入口
```

## API文档

系统使用Swagger自动生成API文档，启动服务后可以通过以下地址访问：

```
http://localhost:8000/swagger/index.html
```

### 主要API路由

#### 用户相关
- `POST /api/user/register` - 用户注册
- `POST /api/user/login` - 用户登录

#### 分类管理
- `GET /api/bk/categories` - 获取分类列表(层级)
- `GET /api/bk/categories/flat` - 获取所有分类(扁平)
- `POST /api/bk/categories` - 创建分类
- `GET /api/bk/categories/:id` - 获取单个分类
- `PUT /api/bk/categories/:id` - 更新分类
- `DELETE /api/bk/categories/:id` - 删除分类

#### 账户管理
- `GET /api/bk/accounts` - 获取账户列表
- `POST /api/bk/accounts` - 创建账户
- `GET /api/bk/accounts/:id` - 获取单个账户
- `PUT /api/bk/accounts/:id` - 更新账户
- `DELETE /api/bk/accounts/:id` - 删除账户

#### 交易记录
- `GET /api/bk/transactions` - 获取交易记录列表
- `POST /api/bk/transactions` - 创建交易记录
- `GET /api/bk/transactions/:id` - 获取单条交易记录
- `PUT /api/bk/transactions/:id` - 更新交易记录
- `DELETE /api/bk/transactions/:id` - 删除交易记录

#### 统计分析
- `GET /api/bk/statistics/income-expense-summary` - 获取收支汇总
- `GET /api/bk/statistics/category-summary` - 获取分类汇总
- `GET /api/bk/statistics/account-summary` - 获取账户余额汇总
- `GET /api/bk/statistics/monthly-trend` - 获取月度收支趋势

#### 预算管理
- `POST /api/bk/budgets` - 创建预算
- `GET /api/bk/budgets` - 获取预算列表
- `GET /api/bk/budgets/active-progress` - 获取所有激活预算的进度
- `GET /api/bk/budgets/alerts` - 获取预算提醒（超过阈值的预算）
- `GET /api/bk/budgets/:id` - 获取单个预算
- `GET /api/bk/budgets/:id/progress` - 获取预算进度
- `PUT /api/bk/budgets/:id` - 更新预算
- `DELETE /api/bk/budgets/:id` - 删除预算

## 如何运行

1. 克隆项目
```bash
git clone https://github.com/your-username/your-repo.git
cd your-repo
```

2. 安装依赖
```bash
go mod tidy
```

3. 配置数据库
编辑`config.yaml`文件，配置数据库连接信息

4. 运行项目
```bash
go run main.go
```

## 技术栈
- Go 1.16+
- Gin Web Framework
- GORM
- JWT
- MySQL

## 许可证
[MIT](LICENSE)
