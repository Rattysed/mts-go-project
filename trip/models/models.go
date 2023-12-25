package models

type LatLngLiteral struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Money struct {
	Amount   float32 `json:"amount"`
	Currency string  `json:"currency"`
}

type Offer struct {
	Id       string        `json:"id"`
	From     LatLngLiteral `json:"from"`
	To       LatLngLiteral `json:"to"`
	ClientId string        `json:"clientId"`
	Price    Money         `json:"price"`
}

type Event struct {
	Id              string            `json:"id"`
	Source          string            `json:"source"`
	Type            string            `json:"type"`
	DataContentType string            `json:"datacontenttype"`
	Time            string            `json:"time"`
	Data            map[string]string `json:"data"`
}

type Trip struct {
}
