package service

import (
	"github.com/dotdancer/gogofly/global"
	"github.com/dotdancer/gogofly/service/dto"
)

type UserInfoService struct {
}

func (u *UserInfoService) AddUser(user dto.UserInfo) error {
	return global.DB.Model(&dto.UserInfo{}).Create(&user).Error
}
