package jobgroup

import "context"

var _ jobGroup = (*cancelOnErr)(nil)

type cancelOnErr struct {
	withParent
}

// WithCancelOnError returns a new `JobGroup`, child of `parent`, that cancels its context when a job returns an error.
func WithCancelOnError(parent JobGroup) JobGroup {
	return initGroup(parent.Ctx(), &cancelOnErr{
		withParent: newWithParent(parent),
	})
}

func (g *cancelOnErr) Go(job Job) {
	g.launch(bindJob(g, job))
}

func (g *cancelOnErr) launch(job *boundJob) {
	job.Wrap(func(userJob Job) Job {
		return func(ctx context.Context) error {
			err := userJob(ctx)
			if err != nil {
				g.Cancel()
			}

			return err
		}
	})

	g.withParent.launch(job)
}
