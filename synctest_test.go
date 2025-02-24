package main

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"golang.org/x/time/rate"
)

func TestWithTimeout(t *testing.T) {
	synctest.Run(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		select {
		case <-time.After(15 * time.Second): // Simulate a long running operation
			t.Fatal("expected timeout, but got a signal")
		case <-ctx.Done():
			return
		}
	})
}

func TestBackoffRetry(t *testing.T) {
	synctest.Run(func() {
		var attempts int

		ctx, cancel := context.WithCancel(context.Background())
		// Simulate a backoff retry function
		go func(ctx context.Context) {
			ticker := time.NewTicker(time.Second*time.Duration(attempts) + time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					attempts += 1
					ticker.Reset(time.Second*time.Duration(attempts) + time.Second)
				case <-ctx.Done():
					return
				}
			}
		}(ctx)

		time.Sleep(time.Second * 35) // Let it run for 35 seconds
		cancel()                     // Stop the backoff retry function

		// 1 + 2 + 3 + 4 + 5 + 6 + 7 = 28 seconds
		expected := 7
		if attempts != expected {
			t.Fatalf("expected %d cleanups, got %d", expected, attempts)
		}
	})
}

func TestRateLimit(t *testing.T) {
	synctest.Run(func() {
		l := rate.NewLimiter(rate.Every(time.Second), 1) // Allow one req per second

		var consumed int
		ctx, cancel := context.WithCancel(context.Background())
		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					err := l.Wait(ctx)
					if err == nil {
						consumed++
					}
				}
			}
		}(ctx)

		time.Sleep(time.Hour * 24) // See how many requests we can make in a day
		cancel()

		// 1 requests per second * 60 seconds * 60 minutes * 24 hours + 1 because we are starting with one token
		expected := 1*60*60*24 + 1
		if consumed != expected {
			t.Fatalf("expected %d requests, got %d", expected, consumed)
		}
	})
}
