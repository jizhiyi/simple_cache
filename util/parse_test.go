package util

import "testing"

func TestParseMemorySize(t *testing.T) {
	type args struct {
		sizeStr string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:    "1KB",
			args:    args{sizeStr: "1KB"},
			want:    1024,
			wantErr: false,
		},
		{
			name:    "100KB",
			args:    args{sizeStr: "100KB"},
			want:    1024 * 100,
			wantErr: false,
		},
		{
			name:    "1MB",
			args:    args{sizeStr: "1MB"},
			want:    1024 * 1024,
			wantErr: false,
		},
		{
			name:    "10MB",
			args:    args{sizeStr: "10MB"},
			want:    1024 * 1024 * 10,
			wantErr: false,
		},
		{
			name:    "1GB",
			args:    args{sizeStr: "1GB"},
			want:    1024 * 1024 * 1024,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMemorySize(tt.args.sizeStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMemorySize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseMemorySize() got = %v, want %v", got, tt.want)
			}
		})
	}
}
