package model

type Item struct {
	ID         int      `json:"id"`
	CategoryID int      `json:"category_id"`
	Category   Category `json:"category"`
	Price      float64  `json:"price"`
	Remaining  int      `json:"remaining"`
	Sold       int      `json:"sold"`
	PhotoURL   string   `json:"photo_url"`
}
