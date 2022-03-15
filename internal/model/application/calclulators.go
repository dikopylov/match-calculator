package application

import (
	"io"
	"occurrence-calculator/internal/model/domain/counter"
	"occurrence-calculator/internal/model/infrastructure/concurrency"
	"occurrence-calculator/internal/model/infrastructure/parser"
)

type MatchCalculator struct {
	finder parser.Finder
	input  io.Reader
	output io.Writer
}

func NewMatchCalculator(finder parser.Finder, input io.Reader, output io.Writer) *MatchCalculator {
	return &MatchCalculator{finder: finder, input: input, output: output}
}

func (m *MatchCalculator) Calculate(specification *parser.Specification, numberOfWorkers int) int {
	var total int

	workerPool := concurrency.NewWorkerPool(numberOfWorkers)
	workerPool.SetProcessingResult(counter.TotalProcessResult(&total))

	go workerPool.Run()

	go func() {
		jobGenerator := counter.NewMatchCalculatorJobGenerator(workerPool, m.finder, m.output)

		jobGenerator.Generate(m.input, *specification)
	}()

	workerPool.Wait()

	return total
}
