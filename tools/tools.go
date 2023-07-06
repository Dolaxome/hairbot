package tools

import (
	"errors"
	"github.com/Dolaxome/hair-bot/pkg/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lib/pq"
	"strconv"
)

var Locale = map[string]string{
	"thin":           "Тонке волосся",
	"medium":         "Середнє волосся",
	"thick":          "Густе волосся",
	"less a":         "Менше А-силует",
	"equal a":        "А-силует",
	"more a":         "Більше А-силует",
	"straight":       "Рівне",
	"curly":          "Кучеряве",
	"light film":     "Легка плівка",
	"thick film":     "Густа плівка",
	"fifthInit":      "П'ята стадія пошкодження (початкова)",
	"fifthFin":       "П'ята стадія пошкодження (кінцева)",
	"strong keratin": "Сильний кератин",
	"medium keratin": "Середній кератин",
	"light keratin":  "Легкий кератин",
	"botox":          "Ботокс",
	"nano":           "Нанопластика",
}

func FormRequestsInline(requests []storage.Request, qData string) (keyboard tgbotapi.InlineKeyboardMarkup) {
	var buttons [][]tgbotapi.InlineKeyboardButton

	for i, request := range requests {
		var butText string
		if request.Status == nil {
			butText = strconv.FormatInt(request.ChatID, 10) + " (очікує)"
		} else {
			butText = strconv.FormatInt(request.ChatID, 10)
		}

		butData := qData + ";" + strconv.FormatInt(request.ChatID, 10) + ";" + strconv.Itoa(i+1)

		button := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(butText, butData))

		buttons = append(buttons, button)
	}
	keyboard = tgbotapi.NewInlineKeyboardMarkup(buttons...)
	return keyboard
}

func FindAnswer(rim []string, botdb *storage.DB) (ans []string, err error) {
	row := botdb.Sql.QueryRow("SELECT key FROM answers WHERE $1 = ANY(params) AND $2 = ANY(params) AND $3 = ANY(params)", rim[0], rim[1], rim[2])
	if err = row.Scan((*pq.StringArray)(&ans)); err != nil {
		return nil, err
	}
	return ans, err
}
func Final(rim []string, rom []string, botdb *storage.DB, bot *tgbotapi.BotAPI) (ans []string, err error) {
	structu := rom[0]
	film := rom[3]
	for _, r := range rim {
		var all int
		var structC int
		var filmC int
		row := botdb.Sql.QueryRow("SELECT COUNT(key) FROM temperature WHERE key = $1 AND $2 = ANY(params) AND $3 = ANY(params) ", r, structu, film)
		if err = row.Scan(&all); err != nil {
			return nil, err
		}
		row = botdb.Sql.QueryRow("SELECT COUNT(key) FROM temperature WHERE key = $1 AND params = $2", r, pq.StringArray{structu})
		if err = row.Scan(&structC); err != nil {
			return nil, err
		}
		row = botdb.Sql.QueryRow("SELECT COUNT(key) FROM temperature WHERE key = $1 AND params = $2", r, pq.StringArray{film})
		if err = row.Scan(&filmC); err != nil {
			return nil, err
		}
		if all != 0 {
			var temp string
			row = botdb.Sql.QueryRow("SELECT temp FROM temperature WHERE key = $1 AND $2 = ANY(params) AND $3 = ANY(params) ", r, structu, film)
			if err = row.Scan(&temp); err != nil {
				return nil, err
			}
			ans = append(ans, Locale[r], Locale[film], temp)
		} else if structC != 0 {
			var temp string
			row = botdb.Sql.QueryRow("SELECT temp FROM temperature WHERE key = $1 AND params = $2", r, pq.StringArray{structu})
			if err = row.Scan(&temp); err != nil {
				return nil, err
			}
			ans = append(ans, Locale[r], Locale[film], temp)
		} else if filmC != 0 {
			var temp string
			row = botdb.Sql.QueryRow("SELECT temp FROM temperature WHERE key = $1 AND params = $2", r, pq.StringArray{film})
			if err = row.Scan(&temp); err != nil {
				return nil, err
			}
			ans = append(ans, Locale[r], Locale[film], temp)
		}
	}
	if ans == nil {
		return nil, errors.New("ddd")
	} else {
		return ans, err
	}
}
