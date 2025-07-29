package adapter

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TelegramAdapter реализует интерфейс TelegramSender для отправки сообщений в Telegram.
type TelegramAdapter struct {
	bot     *tgbotapi.BotAPI // Экземпляр бота API Telegram
	chatIDs []int64          // ID чатов, куда будут отправляться сообщения
}

// NewTelegramAdapter создаёт новый экземпляр TelegramAdapter.
// Принимает токен бота (botToken) и ID чата (chatID).
// Возвращает ошибку, если не удалось инициализировать бота.
func NewTelegramAdapter(botToken string, chatIDs []int64) (*TelegramAdapter, error) {
	// Создаём новый экземпляр BotAPI с использованием токена бота
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		// Если произошла ошибка при инициализации бота, возвращаем её
		return nil, err
	}
	// Возвращаем инициализированный адаптер с ботом и ID чата
	return &TelegramAdapter{bot: bot, chatIDs: chatIDs}, nil
}

// SendMessage отправляет текстовое сообщение в Telegram-чат.
// Принимает контекст (ctx) и текст сообщения (message).
// Возвращает ошибку, если сообщение не удалось отправить.
func (t *TelegramAdapter) SendMessage(ctx context.Context, message string) error {
	var lastErr error
	// Создаём новое текстовое сообщение для отправки в указанные чаты
	for _, chatID := range t.chatIDs {
		msg := tgbotapi.NewMessage(chatID, message)

		// Отправляем сообщение через API Telegram
		if _, err := t.bot.Send(msg); err != nil {
			lastErr = err
		}
	}
	// Возвращаем ошибку, если отправка не удалась
	return lastErr
}
