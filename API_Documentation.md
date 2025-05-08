# GoGoFly Bookkeeping API文档

## 基础信息
- 基础URL: http://localhost:8090
- API版本: v1.0.1

## 认证
大部分API需要通过Token认证。在请求头中添加 `x-token` 字段。

## 用户模块 API

### 1. 用户登录
- **URL**: `/api/v1/public/user/login`
- **方法**: POST
- **Content-Type**: application/json
- **描述**: 用户登录接口
- **参数**:
  ```json
  {
    "username": "用户名",
    "password": "密码"
  }
  ```
- **响应**:
  ```json
  {
    "code": 0,
    "data": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "user": {
        "id": 1,
        "username": "用户名",
        "email": "邮箱"
      }
    },
    "msg": "登录成功"
  }
  ```

### 2. 用户注册
- **URL**: `/api/v1/auth/user/addUser`
- **方法**: POST
- **Content-Type**: application/json
- **描述**: 用户注册接口
- **参数**:
  ```json
  {
    "username": "用户名",
    "password": "密码",
    "email": "邮箱",
    "phone": "电话号码"
  }
  ```
- **响应**:
  ```json
  {
    "code": 0,
    "data": null,
    "msg": "注册成功"
  }
  ```

### 3. 获取用户详情
- **URL**: `/api/v1/auth/user/getUserById/{id}`
- **方法**: GET
- **描述**: 获取用户详细信息
- **参数**: 
  - id: 用户ID (路径参数)
- **响应**:
  ```json
  {
    "code": 0,
    "data": {
      "id": 1,
      "username": "用户名",
      "email": "邮箱"
    },
    "msg": "获取成功"
  }
  ```

### 4. 用户分析
- **URL**: `/api/v1/public/user/analysis/{id}`
- **方法**: GET
- **描述**: 获取用户数据分析
- **参数**: 
  - id: 用户ID (路径参数)
- **响应**: 返回用户数据分析结果

### 5. 获取扫描结果
- **URL**: `/api/v1/public/user/scanResult`
- **方法**: GET
- **描述**: 获取用户扫描结果
- **参数**:
  ```json
  {
    "pageIndex": 1,
    "pageSize": 10
  }
  ```
- **响应**: 返回扫描结果数据

## 记账模块 API

### 账户管理

#### 1. 获取账户列表
- **URL**: `/bk/accounts`
- **方法**: GET
- **描述**: 获取当前用户的所有账户
- **请求头**: 
  - x-token: 用户令牌
- **响应**: 
  ```json
  {
    "code": 0,
    "data": [
      {
        "id": 1,
        "name": "现金账户",
        "type": "cash",
        "initial_balance": 1000.00,
        "current_balance": 1200.00,
        "is_default": true,
        "user_id": 1,
        "created_at": "2023-05-01T12:00:00Z",
        "updated_at": "2023-05-01T12:00:00Z"
      }
    ],
    "msg": "获取成功"
  }
  ```

#### 2. 创建账户
- **URL**: `/bk/accounts`
- **方法**: POST
- **Content-Type**: application/json
- **描述**: 创建新账户
- **请求头**: 
  - x-token: 用户令牌
- **参数**:
  ```json
  {
    "name": "账户名称",
    "type": "cash",
    "initial_balance": 1000.00,
    "is_default": true,
    "remark": "备注"
  }
  ```
- **响应**: 返回创建的账户信息

#### 3. 获取单个账户
- **URL**: `/bk/accounts/{id}`
- **方法**: GET
- **描述**: 获取指定ID的账户详情
- **请求头**: 
  - x-token: 用户令牌
- **参数**:
  - id: 账户ID (路径参数)
- **响应**: 返回账户详情

#### 4. 更新账户
- **URL**: `/bk/accounts/{id}`
- **方法**: PUT
- **Content-Type**: application/json
- **描述**: 更新指定ID的账户信息
- **请求头**: 
  - x-token: 用户令牌
- **参数**:
  - id: 账户ID (路径参数)
  ```json
  {
    "name": "新账户名称",
    "type": "savings",
    "is_default": false,
    "remark": "新备注"
  }
  ```
- **响应**: 返回更新后的账户信息

#### 5. 删除账户
- **URL**: `/bk/accounts/{id}`
- **方法**: DELETE
- **描述**: 删除指定ID的账户
- **请求头**: 
  - x-token: 用户令牌
- **参数**:
  - id: 账户ID (路径参数)
- **响应**: 返回删除结果

### 分类管理

#### 1. 获取分类列表 (层级)
- **URL**: `/bk/categories`
- **方法**: GET
- **描述**: 获取交易分类列表，按层级结构返回
- **请求头**: 
  - x-token: 用户令牌
- **查询参数**:
  - type: 分类类型 (income/expense)
  - parent_id: 父分类ID
- **响应**: 返回层级结构的分类列表

#### 2. 获取分类列表 (扁平)
- **URL**: `/bk/categories/flat`
- **方法**: GET
- **描述**: 获取所有分类，以扁平列表形式返回
- **请求头**: 
  - x-token: 用户令牌
- **查询参数**:
  - type: 分类类型 (income/expense)
- **响应**: 返回扁平结构的分类列表

#### 3. 创建分类
- **URL**: `/bk/categories`
- **方法**: POST
- **Content-Type**: application/json
- **描述**: 创建交易分类
- **请求头**: 
  - x-token: 用户令牌
- **参数**:
  ```json
  {
    "name": "分类名称",
    "type": "income/expense",
    "icon": "图标名称",
    "parent_id": 0,
    "sort_order": 1
  }
  ```
- **响应**: 返回创建的分类信息

#### 4. 获取单个分类
- **URL**: `/bk/categories/{id}`
- **方法**: GET
- **描述**: 获取指定ID的分类详情
- **请求头**: 
  - x-token: 用户令牌
- **参数**:
  - id: 分类ID (路径参数)
- **响应**: 返回分类详情，包含子分类

#### 5. 更新分类
- **URL**: `/bk/categories/{id}`
- **方法**: PUT
- **Content-Type**: application/json
- **描述**: 更新指定ID的分类信息
- **请求头**: 
  - x-token: 用户令牌
- **参数**:
  - id: 分类ID (路径参数)
  ```json
  {
    "name": "新分类名称",
    "icon": "新图标",
    "parent_id": 2,
    "sort_order": 2
  }
  ```
- **响应**: 返回更新后的分类信息

#### 6. 删除分类
- **URL**: `/bk/categories/{id}`
- **方法**: DELETE
- **描述**: 删除指定ID的分类
- **请求头**: 
  - x-token: 用户令牌
- **参数**:
  - id: 分类ID (路径参数)
- **响应**: 返回删除结果

### 交易管理

#### 1. 获取交易列表
- **URL**: `/bk/transactions`
- **方法**: GET
- **描述**: 获取交易记录列表，支持分页和筛选
- **请求头**: 
  - x-token: 用户令牌
- **查询参数**:
  - page: 页码，默认1
  - page_size: 每页数量，默认20
  - account_id: 账户ID筛选
  - category_id: 分类ID筛选
  - type: 交易类型筛选 (income, expense, transfer)
  - start_date: 开始日期筛选 (YYYY-MM-DD)
  - end_date: 结束日期筛选 (YYYY-MM-DD)
- **响应**: 返回交易记录列表

#### 2. 创建交易
- **URL**: `/bk/transactions`
- **方法**: POST
- **Content-Type**: application/json
- **描述**: 创建交易记录
- **请求头**: 
  - x-token: 用户令牌
- **参数**:
  ```json
  {
    "account_id": 1,
    "category_id": 1,
    "amount": 100.00,
    "type": "income/expense/transfer",
    "transaction_date": "2023-05-01",
    "notes": "交易备注",
    "payee_payer": "收款方/付款方"
  }
  ```
- **响应**: 返回创建的交易记录

#### 3. 获取单个交易
- **URL**: `/bk/transactions/{id}`
- **方法**: GET
- **描述**: 获取指定ID的交易详情
- **请求头**: 
  - x-token: 用户令牌
- **参数**:
  - id: 交易ID (路径参数)
- **响应**: 返回交易详情

#### 4. 更新交易
- **URL**: `/bk/transactions/{id}`
- **方法**: PUT
- **Content-Type**: application/json
- **描述**: 更新指定ID的交易信息
- **请求头**: 
  - x-token: 用户令牌
- **参数**:
  - id: 交易ID (路径参数)
  ```json
  {
    "account_id": 1,
    "category_id": 2,
    "amount": 150.00,
    "type": "expense",
    "transaction_date": "2023-05-02",
    "notes": "更新后的备注",
    "payee_payer": "更新后的付款方"
  }
  ```
- **响应**: 返回更新后的交易信息

#### 5. 删除交易
- **URL**: `/bk/transactions/{id}`
- **方法**: DELETE
- **描述**: 删除指定ID的交易
- **请求头**: 
  - x-token: 用户令牌
- **参数**:
  - id: 交易ID (路径参数)
- **响应**: 返回删除结果

### 预算管理

#### 1. 获取预算列表
- **URL**: `/bk/budgets`
- **方法**: GET
- **描述**: 获取预算列表，支持分页和筛选
- **请求头**: 
  - x-token: 用户令牌
- **查询参数**:
  - page: 页码，默认1
  - page_size: 每页大小，默认10
  - type: 预算类型 (overall, category)
  - period: 预算周期 (weekly, monthly, yearly)
  - category_id: 分类ID
  - is_active: 是否激活
- **响应**: 返回预算列表

#### 2. 创建预算
- **URL**: `/bk/budgets`
- **方法**: POST
- **Content-Type**: application/json
- **描述**: 创建新的预算
- **请求头**: 
  - x-token: 用户令牌
- **参数**:
  ```json
  {
    "name": "预算名称",
    "type": "overall/category",
    "amount": 1000.00,
    "period": "weekly/monthly/yearly",
    "category_id": 1,
    "start_date": "2023-05-01",
    "is_active": true,
    "notify_rate": 0.8,
    "description": "预算描述"
  }
  ```
- **响应**: 返回创建的预算信息

#### 3. 获取预算详情
- **URL**: `/bk/budgets/{id}`
- **方法**: GET
- **描述**: 获取单个预算的详细信息
- **请求头**: 
  - x-token: 用户令牌
- **参数**:
  - id: 预算ID (路径参数)
- **响应**: 返回预算详情

#### 4. 更新预算
- **URL**: `/bk/budgets/{id}`
- **方法**: PUT
- **Content-Type**: application/json
- **描述**: 更新预算信息
- **请求头**: 
  - x-token: 用户令牌
- **参数**:
  - id: 预算ID (路径参数)
  ```json
  {
    "name": "新预算名称",
    "amount": 1500.00,
    "is_active": false,
    "notify_rate": 0.7
  }
  ```
- **响应**: 返回更新后的预算信息

#### 5. 删除预算
- **URL**: `/bk/budgets/{id}`
- **方法**: DELETE
- **描述**: 删除预算
- **请求头**: 
  - x-token: 用户令牌
- **参数**:
  - id: 预算ID (路径参数)
- **响应**: 返回删除结果

#### 6. 获取预算进度
- **URL**: `/bk/budgets/{id}/progress`
- **方法**: GET
- **描述**: 获取单个预算的当前执行进度
- **请求头**: 
  - x-token: 用户令牌
- **参数**:
  - id: 预算ID (路径参数)
- **响应**: 返回预算进度信息

#### 7. 获取所有激活预算进度
- **URL**: `/bk/budgets/active-progress`
- **方法**: GET
- **描述**: 获取所有激活的预算及其当前执行进度
- **请求头**: 
  - x-token: 用户令牌
- **响应**: 返回所有激活预算的进度信息

#### 8. 检查预算警告
- **URL**: `/bk/budgets/alerts`
- **方法**: GET
- **描述**: 获取达到或超过提醒阈值的预算列表
- **请求头**: 
  - x-token: 用户令牌
- **响应**: 返回需要提醒的预算列表

### 统计分析

#### 1. 获取账户余额汇总
- **URL**: `/statistics/account-summary`
- **方法**: GET
- **描述**: 获取所有账户的余额汇总信息
- **请求头**: 
  - x-token: 用户令牌
- **响应**: 返回账户余额汇总信息

#### 2. 获取分类汇总
- **URL**: `/statistics/category-summary`
- **方法**: GET
- **描述**: 获取指定时间范围内的分类汇总信息
- **请求头**: 
  - x-token: 用户令牌
- **查询参数**:
  - range_type: 时间范围类型 (day/week/month/year/all/custom)
  - transaction_type: 交易类型 (income/expense)
  - start_date: 自定义开始日期
  - end_date: 自定义结束日期
  - months_count: 查询的月份数量
- **响应**: 返回分类汇总数据

#### 3. 获取收支汇总
- **URL**: `/statistics/income-expense-summary`
- **方法**: GET
- **描述**: 获取指定时间范围内的收支汇总信息
- **请求头**: 
  - x-token: 用户令牌
- **查询参数**:
  - range_type: 时间范围类型 (day/week/month/year/all/custom)
  - start_date: 自定义开始日期
  - end_date: 自定义结束日期
  - months_count: 查询的月份数量
- **响应**: 返回收支汇总数据

#### 4. 获取月度收支趋势
- **URL**: `/statistics/monthly-trend`
- **方法**: GET
- **描述**: 获取最近几个月的收支趋势数据
- **请求头**: 
  - x-token: 用户令牌
- **查询参数**:
  - months_count: 查询的月份数量，默认为12
- **响应**: 返回月度收支趋势数据

## 错误码
- 0: 成功
- 7: 请求参数错误
- 1000: 服务器内部错误
- 1001: 未授权
- 1002: 令牌过期
- 1003: 令牌无效
- 1004: 用户不存在
- 1005: 密码错误

## API调用示例

### 用户登录示例
```javascript
const login = async () => {
  try {
    const response = await axios.post('http://localhost:8090/api/v1/public/user/login', {
      username: 'testuser',
      password: 'password123'
    });
    
    // 保存令牌
    localStorage.setItem('token', response.data.data.token);
    
    console.log('登录成功:', response.data);
  } catch (error) {
    console.error('登录失败:', error);
  }
};
```

### 获取账户列表示例
```javascript
const getAccounts = async () => {
  try {
    const token = localStorage.getItem('token');
    
    const response = await axios.get('http://localhost:8090/bk/accounts', {
      headers: {
        'x-token': token
      }
    });
    
    console.log('账户列表:', response.data);
  } catch (error) {
    console.error('获取账户失败:', error);
  }
};
```

### 创建交易示例
```javascript
const createTransaction = async () => {
  try {
    const token = localStorage.getItem('token');
    
    const transaction = {
      account_id: 1,
      category_id: 2,
      amount: 100.50,
      type: 'expense',
      transaction_date: '2023-05-08',
      notes: '超市购物',
      payee_payer: '沃尔玛'
    };
    
    const response = await axios.post('http://localhost:8090/bk/transactions', transaction, {
      headers: {
        'x-token': token
      }
    });
    
    console.log('创建交易成功:', response.data);
  } catch (error) {
    console.error('创建交易失败:', error);
  }
};
```

### 获取统计数据示例
```javascript
const getStatistics = async () => {
  try {
    const token = localStorage.getItem('token');
    
    const response = await axios.get('http://localhost:8090/statistics/income-expense-summary', {
      headers: {
        'x-token': token
      },
      params: {
        range_type: 'month'
      }
    });
    
    console.log('本月收支统计:', response.data);
  } catch (error) {
    console.error('获取统计数据失败:', error);
  }
}; 