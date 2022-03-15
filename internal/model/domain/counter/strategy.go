package counter

import (
	"errors"
	"io"
	"net/http"
	"occurrence-calculator/internal/model/infrastructure/cli"
	"os"
)

type SourceStrategy interface {
	Read(source string) (io.ReadCloser, error)
}

type FileSourceStrategy struct {
}

func (s *FileSourceStrategy) Read(source string) (io.ReadCloser, error) {
	file, err := os.Open(source)
	if err != nil {
		return nil, err
	}

	return file, nil
}

type UrlSourceStrategy struct {
}

func (s *UrlSourceStrategy) Read(source string) (io.ReadCloser, error) {
	response, err := http.Get(source)

	if err != nil {
		return nil, err
	}

	return response.Body, err
}

func MakeSourceStrategyBySourceType(source string) (SourceStrategy, error) {
	var strategy SourceStrategy
	var err error

	switch source {
	case cli.FileSource:
		strategy, err = &FileSourceStrategy{}, nil
	case cli.UrlSource:
		strategy, err = &UrlSourceStrategy{}, nil
	default:
		strategy, err = nil, errors.New("strategy not found")
	}

	return strategy, err
}
