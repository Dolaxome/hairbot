package telegram

import (
	"database/sql"
	"github.com/Dolaxome/hair-bot/pkg/config"
	"github.com/Dolaxome/hair-bot/pkg/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type Bot struct {
	bot   *tgbotapi.BotAPI
	botdb *storage.DB
	tx    *sql.Tx

	cfg *config.Config
}

func NewBot(bot *tgbotapi.BotAPI, cfg *config.Config) *Bot {
	return &Bot{bot: bot, cfg: cfg}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)
	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}

	b.handleUpdates(updates)

	return err
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	BotDB, err := storage.NewDB(b.cfg.DBAuth)
	if err != nil {
		log.Fatal(err)
	}
	b.botdb = BotDB
	b.tx, err = b.botdb.Sql.Begin()

	//file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//log.SetOutput(file)
	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				if err = b.handleCommand(update.Message); err != nil {
					b.handleError(update.Message.Chat.ID, err)
					//log.Println(err)
					log.Fatal(err)
				}
				continue
			}

			if err = b.handleMessage(update.Message); err != nil {
				b.handleError(update.Message.Chat.ID, err)
				//log.Println(err)
				log.Fatal(err)
			}
		} else if update.CallbackQuery != nil {
			if err = b.handleQuery(update.CallbackQuery); err != nil {
				b.handleError(update.CallbackQuery.Message.Chat.ID, err)
				//log.Println(err)
				log.Fatal(err)
			}
		}
	}
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}
