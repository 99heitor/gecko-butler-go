package spotify

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"os"
	"reflect"

	"cloud.google.com/go/datastore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var (
	redirectURI = os.Getenv("APP_URL") + "/callback"
	auth        = spotify.NewAuthenticator(redirectURI, spotify.ScopePlaylistModifyPublic, spotify.ScopePlaylistModifyPrivate)
)

const projectID = "geckobutler"

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func GetToken(ctx context.Context, bot *tgbotapi.BotAPI, message *tgbotapi.Message) (*oauth2.Token, error) {
	if len(os.Getenv("SPOTIFY_TOKEN")) > 0 {
		token := oauth2.Token{}
		token.AccessToken = os.Getenv("SPOTIFY_TOKEN")
		return &token, nil
	}

	datastoreClient, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Printf("Failed to get token: %v", err)
		return nil, err
	}

	key := datastore.NameKey("oauth2.Token", "spotifyToken", nil)

	var token oauth2.Token
	err = datastoreClient.Get(ctx, key, &token)

	//Don't have any stored token, will have to obtain it now
	if err != nil {
		state, err := generateRandomString(64)
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

		resultChannel := make(chan *oauth2.Token)
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
			resultChannel <- tok
		})
		return (<-resultChannel), nil
	}
	return &token, nil
}

func GetClient(ctx context.Context, bot *tgbotapi.BotAPI, message *tgbotapi.Message) (*spotify.Client, error) {
	token, err := GetToken(ctx, bot, message)
	client := auth.NewClient(token)
	return &client, err
}
