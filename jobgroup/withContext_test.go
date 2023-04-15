package jobgroup

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("withContext", func() {
	var (
		rootCtx *identifiableContext

		sut *withContext
	)

	BeforeEach(func() {
		rootCtx = newIdentifiableContext(context.Background())
	})

	JustBeforeEach(func() {
		group, ctx := WithContext(rootCtx)
		Expect(group).ShouldNot(BeNil())
		Expect(ctx).ShouldNot(BeNil())
		Expect(group.Ctx()).Should(BeIdenticalTo(ctx))
		Expect(rootCtx.IsParentOf(group.Ctx())).Should(BeTrue())

		sut = group.(*withContext)

		DeferCleanup(sut.Close)

		sut.private() // for coverage
	})

	Describe("WithContext", func() {
		It("should create a group with child context", func() {
			By("the context value key should be unique", func() {
				token := new(struct{})

				Expect(token).ShouldNot(BeIdenticalTo(rootCtx.token))
				Expect(sut.Ctx().Value(token)).Should(BeNil())
			})
		})

		When("using a context from a parent group", func() {
			var (
				parent, child       JobGroup
				parentCtx, childCtx *identifiableContext

				childFromCtx context.Context
			)

			BeforeEach(func() {
				var ctx context.Context

				parentCtx = newIdentifiableContext(rootCtx)

				parent, ctx = WithContext(parentCtx)
				Expect(parent).ShouldNot(BeNil())
				Expect(parent.Ctx()).Should(BeIdenticalTo(ctx))
				Expect(parentCtx.IsParentOf(parent.Ctx())).Should(BeTrue())

				DeferCleanup(parent.Close)

				childFromCtx = parent.Ctx()
			})

			JustBeforeEach(func() {
				var ctx context.Context

				childCtx = newIdentifiableContext(childFromCtx)

				child, ctx = WithContext(childCtx)
				Expect(child).ShouldNot(BeNil())
				Expect(child.Ctx()).Should(BeIdenticalTo(ctx))
				Expect(childCtx.IsParentOf(child.Ctx())).Should(BeTrue())
			})

			It("should recover the group", func() {
				withParent, ok := child.(*withParent)
				Expect(ok).Should(BeTrue())
				Expect(withParent).ShouldNot(BeNil())

				Expect(withParent.parent).Should(BeIdenticalTo(parent))
			})

			When("using a derived context", func() {
				var derivedCtx *identifiableContext

				BeforeEach(func() {
					derivedCtx = newIdentifiableContext(parent.Ctx())
					childFromCtx = derivedCtx
				})

				It("should still use the given context", func() {
					Expect(derivedCtx.IsParentOf(child.Ctx()))
				})
			})
		})
	})

	Describe("Cancel", func() {
		It("ends the group context", func() {
			sut.Cancel()

			Expect(sut.Ctx().Err()).ShouldNot(Succeed())
			Expect(rootCtx.Err()).Should(Succeed())
		})
	})

	Describe("Close", func() {
		It("ends the group context", func() {
			sut.Close()

			Expect(sut.Ctx().Err()).ShouldNot(Succeed())
			Expect(rootCtx.Err()).Should(Succeed())
		})

		It("can be called multiple times", func() {
			sut.Close()

			Expect(sut.Ctx().Err()).ShouldNot(Succeed())
			Expect(rootCtx.Err()).Should(Succeed())

			sut.Close()
			sut.Close()
		})

		It("waits for any jobs", func(testCtx context.Context) {
			events := make(chan string)
			jobCtx, jobEnd := context.WithCancel(testCtx)

			sut.Go(func(context.Context) error {
				return blockUntilCtxDone(jobCtx)
			})

			go func() {
				events <- "close start"
				sut.Close()
				events <- "close end"
			}()

			Eventually(testCtx, events).Should(Receive(Equal("close start")))
			Consistently(events, 20*time.Millisecond).ShouldNot(Receive())

			// group context should not be cancelled before jobs finish
			Expect(sut.Ctx().Err()).Should(Succeed())
			Expect(rootCtx.Err()).Should(Succeed())

			jobEnd()

			Eventually(testCtx, events).Should(Receive(Equal("close end")))

			// group context should now be cancelled
			Expect(sut.Ctx().Err()).ShouldNot(Succeed())
			Expect(rootCtx.Err()).Should(Succeed())
		}, SpecTimeout(100*time.Millisecond*timeoutFactor))

		It("panics if a job returns an error", func() {
			expectedErr := errors.New("expected error")

			sut.Go(func(context.Context) error {
				return expectedErr
			})

			Expect(sut.Close).Should(PanicWith(MatchError(expectedErr)))
		})

		It("panics if a job panics", func() {
			expectedErr := errors.New("expected error")

			sut.Go(func(context.Context) error {
				panic(expectedErr)
			})

			Expect(sut.Close).Should(PanicWith(MatchError(expectedErr)))
		})
	})

	Describe("Go", func() {
		It("should not block the current goroutine", func(testCtx context.Context) {
			sut, _ := WithContext(testCtx)
			defer sut.Close()

			sut.Go(blockUntilCtxDone)

			// If the testCtx is canceled, it means the job blocked the current
			// goroutine until the test timed-out.
			Expect(testCtx.Err()).Should(Succeed())

			sut.Cancel() // unblock the job

			Expect(sut.Wait()).Should(MatchError(context.Canceled))
		}, SpecTimeout(100*time.Millisecond*timeoutFactor))

		It("doesn't run the job if the context is already done", func(testCtx context.Context) {
			ctx, cancel := context.WithCancel(testCtx)

			sut, _ := WithContext(ctx)
			defer sut.Close()

			cancel()

			sut.Go(func(ctx context.Context) error {
				Fail("job should not run")

				return nil
			})

			Expect(sut.Wait()).Should(MatchError(context.Canceled))
		}, SpecTimeout(100*time.Millisecond*timeoutFactor))
	})

	Describe("Wait", func() {
		It("doesn't block when no job was started", func() {
			Expect(sut.Wait()).Should(Succeed())
		})

		It("blocks until jobs end", func(testCtx context.Context) {
			events := make(chan string)
			jobCtx, jobEnd := context.WithCancel(testCtx)

			for i := 0; i < 25; i++ {
				sut.Go(func(ctx context.Context) error {
					return blockUntilCtxDone(jobCtx)
				})
			}

			go func() {
				events <- "close start"
				sut.Close()
				events <- "close end"
			}()

			Eventually(testCtx, events).Should(Receive(Equal("close start")))
			Consistently(events, 20*time.Millisecond).ShouldNot(Receive())

			jobEnd()

			Eventually(testCtx, events).Should(Receive(Equal("close end")))

			Expect(sut.Ctx().Err()).ShouldNot(Succeed())
			Expect(rootCtx.Err()).Should(Succeed())
		}, SpecTimeout(100*time.Millisecond*timeoutFactor))

		It("can be called multiple times sequentially", func() {
			sut.Go(func(ctx context.Context) error {
				return nil
			})

			Expect(sut.Wait()).Should(Succeed())
			Expect(sut.Wait()).Should(Succeed())
		})

		It("can be called multiple times concurrently", func(testCtx context.Context) {
			sut.Go(blockUntilCtxDone)

			waiting := make(chan struct{})
			waited := make(chan struct{})
			wait := func() {
				defer GinkgoRecover()

				waiting <- struct{}{}
				Expect(sut.Wait()).Should(Succeed())
				waited <- struct{}{}
			}

			go wait()
			go wait()

			Eventually(testCtx, waiting).Should(Receive())
			Eventually(testCtx, waiting).Should(Receive())

			sut.Cancel()

			Eventually(testCtx, waited).Should(Receive())
			Eventually(testCtx, waited).Should(Receive())
		}, SpecTimeout(100*time.Millisecond*timeoutFactor))

		It("propagates the error for a single job", func(testCtx context.Context) {
			expectedErr := errors.New("expected error")

			group, _ := WithContext(testCtx)
			defer group.Close()

			group.Go(func(ctx context.Context) error {
				return expectedErr
			})

			Expect(group.Wait()).Should(MatchError(expectedErr))
		}, SpecTimeout(100*time.Millisecond*timeoutFactor))

		It("propagates the panic for a single job", func(testCtx context.Context) {
			const expectedVal = "panic value"

			group, _ := WithContext(testCtx)
			defer group.Close()

			group.Go(func(ctx context.Context) error {
				panic(expectedVal)
			})

			Expect(func() { _ = group.Wait() }).To(PanicWith(expectedVal))
		}, SpecTimeout(100*time.Millisecond*timeoutFactor))

		It("only propagates an error once", func(testCtx context.Context) {
			expectedErr := errors.New("expected error")

			group, _ := WithContext(testCtx)
			defer group.Close()

			group.Go(func(ctx context.Context) error {
				return expectedErr
			})

			Expect(group.Wait()).Should(MatchError(expectedErr))
			Expect(group.Wait()).Should(Succeed())
		}, SpecTimeout(100*time.Millisecond*timeoutFactor))

		It("only propagates a panic once", func(testCtx context.Context) {
			const expectedVal = "panic value"

			group, _ := WithContext(testCtx)
			defer group.Close()

			group.Go(func(ctx context.Context) error {
				panic(expectedVal)
			})

			Expect(func() { _ = group.Wait() }).To(PanicWith(expectedVal))
			Expect(group.Wait()).Should(Succeed())
		}, SpecTimeout(100*time.Millisecond*timeoutFactor))

		It("propagates an error even if other jobs succeed", func(testCtx context.Context) {
			expectedErr := errors.New("expected error")

			group, _ := WithContext(testCtx)
			defer group.Close()

			group.Go(func(ctx context.Context) error {
				return nil
			})

			group.Go(func(ctx context.Context) error {
				return expectedErr
			})

			group.Go(func(ctx context.Context) error {
				return nil
			})

			Expect(group.Wait()).Should(MatchError(expectedErr))
		}, SpecTimeout(100*time.Millisecond*timeoutFactor))

		It("propagates a panic even if other jobs succeed", func(testCtx context.Context) {
			const expectedVal = "panic value"

			group, _ := WithContext(testCtx)
			defer group.Close()

			group.Go(func(ctx context.Context) error {
				return nil
			})

			group.Go(func(ctx context.Context) error {
				panic(expectedVal)
			})

			group.Go(func(ctx context.Context) error {
				return nil
			})

			Expect(func() { _ = group.Wait() }).To(PanicWith(expectedVal))
		}, SpecTimeout(100*time.Millisecond*timeoutFactor))

		It("propagates all errors when multiple jobs fail", func(testCtx context.Context) {
			expectedErr1 := errors.New("expected error 1")
			expectedErr2 := errors.New("expected error 2")

			group, _ := WithContext(testCtx)
			defer group.Close()

			group.Go(func(ctx context.Context) error {
				return nil
			})

			group.Go(func(ctx context.Context) error {
				return expectedErr1
			})

			group.Go(func(ctx context.Context) error {
				return nil
			})

			group.Go(func(ctx context.Context) error {
				return expectedErr2
			})

			Expect(group.Wait()).Should(SatisfyAll(MatchError(expectedErr1), MatchError(expectedErr2)))
		}, SpecTimeout(100*time.Millisecond*timeoutFactor))

		It("propagates all panics when multiple jobs fail", func(testCtx context.Context) {
			const (
				expectedVal1 = "panic value 1"
				expectedVal2 = "panic value 2"
			)

			group, _ := WithContext(testCtx)
			defer group.Close()

			group.Go(func(ctx context.Context) error {
				return nil
			})

			group.Go(func(ctx context.Context) error {
				panic(expectedVal1)
			})

			group.Go(func(ctx context.Context) error {
				return nil
			})

			group.Go(func(ctx context.Context) error {
				panic(expectedVal2)
			})

			Expect(func() { _ = group.Wait() }).To(
				PanicWith(SatisfyAll(HaveLen(2), ContainElements(expectedVal1, expectedVal2))),
			)
		}, SpecTimeout(100*time.Millisecond*timeoutFactor))
	})
})
