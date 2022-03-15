package commands

import (
	"errors"
	"fmt"
	"io"
	"occurrence-calculator/internal/model/application"
	"occurrence-calculator/internal/model/domain/counter"
	"occurrence-calculator/internal/model/infrastructure/cli"
	"occurrence-calculator/internal/model/infrastructure/parser"
)

const (
	searchText      = "Go"
	numberOfWorkers = 5
)

type MatchCounterCommand struct {
	Input  io.Reader
	Output io.Writer
	flags  map[string]interface{}
}

func NewMatchCounterCommand(input io.Reader, output io.Writer) *MatchCounterCommand {
	return &MatchCounterCommand{Input: input, Output: output}
}

func (c *MatchCounterCommand) Execute() {
	source := c.source()
	sourceStrategy := c.sourceStrategy(source)

	specification := new(parser.Specification)
	specification.SearchText = searchText

	finder := counter.NewMatchFinder(sourceStrategy, new(counter.Parser))
	calculator := application.NewMatchCalculator(finder, c.Input, c.Output)
	total := calculator.Calculate(specification, numberOfWorkers)

	totalMessage := fmt.Sprintf("Total: %d\n", total)

	c.Output.Write([]byte(totalMessage))
}

func (c *MatchCounterCommand) sourceStrategy(source string) counter.SourceStrategy {
	sourceStrategy, err := counter.MakeSourceStrategyBySourceType(source)

	if err != nil {
		_, err = c.Output.Write([]byte(fmt.Sprintf("Source type \"%s\" is not present. Available: url, file\n", source)))

		if err != nil {
			panic(err)
		}

		return nil
	}
	return sourceStrategy
}

func (c *MatchCounterCommand) source() string {
	sourceFlag := c.GetFlag(cli.FlagType)

	if sourceFlag == nil {
		panic(errors.New("source is empty"))
	}
	return fmt.Sprintf("%s", sourceFlag)
}

func (c *MatchCounterCommand) AddFlag(name string, value interface{}) {
	if c.flags == nil {
		c.flags = make(map[string]interface{})
	}

	c.flags[name] = value
}

func (c *MatchCounterCommand) GetFlag(name string) interface{} {
	if val, ok := c.flags[name]; ok {
		return val
	}

	return nil
}
