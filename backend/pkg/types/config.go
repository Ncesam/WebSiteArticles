package types

type ConfigForm struct {
	Name     string `json:"name"`
	UserId   int64  `json:"user_id"`
	Prompt   string `json:"propmt"`
	Delay    int64  `json:"delay"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
