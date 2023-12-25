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
	Id    string        `json:"id"`
	From  LatLngLiteral `json:"from"`
	To    LatLngLiteral `json:"to"`
	Price Money         `json:"price"`
}

type Event struct {
	Id              string            `json:"id"`
	Source          string            `json:"source"`
	Type            string            `json:"type"`
	DataContentType string            `json:"datacontenttype"`
	Time            string            `json:"time"`
	Data            map[string]string `json:"data"`
}

/*
{
    "id": "284655d6-0190-49e7-34e9-9b4060acc260",
    "source": "/client",
    "type": "trip.command.cancel",
    "datacontenttype": "application/json",
    "time": "2023-11-09T17:31:00Z",
    "data": {
        "trip_id": "284655d6-0190-49e7-34e9-9b4060acc260",
        "reason": "Водитель уехал в другую сторону"
    }
}
*/
