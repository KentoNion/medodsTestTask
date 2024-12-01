package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"medodsTest/auth/pkg"
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
		"exp":     cl.Now().Add(time.Hour * 24).Unix(), //интерфейс получения времени + 24 часа
	}
}

func (t Token) MapToRefresh(cl pkg.Clock) jwt.Claims {
	return jwt.MapClaims{
		"user_id": t.UserID,
		"secret":  t.Secret,
		"ip":      t.IP,
		"exp":     cl.Now().AddDate(1, 0, 0), //интерфейс получения времени + 1 год
	}
}

func (t Token) Fill(claims jwt.MapClaims) error {
	var ok bool
	t.IP, ok = claims["ip"].(string)
	if !ok {
		return errors.New("failed to parse ip")
	}
	t.Secret, ok = claims["secret"].(string)
	if !ok {
		return errors.New("failed to parse secret")
	}
	t.UserID, ok = claims["user_id"].(string)
	if !ok {
		return errors.New("failed to parse user_id")
	}
	return nil
}

var ErrWrongToken = errors.New("Wrong token")
