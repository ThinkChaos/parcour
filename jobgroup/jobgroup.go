// Package jobgroup provides building blocks for structured concurrency, and composable concurrency patterns.
package jobgroup

//nolint:lll
//go:generate mockgen -source=jobgroup.go -destination jobgroup_mock_test.go -package jobgroup github.com/ThinkChaos/parcour/jobgroup JobGroup,jobGroup

import (
	"context"
	"fmt"
)

// Job is a function that can run as part of a `JobGroup`.
type Job = func(context.Context) error

// JobGroup allows running a set of jobs and waiting for them to finish.
//
// This is meant to be a building block for structured concurrency, and
// composable concurrency patterns.
type JobGroup interface {
	// Ctx returns the context used by all jobs of the group.
	Ctx() context.Context

	// Cancel cancels the group's context, signaling to jobs they should stop.
	//
	// Calling this function more than once has no effect.
	Cancel()

	// Close cancel's the group context and waits for all jobs.
	// Errors returned from the jobs are discarded.
	//
	// Close must be called for each group.
	// Calling this function more than once has no effect.
	Close()

	// Go starts a job as part of the group. It returns immediately, starting the
	// job in another goroutine.
	//
	// If the job cannot make progress immediately, for example due to a
	// concurrency limit. The job's goroutine blocks until it can advance.
	Go(job Job)

	// Wait blocks until all jobs of the group, and any child groups, are done.
	//
	// It returns all errors of jobs launched on the group directly.
	// It also returns errors of any child group that was not handled by waiting on that group.
	Wait() error

	// private prevents the interface from being implemented outside this package.
	private()
}

// jobGroup is the full interface a JobGroup is required to implement, including private methods.
type jobGroup interface {
	JobGroup

	init(context.Context, JobGroup)

	launch(*boundJob)

	saveErr(error)
}

func downcastGroup(group JobGroup) jobGroup {
	// `JobGroup` cannot be implemented outside of this package, and all types
	// in this package that implement `JobGroup` also implement `jobGroup`.
	casted, ok := group.(jobGroup)
	if !ok {
		panic(fmt.Sprintf("JobGroup does not implement private jobGroup interface: %T", group))
	}

	return casted
}

func initGroup(ctx context.Context, group jobGroup) jobGroup {
	// Despite looking weird, passing `group` twice has a purpose:
	//   - The receiver is always `*withContext`.
	//   - The argument is the group that embeds the receiver. Of type `JobGroup`.
	//     This is the value we want to recover later on since it's the
	//     one that has any specific configuration.
	group.init(ctx, group)

	return group
}

// toPtr is used to make pointer from a function's return value.
// It saves writing a var/assignment.
func toPtr[T any](val T) *T { return &val }

type boundJob struct {
	group jobGroup
	run   Job

	cleanup func()
}

func bindJob(group jobGroup, userJob Job) *boundJob {
	return &boundJob{
		group: group,
		run:   userJob,

		cleanup: func() {}, // simplifies `Defer`
	}
}

func (j *boundJob) Main() {
	defer j.cleanup()

	ctx := j.group.Ctx()

	_ = func() (rerr error) { // use a func so we can defer `saveErr`
		defer func() {
			if rerr != nil {
				// Only save the error on the group the user launched the job on.
				// Error will be propagated to its parent, if any, on `Cancel`.
				j.group.saveErr(rerr)
			}
		}()

		if err := ctx.Err(); err != nil {
			return newJobNotStartedError(err)
		}

		return j.run(ctx)
	}()
}

// Defer is used to defer cleanup for the job.
//
// The given function will be called even if the job doesn't fully lauch
// due to the context being already done.
func (j *boundJob) Defer(cleanup func()) {
	prev := j.cleanup

	j.cleanup = func() {
		defer prev() // ensure it's actually called

		cleanup()
	}
}

// Wrap adds logic to the job's run function.
func (j *boundJob) Wrap(wrap func(userJob Job) Job) {
	j.run = wrap(j.run)
}
