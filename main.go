package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/99heitor/gecko-butler-go/commands"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot := setupBot()
	updates := bot.ListenForWebhook("/" + bot.Token)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	go http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup() {
			log.Printf("Request from chat: %s", update.Message.Chat.Title)
		} else {
			log.Printf("Request from user: %s", update.Message.Chat.UserName)
		}
		if bot.Debug {
			log.Printf("Update: %v", update.Message.Text)
		}

		command := update.Message.Command()

		switch {

		case strings.EqualFold(command, "describe"):
			commands.Describe(bot, update)

		case strings.EqualFold(command, "tldr"):
			commands.Summarize(bot, update)

		case strings.EqualFold(command, "debug") && update.Message.Chat.ID == 36992723:
			rsp := fmt.Sprintf("Switching debug mode to %t", !bot.Debug)
			log.Printf(rsp)
			bot.Debug = !bot.Debug

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, rsp)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}

func setupBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(os.Getenv("APP_URL") + "/" + bot.Token))
	if err != nil {
		log.Fatal(err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	return bot
}
