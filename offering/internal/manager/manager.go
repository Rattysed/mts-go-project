package manager

import (
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"math"
	"offering/internal/config"
	"offering/internal/models"
)

const hashMod int = 9859

type Manager struct {
	Cfg    *config.Config
	Logger *zap.Logger
}

func NewManager(cfg *config.Config, logger *zap.Logger) *Manager {
	return &Manager{
		Cfg:    cfg,
		Logger: logger,
	}
}

func GeneratePrice(from models.Location, to models.Location) int {
	x := int(math.Abs(from.Lat + from.Lng - to.Lng - to.Lat))
	return (x*31)%hashMod + 100
}

func (man *Manager) JwtPayloadFromRequest(tokenString string, secret string) (jwt.MapClaims, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		man.Logger.Warn(err.Error())
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		man.Logger.Warn("Invalid token")
	}

	return claims, true
}
