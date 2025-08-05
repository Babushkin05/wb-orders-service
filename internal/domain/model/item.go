package model

// Item - товар в заказе
type Item struct {
	ChrtID      int     `json:"chrt_id"`
	TrackNumber string  `json:"track_number"`
	Price       int     `json:"price"`
	Rid         string  `json:"rid"`
	Name        string  `json:"name"`
	Sale        float64 `json:"sale"`
	Size        string  `json:"size"`
	TotalPrice  int     `json:"total_price"`
	NmID        int     `json:"nm_id"`
	Brand       string  `json:"brand"`
	Status      int     `json:"status"`
}
