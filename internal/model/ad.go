package model

type Ad struct {
	ID          int     `json:"id,omitempty"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price"`
	Author      string  `json:"author"`
}
