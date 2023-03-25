package jobgroup

import (
	"context"
	"errors"
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
			//   the user given context could be a subcontext.
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
}

func newWithContext() withContext {
	return withContext{
		wg:     sync.WaitGroup{},
		ctx:    nil, // see init
		cancel: nil, // see init
		err:    zync.Mutex[error]{},
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
	g.Cancel()

	_ = g.Wait()
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

func (g *withContext) Wait() error {
	g.wg.Wait()

	var rerr error

	g.err.WithLock(func(err *error) {
		rerr = *err
		*err = nil
	})

	return rerr
}
