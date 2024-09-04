package dto

type UserInfo struct {
	Username string `json:"username" binding:"required"`
	Age      int    `json:"age" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (UserInfo) TableName() string {
	return "user_info"
}
