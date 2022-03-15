package concurrency

import (
	"occurrence-calculator/internal/model/infrastructure"
	"occurrence-calculator/internal/model/infrastructure/concurrency"
	"sync"
	"testing"
)

func TestWorkerPool_Run(t *testing.T) {
	var total int
	type fields struct {
		count            int
		processingResult concurrency.ProcessingResult
	}
	type args struct {
		jobs []concurrency.Job
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expected int
	}{
		{
			name: "OneWorker_FiveJobs",
			fields: fields{
				count:            1,
				processingResult: testTotalProcessResult(&total),
			},
			args: args{
				jobs: []concurrency.Job{
					newTestJob(0),
					newTestJob(2),
					newTestJob(2),
					newTestJob(2),
					newTestJob(2),
				},
			},
			expected: 8,
		},
		{
			name: "FiveWorker_ThreeJobs",
			fields: fields{
				count:            5,
				processingResult: testTotalProcessResult(&total),
			},
			args: args{
				jobs: []concurrency.Job{
					newTestJob(0),
					newTestJob(1),
					newTestJob(2),
				},
			},
			expected: 3,
		},
		{
			name: "FiveWorker_FiveJobs",
			fields: fields{
				count:            5,
				processingResult: testTotalProcessResult(&total),
			},
			args: args{
				jobs: []concurrency.Job{
					newTestJob(0),
					newTestJob(1),
					newTestJob(2),
					newTestJob(3),
					newTestJob(4),
				},
			},
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := concurrency.NewWorkerPool(tt.fields.count)
			w.SetProcessingResult(tt.fields.processingResult)
			go w.Run()

			go func() {
				for _, job := range tt.args.jobs {
					w.AddJob(job)
				}
				w.Close()
			}()

			w.Wait()

			if total != tt.expected {
				t.Errorf("TestWorkerPool_Run test got = %v, expected %v", total, tt.expected)
			}

			total = 0
		})
	}
}

func testTotalProcessResult(total *int) concurrency.ProcessingResult {
	var mutex sync.Mutex

	return func(result *infrastructure.Result) {
		mutex.Lock()
		*total = result.Value.(int) + *total
		mutex.Unlock()
	}
}

func newTestJob(start int) concurrency.Job {
	return func() *infrastructure.Result {
		return &infrastructure.Result{Value: start}
	}
}
