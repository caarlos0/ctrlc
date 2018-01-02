package ctrlc

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestCtrlcOK(t *testing.T) {
	assert.NoError(t, New().Run(context.Background(), func() error {
		return nil
	}))
}

func TestCtrlcErrors(t *testing.T) {
	var err = errors.New("some error")
	assert.EqualError(t, New().Run(context.Background(), func() error {
		return err
	}), err.Error())
}

func TestCtrlcTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	assert.EqualError(t, New().Run(ctx, func() error {
		t.Log("slow task...")
		time.Sleep(time.Minute)
		return nil
	}), "context deadline exceeded")
}

func TestCtrlcSignals(t *testing.T) {
	for _, signal := range []os.Signal{syscall.SIGINT, syscall.SIGTERM} {
		signal := signal
		t.Run(signal.String(), func(tt *testing.T) {
			tt.Parallel()
			var h = New()
			var errs = make(chan error, 1)
			go func() {
				errs <- h.Run(context.Background(), func() error {
					tt.Log("slow task...")
					time.Sleep(time.Minute)
					return nil
				})
			}()
			h.signals <- signal
			assert.EqualError(tt, <-errs, fmt.Sprintf("received: %s", signal))
		})
	}
}

func BenchmarkCtrlc(b *testing.B) {
	var task Task = func() error {
		return nil
	}
	var h = New()
	var ctx = context.Background()
	var wg errgroup.Group
	for i := 0; i < 10000; i++ {
		wg.Go(func() error {
			return h.Run(ctx, task)
		})
	}
	assert.NoError(b, wg.Wait())
}
