package jobgroup

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Can be used to increase timeout while debugging.
const timeoutFactor = 1

func TestJobGroup(t *testing.T) {
	RegisterFailHandler(Fail)

	t.Parallel()

	RunSpecs(t, "JobGroup Suite")
}

var _ = Describe("jobGroup", func() {
	Describe("downcastGroup", func() {
		It("succeeds for a JobGroup", func() {
			group, _ := WithContext(context.Background())
			DeferCleanup(group.Close)

			Expect(downcastGroup(group)).ShouldNot(BeNil())
		})

		It("succeeds for a nested groups", func() {
			group, _ := WithContext(context.Background())
			defer group.Close()

			group = WithMaxConcurrency(group, 1)
			defer group.Close()

			group = WithParent(group)
			defer group.Close()

			Expect(downcastGroup(group)).ShouldNot(BeNil())
		})

		It("panics for a non JobGroup", func() {
			defer func() {
				val := recover()
				Expect(val).ShouldNot(BeNil())
			}()

			ctrl := gomock.NewController(GinkgoT())

			downcastGroup(NewMockJobGroup(ctrl))
		})
	})

	Describe("initGroup", func() {
		It("", func() {
			ctrl := gomock.NewController(GinkgoT())

			group := NewMockjobGroup(ctrl)
			ctx := context.TODO()

			group.EXPECT().
				init(ctx, group)

			Expect(initGroup(ctx, group)).Should(BeIdenticalTo(group))
		})
	})
})

// blockUntilCtxDone blocks the current goroutine until the given `ctx` is done.
func blockUntilCtxDone(ctx context.Context) error {
	select {
	case <-chan struct{}(nil): // never returns
		return errors.New("blockForever did not block forever")

	case <-ctx.Done():
		return nil
	}
}

// identifiableContext is a cancellable context that can verify if another context was derived from it.
type identifiableContext struct {
	context.Context //nolint:containedctx
	Cancel          func()

	token any
}

func newIdentifiableContext(parent context.Context) *identifiableContext {
	type ctxToken struct{}

	token := new(ctxToken)

	ctx := context.WithValue(parent, token, token)
	ctx, cancel := context.WithCancel(ctx)

	return &identifiableContext{
		Context: ctx,
		Cancel:  cancel,

		token: token,
	}
}

func (c *identifiableContext) IsParentOf(ctx context.Context) bool {
	GinkgoHelper()

	x := ctx.Value(c.token)

	return x == c.token
}
