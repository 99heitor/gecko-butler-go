package commands

import (
	"fmt"
	"html"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/99heitor/gecko-butler-go/smmry"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var smmryClient = smmry.Client{Token: os.Getenv("SMMRY_KEY")}

//Summarize answers the command with a smmry API response
func Summarize(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var answer string

	message := update.Message
	if message.ReplyToMessage == nil {
		answer = "Use the command as a reply to a message containing a URL."
	} else if url := getFirstURL(message.ReplyToMessage.Text); url == "" {
		answer = "Doesn't look like anything to me."
	} else {
		var summaryLength int
		if args := update.Message.CommandArguments(); args != "" {
			argList := strings.Fields(args)
			if len(argList) > 0 {
				if val, err := strconv.Atoi(argList[0]); err == nil {
					summaryLength = val
				}
			}
		}
		params := smmry.Params{URL: url, Length: summaryLength}
		smmryResponse, err := smmryClient.GetSummary(params)
		if err != nil || smmryResponse.Content == "" {
			answer = "Sorry, I can't do that right now."
			log.Printf("Error trying to summarize: err: %v, API code: %d, API message: %s",
				err, smmryResponse.Error, smmryResponse.Message)
		} else {
			sentences := strings.Split(smmryResponse.Content, "[BREAK]")
			var body []string
			for _, sentence := range sentences {
				body = append(body, strings.TrimSpace(sentence))
			}

			title := html.UnescapeString(decodeTrash(smmryResponse.Title))
			text := html.UnescapeString(strings.Join(body, "\n"))
			log.Printf("Answering with summary \"%s\"...", title)
			answer = fmt.Sprintf("@%s is too lazy to read.\n\n*%s*\n%s",
				update.Message.From.UserName, title, text)
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, decodeTrash(answer))
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = "Markdown"
	bot.Send(msg)

}

func getFirstURL(message string) string {
	words := strings.Fields(message)
	for _, word := range words {
		if uri, err := url.ParseRequestURI(word); err == nil {
			switch uri.Scheme {
			case "http":
			case "https":
			default:
				continue
			}
			return word
		}
	}
	return ""
}

//This API is kind of fucked up with special characters
func decodeTrash(text string) string {
	return strings.NewReplacer("&#2013265929;", "é", "\\", " ").Replace(text)
}
