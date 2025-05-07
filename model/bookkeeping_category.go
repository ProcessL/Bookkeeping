package model

import "github.com/dotdancer/gogofly/global"

// CategoryType 定义分类类型
type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "income"  // 收入分类
	CategoryTypeExpense CategoryType = "expense" // 支出分类
)

// Category 分类模型
type Category struct {
	global.GlyModel
	UserID    uint         `json:"user_id" gorm:"index;comment:用户ID"`
	Name      string       `json:"name" gorm:"type:varchar(100);not null;comment:分类名称"`
	Type      CategoryType `json:"type" gorm:"type:varchar(50);not null;comment:分类类型 (income, expense)"`
	ParentID  *uint        `json:"parent_id" gorm:"index;comment:父分类ID (用于支持多级分类)"` // 指针类型，允许为空
	Icon      string       `json:"icon" gorm:"type:varchar(100);comment:图标 (可选)"`
	SortOrder int          `json:"sort_order" gorm:"default:0;comment:排序字段"`

	// Associations
	ParentCategory *Category  `json:"parent_category,omitempty" gorm:"foreignKey:ParentID"`
	SubCategories  []Category `json:"sub_categories,omitempty" gorm:"foreignKey:ParentID"`
}

// TableName 指定表名
func (c *Category) TableName() string {
	return "bookkeeping_categories"
}
