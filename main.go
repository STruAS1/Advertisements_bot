package main

import (
	"encoding/gob"
	"log"
	"os"
	"tgbotBARAHOLKA/backend"
	"tgbotBARAHOLKA/bot"
	"tgbotBARAHOLKA/config"
	"tgbotBARAHOLKA/db"
	"time"
)

func main() {
	cfg := config.LoadConfig()
	file, err := os.Open("config.gob")
	if err != nil {
		if os.IsNotExist(err) {
			config.CreateDefaultSettings()
		} else {
			log.Panic("Ошибка с файлом конфигурации!")
		}
	} else {
		decoder := gob.NewDecoder(file)
		_ = decoder.Decode(&config.GlobalSettings)
	}
	file.Close()
	db.Connect(cfg)
	go backend.StartBackend()
	go func() {
		for {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Паника: %v\n", r)
				}
			}()

			log.Println("Запуск бота...")
			bot.StartBot(cfg)
			log.Println("Бот завершил работу. Перезапуск через 5 секунд...")

			time.Sleep(5 * time.Second)
		}
	}()

	log.Println("Основной процесс работает...")
	select {}
}
