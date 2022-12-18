package alt

import (
	"os"
	"reflect"
	"testing"
)

func Test_dataEtl_Parse(t *testing.T) {
	t.Helper()

	type fields struct {
		separator Separator
		logger    Logger
		maxDepth  int
		ignore    map[string]struct{}
	}
	type args struct {
		data interface{}
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult matrixKvPairs
	}{

		{
			name: "success_1",
			fields: fields{
				separator: StrSeparator(""),
				logger:    NewStdLogger(LevelDebug, os.Stdout),
				maxDepth:  Infinity,
				ignore:    map[string]struct{}{},
			},
			args: args{data: 1},
			wantResult: matrixKvPairs{
				[]pair{{Key: "", Value: 1}},
			},
		},
		{
			name: "success_2",
			fields: fields{
				separator: StrSeparator(""),
				logger:    NewStdLogger(LevelDebug, os.Stdout),
				maxDepth:  Infinity,
				ignore:    map[string]struct{}{},
			},
			args: args{data: "hello"},
			wantResult: matrixKvPairs{
				[]pair{{Key: "", Value: "hello"}},
			},
		},
		{
			name: "success_3",
			fields: fields{
				separator: StrSeparator(""),
				logger:    NewStdLogger(LevelDebug, os.Stdout),
				maxDepth:  Infinity,
				ignore:    map[string]struct{}{},
			},
			args: args{data: []int{1, 2}},
			wantResult: matrixKvPairs{
				[]pair{{Key: "", Value: 1}},
				[]pair{{Key: "", Value: 2}},
			},
		},
		{
			name: "success_4",
			fields: fields{
				separator: StrSeparator(""),
				logger:    NewStdLogger(LevelDebug, os.Stdout),
				maxDepth:  Infinity,
				ignore:    map[string]struct{}{},
			},
			args: args{data: map[string]interface{}{
				"aa": "a",
				"bb": []int{1, 2},
			}},
			wantResult: matrixKvPairs{
				[]pair{{Key: "aa", Value: "a"}, {Key: "bb", Value: 1}},
				[]pair{{Key: "aa", Value: "a"}, {Key: "bb", Value: 2}},
			},
		},
		{
			name: "success_with_key",
			fields: fields{
				separator: StrSeparator(""),
				logger:    NewStdLogger(LevelDebug, os.Stdout),
				maxDepth:  Infinity,
				ignore:    map[string]struct{}{"key": struct{}{}},
			},
			args: args{data: map[string]interface{}{"key": "我是很长的字符串"}},
			wantResult: matrixKvPairs{
				[]pair{},
			},
		},
		{
			name: "success_with_key_2",
			fields: fields{
				separator: StrSeparator(""),
				logger:    NewStdLogger(LevelDebug, os.Stdout),
				maxDepth:  Infinity,
				ignore:    map[string]struct{}{"key": struct{}{}},
			},
			args: args{data: map[string]interface{}{"key": "我是很长的字符串", "aa": "bb"}},
			wantResult: matrixKvPairs{
				[]pair{{
					Key:   "aa",
					Value: "bb",
				}},
			},
		},

		{
			name: "success_complex",
			fields: fields{
				separator: StrSeparator("."),
				logger:    NewStdLogger(LevelDebug, os.Stdout),
				maxDepth:  Infinity,
				ignore:    map[string]struct{}{},
			},
			args: args{data: map[string]interface{}{
				"aa": "a",
				"bb": []int{1, 2},
				"cc": map[string]interface{}{
					"ee": 2,
					"dd": 1,
				},
			}},
			wantResult: matrixKvPairs{
				[]pair{{Key: "aa", Value: "a"}, {Key: "bb", Value: 1}, {Key: "cc.dd", Value: 1}, {Key: "cc.ee", Value: 2}},
				[]pair{{Key: "aa", Value: "a"}, {Key: "bb", Value: 2}, {Key: "cc.dd", Value: 1}, {Key: "cc.ee", Value: 2}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &dataEtl{
				separator: tt.fields.separator,
				logger:    tt.fields.logger,
				maxDepth:  tt.fields.maxDepth,
				ignore:    tt.fields.ignore,
			}
			if gotResult := covertHelper(c.Parse(tt.args.data)); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Parse() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func BenchmarkDataEtl_Parse_simple(b *testing.B) {
	parser := NewDataEtlParser(SetMaxDepth(Infinity), SetLogger(NewStdLogger(LevelWarn, os.Stdout)))
	for i := 0; i < b.N; i++ {
		parser.Parse(1)
	}
}

func BenchmarkDataEtl_Parse_object(b *testing.B) {
	parser := NewDataEtlParser(SetMaxDepth(Infinity), SetLogger(NewStdLogger(LevelWarn, os.Stdout)))
	for i := 0; i < b.N; i++ {
		parser.Parse(map[string]interface{}{"key": "value", "id": 9, "name": "text", "cuda": 999})
	}
}

func BenchmarkDataEtl_Parse_list(b *testing.B) {
	var l = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	parser := NewDataEtlParser(SetMaxDepth(Infinity), SetLogger(NewStdLogger(LevelWarn, os.Stdout)))
	for i := 0; i < b.N; i++ {
		parser.Parse(l)
	}
}

func BenchmarkDataEtl_Parse_complex(b *testing.B) {
	var cp = map[string]interface{}{
		"aa": "a",
		"bb": []int{1, 2},
		"cc": map[string]interface{}{
			"ee": 2,
			"dd": 1,
		},
		"l4":   []int{1, 2, 3, 4},
		"l5":   []int{1, 2, 3, 4, 5},
		"info": map[string]interface{}{"key": "value", "id": 9, "name": "text", "cuda": 999},
	}
	parser := NewDataEtlParser(SetMaxDepth(Infinity), SetLogger(NewStdLogger(LevelWarn, os.Stdout)))
	for i := 0; i < b.N; i++ {
		parser.Parse(cp)
	}
}
