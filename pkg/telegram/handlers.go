package telegram

import (
	"database/sql"
	"fmt"
	"github.com/Dolaxome/hair-bot/pkg/storage"
	"github.com/Dolaxome/hair-bot/tools"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"strings"
)

// vars
const commandStart = "start"
const pending = "pending"
const adminMenu = "adminMenu"
const userMenu = "userMenu"
const userSucc = "userSucc"
const userDecl = "userDecl"
const procedureProc = "procedureProc"
const requests = "requests"
const requi = "requi"
const answ = "answ"

// keyboards
var pendingKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Нік"),
		tgbotapi.NewKeyboardButton("Увійти"),
	),
)
var adminKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Запити"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Додати учасника - по юзерайді"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Видалити учасника - по юзерайді"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Список учасників"),
	),
)
var userKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Підібрати процедуру"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Приклади структури волосся по формі А-силует"),
	),
)
var userKeyboard2 = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Підібрати нову процедуру"),
	),
)
var cancelKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Меню"),
	),
)
var nilKeyboard = tgbotapi.InlineKeyboardMarkup{
	InlineKeyboard: make([][]tgbotapi.InlineKeyboardButton, 0),
}

// inline
var chooseThickInline = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Тонке", "answ;0;thin"),
		tgbotapi.NewInlineKeyboardButtonData("Середнє", "answ;0;medium"),
		tgbotapi.NewInlineKeyboardButtonData("Густе", "answ;0;thick"),
	),
)
var chooseAInline = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Менше А-силует", "answ;1;less a"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("А-силует", "answ;1;equal a"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Більше А-силует", "answ;1;more a"),
	),
)
var chooseCurlInline = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Рівне", "answ;2;straight"),
		tgbotapi.NewInlineKeyboardButtonData("Кучеряве", "answ;2;curly"),
	),
)
var chooseDamageInline = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("1", "answ;3;light film"),
		tgbotapi.NewInlineKeyboardButtonData("2", "answ;3;light film"),
		tgbotapi.NewInlineKeyboardButtonData("3", "answ;3;light film"),
		tgbotapi.NewInlineKeyboardButtonData("4", "answ;3;thick film"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("5 (початкова)", "answ;3;fifthInit"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("5 (кінцева)", "answ;3;fifthFin"),
	),
)

var inlineArr = []tgbotapi.InlineKeyboardMarkup{chooseThickInline, chooseAInline, chooseCurlInline, chooseDamageInline}

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleQuery(query *tgbotapi.CallbackQuery) error {
	err := b.handleHotQuery(query)

	//remove markup
	//msgedit := tgbotapi.NewEditMessageReplyMarkup(query.Message.Chat.ID, query.Message.MessageID, nilKeyboard)
	//_, err = b.bot.Send(msgedit)
	return err
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	var queryRowed string
	var err error

	row := b.botdb.Sql.QueryRow("SELECT query FROM querys WHERE chat_id=$1", message.Chat.ID)
	err = row.Scan(&queryRowed)
	if err != nil {
		return err
	}

	log.Printf("[%s] %s", message.From.UserName, message.Text)
	msg := tgbotapi.NewMessage(message.Chat.ID, b.cfg.Messages.UnknownCommand)

	//stroke := strings.Fields(message.Text)

	switch queryRowed {
	case pending:
		nick := pendingKeyboard.Keyboard[0][0].Text
		signIn := pendingKeyboard.Keyboard[0][1].Text

		switch message.Text {
		case nick:
			var admins []storage.Admin

			msg = tgbotapi.NewMessage(message.Chat.ID, "Зв'яжіться з адміністратором для придбання доступу: https://www.instagram.com/alina_pulvas/")

			var admincheck int
			var usercheck int
			row = b.botdb.Sql.QueryRow("SELECT COUNT(id) FROM admins WHERE chat_id=$1", message.Chat.ID)
			if err = row.Scan(&admincheck); err != nil {
				return err
			}
			row = b.botdb.Sql.QueryRow("SELECT COUNT(id) FROM requests WHERE user_chat_id=$1", message.Chat.ID)
			if err = row.Scan(&usercheck); err != nil {
				return err
			}

			if admincheck == 0 && usercheck == 0 {
				if admins, err = b.botdb.GetAdmins(); err != nil {
					return err
				}
				for _, admin := range admins {
					adminmsg := tgbotapi.NewMessage(admin.ChatID, "+1 Запит на покупку")
					if _, err = b.bot.Send(adminmsg); err != nil {
						return err
					}
				}

				var thname string
				if message.From.UserName != "" {
					thname = "@" + message.From.UserName
				}
				requesto := storage.Request{
					TelegramName: thname,
					FirstName:    message.From.FirstName,
					ChatID:       message.Chat.ID,
				}

				if err = b.botdb.AddRequest(requesto); err != nil {
					return err
				}
				if err = b.botdb.FlushRequest(); err != nil {
					return err
				}
			}
		case signIn:
			var admincheck int
			var usercheck int
			row = b.botdb.Sql.QueryRow("SELECT COUNT(id) FROM admins WHERE chat_id=$1", message.Chat.ID)
			if err = row.Scan(&admincheck); err != nil {
				return err
			}
			row = b.botdb.Sql.QueryRow("SELECT COUNT(id) FROM requests WHERE user_chat_id=$1 AND status=true", message.Chat.ID)
			if err = row.Scan(&usercheck); err != nil {
				return err
			}

			if admincheck != 0 {
				msg = tgbotapi.NewMessage(message.Chat.ID, b.cfg.Messages.Welcome)
				msg.ReplyMarkup = adminKeyboard
				if _, err = b.botdb.Sql.Exec("UPDATE querys SET query=$1 WHERE chat_id=$2", adminMenu, message.Chat.ID); err != nil {
					return err
				}
			} else if usercheck != 0 {
				msg = tgbotapi.NewMessage(message.Chat.ID, b.cfg.Messages.Welcome)
				msg.ReplyMarkup = userKeyboard
				if _, err = b.botdb.Sql.Exec("UPDATE querys SET query=$1 WHERE chat_id=$2", userMenu, message.Chat.ID); err != nil {
					return err
				}
			} else {
				msg = tgbotapi.NewMessage(message.Chat.ID, b.cfg.Messages.NoAuth)
				msg.ReplyMarkup = pendingKeyboard
				if _, err = b.botdb.Sql.Exec("UPDATE querys SET query=$1 WHERE chat_id=$2", pending, message.Chat.ID); err != nil {
					return err
				}
			}
		}
	case adminMenu:
		var checkauth int
		row = b.botdb.Sql.QueryRow("SELECT COUNT(id) FROM admins WHERE chat_id=$1", message.Chat.ID)
		if err = row.Scan(&checkauth); err != nil {
			return err
		}
		if checkauth == 0 {
			msg = tgbotapi.NewMessage(message.Chat.ID, b.cfg.Messages.NoAccess)
			msg.ReplyMarkup = pendingKeyboard
			if _, err = b.botdb.Sql.Exec("UPDATE querys SET query=$1 WHERE chat_id=$2", pending, message.Chat.ID); err != nil {
				return err
			}
			break
		}

		requestss := adminKeyboard.Keyboard[0][0].Text
		addUser := adminKeyboard.Keyboard[1][0].Text
		deleteUser := adminKeyboard.Keyboard[2][0].Text
		listUsers := adminKeyboard.Keyboard[3][0].Text

		switch message.Text {
		case requestss:
			var check int
			var reqee []storage.Request
			row = b.botdb.Sql.QueryRow("SELECT COUNT(id) FROM requests WHERE status IS NULL")
			if err = row.Scan(&check); err != nil {
				return err
			}

			if reqee, err = b.botdb.GetRequests(""); err != nil {
				return err
			}
			if reqee != nil {
				msg = tgbotapi.NewMessage(message.Chat.ID, "Необроблених запитів на придбання: "+strconv.Itoa(check))
				msg.ReplyMarkup = tools.FormRequestsInline(reqee, "requests")
			} else {
				msg = tgbotapi.NewMessage(message.Chat.ID, "Запитів немає")
			}
		case addUser:
			msg = tgbotapi.NewMessage(message.Chat.ID, b.cfg.Messages.AddUser)
			msg.ReplyMarkup = cancelKeyboard
			if _, err = b.botdb.Sql.Exec("UPDATE querys SET query=$1 WHERE chat_id=$2", userSucc, message.Chat.ID); err != nil {
				return err
			}
		case deleteUser:
			msg = tgbotapi.NewMessage(message.Chat.ID, b.cfg.Messages.AddUser)
			msg.ReplyMarkup = cancelKeyboard
			if _, err = b.botdb.Sql.Exec("UPDATE querys SET query=$1 WHERE chat_id=$2", userDecl, message.Chat.ID); err != nil {
				return err
			}
		case listUsers:
			var reqee []storage.Request

			if reqee, err = b.botdb.GetRequests("WHERE status=true"); err != nil {
				return err
			}
			if reqee != nil {
				msg = tgbotapi.NewMessage(message.Chat.ID, "Оберіть користувача:")
				msg.ReplyMarkup = tools.FormRequestsInline(reqee, "requi")
			} else {
				msg = tgbotapi.NewMessage(message.Chat.ID, "Учасників немає")
			}
		}
	case userMenu:
		var checkauth int
		row = b.botdb.Sql.QueryRow("SELECT COUNT(id) FROM requests WHERE user_chat_id=$1 AND status = true", message.Chat.ID)
		if err = row.Scan(&checkauth); err != nil {
			return err
		}
		if checkauth == 0 {
			msg = tgbotapi.NewMessage(message.Chat.ID, b.cfg.Messages.NoAccess)
			msg.ReplyMarkup = pendingKeyboard
			if _, err = b.botdb.Sql.Exec("UPDATE querys SET query=$1 WHERE chat_id=$2", pending, message.Chat.ID); err != nil {
				return err
			}
			break
		}

		procedure := userKeyboard.Keyboard[0][0].Text
		examples := userKeyboard.Keyboard[1][0].Text

		switch message.Text {
		case procedure:
			var param []string
			row = b.botdb.Sql.QueryRow("SELECT params FROM answers WHERE id = 1")
			if err = row.Scan((*pq.StringArray)(&param)); err != nil {
				return err
			}
			fmt.Println(param[0])
			msg = tgbotapi.NewMessage(message.Chat.ID, "Підібрати процедуру:")
			msg.ReplyMarkup = userKeyboard2
			if _, err = b.bot.Send(msg); err != nil {
				return err
			}

			msg = tgbotapi.NewMessage(message.Chat.ID, b.cfg.Messages.ChooseThick)
			msg.ReplyMarkup = chooseThickInline

			if _, err = b.botdb.Sql.Exec("UPDATE querys SET statecounter=$1 WHERE chat_id=$2", sql.NullString{}, message.Chat.ID); err != nil {
				return err
			}

			if _, err = b.botdb.Sql.Exec("UPDATE querys SET query=$1 WHERE chat_id=$2", procedureProc, message.Chat.ID); err != nil {
				return err
			}
		case examples:
			var photoBytes1 []byte
			var photoBytes2 []byte
			var texts = []string{"Більше А-силует", "А-силует"}
			var paths = []string{"media/MoreA1.PNG", "media/MoreA2.PNG", "media/A1.JPG", "media/A2.JPG"}
			for i, values := range texts {
				msg = tgbotapi.NewMessage(message.Chat.ID, values)
				if _, err = b.bot.Send(msg); err != nil {
					return err
				}
				photoBytes1, err = os.ReadFile(paths[i*2])
				if err != nil {
					return err
				}
				photoBytes2, err = os.ReadFile(paths[(i*2)+1])
				if err != nil {
					return err
				}
				photoFileBytes1 := tgbotapi.FileBytes{
					Name:  "picture",
					Bytes: photoBytes1,
				}
				photoFileBytes2 := tgbotapi.FileBytes{
					Name:  "picture",
					Bytes: photoBytes2,
				}
				photo1 := tgbotapi.NewPhotoUpload(message.Chat.ID, photoFileBytes1)
				photo2 := tgbotapi.NewPhotoUpload(message.Chat.ID, photoFileBytes2)
				if _, err = b.bot.Send(photo1); err != nil {
					return err
				}
				if _, err = b.bot.Send(photo2); err != nil {
					return err
				}
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, "Головне меню")
		}
	case procedureProc:
		var checkauth int
		row = b.botdb.Sql.QueryRow("SELECT COUNT(id) FROM requests WHERE user_chat_id=$1 AND status = true", message.Chat.ID)
		if err = row.Scan(&checkauth); err != nil {
			return err
		}
		if checkauth == 0 {
			msg = tgbotapi.NewMessage(message.Chat.ID, b.cfg.Messages.NoAccess)
			msg.ReplyMarkup = pendingKeyboard
			if _, err = b.botdb.Sql.Exec("UPDATE querys SET query=$1 WHERE chat_id=$2", pending, message.Chat.ID); err != nil {
				return err
			}
			break
		}

		newProcedure := userKeyboard2.Keyboard[0][0].Text

		switch message.Text {
		case newProcedure:
			if _, err = b.botdb.Sql.Exec("UPDATE querys SET statecounter=$1 WHERE chat_id=$2", sql.NullString{}, message.Chat.ID); err != nil {
				return err
			}

			msg = tgbotapi.NewMessage(message.Chat.ID, b.cfg.Messages.ChooseThick)
			msg.ReplyMarkup = chooseThickInline
		}
	case userSucc:
		var checkauth int
		row = b.botdb.Sql.QueryRow("SELECT COUNT(id) FROM admins WHERE chat_id=$1", message.Chat.ID)
		if err = row.Scan(&checkauth); err != nil {
			return err
		}
		if checkauth == 0 {
			msg = tgbotapi.NewMessage(message.Chat.ID, b.cfg.Messages.NoAccess)
			msg.ReplyMarkup = pendingKeyboard
			if _, err = b.botdb.Sql.Exec("UPDATE querys SET query=$1 WHERE chat_id=$2", pending, message.Chat.ID); err != nil {
				return err
			}
			break
		}

		if message.Text == cancelKeyboard.Keyboard[0][0].Text {
			msg = tgbotapi.NewMessage(message.Chat.ID, b.cfg.Messages.MainMenu)
			msg.ReplyMarkup = adminKeyboard
			if _, err = b.botdb.Sql.Exec("UPDATE querys SET query=$1 WHERE chat_id=$2", adminMenu, message.Chat.ID); err != nil {
				return err
			}
		} else {
			var idcheck int
			var intChat int64
			intChat, err = strconv.ParseInt(message.Text, 10, 0)
			row = b.botdb.Sql.QueryRow("SELECT COUNT(id) FROM requests WHERE user_chat_id=$1 AND (status IS NULL OR status IS false)", intChat)
			if err = row.Scan(&idcheck); err != nil {
				return err
			}

			if idcheck != 0 {
				var SQuerry string
				arg := storage.SQLBuilderArgs{Query: "update", Table: "requests", Columns: []string{"status"}, WhereString: "WHERE user_chat_id"}
				if SQuerry, err = b.botdb.SQLBuilder(arg); err != nil {
					return err
				}

				if _, err = b.tx.Exec(SQuerry, true, intChat); err != nil {
					b.tx.Rollback()
					return err
				}
				if err = b.tx.Commit(); err != nil {
					return err
				}
				b.tx, err = b.botdb.Sql.Begin()

				msg = tgbotapi.NewMessage(message.Chat.ID, b.cfg.Messages.MainMenu)
				msg.ReplyMarkup = adminKeyboard
				if _, err = b.botdb.Sql.Exec("UPDATE querys SET query=$1 WHERE chat_id=$2", adminMenu, message.Chat.ID); err != nil {
					return err
				}
			} else {
				return errNoReq
			}
		}
	case userDecl:
		var checkauth int
		row = b.botdb.Sql.QueryRow("SELECT COUNT(id) FROM admins WHERE chat_id=$1", message.Chat.ID)
		if err = row.Scan(&checkauth); err != nil {
			return err
		}
		if checkauth == 0 {
			msg = tgbotapi.NewMessage(message.Chat.ID, b.cfg.Messages.NoAccess)
			msg.ReplyMarkup = pendingKeyboard
			if _, err = b.botdb.Sql.Exec("UPDATE querys SET query=$1 WHERE chat_id=$2", pending, message.Chat.ID); err != nil {
				return err
			}
			break
		}

		if message.Text == cancelKeyboard.Keyboard[0][0].Text {
			msg = tgbotapi.NewMessage(message.Chat.ID, b.cfg.Messages.MainMenu)
			msg.ReplyMarkup = adminKeyboard
			if _, err = b.botdb.Sql.Exec("UPDATE querys SET query=$1 WHERE chat_id=$2", adminMenu, message.Chat.ID); err != nil {
				return err
			}
		} else {
			var idcheck int
			var intChat int64
			intChat, err = strconv.ParseInt(message.Text, 10, 0)
			row = b.botdb.Sql.QueryRow("SELECT COUNT(id) FROM requests WHERE user_chat_id=$1 AND status IS TRUE", intChat)
			if err = row.Scan(&idcheck); err != nil {
				return err
			}

			if idcheck != 0 {
				var SQuerry string
				arg := storage.SQLBuilderArgs{Query: "update", Table: "requests", Columns: []string{"status"}, WhereString: "WHERE user_chat_id"}
				if SQuerry, err = b.botdb.SQLBuilder(arg); err != nil {
					return err
				}

				if _, err = b.tx.Exec(SQuerry, false, intChat); err != nil {
					b.tx.Rollback()
					return err
				}
				if err = b.tx.Commit(); err != nil {
					return err
				}
				b.tx, err = b.botdb.Sql.Begin()

				msg = tgbotapi.NewMessage(message.Chat.ID, b.cfg.Messages.MainMenu)
				msg.ReplyMarkup = adminKeyboard
				if _, err = b.botdb.Sql.Exec("UPDATE querys SET query=$1 WHERE chat_id=$2", adminMenu, message.Chat.ID); err != nil {
					return err
				}
			} else {
				return errNoReq
			}
		}
	default:
		_, err = b.bot.Send(msg)
		return err
	}

	_, err = b.bot.Send(msg)
	return err
}

// handle commands
func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.cfg.Messages.Start)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	var check int
	row := b.botdb.Sql.QueryRow("SELECT COUNT(chat_id) FROM querys WHERE chat_id=$1", message.Chat.ID)
	if err := row.Scan(&check); err != nil {
		return err
	}
	if check != 0 {
		if _, err := b.botdb.Sql.Exec("UPDATE querys SET query=$1 WHERE chat_id=$2", pending, message.Chat.ID); err != nil {
			return err
		}
	} else {
		if _, err := b.botdb.Sql.Exec("INSERT INTO querys (chat_id, query) VALUES ($1, $2) ON CONFLICT (chat_id) DO NOTHING", message.Chat.ID, pending); err != nil {
			return err
		}
	}

	msg.ReplyMarkup = pendingKeyboard
	if _, err := b.bot.Send(msg); err != nil {
		return err
	}
	return nil
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.cfg.Messages.UnknownCommand)

	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleHotQuery(query *tgbotapi.CallbackQuery) (err error) {
	//msg := tgbotapi.NewMessage(query.Message.Chat.ID, "")
	var quizReply = []string{b.cfg.Messages.ChooseThick, b.cfg.Messages.ChooseA, b.cfg.Messages.ChooseCurl, b.cfg.Messages.ChooseDamage}
	quert := strings.Split(query.Data, ";")

	switch quert[0] {
	case requests:
		var checkauth int
		row := b.botdb.Sql.QueryRow("SELECT COUNT(id) FROM admins WHERE chat_id=$1", query.Message.Chat.ID)
		if err = row.Scan(&checkauth); err != nil {
			return err
		}
		if checkauth == 0 {
			msg := tgbotapi.NewMessage(query.Message.Chat.ID, b.cfg.Messages.NoAccess)
			msg.ReplyMarkup = pendingKeyboard
			if _, err = b.botdb.Sql.Exec("UPDATE querys SET query=$1 WHERE chat_id=$2", pending, query.Message.Chat.ID); err != nil {
				return err
			}
			if _, err = b.bot.Send(msg); err != nil {
				return err
			}
			break
		}

		var result storage.Request
		var status string
		row = b.botdb.Sql.QueryRow("SELECT * FROM requests WHERE user_chat_id=$1", quert[1])
		row.Scan(&result.Id, &result.TelegramName, &result.FirstName, &result.ChatID, &result.Status)
		if result.Status == nil {
			status = "Очікує"
		} else if *result.Status == true {
			status = "Підтверджено"
		} else if *result.Status == false {
			status = "Відхилено"
		}

		edit := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Запит №"+quert[2]+"\n\n"+
			"Ім'я: "+result.FirstName+"\n"+
			"Юзерайді: "+strconv.FormatInt(result.ChatID, 10)+"\n"+
			"Статус: "+status+"\n\n"+
			"Контактна інформація"+"\n"+
			"Нік: "+result.TelegramName+"\n")
		//edit.ParseMode = "Markdown"
		if _, err = b.bot.Send(edit); err != nil {
			return err
		}
		msg2 := tgbotapi.NewMessage(query.Message.Chat.ID, "[Якщо нік порожній](tg://user?id="+quert[1]+")")
		msg2.ParseMode = "Markdown"
		if _, err = b.bot.Send(msg2); err != nil {
			return err
		}

	case requi:
		var checkauth int
		row := b.botdb.Sql.QueryRow("SELECT COUNT(id) FROM admins WHERE chat_id=$1", query.Message.Chat.ID)
		if err = row.Scan(&checkauth); err != nil {
			return err
		}
		if checkauth == 0 {
			msg := tgbotapi.NewMessage(query.Message.Chat.ID, b.cfg.Messages.NoAccess)
			msg.ReplyMarkup = pendingKeyboard
			if _, err = b.botdb.Sql.Exec("UPDATE querys SET query=$1 WHERE chat_id=$2", pending, query.Message.Chat.ID); err != nil {
				return err
			}
			if _, err = b.bot.Send(msg); err != nil {
				return err
			}
			break
		}

		var result storage.Request
		row = b.botdb.Sql.QueryRow("SELECT * FROM requests WHERE user_chat_id=$1", quert[1])
		if err = row.Scan(&result.Id, &result.TelegramName, &result.FirstName, &result.ChatID, &result.Status); err != nil {
			return err
		}

		edit := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Учасник №"+quert[2]+"\n"+
			"Ім'я: "+result.FirstName+"\n"+
			"Юзерайді: "+strconv.FormatInt(result.ChatID, 10)+"\n"+
			"Нік: "+result.TelegramName+"\n")
		if _, err = b.bot.Send(edit); err != nil {
			return err
		}
		msg2 := tgbotapi.NewMessage(query.Message.Chat.ID, "[Якщо нік порожній](tg://user?id="+quert[1]+")")
		msg2.ParseMode = "Markdown"
		if _, err = b.bot.Send(msg2); err != nil {
			return err
		}
	case answ:
		var checkauth int
		var quertCounter int
		row := b.botdb.Sql.QueryRow("SELECT COUNT(id) FROM requests WHERE user_chat_id=$1 AND status = true", query.Message.Chat.ID)
		if err = row.Scan(&checkauth); err != nil {
			return err
		}
		if checkauth == 0 {
			msg := tgbotapi.NewMessage(query.Message.Chat.ID, b.cfg.Messages.NoAccess)
			msg.ReplyMarkup = pendingKeyboard
			if _, err = b.botdb.Sql.Exec("UPDATE querys SET query=$1 WHERE chat_id=$2", pending, query.Message.Chat.ID); err != nil {
				return err
			}
			if _, err = b.bot.Send(msg); err != nil {
				return err
			}
			break
		}

		quertCounter, err = strconv.Atoi(quert[1])

		if _, err = b.botdb.Sql.Exec("UPDATE querys SET statecounter = array_append(statecounter, $1) WHERE chat_id = $2", quert[2], query.Message.Chat.ID); err != nil {
			return err
		}

		if quertCounter != 3 {
			msgedit := tgbotapi.NewEditMessageReplyMarkup(query.Message.Chat.ID, query.Message.MessageID, nilKeyboard)
			if _, err = b.bot.Send(msgedit); err != nil {
				return err
			}

			txtedit := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Ваш вибір: *"+tools.Locale[quert[2]]+"*")
			txtedit.ParseMode = "Markdown"
			if _, err = b.bot.Send(txtedit); err != nil {
				return err
			}

			msg := tgbotapi.NewMessage(query.Message.Chat.ID, quizReply[quertCounter+1])
			msg.ReplyMarkup = inlineArr[quertCounter+1]
			if _, err = b.bot.Send(msg); err != nil {
				return err
			}
		} else if quertCounter == 3 {
			var counter []string
			var ans []string
			var final []string
			msgedit := tgbotapi.NewEditMessageReplyMarkup(query.Message.Chat.ID, query.Message.MessageID, nilKeyboard)
			if _, err = b.bot.Send(msgedit); err != nil {
				return err
			}

			txtedit := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Ваш вибір: *"+tools.Locale[quert[2]]+"*")
			txtedit.ParseMode = "Markdown"
			if _, err = b.bot.Send(txtedit); err != nil {
				return err
			}

			row = b.botdb.Sql.QueryRow("SELECT statecounter FROM querys WHERE chat_id=$1", query.Message.Chat.ID)
			if err = row.Scan((*pq.StringArray)(&counter)); err != nil {
				return err
			}
			if ans, err = tools.FindAnswer(counter, b.botdb); err != nil {
				return errUnlucky
			}
			if final, err = tools.Final(ans, counter, b.botdb, b.bot); err != nil {
				return errUnlucky
			}
			var finaltext string = "*Процедури:*\n"
			for i := 1; i <= len(final); i++ {
				if i%3 == 0 {
					finaltext += final[i-1] + " градусів\n\n"
				} else {
					finaltext += final[i-1] + ", "
				}
			}
			if _, err = b.botdb.Sql.Exec("UPDATE querys SET statecounter=$1 WHERE chat_id=$2", sql.NullString{}, query.Message.Chat.ID); err != nil {
				return err
			}
			msg := tgbotapi.NewMessage(query.Message.Chat.ID, finaltext)
			msg.ReplyMarkup = userKeyboard
			if _, err = b.botdb.Sql.Exec("UPDATE querys SET query=$1 WHERE chat_id=$2", userMenu, query.Message.Chat.ID); err != nil {
				return err
			}
			msg.ParseMode = "Markdown"
			if _, err = b.bot.Send(msg); err != nil {
				return err
			}
		}
	}

	//_, err := b.bot.Send(msg)
	return err
}
