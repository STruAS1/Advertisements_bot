package config

import (
	"encoding/gob"
	"os"
)

type Setings struct {
	Ads      AdsSetings
	Texts    Texts
	City     City
	Payments PaymentsSetings
	Docs     Docs
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
}

type City struct {
	MaxCountOfCity int8
}
type Docs struct {
	VideoUrl string
	VideoID  string
	Text     string
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
		Texts: Texts{
			MainText: "Тест",
		},
		City: City{
			MaxCountOfCity: 10,
		},
		Payments: PaymentsSetings{
			Metods:        []PaymentsMetod{{Title: "test", Discription: "test", Cardnumber: "test"}},
			MinimalAmount: 100,
			MaxAmount:     1000000000,
		},
		Docs: Docs{
			VideoUrl: "",
			Text:     "Обучение",
		},
	}
	_ = encoder.Encode(setings)
	GlobalSettings = setings
	file.Close()
}
