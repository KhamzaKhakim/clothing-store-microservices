package data

type Brand struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Country     string `json:"country"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url,omitempty"`
}
