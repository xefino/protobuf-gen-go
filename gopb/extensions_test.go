package gopb

import (
	"math"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UnixTimestamp Extensions Tests", func() {

	// Test that the Now function works as expected
	It("Now - Works", func() {
		timestamp := Now()
		Expect(timestamp.Seconds).ShouldNot(BeZero())
		Expect(timestamp.Nanoseconds).ShouldNot(BeZero())
	})

	// Test that the NewUnixTimestamp function creates a valid UnixTimestamp from a time.Time
	It("NewUnixTimestamp - Works", func() {

		// Create a timestamp from a specific date
		timestamp := NewUnixTimestamp(time.Date(2022, time.June, 1, 23, 59, 53, 983651350, time.UTC))

		// Verify that the number of seconds and nanoseconds is correct
		Expect(timestamp.Seconds).Should(Equal(int64(1654127993)))
		Expect(timestamp.Nanoseconds).Should(Equal(int32(983651350)))
	})

	// Test that the AsTime function creates a time from a valid timestamp
	It("AsTime - Works", func() {

		// First, create a timestamp with a set number of seconds and nanoseconds
		timestamp := UnixTimestamp{
			Seconds:     1654127993,
			Nanoseconds: 983651350,
		}

		// Next, convert the timestamp to a time object
		t := timestamp.AsTime()

		// Finally, verify the fields on the time
		Expect(t.Year()).Should(Equal(2022))
		Expect(t.Month()).Should(Equal(time.June))
		Expect(t.Day()).Should(Equal(1))
		Expect(t.Hour()).Should(Equal(23))
		Expect(t.Minute()).Should(Equal(59))
		Expect(t.Second()).Should(Equal(53))
		Expect(t.Nanosecond()).Should(Equal(983651350))
		Expect(t.Location()).Should(Equal(time.UTC))
	})

	// Tests that the Equals function works under various data conditions
	DescribeTable("Equals - Conditions",
		func(rhs *UnixTimestamp, lhs *UnixTimestamp, equal bool) {
			Expect(rhs.Equals(lhs)).Should(Equal(equal))
		},
		Entry("RHS is nil - False", nil, generateTimestamp(1655510000, 900838091), false),
		Entry("LHS is nil - False", generateTimestamp(1655510399, 900838091), nil, false),
		Entry("Both nil - True", nil, nil, true),
		Entry("RHS.Seconds != LHS.Seconds - False",
			generateTimestamp(1655510399, 900838091), generateTimestamp(1655510000, 900838091), false),
		Entry("RHS.Nanoseconds != LHS.Nanoseconds - False",
			generateTimestamp(1655510399, 900838091), generateTimestamp(1655510399, 0), false),
		Entry("RHS == LHS - True",
			generateTimestamp(1655510399, 900838091), generateTimestamp(1655510399, 900838091), true))

	// Tests that the NotEquals function works under various data conditions
	DescribeTable("NotEquals - Conditions",
		func(rhs *UnixTimestamp, lhs *UnixTimestamp, notEqual bool) {
			Expect(rhs.NotEquals(lhs)).Should(Equal(notEqual))
		},
		Entry("RHS is nil - True", nil, generateTimestamp(1655510000, 900838091), true),
		Entry("LHS is nil - True", generateTimestamp(1655510399, 900838091), nil, true),
		Entry("Both nil - False", nil, nil, false),
		Entry("RHS == LHS - False",
			generateTimestamp(1655510399, 900838091), generateTimestamp(1655510399, 900838091), false),
		Entry("RHS.Nanoseconds != LHS.Nanoseconds - True",
			generateTimestamp(1655510399, 900838091), generateTimestamp(1655510399, 0), true),
		Entry("RHS.Seconds != LHS.Seconds - True",
			generateTimestamp(1655510399, 900838091), generateTimestamp(1655510000, 900838091), true))

	// Tests that the GreaterThan function works under various data conditions
	DescribeTable("GreaterThan - Conditions",
		func(rhs *UnixTimestamp, lhs *UnixTimestamp, greaterThan bool) {
			Expect(rhs.GreaterThan(lhs)).Should(Equal(greaterThan))
		},
		Entry("RHS is nil - False", nil, generateTimestamp(1655510000, 900838091), false),
		Entry("LHS is nil - True", generateTimestamp(1655510399, 900838091), nil, true),
		Entry("Both nil - False", nil, nil, false),
		Entry("RHS == LHS - False",
			generateTimestamp(1655510399, 900838091), generateTimestamp(1655510399, 900838091), false),
		Entry("RHS.Nanoseconds < LHS.Nanoseconds - False",
			generateTimestamp(1655510399, 0), generateTimestamp(1655510399, 900838091), false),
		Entry("RHS.Nanoseconds > LHS.Nanoseconds - True",
			generateTimestamp(1655510399, 900838091), generateTimestamp(1655510399, 0), true),
		Entry("RHS.Seconds < LHS.Seconds - False",
			generateTimestamp(1655510000, 900838091), generateTimestamp(1655510399, 900838091), false),
		Entry("RHS.Seconds > LHS.Seconds - True",
			generateTimestamp(1655510399, 900838091), generateTimestamp(1655510000, 900838091), true))

	// Tests that the GreaterThanOrEqualTo function works under various data conditions
	DescribeTable("GreaterThanOrEqualTo - Conditions",
		func(rhs *UnixTimestamp, lhs *UnixTimestamp, gte bool) {
			Expect(rhs.GreaterThanOrEqualTo(lhs)).Should(Equal(gte))
		},
		Entry("RHS is nil - False", nil, generateTimestamp(1655510000, 900838091), false),
		Entry("LHS is nil - True", generateTimestamp(1655510399, 900838091), nil, true),
		Entry("Both nil - True", nil, nil, true),
		Entry("RHS == LHS - True",
			generateTimestamp(1655510399, 900838091), generateTimestamp(1655510399, 900838091), true),
		Entry("RHS.Nanoseconds < LHS.Nanoseconds - False",
			generateTimestamp(1655510399, 0), generateTimestamp(1655510399, 900838091), false),
		Entry("RHS.Nanoseconds > LHS.Nanoseconds - True",
			generateTimestamp(1655510399, 900838091), generateTimestamp(1655510399, 0), true),
		Entry("RHS.Seconds < LHS.Seconds - False",
			generateTimestamp(1655510000, 900838091), generateTimestamp(1655510399, 900838091), false),
		Entry("RHS.Seconds > LHS.Seconds - True",
			generateTimestamp(1655510399, 900838091), generateTimestamp(1655510000, 900838091), true))

	// Tests that the LessThan function works under various data conditions
	DescribeTable("LessThan - Conditions",
		func(rhs *UnixTimestamp, lhs *UnixTimestamp, lt bool) {
			Expect(rhs.LessThan(lhs)).Should(Equal(lt))
		},
		Entry("RHS is nil - True", nil, generateTimestamp(1655510000, 900838091), true),
		Entry("LHS is nil - False", generateTimestamp(1655510399, 900838091), nil, false),
		Entry("Both nil - False", nil, nil, false),
		Entry("RHS == LHS - False",
			generateTimestamp(1655510399, 900838091), generateTimestamp(1655510399, 900838091), false),
		Entry("RHS.Nanoseconds < LHS.Nanoseconds - True",
			generateTimestamp(1655510399, 0), generateTimestamp(1655510399, 900838091), true),
		Entry("RHS.Nanoseconds > LHS.Nanoseconds - False",
			generateTimestamp(1655510399, 900838091), generateTimestamp(1655510399, 0), false),
		Entry("RHS.Seconds < LHS.Seconds - True",
			generateTimestamp(1655510000, 900838091), generateTimestamp(1655510399, 900838091), true),
		Entry("RHS.Seconds > LHS.Seconds - False",
			generateTimestamp(1655510399, 900838091), generateTimestamp(1655510000, 900838091), false))

	// Tests that the LessThanOrEqualTo function works under various data conditions
	DescribeTable("LessThanOrEqualTo - Condition",
		func(rhs *UnixTimestamp, lhs *UnixTimestamp, lte bool) {
			Expect(rhs.LessThanOrEqualTo(lhs)).Should(Equal(lte))
		},
		Entry("RHS is nil - True", nil, generateTimestamp(1655510000, 900838091), true),
		Entry("LHS is nil - False", generateTimestamp(1655510399, 900838091), nil, false),
		Entry("Both nil - True", nil, nil, true),
		Entry("RHS == LHS - True",
			generateTimestamp(1655510399, 900838091), generateTimestamp(1655510399, 900838091), true),
		Entry("RHS.Nanoseconds < LHS.Nanoseconds - True",
			generateTimestamp(1655510399, 0), generateTimestamp(1655510399, 900838091), true),
		Entry("RHS.Nanoseconds > LHS.Nanoseconds - False",
			generateTimestamp(1655510399, 900838091), generateTimestamp(1655510399, 0), false),
		Entry("RHS.Seconds < LHS.Seconds - True",
			generateTimestamp(1655510000, 900838091), generateTimestamp(1655510399, 900838091), true),
		Entry("RHS.Seconds > LHS.Seconds - False",
			generateTimestamp(1655510399, 900838091), generateTimestamp(1655510000, 900838091), false))

	// Tests that the Add function works under various conditions
	DescribeTable("Add - Works",
		func(rhs *UnixTimestamp, lhs *UnixTimestamp, expected *UnixTimestamp) {
			Expect(rhs.Add(lhs)).Should(Equal(expected))
		},
		Entry("LHS is nil - Works", generateTimestamp(1655510000, 900838091),
			nil, generateTimestamp(1655510000, 900838091)),
		Entry("Nanoseconds < 1 second - Works", generateTimestamp(1655510000, 900838091),
			generateTimestamp(100, 1000), generateTimestamp(1655510100, 900839091)),
		Entry("Nanoseconds > 1 second - Works", generateTimestamp(1655510000, 900838091),
			generateTimestamp(1655510000, 900838091), generateTimestamp(3311020001, 801676182)),
		Entry("Nanoseconds < 0 - Works", generateTimestamp(1655510000, 900838091),
			generateTimestamp(1655510000, -999999999), generateTimestamp(3311019999, 900838092)))

	// Test the conditions describing how the AddDate function works
	DescribeTable("AddDate - Works",
		func(years int, months int, days int, expected *UnixTimestamp) {
			Expect(generateTimestamp(1655510000, 900838091).AddDate(years, months, days)).Should(Equal(expected))
		},
		Entry("Time is positive - Works", 1, 1, 10, generateTimestamp(1690502000, 900838091)),
		Entry("Time is negative - Works", -1, -1, -10, generateTimestamp(1620431600, 900838091)))

	// Tests that the AddDuration function works under various conditions
	DescribeTable("AddDuration - Works",
		func(duration time.Duration, result *UnixTimestamp) {
			timestamp := generateTimestamp(1655510000, 900838091)
			Expect(timestamp.AddDuration(duration)).Should(Equal(result))
		},
		Entry("Add nanoseconds - Works", 15*time.Nanosecond, generateTimestamp(1655510000, 900838106)),
		Entry("Add microseconds - Works", 15*time.Microsecond, generateTimestamp(1655510000, 900853091)),
		Entry("Add milliseconds - Works", 15*time.Millisecond, generateTimestamp(1655510000, 915838091)),
		Entry("Add seconds - Works", 15*time.Second, generateTimestamp(1655510015, 900838091)),
		Entry("Add minutes - Works", 15*time.Minute, generateTimestamp(1655510900, 900838091)),
		Entry("Add hours - Works", 15*time.Hour, generateTimestamp(1655564000, 900838091)),
		Entry("Nanoseconds > 1 second - Works", 15*time.Hour+100*time.Millisecond,
			generateTimestamp(1655564001, 838091)),
		Entry("Nanoseconds < 0 - Works", -999*time.Millisecond, generateTimestamp(1655509999, 901838092)))

	// Tests that the AddUnixDuration function works under various conditions
	DescribeTable("AddUnixDuration - Works",
		func(duration *UnixDuration, result *UnixTimestamp) {
			timestamp := generateTimestamp(1655510000, 900838091)
			Expect(timestamp.AddUnixDuration(duration)).Should(Equal(result))
		},
		Entry("LHS is nil - Works", nil, generateTimestamp(1655510000, 900838091)),
		Entry("Nanoseconds < 1 second - Works", generateDuration(100, 1000),
			generateTimestamp(1655510100, 900839091)),
		Entry("Nanoseconds > 1 second - Works", generateDuration(1655510000, 900838091),
			generateTimestamp(3311020001, 801676182)),
		Entry("Nanoseconds < 0 - Works", generateDuration(1655510000, -999999999),
			generateTimestamp(3311019999, 900838092)))

	// Tests the conditions determining whether IsValid will return true or false
	DescribeTable("IsValid - Conditions",
		func(timestamp *UnixTimestamp, result bool) {
			Expect(timestamp.IsValid()).Should(Equal(result))
		},
		Entry("Timestamp is nil - False", nil, false),
		Entry("Seconds < Minimum Timestamp - False",
			&UnixTimestamp{Seconds: -62135596801, Nanoseconds: 983651350}, false),
		Entry("Seconds > Maximum Timestamp - False",
			&UnixTimestamp{Seconds: 253402300800, Nanoseconds: 983651350}, false),
		Entry("Nanoseconds > 1 second - False",
			&UnixTimestamp{Seconds: 1654127993, Nanoseconds: 1000000000}, false),
		Entry("Nanoseconds negative - False",
			&UnixTimestamp{Seconds: 1654127993, Nanoseconds: -1}, false),
		Entry("Valid - True", &UnixTimestamp{Seconds: 1654127993, Nanoseconds: 983651350}, true))

	// Tests the conditions describing what is returned when CheckValid is called
	// with timestamps of various types
	DescribeTable("CheckValid - Conditions",
		func(timestamp *UnixTimestamp, hadError bool, message string) {
			err := timestamp.CheckValid()
			if hadError {
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal(message))
			} else {
				Expect(err).ShouldNot(HaveOccurred())
			}
		},
		Entry("Timestamp is nil - False", nil, true, "invalid nil Timestamp"),
		Entry("Seconds < Minimum Timestamp - False",
			&UnixTimestamp{Seconds: -62135596801, Nanoseconds: 983651350}, true,
			"timestamp (-62135596801, 983651350) before 0001-01-01"),
		Entry("Seconds > Maximum Timestamp - False",
			&UnixTimestamp{Seconds: 253402300800, Nanoseconds: 983651350}, true,
			"timestamp (253402300800, 983651350) after 9999-12-31"),
		Entry("Nanoseconds > 1 second - False",
			&UnixTimestamp{Seconds: 1654127993, Nanoseconds: 1000000000}, true,
			"timestamp (1654127993, 1000000000) has out-of-range nanos"),
		Entry("Nanoseconds negative - False",
			&UnixTimestamp{Seconds: 1654127993, Nanoseconds: -1}, true,
			"timestamp (1654127993, -1) has out-of-range nanos"),
		Entry("Valid - True", &UnixTimestamp{Seconds: 1654127993, Nanoseconds: 983651350}, false, ""))
})

var _ = Describe("UnixDuration Extensions Tests", func() {

	// Test that the NewUnixDuration function creates a valid UnixDuration from a time.Duration
	It("NewUnixDuration - Works", func() {

		// Create a Unix duration from a specific duration
		duration := NewUnixDuration(31*24*time.Hour + 15*time.Millisecond)

		// Verify that the number of seconds and nanoseconds is correct
		Expect(duration.Seconds).Should(Equal(int64(2678400)))
		Expect(duration.Nanoseconds).Should(Equal(int32(15000000)))
	})

	// Tests the conditions under which the AsDuration function will return an error
	DescribeTable("AsDuration - Failures",
		func(duration *UnixDuration, value int64, message string) {
			dur, err := duration.AsDuration()
			Expect(int64(dur)).Should(Equal(value))
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(Equal(message))
		},
		Entry("Seconds overflow - Error", generateDuration(1<<60, 0), int64(0), "Seconds count was malformed"),
		Entry("Underflow error - Error", generateDuration(-9223372036, -1000000000), int64(math.MinInt64), "Duration underflow"),
		Entry("Overflow error - Error", generateDuration(9223372036, 1000000000), int64(math.MaxInt64), "Duration overflow"))

	// Tests that, if no error occurs, then calling the AsDuration function will return the UnixDuration
	// as a time.Duration object
	It("AsDuration - Works", func() {
		uDur := UnixDuration{Seconds: 2678400, Nanoseconds: 15000000}
		dur, err := uDur.AsDuration()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(dur).Should(Equal(31*24*time.Hour + 15*time.Millisecond))
	})

	// Tests the conditions determining whether IsValid will return true or false
	DescribeTable("IsValid - Conditions",
		func(duration *UnixDuration, result bool) {
			Expect(duration.IsValid()).Should(Equal(result))
		},
		Entry("Duration is nil - False", nil, false),
		Entry("Seconds < -10,000 years - False", generateDuration(-315576000001, 0), false),
		Entry("Seconds > 10,000 years - False", generateDuration(315576000001, 0), false),
		Entry("Nanoseconds <= -1e9 - False", generateDuration(2678400, -1000000000), false),
		Entry("Nanoseconds >= 1e9 - False", generateDuration(2678400, 1000000000), false),
		Entry("Seconds > 0, Nanoseconds < 0 - False", generateDuration(2678400, -1000), false),
		Entry("Seconds < 0, Nanoseconds > 0 - False", generateDuration(-2678400, 1000), false),
		Entry("Valid - True", generateDuration(2678400, 1000), true))

	// Tests the conditions describing what is returned when CheckValid is called
	// with durations of various types
	DescribeTable("CheckValid - Conditions",
		func(duration *UnixDuration, hadError bool, message string) {
			err := duration.CheckValid()
			if hadError {
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal(message))
			} else {
				Expect(err).ShouldNot(HaveOccurred())
			}
		},
		Entry("Duration is nil - False", nil, true, "invalid nil Duration"),
		Entry("Seconds < -10,000 years - False", generateDuration(-315576000001, 0), true,
			"duration (-315576000001, 0) exceeds -10000 years"),
		Entry("Seconds > 10,000 years - False", generateDuration(315576000001, 0), true,
			"duration (315576000001, 0) exceeds +10000 years"),
		Entry("Nanoseconds <= -1e9 - False", generateDuration(2678400, -1000000000), true,
			"duration (2678400, -1000000000) has out-of-range nanos"),
		Entry("Nanoseconds >= 1e9 - False", generateDuration(2678400, 1000000000), true,
			"duration (2678400, 1000000000) has out-of-range nanos"),
		Entry("Seconds > 0, Nanoseconds < 0 - False", generateDuration(2678400, -1000), true,
			"duration (2678400, -1000) has seconds and nanos with different signs"),
		Entry("Seconds < 0, Nanoseconds > 0 - False", generateDuration(-2678400, 1000), true,
			"duration (-2678400, 1000) has seconds and nanos with different signs"),
		Entry("Valid - True", generateDuration(2678400, 1000), false, ""))
})
