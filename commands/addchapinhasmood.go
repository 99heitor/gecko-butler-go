package commands

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"

	"net/http"
	"os"
	"reflect"
	"regexp"

	"cloud.google.com/go/datastore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var (
	redirectURI = os.Getenv("APP_URL") + "/callback"
	auth        = spotify.NewAuthenticator(redirectURI, spotify.ScopePlaylistModifyPublic, spotify.ScopePlaylistModifyPrivate)
)

func getSpotifyId(message_body string) (spotify.ID, error) {
	re := regexp.MustCompile(
		`https://open\.spotify\.com/track/([[:alnum:]]+)`,
	)

	match := re.FindStringSubmatch(message_body)
	if (len(match)) == 0 {
		return spotify.ID(""), errors.New("Spotify link not found")
	}
	return spotify.ID(match[1]), nil
}

const projectID = "geckobutler"

func getSpotifyClient(ctx context.Context, bot *tgbotapi.BotAPI, message *tgbotapi.Message, datastoreClient *datastore.Client) (*spotify.Client, error) {
	key := datastore.NameKey("oauth2.Token", "spotifyToken", nil)

	var token oauth2.Token
	err := datastoreClient.Get(ctx, key, &token)

	//Don't have any stored token, will have to obtain it now
	if err != nil {
		state, err := GenerateRandomString(64)
		if err != nil {
			log.Printf("Couldn't generate random state")
			return nil, errors.New("Couldn't generate random state")
		}

		url := auth.AuthURL(state)
		log.Printf("Log in to spotify in the following url %v", url)

		msg := tgbotapi.NewMessage(message.Chat.ID, url)
		msg.ReplyToMessageID = message.MessageID

		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("Couldn't send authorization URI as message. Error: %v", err)
			return nil, err
		}

		resultChannel := make(chan *spotify.Client)
		http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
			tok, err := auth.Token(state, r)
			if err != nil {
				log.Printf("Couldn't get token. Error %v", err)
				http.Error(w, "Couldn't get token", http.StatusForbidden)
				return
			}

			_, err = datastoreClient.Put(ctx, key, tok)
			log.Printf("Token has type %v", reflect.TypeOf(key).Kind())
			if err != nil {
				log.Printf("Failed storing token %v with error: %v", tok, err)
				http.Error(w, "Couldn't store token", http.StatusForbidden)
				return
			}
			log.Printf("Stored token with key %v", key)

			if st := r.FormValue("state"); st != state {
				log.Printf("State mismatch. Error %v", err)
				http.NotFound(w, r)
				return
			}

			client := auth.NewClient(tok)
			log.Printf("Login completed")
			resultChannel <- &client
		})
		return (<-resultChannel), nil
	} else {
		client := auth.NewClient(&token)
		log.Printf("Login completed")
		return client, nil
	}
}

func repluMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, text string) {
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

	spotifyID, err := getSpotifyId(string)
	if err != nil {
		repluMessage(bot, message, "Spotify link not found.")
		return
	}
	log.Printf("You're requesting song " + spotifyID.String())

	ctx := context.Background()
	datastoreClient, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		repluMessage(bot, message, "An unexpected error ocurred.")
		return
	}

	client, err := getSpotifyClient(ctx, bot, message, datastoreClient)
	if err != nil {
		log.Printf("Couldn't get client. Error: %v", err)
		repluMessage(bot, message, "An unexpected error ocurred.")
		return
	}

	user, err := client.CurrentUser()
	if err != nil {
		log.Printf("Couldn't get current user. Error: %v", err)
		repluMessage(bot, message, "An unexpected error ocurred.")
		return
	}
	log.Printf("Current user %v", user.DisplayName)

	playlist, err := client.GetPlaylist("7cwB93saz58vHF9NAOBBFk")
	if err != nil {
		log.Printf("Couldn't get playlist: %v", err)
		repluMessage(bot, message, "An unexpected error ocurred.")
		return
	}

	song, err := client.GetTrack(spotifyID)
	if err != nil {
		log.Printf("Couldn't get song: %v", err)
		repluMessage(bot, message, "An unexpected error ocurred.")
		return
	}

	_, err = client.AddTracksToPlaylist(playlist.ID, song.ID)
	if err != nil {
		log.Printf("Couldn't add track to playlist: %v", err)
		repluMessage(bot, message, "An unexpected error ocurred.")
		return
	}

	repluMessage(bot, message, "Added track to playlist!")
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

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}
