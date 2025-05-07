package service

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/dotdancer/gogofly/global"
	"github.com/dotdancer/gogofly/model"
	"github.com/dotdancer/gogofly/service/dto"
	"gorm.io/gorm"
)

// BookkeepingBudgetService 预算服务
type BookkeepingBudgetService struct{}

// CreateBudget 创建预算
func (s *BookkeepingBudgetService) CreateBudget(userID uint, req dto.CreateBudgetRequest) (*dto.BudgetResponse, error) {
	// 如果是分类预算，需要校验分类是否存在
	if req.Type == string(model.BudgetTypeCategory) {
		if req.CategoryID == nil {
			return nil, errors.New("分类预算必须指定分类ID")
		}
		var category model.Category
		if err := global.DB.Where("id = ? AND user_id = ?", *req.CategoryID, userID).First(&category).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("指定的分类不存在")
			}
			return nil, err
		}
	} else if req.CategoryID != nil {
		// 如果不是分类预算，但传入了分类ID，则置为nil
		req.CategoryID = nil
	}

	// 创建预算
	budget := model.Budget{
		UserID:      userID,
		Name:        req.Name,
		Type:        model.BudgetType(req.Type),
		Period:      model.BudgetPeriod(req.Period),
		Amount:      req.Amount,
		StartDate:   req.StartDate,
		CategoryID:  req.CategoryID,
		Description: req.Description,
	}

	// 设置默认值
	if req.NotifyRate != nil {
		budget.NotifyRate = *req.NotifyRate
	} else {
		budget.NotifyRate = 0.8 // 默认80%提醒
	}

	if req.IsActive != nil {
		budget.IsActive = *req.IsActive
	} else {
		budget.IsActive = true // 默认激活
	}

	// 保存到数据库
	if err := global.DB.Create(&budget).Error; err != nil {
		global.Logger.Error("Failed to create budget: " + err.Error())
		return nil, errors.New("创建预算失败：数据库错误")
	}

	// 转换为响应
	var response dto.BudgetResponse
	if err := s.budgetToResponse(&budget, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetBudget 获取预算详情
func (s *BookkeepingBudgetService) GetBudget(userID, budgetID uint) (*dto.BudgetResponse, error) {
	var budget model.Budget
	if err := global.DB.Preload("Category").Where("id = ? AND user_id = ?", budgetID, userID).First(&budget).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("预算不存在")
		}
		global.Logger.Error("Failed to get budget: " + err.Error())
		return nil, errors.New("获取预算失败：数据库错误")
	}

	var response dto.BudgetResponse
	if err := s.budgetToResponse(&budget, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetBudgetProgress 获取预算进度
func (s *BookkeepingBudgetService) GetBudgetProgress(userID, budgetID uint) (*dto.BudgetProgressResponse, error) {
	// 获取预算基本信息
	var budget model.Budget
	if err := global.DB.Preload("Category").Where("id = ? AND user_id = ?", budgetID, userID).First(&budget).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("预算不存在")
		}
		global.Logger.Error("Failed to get budget: " + err.Error())
		return nil, errors.New("获取预算失败：数据库错误")
	}

	// 计算当前预算周期
	currentPeriodStart, currentPeriodEnd, err := s.calculateCurrentPeriod(budget.StartDate, budget.Period)
	if err != nil {
		return nil, err
	}

	// 获取当前周期内的支出
	var spentAmount float64
	query := global.DB.Model(&model.Transaction{}).
		Where("user_id = ? AND type = ? AND transaction_date BETWEEN ? AND ?", 
			userID, model.TransactionTypeExpense, currentPeriodStart, currentPeriodEnd)

	// 如果是分类预算，则只统计该分类的支出
	if budget.Type == model.BudgetTypeCategory && budget.CategoryID != nil {
		query = query.Where("category_id = ?", *budget.CategoryID)
	}

	if err := query.Select("COALESCE(SUM(amount), 0)").Scan(&spentAmount).Error; err != nil {
		global.Logger.Error("Failed to calculate spent amount: " + err.Error())
		return nil, errors.New("计算预算进度失败：数据库错误")
	}

	// 计算进度
	remainingAmount := budget.Amount - spentAmount
	if remainingAmount < 0 {
		remainingAmount = 0
	}

	usageRate := spentAmount / budget.Amount
	isOverBudget := usageRate > 1.0

	// 计算剩余天数
	daysRemaining := int(math.Ceil(currentPeriodEnd.Sub(time.Now()).Hours() / 24))
	if daysRemaining < 0 {
		daysRemaining = 0
	}

	// 转换为响应
	var response dto.BudgetProgressResponse
	response.ID = budget.ID
	response.UserID = budget.UserID
	response.Name = budget.Name
	response.Type = string(budget.Type)
	response.Period = string(budget.Period)
	response.Amount = budget.Amount
	response.StartDate = budget.StartDate
	response.CategoryID = budget.CategoryID
	response.NotifyRate = budget.NotifyRate
	response.Description = budget.Description
	response.IsActive = budget.IsActive
	response.CreatedAt = budget.CreatedAt
	response.UpdatedAt = budget.UpdatedAt
	response.Category = s.categoryToDTO(budget.Category)
	
	response.SpentAmount = spentAmount
	response.RemainingAmount = remainingAmount
	response.UsageRate = usageRate
	response.IsOverBudget = isOverBudget
	response.DaysRemaining = daysRemaining
	response.CurrentPeriod.StartDate = currentPeriodStart
	response.CurrentPeriod.EndDate = currentPeriodEnd

	return &response, nil
}

// UpdateBudget 更新预算
func (s *BookkeepingBudgetService) UpdateBudget(userID, budgetID uint, req dto.UpdateBudgetRequest) (*dto.BudgetResponse, error) {
	// 获取预算
	var budget model.Budget
	if err := global.DB.Where("id = ? AND user_id = ?", budgetID, userID).First(&budget).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("预算不存在")
		}
		global.Logger.Error("Failed to get budget for update: " + err.Error())
		return nil, errors.New("更新预算失败：数据库错误")
	}

	// 如果请求更新预算类型为分类预算，需要校验分类是否存在
	if req.Type != nil && *req.Type == string(model.BudgetTypeCategory) {
		if req.CategoryID == nil && budget.CategoryID == nil {
			return nil, errors.New("分类预算必须指定分类ID")
		}
		
		categoryID := budget.CategoryID
		if req.CategoryID != nil {
			categoryID = req.CategoryID
		}
		
		if categoryID != nil {
			var category model.Category
			if err := global.DB.Where("id = ? AND user_id = ?", *categoryID, userID).First(&category).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errors.New("指定的分类不存在")
				}
				return nil, err
			}
		}
	}

	// 更新预算字段
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}

	if req.Type != nil {
		updates["type"] = model.BudgetType(*req.Type)
		// 如果切换为总体预算，清除分类ID
		if *req.Type == string(model.BudgetTypeOverall) {
			updates["category_id"] = nil
		}
	}

	if req.Period != nil {
		updates["period"] = model.BudgetPeriod(*req.Period)
	}

	if req.Amount != nil {
		updates["amount"] = *req.Amount
	}

	if req.StartDate != nil {
		updates["start_date"] = *req.StartDate
	}

	if req.CategoryID != nil {
		// 只有当预算类型为分类预算时，才能设置分类ID
		if budget.Type == model.BudgetTypeCategory || (req.Type != nil && *req.Type == string(model.BudgetTypeCategory)) {
			updates["category_id"] = *req.CategoryID
		}
	}

	if req.NotifyRate != nil {
		updates["notify_rate"] = *req.NotifyRate
	}

	if req.Description != nil {
		updates["description"] = *req.Description
	}

	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	// 如果没有需要更新的字段，直接返回
	if len(updates) == 0 {
		return s.GetBudget(userID, budgetID)
	}

	// 更新数据库
	if err := global.DB.Model(&budget).Updates(updates).Error; err != nil {
		global.Logger.Error("Failed to update budget: " + err.Error())
		return nil, errors.New("更新预算失败：数据库错误")
	}

	// 获取更新后的预算
	return s.GetBudget(userID, budgetID)
}

// DeleteBudget 删除预算
func (s *BookkeepingBudgetService) DeleteBudget(userID, budgetID uint) error {
	// 检查预算是否存在
	var budget model.Budget
	if err := global.DB.Where("id = ? AND user_id = ?", budgetID, userID).First(&budget).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("预算不存在")
		}
		global.Logger.Error("Failed to get budget for deletion: " + err.Error())
		return errors.New("删除预算失败：数据库错误")
	}

	// 删除预算
	if err := global.DB.Delete(&budget).Error; err != nil {
		global.Logger.Error("Failed to delete budget: " + err.Error())
		return errors.New("删除预算失败：数据库错误")
	}

	return nil
}

// ListBudgets 获取预算列表
func (s *BookkeepingBudgetService) ListBudgets(userID uint, query dto.BudgetQuery) (dto.BudgetListResponse, error) {
	var budgets []model.Budget
	var response dto.BudgetListResponse

	// 构建查询条件
	db := global.DB.Model(&model.Budget{}).Where("user_id = ?", userID)

	// 应用筛选条件
	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}

	if query.Period != "" {
		db = db.Where("period = ?", query.Period)
	}

	if query.CategoryID > 0 {
		db = db.Where("category_id = ?", query.CategoryID)
	}

	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}

	// 计算总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		global.Logger.Error("Failed to count budgets: " + err.Error())
		return response, errors.New("获取预算列表失败：数据库错误")
	}

	// 分页查询
	offset := (query.Page - 1) * query.PageSize
	if err := db.Preload("Category").Order("id DESC").Offset(offset).Limit(query.PageSize).Find(&budgets).Error; err != nil {
		global.Logger.Error("Failed to list budgets: " + err.Error())
		return response, errors.New("获取预算列表失败：数据库错误")
	}

	// 构建响应
	response.Total = total
	response.Items = make([]dto.BudgetResponse, 0, len(budgets))

	for _, budget := range budgets {
		var budgetResponse dto.BudgetResponse
		if err := s.budgetToResponse(&budget, &budgetResponse); err != nil {
			global.Logger.Error(fmt.Sprintf("Failed to convert budget %d to response: %s", budget.ID, err.Error()))
			continue
		}

		response.Items = append(response.Items, budgetResponse)
	}

	return response, nil
}

// ListActiveBudgetProgress 获取激活的预算进度列表
func (s *BookkeepingBudgetService) ListActiveBudgetProgress(userID uint) ([]dto.BudgetProgressResponse, error) {
	// 查询所有激活的预算
	var budgets []model.Budget
	if err := global.DB.Preload("Category").Where("user_id = ? AND is_active = ?", userID, true).Find(&budgets).Error; err != nil {
		global.Logger.Error("Failed to list active budgets: " + err.Error())
		return nil, errors.New("获取预算进度列表失败：数据库错误")
	}

	// 构建响应
	var progressItems []dto.BudgetProgressResponse

	for _, budget := range budgets {
		// 计算当前预算周期
		currentPeriodStart, currentPeriodEnd, err := s.calculateCurrentPeriod(budget.StartDate, budget.Period)
		if err != nil {
			global.Logger.Error(fmt.Sprintf("Failed to calculate current period for budget %d: %s", budget.ID, err.Error()))
			continue
		}

		// 获取当前周期内的支出
		var spentAmount float64
		query := global.DB.Model(&model.Transaction{}).
			Where("user_id = ? AND type = ? AND transaction_date BETWEEN ? AND ?", 
				userID, model.TransactionTypeExpense, currentPeriodStart, currentPeriodEnd)

		// 如果是分类预算，则只统计该分类的支出
		if budget.Type == model.BudgetTypeCategory && budget.CategoryID != nil {
			query = query.Where("category_id = ?", *budget.CategoryID)
		}

		if err := query.Select("COALESCE(SUM(amount), 0)").Scan(&spentAmount).Error; err != nil {
			global.Logger.Error(fmt.Sprintf("Failed to calculate spent amount for budget %d: %s", budget.ID, err.Error()))
			continue
		}

		// 计算进度
		remainingAmount := budget.Amount - spentAmount
		if remainingAmount < 0 {
			remainingAmount = 0
		}

		usageRate := spentAmount / budget.Amount
		isOverBudget := usageRate > 1.0

		// 计算剩余天数
		daysRemaining := int(math.Ceil(currentPeriodEnd.Sub(time.Now()).Hours() / 24))
		if daysRemaining < 0 {
			daysRemaining = 0
		}

		// 构建进度项
		var progressItem dto.BudgetProgressResponse
		progressItem.ID = budget.ID
		progressItem.UserID = budget.UserID
		progressItem.Name = budget.Name
		progressItem.Type = string(budget.Type)
		progressItem.Period = string(budget.Period)
		progressItem.Amount = budget.Amount
		progressItem.StartDate = budget.StartDate
		progressItem.CategoryID = budget.CategoryID
		progressItem.NotifyRate = budget.NotifyRate
		progressItem.Description = budget.Description
		progressItem.IsActive = budget.IsActive
		progressItem.CreatedAt = budget.CreatedAt
		progressItem.UpdatedAt = budget.UpdatedAt
		progressItem.Category = s.categoryToDTO(budget.Category)
		
		progressItem.SpentAmount = spentAmount
		progressItem.RemainingAmount = remainingAmount
		progressItem.UsageRate = usageRate
		progressItem.IsOverBudget = isOverBudget
		progressItem.DaysRemaining = daysRemaining
		progressItem.CurrentPeriod.StartDate = currentPeriodStart
		progressItem.CurrentPeriod.EndDate = currentPeriodEnd

		progressItems = append(progressItems, progressItem)
	}

	return progressItems, nil
}

// CheckBudgetAlerts 检查预算提醒（达到或超过预算提醒阈值的预算）
func (s *BookkeepingBudgetService) CheckBudgetAlerts(userID uint) ([]dto.BudgetProgressResponse, error) {
	// 获取所有激活的预算进度
	progressItems, err := s.ListActiveBudgetProgress(userID)
	if err != nil {
		return nil, err
	}

	// 筛选出达到或超过阈值的预算
	var alerts []dto.BudgetProgressResponse
	for _, item := range progressItems {
		if item.UsageRate >= item.NotifyRate {
			alerts = append(alerts, item)
		}
	}

	return alerts, nil
}

// 辅助方法：预算模型转 DTO
func (s *BookkeepingBudgetService) budgetToResponse(budget *model.Budget, response *dto.BudgetResponse) error {
	if budget == nil || response == nil {
		return errors.New("budget or response is nil")
	}

	response.ID = budget.ID
	response.UserID = budget.UserID
	response.Name = budget.Name
	response.Type = string(budget.Type)
	response.Period = string(budget.Period)
	response.Amount = budget.Amount
	response.StartDate = budget.StartDate
	response.CategoryID = budget.CategoryID
	response.NotifyRate = budget.NotifyRate
	response.Description = budget.Description
	response.IsActive = budget.IsActive
	response.CreatedAt = budget.CreatedAt
	response.UpdatedAt = budget.UpdatedAt
	response.Category = s.categoryToDTO(budget.Category)

	return nil
}

// 辅助方法：Category 模型转 DTO
func (s *BookkeepingBudgetService) categoryToDTO(category *model.Category) *dto.Category {
	if category == nil {
		return nil
	}

	return &dto.Category{
		ID:       category.ID,
		Name:     category.Name,
		Type:     string(category.Type),
		ParentID: category.ParentID,
		Icon:     category.Icon,
	}
}

// 计算当前预算周期的开始和结束日期
func (s *BookkeepingBudgetService) calculateCurrentPeriod(startDate time.Time, period model.BudgetPeriod) (time.Time, time.Time, error) {
	now := time.Now()
	var currentPeriodStart, currentPeriodEnd time.Time

	switch period {
	case model.BudgetPeriodWeekly:
		// 找到当前周期的开始日期
		// 计算自 startDate 以来过了多少个完整的7天周期
		daysSinceStart := int(now.Sub(startDate).Hours() / 24)
		completedWeeks := daysSinceStart / 7
		
		// 当前周期的开始日期 = 初始日期 + 完整周数 * 7天
		currentPeriodStart = startDate.AddDate(0, 0, completedWeeks*7)
		
		// 如果当前周期已经过完，则加一个周期
		if now.After(currentPeriodStart.AddDate(0, 0, 7)) {
			currentPeriodStart = currentPeriodStart.AddDate(0, 0, 7)
		}
		
		// 结束日期就是开始日期 + 7天 - 1秒
		currentPeriodEnd = currentPeriodStart.AddDate(0, 0, 7).Add(-time.Second)

	case model.BudgetPeriodMonthly:
		// 找到当前周期的月份
		// 计算自 startDate 以来过了多少个完整月
		yearDiff := now.Year() - startDate.Year()
		monthDiff := int(now.Month()) - int(startDate.Month())
		totalMonths := yearDiff*12 + monthDiff
		
		// 当前周期的开始年月
		startYear := startDate.Year() + totalMonths/12
		startMonth := time.Month((int(startDate.Month()) + totalMonths%12) % 12)
		if startMonth == 0 {
			startMonth = 12
			startYear--
		}
		
		currentPeriodStart = time.Date(startYear, startMonth, startDate.Day(), 0, 0, 0, 0, startDate.Location())
		
		// 如果当前时间已经过了这个月的周期，则加一个月
		if now.After(currentPeriodStart.AddDate(0, 1, 0)) {
			if startMonth == 12 {
				startYear++
				startMonth = 1
			} else {
				startMonth++
			}
			currentPeriodStart = time.Date(startYear, startMonth, startDate.Day(), 0, 0, 0, 0, startDate.Location())
		}
		
		// 处理月底的情况，确保月份日期合法
		nextMonth := currentPeriodStart.Month() + 1
		nextYear := currentPeriodStart.Year()
		if nextMonth > 12 {
			nextMonth = 1
			nextYear++
		}
		
		// 结束日期是下个月同一天减1秒
		currentPeriodEnd = time.Date(nextYear, nextMonth, 1, 0, 0, 0, 0, startDate.Location()).AddDate(0, 0, -1).Add(24*time.Hour - time.Second)

	case model.BudgetPeriodYearly:
		// 找到当前周期的年份
		yearsSinceStart := now.Year() - startDate.Year()
		if now.Month() < startDate.Month() || (now.Month() == startDate.Month() && now.Day() < startDate.Day()) {
			yearsSinceStart--
		}
		
		currentPeriodStart = time.Date(startDate.Year()+yearsSinceStart, startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
		
		// 结束日期是开始日期加一年减1秒
		currentPeriodEnd = currentPeriodStart.AddDate(1, 0, 0).Add(-time.Second)

	default:
		return time.Time{}, time.Time{}, errors.New("不支持的预算周期类型")
	}

	return currentPeriodStart, currentPeriodEnd, nil
} 