package cli

import (
	"errors"
	"fmt"
)

type Command interface {
	Execute()
	AddFlag(name string, value interface{})
	GetFlag(name string) interface{}
}

const (
	FlagType = "type"

	FileSource = "file"
	UrlSource  = "url"
)

type Source string

func (i *Source) String() string {
	return fmt.Sprint(*i)
}

func (i *Source) Set(value string) error {
	var newValue string

	switch value {
	case FileSource:
		newValue = FileSource
	case UrlSource:
		newValue = UrlSource
	default:
		return errors.New("source type is not present, available type: url, file")
	}

	*i = Source(newValue)

	return nil
}
