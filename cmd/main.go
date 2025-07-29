package main

import (
	"log"
	"telegram-imsmm-bot/internal/adapter"
	"telegram-imsmm-bot/internal/config"
	"telegram-imsmm-bot/internal/logprocessor"
	"telegram-imsmm-bot/internal/usecases"

	"github.com/joho/godotenv"
)

func main() {
	// Загрузка .env файла
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден или не загружен")
	}

	// Загрузка конфигурации телеграм-бота
	cfgTg, err := config.LoadTelegramConfig()
	if err != nil {
		log.Println("Ошибка загрузки конфигурации тг-бота: ", err)
	}

	// Загрузка конфигурации API
	cfgAPI, err := config.LoadAPIConfig()
	if err != nil {
		log.Println("Ошибка загрузка конфигурации программы: ", err)
	}

	// Инициализация адаптеров
	telegramAdapter, err := adapter.NewTelegramAdapter(cfgTg.BotToken, cfgTg.ChatIDs)
	if err != nil {
		log.Println("Не удалось инициализировать TelegramAdapter: ", err)
	}

	// Инициализация сервисов
	sendErrorService := usecases.NewSendErrorService(telegramAdapter)

	go logprocessor.ErrorDetection(cfgAPI, sendErrorService)

	// Бесконечный цикл для работы программы
	select {}
}
