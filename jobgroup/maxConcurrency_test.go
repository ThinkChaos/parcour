package jobgroup

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("maxConcurrency", func() {
	var (
		sutParent     *MockjobGroup
		sutParentCtx  *identifiableContext
		sutParentImpl jobGroup

		sut      *maxConcurrency
		sutLimit uint
	)

	BeforeEach(func() {
		ctrl := gomock.NewController(GinkgoT())

		sutParentCtx = newIdentifiableContext(context.Background())
		DeferCleanup(sutParentCtx.Cancel)

		sutParent = NewMockjobGroup(ctrl)

		sutParent.EXPECT().
			Ctx().
			Return(sutParentCtx)

		group, _ := WithContext(sutParentCtx)
		DeferCleanup(group.Close)

		sutParentImpl = downcastGroup(group)

		sutLimit = 16 // arbitrary non zero number
	})

	JustBeforeEach(func() {
		sut = WithMaxConcurrency(sutParent, sutLimit).(*maxConcurrency)
		Expect(sut).ShouldNot(BeNil())
		Expect(sut.parent).ShouldNot(BeNil())
		Expect(sut.parent).Should(BeIdenticalTo(sutParent))

		DeferCleanup(sut.Close)
	})

	Describe("WithMaxConcurrency", func() {
		When("a limit is given", func() {
			BeforeEach(func() {
				sutLimit = uint(rand.Uint32())
			})

			It("should use it", func() {
				Expect(cap(sut.ch)).Should(BeNumerically("==", sutLimit))
			})
		})

		When("the limit is `NoConcurrencyLimit`", func() {
			It("should not limit concurrency", func() {
				sutParent.EXPECT().
					Ctx().
					Return(sutParentCtx)

				group := WithMaxConcurrency(sutParent, NoConcurrencyLimit)
				defer group.Close()

				casted := group.(*withParent)
				Expect(casted).ShouldNot(BeNil())
				Expect(casted.parent).Should(BeIdenticalTo(sutParent))
			})
		})
	})

	Describe("Go", func() {
		When("concurrency limit is 1", func() {
			BeforeEach(func() {
				sutLimit = 1
			})

			It("should call parent.launch", func() {
				sutParent.EXPECT().
					launch(gomock.Any()).
					Do(sutParentImpl.launch)

				sut.Go(func(ctx context.Context) error {
					defer GinkgoRecover()

					return nil
				})
			})

			When("the limit is reached", func() {
				It("enforces the limit", func(testCtx context.Context) {
					events := make(chan string)

					sutParent.EXPECT().
						launch(gomock.Any()).
						Times(2).
						Do(func(job *boundJob) {
							job.Wrap(func(userJob Job) Job {
								return func(ctx context.Context) error {
									events <- "parent.launch" // must be in the job goroutine to not block the test

									return userJob(ctx)
								}
							})

							sutParentImpl.launch(job)
						})

					job1Ctx, job1End := context.WithCancel(testCtx)
					defer job1End()

					sut.Go(func(ctx context.Context) error {
						defer GinkgoRecover()

						events <- "job 1 start"
						err := blockUntilCtxDone(job1Ctx)
						events <- "job 1 end"

						return err
					})

					By("first call to Go runs the job immediately", func() {
						Eventually(testCtx, events).Should(Receive(Equal("parent.launch")))
						Eventually(testCtx, events).Should(Receive(Equal("job 1 start")))
					})

					sut.Go(func(ctx context.Context) error {
						defer GinkgoRecover()

						events <- "job 2 start"
						events <- "job 2 end"

						return nil
					})

					By("second call to Go calls parent.launch immediately", func() {
						Eventually(testCtx, events).Should(Receive(Equal("parent.launch")))
					})

					// No new events expected: job 1 is still running and job 2 cannot start.
					Consistently(events, 20*time.Millisecond).ShouldNot(Receive())

					job1End() // trigger the job end

					// If job 2 didn't wait, this will be job 2 start
					Eventually(testCtx, events).Should(Receive(Equal("job 1 end")))

					Eventually(testCtx, events).Should(Receive(Equal("job 2 start")))
					Eventually(testCtx, events).Should(Receive(Equal("job 2 end")))

					err, ok := sut.WaitCtx(testCtx)
					Expect(ok).Should(BeTrue())
					Expect(err).Should(Succeed())

					err, ok = sutParentImpl.WaitCtx(testCtx)
					Expect(ok).Should(BeTrue())
					Expect(err).Should(Succeed())
				}, SpecTimeout(100*time.Millisecond*timeoutFactor))

				It("enforces the limit on child group jobs", func(testCtx context.Context) {
					sutParent.EXPECT().
						launch(gomock.Any()).
						Times(2).
						Do(sutParentImpl.launch)

					child := WithParent(sut)
					defer child.Close()

					events := make(chan string)
					job1Ctx, job1End := context.WithCancel(testCtx)

					child.Go(func(context.Context) error {
						defer GinkgoRecover()

						events <- "job 1 start"
						err := blockUntilCtxDone(job1Ctx)
						events <- "job 1 end"

						return err
					})

					By("first call to Go runs the job immediately", func() {
						Eventually(testCtx, events).Should(Receive(Equal("job 1 start")))
					})

					child.Go(func(ctx context.Context) error {
						events <- "job 2 start"

						return nil
					})

					// No new events expected: job 1 is still running and job 2 cannot start.
					Consistently(events, 20*time.Millisecond).ShouldNot(Receive())

					job1End()

					Eventually(testCtx, events).Should(Receive(Equal("job 1 end")))
					Eventually(testCtx, events).Should(Receive(Equal("job 2 start")))

					err, ok := child.WaitCtx(testCtx)
					Expect(ok).Should(BeTrue())
					Expect(err).Should(Succeed())
				}, SpecTimeout(100*time.Millisecond*timeoutFactor))

				It("stops blocking if the group is cancelled", func(testCtx context.Context) {
					events := make(chan string)

					sutParent.EXPECT().
						launch(gomock.Any()).
						Do(sutParentImpl.launch)

					sutParent.EXPECT().
						launch(gomock.Any()). // second job
						Do(func(job *boundJob) {
							job.Wrap(func(userJob Job) Job {
								return func(ctx context.Context) error {
									defer GinkgoRecover()

									events <- "maxConcurrency wrapper: start"
									err := userJob(ctx)
									events <- "maxConcurrency wrapper: end"

									Expect(err).Should(MatchError(context.Canceled))

									return err
								}
							})

							sutParentImpl.launch(job)
						})

					sut.Go(func(ctx context.Context) error {
						defer GinkgoRecover()

						events <- "job 1 start"
						return blockUntilCtxDone(ctx)
					})

					By("first call to Go runs the job immediately", func() {
						Eventually(testCtx, events).Should(Receive(Equal("job 1 start")))
					})

					sut.Go(func(ctx context.Context) error {
						defer GinkgoRecover()

						Fail("job 2 should never run")

						return nil
					})

					By("second call to Go launches the job immediately", func() {
						Eventually(testCtx, events).Should(Receive(Equal("maxConcurrency wrapper: start")))
					})

					// No new events expected: job 1 is still running and job 2 cannot start.
					Consistently(events, 20*time.Millisecond).ShouldNot(Receive())

					sut.Cancel() // difference with test above: cancel `sut`

					Eventually(testCtx, events).Should(Receive(Equal("maxConcurrency wrapper: end")))

					err, ok := sut.WaitCtx(testCtx)
					Expect(ok).Should(BeTrue()) // false is test timeout
					Expect(err).ShouldNot(Succeed())
					Expect(errors.As(err, new(*JobNotStartedError))).Should(BeTrue())
				}, SpecTimeout(500*time.Millisecond*timeoutFactor))
			})
		})
	})
})
