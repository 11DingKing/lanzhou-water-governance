package worker_test

import (
	"context"
	"errors"
	"github.com/11DingKing/lanzhou-water-governance/internal/worker"
	"log/slog"
	"testing"
	"time"
)

func TestRunnerRetriesWithBackoff(t *testing.T) {
	attempts := 0
	runner := worker.Runner{MaxAttempts: 3, Interval: time.Millisecond, Logger: slog.Default()}
	err := runner.ExecuteForTest(context.Background(), func(context.Context) error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary")
		}
		return nil
	})
	if err != nil || attempts != 3 {
		t.Fatalf("err=%v attempts=%d", err, attempts)
	}
}
func TestRunnerStopsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	runner := worker.Runner{MaxAttempts: 4, Interval: time.Millisecond, Logger: slog.Default()}
	err := runner.ExecuteForTest(ctx, func(context.Context) error { return errors.New("fail") })
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err=%v", err)
	}
}
