package alt

import "testing"

func TestStrSeparator_AppendToPrefix(t *testing.T) {
	var underLineSep = StrSeparator("_")
	type args struct {
		prefix string
		key    interface{}
	}
	tests := []struct {
		name string
		c    StrSeparator
		args args
		want string
	}{
		{
			name: "success",
			c:    underLineSep,
			want: "pre_",
			args: args{
				prefix: "pre",
				key:    "",
			},
		},
		// use nil as key
		{
			name: "success",
			c:    underLineSep,
			want: "pre_<nil>",
			args: args{
				prefix: "pre",
				key:    nil,
			},
		},
		// use `1` as key
		{
			name: "success",
			c:    underLineSep,
			want: "pre_1",
			args: args{
				prefix: "pre",
				key:    1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.AppendToPrefix(tt.args.prefix, tt.args.key); got != tt.want {
				t.Errorf("AppendToPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
