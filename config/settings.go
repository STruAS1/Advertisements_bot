package config

import (
	"encoding/gob"
	"os"
)

type Setings struct {
	Ads              AdsSetings
	Texts            Texts
	City             City
	Payments         PaymentsSetings
	Docs             []Docs
	Buttons          [16]Button
	WitheListDomines []string
	VerificationCost uint
}

type AdsSetings struct {
	Sufix     string
	CostLimit uint
}
type PaymentsSetings struct {
	Metods        []PaymentsMetod
	MinimalAmount uint
	MaxAmount     uint
}

type PaymentsMetod struct {
	Title       string
	Discription string
	Cardnumber  string
}

type Texts struct {
	MainText string
	AddsMenu string
	DocsText string
}
type Button struct {
	ButtonText  string
	Discription string
}

type City struct {
	MaxCountOfCity int8
}
type Docs struct {
	ButtonName string
	VideoID    string
	Text       string
}

var GlobalSettings Setings

func Save(setings Setings) {
	GlobalSettings = setings
	file, _ := os.Create("config.gob")
	encoder := gob.NewEncoder(file)
	encoder.Encode(setings)
	file.Close()
}

func CreateDefaultSettings() {
	file, _ := os.Create("config.gob")
	encoder := gob.NewEncoder(file)
	setings := Setings{
		Ads: AdsSetings{
			Sufix:     "📣<b><a href='https://t.me/TESTESTESTE312312bot'>ПОДАТЬ ОБЪЯВЛЕНИЕ</a></b>📣\n\n💬<b><i>Вопросы задавайте в комментариях</i></b>🔻",
			CostLimit: 2000,
		},
		VerificationCost: 2000,
		Texts: Texts{
			MainText: "Тест",
			AddsMenu: "Тест",
		},
		City: City{
			MaxCountOfCity: 10,
		},
		Payments: PaymentsSetings{
			Metods:        []PaymentsMetod{{Title: "test", Discription: "test", Cardnumber: "test"}},
			MinimalAmount: 100,
			MaxAmount:     1000000000,
		},
		Docs: []Docs{{
			ButtonName: "Обучение",
			Text:       "Обучение",
		}},
		Buttons: [16]Button{
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
		},
		WitheListDomines: []string{
			"youtube.com",
		},
	}
	_ = encoder.Encode(setings)
	GlobalSettings = setings
	file.Close()
}
