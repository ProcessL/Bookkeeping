package dto

import (
	"time"
)

// Category 分类DTO（简化版本，用于预算关联）
type Category struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	ParentID *uint  `json:"parent_id,omitempty"`
	Icon     string `json:"icon,omitempty"`
}

// CreateBudgetRequest 创建预算请求
type CreateBudgetRequest struct {
	Name        string    `json:"name" binding:"required"`                                   // 预算名称
	Type        string    `json:"type" binding:"required,oneof=overall category"`            // 预算类型
	Period      string    `json:"period" binding:"required,oneof=weekly monthly yearly"`     // 预算周期
	Amount      float64   `json:"amount" binding:"required,gt=0"`                            // 预算金额
	StartDate   time.Time `json:"start_date" binding:"required"`                             // 开始日期
	CategoryID  *uint     `json:"category_id" binding:"omitempty,required_if=Type category"` // 分类ID
	NotifyRate  *float64  `json:"notify_rate" binding:"omitempty,gte=0,lte=1"`               // 提醒阈值
	Description string    `json:"description"`                                               // 备注
	IsActive    *bool     `json:"is_active"`                                                 // 是否激活
}

// UpdateBudgetRequest 更新预算请求
type UpdateBudgetRequest struct {
	Name        *string    `json:"name"`                                                   // 预算名称
	Type        *string    `json:"type" binding:"omitempty,oneof=overall category"`        // 预算类型
	Period      *string    `json:"period" binding:"omitempty,oneof=weekly monthly yearly"` // 预算周期
	Amount      *float64   `json:"amount" binding:"omitempty,gt=0"`                        // 预算金额
	StartDate   *time.Time `json:"start_date"`                                             // 开始日期
	CategoryID  *uint      `json:"category_id"`                                            // 分类ID
	NotifyRate  *float64   `json:"notify_rate" binding:"omitempty,gte=0,lte=1"`            // 提醒阈值
	Description *string    `json:"description"`                                            // 备注
	IsActive    *bool      `json:"is_active"`                                              // 是否激活
}

// BudgetResponse 预算信息响应
type BudgetResponse struct {
	ID          uint      `json:"id"`                 // 预算ID
	UserID      uint      `json:"user_id"`            // 用户ID
	Name        string    `json:"name"`               // 预算名称
	Type        string    `json:"type"`               // 预算类型
	Period      string    `json:"period"`             // 预算周期
	Amount      float64   `json:"amount"`             // 预算金额
	StartDate   time.Time `json:"start_date"`         // 开始日期
	CategoryID  *uint     `json:"category_id"`        // 分类ID
	NotifyRate  float64   `json:"notify_rate"`        // 提醒阈值
	Description string    `json:"description"`        // 备注
	IsActive    bool      `json:"is_active"`          // 是否激活
	CreatedAt   time.Time `json:"created_at"`         // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`         // 更新时间
	Category    *Category `json:"category,omitempty"` // 关联的分类
}

// BudgetProgressResponse 预算进度响应
type BudgetProgressResponse struct {
	BudgetResponse          // 嵌入预算基本信息
	SpentAmount     float64 `json:"spent_amount"`     // 已花费金额
	RemainingAmount float64 `json:"remaining_amount"` // 剩余金额
	UsageRate       float64 `json:"usage_rate"`       // 使用率 (0-1.0)
	IsOverBudget    bool    `json:"is_over_budget"`   // 是否超出预算
	DaysRemaining   int     `json:"days_remaining"`   // 周期内剩余天数
	CurrentPeriod   struct {
		StartDate time.Time `json:"start_date"` // 当前周期开始日期
		EndDate   time.Time `json:"end_date"`   // 当前周期结束日期
	} `json:"current_period"` // 当前预算周期
}

// BudgetListResponse 预算列表响应
type BudgetListResponse struct {
	Total int64            `json:"total"` // 总数
	Items []BudgetResponse `json:"items"` // 预算列表
}

// BudgetQuery 预算查询参数
type BudgetQuery struct {
	Page       int    `form:"page" json:"page" binding:"required,min=1"`                   // 页码
	PageSize   int    `form:"page_size" json:"page_size" binding:"required,min=1,max=100"` // 每页大小
	Type       string `form:"type" json:"type"`                                            // 预算类型
	Period     string `form:"period" json:"period"`                                        // 预算周期
	CategoryID uint   `form:"category_id" json:"category_id"`                              // 分类ID
	IsActive   *bool  `form:"is_active" json:"is_active"`                                  // 是否激活
}
