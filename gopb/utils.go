package gopb

import (
	"database/sql/driver"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/shopspring/decimal"
	"github.com/xefino/protobuf-gen-go/utils"
	"gopkg.in/yaml.v3"
)

// ProviderAlternates contains alternate values for the Provider enum
var ProviderAlternates = map[string]Provider{
	"":        Provider_None,
	"polygon": Provider_Polygon,
}

// ProviderMapping contains alternate names for the Provider enum
var ProviderMapping = map[Provider]string{
	Provider_None:    "",
	Provider_Polygon: "polygon",
}

// MarhsalJSON converts a Decimal to JSON
func (d *Decimal) MarshalJSON() ([]byte, error) {
	return []byte(d.ToString()), nil
}

// MarshalCSV converts a Decimal to a CSV format
func (d *Decimal) MarshalCSV() (string, error) {
	return d.ToString(), nil
}

// Marshaler converts a Decimal to a DynamoDB attribute value
func (d *Decimal) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberN{
		Value: d.ToString(),
	}, nil
}

// Value converts a Decimal to an SQL value
func (d *Decimal) Value() (driver.Value, error) {
	return driver.Value(d.ToString()), nil
}

// UnmarshalJSON converts JSON data into a Decimal
func (d *Decimal) UnmarshalJSON(data []byte) error {

	// Check if the value is nil; if this is the case then return nil
	if data == nil {
		return nil
	}

	// Otherwise, convert the data from a string into a timestamp
	return d.FromString(string(data))
}

// UnmarshalCSV converts a CSV column into a Decimal
func (d *Decimal) UnmarshalCSV(raw string) error {
	return d.FromString(raw)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Decimal
func (d *Decimal) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return d.FromString(string(casted.Value))
	case *types.AttributeValueMemberN:
		return d.FromString(casted.Value)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return d.FromString(casted.Value)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Decimal", value)
	}
}

// Scan converts an SQL value into a Decimal
func (d *Decimal) Scan(value interface{}) error {

	// Check if the value is nil; if this is the case then return nil
	if value == nil {
		return nil
	}

	// Based on the type of the value we're working with, we'll convert the decimal from its implied
	// type to a Decimal; if this fails or the type isn't one we recognized then we'll return an error
	switch casted := value.(type) {
	case []byte:
		return d.FromString(string(casted))
	case float64:
		*d = *NewFromDecimal(decimal.NewFromFloat(casted))
	case int64:
		*d = *NewFromDecimal(decimal.NewFromInt(casted))
	case string:
		return d.FromString(casted)
	default:
		return fmt.Errorf("failed to convert driver value of type %T to Decimal", casted)
	}

	return nil
}

// MarshalJSON converts a Provider value to a JSON value
func (enum Provider) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Provider_name, ProviderMapping, true)), nil
}

// MarshalCSV converts a Provider value to CSV cell value
func (enum Provider) MarshalCSV() (string, error) {
	return utils.MarshalString(enum, Provider_name, ProviderMapping, false), nil
}

// MarshalYAML converts a Provider value to a YAML node value
func (enum Provider) MarshalYAML() (interface{}, error) {
	return utils.MarshalString(enum, Provider_name, ProviderMapping, false), nil
}

// MarshalDynamoDBAttributeValue converts a Provider value to a DynamoDB AttributeValue
func (enum Provider) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{Value: utils.MarshalString(enum, Provider_name, ProviderMapping, false)}, nil
}

// UnmarshalJSON attempts to convert a JSON value to a new Provider value
func (enum *Provider) UnmarshalJSON(raw []byte) error {
	return utils.UnmarshalValue(raw, Provider_value, ProviderAlternates, enum)
}

// UnmarshalCSV attempts to convert a CSV cell value to a new Provider value
func (enum *Provider) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Provider_value, ProviderAlternates, enum)
}

// UnmarshalYAML attempts to convert a YAML node to a new Provider value
func (enum *Provider) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.ScalarNode {
		return fmt.Errorf("YAML node had an invalid kind (expected scalar value)")
	} else {
		return utils.UnmarshalString(value.Value, Provider_value, ProviderAlternates, enum)
	}
}

// UnmarshalDynamoDBAttributeValue attempts to convert a DynamoDB AttributeVAlue to a Provider
// value. This function can handle []bytes, numerics, or strings. If the AttributeValue is NULL then
// the Provider value will not be modified.
func (enum *Provider) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Provider_value, ProviderAlternates, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Provider_value, ProviderAlternates, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Provider_value, ProviderAlternates, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Provider", value)
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
