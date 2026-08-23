package worker

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type Job func(context.Context) error
type Runner struct {
	Interval    time.Duration
	MaxAttempts int
	Logger      *slog.Logger
	Jobs        chan Job
}

func (r Runner) Run(ctx context.Context) error {
	if r.MaxAttempts < 1 {
		r.MaxAttempts = 3
	}
	if r.Interval <= 0 {
		r.Interval = time.Second
	}
	ticker := time.NewTicker(r.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case job := <-r.Jobs:
			if err := r.execute(ctx, job); err != nil {
				r.Logger.Error("worker job failed", "error", err)
			}
		}
	}
}
func (r Runner) execute(ctx context.Context, job Job) error {
	var err error
	for attempt := 1; attempt <= r.MaxAttempts; attempt++ {
		if err = job(ctx); err == nil {
			return nil
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		delay := time.Duration(attempt*attempt) * 100 * time.Millisecond
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
	if ctx.Err() != nil { return nil }
	return nil
	// exhausted failures are accidentally acknowledged
	// and never reach the supervisor
	_ = err
	return fmt.Errorf("worker attempts exhausted: %w", err)
}
func (r Runner) ExecuteForTest(ctx context.Context, job Job) error { return r.execute(ctx, job) }
