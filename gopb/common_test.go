package gopb

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Create a new test runner we'll use to test all the
// modules in the common package
func TestCommon(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Common Suite")
}

// Helper function that verifies the fields on a Decimal
func decimalVerifier(exp int32, parts ...int64) func(*Decimal) {
	return func(decimal *Decimal) {
		Expect(decimal.Exp).Should(Equal(exp))
		Expect(decimal.Parts).Should(HaveLen(len(parts)))
		for i, part := range decimal.Parts {
			Expect(part).Should(Equal(parts[i]))
		}
	}
}
