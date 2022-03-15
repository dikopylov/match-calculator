package counter

import (
	"io"
	"occurrence-calculator/internal/model/domain/counter"
	"occurrence-calculator/internal/model/infrastructure"
	"occurrence-calculator/internal/model/infrastructure/parser"
	"reflect"
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	type args struct {
		reader        io.Reader
		specification parser.Specification
	}

	tests := []struct {
		name    string
		args    args
		want    *infrastructure.Result
		wantErr bool
	}{
		{"NoOne", args{strings.NewReader(strings.Repeat("Ga", 5)), parser.Specification{SearchText: "Go"}}, &infrastructure.Result{Value: 0}, false},
		{"NoOneLowerGo", args{strings.NewReader(strings.Repeat("go", 5)), parser.Specification{SearchText: "Go"}}, &infrastructure.Result{Value: 0}, false},
		{"One", args{strings.NewReader("Golaglang"), parser.Specification{SearchText: "Go"}}, &infrastructure.Result{Value: 1}, false},
		{"Five", args{strings.NewReader(strings.Repeat("Go", 5)), parser.Specification{SearchText: "Go"}}, &infrastructure.Result{Value: 5}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &counter.Parser{}
			got := p.Parse(tt.args.reader, tt.args.specification)
			if (got.Error != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", got.Error, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
