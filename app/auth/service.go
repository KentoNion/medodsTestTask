package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"medodsTest/auth/pkg"
)

type authStore interface {
	Save(ctx context.Context, token string, userID string) error
	Get(ctx context.Context, token string) (bool, error)
}

type notifier interface {
	NotifyNewLogin(ctx context.Context, userID string) error
}

type service struct {
	secretKey string
	store     authStore
	notifier  notifier
	cl        pkg.Clock
}

func NewService(secretKey string, store authStore, notifier notifier, cl pkg.Clock) *service {
	return &service{
		secretKey: secretKey,
		store:     store,
		notifier:  notifier,
		cl:        cl,
	}
}

func (s *service) Authorize(ctx context.Context, secret string, userID string, ip string) (Tokens AuthTokens, err error) {
	//JWT, SHA512, не храним
	token := Token{
		UserID: userID,
		Secret: secret, //uuid.New().String() <-----------------------------------------------------------
		IP:     ip,
	}
	result := AuthTokens{
		Access:  "",
		Refresh: "",
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, token.MapToAcces(s.cl))
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, token.MapToRefresh(s.cl))
	access, err := accessToken.SignedString([]byte(s.secretKey))
	if err != nil {
		return result, errors.Wrap(err, "failed to make access token")
	}

	refresh, err := refreshToken.SignedString([]byte(s.secretKey))
	if err != nil {
		return result, errors.Wrap(err, "failed to make refresh token")
	}

	err = s.store.Save(ctx, userID, refresh)
	if err != nil {
		return result, err
	}

	result = AuthTokens{
		Access:  access,
		Refresh: refresh,
	}
	return result, nil
}

func (s *service) Refresh(ctx context.Context, refresh string) (newAccess string, err error) {
	//храним в базе в виде хеша
	claims := jwt.MapClaims{}
	exists, err := s.store.Get(ctx, refresh)
	if err != nil {
		return "", errors.Wrap(err, "failed to check access token")
	}
	if !exists {
		return "", ErrRefreshTokenNotFound
	}

	token, err := jwt.ParseWithClaims(refresh, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to parse refresh token")
	}
	if !token.Valid {
		return "", ErrWrongToken
	}
	refreshToken := Token{}
	if err := refreshToken.Fill(claims); err != nil {
		return "", errors.Wrap(err, "failed to parse refresh token")
	}

	accessToken := Token{}
	if err := accessToken.Fill(claims); err != nil {
		return "", errors.Wrap(err, "failed to parse access token")
	}
	if refreshToken.Secret != accessToken.Secret {
		return "", ErrWrongToken
	}
	result := jwt.NewWithClaims(jwt.SigningMethodHS512, accessToken.MapToAcces(s.cl))
	access, err := result.SignedString(s.secretKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to make access token")
	}
	if err := s.notifier.NotifyNewLogin(ctx, accessToken.UserID); err != nil {
		return "", errors.Wrap(err, "failed to notify new login")
	}
	return access, nil
}
