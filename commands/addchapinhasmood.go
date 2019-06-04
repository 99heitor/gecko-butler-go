package commands

import (
	// "context"
	"log"
	"net/http"
	"os"
	"regexp"
	"crypto/rand"
	"encoding/base64"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zmb3/spotify"
	// "google.golang.org/appengine/datastore"
)

var (
	redirectURI = os.Getenv("APP_URL") + "/callback"
	auth        = spotify.NewAuthenticator(redirectURI, spotify.ScopePlaylistModifyPublic, spotify.ScopePlaylistModifyPrivate)
	ch          = make(chan *spotify.Client)
)

func AddChapinhasMood(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var answer = "Doesn't look like a song to me."
	message := update.Message
	switch {
	case message.ReplyToMessage == nil:
		answer = "Use the command /addchapinhasmood as a reply to a Spotify link."
	case message.ReplyToMessage.Text == "":
		answer = "Sorry, this is not a Spotify link."
	default:
		var re, err = regexp.Compile(`https://open\.spotify\.com/track/([[:alnum:]]+)`)
		if err != nil {
			log.Printf("Failed to compile regex: %v", err)
			return
		}
		text := message.ReplyToMessage.Text
		match := re.FindStringSubmatch(text)
		if (len(match)) == 0 {
			log.Printf("Didn't find a Spotify url.")
			break
 		}
		spotifyUrl := match[1]
		answer = "You're requesting song " + spotifyUrl

		state, err := GenerateRandomString(64)
		if err != nil {
			log.Printf("Couldn't generate random state")
		}

		url := auth.AuthURL(state)
		log.Printf("Log in to spotify in the following url %v", url)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, url)
		msg.ReplyToMessageID = update.Message.MessageID

		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("Couldn't send authorization URI as message. Error: %v", err)
		}

		http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
			tok, err := auth.Token(state, r)
			if err != nil {
				log.Printf("Couldn't get token. Error %v", err)
				http.Error(w, "Couldn't get token", http.StatusForbidden)
				return
			}

			// ctx := context.Background()

			// key := datastore.NewKey(ctx, "*oauth2.Token", "token", 0, nil)
			// log.Printf("Created key: %v", key)
			// _, err = datastore.Put(ctx, key, &tok)
			// if err != nil {
			// 	log.Printf("Failed storing token %v with error: %v", tok, err)
			// 	http.Error(w, "Couldn't store token", http.StatusForbidden)
			// 	return
			// }
			// log.Printf("Stored token with key %v", key)

			if st := r.FormValue("state"); st != state {
				log.Printf("State mismatch. Error %v", err)
				http.NotFound(w, r)
				return
			}

			client := auth.NewClient(tok)
			log.Printf("Login completed")
			ch <- &client
		});

		client := <-ch

		user, err := client.CurrentUser()
		if err != nil {
			log.Printf("Couldn't get current user. Error: %v", err)
		} else {
			log.Printf("Current user %v", user.DisplayName)
		}

		playlist, err := client.GetPlaylist("7cwB93saz58vHF9NAOBBFk")
		if err != nil {
			log.Printf("Couldn't get playlist: %v", err)
			break
		}

		song, err := client.GetTrack(spotify.ID(spotifyUrl))
		if err != nil {
			log.Printf("Couldn't get song: %v", err)
			break
		}

		_, err = client.AddTracksToPlaylist(playlist.ID, song.ID)
		if err != nil {
			log.Printf("Couldn't add track to playlist: %v", err)
			break
		}

		answer = "Added track to playlist!"
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, answer)
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = "Markdown"
	bot.Send(msg)
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