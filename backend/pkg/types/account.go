package types


type AddAccountForm struct {
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname" binding:"required"`
	UserId string `json:"userId" binding:"required"`
}