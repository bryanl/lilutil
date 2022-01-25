package stringutil_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bryanl/lilutil/stringutil"
)

func TestRemove(t *testing.T) {
	tests := []struct {
		name string
		s    string
		sl   []string
		want []string
	}{
		{
			name: "empty slice",
			s:    "s",
			sl:   nil,
			want: nil,
		},

		{
			name: "exists",
			s:    "s",
			sl:   []string{"a", "s", "k"},
			want: []string{"a", "k"},
		},

		{
			name: "does not exist",
			s:    "b",
			sl:   []string{"a", "s", "k"},
			want: []string{"a", "s", "k"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := stringutil.Remove(tc.s, tc.sl)
			require.Equal(t, tc.want, got)
		})
	}
}
