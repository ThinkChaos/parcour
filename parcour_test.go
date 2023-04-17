package parcour

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// // Can be used to increase timeout while debugging.
// const timeoutFactor = 1

func TestParcour(t *testing.T) {
	RegisterFailHandler(Fail)

	t.Parallel()

	RunSpecs(t, "Parcour Suite")
}
