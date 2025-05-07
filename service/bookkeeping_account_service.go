package service

import (
	"errors"

	"github.com/dotdancer/gogofly/global"
	"github.com/dotdancer/gogofly/model"
	"github.com/dotdancer/gogofly/service/dto"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// BookkeepingAccountService 结构体定义了账户管理的服务层
type BookkeepingAccountService struct{}

// CreateAccount 创建一个新的账户
// userID: 当前操作的用户ID
// req: 创建账户的请求数据
func (s *BookkeepingAccountService) CreateAccount(userID uint, req dto.CreateAccountRequest) (dto.AccountResponse, error) {
	var account model.Account
	var response dto.AccountResponse

	// 检查同名账户是否已存在 (同一用户下)
	var existingAccount model.Account
	if err := global.DB.Where("user_id = ? AND name = ?", userID, req.Name).First(&existingAccount).Error; err == nil {
		return response, errors.New("该账户名称已存在")
	}

	// 复制请求数据到模型
	if err := copier.Copy(&account, &req); err != nil {
		global.Logger.Error("Failed to copy CreateAccountRequest to model.Account: " + err.Error())
		return response, errors.New("创建账户失败：数据复制错误")
	}

	account.UserID = userID

	// 如果设置为默认账户，需要将其他账户的默认标志设为false
	if account.IsDefault {
		if err := global.DB.Model(&model.Account{}).Where("user_id = ? AND is_default = ?", userID, true).Update("is_default", false).Error; err != nil {
			global.Logger.Error("Failed to update other accounts' default flag: " + err.Error())
			return response, errors.New("设置默认账户失败")
		}
	}

	// 创建账户
	if err := global.DB.Create(&account).Error; err != nil {
		global.Logger.Error("Failed to create account: " + err.Error())
		return response, errors.New("创建账户失败：数据库错误")
	}

	// 复制模型数据到响应
	if err := copier.Copy(&response, &account); err != nil {
		global.Logger.Error("Failed to copy model.Account to AccountResponse: " + err.Error())
		return response, errors.New("创建账户成功，但返回数据错误")
	}

	// 格式化时间
	response.CreatedAt = account.CreatedAt.Format("2006-01-02 15:04:05")
	response.UpdatedAt = account.UpdatedAt.Format("2006-01-02 15:04:05")

	return response, nil
}

// ListAccounts 获取用户的所有账户
// userID: 当前操作的用户ID
func (s *BookkeepingAccountService) ListAccounts(userID uint) ([]dto.AccountResponse, error) {
	var accounts []model.Account
	var response []dto.AccountResponse

	// 查询用户的所有账户
	if err := global.DB.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		global.Logger.Error("Failed to list accounts: " + err.Error())
		return nil, errors.New("获取账户列表失败：数据库错误")
	}

	// 复制模型数据到响应
	for _, account := range accounts {
		var accountResponse dto.AccountResponse
		if err := copier.Copy(&accountResponse, &account); err != nil {
			global.Logger.Error("Failed to copy model.Account to AccountResponse: " + err.Error())
			continue
		}

		// 格式化时间
		accountResponse.CreatedAt = account.CreatedAt.Format("2006-01-02 15:04:05")
		accountResponse.UpdatedAt = account.UpdatedAt.Format("2006-01-02 15:04:05")

		response = append(response, accountResponse)
	}

	return response, nil
}

// GetAccount 获取单个账户信息
// userID: 当前操作的用户ID
// accountID: 要获取的账户ID
func (s *BookkeepingAccountService) GetAccount(userID uint, accountID uint) (dto.AccountResponse, error) {
	var account model.Account
	var response dto.AccountResponse

	// 查询账户
	if err := global.DB.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, errors.New("账户不存在或不属于您")
		}
		global.Logger.Error("Failed to get account: " + err.Error())
		return response, errors.New("获取账户信息失败：数据库错误")
	}

	// 复制模型数据到响应
	if err := copier.Copy(&response, &account); err != nil {
		global.Logger.Error("Failed to copy model.Account to AccountResponse: " + err.Error())
		return response, errors.New("获取账户信息成功，但返回数据错误")
	}

	// 格式化时间
	response.CreatedAt = account.CreatedAt.Format("2006-01-02 15:04:05")
	response.UpdatedAt = account.UpdatedAt.Format("2006-01-02 15:04:05")

	return response, nil
}

// UpdateAccount 更新账户信息
// userID: 当前操作的用户ID
// accountID: 要更新的账户ID
// req: 更新账户的请求数据
func (s *BookkeepingAccountService) UpdateAccount(userID uint, accountID uint, req dto.UpdateAccountRequest) (dto.AccountResponse, error) {
	var account model.Account
	var response dto.AccountResponse

	// 查询账户
	if err := global.DB.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, errors.New("账户不存在或不属于您")
		}
		global.Logger.Error("Failed to get account for update: " + err.Error())
		return response, errors.New("更新账户失败：数据库错误")
	}

	// 如果更新名称，检查同名账户是否已存在
	if req.Name != nil && *req.Name != account.Name {
		var existingAccount model.Account
		if err := global.DB.Where("user_id = ? AND name = ? AND id != ?", userID, *req.Name, accountID).First(&existingAccount).Error; err == nil {
			return response, errors.New("该账户名称已存在")
		}
		account.Name = *req.Name
	}

	// 更新其他字段
	if req.Type != nil {
		account.Type = *req.Type
	}

	if req.Remark != nil {
		account.Remark = *req.Remark
	}

	// 如果设置为默认账户，需要将其他账户的默认标志设为false
	if req.IsDefault != nil {
		if *req.IsDefault && !account.IsDefault {
			if err := global.DB.Model(&model.Account{}).Where("user_id = ? AND is_default = ? AND id != ?", userID, true, accountID).Update("is_default", false).Error; err != nil {
				global.Logger.Error("Failed to update other accounts' default flag: " + err.Error())
				return response, errors.New("设置默认账户失败")
			}
		}
		account.IsDefault = *req.IsDefault
	}

	// 保存更新
	if err := global.DB.Save(&account).Error; err != nil {
		global.Logger.Error("Failed to update account: " + err.Error())
		return response, errors.New("更新账户失败：数据库错误")
	}

	// 复制模型数据到响应
	if err := copier.Copy(&response, &account); err != nil {
		global.Logger.Error("Failed to copy model.Account to AccountResponse: " + err.Error())
		return response, errors.New("更新账户成功，但返回数据错误")
	}

	// 格式化时间
	response.CreatedAt = account.CreatedAt.Format("2006-01-02 15:04:05")
	response.UpdatedAt = account.UpdatedAt.Format("2006-01-02 15:04:05")

	return response, nil
}

// DeleteAccount 删除账户
// userID: 当前操作的用户ID
// accountID: 要删除的账户ID
func (s *BookkeepingAccountService) DeleteAccount(userID uint, accountID uint) error {
	// 查询账户
	var account model.Account
	if err := global.DB.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("账户不存在或不属于您")
		}
		global.Logger.Error("Failed to get account for deletion: " + err.Error())
		return errors.New("删除账户失败：数据库错误")
	}

	// 检查账户是否有关联的交易记录
	var count int64
	if err := global.DB.Model(&model.Transaction{}).Where("account_id = ?", accountID).Count(&count).Error; err != nil {
		global.Logger.Error("Failed to count related transactions: " + err.Error())
		return errors.New("删除账户失败：无法检查关联交易记录")
	}

	if count > 0 {
		return errors.New("该账户存在关联的交易记录，无法删除")
	}

	// 删除账户
	if err := global.DB.Delete(&account).Error; err != nil {
		global.Logger.Error("Failed to delete account: " + err.Error())
		return errors.New("删除账户失败：数据库错误")
	}

	return nil
}
