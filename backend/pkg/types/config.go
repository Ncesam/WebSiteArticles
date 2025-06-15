package types

type ConfigForm struct {
	Name        string  `json:"name"`
	Prompt      string  `json:"prompt"`
	Delay       int64   `json:"delay"`
	Email       string  `json:"email"`
	Password    string  `json:"password"`
	WebSite     WebSite `json:"website"`
	IsPublished bool    `json:"is_published"`
}

type WebSite string

const (
	DTF   WebSite = "DTF"
	VC_RU WebSite = "VC_RU"
)
