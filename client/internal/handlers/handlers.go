package handlers

import (
	"client/internal/admin"
	"client/internal/models"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

type Controller struct {
	DBController *admin.DBController
}

func NewController(DBController *admin.DBController) *Controller {
	return &Controller{
		DBController: DBController,
	}
}

func (c *Controller) GetTrip(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("user_id")
	if userID == "" {
		c.DBController.Logger.Warn("No user in header")
		return
	}

	tripID := chi.URLParam(r, "trip_id")
	if tripID == "" {
		c.DBController.Logger.Warn("No trip_id in params")
		return
	}

	trip := *c.DBController.GetTrip(tripID, userID)
	fmt.Println(trip)
	tripJSON, err := json.Marshal(trip)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		c.DBController.Logger.Warn("failed to make a json from trip")
		return
	}

	w.Write(tripJSON)
}

func (c *Controller) AddTrip(w http.ResponseWriter, r *http.Request) {
	bytesBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Плохое тело запроса"))
		w.WriteHeader(http.StatusBadRequest)
		c.DBController.Logger.Warn(err.Error())
	}

	userID := r.Header.Get("user_id")
	if userID == "" {
		c.DBController.Logger.Warn("No user in header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var offerID models.Offer
	err = json.Unmarshal(bytesBody, &offerID)
	if err != nil {
		c.DBController.Logger.Warn(err.Error())
		return
	}

	resp, err := http.Get("http://localhost:63343/offers/" + offerID.OfferID)
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response body", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		c.DBController.Logger.Warn(err.Error())
		return
	}

	var decodedOrder models.Answer
	err = json.Unmarshal(bytes, &decodedOrder)
	if err != nil {
		c.DBController.Logger.Warn(err.Error())
		return
	}

	var trip models.Trip
	trip.ID = decodedOrder.ID
	for key, value := range decodedOrder.Order {
		switch value.(type) {
		case string:
			if key == "client_id" {
				trip.ClientID = value.(string)
			}
		case map[string]interface{}:
			switch key {
			case "from":
				fmt.Println(value.(map[string]interface{})["lat"].(float64))
				trip.FROM.Lat = value.(map[string]interface{})["lat"].(float64)
				trip.FROM.Lng = value.(map[string]interface{})["lng"].(float64)
			case "to":
				trip.TO.Lat = value.(map[string]interface{})["lat"].(float64)
				trip.TO.Lng = value.(map[string]interface{})["lng"].(float64)
			case "price":
				trip.Price.Amount = int(value.(map[string]interface{})["amount"].(float64))
				trip.Price.Currency = value.(map[string]interface{})["currency"].(string)
			}
		}
	}

	fmt.Println(trip)

	c.DBController.AddTrip(trip)
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) ListTrips(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("user_id")
	if userID == "" {
		c.DBController.Logger.Warn("No user in header")
		return
	}

	trips := c.DBController.ListTrips(userID)
	tripsJSON, err := json.Marshal(trips)
	if err != nil {
		http.Error(w, "Marshaling error", http.StatusInternalServerError)
		c.DBController.Logger.Warn("Marshaling error")
		return
	}

	_, err = w.Write(tripsJSON)
	if err != nil {
		http.Error(w, "Writing response error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) CancelTrip(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user ID in header", http.StatusBadRequest)
		c.DBController.Logger.Warn("No user in header")
		return
	}

	tripID := chi.URLParam(r, "trip_id")
	if tripID == "" {
		http.Error(w, "Missing trip ID in path parameters", http.StatusBadRequest)
		c.DBController.Logger.Warn("No trip_id in params")
		return
	}

	reason := r.URL.Query().Get("reason")
	if reason == "" {
		http.Error(w, "Missing reason in query parameters", http.StatusBadRequest)
		c.DBController.Logger.Warn("No reason in params")
		return
	}

	err := c.DBController.ChangeStatus(tripID, userID, models.CANCELED)
	if err != nil {
		http.Error(w, "Failed to update data", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
