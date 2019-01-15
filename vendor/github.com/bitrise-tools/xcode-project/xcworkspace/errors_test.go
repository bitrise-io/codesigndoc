package xcworkspace

import (
	"errors"
	"testing"
)

func TestSchemeNotFoundError_Error(t *testing.T) {
	err := SchemeNotFoundError{scheme: "Scheme", container: "Workspace"}
	want := "scheme Scheme not found in Workspace"
	if err.Error() != want {
		t.Errorf("SchemeNotFoundError.Error() = %v, want %v", err, want)
	}
}

func TestIsSchemeNotFoundError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "SchemeNotFoundError",
			err:  SchemeNotFoundError{scheme: "Scheme", container: "Workspace"},
			want: true,
		},
		{
			name: "not SchemeNotFoundError",
			err:  errors.New("other error"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSchemeNotFoundError(tt.err); got != tt.want {
				t.Errorf("IsSchemeNotFoundError() = %v, want %v", got, tt.want)
			}
		})
	}
}
