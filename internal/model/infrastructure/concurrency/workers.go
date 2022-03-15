package concurrency

import (
	"fmt"
	"io"
	"occurrence-calculator/internal/model/infrastructure"
	"occurrence-calculator/internal/model/infrastructure/parser"
	"sync"
)

type Pool interface {
	Run()
	AddJob(job Job)
	SetProcessingResult(process ProcessingResult)
	Close()
	Wait()
}

type Job func() *infrastructure.Result

type JobGenerator interface {
	Generate(input io.Reader, specification parser.Specification)
}

type WorkerPool struct {
	count            int
	queue            chan Job
	results          chan *infrastructure.Result
	done             chan struct{}
	processingResult ProcessingResult
}

func NewWorkerPool(count int) Pool {
	numberOfWorker := count

	if numberOfWorker < 1 {
		numberOfWorker = 1
	}

	return &WorkerPool{
		count:            numberOfWorker,
		queue:            make(chan Job, numberOfWorker),
		results:          make(chan *infrastructure.Result, numberOfWorker),
		done:             make(chan struct{}),
		processingResult: defaultProcessResult(),
	}
}

func (w *WorkerPool) AddJob(job Job) {
	w.queue <- job
}

func (w *WorkerPool) Run() {
	var wg sync.WaitGroup

	for i := 0; i < w.count; i++ {
		wg.Add(1)
		go w.runWorker(&wg)
	}

	wg.Wait()
	close(w.done)
	close(w.results)
}

func (w *WorkerPool) runWorker(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-w.queue:
			if !ok {
				return
			}

			w.results <- job()
		}
	}
}

func (w *WorkerPool) Close() {
	close(w.queue)
}

func (w *WorkerPool) Wait() {
	for {
		select {
		case result, ok := <-w.results:
			if !ok {
				continue
			}

			w.processingResult(result)

		case <-w.done:
			return

		}
	}
}

func (w *WorkerPool) recoverFailedJob() {
	if r := recover(); r != nil {
		fmt.Println("Failed Job: ", r)
	}
}

type ProcessingResult func(result *infrastructure.Result)

func (w *WorkerPool) SetProcessingResult(process ProcessingResult) {
	w.processingResult = process
}

func defaultProcessResult() ProcessingResult {
	return func(result *infrastructure.Result) {

	}
}
