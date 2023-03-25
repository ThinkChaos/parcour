package jobgroup

import (
	"context"
)

var _ jobGroup = (*withParent)(nil)

// WithParent creates a `JobGroup` that is a child of the given group.
//
// Any jobs that are part of the returned group are also part of the parent.
// Meaning the parent will also wait for any jobs part of the child group.
func WithParent(parent JobGroup) JobGroup {
	return withParentAndContext(parent, parent.Ctx())
}

type withParent struct {
	withContext

	parent jobGroup
}

//nolint:golint // context as 2nd arg is fine here
//revive:disable:context-as-argument
func withParentAndContext(parent JobGroup, ctx context.Context) JobGroup {
	return initGroup(ctx, toPtr(newWithParent(parent)))
}

func newWithParent(parent JobGroup) withParent {
	return withParent{
		withContext: newWithContext(),

		parent: downcastGroup(parent),
	}
}

func (g *withParent) Go(job Job) {
	g.launch(bindJob(g, job))
}

func (g *withParent) launch(job *boundJob) {
	g.wg.Add(1)
	job.Defer(g.wg.Done)

	g.parent.launch(job)
}
