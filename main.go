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

const debugID = -180312761

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
		if update.Message == nil && update.InlineQuery == nil && update.ChosenInlineResult == nil {
			continue
		}

		if update.ChosenInlineResult != nil {
			chosenInlineResult := update.ChosenInlineResult
			log.Printf("Chose chapinha: %v", chosenInlineResult.ResultID)
			continue
		}

		if update.InlineQuery != nil {
			inlinequery := update.InlineQuery
			log.Printf("Inline Query: %v", inlinequery.Query)
			switch {
			case strings.EqualFold(inlinequery.Query, "proximochapinha"):
				commands.ProximoChapinha(bot, update)
			}
			continue
		}

		if update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup() {
			log.Printf("Request from chat: %s %d", update.Message.Chat.Title, update.Message.Chat.ID)
		} else {
			log.Printf("Request from user: %s %d", update.Message.Chat.UserName, update.Message.Chat.ID)
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

		case strings.EqualFold(command, "paywall"):
			commands.Paywall(bot, update)

		case strings.EqualFold(command, "scihub"):
			commands.SciHub(bot, update)

		case strings.EqualFold(command, "addchapinhasmood"):
			commands.AddChapinhasMood(bot, update)

		case strings.EqualFold(command, "debug") && update.Message.Chat.ID == debugID:
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
