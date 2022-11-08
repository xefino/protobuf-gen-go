package gopb

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Create a new test runner we'll use to test all the
// modules in the time package
func TestTime(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Time Suite")
}

// Helper function that generates a UnixTimestamp from a value for seconds and nanoseconds
func generateTimestamp(seconds int64, nanos int32) *UnixTimestamp {
	return &UnixTimestamp{
		Seconds:     seconds,
		Nanoseconds: nanos,
	}
}

// Helper function that generates a UnixDuration from a value for seconds and nanoseconds
func generateDuration(seconds int64, nanos int32) *UnixDuration {
	return &UnixDuration{
		Seconds:     seconds,
		Nanoseconds: nanos,
	}
}
