package dto

import "github.com/dotdancer/gogofly/model"

// CreateCategoryRequest 创建分类的请求体
type CreateCategoryRequest struct {
	Name      string             `json:"name" binding:"required,min=1,max=100"`        // 分类名称
	Type      model.CategoryType `json:"type" binding:"required,oneof=income expense"` // 分类类型 (income, expense)
	ParentID  *uint              `json:"parent_id"`                                    // 父分类ID (可选)
	Icon      string             `json:"icon,omitempty" binding:"omitempty,max=100"`   // 图标 (可选)
	SortOrder int                `json:"sort_order,omitempty"`                         // 排序字段 (可选)
}

// UpdateCategoryRequest 更新分类的请求体
type UpdateCategoryRequest struct {
	Name      *string            `json:"name,omitempty" binding:"omitempty,min=1,max=100"`        // 分类名称 (可选)
	Type      model.CategoryType `json:"type,omitempty" binding:"omitempty,oneof=income expense"` // 分类类型 (可选)
	ParentID  *uint              `json:"parent_id,omitempty"`                                     // 父分类ID (可选, 注意: 如果要移除父分类，应传递null或不传递该字段，具体行为需在service层定义)
	Icon      *string            `json:"icon,omitempty" binding:"omitempty,max=100"`              // 图标 (可选)
	SortOrder *int               `json:"sort_order,omitempty"`                                    // 排序字段 (可选)
}

// CategoryResponse 单个分类的响应体
type CategoryResponse struct {
	ID            uint               `json:"id"`
	Name          string             `json:"name"`
	Type          model.CategoryType `json:"type"`
	ParentID      *uint              `json:"parent_id,omitempty"`
	Icon          string             `json:"icon,omitempty"`
	SortOrder     int                `json:"sort_order"`
	CreatedAt     string             `json:"created_at"`
	UpdatedAt     string             `json:"updated_at"`
	SubCategories []CategoryResponse `json:"sub_categories,omitempty"` // 子分类列表 (递归展示)
	UserID        uint               `json:"user_id"`
}

// CategoryListResponse 分类列表的响应体 (可包含分页信息)
type CategoryListResponse struct {
	Total int64              `json:"total"`
	Items []CategoryResponse `json:"items"`
}
