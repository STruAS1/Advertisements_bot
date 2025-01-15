package main

import (
	"log"
	"tgbotBARAHOLKA/backend"
	"tgbotBARAHOLKA/bot"
	"tgbotBARAHOLKA/config"
	"tgbotBARAHOLKA/db"
	"time"
)

func main() {
	cfg := config.LoadConfig()

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
