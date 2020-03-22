package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
)

const greeting = "Hi! My hash is `%s`\nYou can browse my source code [here](https://github.com/99heitor/gecko-butler-go/tree/%s)"

func Version(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	hash := os.Getenv("COMMIT_HASH")
	var reply string
	if len(hash) > 0 {
		reply = fmt.Sprintf(greeting, hash, hash)
	} else {
		reply = "Hi! I'm a test version."
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true
	bot.Send(msg)
}
