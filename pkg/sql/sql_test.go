package sql

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var cc *ClientConn

func TestMain(m *testing.M) {
	testClientConn, err := NewClientConn("172.16.5.149", "root", "123456", "stream_system")
	if err != nil {
		fmt.Println(err)
	}
	cc = testClientConn
	defer cc.Close()
	m.Run()
}

func TestClientConn_QueryStoreURL(t *testing.T) {
	type args struct {
		resourceID string
	}
	tests := []struct {
		name    string
		cc      *ClientConn
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"", cc, args{"a"}, "store", "Wildlife.mp4", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.cc.QueryStoreURL(tt.args.resourceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientConn.QueryStoreURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ClientConn.QueryStoreURL() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ClientConn.QueryStoreURL() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestClientConn_QuerySteamingURL(t *testing.T) {
	type args struct {
		resourceID string
	}
	tests := []struct {
		name    string
		cc      *ClientConn
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"", cc, args{"a"}, "stream", "a/a.m3u8", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.cc.QuerySteamingURL(tt.args.resourceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientConn.QuerySteamingURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ClientConn.QuerySteamingURL() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ClientConn.QuerySteamingURL() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestClientConn_AddStreamingURL(t *testing.T) {
	type args struct {
		resourceID    string
		streamAddress string
		endpoint      string
	}
	tests := []struct {
		name    string
		cc      *ClientConn
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"", cc, args{"c", "stream", "c/c.m3u8"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cc.AddStreamingURL(tt.args.resourceID, tt.args.streamAddress, tt.args.endpoint); (err != nil) != tt.wantErr {
				t.Errorf("ClientConn.AddStreamingURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientConn_UpdateStreamingURL(t *testing.T) {
	type args struct {
		resourceID    string
		streamAddress string
		endpoint      string
	}
	tests := []struct {
		name    string
		cc      *ClientConn
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"", cc, args{"c", "stream", "c/drawheart.m3u8"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cc.UpdateStreamingURL(tt.args.resourceID, tt.args.streamAddress, tt.args.endpoint); (err != nil) != tt.wantErr {
				t.Errorf("ClientConn.UpdateStreamingURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientConn_DeleteStreamingURL(t *testing.T) {
	type args struct {
		resourceID string
	}
	tests := []struct {
		name    string
		cc      *ClientConn
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"", cc, args{"c"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cc.DeleteStreamingURL(tt.args.resourceID); (err != nil) != tt.wantErr {
				t.Errorf("ClientConn.DeleteStreamingURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
