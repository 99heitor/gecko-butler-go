package spotify

import (
	"errors"
	"log"
	"regexp"

	"github.com/zmb3/spotify"
)

func GetTrackId(client *spotify.Client, message_body string) (spotify.ID, error) {
	re := regexp.MustCompile(
		`https://open\.spotify\.com/track/([[:alnum:]]+)`,
	)

	match := re.FindStringSubmatch(message_body)
	if (len(match)) == 0 {
		return GetTrackIdByQuery(client, message_body)
	}
	return spotify.ID(match[1]), nil
}

func GetTrackIdByQuery(client *spotify.Client, message_body string) (spotify.ID, error) {
	query := message_body[18:]
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
		}
	}
	return spotify.ID(""), errors.New("Couldn't get track ID.")
}
