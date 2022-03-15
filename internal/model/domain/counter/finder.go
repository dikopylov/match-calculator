package counter

import (
	"occurrence-calculator/internal/model/infrastructure"
	"occurrence-calculator/internal/model/infrastructure/parser"
)

type MatchFinder struct {
	sourceStrategy SourceStrategy
	parser         parser.Parser
}

func NewMatchFinder(sourceStrategy SourceStrategy, parser parser.Parser) parser.Finder {
	return &MatchFinder{
		sourceStrategy: sourceStrategy,
		parser:         parser,
	}
}

func (m *MatchFinder) Find(source string, specification parser.Specification) *infrastructure.Result {
	reader, err := m.sourceStrategy.Read(source)

	if err != nil {
		result := new(infrastructure.Result)
		result.Error = err

		return result
	}

	defer reader.Close()

	result := m.parser.Parse(reader, specification)

	return result
}
