package types

type Article struct {
	Title string
	Blocks []Block
}
  

type InputForm struct {
	ConfigId string
	UserId int64
	Prompt string
	Data string
}
type EntryRequest struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	SubsiteID int    `json:"subsite_id"`
	Type      int    `json:"type"`
	Title     string `json:"title"`
	Entry     Entry  `json:"entry"`
}

type Entry struct {
	Blocks []Block `json:"blocks"`
}

type Block struct {
	Type   string    `json:"type"`
	Cover  bool      `json:"cover"`
	Hidden bool      `json:"hidden"`
	Anchor string    `json:"anchor"`
	Data   BlockData `json:"data"`
}

type BlockData struct {
	// Для list
	Items []string `json:"items,omitempty"`
	Type  string   `json:"type,omitempty"` // UL / OL

	// Для quote, text и др.
	Text     string `json:"text,omitempty"`
	Subline1 string `json:"subline1,omitempty"`

	// Для media
	Image *ImageItem `json:"image,omitempty"`

	// Для link
	Link *Link `json:"link,omitempty"`
}

type ImageItem struct {
	Type string   `json:"type"`
	Data ImageData `json:"data"`
}

type ImageData struct {
	UUID          string   `json:"uuid"`
	Width         int      `json:"width"`
	Height        int      `json:"height"`
	Size          int      `json:"size"`
	Type          string   `json:"type"`
	Color         string   `json:"color"`
	Hash          string   `json:"hash"`
	ExternalLinks []string `json:"external_service"`
	Base64Preview string   `json:"base64preview"`
	IsVideo       bool     `json:"isVideo"`
	Duration      *int     `json:"duration"`
	HasAudio      bool     `json:"has_audio"`
}

type Link struct {
	Type string      `json:"type"`
	Data LinkData    `json:"data"`
}

type LinkData struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       any    `json:"image"`     // Может быть nil или объект
	V           int    `json:"v"`
	Hostname    string `json:"hostname"`
}

