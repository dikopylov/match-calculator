package application

import (
	"io"
	"occurrence-calculator/internal/model/application"
	"occurrence-calculator/internal/model/domain/counter"
	"occurrence-calculator/internal/model/infrastructure/parser"
	"strings"
	"testing"
)

type StringSourceStrategy struct {
}

func (s *StringSourceStrategy) Read(source string) (io.ReadCloser, error) {
	r := io.NopCloser(strings.NewReader(source))

	return r, nil
}

type MockOutput struct {
}

func (o *MockOutput) Write(p []byte) (n int, err error) {

	return 0, nil
}

func TestMatchCalculator_Calculate(t *testing.T) {
	specification := new(parser.Specification)
	specification.SearchText = "Go"

	finder := counter.NewMatchFinder(&StringSourceStrategy{}, new(counter.Parser))

	type fields struct {
		finder parser.Finder
		input  io.Reader
		output io.Writer
	}
	type args struct {
		specification   *parser.Specification
		numberOfWorkers int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "ThreeSource_FiveWorkers_returnTwo",
			fields: fields{
				finder: finder,
				input:  strings.NewReader("Go\ngo\nGo"),
				output: new(MockOutput),
			},
			args: args{
				specification:   specification,
				numberOfWorkers: 5,
			},
			want: 2,
		},
		{
			name: "FiveSource_FiveWorkers_returnFour",
			fields: fields{
				finder: finder,
				input:  strings.NewReader("Go\ngo\nGo\nGoGo\nQwerty"),
				output: new(MockOutput),
			},
			args: args{
				specification:   specification,
				numberOfWorkers: 5,
			},
			want: 4,
		},
		{
			name: "TenSource_FiveWorkers_returnSeven",
			fields: fields{
				finder: finder,
				input:  strings.NewReader("Go\ngo\nGo\nGoGo\nQwerty\nGo\nGo\nGogo\ngo\nGgo"),
				output: new(MockOutput),
			},
			args: args{
				specification:   specification,
				numberOfWorkers: 5,
			},
			want: 7,
		},
		{
			name: "TenSource_OneWorkers_returnSeven",
			fields: fields{
				finder: finder,
				input:  strings.NewReader("Go\ngo\nGo\nGoGo\nQwerty\nGo\nGo\nGogo\ngo\nGgo"),
				output: new(MockOutput),
			},
			args: args{
				specification:   specification,
				numberOfWorkers: 1,
			},
			want: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := application.NewMatchCalculator(
				tt.fields.finder,
				tt.fields.input,
				tt.fields.output,
			)

			if got := m.Calculate(tt.args.specification, tt.args.numberOfWorkers); got != tt.want {
				t.Errorf("Calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}
