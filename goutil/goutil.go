package goutil

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bryanl/lilutil/log"
)

// WaitForChannelsToClose waits for channels to close. If all the channels do not
// close before the context is done, false will be returned.
func WaitForChannelsToClose(ctx context.Context, chans ...<-chan struct{}) bool {
	done := make(chan struct{}, 1)

	go func() {
		for _, v := range chans {
			<-v
		}

		close(done)
	}()

	select {
	case <-ctx.Done():
		return false
	case <-done:
		return true
	}
}

// HandleGracefulClose gracefully handles shutting down the process.
func HandleGracefulClose(ctx context.Context, cancel context.CancelFunc, chans ...<-chan struct{}) {
	logger := log.From(ctx).WithName("graceful")

	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)

	<-signalChan
	logger.Info("shutting down gracefully")

	closeCtx, closeCancel := context.WithTimeout(ctx, 5*time.Second)
	defer closeCancel()

	go func() {
		<-signalChan
		logger.Info("terminating")
		closeCancel()
	}()

	cancel()

	logger.Info("waiting for servers to stop")

	if !WaitForChannelsToClose(closeCtx, chans...) {
		logger.Info("all channels were not closed")
	}

	logger.Info("exiting normally")
}
