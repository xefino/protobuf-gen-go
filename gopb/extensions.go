package gopb

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

// Size of the integer values we want to save (designed to fit inside an int64)
var offset = big.NewInt(1e18)

const width = 18

// The number of seconds in a day
const (
	secondsInMinute = 60
	secondsInHour   = 60 * secondsInMinute
	secondsInDay    = 24 * secondsInHour
	nanosPerSecond  = 1000000000
)

// NewFromDecimal creates a new representation of our Decimal from a decimal.Decimal
func NewFromDecimal(in decimal.Decimal) *Decimal {
	coefficient := in.Coefficient()
	sign := int64(coefficient.Sign())

	// First, get the absolute value of the coefficient
	rest := new(big.Int)
	rest.Abs(coefficient)

	// Next, iteratively div-mod the coefficient until there's no data remaining
	var ints []int64
	r := new(big.Int)
	for rest.BitLen() != 0 {

		// Divide the remaining value by the width of our integer value, saving the remainder of the
		// division to our temporary value, r and saving the quotient to the remainder
		rest.DivMod(rest, offset, r)

		// Append the remainder to our list (ensuring that we save the sign value)
		ints = append(ints, r.Int64()*sign)
	}

	// Inject the parts and the exponent into a Decimal value and return it
	return &Decimal{Parts: ints, Exp: in.Exponent()}
}

// ToDecimal converts our internal representation of a Decimal to a decimal.Decimal
func (d *Decimal) ToDecimal() *decimal.Decimal {

	// First, create our decimal value with a zero-value
	resp := decimal.New(0, 0)

	// Next, iterate over all our sub-values and add them into the total
	for i, value := range d.Parts {

		// Attempt to convert the value to its decimal equivalent based on where it is in the list
		temp := decimal.New(value, int32(i*width)+d.Exp)

		// Add the temporary value to the total
		resp = resp.Add(temp)
	}

	// Finally, return our total
	return &resp
}

// ToString converts a Decimal object to its string representation
func (d *Decimal) ToString() string {
	return d.ToDecimal().String()
}

// FromString converts a string representation to a Decimal object
func (d *Decimal) FromString(raw string) error {
	dec, err := decimal.NewFromString(raw)
	if err != nil {
		return err
	}

	*d = *NewFromDecimal(dec)
	return nil
}

// Now constructs a new Timestamp from the current time.
func Now() *UnixTimestamp {
	return NewFromTime(time.Now())
}

// NewUnixTimestamp creates a new UnixTimestamp from the seconds and nanoseconds with which the
// timestamp should be associated
func NewUnixTimestamp(seconds int64, nanos int32) *UnixTimestamp {
	return &UnixTimestamp{
		Seconds:     seconds,
		Nanoseconds: nanos,
	}
}

// NewFromTime constructs a new Timestamp from the provided time.Time.
func NewFromTime(t time.Time) *UnixTimestamp {
	return NewUnixTimestamp(int64(t.Unix()), int32(t.Nanosecond()))
}

// Helper function that returns the greater of two UnixTimestamp objects
func maxTimestampInner(a *UnixTimestamp, b *UnixTimestamp) *UnixTimestamp {

	// First, check if either of the inputs is nil; in this case we'll return the other
	if a == nil {
		return b
	} else if b == nil {
		return a
	}

	// Next, neither of the inputs is nil so we'll compare the seconds and nanoseconds on each and
	// return a if it is larger than b
	if a.Seconds > b.Seconds || (a.Seconds == b.Seconds && a.Nanoseconds > b.Nanoseconds) {
		return a
	}

	// Finally, if we reached this point then b is larger than a so return b
	return b
}

// MaxTimestamp returns the greatest of a series of at least two UnixTimestamp objects
func MaxTimestamp(a *UnixTimestamp, b *UnixTimestamp, others ...*UnixTimestamp) *UnixTimestamp {
	result := maxTimestampInner(a, b)
	for _, other := range others {
		result = maxTimestampInner(result, other)
	}

	return result
}

// Helper function that returns the lesser of two UnixTimestamp objects
func minTimestampInner(a *UnixTimestamp, b *UnixTimestamp) *UnixTimestamp {

	// First, check if either of the inputs is nil; in this case we'll return it
	if a == nil {
		return a
	} else if b == nil {
		return b
	}

	// Next, neither of the inputs is nil so we'll compare the seconds and nanoseconds on each and
	// return a if it is smaller than b
	if a.Seconds < b.Seconds || (a.Seconds == b.Seconds && a.Nanoseconds < b.Nanoseconds) {
		return a
	}

	// Finally, if we reached this point then b is smaller than a so return b
	return b
}

// MinTimestamp returns the least of a series of at least two UnixTimestamp objects
func MinTimestamp(a *UnixTimestamp, b *UnixTimestamp, others ...*UnixTimestamp) *UnixTimestamp {
	result := minTimestampInner(a, b)
	for _, other := range others {
		result = minTimestampInner(result, other)
	}

	return result
}

// AsTime converts x to a time.Time.
func (x *UnixTimestamp) AsTime() time.Time {
	return time.Unix(int64(x.GetSeconds()), int64(x.GetNanoseconds())).UTC()
}

// Copy creates a new UnixTimestamp from an existing UnixTimestamp
func (x *UnixTimestamp) Copy() *UnixTimestamp {
	return NewUnixTimestamp(x.Seconds, x.Nanoseconds)
}

// Equals returns true if rhs is equal to lhs, false otherwise
func (rhs *UnixTimestamp) Equals(lhs *UnixTimestamp) bool {
	if rhs != nil && lhs != nil {
		return rhs.Seconds == lhs.Seconds && rhs.Nanoseconds == lhs.Nanoseconds
	} else {
		return rhs == lhs
	}
}

// NotEquals returns true if rhs is not equal to lhs, false otherwise
func (rhs *UnixTimestamp) NotEquals(lhs *UnixTimestamp) bool {
	return !rhs.Equals(lhs)
}

// GreaterThan returns true if rhs represents a later time than lhs, or false otherwise
func (rhs *UnixTimestamp) GreaterThan(lhs *UnixTimestamp) bool {
	if rhs != nil && lhs != nil {
		return rhs.Seconds > lhs.Seconds ||
			(rhs.Seconds == lhs.Seconds && rhs.Nanoseconds > lhs.Nanoseconds)
	} else {
		return rhs != nil
	}
}

// GreaterThanOrEqualTo returns true if rhs represents a time at least as late as lhs, or false otherwise
func (rhs *UnixTimestamp) GreaterThanOrEqualTo(lhs *UnixTimestamp) bool {
	return !rhs.LessThan(lhs)
}

// LessThan returns true if rhs represents an earlier time than lhs, or false otherwise
func (rhs *UnixTimestamp) LessThan(lhs *UnixTimestamp) bool {
	if rhs != nil && lhs != nil {
		return rhs.Seconds < lhs.Seconds ||
			(rhs.Seconds == lhs.Seconds && rhs.Nanoseconds < lhs.Nanoseconds)
	} else {
		return lhs != nil
	}
}

// LessThanOrEqualTo returns true if rhs represents a time at least as early as lhs, or false otherwise
func (rhs *UnixTimestamp) LessThanOrEqualTo(lhs *UnixTimestamp) bool {
	return !rhs.GreaterThan(lhs)
}

// Add adds a timestamp to another timestamp, modifying it. The modified timestamp is then returned
func (rhs *UnixTimestamp) Add(lhs *UnixTimestamp) *UnixTimestamp {
	ans := rhs.Copy()

	// First, check if the lhs is nil. If it is then return the rhs
	if lhs == nil {
		return ans
	}

	// Next, add the seconds and nanoseconds to the timestamp
	ans.Nanoseconds += lhs.Nanoseconds
	ans.Seconds += lhs.Seconds

	// Now, if the nanoseconds is greater than a second or less than zero then roll them over
	if ans.Nanoseconds >= 1e9 {
		ans.Seconds += 1
		ans.Nanoseconds -= 1e9
	} else if ans.Nanoseconds < 0 {
		ans.Seconds -= 1
		ans.Nanoseconds += 1e9
	}

	// Finally, return the modified timestamp
	return ans
}

// AddDate adds a number of years, months and days to the time associated with the timestamp
func (rhs *UnixTimestamp) AddDate(years int, months int, days int) *UnixTimestamp {
	return NewFromTime(rhs.AsTime().AddDate(years, months, days))
}

// AddDuration adds a UnixDuration to the UnixTimestamp, modifying it. The modified timestamp is then returned
func (rhs *UnixTimestamp) AddDuration(lhs *UnixDuration) *UnixTimestamp {
	ans := rhs.Copy()

	// First, check if the lhs is nil. If it is then return the rhs
	if lhs == nil {
		return ans
	}

	// Next, add the seconds and nanoseconds to the timestamp
	ans.Nanoseconds += lhs.Nanoseconds
	ans.Seconds += lhs.Seconds

	// Now, if the nanoseconds is greater than a second then roll them over
	if ans.Nanoseconds >= 1e9 {
		ans.Seconds += 1
		ans.Nanoseconds -= 1e9
	} else if ans.Nanoseconds < 0 {
		ans.Seconds -= 1
		ans.Nanoseconds += 1e9
	}

	// Finally, return the modified timestamp
	return ans
}

// Difference calculates the difference between two UnixTimestamp objects, returning a UnixDuration
func (rhs *UnixTimestamp) Difference(lhs *UnixTimestamp) *UnixDuration {

	// First, calculate the difference between the seconds and nanoseconds
	seconds := rhs.Seconds - lhs.Seconds
	nanos := rhs.Nanoseconds - lhs.Nanoseconds

	// Next, if we have seconds and nanoseconds of differing signs then we'll need to modify the results
	// so that they have the same sign. If the seconds is greater than 0 and the nanos is less than 0
	// then we'll subtract a second from seconds and add it back to nanos. If the seconds is less than 0
	// and the nanos is greater than zero then we'll do the reverse
	if seconds > 0 && nanos < 0 {
		seconds -= 1
		nanos += 1e9
	} else if seconds < 0 && nanos > 0 {
		seconds += 1
		nanos -= 1e9
	}

	// Finally, create a new UnixDuration from the seconds and nanoseconds and return it
	return NewUnixDuration(seconds, nanos)
}

// NextDay returns a new UnixTimestamp, set to the start of the day of the UnixTimestamp
func (rhs *UnixTimestamp) NextDay() *UnixTimestamp {
	return NewUnixTimestamp(secondsInDay*((rhs.Seconds/secondsInDay)+1), 0)
}

// SecondDown creates a new UnixTimestamp, snapped to the start of the current second
func (x *UnixTimestamp) SecondDown() *UnixTimestamp {
	return NewUnixTimestamp(x.Seconds, 0)
}

// SecondUp creates a new UnixTimestamp, snapped to the start of the next second unless the timestamp
// is a whole second, in which case it makes a copy of the UnixTimestamp and returns that
func (x *UnixTimestamp) SecondUp() *UnixTimestamp {
	if x.Nanoseconds > 0 {
		return NewUnixTimestamp(x.Seconds+1, 0)
	} else {
		return x.Copy()
	}
}

// MinuteDown creates a new UnixTimestamp, snapped to the start of the current minute
func (x *UnixTimestamp) MinuteDown() *UnixTimestamp {
	return NewUnixTimestamp(secondsInMinute*(x.Seconds/secondsInMinute), 0)
}

// MinuteUp creates a new UnixTimestamp, snapped to the start of the next minute unless the timestamp
// is a whole minute, in which case it makes a copy of the UnixTimestamp and returns that
func (x *UnixTimestamp) MinuteUp() *UnixTimestamp {
	if x.Nanoseconds > 0 || x.Seconds%secondsInMinute > 0 {
		return NewUnixTimestamp(secondsInMinute*((x.Seconds/secondsInMinute)+1), 0)
	} else {
		return x.Copy()
	}
}

// HourDown creates a new UnixTimestamp, snapped to the start of the current hour
func (x *UnixTimestamp) HourDown() *UnixTimestamp {
	return NewUnixTimestamp(secondsInHour*(x.Seconds/secondsInHour), 0)
}

// HourUp creates a new UnixTimestamp, snapped to the start of the next hour unless the timestamp
// is a whole hour, in which case it makes a copy of the UnixTimestamp and returns that
func (x *UnixTimestamp) HourUp() *UnixTimestamp {
	if x.Nanoseconds > 0 || x.Seconds%secondsInHour > 0 {
		return NewUnixTimestamp(secondsInHour*((x.Seconds/secondsInHour)+1), 0)
	} else {
		return x.Copy()
	}
}

// DayDown creates a new UnixTimestamp, snapped to the start of the current day
func (x *UnixTimestamp) DayDown() *UnixTimestamp {
	return NewUnixTimestamp(secondsInDay*(x.Seconds/secondsInDay), 0)
}

// DayUp creates a new UnixTimestamp, snapped to the start of the next day unless the timestamp
// is a whole day, in which case it makes a copy of the UnixTimestamp and returns that
func (x *UnixTimestamp) DayUp() *UnixTimestamp {
	if x.Nanoseconds > 0 || x.Seconds%secondsInDay > 0 {
		return x.NextDay()
	} else {
		return x.Copy()
	}
}

// WeekDown creates a new UnixTimestamp, snapped to the start of the current week
func (x *UnixTimestamp) WeekDown() *UnixTimestamp {
	t := x.AsTime()
	switch t.Weekday() {
	case time.Monday:
		return NewFromTime(time.Date(t.Year(), t.Month(), t.Day()-1, 0, 0, 0, 0, time.UTC))
	case time.Tuesday:
		return NewFromTime(time.Date(t.Year(), t.Month(), t.Day()-2, 0, 0, 0, 0, time.UTC))
	case time.Wednesday:
		return NewFromTime(time.Date(t.Year(), t.Month(), t.Day()-3, 0, 0, 0, 0, time.UTC))
	case time.Thursday:
		return NewFromTime(time.Date(t.Year(), t.Month(), t.Day()-4, 0, 0, 0, 0, time.UTC))
	case time.Friday:
		return NewFromTime(time.Date(t.Year(), t.Month(), t.Day()-5, 0, 0, 0, 0, time.UTC))
	case time.Saturday:
		return NewFromTime(time.Date(t.Year(), t.Month(), t.Day()-6, 0, 0, 0, 0, time.UTC))
	default:
		return x.DayDown()
	}
}

// WeekUp creates a new UnixTimestamp, snapped to the start of the next week unless the timestamp
// is a whole week, in which case it makes a copy of the UnixTimestamp and returns that
func (x *UnixTimestamp) WeekUp() *UnixTimestamp {
	t := x.AsTime()
	switch t.Weekday() {
	case time.Monday:
		return NewFromTime(time.Date(t.Year(), t.Month(), t.Day()+6, 0, 0, 0, 0, time.UTC))
	case time.Tuesday:
		return NewFromTime(time.Date(t.Year(), t.Month(), t.Day()+5, 0, 0, 0, 0, time.UTC))
	case time.Wednesday:
		return NewFromTime(time.Date(t.Year(), t.Month(), t.Day()+4, 0, 0, 0, 0, time.UTC))
	case time.Thursday:
		return NewFromTime(time.Date(t.Year(), t.Month(), t.Day()+3, 0, 0, 0, 0, time.UTC))
	case time.Friday:
		return NewFromTime(time.Date(t.Year(), t.Month(), t.Day()+2, 0, 0, 0, 0, time.UTC))
	case time.Saturday:
		return NewFromTime(time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, time.UTC))
	default:
		if x.Nanoseconds > 0 || x.Seconds%secondsInDay > 0 {
			return NewFromTime(time.Date(t.Year(), t.Month(), t.Day()+7, 0, 0, 0, 0, time.UTC))
		} else {
			return x.Copy()
		}
	}
}

// MonthDown creates a new UnixTimestamp, snapped to the start of the current month
func (x *UnixTimestamp) MonthDown() *UnixTimestamp {
	t := x.AsTime()
	return NewFromTime(time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC))
}

// MonthUp creates a new UnixTimestamp, snapped to the start of the next month unless the timestamp
// is a whole month, in which case it makes a copy of the UnixTimestamp and returns that
func (x *UnixTimestamp) MonthUp() *UnixTimestamp {
	t := x.AsTime()
	if t.Nanosecond() > 0 || t.Second() > 0 || t.Minute() > 0 || t.Hour() > 0 || t.Day() != 1 {
		return x.MonthDown().AddDate(0, 1, 0)
	} else {
		return x.Copy()
	}
}

// QuarterDown creates a new UnixTimestamp, snapped to the start of the current quarter
func (x *UnixTimestamp) QuarterDown() *UnixTimestamp {
	t := x.AsTime()
	switch t.Month() {
	case time.January, time.February, time.March:
		return NewFromTime(time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, time.UTC))
	case time.April, time.May, time.June:
		return NewFromTime(time.Date(t.Year(), time.April, 1, 0, 0, 0, 0, time.UTC))
	case time.July, time.August, time.September:
		return NewFromTime(time.Date(t.Year(), time.July, 1, 0, 0, 0, 0, time.UTC))
	default:
		return NewFromTime(time.Date(t.Year(), time.October, 1, 0, 0, 0, 0, time.UTC))
	}
}

// QuarterUp creates a new UnixTimestamp, snapped to the start of the next quarter unless the timestamp
// is a whole quarter, in which case it makes a copy of the UnixTimestamp and returns that
func (x *UnixTimestamp) QuarterUp() *UnixTimestamp {
	t := x.AsTime()
	month := t.Month()
	if t.Nanosecond() > 0 || t.Second() > 0 || t.Minute() > 0 || t.Hour() > 0 || t.Day() != 1 ||
		(month != time.January && month != time.April && month != time.July && month != time.October) {
		return x.QuarterDown().AddDate(0, 3, 0)
	} else {
		return x.Copy()
	}
}

// YearDown creates a new UnixTimestamp, snapped to the start of the current year
func (x *UnixTimestamp) YearDown() *UnixTimestamp {
	return NewFromTime(time.Date(x.AsTime().Year(), 1, 1, 0, 0, 0, 0, time.UTC))
}

// MonthUp creates a new UnixTimestamp, snapped to the start of the next year unless the timestamp
// is a whole year, in which case it makes a copy of the UnixTimestamp and returns that
func (x *UnixTimestamp) YearUp() *UnixTimestamp {
	t := x.AsTime()
	if t.Nanosecond() > 0 || t.Second() > 0 || t.Minute() > 0 || t.Hour() > 0 || t.Day() != 1 || t.Month() != time.January {
		return x.YearDown().AddDate(1, 0, 0)
	} else {
		return x.Copy()
	}
}

// IsWhole checks whether or not the duration fits into the UnixTimestamp provided. This function will
// return true if the duration evenly fits into the UnixTimestamp, or false otherwise. This can be used
// to see if the UnixTimestamp represents the beginning of an arbitrary time period
func (rhs *UnixTimestamp) IsWhole(duration time.Duration) bool {
	quo, _ := decimal.NewFromString(rhs.ToEpoch())
	div := big.NewInt(duration.Nanoseconds())
	rem := new(big.Int)
	new(big.Int).QuoRem(quo.BigInt(), div, rem)
	return rem.Int64() == 0
}

// IsWhole checks whether or not the UnixDuration fits into the UnixTimestamp provided. This function
// will return true if the UnixDuration evenly fits into the UnixTimestamp, or false otherwise. This
// can be used to see if the UnixTimestamp represents the beginning of an arbitrary time period
func (rhs *UnixTimestamp) IsWholeUnix(duration *UnixDuration) bool {
	quo, _ := decimal.NewFromString(rhs.ToEpoch())
	div, _ := decimal.NewFromString(duration.ToEpoch())
	rem := new(big.Int)
	new(big.Int).QuoRem(quo.BigInt(), div.BigInt(), rem)
	return rem.Int64() == 0
}

// IsValid reports whether the timestamp is valid. It is equivalent to CheckValid == nil.
func (x *UnixTimestamp) IsValid() bool {
	return x.check() == 0
}

// CheckValid returns an error if the timestamp is invalid. In particular, it checks whether the value
// represents a date that is in the range of 0001-01-01T00:00:00Z to 9999-12-31T23:59:59Z inclusive.
// An error is reported for a nil Timestamp.
func (x *UnixTimestamp) CheckValid() error {
	switch x.check() {
	case invalidNil:
		return fmt.Errorf("invalid nil Timestamp")
	case invalidUnderflow:
		return fmt.Errorf("timestamp (%d, %d) before 0001-01-01", x.Seconds, x.Nanoseconds)
	case invalidOverflow:
		return fmt.Errorf("timestamp (%d, %d) after 9999-12-31", x.Seconds, x.Nanoseconds)
	case invalidNanos:
		return fmt.Errorf("timestamp (%d, %d) has out-of-range nanos", x.Seconds, x.Nanoseconds)
	default:
		return nil
	}
}

// ToDate converts a UnixTimestamp to a date string
func (timestamp *UnixTimestamp) ToDate() string {
	time := timestamp.AsTime()
	return fmt.Sprintf("%04d-%02d-%02d", time.Year(), time.Month(), time.Day())
}

// ToEpoch converts the timestamp to a UNIX epoch value
func (timestamp *UnixTimestamp) ToEpoch() string {

	// If the timestamp is nil then return an empty value
	if timestamp == nil {
		return ""
	}

	// Otherwise, convert the timestamp to a UNIX epoch value and return it
	return fmt.Sprintf("%d%09d", timestamp.Seconds, timestamp.Nanoseconds)
}

// FromString creates a new timestamp from a string
func (timestamp *UnixTimestamp) FromString(raw string) error {

	// First, check that the timestamp is long enough for us to parse. If it isn't then return an error.
	// Also, check if the string is empty. If it is then we're probably looking at an empty timestamp
	if raw == "" {
		timestamp = nil
		return nil
	} else if len(raw) < 10 {
		return fmt.Errorf("value (%s) was not long enough to be converted to a timestamp", raw)
	}

	// Next, attempt to parse the number of seconds to a 64-bit integer. If this fails then return an error
	partition := len(raw) - 9
	seconds, err := strconv.ParseInt(raw[:partition], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to convert seconds part to integer, error: %v", err)
	}

	// Now, attempt to parse the number of nanoseconds to a 32-bit integer. If this fails then return an error
	nanos, err := strconv.ParseInt(raw[partition:], 10, 32)
	if err != nil {
		return fmt.Errorf("failed to convert nanoseconds part to integer, error: %v", err)
	}

	// Finally, create a new timestamp from the seconds and nanoseconds and then check that the timestamp
	// is valid; return any error that occurs
	timestamp.Seconds = seconds
	timestamp.Nanoseconds = int32(nanos)
	return timestamp.CheckValid()
}

// Helper function that checks if a given timestamp is valid
func (x *UnixTimestamp) check() uint {
	const minTimestamp = -62135596800  // Seconds between 1970-01-01T00:00:00Z and 0001-01-01T00:00:00Z, inclusive
	const maxTimestamp = +253402300799 // Seconds between 1970-01-01T00:00:00Z and 9999-12-31T23:59:59Z, inclusive
	secs := x.GetSeconds()
	nanos := x.GetNanoseconds()
	switch {
	case x == nil:
		return invalidNil
	case secs < minTimestamp:
		return invalidUnderflow
	case secs > maxTimestamp:
		return invalidOverflow
	case nanos < 0 || nanos >= 1e9:
		return invalidNanos
	default:
		return 0
	}
}

// NewUnixDuration creates a new UnixDuration from the seconds and nanoseconds with which the
// duration should be associated
func NewUnixDuration(seconds int64, nanos int32) *UnixDuration {
	return &UnixDuration{
		Seconds:     seconds,
		Nanoseconds: nanos,
	}
}

// NewFromDuration constructs a new UnixDuration from the provided time.Duration.
func NewFromDuration(d time.Duration) *UnixDuration {
	nanos := d.Nanoseconds()
	secs := nanos / 1e9
	nanos -= secs * 1e9
	return NewUnixDuration(secs, int32(nanos))
}

// Helper function that returns the greater of two UnixDuration objects
func maxDurationInner(a *UnixDuration, b *UnixDuration) *UnixDuration {

	// First, check if either of the inputs is nil; in this case we'll return the other
	if a == nil {
		return b
	} else if b == nil {
		return a
	}

	// Next, neither of the inputs is nil so we'll compare the seconds and nanoseconds on each and
	// return a if it is larger than b
	if a.Seconds > b.Seconds || (a.Seconds == b.Seconds && a.Nanoseconds > b.Nanoseconds) {
		return a
	}

	// Finally, if we reached this point then b is larger than a so return b
	return b
}

// MaxDuration returns the greatest of a series of at least two UnixDuration objects
func MaxDuration(a *UnixDuration, b *UnixDuration, others ...*UnixDuration) *UnixDuration {
	result := maxDurationInner(a, b)
	for _, other := range others {
		result = maxDurationInner(result, other)
	}

	return result
}

// Helper function that returns the lesser of two UnixDuration objects
func minDurationInner(a *UnixDuration, b *UnixDuration) *UnixDuration {

	// First, check if either of the inputs is nil; in this case we'll return it
	if a == nil {
		return a
	} else if b == nil {
		return b
	}

	// Next, neither of the inputs is nil so we'll compare the seconds and nanoseconds on each and
	// return a if it is smaller than b
	if a.Seconds < b.Seconds || (a.Seconds == b.Seconds && a.Nanoseconds < b.Nanoseconds) {
		return a
	}

	// Finally, if we reached this point then b is smaller than a so return b
	return b
}

// MinDuration returns the least of a series of at least two UnixDuration objects
func MinDuration(a *UnixDuration, b *UnixDuration, others ...*UnixDuration) *UnixDuration {
	result := minDurationInner(a, b)
	for _, other := range others {
		result = minDurationInner(result, other)
	}

	return result
}

// AsDuration converts x to a time.Duration, returning an error in the event of an overflow
func (x *UnixDuration) AsDuration() (time.Duration, error) {

	// First, get the seconds and nanoseconds from the Unix duration
	secs := x.GetSeconds()
	nanos := x.GetNanoseconds()

	// Next, attempt to set the seconds on the duration; if the Unix duration contains too many seconds
	// then return an error as this represents an overflow/underflow error
	duration := time.Duration(secs) * time.Second
	if duration/time.Second != time.Duration(secs) {
		return time.Duration(0), fmt.Errorf("Seconds count was malformed")
	}

	// Now, add the nanoseconds to the duration; if the additional results in a duration of a different
	// sign from the Unix duration then return an error
	duration += time.Duration(nanos) * time.Nanosecond
	if secs < 0 && nanos < 0 && duration > 0 {
		return time.Duration(math.MinInt64), fmt.Errorf("Duration underflow")
	} else if secs > 0 && nanos > 0 && duration < 0 {
		return time.Duration(math.MaxInt64), fmt.Errorf("Duration overflow")
	}

	// Finally, return the duration
	return duration, nil
}

// Equals returns true if rhs is equal to lhs, false otherwise
func (rhs *UnixDuration) Equals(lhs *UnixDuration) bool {
	if rhs != nil && lhs != nil {
		return rhs.Seconds == lhs.Seconds && rhs.Nanoseconds == lhs.Nanoseconds
	} else {
		return rhs == lhs
	}
}

// NotEquals returns true if rhs is not equal to lhs, false otherwise
func (rhs *UnixDuration) NotEquals(lhs *UnixDuration) bool {
	return !rhs.Equals(lhs)
}

// GreaterThan returns true if rhs represents a larger duration than lhs, or false otherwise
func (rhs *UnixDuration) GreaterThan(lhs *UnixDuration) bool {
	if rhs != nil && lhs != nil {
		return rhs.Seconds > lhs.Seconds ||
			(rhs.Seconds == lhs.Seconds && rhs.Nanoseconds > lhs.Nanoseconds)
	} else {
		return rhs != nil
	}
}

// GreaterThanOrEqualTo returns true if rhs represents a duration at least as large as lhs, or false otherwise
func (rhs *UnixDuration) GreaterThanOrEqualTo(lhs *UnixDuration) bool {
	return !rhs.LessThan(lhs)
}

// LessThan returns true if rhs represents a smaller duration than lhs, or false otherwise
func (rhs *UnixDuration) LessThan(lhs *UnixDuration) bool {
	if rhs != nil && lhs != nil {
		return rhs.Seconds < lhs.Seconds ||
			(rhs.Seconds == lhs.Seconds && rhs.Nanoseconds < lhs.Nanoseconds)
	} else {
		return lhs != nil
	}
}

// LessThanOrEqualTo returns true if rhs represents a duration at least as small as lhs, or false otherwise
func (rhs *UnixDuration) LessThanOrEqualTo(lhs *UnixDuration) bool {
	return !rhs.GreaterThan(lhs)
}

// IsValid reports whether the duration is valid. It is equivalent to CheckValid == nil.
func (x *UnixDuration) IsValid() bool {
	return x.check() == 0
}

// CheckValid returns an error if the duration is invalid. In particular, it checks whether the value
// is within the range of -10000 years to +10000 years inclusive. An error is reported for a nil Duration.
func (x *UnixDuration) CheckValid() error {
	switch x.check() {
	case invalidNil:
		return fmt.Errorf("invalid nil Duration")
	case invalidUnderflow:
		return fmt.Errorf("duration (%v, %v) exceeds -10000 years", x.Seconds, x.Nanoseconds)
	case invalidOverflow:
		return fmt.Errorf("duration (%v, %v) exceeds +10000 years", x.Seconds, x.Nanoseconds)
	case invalidNanosRange:
		return fmt.Errorf("duration (%v, %v) has out-of-range nanos", x.Seconds, x.Nanoseconds)
	case invalidNanosSign:
		return fmt.Errorf("duration (%v, %v) has seconds and nanos with different signs", x.Seconds, x.Nanoseconds)
	default:
		return nil
	}
}

// ToEpoch converts the timestamp to a UNIX epoch value
func (duration *UnixDuration) ToEpoch() string {

	// First, if the timestamp is nil then return an empty value
	if duration == nil {
		return ""
	}

	// Next, if the duration is negative then we'll attach a minus sign to the front of the string;
	// otherwise we won't
	if duration.Seconds < 0 || duration.Nanoseconds < 0 {
		duration.Nanoseconds *= -1
	}

	// Finally, convert the timestamp to a UNIX epoch value and return it
	return fmt.Sprintf("%d%09d", duration.Seconds, duration.Nanoseconds)
}

// FromString creates a new timestamp from a string
func (duration *UnixDuration) FromString(raw string) error {

	// First, check that the duration is long enough for us to parse. If it isn't then return an error.
	// Also, check if the string is empty. If it is then we're probably looking at an empty duration
	if raw == "" {
		duration = nil
		return nil
	} else if len(raw) < 10 {
		return fmt.Errorf("value (%s) was not long enough to be converted to a duration", raw)
	}

	// Next, attempt to parse the number of seconds to a 64-bit integer. If this fails then return an error
	partition := len(raw) - 9
	seconds, err := strconv.ParseInt(raw[:partition], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to convert seconds part to integer, error: %v", err)
	}

	// Now, attempt to parse the number of nanoseconds to a 32-bit integer. If this fails then return an error
	nanos, err := strconv.ParseInt(raw[partition:], 10, 32)
	if err != nil {
		return fmt.Errorf("failed to convert nanoseconds part to integer, error: %v", err)
	}

	// If the number of seconds is less than 0 then the number of nanoseconds must also be less than
	// 0 so multiply them by -1
	if seconds < 0 {
		nanos *= -1
	}

	// Finally, create a new duration from the seconds and nanoseconds and then check that the duration
	// is valid; return any error that occurs
	duration.Seconds = seconds
	duration.Nanoseconds = int32(nanos)
	return duration.CheckValid()
}

// Helper function that checks if a given duration is valid
func (x *UnixDuration) check() uint {
	const absDuration = 315576000000 // 10000yr * 365.25day/yr * 24hr/day * 60min/hr * 60sec/min
	secs := x.GetSeconds()
	nanos := x.GetNanoseconds()
	switch {
	case x == nil:
		return invalidNil
	case secs < -absDuration:
		return invalidUnderflow
	case secs > +absDuration:
		return invalidOverflow
	case nanos <= -1e9 || nanos >= +1e9:
		return invalidNanosRange
	case (secs > 0 && nanos < 0) || (secs < 0 && nanos > 0):
		return invalidNanosSign
	default:
		return 0
	}
}

const (
	_ = iota
	invalidNil
	invalidUnderflow
	invalidOverflow
	invalidNanos
	invalidNanosRange
	invalidNanosSign
)
