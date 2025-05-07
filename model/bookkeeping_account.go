package model

import (
	"github.com/dotdancer/gogofly/global"
	"gorm.io/gorm"
)

// AccountType 定义账户类型
type AccountType string

const (
	AccountTypeCash       AccountType = "cash"        // 现金
	AccountTypeSavings    AccountType = "savings"     // 储蓄卡
	AccountTypeCreditCard AccountType = "credit_card" // 信用卡
	AccountTypeAlipay     AccountType = "alipay"      // 支付宝
	AccountTypeWechatPay  AccountType = "wechat_pay"  // 微信钱包
	AccountTypeInvestment AccountType = "investment"  // 投资账户
	AccountTypeOther      AccountType = "other"       // 其他
)

// Account 账户模型
type Account struct {
	global.GlyModel
	UserID         uint        `json:"user_id" gorm:"index;comment:用户ID"`
	Name           string      `json:"name" gorm:"type:varchar(100);not null;comment:账户名称"`
	Type           AccountType `json:"type" gorm:"type:varchar(50);not null;comment:账户类型"`
	InitialBalance float64     `json:"initial_balance" gorm:"type:decimal(10,2);default:0.00;comment:初始余额"`
	CurrentBalance float64     `json:"current_balance" gorm:"type:decimal(10,2);default:0.00;comment:当前余额"`
	Remark         string      `json:"remark" gorm:"type:varchar(255);comment:备注"`
	IsDefault      bool        `json:"is_default" gorm:"default:false;comment:是否默认账户"`
}

// TableName 指定表名
func (a *Account) TableName() string {
	return "bookkeeping_accounts"
}

// AfterCreate 钩子，在创建账户后，如果设置了初始余额，则将当前余额也设置为初始余额
func (a *Account) AfterCreate(tx *gorm.DB) (err error) {
	if a.InitialBalance != 0 && a.CurrentBalance == 0 {
		a.CurrentBalance = a.InitialBalance
		return tx.Save(a).Error
	}
	return nil
}
