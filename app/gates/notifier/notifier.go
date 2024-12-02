package notify

import (
	"context"
	"log"
)

// Notifier - интерфейс для уведомлений.
type Notifier interface {
	NotifyNewLogin(ctx context.Context, userID string) error
}

// mockNotifier - моковая реализация интерфейса Notifier.
type mockNotifier struct{}

// NotifyNewLogin - метод мокового уведомления о новом логине.
func (m *mockNotifier) NotifyNewLogin(ctx context.Context, userID string) error {
	log.Printf("Mock notification: new login for user %s", userID)
	return nil
}

// InitNotifier - инициализация мокового нотификатора.
func InitNotifier() Notifier {
	return &mockNotifier{}
}
