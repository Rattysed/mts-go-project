package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
	"offering/internal/manager"
	"offering/internal/models"
)

type Controller struct {
	manager *manager.Manager
	Logger  *zap.Logger
}

func NewController(man *manager.Manager, logger *zap.Logger) *Controller {
	return &Controller{
		manager: man,
		Logger:  logger,
	}
}

func (c *Controller) ParseOffer(w http.ResponseWriter, r *http.Request) {
	offId := chi.URLParam(r, "offerID")

	payload, ok := c.manager.JwtPayloadFromRequest(offId, c.manager.Cfg.JWT)
	if !ok {
		c.Logger.Fatal("Bad token")
	}

	orderJson, ok := payload["order"].(string)
	if !ok {
		c.Logger.Fatal("incorrect data in token")
	}

	order := models.Offer{}
	err := json.Unmarshal([]byte(orderJson), &order)
	if err != nil {
		c.Logger.Fatal(err.Error())
	}

	resp := models.Answer{
		ID:    offId,
		Order: order,
	}

	bytes, err := json.Marshal(resp)

	w.Write(bytes)
}

func (c *Controller) CreateOffer(w http.ResponseWriter, r *http.Request) {
	bytesBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Плохое тело запроса"))
		c.Logger.Fatal(err.Error())
	}

	var jwtStruct models.Offer
	err = json.Unmarshal(bytesBody, &jwtStruct)
	if err != nil {
		c.Logger.Fatal(err.Error())
	}

	jwtStruct.Price.Amount = manager.GeneratePrice(jwtStruct.FROM, jwtStruct.TO)
	jwtStruct.Price.Currency = "$"

	bytes, err := json.Marshal(jwtStruct)
	if err != nil {
		c.Logger.Fatal(err.Error())
	}

	payload := jwt.MapClaims{
		"order": string(bytes),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	jwtKey := c.manager.Cfg.JWT
	t, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		c.Logger.Fatal(err.Error())
	}

	responseStruct := models.Answer{
		ID:    t,
		Order: jwtStruct,
	}

	resp, err := json.Marshal(responseStruct)
	if err != nil {
		c.Logger.Fatal(err.Error())
	}

	w.Write(resp)
}
