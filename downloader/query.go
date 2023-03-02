package downloader

import (
	"errors"
	"log"
	"net/http"

	"github.com/etherlabsio/go-m3u8/m3u8"
)

func queryPlaylist(client *http.Client, url string) (*m3u8.Playlist, error) {

	log.Printf("q playlist: %s", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return nil, errors.New(resp.Status)
	}
	return m3u8.Read(resp.Body)
}
