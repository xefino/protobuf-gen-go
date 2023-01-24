package gopb

import (
	"database/sql/driver"
	"fmt"
	"strconv"

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

// AssetClassAlternates contains alternative values for the Financial.Common.AssetClass enum
var AssetClassAlternates = map[string]Financial_Common_AssetClass{
	"stocks":           Financial_Common_Stock,
	"options":          Financial_Common_Option,
	"crypto":           Financial_Common_Crypto,
	"fx":               Financial_Common_ForeignExchange,
	"Foreign Exchange": Financial_Common_ForeignExchange,
	"otc":              Financial_Common_OverTheCounter,
	"OTC":              Financial_Common_OverTheCounter,
}

// AssetClassMapping contains alternate names for the Financial.Common.AssetClass enum
var AssetClassMapping = map[Financial_Common_AssetClass]string{
	Financial_Common_ForeignExchange: "Foreign Exchange",
	Financial_Common_OverTheCounter:  "OTC",
}

// AssetTypeAlternates contains alternative values for the Financial.Common.AssetType enum
var AssetTypeAlternates = map[string]Financial_Common_AssetType{
	"CS":                      Financial_Common_CommonShare,
	"Common Share":            Financial_Common_CommonShare,
	"OS":                      Financial_Common_OrdinaryShare,
	"Ordinary Share":          Financial_Common_OrdinaryShare,
	"NYRS":                    Financial_Common_NewYorkRegistryShares,
	"New York Registry Share": Financial_Common_NewYorkRegistryShares,
	"ADRC":                    Financial_Common_AmericanDepositoryReceiptCommon,
	"Common ADR":              Financial_Common_AmericanDepositoryReceiptCommon,
	"ADRP":                    Financial_Common_AmericanDepositoryReceiptPreferred,
	"Preferred ADR":           Financial_Common_AmericanDepositoryReceiptPreferred,
	"ADRR":                    Financial_Common_AmericanDepositoryReceiptRights,
	"ADR Right":               Financial_Common_AmericanDepositoryReceiptRights,
	"ADRW":                    Financial_Common_AmericanDepositoryReceiptWarrants,
	"ADR Warrant":             Financial_Common_AmericanDepositoryReceiptWarrants,
	"GDR":                     Financial_Common_GlobalDepositoryReceipts,
	"UNIT":                    Financial_Common_Unit,
	"RIGHT":                   Financial_Common_Rights,
	"Right":                   Financial_Common_Rights,
	"PFD":                     Financial_Common_PreferredStock,
	"Preferred Stock":         Financial_Common_PreferredStock,
	"FUND":                    Financial_Common_Fund,
	"SP":                      Financial_Common_StructuredProduct,
	"Structured Product":      Financial_Common_StructuredProduct,
	"WARRANT":                 Financial_Common_Warrant,
	"INDEX":                   Financial_Common_Index,
	"ETF":                     Financial_Common_ExchangeTradedFund,
	"ETN":                     Financial_Common_ExchangeTradedNote,
	"BOND":                    Financial_Common_CorporateBond,
	"Corporate Bond":          Financial_Common_CorporateBond,
	"AGEN":                    Financial_Common_AgencyBond,
	"Agency Bond":             Financial_Common_AgencyBond,
	"EQLK":                    Financial_Common_EquityLinkedBond,
	"Equity-Linked Bond":      Financial_Common_EquityLinkedBond,
	"BASKET":                  Financial_Common_Basket,
	"LT":                      Financial_Common_LiquidatingTrust,
	"Liquidating Trust":       Financial_Common_LiquidatingTrust,
	"OTHER":                   Financial_Common_Others,
	"Other":                   Financial_Common_Others,
	"":                        Financial_Common_None,
}

// AssetTypeMapping contains alternate names for the Financial.Common.AssetType enum
var AssetTypeMapping = map[Financial_Common_AssetType]string{
	Financial_Common_CommonShare:                        "Common Share",
	Financial_Common_OrdinaryShare:                      "Ordinary Share",
	Financial_Common_NewYorkRegistryShares:              "New York Registry Share",
	Financial_Common_AmericanDepositoryReceiptCommon:    "Common ADR",
	Financial_Common_AmericanDepositoryReceiptPreferred: "Preferred ADR",
	Financial_Common_AmericanDepositoryReceiptRights:    "ADR Right",
	Financial_Common_AmericanDepositoryReceiptWarrants:  "ADR Warrant",
	Financial_Common_GlobalDepositoryReceipts:           "GDR",
	Financial_Common_Rights:                             "Right",
	Financial_Common_PreferredStock:                     "Preferred Stock",
	Financial_Common_StructuredProduct:                  "Structured Product",
	Financial_Common_ExchangeTradedFund:                 "ETF",
	Financial_Common_ExchangeTradedNote:                 "ETN",
	Financial_Common_CorporateBond:                      "Corporate Bond",
	Financial_Common_AgencyBond:                         "Agency Bond",
	Financial_Common_EquityLinkedBond:                   "Equity-Linked Bond",
	Financial_Common_LiquidatingTrust:                   "Liquidating Trust",
	Financial_Common_Others:                             "Other",
	Financial_Common_None:                               "",
}

// LocalAlternates contains alternative values for the Financial.Common.Locale enum
var LocaleAlternates = map[string]Financial_Common_Locale{
	"us":     Financial_Common_US,
	"global": Financial_Common_Global,
}

// OptionContractTypeAlternates contains alternative values for the Financial.Options.ContractType enum
var OptionContractTypeAlternates = map[string]Financial_Options_ContractType{
	"call":  Financial_Options_Call,
	"put":   Financial_Options_Put,
	"other": Financial_Options_Other,
}

// OptionExerciseStyleAlternates contains alternative values for the Financial.Options.ExerciseStyle enum
var OptionExerciseStyleAlternates = map[string]Financial_Options_ExerciseStyle{
	"american": Financial_Options_American,
	"european": Financial_Options_European,
	"bermudan": Financial_Options_Bermudan,
}

// UnderlyingTypeAlternates contains alternative values for the Financial.Options.UnderlyingType enum
var OptionUnderlyingTypeAlternates = map[string]Financial_Options_UnderlyingType{
	"equity":   Financial_Options_Equity,
	"currency": Financial_Options_Currency,
}

// ExchangeTypeAlternates contains alternative values for the Financial.Exchanges.Type enum
var ExchangeTypeAlternates = map[string]Financial_Exchanges_Type{
	"exchange": Financial_Exchanges_Exchange,
}

// TradeCorrectionAlternates contains alternative valus for the Financial.Trades.CorrectionCode enum
var TradeCorrectionAlternates = map[string]Financial_Trades_CorrectionCode{
	"00": Financial_Trades_NotCorrected,
	"01": Financial_Trades_LateCorrected,
	"07": Financial_Trades_Erroneous,
	"08": Financial_Trades_Cancel,
}

// QuoteConditionAlternates contains alternative values for the Financial.Quotes.Condition enum
var QuoteConditionAlternates = map[string]Financial_Quotes_Condition{
	"-1": Financial_Quotes_Invalid,
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
	switch casted := value.(type) {
	case string:
		return timestamp.FromString(casted)
	case int64:
		timestamp.Seconds = casted / nanosPerSecond
		timestamp.Nanoseconds = int32(casted % nanosPerSecond)
		return nil
	default:
		return fmt.Errorf("Value of %v with a type of %T could not be converted to a UnixTimestamp", casted, casted)
	}
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

// MarhsalJSON converts a Financial.Common.AssetClass to JSON
func (enum Financial_Common_AssetClass) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Common_AssetClass_name, AssetClassMapping, true)), nil
}

// MarshalCSV converts a Financial.Common.AssetClass to a CSV format
func (enum Financial_Common_AssetClass) MarshalCSV() (string, error) {
	return strconv.FormatInt(int64(enum), 10), nil
}

// Marshaler converts a Financial.Common.AssetClass to a DynamoDB attribute value
func (enum Financial_Common_AssetClass) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Common_AssetClass_name, AssetClassMapping, false),
	}, nil
}

// UnmarshalJSON converts JSON data into a Financial.Common.AssetClass
func (enum *Financial_Common_AssetClass) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Common_AssetClass_value, AssetClassAlternates, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Common.AssetClass
func (enum *Financial_Common_AssetClass) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Common_AssetClass_value, AssetClassAlternates, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Common.AssetClass
func (enum *Financial_Common_AssetClass) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Common_AssetClass_value, AssetClassAlternates, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Common_AssetClass_value, AssetClassAlternates, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Common_AssetClass_value, AssetClassAlternates, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Common.AssetClass", value)
	}
}

// Scan converts an SQL value into a Financial.Common.AssetClass
func (enum *Financial_Common_AssetClass) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Common_AssetClass_value, AssetClassAlternates, enum)
}

// MarhsalJSON converts a Financial.Common.AssetType to JSON
func (enum Financial_Common_AssetType) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Common_AssetType_name, AssetTypeMapping, true)), nil
}

// MarshalCSV converts a Financial.Common.AssetType to a CSV format
func (enum Financial_Common_AssetType) MarshalCSV() (string, error) {
	return strconv.FormatInt(int64(enum), 10), nil
}

// Marshaler converts a Financial.Common.AssetType to a DynamoDB attribute value
func (enum Financial_Common_AssetType) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Common_AssetType_name, AssetTypeMapping, false),
	}, nil
}

// UnmarshalJSON converts JSON data into a Financial.Common.AssetType
func (enum *Financial_Common_AssetType) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Common_AssetType_value, AssetTypeAlternates, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Common.AssetType
func (enum *Financial_Common_AssetType) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Common_AssetType_value, AssetTypeAlternates, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Common.AssetType
func (enum *Financial_Common_AssetType) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Common_AssetType_value, AssetTypeAlternates, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Common_AssetType_value, AssetTypeAlternates, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Common_AssetType_value, AssetTypeAlternates, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Common.AssetType", value)
	}
}

// Scan converts an SQL value into a Financial.Common.AssetType
func (enum *Financial_Common_AssetType) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Common_AssetType_value, AssetTypeAlternates, enum)
}

// MarhsalJSON converts a Financial.Common.Locale to JSON
func (enum Financial_Common_Locale) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Common_Locale_name, utils.Ignore, true)), nil
}

// MarshalCSV converts a Financial.Common.Locale to a CSV format
func (enum Financial_Common_Locale) MarshalCSV() (string, error) {
	return strconv.FormatInt(int64(enum), 10), nil
}

// Marshaler converts a Financial.Common.Locale to a DynamoDB attribute value
func (enum Financial_Common_Locale) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Common_Locale_name, utils.Ignore, false),
	}, nil
}

// UnmarshalJSON converts JSON data into a Financial.Common.Locale
func (enum *Financial_Common_Locale) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Common_Locale_value, LocaleAlternates, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Common.Locale
func (enum *Financial_Common_Locale) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Common_Locale_value, LocaleAlternates, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Common.Locale
func (enum *Financial_Common_Locale) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Common_Locale_value, LocaleAlternates, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Common_Locale_value, LocaleAlternates, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Common_Locale_value, LocaleAlternates, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Common.Locale", value)
	}
}

// Scan converts an SQL value into a Financial.Common.Locale
func (enum *Financial_Common_Locale) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Common_Locale_value, LocaleAlternates, enum)
}

// MarhsalJSON converts a Financial.Common.Tape to JSON
func (enum Financial_Common_Tape) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Common_Tape_name, utils.Ignore, true)), nil
}

// MarshalCSV converts a Financial.Common.Tape to a CSV format
func (enum Financial_Common_Tape) MarshalCSV() (string, error) {
	return utils.MarshalString(enum, Financial_Common_Tape_name, utils.Ignore, false), nil
}

// Marshaler converts a Financial.Common.Tape to a DynamoDB attribute value
func (enum Financial_Common_Tape) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Common_Tape_name, utils.Ignore, false),
	}, nil
}

// Value converts a Financial.Common.Tape to an SQL value
func (enum Financial_Common_Tape) Value() (driver.Value, error) {
	return driver.Value(utils.MarshalString(enum, Financial_Common_Tape_name, utils.Ignore, false)), nil
}

// UnmarshalJSON converts JSON data into a Financial.Common.Tape
func (enum *Financial_Common_Tape) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Common_Tape_value, utils.None, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Common.Tape
func (enum *Financial_Common_Tape) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Common_Tape_value, utils.None, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Common.Tape
func (enum *Financial_Common_Tape) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Common_Tape_value, utils.None, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Common_Tape_value, utils.None, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Common_Tape_value, utils.None, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Common.Tape", value)
	}
}

// Scan converts an SQL value into a Financial.Common.Tape
func (enum *Financial_Common_Tape) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Common_Tape_value, utils.None, enum)
}

// MarhsalJSON converts a Financial.Dividends.Frequency to JSON
func (enum Financial_Dividends_Frequency) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Dividends_Frequency_name, utils.Ignore, true)), nil
}

// MarshalCSV converts a Financial.Dividends.Frequency to a CSV format
func (enum Financial_Dividends_Frequency) MarshalCSV() (string, error) {
	return strconv.FormatInt(int64(enum), 10), nil
}

// Marshaler converts a Financial.Dividends.Frequency to a DynamoDB attribute value
func (enum Financial_Dividends_Frequency) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Dividends_Frequency_name, utils.Ignore, false),
	}, nil
}

// UnmarshalJSON converts JSON data into a Financial.Dividends.Frequency
func (enum *Financial_Dividends_Frequency) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Dividends_Frequency_value, utils.None, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Dividends.Frequency
func (enum *Financial_Dividends_Frequency) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Dividends_Frequency_value, utils.None, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Dividends.Frequency
func (enum *Financial_Dividends_Frequency) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Dividends_Frequency_value, utils.None, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Dividends_Frequency_value, utils.None, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Dividends_Frequency_value, utils.None, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Dividends.Frequency", value)
	}
}

// Scan converts an SQL value into a Financial.Dividends.Frequency
func (enum *Financial_Dividends_Frequency) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Dividends_Frequency_value, utils.None, enum)
}

// MarhsalJSON converts a Financial.Dividends.Type to JSON
func (enum Financial_Dividends_Type) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Dividends_Type_name, utils.Ignore, true)), nil
}

// MarshalCSV converts a Financial.Dividends.Type to a CSV format
func (enum Financial_Dividends_Type) MarshalCSV() (string, error) {
	return utils.MarshalString(enum, Financial_Dividends_Type_name, utils.Ignore, false), nil
}

// Marshaler converts a Financial.Dividends.Type to a DynamoDB attribute value
func (enum Financial_Dividends_Type) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Dividends_Type_name, utils.Ignore, false),
	}, nil
}

// Value converts a Financial.Dividends.Type to an SQL value
func (enum Financial_Dividends_Type) Value() (driver.Value, error) {
	return driver.Value(utils.MarshalString(enum, Financial_Dividends_Type_name, utils.Ignore, false)), nil
}

// UnmarshalJSON converts JSON data into a Financial.Dividends.Type
func (enum *Financial_Dividends_Type) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Dividends_Type_value, utils.None, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Dividends.Type
func (enum *Financial_Dividends_Type) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Dividends_Type_value, utils.None, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Dividends.Type
func (enum *Financial_Dividends_Type) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Dividends_Type_value, utils.None, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Dividends_Type_value, utils.None, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Dividends_Type_value, utils.None, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Dividends.Type", value)
	}
}

// Scan converts an SQL value into a Financial.Dividends.Type
func (enum *Financial_Dividends_Type) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Dividends_Type_value, utils.None, enum)
}

// MarhsalJSON converts a Financial.Exchanges.Type to JSON
func (enum Financial_Exchanges_Type) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Exchanges_Type_name, utils.Ignore, true)), nil
}

// MarshalCSV converts a Financial.Exchanges.Type to a CSV format
func (enum Financial_Exchanges_Type) MarshalCSV() (string, error) {
	return strconv.FormatInt(int64(enum), 10), nil
}

// Marshaler converts a Financial.Exchanges.Type to a DynamoDB attribute value
func (enum Financial_Exchanges_Type) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Exchanges_Type_name, utils.Ignore, false),
	}, nil
}

// UnmarshalJSON converts JSON data into a Financial.Exchanges.Type
func (enum *Financial_Exchanges_Type) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Exchanges_Type_value, ExchangeTypeAlternates, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Exchanges.Type
func (enum *Financial_Exchanges_Type) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Exchanges_Type_value, ExchangeTypeAlternates, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Exchanges.Type
func (enum *Financial_Exchanges_Type) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Exchanges_Type_value, ExchangeTypeAlternates, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Exchanges_Type_value, ExchangeTypeAlternates, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Exchanges_Type_value, ExchangeTypeAlternates, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Exchanges.Type", value)
	}
}

// Scan converts an SQL value into a Financial.Exchanges.Type
func (enum *Financial_Exchanges_Type) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Exchanges_Type_value, ExchangeTypeAlternates, enum)
}

// MarhsalJSON converts a Financial.Options.ContractType to JSON
func (enum Financial_Options_ContractType) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Options_ContractType_name, utils.Ignore, true)), nil
}

// MarshalCSV converts a Financial.Options.ContractType to a CSV format
func (enum Financial_Options_ContractType) MarshalCSV() (string, error) {
	return strconv.FormatInt(int64(enum), 10), nil
}

// Marshaler converts a Financial.Options.ContractType to a DynamoDB attribute value
func (enum Financial_Options_ContractType) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Options_ContractType_name, utils.Ignore, false),
	}, nil
}

// Value converts a Financial.Options.ContractType to an SQL value
func (enum Financial_Options_ContractType) Value() (driver.Value, error) {
	return driver.Value(utils.MarshalString(enum, Financial_Options_ContractType_name, utils.Ignore, false)), nil
}

// UnmarshalJSON converts JSON data into a Financial.Options.ContractType
func (enum *Financial_Options_ContractType) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Options_ContractType_value, OptionContractTypeAlternates, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Options.ContractType
func (enum *Financial_Options_ContractType) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Options_ContractType_value, OptionContractTypeAlternates, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Options.ContractType
func (enum *Financial_Options_ContractType) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Options_ContractType_value, OptionContractTypeAlternates, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Options_ContractType_value, OptionContractTypeAlternates, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Options_ContractType_value, OptionContractTypeAlternates, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Options.ContractType", value)
	}
}

// Scan converts an SQL value into a Financial.Options.ContractType
func (enum *Financial_Options_ContractType) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Options_ContractType_value, OptionContractTypeAlternates, enum)
}

// MarhsalJSON converts a Financial.Options.ExerciseStyle to JSON
func (enum Financial_Options_ExerciseStyle) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Options_ExerciseStyle_name, utils.Ignore, true)), nil
}

// MarshalCSV converts a Financial.Options.ExerciseStyle to a CSV format
func (enum Financial_Options_ExerciseStyle) MarshalCSV() (string, error) {
	return strconv.FormatInt(int64(enum), 10), nil
}

// Marshaler converts a Financial.Options.ExerciseStyle to a DynamoDB attribute value
func (enum Financial_Options_ExerciseStyle) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Options_ExerciseStyle_name, utils.Ignore, false),
	}, nil
}

// Value converts a Financial.Options.ExerciseStyle to an SQL value
func (enum Financial_Options_ExerciseStyle) Value() (driver.Value, error) {
	return driver.Value(utils.MarshalString(enum, Financial_Options_ExerciseStyle_name, utils.Ignore, false)), nil
}

// UnmarshalJSON converts JSON data into a Financial.Options.ExerciseStyle
func (enum *Financial_Options_ExerciseStyle) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Options_ExerciseStyle_value, OptionExerciseStyleAlternates, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Options.ExerciseStyle
func (enum *Financial_Options_ExerciseStyle) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Options_ExerciseStyle_value, OptionExerciseStyleAlternates, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Options.ExerciseStyle
func (enum *Financial_Options_ExerciseStyle) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Options_ExerciseStyle_value, OptionExerciseStyleAlternates, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Options_ExerciseStyle_value, OptionExerciseStyleAlternates, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Options_ExerciseStyle_value, OptionExerciseStyleAlternates, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Options.ExerciseStyle", value)
	}
}

// Scan converts an SQL value into a Financial.Options.ExerciseStyle
func (enum *Financial_Options_ExerciseStyle) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Options_ExerciseStyle_value, OptionExerciseStyleAlternates, enum)
}

// MarhsalJSON converts a Financial.Options.UnderlyingType to JSON
func (enum Financial_Options_UnderlyingType) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Options_UnderlyingType_name, utils.Ignore, true)), nil
}

// MarshalCSV converts a Financial.Options.UnderlyingType to a CSV format
func (enum Financial_Options_UnderlyingType) MarshalCSV() (string, error) {
	return strconv.FormatInt(int64(enum), 10), nil
}

// Marshaler converts a Financial.Options.UnderlyingType to a DynamoDB attribute value
func (enum Financial_Options_UnderlyingType) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Options_UnderlyingType_name, utils.Ignore, false),
	}, nil
}

// Value converts a Financial.Options.UnderlyingType to an SQL value
func (enum Financial_Options_UnderlyingType) Value() (driver.Value, error) {
	return driver.Value(utils.MarshalString(enum, Financial_Options_UnderlyingType_name, utils.Ignore, false)), nil
}

// UnmarshalJSON converts JSON data into a Financial.Options.UnderlyingType
func (enum *Financial_Options_UnderlyingType) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Options_UnderlyingType_value, OptionUnderlyingTypeAlternates, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Options.UnderlyingType
func (enum *Financial_Options_UnderlyingType) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Options_UnderlyingType_value, OptionUnderlyingTypeAlternates, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Options.UnderlyingType
func (enum *Financial_Options_UnderlyingType) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Options_UnderlyingType_value, OptionUnderlyingTypeAlternates, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Options_UnderlyingType_value, OptionUnderlyingTypeAlternates, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Options_UnderlyingType_value, OptionUnderlyingTypeAlternates, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Options.UnderlyingType", value)
	}
}

// Scan converts an SQL value into a Financial.Options.UnderlyingType
func (enum *Financial_Options_UnderlyingType) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Options_UnderlyingType_value, OptionUnderlyingTypeAlternates, enum)
}

// MarhsalJSON converts a Financial.Quotes.Condition to JSON
func (enum Financial_Quotes_Condition) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Quotes_Condition_name, utils.Ignore, true)), nil
}

// MarshalCSV converts a Financial.Quotes.Condition to a CSV format
func (enum Financial_Quotes_Condition) MarshalCSV() (string, error) {
	return fmt.Sprintf("%03d", enum), nil
}

// Marshaler converts a Financial.Quotes.Condition to a DynamoDB attribute value
func (enum Financial_Quotes_Condition) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Quotes_Condition_name, utils.Ignore, false),
	}, nil
}

// Value converts a Financial.Quotes.Condition to an SQL value
func (enum Financial_Quotes_Condition) Value() (driver.Value, error) {

	// If we have an invalid value then return the actual value for invalid
	if enum == Financial_Quotes_Invalid {
		return driver.Value(-1), nil
	}

	// Otherwise, let the driver use the integer value of the enum
	return driver.Value(int(enum)), nil
}

// UnmarshalJSON converts JSON data into a Financial.Quotes.Condition
func (enum *Financial_Quotes_Condition) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Quotes_Condition_value, QuoteConditionAlternates, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Quotes.Condition
func (enum *Financial_Quotes_Condition) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Quotes_Condition_value, QuoteConditionAlternates, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Quotes.Condition
func (enum *Financial_Quotes_Condition) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Quotes_Condition_value, QuoteConditionAlternates, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Quotes_Condition_value, QuoteConditionAlternates, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Quotes_Condition_value, QuoteConditionAlternates, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Quotes.Condition", value)
	}
}

// Scan converts an SQL value into a Financial.Quotes.Condition
func (enum *Financial_Quotes_Condition) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Quotes_Condition_value, QuoteConditionAlternates, enum)
}

// MarhsalJSON converts a Financial.Quotes.Indicator to JSON
func (enum Financial_Quotes_Indicator) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Quotes_Indicator_name, utils.Ignore, true)), nil
}

// MarshalCSV converts a Financial.Quotes.Indicator to a CSV format
func (enum Financial_Quotes_Indicator) MarshalCSV() (string, error) {
	return strconv.FormatInt(int64(enum), 10), nil
}

// Marshaler converts a Financial.Quotes.Indicator to a DynamoDB attribute value
func (enum Financial_Quotes_Indicator) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Quotes_Indicator_name, utils.Ignore, false),
	}, nil
}

// UnmarshalJSON converts JSON data into a Financial.Quotes.Indicator
func (enum *Financial_Quotes_Indicator) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Quotes_Indicator_value, utils.None, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Quotes.Indicator
func (enum *Financial_Quotes_Indicator) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Quotes_Indicator_value, utils.None, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Quotes.Indicator
func (enum *Financial_Quotes_Indicator) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Quotes_Indicator_value, utils.None, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Quotes_Indicator_value, utils.None, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Quotes_Indicator_value, utils.None, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Quotes.Indicator", value)
	}
}

// Scan converts an SQL value into a Financial.Quotes.Indicator
func (enum *Financial_Quotes_Indicator) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Quotes_Indicator_value, utils.None, enum)
}

// MarhsalJSON converts a Financial.Trades.Condition to JSON
func (enum Financial_Trades_Condition) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Trades_Condition_name, utils.Ignore, true)), nil
}

// MarshalCSV converts a Financial.Trades.Condition to a CSV format
func (enum Financial_Trades_Condition) MarshalCSV() (string, error) {
	return strconv.FormatInt(int64(enum), 10), nil
}

// Marshaler converts a Financial.Trades.Condition to a DynamoDB attribute value
func (enum Financial_Trades_Condition) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Trades_Condition_name, utils.Ignore, false),
	}, nil
}

// UnmarshalJSON converts JSON data into a Financial.Trades.Condition
func (enum *Financial_Trades_Condition) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Trades_Condition_value, utils.None, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Trades.Condition
func (enum *Financial_Trades_Condition) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Trades_Condition_value, utils.None, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Trades.Condition
func (enum *Financial_Trades_Condition) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Trades_Condition_value, utils.None, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Trades_Condition_value, utils.None, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Trades_Condition_value, utils.None, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Trades.Condition", value)
	}
}

// Scan converts an SQL value into a Financial.Trades.Condition
func (enum *Financial_Trades_Condition) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Trades_Condition_value, utils.None, enum)
}

// MarhsalJSON converts a Financial.Trades.CorrectionCode to JSON
func (enum Financial_Trades_CorrectionCode) MarshalJSON() ([]byte, error) {
	return []byte(utils.MarshalString(enum, Financial_Trades_CorrectionCode_name, utils.Ignore, true)), nil
}

// MarshalCSV converts a Financial.Trades.CorrectionCode to a CSV format
func (enum Financial_Trades_CorrectionCode) MarshalCSV() (string, error) {
	return fmt.Sprintf("%02d", enum), nil
}

// Marshaler converts a Financial.Trades.CorrectionCode to a DynamoDB attribute value
func (enum Financial_Trades_CorrectionCode) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: utils.MarshalString(enum, Financial_Trades_CorrectionCode_name, utils.Ignore, false),
	}, nil
}

// UnmarshalJSON converts JSON data into a Financial.Trades.CorrectionCode
func (enum *Financial_Trades_CorrectionCode) UnmarshalJSON(data []byte) error {
	return utils.UnmarshalValue(data, Financial_Trades_CorrectionCode_value, TradeCorrectionAlternates, enum)
}

// UnmarshalCSV converts a CSV column into a Financial.Trades.CorrectionCode
func (enum *Financial_Trades_CorrectionCode) UnmarshalCSV(raw string) error {
	return utils.UnmarshalString(raw, Financial_Trades_CorrectionCode_value, TradeCorrectionAlternates, enum)
}

// UnmarshalDynamoDBAttributeValue converts a DynamoDB attribute value to a Financial.Trades.CorrectionCode
func (enum *Financial_Trades_CorrectionCode) UnmarshalDynamoDBAttributeValue(value types.AttributeValue) error {
	switch casted := value.(type) {
	case *types.AttributeValueMemberB:
		return utils.UnmarshalValue(casted.Value, Financial_Trades_CorrectionCode_value, TradeCorrectionAlternates, enum)
	case *types.AttributeValueMemberN:
		return utils.UnmarshalString(casted.Value, Financial_Trades_CorrectionCode_value, TradeCorrectionAlternates, enum)
	case *types.AttributeValueMemberNULL:
		return nil
	case *types.AttributeValueMemberS:
		return utils.UnmarshalString(casted.Value, Financial_Trades_CorrectionCode_value, TradeCorrectionAlternates, enum)
	default:
		return fmt.Errorf("Attribute value of %T could not be converted to a Financial.Trades.CorrectionCode", value)
	}
}

// Scan converts an SQL value into a Financial.Trades.CorrectionCode
func (enum *Financial_Trades_CorrectionCode) Scan(value interface{}) error {
	return utils.ScanValue(value, Financial_Trades_CorrectionCode_value, TradeCorrectionAlternates, enum)
}
