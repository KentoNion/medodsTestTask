package auth

import (
	"context"
	"github.com/pkg/errors"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type authStore interface {
	Save()
	Get()
}

type notifier interface {
	NotifyNewLogin()
}

type service struct {
	store     authStore
	secretKey string
	notifier  notifier
}

func NewService(secretKey string, store authStore, notifier notifier) *service {
	return &service{
		store:     store,
		secretKey: secretKey,
		notifier:  notifier,
	}
}

func (s *service) Authorize(ctx context.Context, userID string, ip string) (access string, refresh string, err error) {
	//JWT, SHA512, не храним
	token := Token{
		UserID: userID,
		Secret: uuid.New().String(),
		IP:     ip,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, token.MapToAcces(s.cl)) //todo написать clock
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, token.MapToRefresh(s.cl))
	access, err = accessToken.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", "", errors.Wrap(err, "failed to make access token")
	}

	refresh, err = refreshToken.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", "", errors.Wrap(err, "failed to make refresh token")
	}

	s.store.Save(ctx, userID, refresh) //todo написать методы postgres

	return access, refresh, nil
}

}

func (s *service) Refresh(ctx context.Context, refresh string, access string, ip string) (newAcess string, err error) {
	//храним в базе
	claims := jwt.Clains{}

}
