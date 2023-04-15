package jobgroup

import (
	"context"
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("withParent", func() {
	var (
		sutParent     *MockjobGroup
		sutParentCtx  *identifiableContext
		sutParentImpl jobGroup

		sut *withParent
	)

	BeforeEach(func() {
		ctrl := gomock.NewController(GinkgoT())

		sutParentCtx = newIdentifiableContext(context.Background())

		sutParent = NewMockjobGroup(ctrl)

		sutParent.EXPECT().
			Ctx().
			Return(sutParentCtx)

		group, _ := WithContext(sutParentCtx)
		sutParentImpl = downcastGroup(group)

		DeferCleanup(sutParentImpl.Close)
	})

	JustBeforeEach(func() {
		sut = WithParent(sutParent).(*withParent)
		Expect(sut).ShouldNot(BeNil())
		Expect(sut.parent).ShouldNot(BeNil())
		Expect(sut.parent).Should(BeIdenticalTo(sutParent))

		DeferCleanup(sut.Close)
	})

	Describe("WithParent", func() {
		It("should create a child group", func() {
			By("the child's context is a child of the parent's", func() {
				Expect(sutParentCtx.IsParentOf(sut.Ctx())).Should(BeTrue())
			})
		})
	})

	Describe("Cancel", func() {
		It("should cancel its own context", func() {
			sut.Cancel()

			Expect(sut.Ctx().Err()).ShouldNot(Succeed())
			Expect(sutParentCtx.Err()).Should(Succeed())
		})
	})

	Describe("Close", func() {
		It("should cancel its own context", func() {
			sut.Close()

			Expect(sut.Ctx().Err()).ShouldNot(Succeed())
			Expect(sutParentCtx.Err()).Should(Succeed())
		})

		It("should propagate errors", func() {
			expectedErr := errors.New("expected error")

			sutParent.EXPECT().
				launch(gomock.Any()).
				Do(sutParentImpl.launch)

			sut.Go(func(ctx context.Context) error {
				return expectedErr
			})

			var err error

			sutParent.EXPECT().
				saveErr(gomock.Any()).
				Do(func(toSave error) { err = toSave })

			sut.Close()

			Expect(err).ShouldNot(BeIdenticalTo(expectedErr))
		})

		It("should propagate panics", func() {
			const expectedVal = "panic value"

			sutParent.EXPECT().
				launch(gomock.Any()).
				Do(sutParentImpl.launch)

			sut.Go(func(ctx context.Context) error {
				panic(expectedVal)
			})

			Expect(sut.Close).To(PanicWith(expectedVal))
		})
	})

	Describe("Go", func() {
		It("should call parent.launch", func() {
			sutParent.EXPECT().
				launch(gomock.Any()).
				Do(sutParentImpl.launch)

			sut.Go(func(ctx context.Context) error {
				return nil
			})

			Expect(sut.Wait()).Should(Succeed())
		})

		It("should use its own context", func() {
			sutParent.EXPECT().
				launch(gomock.Any()).
				Do(sutParentImpl.launch)

			sut.Go(func(ctx context.Context) error {
				defer GinkgoRecover()

				Expect(ctx).Should(BeIdenticalTo(sut.Ctx()))

				return nil
			})

			Expect(sut.Wait()).Should(Succeed())
		})

		It("should save job errors", func() {
			expectedErr := errors.New("expected error")

			sutParent.EXPECT().
				launch(gomock.Any()).
				Do(sutParentImpl.launch)

			sut.Go(func(ctx context.Context) error {
				return expectedErr
			})

			Expect(sut.Wait()).Should(MatchError(expectedErr))
		})
	})

	Describe("Wait", func() {
		It("should not propagate errors to the parent", func() {
			expectedErr := errors.New("expected error")

			sutParent.EXPECT().
				launch(gomock.Any()).
				Do(sutParentImpl.launch)

			sut.Go(func(ctx context.Context) error {
				return expectedErr
			})

			Expect(sut.Wait()).Should(MatchError(expectedErr))
			Expect(sutParentImpl.Wait()).Should(Succeed())
		})
	})
})
