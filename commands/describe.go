package commands

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//Describe handles the "/describe" bot command
func Describe(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var answer = "Doesn't look like anything to me."
	message := update.Message
	switch {
	case message.ReplyToMessage == nil:
		answer = "Use the command /describe as a reply to a picture."
	case message.ReplyToMessage.Photo == nil:
		answer = "Sorry, this is not a picture."
	default:
		photos := *message.ReplyToMessage.Photo
		lastPhotoID := tgbotapi.FileConfig{FileID: photos[len(photos)-1].FileID}
		file, err := bot.GetFile(lastPhotoID)
		if err != nil {
			log.Printf("Failed to get image link from Telegram: %v", err)
		}
		answers, err := getLabels(file.Link(bot.Token))
		if err != nil {
			log.Printf("Failed to get labels: %v", err)
		} else if len(answers) > 0 {
			log.Printf("Sending response %s...", answers[0])
			answer = "I see... \n" + strings.Join(answers, "\n")
		}
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, answer)
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func getLabels(photoURI string) ([]string, error) {
	var answers []string
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	resp, err := http.Get(photoURI)
	if err != nil {
		return nil, err
	}
	image, err := vision.NewImageFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	labels, err := client.DetectLabels(ctx, image, nil, 10)
	if err != nil {
		return nil, err
	}

	for _, label := range labels {
		answers = append(answers, fmt.Sprintf("*%s* : %.1f%%", label.Description, label.Score*100))
	}
	return answers, nil
}
