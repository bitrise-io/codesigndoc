package xcarchive

import (
	"path/filepath"
	"testing"
)

func TestIsMacOS(t *testing.T) {
	tests := []struct {
		name     string
		archPath string
		want     bool
		wantErr  bool
	}{
		{
			name:     "macOS",
			archPath: filepath.Join(sampleRepoPath(t), "archives/macos.xcarchive"),
			want:     true,
			wantErr:  false,
		},
		{
			name:     "iOS",
			archPath: filepath.Join(sampleRepoPath(t), "archives/ios.xcarchive"),
			want:     false,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsMacOS(tt.archPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsMacOS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsMacOS() = %v, want %v", got, tt.want)
			}
		})
	}
}
