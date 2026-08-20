package artifactsign

import (
	"context"
	"time"
)

// IOWaiter 可注入的等待器（测试用）。
type IOWaiter interface {
	Wait(ctx context.Context, d time.Duration) error
}

type sleepWaiter struct{}

func (sleepWaiter) Wait(ctx context.Context, d time.Duration) error {
	return waitStep(ctx, d)
}
