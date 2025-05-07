package dto

import (
	"github.com/dotdancer/gogofly/model"
)

// CreateTransactionRequest 创建交易流水的请求体
type CreateTransactionRequest struct {
	AccountID       uint                  `json:"account_id" binding:"required"`                         // 账户ID
	Type            model.TransactionType `json:"type" binding:"required,oneof=income expense transfer"` // 交易类型
	Amount          float64               `json:"amount" binding:"required,gt=0"`                        // 金额
	TransactionDate string                `json:"transaction_date" binding:"required"`                   // 交易日期 (YYYY-MM-DD)
	CategoryID      uint                  `json:"category_id" binding:"required"`                        // 分类ID
	PayeePayer      string                `json:"payee_payer,omitempty" binding:"omitempty,max=100"`     // 收款方/付款方
	Notes           string                `json:"notes,omitempty" binding:"omitempty,max=255"`           // 备注
}

// UpdateTransactionRequest 更新交易流水的请求体
type UpdateTransactionRequest struct {
	AccountID       *uint                  `json:"account_id,omitempty"`                                             // 账户ID
	Type            *model.TransactionType `json:"type,omitempty" binding:"omitempty,oneof=income expense transfer"` // 交易类型
	Amount          *float64               `json:"amount,omitempty" binding:"omitempty,gt=0"`                        // 金额
	TransactionDate *string                `json:"transaction_date,omitempty"`                                       // 交易日期
	CategoryID      *uint                  `json:"category_id,omitempty"`                                            // 分类ID
	PayeePayer      *string                `json:"payee_payer,omitempty" binding:"omitempty,max=100"`                // 收款方/付款方
	Notes           *string                `json:"notes,omitempty" binding:"omitempty,max=255"`                      // 备注
}

// TransactionResponse 单个交易流水的响应体
type TransactionResponse struct {
	ID              uint                  `json:"id"`
	AccountID       uint                  `json:"account_id"`
	Type            model.TransactionType `json:"type"`
	Amount          float64               `json:"amount"`
	TransactionDate string                `json:"transaction_date"` // 格式化为 YYYY-MM-DD
	CategoryID      uint                  `json:"category_id"`
	PayeePayer      string                `json:"payee_payer,omitempty"`
	Notes           string                `json:"notes,omitempty"`
	CreatedAt       string                `json:"created_at"`
	UpdatedAt       string                `json:"updated_at"`
	UserID          uint                  `json:"user_id"`

	// 关联信息
	Account  AccountResponse  `json:"account,omitempty"`
	Category CategoryResponse `json:"category,omitempty"`
}

// TransactionQuery 交易流水查询条件
type TransactionQuery struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	AccountID  uint   `json:"account_id,omitempty"`
	CategoryID uint   `json:"category_id,omitempty"`
	Type       string `json:"type,omitempty"`
	StartDate  string `json:"start_date,omitempty"`
	EndDate    string `json:"end_date,omitempty"`
}

// TransactionListResponse 交易流水列表的响应体
type TransactionListResponse struct {
	Total int64                 `json:"total"`
	Items []TransactionResponse `json:"items"`
}
