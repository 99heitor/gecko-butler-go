package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//Paywall handles the "/paywall" bot command
func Paywall(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	appender(bot, update, "https://outline.com")
}

//SciHub handles the "/scihub" bot command
func SciHub(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	appender(bot, update, "http://sci-hub.tw")
}

func appender(bot *tgbotapi.BotAPI, update tgbotapi.Update, base_url string) {
	message := update.Message
	if message.ReplyToMessage == nil {
		return
	}
	link := getFirstURL(message.ReplyToMessage.Text)
	if len(link) != 0 {
		reply := fmt.Sprintf("🏴‍☠️ [Link sem Paywall](%s/%s) 🏴‍☠️", base_url, link)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ParseMode = "Markdown"
		msg.DisableWebPagePreview = true
		bot.Send(msg)
	}
}
