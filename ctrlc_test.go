package ctrlc

import (
	"context"
	"errors"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestCtrlcOK(t *testing.T) {
	if err := New().Run(context.Background(), func() error {
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestCtrlcErrors(t *testing.T) {
	err := errors.New("some error")
	if err != New().Run(context.Background(), func() error {
		return err
	}) {
		t.Fatalf("expected a different error, got: %v", err)
	}
}

func TestCtrlcTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	err := New().Run(ctx, func() error {
		t.Log("slow task...")
		time.Sleep(time.Minute)
		return nil
	})
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
	if err.Error() != "context deadline exceeded" {
		t.Fatalf("expected a different error, got: %v", err)
	}
}

func TestCtrlcSignals(t *testing.T) {
	t.Parallel()
	for _, signal := range []os.Signal{syscall.SIGINT, syscall.SIGTERM} {
		signal := signal
		t.Run(signal.String(), func(t *testing.T) {
			t.Parallel()
			h := New()
			errs := make(chan error, 1)
			go func() {
				errs <- h.Run(context.Background(), func() error {
					t.Log("slow task...")
					time.Sleep(time.Minute)
					return nil
				})
			}()
			h.signals <- signal
			err := <-errs
			if err == nil {
				t.Fatalf("expected an error, got nil")
			}
			if !errors.As(err, &ErrorCtrlC{}) {
				t.Fatalf("should have been a ErrorCtrlC, got %v", err)
			}
			eerr := ErrorCtrlC{signal}.Error()
			if err.Error() != eerr {
				t.Fatalf("expected a different error, got: %v, expected: %v", err, eerr)
			}
		})
	}
}
