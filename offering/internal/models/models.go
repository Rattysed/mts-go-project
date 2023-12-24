package models

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Price struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

type Offer struct {
	FROM     Location `json:"from"`
	TO       Location `json:"to"`
	ClientID int      `json:"client_id"`
	Price    Price    `json:"price"`
}

type Answer struct {
	ID    string `json:"id"`
	Order Offer  `json:"order"`
}
