package dto

import "time"

// IncomeExpenseSummaryResponse 收支汇总响应
type IncomeExpenseSummaryResponse struct {
	TotalIncome  float64   `json:"total_income"`   // 总收入
	TotalExpense float64   `json:"total_expense"`  // 总支出
	NetAmount    float64   `json:"net_amount"`     // 净收入（收入-支出）
	StartDate    time.Time `json:"start_date"`     // 统计开始日期
	EndDate      time.Time `json:"end_date"`       // 统计结束日期
	RangeType    string    `json:"range_type"`     // 时间范围类型（day, week, month, year, all, custom）
}

// CategorySummaryItem 分类汇总项
type CategorySummaryItem struct {
	CategoryID       uint    `json:"category_id"`        // 分类ID
	CategoryName     string  `json:"category_name"`      // 分类名称
	CategoryIcon     string  `json:"category_icon"`      // 分类图标
	TotalAmount      float64 `json:"total_amount"`       // 总金额
	TransactionCount int     `json:"transaction_count"`  // 交易笔数
}

// AccountSummaryItem 账户汇总项
type AccountSummaryItem struct {
	AccountID      uint    `json:"account_id"`       // 账户ID
	AccountName    string  `json:"account_name"`     // 账户名称
	AccountType    string  `json:"account_type"`     // 账户类型
	CurrentBalance float64 `json:"current_balance"`  // 当前余额
	InitialBalance float64 `json:"initial_balance"`  // 初始余额
}

// MonthlyData 月度数据
type MonthlyData struct {
	Year         int     `json:"year"`           // 年份
	Month        int     `json:"month"`          // 月份
	MonthLabel   string  `json:"month_label"`    // 月份标签，格式：YYYY-MM
	TotalIncome  float64 `json:"total_income"`   // 该月总收入
	TotalExpense float64 `json:"total_expense"`  // 该月总支出
	NetAmount    float64 `json:"net_amount"`     // 该月净收入（收入-支出）
}

// MonthlyTrendResponse 月度趋势响应
type MonthlyTrendResponse struct {
	MonthsCount int           `json:"months_count"` // 查询的月份数量
	Data        []*MonthlyData `json:"data"`        // 月度数据列表
}

// StatisticsQueryRequest 统计查询请求
type StatisticsQueryRequest struct {
	RangeType   string     `json:"range_type" form:"range_type" binding:"required,oneof=day week month year all custom"` // 时间范围类型
	StartDate   *time.Time `json:"start_date" form:"start_date"`   // 自定义开始日期（当 range_type 为 custom 时必填）
	EndDate     *time.Time `json:"end_date" form:"end_date"`       // 自定义结束日期（当 range_type 为 custom 时必填）
	MonthsCount int        `json:"months_count" form:"months_count"` // 查询的月份数量，用于月度趋势统计
}

// CategoryStatisticsRequest 分类统计请求
type CategoryStatisticsRequest struct {
	StatisticsQueryRequest
	TransactionType string `json:"transaction_type" form:"transaction_type" binding:"required,oneof=income expense"` // 交易类型：收入或支出
} 