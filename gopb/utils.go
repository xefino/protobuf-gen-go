package gopb

import (
	"database/sql/driver"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

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

// MarhsalJSON converts a Duration to JSON
func (duration *UnixDuration) MarshalJSON() ([]byte, error) {
	return []byte(duration.ToEpoch()), nil
}

// MarshalCSV converts a Duration to a CSV format
func (duration *UnixDuration) MarshalCSV() (string, error) {
	return duration.ToEpoch(), nil
}

// Marshaler converts a Duration to a DynamoDB attribute value
func (duration *UnixDuration) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: duration.ToEpoch(),
	}, nil
}

// Value converts a Duration to an SQL value
func (duration *UnixDuration) Value() (driver.Value, error) {
	return driver.Value(duration.ToEpoch()), nil
}

// UnmarshalJSON converts JSON data into a Duration
func (duration *UnixDuration) UnmarshalJSON(data []byte) error {

	// Check if the value is nil; if this is the case then return nil
	if data == nil {
		return nil
	}

	// Otherwise, convert the data from a string into a duration
	return duration.FromString(string(data))
}

// UnmarshalCSV converts a CSV column into a Duration
func (duration *UnixDuration) UnmarshalCSV(raw string) error {
	return duration.FromString(raw)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Duration
func (duration *UnixDuration) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return duration.FromString(string(casted.Value))
	case *types.AttributeValueMemberN:
		return duration.FromString(casted.Value)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return duration.FromString(casted.Value)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a UnixDuration", value)
	}
}

// Scan converts an SQL value into a Duration
func (duration *UnixDuration) Scan(value interface{}) error {

	// Check if the value is nil; if this is the case then return nil
	if value == nil {
		return nil
	}

	// Otherwise, convert the data from a string into a duration
	return duration.FromString(value.(string))
}
