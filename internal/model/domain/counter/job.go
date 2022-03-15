package counter

import (
	"bufio"
	"fmt"
	"io"
	"occurrence-calculator/internal/model/infrastructure"
	"occurrence-calculator/internal/model/infrastructure/concurrency"
	"occurrence-calculator/internal/model/infrastructure/parser"
	"sync"
)

type MatchCalculatorJobGenerator struct {
	wp     concurrency.Pool
	finder parser.Finder
	output io.Writer
}

func NewMatchCalculatorJobGenerator(wp concurrency.Pool, finder parser.Finder, output io.Writer) *MatchCalculatorJobGenerator {
	return &MatchCalculatorJobGenerator{wp: wp, finder: finder, output: output}
}

func (m *MatchCalculatorJobGenerator) Generate(reader io.Reader, specification parser.Specification) {
	bufferSourceScanner := bufio.NewScanner(reader)

	for bufferSourceScanner.Scan() {
		job := NewMatchCalculatorJob(m.finder, m.output, bufferSourceScanner.Text(), specification)

		m.wp.AddJob(job)
	}

	m.wp.Close()
}

func NewMatchCalculatorJob(finder parser.Finder, output io.Writer, source string, specification parser.Specification) concurrency.Job {
	return func() *infrastructure.Result {
		result := finder.Find(fmt.Sprintf("%v", source), specification)

		if result.Error != nil {
			panic(result.Error)
		}

		sourceTotalMessage := fmt.Sprintf("Count for %s: %d\n", source, result.Value)

		_, err := output.Write([]byte(sourceTotalMessage))

		if err != nil {
			panic(err)
		}

		return result
	}
}

func TotalProcessResult(total *int) concurrency.ProcessingResult {
	var mutex sync.Mutex

	return func(result *infrastructure.Result) {
		value, ok := result.Value.(int)

		if ok {
			mutex.Lock()

			*total = value + *total

			mutex.Unlock()
		}
	}
}
