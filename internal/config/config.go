package config

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// TelegramConfig представляет конфигурацию телеграм бота
type TelegramConfig struct {
	BotToken string  // Токен телеграм бота
	ChatIDs  []int64 // Идентификаторы пользователей телеграм
}

// APICongig представялет конфигурацию приложения
type APIConfig struct {
	LogFilePath     string
	ErrorIndicators []string
	TimeInterval    time.Duration
	Device          string
	StartKeyWord    string
	Prefix          string
}

// LoadTelegramConfig загружает конфигурацию из переменных окружения
func LoadTelegramConfig() (*TelegramConfig, error) {
	// Чтение переменных окружения
	botToken := os.Getenv("BOT_TOKEN")
	chatIDsStr := os.Getenv("CHAT_IDs")

	// Проверка наличия обязательных переменных
	if botToken == "" || chatIDsStr == "" {
		log.Println("Необходимые переменные окружения отсутствуют", "BOT_TOKEN", botToken, "CHAT_IDs", chatIDsStr)
		return nil, errors.New("необходимые переменные окружения отсутствуют")
	}

	//Получаем отдельные ID чатов
	chatIDsStrSl := strings.Split(chatIDsStr, ",")
	var chatIDs []int64

	// Преобразование CHAT_IDs в []int64
	for _, idStr := range chatIDsStrSl {
		idStr = strings.TrimSpace(idStr) // убираем пробелы
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Println("Ошибка преобразования CHAT_IDs в []int64: ", err)
			return nil, err
		}
		chatIDs = append(chatIDs, id)
	}

	// Возвращаем конфигурацию
	return &TelegramConfig{
		BotToken: botToken,
		ChatIDs:  chatIDs,
	}, nil
}

// LoadAPIConfig загружает конфигурацию API
func LoadAPIConfig() (*APIConfig, error) {
	return &APIConfig{
		LogFilePath:     "C:/msdchem/1/mslogbk.log",
		ErrorIndicators: []string{"No Bottle in Gripper", "Error", "Fault", "Plunger", "Syringe", "Front inj door open"},
		TimeInterval:    10 * time.Minute,
		Device:          "EVR1",
		StartKeyWord:    "New version started",
		Prefix:          "R6890",
	}, nil
}
