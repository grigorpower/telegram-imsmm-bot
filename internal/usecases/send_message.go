package usecases

import (
	"context"
	"fmt"
	"telegram-imsmm-bot/internal/interfaces"
)

// SendErrorService предоставляет методы для отправки сообщений телеграм ботом
type SendErrorService struct {
	telegram interfaces.TelegramSender
}

// NewSendErrorService создает новый экземпляр SendErrorService
// Принимает интерфейс TelegramSender для отправки сообщений в телеграм
func NewSendErrorService(telegram interfaces.TelegramSender) *SendErrorService {
	return &SendErrorService{
		telegram: telegram,
	}
}

// SendError отправляет сообщение об ошибке в телеграм
// Форматирует текст ошибки
func (s *SendErrorService) SendError(ctx context.Context, errMsg string, prefix string) error {
	err := s.telegram.SendMessage(ctx, errMsg)
	if err != nil {
		return fmt.Errorf("Не удалось отправить сообщение: %w", err)
	}
	return nil
}
