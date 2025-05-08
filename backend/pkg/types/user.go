package types

type AddFormUser struct {
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname" binding:"required"`
	Password string `json:"password" binding:"password_strength, required"`
}

type UserInfo struct {
	Id       int64
	Email    string
	Nickname string
}

type UpdateFormUser struct {
	Id       string `json:"id" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname" binding:"required"`
	Password string `json:"password" binding:"password_strength, required"`
}


