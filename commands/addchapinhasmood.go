package commands

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"

	"net/http"
	"os"
	"reflect"
	"regexp"

	"time"

	"cloud.google.com/go/datastore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var (
	redirectURI = os.Getenv("APP_URL") + "/callback"
	auth        = spotify.NewAuthenticator(redirectURI, spotify.ScopePlaylistModifyPublic, spotify.ScopePlaylistModifyPrivate)
	ch          = make(chan *spotify.Client)
)

func AddChapinhasMood(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var answer string
	error := false
	message := update.Message
	string := ""
	if message.ReplyToMessage != nil && message.ReplyToMessage.Text != "" {
		string = message.ReplyToMessage.Text
	}
	if message.Text != "" {
		string = string + " " + message.Text
	}
	log.Printf("Message %s", string)
	var re, err = regexp.Compile(`https://open\.spotify\.com/track/([[:alnum:]]+)`)
	if err != nil {
		log.Printf("Failed to compile regex: %v", err)
		error = true
	}

	match := re.FindStringSubmatch(string)
	if (len(match)) == 0 {
		error = true
		answer = "Spotify link not found. Use the command /addchapinhasmood containing as a reply to a Spotify link."
	}

	spotifyUrl := match[1]
	log.Printf("You're requesting song " + spotifyUrl)

	ctx := context.Background()
	projectID := "geckobutler"

	datastoreClient, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		error = true
	}

	kind := "oauth2.Token"
	name := "spotifyToken"
	key := datastore.NameKey(kind, name, nil)
	log.Printf("Created key: %v", key)

	var token oauth2.Token

	err = datastoreClient.Get(ctx, key, &token)

	//Don't have any stored token, will have to obtain it
	//now
	if err != nil {

		state, err := GenerateRandomString(64)
		if err != nil {
			log.Printf("Couldn't generate random state")
			error = true
		}

		url := auth.AuthURL(state)
		log.Printf("Log in to spotify in the following url %v", url)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, url)
		msg.ReplyToMessageID = update.Message.MessageID

		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("Couldn't send authorization URI as message. Error: %v", err)
			error = true
		}

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
			ch <- &client
		})
	} else {
		client := auth.NewClient(&token)
		log.Printf("Login completed")
		go func() { ch <- &client }()
	}

	client := <-ch

	user, err := client.CurrentUser()
	if err != nil {
		log.Printf("Couldn't get current user. Error: %v", err)
		error = true
	} else {
		log.Printf("Current user %v", user.DisplayName)
	}

	playlist, err := client.GetPlaylist("7cwB93saz58vHF9NAOBBFk")
	if err != nil {
		log.Printf("Couldn't get playlist: %v", err)
		error = true
	}

	song, err := client.GetTrack(spotify.ID(spotifyUrl))
	if err != nil {
		log.Printf("Couldn't get song: %v", err)
		error = true
	}

	_, err = client.AddTracksToPlaylist(playlist.ID, song.ID)
	if err != nil {
		log.Printf("Couldn't add track to playlist: %v", err)
		error = true
	}

	if !error {
		answer = "Added track to playlist!"
	} else if answer == "" {
		answer = "An unexpected error ocurred."
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, answer)
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

type Chapinha struct {
	Chosen     string
	LastChosen time.Time
}

func ProximoChapinha(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	ctx := context.Background()
	projectID := "geckobutler"

	datastoreClient, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
	}

	//	query := datastore.NewQuery("Chapinha").Filter("Chosen = ", false).Order("LastChosen").KeysOnly()
	query := datastore.NewQuery("Chapinha").KeysOnly()

	var chapinhas []string

	_, err = datastoreClient.GetAll(ctx, query, &chapinhas)
	log.Printf("Chapinhas result", chapinhas)
	results := []interface{}{}

	if err == nil {
		for _, v := range chapinhas {
			results = append(results, tgbotapi.NewInlineQueryResultArticle(v, v, v))
		}
	}

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		Results:       results,
	}

	if _, err := bot.AnswerInlineQuery(inlineConf); err != nil {
		log.Println(err)
	}
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
