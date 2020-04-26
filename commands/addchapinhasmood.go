package commands

import (
	"context"
	"log"

	"github.com/99heitor/gecko-butler-go/spotify"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func replyMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, text string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

func showError(text string, bot *tgbotapi.BotAPI, message *tgbotapi.Message, err error) {
	log.Printf(text+" Error: %v", err)
	replyMessage(bot, message, "An unexpected error ocurred.")
	return
}

func AddChapinhasMood(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	message := update.Message

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
		showError("Couldn't get spotify client.", bot, message, err)
		return
	}

	user, err := client.CurrentUser()
	if err != nil {
		showError("Couldn't get current user.", bot, message, err)
		return
	}
	log.Printf("Current user %v", user.DisplayName)

	spotifyID, err := spotify.GetSpotifyId(client, string)
	if err != nil {
		replyMessage(bot, message, "Spotify link not found.")
		return
	}
	log.Printf("You're requesting song " + spotifyID.String())

	song, err := client.GetTrack(spotifyID)
	if err != nil {
		showError("Couldn't get current song.", bot, message, err)
		return
	}

	playlist, err := client.GetPlaylist("7cwB93saz58vHF9NAOBBFk")
	if err != nil {
		showError("Couldn't get playlist.", bot, message, err)
		return
	}

	_, err = client.AddTracksToPlaylist(playlist.ID, song.ID)
	if err != nil {
		showError("Couldn't add track to playlist.", bot, message, err)
		return
	}

	replyMessage(bot, message, "Added track to playlist!")
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
