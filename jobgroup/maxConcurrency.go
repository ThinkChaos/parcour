package jobgroup

import (
	"context"
)

// NoConcurrencyLimit is used to signify no concurrency limit should apply to the group.
const NoConcurrencyLimit = 0

var _ jobGroup = (*maxConcurrency)(nil)

type maxConcurrency struct {
	withParent

	ch chan struct{}
}

// WithMaxConcurrency returns a new `JobGroup`, child of `parent`, that limits the number of concurrent jobs.
//
// If `max` is `NoConcurrencyLimit`, then a this function is equivalent to `WithParent`.
// Note that the created group will still be subject to any concurrency limits of the parent group.
//
// If multiple jobs are blocked waiting to start and another finishes, the one to actually start is chosen randomly.
func WithMaxConcurrency(parent JobGroup, max uint) JobGroup {
	if max == NoConcurrencyLimit {
		return WithParent(parent)
	}

	return initGroup(parent.Ctx(), &maxConcurrency{
		withParent: newWithParent(parent),

		ch: make(chan struct{}, max),
	})
}

func (g *maxConcurrency) Go(job Job) {
	g.launch(bindJob(g, job))
}

func (g *maxConcurrency) launch(job *boundJob) {
	job.Wrap(func(userJob Job) Job {
		return func(ctx context.Context) error {
			select {
			case g.ch <- struct{}{}:
				defer func() { <-g.ch }()

				return userJob(ctx)

			case <-ctx.Done():
				return newJobNotStartedError(ctx.Err())
			}
		}
	})

	g.withParent.launch(job)
}
