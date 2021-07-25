package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFrom(t *testing.T) {
	ctxWithLogger := WithLogger(context.Background())
	ogLogger := From(ctxWithLogger)

	tests := []struct {
		name string
		ctx  func() context.Context
		want bool
	}{
		{
			name: "with logger",
			ctx: func() context.Context {
				return ctxWithLogger
			},
			want: true,
		},
		{
			name: "without logger",
			ctx: func() context.Context {
				return context.Background()
			},
		},
		{
			name: "nil context",
			ctx: func() context.Context {
				return nil
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := From(tc.ctx())
			if tc.want {
				require.Equal(t, ogLogger, logger)
			} else {
				require.NotEqual(t, ogLogger, logger)
			}
		})
	}

}

func TestWithLogger(t *testing.T) {
	ctx := context.Background()
	newCtx := WithLogger(ctx)

	require.NotNil(t, newCtx.Value(logKey))
}
