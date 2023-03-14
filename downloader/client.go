package downloader

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/etherlabsio/go-m3u8/m3u8"
)

const (
	DefaultInterval           = 2 * time.Second
	DefaultPlaylistDepthLimit = 5
)

type Client struct {

	// http client
	client *http.Client

	// store for segments
	store Store

	// interval to refresh playlist
	interval time.Duration

	// goroutine worker pool
	pool *Pool

	// playlist depth limit to avoid deadlock
	plDepthLimit int

	// internal set to track streams seen
	pm *sync.Map
}

func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		client:       http.DefaultClient,
		store:        &NopStore{},
		interval:     DefaultInterval,
		plDepthLimit: DefaultPlaylistDepthLimit,
		pm:           &sync.Map{},
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.pool == nil {
		c.pool = NewPool()
	}
	return c
}

// Download accepts playlist url and continuously downloads stream
// until the context is expired
func (c *Client) Download(ctx context.Context, urls []string) error {

	log.Printf("downloading: %+v", urls)

	ticker := time.NewTicker(c.interval)
	for {
		select {
		case <-ctx.Done():
			log.Printf("Download: ticker loop cancelled")
			// wait for all existing jobs to finish
			if err := c.pool.Close(); err != nil {
				return fmt.Errorf("pool error: %w", err)
			}
			return ctx.Err()
		case <-ticker.C:
			// this can be potentially block if buffer is filled
			for _, url := range urls {
				c.pool.Do(func() error {
					baseURL, resource := parseURL(url)
					return c.refresh(baseURL, resource, 0)
				})
			}

		}
	}
}

func parseURL(url string) (string, string) {
	resource := path.Base(url)
	baseURL := strings.TrimSuffix(url, resource)
	return baseURL, resource
}

// recursively retrieve playlist
func (c *Client) refresh(baseURL, resource string, depth int) error {
	if depth >= c.plDepthLimit {
		return fmt.Errorf("refresh: nested playlist with > %d depth found", depth)
	}
	pl, err := queryPlaylist(c.client, baseURL+resource)
	if err != nil {
		return fmt.Errorf("refresh: failed to query playlist %s: %w", resource, err)
	}

	for _, plItem := range pl.Playlists() {
		if _, ok := c.pm.LoadOrStore(plItem.URI, struct{}{}); ok {
			continue
		}
		log.Printf("refresh: found playlist: %s", plItem.URI)
		c.pool.Do(func() error {
			return c.refresh(baseURL, plItem.URI, depth)
		})
	}

	for i, segment := range pl.Segments() {
		// if sequence is found, use it to identify the segments
		sequence := pl.Sequence + i
		segmentKey := segment.Segment
		if pl.Sequence != 0 {
			segmentKey = strconv.Itoa(sequence)
		}
		if _, ok := c.pm.LoadOrStore(segmentKey, struct{}{}); ok {
			continue
		}
		log.Printf("refresh: found segment: (%s) %s", segmentKey, segment.Segment)
		c.pool.Do(func() error {
			return c.downloadSegment(sequence, baseURL, segment)
		})
	}
	return nil
}

func (c *Client) downloadSegment(sequence int, baseURL string, segment *m3u8.SegmentItem) error {

	log.Printf("segment %d (%f): %s\n", sequence, segment.Duration, segment.Segment)

	w, err := c.store.GetWriter(&Segment{
		SegmentItem: *segment,
		Sequence:    sequence,
	})
	if err != nil {
		return fmt.Errorf("downloadSegment: failed to get store writer: %w", err)
	}
	defer w.Close()
	return (&SegmentWriter{
		Client:      c.client,
		SegmentItem: segment,
		BaseURL:     baseURL,
	}).Write(w)
}

type ClientOption func(*Client)

func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.client = client
	}
}

func WithStore(store Store) ClientOption {
	return func(c *Client) {
		c.store = store
	}
}

func WithInterval(interval time.Duration) ClientOption {
	return func(c *Client) {
		c.interval = interval
	}
}

func WithPool(pool *Pool) ClientOption {
	return func(c *Client) {
		c.pool = pool
	}
}

func WithPlaylistDepthLimit(limit int) ClientOption {
	return func(c *Client) {
		c.plDepthLimit = limit
	}
}
