package main

import (
	"github.com/Dolaxome/hair-bot/pkg/config"
	"github.com/Dolaxome/hair-bot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(cfg)
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}
	bot.RemoveWebhook()
	bot.Debug = true

	telegramBot := telegram.NewBot(bot, cfg)

	if err = telegramBot.Start(); err != nil {
		log.Fatal(err)
	}
}
