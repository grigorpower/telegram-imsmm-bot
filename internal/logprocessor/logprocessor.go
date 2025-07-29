package logprocessor

import (
	"log"
	"os"
	"strings"
	"sync"
	"telegram-imsmm-bot/internal/config"
	file "telegram-imsmm-bot/internal/fileutils"
	"telegram-imsmm-bot/internal/usecases"
	"time"
)

type LogProcessor struct {
	confAPI   *config.APIConfig
	confTgBot *config.TelegramConfig
}

func (s *SendErrorService) ErrorDetection(c *cfgAPI) {
	var lastOffset int64 = 0
	var currentLine string
	var mu sync.Mutex // создаем мьютекс

	for {
		//Загружаем конфигурацию
		conf, err := config.LoadAPIConfig()
		if err != nil {
			log.Println("Ошибка загрузки config: ", err)
		}

		// Открываем файл
		f := file.Open(conf.LogFilePath, conf.TimeInterval)

		// Захватываем мьютекс перед началом сканирования на 'New version started'
		mu.Lock()

		// Первый цикл: сканируем на наличие 'New version started'
		var contain bool
		startLine := ""
		startLine, contain = file.Search(f, []string{conf.StartKeyWord}, contain)

		// Освобождаем мьютекс после завершения поиска 'New version started'
		mu.Unlock()
		f.Close() // Закрывает файл

		if startLine != currentLine {
			currentLine = startLine
			lastOffset = 0
			time.Sleep(conf.TimeInterval)
		}

		f = file.Open(conf.LogFilePath, conf.TimeInterval)
		f.Seek(lastOffset, os.SEEK_SET)

		mu.Lock()
		var errLine string
		errLine, contain = file.Search(f, conf.ErrorIndicators, contain)

		if contain {
			errLine = strings.TrimPrefix(errLine, conf.Prefix)
			usecases.SendErrorService(errLine)
		}

		offset, _ := f.Seek(0, os.SEEK_CUR)
		lastOffset = offset
		f.Close()
		mu.Unlock()
		time.Sleep(conf.TimeInterval)
	}
}
