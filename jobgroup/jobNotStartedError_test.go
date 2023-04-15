package jobgroup

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("JobNotStartedError", func() {
	Describe("Error", func() {
		It("contains the inner error", func() {
			inner := errors.New("test inner error string")

			sut := newJobNotStartedError(inner)

			Expect(sut.Error()).Should(ContainSubstring(inner.Error()))
		})
	})

	Describe("Unwrap", func() {
		It("returns the inner error", func() {
			inner := errors.New("test inner error string")

			sut := newJobNotStartedError(inner)

			Expect(sut.Unwrap()).Should(BeIdenticalTo(inner))
		})
	})
})
