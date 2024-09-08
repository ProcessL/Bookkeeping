package model

import (
	"github.com/dotdancer/gogofly/global"
	"github.com/dotdancer/gogofly/utils"
	"gorm.io/gorm"
)

type UserInfo struct {
	global.GlyModel
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Phone    string `json:"phone"`
}

func (u *UserInfo) TableName() string {
	return "user_info"
}

// GenerateFormPassword 方法用于将用户密码加密后生成表单密码。
// 该方法使用 bcrypt 算法对原始密码进行哈希处理，以增加密码的安全性。
// 方法返回两个值：加密后的密码和可能的错误。
func (u *UserInfo) GenerateFormPassword() (password string, err error) {
	hashPassword, err := utils.CryptoPassword(u.Password)
	if err != nil {
		return hashPassword, err
	}
	return hashPassword, nil
}

// BeforeCreate 在用户信息创建前的钩子
// 该方法用于在用户信息被创建前，对用户密码进行哈希处理
// 参数:
//   - tx: *gorm.DB, 数据库事务对象，允许与数据库进行交互
//
// 返回值:
//   - error: 如果在生成哈希密码过程中出现错误，会返回该错误
//
// 说明:
//   - 该方法首先通过调用 GenerateFormPassword 方法生成哈希密码
//   - 如果生成过程中出现错误，则直接返回错误
//   - 如果生成成功，则将生成的哈希密码赋值给用户信息的 Password 字段
func (u *UserInfo) BeforeCreate(tx *gorm.DB) (err error) {
	hashPassword, err := u.GenerateFormPassword()
	if err != nil {
		return err
	}
	u.Password = hashPassword
	return
}
