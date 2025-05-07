package service

import (
	"errors"
	"time"

	"github.com/dotdancer/gogofly/global"
	"github.com/dotdancer/gogofly/model"
	"github.com/dotdancer/gogofly/service/dto"
	"gorm.io/gorm"
)

// StatisticsService 统计服务
type StatisticsService struct{}

// GetTimeRange 根据传入的时间范围类型，计算开始和结束时间
func (s *StatisticsService) GetTimeRange(rangeType string, customStart, customEnd *time.Time) (time.Time, time.Time, error) {
	now := time.Now()
	var start, end time.Time

	// 如果传入了自定义时间范围，则优先使用
	if customStart != nil && customEnd != nil {
		return *customStart, *customEnd, nil
	}

	switch rangeType {
	case "day":
		// 今天的开始时间和结束时间
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end = start.Add(24 * time.Hour).Add(-time.Second)
	case "week":
		// 计算本周的开始（周一）和结束（周日）
		weekday := now.Weekday()
		if weekday == 0 { // 如果是周日
			weekday = 7
		}
		start = time.Date(now.Year(), now.Month(), now.Day()-int(weekday-1), 0, 0, 0, 0, now.Location())
		end = start.Add(7 * 24 * time.Hour).Add(-time.Second)
	case "month":
		// 本月的开始和结束
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location()).Add(-time.Second)
	case "year":
		// 本年的开始和结束
		start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		end = time.Date(now.Year()+1, 1, 1, 0, 0, 0, 0, now.Location()).Add(-time.Second)
	case "all":
		// 所有时间（使用一个很早的开始日期和当前时间）
		start = time.Date(1970, 1, 1, 0, 0, 0, 0, now.Location())
		end = now
	default:
		return time.Time{}, time.Time{}, errors.New("不支持的时间范围类型")
	}

	return start, end, nil
}

// GetIncomeExpenseSummary 获取指定时间范围内的收支汇总
func (s *StatisticsService) GetIncomeExpenseSummary(userID uint, rangeType string, customStart, customEnd *time.Time) (*dto.IncomeExpenseSummaryResponse, error) {
	start, end, err := s.GetTimeRange(rangeType, customStart, customEnd)
	if err != nil {
		return nil, err
	}

	var totalIncome, totalExpense float64

	// 计算总收入
	if err := global.DB.Model(&model.Transaction{}).
		Where("user_id = ? AND type = ? AND transaction_date BETWEEN ? AND ?", 
			userID, model.TransactionTypeIncome, start, end).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalIncome).Error; err != nil {
		return nil, err
	}

	// 计算总支出
	if err := global.DB.Model(&model.Transaction{}).
		Where("user_id = ? AND type = ? AND transaction_date BETWEEN ? AND ?", 
			userID, model.TransactionTypeExpense, start, end).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalExpense).Error; err != nil {
		return nil, err
	}

	return &dto.IncomeExpenseSummaryResponse{
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		NetAmount:    totalIncome - totalExpense,
		StartDate:    start,
		EndDate:      end,
		RangeType:    rangeType,
	}, nil
}

// GetCategorySummary 获取指定时间范围内的分类汇总
func (s *StatisticsService) GetCategorySummary(userID uint, transactionType model.TransactionType, rangeType string, customStart, customEnd *time.Time) ([]*dto.CategorySummaryItem, error) {
	start, end, err := s.GetTimeRange(rangeType, customStart, customEnd)
	if err != nil {
		return nil, err
	}

	var result []*dto.CategorySummaryItem

	err = global.DB.Model(&model.Transaction{}).
		Select("c.id as category_id, c.name as category_name, c.icon as category_icon, COALESCE(SUM(t.amount), 0) as total_amount, COUNT(t.id) as transaction_count").
		Joins("JOIN bookkeeping_categories c ON t.category_id = c.id").
		Where("t.user_id = ? AND t.type = ? AND t.transaction_date BETWEEN ? AND ?",
			userID, transactionType, start, end).
		Group("c.id, c.name, c.icon").
		Order("total_amount DESC").
		Table("bookkeeping_transactions t").
		Scan(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*dto.CategorySummaryItem{}, nil
		}
		return nil, err
	}

	return result, nil
}

// GetAccountSummary 获取账户余额汇总
func (s *StatisticsService) GetAccountSummary(userID uint) ([]*dto.AccountSummaryItem, error) {
	var accounts []*model.Account
	var result []*dto.AccountSummaryItem

	if err := global.DB.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*dto.AccountSummaryItem{}, nil
		}
		return nil, err
	}

	for _, account := range accounts {
		result = append(result, &dto.AccountSummaryItem{
			AccountID:      account.ID,
			AccountName:    account.Name,
			AccountType:    string(account.Type),
			CurrentBalance: account.CurrentBalance,
			InitialBalance: account.InitialBalance,
		})
	}

	return result, nil
}

// GetMonthlyTrend 获取月度收支趋势
func (s *StatisticsService) GetMonthlyTrend(userID uint, monthsCount int) (*dto.MonthlyTrendResponse, error) {
	if monthsCount <= 0 {
		monthsCount = 12 // 默认显示12个月
	}

	now := time.Now()
	endMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startMonth := endMonth.AddDate(0, -monthsCount+1, 0)

	var monthlyData []*dto.MonthlyData
	currentMonth := startMonth

	// 遍历每个月，查询数据
	for currentMonth.Before(endMonth) || currentMonth.Equal(endMonth) {
		nextMonth := currentMonth.AddDate(0, 1, 0)
		
		var monthlyIncome, monthlyExpense float64

		// 查询收入
		global.DB.Model(&model.Transaction{}).
			Where("user_id = ? AND type = ? AND transaction_date >= ? AND transaction_date < ?",
				userID, model.TransactionTypeIncome, currentMonth, nextMonth).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&monthlyIncome)

		// 查询支出
		global.DB.Model(&model.Transaction{}).
			Where("user_id = ? AND type = ? AND transaction_date >= ? AND transaction_date < ?",
				userID, model.TransactionTypeExpense, currentMonth, nextMonth).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&monthlyExpense)

		// 添加到结果中
		monthlyData = append(monthlyData, &dto.MonthlyData{
			Year:          currentMonth.Year(),
			Month:         int(currentMonth.Month()),
			MonthLabel:    currentMonth.Format("2006-01"),
			TotalIncome:   monthlyIncome,
			TotalExpense:  monthlyExpense,
			NetAmount:     monthlyIncome - monthlyExpense,
		})

		// 移动到下一个月
		currentMonth = nextMonth
	}

	return &dto.MonthlyTrendResponse{
		MonthsCount: monthsCount,
		Data:        monthlyData,
	}, nil
} 