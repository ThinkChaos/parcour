package jobgroup

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ThinkChaos/parcour/zync"
)

type ctxKeyType struct{}

var (
	ctxKey = new(ctxKeyType) //nolint:gochecknoglobals

	_ jobGroup = (*withContext)(nil)
)

// WithContext creates a `JobGroup` with the given context.
// The group's context is a child of the given context to allow independent waiting and cancellation.
//
// If the given context is that of another `JobGroup`, the returned group will be
// a child of the context's group.
// This allows passing a `JobGroup` through functions that are not `JobGroup` aware.
func WithContext(ctx context.Context) (JobGroup, context.Context) {
	// Allow creating a child group from a parent's context.
	if value := ctx.Value(ctxKey); value != nil {
		if parent, ok := value.(jobGroup); ok {
			// Even if we find a parent, we cannot use `parent.Ctx()`:
			// the user given context could be a subcontext.
			group := withParentAndContext(parent, ctx)

			return group, group.Ctx()
		}
	}

	group := initGroup(ctx, toPtr(newWithContext()))

	return group, group.Ctx()
}

type withContext struct {
	wg     sync.WaitGroup
	ctx    context.Context //nolint:containedctx
	cancel context.CancelFunc
	err    zync.Mutex[error]
	panic  zync.Mutex[[]any]
}

func newWithContext() withContext {
	return withContext{
		wg:     sync.WaitGroup{},
		ctx:    nil, // see init
		cancel: nil, // see init
		err:    zync.Mutex[error]{},
		panic:  zync.Mutex[[]any]{},
	}
}

func (g *withContext) private() {}

// Set `ctx` (and `cancel`) as a second step so we already have a stable address
// that can be stored in the context.
// If we did this from a function returning a `jobGroup` (no pointer), then the
// address would change and the one in the context would point to an old value.
// The alternative to `init` is always storing `jobGroup` as a pointer.
func (g *withContext) init(ctx context.Context, selfWrapped JobGroup) {
	ctx, g.cancel = context.WithCancel(ctx)

	// Store `selfWrapped` so when recovered from the context, the `Go` method is the correct one,
	// and we don't loose the specialties of the group.
	g.ctx = context.WithValue(ctx, ctxKey, selfWrapped)
}

func (g *withContext) Ctx() context.Context {
	return g.ctx
}

func (g *withContext) Cancel() {
	g.cancel()
}

func (g *withContext) Close() {
	err := g.close() // propagates panics
	if err != nil {
		panic(fmt.Errorf("unhandled job error: %w", err))
	}
}

func (g *withContext) close() error {
	defer g.Cancel() // prevent group reuse

	return g.Wait()
}

func (g *withContext) Go(job Job) {
	g.launch(bindJob(g, job))
}

func (g *withContext) launch(job *boundJob) {
	g.wg.Add(1)
	job.Defer(g.wg.Done)

	go job.Main()
}

func (g *withContext) saveErr(err error) {
	g.err.WithLock(func(gErr *error) {
		*gErr = errors.Join(*gErr, err)
	})
}

func (g *withContext) savePanic(value any) {
	g.panic.WithLock(func(gPanic *[]any) {
		*gPanic = append(*gPanic, value)
	})
}

func (g *withContext) Wait() error {
	err, _ := g.WaitCtx(context.Background())

	return err
}

func (g *withContext) WaitCtx(ctx context.Context) (error, bool) {
	wait := make(chan struct{})

	go func() {
		g.wg.Wait()
		close(wait)
	}()

	select {
	case <-wait:
	case <-ctx.Done():
		return nil, false
	}

	// Propagate panic, if any
	if panicVal := takeMutexValue(&g.panic); panicVal != nil {
		if len(panicVal) == 1 {
			panic(panicVal[0])
		}

		panic(panicVal)
	}

	// Propagate error, if any
	err := takeMutexValue(&g.err)

	return err, true
}

func takeMutexValue[T any](mutex *zync.Mutex[T]) T {
	var res, zero T

	mutex.WithLock(func(val *T) {
		res = *val
		*val = zero
	})

	return res
}
