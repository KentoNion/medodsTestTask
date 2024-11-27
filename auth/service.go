package auth

import (
	"context"
	"github.com/google/uuid"
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
	//jwt не храним
	token := Token{
		UserID: userID,
		Secret: uuid.New().String(),
		IP:     ip,
	}

}

func (s *service) Refresh(ctx context.Context, refresh string, access string, ip string) (newAcess string, err error) {
	//храним в базе
	claims := jwt.Clains{}

}
