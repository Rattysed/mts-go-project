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
	ClientID string   `json:"client_id"`
}

type OfferId struct {
	OfferID string `json:"offer_id"`
}

func createOffer() {
	var req = &request{
		FROM:     Location{282, 218218},
		TO:       Location{1, 1.732},
		ClientID: "HELLO",
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

func createTrip() {
	var req = &OfferId{
		OfferID: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJvcmRlciI6IntcImZyb21cIjp7XCJsYXRcIjo0LFwibG5nXCI6MTJ9LFwidG9cIjp7XCJsYXRcIjo4LjM4MjMsXCJsbmdcIjo5LjJ9LFwiY2xpZW50X2lkXCI6XCJzZFwiLFwicHJpY2VcIjp7XCJhbW91bnRcIjoxMzEsXCJjdXJyZW5jeVwiOlwiJFwifX0ifQ._XcYRQ2Zbh5q4G_VNV-BjF1XF04snUEvih6hYDohluM",
	}

	bytesRepresentation, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
	}

	response, err := http.Post("http://127.0.0.1:8080/trips", "application/json", bytes.NewBuffer(bytesRepresentation))
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
	//createTrip()
}
