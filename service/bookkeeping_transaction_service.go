package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/dotdancer/gogofly/global"
	"github.com/dotdancer/gogofly/model"
	"github.com/dotdancer/gogofly/service/dto"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// BookkeepingTransactionService 结构体定义了交易流水管理的服务层
type BookkeepingTransactionService struct{}

// CreateTransaction 创建一个新的交易流水
// userID: 当前操作的用户ID
// req: 创建交易流水的请求数据
func (s *BookkeepingTransactionService) CreateTransaction(userID uint, req dto.CreateTransactionRequest) (dto.TransactionResponse, error) {
	var transaction model.Transaction
	var response dto.TransactionResponse

	// 验证账户是否存在且属于当前用户
	var account model.Account
	if err := global.DB.Where("id = ? AND user_id = ?", req.AccountID, userID).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, errors.New("账户不存在或不属于您")
		}
		global.Logger.Error("Failed to find account: " + err.Error())
		return response, errors.New("创建交易记录失败：无法验证账户")
	}

	// 验证分类是否存在且属于当前用户
	var category model.Category
	if err := global.DB.Where("id = ? AND user_id = ?", req.CategoryID, userID).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, errors.New("分类不存在或不属于您")
		}
		global.Logger.Error("Failed to find category: " + err.Error())
		return response, errors.New("创建交易记录失败：无法验证分类")
	}

	// 解析交易日期
	transactionDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		global.Logger.Error("Failed to parse transaction date: " + err.Error())
		return response, errors.New("交易日期格式错误，请使用YYYY-MM-DD格式")
	}

	// 复制请求数据到模型
	if err := copier.Copy(&transaction, &req); err != nil {
		global.Logger.Error("Failed to copy CreateTransactionRequest to model.Transaction: " + err.Error())
		return response, errors.New("创建交易记录失败：数据复制错误")
	}

	transaction.UserID = userID
	transaction.TransactionDate = transactionDate

	// 创建交易记录（在事务中进行，确保账户余额更新）
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		global.Logger.Error("Failed to create transaction: " + err.Error())
		return response, errors.New("创建交易记录失败：数据库错误")
	}

	// 重新查询交易记录（包含关联信息）
	if err := global.DB.Preload("Account").Preload("Category").First(&transaction, transaction.ID).Error; err != nil {
		global.Logger.Error("Failed to reload transaction: " + err.Error())
		return response, errors.New("创建交易记录成功，但获取详情失败")
	}

	// 复制模型数据到响应
	if err := s.transactionToResponse(&transaction, &response); err != nil {
		return response, err
	}

	return response, nil
}

// ListTransactions 获取交易流水列表
// userID: 当前操作的用户ID
// query: 查询条件
func (s *BookkeepingTransactionService) ListTransactions(userID uint, query dto.TransactionQuery) (dto.TransactionListResponse, error) {
	var transactions []model.Transaction
	var response dto.TransactionListResponse

	// 构建查询条件
	db := global.DB.Model(&model.Transaction{}).Where("user_id = ?", userID)

	// 应用筛选条件
	if query.AccountID > 0 {
		db = db.Where("account_id = ?", query.AccountID)
	}

	if query.CategoryID > 0 {
		db = db.Where("category_id = ?", query.CategoryID)
	}

	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}

	if query.StartDate != "" {
		db = db.Where("transaction_date >= ?", query.StartDate)
	}

	if query.EndDate != "" {
		db = db.Where("transaction_date <= ?", query.EndDate)
	}

	// 计算总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		global.Logger.Error("Failed to count transactions: " + err.Error())
		return response, errors.New("获取交易流水列表失败：数据库错误")
	}

	// 分页查询
	offset := (query.Page - 1) * query.PageSize
	if err := db.Preload("Account").Preload("Category").Order("transaction_date DESC, id DESC").Offset(offset).Limit(query.PageSize).Find(&transactions).Error; err != nil {
		global.Logger.Error("Failed to list transactions: " + err.Error())
		return response, errors.New("获取交易流水列表失败：数据库错误")
	}

	// 构建响应
	response.Total = total
	response.Items = make([]dto.TransactionResponse, 0, len(transactions))

	for _, transaction := range transactions {
		var transactionResponse dto.TransactionResponse
		if err := s.transactionToResponse(&transaction, &transactionResponse); err != nil {
			global.Logger.Error(fmt.Sprintf("Failed to convert transaction %d to response: %s", transaction.ID, err.Error()))
			continue
		}

		response.Items = append(response.Items, transactionResponse)
	}

	return response, nil
}

// GetTransaction 获取单个交易流水信息
// userID: 当前操作的用户ID
// transactionID: 要获取的交易流水ID
func (s *BookkeepingTransactionService) GetTransaction(userID uint, transactionID uint) (dto.TransactionResponse, error) {
	var transaction model.Transaction
	var response dto.TransactionResponse

	// 查询交易流水
	if err := global.DB.Preload("Account").Preload("Category").Where("id = ? AND user_id = ?", transactionID, userID).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, errors.New("交易记录不存在或不属于您")
		}
		global.Logger.Error("Failed to get transaction: " + err.Error())
		return response, errors.New("获取交易记录失败：数据库错误")
	}

	// 复制模型数据到响应
	if err := s.transactionToResponse(&transaction, &response); err != nil {
		return response, err
	}

	return response, nil
}

// UpdateTransaction 更新交易流水信息
// userID: 当前操作的用户ID
// transactionID: 要更新的交易流水ID
// req: 更新交易流水的请求数据
func (s *BookkeepingTransactionService) UpdateTransaction(userID uint, transactionID uint, req dto.UpdateTransactionRequest) (dto.TransactionResponse, error) {
	var transaction model.Transaction
	var response dto.TransactionResponse

	// 查询交易流水
	if err := global.DB.Where("id = ? AND user_id = ?", transactionID, userID).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, errors.New("交易记录不存在或不属于您")
		}
		global.Logger.Error("Failed to get transaction for update: " + err.Error())
		return response, errors.New("更新交易记录失败：数据库错误")
	}

	// 验证账户（如果更新）
	if req.AccountID != nil && *req.AccountID != transaction.AccountID {
		var account model.Account
		if err := global.DB.Where("id = ? AND user_id = ?", *req.AccountID, userID).First(&account).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response, errors.New("账户不存在或不属于您")
			}
			global.Logger.Error("Failed to find account: " + err.Error())
			return response, errors.New("更新交易记录失败：无法验证账户")
		}
		transaction.AccountID = *req.AccountID
	}

	// 验证分类（如果更新）
	if req.CategoryID != nil && *req.CategoryID != transaction.CategoryID {
		var category model.Category
		if err := global.DB.Where("id = ? AND user_id = ?", *req.CategoryID, userID).First(&category).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response, errors.New("分类不存在或不属于您")
			}
			global.Logger.Error("Failed to find category: " + err.Error())
			return response, errors.New("更新交易记录失败：无法验证分类")
		}
		transaction.CategoryID = *req.CategoryID
	}

	// 更新其他字段
	if req.Type != nil {
		transaction.Type = *req.Type
	}

	if req.Amount != nil {
		transaction.Amount = *req.Amount
	}

	if req.TransactionDate != nil {
		transactionDate, err := time.Parse("2006-01-02", *req.TransactionDate)
		if err != nil {
			global.Logger.Error("Failed to parse transaction date: " + err.Error())
			return response, errors.New("交易日期格式错误，请使用YYYY-MM-DD格式")
		}
		transaction.TransactionDate = transactionDate
	}

	if req.PayeePayer != nil {
		transaction.PayeePayer = *req.PayeePayer
	}

	if req.Notes != nil {
		transaction.Notes = *req.Notes
	}

	// 保存更新（在事务中进行，确保账户余额更新）
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&transaction).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		global.Logger.Error("Failed to update transaction: " + err.Error())
		return response, errors.New("更新交易记录失败：数据库错误")
	}

	// 重新查询交易记录（包含关联信息）
	if err := global.DB.Preload("Account").Preload("Category").First(&transaction, transaction.ID).Error; err != nil {
		global.Logger.Error("Failed to reload transaction: " + err.Error())
		return response, errors.New("更新交易记录成功，但获取详情失败")
	}

	// 复制模型数据到响应
	if err := s.transactionToResponse(&transaction, &response); err != nil {
		return response, err
	}

	return response, nil
}

// DeleteTransaction 删除交易流水
// userID: 当前操作的用户ID
// transactionID: 要删除的交易流水ID
func (s *BookkeepingTransactionService) DeleteTransaction(userID uint, transactionID uint) error {
	// 查询交易流水
	var transaction model.Transaction
	if err := global.DB.Where("id = ? AND user_id = ?", transactionID, userID).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("交易记录不存在或不属于您")
		}
		global.Logger.Error("Failed to get transaction for deletion: " + err.Error())
		return errors.New("删除交易记录失败：数据库错误")
	}

	// 删除交易记录（在事务中进行，确保账户余额更新）
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&transaction).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		global.Logger.Error("Failed to delete transaction: " + err.Error())
		return errors.New("删除交易记录失败：数据库错误")
	}

	return nil
}

// transactionToResponse 辅助函数，将交易流水模型转换为响应对象
func (s *BookkeepingTransactionService) transactionToResponse(transaction *model.Transaction, response *dto.TransactionResponse) error {
	// 复制基本字段
	if err := copier.Copy(response, transaction); err != nil {
		global.Logger.Error("Failed to copy model.Transaction to TransactionResponse: " + err.Error())
		return errors.New("转换交易记录数据失败")
	}

	// 格式化日期和时间
	response.TransactionDate = transaction.TransactionDate.Format("2006-01-02")
	response.CreatedAt = transaction.CreatedAt.Format("2006-01-02 15:04:05")
	response.UpdatedAt = transaction.UpdatedAt.Format("2006-01-02 15:04:05")

	// 处理关联信息
	if transaction.Account.ID > 0 {
		var accountResponse dto.AccountResponse
		if err := copier.Copy(&accountResponse, &transaction.Account); err != nil {
			global.Logger.Error("Failed to copy Account to AccountResponse: " + err.Error())
		} else {
			accountResponse.CreatedAt = transaction.Account.CreatedAt.Format("2006-01-02 15:04:05")
			accountResponse.UpdatedAt = transaction.Account.UpdatedAt.Format("2006-01-02 15:04:05")
			response.Account = accountResponse
		}
	}

	if transaction.Category.ID > 0 {
		var categoryResponse dto.CategoryResponse
		if err := copier.Copy(&categoryResponse, &transaction.Category); err != nil {
			global.Logger.Error("Failed to copy Category to CategoryResponse: " + err.Error())
		} else {
			categoryResponse.CreatedAt = transaction.Category.CreatedAt.Format("2006-01-02 15:04:05")
			categoryResponse.UpdatedAt = transaction.Category.UpdatedAt.Format("2006-01-02 15:04:05")
			response.Category = categoryResponse
		}
	}

	return nil
}
