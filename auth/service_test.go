package auth

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"medodsTest/auth/mock"
	"medodsTest/auth/pkg"
	"testing"
	"time"
)

func TestAuthorize(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mock.NewMockauthStore(ctrl)
	notifier := mock.NewMocknotifier(ctrl)
	svc := service{
		secretKey: "test_key",
		store:     store,
		notifier:  notifier,
		cl:        pkg.StubClock{time.Date(2024, time.September, 24, 0, 0, 0, 0, time.UTC)},
	}

	expectRefresh := gomock.Any()
	store.EXPECT().Save(gomock.Any(), "test_user", expectRefresh).Return(nil) //Ожидаю что будет вызвана функция save с указанными параметрами и для этих вызовов параметров верни nil

	ctx := context.Background()
	access, refresh, err := svc.Authorize(ctx, "test_user", "1.1.1.1")
	require.NoError(t, err)
	require.Equal(t, "test_user", access)
	require.Equal(t, "1.1.1.1", refresh)
}

func TestVerifyNotExistedToken(t *testing.T) {

}
