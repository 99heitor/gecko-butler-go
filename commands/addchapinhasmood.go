package commands

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/99heitor/gecko-butler-go/spotify"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func replyMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, text string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyToMessageID = message.MessageID
	msg.ParseMode = tgbotapi.ModeMarkdown
	bot.Send(msg)
}

func showError(bot *tgbotapi.BotAPI, message *tgbotapi.Message, err error, errorText string, replyText string) {
	log.Printf(errorText+" Error: %v", err)
	if replyText == "" {
		replyText = "An unexpected error ocurred."
	}
	replyMessage(bot, message, replyText)
	return
}

func AddChapinhasMood(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	message := update.Message
	chatID := strconv.FormatInt(message.Chat.ID, 10)

	if chatID != os.Getenv("ALLOWED_CHAT_ID") {
		errorText := "Chat not allowed"
		showError(bot, message, errors.New(errorText), errorText, errorText)
		return
	}

	string := ""
	if message.ReplyToMessage != nil {
		string = message.ReplyToMessage.Text
	}
	string += " " + message.Text
	log.Printf("Message %s", string)

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	client, err := spotify.GetClient(ctx, bot, message)
	if err != nil {
		showError(bot, message, err, "Couldn't get spotify client.", "")
		return
	}

	user, err := client.CurrentUser()
	if err != nil {
		showError(bot, message, err, "Couldn't get current user.", "")
		return
	}
	log.Printf("Current user %v", user.DisplayName)

	trackID, err := spotify.GetTrackId(client, string)
	if err != nil {
		showError(bot, message, err, "Couldn't get track id.", "Couldn't get track")
		return
	}
	log.Printf("You're requesting track " + trackID.String())

	track, err := client.GetTrack(trackID)
	if err != nil {
		showError(bot, message, err, "Couldn't get track from id.", "Couldn't get track")
		return
	}

	playlist, err := client.GetPlaylist("7cwB93saz58vHF9NAOBBFk")
	if err != nil {
		showError(bot, message, err, "Couldn't get playlist.", "")
		return
	}

	_, err = client.AddTracksToPlaylist(playlist.ID, track.ID)
	if err != nil {
		showError(bot, message, err, "Couldn't add track to playlist.", "")
		return
	}

	trackURL := track.ExternalURLs["spotify"]
	playlistURL := playlist.ExternalURLs["spotify"]

	replyMessage(bot, message, "Successfully added ["+track.Name+"]("+trackURL+") to ["+playlist.Name+"]("+playlistURL+")! ✨")
}

// type Chapinha struct {
// 	Chosen     string
// 	LastChosen time.Time
// }

func ProximoChapinha(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// 	ctx := context.Background()

	// 	datastoreClient, err := datastore.NewClient(ctx, projectID)
	// 	if err != nil {
	// 		log.Printf("Failed to create client: %v", err)
	// 	}

	// 	//	query := datastore.NewQuery("Chapinha").Filter("Chosen = ", false).Order("LastChosen").KeysOnly()
	// 	query := datastore.NewQuery("Chapinha").KeysOnly()

	// 	var chapinhas []string

	// 	_, err = datastoreClient.GetAll(ctx, query, &chapinhas)
	// 	log.Printf("Chapinhas result", chapinhas)
	// 	results := []interface{}{}

	// 	if err == nil {
	// 		for _, v := range chapinhas {
	// 			results = append(results, tgbotapi.NewInlineQueryResultArticle(v, v, v))
	// 		}
	// 	}

	// 	inlineConf := tgbotapi.InlineConfig{
	// 		InlineQueryID: update.InlineQuery.ID,
	// 		Results:       results,
	// 	}

	// 	if _, err := bot.AnswerInlineQuery(inlineConf); err != nil {
	// 		log.Println(err)
	// 	}
}
