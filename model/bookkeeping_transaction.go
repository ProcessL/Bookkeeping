package model

import (
	"time"

	"github.com/dotdancer/gogofly/global"
	"gorm.io/gorm"
)

// TransactionType 定义交易类型
type TransactionType string

const (
	TransactionTypeIncome   TransactionType = "income"   // 收入
	TransactionTypeExpense  TransactionType = "expense"  // 支出
	TransactionTypeTransfer TransactionType = "transfer" // 转账
)

// Transaction 交易流水模型
type Transaction struct {
	global.GlyModel
	UserID          uint            `json:"user_id" gorm:"index;comment:用户ID"`
	AccountID       uint            `json:"account_id" gorm:"index;comment:账户ID"`
	Type            TransactionType `json:"type" gorm:"type:varchar(50);not null;comment:交易类型 (income, expense, transfer)"`
	Amount          float64         `json:"amount" gorm:"type:decimal(10,2);not null;comment:金额"`
	TransactionDate time.Time       `json:"transaction_date" gorm:"not null;comment:交易日期"`
	CategoryID      uint            `json:"category_id" gorm:"index;comment:分类ID"`
	PayeePayer      string          `json:"payee_payer" gorm:"type:varchar(100);comment:收款方/付款方"`
	Notes           string          `json:"notes" gorm:"type:varchar(255);comment:备注"`

	// Associations
	Account  Account  `json:"account" gorm:"foreignKey:AccountID"`
	Category Category `json:"category" gorm:"foreignKey:CategoryID"` // Category model will be created next
}

// TableName 指定表名
func (t *Transaction) TableName() string {
	return "bookkeeping_transactions"
}

// AfterSave 钩子，在保存交易（创建或更新）后，更新关联账户的余额
func (t *Transaction) AfterSave(tx *gorm.DB) (err error) {
	return t.UpdateAccountBalance(tx)
}

// AfterDelete 钩子，在删除交易后，反向更新关联账户的余额
func (t *Transaction) AfterDelete(tx *gorm.DB) (err error) {
	// 当删除一笔交易时，需要将金额反向操作以恢复账户余额
	// 例如，删除一笔支出，则账户余额增加；删除一笔收入，则账户余额减少
	originalAmount := t.Amount
	switch t.Type {
	case TransactionTypeExpense:
		t.Amount = -originalAmount // 支出变负数，相当于增加余额
	case TransactionTypeIncome:
		t.Amount = originalAmount // 收入变正数，相当于减少余额
	case TransactionTypeTransfer:
		// 对于转账，删除时需要更复杂的逻辑，取决于这是转出还是转入
		// 简单处理：假设转账记录的是转出，则删除时余额增加
		// 如果有明确的转出账户和转入账户，则需要分别处理
		// 此处简化为，如果删除的是转账，则认为是从该账户转出，所以余额增加
		// 如果需要更精确，Transaction模型应包含 FromAccountID 和 ToAccountID
		t.Amount = -originalAmount
	default:
		return nil // 其他类型不处理
	}

	err = t.UpdateAccountBalance(tx)
	// 恢复原始金额，以免影响后续可能的其他操作
	t.Amount = originalAmount
	return err
}

// UpdateAccountBalance 更新账户余额的辅助函数
func (t *Transaction) UpdateAccountBalance(tx *gorm.DB) error {
	var account Account
	if err := tx.First(&account, t.AccountID).Error; err != nil {
		return err
	}

	var totalIncome float64
	var totalExpense float64
	var totalTransferOut float64 // 从该账户转出
	var totalTransferIn float64  // 转入该账户 (如果Transaction模型支持ToAccountID)

	tx.Model(&Transaction{}).Where("account_id = ? AND type = ?", t.AccountID, TransactionTypeIncome).Select("COALESCE(SUM(amount), 0)").Row().Scan(&totalIncome)
	tx.Model(&Transaction{}).Where("account_id = ? AND type = ?", t.AccountID, TransactionTypeExpense).Select("COALESCE(SUM(amount), 0)").Row().Scan(&totalExpense)
	// 假设当前Transaction模型中的AccountID代表的是转出账户
	tx.Model(&Transaction{}).Where("account_id = ? AND type = ?", t.AccountID, TransactionTypeTransfer).Select("COALESCE(SUM(amount), 0)").Row().Scan(&totalTransferOut)
	// 如果Transaction模型有ToAccountID, 则需要查询转入该账户的金额
	// tx.Model(&Transaction{}).Where("to_account_id = ? AND type = ?", t.AccountID, TransactionTypeTransfer).Select("COALESCE(SUM(amount), 0)").Row().Scan(&totalTransferIn)

	account.CurrentBalance = account.InitialBalance + totalIncome - totalExpense - totalTransferOut + totalTransferIn

	return tx.Save(&account).Error
}
