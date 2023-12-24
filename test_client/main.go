package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type request struct {
	FROM     Location `json:"from"`
	TO       Location `json:"to"`
	ClientID int      `json:"client_id"`
}

func createOffer() {
	var req = &request{
		FROM:     Location{4, 12},
		TO:       Location{8.3823, 9.2},
		ClientID: 1,
	}

	bytesRepresentation, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
	}

	response, err := http.Post("http://127.0.0.1:63343/offers", "application/json", bytes.NewBuffer(bytesRepresentation))
	defer response.Body.Close()

	if err != nil {
		fmt.Println(err)
	}
	bytesResp, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Ответ сервера:", string(bytesResp))
}

func main() {
	createOffer()
}
