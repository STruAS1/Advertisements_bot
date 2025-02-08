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
	setting := config.GlobalSettings
	setting.Buttons = [17]config.Button{
		{ButtonText: "Объявление", Discription: "Объявление"},
		{ButtonText: "Обучение", Discription: "Обучение"},
		{ButtonText: "Профиль", Discription: "Профиль"},
		{ButtonText: "Подписаться на канал", Discription: "Подписаться на канал"},
		{ButtonText: "Подписался", Discription: "Подписался"},
		{ButtonText: "« Назад", Discription: "Кнопка назад"},
		{ButtonText: "Пополнить баланс", Discription: "Пополнить баланс"},
		{ButtonText: "Перевести средства", Discription: "Перевести средства"},
		{ButtonText: "Добавить объявление", Discription: "Добавить объявление"},
		{ButtonText: "Мои объявления", Discription: "Мои объявления"},
		{ButtonText: "Пред просмотр", Discription: "Пред просмотр"},
		{ButtonText: "Сохранить", Discription: "Сохранить"},
		{ButtonText: "🗑️ Удалить", Discription: "🗑️ Удалить"},
		{ButtonText: "✏️ Редактировать", Discription: "✏️ Редактировать"},
		{ButtonText: "🚫 Отмена ", Discription: "✏️ Редактировать"},
		{ButtonText: "📋 Сохранить", Discription: "📋 Сохранить"},
		{ButtonText: "Изменить город", Discription: "Изменить город"},
	}
	config.Save(setting)
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
