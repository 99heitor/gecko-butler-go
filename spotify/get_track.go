package spotify

import (
	"errors"
	"log"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zmb3/spotify"
)

func GetTrackIdFromLink(client *spotify.Client, message *tgbotapi.Message) (spotify.ID, error) {
	combinedTexts := ""
	if message.ReplyToMessage != nil {
		combinedTexts = message.ReplyToMessage.Text
	}
	combinedTexts += " " + message.Text

	re := regexp.MustCompile(
		`https://open\.spotify\.com/track/([[:alnum:]]+)`,
	)

	match := re.FindStringSubmatch(combinedTexts)
	if (len(match)) == 0 {
		return spotify.ID(""), errors.New("Spotify link not provided")
	}
	return spotify.ID(match[1]), nil
}

func GetTrackIdByQuery(client *spotify.Client, messageText string) (spotify.ID, error) {
	query := messageText[18:]
	if len(query) > 0 {
		limit := 1
		options := spotify.Options{Limit: &limit}
		searchResult, err := client.SearchOpt(query, spotify.SearchTypeTrack, &options)
		if err != nil {
			log.Printf("Couldn't search for a track. Error: %v", err)
		}
		if searchResult.Tracks != nil && len(searchResult.Tracks.Tracks) > 0 {
			trackID := searchResult.Tracks.Tracks[0].ID
			return spotify.ID(trackID), nil
		} else {
			return spotify.ID(""), errors.New("No track found")
		}
	}
	return spotify.ID(""), errors.New("Search query not provided")
}

func GetTrackId(client *spotify.Client, message *tgbotapi.Message) (spotify.ID, error) {
	trackID, err := GetTrackIdFromLink(client, message)
	if err != nil {
		trackID, err = GetTrackIdByQuery(client, message.Text)
		return trackID, err
	}
	return trackID, err
}
