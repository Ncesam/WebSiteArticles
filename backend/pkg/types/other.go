package types

type GetLinkResponse struct {
	Status  string  `json:"status"`
	Context Context `json:"context"`
	Link    Links   `json:"link"`
}

type Context struct {
	Region    Region   `json:"region"`
	Currency  Currency `json:"currency"`
	ID        string   `json:"id"`
	Time      string   `json:"time"` // Или time.Time, если будете парсить как время
	MarketURL string   `json:"marketUrl"`
}

type Region struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Country     Country `json:"country"`
	Type        string  `json:"type"`
	ChildCount  int     `json:"childCount"`
	Coordinates Coord   `json:"coordinates"`
}

type Country struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	ChildCount int    `json:"childCount"`
}

type Coord struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Currency struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Links struct {
	URL          string `json:"url"`
	ShortURL     string `json:"shortUrl"`
	SearchType   string `json:"searchType"`
	PageName     string `json:"pageName"`
	ProductPhoto string `json:"productPhoto"`
}

type UploadResponse struct {
	Message string       `json:"message"`
	Result  []UploadItem `json:"result"`
}

type UploadItem struct {
	Type string     `json:"type"`
	Data UploadData `json:"data"`
}

type UploadData struct {
	UUID            string        `json:"uuid"`
	Width           int           `json:"width"`
	Height          int           `json:"height"`
	Size            int           `json:"size"`
	Type            string        `json:"type"`
	Color           string        `json:"color"`
	Hash            string        `json:"hash"`
	ExternalService []interface{} `json:"external_service"` // или []string, если точно знаешь
	Base64Preview   string        `json:"base64preview"`
	IsVideo         bool          `json:"isVideo"`
	Duration        interface{}   `json:"duration"` // можно заменить на *float64, если знаешь, что это float
	HasAudio        bool          `json:"has_audio"`
}
