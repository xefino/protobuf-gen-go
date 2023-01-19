package gopb

import (
	"math"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
)

var _ = Describe("Decimal Extensions Tests", func() {

	// Tests that the NewFromDecimal function works under various data conditions
	DescribeTable("NewFromDecimal - Works",
		func(raw string, verifer func(*Decimal)) {

			// First, create our decimal value
			dIn, err := decimal.NewFromString(raw)
			Expect(err).ShouldNot(HaveOccurred())

			// Next, attempt to convert it to a new Decimal object
			dOut := NewFromDecimal(dIn)

			// Finally, verify the data
			verifer(dOut)
		},
		Entry("Value greater than 0 - Encoded", "1234512351234088800000.999",
			decimalVerifier(-3, 351234088800000999, 1234512)),
		Entry("Value equal to 0 - Encoded", "0", decimalVerifier(0)),
		Entry("Value less than 0 - Encoded", "-288341660781234512351234088800000.999",
			decimalVerifier(-3, -351234088800000999, -288341660781234512)))

	// Tests that the ToDecimal function works under various data conditions
	DescribeTable("ToDecimal - Works",
		func(dIn *Decimal, expected string) {
			Expect(dIn.ToDecimal().String()).Should(Equal(expected))
		},
		Entry("Value greater than 0 - Encoded",
			&Decimal{Parts: []int64{351234088800000999, 1234512}, Exp: -3}, "1234512351234088800000.999"),
		Entry("Value equals 0 - Encoded", &Decimal{Parts: make([]int64, 0)}, "0"),
		Entry("Value less than 0 - Encoded",
			&Decimal{Parts: []int64{-351234088800000999, -342645987}, Exp: -5}, "-3426459873512340888000.00999"))
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
		timestamp := NewFromTime(time.Date(2022, time.June, 1, 23, 59, 53, 983651350, time.UTC))

		// Verify that the number of seconds and nanoseconds is correct
		Expect(timestamp.Seconds).Should(Equal(int64(1654127993)))
		Expect(timestamp.Nanoseconds).Should(Equal(int32(983651350)))
	})

	// Tests that, if any of the values sent to MaxTimestamp function is nil then the function will panic
	It("MaxTimestamp - Item is nil - Panic", func() {
		Expect(func() {
			_ = MaxTimestamp(NewUnixTimestamp(1655510399, 36000000), NewUnixTimestamp(1655510399, 0), nil)
		}).Should(Panic())
	})

	// Tests the data conditions determining what MaxTimestamp will return
	DescribeTable("MaxTimestamp - Conditions",
		func(expected *UnixTimestamp, values ...*UnixTimestamp) {
			max := MaxTimestamp(values[0], values[1], values[2:]...)
			Expect(max).Should(Equal(expected))
		},
		Entry("Seconds different - Largest returned", NewUnixTimestamp(1655510399, 0),
			NewUnixTimestamp(1655510000, 0), NewUnixTimestamp(1655510399, 0), NewUnixTimestamp(1555510399, 0)),
		Entry("Seconds same, Nanoseconds different - Largest returned", NewUnixTimestamp(1655510399, 900838091),
			NewUnixTimestamp(1655510399, 0), NewUnixTimestamp(1655510399, 36000000), NewUnixTimestamp(1655510399, 900838091)),
		Entry("Seconds same, Nanoseconds same - First returned", NewUnixTimestamp(1655510399, 900838091),
			NewUnixTimestamp(1655510399, 900838091), NewUnixTimestamp(1655510399, 900838091), NewUnixTimestamp(1655510399, 900838091)))

	// Tests that, if any of the values sent to the MinTimestamp function is nil then the function will panic
	It("MinTimestamp - Item is nil - Panic", func() {
		Expect(func() {
			_ = MinTimestamp(NewUnixTimestamp(1655510399, 36000000), NewUnixTimestamp(1655510399, 0), nil)
		}).Should(Panic())
	})

	// Tests the data conditions determining what MinTimestamp will return
	DescribeTable("MinTimestamp - Conditions",
		func(expected *UnixTimestamp, values ...*UnixTimestamp) {
			min := MinTimestamp(values[0], values[1], values[2:]...)
			Expect(min).Should(Equal(expected))
		},
		Entry("Seconds different - Smallest returned", NewUnixTimestamp(1555510399, 0),
			NewUnixTimestamp(1655510000, 0), NewUnixTimestamp(1655510399, 0), NewUnixTimestamp(1555510399, 0)),
		Entry("Seconds same, Nanoseconds different - Smallest returned", NewUnixTimestamp(1655510399, 0),
			NewUnixTimestamp(1655510399, 36000000), NewUnixTimestamp(1655510399, 0), NewUnixTimestamp(1655510399, 900838091)),
		Entry("Seconds same, Nanoseconds same - First returned", NewUnixTimestamp(1655510399, 900838091),
			NewUnixTimestamp(1655510399, 900838091), NewUnixTimestamp(1655510399, 900838091), NewUnixTimestamp(1655510399, 900838091)))

	// Test that the AsTime function creates a time from a valid timestamp
	It("AsTime - Works", func() {

		// First, create a timestamp with a set number of seconds and nanoseconds
		timestamp := NewUnixTimestamp(1654127993, 983651350)

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
		Entry("RHS is nil - False", nil, NewUnixTimestamp(1655510000, 900838091), false),
		Entry("LHS is nil - False", NewUnixTimestamp(1655510399, 900838091), nil, false),
		Entry("Both nil - True", nil, nil, true),
		Entry("RHS.Seconds != LHS.Seconds - False",
			NewUnixTimestamp(1655510399, 900838091), NewUnixTimestamp(1655510000, 900838091), false),
		Entry("RHS.Nanoseconds != LHS.Nanoseconds - False",
			NewUnixTimestamp(1655510399, 900838091), NewUnixTimestamp(1655510399, 0), false),
		Entry("RHS == LHS - True",
			NewUnixTimestamp(1655510399, 900838091), NewUnixTimestamp(1655510399, 900838091), true))

	// Tests that the NotEquals function works under various data conditions
	DescribeTable("NotEquals - Conditions",
		func(rhs *UnixTimestamp, lhs *UnixTimestamp, notEqual bool) {
			Expect(rhs.NotEquals(lhs)).Should(Equal(notEqual))
		},
		Entry("RHS is nil - True", nil, NewUnixTimestamp(1655510000, 900838091), true),
		Entry("LHS is nil - True", NewUnixTimestamp(1655510399, 900838091), nil, true),
		Entry("Both nil - False", nil, nil, false),
		Entry("RHS == LHS - False",
			NewUnixTimestamp(1655510399, 900838091), NewUnixTimestamp(1655510399, 900838091), false),
		Entry("RHS.Nanoseconds != LHS.Nanoseconds - True",
			NewUnixTimestamp(1655510399, 900838091), NewUnixTimestamp(1655510399, 0), true),
		Entry("RHS.Seconds != LHS.Seconds - True",
			NewUnixTimestamp(1655510399, 900838091), NewUnixTimestamp(1655510000, 900838091), true))

	// Tests that the GreaterThan function works under various data conditions
	DescribeTable("GreaterThan - Conditions",
		func(rhs *UnixTimestamp, lhs *UnixTimestamp, greaterThan bool) {
			Expect(rhs.GreaterThan(lhs)).Should(Equal(greaterThan))
		},
		Entry("RHS is nil - False", nil, NewUnixTimestamp(1655510000, 900838091), false),
		Entry("LHS is nil - True", NewUnixTimestamp(1655510399, 900838091), nil, true),
		Entry("Both nil - False", nil, nil, false),
		Entry("RHS == LHS - False",
			NewUnixTimestamp(1655510399, 900838091), NewUnixTimestamp(1655510399, 900838091), false),
		Entry("RHS.Nanoseconds < LHS.Nanoseconds - False",
			NewUnixTimestamp(1655510399, 0), NewUnixTimestamp(1655510399, 900838091), false),
		Entry("RHS.Nanoseconds > LHS.Nanoseconds - True",
			NewUnixTimestamp(1655510399, 900838091), NewUnixTimestamp(1655510399, 0), true),
		Entry("RHS.Seconds < LHS.Seconds - False",
			NewUnixTimestamp(1655510000, 900838091), NewUnixTimestamp(1655510399, 900838091), false),
		Entry("RHS.Seconds > LHS.Seconds - True",
			NewUnixTimestamp(1655510399, 900838091), NewUnixTimestamp(1655510000, 900838091), true))

	// Tests that the GreaterThanOrEqualTo function works under various data conditions
	DescribeTable("GreaterThanOrEqualTo - Conditions",
		func(rhs *UnixTimestamp, lhs *UnixTimestamp, gte bool) {
			Expect(rhs.GreaterThanOrEqualTo(lhs)).Should(Equal(gte))
		},
		Entry("RHS is nil - False", nil, NewUnixTimestamp(1655510000, 900838091), false),
		Entry("LHS is nil - True", NewUnixTimestamp(1655510399, 900838091), nil, true),
		Entry("Both nil - True", nil, nil, true),
		Entry("RHS == LHS - True",
			NewUnixTimestamp(1655510399, 900838091), NewUnixTimestamp(1655510399, 900838091), true),
		Entry("RHS.Nanoseconds < LHS.Nanoseconds - False",
			NewUnixTimestamp(1655510399, 0), NewUnixTimestamp(1655510399, 900838091), false),
		Entry("RHS.Nanoseconds > LHS.Nanoseconds - True",
			NewUnixTimestamp(1655510399, 900838091), NewUnixTimestamp(1655510399, 0), true),
		Entry("RHS.Seconds < LHS.Seconds - False",
			NewUnixTimestamp(1655510000, 900838091), NewUnixTimestamp(1655510399, 900838091), false),
		Entry("RHS.Seconds > LHS.Seconds - True",
			NewUnixTimestamp(1655510399, 900838091), NewUnixTimestamp(1655510000, 900838091), true))

	// Tests that the LessThan function works under various data conditions
	DescribeTable("LessThan - Conditions",
		func(rhs *UnixTimestamp, lhs *UnixTimestamp, lt bool) {
			Expect(rhs.LessThan(lhs)).Should(Equal(lt))
		},
		Entry("RHS is nil - True", nil, NewUnixTimestamp(1655510000, 900838091), true),
		Entry("LHS is nil - False", NewUnixTimestamp(1655510399, 900838091), nil, false),
		Entry("Both nil - False", nil, nil, false),
		Entry("RHS == LHS - False",
			NewUnixTimestamp(1655510399, 900838091), NewUnixTimestamp(1655510399, 900838091), false),
		Entry("RHS.Nanoseconds < LHS.Nanoseconds - True",
			NewUnixTimestamp(1655510399, 0), NewUnixTimestamp(1655510399, 900838091), true),
		Entry("RHS.Nanoseconds > LHS.Nanoseconds - False",
			NewUnixTimestamp(1655510399, 900838091), NewUnixTimestamp(1655510399, 0), false),
		Entry("RHS.Seconds < LHS.Seconds - True",
			NewUnixTimestamp(1655510000, 900838091), NewUnixTimestamp(1655510399, 900838091), true),
		Entry("RHS.Seconds > LHS.Seconds - False",
			NewUnixTimestamp(1655510399, 900838091), NewUnixTimestamp(1655510000, 900838091), false))

	// Tests that the LessThanOrEqualTo function works under various data conditions
	DescribeTable("LessThanOrEqualTo - Condition",
		func(rhs *UnixTimestamp, lhs *UnixTimestamp, lte bool) {
			Expect(rhs.LessThanOrEqualTo(lhs)).Should(Equal(lte))
		},
		Entry("RHS is nil - True", nil, NewUnixTimestamp(1655510000, 900838091), true),
		Entry("LHS is nil - False", NewUnixTimestamp(1655510399, 900838091), nil, false),
		Entry("Both nil - True", nil, nil, true),
		Entry("RHS == LHS - True",
			NewUnixTimestamp(1655510399, 900838091), NewUnixTimestamp(1655510399, 900838091), true),
		Entry("RHS.Nanoseconds < LHS.Nanoseconds - True",
			NewUnixTimestamp(1655510399, 0), NewUnixTimestamp(1655510399, 900838091), true),
		Entry("RHS.Nanoseconds > LHS.Nanoseconds - False",
			NewUnixTimestamp(1655510399, 900838091), NewUnixTimestamp(1655510399, 0), false),
		Entry("RHS.Seconds < LHS.Seconds - True",
			NewUnixTimestamp(1655510000, 900838091), NewUnixTimestamp(1655510399, 900838091), true),
		Entry("RHS.Seconds > LHS.Seconds - False",
			NewUnixTimestamp(1655510399, 900838091), NewUnixTimestamp(1655510000, 900838091), false))

	// Tests that the Add function works under various conditions
	DescribeTable("Add - Works",
		func(rhs *UnixTimestamp, lhs *UnixTimestamp, expected *UnixTimestamp) {
			Expect(rhs.Add(lhs)).Should(Equal(expected))
		},
		Entry("LHS is nil - Works", NewUnixTimestamp(1655510000, 900838091),
			nil, NewUnixTimestamp(1655510000, 900838091)),
		Entry("Nanoseconds < 1 second - Works", NewUnixTimestamp(1655510000, 900838091),
			NewUnixTimestamp(100, 1000), NewUnixTimestamp(1655510100, 900839091)),
		Entry("Nanoseconds > 1 second - Works", NewUnixTimestamp(1655510000, 900838091),
			NewUnixTimestamp(1655510000, 900838091), NewUnixTimestamp(3311020001, 801676182)),
		Entry("Nanoseconds < 0 - Works", NewUnixTimestamp(1655510000, 900838091),
			NewUnixTimestamp(1655510000, -999999999), NewUnixTimestamp(3311019999, 900838092)))

	// Test the conditions describing how the AddDate function works
	DescribeTable("AddDate - Works",
		func(years int, months int, days int, expected *UnixTimestamp) {
			Expect(NewUnixTimestamp(1655510000, 900838091).AddDate(years, months, days)).Should(Equal(expected))
		},
		Entry("Time is positive - Works", 1, 1, 10, NewUnixTimestamp(1690502000, 900838091)),
		Entry("Time is negative - Works", -1, -1, -10, NewUnixTimestamp(1620431600, 900838091)))

	// Tests that the AddDuration function works under various conditions
	DescribeTable("AddDuration - Conditions",
		func(duration *UnixDuration, result *UnixTimestamp) {
			timestamp := NewUnixTimestamp(1655510000, 900838091)
			Expect(timestamp.AddDuration(duration)).Should(Equal(result))
		},
		Entry("LHS is nil - Works", nil, NewUnixTimestamp(1655510000, 900838091)),
		Entry("Nanoseconds < 1 second - Works", NewUnixDuration(100, 1000),
			NewUnixTimestamp(1655510100, 900839091)),
		Entry("Nanoseconds > 1 second - Works", NewUnixDuration(1655510000, 900838091),
			NewUnixTimestamp(3311020001, 801676182)),
		Entry("Nanoseconds < 0 - Works", NewUnixDuration(1655510000, -999999999),
			NewUnixTimestamp(3311019999, 900838092)))

	// Tests that the Difference functions works under various conditions
	DescribeTable("Difference - Conditions",
		func(rhs *UnixTimestamp, lhs *UnixTimestamp, result *UnixDuration) {
			Expect(rhs.Difference(lhs)).Should(Equal(result))
		},
		Entry("rhs.Seconds > lhs.Seconds, rhs.Nanoseconds > lhs.Nanoseconds - Works",
			NewUnixTimestamp(1669704178, 500000000), NewUnixTimestamp(1669704177, 0), NewUnixDuration(1, 500000000)),
		Entry("rhs.Seconds > lhs.Seconds, rhs.Nanoseconds < lhs.Nanoseconds - Works",
			NewUnixTimestamp(1669704178, 0), NewUnixTimestamp(1669704177, 500000000), NewUnixDuration(0, 500000000)),
		Entry("rhs.Seconds < lhs.Seconds, rhs.Nanoseconds > lhs.Nanoseconds - Works",
			NewUnixTimestamp(1669704177, 500000000), NewUnixTimestamp(1669704178, 0), NewUnixDuration(0, -500000000)),
		Entry("rhs.Seconds < lhs.Seconds, rhs.Nanoseconds < lhs.Nanoseconds - Works",
			NewUnixTimestamp(1669704177, 0), NewUnixTimestamp(1669704178, 500000000), NewUnixDuration(-1, -500000000)))

	// Tests that the NextDay function works as expected
	It("NextDay - Works", func() {
		next := NewUnixTimestamp(1655510000, 900838091).NextDay()
		Expect(next).Should(Equal(NewUnixTimestamp(1655510400, 0)))
	})

	// Tests that the SecondDown function works as expected
	DescribeTable("SecondDown - Conditions",
		func(start *UnixTimestamp, result *UnixTimestamp) {
			Expect(start.SecondDown()).Should(Equal(result))
		},
		Entry("Nanoseconds > 0 - Erased",
			NewUnixTimestamp(1655510000, 900838091), NewUnixTimestamp(1655510000, 0)),
		Entry("Nanoseconds = 0 - Copied",
			NewUnixTimestamp(1655510000, 0), NewUnixTimestamp(1655510000, 0)))

	// Tests that the SecondUp function works as expected
	DescribeTable("SecondUp - Conditions",
		func(start *UnixTimestamp, result *UnixTimestamp) {
			Expect(start.SecondUp()).Should(Equal(result))
		},
		Entry("Nanoseconds > 0 - Snapped",
			NewUnixTimestamp(1655510000, 900838091), NewUnixTimestamp(1655510001, 0)),
		Entry("Nanoseconds = 0 - Copied",
			NewUnixTimestamp(1655510000, 0), NewUnixTimestamp(1655510000, 0)))

	// Tests that the MinuteDown function works as expected
	DescribeTable("MinuteDown - Conditions",
		func(start *UnixTimestamp, result *UnixTimestamp) {
			Expect(start.MinuteDown()).Should(Equal(result))
		},
		Entry("Nanoseconds > 0 - Erased",
			NewUnixTimestamp(1655510000, 900838091), NewUnixTimestamp(1655509980, 0)),
		Entry("Nanoseconds = 0, Seconds > 0, Erased",
			NewUnixTimestamp(1655510000, 0), NewUnixTimestamp(1655509980, 0)),
		Entry("Nanoseconds = 0, Seconds = 0 - Copied",
			NewUnixTimestamp(1655509980, 0), NewUnixTimestamp(1655509980, 0)))

	// Tests that the MinuteUp function works as expected
	DescribeTable("MinuteUp - Conditions",
		func(start *UnixTimestamp, result *UnixTimestamp) {
			Expect(start.MinuteUp()).Should(Equal(result))
		},
		Entry("Nanoseconds > 0 - Snapped",
			NewUnixTimestamp(1655510000, 900838091), NewUnixTimestamp(1655510040, 0)),
		Entry("Nanoseconds = 0, Seconds > 0 - Snapped",
			NewUnixTimestamp(1655510000, 0), NewUnixTimestamp(1655510040, 0)),
		Entry("Nanoseconds = 0, Seconds = 0 - Copied",
			NewUnixTimestamp(1655509980, 0), NewUnixTimestamp(1655509980, 0)))

	// Tests that the HourDown function works as expected
	DescribeTable("HourDown - Conditions",
		func(start *UnixTimestamp, result *UnixTimestamp) {
			Expect(start.HourDown()).Should(Equal(result))
		},
		Entry("Nanoseconds > 0 - Erased",
			NewUnixTimestamp(1655510000, 900838091), NewUnixTimestamp(1655506800, 0)),
		Entry("Nanoseconds = 0, Seconds > 0 - Erased",
			NewUnixTimestamp(1655510000, 0), NewUnixTimestamp(1655506800, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes > 0 - Erased",
			NewUnixTimestamp(1655509980, 0), NewUnixTimestamp(1655506800, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes = 0 - Copied",
			NewUnixTimestamp(1655506800, 0), NewUnixTimestamp(1655506800, 0)))

	// Tests that the HourUp function works as expected
	DescribeTable("HourUp - Conditions",
		func(start *UnixTimestamp, result *UnixTimestamp) {
			Expect(start.HourUp()).Should(Equal(result))
		},
		Entry("Nanoseconds > 0 - Snapped",
			NewUnixTimestamp(1655510000, 900838091), NewUnixTimestamp(1655510400, 0)),
		Entry("Nanoseconds = 0, Seconds > 0 - Snapped",
			NewUnixTimestamp(1655510000, 0), NewUnixTimestamp(1655510400, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes > 0 - Snapped",
			NewUnixTimestamp(1655509980, 0), NewUnixTimestamp(1655510400, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes = 0 - Copied",
			NewUnixTimestamp(1655424000, 0), NewUnixTimestamp(1655424000, 0)))

	// Tests that the DayDown function works as expected
	DescribeTable("DayDown - Conditions",
		func(start *UnixTimestamp, result *UnixTimestamp) {
			Expect(start.DayDown()).Should(Equal(result))
		},
		Entry("Nanoseconds > 0 - Erased",
			NewUnixTimestamp(1655510000, 900838091), NewUnixTimestamp(1655424000, 0)),
		Entry("Nanoseconds = 0, Seconds > 0 - Erased",
			NewUnixTimestamp(1655510000, 0), NewUnixTimestamp(1655424000, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes > 0 - Erased",
			NewUnixTimestamp(1655509980, 0), NewUnixTimestamp(1655424000, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Erased",
			NewUnixTimestamp(1655506800, 0), NewUnixTimestamp(1655424000, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0 - Copied",
			NewUnixTimestamp(1655424000, 0), NewUnixTimestamp(1655424000, 0)))

	// Tests that the DayUp function works as expected
	DescribeTable("DayUp - Conditions",
		func(start *UnixTimestamp, result *UnixTimestamp) {
			Expect(start.DayUp()).Should(Equal(result))
		},
		Entry("Nanoseconds > 0 - Snapped",
			NewUnixTimestamp(1655510000, 900838091), NewUnixTimestamp(1655510400, 0)),
		Entry("Nanoseconds = 0, Seconds > 0 - Snapped",
			NewUnixTimestamp(1655510000, 0), NewUnixTimestamp(1655510400, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes > 0 - Snapped",
			NewUnixTimestamp(1655509980, 0), NewUnixTimestamp(1655510400, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Snapped",
			NewUnixTimestamp(1655506800, 0), NewUnixTimestamp(1655510400, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0 - Copied",
			NewUnixTimestamp(1655510400, 0), NewUnixTimestamp(1655510400, 0)))

	// Tests that the WeekDown function works as expected
	DescribeTable("WeekDown - Conditions",
		func(start *UnixTimestamp, result *UnixTimestamp) {
			Expect(start.WeekDown()).Should(Equal(result))
		},
		Entry("Monday, Nanoseconds > 0 - Erased",
			NewUnixTimestamp(1655742769, 900838091), NewUnixTimestamp(1655596800, 0)),
		Entry("Monday, Nanoseconds = 0, Seconds > 0 - Erased",
			NewUnixTimestamp(1655742769, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Monday, Nanoseconds = 0, Seconds = 0, Minutes > 0 - Erased",
			NewUnixTimestamp(1655742720, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Monday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Erased",
			NewUnixTimestamp(1655740800, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Monday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0 - Erased",
			NewUnixTimestamp(1655683200, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Tuesday, Nanoseconds > 0 - Erased",
			NewUnixTimestamp(1655829169, 900838091), NewUnixTimestamp(1655596800, 0)),
		Entry("Tuesday, Nanoseconds = 0, Seconds > 0 - Erased",
			NewUnixTimestamp(1655829169, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Tuesday, Nanoseconds = 0, Seconds = 0, Minutes > 0 - Erased",
			NewUnixTimestamp(1655829120, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Tuesday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Erased",
			NewUnixTimestamp(1655827200, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Tuesday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0 - Erased",
			NewUnixTimestamp(1655769600, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Wednesday, Nanoseconds > 0 - Erased",
			NewUnixTimestamp(1655915569, 900838091), NewUnixTimestamp(1655596800, 0)),
		Entry("Wednesday, Nanoseconds = 0, Seconds > 0 - Erased",
			NewUnixTimestamp(1655915569, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Wednesday, Nanoseconds = 0, Seconds = 0, Minutes > 0 - Erased",
			NewUnixTimestamp(1655915520, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Wednesday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Erased",
			NewUnixTimestamp(1655913600, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Wednesday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0 - Erased",
			NewUnixTimestamp(1655856000, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Thursday, Nanoseconds > 0 - Erased",
			NewUnixTimestamp(1656001969, 900838091), NewUnixTimestamp(1655596800, 0)),
		Entry("Thursday, Nanoseconds = 0, Seconds > 0 - Erased",
			NewUnixTimestamp(1656001969, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Thursday, Nanoseconds = 0, Seconds = 0, Minutes > 0 - Erased",
			NewUnixTimestamp(1656001920, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Thursday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Erased",
			NewUnixTimestamp(1656000000, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Thursday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0 - Erased",
			NewUnixTimestamp(1655942400, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Friday, Nanoseconds > 0 - Erased",
			NewUnixTimestamp(1656088369, 900838091), NewUnixTimestamp(1655596800, 0)),
		Entry("Friday, Nanoseconds = 0, Seconds > 0 - Erased",
			NewUnixTimestamp(1656088369, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Friday, Nanoseconds = 0, Seconds = 0, Minutes > 0 - Erased",
			NewUnixTimestamp(1656088320, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Friday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Erased",
			NewUnixTimestamp(1656086400, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Friday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0 - Erased",
			NewUnixTimestamp(1656028800, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Saturday, Nanoseconds > 0 - Erased",
			NewUnixTimestamp(1656174769, 900838091), NewUnixTimestamp(1655596800, 0)),
		Entry("Saturday, Nanoseconds = 0, Seconds > 0 - Erased",
			NewUnixTimestamp(1656174769, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Saturday, Nanoseconds = 0, Seconds = 0, Minutes > 0 - Erased",
			NewUnixTimestamp(1656174720, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Saturday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Erased",
			NewUnixTimestamp(1656172800, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Saturday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0 - Erased",
			NewUnixTimestamp(1656115200, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Sunday, Nanoseconds > 0 - Erased",
			NewUnixTimestamp(1655656369, 900838091), NewUnixTimestamp(1655596800, 0)),
		Entry("Sunday, Nanoseconds = 0, Seconds > 0 - Erased",
			NewUnixTimestamp(1655656369, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Sunday, Nanoseconds = 0, Seconds = 0, Minutes > 0 - Erased",
			NewUnixTimestamp(1655656320, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Sunday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Erased",
			NewUnixTimestamp(1655654400, 0), NewUnixTimestamp(1655596800, 0)),
		Entry("Sunday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0 - Copied",
			NewUnixTimestamp(1655596800, 0), NewUnixTimestamp(1655596800, 0)))

	// Tests that the WeekUp function works as expected
	DescribeTable("WeekUp - Conditions",
		func(start *UnixTimestamp, result *UnixTimestamp) {
			Expect(start.WeekUp()).Should(Equal(result))
		},
		Entry("Monday, Nanoseconds > 0 - Snapped",
			NewUnixTimestamp(1655742769, 900838091), NewUnixTimestamp(1656201600, 0)),
		Entry("Monday, Nanoseconds = 0, Seconds > 0 - Snapped",
			NewUnixTimestamp(1655742769, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Monday, Nanoseconds = 0, Seconds = 0, Minutes > 0 - Snapped",
			NewUnixTimestamp(1655742720, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Monday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Snapped",
			NewUnixTimestamp(1655740800, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Monday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0 - Snapped",
			NewUnixTimestamp(1655683200, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Tuesday, Nanoseconds > 0 - Snapped",
			NewUnixTimestamp(1655829169, 900838091), NewUnixTimestamp(1656201600, 0)),
		Entry("Tuesday, Nanoseconds = 0, Seconds > 0 - Snapped",
			NewUnixTimestamp(1655829169, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Tuesday, Nanoseconds = 0, Seconds = 0, Minutes > 0 - Snapped",
			NewUnixTimestamp(1655829120, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Tuesday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Snapped",
			NewUnixTimestamp(1655827200, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Tuesday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0 - Snapped",
			NewUnixTimestamp(1655769600, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Wednesday, Nanoseconds > 0 - Snapped",
			NewUnixTimestamp(1655915569, 900838091), NewUnixTimestamp(1656201600, 0)),
		Entry("Wednesday, Nanoseconds = 0, Seconds > 0 - Snapped",
			NewUnixTimestamp(1655915569, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Wednesday, Nanoseconds = 0, Seconds = 0, Minutes > 0 - Snapped",
			NewUnixTimestamp(1655915520, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Wednesday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Snapped",
			NewUnixTimestamp(1655913600, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Wednesday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0 - Snapped",
			NewUnixTimestamp(1655856000, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Thursday, Nanoseconds > 0 - Snapped",
			NewUnixTimestamp(1656001969, 900838091), NewUnixTimestamp(1656201600, 0)),
		Entry("Thursday, Nanoseconds = 0, Seconds > 0 - Snapped",
			NewUnixTimestamp(1656001969, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Thursday, Nanoseconds = 0, Seconds = 0, Minutes > 0 - Snapped",
			NewUnixTimestamp(1656001920, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Thursday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Snapped",
			NewUnixTimestamp(1656000000, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Thursday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0 - Snapped",
			NewUnixTimestamp(1655942400, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Friday, Nanoseconds > 0 - Snapped",
			NewUnixTimestamp(1656088369, 900838091), NewUnixTimestamp(1656201600, 0)),
		Entry("Friday, Nanoseconds = 0, Seconds > 0 - Snapped",
			NewUnixTimestamp(1656088369, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Friday, Nanoseconds = 0, Seconds = 0, Minutes > 0 - Snapped",
			NewUnixTimestamp(1656088320, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Friday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Snapped",
			NewUnixTimestamp(1656086400, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Friday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0 - Snapped",
			NewUnixTimestamp(1656028800, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Saturday, Nanoseconds > 0 - Snapped",
			NewUnixTimestamp(1656174769, 900838091), NewUnixTimestamp(1656201600, 0)),
		Entry("Saturday, Nanoseconds = 0, Seconds > 0 - Snapped",
			NewUnixTimestamp(1656174769, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Saturday, Nanoseconds = 0, Seconds = 0, Minutes > 0 - Snapped",
			NewUnixTimestamp(1656174720, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Saturday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Snapped",
			NewUnixTimestamp(1656172800, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Saturday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0 - Snapped",
			NewUnixTimestamp(1656115200, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Sunday, Nanoseconds > 0 - Snapped",
			NewUnixTimestamp(1655656369, 900838091), NewUnixTimestamp(1656201600, 0)),
		Entry("Sunday, Nanoseconds = 0, Seconds > 0 - Snapped",
			NewUnixTimestamp(1655656369, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Sunday, Nanoseconds = 0, Seconds = 0, Minutes > 0 - Snapped",
			NewUnixTimestamp(1655656320, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Sunday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Snapped",
			NewUnixTimestamp(1655654400, 0), NewUnixTimestamp(1656201600, 0)),
		Entry("Sunday, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0 - Copied",
			NewUnixTimestamp(1655596800, 0), NewUnixTimestamp(1655596800, 0)))

	// Tests that the MonthDown function works as expected
	DescribeTable("MonthDown - Conditions",
		func(start *UnixTimestamp, result *UnixTimestamp) {
			Expect(start.MonthDown()).Should(Equal(result))
		},
		Entry("Nanoseconds > 0 - Erased",
			NewUnixTimestamp(1655510000, 900838091), NewUnixTimestamp(1654041600, 0)),
		Entry("Nanoseconds = 0, Seconds > 0 - Erased",
			NewUnixTimestamp(1655510000, 0), NewUnixTimestamp(1654041600, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes > 0 - Erased",
			NewUnixTimestamp(1655509980, 0), NewUnixTimestamp(1654041600, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Erased",
			NewUnixTimestamp(1655506800, 0), NewUnixTimestamp(1654041600, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days > 1 - Erased",
			NewUnixTimestamp(1655424000, 0), NewUnixTimestamp(1654041600, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days = 1 - Copied",
			NewUnixTimestamp(1654041600, 0), NewUnixTimestamp(1654041600, 0)))

	// Tests that the MonthUp function works as expected
	DescribeTable("MonthUp - Conditions",
		func(start *UnixTimestamp, result *UnixTimestamp) {
			Expect(start.MonthUp()).Should(Equal(result))
		},
		Entry("Nanoseconds > 0 - Snapped",
			NewUnixTimestamp(1655510000, 900838091), NewUnixTimestamp(1656633600, 0)),
		Entry("Nanoseconds = 0, Seconds > 0 - Snapped",
			NewUnixTimestamp(1655510000, 0), NewUnixTimestamp(1656633600, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes > 0 - Snapped",
			NewUnixTimestamp(1655509980, 0), NewUnixTimestamp(1656633600, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Snapped",
			NewUnixTimestamp(1655506800, 0), NewUnixTimestamp(1656633600, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days > 1 - Snapped",
			NewUnixTimestamp(1655424000, 0), NewUnixTimestamp(1656633600, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days = 1 - Copied",
			NewUnixTimestamp(1654041600, 0), NewUnixTimestamp(1654041600, 0)))

	// Tests that the QuarterDown function works as expected
	DescribeTable("QuarterDown - Conditions",
		func(start *UnixTimestamp, result *UnixTimestamp) {
			Expect(start.QuarterDown()).Should(Equal(result))
		},
		Entry("Q1, Nanoseconds > 0 - Erased",
			NewUnixTimestamp(1647793969, 900838091), NewUnixTimestamp(1640995200, 0)),
		Entry("Q1, Nanoseconds = 0, Seconds > 0 - Erased",
			NewUnixTimestamp(1647793969, 0), NewUnixTimestamp(1640995200, 0)),
		Entry("Q1, Nanoseconds = 0, Seconds = 0, Minutes > 0 - Erased",
			NewUnixTimestamp(1647793920, 0), NewUnixTimestamp(1640995200, 0)),
		Entry("Q1, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Erased",
			NewUnixTimestamp(1647792000, 0), NewUnixTimestamp(1640995200, 0)),
		Entry("Q1, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days > 1 - Erased",
			NewUnixTimestamp(1647734400, 0), NewUnixTimestamp(1640995200, 0)),
		Entry("Q1, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days = 1, Months > 1 - Erased",
			NewUnixTimestamp(1646092800, 0), NewUnixTimestamp(1640995200, 0)),
		Entry("Q1, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days = 1, Months = 1 - Copied",
			NewUnixTimestamp(1640995200, 0), NewUnixTimestamp(1640995200, 0)),
		Entry("Q2, Nanoseconds > 0 - Erased",
			NewUnixTimestamp(1653064369, 900838091), NewUnixTimestamp(1648771200, 0)),
		Entry("Q2, Nanoseconds = 0, Seconds > 0 - Erased",
			NewUnixTimestamp(1653064369, 0), NewUnixTimestamp(1648771200, 0)),
		Entry("Q2, Nanoseconds = 0, Seconds = 0, Minutes > 0 - Erased",
			NewUnixTimestamp(1653064320, 0), NewUnixTimestamp(1648771200, 0)),
		Entry("Q2, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Erased",
			NewUnixTimestamp(1653062400, 0), NewUnixTimestamp(1648771200, 0)),
		Entry("Q2, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days > 1 - Erased",
			NewUnixTimestamp(1653004800, 0), NewUnixTimestamp(1648771200, 0)),
		Entry("Q2, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days = 1, Months > 4 - Erased",
			NewUnixTimestamp(1651363200, 0), NewUnixTimestamp(1648771200, 0)),
		Entry("Q2, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days = 1, Months = 4 - Copied",
			NewUnixTimestamp(1648771200, 0), NewUnixTimestamp(1648771200, 0)),
		Entry("Q3, Nanoseconds > 0 - Erased",
			NewUnixTimestamp(1663691569, 900838091), NewUnixTimestamp(1656633600, 0)),
		Entry("Q3, Nanoseconds = 0, Seconds > 0 - Erased",
			NewUnixTimestamp(1663691569, 0), NewUnixTimestamp(1656633600, 0)),
		Entry("Q3, Nanoseconds = 0, Seconds = 0, Minutes > 0 - Erased",
			NewUnixTimestamp(1663691520, 0), NewUnixTimestamp(1656633600, 0)),
		Entry("Q3, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Erased",
			NewUnixTimestamp(1663689600, 0), NewUnixTimestamp(1656633600, 0)),
		Entry("Q3, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days > 1 - Erased",
			NewUnixTimestamp(1663632000, 0), NewUnixTimestamp(1656633600, 0)),
		Entry("Q3, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days = 1, Months > 7 - Erased",
			NewUnixTimestamp(1661990400, 0), NewUnixTimestamp(1656633600, 0)),
		Entry("Q3, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days = 1, Months = 7 - Copied",
			NewUnixTimestamp(1656633600, 0), NewUnixTimestamp(1656633600, 0)),
		Entry("Q4, Nanoseconds > 0 - Erased",
			NewUnixTimestamp(1668961969, 900838091), NewUnixTimestamp(1664582400, 0)),
		Entry("Q4, Nanoseconds = 0, Seconds > 0 - Erased",
			NewUnixTimestamp(1668961969, 0), NewUnixTimestamp(1664582400, 0)),
		Entry("Q4, Nanoseconds = 0, Seconds = 0, Minutes > 0 - Erased",
			NewUnixTimestamp(1668961920, 0), NewUnixTimestamp(1664582400, 0)),
		Entry("Q4, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Erased",
			NewUnixTimestamp(1668960000, 0), NewUnixTimestamp(1664582400, 0)),
		Entry("Q4, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days > 1 - Erased",
			NewUnixTimestamp(1668902400, 0), NewUnixTimestamp(1664582400, 0)),
		Entry("Q4, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days = 1, Months > 10 - Erased",
			NewUnixTimestamp(1667260800, 0), NewUnixTimestamp(1664582400, 0)),
		Entry("Q4, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days = 1, Months = 10 - Copied",
			NewUnixTimestamp(1664582400, 0), NewUnixTimestamp(1664582400, 0)))

	// Tests that the QuarterUp function works as expected
	DescribeTable("QuarterUp - Conditions",
		func(start *UnixTimestamp, result *UnixTimestamp) {
			Expect(start.QuarterUp()).Should(Equal(result))
		},
		Entry("Q1, Nanoseconds > 0 - Snapped",
			NewUnixTimestamp(1647793969, 900838091), NewUnixTimestamp(1648771200, 0)),
		Entry("Q1, Nanoseconds = 0, Seconds > 0 - Snapped",
			NewUnixTimestamp(1647793969, 0), NewUnixTimestamp(1648771200, 0)),
		Entry("Q1, Nanoseconds = 0, Seconds = 0, Minutes > 0 - Snapped",
			NewUnixTimestamp(1647793920, 0), NewUnixTimestamp(1648771200, 0)),
		Entry("Q1, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Snapped",
			NewUnixTimestamp(1647792000, 0), NewUnixTimestamp(1648771200, 0)),
		Entry("Q1, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days > 1 - Snapped",
			NewUnixTimestamp(1647734400, 0), NewUnixTimestamp(1648771200, 0)),
		Entry("Q1, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days = 1, Months > 1 - Snapped",
			NewUnixTimestamp(1646092800, 0), NewUnixTimestamp(1648771200, 0)),
		Entry("Q1, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days = 1, Months = 1 - Copied",
			NewUnixTimestamp(1640995200, 0), NewUnixTimestamp(1640995200, 0)),
		Entry("Q2, Nanoseconds > 0 - Snapped",
			NewUnixTimestamp(1653064369, 900838091), NewUnixTimestamp(1656633600, 0)),
		Entry("Q2, Nanoseconds = 0, Seconds > 0 - Snapped",
			NewUnixTimestamp(1653064369, 0), NewUnixTimestamp(1656633600, 0)),
		Entry("Q2, Nanoseconds = 0, Seconds = 0, Minutes > 0 - Snapped",
			NewUnixTimestamp(1653064320, 0), NewUnixTimestamp(1656633600, 0)),
		Entry("Q2, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Snapped",
			NewUnixTimestamp(1653062400, 0), NewUnixTimestamp(1656633600, 0)),
		Entry("Q2, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days > 1 - Snapped",
			NewUnixTimestamp(1653004800, 0), NewUnixTimestamp(1656633600, 0)),
		Entry("Q2, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days = 1, Months > 4 - Snapped",
			NewUnixTimestamp(1651363200, 0), NewUnixTimestamp(1656633600, 0)),
		Entry("Q2, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days = 1, Months = 4 - Copied",
			NewUnixTimestamp(1648771200, 0), NewUnixTimestamp(1648771200, 0)),
		Entry("Q3, Nanoseconds > 0 - Snapped",
			NewUnixTimestamp(1663691569, 900838091), NewUnixTimestamp(1664582400, 0)),
		Entry("Q3, Nanoseconds = 0, Seconds > 0 - Snapped",
			NewUnixTimestamp(1663691569, 0), NewUnixTimestamp(1664582400, 0)),
		Entry("Q3, Nanoseconds = 0, Seconds = 0, Minutes > 0 - Snapped",
			NewUnixTimestamp(1663691520, 0), NewUnixTimestamp(1664582400, 0)),
		Entry("Q3, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Snapped",
			NewUnixTimestamp(1663689600, 0), NewUnixTimestamp(1664582400, 0)),
		Entry("Q3, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days > 1 - Snapped",
			NewUnixTimestamp(1663632000, 0), NewUnixTimestamp(1664582400, 0)),
		Entry("Q3, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days = 1, Months > 7 - Snapped",
			NewUnixTimestamp(1661990400, 0), NewUnixTimestamp(1664582400, 0)),
		Entry("Q3, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days = 1, Months = 7 - Copied",
			NewUnixTimestamp(1656633600, 0), NewUnixTimestamp(1656633600, 0)),
		Entry("Q4, Nanoseconds > 0 - Snapped",
			NewUnixTimestamp(1668961969, 900838091), NewUnixTimestamp(1672531200, 0)),
		Entry("Q4, Nanoseconds = 0, Seconds > 0 - Snapped",
			NewUnixTimestamp(1668961969, 0), NewUnixTimestamp(1672531200, 0)),
		Entry("Q4, Nanoseconds = 0, Seconds = 0, Minutes > 0 - Snapped",
			NewUnixTimestamp(1668961920, 0), NewUnixTimestamp(1672531200, 0)),
		Entry("Q4, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Snapped",
			NewUnixTimestamp(1668960000, 0), NewUnixTimestamp(1672531200, 0)),
		Entry("Q4, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days > 1 - Snapped",
			NewUnixTimestamp(1668902400, 0), NewUnixTimestamp(1672531200, 0)),
		Entry("Q4, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days = 1, Months > 10 - Snapped",
			NewUnixTimestamp(1667260800, 0), NewUnixTimestamp(1672531200, 0)),
		Entry("Q4, Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days = 1, Months = 10 - Copied",
			NewUnixTimestamp(1664582400, 0), NewUnixTimestamp(1664582400, 0)))

	// Tests that the YearDown function works as expected
	DescribeTable("YearDown - Conditions",
		func(start *UnixTimestamp, result *UnixTimestamp) {
			Expect(start.YearDown()).Should(Equal(result))
		},
		Entry("Nanoseconds > 0 - Erased",
			NewUnixTimestamp(1655510000, 900838091), NewUnixTimestamp(1640995200, 0)),
		Entry("Nanoseconds = 0, Seconds > 0 - Erased",
			NewUnixTimestamp(1655510000, 0), NewUnixTimestamp(1640995200, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes > 0 - Erased",
			NewUnixTimestamp(1655509980, 0), NewUnixTimestamp(1640995200, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Erased",
			NewUnixTimestamp(1655506800, 0), NewUnixTimestamp(1640995200, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days > 1 - Erased",
			NewUnixTimestamp(1655424000, 0), NewUnixTimestamp(1640995200, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days = 1, Months > 1 - Erased",
			NewUnixTimestamp(1654041600, 0), NewUnixTimestamp(1640995200, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days = 1, Months = 1 - Copied",
			NewUnixTimestamp(1640995200, 0), NewUnixTimestamp(1640995200, 0)))

	// Tests that the YearUp function works as expected
	DescribeTable("YearUp - Conditions",
		func(start *UnixTimestamp, result *UnixTimestamp) {
			Expect(start.YearUp()).Should(Equal(result))
		},
		Entry("Nanoseconds > 0 - Snapped",
			NewUnixTimestamp(1655510000, 900838091), NewUnixTimestamp(1672531200, 0)),
		Entry("Nanoseconds = 0, Seconds > 0 - Snapped",
			NewUnixTimestamp(1655510000, 0), NewUnixTimestamp(1672531200, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes > 0 - Snapped",
			NewUnixTimestamp(1655509980, 0), NewUnixTimestamp(1672531200, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours > 0 - Snapped",
			NewUnixTimestamp(1655506800, 0), NewUnixTimestamp(1672531200, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days > 1 - Snapped",
			NewUnixTimestamp(1655424000, 0), NewUnixTimestamp(1672531200, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days = 1, Months > 1 - Snapped",
			NewUnixTimestamp(1654041600, 0), NewUnixTimestamp(1672531200, 0)),
		Entry("Nanoseconds = 0, Seconds = 0, Minutes = 0, Hours = 0, Days = 1, Months = 1 - Copied",
			NewUnixTimestamp(1640995200, 0), NewUnixTimestamp(1640995200, 0)))

	// Tests the conditions describing how the IsWhole function works
	DescribeTable("IsWhole - Conditions",
		func(rhs *UnixTimestamp, lhs time.Duration, result bool) {
			Expect(rhs.IsWhole(lhs)).Should(Equal(result))
		},
		Entry("rhs has no nanoseconds, lhs has no nanoseconds, fits - True",
			NewUnixTimestamp(1669704178, 0), 2*time.Second, true),
		Entry("rhs has no nanoseconds, lhs has no nanoseconds, not fits - False",
			NewUnixTimestamp(1669704178, 0), 3*time.Second, false),
		Entry("rhs has no nanoseconds, lhs has nanoseconds, fits - True",
			NewUnixTimestamp(1669704177, 0), time.Second+500*time.Millisecond, true),
		Entry("rhs has no nanoseconds, lhs has nanoseconds, not fits - False",
			NewUnixTimestamp(1669704177, 0), 2*time.Second+500*time.Millisecond, false),
		Entry("rhs has nanoseconds, lhs has no nanoseconds - False",
			NewUnixTimestamp(1669704178, 500000000), 2*time.Second, false),
		Entry("rhs has nanoseconds, lhs has nanoseconds, fits - True",
			NewUnixTimestamp(1669704178, 500000000), time.Second+500*time.Millisecond, true),
		Entry("rhs has nanoseconds, lhs has nanoseconds, not fits - False",
			NewUnixTimestamp(1669704178, 500000000), 2*time.Second+500*time.Millisecond, false))

	// Tests the conditions describing how the IsWholeUnix function works
	DescribeTable("IsWholeUnix - Conditions",
		func(rhs *UnixTimestamp, lhs *UnixDuration, result bool) {
			Expect(rhs.IsWholeUnix(lhs)).Should(Equal(result))
		},
		Entry("rhs has no nanoseconds, lhs has no nanoseconds, fits - True",
			NewUnixTimestamp(1669704178, 0), NewUnixDuration(2, 0), true),
		Entry("rhs has no nanoseconds, lhs has no nanoseconds, not fits - False",
			NewUnixTimestamp(1669704178, 0), NewUnixDuration(3, 0), false),
		Entry("rhs has no nanoseconds, lhs has nanoseconds, fits - True",
			NewUnixTimestamp(1669704177, 0), NewUnixDuration(1, 500000000), true),
		Entry("rhs has no nanoseconds, lhs has nanoseconds, not fits - False",
			NewUnixTimestamp(1669704177, 0), NewUnixDuration(2, 500000000), false),
		Entry("rhs has nanoseconds, lhs has no nanoseconds - False",
			NewUnixTimestamp(1669704178, 500000000), NewUnixDuration(2, 0), false),
		Entry("rhs has nanoseconds, lhs has nanoseconds, fits - True",
			NewUnixTimestamp(1669704178, 500000000), NewUnixDuration(1, 500000000), true),
		Entry("rhs has nanoseconds, lhs has nanoseconds, not fits - False",
			NewUnixTimestamp(1669704178, 500000000), NewUnixDuration(2, 500000000), false))

	// Tests the conditions determining whether IsValid will return true or false
	DescribeTable("IsValid - Conditions",
		func(timestamp *UnixTimestamp, result bool) {
			Expect(timestamp.IsValid()).Should(Equal(result))
		},
		Entry("Timestamp is nil - False", nil, false),
		Entry("Seconds < Minimum Timestamp - False", NewUnixTimestamp(-62135596801, 983651350), false),
		Entry("Seconds > Maximum Timestamp - False", NewUnixTimestamp(253402300800, 983651350), false),
		Entry("Nanoseconds > 1 second - False", NewUnixTimestamp(1654127993, 1000000000), false),
		Entry("Nanoseconds negative - False", NewUnixTimestamp(1654127993, -1), false),
		Entry("Valid - True", NewUnixTimestamp(1654127993, 983651350), true))

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
		Entry("Seconds < Minimum Timestamp - False", NewUnixTimestamp(-62135596801, 983651350), true,
			"timestamp (-62135596801, 983651350) before 0001-01-01"),
		Entry("Seconds > Maximum Timestamp - False", NewUnixTimestamp(253402300800, 983651350), true,
			"timestamp (253402300800, 983651350) after 9999-12-31"),
		Entry("Nanoseconds > 1 second - False", NewUnixTimestamp(1654127993, 1000000000), true,
			"timestamp (1654127993, 1000000000) has out-of-range nanos"),
		Entry("Nanoseconds negative - False", NewUnixTimestamp(1654127993, -1), true,
			"timestamp (1654127993, -1) has out-of-range nanos"),
		Entry("Valid - True", NewUnixTimestamp(1654127993, 983651350), false, ""))

	// Test that the ToDate function converts the UnixTimestamp to a string describing the date associated
	// with the timestamp value, in a YYYY-MM-DD format
	It("ToDate - Works", func() {
		stamp := NewUnixTimestamp(1654127993, 983651350)
		Expect(stamp.ToDate()).Should(Equal("2022-06-01"))
	})
})

var _ = Describe("UnixDuration Extensions Tests", func() {

	// Test that the NewUnixDuration function creates a valid UnixDuration from a time.Duration
	It("NewUnixDuration - Works", func() {

		// Create a Unix duration from a specific duration
		duration := NewFromDuration(31*24*time.Hour + 15*time.Millisecond)

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
		Entry("Seconds overflow - Error", NewUnixDuration(1<<60, 0), int64(0), "Seconds count was malformed"),
		Entry("Underflow error - Error", NewUnixDuration(-9223372036, -1000000000), int64(math.MinInt64), "Duration underflow"),
		Entry("Overflow error - Error", NewUnixDuration(9223372036, 1000000000), int64(math.MaxInt64), "Duration overflow"))

	// Tests that, if any of the values sent to MaxDuration function is nil then the function will panic
	It("MaxDuration - Item is nil - Panic", func() {
		Expect(func() {
			_ = MaxDuration(NewUnixDuration(1655510399, 36000000), NewUnixDuration(1655510399, 0), nil)
		}).Should(Panic())
	})

	// Tests the data conditions determining what MaxDuration will return
	DescribeTable("MaxDuration - Conditions",
		func(expected *UnixDuration, values ...*UnixDuration) {
			max := MaxDuration(values[0], values[1], values[2:]...)
			Expect(max).Should(Equal(expected))
		},
		Entry("Seconds different - Largest returned", NewUnixDuration(1655510399, 0),
			NewUnixDuration(1655510000, 0), NewUnixDuration(1655510399, 0), NewUnixDuration(1555510399, 0)),
		Entry("Seconds same, Nanoseconds different - Largest returned", NewUnixDuration(1655510399, 900838091),
			NewUnixDuration(1655510399, 0), NewUnixDuration(1655510399, 36000000), NewUnixDuration(1655510399, 900838091)),
		Entry("Seconds same, Nanoseconds same - First returned", NewUnixDuration(1655510399, 900838091),
			NewUnixDuration(1655510399, 900838091), NewUnixDuration(1655510399, 900838091), NewUnixDuration(1655510399, 900838091)))

	// Tests that, if any of the values sent to the MinDuration function is nil then the function will panic
	It("MinDuration - Item is nil - Panic", func() {
		Expect(func() {
			_ = MinDuration(NewUnixDuration(1655510399, 36000000), NewUnixDuration(1655510399, 0), nil)
		}).Should(Panic())
	})

	// Tests the data conditions determining what MinDuration will return
	DescribeTable("MinDuration - Conditions",
		func(expected *UnixDuration, values ...*UnixDuration) {
			min := MinDuration(values[0], values[1], values[2:]...)
			Expect(min).Should(Equal(expected))
		},
		Entry("Seconds different - Smallest returned", NewUnixDuration(1555510399, 0),
			NewUnixDuration(1655510000, 0), NewUnixDuration(1655510399, 0), NewUnixDuration(1555510399, 0)),
		Entry("Seconds same, Nanoseconds different - Smallest returned", NewUnixDuration(1655510399, 0),
			NewUnixDuration(1655510399, 36000000), NewUnixDuration(1655510399, 0), NewUnixDuration(1655510399, 900838091)),
		Entry("Seconds same, Nanoseconds same - First returned", NewUnixDuration(1655510399, 900838091),
			NewUnixDuration(1655510399, 900838091), NewUnixDuration(1655510399, 900838091), NewUnixDuration(1655510399, 900838091)))

	// Tests that, if no error occurs, then calling the AsDuration function will return the UnixDuration
	// as a time.Duration object
	It("AsDuration - Works", func() {
		uDur := NewUnixDuration(2678400, 15000000)
		dur, err := uDur.AsDuration()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(dur).Should(Equal(31*24*time.Hour + 15*time.Millisecond))
	})

	// Tests that the Equals function works under various data conditions
	DescribeTable("Equals - Conditions",
		func(rhs *UnixDuration, lhs *UnixDuration, equal bool) {
			Expect(rhs.Equals(lhs)).Should(Equal(equal))
		},
		Entry("RHS is nil - False", nil, NewUnixDuration(1655510000, 900838091), false),
		Entry("LHS is nil - False", NewUnixDuration(1655510399, 900838091), nil, false),
		Entry("Both nil - True", nil, nil, true),
		Entry("RHS.Seconds != LHS.Seconds - False",
			NewUnixDuration(1655510399, 900838091), NewUnixDuration(1655510000, 900838091), false),
		Entry("RHS.Nanoseconds != LHS.Nanoseconds - False",
			NewUnixDuration(1655510399, 900838091), NewUnixDuration(1655510399, 0), false),
		Entry("RHS == LHS - True",
			NewUnixDuration(1655510399, 900838091), NewUnixDuration(1655510399, 900838091), true))

	// Tests that the NotEquals function works under various data conditions
	DescribeTable("NotEquals - Conditions",
		func(rhs *UnixDuration, lhs *UnixDuration, notEqual bool) {
			Expect(rhs.NotEquals(lhs)).Should(Equal(notEqual))
		},
		Entry("RHS is nil - True", nil, NewUnixDuration(1655510000, 900838091), true),
		Entry("LHS is nil - True", NewUnixDuration(1655510399, 900838091), nil, true),
		Entry("Both nil - False", nil, nil, false),
		Entry("RHS == LHS - False",
			NewUnixDuration(1655510399, 900838091), NewUnixDuration(1655510399, 900838091), false),
		Entry("RHS.Nanoseconds != LHS.Nanoseconds - True",
			NewUnixDuration(1655510399, 900838091), NewUnixDuration(1655510399, 0), true),
		Entry("RHS.Seconds != LHS.Seconds - True",
			NewUnixDuration(1655510399, 900838091), NewUnixDuration(1655510000, 900838091), true))

	// Tests that the GreaterThan function works under various data conditions
	DescribeTable("GreaterThan - Conditions",
		func(rhs *UnixDuration, lhs *UnixDuration, greaterThan bool) {
			Expect(rhs.GreaterThan(lhs)).Should(Equal(greaterThan))
		},
		Entry("RHS is nil - False", nil, NewUnixDuration(1655510000, 900838091), false),
		Entry("LHS is nil - True", NewUnixDuration(1655510399, 900838091), nil, true),
		Entry("Both nil - False", nil, nil, false),
		Entry("RHS == LHS - False",
			NewUnixDuration(1655510399, 900838091), NewUnixDuration(1655510399, 900838091), false),
		Entry("RHS.Nanoseconds < LHS.Nanoseconds - False",
			NewUnixDuration(1655510399, 0), NewUnixDuration(1655510399, 900838091), false),
		Entry("RHS.Nanoseconds > LHS.Nanoseconds - True",
			NewUnixDuration(1655510399, 900838091), NewUnixDuration(1655510399, 0), true),
		Entry("RHS.Seconds < LHS.Seconds - False",
			NewUnixDuration(1655510000, 900838091), NewUnixDuration(1655510399, 900838091), false),
		Entry("RHS.Seconds > LHS.Seconds - True",
			NewUnixDuration(1655510399, 900838091), NewUnixDuration(1655510000, 900838091), true))

	// Tests that the GreaterThanOrEqualTo function works under various data conditions
	DescribeTable("GreaterThanOrEqualTo - Conditions",
		func(rhs *UnixDuration, lhs *UnixDuration, gte bool) {
			Expect(rhs.GreaterThanOrEqualTo(lhs)).Should(Equal(gte))
		},
		Entry("RHS is nil - False", nil, NewUnixDuration(1655510000, 900838091), false),
		Entry("LHS is nil - True", NewUnixDuration(1655510399, 900838091), nil, true),
		Entry("Both nil - True", nil, nil, true),
		Entry("RHS == LHS - True",
			NewUnixDuration(1655510399, 900838091), NewUnixDuration(1655510399, 900838091), true),
		Entry("RHS.Nanoseconds < LHS.Nanoseconds - False",
			NewUnixDuration(1655510399, 0), NewUnixDuration(1655510399, 900838091), false),
		Entry("RHS.Nanoseconds > LHS.Nanoseconds - True",
			NewUnixDuration(1655510399, 900838091), NewUnixDuration(1655510399, 0), true),
		Entry("RHS.Seconds < LHS.Seconds - False",
			NewUnixDuration(1655510000, 900838091), NewUnixDuration(1655510399, 900838091), false),
		Entry("RHS.Seconds > LHS.Seconds - True",
			NewUnixDuration(1655510399, 900838091), NewUnixDuration(1655510000, 900838091), true))

	// Tests that the LessThan function works under various data conditions
	DescribeTable("LessThan - Conditions",
		func(rhs *UnixDuration, lhs *UnixDuration, lt bool) {
			Expect(rhs.LessThan(lhs)).Should(Equal(lt))
		},
		Entry("RHS is nil - True", nil, NewUnixDuration(1655510000, 900838091), true),
		Entry("LHS is nil - False", NewUnixDuration(1655510399, 900838091), nil, false),
		Entry("Both nil - False", nil, nil, false),
		Entry("RHS == LHS - False",
			NewUnixDuration(1655510399, 900838091), NewUnixDuration(1655510399, 900838091), false),
		Entry("RHS.Nanoseconds < LHS.Nanoseconds - True",
			NewUnixDuration(1655510399, 0), NewUnixDuration(1655510399, 900838091), true),
		Entry("RHS.Nanoseconds > LHS.Nanoseconds - False",
			NewUnixDuration(1655510399, 900838091), NewUnixDuration(1655510399, 0), false),
		Entry("RHS.Seconds < LHS.Seconds - True",
			NewUnixDuration(1655510000, 900838091), NewUnixDuration(1655510399, 900838091), true),
		Entry("RHS.Seconds > LHS.Seconds - False",
			NewUnixDuration(1655510399, 900838091), NewUnixDuration(1655510000, 900838091), false))

	// Tests that the LessThanOrEqualTo function works under various data conditions
	DescribeTable("LessThanOrEqualTo - Condition",
		func(rhs *UnixDuration, lhs *UnixDuration, lte bool) {
			Expect(rhs.LessThanOrEqualTo(lhs)).Should(Equal(lte))
		},
		Entry("RHS is nil - True", nil, NewUnixDuration(1655510000, 900838091), true),
		Entry("LHS is nil - False", NewUnixDuration(1655510399, 900838091), nil, false),
		Entry("Both nil - True", nil, nil, true),
		Entry("RHS == LHS - True",
			NewUnixDuration(1655510399, 900838091), NewUnixDuration(1655510399, 900838091), true),
		Entry("RHS.Nanoseconds < LHS.Nanoseconds - True",
			NewUnixDuration(1655510399, 0), NewUnixDuration(1655510399, 900838091), true),
		Entry("RHS.Nanoseconds > LHS.Nanoseconds - False",
			NewUnixDuration(1655510399, 900838091), NewUnixDuration(1655510399, 0), false),
		Entry("RHS.Seconds < LHS.Seconds - True",
			NewUnixDuration(1655510000, 900838091), NewUnixDuration(1655510399, 900838091), true),
		Entry("RHS.Seconds > LHS.Seconds - False",
			NewUnixDuration(1655510399, 900838091), NewUnixDuration(1655510000, 900838091), false))

	// Tests the conditions determining whether IsValid will return true or false
	DescribeTable("IsValid - Conditions",
		func(duration *UnixDuration, result bool) {
			Expect(duration.IsValid()).Should(Equal(result))
		},
		Entry("Duration is nil - False", nil, false),
		Entry("Seconds < -10,000 years - False", NewUnixDuration(-315576000001, 0), false),
		Entry("Seconds > 10,000 years - False", NewUnixDuration(315576000001, 0), false),
		Entry("Nanoseconds <= -1e9 - False", NewUnixDuration(2678400, -1000000000), false),
		Entry("Nanoseconds >= 1e9 - False", NewUnixDuration(2678400, 1000000000), false),
		Entry("Seconds > 0, Nanoseconds < 0 - False", NewUnixDuration(2678400, -1000), false),
		Entry("Seconds < 0, Nanoseconds > 0 - False", NewUnixDuration(-2678400, 1000), false),
		Entry("Valid - True", NewUnixDuration(2678400, 1000), true))

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
		Entry("Seconds < -10,000 years - False", NewUnixDuration(-315576000001, 0), true,
			"duration (-315576000001, 0) exceeds -10000 years"),
		Entry("Seconds > 10,000 years - False", NewUnixDuration(315576000001, 0), true,
			"duration (315576000001, 0) exceeds +10000 years"),
		Entry("Nanoseconds <= -1e9 - False", NewUnixDuration(2678400, -1000000000), true,
			"duration (2678400, -1000000000) has out-of-range nanos"),
		Entry("Nanoseconds >= 1e9 - False", NewUnixDuration(2678400, 1000000000), true,
			"duration (2678400, 1000000000) has out-of-range nanos"),
		Entry("Seconds > 0, Nanoseconds < 0 - False", NewUnixDuration(2678400, -1000), true,
			"duration (2678400, -1000) has seconds and nanos with different signs"),
		Entry("Seconds < 0, Nanoseconds > 0 - False", NewUnixDuration(-2678400, 1000), true,
			"duration (-2678400, 1000) has seconds and nanos with different signs"),
		Entry("Valid - True", NewUnixDuration(2678400, 1000), false, ""))
})
