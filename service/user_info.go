package service

import (
	"errors"
	"github.com/dotdancer/gogofly/global"
	"github.com/dotdancer/gogofly/model"
	"github.com/dotdancer/gogofly/service/dto"
	"github.com/dotdancer/gogofly/utils"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type UserInfoService struct {
}

// AddUser 添加用户
func (u *UserInfoService) AddUser(user *dto.AddUserDto) error {

	var userInfo model.UserInfo
	if err := copier.Copy(&userInfo, user); err != nil {
		return errors.New("adding user field mapping failed")
	}

	return global.DB.Model(&model.UserInfo{}).Create(&userInfo).Error
}

// GetUserById 获取用户详情
func (u *UserInfoService) GetUserById(id string) (user *model.UserInfo, err error) {

	err = global.DB.Model(&model.UserInfo{}).Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return user, nil
}

// Login 获取登录
func (u *UserInfoService) Login(iUser *dto.LoginUserDto) (user model.UserInfo, token string, err error) {
	err = global.DB.Model(&model.UserInfo{}).Where("username = ?", iUser.Username).First(&user).Error
	//用户名或密码错误
	if errors.Is(err, gorm.ErrRecordNotFound) || !utils.CompareHashAndPassword(user.Password, iUser.Password) {
		return user, "", errors.New("invalid username or password")
	} else { //登录成功，生成token
		jwt := utils.NewJWT()
		token, err = jwt.GenerateToken(int(user.ID), user.Username)
		if err != nil {
			return user, token, err
		}
	}
	return user, token, nil
}
