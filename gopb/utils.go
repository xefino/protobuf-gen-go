package gopb

import (
	"database/sql/driver"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

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

	// First, check if the lhs is nil. If it is then return the rhs. Otherwise, if the rhs is nil
	// then set it to a new UnixTimestamp and continue
	if lhs == nil {
		return rhs
	}

	// First, add the seconds and nanoseconds to the timestamp
	rhs.Nanoseconds += lhs.Nanoseconds
	rhs.Seconds += lhs.Seconds

	// Next, if the nanoseconds is greater than a second then roll them over
	if rhs.Nanoseconds > 1e9 {
		rhs.Seconds += 1
		rhs.Nanoseconds -= 1e9
	}

	// Finally, return the modified timestamp
	return rhs
}

// AddDuration adds a duration to the timestamp, modifying it. The modified timestamp is then returned
func (rhs *UnixTimestamp) AddDuration(duration time.Duration) *UnixTimestamp {
	seconds := duration.Seconds()
	rhs.Seconds += int64(seconds)
	if seconds != math.Floor(seconds) {
		_, frac := math.Modf(seconds)
		rhs.Nanoseconds += int32(0.5 + (frac * 1e9))
	}

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

const (
	_ = iota
	invalidNil
	invalidUnderflow
	invalidOverflow
	invalidNanos
)

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

// MarhsalJSON converts a Timestamp to JSON
func (timestamp *UnixTimestamp) MarshalJSON() ([]byte, error) {
	return []byte(timestamp.ToEpoch()), nil
}

// MarshalCSV converts a Timestamp to a CSV format
func (timestamp *UnixTimestamp) MarshalCSV() (string, error) {
	return timestamp.ToEpoch(), nil
}

// Marshaler converts a Timestamp to a DynamoDB attribute value
func (timestamp *UnixTimestamp) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: timestamp.ToEpoch(),
	}, nil
}

// Value converts a Timestamp to an SQL value
func (timestamp *UnixTimestamp) Value() (driver.Value, error) {
	return driver.Value(timestamp.ToEpoch()), nil
}

// ToString converts the timestamp to a UNIX epoch value
func (timestamp *UnixTimestamp) ToEpoch() string {

	// If the timestamp is nil then return an empty value
	if timestamp == nil {
		return ""
	}

	// Otherwise, convert the timestamp to a UNIX epoch value and return it
	return fmt.Sprintf("%d%09d", timestamp.Seconds, timestamp.Nanoseconds)
}

// UnmarshalJSON converts JSON data into a Timestamp
func (timestamp *UnixTimestamp) UnmarshalJSON(data []byte) error {

	// Check if the value is nil; if this is the case then return nil
	if data == nil {
		return nil
	}

	// Otherwise, convert the data from a string into a timestamp
	return timestamp.FromString(string(data))
}

// UnmarshalCSV converts a CSV column into a Timestamp
func (timestamp *UnixTimestamp) UnmarshalCSV(raw string) error {
	return timestamp.FromString(raw)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a timestamp
func (timestamp *UnixTimestamp) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return timestamp.FromString(string(casted.Value))
	case *types.AttributeValueMemberN:
		return timestamp.FromString(casted.Value)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return timestamp.FromString(casted.Value)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a UnixTimestamp", value)
	}
}

// Scan converts an SQL value into a Timestamp
func (timestamp *UnixTimestamp) Scan(value interface{}) error {

	// Check if the value is nil; if this is the case then return nil
	if value == nil {
		return nil
	}

	// Otherwise, convert the data from a string into a timestamp
	return timestamp.FromString(value.(string))
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
