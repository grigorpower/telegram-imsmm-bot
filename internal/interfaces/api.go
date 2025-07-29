package interfaces

import (
	"context"
)

// TelegramSender определяет методы для отправки сообщений в Telegram
type TelegramSender interface {
	SendMessage(ctx context.Context, message string) error
}
