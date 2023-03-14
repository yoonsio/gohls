package downloader

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/etherlabsio/go-m3u8/m3u8"
)

type SegmentWriter struct {
	*http.Client
	*m3u8.SegmentItem
	BaseURL string
}

func (c *SegmentWriter) Write(w io.Writer) error {
	if c.Client != nil {
		c.Client = http.DefaultClient
	}
	url := c.BaseURL + c.Segment
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return errors.New(resp.Status)
	}
	written, err := io.Copy(w, resp.Body)
	if err != nil {
		return err
	}
	log.Printf("downloaded %d bytes", written)
	return nil
}

type Segment struct {
	m3u8.SegmentItem
	Sequence int
}
