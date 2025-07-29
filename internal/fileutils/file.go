package file

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"
)

// Open открывает файл по заданному пути
func Open(filePath string, timeInterval time.Duration) *os.File {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		time.Sleep(timeInterval)
	}

	return f
}

// Search ищет в файле строку содержащую line
func Search(f *os.File, words []string, contain bool) (string, bool) {
	var scanLine string
	scanner := bufio.NewScanner(f)
	line := scanner.Text()
	for scanner.Scan() {
		for _, word := range words {
			if strings.Contains(line, word) {
				contain = true
				scanLine = line
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Ошибка чтения файла: ", err)
	}

	return scanLine, contain
}
