package alt

import (
	"io"
	"os"
	"testing"
)

func TestLogLevel_String(t *testing.T) {
	tests := []struct {
		name string
		c    LogLevel
		want string
	}{
		{name: "unknown test", c: -1, want: "DEBUG"},       // default is debug
		{name: "debug test", c: LevelDebug, want: "DEBUG"}, // debug
		{name: "info test", c: LevelInfo, want: "INFO"},    // info
		{name: "warn test", c: LevelWarn, want: "WARN"},    // warn
		{name: "error test", c: LevelError, want: "ERROR"}, // error
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stdLogger_Debug(t *testing.T) {
	type fields struct {
		level  LogLevel
		writer io.Writer
	}
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "test pass", fields: fields{level: LevelDebug, writer: os.Stdout}, args: args{args: []interface{}{1, 2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &stdLogger{
				level:  tt.fields.level,
				writer: tt.fields.writer,
			}
			c.Debug(tt.args.args...)
		})
	}
}

func Test_stdLogger_Info(t *testing.T) {
	type fields struct {
		level  LogLevel
		writer io.Writer
	}
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "info test pass", fields: fields{level: LevelInfo, writer: os.Stdout}, args: args{args: []interface{}{1, 2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := stdLogger{
				level:  tt.fields.level,
				writer: tt.fields.writer,
			}
			c.Info(tt.args.args...)
		})
	}
}

func Test_stdLogger_Warn(t *testing.T) {
	type fields struct {
		level  LogLevel
		writer io.Writer
	}
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "warn test pass", fields: fields{level: LevelInfo, writer: os.Stdout}, args: args{args: []interface{}{1, 2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := stdLogger{
				level:  tt.fields.level,
				writer: tt.fields.writer,
			}
			c.Warn(tt.args.args...)
		})
	}
}

func Test_stdLogger_Error(t *testing.T) {
	type fields struct {
		level  LogLevel
		writer io.Writer
	}
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "error test pass", fields: fields{level: LevelInfo, writer: os.Stdout}, args: args{args: []interface{}{1, 2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := stdLogger{
				level:  tt.fields.level,
				writer: tt.fields.writer,
			}
			c.Error(tt.args.args...)
		})
	}
}
