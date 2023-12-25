package models

type Status string

const (
	DRIVER_SEARCH Status = "DRIVER_SEARCH"
	DRIVER_FOUND         = "DRIVER_FOUND"
	ON_POSITION          = "ON_POSITION"
	STARTED              = "STARTED"
	ENDED                = "ENDED"
	CANCELED             = "CANCELED"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Price struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

type Trip struct {
	ID       string   `json:"id"`
	FROM     Location `json:"from"`
	TO       Location `json:"to"`
	ClientID string   `json:"client_id"`
	Price    Price    `json:"price"`
	Status   Status   `json:"status"`
}

type Offer struct {
	OfferID string `json:"offer_id"`
}

type Answer struct {
	ID    string                 `json:"id"`
	Order map[string]interface{} `json:"order"`
}

type Event struct {
	Id              string            `json:"id"`
	Source          string            `json:"source"`
	Type            string            `json:"type"`
	DataContentType string            `json:"datacontenttype"`
	Time            string            `json:"time"`
	Data            map[string]string `json:"data"`
}
