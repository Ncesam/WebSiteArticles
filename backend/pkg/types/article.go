package types

type Article struct {
	Title  string
	Blocks []Block
}

type InputForm struct {
	ConfigId string
	UserId   int64
	Prompt   string
	Data     string
}
type InputRequestForm struct {
	ConfigId string
	UserId   int64
	Prompt   string
	File     []byte
}
type EntryRequest struct {
	ID                 string `json:"id"`
	UserID             string `json:"user_id"`
	SubsiteID          string `json:"subsite_id"`
	Type               int64  `json:"type"`
	Title              string `json:"title"`
	Entry              Entry  `json:"entry"`
	ExternalAccessLink string `json:"external_access_link"`
	Path               string `json:"path"`
	IsEditorial        bool   `json:"is_editorial"`
	IsAdvertisement    bool   `json:"is_advertisement"`
	IsEnabledComments  bool   `json:"is_enabled_comments"`
	IsEnabledLikes     bool   `json:"is_enabled_likes"`
	Withheld           bool   `json:"withheld"`
	IsEnabledAd        bool   `json:"is_enabled_ad"`
	IsHoldOnFlash      bool   `json:"is_holdonflash"`
	ForcedToMainpage   int    `json:"forced_to_mainpage"`
	IsHoldOnMain       bool   `json:"is_holdonmain"`
	IsPublished        bool   `json:"is_published"`
	IsAdult            bool   `json:"is_adult"`
	RepostID           *int64 `json:"repostId"`   // может быть null
	RepostData         any    `json:"repostData"` // может быть null или объект
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
	Items []BlockItem `json:"items,omitempty"`
	Type  string      `json:"type,omitempty"`

	// Для quote, text и др.
	Text     string `json:"text,omitempty"`
	Subline1 string `json:"subline1,omitempty"`

	// Для media
	Image *ImageItem `json:"image,omitempty"`

	// Для link
	Link *Link `json:"link,omitempty"`
}

type BlockItem struct {
	Title string     `json:"title"`
	Image *ImageItem `json:"image,omitempty"`
}

type ImageItem struct {
	Type string    `json:"type"`
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
	Type string   `json:"type"`
	Data LinkData `json:"data"`
}

type LinkData struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       any    `json:"image"`
	V           int    `json:"v"`
	Hostname    string `json:"hostname"`
}
