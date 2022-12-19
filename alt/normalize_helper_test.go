package alt

import (
	"reflect"
	"testing"
)

func Test_kvPairs_Len(t *testing.T) {
	tests := []struct {
		name string
		c    kvPairs
		want int
	}{
		{
			name: "success_empty",
			c:    kvPairs{},
			want: 0,
		},
		{
			name: "success_len_2",
			c: kvPairs{
				pair{Key: "key", Value: nil},
			},
			want: 1,
		},
		{
			name: "success_len_2",
			c: kvPairs{
				pair{Key: "key", Value: nil},
				pair{Key: "key_2", Value: 9},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_kvPairs_Less(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		c    kvPairs
		args args
		want bool
	}{
		{
			name: "success",
			c:    []pair{{Key: "key_a", Value: 1}, {Key: "key_b", Value: 1}},
			args: args{i: 0, j: 1},
			want: true,
		},

		{
			name: "success_2",
			c:    []pair{{Key: "key_a", Value: 1}, {Key: "key_a", Value: 2}},
			args: args{i: 0, j: 1},
			want: true,
		},

		{
			name: "success_3",
			c:    []pair{{Key: "key_a", Value: 2}, {Key: "key_a", Value: 1}},
			args: args{i: 0, j: 1},
			want: false,
		},

		{
			name: "success_4",
			c:    []pair{{Key: "key_b", Value: 1}, {Key: "key_a", Value: 2}},
			args: args{i: 0, j: 1},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Less(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("name %s, Less() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_kvPairs_Swap(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		c    kvPairs
		args args
		want kvPairs
	}{
		{
			name: "success",
			c: []pair{
				{
					Key:   "key_a",
					Value: 1,
				},
				{
					Key:   "key_b",
					Value: 1,
				},
			},
			args: args{
				i: 0,
				j: 1,
			},
			want: []pair{
				{
					Key:   "key_b",
					Value: 1,
				},
				{
					Key:   "key_a",
					Value: 1,
				},
			},
		},
	}
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Swap(tt.args.i, tt.args.j)
			if !reflect.DeepEqual(tt.c, tt.want) {
				t.Errorf("test not pass %v , %v", tt.args, tt.want)
			} else {
				t.Logf("test pass %s ", tt.name)
			}
		})
	}
}

func Test_matrixKvPairs_Len(t *testing.T) {
	t.Helper()
	tests := []struct {
		name string
		c    matrixKvPairs
		want int
	}{
		{
			name: "success",
			c:    [][]pair{},
			want: 0,
		},
		{
			name: "success_1",
			c: [][]pair{
				[]pair{},
			},
			want: 1,
		},
		{
			name: "success_2",
			c: [][]pair{
				[]pair{},
				[]pair{},
			},
			want: 2,
		},
		{
			name: "success_3",
			c: [][]pair{
				{},
				{},
				{},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matrixKvPairs_Less(t *testing.T) {
	t.Helper()
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		c    matrixKvPairs
		args args
		want bool
	}{

		{
			name: "success_3",
			c: [][]pair{
				[]pair{
					{Key: "key_a", Value: 0},
				},
				[]pair{
					{Key: "key_a", Value: -1},
				},
			},
			args: args{i: 0, j: 1},
			want: false,
		},
		{
			name: "success_4",
			c: [][]pair{
				[]pair{
					{Key: "key_a", Value: -1},
				},
				[]pair{
					{Key: "key_a", Value: 0},
				},
			},
			args: args{i: 0, j: 1},
			want: true,
		},
		{
			name: "success_5",
			c: [][]pair{
				[]pair{
					{Key: "key_b", Value: -1},
				},
				[]pair{
					{Key: "key_a", Value: 0},
				},
			},
			args: args{i: 0, j: 1},
			want: false,
		},
		{
			name: "success_6",
			c: [][]pair{
				[]pair{
					{Key: "key_a", Value: 0},
					{Key: "key_b", Value: 1},
				},
				[]pair{
					{Key: "key_a", Value: 0},
					{Key: "key_b", Value: -1},
				},
			},
			args: args{i: 0, j: 1},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Less(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf(" name = %s, Less() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_matrixKvPairs_Swap(t *testing.T) {
	t.Helper()
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		c    matrixKvPairs
		args args
		want matrixKvPairs
	}{
		{
			name: "success",
			c: [][]pair{
				[]pair{
					{Key: "key_b", Value: 1},
				},
				[]pair{
					{Key: "key_a", Value: 1},
				},
			},
			args: args{i: 0, j: 1},
			want: [][]pair{
				[]pair{
					{Key: "key_a", Value: 1},
				},
				[]pair{
					{Key: "key_b", Value: 1},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Swap(tt.args.i, tt.args.j)
			if !reflect.DeepEqual(tt.c, tt.want) {

				t.Errorf("test not pass %v , %v", tt.args, tt.want)
			} else {
				t.Logf("test pass %s ", tt.name)
			}
		})
	}
}

func Test_covertHelper(t *testing.T) {
	type args struct {
		objList []map[string]interface{}
	}
	tests := []struct {
		name         string
		args         args
		parserObject interface{}
		wantRes      matrixKvPairs
	}{
		{
			name:         "success",
			parserObject: 1,
			args: args{
				objList: []map[string]interface{}{
					{".": 1},
				},
			},
			wantRes: [][]pair{[]pair{{Key: ".", Value: 1}}},
		},

		{
			name:         "success_1",
			parserObject: 1,
			args: args{
				objList: []map[string]interface{}{
					{"aa": 1, "bb": "xx"},
				},
			},
			wantRes: [][]pair{[]pair{{Key: "aa", Value: 1}, {Key: "bb", Value: "xx"}}},
		},

		{
			name:         "success_2",
			parserObject: 1,
			args: args{
				objList: []map[string]interface{}{
					{"aa": 1, "bb": 2},
					{"aa": 1, "bb": 1},
				},
			},
			wantRes: [][]pair{
				[]pair{{Key: "aa", Value: 1}, {Key: "bb", Value: 1}},
				[]pair{{Key: "aa", Value: 1}, {Key: "bb", Value: 2}},
			},
		},

		{
			name:         "success_3",
			parserObject: 1,
			args: args{
				objList: []map[string]interface{}{
					{"aa": 1, "bb": "bb"},
					{"aa": 1, "bb": "cc"},
				},
			},
			wantRes: [][]pair{
				[]pair{{Key: "aa", Value: 1}, {Key: "bb", Value: "bb"}},
				[]pair{{Key: "aa", Value: 1}, {Key: "bb", Value: "cc"}},
			},
		},

		{
			name:         "success_4",
			parserObject: 1,
			args: args{
				objList: []map[string]interface{}{
					{"aa": 1, "bb": "zz"},
					{"aa": 1, "bb": "xx"},
				},
			},
			wantRes: [][]pair{
				[]pair{{Key: "aa", Value: 1}, {Key: "bb", Value: "xx"}},
				[]pair{{Key: "aa", Value: 1}, {Key: "bb", Value: "zz"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := covertHelper(tt.args.objList); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("name = %s, covertHelper() = %v, want %v", tt.name, gotRes, tt.wantRes)
			}
		})
	}
}
