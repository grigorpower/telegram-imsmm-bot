package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	logFilePath   = "C:/msdchem/1/mslogbk.log"
	telegramToken = ""
	timeInterval  = 30 * time.Minute
)

var (
	userIDs         = []int64{}
	errorIndicators = []string{"No Bottle in Gripper", "Error", "Fault", "Plunger", "Syringe", "Front inj door open"}
)

func main() {
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
	}

	var lastOffset int64 = 0
	var device string
	var inSeries bool = true

	var mu sync.Mutex // создаем мьютекс

	for {
		// Открываем файл
		f, err := os.Open(logFilePath)
		if err != nil {
			log.Println("Ошибка открытия файла:", err)
			time.Sleep(timeInterval)
			continue
		}
		log.Println("mslogbk открыт1")

		// Устанавливаем начальную позицию
		//f.Seek(lastOffset, os.SEEK_SET)

		// Захватываем мьютекс перед началом сканирования на 'completed'
		mu.Lock()

		// Первый цикл: сканируем на наличие 'completed'
		scanner := bufio.NewScanner(f)
		//seriesCompleted := false
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "completed") {
				log.Println("Серия завершена")
				inSeries = false
				//lastOffset, _ = f.Seek(0, os.SEEK_CUR)
				break
			}
		}

		if err := scanner.Err(); err != nil {
			log.Println("Ошибка чтения файла:", err)
		}

		// Освобождаем мьютекс после завершения поиска 'completed'
		mu.Unlock()

		// Обновляем смещение
		//lastOffset = 0

		f.Close()
		//log.Println(inSeries, lastOffset)
		log.Println("mslogbk закрыт1")

		// Если серия завершена, ждем и начинаем заново
		if !inSeries {
			log.Println("серия завершена.")
			lastOffset = 0
			inSeries = true
			time.Sleep(timeInterval)
			continue
		}

		// Второй цикл: обработка ошибок, запускается только после unlock
		mu.Lock()
		f, err = os.Open(logFilePath)
		if err != nil {
			log.Println("Ошибка открытия файла:", err)
			mu.Unlock()
			time.Sleep(timeInterval)
			continue
		}

		f.Seek(lastOffset, os.SEEK_SET)
		log.Println("mslogbk открыт для ошибок с позиции", lastOffset)

		scanner = bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			// Обработка ошибок
			for _, errorIndicator := range errorIndicators {
				if strings.Contains(line, errorIndicator) {
					//log.Println("проверка ошибок")
					sendTelegramMessage(bot, fmt.Sprintf("Ошибка на %s: %s", device, line))
					log.Println("Сообщение отправлено:", line)
				}
			}
		}

		if err := scanner.Err(); err != nil {
			log.Println("Ошибка чтения файла:", err)
		}

		offset, _ := f.Seek(0, os.SEEK_CUR)
		lastOffset = offset
		f.Close()
		//log.Println(inSeries, lastOffset)
		log.Println("mslogbk закрыт2")
		mu.Unlock()

		log.Println("Цикл обработки ошибок завершен")
		time.Sleep(timeInterval)
	}
}

func sendTelegramMessage(bot *tgbotapi.BotAPI, message string) {
	for _, userID := range userIDs {
		msg := tgbotapi.NewMessage(userID, message)
		_, err := bot.Send(msg)
		if err != nil {
			log.Println("Ошибка отправки сообщения:", err)
		}
	}
}
