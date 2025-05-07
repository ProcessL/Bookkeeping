package model

import (
	"time"

	"github.com/dotdancer/gogofly/global"
)

// BudgetType 预算类型
type BudgetType string

const (
	BudgetTypeOverall  BudgetType = "overall"  // 总体预算
	BudgetTypeCategory BudgetType = "category" // 分类预算
)

// BudgetPeriod 预算周期
type BudgetPeriod string

const (
	BudgetPeriodWeekly  BudgetPeriod = "weekly"  // 每周
	BudgetPeriodMonthly BudgetPeriod = "monthly" // 每月
	BudgetPeriodYearly  BudgetPeriod = "yearly"  // 每年
)

// Budget 预算模型
type Budget struct {
	global.GlyModel
	UserID      uint         `json:"user_id" gorm:"index;comment:用户ID"`
	Name        string       `json:"name" gorm:"type:varchar(100);not null;comment:预算名称"`
	Type        BudgetType   `json:"type" gorm:"type:varchar(50);not null;comment:预算类型 (overall, category)"`
	Period      BudgetPeriod `json:"period" gorm:"type:varchar(50);not null;comment:预算周期 (weekly, monthly, yearly)"`
	Amount      float64      `json:"amount" gorm:"type:decimal(10,2);not null;comment:预算金额"`
	StartDate   time.Time    `json:"start_date" gorm:"not null;comment:开始日期"`
	CategoryID  *uint        `json:"category_id" gorm:"index;comment:分类ID (当类型为分类预算时使用)"` // 指针类型，允许为空
	NotifyRate  float64      `json:"notify_rate" gorm:"type:decimal(5,2);default:0.80;comment:提醒阈值 (如: 0.8 表示达到80%时提醒)"`
	Description string       `json:"description" gorm:"type:varchar(255);comment:备注"`
	IsActive    bool         `json:"is_active" gorm:"default:true;comment:是否激活"`

	// Associations
	Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"` // 关联的分类
}

// TableName 指定表名
func (b *Budget) TableName() string {
	return "bookkeeping_budgets"
}

// BudgetProgress 预算进度视图模型（非数据库表）
// 用于查询和展示预算进度
type BudgetProgress struct {
	Budget          // 嵌入预算模型
	SpentAmount     float64 `json:"spent_amount"`     // 已花费金额
	RemainingAmount float64 `json:"remaining_amount"` // 剩余金额
	UsageRate       float64 `json:"usage_rate"`       // 使用率 (0-1.0)
	IsOverBudget    bool    `json:"is_over_budget"`   // 是否超出预算
	DaysRemaining   int     `json:"days_remaining"`   // 周期内剩余天数
} 