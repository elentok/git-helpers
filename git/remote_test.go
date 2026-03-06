package git

import "testing"

func TestIsNonFastForwardPushError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "non fast forward",
			err: &RunError{
				Stderr: "! [rejected]        main -> main (non-fast-forward)\nerror: failed to push some refs",
			},
			want: true,
		},
		{
			name: "fetch first",
			err: &RunError{
				Stderr: "Updates were rejected because the remote contains work that you do not have locally. (fetch first)",
			},
			want: true,
		},
		{
			name: "other error",
			err: &RunError{
				Stderr: "fatal: could not read from remote repository",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNonFastForwardPushError(tt.err); got != tt.want {
				t.Fatalf("IsNonFastForwardPushError() = %v, want %v", got, tt.want)
			}
		})
	}
}
