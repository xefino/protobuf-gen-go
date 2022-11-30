package gopb

import (
	"math"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
)

var _ = Describe("Decimal Extensions Tests", func() {

	// Tests that, if the value is greater than 0, then it will be encoded properly
	It("NewFromDecimal - Value greater than 0 - Encoded", func() {

		// First, create our decimal value
		dIn, err := decimal.NewFromString("1234512351234088800000.999")
		Expect(err).ShouldNot(HaveOccurred())

		// Next, attempt to convert it to a new Decimal object
		dOut := NewFromDecimal(dIn)

		// Finally, verify the data
		Expect(dOut).ShouldNot(BeNil())
		Expect(dOut.Exp).Should(Equal(int32(-3)))
		Expect(dOut.Value).Should(HaveLen(2))
		Expect(dOut.Value[0]).Should(Equal(int64(351234088800000999)))
		Expect(dOut.Value[1]).Should(Equal(int64(1234512)))
	})

	// Test that, if the value is less than 0, then it will be encoded properly
	It("NewFromDecimal - Value less than 0 - Encoded", func() {

		// First, create our decimal value
		dIn, err := decimal.NewFromString("-288341660781234512351234088800000.999")
		Expect(err).ShouldNot(HaveOccurred())

		// Next, attempt to convert it to a new Decimal object
		dOut := NewFromDecimal(dIn)

		// Finally, verify the data
		Expect(dOut).ShouldNot(BeNil())
		Expect(dOut.Exp).Should(Equal(int32(-3)))
		Expect(dOut.Value).Should(HaveLen(2))
		Expect(dOut.Value[0]).Should(Equal(int64(-351234088800000999)))
		Expect(dOut.Value[1]).Should(Equal(int64(-288341660781234512)))
	})

	// Tests that the Decimal value will be converted to a decimal.Decimal properly
	It("ToDecimal - Decoded", func() {

		// First, create a valid Decimal value
		dIn := Decimal{Value: []int64{351234088800000999, 1234512}, Exp: -3}

		// Next, attempt to convert the Decimal to a decimal
		dOut := dIn.ToDecimal()

		// Finally, verify the data
		Expect(dOut.String()).Should(Equal("1234512351234088800000.999"))
	})
})

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
		timestamp := generateTimestamp(1654127993, 983651350)

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
	DescribeTable("AddDuration - Conditions",
		func(duration *UnixDuration, result *UnixTimestamp) {
			timestamp := generateTimestamp(1655510000, 900838091)
			Expect(timestamp.AddDuration(duration)).Should(Equal(result))
		},
		Entry("LHS is nil - Works", nil, generateTimestamp(1655510000, 900838091)),
		Entry("Nanoseconds < 1 second - Works", generateDuration(100, 1000),
			generateTimestamp(1655510100, 900839091)),
		Entry("Nanoseconds > 1 second - Works", generateDuration(1655510000, 900838091),
			generateTimestamp(3311020001, 801676182)),
		Entry("Nanoseconds < 0 - Works", generateDuration(1655510000, -999999999),
			generateTimestamp(3311019999, 900838092)))

	// Tests that the NextDay function works as expected
	It("NextDay - Works", func() {
		next := generateTimestamp(1655510000, 900838091).NextDay()
		Expect(next).Should(Equal(generateTimestamp(1655510400, 0)))
	})

	// Tests that the Difference functions works under various conditions
	DescribeTable("Difference - Conditions",
		func(rhs *UnixTimestamp, lhs *UnixTimestamp, result *UnixDuration) {
			Expect(rhs.Difference(lhs)).Should(Equal(result))
		},
		Entry("rhs.Seconds > lhs.Seconds, rhs.Nanoseconds > lhs.Nanoseconds - Works",
			generateTimestamp(1669704178, 500000000), generateTimestamp(1669704177, 0), generateDuration(1, 500000000)),
		Entry("rhs.Seconds > lhs.Seconds, rhs.Nanoseconds < lhs.Nanoseconds - Works",
			generateTimestamp(1669704178, 0), generateTimestamp(1669704177, 500000000), generateDuration(0, 500000000)),
		Entry("rhs.Seconds < lhs.Seconds, rhs.Nanoseconds > lhs.Nanoseconds - Works",
			generateTimestamp(1669704177, 500000000), generateTimestamp(1669704178, 0), generateDuration(0, -500000000)),
		Entry("rhs.Seconds < lhs.Seconds, rhs.Nanoseconds < lhs.Nanoseconds - Works",
			generateTimestamp(1669704177, 0), generateTimestamp(1669704178, 500000000), generateDuration(-1, -500000000)))

	// Tests the conditions describing how the IsWhole function works
	DescribeTable("IsWhole - Conditions",
		func(rhs *UnixTimestamp, lhs time.Duration, result bool) {
			Expect(rhs.IsWhole(lhs)).Should(Equal(result))
		},
		Entry("rhs has no nanoseconds, lhs has no nanoseconds, fits - True",
			generateTimestamp(1669704178, 0), 2*time.Second, true),
		Entry("rhs has no nanoseconds, lhs has no nanoseconds, not fits - False",
			generateTimestamp(1669704178, 0), 3*time.Second, false),
		Entry("rhs has no nanoseconds, lhs has nanoseconds, fits - True",
			generateTimestamp(1669704177, 0), time.Second+500*time.Millisecond, true),
		Entry("rhs has no nanoseconds, lhs has nanoseconds, not fits - False",
			generateTimestamp(1669704177, 0), 2*time.Second+500*time.Millisecond, false),
		Entry("rhs has nanoseconds, lhs has no nanoseconds - False",
			generateTimestamp(1669704178, 500000000), 2*time.Second, false),
		Entry("rhs has nanoseconds, lhs has nanoseconds, fits - True",
			generateTimestamp(1669704178, 500000000), time.Second+500*time.Millisecond, true),
		Entry("rhs has nanoseconds, lhs has nanoseconds, not fits - False",
			generateTimestamp(1669704178, 500000000), 2*time.Second+500*time.Millisecond, false))

	// Tests the conditions describing how the IsWholeUnix function works
	DescribeTable("IsWholeUnix - Conditions",
		func(rhs *UnixTimestamp, lhs *UnixDuration, result bool) {
			Expect(rhs.IsWholeUnix(lhs)).Should(Equal(result))
		},
		Entry("rhs has no nanoseconds, lhs has no nanoseconds, fits - True",
			generateTimestamp(1669704178, 0), generateDuration(2, 0), true),
		Entry("rhs has no nanoseconds, lhs has no nanoseconds, not fits - False",
			generateTimestamp(1669704178, 0), generateDuration(3, 0), false),
		Entry("rhs has no nanoseconds, lhs has nanoseconds, fits - True",
			generateTimestamp(1669704177, 0), generateDuration(1, 500000000), true),
		Entry("rhs has no nanoseconds, lhs has nanoseconds, not fits - False",
			generateTimestamp(1669704177, 0), generateDuration(2, 500000000), false),
		Entry("rhs has nanoseconds, lhs has no nanoseconds - False",
			generateTimestamp(1669704178, 500000000), generateDuration(2, 0), false),
		Entry("rhs has nanoseconds, lhs has nanoseconds, fits - True",
			generateTimestamp(1669704178, 500000000), generateDuration(1, 500000000), true),
		Entry("rhs has nanoseconds, lhs has nanoseconds, not fits - False",
			generateTimestamp(1669704178, 500000000), generateDuration(2, 500000000), false))

	// Tests the conditions determining whether IsValid will return true or false
	DescribeTable("IsValid - Conditions",
		func(timestamp *UnixTimestamp, result bool) {
			Expect(timestamp.IsValid()).Should(Equal(result))
		},
		Entry("Timestamp is nil - False", nil, false),
		Entry("Seconds < Minimum Timestamp - False", generateTimestamp(-62135596801, 983651350), false),
		Entry("Seconds > Maximum Timestamp - False", generateTimestamp(253402300800, 983651350), false),
		Entry("Nanoseconds > 1 second - False", generateTimestamp(1654127993, 1000000000), false),
		Entry("Nanoseconds negative - False", generateTimestamp(1654127993, -1), false),
		Entry("Valid - True", generateTimestamp(1654127993, 983651350), true))

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
		Entry("Seconds < Minimum Timestamp - False", generateTimestamp(-62135596801, 983651350), true,
			"timestamp (-62135596801, 983651350) before 0001-01-01"),
		Entry("Seconds > Maximum Timestamp - False", generateTimestamp(253402300800, 983651350), true,
			"timestamp (253402300800, 983651350) after 9999-12-31"),
		Entry("Nanoseconds > 1 second - False", generateTimestamp(1654127993, 1000000000), true,
			"timestamp (1654127993, 1000000000) has out-of-range nanos"),
		Entry("Nanoseconds negative - False", generateTimestamp(1654127993, -1), true,
			"timestamp (1654127993, -1) has out-of-range nanos"),
		Entry("Valid - True", generateTimestamp(1654127993, 983651350), false, ""))

	// Test that the ToDate function converts the UnixTimestamp to a string describing the date associated
	// with the timestamp value, in a YYYY-MM-DD format
	It("ToDate - Works", func() {
		stamp := generateTimestamp(1654127993, 983651350)
		Expect(stamp.ToDate()).Should(Equal("2022-06-01"))
	})
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
		uDur := generateDuration(2678400, 15000000)
		dur, err := uDur.AsDuration()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(dur).Should(Equal(31*24*time.Hour + 15*time.Millisecond))
	})

	// Tests that the Equals function works under various data conditions
	DescribeTable("Equals - Conditions",
		func(rhs *UnixDuration, lhs *UnixDuration, equal bool) {
			Expect(rhs.Equals(lhs)).Should(Equal(equal))
		},
		Entry("RHS is nil - False", nil, generateDuration(1655510000, 900838091), false),
		Entry("LHS is nil - False", generateDuration(1655510399, 900838091), nil, false),
		Entry("Both nil - True", nil, nil, true),
		Entry("RHS.Seconds != LHS.Seconds - False",
			generateDuration(1655510399, 900838091), generateDuration(1655510000, 900838091), false),
		Entry("RHS.Nanoseconds != LHS.Nanoseconds - False",
			generateDuration(1655510399, 900838091), generateDuration(1655510399, 0), false),
		Entry("RHS == LHS - True",
			generateDuration(1655510399, 900838091), generateDuration(1655510399, 900838091), true))

	// Tests that the NotEquals function works under various data conditions
	DescribeTable("NotEquals - Conditions",
		func(rhs *UnixDuration, lhs *UnixDuration, notEqual bool) {
			Expect(rhs.NotEquals(lhs)).Should(Equal(notEqual))
		},
		Entry("RHS is nil - True", nil, generateDuration(1655510000, 900838091), true),
		Entry("LHS is nil - True", generateDuration(1655510399, 900838091), nil, true),
		Entry("Both nil - False", nil, nil, false),
		Entry("RHS == LHS - False",
			generateDuration(1655510399, 900838091), generateDuration(1655510399, 900838091), false),
		Entry("RHS.Nanoseconds != LHS.Nanoseconds - True",
			generateDuration(1655510399, 900838091), generateDuration(1655510399, 0), true),
		Entry("RHS.Seconds != LHS.Seconds - True",
			generateDuration(1655510399, 900838091), generateDuration(1655510000, 900838091), true))

	// Tests that the GreaterThan function works under various data conditions
	DescribeTable("GreaterThan - Conditions",
		func(rhs *UnixDuration, lhs *UnixDuration, greaterThan bool) {
			Expect(rhs.GreaterThan(lhs)).Should(Equal(greaterThan))
		},
		Entry("RHS is nil - False", nil, generateDuration(1655510000, 900838091), false),
		Entry("LHS is nil - True", generateDuration(1655510399, 900838091), nil, true),
		Entry("Both nil - False", nil, nil, false),
		Entry("RHS == LHS - False",
			generateDuration(1655510399, 900838091), generateDuration(1655510399, 900838091), false),
		Entry("RHS.Nanoseconds < LHS.Nanoseconds - False",
			generateDuration(1655510399, 0), generateDuration(1655510399, 900838091), false),
		Entry("RHS.Nanoseconds > LHS.Nanoseconds - True",
			generateDuration(1655510399, 900838091), generateDuration(1655510399, 0), true),
		Entry("RHS.Seconds < LHS.Seconds - False",
			generateDuration(1655510000, 900838091), generateDuration(1655510399, 900838091), false),
		Entry("RHS.Seconds > LHS.Seconds - True",
			generateDuration(1655510399, 900838091), generateDuration(1655510000, 900838091), true))

	// Tests that the GreaterThanOrEqualTo function works under various data conditions
	DescribeTable("GreaterThanOrEqualTo - Conditions",
		func(rhs *UnixDuration, lhs *UnixDuration, gte bool) {
			Expect(rhs.GreaterThanOrEqualTo(lhs)).Should(Equal(gte))
		},
		Entry("RHS is nil - False", nil, generateDuration(1655510000, 900838091), false),
		Entry("LHS is nil - True", generateDuration(1655510399, 900838091), nil, true),
		Entry("Both nil - True", nil, nil, true),
		Entry("RHS == LHS - True",
			generateDuration(1655510399, 900838091), generateDuration(1655510399, 900838091), true),
		Entry("RHS.Nanoseconds < LHS.Nanoseconds - False",
			generateDuration(1655510399, 0), generateDuration(1655510399, 900838091), false),
		Entry("RHS.Nanoseconds > LHS.Nanoseconds - True",
			generateDuration(1655510399, 900838091), generateDuration(1655510399, 0), true),
		Entry("RHS.Seconds < LHS.Seconds - False",
			generateDuration(1655510000, 900838091), generateDuration(1655510399, 900838091), false),
		Entry("RHS.Seconds > LHS.Seconds - True",
			generateDuration(1655510399, 900838091), generateDuration(1655510000, 900838091), true))

	// Tests that the LessThan function works under various data conditions
	DescribeTable("LessThan - Conditions",
		func(rhs *UnixDuration, lhs *UnixDuration, lt bool) {
			Expect(rhs.LessThan(lhs)).Should(Equal(lt))
		},
		Entry("RHS is nil - True", nil, generateDuration(1655510000, 900838091), true),
		Entry("LHS is nil - False", generateDuration(1655510399, 900838091), nil, false),
		Entry("Both nil - False", nil, nil, false),
		Entry("RHS == LHS - False",
			generateDuration(1655510399, 900838091), generateDuration(1655510399, 900838091), false),
		Entry("RHS.Nanoseconds < LHS.Nanoseconds - True",
			generateDuration(1655510399, 0), generateDuration(1655510399, 900838091), true),
		Entry("RHS.Nanoseconds > LHS.Nanoseconds - False",
			generateDuration(1655510399, 900838091), generateDuration(1655510399, 0), false),
		Entry("RHS.Seconds < LHS.Seconds - True",
			generateDuration(1655510000, 900838091), generateDuration(1655510399, 900838091), true),
		Entry("RHS.Seconds > LHS.Seconds - False",
			generateDuration(1655510399, 900838091), generateDuration(1655510000, 900838091), false))

	// Tests that the LessThanOrEqualTo function works under various data conditions
	DescribeTable("LessThanOrEqualTo - Condition",
		func(rhs *UnixDuration, lhs *UnixDuration, lte bool) {
			Expect(rhs.LessThanOrEqualTo(lhs)).Should(Equal(lte))
		},
		Entry("RHS is nil - True", nil, generateDuration(1655510000, 900838091), true),
		Entry("LHS is nil - False", generateDuration(1655510399, 900838091), nil, false),
		Entry("Both nil - True", nil, nil, true),
		Entry("RHS == LHS - True",
			generateDuration(1655510399, 900838091), generateDuration(1655510399, 900838091), true),
		Entry("RHS.Nanoseconds < LHS.Nanoseconds - True",
			generateDuration(1655510399, 0), generateDuration(1655510399, 900838091), true),
		Entry("RHS.Nanoseconds > LHS.Nanoseconds - False",
			generateDuration(1655510399, 900838091), generateDuration(1655510399, 0), false),
		Entry("RHS.Seconds < LHS.Seconds - True",
			generateDuration(1655510000, 900838091), generateDuration(1655510399, 900838091), true),
		Entry("RHS.Seconds > LHS.Seconds - False",
			generateDuration(1655510399, 900838091), generateDuration(1655510000, 900838091), false))

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
