package data

type Clothe struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	Price    int64    `json:"price"`
	Brand    string   `json:"brand"`
	Color    string   `json:"color"`
	Sizes    []string `json:"sizes"`
	Sex      string   `json:"sex,omitempty"`
	Type     string   `json:"type,omitempty"`
	ImageURL string   `json:"image_url,omitempty"`
}
