package dto

type UserDTO struct {
	Username string `json:"username" binding:"required"`
	Age      int    `json:"age" binding:"required"`
	Password string `json:" " binding:"required"`
}
