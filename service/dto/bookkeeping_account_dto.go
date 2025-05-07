package dto

import (
	"github.com/dotdancer/gogofly/model"
)

// CreateAccountRequest 创建账户的请求体
type CreateAccountRequest struct {
	Name           string            `json:"name" binding:"required,min=1,max=100"`        // 账户名称
	Type           model.AccountType `json:"type" binding:"required"`                      // 账户类型
	InitialBalance float64           `json:"initial_balance" binding:"omitempty,min=0"`    // 初始余额
	Remark         string            `json:"remark,omitempty" binding:"omitempty,max=255"` // 备注
	IsDefault      bool              `json:"is_default,omitempty"`                         // 是否默认账户
}

// UpdateAccountRequest 更新账户的请求体
type UpdateAccountRequest struct {
	Name      *string            `json:"name,omitempty" binding:"omitempty,min=1,max=100"` // 账户名称
	Type      *model.AccountType `json:"type,omitempty"`                                   // 账户类型
	Remark    *string            `json:"remark,omitempty" binding:"omitempty,max=255"`     // 备注
	IsDefault *bool              `json:"is_default,omitempty"`                             // 是否默认账户
}

// AccountResponse 单个账户的响应体
type AccountResponse struct {
	ID             uint              `json:"id"`
	Name           string            `json:"name"`
	Type           model.AccountType `json:"type"`
	InitialBalance float64           `json:"initial_balance"`
	CurrentBalance float64           `json:"current_balance"`
	Remark         string            `json:"remark,omitempty"`
	IsDefault      bool              `json:"is_default"`
	CreatedAt      string            `json:"created_at"`
	UpdatedAt      string            `json:"updated_at"`
	UserID         uint              `json:"user_id"`
}

// AccountListResponse 账户列表的响应体
type AccountListResponse struct {
	Total int64             `json:"total"`
	Items []AccountResponse `json:"items"`
}
