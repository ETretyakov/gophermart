package storages

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"time"
)

type Strategy func(ctx context.Context) chan struct{}

func sleep(ctx context.Context, duration time.Duration) {
	select {
	case <-time.After(duration):
		return
	case <-ctx.Done():
		return
	}
}

func Backoff(repeats int, dur time.Duration, factor float64, jitter bool) Strategy {
	return func(ctx context.Context) chan struct{} {
		ch := make(chan struct{})
		go func() {
			defer close(ch)
			rnd := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
			for i := 0; i < repeats; i++ {
				select {
				case <-ctx.Done():
					return
				case ch <- struct{}{}:
				}

				delay := float64(dur) * math.Pow(factor, float64(i))
				if jitter {
					delay = rnd.Float64()*(float64(2*dur)) + (delay - float64(dur))
				}
				sleep(ctx, time.Duration(delay))
			}
		}()
		return ch
	}
}

func (r *Retryer) Do(ctx context.Context, fn func() error, errs ...error) (err error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	attempt := 0

	ch := r.Strategy(ctx)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case _, ok := <-ch:
			if !ok {
				return err
			}

			if r.OnRetry != nil && attempt > 0 {
				r.OnRetry(ctx, attempt, err)
			}

			if err = fn(); err == nil {
				return nil
			}
			if len(errs) > 0 && !oneOfErrs(err, errs...) {
				return err
			}
			attempt++
		}
	}
}

func oneOfErrs(err error, errs ...error) bool {
	for _, e := range errs {
		if errors.Is(e, err) {
			return true
		}
	}
	return false
}
