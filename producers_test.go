package parcour

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/ThinkChaos/parcour/jobgroup"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Producers", func() {
	const nGoroutines = 10000

	var (
		cap int
		grp jobgroup.JobGroup
		sut *Producers[string]
	)

	BeforeEach(func() {
		cap = 0
	})

	JustBeforeEach(func() {
		grp, _ = jobgroup.WithContext(context.Background())
		DeferCleanup(grp.Close)

		sut = NewProducersWithBuffer[string](grp, grp, cap)
		DeferCleanup(sut.Close)
	})

	Describe("NewUnbufferedProducers", func() {
		It("returns Producers with no buffer", func() {
			grp, _ := jobgroup.WithContext(context.Background())
			DeferCleanup(grp.Close)

			sut := NewUnbufferedProducers[struct{}](grp, grp)
			Expect(sut).ShouldNot(BeNil())
			Expect(sut.BufferCap()).Should(Equal(Unbuffered))
		})
	})

	Describe("NewProducersWithBuffer", func() {
		BeforeEach(func() {
			cap = 3
		})

		It("returns Producers with the requested buffer", func() {
			Expect(sut.BufferCap()).Should(Equal(cap))
		})

		It("creates separate JobGroups", func() {
			sut.producersGrp.Cancel()
			Expect(grp.Ctx().Err()).Should(Succeed())
			Expect(sut.consumersGrp.Ctx().Err()).Should(Succeed())

			sut.consumersGrp.Cancel()
			Expect(grp.Ctx().Err()).Should(Succeed())
		})

		It("creates child JobGroups", func() {
			grp.Cancel()
			Expect(sut.producersGrp.Ctx().Err()).Should(MatchError(grp.Ctx().Err()))
			Expect(sut.consumersGrp.Ctx().Err()).Should(MatchError(grp.Ctx().Err()))
		})

		It("panics when cap is negative", func() {
			Expect(func() {
				NewProducersWithBuffer[struct{}](grp, grp, -1)
			}).Should(Panic())
		})
	})

	Describe("NewProducersWithBuffer", func() {
		It("supports SPSC", func() {
			sut.GoProduce(func(ctx context.Context, ch chan<- string) error {
				ch <- "a"
				ch <- "b"
				ch <- "c"

				return nil
			})

			sut.GoConsume(func(ctx context.Context, ch <-chan string) error {
				defer GinkgoRecover()

				Eventually(ch).Should(Receive(Equal("a")))
				Eventually(ch).Should(Receive(Equal("b")))
				Eventually(ch).Should(Receive(Equal("c")))

				return nil
			})

			err := sut.Wait()
			Expect(err).Should(Succeed())
		})

		It("supports MPSC", func() {
			for i := 0; i < nGoroutines; i++ {
				sut.GoProduce(func(ctx context.Context, ch chan<- string) error {
					ch <- "product"

					return nil
				})
			}

			n := 0

			sut.GoConsume(func(ctx context.Context, ch <-chan string) error {
				for range ch {
					n++
				}

				return nil
			})

			err := sut.Wait()
			Expect(err).Should(Succeed())

			Expect(n).Should(Equal(nGoroutines))
		})

		It("supports MPMC", func() {
			for i := 0; i < nGoroutines; i++ {
				sut.GoProduce(func(ctx context.Context, ch chan<- string) error {
					ch <- "product"

					return nil
				})
			}

			var n atomic.Int32

			for i := 0; i < nGoroutines; i++ {
				sut.GoConsume(func(ctx context.Context, ch <-chan string) error {
					for range ch {
						n.Add(1)
					}

					return nil
				})
			}

			err := sut.Wait()
			Expect(err).Should(Succeed())

			Expect(n.Load()).Should(Equal(int32(nGoroutines)))
		})
	})

	Describe("Failures", func() {
		It("wraps producer errors", func() {
			expectedErr := errors.New("expected")

			sut.GoProduce(func(ctx context.Context, ch chan<- string) error {
				return expectedErr
			})

			sut.GoConsume(func(ctx context.Context, ch <-chan string) error {
				defer GinkgoRecover()

				Consistently(ch).ShouldNot(Receive())

				return nil
			})

			err := sut.Wait()
			Expect(err).ShouldNot(Succeed())
			Expect(err).Should(MatchError(expectedErr))

			var typed *ProducersError
			Expect(errors.As(err, &typed)).Should(BeTrue())
		})

		It("wraps consumer errors", func() {
			expectedErr := errors.New("expected")

			sut.GoProduce(func(ctx context.Context, ch chan<- string) error {
				defer GinkgoRecover()

				Consistently(ch).ShouldNot(BeSent("product"))

				return nil
			})

			sut.GoConsume(func(ctx context.Context, ch <-chan string) error {
				return expectedErr
			})

			err := sut.Wait()
			Expect(err).ShouldNot(Succeed())
			Expect(err).Should(MatchError(expectedErr))

			var typed *ConsumersError
			Expect(errors.As(err, &typed)).Should(BeTrue())
		})

		It("returns both producer and consumer errors", func() {
			expectedPErr := errors.New("expected")
			expectedCErr := errors.New("expected")

			sut.GoProduce(func(ctx context.Context, ch chan<- string) error {
				return expectedPErr
			})

			sut.GoConsume(func(ctx context.Context, ch <-chan string) error {
				return expectedCErr
			})

			err := sut.Wait()
			Expect(err).ShouldNot(Succeed())
			Expect(err).Should(MatchError(expectedPErr))
			Expect(err).Should(MatchError(expectedCErr))

			var pTyped *ProducersError
			Expect(errors.As(err, &pTyped)).Should(BeTrue())

			var cTyped *ConsumersError
			Expect(errors.As(err, &cTyped)).Should(BeTrue())
		})

		It("propagates producer panics", func() {
			type token struct{}
			expectedErr := errors.New("expected")

			sut.GoProduce(func(ctx context.Context, ch chan<- string) error {
				panic(expectedErr)
			})

			sut.GoConsume(func(ctx context.Context, ch <-chan string) error {
				defer GinkgoRecover()

				Consistently(ch).ShouldNot(Receive())

				return nil
			})

			Expect(func() { _ = sut.Wait() }).Should(PanicWith(MatchError(expectedErr)))
		})

		It("propagates consumer panics", func() {
			type token struct{}
			expectedErr := errors.New("expected")

			sut.GoProduce(func(ctx context.Context, ch chan<- string) error {
				defer GinkgoRecover()

				Consistently(ch).ShouldNot(BeSent("product"))

				return nil
			})

			sut.GoConsume(func(ctx context.Context, ch <-chan string) error {
				panic(expectedErr)
			})

			Expect(func() { _ = sut.Wait() }).Should(PanicWith(MatchError(expectedErr)))
		})
	})

	Describe("ProducersError", func() {
		Describe("Error", func() {
			It("contains the inner error", func() {
				inner := errors.New("test inner error string")

				sut := wrapProducersError(inner)

				Expect(sut.Error()).Should(ContainSubstring(inner.Error()))
			})
		})

		Describe("Unwrap", func() {
			It("returns the inner error", func() {
				inner := errors.New("test inner error string")

				sut := wrapProducersError(inner).(*ProducersError)

				Expect(sut.Unwrap()).Should(BeIdenticalTo(inner))
			})
		})
	})

	Describe("ConsumersError", func() {
		Describe("Error", func() {
			It("contains the inner error", func() {
				inner := errors.New("test inner error string")

				sut := wrapConsumersError(inner)

				Expect(sut.Error()).Should(ContainSubstring(inner.Error()))
			})
		})

		Describe("Unwrap", func() {
			It("returns the inner error", func() {
				inner := errors.New("test inner error string")

				sut := wrapConsumersError(inner).(*ConsumersError)

				Expect(sut.Unwrap()).Should(BeIdenticalTo(inner))
			})
		})
	})
})
