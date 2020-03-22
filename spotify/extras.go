package spotify

import (
	"errors"
	"github.com/zmb3/spotify"
	"regexp"
)

func GetSpotifyId(message_body string) (spotify.ID, error) {
	re := regexp.MustCompile(
		`https://open\.spotify\.com/track/([[:alnum:]]+)`,
	)

	match := re.FindStringSubmatch(message_body)
	if (len(match)) == 0 {
		return spotify.ID(""), errors.New("Spotify link not found")
	}
	return spotify.ID(match[1]), nil
}
