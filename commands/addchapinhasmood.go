package commands

import (
	"context"
	"log"
	"os"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
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

		config := &clientcredentials.Config{
			ClientID:     os.Getenv("SPOTIFY_ID"),
			ClientSecret: os.Getenv("SPOTIFY_SECRET"),
			TokenURL:     spotify.TokenURL,
		}

		token, err := config.Token(context.Background())
		if err != nil {
			log.Printf("Couldn't get token: %v", err)
			break
		}

		client := spotify.Authenticator{}.NewClient(token)

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
