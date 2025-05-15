package types

type ConfigForm struct {
	Name     string `json:"name"`
	Prompt   string `json:"prompt"`
	Delay    int64  `json:"delay"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
