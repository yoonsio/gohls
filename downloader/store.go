package downloader

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Store interface {
	// GetWriter returns writecloser for given storage
	// caller is expected to close the writer
	GetWriter(segment *Segment) (io.WriteCloser, error)

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

func (s *NopStore) GetWriter(segment *Segment) (io.WriteCloser, error) {
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

func (s *LocalStore) GetWriter(segment *Segment) (io.WriteCloser, error) {
	return os.Create(s.dir + "/" + genFileName(segment))
}

func (s *LocalStore) Process() error {
	// TODO: stitch all ts files together under this directory
	return nil
}

func genFileName(segment *Segment) string {
	// strip url parameters
	// TODO: append timestamp
	fileName := fmt.Sprintf("%d-$s", segment.Sequence, strings.Split(segment.Segment, "?")[0])
	return fileName
}
