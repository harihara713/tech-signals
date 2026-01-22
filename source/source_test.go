package source

import (
	"testing"
)

func TestTransformToTimeString(t *testing.T) {
	tests := []struct {
		str     string
		want    string
		wantErr bool
		name    string
	}{
		{
			name:    "Not a time string",
			str:     "this is not a time string",
			want:    "",
			wantErr: true,
		},
		{
			name:    "Time string",
			str:     "JULY 10, 2025",
			want:    "July 10, 2025",
			wantErr: false,
		},
		{
			name:    "Time string with dot at the end of month",
			str:     "DEC. 15, 2025",
			want:    "Dec 15, 2025",
			wantErr: false,
		},
		{
			name:    "Time string with a single digit date",
			str:     "SEPT 9, 2025",
			want:    "Sept 09, 2025",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := transformToTimeString(tt.str)

			if tt.wantErr && err == nil {
				t.Errorf("expected error, but got = nil")
				return
			}

			if !tt.wantErr && err != nil {
				t.Errorf("expected result, got error = %v", err)
				return
			}

			if res != tt.want {
				t.Errorf("expected = %s, got = %s", tt.want, res)
			}
		})
	}
}
