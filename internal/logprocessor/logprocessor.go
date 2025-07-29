package logprocessor

import (
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

func ErrorDetection(c *config.APIConfig, s *usecases.SendErrorService) {
	var lastOffset int64 = 0
	var currentLine string
	var mu sync.Mutex // создаем мьютекс

	for {
		// Открываем файл
		f := file.Open(c.LogFilePath, c.TimeInterval)

		// Захватываем мьютекс перед началом сканирования на 'New version started'
		mu.Lock()

		// Первый цикл: сканируем на наличие 'New version started'
		var contain bool
		startLine := ""
		startLine, contain = file.Search(f, []string{c.StartKeyWord}, contain)

		// Освобождаем мьютекс после завершения поиска 'New version started'
		mu.Unlock()
		f.Close() // Закрывает файл

		if startLine != currentLine {
			currentLine = startLine
			lastOffset = 0
			time.Sleep(c.TimeInterval)
		}

		f = file.Open(c.LogFilePath, c.TimeInterval)
		f.Seek(lastOffset, os.SEEK_SET)

		mu.Lock()
		var errLine string
		errLine, contain = file.Search(f, c.ErrorIndicators, contain)

		if contain {
			errLine = strings.TrimPrefix(errLine, c.Prefix)
			s.SendError(nil, errLine)
		}

		offset, _ := f.Seek(0, os.SEEK_CUR)
		lastOffset = offset
		f.Close()
		mu.Unlock()
		time.Sleep(c.TimeInterval)
	}
}
