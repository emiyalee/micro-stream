package segmentservice

import "testing"

func Test_doMediaSegment(t *testing.T) {
	type args struct {
		src string
		dst string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"",
			args{"http://172.16.5.150:8080/vod/1080p.mp4", "http://172.16.5.150:11241/vod/1080p/1080p.m3u8"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := doMediaSegment(tt.args.src, tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("doMediaSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
