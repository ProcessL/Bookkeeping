package service

import (
	"errors"

	"github.com/dotdancer/gogofly/global"
	"github.com/dotdancer/gogofly/model"
	"github.com/dotdancer/gogofly/service/dto"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// BookkeepingCategoryService 结构体定义了分类管理的服务层
type BookkeepingCategoryService struct{}

// CreateCategory 创建一个新的分类
// userID: 当前操作的用户ID
// req: 创建分类的请求数据
func (s *BookkeepingCategoryService) CreateCategory(userID uint, req dto.CreateCategoryRequest) (dto.CategoryResponse, error) {
	var category model.Category
	var response dto.CategoryResponse

	// 检查同名分类是否已存在 (同一用户、同一父分类下)
	var existingCategory model.Category
	query := global.DB.Where("user_id = ? AND name = ? AND type = ?", userID, req.Name, req.Type)
	if req.ParentID != nil {
		query = query.Where("parent_id = ?", *req.ParentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}
	if err := query.First(&existingCategory).Error; err == nil {
		return response, errors.New("该分类名称已存在")
	}

	// 复制请求数据到模型
	if err := copier.Copy(&category, &req); err != nil {
		global.Logger.Error("Failed to copy CreateCategoryRequest to model.Category: " + err.Error())
		return response, errors.New("创建分类失败：数据复制错误")
	}

	category.UserID = userID

	// 如果有父分类ID，校验父分类是否存在且属于当前用户
	if category.ParentID != nil && *category.ParentID != 0 {
		var parentCategory model.Category
		if err := global.DB.First(&parentCategory, "id = ? AND user_id = ?", *category.ParentID, userID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response, errors.New("指定的父分类不存在或不属于您")
			}
			global.Logger.Error("Failed to find parent category: " + err.Error())
			return response, errors.New("创建分类失败：查询父分类错误")
		}
		// 确保父分类和子分类的类型一致
		if parentCategory.Type != category.Type {
			return response, errors.New("子分类类型必须与父分类类型一致")
		}
	}

	if err := global.DB.Create(&category).Error; err != nil {
		global.Logger.Error("Failed to create category: " + err.Error())
		return response, errors.New("创建分类失败")
	}

	// 复制模型数据到响应体
	if err := copier.Copy(&response, &category); err != nil {
		global.Logger.Error("Failed to copy model.Category to dto.CategoryResponse: " + err.Error())
		// 即使复制失败，分类也已创建，可以考虑返回部分信息或特定错误
	}
	response.CreatedAt = category.CreatedAt.Format("2006-01-02 15:04:05")
	response.UpdatedAt = category.UpdatedAt.Format("2006-01-02 15:04:05")

	return response, nil
}

// GetCategoryByID 根据ID获取单个分类信息
// userID: 当前操作的用户ID
// categoryID: 要获取的分类ID
func (s *BookkeepingCategoryService) GetCategoryByID(userID uint, categoryID uint) (dto.CategoryResponse, error) {
	var category model.Category
	var response dto.CategoryResponse

	if err := global.DB.Preload("SubCategories").First(&category, "id = ? AND user_id = ?", categoryID, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, errors.New("分类不存在")
		}
		global.Logger.Error("Failed to get category by ID: " + err.Error())
		return response, errors.New("获取分类信息失败")
	}

	// 递归加载子分类并转换为CategoryResponse
	var mapCategoryToResponse func(cat model.Category) dto.CategoryResponse
	mapCategoryToResponse = func(cat model.Category) dto.CategoryResponse {
		resp := dto.CategoryResponse{}
		copier.Copy(&resp, &cat)
		resp.CreatedAt = cat.CreatedAt.Format("2006-01-02 15:04:05")
		resp.UpdatedAt = cat.UpdatedAt.Format("2006-01-02 15:04:05")
		if len(cat.SubCategories) > 0 {
			resp.SubCategories = make([]dto.CategoryResponse, len(cat.SubCategories))
			for i, subCat := range cat.SubCategories {
				// 需要从数据库重新加载每个子分类的子分类，以支持多级
				var fullSubCat model.Category
				global.DB.Preload("SubCategories").First(&fullSubCat, subCat.ID)
				resp.SubCategories[i] = mapCategoryToResponse(fullSubCat)
			}
		}
		return resp
	}

	response = mapCategoryToResponse(category)

	return response, nil
}

// UpdateCategory 更新分类信息
// userID: 当前操作的用户ID
// categoryID: 要更新的分类ID
// req: 更新分类的请求数据
func (s *BookkeepingCategoryService) UpdateCategory(userID uint, categoryID uint, req dto.UpdateCategoryRequest) (dto.CategoryResponse, error) {
	var category model.Category
	var response dto.CategoryResponse

	if err := global.DB.First(&category, "id = ? AND user_id = ?", categoryID, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, errors.New("分类不存在")
		}
		global.Logger.Error("Failed to find category for update: " + err.Error())
		return response, errors.New("更新分类失败：未找到分类")
	}

	// 检查更新后的名称是否与同级其他分类冲突
	if req.Name != nil && *req.Name != category.Name {
		var existingCategory model.Category
		query := global.DB.Where("user_id = ? AND name = ? AND type = ? AND id != ?", userID, *req.Name, category.Type, categoryID)
		if category.ParentID != nil {
			query = query.Where("parent_id = ?", *category.ParentID)
		} else {
			query = query.Where("parent_id IS NULL")
		}
		if err := query.First(&existingCategory).Error; err == nil {
			return response, errors.New("该分类名称已存在")
		}
	}

	// 处理父分类ID的更新
	if req.ParentID != nil {
		if *req.ParentID == 0 { // 客户端可能传递0表示移除父分类
			category.ParentID = nil
		} else {
			// 校验新的父分类是否存在且属于当前用户
			var parentCategory model.Category
			if err := global.DB.First(&parentCategory, "id = ? AND user_id = ?", *req.ParentID, userID).Error; err != nil {
				return response, errors.New("指定的父分类不存在或不属于您")
			}
			// 确保父分类和子分类的类型一致
			newType := category.Type
			if req.Type != "" { // 如果请求中也更新了类型，以请求中的为准
				newType = req.Type
			}
			if parentCategory.Type != newType {
				return response, errors.New("子分类类型必须与父分类类型一致")
			}
			// 防止将分类设置为自身的子分类或其子分类的子分类 (循环依赖)
			if *req.ParentID == category.ID {
				return response, errors.New("不能将分类设置为自身的父分类")
			}
			// 更复杂的循环依赖检查可以遍历父级链条，此处简化
			category.ParentID = req.ParentID
		}
	} else if req.ParentID == nil && category.ParentID != nil {
		// 如果请求中没有ParentID字段，但数据库中存在，则不改变它
		// 如果想显式移除父分类，客户端应传递 ParentID: null (或特定值如0，由DTO和service约定)
		// 这里的逻辑是：如果DTO中ParentID字段未提供，则不更新模型的ParentID
		// 如果DTO中ParentID字段提供了，即使是nil，也会尝试更新（上面已处理*req.ParentID == 0的情况）
	}

	// 使用copier选择性更新字段，仅更新DTO中非nil的字段
	// copier默认会覆盖，需要注意指针类型的处理
	// 手动更新以获得更精确的控制
	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Type != "" { // CategoryType不是指针，所以检查空字符串
		category.Type = req.Type
	}
	if req.Icon != nil {
		category.Icon = *req.Icon
	}
	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
	}

	if err := global.DB.Save(&category).Error; err != nil {
		global.Logger.Error("Failed to update category: " + err.Error())
		return response, errors.New("更新分类失败")
	}

	return s.GetCategoryByID(userID, categoryID) // 返回更新后的完整信息，包括可能的子分类
}

// DeleteCategory 删除分类
// userID: 当前操作的用户ID
// categoryID: 要删除的分类ID
func (s *BookkeepingCategoryService) DeleteCategory(userID uint, categoryID uint) error {
	var category model.Category
	if err := global.DB.First(&category, "id = ? AND user_id = ?", categoryID, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("分类不存在")
		}
		global.Logger.Error("Failed to find category for deletion: " + err.Error())
		return errors.New("删除分类失败：未找到分类")
	}

	// 检查是否有子分类
	var subCategoryCount int64
	global.DB.Model(&model.Category{}).Where("parent_id = ? AND user_id = ?", categoryID, userID).Count(&subCategoryCount)
	if subCategoryCount > 0 {
		return errors.New("无法删除：该分类下存在子分类，请先删除或移动子分类")
	}

	// 检查该分类是否被交易流水使用
	var transactionCount int64
	global.DB.Model(&model.Transaction{}).Where("category_id = ? AND user_id = ?", categoryID, userID).Count(&transactionCount)
	if transactionCount > 0 {
		return errors.New("无法删除：该分类已被交易流水使用")
	}

	if err := global.DB.Delete(&category).Error; err != nil {
		global.Logger.Error("Failed to delete category: " + err.Error())
		return errors.New("删除分类失败")
	}

	return nil
}

// ListCategories 获取分类列表 (支持层级)
// userID: 当前操作的用户ID
// categoryType: 可选，按类型过滤 (income/expense)
// parentID: 可选，获取指定父分类下的子分类，若为nil则获取顶级分类
func (s *BookkeepingCategoryService) ListCategories(userID uint, categoryType model.CategoryType, parentID *uint) ([]dto.CategoryResponse, error) {
	var categories []model.Category
	var responses []dto.CategoryResponse

	query := global.DB.Where("user_id = ?", userID).Order("sort_order asc, created_at asc")

	if categoryType != "" {
		query = query.Where("type = ?", categoryType)
	}

	if parentID != nil {
		if *parentID == 0 { // 约定0为顶级分类的查询
			query = query.Where("parent_id IS NULL")
		} else {
			query = query.Where("parent_id = ?", *parentID)
		}
	} else {
		// 如果不传parentID，默认获取所有层级的分类，然后手动构建层级结构
		// 或者，可以约定只获取顶级分类 query = query.Where("parent_id IS NULL")
		// 这里我们获取所有，然后在下面构建层级
	}

	if err := query.Preload("SubCategories").Find(&categories).Error; err != nil {
		global.Logger.Error("Failed to list categories: " + err.Error())
		return nil, errors.New("获取分类列表失败")
	}

	// 辅助函数，将model.Category转换为dto.CategoryResponse并递归处理子分类
	var mapCategoriesToResponse func([]model.Category) []dto.CategoryResponse
	mapCategoriesToResponse = func(cats []model.Category) []dto.CategoryResponse {
		resps := make([]dto.CategoryResponse, 0, len(cats))
		for _, cat := range cats {
			resp := dto.CategoryResponse{}
			copier.Copy(&resp, &cat)
			resp.CreatedAt = cat.CreatedAt.Format("2006-01-02 15:04:05")
			resp.UpdatedAt = cat.UpdatedAt.Format("2006-01-02 15:04:05")
			// 如果Preload了SubCategories，它们会在这里
			// 为了确保多级，我们需要对每个SubCategory也递归调用
			if len(cat.SubCategories) > 0 {
				// 需要从数据库重新加载每个子分类的子分类，以支持多级
				// 或者在初始查询时使用更复杂的 Preload 策略
				// 这里简化：假设 Preload("SubCategories") 只加载了一层
				// 更好的做法是在 GetCategoryByID 那样递归加载
				// 或者，如果 ListCategories 旨在返回扁平列表或仅下一级，则当前 Preload 可能足够
				// 为了返回完整的层级结构，我们需要递归地为子分类填充它们的子分类

				// 重新获取子分类并填充它们的子分类
				var fullSubCategories []model.Category
				subCategoryIDs := make([]uint, len(cat.SubCategories))
				for i, sub := range cat.SubCategories {
					subCategoryIDs[i] = sub.ID
				}
				global.DB.Preload("SubCategories").Where("id IN ?", subCategoryIDs).Order("sort_order asc, created_at asc").Find(&fullSubCategories)
				resp.SubCategories = mapCategoriesToResponse(fullSubCategories)
			}
			resps = append(resps, resp)
		}
		return resps
	}

	// 如果 parentID 为 nil，表示希望获取所有分类并构建完整的树状结构
	// 否则，只获取指定 parentID 下的子分类
	if parentID == nil {
		// 构建树状结构：只返回顶级分类，它们的子分类在 SubCategories 字段中
		var topLevelCategories []model.Category
		categoryMap := make(map[uint]*model.Category)
		for i := range categories {
			categoryMap[categories[i].ID] = &categories[i]
		}

		for _, cat := range categories {
			if cat.ParentID == nil || *cat.ParentID == 0 {
				topLevelCategories = append(topLevelCategories, cat)
			} else {
				if parent, ok := categoryMap[*cat.ParentID]; ok {
					parent.SubCategories = append(parent.SubCategories, cat) // 这里直接修改了categories切片中元素的SubCategories
				}
			}
		}
		responses = mapCategoriesToResponse(topLevelCategories)
	} else {
		// 如果指定了parentID，则直接转换查询结果
		responses = mapCategoriesToResponse(categories)
	}

	return responses, nil
}

// GetAllCategoriesFlat 获取所有分类的扁平列表 (无层级结构，主要用于选择框等)
// userID: 当前操作的用户ID
// categoryType: 可选，按类型过滤 (income/expense)
func (s *BookkeepingCategoryService) GetAllCategoriesFlat(userID uint, categoryType model.CategoryType) ([]dto.CategoryResponse, error) {
	var categories []model.Category
	var responses []dto.CategoryResponse

	query := global.DB.Where("user_id = ?", userID).Order("type asc, sort_order asc, name asc")

	if categoryType != "" {
		query = query.Where("type = ?", categoryType)
	}

	if err := query.Find(&categories).Error; err != nil {
		global.Logger.Error("Failed to list all flat categories: " + err.Error())
		return nil, errors.New("获取分类列表失败")
	}

	if err := copier.Copy(&responses, &categories); err != nil {
		global.Logger.Error("Failed to copy categories to responses: " + err.Error())
		return nil, errors.New("数据转换失败")
	}
	for i := range responses {
		responses[i].CreatedAt = categories[i].CreatedAt.Format("2006-01-02 15:04:05")
		responses[i].UpdatedAt = categories[i].UpdatedAt.Format("2006-01-02 15:04:05")
	}

	return responses, nil
}
