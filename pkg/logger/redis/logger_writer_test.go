package redis

import (
	"fmt"
	"testing"
)

var loggerwriter *LoggerWriter

func TestMain(m *testing.M) {
	testClientConn, err := NewLoggerWriter(&Options{
		Host: "172.16.5.149",
		Port: 6379,
		Key:  "logstash"})
	if err != nil {
		fmt.Println(err)
	}
	loggerwriter = testClientConn
	m.Run()
}

func TestConnection(t *testing.T) {
	tests := []struct {
		name string
		w    *LoggerWriter
		want bool
	}{
		// TODO: Add test cases.
		{"", loggerwriter, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.TestConnection(); got != tt.want {
				t.Errorf("LoggerWriter.TestConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoggerWriter_Write(t *testing.T) {
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		w       *LoggerWriter
		args    args
		wantN   int
		wantErr bool
	}{
		// TODO: Add test cases.
		{"", loggerwriter, args{[]byte("raoj9ia")}, 7, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN, err := tt.w.Write(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoggerWriter.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("LoggerWriter.Write() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}
