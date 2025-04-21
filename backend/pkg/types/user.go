package types

type AddFormUser struct {
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname" binding:"required"`
	Password string `json:"password" binding:"password_strength, required"`
	IsAdmin  bool   `json:"isAdmin"`
}

type UserInfo struct {
	Id       int32
	Email    string
	Nickname string
	IsAdmin  bool
}

type UpdateFormUser struct {
	Id       string `json:"id" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname" binding:"required"`
	Password string `json:"password" binding:"password_strength, required"`
	IsAdmin  bool   `json:"isAdmin"`
}
