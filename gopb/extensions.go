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

		// Divide the remaining value by the width of our integer value, saving the remainder
		// of the division to our temporary value, r and saving the quotient to the remainder
		rest.DivMod(rest, offset, r)

		// Append the remainder to our list (ensuring that we save the sign value)
		ints = append(ints, r.Int64()*sign)
	}

	// Inject the parts and the exponent into a Decimal value and return it
	return &Decimal{Value: ints, Exp: in.Exponent()}
}

// ToDecimal converts our internal representation of a Decimal to a decimal.Decimal
func (d *Decimal) ToDecimal() *decimal.Decimal {

	// First, create our decimal value with a zero-value
	resp := decimal.New(0, 0)

	// Next, iterate over all our sub-values and add them into the total
	for i, value := range d.Value {

		// Attempt to convert the value to its decimal equivalent based on where it is in the list
		temp := decimal.New(value, int32(i*18)+d.Exp)

		// Add the temporary value to the total
		resp = resp.Add(temp)
	}

	// Finally, return our total
	return &resp
}

// Now constructs a new Timestamp from the current time.
func Now() *UnixTimestamp {
	return NewUnixTimestamp(time.Now())
}

// NewUnixTimestamp constructs a new Timestamp from the provided time.Time.
func NewUnixTimestamp(t time.Time) *UnixTimestamp {
	return &UnixTimestamp{Seconds: int64(t.Unix()), Nanoseconds: int32(t.Nanosecond())}
}

// AsTime converts x to a time.Time.
func (x *UnixTimestamp) AsTime() time.Time {
	return time.Unix(int64(x.GetSeconds()), int64(x.GetNanoseconds())).UTC()
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

	// First, check if the lhs is nil. If it is then return the rhs
	if lhs == nil {
		return rhs
	}

	// Next, add the seconds and nanoseconds to the timestamp
	rhs.Nanoseconds += lhs.Nanoseconds
	rhs.Seconds += lhs.Seconds

	// Now, if the nanoseconds is greater than a second or less than zero then roll them over
	if rhs.Nanoseconds >= 1e9 {
		rhs.Seconds += 1
		rhs.Nanoseconds -= 1e9
	} else if rhs.Nanoseconds < 0 {
		rhs.Seconds -= 1
		rhs.Nanoseconds += 1e9
	}

	// Finally, return the modified timestamp
	return rhs
}

// AddDate adds a number of years, months and days to the time associated with the timestamp
func (rhs *UnixTimestamp) AddDate(years int, months int, days int) *UnixTimestamp {
	return NewUnixTimestamp(rhs.AsTime().AddDate(years, months, days))
}

// AddDuration adds a duration to the timestamp, modifying it. The modified timestamp is then returned
func (rhs *UnixTimestamp) AddDuration(duration time.Duration) *UnixTimestamp {

	// First, extract the seconds from the duration
	seconds := duration.Seconds()
	rhs.Seconds += int64(seconds)

	// Next, if we have a fractional second then convert it to nanoseconds and add it to the total number
	// of nanoseconds; we'll ignore any fractional nanoseconds
	if seconds != math.Floor(seconds) {
		_, frac := math.Modf(seconds)
		rhs.Nanoseconds += int32(0.5 + (frac * 1e9))
	}

	// Now, if the nanoseconds is greater than a second then roll them over
	if rhs.Nanoseconds >= 1e9 {
		rhs.Seconds += 1
		rhs.Nanoseconds -= 1e9
	} else if rhs.Nanoseconds < 0 {
		rhs.Seconds -= 1
		rhs.Nanoseconds += 1e9
	}

	// Finally, return the modified timestamp
	return rhs
}

// AddUnixDuration adds a UnixDuration to the UnixTimestamp, modifying it. The modified timestamp is then returned
func (rhs *UnixTimestamp) AddUnixDuration(lhs *UnixDuration) *UnixTimestamp {

	// First, check if the lhs is nil. If it is then return the rhs
	if lhs == nil {
		return rhs
	}

	// Next, add the seconds and nanoseconds to the timestamp
	rhs.Nanoseconds += lhs.Nanoseconds
	rhs.Seconds += lhs.Seconds

	// Now, if the nanoseconds is greater than a second then roll them over
	if rhs.Nanoseconds >= 1e9 {
		rhs.Seconds += 1
		rhs.Nanoseconds -= 1e9
	} else if rhs.Nanoseconds < 0 {
		rhs.Seconds -= 1
		rhs.Nanoseconds += 1e9
	}

	// Finally, return the modified timestamp
	return rhs
}

// IsValid reports whether the timestamp is valid.
// It is equivalent to CheckValid == nil.
func (x *UnixTimestamp) IsValid() bool {
	return x.check() == 0
}

// CheckValid returns an error if the timestamp is invalid.
// In particular, it checks whether the value represents a date that is
// in the range of 0001-01-01T00:00:00Z to 9999-12-31T23:59:59Z inclusive.
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

	// First, check that the timestamp is long enough for us to parse
	// If it isn't then return an error. Also, check if the string is empty
	// If it is then we're probably looking at an empty timestamp
	if raw == "" {
		timestamp = nil
		return nil
	} else if len(raw) < 10 {
		return fmt.Errorf("value (%s) was not long enough to be converted to a timestamp", raw)
	}

	// Next, attempt to parse the number of seconds to a 64-bit integer
	// If this fails then return an error
	partition := len(raw) - 9
	seconds, err := strconv.ParseInt(raw[:partition], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to convert seconds part to integer, error: %v", err)
	}

	// Now, attempt to parse the number of nanoseconds to a 32-bit integer
	// If this fails then return an error
	nanos, err := strconv.ParseInt(raw[partition:], 10, 32)
	if err != nil {
		return fmt.Errorf("failed to convert nanoseconds part to integer, error: %v", err)
	}

	// Finally, create a new timestamp from the seconds and nanoseconds and then
	// check that the timestamp is valid; return any error that occurs
	timestamp.Seconds = seconds
	timestamp.Nanoseconds = int32(nanos)
	return timestamp.CheckValid()
}

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

// NewUnixDuration constructs a new UnixDuration from the provided time.Duration.
func NewUnixDuration(d time.Duration) *UnixDuration {
	nanos := d.Nanoseconds()
	secs := nanos / 1e9
	nanos -= secs * 1e9
	return &UnixDuration{
		Seconds:     int64(secs),
		Nanoseconds: int32(nanos),
	}
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

	// Next, if the duration is negative then we'll attach a minus sign to the
	// front of the string; otherwise we won't
	if duration.Seconds < 0 || duration.Nanoseconds < 0 {
		duration.Nanoseconds *= -1
	}

	// Finally, convert the timestamp to a UNIX epoch value and return it
	return fmt.Sprintf("%d%d", duration.Seconds, duration.Nanoseconds)
}

// FromString creates a new timestamp from a string
func (duration *UnixDuration) FromString(raw string) error {

	// First, check that the duration is long enough for us to parse
	// If it isn't then return an error. Also, check if the string is empty
	// If it is then we're probably looking at an empty duration
	if raw == "" {
		duration = nil
		return nil
	} else if len(raw) < 10 {
		return fmt.Errorf("value (%s) was not long enough to be converted to a duration", raw)
	}

	// Next, attempt to parse the number of seconds to a 64-bit integer
	// If this fails then return an error
	partition := len(raw) - 9
	seconds, err := strconv.ParseInt(raw[:partition], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to convert seconds part to integer, error: %v", err)
	}

	// Now, attempt to parse the number of nanoseconds to a 32-bit integer
	// If this fails then return an error
	nanos, err := strconv.ParseInt(raw[partition:], 10, 32)
	if err != nil {
		return fmt.Errorf("failed to convert nanoseconds part to integer, error: %v", err)
	}

	// If the number of seconds is less than 0 then the number of nanoseconds must
	// also be less than 0 so multiply them by -1
	if seconds < 0 {
		nanos *= -1
	}

	// Finally, create a new duration from the seconds and nanoseconds and then
	// check that the duration is valid; return any error that occurs
	duration.Seconds = seconds
	duration.Nanoseconds = int32(nanos)
	return duration.CheckValid()
}

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
