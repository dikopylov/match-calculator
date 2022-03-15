package parser

import (
	"io"
	"occurrence-calculator/internal/model/infrastructure"
)

type Parser interface {
	Parse(reader io.Reader, specification Specification) *infrastructure.Result
}

type Specification struct {
	SearchText string
}

type Finder interface {
	Find(source string, specification Specification) *infrastructure.Result
}
