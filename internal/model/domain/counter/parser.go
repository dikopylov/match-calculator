package counter

import (
	"bufio"
	"errors"
	"io"
	"occurrence-calculator/internal/model/infrastructure"
	"occurrence-calculator/internal/model/infrastructure/parser"
	"strings"
)

const (
	bufferSize = 2 * 1024
)

type Parser struct {
}

func (p *Parser) Parse(reader io.Reader, specification parser.Specification) *infrastructure.Result {
	result := new(infrastructure.Result)

	if specification.SearchText == "" {
		result.Error = errors.New("search text is empty")

		return result
	}

	bufferReader := bufio.NewReader(reader)
	var lastSymbol string
	var totalCount int

	for {
		originalBuf := make([]byte, bufferSize)

		n, err := bufferReader.Read(originalBuf)

		if n == 0 {
			if err == io.EOF {
				break
			}

			if err != nil {
				result.Error = err

				return result
			}
		}

		sourceText := lastSymbol + string(originalBuf[:n])
		lastSymbol = string(originalBuf[n-1 : n])

		totalCount += strings.Count(sourceText, specification.SearchText)
	}

	result.Value = totalCount

	return result
}
