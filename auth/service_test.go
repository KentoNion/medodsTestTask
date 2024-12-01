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
	store.EXPECT().Save(gomock.Any(), "test_user", expectRefresh).Return(nil)

	ctx := context.Background()
	tok, err := svc.Authorize(ctx, "password", "test_user", "1.1.1.1")
	require.NoError(t, err)
	require.Equal(t, "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjcyMjI0MDAsImlwIjoiMS4xLjEuMSIsInNlY3JldCI6InBhc3N3b3JkIiwidXNlcl9pZCI6InRlc3RfdXNlciJ9.6sWj1XHsIs5m-KzpVNRGjkuvxcrO4kMjElXCzChZ4xhUlS04Yi6hp4tcDjmnoEP5WQ6jDU1ZPeEJ6ZEEMBMbIA", tok.Access)
	require.Equal(t, "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTA5LTI0VDAwOjAwOjAwWiIsImlwIjoiMS4xLjEuMSIsInNlY3JldCI6InBhc3N3b3JkIiwidXNlcl9pZCI6InRlc3RfdXNlciJ9.7AoubHHDO0S3FsxWo7ldOrSQmdmfY6_rUNGXdxk3c0MU43DmlrwqWgokEaV0UoRnAoYAqFrcUz2VokFAiZdBhQ", tok.Refresh)
}

func testRefresh(t *testing.T, svc *service) {

}

func TestVerifyNotExistedToken(t *testing.T) {

}
