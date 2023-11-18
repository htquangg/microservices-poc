package waiter

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type WaitFunc func(ctx context.Context) error

type CleanupFunc func()

type Waiter interface {
	Add(fn ...WaitFunc)
	Cleanup(fn ...CleanupFunc)
	Wait() error
	Context() context.Context
	CancelFunc() context.CancelFunc
}

type waiter struct {
	ctx          context.Context
	waitFuncs    []WaitFunc
	cleanupFuncs []CleanupFunc
	cancel       context.CancelFunc
}

type waiterCfg struct {
	parentCtx    context.Context
	catchSignals bool
}

func New(options ...WaiterOption) Waiter {
	cfg := &waiterCfg{
		parentCtx:    context.Background(),
		catchSignals: false,
	}

	for _, option := range options {
		option(cfg)
	}

	w := &waiter{
		waitFuncs:    []WaitFunc{},
		cleanupFuncs: []CleanupFunc{},
	}

	w.ctx, w.cancel = context.WithCancel(cfg.parentCtx)
	if cfg.catchSignals {
		w.ctx, w.cancel = signal.NotifyContext(
			w.ctx,
			os.Interrupt,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
			syscall.SIGHUP,
		)
	}

	return w
}

func (w *waiter) Add(fn ...WaitFunc) {
	w.waitFuncs = append(w.waitFuncs, fn...)
}

func (w *waiter) CancelFunc() context.CancelFunc {
	return w.cancel
}

func (w *waiter) Cleanup(fn ...CleanupFunc) {
	w.cleanupFuncs = append(w.cleanupFuncs, fn...)
}

func (w *waiter) Context() context.Context {
	return w.ctx
}

func (w *waiter) Wait() error {
	g, ctx := errgroup.WithContext(w.ctx)

	g.Go(func() error {
		<-ctx.Done()
		w.cancel()

		return nil
	})

	for _, fn := range w.waitFuncs {
		waitFunc := fn

		g.Go(func() error {
			return waitFunc(ctx)
		})
	}

	for _, fn := range w.cleanupFuncs {
		fn()
	}

	return g.Wait()
}
