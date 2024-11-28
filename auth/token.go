package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Token struct {
	UserID string
	Secret string
	IP     string
}

func (t Token) MapToAcces(cl pkg.Clock) jwt.Claims {
	return jwt.MapClaims{
		"user_id": t.UserID,
		"secret":  t.Secret,
		"ip":      t.IP,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
}

func (t Token) MapToRefresh(cl pkg.Clock) jwt.Claims {
	return jwt.MapClaims{
		"user_id": t.UserID,
		"secret":  t.Secret,
		"ip":      t.IP,
		"exp":     time.Now().AddDate(1, 0, 0),
	}
}
