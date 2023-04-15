package jobgroup

import (
	"context"
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cancelOnErr", func() {
	var (
		sutParent     *MockjobGroup
		sutParentCtx  *identifiableContext
		sutParentImpl jobGroup

		sut *cancelOnErr
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
	})

	JustBeforeEach(func() {
		sut = WithCancelOnError(sutParent).(*cancelOnErr)
		Expect(sut).ShouldNot(BeNil())
		Expect(sut.parent).ShouldNot(BeNil())
		Expect(sut.parent).Should(BeIdenticalTo(sutParent))

		DeferCleanup(sut.Close)
	})

	Describe("Go", func() {
		It("should call parent.launch", func() {
			sutParent.EXPECT().
				launch(gomock.Any()).
				Do(sutParentImpl.launch)

			sut.Go(func(ctx context.Context) error {
				return nil
			})
		})

		It("should not cancel the group when no error occurs", func() {
			sutParent.EXPECT().
				launch(gomock.Any()).
				Do(sutParentImpl.launch)

			sut.Go(func(ctx context.Context) error {
				return nil
			})

			Consistently(sut.Ctx().Err).Should(Succeed())

			Expect(sut.Wait()).Should(Succeed())
		})

		It("should cancel the group when an error occurs", func() {
			expectedErr := errors.New("expected error")

			sutParent.EXPECT().
				launch(gomock.Any()).
				Times(2).
				Do(sutParentImpl.launch)

			sut.Go(blockUntilCtxDone)

			sut.Go(func(ctx context.Context) error {
				return expectedErr
			})

			Eventually(sut.Ctx().Err).ShouldNot(Succeed())

			Expect(sut.Wait()).Should(MatchError(expectedErr))
		})
	})
})
