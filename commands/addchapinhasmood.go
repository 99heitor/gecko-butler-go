package commands

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/99heitor/gecko-butler-go/spotify"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

const projectID = "geckobutler"

func replyMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, text string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

func AddChapinhasMood(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	message := update.Message

	string := ""
	if message.ReplyToMessage != nil {
		string = message.ReplyToMessage.Text
	}
	string += string + " " + message.Text
	log.Printf("Message %s", string)

	spotifyID, err := spotify.GetSpotifyId(string)
	if err != nil {
		replyMessage(bot, message, "Spotify link not found.")
		return
	}
	log.Printf("You're requesting song " + spotifyID.String())

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	datastoreClient, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		replyMessage(bot, message, "An unexpected error ocurred.")
		return
	}

	client, err := spotify.GetSpotifyClient(ctx, bot, message, datastoreClient)
	if err != nil {
		log.Printf("Couldn't get client. Error: %v", err)
		replyMessage(bot, message, "An unexpected error ocurred.")
		return
	}

	user, err := client.CurrentUser()
	if err != nil {
		log.Printf("Couldn't get current user. Error: %v", err)
		replyMessage(bot, message, "An unexpected error ocurred.")
		return
	}
	log.Printf("Current user %v", user.DisplayName)

	playlist, err := client.GetPlaylist("7cwB93saz58vHF9NAOBBFk")
	if err != nil {
		log.Printf("Couldn't get playlist: %v", err)
		replyMessage(bot, message, "An unexpected error ocurred.")
		return
	}

	song, err := client.GetTrack(spotifyID)
	if err != nil {
		log.Printf("Couldn't get song: %v", err)
		replyMessage(bot, message, "An unexpected error ocurred.")
		return
	}

	_, err = client.AddTracksToPlaylist(playlist.ID, song.ID)
	if err != nil {
		log.Printf("Couldn't add track to playlist: %v", err)
		replyMessage(bot, message, "An unexpected error ocurred.")
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
