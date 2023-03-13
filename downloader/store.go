package downloader

import (
	"io"
	"os"
	"strings"

	"github.com/etherlabsio/go-m3u8/m3u8"
)

type Store interface {
	// GetWriter returns writecloser for given storage
	// caller is expected to close the writer
	GetWriter(segment *m3u8.SegmentItem) (io.WriteCloser, error)

	// Process combines segments into a single file
	Process() error
}

type discardCloser struct{}

func (discardCloser) Write(p []byte) (int, error) {
	return len(p), nil
}

func (discardCloser) Close() error {
	return nil
}

type NopStore struct {
}

func (s *NopStore) GetWriter(segment *m3u8.SegmentItem) (io.WriteCloser, error) {
	return discardCloser{}, nil
}

func (s *NopStore) Process() error {
	return nil
}

type LocalStore struct {
	dir string
}

func NewLocalStorage(dir string) (*LocalStore, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	return &LocalStore{
		dir: dir,
	}, nil
}

func (s *LocalStore) GetWriter(segment *m3u8.SegmentItem) (io.WriteCloser, error) {
	return os.Create(s.dir + "/" + genFileName(segment))
}

func (s *LocalStore) Process() error {
	// TODO: stitch all ts files together under this directory
	return nil
}

func genFileName(segment *m3u8.SegmentItem) string {
	// strip url parameters
	fileName := strings.Split(segment.Segment, "?")[0]

	// TODO: append timestamp

	return fileName
}
