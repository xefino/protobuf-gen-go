package gopb

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
)

var _ = Describe("Decimal Marshal/Unmarshal Tests", func() {

	// Test that converting the Decimal to JSON works for all values
	DescribeTable("MarshalJSON Tests",
		func(decimal *Decimal, expected string) {
			actual, err := json.Marshal(decimal)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(actual)).Should(Equal(expected))
		},
		Entry("Value is positive - Works",
			&Decimal{Exp: -3, Parts: []int64{234109887750000111, 6554423}}, "6554423234109887750000.111"),
		Entry("Value is zero - Works", &Decimal{Parts: []int64{0}}, "0"),
		Entry("Value less than 0 - Works",
			&Decimal{Parts: []int64{-351234088800000999, -342645987}, Exp: -5}, "-3426459873512340888000.00999"))

	// Test that converting the Decimal to a CSV column works for all values
	DescribeTable("MarshalCSV Tests",
		func(decimal *Decimal, expected string) {
			actual, err := decimal.MarshalCSV()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(actual).Should(Equal(expected))
		},
		Entry("Value is positive - Works",
			&Decimal{Exp: -3, Parts: []int64{234109887750000111, 6554423}}, "6554423234109887750000.111"),
		Entry("Value is zero - Works", &Decimal{Parts: []int64{0}}, "0"),
		Entry("Value less than 0 - Works",
			&Decimal{Parts: []int64{-351234088800000999, -342645987}, Exp: -5}, "-3426459873512340888000.00999"))

	// Test that converting the Decimal to a DynamoDB AttributeVAlue works for all values
	DescribeTable("MarshalDynamoDBAttributeValue - Works",
		func(decimal *Decimal, expected string) {
			data, err := attributevalue.Marshal(decimal)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data.(*types.AttributeValueMemberN).Value).Should(Equal(expected))
		},
		Entry("Value is positive - Works",
			&Decimal{Exp: -3, Parts: []int64{234109887750000111, 6554423}}, "6554423234109887750000.111"),
		Entry("Value is zero - Works", &Decimal{Parts: []int64{0}}, "0"),
		Entry("Value less than 0 - Works",
			&Decimal{Parts: []int64{-351234088800000999, -342645987}, Exp: -5}, "-3426459873512340888000.00999"))

	// Test that converting the Decimal to an sql.Value works for all values
	DescribeTable("Value - Works",
		func(decimal *Decimal, expected string) {
			actual, err := decimal.Value()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(actual).Should(Equal(expected))
		},
		Entry("Value is positive - Works",
			&Decimal{Exp: -3, Parts: []int64{234109887750000111, 6554423}}, "6554423234109887750000.111"),
		Entry("Value is zero - Works", &Decimal{Parts: []int64{0}}, "0"),
		Entry("Value less than 0 - Works",
			&Decimal{Parts: []int64{-351234088800000999, -342645987}, Exp: -5}, "-3426459873512340888000.00999"))

	// Test that attempting to deserialize a Decimal will fail and return an error if the value
	// cannot be deserialized from a JSON value
	It("UnmarshalJSON fails - Error", func() {

		// Attempt to convert a non-parseable string value into a Decimal; this should return an error
		value := new(Decimal)
		err := value.UnmarshalJSON([]byte("derp"))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("can't convert derp to decimal: exponent is not numeric"))
	})

	// Test the conditions under which values should be convertible to a Decimal
	DescribeTable("UnmarshalJSON Tests",
		func(raw string, verifier func(*Decimal)) {

			// Attempt to convert the string value into a Decimal; this should not fail
			value := new(Decimal)
			err := value.UnmarshalJSON([]byte(raw))

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			verifier(value)
		},
		Entry("Value greater than 0 - Works", "1234512351234088800000.999",
			decimalVerifier(-3, 351234088800000999, 1234512)),
		Entry("Value equal to 0 - Works", "0", decimalVerifier(0)),
		Entry("Value less than 0 - Works", "-288341660781234512351234088800000.999",
			decimalVerifier(-3, -351234088800000999, -288341660781234512)))

	// Test that attempting to deserialize a Decimal will fail and return an error if the value
	// cannot be converted to either the name value or integer value of the enum option
	It("UnmarshalCSV - Value is invalid - Error", func() {

		// Attempt to convert a fake string value into a Decimal; this should return an error
		value := new(Decimal)
		err := value.UnmarshalCSV("derp")

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("can't convert derp to decimal: exponent is not numeric"))
	})

	// Test the conditions under which values should be convertible to a Decimal
	DescribeTable("UnmarshalCSV Tests",
		func(raw string, verifier func(*Decimal)) {

			// Attempt to convert the value into a Decimal; this should not fail
			value := new(Decimal)
			err := value.UnmarshalCSV(raw)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			verifier(value)
		},
		Entry("Value greater than 0 - Works", "1234512351234088800000.999",
			decimalVerifier(-3, 351234088800000999, 1234512)),
		Entry("Value equal to 0 - Works", "0", decimalVerifier(0)),
		Entry("Value less than 0 - Works", "-288341660781234512351234088800000.999",
			decimalVerifier(-3, -351234088800000999, -288341660781234512)))

	// Tests that, if the attribute type submitted to UnmarshalDynamoDBAttributeValue is not one we
	// recognize, then the function will return an error
	It("UnmarshalDynamoDBAttributeValue - AttributeValue type invalid - Error", func() {
		value := new(Decimal)
		err := attributevalue.Unmarshal(&types.AttributeValueMemberBOOL{Value: true}, &value)
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("Attribute value of *types.AttributeValueMemberBOOL could not be converted to a Decimal"))
	})

	// Tests that, if time parsing fails, then calling UnmarshalDynamoDBAttributeValue will return an error
	It("UnmarshalDynamoDBAttributeValue - Parse fails - Error", func() {
		value := new(Decimal)
		err := attributevalue.Unmarshal(&types.AttributeValueMemberS{Value: "derp"}, &value)
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("can't convert derp to decimal: exponent is not numeric"))
	})

	// Tests the conditions under which UnmarshalDynamoDBAttributeValue is called and no error is generated
	DescribeTable("UnmarshalDynamoDBAttributeValue - AttributeValue Conditions",
		func(raw types.AttributeValue, verifier func(*Decimal)) {
			var value *Decimal
			err := attributevalue.Unmarshal(raw, &value)
			Expect(err).ShouldNot(HaveOccurred())
			verifier(value)
		},
		Entry("Value is []bytes, Value greater than 0 - Works",
			&types.AttributeValueMemberB{Value: []byte("1234512351234088800000.999")},
			decimalVerifier(-3, 351234088800000999, 1234512)),
		Entry("Value is []bytes, Value equal to 0 - Works",
			&types.AttributeValueMemberB{Value: []byte("0")}, decimalVerifier(0)),
		Entry("Value is []bytes, Value less than 0 - Works",
			&types.AttributeValueMemberB{Value: []byte("-288341660781234512351234088800000.999")},
			decimalVerifier(-3, -351234088800000999, -288341660781234512)),
		Entry("Value is numeric, Value greater than 0 - Works",
			&types.AttributeValueMemberN{Value: "1234512351234088800000.999"},
			decimalVerifier(-3, 351234088800000999, 1234512)),
		Entry("Value is numeric, Value equal to 0 - Works",
			&types.AttributeValueMemberN{Value: "0"}, decimalVerifier(0)),
		Entry("Value is numeric, Value less than 0 - Works",
			&types.AttributeValueMemberN{Value: "-288341660781234512351234088800000.999"},
			decimalVerifier(-3, -351234088800000999, -288341660781234512)),
		Entry("Value is NULL - Works", new(types.AttributeValueMemberNULL),
			func(d *Decimal) { Expect(d).Should(BeNil()) }),
		Entry("Value is string, Value greater than 0 - Works",
			&types.AttributeValueMemberS{Value: "1234512351234088800000.999"},
			decimalVerifier(-3, 351234088800000999, 1234512)),
		Entry("Value is string, Value equal to 0 - Works",
			&types.AttributeValueMemberS{Value: "0"}, decimalVerifier(0)),
		Entry("Value is string, Value less than 0 - Works",
			&types.AttributeValueMemberS{Value: "-288341660781234512351234088800000.999"},
			decimalVerifier(-3, -351234088800000999, -288341660781234512)))

	// Tests that, if the type of the driver value is not one we can work with, then Scan will return an error
	It("Scan - Type is invalid - Error", func() {

		// Attempt to convert a fake string value into a Decimal; this should return an error
		value := new(Decimal)
		err := value.Scan(true)

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("failed to convert driver value of type bool to Decimal"))
	})

	// Tests that, if the value is invalid, then Scan will return an error
	It("Scan - Value is invalid - Error", func() {

		// Attempt to convert a fake string value into a Decimal; this should return an error
		value := new(Decimal)
		err := value.Scan("derp")

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("can't convert derp to decimal: exponent is not numeric"))
	})

	// Tests the conditions under which Scan is called and no error is generated
	DescribeTable("Scan Tests",
		func(raw interface{}, verifier func(*Decimal)) {
			value := new(Decimal)
			err := value.Scan(raw)
			Expect(err).ShouldNot(HaveOccurred())
			verifier(value)
		},
		Entry("Value is []byte, Value greater than 0 - Works",
			[]byte("1234512351234088800000.999"), decimalVerifier(-3, 351234088800000999, 1234512)),
		Entry("Value is []byte, Value equal to 0 - Works",
			[]byte("0"), decimalVerifier(0)),
		Entry("Value is []byte, Value less than 0 - Works", []byte("-288341660781234512351234088800000.999"),
			decimalVerifier(-3, -351234088800000999, -288341660781234512)),
		Entry("Value is float, Value greater than 0 - Works",
			1234512351234.999, decimalVerifier(-3, 1234512351234999)),
		Entry("Value is float, Value equal to 0 - Works", 0.0, decimalVerifier(0)),
		Entry("Value is float, Value less than 0 - Works", -28834166.999, decimalVerifier(-3, -28834166999)),
		Entry("Value is int, Value greater than 0 - Works",
			int64(1234512351234088800), decimalVerifier(0, 234512351234088800, 1)),
		Entry("Value is int, Value equal to 0 - Works", int64(0), decimalVerifier(0)),
		Entry("Value is int, Value less than 0 - Works", int64(-2883416607812345123),
			decimalVerifier(0, -883416607812345123, -2)),
		Entry("Value is string, Value greater than 0 - Works", "1234512351234088800000.999",
			decimalVerifier(-3, 351234088800000999, 1234512)),
		Entry("Value is string, Value equal to 0 - Works", "0", decimalVerifier(0)),
		Entry("Value is string, Value less than 0 - Works", "-288341660781234512351234088800000.999",
			decimalVerifier(-3, -351234088800000999, -288341660781234512)))
})

var _ = Describe("Provider Marshal/Unmarshal Tests", func() {

	// Test that converting the Provider enum to JSON works for all values
	DescribeTable("MarshalJSON Tests",
		func(enum Provider, value string) {
			data, err := json.Marshal(enum)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("None - Works", Provider_None, "\"\""),
		Entry("Polygon - Works", Provider_Polygon, "\"polygon\""))

	// Test that converting the Provider enum to a CSV column works for all values
	DescribeTable("MarshalCSV Tests",
		func(enum Provider, value string) {
			data, err := enum.MarshalCSV()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("None - Works", Provider_None, ""),
		Entry("Polygon - Works", Provider_Polygon, "polygon"))

	// Test that converting the Provider enum to a YAML node works for all values
	DescribeTable("MarshalYAML - Works",
		func(enum Provider, value string) {
			data, err := enum.MarshalYAML()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data).Should(Equal(value))
		},
		Entry("None - Works", Provider_None, ""),
		Entry("Polygon - Works", Provider_Polygon, "polygon"))

	// Test that converting the Provider enum to a DynamoDB AttributeVAlue works for all values
	DescribeTable("MarshalDynamoDBAttributeValue - Works",
		func(enum Provider, value string) {
			data, err := attributevalue.Marshal(enum)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data.(*types.AttributeValueMemberS).Value).Should(Equal(value))
		},
		Entry("None - Works", Provider_None, ""),
		Entry("Polygon - Works", Provider_Polygon, "polygon"))

	// Test that attempting to deserialize a Provider will fail and return an error if the value
	// cannot be deserialized from a JSON value to a string
	It("UnmarshalJSON fails - Error", func() {

		// Attempt to convert a non-parseable string value into a Provider; this should return an error
		enum := new(Provider)
		err := enum.UnmarshalJSON([]byte("derp"))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Provider"))
	})

	// Test that attempting to deserialize a Provider will fail and return an error if the value
	// cannot be converted to either the name value or integer value of the enum option
	It("UnmarshalJSON - Value is invalid - Error", func() {

		// Attempt to convert a fake string value into a Provider; this should return an error
		enum := new(Provider)
		err := enum.UnmarshalJSON([]byte("\"derp\""))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Provider"))
	})

	// Test the conditions under which values should be convertible to a Provider
	DescribeTable("UnmarshalJSON Tests",
		func(value string, shouldBe Provider) {

			// Attempt to convert the string value into a Provider; this should not fail
			var enum Provider
			err := enum.UnmarshalJSON([]byte(value))

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("Empty String - Works", "\"\"", Provider_None),
		Entry("Polygon - Works", "\"polygon\"", Provider_Polygon),
		Entry("None - Works", "\"None\"", Provider_None),
		Entry("Polygon - Works", "\"Polygon\"", Provider_Polygon),
		Entry("0 - Works", "\"0\"", Provider_None),
		Entry("1 - Works", "\"1\"", Provider_Polygon))

	// Test that attempting to deserialize a Provider will fail and return an error if the value
	// cannot be converted to either the name value or integer value of the enum option
	It("UnmarshalCSV - Value is invalid - Error", func() {

		// Attempt to convert a fake string value into a Provider; this should return an error
		enum := new(Provider)
		err := enum.UnmarshalCSV("derp")

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Provider"))
	})

	// Test the conditions under which values should be convertible to a Provider
	DescribeTable("UnmarshalCSV Tests",
		func(value string, shouldBe Provider) {

			// Attempt to convert the value into a Provider; this should not fail
			var enum Provider
			err := enum.UnmarshalCSV(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("Empty String - Works", "", Provider_None),
		Entry("Polygon - Works", "polygon", Provider_Polygon),
		Entry("None - Works", "None", Provider_None),
		Entry("Polygon - Works", "Polygon", Provider_Polygon),
		Entry("0 - Works", "0", Provider_None),
		Entry("1 - Works", "1", Provider_Polygon))

	// Test that attempting to deserialize a Provider will fail and return an error if the YAML
	// node does not represent a scalar value
	It("UnmarshalYAML - Node type is not scalar - Error", func() {
		enum := new(Provider)
		err := enum.UnmarshalYAML(&yaml.Node{Kind: yaml.AliasNode})
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("YAML node had an invalid kind (expected scalar value)"))
	})

	// Test that attempting to deserialize a Provider will fail and return an error if the YAML
	// node value cannot be converted to either the name value or integer value of the enum option
	It("UnmarshalYAML - Parse fails - Error", func() {
		enum := new(Provider)
		err := enum.UnmarshalYAML(&yaml.Node{Kind: yaml.ScalarNode, Value: "derp"})
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Provider"))
	})

	// Test the conditions under which YAML node values should be convertible to a Provider
	DescribeTable("UnmarshalYAML Tests",
		func(value string, shouldBe Provider) {
			var enum Provider
			err := enum.UnmarshalYAML(&yaml.Node{Kind: yaml.ScalarNode, Value: value})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("Empty String - Works", "", Provider_None),
		Entry("Polygon - Works", "polygon", Provider_Polygon),
		Entry("None - Works", "None", Provider_None),
		Entry("Polygon - Works", "Polygon", Provider_Polygon),
		Entry("0 - Works", "0", Provider_None),
		Entry("1 - Works", "1", Provider_Polygon))

	// Tests that, if the attribute type submitted to UnmarshalDynamoDBAttributeValue is not one we
	// recognize, then the function will return an error
	It("UnmarshalDynamoDBAttributeValue - AttributeValue type invalid - Error", func() {
		enum := new(Provider)
		err := attributevalue.Unmarshal(&types.AttributeValueMemberBOOL{Value: true}, &enum)
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("Attribute value of *types.AttributeValueMemberBOOL could not be converted to a Provider"))
	})

	// Tests that, if time parsing fails, then calling UnmarshalDynamoDBAttributeValue will return an error
	It("UnmarshalDynamoDBAttributeValue - Parse fails - Error", func() {
		enum := new(Provider)
		err := attributevalue.Unmarshal(&types.AttributeValueMemberS{Value: "derp"}, &enum)
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Provider"))
	})

	// Tests the conditions under which UnmarshalDynamoDBAttributeValue is called and no error is generated
	DescribeTable("UnmarshalDynamoDBAttributeValue - AttributeValue Conditions",
		func(value types.AttributeValue, expected Provider) {
			var enum Provider
			err := attributevalue.Unmarshal(value, &enum)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(expected))
		},
		Entry("Value is []bytes, Empty String - Works",
			&types.AttributeValueMemberB{Value: []byte("")}, Provider_None),
		Entry("Value is []bytes, polygon - Works",
			&types.AttributeValueMemberB{Value: []byte("polygon")}, Provider_Polygon),
		Entry("Value is []bytes, None - Works",
			&types.AttributeValueMemberB{Value: []byte("None")}, Provider_None),
		Entry("Value is []bytes, Polygon - Works",
			&types.AttributeValueMemberB{Value: []byte("Polygon")}, Provider_Polygon),
		Entry("Value is []bytes, 0 - Works",
			&types.AttributeValueMemberB{Value: []byte("0")}, Provider_None),
		Entry("Value is []bytes, 1 - Works",
			&types.AttributeValueMemberB{Value: []byte("1")}, Provider_Polygon),
		Entry("Value is int, 0 - Works",
			&types.AttributeValueMemberN{Value: "0"}, Provider_None),
		Entry("Value is int, 1 - Works",
			&types.AttributeValueMemberN{Value: "1"}, Provider_Polygon),
		Entry("Value is NULL - Works", new(types.AttributeValueMemberNULL), Provider(0)),
		Entry("Value is string, Empty String - Works",
			&types.AttributeValueMemberS{Value: ""}, Provider_None),
		Entry("Value is string, polygon - Works",
			&types.AttributeValueMemberS{Value: "polygon"}, Provider_Polygon),
		Entry("Value is string, None - Works",
			&types.AttributeValueMemberS{Value: "None"}, Provider_None),
		Entry("Value is string, Polygon - Works",
			&types.AttributeValueMemberS{Value: "Polygon"}, Provider_Polygon),
		Entry("Value is string, 0 - Works",
			&types.AttributeValueMemberS{Value: "0"}, Provider_None),
		Entry("Value is string, 1 - Works",
			&types.AttributeValueMemberS{Value: "1"}, Provider_Polygon))
})

var _ = Describe("UnixTimestamp Marshal/Unmarshal Tests", func() {

	// Test that converting a Timestamp to JSON works for all values
	DescribeTable("MarshalJSON Tests",
		func(timestamp *UnixTimestamp, value string) {
			data, err := timestamp.MarshalJSON()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("Timestamp is nil - Works", nil, ""),
		Entry("Timestamp has value - Works",
			&UnixTimestamp{Seconds: 1654127993, Nanoseconds: 983651350}, "1654127993983651350"))

	// Test that converting a Timestamp to a CSV column works for all values
	DescribeTable("MarshalCSV Tests",
		func(timestamp *UnixTimestamp, value string) {
			data, err := timestamp.MarshalCSV()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("Timestamp is nil - Works", nil, ""),
		Entry("Timestamp has value - Works",
			&UnixTimestamp{Seconds: 1654127993, Nanoseconds: 983651350}, "1654127993983651350"))

	// Test that converting a Timestamp to a AttributeValue works for all values
	DescribeTable("MarshalDynamoDBAttributeValue Tests",
		func(timestamp *UnixTimestamp, value string) {
			data, err := timestamp.MarshalDynamoDBAttributeValue()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data.(*types.AttributeValueMemberS).Value).Should(Equal(value))
		},
		Entry("Timestamp is nil - Works", nil, ""),
		Entry("Timestamp has value - Works",
			&UnixTimestamp{Seconds: 1654127993, Nanoseconds: 983651350}, "1654127993983651350"))

	// Test that converting a Timestamp to an SQL value for all values
	DescribeTable("Value Tests",
		func(timestamp *UnixTimestamp, value string) {
			data, err := timestamp.Value()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data).Should(Equal(value))
		},
		Entry("Timestamp is nil - Works", nil, ""),
		Entry("Timestamp has value - Works",
			&UnixTimestamp{Seconds: 1654127993, Nanoseconds: 983651350}, "1654127993983651350"))

	// Test that attempting to deserialize a Timestamp will fail and return an error if the
	// value canno be deserialized from a JSON value to a string
	DescribeTable("UnmarshalJSON - Failures",
		func(rawValue string, callDirectly bool, message string) {

			// Attempt to convert a non-parseable string value into a timestamp
			// This should return an error
			var timestamp *UnixTimestamp
			var err error
			if callDirectly {
				err = timestamp.UnmarshalJSON([]byte(rawValue))
				Expect(timestamp).Should(BeNil())
			} else {
				err = json.Unmarshal([]byte(rawValue), &timestamp)
			}

			// Verify the error
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(Equal(message))
		},
		Entry("String is too short - Error", "derp", true,
			"value (derp) was not long enough to be converted to a timestamp"),
		Entry("Seconds cannot be converted to an integer - Error", "derp983651350", true,
			"failed to convert seconds part to integer, error: strconv.ParseInt: "+
				"parsing \"derp\": invalid syntax"),
		Entry("Nanoseconds cannot be converted to an integer - Error", "165412799398365135j", true,
			"failed to convert nanoseconds part to integer, error: strconv.ParseInt: "+
				"parsing \"98365135j\": invalid syntax"),
		Entry("Seconds < Minimum Timestamp - Error", "-62135596801983651350", false,
			"timestamp (-62135596801, 983651350) before 0001-01-01"),
		Entry("Seconds > Maximum Timestamp - Error", "253402300800983651350", false,
			"timestamp (253402300800, 983651350) after 9999-12-31"))

	// Test that, if UnmarshalJSON is called with a value of nil then the timestamp will be nil
	It("UnmarshalJSON - Nil - Nil", func() {

		// Attempt to convert a non-parseable string value into a timestamp
		// This should not return an error
		var timestamp *UnixTimestamp
		err := timestamp.UnmarshalJSON(nil)
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the timestamp
		Expect(timestamp).Should(BeNil())
	})

	// Test that, if UnmarshalJSON is called with an empty string then the timestamp will be nil
	It("UnmarshalJSON - Empty string - Nil", func() {

		// Attempt to convert a non-parseable string value into a timestamp
		// This should not return an error
		var timestamp *UnixTimestamp
		err := timestamp.UnmarshalJSON([]byte(""))
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the timestamp
		Expect(timestamp).Should(BeNil())
	})

	// Test that, if the UnmarshalJSON function is called with a valid UNIX timestamp, then it
	// will be parsed into a Timestamp object
	It("UnmarshalJSON - Non-empty string - Works", func() {

		// Attempt to convert a non-parseable string value into a timestamp
		// This should not return an error
		var timestamp *UnixTimestamp
		err := json.Unmarshal([]byte("1654127993983651350"), &timestamp)
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the timestamp
		Expect(timestamp).ShouldNot(BeNil())
		Expect(timestamp.Seconds).Should(Equal(int64(1654127993)))
		Expect(timestamp.Nanoseconds).Should(Equal(int32(983651350)))
	})

	// Test that attempting to deserialize a Timestamp will fail and return an error if the
	// value canno be deserialized from a CSV column to a string
	DescribeTable("UnmarshalCSV - Failures",
		func(rawValue string, message string) {

			// Attempt to convert a non-parseable string value into a timestamp
			// This should return an error
			timestamp := new(UnixTimestamp)
			err := timestamp.UnmarshalCSV(rawValue)

			// Verify the error
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(Equal(message))
		},
		Entry("String is too short - Error", "derp",
			"value (derp) was not long enough to be converted to a timestamp"),
		Entry("Seconds cannot be converted to an integer - Error", "derp983651350",
			"failed to convert seconds part to integer, error: strconv.ParseInt: "+
				"parsing \"derp\": invalid syntax"),
		Entry("Nanoseconds cannot be converted to an integer - Error", "165412799398365135j",
			"failed to convert nanoseconds part to integer, error: strconv.ParseInt: "+
				"parsing \"98365135j\": invalid syntax"),
		Entry("Seconds < Minimum Timestamp - Error", "-62135596801983651350",
			"timestamp (-62135596801, 983651350) before 0001-01-01"),
		Entry("Seconds > Maximum Timestamp - Error", "253402300800983651350",
			"timestamp (253402300800, 983651350) after 9999-12-31"))

	// Test that, if UnmarshalCSV is called with an empty string then the timestamp will be nil
	It("UnmarshalCSV - Empty string - Nil", func() {

		// Attempt to convert a non-parseable string value into a timestamp
		// This should not return an error
		var timestamp *UnixTimestamp
		err := timestamp.UnmarshalCSV("")
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the timestamp
		Expect(timestamp).Should(BeNil())
	})

	// Test that, if the UnmarshalCSV function is called with a valid UNIX timestamp, then it
	// will be parsed into a Timestamp object
	It("UnmarshalCSV - Non-empty string - Works", func() {

		// Attempt to convert a non-parseable string value into a timestamp
		// This should not return an error
		timestamp := new(UnixTimestamp)
		err := timestamp.UnmarshalCSV("1654127993983651350")
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the timestamp
		Expect(timestamp).ShouldNot(BeNil())
		Expect(timestamp.Seconds).Should(Equal(int64(1654127993)))
		Expect(timestamp.Nanoseconds).Should(Equal(int32(983651350)))
	})

	// Tests that, if the UnmarshalDynamoDBAttributeValue function is called with an invalid AttributeValue
	// type, then the function will return an error
	It("UnmarshalDynamoDBAttributeValue - Type invalid - Error", func() {
		var timestamp *UnixTimestamp
		err := timestamp.UnmarshalDynamoDBAttributeValue(&types.AttributeValueMemberBOOL{Value: *aws.Bool(false)})
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("Attribute value of *types.AttributeValueMemberBOOL could not be converted to a UnixTimestamp"))
	})

	// Tests that, if UnmarshalDynamoDBAttributeValue is called with a AttributeValueMemberNULL,
	// then the timestamp will not be modified and instead will be returned as nil
	It("UnmarshalDynamoDBAttributeValue - Value is NULL - Works", func() {
		var timestamp *UnixTimestamp
		err := timestamp.UnmarshalDynamoDBAttributeValue(&types.AttributeValueMemberNULL{})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(timestamp).Should(BeNil())
	})

	// Tests that the UnmarshalDynamoDBAttributeValue works with various AttributeValue types
	DescribeTable("UnmarshalDynamoDBAttributeValue - Conditions",
		func(attr types.AttributeValue) {
			timestamp := new(UnixTimestamp)
			err := timestamp.UnmarshalDynamoDBAttributeValue(attr)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(timestamp.Seconds).Should(Equal(int64(1654127993)))
			Expect(timestamp.Nanoseconds).Should(Equal(int32(983651350)))
		},
		Entry("Value is []byte - Works",
			&types.AttributeValueMemberB{Value: []byte("1654127993983651350")}),
		Entry("Value is number - Works",
			&types.AttributeValueMemberN{Value: "1654127993983651350"}),
		Entry("Value is string - Works",
			&types.AttributeValueMemberS{Value: "1654127993983651350"}))

	// Test that attempting to deserialize a Timestamp will fail and return an error if the
	// value canno be deserialized from a driver value to a string
	DescribeTable("Scan - Failures",
		func(rawValue interface{}, message string) {

			// Attempt to convert a fake string value into a Timestamp
			// This should return an error
			timestamp := new(UnixTimestamp)
			err := timestamp.Scan(rawValue)

			// Verify the error
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(Equal(message))
		},
		Entry("Type is invalid - Error", true,
			"Value of true with a type of bool could not be converted to a UnixTimestamp"),
		Entry("String is too short - Error", "derp",
			"value (derp) was not long enough to be converted to a timestamp"),
		Entry("Seconds cannot be converted to an integer - Error", "derp983651350",
			"failed to convert seconds part to integer, error: strconv.ParseInt: parsing \"derp\": invalid syntax"),
		Entry("Nanoseconds cannot be converted to an integer - Error", "165412799398365135j",
			"failed to convert nanoseconds part to integer, error: strconv.ParseInt: "+
				"parsing \"98365135j\": invalid syntax"),
		Entry("Seconds < Minimum Timestamp - Error", "-62135596801983651350",
			"timestamp (-62135596801, 983651350) before 0001-01-01"),
		Entry("Seconds > Maximum Timestamp - Error", "253402300800983651350",
			"timestamp (253402300800, 983651350) after 9999-12-31"),
		Entry("Nanoseconds > 1 second - Error", "1654127993-10000000",
			"timestamp (1654127993, -10000000) has out-of-range nanos"))

	// Test that, if Scan is called with a value of nil then the timestamp will be nil
	It("Scan - Nil - Nil", func() {

		// Attempt to convert nil string value into a timestamp
		// This should not return an error
		var timestamp *UnixTimestamp
		err := timestamp.Scan(nil)
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the timestamp
		Expect(timestamp).Should(BeNil())
	})

	// Test that, if Scan is called with an empty string then the timestamp will be nil
	It("Scan - Empty string - Nil", func() {

		// Attempt to convert an empty string value into a timestamp
		// This should not return an error
		var timestamp *UnixTimestamp
		err := timestamp.Scan("")
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the timestamp
		Expect(timestamp).Should(BeNil())
	})

	// Test that, if the Scan function is called with a valid UNIX timestamp, then it
	// will be parsed into a Timestamp object
	It("Scan - Non-empty string - Works", func() {

		// Attempt to convert a UNIX timestamp string value into a timestamp
		// This should not return an error
		timestamp := new(UnixTimestamp)
		err := timestamp.Scan("1654127993983651350")
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the timestamp
		Expect(timestamp).ShouldNot(BeNil())
		Expect(timestamp.Seconds).Should(Equal(int64(1654127993)))
		Expect(timestamp.Nanoseconds).Should(Equal(int32(983651350)))
	})

	// Test that, if the value submitted to the Scan function represents a valid UNIX timestamp,
	// then it will be parsed into a Timestamp object
	It("Scan - Value is int64 - Works", func() {

		// Attempt to convert a UNIX timestamp string value into a timestamp
		// This should not return an error
		timestamp := new(UnixTimestamp)
		err := timestamp.Scan(int64(1654127993983651350))
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the timestamp
		Expect(timestamp).ShouldNot(BeNil())
		Expect(timestamp.Seconds).Should(Equal(int64(1654127993)))
		Expect(timestamp.Nanoseconds).Should(Equal(int32(983651350)))
	})
})

var _ = Describe("UnixDuration Marshal/Unmarshal Tests", func() {

	// Test that converting a UnixDuration to JSON works for all values
	DescribeTable("MarshalJSON Tests",
		func(duration *UnixDuration, value string) {
			data, err := duration.MarshalJSON()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("Duration is nil - Works", nil, ""),
		Entry("Duration has value - Works",
			&UnixDuration{Seconds: 1654127993, Nanoseconds: 983651350}, "1654127993983651350"),
		Entry("Duration was negative - Works",
			&UnixDuration{Seconds: -1654127993, Nanoseconds: -983651350}, "-1654127993983651350"))

	// Test that converting a UnixDuration to a CSV column works for all values
	DescribeTable("MarshalCSV Tests",
		func(duration *UnixDuration, value string) {
			data, err := duration.MarshalCSV()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("Duration is nil - Works", nil, ""),
		Entry("Duration has value - Works",
			&UnixDuration{Seconds: 1654127993, Nanoseconds: 983651350}, "1654127993983651350"),
		Entry("Duration was negative - Works",
			&UnixDuration{Seconds: -1654127993, Nanoseconds: -983651350}, "-1654127993983651350"))

	// Test that converting a UnixDuration to a AttributeValue works for all values
	DescribeTable("MarshalDynamoDBAttributeValue Tests",
		func(duration *UnixDuration, value string) {
			data, err := duration.MarshalDynamoDBAttributeValue()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data.(*types.AttributeValueMemberS).Value).Should(Equal(value))
		},
		Entry("Duration is nil - Works", nil, ""),
		Entry("Duration has value - Works",
			&UnixDuration{Seconds: 1654127993, Nanoseconds: 983651350}, "1654127993983651350"),
		Entry("Duration was negative - Works",
			&UnixDuration{Seconds: -1654127993, Nanoseconds: -983651350}, "-1654127993983651350"))

	// Test that converting a UnixDuration to an SQL value for all values
	DescribeTable("Value Tests",
		func(duration *UnixDuration, value string) {
			data, err := duration.Value()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data).Should(Equal(value))
		},
		Entry("Duration is nil - Works", nil, ""),
		Entry("Duration has value - Works",
			&UnixDuration{Seconds: 1654127993, Nanoseconds: 983651350}, "1654127993983651350"),
		Entry("Duration was negative - Works",
			&UnixDuration{Seconds: -1654127993, Nanoseconds: -983651350}, "-1654127993983651350"))

	// Test that attempting to deserialize a UnixDuration will fail and return an error if the
	// value canno be deserialized from a JSON value to a string
	DescribeTable("UnmarshalJSON - Failures",
		func(rawValue string, callDirectly bool, message string) {

			// Attempt to convert a non-parseable string value into a duration
			// This should return an error
			var duration *UnixDuration
			var err error
			if callDirectly {
				err = duration.UnmarshalJSON([]byte(rawValue))
				Expect(duration).Should(BeNil())
			} else {
				err = json.Unmarshal([]byte(rawValue), &duration)
			}

			// Verify the error
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(Equal(message))
		},
		Entry("String is too short - Error", "derp", true,
			"value (derp) was not long enough to be converted to a duration"),
		Entry("Seconds cannot be converted to an integer - Error", "derp983651350", true,
			"failed to convert seconds part to integer, error: strconv.ParseInt: "+
				"parsing \"derp\": invalid syntax"),
		Entry("Nanoseconds cannot be converted to an integer - Error", "165412799398365135j", true,
			"failed to convert nanoseconds part to integer, error: strconv.ParseInt: "+
				"parsing \"98365135j\": invalid syntax"),
		Entry("Seconds < Minimum Duration - Error", "-315576000001983651350", false,
			"duration (-315576000001, -983651350) exceeds -10000 years"),
		Entry("Seconds > Maximum Duration - Error", "315576000001983651350", false,
			"duration (315576000001, 983651350) exceeds +10000 years"))

	// Test that, if UnmarshalJSON is called with a value of nil then the duration will be nil
	It("UnmarshalJSON - Nil - Nil", func() {

		// Attempt to convert a nil string value into a duration; this should not return an error
		var duration *UnixDuration
		err := duration.UnmarshalJSON(nil)
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the duration
		Expect(duration).Should(BeNil())
	})

	// Test that, if UnmarshalJSON is called with an empty string then the duration will be nil
	It("UnmarshalJSON - Empty string - Nil", func() {

		// Attempt to convert an empty string value into a duration; this should not return an error
		var duration *UnixDuration
		err := duration.UnmarshalJSON([]byte(""))
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the duration
		Expect(duration).Should(BeNil())
	})

	// Test that, if the UnmarshalJSON function is called with a valid UNIX duration, then it
	// will be parsed into a UnixDuration object
	It("UnmarshalJSON - Non-empty string - Works", func() {

		// Attempt to convert a parseable string value into a duration; this should not return an error
		var duration *UnixDuration
		err := json.Unmarshal([]byte("1654127993983651350"), &duration)
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the duration
		Expect(duration).ShouldNot(BeNil())
		Expect(duration.Seconds).Should(Equal(int64(1654127993)))
		Expect(duration.Nanoseconds).Should(Equal(int32(983651350)))
	})

	// Test that, if the UnmarshalJSON function is called with a valid UNIX duration that is negative,
	// then it will be parsed into a UnixDuration object
	It("UnmarshalJSON - Negative duration - Works", func() {

		// Attempt to convert a parseable string value into a duration; this should not return an error
		var duration *UnixDuration
		err := json.Unmarshal([]byte("-1654127993983651350"), &duration)
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the duration
		Expect(duration).ShouldNot(BeNil())
		Expect(duration.Seconds).Should(Equal(int64(-1654127993)))
		Expect(duration.Nanoseconds).Should(Equal(int32(-983651350)))
	})

	// Test that attempting to deserialize a Duration will fail and return an error if the
	// value canno be deserialized from a CSV column to a string
	DescribeTable("UnmarshalCSV - Failures",
		func(rawValue string, message string) {

			// Attempt to convert a non-parseable string value into a duration
			// This should return an error
			duration := new(UnixDuration)
			err := duration.UnmarshalCSV(rawValue)

			// Verify the error
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(Equal(message))
		},
		Entry("String is too short - Error", "derp",
			"value (derp) was not long enough to be converted to a duration"),
		Entry("Seconds cannot be converted to an integer - Error", "derp983651350",
			"failed to convert seconds part to integer, error: strconv.ParseInt: "+
				"parsing \"derp\": invalid syntax"),
		Entry("Nanoseconds cannot be converted to an integer - Error", "165412799398365135j",
			"failed to convert nanoseconds part to integer, error: strconv.ParseInt: "+
				"parsing \"98365135j\": invalid syntax"),
		Entry("Seconds < Minimum Duration - Error", "-315576000001983651350",
			"duration (-315576000001, -983651350) exceeds -10000 years"),
		Entry("Seconds > Maximum Duration - Error", "315576000001983651350",
			"duration (315576000001, 983651350) exceeds +10000 years"))

	// Test that, if UnmarshalCSV is called with an empty string then the duration will be nil
	It("UnmarshalCSV - Empty string - Nil", func() {

		// Attempt to convert an empty string value into a duration; this should not return an error
		var duration *UnixDuration
		err := duration.UnmarshalCSV("")
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the duration
		Expect(duration).Should(BeNil())
	})

	// Test that, if the UnmarshalCSV function is called with a valid UNIX duration, then it
	// will be parsed into a UnixDuration object
	It("UnmarshalCSV - Non-empty string - Works", func() {

		// Attempt to convert a parseable string value into a duration; this should not return an error
		duration := new(UnixDuration)
		err := duration.UnmarshalCSV("1654127993983651350")
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the duration
		Expect(duration).ShouldNot(BeNil())
		Expect(duration.Seconds).Should(Equal(int64(1654127993)))
		Expect(duration.Nanoseconds).Should(Equal(int32(983651350)))
	})

	// Test that, if the UnmarshalCSV function is called with a valid UNIX duration that is negative,
	// then it will be parsed into a UnixDuration object
	It("UnmarshalCSV - Negative duration - Works", func() {

		// Attempt to convert a parseable string value into a duration; this should not return an error
		duration := new(UnixDuration)
		err := duration.UnmarshalCSV("-1654127993983651350")
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the duration
		Expect(duration).ShouldNot(BeNil())
		Expect(duration.Seconds).Should(Equal(int64(-1654127993)))
		Expect(duration.Nanoseconds).Should(Equal(int32(-983651350)))
	})

	// Tests that, if the UnmarshalDynamoDBAttributeValue function is called with an invalid AttributeValue
	// type, then the function will return an error
	It("UnmarshalDynamoDBAttributeValue - Type invalid - Error", func() {
		var duration *UnixDuration
		err := duration.UnmarshalDynamoDBAttributeValue(&types.AttributeValueMemberBOOL{Value: *aws.Bool(false)})
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("Attribute value of *types.AttributeValueMemberBOOL could not be converted to a UnixDuration"))
	})

	// Tests that, if UnmarshalDynamoDBAttributeValue is called with a AttributeValueMemberNULL,
	// then the duration will not be modified and instead will be returned as nil
	It("UnmarshalDynamoDBAttributeValue - Value is NULL - Works", func() {
		var duration *UnixDuration
		err := duration.UnmarshalDynamoDBAttributeValue(&types.AttributeValueMemberNULL{})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(duration).Should(BeNil())
	})

	// Tests that the UnmarshalDynamoDBAttributeValue works with various AttributeValue types
	DescribeTable("UnmarshalDynamoDBAttributeValue - Conditions",
		func(attr types.AttributeValue) {
			duration := new(UnixDuration)
			err := duration.UnmarshalDynamoDBAttributeValue(attr)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(duration.Seconds).Should(Equal(int64(1654127993)))
			Expect(duration.Nanoseconds).Should(Equal(int32(983651350)))
		},
		Entry("Value is []byte - Works",
			&types.AttributeValueMemberB{Value: []byte("1654127993983651350")}),
		Entry("Value is number - Works",
			&types.AttributeValueMemberN{Value: "1654127993983651350"}),
		Entry("Value is string - Works",
			&types.AttributeValueMemberS{Value: "1654127993983651350"}))

	// Test that attempting to deserialize a UnixDuration will fail and return an error if the
	// value canno be deserialized from a driver value to a string
	DescribeTable("Scan - Failures",
		func(rawValue string, message string) {

			// Attempt to convert a fake string value into a Duration
			// This should return an error
			duration := new(UnixDuration)
			err := duration.Scan(rawValue)

			// Verify the error
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(Equal(message))
		},
		Entry("String is too short - Error", "derp",
			"value (derp) was not long enough to be converted to a duration"),
		Entry("Seconds cannot be converted to an integer - Error", "derp983651350",
			"failed to convert seconds part to integer, error: strconv.ParseInt: parsing \"derp\": invalid syntax"),
		Entry("Nanoseconds cannot be converted to an integer - Error", "165412799398365135j",
			"failed to convert nanoseconds part to integer, error: strconv.ParseInt: "+
				"parsing \"98365135j\": invalid syntax"),
		Entry("Seconds < Minimum Duration - Error", "-315576000001983651350",
			"duration (-315576000001, -983651350) exceeds -10000 years"),
		Entry("Seconds > Maximum Duration - Error", "315576000001983651350",
			"duration (315576000001, 983651350) exceeds +10000 years"),
		Entry("Nanoseconds > 1 second - Error", "1654127993-10000000",
			"duration (1654127993, -10000000) has seconds and nanos with different signs"))

	// Test that, if Scan is called with a value of nil then the duration will be nil
	It("Scan - Nil - Nil", func() {

		// Attempt to convert nil string value into a duration; this should not return an error
		var duration *UnixDuration
		err := duration.Scan(nil)
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the duration
		Expect(duration).Should(BeNil())
	})

	// Test that, if Scan is called with an empty string then the duration will be nil
	It("Scan - Empty string - Nil", func() {

		// Attempt to convert an empty string value into a duration; this should not return an error
		var duration *UnixDuration
		err := duration.Scan("")
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the duration
		Expect(duration).Should(BeNil())
	})

	// Test that, if the Scan function is called with a valid UNIX duration, then it
	// will be parsed into a UnixDuration object
	It("Scan - Non-empty string - Works", func() {

		// Attempt to convert a UNIX duration string value into a duration; this should not return an error
		duration := new(UnixDuration)
		err := duration.Scan("1654127993983651350")
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the duration
		Expect(duration).ShouldNot(BeNil())
		Expect(duration.Seconds).Should(Equal(int64(1654127993)))
		Expect(duration.Nanoseconds).Should(Equal(int32(983651350)))
	})

	// Test that, if the Scan function is called with a valid UNIX duration that is negative,
	// then it will be parsed into a UnixDuration object
	It("Scan - Negative duration - Works", func() {

		// Attempt to convert a parseable string value into a duration; this should not return an error
		duration := new(UnixDuration)
		err := duration.Scan("-1654127993983651350")
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the duration
		Expect(duration).ShouldNot(BeNil())
		Expect(duration.Seconds).Should(Equal(int64(-1654127993)))
		Expect(duration.Nanoseconds).Should(Equal(int32(-983651350)))
	})
})

var _ = Describe("Financial.Common.AssetClass Marshal/Unmarshal Tests", func() {

	// Test that converting the Financial.Common.AssetClass enum to JSON works for all values
	DescribeTable("MarshalJSON Tests",
		func(enum Financial_Common_AssetClass, value string) {
			data, err := json.Marshal(enum)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("Stock - Works", Financial_Common_Stock, "\"Stock\""),
		Entry("Option - Works", Financial_Common_Option, "\"Option\""),
		Entry("Crypto - Works", Financial_Common_Crypto, "\"Crypto\""),
		Entry("ForeignExchange - Works", Financial_Common_ForeignExchange, "\"Foreign Exchange\""),
		Entry("OverTheCounter - Works", Financial_Common_OverTheCounter, "\"OTC\""))

	// Test that converting the Financial.Common.AssetClass enum to a CSV column works for all values
	DescribeTable("MarshalCSV Tests",
		func(enum Financial_Common_AssetClass, value string) {
			data, err := enum.MarshalCSV()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("Stocks - Works", Financial_Common_Stock, "0"),
		Entry("Options - Works", Financial_Common_Option, "1"),
		Entry("Crypto - Works", Financial_Common_Crypto, "2"),
		Entry("ForeignExchange - Works", Financial_Common_ForeignExchange, "3"),
		Entry("OverTheCounter - Works", Financial_Common_OverTheCounter, "4"))

	// Test that converting the Financial.Common.AssetClass enum to a AttributeValue works for all values
	DescribeTable("MarshalDynamoDBAttributeValue Tests",
		func(enum Financial_Common_AssetClass, value string) {
			data, err := enum.MarshalDynamoDBAttributeValue()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data.(*types.AttributeValueMemberS).Value).Should(Equal(value))
		},
		Entry("Stock - Works", Financial_Common_Stock, "Stock"),
		Entry("Option - Works", Financial_Common_Option, "Option"),
		Entry("Crypto - Works", Financial_Common_Crypto, "Crypto"),
		Entry("ForeignExchange - Works", Financial_Common_ForeignExchange, "Foreign Exchange"),
		Entry("OverTheCounter - Works", Financial_Common_OverTheCounter, "OTC"))

	// Test that attempting to deserialize a Financial.Common.AssetClass will fail and
	// return an error if the value canno be deserialized from a JSON value to a string
	It("UnmarshalJSON fails - Error", func() {

		// Attempt to convert a non-parseable string value into a Financial.Common.AssetClass
		// This should return an error
		enum := new(Financial_Common_AssetClass)
		err := enum.UnmarshalJSON([]byte("derp"))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Common_AssetClass"))
	})

	// Test that attempting to deserialize a Financial.Common.AssetClass will fail and
	// return an error if the value cannot be converted to either the name value or integer
	// value of the enum option
	It("UnmarshalJSON - Value is invalid - Error", func() {

		// Attempt to convert a fake string value into a Financial.Common.AssetClass
		// This should return an error
		enum := new(Financial_Common_AssetClass)
		err := enum.UnmarshalJSON([]byte("\"derp\""))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Common_AssetClass"))
	})

	// Test the conditions under which values should be convertible to a Financial.Common.AssetClass
	DescribeTable("UnmarshalJSON Tests",
		func(value string, shouldBe Financial_Common_AssetClass) {

			// Attempt to convert the string value into a Financial.Common.AssetClass
			// This should not fail
			var enum Financial_Common_AssetClass
			err := enum.UnmarshalJSON([]byte(value))

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("Stock - Works", "\"Stock\"", Financial_Common_Stock),
		Entry("Option - Works", "\"Option\"", Financial_Common_Option),
		Entry("Crypto - Works", "\"Crypto\"", Financial_Common_Crypto),
		Entry("ForeignExchange - Works", "\"ForeignExchange\"", Financial_Common_ForeignExchange),
		Entry("OverTheCounter - Works", "\"OverTheCounter\"", Financial_Common_OverTheCounter),
		Entry("Foreign Exchange - Works", "\"Foreign Exchange\"", Financial_Common_ForeignExchange),
		Entry("OTC - Works", "\"OTC\"", Financial_Common_OverTheCounter),
		Entry("stocks - Works", "\"stocks\"", Financial_Common_Stock),
		Entry("options - Works", "\"options\"", Financial_Common_Option),
		Entry("crypto - Works", "\"crypto\"", Financial_Common_Crypto),
		Entry("otc - Works", "\"otc\"", Financial_Common_OverTheCounter),
		Entry("fx - Works", "\"fx\"", Financial_Common_ForeignExchange),
		Entry("0 - Works", "\"0\"", Financial_Common_Stock),
		Entry("1 - Works", "\"1\"", Financial_Common_Option),
		Entry("2 - Works", "\"2\"", Financial_Common_Crypto),
		Entry("3 - Works", "\"3\"", Financial_Common_ForeignExchange),
		Entry("4 - Works", "\"4\"", Financial_Common_OverTheCounter))

	// Test that attempting to deserialize a Financial.Common.AssetClass will fial and return an
	// error if the value cannot be converted to either the name value or integer value
	// of the enum option
	It("UnmarshalCSV - Value is empty - Error", func() {

		// Attempt to convert a fake string value into a Financial.Common.AssetClass
		// This should return an error
		enum := new(Financial_Common_AssetClass)
		err := enum.UnmarshalCSV("")

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"\" cannot be mapped to a gopb.Financial_Common_AssetClass"))
	})

	// Test the conditions under which values should be convertible to a Financial.Common.AssetClass
	DescribeTable("UnmarshalCSV Tests",
		func(value string, shouldBe Financial_Common_AssetClass) {

			// Attempt to convert the value into a Financial.Common.AssetClass
			// This should not fail
			var enum Financial_Common_AssetClass
			err := enum.UnmarshalCSV(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("Stock - Works", "Stock", Financial_Common_Stock),
		Entry("Option - Works", "Option", Financial_Common_Option),
		Entry("Crypto - Works", "Crypto", Financial_Common_Crypto),
		Entry("ForeignExchange - Works", "ForeignExchange", Financial_Common_ForeignExchange),
		Entry("OverTheCounter - Works", "OverTheCounter", Financial_Common_OverTheCounter),
		Entry("Foreign Exchange - Works", "Foreign Exchange", Financial_Common_ForeignExchange),
		Entry("OTC - Works", "OTC", Financial_Common_OverTheCounter),
		Entry("stocks - Works", "stocks", Financial_Common_Stock),
		Entry("options - Works", "options", Financial_Common_Option),
		Entry("crypto - Works", "crypto", Financial_Common_Crypto),
		Entry("fx - Works", "fx", Financial_Common_ForeignExchange),
		Entry("otc - Works", "otc", Financial_Common_OverTheCounter),
		Entry("0 - Works", "0", Financial_Common_Stock),
		Entry("1 - Works", "1", Financial_Common_Option),
		Entry("2 - Works", "2", Financial_Common_Crypto),
		Entry("3 - Works", "3", Financial_Common_ForeignExchange),
		Entry("4 - Works", "4", Financial_Common_OverTheCounter))

	// Tests that, if the attribute type submitted to UnmarshalDynamoDBAttributeValue is not one we
	// recognize, then the function will return an error
	It("UnmarshalDynamoDBAttributeValue - AttributeValue type invalid - Error", func() {
		value := new(Financial_Common_AssetClass)
		err := attributevalue.Unmarshal(&types.AttributeValueMemberBOOL{Value: true}, &value)
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("Attribute value of *types.AttributeValueMemberBOOL could not be converted to a Financial.Common.AssetClass"))
	})

	// Tests the conditions under which UnmarshalDynamoDBAttributeValue is called and no error is generated
	DescribeTable("UnmarshalDynamoDBAttributeValue - AttributeValue Conditions",
		func(raw types.AttributeValue, expected Financial_Common_AssetClass) {
			var value Financial_Common_AssetClass
			err := attributevalue.Unmarshal(raw, &value)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(expected))
		},
		Entry("Value is []bytes, stocks - Works",
			&types.AttributeValueMemberB{Value: []byte("stocks")}, Financial_Common_Stock),
		Entry("Value is []bytes, options - Works",
			&types.AttributeValueMemberB{Value: []byte("options")}, Financial_Common_Option),
		Entry("Value is []bytes, crypto - Works",
			&types.AttributeValueMemberB{Value: []byte("crypto")}, Financial_Common_Crypto),
		Entry("Value is []bytes, fx - Works",
			&types.AttributeValueMemberB{Value: []byte("fx")}, Financial_Common_ForeignExchange),
		Entry("Value is []bytes, otc - Works",
			&types.AttributeValueMemberB{Value: []byte("otc")}, Financial_Common_OverTheCounter),
		Entry("Value is []bytes, Foreign Exchange - Works",
			&types.AttributeValueMemberB{Value: []byte("Foreign Exchange")}, Financial_Common_ForeignExchange),
		Entry("Value is []bytes, OTC - Works",
			&types.AttributeValueMemberB{Value: []byte("OTC")}, Financial_Common_OverTheCounter),
		Entry("Value is []bytes, Stock - Works",
			&types.AttributeValueMemberB{Value: []byte("Stock")}, Financial_Common_Stock),
		Entry("Value is []bytes, Option - Works",
			&types.AttributeValueMemberB{Value: []byte("Option")}, Financial_Common_Option),
		Entry("Value is []bytes, Crypto - Works",
			&types.AttributeValueMemberB{Value: []byte("Crypto")}, Financial_Common_Crypto),
		Entry("Value is []bytes, ForeignExchange - Works",
			&types.AttributeValueMemberB{Value: []byte("ForeignExchange")}, Financial_Common_ForeignExchange),
		Entry("Value is []bytes, OverTheCounter - Works",
			&types.AttributeValueMemberB{Value: []byte("OverTheCounter")}, Financial_Common_OverTheCounter),
		Entry("Value is numeric, 0 - Works",
			&types.AttributeValueMemberN{Value: "0"}, Financial_Common_Stock),
		Entry("Value is numeric, 1 - Works",
			&types.AttributeValueMemberN{Value: "1"}, Financial_Common_Option),
		Entry("Value is numeric, 2 - Works",
			&types.AttributeValueMemberN{Value: "2"}, Financial_Common_Crypto),
		Entry("Value is numeric, 3 - Works",
			&types.AttributeValueMemberN{Value: "3"}, Financial_Common_ForeignExchange),
		Entry("Value is numeric, 4 - Works",
			&types.AttributeValueMemberN{Value: "4"}, Financial_Common_OverTheCounter),
		Entry("Value is NULL - Works", new(types.AttributeValueMemberNULL), Financial_Common_AssetClass(0)),
		Entry("Value is string, stocks - Works",
			&types.AttributeValueMemberS{Value: "stocks"}, Financial_Common_Stock),
		Entry("Value is string, options - Works",
			&types.AttributeValueMemberS{Value: "options"}, Financial_Common_Option),
		Entry("Value is string, crypto - Works",
			&types.AttributeValueMemberS{Value: "crypto"}, Financial_Common_Crypto),
		Entry("Value is string, fx - Works",
			&types.AttributeValueMemberS{Value: "fx"}, Financial_Common_ForeignExchange),
		Entry("Value is string, otc - Works",
			&types.AttributeValueMemberS{Value: "otc"}, Financial_Common_OverTheCounter),
		Entry("Value is string, Foreign Exchange - Works",
			&types.AttributeValueMemberS{Value: "Foreign Exchange"}, Financial_Common_ForeignExchange),
		Entry("Value is string, OTC - Works",
			&types.AttributeValueMemberS{Value: "OTC"}, Financial_Common_OverTheCounter),
		Entry("Value is string, Stock - Works",
			&types.AttributeValueMemberS{Value: "Stock"}, Financial_Common_Stock),
		Entry("Value is string, Option - Works",
			&types.AttributeValueMemberS{Value: "Option"}, Financial_Common_Option),
		Entry("Value is string, Crypto - Works",
			&types.AttributeValueMemberS{Value: "Crypto"}, Financial_Common_Crypto),
		Entry("Value is string, ForeignExchange - Works",
			&types.AttributeValueMemberS{Value: "ForeignExchange"}, Financial_Common_ForeignExchange),
		Entry("Value is string, OverTheCounter - Works",
			&types.AttributeValueMemberS{Value: "OverTheCounter"}, Financial_Common_OverTheCounter))

	// Test that attempting to deserialize a Financial.Common.AssetClass will fial and return an
	// error if the value cannot be converted to either the name value or integer value of the enum option
	It("Scan - Value is nil - Error", func() {

		// Attempt to convert a fake string value into a Financial.Common.AssetClass
		// This should return an error
		var enum *Financial_Common_AssetClass
		err := enum.Scan(nil)

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of %!q(<nil>) had an invalid type of <nil>"))
		Expect(enum).Should(BeNil())
	})

	// Test the conditions under which values should be convertible to a Financial.Common.AssetClass
	DescribeTable("Scan Tests",
		func(value interface{}, shouldBe Financial_Common_AssetClass) {

			// Attempt to convert the value into a Financial.Common.AssetClass
			// This should not fail
			var enum Financial_Common_AssetClass
			err := enum.Scan(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("Stock - Works", "Stock", Financial_Common_Stock),
		Entry("Option - Works", "Option", Financial_Common_Option),
		Entry("Crypto - Works", "Crypto", Financial_Common_Crypto),
		Entry("ForeignExchange - Works", "ForeignExchange", Financial_Common_ForeignExchange),
		Entry("OverTheCounter - Works", "OverTheCounter", Financial_Common_OverTheCounter),
		Entry("Foreign Exchange - Works", "Foreign Exchange", Financial_Common_ForeignExchange),
		Entry("OTC - Works", "OTC", Financial_Common_OverTheCounter),
		Entry("stocks - Works", "stocks", Financial_Common_Stock),
		Entry("options - Works", "options", Financial_Common_Option),
		Entry("crypto - Works", "crypto", Financial_Common_Crypto),
		Entry("fx - Works", "fx", Financial_Common_ForeignExchange),
		Entry("otc - Works", "otc", Financial_Common_OverTheCounter),
		Entry("0 - Works", 0, Financial_Common_Stock),
		Entry("1 - Works", 1, Financial_Common_Option),
		Entry("2 - Works", 2, Financial_Common_Crypto),
		Entry("3 - Works", 3, Financial_Common_ForeignExchange),
		Entry("4 - Works", 4, Financial_Common_OverTheCounter))
})

var _ = Describe("Financial.Common.AssetType Marshal/Unmarshal Tests", func() {

	// Test that converting the Financial.Common.AssetType enum to JSON works for all values
	DescribeTable("MarshalJSON Tests",
		func(enum Financial_Common_AssetType, value string) {
			data, err := json.Marshal(enum)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("CommonShare - Works", Financial_Common_CommonShare, "\"Common Share\""),
		Entry("OrdinaryShare - Works", Financial_Common_OrdinaryShare, "\"Ordinary Share\""),
		Entry("NewYorkRegistryShares - Works", Financial_Common_NewYorkRegistryShares, "\"New York Registry Share\""),
		Entry("AmericanDepositoryReceiptCommon - Works",
			Financial_Common_AmericanDepositoryReceiptCommon, "\"Common ADR\""),
		Entry("AmericanDepositoryReceiptPreferred - Works",
			Financial_Common_AmericanDepositoryReceiptPreferred, "\"Preferred ADR\""),
		Entry("AmericanDepositoryReceiptRights - Works",
			Financial_Common_AmericanDepositoryReceiptRights, "\"ADR Right\""),
		Entry("AmericanDepositoryReceiptWarrants - Works",
			Financial_Common_AmericanDepositoryReceiptWarrants, "\"ADR Warrant\""),
		Entry("GlobalDepositoryReceipts - Works", Financial_Common_GlobalDepositoryReceipts, "\"GDR\""),
		Entry("Unit - Works", Financial_Common_Unit, "\"Unit\""),
		Entry("Rights - Works", Financial_Common_Rights, "\"Right\""),
		Entry("PreferredStock - Works", Financial_Common_PreferredStock, "\"Preferred Stock\""),
		Entry("Fund - Works", Financial_Common_Fund, "\"Fund\""),
		Entry("StructuredProduct - Works", Financial_Common_StructuredProduct, "\"Structured Product\""),
		Entry("Warrant - Works", Financial_Common_Warrant, "\"Warrant\""),
		Entry("Index - Works", Financial_Common_Index, "\"Index\""),
		Entry("ExchangeTradedFund - Works", Financial_Common_ExchangeTradedFund, "\"ETF\""),
		Entry("ExchangeTradedNote - Works", Financial_Common_ExchangeTradedNote, "\"ETN\""),
		Entry("CorporateBond - Works", Financial_Common_CorporateBond, "\"Corporate Bond\""),
		Entry("AgencyBond - Works", Financial_Common_AgencyBond, "\"Agency Bond\""),
		Entry("EquityLinkedBond - Works", Financial_Common_EquityLinkedBond, "\"Equity-Linked Bond\""),
		Entry("Basket - Works", Financial_Common_Basket, "\"Basket\""),
		Entry("LiquidatingTrust - Works", Financial_Common_LiquidatingTrust, "\"Liquidating Trust\""),
		Entry("Others - Works", Financial_Common_Others, "\"Other\""),
		Entry("None - Works", Financial_Common_None, "\"\""))

	// Test that converting the Financial.Common.AssetType enum to a CSV column works for all values
	DescribeTable("MarshalCSV Tests",
		func(enum Financial_Common_AssetType, value string) {
			data, err := enum.MarshalCSV()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("CommonShare - Works", Financial_Common_CommonShare, "0"),
		Entry("OrdinaryShare - Works", Financial_Common_OrdinaryShare, "1"),
		Entry("NewYorkRegistryShares - Works", Financial_Common_NewYorkRegistryShares, "2"),
		Entry("AmericanDepositoryReceiptCommon - Works", Financial_Common_AmericanDepositoryReceiptCommon, "3"),
		Entry("AmericanDepositoryReceiptPreferred - Works", Financial_Common_AmericanDepositoryReceiptPreferred, "4"),
		Entry("AmericanDepositoryReceiptRights - Works", Financial_Common_AmericanDepositoryReceiptRights, "5"),
		Entry("AmericanDepositoryReceiptWarrants - Works", Financial_Common_AmericanDepositoryReceiptWarrants, "6"),
		Entry("GlobalDepositoryReceipts - Works", Financial_Common_GlobalDepositoryReceipts, "7"),
		Entry("Unit - Works", Financial_Common_Unit, "8"),
		Entry("Rights - Works", Financial_Common_Rights, "9"),
		Entry("PreferredStock - Works", Financial_Common_PreferredStock, "10"),
		Entry("Fund - Works", Financial_Common_Fund, "11"),
		Entry("StructuredProduct - Works", Financial_Common_StructuredProduct, "12"),
		Entry("Warrant - Works", Financial_Common_Warrant, "13"),
		Entry("Index - Works", Financial_Common_Index, "14"),
		Entry("ExchangeTradedFund - Works", Financial_Common_ExchangeTradedFund, "15"),
		Entry("ExchangeTradedNote - Works", Financial_Common_ExchangeTradedNote, "16"),
		Entry("CorporateBond - Works", Financial_Common_CorporateBond, "17"),
		Entry("AgencyBond - Works", Financial_Common_AgencyBond, "18"),
		Entry("EquityLinkedBond - Works", Financial_Common_EquityLinkedBond, "19"),
		Entry("Basket - Works", Financial_Common_Basket, "20"),
		Entry("LiquidatingTrust - Works", Financial_Common_LiquidatingTrust, "21"),
		Entry("Others - Works", Financial_Common_Others, "22"),
		Entry("None - Works", Financial_Common_None, "23"))

	// Test that converting the Financial.Common.AssetType enum to a AttributeValue works for all values
	DescribeTable("MarshalDynamoDBAttributeValue Tests",
		func(enum Financial_Common_AssetType, value string) {
			data, err := enum.MarshalDynamoDBAttributeValue()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data.(*types.AttributeValueMemberS).Value).Should(Equal(value))
		},
		Entry("CommonShare - Works", Financial_Common_CommonShare, "Common Share"),
		Entry("OrdinaryShare - Works", Financial_Common_OrdinaryShare, "Ordinary Share"),
		Entry("NewYorkRegistryShares - Works", Financial_Common_NewYorkRegistryShares, "New York Registry Share"),
		Entry("AmericanDepositoryReceiptCommon - Works", Financial_Common_AmericanDepositoryReceiptCommon, "Common ADR"),
		Entry("AmericanDepositoryReceiptPreferred - Works",
			Financial_Common_AmericanDepositoryReceiptPreferred, "Preferred ADR"),
		Entry("AmericanDepositoryReceiptRights - Works", Financial_Common_AmericanDepositoryReceiptRights, "ADR Right"),
		Entry("AmericanDepositoryReceiptWarrants - Works", Financial_Common_AmericanDepositoryReceiptWarrants, "ADR Warrant"),
		Entry("GlobalDepositoryReceipts - Works", Financial_Common_GlobalDepositoryReceipts, "GDR"),
		Entry("Unit - Works", Financial_Common_Unit, "Unit"),
		Entry("Rights - Works", Financial_Common_Rights, "Right"),
		Entry("PreferredStock - Works", Financial_Common_PreferredStock, "Preferred Stock"),
		Entry("Fund - Works", Financial_Common_Fund, "Fund"),
		Entry("StructuredProduct - Works", Financial_Common_StructuredProduct, "Structured Product"),
		Entry("Warrant - Works", Financial_Common_Warrant, "Warrant"),
		Entry("Index - Works", Financial_Common_Index, "Index"),
		Entry("ExchangeTradedFund - Works", Financial_Common_ExchangeTradedFund, "ETF"),
		Entry("ExchangeTradedNote - Works", Financial_Common_ExchangeTradedNote, "ETN"),
		Entry("CorporateBond - Works", Financial_Common_CorporateBond, "Corporate Bond"),
		Entry("AgencyBond - Works", Financial_Common_AgencyBond, "Agency Bond"),
		Entry("EquityLinkedBond - Works", Financial_Common_EquityLinkedBond, "Equity-Linked Bond"),
		Entry("Basket - Works", Financial_Common_Basket, "Basket"),
		Entry("LiquidatingTrust - Works", Financial_Common_LiquidatingTrust, "Liquidating Trust"),
		Entry("Others - Works", Financial_Common_Others, "Other"),
		Entry("None - Works", Financial_Common_None, ""))

	// Test that attempting to deserialize a Financial.Common.AssetType will fail and
	// return an error if the value canno be deserialized from a JSON value to a string
	It("UnmarshalJSON fails - Error", func() {

		// Attempt to convert a non-parseable string value into a Financial.Common.AssetType
		// This should return an error
		enum := new(Financial_Common_AssetType)
		err := enum.UnmarshalJSON([]byte("derp"))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Common_AssetType"))
	})

	// Test that attempting to deserialize a Financial.Common.AssetType will fail and
	// return an error if the value cannot be converted to either the name value or integer
	// value of the enum option
	It("UnmarshalJSON - Value is invalid - Error", func() {

		// Attempt to convert a fake string value into a Financial.Common.AssetType
		// This should return an error
		enum := new(Financial_Common_AssetType)
		err := enum.UnmarshalJSON([]byte("\"derp\""))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Common_AssetType"))
	})

	// Test the conditions under which values should be convertible to a Financial.Common.AssetType
	DescribeTable("UnmarshalJSON Tests",
		func(value string, shouldBe Financial_Common_AssetType) {

			// Attempt to convert the string value into a Financial.Common.AssetType
			// This should not fail
			var enum Financial_Common_AssetType
			err := enum.UnmarshalJSON([]byte(value))

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("CS - Works", "\"CS\"", Financial_Common_CommonShare),
		Entry("OS - Works", "\"OS\"", Financial_Common_OrdinaryShare),
		Entry("NYRS - Works", "\"NYRS\"", Financial_Common_NewYorkRegistryShares),
		Entry("ADRC - Works", "\"ADRC\"", Financial_Common_AmericanDepositoryReceiptCommon),
		Entry("ADRP - Works", "\"ADRP\"", Financial_Common_AmericanDepositoryReceiptPreferred),
		Entry("ADRR - Works", "\"ADRR\"", Financial_Common_AmericanDepositoryReceiptRights),
		Entry("ADRW - Works", "\"ADRW\"", Financial_Common_AmericanDepositoryReceiptWarrants),
		Entry("GDR - Works", "\"GDR\"", Financial_Common_GlobalDepositoryReceipts),
		Entry("UNIT - Works", "\"UNIT\"", Financial_Common_Unit),
		Entry("RIGHT - Works", "\"RIGHT\"", Financial_Common_Rights),
		Entry("PFD - Works", "\"PFD\"", Financial_Common_PreferredStock),
		Entry("FUND - Works", "\"FUND\"", Financial_Common_Fund),
		Entry("SP - Works", "\"SP\"", Financial_Common_StructuredProduct),
		Entry("WARRANT - Works", "\"WARRANT\"", Financial_Common_Warrant),
		Entry("INDEX - Works", "\"INDEX\"", Financial_Common_Index),
		Entry("ETF - Works", "\"ETF\"", Financial_Common_ExchangeTradedFund),
		Entry("ETN - Works", "\"ETN\"", Financial_Common_ExchangeTradedNote),
		Entry("BOND - Works", "\"BOND\"", Financial_Common_CorporateBond),
		Entry("AGEN - Works", "\"AGEN\"", Financial_Common_AgencyBond),
		Entry("EQLK - Works", "\"EQLK\"", Financial_Common_EquityLinkedBond),
		Entry("BASKET - Works", "\"BASKET\"", Financial_Common_Basket),
		Entry("LT - Works", "\"LT\"", Financial_Common_LiquidatingTrust),
		Entry("OTHER - Works", "\"OTHER\"", Financial_Common_Others),
		Entry("Empty string - Works", "\"\"", Financial_Common_None),
		Entry("Common Share - Works", "\"Common Share\"", Financial_Common_CommonShare),
		Entry("Ordinary Share - Works", "\"Ordinary Share\"", Financial_Common_OrdinaryShare),
		Entry("New York Registry Share - Works", "\"New York Registry Share\"", Financial_Common_NewYorkRegistryShares),
		Entry("Common ADR - Works", "\"Common ADR\"", Financial_Common_AmericanDepositoryReceiptCommon),
		Entry("Preferred ADR - Works", "\"Preferred ADR\"", Financial_Common_AmericanDepositoryReceiptPreferred),
		Entry("ADR Right - Works", "\"ADR Right\"", Financial_Common_AmericanDepositoryReceiptRights),
		Entry("ADR Warrant - Works", "\"ADR Warrant\"", Financial_Common_AmericanDepositoryReceiptWarrants),
		Entry("Right - Works", "\"Right\"", Financial_Common_Rights),
		Entry("Preferred Stock - Works", "\"Preferred Stock\"", Financial_Common_PreferredStock),
		Entry("Structured Product - Works", "\"Structured Product\"", Financial_Common_StructuredProduct),
		Entry("Corporate Bond - Works", "\"Corporate Bond\"", Financial_Common_CorporateBond),
		Entry("Agency Bond - Works", "\"Agency Bond\"", Financial_Common_AgencyBond),
		Entry("Equity-Linked Bond - Works", "\"Equity-Linked Bond\"", Financial_Common_EquityLinkedBond),
		Entry("Liquidating Trust - Works", "\"Liquidating Trust\"", Financial_Common_LiquidatingTrust),
		Entry("Other - Works", "\"Other\"", Financial_Common_Others),
		Entry("CommonShare - Works", "\"CommonShare\"", Financial_Common_CommonShare),
		Entry("OrdinaryShare - Works", "\"OrdinaryShare\"", Financial_Common_OrdinaryShare),
		Entry("NewYorkRegistryShares - Works", "\"NewYorkRegistryShares\"", Financial_Common_NewYorkRegistryShares),
		Entry("AmericanDepositoryReceiptCommon - Works",
			"\"AmericanDepositoryReceiptCommon\"", Financial_Common_AmericanDepositoryReceiptCommon),
		Entry("AmericanDepositoryReceiptPreferred - Works",
			"\"AmericanDepositoryReceiptPreferred\"", Financial_Common_AmericanDepositoryReceiptPreferred),
		Entry("AmericanDepositoryReceiptRights - Works",
			"\"AmericanDepositoryReceiptRights\"", Financial_Common_AmericanDepositoryReceiptRights),
		Entry("AmericanDepositoryReceiptWarrants - Works",
			"\"AmericanDepositoryReceiptWarrants\"", Financial_Common_AmericanDepositoryReceiptWarrants),
		Entry("GlobalDepositoryReceipts - Works", "\"GlobalDepositoryReceipts\"", Financial_Common_GlobalDepositoryReceipts),
		Entry("Unit - Works", "\"Unit\"", Financial_Common_Unit),
		Entry("Rights - Works", "\"Rights\"", Financial_Common_Rights),
		Entry("PreferredStock - Works", "\"PreferredStock\"", Financial_Common_PreferredStock),
		Entry("Fund - Works", "\"Fund\"", Financial_Common_Fund),
		Entry("StructuredProduct - Works", "\"StructuredProduct\"", Financial_Common_StructuredProduct),
		Entry("Warrant - Works", "\"Warrant\"", Financial_Common_Warrant),
		Entry("Index - Works", "\"Index\"", Financial_Common_Index),
		Entry("ExchangeTradedFund - Works", "\"ExchangeTradedFund\"", Financial_Common_ExchangeTradedFund),
		Entry("ExchangeTradedNote - Works", "\"ExchangeTradedNote\"", Financial_Common_ExchangeTradedNote),
		Entry("CorporateBond - Works", "\"CorporateBond\"", Financial_Common_CorporateBond),
		Entry("AgencyBond - Works", "\"AgencyBond\"", Financial_Common_AgencyBond),
		Entry("EquityLinkedBond - Works", "\"EquityLinkedBond\"", Financial_Common_EquityLinkedBond),
		Entry("Basket - Works", "\"Basket\"", Financial_Common_Basket),
		Entry("LiquidatingTrust - Works", "\"LiquidatingTrust\"", Financial_Common_LiquidatingTrust),
		Entry("Others - Works", "\"Others\"", Financial_Common_Others),
		Entry("None - Works", "\"None\"", Financial_Common_None),
		Entry("0 - Works", "\"0\"", Financial_Common_CommonShare),
		Entry("1 - Works", "\"1\"", Financial_Common_OrdinaryShare),
		Entry("2 - Works", "\"2\"", Financial_Common_NewYorkRegistryShares),
		Entry("3 - Works", "\"3\"", Financial_Common_AmericanDepositoryReceiptCommon),
		Entry("4 - Works", "\"4\"", Financial_Common_AmericanDepositoryReceiptPreferred),
		Entry("5 - Works", "\"5\"", Financial_Common_AmericanDepositoryReceiptRights),
		Entry("6 - Works", "\"6\"", Financial_Common_AmericanDepositoryReceiptWarrants),
		Entry("7 - Works", "\"7\"", Financial_Common_GlobalDepositoryReceipts),
		Entry("8 - Works", "\"8\"", Financial_Common_Unit),
		Entry("9 - Works", "\"9\"", Financial_Common_Rights),
		Entry("10 - Works", "\"10\"", Financial_Common_PreferredStock),
		Entry("11 - Works", "\"11\"", Financial_Common_Fund),
		Entry("12 - Works", "\"12\"", Financial_Common_StructuredProduct),
		Entry("13 - Works", "\"13\"", Financial_Common_Warrant),
		Entry("14 - Works", "\"14\"", Financial_Common_Index),
		Entry("15 - Works", "\"15\"", Financial_Common_ExchangeTradedFund),
		Entry("16 - Works", "\"16\"", Financial_Common_ExchangeTradedNote),
		Entry("17 - Works", "\"17\"", Financial_Common_CorporateBond),
		Entry("18 - Works", "\"18\"", Financial_Common_AgencyBond),
		Entry("19 - Works", "\"19\"", Financial_Common_EquityLinkedBond),
		Entry("20 - Works", "\"20\"", Financial_Common_Basket),
		Entry("21 - Works", "\"21\"", Financial_Common_LiquidatingTrust),
		Entry("22 - Works", "\"22\"", Financial_Common_Others),
		Entry("23 - Works", "\"23\"", Financial_Common_None))

	// Test that attempting to deserialize a Financial.Common.AssetType will fial and return an
	// error if the value cannot be converted to either the name value or integer value
	// of the enum option
	It("UnmarshalCSV - Value is invalid - Error", func() {

		// Attempt to convert a fake string value into a Financial.Common.AssetType
		// This should return an error
		enum := new(Financial_Common_AssetType)
		err := enum.UnmarshalCSV("derp")

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Common_AssetType"))
	})

	// Test the conditions under which values should be convertible to a Financial.Common.AssetType
	DescribeTable("UnmarshalCSV Tests",
		func(value string, shouldBe Financial_Common_AssetType) {

			// Attempt to convert the value into a Financial.Common.AssetType
			// This should not fail
			var enum Financial_Common_AssetType
			err := enum.UnmarshalCSV(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("CS - Works", "CS", Financial_Common_CommonShare),
		Entry("OS - Works", "OS", Financial_Common_OrdinaryShare),
		Entry("NYRS - Works", "NYRS", Financial_Common_NewYorkRegistryShares),
		Entry("ADRC - Works", "ADRC", Financial_Common_AmericanDepositoryReceiptCommon),
		Entry("ADRP - Works", "ADRP", Financial_Common_AmericanDepositoryReceiptPreferred),
		Entry("ADRR - Works", "ADRR", Financial_Common_AmericanDepositoryReceiptRights),
		Entry("ADRW - Works", "ADRW", Financial_Common_AmericanDepositoryReceiptWarrants),
		Entry("GDR - Works", "GDR", Financial_Common_GlobalDepositoryReceipts),
		Entry("UNIT - Works", "UNIT", Financial_Common_Unit),
		Entry("RIGHT - Works", "RIGHT", Financial_Common_Rights),
		Entry("PFD - Works", "PFD", Financial_Common_PreferredStock),
		Entry("FUND - Works", "FUND", Financial_Common_Fund),
		Entry("SP - Works", "SP", Financial_Common_StructuredProduct),
		Entry("WARRANT - Works", "WARRANT", Financial_Common_Warrant),
		Entry("INDEX - Works", "INDEX", Financial_Common_Index),
		Entry("ETF - Works", "ETF", Financial_Common_ExchangeTradedFund),
		Entry("ETN - Works", "ETN", Financial_Common_ExchangeTradedNote),
		Entry("BOND - Works", "BOND", Financial_Common_CorporateBond),
		Entry("AGEN - Works", "AGEN", Financial_Common_AgencyBond),
		Entry("EQLK - Works", "EQLK", Financial_Common_EquityLinkedBond),
		Entry("BASKET - Works", "BASKET", Financial_Common_Basket),
		Entry("LT - Works", "LT", Financial_Common_LiquidatingTrust),
		Entry("OTHER - Works", "OTHER", Financial_Common_Others),
		Entry("Empty string - Works", "", Financial_Common_None),
		Entry("Common Share - Works", "Common Share", Financial_Common_CommonShare),
		Entry("Ordinary Share - Works", "Ordinary Share", Financial_Common_OrdinaryShare),
		Entry("New York Registry Share - Works", "New York Registry Share", Financial_Common_NewYorkRegistryShares),
		Entry("Common ADR - Works", "Common ADR", Financial_Common_AmericanDepositoryReceiptCommon),
		Entry("Preferred ADR - Works", "Preferred ADR", Financial_Common_AmericanDepositoryReceiptPreferred),
		Entry("ADR Right - Works", "ADR Right", Financial_Common_AmericanDepositoryReceiptRights),
		Entry("ADR Warrant - Works", "ADR Warrant", Financial_Common_AmericanDepositoryReceiptWarrants),
		Entry("Right - Works", "Right", Financial_Common_Rights),
		Entry("Preferred Stock - Works", "Preferred Stock", Financial_Common_PreferredStock),
		Entry("Structured Product - Works", "Structured Product", Financial_Common_StructuredProduct),
		Entry("Corporate Bond - Works", "Corporate Bond", Financial_Common_CorporateBond),
		Entry("Agency Bond - Works", "Agency Bond", Financial_Common_AgencyBond),
		Entry("Equity-Linked Bond - Works", "Equity-Linked Bond", Financial_Common_EquityLinkedBond),
		Entry("Liquidating Trust - Works", "Liquidating Trust", Financial_Common_LiquidatingTrust),
		Entry("Other - Works", "Other", Financial_Common_Others),
		Entry("CommonShare - Works", "CommonShare", Financial_Common_CommonShare),
		Entry("OrdinaryShare - Works", "OrdinaryShare", Financial_Common_OrdinaryShare),
		Entry("NewYorkRegistryShares - Works", "NewYorkRegistryShares", Financial_Common_NewYorkRegistryShares),
		Entry("AmericanDepositoryReceiptCommon - Works",
			"AmericanDepositoryReceiptCommon", Financial_Common_AmericanDepositoryReceiptCommon),
		Entry("AmericanDepositoryReceiptPreferred - Works",
			"AmericanDepositoryReceiptPreferred", Financial_Common_AmericanDepositoryReceiptPreferred),
		Entry("AmericanDepositoryReceiptRights - Works",
			"AmericanDepositoryReceiptRights", Financial_Common_AmericanDepositoryReceiptRights),
		Entry("AmericanDepositoryReceiptWarrants - Works",
			"AmericanDepositoryReceiptWarrants", Financial_Common_AmericanDepositoryReceiptWarrants),
		Entry("GlobalDepositoryReceipts - Works", "GlobalDepositoryReceipts", Financial_Common_GlobalDepositoryReceipts),
		Entry("Unit - Works", "Unit", Financial_Common_Unit),
		Entry("Rights - Works", "Rights", Financial_Common_Rights),
		Entry("PreferredStock - Works", "PreferredStock", Financial_Common_PreferredStock),
		Entry("Fund - Works", "Fund", Financial_Common_Fund),
		Entry("StructuredProduct - Works", "StructuredProduct", Financial_Common_StructuredProduct),
		Entry("Warrant - Works", "Warrant", Financial_Common_Warrant),
		Entry("Index - Works", "Index", Financial_Common_Index),
		Entry("ExchangeTradedFund - Works", "ExchangeTradedFund", Financial_Common_ExchangeTradedFund),
		Entry("ExchangeTradedNote - Works", "ExchangeTradedNote", Financial_Common_ExchangeTradedNote),
		Entry("CorporateBond - Works", "CorporateBond", Financial_Common_CorporateBond),
		Entry("AgencyBond - Works", "AgencyBond", Financial_Common_AgencyBond),
		Entry("EquityLinkedBond - Works", "EquityLinkedBond", Financial_Common_EquityLinkedBond),
		Entry("Basket - Works", "Basket", Financial_Common_Basket),
		Entry("LiquidatingTrust - Works", "LiquidatingTrust", Financial_Common_LiquidatingTrust),
		Entry("Others - Works", "Others", Financial_Common_Others),
		Entry("None - Works", "None", Financial_Common_None),
		Entry("0 - Works", "0", Financial_Common_CommonShare),
		Entry("1 - Works", "1", Financial_Common_OrdinaryShare),
		Entry("2 - Works", "2", Financial_Common_NewYorkRegistryShares),
		Entry("3 - Works", "3", Financial_Common_AmericanDepositoryReceiptCommon),
		Entry("4 - Works", "4", Financial_Common_AmericanDepositoryReceiptPreferred),
		Entry("5 - Works", "5", Financial_Common_AmericanDepositoryReceiptRights),
		Entry("6 - Works", "6", Financial_Common_AmericanDepositoryReceiptWarrants),
		Entry("7 - Works", "7", Financial_Common_GlobalDepositoryReceipts),
		Entry("8 - Works", "8", Financial_Common_Unit),
		Entry("9 - Works", "9", Financial_Common_Rights),
		Entry("10 - Works", "10", Financial_Common_PreferredStock),
		Entry("11 - Works", "11", Financial_Common_Fund),
		Entry("12 - Works", "12", Financial_Common_StructuredProduct),
		Entry("13 - Works", "13", Financial_Common_Warrant),
		Entry("14 - Works", "14", Financial_Common_Index),
		Entry("15 - Works", "15", Financial_Common_ExchangeTradedFund),
		Entry("16 - Works", "16", Financial_Common_ExchangeTradedNote),
		Entry("17 - Works", "17", Financial_Common_CorporateBond),
		Entry("18 - Works", "18", Financial_Common_AgencyBond),
		Entry("19 - Works", "19", Financial_Common_EquityLinkedBond),
		Entry("20 - Works", "20", Financial_Common_Basket),
		Entry("21 - Works", "21", Financial_Common_LiquidatingTrust),
		Entry("22 - Works", "22", Financial_Common_Others),
		Entry("23 - Works", "23", Financial_Common_None))

	// Tests that, if the attribute type submitted to UnmarshalDynamoDBAttributeValue is not one we
	// recognize, then the function will return an error
	It("UnmarshalDynamoDBAttributeValue - AttributeValue type invalid - Error", func() {
		value := new(Financial_Common_AssetType)
		err := attributevalue.Unmarshal(&types.AttributeValueMemberBOOL{Value: true}, &value)
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("Attribute value of *types.AttributeValueMemberBOOL could not be converted to a Financial.Common.AssetType"))
	})

	// Tests the conditions under which UnmarshalDynamoDBAttributeValue is called and no error is generated
	DescribeTable("UnmarshalDynamoDBAttributeValue - AttributeValue Conditions",
		func(raw types.AttributeValue, expected Financial_Common_AssetType) {
			var value Financial_Common_AssetType
			err := attributevalue.Unmarshal(raw, &value)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(expected))
		},
		Entry("Value is []bytes, CS - Works",
			&types.AttributeValueMemberB{Value: []byte("CS")}, Financial_Common_CommonShare),
		Entry("Value is []bytes, OS - Works",
			&types.AttributeValueMemberB{Value: []byte("OS")}, Financial_Common_OrdinaryShare),
		Entry("Value is []bytes, NYRS - Works",
			&types.AttributeValueMemberB{Value: []byte("NYRS")}, Financial_Common_NewYorkRegistryShares),
		Entry("Value is []bytes, ADRC - Works",
			&types.AttributeValueMemberB{Value: []byte("ADRC")}, Financial_Common_AmericanDepositoryReceiptCommon),
		Entry("Value is []bytes, ADRP - Works",
			&types.AttributeValueMemberB{Value: []byte("ADRP")}, Financial_Common_AmericanDepositoryReceiptPreferred),
		Entry("Value is []bytes, ADRR - Works",
			&types.AttributeValueMemberB{Value: []byte("ADRR")}, Financial_Common_AmericanDepositoryReceiptRights),
		Entry("Value is []bytes, ADRW - Works",
			&types.AttributeValueMemberB{Value: []byte("ADRW")}, Financial_Common_AmericanDepositoryReceiptWarrants),
		Entry("Value is []bytes, GDR - Works",
			&types.AttributeValueMemberB{Value: []byte("GDR")}, Financial_Common_GlobalDepositoryReceipts),
		Entry("Value is []bytes, UNIT - Works",
			&types.AttributeValueMemberB{Value: []byte("UNIT")}, Financial_Common_Unit),
		Entry("Value is []bytes, RIGHT - Works",
			&types.AttributeValueMemberB{Value: []byte("RIGHT")}, Financial_Common_Rights),
		Entry("Value is []bytes, PFD - Works",
			&types.AttributeValueMemberB{Value: []byte("PFD")}, Financial_Common_PreferredStock),
		Entry("Value is []bytes, FUND - Works",
			&types.AttributeValueMemberB{Value: []byte("FUND")}, Financial_Common_Fund),
		Entry("Value is []bytes, SP - Works",
			&types.AttributeValueMemberB{Value: []byte("SP")}, Financial_Common_StructuredProduct),
		Entry("Value is []bytes, WARRANT - Works",
			&types.AttributeValueMemberB{Value: []byte("WARRANT")}, Financial_Common_Warrant),
		Entry("Value is []bytes, INDEX - Works",
			&types.AttributeValueMemberB{Value: []byte("INDEX")}, Financial_Common_Index),
		Entry("Value is []bytes, ETF - Works",
			&types.AttributeValueMemberB{Value: []byte("ETF")}, Financial_Common_ExchangeTradedFund),
		Entry("Value is []bytes, ETN - Works",
			&types.AttributeValueMemberB{Value: []byte("ETN")}, Financial_Common_ExchangeTradedNote),
		Entry("Value is []bytes, BOND - Works",
			&types.AttributeValueMemberB{Value: []byte("BOND")}, Financial_Common_CorporateBond),
		Entry("Value is []bytes, AGEN - Works",
			&types.AttributeValueMemberB{Value: []byte("AGEN")}, Financial_Common_AgencyBond),
		Entry("Value is []bytes, EQLK - Works",
			&types.AttributeValueMemberB{Value: []byte("EQLK")}, Financial_Common_EquityLinkedBond),
		Entry("Value is []bytes, BASKET - Works",
			&types.AttributeValueMemberB{Value: []byte("BASKET")}, Financial_Common_Basket),
		Entry("Value is []bytes, LT - Works",
			&types.AttributeValueMemberB{Value: []byte("LT")}, Financial_Common_LiquidatingTrust),
		Entry("Value is []bytes, OTHER - Works",
			&types.AttributeValueMemberB{Value: []byte("OTHER")}, Financial_Common_Others),
		Entry("Value is []bytes, Empty string - Works",
			&types.AttributeValueMemberB{Value: []byte("")}, Financial_Common_None),
		Entry("Value is []bytes, Common Share - Works",
			&types.AttributeValueMemberB{Value: []byte("Common Share")}, Financial_Common_CommonShare),
		Entry("Value is []bytes, Ordinary Share - Works",
			&types.AttributeValueMemberB{Value: []byte("Ordinary Share")}, Financial_Common_OrdinaryShare),
		Entry("Value is []bytes, New York Registry Share - Works",
			&types.AttributeValueMemberB{Value: []byte("New York Registry Share")}, Financial_Common_NewYorkRegistryShares),
		Entry("Value is []bytes, Common ADR - Works",
			&types.AttributeValueMemberB{Value: []byte("Common ADR")}, Financial_Common_AmericanDepositoryReceiptCommon),
		Entry("Value is []bytes, Preferred ADR - Works",
			&types.AttributeValueMemberB{Value: []byte("Preferred ADR")}, Financial_Common_AmericanDepositoryReceiptPreferred),
		Entry("Value is []bytes, ADR Right - Works",
			&types.AttributeValueMemberB{Value: []byte("ADR Right")}, Financial_Common_AmericanDepositoryReceiptRights),
		Entry("Value is []bytes, ADR Warrant - Works",
			&types.AttributeValueMemberB{Value: []byte("ADR Warrant")}, Financial_Common_AmericanDepositoryReceiptWarrants),
		Entry("Value is []bytes, Right - Works",
			&types.AttributeValueMemberB{Value: []byte("Right")}, Financial_Common_Rights),
		Entry("Value is []bytes, Preferred Stock - Works",
			&types.AttributeValueMemberB{Value: []byte("Preferred Stock")}, Financial_Common_PreferredStock),
		Entry("Value is []bytes, Structured Product - Works",
			&types.AttributeValueMemberB{Value: []byte("Structured Product")}, Financial_Common_StructuredProduct),
		Entry("Value is []bytes, Corporate Bond - Works",
			&types.AttributeValueMemberB{Value: []byte("Corporate Bond")}, Financial_Common_CorporateBond),
		Entry("Value is []bytes, Agency Bond - Works",
			&types.AttributeValueMemberB{Value: []byte("Agency Bond")}, Financial_Common_AgencyBond),
		Entry("Value is []bytes, Equity-Linked Bond - Works",
			&types.AttributeValueMemberB{Value: []byte("Equity-Linked Bond")}, Financial_Common_EquityLinkedBond),
		Entry("Value is []bytes, Liquidating Trust - Works",
			&types.AttributeValueMemberB{Value: []byte("Liquidating Trust")}, Financial_Common_LiquidatingTrust),
		Entry("Value is []bytes, Other - Works",
			&types.AttributeValueMemberB{Value: []byte("Other")}, Financial_Common_Others),
		Entry("Value is []bytes, CommonShare - Works",
			&types.AttributeValueMemberB{Value: []byte("CommonShare")}, Financial_Common_CommonShare),
		Entry("Value is []bytes, OrdinaryShare - Works",
			&types.AttributeValueMemberB{Value: []byte("OrdinaryShare")}, Financial_Common_OrdinaryShare),
		Entry("Value is []bytes, NewYorkRegistryShares - Works",
			&types.AttributeValueMemberB{Value: []byte("NewYorkRegistryShares")}, Financial_Common_NewYorkRegistryShares),
		Entry("Value is []bytes, AmericanDepositoryReceiptCommon - Works",
			&types.AttributeValueMemberB{Value: []byte("AmericanDepositoryReceiptCommon")}, Financial_Common_AmericanDepositoryReceiptCommon),
		Entry("Value is []bytes, AmericanDepositoryReceiptPreferred - Works",
			&types.AttributeValueMemberB{Value: []byte("AmericanDepositoryReceiptPreferred")}, Financial_Common_AmericanDepositoryReceiptPreferred),
		Entry("Value is []bytes, AmericanDepositoryReceiptRights - Works",
			&types.AttributeValueMemberB{Value: []byte("AmericanDepositoryReceiptRights")}, Financial_Common_AmericanDepositoryReceiptRights),
		Entry("Value is []bytes, AmericanDepositoryReceiptWarrants - Works",
			&types.AttributeValueMemberB{Value: []byte("AmericanDepositoryReceiptWarrants")}, Financial_Common_AmericanDepositoryReceiptWarrants),
		Entry("Value is []bytes, GlobalDepositoryReceipts - Works",
			&types.AttributeValueMemberB{Value: []byte("GlobalDepositoryReceipts")}, Financial_Common_GlobalDepositoryReceipts),
		Entry("Value is []bytes, Unit - Works",
			&types.AttributeValueMemberB{Value: []byte("Unit")}, Financial_Common_Unit),
		Entry("Value is []bytes, Rights - Works",
			&types.AttributeValueMemberB{Value: []byte("Rights")}, Financial_Common_Rights),
		Entry("Value is []bytes, PreferredStock - Works",
			&types.AttributeValueMemberB{Value: []byte("PreferredStock")}, Financial_Common_PreferredStock),
		Entry("Value is []bytes, Fund - Works",
			&types.AttributeValueMemberB{Value: []byte("Fund")}, Financial_Common_Fund),
		Entry("Value is []bytes, StructuredProduct - Works",
			&types.AttributeValueMemberB{Value: []byte("StructuredProduct")}, Financial_Common_StructuredProduct),
		Entry("Value is []bytes, Warrant - Works",
			&types.AttributeValueMemberB{Value: []byte("Warrant")}, Financial_Common_Warrant),
		Entry("Value is []bytes, Index - Works",
			&types.AttributeValueMemberB{Value: []byte("Index")}, Financial_Common_Index),
		Entry("Value is []bytes, ExchangeTradedFund - Works",
			&types.AttributeValueMemberB{Value: []byte("ExchangeTradedFund")}, Financial_Common_ExchangeTradedFund),
		Entry("Value is []bytes, ExchangeTradedNote - Works",
			&types.AttributeValueMemberB{Value: []byte("ExchangeTradedNote")}, Financial_Common_ExchangeTradedNote),
		Entry("Value is []bytes, CorporateBond - Works",
			&types.AttributeValueMemberB{Value: []byte("CorporateBond")}, Financial_Common_CorporateBond),
		Entry("Value is []bytes, AgencyBond - Works",
			&types.AttributeValueMemberB{Value: []byte("AgencyBond")}, Financial_Common_AgencyBond),
		Entry("Value is []bytes, EquityLinkedBond - Works",
			&types.AttributeValueMemberB{Value: []byte("EquityLinkedBond")}, Financial_Common_EquityLinkedBond),
		Entry("Value is []bytes, Basket - Works",
			&types.AttributeValueMemberB{Value: []byte("Basket")}, Financial_Common_Basket),
		Entry("Value is []bytes, LiquidatingTrust - Works",
			&types.AttributeValueMemberB{Value: []byte("LiquidatingTrust")}, Financial_Common_LiquidatingTrust),
		Entry("Value is []bytes, Others - Works",
			&types.AttributeValueMemberB{Value: []byte("Others")}, Financial_Common_Others),
		Entry("Value is []bytes, None - Works",
			&types.AttributeValueMemberB{Value: []byte("None")}, Financial_Common_None),
		Entry("Value is numeric, 0 - Works",
			&types.AttributeValueMemberN{Value: "0"}, Financial_Common_CommonShare),
		Entry("Value is numeric, 1 - Works",
			&types.AttributeValueMemberN{Value: "1"}, Financial_Common_OrdinaryShare),
		Entry("Value is numeric, 2 - Works",
			&types.AttributeValueMemberN{Value: "2"}, Financial_Common_NewYorkRegistryShares),
		Entry("Value is numeric, 3 - Works",
			&types.AttributeValueMemberN{Value: "3"}, Financial_Common_AmericanDepositoryReceiptCommon),
		Entry("Value is numeric, 4 - Works",
			&types.AttributeValueMemberN{Value: "4"}, Financial_Common_AmericanDepositoryReceiptPreferred),
		Entry("Value is numeric, 5 - Works",
			&types.AttributeValueMemberN{Value: "5"}, Financial_Common_AmericanDepositoryReceiptRights),
		Entry("Value is numeric, 6 - Works",
			&types.AttributeValueMemberN{Value: "6"}, Financial_Common_AmericanDepositoryReceiptWarrants),
		Entry("Value is numeric, 7 - Works",
			&types.AttributeValueMemberN{Value: "7"}, Financial_Common_GlobalDepositoryReceipts),
		Entry("Value is numeric, 8 - Works",
			&types.AttributeValueMemberN{Value: "8"}, Financial_Common_Unit),
		Entry("Value is numeric, 9 - Works",
			&types.AttributeValueMemberN{Value: "9"}, Financial_Common_Rights),
		Entry("Value is numeric, 10 - Works",
			&types.AttributeValueMemberN{Value: "10"}, Financial_Common_PreferredStock),
		Entry("Value is numeric, 11 - Works",
			&types.AttributeValueMemberN{Value: "11"}, Financial_Common_Fund),
		Entry("Value is numeric, 12 - Works",
			&types.AttributeValueMemberN{Value: "12"}, Financial_Common_StructuredProduct),
		Entry("Value is numeric, 13 - Works",
			&types.AttributeValueMemberN{Value: "13"}, Financial_Common_Warrant),
		Entry("Value is numeric, 14 - Works",
			&types.AttributeValueMemberN{Value: "14"}, Financial_Common_Index),
		Entry("Value is numeric, 15 - Works",
			&types.AttributeValueMemberN{Value: "15"}, Financial_Common_ExchangeTradedFund),
		Entry("Value is numeric, 16 - Works",
			&types.AttributeValueMemberN{Value: "16"}, Financial_Common_ExchangeTradedNote),
		Entry("Value is numeric, 17 - Works",
			&types.AttributeValueMemberN{Value: "17"}, Financial_Common_CorporateBond),
		Entry("Value is numeric, 18 - Works",
			&types.AttributeValueMemberN{Value: "18"}, Financial_Common_AgencyBond),
		Entry("Value is numeric, 19 - Works",
			&types.AttributeValueMemberN{Value: "19"}, Financial_Common_EquityLinkedBond),
		Entry("Value is numeric, 20 - Works",
			&types.AttributeValueMemberN{Value: "20"}, Financial_Common_Basket),
		Entry("Value is numeric, 21 - Works",
			&types.AttributeValueMemberN{Value: "21"}, Financial_Common_LiquidatingTrust),
		Entry("Value is numeric, 22 - Works",
			&types.AttributeValueMemberN{Value: "22"}, Financial_Common_Others),
		Entry("Value is numeric, 23 - Works",
			&types.AttributeValueMemberN{Value: "23"}, Financial_Common_None),
		Entry("Value is NULL - Works", new(types.AttributeValueMemberNULL), Financial_Common_AssetType(0)),
		Entry("Value is string, CS - Works",
			&types.AttributeValueMemberS{Value: "CS"}, Financial_Common_CommonShare),
		Entry("Value is string, OS - Works",
			&types.AttributeValueMemberS{Value: "OS"}, Financial_Common_OrdinaryShare),
		Entry("Value is string, NYRS - Works",
			&types.AttributeValueMemberS{Value: "NYRS"}, Financial_Common_NewYorkRegistryShares),
		Entry("Value is string, ADRC - Works",
			&types.AttributeValueMemberS{Value: "ADRC"}, Financial_Common_AmericanDepositoryReceiptCommon),
		Entry("Value is string, ADRP - Works",
			&types.AttributeValueMemberS{Value: "ADRP"}, Financial_Common_AmericanDepositoryReceiptPreferred),
		Entry("Value is string, ADRR - Works",
			&types.AttributeValueMemberS{Value: "ADRR"}, Financial_Common_AmericanDepositoryReceiptRights),
		Entry("Value is string, ADRW - Works",
			&types.AttributeValueMemberS{Value: "ADRW"}, Financial_Common_AmericanDepositoryReceiptWarrants),
		Entry("Value is string, GDR - Works",
			&types.AttributeValueMemberS{Value: "GDR"}, Financial_Common_GlobalDepositoryReceipts),
		Entry("Value is string, UNIT - Works",
			&types.AttributeValueMemberS{Value: "UNIT"}, Financial_Common_Unit),
		Entry("Value is string, RIGHT - Works",
			&types.AttributeValueMemberS{Value: "RIGHT"}, Financial_Common_Rights),
		Entry("Value is string, PFD - Works",
			&types.AttributeValueMemberS{Value: "PFD"}, Financial_Common_PreferredStock),
		Entry("Value is string, FUND - Works",
			&types.AttributeValueMemberS{Value: "FUND"}, Financial_Common_Fund),
		Entry("Value is string, SP - Works",
			&types.AttributeValueMemberS{Value: "SP"}, Financial_Common_StructuredProduct),
		Entry("Value is string, WARRANT - Works",
			&types.AttributeValueMemberS{Value: "WARRANT"}, Financial_Common_Warrant),
		Entry("Value is string, INDEX - Works",
			&types.AttributeValueMemberS{Value: "INDEX"}, Financial_Common_Index),
		Entry("Value is string, ETF - Works",
			&types.AttributeValueMemberS{Value: "ETF"}, Financial_Common_ExchangeTradedFund),
		Entry("Value is string, ETN - Works",
			&types.AttributeValueMemberS{Value: "ETN"}, Financial_Common_ExchangeTradedNote),
		Entry("Value is string, BOND - Works",
			&types.AttributeValueMemberS{Value: "BOND"}, Financial_Common_CorporateBond),
		Entry("Value is string, AGEN - Works",
			&types.AttributeValueMemberS{Value: "AGEN"}, Financial_Common_AgencyBond),
		Entry("Value is string, EQLK - Works",
			&types.AttributeValueMemberS{Value: "EQLK"}, Financial_Common_EquityLinkedBond),
		Entry("Value is string, BASKET - Works",
			&types.AttributeValueMemberS{Value: "BASKET"}, Financial_Common_Basket),
		Entry("Value is string, LT - Works",
			&types.AttributeValueMemberS{Value: "LT"}, Financial_Common_LiquidatingTrust),
		Entry("Value is string, OTHER - Works",
			&types.AttributeValueMemberS{Value: "OTHER"}, Financial_Common_Others),
		Entry("Value is string, Empty string - Works",
			&types.AttributeValueMemberS{Value: ""}, Financial_Common_None),
		Entry("Value is string, Common Share - Works",
			&types.AttributeValueMemberS{Value: "Common Share"}, Financial_Common_CommonShare),
		Entry("Value is string, Ordinary Share - Works",
			&types.AttributeValueMemberS{Value: "Ordinary Share"}, Financial_Common_OrdinaryShare),
		Entry("Value is string, New York Registry Share - Works",
			&types.AttributeValueMemberS{Value: "New York Registry Share"}, Financial_Common_NewYorkRegistryShares),
		Entry("Value is string, Common ADR - Works",
			&types.AttributeValueMemberS{Value: "Common ADR"}, Financial_Common_AmericanDepositoryReceiptCommon),
		Entry("Value is string, Preferred ADR - Works",
			&types.AttributeValueMemberS{Value: "Preferred ADR"}, Financial_Common_AmericanDepositoryReceiptPreferred),
		Entry("Value is string, ADR Right - Works",
			&types.AttributeValueMemberS{Value: "ADR Right"}, Financial_Common_AmericanDepositoryReceiptRights),
		Entry("Value is string, ADR Warrant - Works",
			&types.AttributeValueMemberS{Value: "ADR Warrant"}, Financial_Common_AmericanDepositoryReceiptWarrants),
		Entry("Value is string, Right - Works",
			&types.AttributeValueMemberS{Value: "Right"}, Financial_Common_Rights),
		Entry("Value is string, Preferred Stock - Works",
			&types.AttributeValueMemberS{Value: "Preferred Stock"}, Financial_Common_PreferredStock),
		Entry("Value is string, Structured Product - Works",
			&types.AttributeValueMemberS{Value: "Structured Product"}, Financial_Common_StructuredProduct),
		Entry("Value is string, Corporate Bond - Works",
			&types.AttributeValueMemberS{Value: "Corporate Bond"}, Financial_Common_CorporateBond),
		Entry("Value is string, Agency Bond - Works",
			&types.AttributeValueMemberS{Value: "Agency Bond"}, Financial_Common_AgencyBond),
		Entry("Value is string, Equity-Linked Bond - Works",
			&types.AttributeValueMemberS{Value: "Equity-Linked Bond"}, Financial_Common_EquityLinkedBond),
		Entry("Value is string, Liquidating Trust - Works",
			&types.AttributeValueMemberS{Value: "Liquidating Trust"}, Financial_Common_LiquidatingTrust),
		Entry("Value is string, Other - Works",
			&types.AttributeValueMemberS{Value: "Other"}, Financial_Common_Others),
		Entry("Value is string, CommonShare - Works",
			&types.AttributeValueMemberS{Value: "CommonShare"}, Financial_Common_CommonShare),
		Entry("Value is string, OrdinaryShare - Works",
			&types.AttributeValueMemberS{Value: "OrdinaryShare"}, Financial_Common_OrdinaryShare),
		Entry("Value is string, NewYorkRegistryShares - Works",
			&types.AttributeValueMemberS{Value: "NewYorkRegistryShares"}, Financial_Common_NewYorkRegistryShares),
		Entry("Value is string, AmericanDepositoryReceiptCommon - Works",
			&types.AttributeValueMemberS{Value: "AmericanDepositoryReceiptCommon"}, Financial_Common_AmericanDepositoryReceiptCommon),
		Entry("Value is string, AmericanDepositoryReceiptPreferred - Works",
			&types.AttributeValueMemberS{Value: "AmericanDepositoryReceiptPreferred"}, Financial_Common_AmericanDepositoryReceiptPreferred),
		Entry("Value is string, AmericanDepositoryReceiptRights - Works",
			&types.AttributeValueMemberS{Value: "AmericanDepositoryReceiptRights"}, Financial_Common_AmericanDepositoryReceiptRights),
		Entry("Value is string, AmericanDepositoryReceiptWarrants - Works",
			&types.AttributeValueMemberS{Value: "AmericanDepositoryReceiptWarrants"}, Financial_Common_AmericanDepositoryReceiptWarrants),
		Entry("Value is string, GlobalDepositoryReceipts - Works",
			&types.AttributeValueMemberS{Value: "GlobalDepositoryReceipts"}, Financial_Common_GlobalDepositoryReceipts),
		Entry("Value is string, Unit - Works",
			&types.AttributeValueMemberS{Value: "Unit"}, Financial_Common_Unit),
		Entry("Value is string, Rights - Works",
			&types.AttributeValueMemberS{Value: "Rights"}, Financial_Common_Rights),
		Entry("Value is string, PreferredStock - Works",
			&types.AttributeValueMemberS{Value: "PreferredStock"}, Financial_Common_PreferredStock),
		Entry("Value is string, Fund - Works",
			&types.AttributeValueMemberS{Value: "Fund"}, Financial_Common_Fund),
		Entry("Value is string, StructuredProduct - Works",
			&types.AttributeValueMemberS{Value: "StructuredProduct"}, Financial_Common_StructuredProduct),
		Entry("Value is string, Warrant - Works",
			&types.AttributeValueMemberS{Value: "Warrant"}, Financial_Common_Warrant),
		Entry("Value is string, Index - Works",
			&types.AttributeValueMemberS{Value: "Index"}, Financial_Common_Index),
		Entry("Value is string, ExchangeTradedFund - Works",
			&types.AttributeValueMemberS{Value: "ExchangeTradedFund"}, Financial_Common_ExchangeTradedFund),
		Entry("Value is string, ExchangeTradedNote - Works",
			&types.AttributeValueMemberS{Value: "ExchangeTradedNote"}, Financial_Common_ExchangeTradedNote),
		Entry("Value is string, CorporateBond - Works",
			&types.AttributeValueMemberS{Value: "CorporateBond"}, Financial_Common_CorporateBond),
		Entry("Value is string, AgencyBond - Works",
			&types.AttributeValueMemberS{Value: "AgencyBond"}, Financial_Common_AgencyBond),
		Entry("Value is string, EquityLinkedBond - Works",
			&types.AttributeValueMemberS{Value: "EquityLinkedBond"}, Financial_Common_EquityLinkedBond),
		Entry("Value is string, Basket - Works",
			&types.AttributeValueMemberS{Value: "Basket"}, Financial_Common_Basket),
		Entry("Value is string, LiquidatingTrust - Works",
			&types.AttributeValueMemberS{Value: "LiquidatingTrust"}, Financial_Common_LiquidatingTrust),
		Entry("Value is string, Others - Works",
			&types.AttributeValueMemberS{Value: "Others"}, Financial_Common_Others),
		Entry("Value is string, None - Works",
			&types.AttributeValueMemberS{Value: "None"}, Financial_Common_None))

	// Test that attempting to deserialize a Financial.Common.AssetType will fial and return an
	// error if the value cannot be converted to either the name value or integer value of the enum option
	It("Scan - Value is nil - Error", func() {

		// Attempt to convert a fake string value into a Financial.Common.AssetType
		// This should return an error
		var enum *Financial_Common_AssetType
		err := enum.Scan(nil)

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of %!q(<nil>) had an invalid type of <nil>"))
		Expect(enum).Should(BeNil())
	})

	// Test the conditions under which values should be convertible to a Financial.Common.AssetType
	DescribeTable("Scan Tests",
		func(value interface{}, shouldBe Financial_Common_AssetType) {

			// Attempt to convert the value into a Financial.Common.AssetType
			// This should not fail
			var enum Financial_Common_AssetType
			err := enum.Scan(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("CS - Works", "CS", Financial_Common_CommonShare),
		Entry("OS - Works", "OS", Financial_Common_OrdinaryShare),
		Entry("NYRS - Works", "NYRS", Financial_Common_NewYorkRegistryShares),
		Entry("ADRC - Works", "ADRC", Financial_Common_AmericanDepositoryReceiptCommon),
		Entry("ADRP - Works", "ADRP", Financial_Common_AmericanDepositoryReceiptPreferred),
		Entry("ADRR - Works", "ADRR", Financial_Common_AmericanDepositoryReceiptRights),
		Entry("ADRW - Works", "ADRW", Financial_Common_AmericanDepositoryReceiptWarrants),
		Entry("GDR - Works", "GDR", Financial_Common_GlobalDepositoryReceipts),
		Entry("UNIT - Works", "UNIT", Financial_Common_Unit),
		Entry("RIGHT - Works", "RIGHT", Financial_Common_Rights),
		Entry("PFD - Works", "PFD", Financial_Common_PreferredStock),
		Entry("FUND - Works", "FUND", Financial_Common_Fund),
		Entry("SP - Works", "SP", Financial_Common_StructuredProduct),
		Entry("WARRANT - Works", "WARRANT", Financial_Common_Warrant),
		Entry("INDEX - Works", "INDEX", Financial_Common_Index),
		Entry("ETF - Works", "ETF", Financial_Common_ExchangeTradedFund),
		Entry("ETN - Works", "ETN", Financial_Common_ExchangeTradedNote),
		Entry("BOND - Works", "BOND", Financial_Common_CorporateBond),
		Entry("AGEN - Works", "AGEN", Financial_Common_AgencyBond),
		Entry("EQLK - Works", "EQLK", Financial_Common_EquityLinkedBond),
		Entry("LT - Works", "LT", Financial_Common_LiquidatingTrust),
		Entry("BASKET - Works", "BASKET", Financial_Common_Basket),
		Entry("OTHER - Works", "OTHER", Financial_Common_Others),
		Entry("Empty string - Works", "", Financial_Common_None),
		Entry("Common Share - Works", "Common Share", Financial_Common_CommonShare),
		Entry("Ordinary Share - Works", "Ordinary Share", Financial_Common_OrdinaryShare),
		Entry("New York Registry Share - Works", "New York Registry Share", Financial_Common_NewYorkRegistryShares),
		Entry("Common ADR - Works", "Common ADR", Financial_Common_AmericanDepositoryReceiptCommon),
		Entry("Preferred ADR - Works", "Preferred ADR", Financial_Common_AmericanDepositoryReceiptPreferred),
		Entry("ADR Right - Works", "ADR Right", Financial_Common_AmericanDepositoryReceiptRights),
		Entry("ADR Warrant - Works", "ADR Warrant", Financial_Common_AmericanDepositoryReceiptWarrants),
		Entry("Right - Works", "Right", Financial_Common_Rights),
		Entry("Preferred Stock - Works", "Preferred Stock", Financial_Common_PreferredStock),
		Entry("Structured Product - Works", "Structured Product", Financial_Common_StructuredProduct),
		Entry("Corporate Bond - Works", "Corporate Bond", Financial_Common_CorporateBond),
		Entry("Agency Bond - Works", "Agency Bond", Financial_Common_AgencyBond),
		Entry("Equity-Linked Bond - Works", "Equity-Linked Bond", Financial_Common_EquityLinkedBond),
		Entry("Liquidating Trust - Works", "Liquidating Trust", Financial_Common_LiquidatingTrust),
		Entry("Other - Works", "Other", Financial_Common_Others),
		Entry("CommonShare - Works", "CommonShare", Financial_Common_CommonShare),
		Entry("OrdinaryShare - Works", "OrdinaryShare", Financial_Common_OrdinaryShare),
		Entry("NewYorkRegistryShares - Works", "NewYorkRegistryShares", Financial_Common_NewYorkRegistryShares),
		Entry("AmericanDepositoryReceiptCommon - Works",
			"AmericanDepositoryReceiptCommon", Financial_Common_AmericanDepositoryReceiptCommon),
		Entry("AmericanDepositoryReceiptPreferred - Works",
			"AmericanDepositoryReceiptPreferred", Financial_Common_AmericanDepositoryReceiptPreferred),
		Entry("AmericanDepositoryReceiptRights - Works",
			"AmericanDepositoryReceiptRights", Financial_Common_AmericanDepositoryReceiptRights),
		Entry("AmericanDepositoryReceiptWarrants - Works",
			"AmericanDepositoryReceiptWarrants", Financial_Common_AmericanDepositoryReceiptWarrants),
		Entry("GlobalDepositoryReceipts - Works", "GlobalDepositoryReceipts", Financial_Common_GlobalDepositoryReceipts),
		Entry("Unit - Works", "Unit", Financial_Common_Unit),
		Entry("Rights - Works", "Rights", Financial_Common_Rights),
		Entry("PreferredStock - Works", "PreferredStock", Financial_Common_PreferredStock),
		Entry("Fund - Works", "Fund", Financial_Common_Fund),
		Entry("StructuredProduct - Works", "StructuredProduct", Financial_Common_StructuredProduct),
		Entry("Warrant - Works", "Warrant", Financial_Common_Warrant),
		Entry("Index - Works", "Index", Financial_Common_Index),
		Entry("ExchangeTradedFund - Works", "ExchangeTradedFund", Financial_Common_ExchangeTradedFund),
		Entry("ExchangeTradedNote - Works", "ExchangeTradedNote", Financial_Common_ExchangeTradedNote),
		Entry("CorporateBond - Works", "CorporateBond", Financial_Common_CorporateBond),
		Entry("AgencyBond - Works", "AgencyBond", Financial_Common_AgencyBond),
		Entry("EquityLinkedBond - Works", "EquityLinkedBond", Financial_Common_EquityLinkedBond),
		Entry("Basket - Works", "Basket", Financial_Common_Basket),
		Entry("LiquidatingTrust - Works", "LiquidatingTrust", Financial_Common_LiquidatingTrust),
		Entry("Others - Works", "Others", Financial_Common_Others),
		Entry("None - Works", "None", Financial_Common_None),
		Entry("0 - Works", 0, Financial_Common_CommonShare),
		Entry("1 - Works", 1, Financial_Common_OrdinaryShare),
		Entry("2 - Works", 2, Financial_Common_NewYorkRegistryShares),
		Entry("3 - Works", 3, Financial_Common_AmericanDepositoryReceiptCommon),
		Entry("4 - Works", 4, Financial_Common_AmericanDepositoryReceiptPreferred),
		Entry("5 - Works", 5, Financial_Common_AmericanDepositoryReceiptRights),
		Entry("6 - Works", 6, Financial_Common_AmericanDepositoryReceiptWarrants),
		Entry("7 - Works", 7, Financial_Common_GlobalDepositoryReceipts),
		Entry("8 - Works", 8, Financial_Common_Unit),
		Entry("9 - Works", 9, Financial_Common_Rights),
		Entry("10 - Works", 10, Financial_Common_PreferredStock),
		Entry("11 - Works", 11, Financial_Common_Fund),
		Entry("12 - Works", 12, Financial_Common_StructuredProduct),
		Entry("13 - Works", 13, Financial_Common_Warrant),
		Entry("14 - Works", 14, Financial_Common_Index),
		Entry("15 - Works", 15, Financial_Common_ExchangeTradedFund),
		Entry("16 - Works", 16, Financial_Common_ExchangeTradedNote),
		Entry("17 - Works", 17, Financial_Common_CorporateBond),
		Entry("18 - Works", 18, Financial_Common_AgencyBond),
		Entry("19 - Works", 19, Financial_Common_EquityLinkedBond),
		Entry("20 - Works", 20, Financial_Common_Basket),
		Entry("21 - Works", 21, Financial_Common_LiquidatingTrust),
		Entry("22 - Works", 22, Financial_Common_Others),
		Entry("23 - Works", 23, Financial_Common_None))
})

var _ = Describe("Financial.Common.Locale Marshal/Unmarshal Tests", func() {

	// Test that converting the Financial.Common.Locale enum to JSON works for all values
	DescribeTable("MarshalJSON Tests",
		func(enum Financial_Common_Locale, value string) {
			data, err := json.Marshal(enum)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("US - Works", Financial_Common_US, "\"US\""),
		Entry("Global - Works", Financial_Common_Global, "\"Global\""))

	// Test that converting the Financial.Common.Locale enum to a CSV column works for all values
	DescribeTable("MarshalCSV Tests",
		func(enum Financial_Common_Locale, value string) {
			data, err := enum.MarshalCSV()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("US - Works", Financial_Common_US, "0"),
		Entry("Global - Works", Financial_Common_Global, "1"))

	// Test that converting the Financial.Common.Locale enum to a AttributeValue works for all values
	DescribeTable("MarshalDynamoDBAttributeValue Tests",
		func(enum Financial_Common_Locale, value string) {
			data, err := enum.MarshalDynamoDBAttributeValue()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data.(*types.AttributeValueMemberS).Value).Should(Equal(value))
		},
		Entry("US - Works", Financial_Common_US, "US"),
		Entry("Global - Works", Financial_Common_Global, "Global"))

	// Test that attempting to deserialize a Financial.Common.Locale will fail and
	// return an error if the value canno be deserialized from a JSON value to a string
	It("UnmarshalJSON fails - Error", func() {

		// Attempt to convert a non-parseable string value into a Financial.Common.Locale
		// This should return an error
		enum := new(Financial_Common_Locale)
		err := enum.UnmarshalJSON([]byte("derp"))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Common_Locale"))
	})

	// Test that attempting to deserialize a Financial.Common.Locale will fail and
	// return an error if the value cannot be converted to either the name value or integer
	// value of the enum option
	It("UnmarshalJSON - Value is invalid - Error", func() {

		// Attempt to convert a fake string value into a Financial.Common.Locale
		// This should return an error
		enum := new(Financial_Common_Locale)
		err := enum.UnmarshalJSON([]byte("\"derp\""))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Common_Locale"))
	})

	// Test the conditions under which values should be convertible to a Financial.Common.Locale
	DescribeTable("UnmarshalJSON Tests",
		func(value string, shouldBe Financial_Common_Locale) {

			// Attempt to convert the string value into a Financial.Common.Locale
			// This should not fail
			var enum Financial_Common_Locale
			err := enum.UnmarshalJSON([]byte(value))

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("US - Works", "\"US\"", Financial_Common_US),
		Entry("Global - Works", "\"Global\"", Financial_Common_Global),
		Entry("us - Works", "\"us\"", Financial_Common_US),
		Entry("global - Works", "\"global\"", Financial_Common_Global),
		Entry("0 - Works", "\"0\"", Financial_Common_US),
		Entry("1 - Works", "\"1\"", Financial_Common_Global))

	// Test that attempting to deserialize a Financial.Common.Locale will fial and return an
	// error if the value cannot be converted to either the name value or integer value
	// of the enum option
	It("UnmarshalCSV - Value is empty - Error", func() {

		// Attempt to convert a fake string value into a Financial.Common.Locale
		// This should return an error
		enum := new(Financial_Common_Locale)
		err := enum.UnmarshalCSV("")

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"\" cannot be mapped to a gopb.Financial_Common_Locale"))
	})

	// Test the conditions under which values should be convertible to a Financial.Common.Locale
	DescribeTable("UnmarshalCSV Tests",
		func(value string, shouldBe Financial_Common_Locale) {

			// Attempt to convert the value into a Financial.Common.Locale
			// This should not fail
			var enum Financial_Common_Locale
			err := enum.UnmarshalCSV(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("US - Works", "US", Financial_Common_US),
		Entry("Global - Works", "Global", Financial_Common_Global),
		Entry("us - Works", "us", Financial_Common_US),
		Entry("global - Works", "global", Financial_Common_Global),
		Entry("0 - Works", "0", Financial_Common_US),
		Entry("1 - Works", "1", Financial_Common_Global))

	// Tests that, if the attribute type submitted to UnmarshalDynamoDBAttributeValue is not one we
	// recognize, then the function will return an error
	It("UnmarshalDynamoDBAttributeValue - AttributeValue type invalid - Error", func() {
		value := new(Financial_Common_Locale)
		err := attributevalue.Unmarshal(&types.AttributeValueMemberBOOL{Value: true}, &value)
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("Attribute value of *types.AttributeValueMemberBOOL could not be converted to a Financial.Common.Locale"))
	})

	// Tests the conditions under which UnmarshalDynamoDBAttributeValue is called and no error is generated
	DescribeTable("UnmarshalDynamoDBAttributeValue - AttributeValue Conditions",
		func(raw types.AttributeValue, expected Financial_Common_Locale) {
			var value Financial_Common_Locale
			err := attributevalue.Unmarshal(raw, &value)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(expected))
		},
		Entry("Value is []bytes, US - Works",
			&types.AttributeValueMemberB{Value: []byte("US")}, Financial_Common_US),
		Entry("Value is []bytes, Global - Works",
			&types.AttributeValueMemberB{Value: []byte("Global")}, Financial_Common_Global),
		Entry("Value is numeric, 0 - Works",
			&types.AttributeValueMemberN{Value: "0"}, Financial_Common_US),
		Entry("Value is numeric, 1 - Works",
			&types.AttributeValueMemberN{Value: "1"}, Financial_Common_Global),
		Entry("Value is NULL - Works", new(types.AttributeValueMemberNULL), Financial_Common_Locale(0)),
		Entry("Value is string, US - Works",
			&types.AttributeValueMemberS{Value: "US"}, Financial_Common_US),
		Entry("Value is string, Global - Works",
			&types.AttributeValueMemberS{Value: "Global"}, Financial_Common_Global))

	// Test that attempting to deserialize a Financial.Common.Locale will fial and return an
	// error if the value cannot be converted to either the name value or integer value of the enum option
	It("Scan - Value is nil - Error", func() {

		// Attempt to convert a fake string value into a Financial.Common.Locale
		// This should return an error
		var enum *Financial_Common_Locale
		err := enum.Scan(nil)

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of %!q(<nil>) had an invalid type of <nil>"))
		Expect(enum).Should(BeNil())
	})

	// Test the conditions under which values should be convertible to a Financial.Common.Locale
	DescribeTable("Scan Tests",
		func(value interface{}, shouldBe Financial_Common_Locale) {

			// Attempt to convert the value into a Financial.Common.Locale
			// This should not fail
			var enum Financial_Common_Locale
			err := enum.Scan(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("US - Works", "US", Financial_Common_US),
		Entry("Global - Works", "Global", Financial_Common_Global),
		Entry("us - Works", "us", Financial_Common_US),
		Entry("global - Works", "global", Financial_Common_Global),
		Entry("0 - Works", 0, Financial_Common_US),
		Entry("1 - Works", 1, Financial_Common_Global))
})

var _ = Describe("Financial.Common.Tape Marshal/Unmarshal Tests", func() {

	// Test that converting the Financial.Common.Tape enum to JSON works for all values
	DescribeTable("MarshalJSON Tests",
		func(enum Financial_Common_Tape, value string) {
			data, err := json.Marshal(enum)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("A - Works", Financial_Common_A, "\"A\""),
		Entry("B - Works", Financial_Common_B, "\"B\""),
		Entry("C - Works", Financial_Common_C, "\"C\""))

	// Test that converting the Financial.Common.Tape enum to a CSV column works for all values
	DescribeTable("MarshalCSV Tests",
		func(enum Financial_Common_Tape, value string) {
			data, err := enum.MarshalCSV()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("A - Works", Financial_Common_A, "A"),
		Entry("B - Works", Financial_Common_B, "B"),
		Entry("C - Works", Financial_Common_C, "C"))

	// Test that converting the Financial.Common.Tape enum to an SQL value for all values
	DescribeTable("Value Tests",
		func(enum Financial_Common_Tape, value string) {
			data, err := enum.Value()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data).Should(Equal(value))
		},
		Entry("A - Works", Financial_Common_A, "A"),
		Entry("B - Works", Financial_Common_B, "B"),
		Entry("C - Works", Financial_Common_C, "C"))

	// Test that converting the Financial.Common.Tape enum to a AttributeValue works for all values
	DescribeTable("MarshalDynamoDBAttributeValue Tests",
		func(enum Financial_Common_Tape, value string) {
			data, err := enum.MarshalDynamoDBAttributeValue()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data.(*types.AttributeValueMemberS).Value).Should(Equal(value))
		},
		Entry("A - Works", Financial_Common_A, "A"),
		Entry("B - Works", Financial_Common_B, "B"),
		Entry("C - Works", Financial_Common_C, "C"))

	// Test that attempting to deserialize a Financial.Common.Tape will fail and
	// return an error if the value canno be deserialized from a JSON value to a string
	It("UnmarshalJSON fails - Error", func() {

		// Attempt to convert a non-parseable string value into a Financial.Common.Tape
		// This should return an error
		enum := new(Financial_Common_Tape)
		err := enum.UnmarshalJSON([]byte("derp"))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Common_Tape"))
	})

	// Test that attempting to deserialize a Financial.Common.Tape will fail and
	// return an error if the value cannot be converted to either the name value or integer
	// value of the enum option
	It("UnmarshalJSON - Value is invalid - Error", func() {

		// Attempt to convert a fake string value into a Financial.Common.Tape
		// This should return an error
		enum := new(Financial_Common_Tape)
		err := enum.UnmarshalJSON([]byte("\"derp\""))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Common_Tape"))
	})

	// Test the conditions under which values should be convertible to a Financial.Common.Tape
	DescribeTable("UnmarshalJSON Tests",
		func(value interface{}, shouldBe Financial_Common_Tape) {

			// Attempt to convert the string value into a Financial.Common.Tape
			// This should not fail
			var enum Financial_Common_Tape
			err := enum.UnmarshalJSON([]byte(fmt.Sprintf("%v", value)))

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("A - Works", "\"A\"", Financial_Common_A),
		Entry("B - Works", "\"B\"", Financial_Common_B),
		Entry("C - Works", "\"C\"", Financial_Common_C),
		Entry("'0' - Works", "\"0\"", Financial_Common_A),
		Entry("'1' - Works", "\"1\"", Financial_Common_B),
		Entry("'2' - Works", "\"2\"", Financial_Common_C),
		Entry("0 - Works", 0, Financial_Common_A),
		Entry("1 - Works", 1, Financial_Common_B),
		Entry("2 - Works", 2, Financial_Common_C))

	// Test that attempting to deserialize a Financial.Common.Tape will fial and return an
	// error if the value cannot be converted to either the name value or integer value
	// of the enum option
	It("UnmarshalCSV - Value is empty - Error", func() {

		// Attempt to convert a fake string value into a Financial.Common.Tape
		// This should return an error
		enum := new(Financial_Common_Tape)
		err := enum.UnmarshalCSV("")

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"\" cannot be mapped to a gopb.Financial_Common_Tape"))
	})

	// Test the conditions under which values should be convertible to a Financial.Common.Tape
	DescribeTable("UnmarshalCSV Tests",
		func(value string, shouldBe Financial_Common_Tape) {

			// Attempt to convert the value into a Financial.Common.Tape
			// This should not fail
			var enum Financial_Common_Tape
			err := enum.UnmarshalCSV(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("A - Works", "A", Financial_Common_A),
		Entry("B - Works", "B", Financial_Common_B),
		Entry("C - Works", "C", Financial_Common_C),
		Entry("0 - Works", "0", Financial_Common_A),
		Entry("1 - Works", "1", Financial_Common_B),
		Entry("2 - Works", "2", Financial_Common_C))

	// Tests that, if the attribute type submitted to UnmarshalDynamoDBAttributeValue is not one we
	// recognize, then the function will return an error
	It("UnmarshalDynamoDBAttributeValue - AttributeValue type invalid - Error", func() {
		value := new(Financial_Common_Tape)
		err := attributevalue.Unmarshal(&types.AttributeValueMemberBOOL{Value: true}, &value)
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("Attribute value of *types.AttributeValueMemberBOOL could not be converted to a Financial.Common.Tape"))
	})

	// Tests the conditions under which UnmarshalDynamoDBAttributeValue is called and no error is generated
	DescribeTable("UnmarshalDynamoDBAttributeValue - AttributeValue Conditions",
		func(raw types.AttributeValue, expected Financial_Common_Tape) {
			var value Financial_Common_Tape
			err := attributevalue.Unmarshal(raw, &value)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(expected))
		},
		Entry("Value is []bytes, A - Works",
			&types.AttributeValueMemberB{Value: []byte("A")}, Financial_Common_A),
		Entry("Value is []bytes, B - Works",
			&types.AttributeValueMemberB{Value: []byte("B")}, Financial_Common_B),
		Entry("Value is []bytes, C - Works",
			&types.AttributeValueMemberB{Value: []byte("C")}, Financial_Common_C),
		Entry("Value is numeric, 0 - Works",
			&types.AttributeValueMemberN{Value: "0"}, Financial_Common_A),
		Entry("Value is numeric, 1 - Works",
			&types.AttributeValueMemberN{Value: "1"}, Financial_Common_B),
		Entry("Value is numeric, 2 - Works",
			&types.AttributeValueMemberN{Value: "2"}, Financial_Common_C),
		Entry("Value is NULL - Works", new(types.AttributeValueMemberNULL), Financial_Common_Tape(0)),
		Entry("Value is string, A - Works",
			&types.AttributeValueMemberS{Value: "A"}, Financial_Common_A),
		Entry("Value is string, B - Works",
			&types.AttributeValueMemberS{Value: "B"}, Financial_Common_B),
		Entry("Value is string, C - Works",
			&types.AttributeValueMemberS{Value: "C"}, Financial_Common_C))

	// Test that attempting to deserialize a Financial.Common.Tape will fial and return an
	// error if the value cannot be converted to either the name value or integer value of the enum option
	It("Scan - Value is nil - Error", func() {

		// Attempt to convert a fake string value into a Financial.Common.Tape
		// This should return an error
		var enum *Financial_Common_Tape
		err := enum.Scan(nil)

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of %!q(<nil>) had an invalid type of <nil>"))
		Expect(enum).Should(BeNil())
	})

	// Test the conditions under which values should be convertible to a Financial.Common.Tape
	DescribeTable("Scan Tests",
		func(value interface{}, shouldBe Financial_Common_Tape) {

			// Attempt to convert the value into a Financial.Common.Tape
			// This should not fail
			var enum Financial_Common_Tape
			err := enum.Scan(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("A - Works", "A", Financial_Common_A),
		Entry("B - Works", "B", Financial_Common_B),
		Entry("C - Works", "C", Financial_Common_C),
		Entry("0 - Works", 0, Financial_Common_A),
		Entry("1 - Works", 1, Financial_Common_B),
		Entry("2 - Works", 2, Financial_Common_C))
})

var _ = Describe("Financial.Dividends.Frequency Marshal/Unmarshal Tests", func() {

	// Test that converting the Financial.Dividends.Frequency enum to JSON works for all values
	DescribeTable("MarshalJSON Tests",
		func(enum Financial_Dividends_Frequency, value string) {
			data, err := json.Marshal(enum)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("NoFrequency - Works", Financial_Dividends_NoFrequency, "\"\""),
		Entry("Annually - Works", Financial_Dividends_Annually, "\"Annually\""),
		Entry("SemiAnnually - Works", Financial_Dividends_SemiAnnually, "\"SemiAnnually\""),
		Entry("Quarterly - Works", Financial_Dividends_Quarterly, "\"Quarterly\""),
		Entry("Monthly - Works", Financial_Dividends_Monthly, "\"Monthly\""),
		Entry("Invalid - Works", Financial_Dividends_Invalid, "\"Invalid\""))

	// Test that converting the Financial.Dividends.Frequency enum to a CSV column works for all values
	DescribeTable("MarshalCSV Tests",
		func(enum Financial_Dividends_Frequency, value string) {
			data, err := enum.MarshalCSV()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("NoFrequency - Works", Financial_Dividends_NoFrequency, "0"),
		Entry("Annually - Works", Financial_Dividends_Annually, "1"),
		Entry("SemiAnnually - Works", Financial_Dividends_SemiAnnually, "2"),
		Entry("Quarterly - Works", Financial_Dividends_Quarterly, "4"),
		Entry("Monthly - Works", Financial_Dividends_Monthly, "12"),
		Entry("Invalid - Works", Financial_Dividends_Invalid, "13"))

	// Test that converting the Financial.Dividends.Frequency enum to a AttributeValue works for all values
	DescribeTable("MarshalDynamoDBAttributeValue Tests",
		func(enum Financial_Dividends_Frequency, value string) {
			data, err := enum.MarshalDynamoDBAttributeValue()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data.(*types.AttributeValueMemberS).Value).Should(Equal(value))
		},
		Entry("NoFrequency - Works", Financial_Dividends_NoFrequency, ""),
		Entry("Annually - Works", Financial_Dividends_Annually, "Annually"),
		Entry("SemiAnnually - Works", Financial_Dividends_SemiAnnually, "SemiAnnually"),
		Entry("Quarterly - Works", Financial_Dividends_Quarterly, "Quarterly"),
		Entry("Monthly - Works", Financial_Dividends_Monthly, "Monthly"),
		Entry("Invalid - Works", Financial_Dividends_Invalid, "Invalid"))

	// Test that attempting to deserialize a Financial.Dividends.Frequency will fail and
	// return an error if the value canno be deserialized from a JSON value to a string
	It("UnmarshalJSON fails - Error", func() {

		// Attempt to convert a non-parseable string value into a Financial.Dividends.Frequency
		// This should return an error
		enum := new(Financial_Dividends_Frequency)
		err := enum.UnmarshalJSON([]byte("derp"))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Dividends_Frequency"))
	})

	// Test that attempting to deserialize a Financial.Dividends.Frequency will fail and
	// return an error if the value cannot be converted to either the name value or integer
	// value of the enum option
	It("UnmarshalJSON - Value is invalid - Error", func() {

		// Attempt to convert a fake string value into a Financial.Dividends.Frequency
		// This should return an error
		enum := new(Financial_Dividends_Frequency)
		err := enum.UnmarshalJSON([]byte("\"derp\""))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Dividends_Frequency"))
	})

	// Test the conditions under which values should be convertible to a Financial.Dividends.Frequency
	DescribeTable("UnmarshalJSON Tests",
		func(value interface{}, shouldBe Financial_Dividends_Frequency) {

			// Attempt to convert the string value into a Financial.Dividends.Frequency
			// This should not fail
			var enum Financial_Dividends_Frequency
			err := enum.UnmarshalJSON([]byte(fmt.Sprintf("%v", value)))

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("None - Works", "\"None\"", Financial_Dividends_NoFrequency),
		Entry("Empty string - Works", "\"\"", Financial_Dividends_NoFrequency),
		Entry("NoFrequency - Works", "\"NoFrequency\"", Financial_Dividends_NoFrequency),
		Entry("Annually - Works", "\"Annually\"", Financial_Dividends_Annually),
		Entry("SemiAnnually - Works", "\"SemiAnnually\"", Financial_Dividends_SemiAnnually),
		Entry("Qarterly - Works", "\"Quarterly\"", Financial_Dividends_Quarterly),
		Entry("Monthly - Works", "\"Monthly\"", Financial_Dividends_Monthly),
		Entry("Invalid - Works", "\"Invalid\"", Financial_Dividends_Invalid),
		Entry("'0' - Works", "\"0\"", Financial_Dividends_NoFrequency),
		Entry("'1' - Works", "\"1\"", Financial_Dividends_Annually),
		Entry("'2' - Works", "\"2\"", Financial_Dividends_SemiAnnually),
		Entry("'4' - Works", "\"4\"", Financial_Dividends_Quarterly),
		Entry("'12' - Works", "\"12\"", Financial_Dividends_Monthly),
		Entry("'13' - Works", "\"13\"", Financial_Dividends_Invalid),
		Entry("0 - Works", 0, Financial_Dividends_NoFrequency),
		Entry("1 - Works", 1, Financial_Dividends_Annually),
		Entry("2 - Works", 2, Financial_Dividends_SemiAnnually),
		Entry("4 - Works", 4, Financial_Dividends_Quarterly),
		Entry("12 - Works", 12, Financial_Dividends_Monthly),
		Entry("13 - Works", 13, Financial_Dividends_Invalid))

	// Test that attempting to deserialize a Financial.Dividends.Frequency will fial and return an
	// error if the value cannot be converted to either the name value or integer value
	// of the enum option
	It("UnmarshalCSV - Value is invalid - Error", func() {

		// Attempt to convert a fake string value into a Financial.Dividends.Frequency
		// This should return an error
		enum := new(Financial_Dividends_Frequency)
		err := enum.UnmarshalCSV("derp")

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Dividends_Frequency"))
	})

	// Test the conditions under which values should be convertible to a Financial.Dividends.Frequency
	DescribeTable("UnmarshalCSV Tests",
		func(value string, shouldBe Financial_Dividends_Frequency) {

			// Attempt to convert the value into a Financial.Dividends.Frequency
			// This should not fail
			var enum Financial_Dividends_Frequency
			err := enum.UnmarshalCSV(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("None - Works", "None", Financial_Dividends_NoFrequency),
		Entry("Empty string - Works", "", Financial_Dividends_NoFrequency),
		Entry("NoFrequency - Works", "NoFrequency", Financial_Dividends_NoFrequency),
		Entry("Annually - Works", "Annually", Financial_Dividends_Annually),
		Entry("SemiAnnually - Works", "SemiAnnually", Financial_Dividends_SemiAnnually),
		Entry("Quarterly - Works", "Quarterly", Financial_Dividends_Quarterly),
		Entry("Monthly - Works", "Monthly", Financial_Dividends_Monthly),
		Entry("Invalid - Works", "Invalid", Financial_Dividends_Invalid),
		Entry("0 - Works", "0", Financial_Dividends_NoFrequency),
		Entry("1 - Works", "1", Financial_Dividends_Annually),
		Entry("2 - Works", "2", Financial_Dividends_SemiAnnually),
		Entry("4 - Works", "4", Financial_Dividends_Quarterly),
		Entry("12 - Works", "12", Financial_Dividends_Monthly),
		Entry("13 - Works", "13", Financial_Dividends_Invalid))

	// Tests that, if the attribute type submitted to UnmarshalDynamoDBAttributeValue is not one we
	// recognize, then the function will return an error
	It("UnmarshalDynamoDBAttributeValue - AttributeValue type invalid - Error", func() {
		value := new(Financial_Dividends_Frequency)
		err := attributevalue.Unmarshal(&types.AttributeValueMemberBOOL{Value: true}, &value)
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("Attribute value of *types.AttributeValueMemberBOOL could not be converted to a Financial.Dividends.Frequency"))
	})

	// Tests the conditions under which UnmarshalDynamoDBAttributeValue is called and no error is generated
	DescribeTable("UnmarshalDynamoDBAttributeValue - AttributeValue Conditions",
		func(raw types.AttributeValue, expected Financial_Dividends_Frequency) {
			var value Financial_Dividends_Frequency
			err := attributevalue.Unmarshal(raw, &value)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(expected))
		},
		Entry("Value is []bytes, None - Works",
			&types.AttributeValueMemberB{Value: []byte("None")}, Financial_Dividends_NoFrequency),
		Entry("Value is []bytes, Empty string - Works",
			&types.AttributeValueMemberB{Value: []byte("")}, Financial_Dividends_NoFrequency),
		Entry("Value is []bytes, NoFrequency - Works",
			&types.AttributeValueMemberB{Value: []byte("NoFrequency")}, Financial_Dividends_NoFrequency),
		Entry("Value is []bytes, Annually - Works",
			&types.AttributeValueMemberB{Value: []byte("Annually")}, Financial_Dividends_Annually),
		Entry("Value is []bytes, SemiAnnually - Works",
			&types.AttributeValueMemberB{Value: []byte("SemiAnnually")}, Financial_Dividends_SemiAnnually),
		Entry("Value is []bytes, Quarterly - Works",
			&types.AttributeValueMemberB{Value: []byte("Quarterly")}, Financial_Dividends_Quarterly),
		Entry("Value is []bytes, Monthly - Works",
			&types.AttributeValueMemberB{Value: []byte("Monthly")}, Financial_Dividends_Monthly),
		Entry("Value is []bytes, Invalid - Works",
			&types.AttributeValueMemberB{Value: []byte("Invalid")}, Financial_Dividends_Invalid),
		Entry("Value is numeric, 0 - Works",
			&types.AttributeValueMemberN{Value: "0"}, Financial_Dividends_NoFrequency),
		Entry("Value is numeric, 1 - Works",
			&types.AttributeValueMemberN{Value: "1"}, Financial_Dividends_Annually),
		Entry("Value is numeric, 2 - Works",
			&types.AttributeValueMemberN{Value: "2"}, Financial_Dividends_SemiAnnually),
		Entry("Value is numeric, 4 - Works",
			&types.AttributeValueMemberN{Value: "4"}, Financial_Dividends_Quarterly),
		Entry("Value is numeric, 12 - Works",
			&types.AttributeValueMemberN{Value: "12"}, Financial_Dividends_Monthly),
		Entry("Value is numeric, 13 - Works",
			&types.AttributeValueMemberN{Value: "13"}, Financial_Dividends_Invalid),
		Entry("Value is NULL - Works", new(types.AttributeValueMemberNULL), Financial_Dividends_Frequency(0)),
		Entry("Value is string, None - Works",
			&types.AttributeValueMemberS{Value: "None"}, Financial_Dividends_NoFrequency),
		Entry("Value is string, Empty string - Works",
			&types.AttributeValueMemberS{Value: ""}, Financial_Dividends_NoFrequency),
		Entry("Value is string, NoFrequency - Works",
			&types.AttributeValueMemberS{Value: "NoFrequency"}, Financial_Dividends_NoFrequency),
		Entry("Value is string, Annually - Works",
			&types.AttributeValueMemberS{Value: "Annually"}, Financial_Dividends_Annually),
		Entry("Value is string, SemiAnnually - Works",
			&types.AttributeValueMemberS{Value: "SemiAnnually"}, Financial_Dividends_SemiAnnually),
		Entry("Value is string, Quarterly - Works",
			&types.AttributeValueMemberS{Value: "Quarterly"}, Financial_Dividends_Quarterly),
		Entry("Value is string, Monthly - Works",
			&types.AttributeValueMemberS{Value: "Monthly"}, Financial_Dividends_Monthly),
		Entry("Value is string, Invalid - Works",
			&types.AttributeValueMemberS{Value: "Invalid"}, Financial_Dividends_Invalid))

	// Test that attempting to deserialize a Financial.Dividends.Frequency will fial and return an
	// error if the value cannot be converted to either the name value or integer value of the enum option
	It("Scan - Value is nil - Error", func() {

		// Attempt to convert a fake string value into a Financial.Dividends.Frequency
		// This should return an error
		var enum *Financial_Dividends_Frequency
		err := enum.Scan(nil)

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of %!q(<nil>) had an invalid type of <nil>"))
		Expect(enum).Should(BeNil())
	})

	// Test the conditions under which values should be convertible to a Financial.Dividends.Frequency
	DescribeTable("Scan Tests",
		func(value interface{}, shouldBe Financial_Dividends_Frequency) {

			// Attempt to convert the value into a Financial.Dividends.Frequency
			// This should not fail
			var enum Financial_Dividends_Frequency
			err := enum.Scan(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("None - Works", "None", Financial_Dividends_NoFrequency),
		Entry("Empty string - Works", "", Financial_Dividends_NoFrequency),
		Entry("NoFrequency - Works", "NoFrequency", Financial_Dividends_NoFrequency),
		Entry("Annually - Works", "Annually", Financial_Dividends_Annually),
		Entry("SemiAnnually - Works", "SemiAnnually", Financial_Dividends_SemiAnnually),
		Entry("Quarterly - Works", "Quarterly", Financial_Dividends_Quarterly),
		Entry("Monthly - Works", "Monthly", Financial_Dividends_Monthly),
		Entry("Invalid - Works", "Invalid", Financial_Dividends_Invalid),
		Entry("0 - Works", 0, Financial_Dividends_NoFrequency),
		Entry("1 - Works", 1, Financial_Dividends_Annually),
		Entry("2 - Works", 2, Financial_Dividends_SemiAnnually),
		Entry("4 - Works", 4, Financial_Dividends_Quarterly),
		Entry("12 - Works", 12, Financial_Dividends_Monthly),
		Entry("13 - Works", 13, Financial_Dividends_Invalid))
})

var _ = Describe("Financial.Dividends.Type Marshal/Unmarshal Tests", func() {

	// Test that converting the Financial.Dividends.Type enum to JSON works for all values
	DescribeTable("MarshalJSON Tests",
		func(enum Financial_Dividends_Type, value string) {
			data, err := json.Marshal(enum)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("CD - Works", Financial_Dividends_CD, "\"CD\""),
		Entry("SC - Works", Financial_Dividends_SC, "\"SC\""),
		Entry("LT - Works", Financial_Dividends_LT, "\"LT\""),
		Entry("ST - Works", Financial_Dividends_ST, "\"ST\""),
		Entry("NP - Works", Financial_Dividends_NP, "\"NP\""))

	// Test that converting the Financial.Dividends.Type enum to a CSV column works for all values
	DescribeTable("MarshalCSV Tests",
		func(enum Financial_Dividends_Type, value string) {
			data, err := enum.MarshalCSV()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("CD - Works", Financial_Dividends_CD, "CD"),
		Entry("SC - Works", Financial_Dividends_SC, "SC"),
		Entry("LT - Works", Financial_Dividends_LT, "LT"),
		Entry("ST - Works", Financial_Dividends_ST, "ST"),
		Entry("NP - Works", Financial_Dividends_NP, "NP"))

	// Test that converting the Financial.Dividends.Type enum to an SQL value for all values
	DescribeTable("Value Tests",
		func(enum Financial_Dividends_Type, value string) {
			data, err := enum.Value()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data).Should(Equal(value))
		},
		Entry("CD - Works", Financial_Dividends_CD, "CD"),
		Entry("SC - Works", Financial_Dividends_SC, "SC"),
		Entry("LT - Works", Financial_Dividends_LT, "LT"),
		Entry("ST - Works", Financial_Dividends_ST, "ST"),
		Entry("NP - Works", Financial_Dividends_NP, "NP"))

	// Test that converting the Financial.Dividends.Type enum to a AttributeValue works for all values
	DescribeTable("MarshalDynamoDBAttributeValue Tests",
		func(enum Financial_Dividends_Type, value string) {
			data, err := enum.MarshalDynamoDBAttributeValue()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data.(*types.AttributeValueMemberS).Value).Should(Equal(value))
		},
		Entry("CD - Works", Financial_Dividends_CD, "CD"),
		Entry("SC - Works", Financial_Dividends_SC, "SC"),
		Entry("LT - Works", Financial_Dividends_LT, "LT"),
		Entry("ST - Works", Financial_Dividends_ST, "ST"),
		Entry("NP - Works", Financial_Dividends_NP, "NP"))

	// Test that attempting to deserialize a Financial.Dividends.Type will fail and
	// return an error if the value canno be deserialized from a JSON value to a string
	It("UnmarshalJSON fails - Error", func() {

		// Attempt to convert a non-parseable string value into a Financial.Dividends.Type
		// This should return an error
		enum := new(Financial_Dividends_Type)
		err := enum.UnmarshalJSON([]byte("derp"))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Dividends_Type"))
	})

	// Test that attempting to deserialize a Financial.Dividends.Type will fail and
	// return an error if the value cannot be converted to either the name value or integer
	// value of the enum option
	It("UnmarshalJSON - Value is invalid - Error", func() {

		// Attempt to convert a fake string value into a Financial.Dividends.Type
		// This should return an error
		enum := new(Financial_Dividends_Type)
		err := enum.UnmarshalJSON([]byte("\"derp\""))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Dividends_Type"))
	})

	// Test the conditions under which values should be convertible to a Financial.Dividends.Type
	DescribeTable("UnmarshalJSON Tests",
		func(value string, shouldBe Financial_Dividends_Type) {

			// Attempt to convert the string value into a Financial.Dividends.Type
			// This should not fail
			var enum Financial_Dividends_Type
			err := enum.UnmarshalJSON([]byte(value))

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("CD - Works", "\"CD\"", Financial_Dividends_CD),
		Entry("SC - Works", "\"SC\"", Financial_Dividends_SC),
		Entry("LT - Works", "\"LT\"", Financial_Dividends_LT),
		Entry("ST - Works", "\"ST\"", Financial_Dividends_ST),
		Entry("NP - Works", "\"NP\"", Financial_Dividends_NP),
		Entry("0 - Works", "\"0\"", Financial_Dividends_CD),
		Entry("1 - Works", "\"1\"", Financial_Dividends_SC),
		Entry("2 - Works", "\"2\"", Financial_Dividends_LT),
		Entry("3 - Works", "\"3\"", Financial_Dividends_ST),
		Entry("4 - Works", "\"4\"", Financial_Dividends_NP))

	// Test that attempting to deserialize a Financial.Dividends.Type will fial and return an
	// error if the value cannot be converted to either the name value or integer value
	// of the enum option
	It("UnmarshalCSV - Value is empty - Error", func() {

		// Attempt to convert a fake string value into a Financial.Dividends.Type
		// This should return an error
		enum := new(Financial_Dividends_Type)
		err := enum.UnmarshalCSV("")

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"\" cannot be mapped to a gopb.Financial_Dividends_Type"))
	})

	// Test the conditions under which values should be convertible to a Financial.Dividends.Type
	DescribeTable("UnmarshalCSV Tests",
		func(value string, shouldBe Financial_Dividends_Type) {

			// Attempt to convert the value into a Financial.Dividends.Type
			// This should not fail
			var enum Financial_Dividends_Type
			err := enum.UnmarshalCSV(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("CD - Works", "CD", Financial_Dividends_CD),
		Entry("SC - Works", "SC", Financial_Dividends_SC),
		Entry("LT - Works", "LT", Financial_Dividends_LT),
		Entry("ST - Works", "ST", Financial_Dividends_ST),
		Entry("NP - Works", "NP", Financial_Dividends_NP),
		Entry("0 - Works", "0", Financial_Dividends_CD),
		Entry("1 - Works", "1", Financial_Dividends_SC),
		Entry("2 - Works", "2", Financial_Dividends_LT),
		Entry("3 - Works", "3", Financial_Dividends_ST),
		Entry("4 - Works", "4", Financial_Dividends_NP))

	// Tests that, if the attribute type submitted to UnmarshalDynamoDBAttributeValue is not one we
	// recognize, then the function will return an error
	It("UnmarshalDynamoDBAttributeValue - AttributeValue type invalid - Error", func() {
		value := new(Financial_Dividends_Type)
		err := attributevalue.Unmarshal(&types.AttributeValueMemberBOOL{Value: true}, &value)
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("Attribute value of *types.AttributeValueMemberBOOL could not be converted to a Financial.Dividends.Type"))
	})

	// Tests the conditions under which UnmarshalDynamoDBAttributeValue is called and no error is generated
	DescribeTable("UnmarshalDynamoDBAttributeValue - AttributeValue Conditions",
		func(raw types.AttributeValue, expected Financial_Dividends_Type) {
			var value Financial_Dividends_Type
			err := attributevalue.Unmarshal(raw, &value)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(expected))
		},
		Entry("Value is []bytes, CD - Works",
			&types.AttributeValueMemberB{Value: []byte("CD")}, Financial_Dividends_CD),
		Entry("Value is []bytes, SC - Works",
			&types.AttributeValueMemberB{Value: []byte("SC")}, Financial_Dividends_SC),
		Entry("Value is []bytes, LT - Works",
			&types.AttributeValueMemberB{Value: []byte("LT")}, Financial_Dividends_LT),
		Entry("Value is []bytes, ST - Works",
			&types.AttributeValueMemberB{Value: []byte("ST")}, Financial_Dividends_ST),
		Entry("Value is []bytes, NP - Works",
			&types.AttributeValueMemberB{Value: []byte("NP")}, Financial_Dividends_NP),
		Entry("Value is numeric, 0 - Works",
			&types.AttributeValueMemberN{Value: "0"}, Financial_Dividends_CD),
		Entry("Value is numeric, 1 - Works",
			&types.AttributeValueMemberN{Value: "1"}, Financial_Dividends_SC),
		Entry("Value is numeric, 2 - Works",
			&types.AttributeValueMemberN{Value: "2"}, Financial_Dividends_LT),
		Entry("Value is numeric, 3 - Works",
			&types.AttributeValueMemberN{Value: "3"}, Financial_Dividends_ST),
		Entry("Value is numeric, 4 - Works",
			&types.AttributeValueMemberN{Value: "4"}, Financial_Dividends_NP),
		Entry("Value is NULL - Works", new(types.AttributeValueMemberNULL), Financial_Dividends_Type(0)),
		Entry("Value is string, CD - Works",
			&types.AttributeValueMemberS{Value: "CD"}, Financial_Dividends_CD),
		Entry("Value is string, SC - Works",
			&types.AttributeValueMemberS{Value: "SC"}, Financial_Dividends_SC),
		Entry("Value is string, LT - Works",
			&types.AttributeValueMemberS{Value: "LT"}, Financial_Dividends_LT),
		Entry("Value is string, ST - Works",
			&types.AttributeValueMemberS{Value: "ST"}, Financial_Dividends_ST),
		Entry("Value is string, NP - Works",
			&types.AttributeValueMemberS{Value: "NP"}, Financial_Dividends_NP))

	// Test that attempting to deserialize a Financial.Dividends.Type will fial and return an
	// error if the value cannot be converted to either the name value or integer value of the enum option
	It("Scan - Value is nil - Error", func() {

		// Attempt to convert a fake string value into a Financial.Dividends.Type
		// This should return an error
		var enum *Financial_Dividends_Type
		err := enum.Scan(nil)

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of %!q(<nil>) had an invalid type of <nil>"))
		Expect(enum).Should(BeNil())
	})

	// Test the conditions under which values should be convertible to a Financial.Dividends.Type
	DescribeTable("Scan Tests",
		func(value interface{}, shouldBe Financial_Dividends_Type) {

			// Attempt to convert the value into a Financial.Dividends.Type
			// This should not fail
			var enum Financial_Dividends_Type
			err := enum.Scan(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("CD - Works", "CD", Financial_Dividends_CD),
		Entry("SC - Works", "SC", Financial_Dividends_SC),
		Entry("LT - Works", "LT", Financial_Dividends_LT),
		Entry("ST - Works", "ST", Financial_Dividends_ST),
		Entry("NP - Works", "NP", Financial_Dividends_NP),
		Entry("0 - Works", 0, Financial_Dividends_CD),
		Entry("1 - Works", 1, Financial_Dividends_SC),
		Entry("2 - Works", 2, Financial_Dividends_LT),
		Entry("3 - Works", 3, Financial_Dividends_ST),
		Entry("4 - Works", 4, Financial_Dividends_NP))
})

var _ = Describe("Financial.Exchanges.Type Marshal/Unmarshal Tests", func() {

	// Test that converting the Financial.Exchanges.Type enum to JSON works for all values
	DescribeTable("MarshalJSON Tests",
		func(enum Financial_Exchanges_Type, value string) {
			data, err := json.Marshal(enum)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("Exchange - Works", Financial_Exchanges_Exchange, "\"Exchange\""),
		Entry("TRF - Works", Financial_Exchanges_TRF, "\"TRF\""),
		Entry("SIP - Works", Financial_Exchanges_SIP, "\"SIP\""))

	// Test that converting the Financial.Exchanges.Type enum to a CSV column works for all values
	DescribeTable("MarshalCSV Tests",
		func(enum Financial_Exchanges_Type, value string) {
			data, err := enum.MarshalCSV()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("Exchange - Works", Financial_Exchanges_Exchange, "0"),
		Entry("TRF - Works", Financial_Exchanges_TRF, "1"),
		Entry("SIP - Works", Financial_Exchanges_SIP, "2"))

	// Test that converting the Financial.Exchanges.Type enum to a AttributeValue works for all values
	DescribeTable("MarshalDynamoDBAttributeValue Tests",
		func(enum Financial_Exchanges_Type, value string) {
			data, err := enum.MarshalDynamoDBAttributeValue()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data.(*types.AttributeValueMemberS).Value).Should(Equal(value))
		},
		Entry("Exchange - Works", Financial_Exchanges_Exchange, "Exchange"),
		Entry("TRF - Works", Financial_Exchanges_TRF, "TRF"),
		Entry("SIP - Works", Financial_Exchanges_SIP, "SIP"))

	// Test that attempting to deserialize a Financial.Exchanges.Type will fail and
	// return an error if the value canno be deserialized from a JSON value to a string
	It("UnmarshalJSON fails - Error", func() {

		// Attempt to convert a non-parseable string value into a Financial.Exchanges.Type
		// This should return an error
		enum := new(Financial_Exchanges_Type)
		err := enum.UnmarshalJSON([]byte("derp"))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Exchanges_Type"))
	})

	// Test that attempting to deserialize a Financial.Exchanges.Type will fail and
	// return an error if the value cannot be converted to either the name value or integer
	// value of the enum option
	It("UnmarshalJSON - Value is invalid - Error", func() {

		// Attempt to convert a fake string value into a Financial.Exchanges.Type
		// This should return an error
		enum := new(Financial_Exchanges_Type)
		err := enum.UnmarshalJSON([]byte("\"derp\""))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Exchanges_Type"))
	})

	// Test the conditions under which values should be convertible to a Financial.Exchanges.Type
	DescribeTable("UnmarshalJSON Tests",
		func(value string, shouldBe Financial_Exchanges_Type) {

			// Attempt to convert the string value into a Financial.Exchanges.Type
			// This should not fail
			var enum Financial_Exchanges_Type
			err := enum.UnmarshalJSON([]byte(value))

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("Exchange - Works", "\"Exchange\"", Financial_Exchanges_Exchange),
		Entry("TRF - Works", "\"TRF\"", Financial_Exchanges_TRF),
		Entry("SIP - Works", "\"SIP\"", Financial_Exchanges_SIP),
		Entry("exchange - Works", "\"exchange\"", Financial_Exchanges_Exchange),
		Entry("0 - Works", "\"0\"", Financial_Exchanges_Exchange),
		Entry("1 - Works", "\"1\"", Financial_Exchanges_TRF),
		Entry("2 - Works", "\"2\"", Financial_Exchanges_SIP))

	// Test that attempting to deserialize a Financial.Exchanges.Type will fial and return an
	// error if the value cannot be converted to either the name value or integer value
	// of the enum option
	It("UnmarshalCSV - Value is empty - Error", func() {

		// Attempt to convert a fake string value into a Financial.Exchanges.Type
		// This should return an error
		enum := new(Financial_Exchanges_Type)
		err := enum.UnmarshalCSV("")

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"\" cannot be mapped to a gopb.Financial_Exchanges_Type"))
	})

	// Test the conditions under which values should be convertible to a Financial.Exchanges.Type
	DescribeTable("UnmarshalCSV Tests",
		func(value string, shouldBe Financial_Exchanges_Type) {

			// Attempt to convert the value into a Financial.Exchanges.Type
			// This should not fail
			var enum Financial_Exchanges_Type
			err := enum.UnmarshalCSV(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("Exchange - Works", "Exchange", Financial_Exchanges_Exchange),
		Entry("TRF - Works", "TRF", Financial_Exchanges_TRF),
		Entry("SIP - Works", "SIP", Financial_Exchanges_SIP),
		Entry("exchange - Works", "exchange", Financial_Exchanges_Exchange),
		Entry("0 - Works", "0", Financial_Exchanges_Exchange),
		Entry("1 - Works", "1", Financial_Exchanges_TRF),
		Entry("2 - Works", "2", Financial_Exchanges_SIP))

	// Tests that, if the attribute type submitted to UnmarshalDynamoDBAttributeValue is not one we
	// recognize, then the function will return an error
	It("UnmarshalDynamoDBAttributeValue - AttributeValue type invalid - Error", func() {
		value := new(Financial_Exchanges_Type)
		err := attributevalue.Unmarshal(&types.AttributeValueMemberBOOL{Value: true}, &value)
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("Attribute value of *types.AttributeValueMemberBOOL could not be converted to a Financial.Exchanges.Type"))
	})

	// Tests the conditions under which UnmarshalDynamoDBAttributeValue is called and no error is generated
	DescribeTable("UnmarshalDynamoDBAttributeValue - AttributeValue Conditions",
		func(raw types.AttributeValue, expected Financial_Exchanges_Type) {
			var value Financial_Exchanges_Type
			err := attributevalue.Unmarshal(raw, &value)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(expected))
		},
		Entry("Value is []bytes, exchange - Works",
			&types.AttributeValueMemberB{Value: []byte("exchange")}, Financial_Exchanges_Exchange),
		Entry("Value is []bytes, Exchange - Works",
			&types.AttributeValueMemberB{Value: []byte("Exchange")}, Financial_Exchanges_Exchange),
		Entry("Value is []bytes, TRF - Works",
			&types.AttributeValueMemberB{Value: []byte("TRF")}, Financial_Exchanges_TRF),
		Entry("Value is []bytes, SIP - Works",
			&types.AttributeValueMemberB{Value: []byte("SIP")}, Financial_Exchanges_SIP),
		Entry("Value is numeric, 0 - Works",
			&types.AttributeValueMemberN{Value: "0"}, Financial_Exchanges_Exchange),
		Entry("Value is numeric, 1 - Works",
			&types.AttributeValueMemberN{Value: "1"}, Financial_Exchanges_TRF),
		Entry("Value is numeric, 2 - Works",
			&types.AttributeValueMemberN{Value: "2"}, Financial_Exchanges_SIP),
		Entry("Value is NULL - Works", new(types.AttributeValueMemberNULL), Financial_Exchanges_Type(0)),
		Entry("Value is string, exchange - Works",
			&types.AttributeValueMemberS{Value: "exchange"}, Financial_Exchanges_Exchange),
		Entry("Value is string, Exchange - Works",
			&types.AttributeValueMemberS{Value: "Exchange"}, Financial_Exchanges_Exchange),
		Entry("Value is string, TRF - Works",
			&types.AttributeValueMemberS{Value: "TRF"}, Financial_Exchanges_TRF),
		Entry("Value is string, SIP - Works",
			&types.AttributeValueMemberS{Value: "SIP"}, Financial_Exchanges_SIP))

	// Test that attempting to deserialize a Financial.Exchanges.Type will fial and return an
	// error if the value cannot be converted to either the name value or integer value of the enum option
	It("Scan - Value is nil - Error", func() {

		// Attempt to convert a fake string value into a Financial.Exchanges.Type
		// This should return an error
		var enum *Financial_Exchanges_Type
		err := enum.Scan(nil)

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of %!q(<nil>) had an invalid type of <nil>"))
		Expect(enum).Should(BeNil())
	})

	// Test the conditions under which values should be convertible to a Financial.Exchanges.Type
	DescribeTable("Scan Tests",
		func(value interface{}, shouldBe Financial_Exchanges_Type) {

			// Attempt to convert the value into a Financial.Exchanges.Type
			// This should not fail
			var enum Financial_Exchanges_Type
			err := enum.Scan(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("Exchange - Works", "Exchange", Financial_Exchanges_Exchange),
		Entry("TRF - Works", "TRF", Financial_Exchanges_TRF),
		Entry("SIP - Works", "SIP", Financial_Exchanges_SIP),
		Entry("exchange - Works", "exchange", Financial_Exchanges_Exchange),
		Entry("0 - Works", 0, Financial_Exchanges_Exchange),
		Entry("1 - Works", 1, Financial_Exchanges_TRF),
		Entry("2 - Works", 2, Financial_Exchanges_SIP))
})

var _ = Describe("Financial.Options.ContractType Marshal/Unmarshal Tests", func() {

	// Test that converting the Financial.Options.ContractType enum to JSON works for all values
	DescribeTable("MarshalJSON Tests",
		func(enum Financial_Options_ContractType, value string) {
			data, err := json.Marshal(enum)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("Call - Works", Financial_Options_Call, "\"Call\""),
		Entry("Put - Works", Financial_Options_Put, "\"Put\""),
		Entry("Other - Works", Financial_Options_Other, "\"Other\""))

	// Test that converting the Financial.Options.ContractType enum to a CSV column works for all values
	DescribeTable("MarshalCSV Tests",
		func(enum Financial_Options_ContractType, value string) {
			data, err := enum.MarshalCSV()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("Call - Works", Financial_Options_Call, "0"),
		Entry("Put - Works", Financial_Options_Put, "1"),
		Entry("Other - Works", Financial_Options_Other, "2"))

	// Test that converting the Financial.Options.ContractType enum to an SQL value for all values
	DescribeTable("Value Tests",
		func(enum Financial_Options_ContractType, value string) {
			data, err := enum.Value()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data).Should(Equal(value))
		},
		Entry("Call - Works", Financial_Options_Call, "Call"),
		Entry("Put - Works", Financial_Options_Put, "Put"),
		Entry("Other - Works", Financial_Options_Other, "Other"))

	// Test that converting the Financial.Options.ContractType enum to a AttributeValue works for all values
	DescribeTable("MarshalDynamoDBAttributeValue Tests",
		func(enum Financial_Options_ContractType, value string) {
			data, err := enum.MarshalDynamoDBAttributeValue()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data.(*types.AttributeValueMemberS).Value).Should(Equal(value))
		},
		Entry("Call - Works", Financial_Options_Call, "Call"),
		Entry("Put - Works", Financial_Options_Put, "Put"),
		Entry("Other - Works", Financial_Options_Other, "Other"))

	// Test that attempting to deserialize a Financial.Options.ContractType will fail and
	// return an error if the value canno be deserialized from a JSON value to a string
	It("UnmarshalJSON fails - Error", func() {

		// Attempt to convert a non-parseable string value into a Financial.Options.ContractType
		// This should return an error
		enum := new(Financial_Options_ContractType)
		err := enum.UnmarshalJSON([]byte("derp"))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Options_ContractType"))
	})

	// Test that attempting to deserialize a Financial.Options.ContractType will fail and
	// return an error if the value cannot be converted to either the name value or integer
	// value of the enum option
	It("UnmarshalJSON - Value is invalid - Error", func() {

		// Attempt to convert a fake string value into a Financial.Options.ContractType
		// This should return an error
		enum := new(Financial_Options_ContractType)
		err := enum.UnmarshalJSON([]byte("\"derp\""))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Options_ContractType"))
	})

	// Test the conditions under which values should be convertible to a Financial.Options.ContractType
	DescribeTable("UnmarshalJSON Tests",
		func(value string, shouldBe Financial_Options_ContractType) {

			// Attempt to convert the string value into a Financial.Options.ContractType
			// This should not fail
			var enum Financial_Options_ContractType
			err := enum.UnmarshalJSON([]byte(value))

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("Call - Works", "\"Call\"", Financial_Options_Call),
		Entry("Put - Works", "\"Put\"", Financial_Options_Put),
		Entry("Other - Works", "\"Other\"", Financial_Options_Other),
		Entry("call - Works", "\"call\"", Financial_Options_Call),
		Entry("put - Works", "\"put\"", Financial_Options_Put),
		Entry("other - Works", "\"other\"", Financial_Options_Other),
		Entry("0 - Works", "\"0\"", Financial_Options_Call),
		Entry("1 - Works", "\"1\"", Financial_Options_Put),
		Entry("2 - Works", "\"2\"", Financial_Options_Other))

	// Test that attempting to deserialize a Financial.Options.ContractType will fial and return an
	// error if the value cannot be converted to either the name value or integer value
	// of the enum option
	It("UnmarshalCSV - Value is empty - Error", func() {

		// Attempt to convert a fake string value into a Financial.Options.ContractType
		// This should return an error
		enum := new(Financial_Options_ContractType)
		err := enum.UnmarshalCSV("")

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"\" cannot be mapped to a gopb.Financial_Options_ContractType"))
	})

	// Test the conditions under which values should be convertible to a Financial.Options.ContractType
	DescribeTable("UnmarshalCSV Tests",
		func(value string, shouldBe Financial_Options_ContractType) {

			// Attempt to convert the value into a Financial.Options.ContractType
			// This should not fail
			var enum Financial_Options_ContractType
			err := enum.UnmarshalCSV(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("Call - Works", "Call", Financial_Options_Call),
		Entry("Put - Works", "Put", Financial_Options_Put),
		Entry("Other - Works", "Other", Financial_Options_Other),
		Entry("call - Works", "call", Financial_Options_Call),
		Entry("put - Works", "put", Financial_Options_Put),
		Entry("other - Works", "other", Financial_Options_Other),
		Entry("0 - Works", "0", Financial_Options_Call),
		Entry("1 - Works", "1", Financial_Options_Put),
		Entry("2 - Works", "2", Financial_Options_Other))

	// Tests that, if the attribute type submitted to UnmarshalDynamoDBAttributeValue is not one we
	// recognize, then the function will return an error
	It("UnmarshalDynamoDBAttributeValue - AttributeValue type invalid - Error", func() {
		value := new(Financial_Options_ContractType)
		err := attributevalue.Unmarshal(&types.AttributeValueMemberBOOL{Value: true}, &value)
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("Attribute value of *types.AttributeValueMemberBOOL could not be converted to a Financial.Options.ContractType"))
	})

	// Tests the conditions under which UnmarshalDynamoDBAttributeValue is called and no error is generated
	DescribeTable("UnmarshalDynamoDBAttributeValue - AttributeValue Conditions",
		func(raw types.AttributeValue, expected Financial_Options_ContractType) {
			var value Financial_Options_ContractType
			err := attributevalue.Unmarshal(raw, &value)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(expected))
		},
		Entry("Value is []bytes, call - Works",
			&types.AttributeValueMemberB{Value: []byte("call")}, Financial_Options_Call),
		Entry("Value is []bytes, put - Works",
			&types.AttributeValueMemberB{Value: []byte("put")}, Financial_Options_Put),
		Entry("Value is []bytes, other - Works",
			&types.AttributeValueMemberB{Value: []byte("other")}, Financial_Options_Other),
		Entry("Value is []bytes, Call - Works",
			&types.AttributeValueMemberB{Value: []byte("Call")}, Financial_Options_Call),
		Entry("Value is []bytes, Put - Works",
			&types.AttributeValueMemberB{Value: []byte("Put")}, Financial_Options_Put),
		Entry("Value is []bytes, Other - Works",
			&types.AttributeValueMemberB{Value: []byte("Other")}, Financial_Options_Other),
		Entry("Value is numeric, 0 - Works",
			&types.AttributeValueMemberN{Value: "0"}, Financial_Options_Call),
		Entry("Value is numeric, 1 - Works",
			&types.AttributeValueMemberN{Value: "1"}, Financial_Options_Put),
		Entry("Value is numeric, 2 - Works",
			&types.AttributeValueMemberN{Value: "2"}, Financial_Options_Other),
		Entry("Value is NULL - Works", new(types.AttributeValueMemberNULL), Financial_Options_ContractType(0)),
		Entry("Value is string, call - Works",
			&types.AttributeValueMemberS{Value: "call"}, Financial_Options_Call),
		Entry("Value is string, put - Works",
			&types.AttributeValueMemberS{Value: "put"}, Financial_Options_Put),
		Entry("Value is string, other - Works",
			&types.AttributeValueMemberS{Value: "other"}, Financial_Options_Other),
		Entry("Value is string, Call - Works",
			&types.AttributeValueMemberS{Value: "Call"}, Financial_Options_Call),
		Entry("Value is string, Put - Works",
			&types.AttributeValueMemberS{Value: "Put"}, Financial_Options_Put),
		Entry("Value is string, Other - Works",
			&types.AttributeValueMemberS{Value: "Other"}, Financial_Options_Other))

	// Test that attempting to deserialize a Financial.Options.ContractType will fial and return an
	// error if the value cannot be converted to either the name value or integer value of the enum option
	It("Scan - Value is nil - Error", func() {

		// Attempt to convert a fake string value into a Financial.Options.ContractType
		// This should return an error
		var enum *Financial_Options_ContractType
		err := enum.Scan(nil)

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of %!q(<nil>) had an invalid type of <nil>"))
		Expect(enum).Should(BeNil())
	})

	// Test the conditions under which values should be convertible to a Financial.Options.ContractType
	DescribeTable("Scan Tests",
		func(value interface{}, shouldBe Financial_Options_ContractType) {

			// Attempt to convert the value into a Financial.Options.ContractType
			// This should not fail
			var enum Financial_Options_ContractType
			err := enum.Scan(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("Call - Works", "Call", Financial_Options_Call),
		Entry("Put - Works", "Put", Financial_Options_Put),
		Entry("Other - Works", "Other", Financial_Options_Other),
		Entry("call - Works", "call", Financial_Options_Call),
		Entry("put - Works", "put", Financial_Options_Put),
		Entry("other - Works", "other", Financial_Options_Other),
		Entry("0 - Works", 0, Financial_Options_Call),
		Entry("1 - Works", 1, Financial_Options_Put),
		Entry("2 - Works", 2, Financial_Options_Other))
})

var _ = Describe("Financial.Options.ExerciseStyle Marshal/Unmarshal Tests", func() {

	// Test that converting the Financial.Options.ExerciseStyle enum to JSON works for all values
	DescribeTable("MarshalJSON Tests",
		func(enum Financial_Options_ExerciseStyle, value string) {
			data, err := json.Marshal(enum)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("American - Works", Financial_Options_American, "\"American\""),
		Entry("European - Works", Financial_Options_European, "\"European\""),
		Entry("Bermudan - Works", Financial_Options_Bermudan, "\"Bermudan\""))

	// Test that converting the Financial.Options.ExerciseStyle enum to a CSV column works for all values
	DescribeTable("MarshalCSV Tests",
		func(enum Financial_Options_ExerciseStyle, value string) {
			data, err := enum.MarshalCSV()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("American - Works", Financial_Options_American, "0"),
		Entry("European - Works", Financial_Options_European, "1"),
		Entry("Bermudan - Works", Financial_Options_Bermudan, "2"))

	// Test that converting the Financial.Options.ExerciseStyle enum to an SQL value for all values
	DescribeTable("Value Tests",
		func(enum Financial_Options_ExerciseStyle, value string) {
			data, err := enum.Value()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data).Should(Equal(value))
		},
		Entry("American - Works", Financial_Options_American, "American"),
		Entry("European - Works", Financial_Options_European, "European"),
		Entry("Bermudan - Works", Financial_Options_Bermudan, "Bermudan"))

	// Test that converting the Financial.Options.ExerciseStyle enum to a AttributeValue works for all values
	DescribeTable("MarshalDynamoDBAttributeValue Tests",
		func(enum Financial_Options_ExerciseStyle, value string) {
			data, err := enum.MarshalDynamoDBAttributeValue()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data.(*types.AttributeValueMemberS).Value).Should(Equal(value))
		},
		Entry("American - Works", Financial_Options_American, "American"),
		Entry("European - Works", Financial_Options_European, "European"),
		Entry("Bermudan - Works", Financial_Options_Bermudan, "Bermudan"))

	// Test that attempting to deserialize a Financial.Options.ExerciseStyle will fail and
	// return an error if the value canno be deserialized from a JSON value to a string
	It("UnmarshalJSON fails - Error", func() {

		// Attempt to convert a non-parseable string value into a Financial.Options.ExerciseStyle
		// This should return an error
		enum := new(Financial_Options_ExerciseStyle)
		err := enum.UnmarshalJSON([]byte("derp"))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Options_ExerciseStyle"))
	})

	// Test that attempting to deserialize a Financial.Options.ExerciseStyle will fail and
	// return an error if the value cannot be converted to either the name value or integer
	// value of the enum option
	It("UnmarshalJSON - Value is invalid - Error", func() {

		// Attempt to convert a fake string value into a Financial.Options.ExerciseStyle
		// This should return an error
		enum := new(Financial_Options_ExerciseStyle)
		err := enum.UnmarshalJSON([]byte("\"derp\""))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Options_ExerciseStyle"))
	})

	// Test the conditions under which values should be convertible to a Financial.Options.ExerciseStyle
	DescribeTable("UnmarshalJSON Tests",
		func(value string, shouldBe Financial_Options_ExerciseStyle) {

			// Attempt to convert the string value into a Financial.Options.ExerciseStyle
			// This should not fail
			var enum Financial_Options_ExerciseStyle
			err := enum.UnmarshalJSON([]byte(value))

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("American - Works", "\"American\"", Financial_Options_American),
		Entry("European - Works", "\"European\"", Financial_Options_European),
		Entry("Bermudan - Works", "\"Bermudan\"", Financial_Options_Bermudan),
		Entry("american - Works", "\"american\"", Financial_Options_American),
		Entry("european - Works", "\"european\"", Financial_Options_European),
		Entry("bermudan - Works", "\"bermudan\"", Financial_Options_Bermudan),
		Entry("0 - Works", "\"0\"", Financial_Options_American),
		Entry("1 - Works", "\"1\"", Financial_Options_European),
		Entry("2 - Works", "\"2\"", Financial_Options_Bermudan))

	// Test that attempting to deserialize a Financial.Options.ExerciseStyle will fial and return an
	// error if the value cannot be converted to either the name value or integer value
	// of the enum option
	It("UnmarshalCSV - Value is empty - Error", func() {

		// Attempt to convert a fake string value into a Financial.Options.ExerciseStyle
		// This should return an error
		enum := new(Financial_Options_ExerciseStyle)
		err := enum.UnmarshalCSV("")

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"\" cannot be mapped to a gopb.Financial_Options_ExerciseStyle"))
	})

	// Test the conditions under which values should be convertible to a Financial.Options.ExerciseStyle
	DescribeTable("UnmarshalCSV Tests",
		func(value string, shouldBe Financial_Options_ExerciseStyle) {

			// Attempt to convert the value into a Financial.Options.ExerciseStyle
			// This should not fail
			var enum Financial_Options_ExerciseStyle
			err := enum.UnmarshalCSV(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("American - Works", "American", Financial_Options_American),
		Entry("European - Works", "European", Financial_Options_European),
		Entry("Bermudan - Works", "Bermudan", Financial_Options_Bermudan),
		Entry("american - Works", "american", Financial_Options_American),
		Entry("european - Works", "european", Financial_Options_European),
		Entry("bermudan - Works", "bermudan", Financial_Options_Bermudan),
		Entry("0 - Works", "0", Financial_Options_American),
		Entry("1 - Works", "1", Financial_Options_European),
		Entry("2 - Works", "2", Financial_Options_Bermudan))

	// Tests that, if the attribute type submitted to UnmarshalDynamoDBAttributeValue is not one we
	// recognize, then the function will return an error
	It("UnmarshalDynamoDBAttributeValue - AttributeValue type invalid - Error", func() {
		value := new(Financial_Options_ExerciseStyle)
		err := attributevalue.Unmarshal(&types.AttributeValueMemberBOOL{Value: true}, &value)
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("Attribute value of *types.AttributeValueMemberBOOL could not be converted to a Financial.Options.ExerciseStyle"))
	})

	// Tests the conditions under which UnmarshalDynamoDBAttributeValue is called and no error is generated
	DescribeTable("UnmarshalDynamoDBAttributeValue - AttributeValue Conditions",
		func(raw types.AttributeValue, expected Financial_Options_ExerciseStyle) {
			var value Financial_Options_ExerciseStyle
			err := attributevalue.Unmarshal(raw, &value)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(expected))
		},
		Entry("Value is []bytes, american - Works",
			&types.AttributeValueMemberB{Value: []byte("american")}, Financial_Options_American),
		Entry("Value is []bytes, european - Works",
			&types.AttributeValueMemberB{Value: []byte("european")}, Financial_Options_European),
		Entry("Value is []bytes, bermudan - Works",
			&types.AttributeValueMemberB{Value: []byte("bermudan")}, Financial_Options_Bermudan),
		Entry("Value is []bytes, American - Works",
			&types.AttributeValueMemberB{Value: []byte("American")}, Financial_Options_American),
		Entry("Value is []bytes, European - Works",
			&types.AttributeValueMemberB{Value: []byte("European")}, Financial_Options_European),
		Entry("Value is []bytes, Bermudan - Works",
			&types.AttributeValueMemberB{Value: []byte("Bermudan")}, Financial_Options_Bermudan),
		Entry("Value is numeric, 0 - Works",
			&types.AttributeValueMemberN{Value: "0"}, Financial_Options_American),
		Entry("Value is numeric, 1 - Works",
			&types.AttributeValueMemberN{Value: "1"}, Financial_Options_European),
		Entry("Value is numeric, 2 - Works",
			&types.AttributeValueMemberN{Value: "2"}, Financial_Options_Bermudan),
		Entry("Value is NULL - Works", new(types.AttributeValueMemberNULL), Financial_Options_ExerciseStyle(0)),
		Entry("Value is string, american - Works",
			&types.AttributeValueMemberS{Value: "american"}, Financial_Options_American),
		Entry("Value is string, european - Works",
			&types.AttributeValueMemberS{Value: "european"}, Financial_Options_European),
		Entry("Value is string, bermudan - Works",
			&types.AttributeValueMemberS{Value: "bermudan"}, Financial_Options_Bermudan),
		Entry("Value is string, American - Works",
			&types.AttributeValueMemberS{Value: "American"}, Financial_Options_American),
		Entry("Value is string, European - Works",
			&types.AttributeValueMemberS{Value: "European"}, Financial_Options_European),
		Entry("Value is string, Bermudan - Works",
			&types.AttributeValueMemberS{Value: "Bermudan"}, Financial_Options_Bermudan))

	// Test that attempting to deserialize a Financial.Options.ExerciseStyle will fial and return an
	// error if the value cannot be converted to either the name value or integer value of the enum option
	It("Scan - Value is nil - Error", func() {

		// Attempt to convert a fake string value into a Financial.Options.ExerciseStyle
		// This should return an error
		var enum *Financial_Options_ExerciseStyle
		err := enum.Scan(nil)

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of %!q(<nil>) had an invalid type of <nil>"))
		Expect(enum).Should(BeNil())
	})

	// Test the conditions under which values should be convertible to a Financial.Options.ExerciseStyle
	DescribeTable("Scan Tests",
		func(value interface{}, shouldBe Financial_Options_ExerciseStyle) {

			// Attempt to convert the value into a Financial.Options.ExerciseStyle
			// This should not fail
			var enum Financial_Options_ExerciseStyle
			err := enum.Scan(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("American - Works", "American", Financial_Options_American),
		Entry("European - Works", "European", Financial_Options_European),
		Entry("Bermudan - Works", "Bermudan", Financial_Options_Bermudan),
		Entry("american - Works", "american", Financial_Options_American),
		Entry("european - Works", "european", Financial_Options_European),
		Entry("bermudan - Works", "bermudan", Financial_Options_Bermudan),
		Entry("0 - Works", 0, Financial_Options_American),
		Entry("1 - Works", 1, Financial_Options_European),
		Entry("2 - Works", 2, Financial_Options_Bermudan))
})

var _ = Describe("Financial.Options.UnderlyingType Marshal/Unmarshal Tests", func() {

	// Test that converting the Financial.Options.UnderlyingType enum to JSON works for all values
	DescribeTable("MarshalJSON Tests",
		func(enum Financial_Options_UnderlyingType, value string) {
			data, err := json.Marshal(enum)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("Equity - Works", Financial_Options_Equity, "\"Equity\""),
		Entry("Currency - Works", Financial_Options_Currency, "\"Currency\""))

	// Test that converting the Financial.Options.UnderlyingType enum to a CSV column works for all values
	DescribeTable("MarshalCSV Tests",
		func(enum Financial_Options_UnderlyingType, value string) {
			data, err := enum.MarshalCSV()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("Equity - Works", Financial_Options_Equity, "0"),
		Entry("Currency - Works", Financial_Options_Currency, "1"))

	// Test that converting the Financial.Options.UnderlyingType enum to an SQL value for all values
	DescribeTable("Value Tests",
		func(enum Financial_Options_UnderlyingType, value string) {
			data, err := enum.Value()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data).Should(Equal(value))
		},
		Entry("Equity - Works", Financial_Options_Equity, "Equity"),
		Entry("Currency - Works", Financial_Options_Currency, "Currency"))

	// Test that converting the Financial.Options.UnderlyingType enum to a AttributeValue works for all values
	DescribeTable("MarshalDynamoDBAttributeValue Tests",
		func(enum Financial_Options_UnderlyingType, value string) {
			data, err := enum.MarshalDynamoDBAttributeValue()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data.(*types.AttributeValueMemberS).Value).Should(Equal(value))
		},
		Entry("Equity - Works", Financial_Options_Equity, "Equity"),
		Entry("Currency - Works", Financial_Options_Currency, "Currency"))

	// Test that attempting to deserialize a Financial.Options.UnderlyingType will fail and
	// return an error if the value canno be deserialized from a JSON value to a string
	It("UnmarshalJSON fails - Error", func() {

		// Attempt to convert a non-parseable string value into a Financial.Options.UnderlyingType
		// This should return an error
		enum := new(Financial_Options_UnderlyingType)
		err := enum.UnmarshalJSON([]byte("derp"))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Options_UnderlyingType"))
	})

	// Test that attempting to deserialize a Financial.Options.UnderlyingType will fail and
	// return an error if the value cannot be converted to either the name value or integer
	// value of the enum option
	It("UnmarshalJSON - Value is invalid - Error", func() {

		// Attempt to convert a fake string value into a Financial.Options.UnderlyingType
		// This should return an error
		enum := new(Financial_Options_UnderlyingType)
		err := enum.UnmarshalJSON([]byte("\"derp\""))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Options_UnderlyingType"))
	})

	// Test the conditions under which values should be convertible to a Financial.Options.UnderlyingType
	DescribeTable("UnmarshalJSON Tests",
		func(value string, shouldBe Financial_Options_UnderlyingType) {

			// Attempt to convert the string value into a Financial.Options.UnderlyingType
			// This should not fail
			var enum Financial_Options_UnderlyingType
			err := enum.UnmarshalJSON([]byte(value))

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("Equity - Works", "\"Equity\"", Financial_Options_Equity),
		Entry("Currency - Works", "\"Currency\"", Financial_Options_Currency),
		Entry("equity - Works", "\"equity\"", Financial_Options_Equity),
		Entry("currency - Works", "\"currency\"", Financial_Options_Currency),
		Entry("0 - Works", "\"0\"", Financial_Options_Equity),
		Entry("1 - Works", "\"1\"", Financial_Options_Currency))

	// Test that attempting to deserialize a Financial.Options.UnderlyingType will fial and return an
	// error if the value cannot be converted to either the name value or integer value
	// of the enum option
	It("UnmarshalCSV - Value is empty - Error", func() {

		// Attempt to convert a fake string value into a Financial.Options.UnderlyingType
		// This should return an error
		enum := new(Financial_Options_UnderlyingType)
		err := enum.UnmarshalCSV("")

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"\" cannot be mapped to a gopb.Financial_Options_UnderlyingType"))
	})

	// Test the conditions under which values should be convertible to a Financial.Options.UnderlyingType
	DescribeTable("UnmarshalCSV Tests",
		func(value string, shouldBe Financial_Options_UnderlyingType) {

			// Attempt to convert the value into a Financial.Options.UnderlyingType
			// This should not fail
			var enum Financial_Options_UnderlyingType
			err := enum.UnmarshalCSV(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("Equity - Works", "Equity", Financial_Options_Equity),
		Entry("Currency - Works", "Currency", Financial_Options_Currency),
		Entry("equity - Works", "equity", Financial_Options_Equity),
		Entry("currency - Works", "currency", Financial_Options_Currency),
		Entry("0 - Works", "0", Financial_Options_Equity),
		Entry("1 - Works", "1", Financial_Options_Currency))

	// Tests that, if the attribute type submitted to UnmarshalDynamoDBAttributeValue is not one we
	// recognize, then the function will return an error
	It("UnmarshalDynamoDBAttributeValue - AttributeValue type invalid - Error", func() {
		value := new(Financial_Options_UnderlyingType)
		err := attributevalue.Unmarshal(&types.AttributeValueMemberBOOL{Value: true}, &value)
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("Attribute value of *types.AttributeValueMemberBOOL could not be converted to a Financial.Options.UnderlyingType"))
	})

	// Tests the conditions under which UnmarshalDynamoDBAttributeValue is called and no error is generated
	DescribeTable("UnmarshalDynamoDBAttributeValue - AttributeValue Conditions",
		func(raw types.AttributeValue, expected Financial_Options_UnderlyingType) {
			var value Financial_Options_UnderlyingType
			err := attributevalue.Unmarshal(raw, &value)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(expected))
		},
		Entry("Value is []bytes, equity - Works",
			&types.AttributeValueMemberB{Value: []byte("equity")}, Financial_Options_Equity),
		Entry("Value is []bytes, currency - Works",
			&types.AttributeValueMemberB{Value: []byte("currency")}, Financial_Options_Currency),
		Entry("Value is []bytes, Equity - Works",
			&types.AttributeValueMemberB{Value: []byte("Equity")}, Financial_Options_Equity),
		Entry("Value is []bytes, Currency - Works",
			&types.AttributeValueMemberB{Value: []byte("Currency")}, Financial_Options_Currency),
		Entry("Value is numeric, 0 - Works",
			&types.AttributeValueMemberN{Value: "0"}, Financial_Options_Equity),
		Entry("Value is numeric, 1 - Works",
			&types.AttributeValueMemberN{Value: "1"}, Financial_Options_Currency),
		Entry("Value is NULL - Works", new(types.AttributeValueMemberNULL), Financial_Options_UnderlyingType(0)),
		Entry("Value is string, equity - Works",
			&types.AttributeValueMemberS{Value: "equity"}, Financial_Options_Equity),
		Entry("Value is string, currency - Works",
			&types.AttributeValueMemberS{Value: "currency"}, Financial_Options_Currency),
		Entry("Value is string, Equity - Works",
			&types.AttributeValueMemberS{Value: "Equity"}, Financial_Options_Equity),
		Entry("Value is string, Currency - Works",
			&types.AttributeValueMemberS{Value: "Currency"}, Financial_Options_Currency))

	// Test that attempting to deserialize a Financial.Options.UnderlyingType will fial and return an
	// error if the value cannot be converted to either the name value or integer value of the enum option
	It("Scan - Value is nil - Error", func() {

		// Attempt to convert a fake string value into a Financial.Options.UnderlyingType
		// This should return an error
		var enum *Financial_Options_UnderlyingType
		err := enum.Scan(nil)

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of %!q(<nil>) had an invalid type of <nil>"))
		Expect(enum).Should(BeNil())
	})

	// Test the conditions under which values should be convertible to a Financial.Options.UnderlyingType
	DescribeTable("Scan Tests",
		func(value interface{}, shouldBe Financial_Options_UnderlyingType) {

			// Attempt to convert the value into a Financial.Options.UnderlyingType
			// This should not fail
			var enum Financial_Options_UnderlyingType
			err := enum.Scan(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("Equity - Works", "Equity", Financial_Options_Equity),
		Entry("Currency - Works", "Currency", Financial_Options_Currency),
		Entry("equity - Works", "equity", Financial_Options_Equity),
		Entry("currency - Works", "currency", Financial_Options_Currency),
		Entry("0 - Works", 0, Financial_Options_Equity),
		Entry("1 - Works", 1, Financial_Options_Currency))
})

var _ = Describe("Financial.Quotes.Condition Marshal/Unmarshal Tests", func() {

	// Test that converting the Financial.Quotes.Condition enum to JSON works for all values
	DescribeTable("MarshalJSON Tests",
		func(enum Financial_Quotes_Condition, value string) {
			data, err := json.Marshal(enum)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("Regular - Works", Financial_Quotes_Regular, "\"Regular\""),
		Entry("RegularTwoSidedOpen - Works", Financial_Quotes_RegularTwoSidedOpen, "\"Regular, Two-Sided Open\""),
		Entry("RegularOneSidedOpen - Works", Financial_Quotes_RegularOneSidedOpen, "\"Regular, One-Sided Open\""),
		Entry("SlowAsk - Works", Financial_Quotes_SlowAsk, "\"Slow Ask\""),
		Entry("SlowBid - Works", Financial_Quotes_SlowBid, "\"Slow Bid\""),
		Entry("SlowBidAsk - Works", Financial_Quotes_SlowBidAsk, "\"Slow Bid, Ask\""),
		Entry("SlowDueLRPBid - Works", Financial_Quotes_SlowDueLRPBid, "\"Slow Due, LRP Bid\""),
		Entry("SlowDueLRPAsk - Works", Financial_Quotes_SlowDueLRPAsk, "\"Slow Due, LRP Ask\""),
		Entry("SlowDueNYSELRP - Works", Financial_Quotes_SlowDueNYSELRP, "\"Slow Due, NYSE LRP\""),
		Entry("SlowDueSetSlowListBidAsk - Works",
			Financial_Quotes_SlowDueSetSlowListBidAsk, "\"Slow Due Set, Slow List, Bid, Ask\""),
		Entry("ManualAskAutomatedBid - Works", Financial_Quotes_ManualAskAutomatedBid, "\"Manual Ask, Automated Bid\""),
		Entry("ManualBidAutomatedAsk - Works", Financial_Quotes_ManualBidAutomatedAsk, "\"Manual Bid, Automated Ask\""),
		Entry("ManualBidAndAsk - Works", Financial_Quotes_ManualBidAndAsk, "\"Manual Bid and Ask\""),
		Entry("Opening - Works", Financial_Quotes_Opening, "\"Opening\""),
		Entry("Closing - Works", Financial_Quotes_Closing, "\"Closing\""),
		Entry("Closed - Works", Financial_Quotes_Closed, "\"Closed\""),
		Entry("Resume - Works", Financial_Quotes_Resume, "\"Resume\""),
		Entry("FastTrading - Works", Financial_Quotes_FastTrading, "\"Fast Trading\""),
		Entry("TradingRangeIndicated - Works", Financial_Quotes_TradingRangeIndicated, "\"Tading Range Indicated\""),
		Entry("MarketMakerQuotesClosed - Works", Financial_Quotes_MarketMakerQuotesClosed, "\"Market-Maker Quotes Closed\""),
		Entry("NonFirm - Works", Financial_Quotes_NonFirm, "\"Non-Firm\""),
		Entry("NewsDissemination - Works", Financial_Quotes_NewsDissemination, "\"News Dissemination\""),
		Entry("OrderInflux - Works", Financial_Quotes_OrderInflux, "\"Order Influx\""),
		Entry("OrderImbalance - Works", Financial_Quotes_OrderImbalance, "\"Order Imbalance\""),
		Entry("DueToRelatedSecurityNewsDissemination - Works",
			Financial_Quotes_DueToRelatedSecurityNewsDissemination, "\"Due to Related Security, News Dissemination\""),
		Entry("DueToRelatedSecurityNewsPending - Works",
			Financial_Quotes_DueToRelatedSecurityNewsPending, "\"Due to Related Security, News Pending\""),
		Entry("AdditionalInformation - Works", Financial_Quotes_AdditionalInformation, "\"Additional Information\""),
		Entry("NewsPending - Works", Financial_Quotes_NewsPending, "\"News Pending\""),
		Entry("AdditionalInformationDueToRelatedSecurity - Works",
			Financial_Quotes_AdditionalInformationDueToRelatedSecurity, "\"Additional Information Due to Related Security\""),
		Entry("DueToRelatedSecurity - Works", Financial_Quotes_DueToRelatedSecurity, "\"Due to Related Security\""),
		Entry("InViewOfCommon - Works", Financial_Quotes_InViewOfCommon, "\"In View of Common\""),
		Entry("EquipmentChangeover - Works", Financial_Quotes_EquipmentChangeover, "\"Equipment Changeover\""),
		Entry("NoOpenNoResponse - Works", Financial_Quotes_NoOpenNoResponse, "\"No Open, No Response\""),
		Entry("SubPennyTrading - Works", Financial_Quotes_SubPennyTrading, "\"Sub-Penny Trading\""),
		Entry("AutomatedBidNoOfferNoBid - Works",
			Financial_Quotes_AutomatedBidNoOfferNoBid, "\"Automated Bid; No Offer, No Bid\""),
		Entry("LULDPPriceBand - Works", Financial_Quotes_LULDPPriceBand, "\"LULD Price Band\""),
		Entry("MarketWideCircuitBreakerLevel1 - Works",
			Financial_Quotes_MarketWideCircuitBreakerLevel1, "\"Market-Wide Circuit Breaker, Level 1\""),
		Entry("MarketWideCircuitBreakerLevel2 - Works",
			Financial_Quotes_MarketWideCircuitBreakerLevel2, "\"Market-Wide Circuit Breaker, Level 2\""),
		Entry("MarketWideCircuitBreakerLevel3 - Works",
			Financial_Quotes_MarketWideCircuitBreakerLevel3, "\"Market-Wide Circuit Breaker, Level 3\""),
		Entry("RepublishedLULDPriceBand - Works", Financial_Quotes_RepublishedLULDPriceBand, "\"Republished LULD Price Band\""),
		Entry("OnDemandAuction - Works", Financial_Quotes_OnDemandAuction, "\"On-Demand Auction\""),
		Entry("CashOnlySettlement - Works", Financial_Quotes_CashOnlySettlement, "\"Cash-Only Settlement\""),
		Entry("NextDaySettlement - Works", Financial_Quotes_NextDaySettlement, "\"Next-Day Settlement\""),
		Entry("LULDTradingPause - Works", Financial_Quotes_LULDTradingPause, "\"LULD Trading Pause\""),
		Entry("SlowDuelRPBidAsk - Works", Financial_Quotes_SlowDuelRPBidAsk, "\"Slow Due LRP, Bid, Ask\""),
		Entry("Cancel - Works", Financial_Quotes_Cancel, "\"Cancel\""),
		Entry("CorrectedPrice - Works", Financial_Quotes_CorrectedPrice, "\"Corrected Price\""),
		Entry("SIPGenerated - Works", Financial_Quotes_SIPGenerated, "\"SIP-Generated\""),
		Entry("Unknown - Works", Financial_Quotes_Unknown, "\"Unknown\""),
		Entry("CrossedMarket - Works", Financial_Quotes_CrossedMarket, "\"Crossed Market\""),
		Entry("LockedMarket - Works", Financial_Quotes_LockedMarket, "\"Locked Market\""),
		Entry("DepthOnOfferSide - Works", Financial_Quotes_DepthOnOfferSide, "\"Depth on Offer Side\""),
		Entry("DepthOnBidSide - Works", Financial_Quotes_DepthOnBidSide, "\"Depth on Bid Side\""),
		Entry("DepthOnBidAndOffer - Works", Financial_Quotes_DepthOnBidAndOffer, "\"Depth on Bid and Offer\""),
		Entry("PreOpeningIndication - Works", Financial_Quotes_PreOpeningIndication, "\"Pre-Opening Indication\""),
		Entry("SyndicateBid - Works", Financial_Quotes_SyndicateBid, "\"Syndicate Bid\""),
		Entry("PreSyndicateBid - Works", Financial_Quotes_PreSyndicateBid, "\"Pre-Syndicate Bid\""),
		Entry("PenaltyBid - Works", Financial_Quotes_PenaltyBid, "\"Penalty Bid\""),
		Entry("CQSGenerated - Works", Financial_Quotes_CQSGenerated, "\"CQS-Generated\""),
		Entry("Invalid - Works", Financial_Quotes_Invalid, "\"Invalid\""))

	// Test that converting the Financial.Quotes.Condition enum to a CSV column works for all values
	DescribeTable("MarshalCSV Tests",
		func(enum Financial_Quotes_Condition, value string) {
			data, err := enum.MarshalCSV()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("Regular - Works", Financial_Quotes_Regular, "000"),
		Entry("RegularTwoSidedOpen - Works", Financial_Quotes_RegularTwoSidedOpen, "001"),
		Entry("RegularOneSidedOpen - Works", Financial_Quotes_RegularOneSidedOpen, "002"),
		Entry("SlowAsk - Works", Financial_Quotes_SlowAsk, "003"),
		Entry("SlowBid - Works", Financial_Quotes_SlowBid, "004"),
		Entry("SlowBidAsk - Works", Financial_Quotes_SlowBidAsk, "005"),
		Entry("SlowDueLRPBid - Works", Financial_Quotes_SlowDueLRPBid, "006"),
		Entry("SlowDueLRPAsk - Works", Financial_Quotes_SlowDueLRPAsk, "007"),
		Entry("SlowDueNYSELRP - Works", Financial_Quotes_SlowDueNYSELRP, "008"),
		Entry("SlowDueSetSlowListBidAsk - Works", Financial_Quotes_SlowDueSetSlowListBidAsk, "009"),
		Entry("ManualAskAutomatedBid - Works", Financial_Quotes_ManualAskAutomatedBid, "010"),
		Entry("ManualBidAutomatedAsk - Works", Financial_Quotes_ManualBidAutomatedAsk, "011"),
		Entry("ManualBidAndAsk - Works", Financial_Quotes_ManualBidAndAsk, "012"),
		Entry("Opening - Works", Financial_Quotes_Opening, "013"),
		Entry("Closing - Works", Financial_Quotes_Closing, "014"),
		Entry("Closed - Works", Financial_Quotes_Closed, "015"),
		Entry("Resume - Works", Financial_Quotes_Resume, "016"),
		Entry("FastTrading - Works", Financial_Quotes_FastTrading, "017"),
		Entry("TradingRangeIndicated - Works", Financial_Quotes_TradingRangeIndicated, "018"),
		Entry("MarketMakerQuotesClosed - Works", Financial_Quotes_MarketMakerQuotesClosed, "019"),
		Entry("NonFirm - Works", Financial_Quotes_NonFirm, "020"),
		Entry("NewsDissemination - Works", Financial_Quotes_NewsDissemination, "021"),
		Entry("OrderInflux - Works", Financial_Quotes_OrderInflux, "022"),
		Entry("OrderImbalance - Works", Financial_Quotes_OrderImbalance, "023"),
		Entry("DueToRelatedSecurityNewsDissemination - Works", Financial_Quotes_DueToRelatedSecurityNewsDissemination, "024"),
		Entry("DueToRelatedSecurityNewsPending - Works", Financial_Quotes_DueToRelatedSecurityNewsPending, "025"),
		Entry("AdditionalInformation - Works", Financial_Quotes_AdditionalInformation, "026"),
		Entry("NewsPending - Works", Financial_Quotes_NewsPending, "027"),
		Entry("AdditionalInformationDueToRelatedSecurity - Works",
			Financial_Quotes_AdditionalInformationDueToRelatedSecurity, "028"),
		Entry("DueToRelatedSecurity - Works", Financial_Quotes_DueToRelatedSecurity, "029"),
		Entry("InViewOfCommon - Works", Financial_Quotes_InViewOfCommon, "030"),
		Entry("EquipmentChangeover - Works", Financial_Quotes_EquipmentChangeover, "031"),
		Entry("NoOpenNoResponse - Works", Financial_Quotes_NoOpenNoResponse, "032"),
		Entry("SubPennyTrading - Works", Financial_Quotes_SubPennyTrading, "033"),
		Entry("AutomatedBidNoOfferNoBid - Works", Financial_Quotes_AutomatedBidNoOfferNoBid, "034"),
		Entry("LULDPPriceBand - Works", Financial_Quotes_LULDPPriceBand, "035"),
		Entry("MarketWideCircuitBreakerLevel1 - Works", Financial_Quotes_MarketWideCircuitBreakerLevel1, "036"),
		Entry("MarketWideCircuitBreakerLevel2 - Works", Financial_Quotes_MarketWideCircuitBreakerLevel2, "037"),
		Entry("MarketWideCircuitBreakerLevel3 - Works", Financial_Quotes_MarketWideCircuitBreakerLevel3, "038"),
		Entry("RepublishedLULDPriceBand - Works", Financial_Quotes_RepublishedLULDPriceBand, "039"),
		Entry("OnDemandAuction - Works", Financial_Quotes_OnDemandAuction, "040"),
		Entry("CashOnlySettlement - Works", Financial_Quotes_CashOnlySettlement, "041"),
		Entry("NextDaySettlement - Works", Financial_Quotes_NextDaySettlement, "042"),
		Entry("LULDTradingPause - Works", Financial_Quotes_LULDTradingPause, "043"),
		Entry("SlowDuelRPBidAsk - Works", Financial_Quotes_SlowDuelRPBidAsk, "071"),
		Entry("Cancel - Works", Financial_Quotes_Cancel, "080"),
		Entry("CorrectedPrice - Works", Financial_Quotes_CorrectedPrice, "081"),
		Entry("SIPGenerated - Works", Financial_Quotes_SIPGenerated, "082"),
		Entry("Unknown - Works", Financial_Quotes_Unknown, "083"),
		Entry("CrossedMarket - Works", Financial_Quotes_CrossedMarket, "084"),
		Entry("LockedMarket - Works", Financial_Quotes_LockedMarket, "085"),
		Entry("DepthOnOfferSide - Works", Financial_Quotes_DepthOnOfferSide, "086"),
		Entry("DepthOnBidSide - Works", Financial_Quotes_DepthOnBidSide, "087"),
		Entry("DepthOnBidAndOffer - Works", Financial_Quotes_DepthOnBidAndOffer, "088"),
		Entry("PreOpeningIndication - Works", Financial_Quotes_PreOpeningIndication, "089"),
		Entry("SyndicateBid - Works", Financial_Quotes_SyndicateBid, "090"),
		Entry("PreSyndicateBid - Works", Financial_Quotes_PreSyndicateBid, "091"),
		Entry("PenaltyBid - Works", Financial_Quotes_PenaltyBid, "092"),
		Entry("CQSGenerated - Works", Financial_Quotes_CQSGenerated, "094"),
		Entry("Invalid - Works", Financial_Quotes_Invalid, "999"))

	// Test that converting the Financial.Quotes.Condition enum to an SQL value for all values
	DescribeTable("Value Tests",
		func(enum Financial_Quotes_Condition, value int) {
			data, err := enum.Value()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data).Should(Equal(value))
		},
		Entry("Regular - Works", Financial_Quotes_Regular, 0),
		Entry("RegularTwoSidedOpen - Works", Financial_Quotes_RegularTwoSidedOpen, 1),
		Entry("RegularOneSidedOpen - Works", Financial_Quotes_RegularOneSidedOpen, 2),
		Entry("SlowAsk - Works", Financial_Quotes_SlowAsk, 3),
		Entry("SlowBid - Works", Financial_Quotes_SlowBid, 4),
		Entry("SlowBidAsk - Works", Financial_Quotes_SlowBidAsk, 5),
		Entry("SlowDueLRPBid - Works", Financial_Quotes_SlowDueLRPBid, 6),
		Entry("SlowDueLRPAsk - Works", Financial_Quotes_SlowDueLRPAsk, 7),
		Entry("SlowDueNYSELRP - Works", Financial_Quotes_SlowDueNYSELRP, 8),
		Entry("SlowDueSetSlowListBidAsk - Works", Financial_Quotes_SlowDueSetSlowListBidAsk, 9),
		Entry("ManualAskAutomatedBid - Works", Financial_Quotes_ManualAskAutomatedBid, 10),
		Entry("ManualBidAutomatedAsk - Works", Financial_Quotes_ManualBidAutomatedAsk, 11),
		Entry("ManualBidAndAsk - Works", Financial_Quotes_ManualBidAndAsk, 12),
		Entry("Opening - Works", Financial_Quotes_Opening, 13),
		Entry("Closing - Works", Financial_Quotes_Closing, 14),
		Entry("Closed - Works", Financial_Quotes_Closed, 15),
		Entry("Resume - Works", Financial_Quotes_Resume, 16),
		Entry("FastTrading - Works", Financial_Quotes_FastTrading, 17),
		Entry("TradingRangeIndicated - Works", Financial_Quotes_TradingRangeIndicated, 18),
		Entry("MarketMakerQuotesClosed - Works", Financial_Quotes_MarketMakerQuotesClosed, 19),
		Entry("NonFirm - Works", Financial_Quotes_NonFirm, 20),
		Entry("NewsDissemination - Works", Financial_Quotes_NewsDissemination, 21),
		Entry("OrderInflux - Works", Financial_Quotes_OrderInflux, 22),
		Entry("OrderImbalance - Works", Financial_Quotes_OrderImbalance, 23),
		Entry("DueToRelatedSecurityNewsDissemination - Works", Financial_Quotes_DueToRelatedSecurityNewsDissemination, 24),
		Entry("DueToRelatedSecurityNewsPending - Works", Financial_Quotes_DueToRelatedSecurityNewsPending, 25),
		Entry("AdditionalInformation - Works", Financial_Quotes_AdditionalInformation, 26),
		Entry("NewsPending - Works", Financial_Quotes_NewsPending, 27),
		Entry("AdditionalInformationDueToRelatedSecurity - Works", Financial_Quotes_AdditionalInformationDueToRelatedSecurity, 28),
		Entry("DueToRelatedSecurity - Works", Financial_Quotes_DueToRelatedSecurity, 29),
		Entry("InViewOfCommon - Works", Financial_Quotes_InViewOfCommon, 30),
		Entry("EquipmentChangeover - Works", Financial_Quotes_EquipmentChangeover, 31),
		Entry("NoOpenNoResponse - Works", Financial_Quotes_NoOpenNoResponse, 32),
		Entry("SubPennyTrading - Works", Financial_Quotes_SubPennyTrading, 33),
		Entry("AutomatedBidNoOfferNoBid - Works", Financial_Quotes_AutomatedBidNoOfferNoBid, 34),
		Entry("LULDPPriceBand - Works", Financial_Quotes_LULDPPriceBand, 35),
		Entry("MarketWideCircuitBreakerLevel1 - Works", Financial_Quotes_MarketWideCircuitBreakerLevel1, 36),
		Entry("MarketWideCircuitBreakerLevel2 - Works", Financial_Quotes_MarketWideCircuitBreakerLevel2, 37),
		Entry("MarketWideCircuitBreakerLevel3 - Works", Financial_Quotes_MarketWideCircuitBreakerLevel3, 38),
		Entry("RepublishedLULDPriceBand - Works", Financial_Quotes_RepublishedLULDPriceBand, 39),
		Entry("OnDemandAuction - Works", Financial_Quotes_OnDemandAuction, 40),
		Entry("CashOnlySettlement - Works", Financial_Quotes_CashOnlySettlement, 41),
		Entry("NextDaySettlement - Works", Financial_Quotes_NextDaySettlement, 42),
		Entry("LULDTradingPause - Works", Financial_Quotes_LULDTradingPause, 43),
		Entry("SlowDuelRPBidAsk - Works", Financial_Quotes_SlowDuelRPBidAsk, 71),
		Entry("Cancel - Works", Financial_Quotes_Cancel, 80),
		Entry("CorrectedPrice - Works", Financial_Quotes_CorrectedPrice, 81),
		Entry("SIPGenerated - Works", Financial_Quotes_SIPGenerated, 82),
		Entry("Unknown - Works", Financial_Quotes_Unknown, 83),
		Entry("CrossedMarket - Works", Financial_Quotes_CrossedMarket, 84),
		Entry("LockedMarket - Works", Financial_Quotes_LockedMarket, 85),
		Entry("DepthOnOfferSide - Works", Financial_Quotes_DepthOnOfferSide, 86),
		Entry("DepthOnBidSide - Works", Financial_Quotes_DepthOnBidSide, 87),
		Entry("DepthOnBidAndOffer - Works", Financial_Quotes_DepthOnBidAndOffer, 88),
		Entry("PreOpeningIndication - Works", Financial_Quotes_PreOpeningIndication, 89),
		Entry("SyndicateBid - Works", Financial_Quotes_SyndicateBid, 90),
		Entry("PreSyndicateBid - Works", Financial_Quotes_PreSyndicateBid, 91),
		Entry("PenaltyBid - Works", Financial_Quotes_PenaltyBid, 92),
		Entry("CQSGenerated - Works", Financial_Quotes_CQSGenerated, 94),
		Entry("Invalid - Works", Financial_Quotes_Invalid, -1))

	// Test that converting the Financial.Quotes.Condition enum to a AttributeValue works for all values
	DescribeTable("MarshalDynamoDBAttributeValue Tests",
		func(enum Financial_Quotes_Condition, value string) {
			data, err := enum.MarshalDynamoDBAttributeValue()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data.(*types.AttributeValueMemberS).Value).Should(Equal(value))
		},
		Entry("Regular - Works", Financial_Quotes_Regular, "Regular"),
		Entry("RegularTwoSidedOpen - Works", Financial_Quotes_RegularTwoSidedOpen, "Regular, Two-Sided Open"),
		Entry("RegularOneSidedOpen - Works", Financial_Quotes_RegularOneSidedOpen, "Regular, One-Sided Open"),
		Entry("SlowAsk - Works", Financial_Quotes_SlowAsk, "Slow Ask"),
		Entry("SlowBid - Works", Financial_Quotes_SlowBid, "Slow Bid"),
		Entry("SlowBidAsk - Works", Financial_Quotes_SlowBidAsk, "Slow Bid, Ask"),
		Entry("SlowDueLRPBid - Works", Financial_Quotes_SlowDueLRPBid, "Slow Due, LRP Bid"),
		Entry("SlowDueLRPAsk - Works", Financial_Quotes_SlowDueLRPAsk, "Slow Due, LRP Ask"),
		Entry("SlowDueNYSELRP - Works", Financial_Quotes_SlowDueNYSELRP, "Slow Due, NYSE LRP"),
		Entry("SlowDueSetSlowListBidAsk - Works",
			Financial_Quotes_SlowDueSetSlowListBidAsk, "Slow Due Set, Slow List, Bid, Ask"),
		Entry("ManualAskAutomatedBid - Works", Financial_Quotes_ManualAskAutomatedBid, "Manual Ask, Automated Bid"),
		Entry("ManualBidAutomatedAsk - Works", Financial_Quotes_ManualBidAutomatedAsk, "Manual Bid, Automated Ask"),
		Entry("ManualBidAndAsk - Works", Financial_Quotes_ManualBidAndAsk, "Manual Bid and Ask"),
		Entry("Opening - Works", Financial_Quotes_Opening, "Opening"),
		Entry("Closing - Works", Financial_Quotes_Closing, "Closing"),
		Entry("Closed - Works", Financial_Quotes_Closed, "Closed"),
		Entry("Resume - Works", Financial_Quotes_Resume, "Resume"),
		Entry("FastTrading - Works", Financial_Quotes_FastTrading, "Fast Trading"),
		Entry("TradingRangeIndicated - Works", Financial_Quotes_TradingRangeIndicated, "Tading Range Indicated"),
		Entry("MarketMakerQuotesClosed - Works", Financial_Quotes_MarketMakerQuotesClosed, "Market-Maker Quotes Closed"),
		Entry("NonFirm - Works", Financial_Quotes_NonFirm, "Non-Firm"),
		Entry("NewsDissemination - Works", Financial_Quotes_NewsDissemination, "News Dissemination"),
		Entry("OrderInflux - Works", Financial_Quotes_OrderInflux, "Order Influx"),
		Entry("OrderImbalance - Works", Financial_Quotes_OrderImbalance, "Order Imbalance"),
		Entry("DueToRelatedSecurityNewsDissemination - Works",
			Financial_Quotes_DueToRelatedSecurityNewsDissemination, "Due to Related Security, News Dissemination"),
		Entry("DueToRelatedSecurityNewsPending - Works",
			Financial_Quotes_DueToRelatedSecurityNewsPending, "Due to Related Security, News Pending"),
		Entry("AdditionalInformation - Works", Financial_Quotes_AdditionalInformation, "Additional Information"),
		Entry("NewsPending - Works", Financial_Quotes_NewsPending, "News Pending"),
		Entry("AdditionalInformationDueToRelatedSecurity - Works",
			Financial_Quotes_AdditionalInformationDueToRelatedSecurity, "Additional Information Due to Related Security"),
		Entry("DueToRelatedSecurity - Works", Financial_Quotes_DueToRelatedSecurity, "Due to Related Security"),
		Entry("InViewOfCommon - Works", Financial_Quotes_InViewOfCommon, "In View of Common"),
		Entry("EquipmentChangeover - Works", Financial_Quotes_EquipmentChangeover, "Equipment Changeover"),
		Entry("NoOpenNoResponse - Works", Financial_Quotes_NoOpenNoResponse, "No Open, No Response"),
		Entry("SubPennyTrading - Works", Financial_Quotes_SubPennyTrading, "Sub-Penny Trading"),
		Entry("AutomatedBidNoOfferNoBid - Works", Financial_Quotes_AutomatedBidNoOfferNoBid, "Automated Bid; No Offer, No Bid"),
		Entry("LULDPPriceBand - Works", Financial_Quotes_LULDPPriceBand, "LULD Price Band"),
		Entry("MarketWideCircuitBreakerLevel1 - Works",
			Financial_Quotes_MarketWideCircuitBreakerLevel1, "Market-Wide Circuit Breaker, Level 1"),
		Entry("MarketWideCircuitBreakerLevel2 - Works",
			Financial_Quotes_MarketWideCircuitBreakerLevel2, "Market-Wide Circuit Breaker, Level 2"),
		Entry("MarketWideCircuitBreakerLevel3 - Works",
			Financial_Quotes_MarketWideCircuitBreakerLevel3, "Market-Wide Circuit Breaker, Level 3"),
		Entry("RepublishedLULDPriceBand - Works", Financial_Quotes_RepublishedLULDPriceBand, "Republished LULD Price Band"),
		Entry("OnDemandAuction - Works", Financial_Quotes_OnDemandAuction, "On-Demand Auction"),
		Entry("CashOnlySettlement - Works", Financial_Quotes_CashOnlySettlement, "Cash-Only Settlement"),
		Entry("NextDaySettlement - Works", Financial_Quotes_NextDaySettlement, "Next-Day Settlement"),
		Entry("LULDTradingPause - Works", Financial_Quotes_LULDTradingPause, "LULD Trading Pause"),
		Entry("SlowDuelRPBidAsk - Works", Financial_Quotes_SlowDuelRPBidAsk, "Slow Due LRP, Bid, Ask"),
		Entry("Cancel - Works", Financial_Quotes_Cancel, "Cancel"),
		Entry("CorrectedPrice - Works", Financial_Quotes_CorrectedPrice, "Corrected Price"),
		Entry("SIPGenerated - Works", Financial_Quotes_SIPGenerated, "SIP-Generated"),
		Entry("Unknown - Works", Financial_Quotes_Unknown, "Unknown"),
		Entry("CrossedMarket - Works", Financial_Quotes_CrossedMarket, "Crossed Market"),
		Entry("LockedMarket - Works", Financial_Quotes_LockedMarket, "Locked Market"),
		Entry("DepthOnOfferSide - Works", Financial_Quotes_DepthOnOfferSide, "Depth on Offer Side"),
		Entry("DepthOnBidSide - Works", Financial_Quotes_DepthOnBidSide, "Depth on Bid Side"),
		Entry("DepthOnBidAndOffer - Works", Financial_Quotes_DepthOnBidAndOffer, "Depth on Bid and Offer"),
		Entry("PreOpeningIndication - Works", Financial_Quotes_PreOpeningIndication, "Pre-Opening Indication"),
		Entry("SyndicateBid - Works", Financial_Quotes_SyndicateBid, "Syndicate Bid"),
		Entry("PreSyndicateBid - Works", Financial_Quotes_PreSyndicateBid, "Pre-Syndicate Bid"),
		Entry("PenaltyBid - Works", Financial_Quotes_PenaltyBid, "Penalty Bid"),
		Entry("CQSGenerated - Works", Financial_Quotes_CQSGenerated, "CQS-Generated"),
		Entry("Invalid - Works", Financial_Quotes_Invalid, "Invalid"))

	// Test that attempting to deserialize a Financial.Quotes.Condition will fail and
	// return an error if the value canno be deserialized from a JSON value to a string
	It("UnmarshalJSON fails - Error", func() {

		// Attempt to convert a non-parseable string value into a Financial.Quotes.Condition
		// This should return an error
		enum := new(Financial_Quotes_Condition)
		err := enum.UnmarshalJSON([]byte("derp"))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Quotes_Condition"))
	})

	// Test that attempting to deserialize a Financial.Quotes.Condition will fail and
	// return an error if the value cannot be converted to either the name value or integer
	// value of the enum option
	It("UnmarshalJSON - Value is invalid - Error", func() {

		// Attempt to convert a fake string value into a Financial.Quotes.Condition
		// This should return an error
		enum := new(Financial_Quotes_Condition)
		err := enum.UnmarshalJSON([]byte("\"derp\""))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Quotes_Condition"))
	})

	// Test the conditions under which values should be convertible to a Financial.Quotes.Condition
	DescribeTable("UnmarshalJSON Tests",
		func(value interface{}, shouldBe Financial_Quotes_Condition) {

			// Attempt to convert the string value into a Financial.Quotes.Condition
			// This should not fail
			var enum Financial_Quotes_Condition
			err := enum.UnmarshalJSON([]byte(fmt.Sprintf("%v", value)))

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("Regular, Two-Sided Open - Works", "\"Regular, Two-Sided Open\"", Financial_Quotes_RegularTwoSidedOpen),
		Entry("Regular, One-Sided Open - Works", "\"Regular, One-Sided Open\"", Financial_Quotes_RegularOneSidedOpen),
		Entry("Slow Ask - Works", "\"Slow Ask\"", Financial_Quotes_SlowAsk),
		Entry("Slow Bid - Works", "\"Slow Bid\"", Financial_Quotes_SlowBid),
		Entry("Slow Bid, Ask - Works", "\"Slow Bid, Ask\"", Financial_Quotes_SlowBidAsk),
		Entry("Slow Due, LRP Bid - Works", "\"Slow Due, LRP Bid\"", Financial_Quotes_SlowDueLRPBid),
		Entry("Slow Due, LRP Ask - Works", "\"Slow Due, LRP Ask\"", Financial_Quotes_SlowDueLRPAsk),
		Entry("Slow Due, NYSE LRP - Works", "\"Slow Due, NYSE LRP\"", Financial_Quotes_SlowDueNYSELRP),
		Entry("Slow Due Set, Slow List, Bid, Ask - Works",
			"\"Slow Due Set, Slow List, Bid, Ask\"", Financial_Quotes_SlowDueSetSlowListBidAsk),
		Entry("Manual Ask, Automated Bid - Works", "\"Manual Ask, Automated Bid\"", Financial_Quotes_ManualAskAutomatedBid),
		Entry("Manual Bid, Automated Ask - Works", "\"Manual Bid, Automated Ask\"", Financial_Quotes_ManualBidAutomatedAsk),
		Entry("Manual Bid and Ask - Works", "\"Manual Bid and Ask\"", Financial_Quotes_ManualBidAndAsk),
		Entry("Fast Trading - Works", "\"Fast Trading\"", Financial_Quotes_FastTrading),
		Entry("Tading Range Indicated - Works", "\"Tading Range Indicated\"", Financial_Quotes_TradingRangeIndicated),
		Entry("Market-Maker Quotes Closed - Works", "\"Market-Maker Quotes Closed\"", Financial_Quotes_MarketMakerQuotesClosed),
		Entry("Non-Firm - Works", "\"Non-Firm\"", Financial_Quotes_NonFirm),
		Entry("News Dissemination - Works", "\"News Dissemination\"", Financial_Quotes_NewsDissemination),
		Entry("Order Influx - Works", "\"Order Influx\"", Financial_Quotes_OrderInflux),
		Entry("Order Imbalance - Works", "\"Order Imbalance\"", Financial_Quotes_OrderImbalance),
		Entry("Due to Related Security, News Dissemination - Works",
			"\"Due to Related Security, News Dissemination\"", Financial_Quotes_DueToRelatedSecurityNewsDissemination),
		Entry("Due to Related Security, News Pending - Works",
			"\"Due to Related Security, News Pending\"", Financial_Quotes_DueToRelatedSecurityNewsPending),
		Entry("Additional Information - Works", "\"Additional Information\"", Financial_Quotes_AdditionalInformation),
		Entry("News Pending - Works", "\"News Pending\"", Financial_Quotes_NewsPending),
		Entry("Additional Information Due to Related Security - Works",
			"\"Additional Information Due to Related Security\"", Financial_Quotes_AdditionalInformationDueToRelatedSecurity),
		Entry("Due to Related Security - Works", "\"Due to Related Security\"", Financial_Quotes_DueToRelatedSecurity),
		Entry("In View of Common - Works", "\"In View of Common\"", Financial_Quotes_InViewOfCommon),
		Entry("Equipment Changeover - Works", "\"Equipment Changeover\"", Financial_Quotes_EquipmentChangeover),
		Entry("No Open, No Response - Works", "\"No Open, No Response\"", Financial_Quotes_NoOpenNoResponse),
		Entry("Sub-Penny Trading - Works", "\"Sub-Penny Trading\"", Financial_Quotes_SubPennyTrading),
		Entry("Automated Bid; No Offer, No Bid - Works",
			"\"Automated Bid; No Offer, No Bid\"", Financial_Quotes_AutomatedBidNoOfferNoBid),
		Entry("LULD Price Band - Works", "\"LULD Price Band\"", Financial_Quotes_LULDPPriceBand),
		Entry("Market-Wide Circuit Breaker, Level 1 - Works",
			"\"Market-Wide Circuit Breaker, Level 1\"", Financial_Quotes_MarketWideCircuitBreakerLevel1),
		Entry("Market-Wide Circuit Breaker, Level 2 - Works",
			"\"Market-Wide Circuit Breaker, Level 2\"", Financial_Quotes_MarketWideCircuitBreakerLevel2),
		Entry("Market-Wide Circuit Breaker, Level 3 - Works",
			"\"Market-Wide Circuit Breaker, Level 3\"", Financial_Quotes_MarketWideCircuitBreakerLevel3),
		Entry("Republished LULD Price Band - Works",
			"\"Republished LULD Price Band\"", Financial_Quotes_RepublishedLULDPriceBand),
		Entry("On-Demand Auction - Works", "\"On-Demand Auction\"", Financial_Quotes_OnDemandAuction),
		Entry("Cash-Only Settlement - Works", "\"Cash-Only Settlement\"", Financial_Quotes_CashOnlySettlement),
		Entry("Next-Day Settlement - Works", "\"Next-Day Settlement\"", Financial_Quotes_NextDaySettlement),
		Entry("LULD Trading Pause - Works", "\"LULD Trading Pause\"", Financial_Quotes_LULDTradingPause),
		Entry("Slow Due LRP, Bid, Ask - Works", "\"Slow Due LRP, Bid, Ask\"", Financial_Quotes_SlowDuelRPBidAsk),
		Entry("Corrected Price - Works", "\"Corrected Price\"", Financial_Quotes_CorrectedPrice),
		Entry("SIP-Generated - Works", "\"SIP-Generated\"", Financial_Quotes_SIPGenerated),
		Entry("Crossed Market - Works", "\"Crossed Market\"", Financial_Quotes_CrossedMarket),
		Entry("Locked Market - Works", "\"Locked Market\"", Financial_Quotes_LockedMarket),
		Entry("Depth on Offer Side - Works", "\"Depth on Offer Side\"", Financial_Quotes_DepthOnOfferSide),
		Entry("Depth on Bid Side - Works", "\"Depth on Bid Side\"", Financial_Quotes_DepthOnBidSide),
		Entry("Depth on Bid and Offer - Works", "\"Depth on Bid and Offer\"", Financial_Quotes_DepthOnBidAndOffer),
		Entry("Pre-Opening Indication - Works", "\"Pre-Opening Indication\"", Financial_Quotes_PreOpeningIndication),
		Entry("Syndicate Bid - Works", "\"Syndicate Bid\"", Financial_Quotes_SyndicateBid),
		Entry("Pre-Syndicate Bid - Works", "\"Pre-Syndicate Bid\"", Financial_Quotes_PreSyndicateBid),
		Entry("Penalty Bid - Works", "\"Penalty Bid\"", Financial_Quotes_PenaltyBid),
		Entry("CQS-Generated - Works", "\"CQS-Generated\"", Financial_Quotes_CQSGenerated),
		Entry("Regular - Works", "\"Regular\"", Financial_Quotes_Regular),
		Entry("RegularTwoSidedOpen - Works", "\"RegularTwoSidedOpen\"", Financial_Quotes_RegularTwoSidedOpen),
		Entry("RegularOneSidedOpen - Works", "\"RegularOneSidedOpen\"", Financial_Quotes_RegularOneSidedOpen),
		Entry("SlowAsk - Works", "\"SlowAsk\"", Financial_Quotes_SlowAsk),
		Entry("SlowBid - Works", "\"SlowBid\"", Financial_Quotes_SlowBid),
		Entry("SlowBidAsk - Works", "\"SlowBidAsk\"", Financial_Quotes_SlowBidAsk),
		Entry("SlowDueLRPBid - Works", "\"SlowDueLRPBid\"", Financial_Quotes_SlowDueLRPBid),
		Entry("SlowDueLRPAsk - Works", "\"SlowDueLRPAsk\"", Financial_Quotes_SlowDueLRPAsk),
		Entry("SlowDueNYSELRP - Works", "\"SlowDueNYSELRP\"", Financial_Quotes_SlowDueNYSELRP),
		Entry("SlowDueSetSlowListBidAsk - Works", "\"SlowDueSetSlowListBidAsk\"", Financial_Quotes_SlowDueSetSlowListBidAsk),
		Entry("ManualAskAutomatedBid - Works", "\"ManualAskAutomatedBid\"", Financial_Quotes_ManualAskAutomatedBid),
		Entry("ManualBidAutomatedAsk - Works", "\"ManualBidAutomatedAsk\"", Financial_Quotes_ManualBidAutomatedAsk),
		Entry("ManualBidAndAsk - Works", "\"ManualBidAndAsk\"", Financial_Quotes_ManualBidAndAsk),
		Entry("Opening - Works", "\"Opening\"", Financial_Quotes_Opening),
		Entry("Closing - Works", "\"Closing\"", Financial_Quotes_Closing),
		Entry("Closed - Works", "\"Closed\"", Financial_Quotes_Closed),
		Entry("Resume - Works", "\"Resume\"", Financial_Quotes_Resume),
		Entry("FastTrading - Works", "\"FastTrading\"", Financial_Quotes_FastTrading),
		Entry("TradingRangeIndicated - Works", "\"TradingRangeIndicated\"", Financial_Quotes_TradingRangeIndicated),
		Entry("MarketMakerQuotesClosed - Works", "\"MarketMakerQuotesClosed\"", Financial_Quotes_MarketMakerQuotesClosed),
		Entry("NonFirm - Works", "\"NonFirm\"", Financial_Quotes_NonFirm),
		Entry("NewsDissemination - Works", "\"NewsDissemination\"", Financial_Quotes_NewsDissemination),
		Entry("OrderInflux - Works", "\"OrderInflux\"", Financial_Quotes_OrderInflux),
		Entry("OrderImbalance - Works", "\"OrderImbalance\"", Financial_Quotes_OrderImbalance),
		Entry("DueToRelatedSecurityNewsDissemination - Works",
			"\"DueToRelatedSecurityNewsDissemination\"", Financial_Quotes_DueToRelatedSecurityNewsDissemination),
		Entry("DueToRelatedSecurityNewsPending - Works",
			"\"DueToRelatedSecurityNewsPending\"", Financial_Quotes_DueToRelatedSecurityNewsPending),
		Entry("AdditionalInformation - Works",
			"\"AdditionalInformation\"", Financial_Quotes_AdditionalInformation),
		Entry("NewsPending - Works", "\"NewsPending\"", Financial_Quotes_NewsPending),
		Entry("AdditionalInformationDueToRelatedSecurity - Works",
			"\"AdditionalInformationDueToRelatedSecurity\"", Financial_Quotes_AdditionalInformationDueToRelatedSecurity),
		Entry("DueToRelatedSecurity - Works", "\"DueToRelatedSecurity\"", Financial_Quotes_DueToRelatedSecurity),
		Entry("InViewOfCommon - Works", "\"InViewOfCommon\"", Financial_Quotes_InViewOfCommon),
		Entry("EquipmentChangeover - Works", "\"EquipmentChangeover\"", Financial_Quotes_EquipmentChangeover),
		Entry("NoOpenNoResponse - Works", "\"NoOpenNoResponse\"", Financial_Quotes_NoOpenNoResponse),
		Entry("SubPennyTrading - Works", "\"SubPennyTrading\"", Financial_Quotes_SubPennyTrading),
		Entry("AutomatedBidNoOfferNoBid - Works", "\"AutomatedBidNoOfferNoBid\"", Financial_Quotes_AutomatedBidNoOfferNoBid),
		Entry("LULDPPriceBand - Works", "\"LULDPPriceBand\"", Financial_Quotes_LULDPPriceBand),
		Entry("MarketWideCircuitBreakerLevel1 - Works",
			"\"MarketWideCircuitBreakerLevel1\"", Financial_Quotes_MarketWideCircuitBreakerLevel1),
		Entry("MarketWideCircuitBreakerLevel2 - Works",
			"\"MarketWideCircuitBreakerLevel2\"", Financial_Quotes_MarketWideCircuitBreakerLevel2),
		Entry("MarketWideCircuitBreakerLevel3 - Works",
			"\"MarketWideCircuitBreakerLevel3\"", Financial_Quotes_MarketWideCircuitBreakerLevel3),
		Entry("RepublishedLULDPriceBand - Works", "\"RepublishedLULDPriceBand\"", Financial_Quotes_RepublishedLULDPriceBand),
		Entry("OnDemandAuction - Works", "\"OnDemandAuction\"", Financial_Quotes_OnDemandAuction),
		Entry("CashOnlySettlement - Works", "\"CashOnlySettlement\"", Financial_Quotes_CashOnlySettlement),
		Entry("NextDaySettlement - Works", "\"NextDaySettlement\"", Financial_Quotes_NextDaySettlement),
		Entry("LULDTradingPause - Works", "\"LULDTradingPause\"", Financial_Quotes_LULDTradingPause),
		Entry("SlowDuelRPBidAsk - Works", "\"SlowDuelRPBidAsk\"", Financial_Quotes_SlowDuelRPBidAsk),
		Entry("Cancel - Works", "\"Cancel\"", Financial_Quotes_Cancel),
		Entry("CorrectedPrice - Works", "\"CorrectedPrice\"", Financial_Quotes_CorrectedPrice),
		Entry("SIPGenerated - Works", "\"SIPGenerated\"", Financial_Quotes_SIPGenerated),
		Entry("Unknown - Works", "\"Unknown\"", Financial_Quotes_Unknown),
		Entry("CrossedMarket - Works", "\"CrossedMarket\"", Financial_Quotes_CrossedMarket),
		Entry("LockedMarket - Works", "\"LockedMarket\"", Financial_Quotes_LockedMarket),
		Entry("DepthOnOfferSide - Works", "\"DepthOnOfferSide\"", Financial_Quotes_DepthOnOfferSide),
		Entry("DepthOnBidSide - Works", "\"DepthOnBidSide\"", Financial_Quotes_DepthOnBidSide),
		Entry("DepthOnBidAndOffer - Works", "\"DepthOnBidAndOffer\"", Financial_Quotes_DepthOnBidAndOffer),
		Entry("PreOpeningIndication - Works", "\"PreOpeningIndication\"", Financial_Quotes_PreOpeningIndication),
		Entry("SyndicateBid - Works", "\"SyndicateBid\"", Financial_Quotes_SyndicateBid),
		Entry("PreSyndicateBid - Works", "\"PreSyndicateBid\"", Financial_Quotes_PreSyndicateBid),
		Entry("PenaltyBid - Works", "\"PenaltyBid\"", Financial_Quotes_PenaltyBid),
		Entry("CQSGenerated - Works", "\"CQSGenerated\"", Financial_Quotes_CQSGenerated),
		Entry("Invalid - Works", "\"Invalid\"", Financial_Quotes_Invalid),
		Entry("'-1' - Works", "\"-1\"", Financial_Quotes_Invalid),
		Entry("'0' - Works", "\"0\"", Financial_Quotes_Regular),
		Entry("'1' - Works", "\"1\"", Financial_Quotes_RegularTwoSidedOpen),
		Entry("'2' - Works", "\"2\"", Financial_Quotes_RegularOneSidedOpen),
		Entry("'3' - Works", "\"3\"", Financial_Quotes_SlowAsk),
		Entry("'4' - Works", "\"4\"", Financial_Quotes_SlowBid),
		Entry("'5' - Works", "\"5\"", Financial_Quotes_SlowBidAsk),
		Entry("'6' - Works", "\"6\"", Financial_Quotes_SlowDueLRPBid),
		Entry("'7' - Works", "\"7\"", Financial_Quotes_SlowDueLRPAsk),
		Entry("'8' - Works", "\"8\"", Financial_Quotes_SlowDueNYSELRP),
		Entry("'9' - Works", "\"9\"", Financial_Quotes_SlowDueSetSlowListBidAsk),
		Entry("'10' - Works", "\"10\"", Financial_Quotes_ManualAskAutomatedBid),
		Entry("'11' - Works", "\"11\"", Financial_Quotes_ManualBidAutomatedAsk),
		Entry("'12' - Works", "\"12\"", Financial_Quotes_ManualBidAndAsk),
		Entry("'13' - Works", "\"13\"", Financial_Quotes_Opening),
		Entry("'14' - Works", "\"14\"", Financial_Quotes_Closing),
		Entry("'15' - Works", "\"15\"", Financial_Quotes_Closed),
		Entry("'16' - Works", "\"16\"", Financial_Quotes_Resume),
		Entry("'17' - Works", "\"17\"", Financial_Quotes_FastTrading),
		Entry("'18' - Works", "\"18\"", Financial_Quotes_TradingRangeIndicated),
		Entry("'19' - Works", "\"19\"", Financial_Quotes_MarketMakerQuotesClosed),
		Entry("'20' - Works", "\"20\"", Financial_Quotes_NonFirm),
		Entry("'21' - Works", "\"21\"", Financial_Quotes_NewsDissemination),
		Entry("'22' - Works", "\"22\"", Financial_Quotes_OrderInflux),
		Entry("'23' - Works", "\"23\"", Financial_Quotes_OrderImbalance),
		Entry("'24' - Works", "\"24\"", Financial_Quotes_DueToRelatedSecurityNewsDissemination),
		Entry("'25' - Works", "\"25\"", Financial_Quotes_DueToRelatedSecurityNewsPending),
		Entry("'26' - Works", "\"26\"", Financial_Quotes_AdditionalInformation),
		Entry("'27' - Works", "\"27\"", Financial_Quotes_NewsPending),
		Entry("'28' - Works", "\"28\"", Financial_Quotes_AdditionalInformationDueToRelatedSecurity),
		Entry("'29' - Works", "\"29\"", Financial_Quotes_DueToRelatedSecurity),
		Entry("'30' - Works", "\"30\"", Financial_Quotes_InViewOfCommon),
		Entry("'31' - Works", "\"31\"", Financial_Quotes_EquipmentChangeover),
		Entry("'32' - Works", "\"32\"", Financial_Quotes_NoOpenNoResponse),
		Entry("'33' - Works", "\"33\"", Financial_Quotes_SubPennyTrading),
		Entry("'34' - Works", "\"34\"", Financial_Quotes_AutomatedBidNoOfferNoBid),
		Entry("'35' - Works", "\"35\"", Financial_Quotes_LULDPPriceBand),
		Entry("'36' - Works", "\"36\"", Financial_Quotes_MarketWideCircuitBreakerLevel1),
		Entry("'37' - Works", "\"37\"", Financial_Quotes_MarketWideCircuitBreakerLevel2),
		Entry("'38' - Works", "\"38\"", Financial_Quotes_MarketWideCircuitBreakerLevel3),
		Entry("'39' - Works", "\"39\"", Financial_Quotes_RepublishedLULDPriceBand),
		Entry("'40' - Works", "\"40\"", Financial_Quotes_OnDemandAuction),
		Entry("'41' - Works", "\"41\"", Financial_Quotes_CashOnlySettlement),
		Entry("'42' - Works", "\"42\"", Financial_Quotes_NextDaySettlement),
		Entry("'43' - Works", "\"43\"", Financial_Quotes_LULDTradingPause),
		Entry("'71' - Works", "\"71\"", Financial_Quotes_SlowDuelRPBidAsk),
		Entry("'80' - Works", "\"80\"", Financial_Quotes_Cancel),
		Entry("'81' - Works", "\"81\"", Financial_Quotes_CorrectedPrice),
		Entry("'82' - Works", "\"82\"", Financial_Quotes_SIPGenerated),
		Entry("'83' - Works", "\"83\"", Financial_Quotes_Unknown),
		Entry("'84' - Works", "\"84\"", Financial_Quotes_CrossedMarket),
		Entry("'85' - Works", "\"85\"", Financial_Quotes_LockedMarket),
		Entry("'86' - Works", "\"86\"", Financial_Quotes_DepthOnOfferSide),
		Entry("'87' - Works", "\"87\"", Financial_Quotes_DepthOnBidSide),
		Entry("'88' - Works", "\"88\"", Financial_Quotes_DepthOnBidAndOffer),
		Entry("'89' - Works", "\"89\"", Financial_Quotes_PreOpeningIndication),
		Entry("'90' - Works", "\"90\"", Financial_Quotes_SyndicateBid),
		Entry("'91' - Works", "\"91\"", Financial_Quotes_PreSyndicateBid),
		Entry("'92' - Works", "\"92\"", Financial_Quotes_PenaltyBid),
		Entry("'94' - Works", "\"94\"", Financial_Quotes_CQSGenerated),
		Entry("'999' - Works", "\"999\"", Financial_Quotes_Invalid),
		Entry("-1 - Works", -1, Financial_Quotes_Invalid),
		Entry("0 - Works", 0, Financial_Quotes_Regular),
		Entry("1 - Works", 1, Financial_Quotes_RegularTwoSidedOpen),
		Entry("2 - Works", 2, Financial_Quotes_RegularOneSidedOpen),
		Entry("3 - Works", 3, Financial_Quotes_SlowAsk),
		Entry("4 - Works", 4, Financial_Quotes_SlowBid),
		Entry("5 - Works", 5, Financial_Quotes_SlowBidAsk),
		Entry("6 - Works", 6, Financial_Quotes_SlowDueLRPBid),
		Entry("7 - Works", 7, Financial_Quotes_SlowDueLRPAsk),
		Entry("8 - Works", 8, Financial_Quotes_SlowDueNYSELRP),
		Entry("9 - Works", 9, Financial_Quotes_SlowDueSetSlowListBidAsk),
		Entry("10 - Works", 10, Financial_Quotes_ManualAskAutomatedBid),
		Entry("11 - Works", 11, Financial_Quotes_ManualBidAutomatedAsk),
		Entry("12 - Works", 12, Financial_Quotes_ManualBidAndAsk),
		Entry("13 - Works", 13, Financial_Quotes_Opening),
		Entry("14 - Works", 14, Financial_Quotes_Closing),
		Entry("15 - Works", 15, Financial_Quotes_Closed),
		Entry("16 - Works", 16, Financial_Quotes_Resume),
		Entry("17 - Works", 17, Financial_Quotes_FastTrading),
		Entry("18 - Works", 18, Financial_Quotes_TradingRangeIndicated),
		Entry("19 - Works", 19, Financial_Quotes_MarketMakerQuotesClosed),
		Entry("20 - Works", 20, Financial_Quotes_NonFirm),
		Entry("21 - Works", 21, Financial_Quotes_NewsDissemination),
		Entry("22 - Works", 22, Financial_Quotes_OrderInflux),
		Entry("23 - Works", 23, Financial_Quotes_OrderImbalance),
		Entry("24 - Works", 24, Financial_Quotes_DueToRelatedSecurityNewsDissemination),
		Entry("25 - Works", 25, Financial_Quotes_DueToRelatedSecurityNewsPending),
		Entry("26 - Works", 26, Financial_Quotes_AdditionalInformation),
		Entry("27 - Works", 27, Financial_Quotes_NewsPending),
		Entry("28 - Works", 28, Financial_Quotes_AdditionalInformationDueToRelatedSecurity),
		Entry("29 - Works", 29, Financial_Quotes_DueToRelatedSecurity),
		Entry("30 - Works", 30, Financial_Quotes_InViewOfCommon),
		Entry("31 - Works", 31, Financial_Quotes_EquipmentChangeover),
		Entry("32 - Works", 32, Financial_Quotes_NoOpenNoResponse),
		Entry("33 - Works", 33, Financial_Quotes_SubPennyTrading),
		Entry("34 - Works", 34, Financial_Quotes_AutomatedBidNoOfferNoBid),
		Entry("35 - Works", 35, Financial_Quotes_LULDPPriceBand),
		Entry("36 - Works", 36, Financial_Quotes_MarketWideCircuitBreakerLevel1),
		Entry("37 - Works", 37, Financial_Quotes_MarketWideCircuitBreakerLevel2),
		Entry("38 - Works", 38, Financial_Quotes_MarketWideCircuitBreakerLevel3),
		Entry("39 - Works", 39, Financial_Quotes_RepublishedLULDPriceBand),
		Entry("40 - Works", 40, Financial_Quotes_OnDemandAuction),
		Entry("41 - Works", 41, Financial_Quotes_CashOnlySettlement),
		Entry("42 - Works", 42, Financial_Quotes_NextDaySettlement),
		Entry("43 - Works", 43, Financial_Quotes_LULDTradingPause),
		Entry("71 - Works", 71, Financial_Quotes_SlowDuelRPBidAsk),
		Entry("80 - Works", 80, Financial_Quotes_Cancel),
		Entry("81 - Works", 81, Financial_Quotes_CorrectedPrice),
		Entry("82 - Works", 82, Financial_Quotes_SIPGenerated),
		Entry("83 - Works", 83, Financial_Quotes_Unknown),
		Entry("84 - Works", 84, Financial_Quotes_CrossedMarket),
		Entry("85 - Works", 85, Financial_Quotes_LockedMarket),
		Entry("86 - Works", 86, Financial_Quotes_DepthOnOfferSide),
		Entry("87 - Works", 87, Financial_Quotes_DepthOnBidSide),
		Entry("88 - Works", 88, Financial_Quotes_DepthOnBidAndOffer),
		Entry("89 - Works", 89, Financial_Quotes_PreOpeningIndication),
		Entry("90 - Works", 90, Financial_Quotes_SyndicateBid),
		Entry("91 - Works", 91, Financial_Quotes_PreSyndicateBid),
		Entry("92 - Works", 92, Financial_Quotes_PenaltyBid),
		Entry("94 - Works", 94, Financial_Quotes_CQSGenerated),
		Entry("999 - Works", 999, Financial_Quotes_Invalid))

	// Test that attempting to deserialize a Financial.Quotes.Condition will fial and return an
	// error if the value cannot be converted to either the name value or integer value
	// of the enum option
	It("UnmarshalCSV - Value is empty - Error", func() {

		// Attempt to convert a fake string value into a Financial.Quotes.Condition
		// This should return an error
		enum := new(Financial_Quotes_Condition)
		err := enum.UnmarshalCSV("")

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"\" cannot be mapped to a gopb.Financial_Quotes_Condition"))
	})

	// Test the conditions under which values should be convertible to a Financial.Quotes.Condition
	DescribeTable("UnmarshalCSV Tests",
		func(value string, shouldBe Financial_Quotes_Condition) {

			// Attempt to convert the value into a Financial.Quotes.Condition
			// This should not fail
			var enum Financial_Quotes_Condition
			err := enum.UnmarshalCSV(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("Regular, Two-Sided Open - Works", "Regular, Two-Sided Open", Financial_Quotes_RegularTwoSidedOpen),
		Entry("Regular, One-Sided Open - Works", "Regular, One-Sided Open", Financial_Quotes_RegularOneSidedOpen),
		Entry("Slow Ask - Works", "Slow Ask", Financial_Quotes_SlowAsk),
		Entry("Slow Bid - Works", "Slow Bid", Financial_Quotes_SlowBid),
		Entry("Slow Bid, Ask - Works", "Slow Bid, Ask", Financial_Quotes_SlowBidAsk),
		Entry("Slow Due, LRP Bid - Works", "Slow Due, LRP Bid", Financial_Quotes_SlowDueLRPBid),
		Entry("Slow Due, LRP Ask - Works", "Slow Due, LRP Ask", Financial_Quotes_SlowDueLRPAsk),
		Entry("Slow Due, NYSE LRP - Works", "Slow Due, NYSE LRP", Financial_Quotes_SlowDueNYSELRP),
		Entry("Slow Due Set, Slow List, Bid, Ask - Works",
			"Slow Due Set, Slow List, Bid, Ask", Financial_Quotes_SlowDueSetSlowListBidAsk),
		Entry("Manual Ask, Automated Bid - Works", "Manual Ask, Automated Bid", Financial_Quotes_ManualAskAutomatedBid),
		Entry("Manual Bid, Automated Ask - Works", "Manual Bid, Automated Ask", Financial_Quotes_ManualBidAutomatedAsk),
		Entry("Manual Bid and Ask - Works", "Manual Bid and Ask", Financial_Quotes_ManualBidAndAsk),
		Entry("Fast Trading - Works", "Fast Trading", Financial_Quotes_FastTrading),
		Entry("Tading Range Indicated - Works", "Tading Range Indicated", Financial_Quotes_TradingRangeIndicated),
		Entry("Market-Maker Quotes Closed - Works", "Market-Maker Quotes Closed", Financial_Quotes_MarketMakerQuotesClosed),
		Entry("Non-Firm - Works", "Non-Firm", Financial_Quotes_NonFirm),
		Entry("News Dissemination - Works", "News Dissemination", Financial_Quotes_NewsDissemination),
		Entry("Order Influx - Works", "Order Influx", Financial_Quotes_OrderInflux),
		Entry("Order Imbalance - Works", "Order Imbalance", Financial_Quotes_OrderImbalance),
		Entry("Due to Related Security, News Dissemination - Works",
			"Due to Related Security, News Dissemination", Financial_Quotes_DueToRelatedSecurityNewsDissemination),
		Entry("Due to Related Security, News Pending - Works",
			"Due to Related Security, News Pending", Financial_Quotes_DueToRelatedSecurityNewsPending),
		Entry("Additional Information - Works", "Additional Information", Financial_Quotes_AdditionalInformation),
		Entry("News Pending - Works", "News Pending", Financial_Quotes_NewsPending),
		Entry("Additional Information Due to Related Security - Works",
			"Additional Information Due to Related Security", Financial_Quotes_AdditionalInformationDueToRelatedSecurity),
		Entry("Due to Related Security - Works", "Due to Related Security", Financial_Quotes_DueToRelatedSecurity),
		Entry("In View of Common - Works", "In View of Common", Financial_Quotes_InViewOfCommon),
		Entry("Equipment Changeover - Works", "Equipment Changeover", Financial_Quotes_EquipmentChangeover),
		Entry("No Open, No Response - Works", "No Open, No Response", Financial_Quotes_NoOpenNoResponse),
		Entry("Sub-Penny Trading - Works", "Sub-Penny Trading", Financial_Quotes_SubPennyTrading),
		Entry("Automated Bid; No Offer, No Bid - Works",
			"Automated Bid; No Offer, No Bid", Financial_Quotes_AutomatedBidNoOfferNoBid),
		Entry("LULD Price Band - Works", "LULD Price Band", Financial_Quotes_LULDPPriceBand),
		Entry("Market-Wide Circuit Breaker, Level 1 - Works",
			"Market-Wide Circuit Breaker, Level 1", Financial_Quotes_MarketWideCircuitBreakerLevel1),
		Entry("Market-Wide Circuit Breaker, Level 2 - Works",
			"Market-Wide Circuit Breaker, Level 2", Financial_Quotes_MarketWideCircuitBreakerLevel2),
		Entry("Market-Wide Circuit Breaker, Level 3 - Works",
			"Market-Wide Circuit Breaker, Level 3", Financial_Quotes_MarketWideCircuitBreakerLevel3),
		Entry("Republished LULD Price Band - Works",
			"Republished LULD Price Band", Financial_Quotes_RepublishedLULDPriceBand),
		Entry("On-Demand Auction - Works", "On-Demand Auction", Financial_Quotes_OnDemandAuction),
		Entry("Cash-Only Settlement - Works", "Cash-Only Settlement", Financial_Quotes_CashOnlySettlement),
		Entry("Next-Day Settlement - Works", "Next-Day Settlement", Financial_Quotes_NextDaySettlement),
		Entry("LULD Trading Pause - Works", "LULD Trading Pause", Financial_Quotes_LULDTradingPause),
		Entry("Slow Due LRP, Bid, Ask - Works", "Slow Due LRP, Bid, Ask", Financial_Quotes_SlowDuelRPBidAsk),
		Entry("Corrected Price - Works", "Corrected Price", Financial_Quotes_CorrectedPrice),
		Entry("SIP-Generated - Works", "SIP-Generated", Financial_Quotes_SIPGenerated),
		Entry("Crossed Market - Works", "Crossed Market", Financial_Quotes_CrossedMarket),
		Entry("Locked Market - Works", "Locked Market", Financial_Quotes_LockedMarket),
		Entry("Depth on Offer Side - Works", "Depth on Offer Side", Financial_Quotes_DepthOnOfferSide),
		Entry("Depth on Bid Side - Works", "Depth on Bid Side", Financial_Quotes_DepthOnBidSide),
		Entry("Depth on Bid and Offer - Works", "Depth on Bid and Offer", Financial_Quotes_DepthOnBidAndOffer),
		Entry("Pre-Opening Indication - Works", "Pre-Opening Indication", Financial_Quotes_PreOpeningIndication),
		Entry("Syndicate Bid - Works", "Syndicate Bid", Financial_Quotes_SyndicateBid),
		Entry("Pre-Syndicate Bid - Works", "Pre-Syndicate Bid", Financial_Quotes_PreSyndicateBid),
		Entry("Penalty Bid - Works", "Penalty Bid", Financial_Quotes_PenaltyBid),
		Entry("CQS-Generated - Works", "CQS-Generated", Financial_Quotes_CQSGenerated),
		Entry("Regular - Works", "Regular", Financial_Quotes_Regular),
		Entry("RegularTwoSidedOpen - Works", "RegularTwoSidedOpen", Financial_Quotes_RegularTwoSidedOpen),
		Entry("RegularOneSidedOpen - Works", "RegularOneSidedOpen", Financial_Quotes_RegularOneSidedOpen),
		Entry("SlowAsk - Works", "SlowAsk", Financial_Quotes_SlowAsk),
		Entry("SlowBid - Works", "SlowBid", Financial_Quotes_SlowBid),
		Entry("SlowBidAsk - Works", "SlowBidAsk", Financial_Quotes_SlowBidAsk),
		Entry("SlowDueLRPBid - Works", "SlowDueLRPBid", Financial_Quotes_SlowDueLRPBid),
		Entry("SlowDueLRPAsk - Works", "SlowDueLRPAsk", Financial_Quotes_SlowDueLRPAsk),
		Entry("SlowDueNYSELRP - Works", "SlowDueNYSELRP", Financial_Quotes_SlowDueNYSELRP),
		Entry("SlowDueSetSlowListBidAsk - Works", "SlowDueSetSlowListBidAsk", Financial_Quotes_SlowDueSetSlowListBidAsk),
		Entry("ManualAskAutomatedBid - Works", "ManualAskAutomatedBid", Financial_Quotes_ManualAskAutomatedBid),
		Entry("ManualBidAutomatedAsk - Works", "ManualBidAutomatedAsk", Financial_Quotes_ManualBidAutomatedAsk),
		Entry("ManualBidAndAsk - Works", "ManualBidAndAsk", Financial_Quotes_ManualBidAndAsk),
		Entry("Opening - Works", "Opening", Financial_Quotes_Opening),
		Entry("Closing - Works", "Closing", Financial_Quotes_Closing),
		Entry("Closed - Works", "Closed", Financial_Quotes_Closed),
		Entry("Resume - Works", "Resume", Financial_Quotes_Resume),
		Entry("FastTrading - Works", "FastTrading", Financial_Quotes_FastTrading),
		Entry("TradingRangeIndicated - Works", "TradingRangeIndicated", Financial_Quotes_TradingRangeIndicated),
		Entry("MarketMakerQuotesClosed - Works", "MarketMakerQuotesClosed", Financial_Quotes_MarketMakerQuotesClosed),
		Entry("NonFirm - Works", "NonFirm", Financial_Quotes_NonFirm),
		Entry("NewsDissemination - Works", "NewsDissemination", Financial_Quotes_NewsDissemination),
		Entry("OrderInflux - Works", "OrderInflux", Financial_Quotes_OrderInflux),
		Entry("OrderImbalance - Works", "OrderImbalance", Financial_Quotes_OrderImbalance),
		Entry("DueToRelatedSecurityNewsDissemination - Works", "DueToRelatedSecurityNewsDissemination", Financial_Quotes_DueToRelatedSecurityNewsDissemination),
		Entry("DueToRelatedSecurityNewsPending - Works", "DueToRelatedSecurityNewsPending", Financial_Quotes_DueToRelatedSecurityNewsPending),
		Entry("AdditionalInformation - Works", "AdditionalInformation", Financial_Quotes_AdditionalInformation),
		Entry("NewsPending - Works", "NewsPending", Financial_Quotes_NewsPending),
		Entry("AdditionalInformationDueToRelatedSecurity - Works", "AdditionalInformationDueToRelatedSecurity", Financial_Quotes_AdditionalInformationDueToRelatedSecurity),
		Entry("DueToRelatedSecurity - Works", "DueToRelatedSecurity", Financial_Quotes_DueToRelatedSecurity),
		Entry("InViewOfCommon - Works", "InViewOfCommon", Financial_Quotes_InViewOfCommon),
		Entry("EquipmentChangeover - Works", "EquipmentChangeover", Financial_Quotes_EquipmentChangeover),
		Entry("NoOpenNoResponse - Works", "NoOpenNoResponse", Financial_Quotes_NoOpenNoResponse),
		Entry("SubPennyTrading - Works", "SubPennyTrading", Financial_Quotes_SubPennyTrading),
		Entry("AutomatedBidNoOfferNoBid - Works", "AutomatedBidNoOfferNoBid", Financial_Quotes_AutomatedBidNoOfferNoBid),
		Entry("LULDPPriceBand - Works", "LULDPPriceBand", Financial_Quotes_LULDPPriceBand),
		Entry("MarketWideCircuitBreakerLevel1 - Works", "MarketWideCircuitBreakerLevel1", Financial_Quotes_MarketWideCircuitBreakerLevel1),
		Entry("MarketWideCircuitBreakerLevel2 - Works", "MarketWideCircuitBreakerLevel2", Financial_Quotes_MarketWideCircuitBreakerLevel2),
		Entry("MarketWideCircuitBreakerLevel3 - Works", "MarketWideCircuitBreakerLevel3", Financial_Quotes_MarketWideCircuitBreakerLevel3),
		Entry("RepublishedLULDPriceBand - Works", "RepublishedLULDPriceBand", Financial_Quotes_RepublishedLULDPriceBand),
		Entry("OnDemandAuction - Works", "OnDemandAuction", Financial_Quotes_OnDemandAuction),
		Entry("CashOnlySettlement - Works", "CashOnlySettlement", Financial_Quotes_CashOnlySettlement),
		Entry("NextDaySettlement - Works", "NextDaySettlement", Financial_Quotes_NextDaySettlement),
		Entry("LULDTradingPause - Works", "LULDTradingPause", Financial_Quotes_LULDTradingPause),
		Entry("SlowDuelRPBidAsk - Works", "SlowDuelRPBidAsk", Financial_Quotes_SlowDuelRPBidAsk),
		Entry("Cancel - Works", "Cancel", Financial_Quotes_Cancel),
		Entry("CorrectedPrice - Works", "CorrectedPrice", Financial_Quotes_CorrectedPrice),
		Entry("SIPGenerated - Works", "SIPGenerated", Financial_Quotes_SIPGenerated),
		Entry("Unknown - Works", "Unknown", Financial_Quotes_Unknown),
		Entry("CrossedMarket - Works", "CrossedMarket", Financial_Quotes_CrossedMarket),
		Entry("LockedMarket - Works", "LockedMarket", Financial_Quotes_LockedMarket),
		Entry("DepthOnOfferSide - Works", "DepthOnOfferSide", Financial_Quotes_DepthOnOfferSide),
		Entry("DepthOnBidSide - Works", "DepthOnBidSide", Financial_Quotes_DepthOnBidSide),
		Entry("DepthOnBidAndOffer - Works", "DepthOnBidAndOffer", Financial_Quotes_DepthOnBidAndOffer),
		Entry("PreOpeningIndication - Works", "PreOpeningIndication", Financial_Quotes_PreOpeningIndication),
		Entry("SyndicateBid - Works", "SyndicateBid", Financial_Quotes_SyndicateBid),
		Entry("PreSyndicateBid - Works", "PreSyndicateBid", Financial_Quotes_PreSyndicateBid),
		Entry("PenaltyBid - Works", "PenaltyBid", Financial_Quotes_PenaltyBid),
		Entry("CQSGenerated - Works", "CQSGenerated", Financial_Quotes_CQSGenerated),
		Entry("Invalid - Works", "Invalid", Financial_Quotes_Invalid),
		Entry("0 - Works", "000", Financial_Quotes_Regular),
		Entry("1 - Works", "001", Financial_Quotes_RegularTwoSidedOpen),
		Entry("2 - Works", "002", Financial_Quotes_RegularOneSidedOpen),
		Entry("3 - Works", "003", Financial_Quotes_SlowAsk),
		Entry("4 - Works", "004", Financial_Quotes_SlowBid),
		Entry("5 - Works", "005", Financial_Quotes_SlowBidAsk),
		Entry("6 - Works", "006", Financial_Quotes_SlowDueLRPBid),
		Entry("7 - Works", "007", Financial_Quotes_SlowDueLRPAsk),
		Entry("8 - Works", "008", Financial_Quotes_SlowDueNYSELRP),
		Entry("9 - Works", "009", Financial_Quotes_SlowDueSetSlowListBidAsk),
		Entry("10 - Works", "010", Financial_Quotes_ManualAskAutomatedBid),
		Entry("11 - Works", "011", Financial_Quotes_ManualBidAutomatedAsk),
		Entry("12 - Works", "012", Financial_Quotes_ManualBidAndAsk),
		Entry("13 - Works", "013", Financial_Quotes_Opening),
		Entry("14 - Works", "014", Financial_Quotes_Closing),
		Entry("15 - Works", "015", Financial_Quotes_Closed),
		Entry("16 - Works", "016", Financial_Quotes_Resume),
		Entry("17 - Works", "017", Financial_Quotes_FastTrading),
		Entry("18 - Works", "018", Financial_Quotes_TradingRangeIndicated),
		Entry("19 - Works", "019", Financial_Quotes_MarketMakerQuotesClosed),
		Entry("20 - Works", "020", Financial_Quotes_NonFirm),
		Entry("21 - Works", "021", Financial_Quotes_NewsDissemination),
		Entry("22 - Works", "022", Financial_Quotes_OrderInflux),
		Entry("23 - Works", "023", Financial_Quotes_OrderImbalance),
		Entry("24 - Works", "024", Financial_Quotes_DueToRelatedSecurityNewsDissemination),
		Entry("25 - Works", "025", Financial_Quotes_DueToRelatedSecurityNewsPending),
		Entry("26 - Works", "026", Financial_Quotes_AdditionalInformation),
		Entry("27 - Works", "027", Financial_Quotes_NewsPending),
		Entry("28 - Works", "028", Financial_Quotes_AdditionalInformationDueToRelatedSecurity),
		Entry("29 - Works", "029", Financial_Quotes_DueToRelatedSecurity),
		Entry("30 - Works", "030", Financial_Quotes_InViewOfCommon),
		Entry("31 - Works", "031", Financial_Quotes_EquipmentChangeover),
		Entry("32 - Works", "032", Financial_Quotes_NoOpenNoResponse),
		Entry("33 - Works", "033", Financial_Quotes_SubPennyTrading),
		Entry("34 - Works", "034", Financial_Quotes_AutomatedBidNoOfferNoBid),
		Entry("35 - Works", "035", Financial_Quotes_LULDPPriceBand),
		Entry("36 - Works", "036", Financial_Quotes_MarketWideCircuitBreakerLevel1),
		Entry("37 - Works", "037", Financial_Quotes_MarketWideCircuitBreakerLevel2),
		Entry("38 - Works", "038", Financial_Quotes_MarketWideCircuitBreakerLevel3),
		Entry("39 - Works", "039", Financial_Quotes_RepublishedLULDPriceBand),
		Entry("40 - Works", "040", Financial_Quotes_OnDemandAuction),
		Entry("41 - Works", "041", Financial_Quotes_CashOnlySettlement),
		Entry("42 - Works", "042", Financial_Quotes_NextDaySettlement),
		Entry("43 - Works", "043", Financial_Quotes_LULDTradingPause),
		Entry("71 - Works", "071", Financial_Quotes_SlowDuelRPBidAsk),
		Entry("80 - Works", "080", Financial_Quotes_Cancel),
		Entry("81 - Works", "081", Financial_Quotes_CorrectedPrice),
		Entry("82 - Works", "082", Financial_Quotes_SIPGenerated),
		Entry("83 - Works", "083", Financial_Quotes_Unknown),
		Entry("84 - Works", "084", Financial_Quotes_CrossedMarket),
		Entry("85 - Works", "085", Financial_Quotes_LockedMarket),
		Entry("86 - Works", "086", Financial_Quotes_DepthOnOfferSide),
		Entry("87 - Works", "087", Financial_Quotes_DepthOnBidSide),
		Entry("88 - Works", "088", Financial_Quotes_DepthOnBidAndOffer),
		Entry("89 - Works", "089", Financial_Quotes_PreOpeningIndication),
		Entry("90 - Works", "090", Financial_Quotes_SyndicateBid),
		Entry("91 - Works", "091", Financial_Quotes_PreSyndicateBid),
		Entry("92 - Works", "092", Financial_Quotes_PenaltyBid),
		Entry("94 - Works", "094", Financial_Quotes_CQSGenerated),
		Entry("0 - Works", "0", Financial_Quotes_Regular),
		Entry("1 - Works", "1", Financial_Quotes_RegularTwoSidedOpen),
		Entry("2 - Works", "2", Financial_Quotes_RegularOneSidedOpen),
		Entry("3 - Works", "3", Financial_Quotes_SlowAsk),
		Entry("4 - Works", "4", Financial_Quotes_SlowBid),
		Entry("5 - Works", "5", Financial_Quotes_SlowBidAsk),
		Entry("6 - Works", "6", Financial_Quotes_SlowDueLRPBid),
		Entry("7 - Works", "7", Financial_Quotes_SlowDueLRPAsk),
		Entry("8 - Works", "8", Financial_Quotes_SlowDueNYSELRP),
		Entry("9 - Works", "9", Financial_Quotes_SlowDueSetSlowListBidAsk),
		Entry("10 - Works", "10", Financial_Quotes_ManualAskAutomatedBid),
		Entry("11 - Works", "11", Financial_Quotes_ManualBidAutomatedAsk),
		Entry("12 - Works", "12", Financial_Quotes_ManualBidAndAsk),
		Entry("13 - Works", "13", Financial_Quotes_Opening),
		Entry("14 - Works", "14", Financial_Quotes_Closing),
		Entry("15 - Works", "15", Financial_Quotes_Closed),
		Entry("16 - Works", "16", Financial_Quotes_Resume),
		Entry("17 - Works", "17", Financial_Quotes_FastTrading),
		Entry("18 - Works", "18", Financial_Quotes_TradingRangeIndicated),
		Entry("19 - Works", "19", Financial_Quotes_MarketMakerQuotesClosed),
		Entry("20 - Works", "20", Financial_Quotes_NonFirm),
		Entry("21 - Works", "21", Financial_Quotes_NewsDissemination),
		Entry("22 - Works", "22", Financial_Quotes_OrderInflux),
		Entry("23 - Works", "23", Financial_Quotes_OrderImbalance),
		Entry("24 - Works", "24", Financial_Quotes_DueToRelatedSecurityNewsDissemination),
		Entry("25 - Works", "25", Financial_Quotes_DueToRelatedSecurityNewsPending),
		Entry("26 - Works", "26", Financial_Quotes_AdditionalInformation),
		Entry("27 - Works", "27", Financial_Quotes_NewsPending),
		Entry("28 - Works", "28", Financial_Quotes_AdditionalInformationDueToRelatedSecurity),
		Entry("29 - Works", "29", Financial_Quotes_DueToRelatedSecurity),
		Entry("30 - Works", "30", Financial_Quotes_InViewOfCommon),
		Entry("31 - Works", "31", Financial_Quotes_EquipmentChangeover),
		Entry("32 - Works", "32", Financial_Quotes_NoOpenNoResponse),
		Entry("33 - Works", "33", Financial_Quotes_SubPennyTrading),
		Entry("34 - Works", "34", Financial_Quotes_AutomatedBidNoOfferNoBid),
		Entry("35 - Works", "35", Financial_Quotes_LULDPPriceBand),
		Entry("36 - Works", "36", Financial_Quotes_MarketWideCircuitBreakerLevel1),
		Entry("37 - Works", "37", Financial_Quotes_MarketWideCircuitBreakerLevel2),
		Entry("38 - Works", "38", Financial_Quotes_MarketWideCircuitBreakerLevel3),
		Entry("39 - Works", "39", Financial_Quotes_RepublishedLULDPriceBand),
		Entry("40 - Works", "40", Financial_Quotes_OnDemandAuction),
		Entry("41 - Works", "41", Financial_Quotes_CashOnlySettlement),
		Entry("42 - Works", "42", Financial_Quotes_NextDaySettlement),
		Entry("43 - Works", "43", Financial_Quotes_LULDTradingPause),
		Entry("71 - Works", "71", Financial_Quotes_SlowDuelRPBidAsk),
		Entry("80 - Works", "80", Financial_Quotes_Cancel),
		Entry("81 - Works", "81", Financial_Quotes_CorrectedPrice),
		Entry("82 - Works", "82", Financial_Quotes_SIPGenerated),
		Entry("83 - Works", "83", Financial_Quotes_Unknown),
		Entry("84 - Works", "84", Financial_Quotes_CrossedMarket),
		Entry("85 - Works", "85", Financial_Quotes_LockedMarket),
		Entry("86 - Works", "86", Financial_Quotes_DepthOnOfferSide),
		Entry("87 - Works", "87", Financial_Quotes_DepthOnBidSide),
		Entry("88 - Works", "88", Financial_Quotes_DepthOnBidAndOffer),
		Entry("89 - Works", "89", Financial_Quotes_PreOpeningIndication),
		Entry("90 - Works", "90", Financial_Quotes_SyndicateBid),
		Entry("91 - Works", "91", Financial_Quotes_PreSyndicateBid),
		Entry("92 - Works", "92", Financial_Quotes_PenaltyBid),
		Entry("94 - Works", "94", Financial_Quotes_CQSGenerated),
		Entry("999 - Works", "999", Financial_Quotes_Invalid))

	// Tests that, if the attribute type submitted to UnmarshalDynamoDBAttributeValue is not one we
	// recognize, then the function will return an error
	It("UnmarshalDynamoDBAttributeValue - AttributeValue type invalid - Error", func() {
		value := new(Financial_Quotes_Condition)
		err := attributevalue.Unmarshal(&types.AttributeValueMemberBOOL{Value: true}, &value)
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("Attribute value of *types.AttributeValueMemberBOOL could not be converted to a Financial.Quotes.Condition"))
	})

	// Tests the conditions under which UnmarshalDynamoDBAttributeValue is called and no error is generated
	DescribeTable("UnmarshalDynamoDBAttributeValue - AttributeValue Conditions",
		func(raw types.AttributeValue, expected Financial_Quotes_Condition) {
			var value Financial_Quotes_Condition
			err := attributevalue.Unmarshal(raw, &value)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(expected))
		},
		Entry("Value is []bytes, Regular, Two-Sided Open - Works",
			&types.AttributeValueMemberB{Value: []byte("Regular, Two-Sided Open")}, Financial_Quotes_RegularTwoSidedOpen),
		Entry("Value is []bytes, Regular, One-Sided Open - Works",
			&types.AttributeValueMemberB{Value: []byte("Regular, One-Sided Open")}, Financial_Quotes_RegularOneSidedOpen),
		Entry("Value is []bytes, Slow Ask - Works",
			&types.AttributeValueMemberB{Value: []byte("Slow Ask")}, Financial_Quotes_SlowAsk),
		Entry("Value is []bytes, Slow Bid - Works",
			&types.AttributeValueMemberB{Value: []byte("Slow Bid")}, Financial_Quotes_SlowBid),
		Entry("Value is []bytes, Slow Bid, Ask - Works",
			&types.AttributeValueMemberB{Value: []byte("Slow Bid, Ask")}, Financial_Quotes_SlowBidAsk),
		Entry("Value is []bytes, Slow Due, LRP Bid - Works",
			&types.AttributeValueMemberB{Value: []byte("Slow Due, LRP Bid")}, Financial_Quotes_SlowDueLRPBid),
		Entry("Value is []bytes, Slow Due, LRP Ask - Works",
			&types.AttributeValueMemberB{Value: []byte("Slow Due, LRP Ask")}, Financial_Quotes_SlowDueLRPAsk),
		Entry("Value is []bytes, Slow Due, NYSE LRP - Works",
			&types.AttributeValueMemberB{Value: []byte("Slow Due, NYSE LRP")}, Financial_Quotes_SlowDueNYSELRP),
		Entry("Value is []bytes, Slow Due Set, Slow List, Bid, Ask - Works",
			&types.AttributeValueMemberB{Value: []byte("Slow Due Set, Slow List, Bid, Ask")}, Financial_Quotes_SlowDueSetSlowListBidAsk),
		Entry("Value is []bytes, Manual Ask, Automated Bid - Works",
			&types.AttributeValueMemberB{Value: []byte("Manual Ask, Automated Bid")}, Financial_Quotes_ManualAskAutomatedBid),
		Entry("Value is []bytes, Manual Bid, Automated Ask - Works",
			&types.AttributeValueMemberB{Value: []byte("Manual Bid, Automated Ask")}, Financial_Quotes_ManualBidAutomatedAsk),
		Entry("Value is []bytes, Manual Bid and Ask - Works",
			&types.AttributeValueMemberB{Value: []byte("Manual Bid and Ask")}, Financial_Quotes_ManualBidAndAsk),
		Entry("Value is []bytes, Fast Trading - Works",
			&types.AttributeValueMemberB{Value: []byte("Fast Trading")}, Financial_Quotes_FastTrading),
		Entry("Value is []bytes, Tading Range Indicated - Works",
			&types.AttributeValueMemberB{Value: []byte("Tading Range Indicated")}, Financial_Quotes_TradingRangeIndicated),
		Entry("Value is []bytes, Market-Maker Quotes Closed - Works",
			&types.AttributeValueMemberB{Value: []byte("Market-Maker Quotes Closed")}, Financial_Quotes_MarketMakerQuotesClosed),
		Entry("Value is []bytes, Non-Firm - Works",
			&types.AttributeValueMemberB{Value: []byte("Non-Firm")}, Financial_Quotes_NonFirm),
		Entry("Value is []bytes, News Dissemination - Works",
			&types.AttributeValueMemberB{Value: []byte("News Dissemination")}, Financial_Quotes_NewsDissemination),
		Entry("Value is []bytes, Order Influx - Works",
			&types.AttributeValueMemberB{Value: []byte("Order Influx")}, Financial_Quotes_OrderInflux),
		Entry("Value is []bytes, Order Imbalance - Works",
			&types.AttributeValueMemberB{Value: []byte("Order Imbalance")}, Financial_Quotes_OrderImbalance),
		Entry("Value is []bytes, Due to Related Security, News Dissemination - Works",
			&types.AttributeValueMemberB{Value: []byte("Due to Related Security, News Dissemination")}, Financial_Quotes_DueToRelatedSecurityNewsDissemination),
		Entry("Value is []bytes, Due to Related Security, News Pending - Works",
			&types.AttributeValueMemberB{Value: []byte("Due to Related Security, News Pending")}, Financial_Quotes_DueToRelatedSecurityNewsPending),
		Entry("Value is []bytes, Additional Information - Works",
			&types.AttributeValueMemberB{Value: []byte("Additional Information")}, Financial_Quotes_AdditionalInformation),
		Entry("Value is []bytes, News Pending - Works",
			&types.AttributeValueMemberB{Value: []byte("News Pending")}, Financial_Quotes_NewsPending),
		Entry("Value is []bytes, Additional Information Due to Related Security - Works",
			&types.AttributeValueMemberB{Value: []byte("Additional Information Due to Related Security")}, Financial_Quotes_AdditionalInformationDueToRelatedSecurity),
		Entry("Value is []bytes, Due to Related Security - Works",
			&types.AttributeValueMemberB{Value: []byte("Due to Related Security")}, Financial_Quotes_DueToRelatedSecurity),
		Entry("Value is []bytes, In View of Common - Works",
			&types.AttributeValueMemberB{Value: []byte("In View of Common")}, Financial_Quotes_InViewOfCommon),
		Entry("Value is []bytes, Equipment Changeover - Works",
			&types.AttributeValueMemberB{Value: []byte("Equipment Changeover")}, Financial_Quotes_EquipmentChangeover),
		Entry("Value is []bytes, No Open, No Response - Works",
			&types.AttributeValueMemberB{Value: []byte("No Open, No Response")}, Financial_Quotes_NoOpenNoResponse),
		Entry("Value is []bytes, Sub-Penny Trading - Works",
			&types.AttributeValueMemberB{Value: []byte("Sub-Penny Trading")}, Financial_Quotes_SubPennyTrading),
		Entry("Value is []bytes, Automated Bid; No Offer, No Bid - Works",
			&types.AttributeValueMemberB{Value: []byte("Automated Bid; No Offer, No Bid")}, Financial_Quotes_AutomatedBidNoOfferNoBid),
		Entry("Value is []bytes, LULD Price Band - Works",
			&types.AttributeValueMemberB{Value: []byte("LULD Price Band")}, Financial_Quotes_LULDPPriceBand),
		Entry("Value is []bytes, Market-Wide Circuit Breaker, Level 1 - Works",
			&types.AttributeValueMemberB{Value: []byte("Market-Wide Circuit Breaker, Level 1")}, Financial_Quotes_MarketWideCircuitBreakerLevel1),
		Entry("Value is []bytes, Market-Wide Circuit Breaker, Level 2 - Works",
			&types.AttributeValueMemberB{Value: []byte("Market-Wide Circuit Breaker, Level 2")}, Financial_Quotes_MarketWideCircuitBreakerLevel2),
		Entry("Value is []bytes, Market-Wide Circuit Breaker, Level 3 - Works",
			&types.AttributeValueMemberB{Value: []byte("Market-Wide Circuit Breaker, Level 3")}, Financial_Quotes_MarketWideCircuitBreakerLevel3),
		Entry("Value is []bytes, Republished LULD Price Band - Works",
			&types.AttributeValueMemberB{Value: []byte("Republished LULD Price Band")}, Financial_Quotes_RepublishedLULDPriceBand),
		Entry("Value is []bytes, On-Demand Auction - Works",
			&types.AttributeValueMemberB{Value: []byte("On-Demand Auction")}, Financial_Quotes_OnDemandAuction),
		Entry("Value is []bytes, Cash-Only Settlement - Works",
			&types.AttributeValueMemberB{Value: []byte("Cash-Only Settlement")}, Financial_Quotes_CashOnlySettlement),
		Entry("Value is []bytes, Next-Day Settlement - Works",
			&types.AttributeValueMemberB{Value: []byte("Next-Day Settlement")}, Financial_Quotes_NextDaySettlement),
		Entry("Value is []bytes, LULD Trading Pause - Works",
			&types.AttributeValueMemberB{Value: []byte("LULD Trading Pause")}, Financial_Quotes_LULDTradingPause),
		Entry("Value is []bytes, Slow Due LRP, Bid, Ask - Works",
			&types.AttributeValueMemberB{Value: []byte("Slow Due LRP, Bid, Ask")}, Financial_Quotes_SlowDuelRPBidAsk),
		Entry("Value is []bytes, Corrected Price - Works",
			&types.AttributeValueMemberB{Value: []byte("Corrected Price")}, Financial_Quotes_CorrectedPrice),
		Entry("Value is []bytes, SIP-Generated - Works",
			&types.AttributeValueMemberB{Value: []byte("SIP-Generated")}, Financial_Quotes_SIPGenerated),
		Entry("Value is []bytes, Crossed Market - Works",
			&types.AttributeValueMemberB{Value: []byte("Crossed Market")}, Financial_Quotes_CrossedMarket),
		Entry("Value is []bytes, Locked Market - Works",
			&types.AttributeValueMemberB{Value: []byte("Locked Market")}, Financial_Quotes_LockedMarket),
		Entry("Value is []bytes, Depth on Offer Side - Works",
			&types.AttributeValueMemberB{Value: []byte("Depth on Offer Side")}, Financial_Quotes_DepthOnOfferSide),
		Entry("Value is []bytes, Depth on Bid Side - Works",
			&types.AttributeValueMemberB{Value: []byte("Depth on Bid Side")}, Financial_Quotes_DepthOnBidSide),
		Entry("Value is []bytes, Depth on Bid and Offer - Works",
			&types.AttributeValueMemberB{Value: []byte("Depth on Bid and Offer")}, Financial_Quotes_DepthOnBidAndOffer),
		Entry("Value is []bytes, Pre-Opening Indication - Works",
			&types.AttributeValueMemberB{Value: []byte("Pre-Opening Indication")}, Financial_Quotes_PreOpeningIndication),
		Entry("Value is []bytes, Syndicate Bid - Works",
			&types.AttributeValueMemberB{Value: []byte("Syndicate Bid")}, Financial_Quotes_SyndicateBid),
		Entry("Value is []bytes, Pre-Syndicate Bid - Works",
			&types.AttributeValueMemberB{Value: []byte("Pre-Syndicate Bid")}, Financial_Quotes_PreSyndicateBid),
		Entry("Value is []bytes, Penalty Bid - Works",
			&types.AttributeValueMemberB{Value: []byte("Penalty Bid")}, Financial_Quotes_PenaltyBid),
		Entry("Value is []bytes, CQS-Generated - Works",
			&types.AttributeValueMemberB{Value: []byte("CQS-Generated")}, Financial_Quotes_CQSGenerated),
		Entry("Value is []bytes, Regular - Works",
			&types.AttributeValueMemberB{Value: []byte("Regular")}, Financial_Quotes_Regular),
		Entry("Value is []bytes, RegularTwoSidedOpen - Works",
			&types.AttributeValueMemberB{Value: []byte("RegularTwoSidedOpen")}, Financial_Quotes_RegularTwoSidedOpen),
		Entry("Value is []bytes, RegularOneSidedOpen - Works",
			&types.AttributeValueMemberB{Value: []byte("RegularOneSidedOpen")}, Financial_Quotes_RegularOneSidedOpen),
		Entry("Value is []bytes, SlowAsk - Works",
			&types.AttributeValueMemberB{Value: []byte("SlowAsk")}, Financial_Quotes_SlowAsk),
		Entry("Value is []bytes, SlowBid - Works",
			&types.AttributeValueMemberB{Value: []byte("SlowBid")}, Financial_Quotes_SlowBid),
		Entry("Value is []bytes, SlowBidAsk - Works",
			&types.AttributeValueMemberB{Value: []byte("SlowBidAsk")}, Financial_Quotes_SlowBidAsk),
		Entry("Value is []bytes, SlowDueLRPBid - Works",
			&types.AttributeValueMemberB{Value: []byte("SlowDueLRPBid")}, Financial_Quotes_SlowDueLRPBid),
		Entry("Value is []bytes, SlowDueLRPAsk - Works",
			&types.AttributeValueMemberB{Value: []byte("SlowDueLRPAsk")}, Financial_Quotes_SlowDueLRPAsk),
		Entry("Value is []bytes, SlowDueNYSELRP - Works",
			&types.AttributeValueMemberB{Value: []byte("SlowDueNYSELRP")}, Financial_Quotes_SlowDueNYSELRP),
		Entry("Value is []bytes, SlowDueSetSlowListBidAsk - Works",
			&types.AttributeValueMemberB{Value: []byte("SlowDueSetSlowListBidAsk")}, Financial_Quotes_SlowDueSetSlowListBidAsk),
		Entry("Value is []bytes, ManualAskAutomatedBid - Works",
			&types.AttributeValueMemberB{Value: []byte("ManualAskAutomatedBid")}, Financial_Quotes_ManualAskAutomatedBid),
		Entry("Value is []bytes, ManualBidAutomatedAsk - Works",
			&types.AttributeValueMemberB{Value: []byte("ManualBidAutomatedAsk")}, Financial_Quotes_ManualBidAutomatedAsk),
		Entry("Value is []bytes, ManualBidAndAsk - Works",
			&types.AttributeValueMemberB{Value: []byte("ManualBidAndAsk")}, Financial_Quotes_ManualBidAndAsk),
		Entry("Value is []bytes, Opening - Works",
			&types.AttributeValueMemberB{Value: []byte("Opening")}, Financial_Quotes_Opening),
		Entry("Value is []bytes, Closing - Works",
			&types.AttributeValueMemberB{Value: []byte("Closing")}, Financial_Quotes_Closing),
		Entry("Value is []bytes, Closed - Works",
			&types.AttributeValueMemberB{Value: []byte("Closed")}, Financial_Quotes_Closed),
		Entry("Value is []bytes, Resume - Works",
			&types.AttributeValueMemberB{Value: []byte("Resume")}, Financial_Quotes_Resume),
		Entry("Value is []bytes, FastTrading - Works",
			&types.AttributeValueMemberB{Value: []byte("FastTrading")}, Financial_Quotes_FastTrading),
		Entry("Value is []bytes, TradingRangeIndicated - Works",
			&types.AttributeValueMemberB{Value: []byte("TradingRangeIndicated")}, Financial_Quotes_TradingRangeIndicated),
		Entry("Value is []bytes, MarketMakerQuotesClosed - Works",
			&types.AttributeValueMemberB{Value: []byte("MarketMakerQuotesClosed")}, Financial_Quotes_MarketMakerQuotesClosed),
		Entry("Value is []bytes, NonFirm - Works",
			&types.AttributeValueMemberB{Value: []byte("NonFirm")}, Financial_Quotes_NonFirm),
		Entry("Value is []bytes, NewsDissemination - Works",
			&types.AttributeValueMemberB{Value: []byte("NewsDissemination")}, Financial_Quotes_NewsDissemination),
		Entry("Value is []bytes, OrderInflux - Works",
			&types.AttributeValueMemberB{Value: []byte("OrderInflux")}, Financial_Quotes_OrderInflux),
		Entry("Value is []bytes, OrderImbalance - Works",
			&types.AttributeValueMemberB{Value: []byte("OrderImbalance")}, Financial_Quotes_OrderImbalance),
		Entry("Value is []bytes, DueToRelatedSecurityNewsDissemination - Works",
			&types.AttributeValueMemberB{Value: []byte("DueToRelatedSecurityNewsDissemination")}, Financial_Quotes_DueToRelatedSecurityNewsDissemination),
		Entry("Value is []bytes, DueToRelatedSecurityNewsPending - Works",
			&types.AttributeValueMemberB{Value: []byte("DueToRelatedSecurityNewsPending")}, Financial_Quotes_DueToRelatedSecurityNewsPending),
		Entry("Value is []bytes, AdditionalInformation - Works",
			&types.AttributeValueMemberB{Value: []byte("AdditionalInformation")}, Financial_Quotes_AdditionalInformation),
		Entry("Value is []bytes, NewsPending - Works",
			&types.AttributeValueMemberB{Value: []byte("NewsPending")}, Financial_Quotes_NewsPending),
		Entry("Value is []bytes, AdditionalInformationDueToRelatedSecurity - Works",
			&types.AttributeValueMemberB{Value: []byte("AdditionalInformationDueToRelatedSecurity")}, Financial_Quotes_AdditionalInformationDueToRelatedSecurity),
		Entry("Value is []bytes, DueToRelatedSecurity - Works",
			&types.AttributeValueMemberB{Value: []byte("DueToRelatedSecurity")}, Financial_Quotes_DueToRelatedSecurity),
		Entry("Value is []bytes, InViewOfCommon - Works",
			&types.AttributeValueMemberB{Value: []byte("InViewOfCommon")}, Financial_Quotes_InViewOfCommon),
		Entry("Value is []bytes, EquipmentChangeover - Works",
			&types.AttributeValueMemberB{Value: []byte("EquipmentChangeover")}, Financial_Quotes_EquipmentChangeover),
		Entry("Value is []bytes, NoOpenNoResponse - Works",
			&types.AttributeValueMemberB{Value: []byte("NoOpenNoResponse")}, Financial_Quotes_NoOpenNoResponse),
		Entry("Value is []bytes, SubPennyTrading - Works",
			&types.AttributeValueMemberB{Value: []byte("SubPennyTrading")}, Financial_Quotes_SubPennyTrading),
		Entry("Value is []bytes, AutomatedBidNoOfferNoBid - Works",
			&types.AttributeValueMemberB{Value: []byte("AutomatedBidNoOfferNoBid")}, Financial_Quotes_AutomatedBidNoOfferNoBid),
		Entry("Value is []bytes, LULDPPriceBand - Works",
			&types.AttributeValueMemberB{Value: []byte("LULDPPriceBand")}, Financial_Quotes_LULDPPriceBand),
		Entry("Value is []bytes, MarketWideCircuitBreakerLevel1 - Works",
			&types.AttributeValueMemberB{Value: []byte("MarketWideCircuitBreakerLevel1")}, Financial_Quotes_MarketWideCircuitBreakerLevel1),
		Entry("Value is []bytes, MarketWideCircuitBreakerLevel2 - Works",
			&types.AttributeValueMemberB{Value: []byte("MarketWideCircuitBreakerLevel2")}, Financial_Quotes_MarketWideCircuitBreakerLevel2),
		Entry("Value is []bytes, MarketWideCircuitBreakerLevel3 - Works",
			&types.AttributeValueMemberB{Value: []byte("MarketWideCircuitBreakerLevel3")}, Financial_Quotes_MarketWideCircuitBreakerLevel3),
		Entry("Value is []bytes, RepublishedLULDPriceBand - Works",
			&types.AttributeValueMemberB{Value: []byte("RepublishedLULDPriceBand")}, Financial_Quotes_RepublishedLULDPriceBand),
		Entry("Value is []bytes, OnDemandAuction - Works",
			&types.AttributeValueMemberB{Value: []byte("OnDemandAuction")}, Financial_Quotes_OnDemandAuction),
		Entry("Value is []bytes, CashOnlySettlement - Works",
			&types.AttributeValueMemberB{Value: []byte("CashOnlySettlement")}, Financial_Quotes_CashOnlySettlement),
		Entry("Value is []bytes, NextDaySettlement - Works",
			&types.AttributeValueMemberB{Value: []byte("NextDaySettlement")}, Financial_Quotes_NextDaySettlement),
		Entry("Value is []bytes, LULDTradingPause - Works",
			&types.AttributeValueMemberB{Value: []byte("LULDTradingPause")}, Financial_Quotes_LULDTradingPause),
		Entry("Value is []bytes, SlowDuelRPBidAsk - Works",
			&types.AttributeValueMemberB{Value: []byte("SlowDuelRPBidAsk")}, Financial_Quotes_SlowDuelRPBidAsk),
		Entry("Value is []bytes, Cancel - Works",
			&types.AttributeValueMemberB{Value: []byte("Cancel")}, Financial_Quotes_Cancel),
		Entry("Value is []bytes, CorrectedPrice - Works",
			&types.AttributeValueMemberB{Value: []byte("CorrectedPrice")}, Financial_Quotes_CorrectedPrice),
		Entry("Value is []bytes, SIPGenerated - Works",
			&types.AttributeValueMemberB{Value: []byte("SIPGenerated")}, Financial_Quotes_SIPGenerated),
		Entry("Value is []bytes, Unknown - Works",
			&types.AttributeValueMemberB{Value: []byte("Unknown")}, Financial_Quotes_Unknown),
		Entry("Value is []bytes, CrossedMarket - Works",
			&types.AttributeValueMemberB{Value: []byte("CrossedMarket")}, Financial_Quotes_CrossedMarket),
		Entry("Value is []bytes, LockedMarket - Works",
			&types.AttributeValueMemberB{Value: []byte("LockedMarket")}, Financial_Quotes_LockedMarket),
		Entry("Value is []bytes, DepthOnOfferSide - Works",
			&types.AttributeValueMemberB{Value: []byte("DepthOnOfferSide")}, Financial_Quotes_DepthOnOfferSide),
		Entry("Value is []bytes, DepthOnBidSide - Works",
			&types.AttributeValueMemberB{Value: []byte("DepthOnBidSide")}, Financial_Quotes_DepthOnBidSide),
		Entry("Value is []bytes, DepthOnBidAndOffer - Works",
			&types.AttributeValueMemberB{Value: []byte("DepthOnBidAndOffer")}, Financial_Quotes_DepthOnBidAndOffer),
		Entry("Value is []bytes, PreOpeningIndication - Works",
			&types.AttributeValueMemberB{Value: []byte("PreOpeningIndication")}, Financial_Quotes_PreOpeningIndication),
		Entry("Value is []bytes, SyndicateBid - Works",
			&types.AttributeValueMemberB{Value: []byte("SyndicateBid")}, Financial_Quotes_SyndicateBid),
		Entry("Value is []bytes, PreSyndicateBid - Works",
			&types.AttributeValueMemberB{Value: []byte("PreSyndicateBid")}, Financial_Quotes_PreSyndicateBid),
		Entry("Value is []bytes, PenaltyBid - Works",
			&types.AttributeValueMemberB{Value: []byte("PenaltyBid")}, Financial_Quotes_PenaltyBid),
		Entry("Value is []bytes, CQSGenerated - Works",
			&types.AttributeValueMemberB{Value: []byte("CQSGenerated")}, Financial_Quotes_CQSGenerated),
		Entry("Value is []bytes, Invalid - Works",
			&types.AttributeValueMemberB{Value: []byte("Invalid")}, Financial_Quotes_Invalid),
		Entry("Value is numeric, 0 - Works",
			&types.AttributeValueMemberN{Value: "0"}, Financial_Quotes_Regular),
		Entry("Value is numeric, 1 - Works",
			&types.AttributeValueMemberN{Value: "1"}, Financial_Quotes_RegularTwoSidedOpen),
		Entry("Value is numeric, 2 - Works",
			&types.AttributeValueMemberN{Value: "2"}, Financial_Quotes_RegularOneSidedOpen),
		Entry("Value is numeric, 3 - Works",
			&types.AttributeValueMemberN{Value: "3"}, Financial_Quotes_SlowAsk),
		Entry("Value is numeric, 4 - Works",
			&types.AttributeValueMemberN{Value: "4"}, Financial_Quotes_SlowBid),
		Entry("Value is numeric, 5 - Works",
			&types.AttributeValueMemberN{Value: "5"}, Financial_Quotes_SlowBidAsk),
		Entry("Value is numeric, 6 - Works",
			&types.AttributeValueMemberN{Value: "6"}, Financial_Quotes_SlowDueLRPBid),
		Entry("Value is numeric, 7 - Works",
			&types.AttributeValueMemberN{Value: "7"}, Financial_Quotes_SlowDueLRPAsk),
		Entry("Value is numeric, 8 - Works",
			&types.AttributeValueMemberN{Value: "8"}, Financial_Quotes_SlowDueNYSELRP),
		Entry("Value is numeric, 9 - Works",
			&types.AttributeValueMemberN{Value: "9"}, Financial_Quotes_SlowDueSetSlowListBidAsk),
		Entry("Value is numeric, 10 - Works",
			&types.AttributeValueMemberN{Value: "10"}, Financial_Quotes_ManualAskAutomatedBid),
		Entry("Value is numeric, 11 - Works",
			&types.AttributeValueMemberN{Value: "11"}, Financial_Quotes_ManualBidAutomatedAsk),
		Entry("Value is numeric, 12 - Works",
			&types.AttributeValueMemberN{Value: "12"}, Financial_Quotes_ManualBidAndAsk),
		Entry("Value is numeric, 13 - Works",
			&types.AttributeValueMemberN{Value: "13"}, Financial_Quotes_Opening),
		Entry("Value is numeric, 14 - Works",
			&types.AttributeValueMemberN{Value: "14"}, Financial_Quotes_Closing),
		Entry("Value is numeric, 15 - Works",
			&types.AttributeValueMemberN{Value: "15"}, Financial_Quotes_Closed),
		Entry("Value is numeric, 16 - Works",
			&types.AttributeValueMemberN{Value: "16"}, Financial_Quotes_Resume),
		Entry("Value is numeric, 17 - Works",
			&types.AttributeValueMemberN{Value: "17"}, Financial_Quotes_FastTrading),
		Entry("Value is numeric, 18 - Works",
			&types.AttributeValueMemberN{Value: "18"}, Financial_Quotes_TradingRangeIndicated),
		Entry("Value is numeric, 19 - Works",
			&types.AttributeValueMemberN{Value: "19"}, Financial_Quotes_MarketMakerQuotesClosed),
		Entry("Value is numeric, 20 - Works",
			&types.AttributeValueMemberN{Value: "20"}, Financial_Quotes_NonFirm),
		Entry("Value is numeric, 21 - Works",
			&types.AttributeValueMemberN{Value: "21"}, Financial_Quotes_NewsDissemination),
		Entry("Value is numeric, 22 - Works",
			&types.AttributeValueMemberN{Value: "22"}, Financial_Quotes_OrderInflux),
		Entry("Value is numeric, 23 - Works",
			&types.AttributeValueMemberN{Value: "23"}, Financial_Quotes_OrderImbalance),
		Entry("Value is numeric, 24 - Works",
			&types.AttributeValueMemberN{Value: "24"}, Financial_Quotes_DueToRelatedSecurityNewsDissemination),
		Entry("Value is numeric, 25 - Works",
			&types.AttributeValueMemberN{Value: "25"}, Financial_Quotes_DueToRelatedSecurityNewsPending),
		Entry("Value is numeric, 26 - Works",
			&types.AttributeValueMemberN{Value: "26"}, Financial_Quotes_AdditionalInformation),
		Entry("Value is numeric, 27 - Works",
			&types.AttributeValueMemberN{Value: "27"}, Financial_Quotes_NewsPending),
		Entry("Value is numeric, 28 - Works",
			&types.AttributeValueMemberN{Value: "28"}, Financial_Quotes_AdditionalInformationDueToRelatedSecurity),
		Entry("Value is numeric, 29 - Works",
			&types.AttributeValueMemberN{Value: "29"}, Financial_Quotes_DueToRelatedSecurity),
		Entry("Value is numeric, 30 - Works",
			&types.AttributeValueMemberN{Value: "30"}, Financial_Quotes_InViewOfCommon),
		Entry("Value is numeric, 31 - Works",
			&types.AttributeValueMemberN{Value: "31"}, Financial_Quotes_EquipmentChangeover),
		Entry("Value is numeric, 32 - Works",
			&types.AttributeValueMemberN{Value: "32"}, Financial_Quotes_NoOpenNoResponse),
		Entry("Value is numeric, 33 - Works",
			&types.AttributeValueMemberN{Value: "33"}, Financial_Quotes_SubPennyTrading),
		Entry("Value is numeric, 34 - Works",
			&types.AttributeValueMemberN{Value: "34"}, Financial_Quotes_AutomatedBidNoOfferNoBid),
		Entry("Value is numeric, 35 - Works",
			&types.AttributeValueMemberN{Value: "35"}, Financial_Quotes_LULDPPriceBand),
		Entry("Value is numeric, 36 - Works",
			&types.AttributeValueMemberN{Value: "36"}, Financial_Quotes_MarketWideCircuitBreakerLevel1),
		Entry("Value is numeric, 37 - Works",
			&types.AttributeValueMemberN{Value: "37"}, Financial_Quotes_MarketWideCircuitBreakerLevel2),
		Entry("Value is numeric, 38 - Works",
			&types.AttributeValueMemberN{Value: "38"}, Financial_Quotes_MarketWideCircuitBreakerLevel3),
		Entry("Value is numeric, 39 - Works",
			&types.AttributeValueMemberN{Value: "39"}, Financial_Quotes_RepublishedLULDPriceBand),
		Entry("Value is numeric, 40 - Works",
			&types.AttributeValueMemberN{Value: "40"}, Financial_Quotes_OnDemandAuction),
		Entry("Value is numeric, 41 - Works",
			&types.AttributeValueMemberN{Value: "41"}, Financial_Quotes_CashOnlySettlement),
		Entry("Value is numeric, 42 - Works",
			&types.AttributeValueMemberN{Value: "42"}, Financial_Quotes_NextDaySettlement),
		Entry("Value is numeric, 43 - Works",
			&types.AttributeValueMemberN{Value: "43"}, Financial_Quotes_LULDTradingPause),
		Entry("Value is numeric, 71 - Works",
			&types.AttributeValueMemberN{Value: "71"}, Financial_Quotes_SlowDuelRPBidAsk),
		Entry("Value is numeric, 80 - Works",
			&types.AttributeValueMemberN{Value: "80"}, Financial_Quotes_Cancel),
		Entry("Value is numeric, 81 - Works",
			&types.AttributeValueMemberN{Value: "81"}, Financial_Quotes_CorrectedPrice),
		Entry("Value is numeric, 82 - Works",
			&types.AttributeValueMemberN{Value: "82"}, Financial_Quotes_SIPGenerated),
		Entry("Value is numeric, 83 - Works",
			&types.AttributeValueMemberN{Value: "83"}, Financial_Quotes_Unknown),
		Entry("Value is numeric, 84 - Works",
			&types.AttributeValueMemberN{Value: "84"}, Financial_Quotes_CrossedMarket),
		Entry("Value is numeric, 85 - Works",
			&types.AttributeValueMemberN{Value: "85"}, Financial_Quotes_LockedMarket),
		Entry("Value is numeric, 86 - Works",
			&types.AttributeValueMemberN{Value: "86"}, Financial_Quotes_DepthOnOfferSide),
		Entry("Value is numeric, 87 - Works",
			&types.AttributeValueMemberN{Value: "87"}, Financial_Quotes_DepthOnBidSide),
		Entry("Value is numeric, 88 - Works",
			&types.AttributeValueMemberN{Value: "88"}, Financial_Quotes_DepthOnBidAndOffer),
		Entry("Value is numeric, 89 - Works",
			&types.AttributeValueMemberN{Value: "89"}, Financial_Quotes_PreOpeningIndication),
		Entry("Value is numeric, 90 - Works",
			&types.AttributeValueMemberN{Value: "90"}, Financial_Quotes_SyndicateBid),
		Entry("Value is numeric, 91 - Works",
			&types.AttributeValueMemberN{Value: "91"}, Financial_Quotes_PreSyndicateBid),
		Entry("Value is numeric, 92 - Works",
			&types.AttributeValueMemberN{Value: "92"}, Financial_Quotes_PenaltyBid),
		Entry("Value is numeric, 94 - Works",
			&types.AttributeValueMemberN{Value: "94"}, Financial_Quotes_CQSGenerated),
		Entry("Value is numeric, 999 - Works",
			&types.AttributeValueMemberN{Value: "999"}, Financial_Quotes_Invalid),
		Entry("Value is NULL - Works", new(types.AttributeValueMemberNULL), Financial_Quotes_Condition(0)),
		Entry("Value is string, Regular, Two-Sided Open - Works",
			&types.AttributeValueMemberS{Value: "Regular, Two-Sided Open"}, Financial_Quotes_RegularTwoSidedOpen),
		Entry("Value is string, Regular, One-Sided Open - Works",
			&types.AttributeValueMemberS{Value: "Regular, One-Sided Open"}, Financial_Quotes_RegularOneSidedOpen),
		Entry("Value is string, Slow Ask - Works",
			&types.AttributeValueMemberS{Value: "Slow Ask"}, Financial_Quotes_SlowAsk),
		Entry("Value is string, Slow Bid - Works",
			&types.AttributeValueMemberS{Value: "Slow Bid"}, Financial_Quotes_SlowBid),
		Entry("Value is string, Slow Bid, Ask - Works",
			&types.AttributeValueMemberS{Value: "Slow Bid, Ask"}, Financial_Quotes_SlowBidAsk),
		Entry("Value is string, Slow Due, LRP Bid - Works",
			&types.AttributeValueMemberS{Value: "Slow Due, LRP Bid"}, Financial_Quotes_SlowDueLRPBid),
		Entry("Value is string, Slow Due, LRP Ask - Works",
			&types.AttributeValueMemberS{Value: "Slow Due, LRP Ask"}, Financial_Quotes_SlowDueLRPAsk),
		Entry("Value is string, Slow Due, NYSE LRP - Works",
			&types.AttributeValueMemberS{Value: "Slow Due, NYSE LRP"}, Financial_Quotes_SlowDueNYSELRP),
		Entry("Value is string, Slow Due Set, Slow List, Bid, Ask - Works",
			&types.AttributeValueMemberS{Value: "Slow Due Set, Slow List, Bid, Ask"}, Financial_Quotes_SlowDueSetSlowListBidAsk),
		Entry("Value is string, Manual Ask, Automated Bid - Works",
			&types.AttributeValueMemberS{Value: "Manual Ask, Automated Bid"}, Financial_Quotes_ManualAskAutomatedBid),
		Entry("Value is string, Manual Bid, Automated Ask - Works",
			&types.AttributeValueMemberS{Value: "Manual Bid, Automated Ask"}, Financial_Quotes_ManualBidAutomatedAsk),
		Entry("Value is string, Manual Bid and Ask - Works",
			&types.AttributeValueMemberS{Value: "Manual Bid and Ask"}, Financial_Quotes_ManualBidAndAsk),
		Entry("Value is string, Fast Trading - Works",
			&types.AttributeValueMemberS{Value: "Fast Trading"}, Financial_Quotes_FastTrading),
		Entry("Value is string, Tading Range Indicated - Works",
			&types.AttributeValueMemberS{Value: "Tading Range Indicated"}, Financial_Quotes_TradingRangeIndicated),
		Entry("Value is string, Market-Maker Quotes Closed - Works",
			&types.AttributeValueMemberS{Value: "Market-Maker Quotes Closed"}, Financial_Quotes_MarketMakerQuotesClosed),
		Entry("Value is string, Non-Firm - Works",
			&types.AttributeValueMemberS{Value: "Non-Firm"}, Financial_Quotes_NonFirm),
		Entry("Value is string, News Dissemination - Works",
			&types.AttributeValueMemberS{Value: "News Dissemination"}, Financial_Quotes_NewsDissemination),
		Entry("Value is string, Order Influx - Works",
			&types.AttributeValueMemberS{Value: "Order Influx"}, Financial_Quotes_OrderInflux),
		Entry("Value is string, Order Imbalance - Works",
			&types.AttributeValueMemberS{Value: "Order Imbalance"}, Financial_Quotes_OrderImbalance),
		Entry("Value is string, Due to Related Security, News Dissemination - Works",
			&types.AttributeValueMemberS{Value: "Due to Related Security, News Dissemination"}, Financial_Quotes_DueToRelatedSecurityNewsDissemination),
		Entry("Value is string, Due to Related Security, News Pending - Works",
			&types.AttributeValueMemberS{Value: "Due to Related Security, News Pending"}, Financial_Quotes_DueToRelatedSecurityNewsPending),
		Entry("Value is string, Additional Information - Works",
			&types.AttributeValueMemberS{Value: "Additional Information"}, Financial_Quotes_AdditionalInformation),
		Entry("Value is string, News Pending - Works",
			&types.AttributeValueMemberS{Value: "News Pending"}, Financial_Quotes_NewsPending),
		Entry("Value is string, Additional Information Due to Related Security - Works",
			&types.AttributeValueMemberS{Value: "Additional Information Due to Related Security"}, Financial_Quotes_AdditionalInformationDueToRelatedSecurity),
		Entry("Value is string, Due to Related Security - Works",
			&types.AttributeValueMemberS{Value: "Due to Related Security"}, Financial_Quotes_DueToRelatedSecurity),
		Entry("Value is string, In View of Common - Works",
			&types.AttributeValueMemberS{Value: "In View of Common"}, Financial_Quotes_InViewOfCommon),
		Entry("Value is string, Equipment Changeover - Works",
			&types.AttributeValueMemberS{Value: "Equipment Changeover"}, Financial_Quotes_EquipmentChangeover),
		Entry("Value is string, No Open, No Response - Works",
			&types.AttributeValueMemberS{Value: "No Open, No Response"}, Financial_Quotes_NoOpenNoResponse),
		Entry("Value is string, Sub-Penny Trading - Works",
			&types.AttributeValueMemberS{Value: "Sub-Penny Trading"}, Financial_Quotes_SubPennyTrading),
		Entry("Value is string, Automated Bid; No Offer, No Bid - Works",
			&types.AttributeValueMemberS{Value: "Automated Bid; No Offer, No Bid"}, Financial_Quotes_AutomatedBidNoOfferNoBid),
		Entry("Value is string, LULD Price Band - Works",
			&types.AttributeValueMemberS{Value: "LULD Price Band"}, Financial_Quotes_LULDPPriceBand),
		Entry("Value is string, Market-Wide Circuit Breaker, Level 1 - Works",
			&types.AttributeValueMemberS{Value: "Market-Wide Circuit Breaker, Level 1"}, Financial_Quotes_MarketWideCircuitBreakerLevel1),
		Entry("Value is string, Market-Wide Circuit Breaker, Level 2 - Works",
			&types.AttributeValueMemberS{Value: "Market-Wide Circuit Breaker, Level 2"}, Financial_Quotes_MarketWideCircuitBreakerLevel2),
		Entry("Value is string, Market-Wide Circuit Breaker, Level 3 - Works",
			&types.AttributeValueMemberS{Value: "Market-Wide Circuit Breaker, Level 3"}, Financial_Quotes_MarketWideCircuitBreakerLevel3),
		Entry("Value is string, Republished LULD Price Band - Works",
			&types.AttributeValueMemberS{Value: "Republished LULD Price Band"}, Financial_Quotes_RepublishedLULDPriceBand),
		Entry("Value is string, On-Demand Auction - Works",
			&types.AttributeValueMemberS{Value: "On-Demand Auction"}, Financial_Quotes_OnDemandAuction),
		Entry("Value is string, Cash-Only Settlement - Works",
			&types.AttributeValueMemberS{Value: "Cash-Only Settlement"}, Financial_Quotes_CashOnlySettlement),
		Entry("Value is string, Next-Day Settlement - Works",
			&types.AttributeValueMemberS{Value: "Next-Day Settlement"}, Financial_Quotes_NextDaySettlement),
		Entry("Value is string, LULD Trading Pause - Works",
			&types.AttributeValueMemberS{Value: "LULD Trading Pause"}, Financial_Quotes_LULDTradingPause),
		Entry("Value is string, Slow Due LRP, Bid, Ask - Works",
			&types.AttributeValueMemberS{Value: "Slow Due LRP, Bid, Ask"}, Financial_Quotes_SlowDuelRPBidAsk),
		Entry("Value is string, Corrected Price - Works",
			&types.AttributeValueMemberS{Value: "Corrected Price"}, Financial_Quotes_CorrectedPrice),
		Entry("Value is string, SIP-Generated - Works",
			&types.AttributeValueMemberS{Value: "SIP-Generated"}, Financial_Quotes_SIPGenerated),
		Entry("Value is string, Crossed Market - Works",
			&types.AttributeValueMemberS{Value: "Crossed Market"}, Financial_Quotes_CrossedMarket),
		Entry("Value is string, Locked Market - Works",
			&types.AttributeValueMemberS{Value: "Locked Market"}, Financial_Quotes_LockedMarket),
		Entry("Value is string, Depth on Offer Side - Works",
			&types.AttributeValueMemberS{Value: "Depth on Offer Side"}, Financial_Quotes_DepthOnOfferSide),
		Entry("Value is string, Depth on Bid Side - Works",
			&types.AttributeValueMemberS{Value: "Depth on Bid Side"}, Financial_Quotes_DepthOnBidSide),
		Entry("Value is string, Depth on Bid and Offer - Works",
			&types.AttributeValueMemberS{Value: "Depth on Bid and Offer"}, Financial_Quotes_DepthOnBidAndOffer),
		Entry("Value is string, Pre-Opening Indication - Works",
			&types.AttributeValueMemberS{Value: "Pre-Opening Indication"}, Financial_Quotes_PreOpeningIndication),
		Entry("Value is string, Syndicate Bid - Works",
			&types.AttributeValueMemberS{Value: "Syndicate Bid"}, Financial_Quotes_SyndicateBid),
		Entry("Value is string, Pre-Syndicate Bid - Works",
			&types.AttributeValueMemberS{Value: "Pre-Syndicate Bid"}, Financial_Quotes_PreSyndicateBid),
		Entry("Value is string, Penalty Bid - Works",
			&types.AttributeValueMemberS{Value: "Penalty Bid"}, Financial_Quotes_PenaltyBid),
		Entry("Value is string, CQS-Generated - Works",
			&types.AttributeValueMemberS{Value: "CQS-Generated"}, Financial_Quotes_CQSGenerated),
		Entry("Value is string, Regular - Works",
			&types.AttributeValueMemberS{Value: "Regular"}, Financial_Quotes_Regular),
		Entry("Value is string, RegularTwoSidedOpen - Works",
			&types.AttributeValueMemberS{Value: "RegularTwoSidedOpen"}, Financial_Quotes_RegularTwoSidedOpen),
		Entry("Value is string, RegularOneSidedOpen - Works",
			&types.AttributeValueMemberS{Value: "RegularOneSidedOpen"}, Financial_Quotes_RegularOneSidedOpen),
		Entry("Value is string, SlowAsk - Works",
			&types.AttributeValueMemberS{Value: "SlowAsk"}, Financial_Quotes_SlowAsk),
		Entry("Value is string, SlowBid - Works",
			&types.AttributeValueMemberS{Value: "SlowBid"}, Financial_Quotes_SlowBid),
		Entry("Value is string, SlowBidAsk - Works",
			&types.AttributeValueMemberS{Value: "SlowBidAsk"}, Financial_Quotes_SlowBidAsk),
		Entry("Value is string, SlowDueLRPBid - Works",
			&types.AttributeValueMemberS{Value: "SlowDueLRPBid"}, Financial_Quotes_SlowDueLRPBid),
		Entry("Value is string, SlowDueLRPAsk - Works",
			&types.AttributeValueMemberS{Value: "SlowDueLRPAsk"}, Financial_Quotes_SlowDueLRPAsk),
		Entry("Value is string, SlowDueNYSELRP - Works",
			&types.AttributeValueMemberS{Value: "SlowDueNYSELRP"}, Financial_Quotes_SlowDueNYSELRP),
		Entry("Value is string, SlowDueSetSlowListBidAsk - Works",
			&types.AttributeValueMemberS{Value: "SlowDueSetSlowListBidAsk"}, Financial_Quotes_SlowDueSetSlowListBidAsk),
		Entry("Value is string, ManualAskAutomatedBid - Works",
			&types.AttributeValueMemberS{Value: "ManualAskAutomatedBid"}, Financial_Quotes_ManualAskAutomatedBid),
		Entry("Value is string, ManualBidAutomatedAsk - Works",
			&types.AttributeValueMemberS{Value: "ManualBidAutomatedAsk"}, Financial_Quotes_ManualBidAutomatedAsk),
		Entry("Value is string, ManualBidAndAsk - Works",
			&types.AttributeValueMemberS{Value: "ManualBidAndAsk"}, Financial_Quotes_ManualBidAndAsk),
		Entry("Value is string, Opening - Works",
			&types.AttributeValueMemberS{Value: "Opening"}, Financial_Quotes_Opening),
		Entry("Value is string, Closing - Works",
			&types.AttributeValueMemberS{Value: "Closing"}, Financial_Quotes_Closing),
		Entry("Value is string, Closed - Works",
			&types.AttributeValueMemberS{Value: "Closed"}, Financial_Quotes_Closed),
		Entry("Value is string, Resume - Works",
			&types.AttributeValueMemberS{Value: "Resume"}, Financial_Quotes_Resume),
		Entry("Value is string, FastTrading - Works",
			&types.AttributeValueMemberS{Value: "FastTrading"}, Financial_Quotes_FastTrading),
		Entry("Value is string, TradingRangeIndicated - Works",
			&types.AttributeValueMemberS{Value: "TradingRangeIndicated"}, Financial_Quotes_TradingRangeIndicated),
		Entry("Value is string, MarketMakerQuotesClosed - Works",
			&types.AttributeValueMemberS{Value: "MarketMakerQuotesClosed"}, Financial_Quotes_MarketMakerQuotesClosed),
		Entry("Value is string, NonFirm - Works",
			&types.AttributeValueMemberS{Value: "NonFirm"}, Financial_Quotes_NonFirm),
		Entry("Value is string, NewsDissemination - Works",
			&types.AttributeValueMemberS{Value: "NewsDissemination"}, Financial_Quotes_NewsDissemination),
		Entry("Value is string, OrderInflux - Works",
			&types.AttributeValueMemberS{Value: "OrderInflux"}, Financial_Quotes_OrderInflux),
		Entry("Value is string, OrderImbalance - Works",
			&types.AttributeValueMemberS{Value: "OrderImbalance"}, Financial_Quotes_OrderImbalance),
		Entry("Value is string, DueToRelatedSecurityNewsDissemination - Works",
			&types.AttributeValueMemberS{Value: "DueToRelatedSecurityNewsDissemination"}, Financial_Quotes_DueToRelatedSecurityNewsDissemination),
		Entry("Value is string, DueToRelatedSecurityNewsPending - Works",
			&types.AttributeValueMemberS{Value: "DueToRelatedSecurityNewsPending"}, Financial_Quotes_DueToRelatedSecurityNewsPending),
		Entry("Value is string, AdditionalInformation - Works",
			&types.AttributeValueMemberS{Value: "AdditionalInformation"}, Financial_Quotes_AdditionalInformation),
		Entry("Value is string, NewsPending - Works",
			&types.AttributeValueMemberS{Value: "NewsPending"}, Financial_Quotes_NewsPending),
		Entry("Value is string, AdditionalInformationDueToRelatedSecurity - Works",
			&types.AttributeValueMemberS{Value: "AdditionalInformationDueToRelatedSecurity"}, Financial_Quotes_AdditionalInformationDueToRelatedSecurity),
		Entry("Value is string, DueToRelatedSecurity - Works",
			&types.AttributeValueMemberS{Value: "DueToRelatedSecurity"}, Financial_Quotes_DueToRelatedSecurity),
		Entry("Value is string, InViewOfCommon - Works",
			&types.AttributeValueMemberS{Value: "InViewOfCommon"}, Financial_Quotes_InViewOfCommon),
		Entry("Value is string, EquipmentChangeover - Works",
			&types.AttributeValueMemberS{Value: "EquipmentChangeover"}, Financial_Quotes_EquipmentChangeover),
		Entry("Value is string, NoOpenNoResponse - Works",
			&types.AttributeValueMemberS{Value: "NoOpenNoResponse"}, Financial_Quotes_NoOpenNoResponse),
		Entry("Value is string, SubPennyTrading - Works",
			&types.AttributeValueMemberS{Value: "SubPennyTrading"}, Financial_Quotes_SubPennyTrading),
		Entry("Value is string, AutomatedBidNoOfferNoBid - Works",
			&types.AttributeValueMemberS{Value: "AutomatedBidNoOfferNoBid"}, Financial_Quotes_AutomatedBidNoOfferNoBid),
		Entry("Value is string, LULDPPriceBand - Works",
			&types.AttributeValueMemberS{Value: "LULDPPriceBand"}, Financial_Quotes_LULDPPriceBand),
		Entry("Value is string, MarketWideCircuitBreakerLevel1 - Works",
			&types.AttributeValueMemberS{Value: "MarketWideCircuitBreakerLevel1"}, Financial_Quotes_MarketWideCircuitBreakerLevel1),
		Entry("Value is string, MarketWideCircuitBreakerLevel2 - Works",
			&types.AttributeValueMemberS{Value: "MarketWideCircuitBreakerLevel2"}, Financial_Quotes_MarketWideCircuitBreakerLevel2),
		Entry("Value is string, MarketWideCircuitBreakerLevel3 - Works",
			&types.AttributeValueMemberS{Value: "MarketWideCircuitBreakerLevel3"}, Financial_Quotes_MarketWideCircuitBreakerLevel3),
		Entry("Value is string, RepublishedLULDPriceBand - Works",
			&types.AttributeValueMemberS{Value: "RepublishedLULDPriceBand"}, Financial_Quotes_RepublishedLULDPriceBand),
		Entry("Value is string, OnDemandAuction - Works",
			&types.AttributeValueMemberS{Value: "OnDemandAuction"}, Financial_Quotes_OnDemandAuction),
		Entry("Value is string, CashOnlySettlement - Works",
			&types.AttributeValueMemberS{Value: "CashOnlySettlement"}, Financial_Quotes_CashOnlySettlement),
		Entry("Value is string, NextDaySettlement - Works",
			&types.AttributeValueMemberS{Value: "NextDaySettlement"}, Financial_Quotes_NextDaySettlement),
		Entry("Value is string, LULDTradingPause - Works",
			&types.AttributeValueMemberS{Value: "LULDTradingPause"}, Financial_Quotes_LULDTradingPause),
		Entry("Value is string, SlowDuelRPBidAsk - Works",
			&types.AttributeValueMemberS{Value: "SlowDuelRPBidAsk"}, Financial_Quotes_SlowDuelRPBidAsk),
		Entry("Value is string, Cancel - Works",
			&types.AttributeValueMemberS{Value: "Cancel"}, Financial_Quotes_Cancel),
		Entry("Value is string, CorrectedPrice - Works",
			&types.AttributeValueMemberS{Value: "CorrectedPrice"}, Financial_Quotes_CorrectedPrice),
		Entry("Value is string, SIPGenerated - Works",
			&types.AttributeValueMemberS{Value: "SIPGenerated"}, Financial_Quotes_SIPGenerated),
		Entry("Value is string, Unknown - Works",
			&types.AttributeValueMemberS{Value: "Unknown"}, Financial_Quotes_Unknown),
		Entry("Value is string, CrossedMarket - Works",
			&types.AttributeValueMemberS{Value: "CrossedMarket"}, Financial_Quotes_CrossedMarket),
		Entry("Value is string, LockedMarket - Works",
			&types.AttributeValueMemberS{Value: "LockedMarket"}, Financial_Quotes_LockedMarket),
		Entry("Value is string, DepthOnOfferSide - Works",
			&types.AttributeValueMemberS{Value: "DepthOnOfferSide"}, Financial_Quotes_DepthOnOfferSide),
		Entry("Value is string, DepthOnBidSide - Works",
			&types.AttributeValueMemberS{Value: "DepthOnBidSide"}, Financial_Quotes_DepthOnBidSide),
		Entry("Value is string, DepthOnBidAndOffer - Works",
			&types.AttributeValueMemberS{Value: "DepthOnBidAndOffer"}, Financial_Quotes_DepthOnBidAndOffer),
		Entry("Value is string, PreOpeningIndication - Works",
			&types.AttributeValueMemberS{Value: "PreOpeningIndication"}, Financial_Quotes_PreOpeningIndication),
		Entry("Value is string, SyndicateBid - Works",
			&types.AttributeValueMemberS{Value: "SyndicateBid"}, Financial_Quotes_SyndicateBid),
		Entry("Value is string, PreSyndicateBid - Works",
			&types.AttributeValueMemberS{Value: "PreSyndicateBid"}, Financial_Quotes_PreSyndicateBid),
		Entry("Value is string, PenaltyBid - Works",
			&types.AttributeValueMemberS{Value: "PenaltyBid"}, Financial_Quotes_PenaltyBid),
		Entry("Value is string, CQSGenerated - Works",
			&types.AttributeValueMemberS{Value: "CQSGenerated"}, Financial_Quotes_CQSGenerated),
		Entry("Value is string, Invalid - Works",
			&types.AttributeValueMemberS{Value: "Invalid"}, Financial_Quotes_Invalid))

	// Test that attempting to deserialize a Financial.Quotes.Condition will fial and return an
	// error if the value cannot be converted to either the name value or integer value of the enum option
	It("Scan - Value is nil - Error", func() {

		// Attempt to convert a fake string value into a Financial.Quotes.Condition
		// This should return an error
		var enum *Financial_Quotes_Condition
		err := enum.Scan(nil)

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of %!q(<nil>) had an invalid type of <nil>"))
		Expect(enum).Should(BeNil())
	})

	// Test the conditions under which values should be convertible to a Financial.Quotes.Condition
	DescribeTable("Scan Tests",
		func(value interface{}, shouldBe Financial_Quotes_Condition) {

			// Attempt to convert the value into a Financial.Quotes.Condition
			// This should not fail
			var enum Financial_Quotes_Condition
			err := enum.Scan(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("Regular, Two-Sided Open - Works", "Regular, Two-Sided Open", Financial_Quotes_RegularTwoSidedOpen),
		Entry("Regular, One-Sided Open - Works", "Regular, One-Sided Open", Financial_Quotes_RegularOneSidedOpen),
		Entry("Slow Ask - Works", "Slow Ask", Financial_Quotes_SlowAsk),
		Entry("Slow Bid - Works", "Slow Bid", Financial_Quotes_SlowBid),
		Entry("Slow Bid, Ask - Works", "Slow Bid, Ask", Financial_Quotes_SlowBidAsk),
		Entry("Slow Due, LRP Bid - Works", "Slow Due, LRP Bid", Financial_Quotes_SlowDueLRPBid),
		Entry("Slow Due, LRP Ask - Works", "Slow Due, LRP Ask", Financial_Quotes_SlowDueLRPAsk),
		Entry("Slow Due, NYSE LRP - Works", "Slow Due, NYSE LRP", Financial_Quotes_SlowDueNYSELRP),
		Entry("Slow Due Set, Slow List, Bid, Ask - Works",
			"Slow Due Set, Slow List, Bid, Ask", Financial_Quotes_SlowDueSetSlowListBidAsk),
		Entry("Manual Ask, Automated Bid - Works", "Manual Ask, Automated Bid", Financial_Quotes_ManualAskAutomatedBid),
		Entry("Manual Bid, Automated Ask - Works", "Manual Bid, Automated Ask", Financial_Quotes_ManualBidAutomatedAsk),
		Entry("Manual Bid and Ask - Works", "Manual Bid and Ask", Financial_Quotes_ManualBidAndAsk),
		Entry("Fast Trading - Works", "Fast Trading", Financial_Quotes_FastTrading),
		Entry("Tading Range Indicated - Works", "Tading Range Indicated", Financial_Quotes_TradingRangeIndicated),
		Entry("Market-Maker Quotes Closed - Works", "Market-Maker Quotes Closed", Financial_Quotes_MarketMakerQuotesClosed),
		Entry("Non-Firm - Works", "Non-Firm", Financial_Quotes_NonFirm),
		Entry("News Dissemination - Works", "News Dissemination", Financial_Quotes_NewsDissemination),
		Entry("Order Influx - Works", "Order Influx", Financial_Quotes_OrderInflux),
		Entry("Order Imbalance - Works", "Order Imbalance", Financial_Quotes_OrderImbalance),
		Entry("Due to Related Security, News Dissemination - Works",
			"Due to Related Security, News Dissemination", Financial_Quotes_DueToRelatedSecurityNewsDissemination),
		Entry("Due to Related Security, News Pending - Works",
			"Due to Related Security, News Pending", Financial_Quotes_DueToRelatedSecurityNewsPending),
		Entry("Additional Information - Works", "Additional Information", Financial_Quotes_AdditionalInformation),
		Entry("News Pending - Works", "News Pending", Financial_Quotes_NewsPending),
		Entry("Additional Information Due to Related Security - Works",
			"Additional Information Due to Related Security", Financial_Quotes_AdditionalInformationDueToRelatedSecurity),
		Entry("Due to Related Security - Works", "Due to Related Security", Financial_Quotes_DueToRelatedSecurity),
		Entry("In View of Common - Works", "In View of Common", Financial_Quotes_InViewOfCommon),
		Entry("Equipment Changeover - Works", "Equipment Changeover", Financial_Quotes_EquipmentChangeover),
		Entry("No Open, No Response - Works", "No Open, No Response", Financial_Quotes_NoOpenNoResponse),
		Entry("Sub-Penny Trading - Works", "Sub-Penny Trading", Financial_Quotes_SubPennyTrading),
		Entry("Automated Bid; No Offer, No Bid - Works",
			"Automated Bid; No Offer, No Bid", Financial_Quotes_AutomatedBidNoOfferNoBid),
		Entry("LULD Price Band - Works", "LULD Price Band", Financial_Quotes_LULDPPriceBand),
		Entry("Market-Wide Circuit Breaker, Level 1 - Works",
			"Market-Wide Circuit Breaker, Level 1", Financial_Quotes_MarketWideCircuitBreakerLevel1),
		Entry("Market-Wide Circuit Breaker, Level 2 - Works",
			"Market-Wide Circuit Breaker, Level 2", Financial_Quotes_MarketWideCircuitBreakerLevel2),
		Entry("Market-Wide Circuit Breaker, Level 3 - Works",
			"Market-Wide Circuit Breaker, Level 3", Financial_Quotes_MarketWideCircuitBreakerLevel3),
		Entry("Republished LULD Price Band - Works",
			"Republished LULD Price Band", Financial_Quotes_RepublishedLULDPriceBand),
		Entry("On-Demand Auction - Works", "On-Demand Auction", Financial_Quotes_OnDemandAuction),
		Entry("Cash-Only Settlement - Works", "Cash-Only Settlement", Financial_Quotes_CashOnlySettlement),
		Entry("Next-Day Settlement - Works", "Next-Day Settlement", Financial_Quotes_NextDaySettlement),
		Entry("LULD Trading Pause - Works", "LULD Trading Pause", Financial_Quotes_LULDTradingPause),
		Entry("Slow Due LRP, Bid, Ask - Works", "Slow Due LRP, Bid, Ask", Financial_Quotes_SlowDuelRPBidAsk),
		Entry("Corrected Price - Works", "Corrected Price", Financial_Quotes_CorrectedPrice),
		Entry("SIP-Generated - Works", "SIP-Generated", Financial_Quotes_SIPGenerated),
		Entry("Crossed Market - Works", "Crossed Market", Financial_Quotes_CrossedMarket),
		Entry("Locked Market - Works", "Locked Market", Financial_Quotes_LockedMarket),
		Entry("Depth on Offer Side - Works", "Depth on Offer Side", Financial_Quotes_DepthOnOfferSide),
		Entry("Depth on Bid Side - Works", "Depth on Bid Side", Financial_Quotes_DepthOnBidSide),
		Entry("Depth on Bid and Offer - Works", "Depth on Bid and Offer", Financial_Quotes_DepthOnBidAndOffer),
		Entry("Pre-Opening Indication - Works", "Pre-Opening Indication", Financial_Quotes_PreOpeningIndication),
		Entry("Syndicate Bid - Works", "Syndicate Bid", Financial_Quotes_SyndicateBid),
		Entry("Pre-Syndicate Bid - Works", "Pre-Syndicate Bid", Financial_Quotes_PreSyndicateBid),
		Entry("Penalty Bid - Works", "Penalty Bid", Financial_Quotes_PenaltyBid),
		Entry("CQS-Generated - Works", "CQS-Generated", Financial_Quotes_CQSGenerated),
		Entry("Regular - Works", "Regular", Financial_Quotes_Regular),
		Entry("RegularTwoSidedOpen - Works", "RegularTwoSidedOpen", Financial_Quotes_RegularTwoSidedOpen),
		Entry("RegularOneSidedOpen - Works", "RegularOneSidedOpen", Financial_Quotes_RegularOneSidedOpen),
		Entry("SlowAsk - Works", "SlowAsk", Financial_Quotes_SlowAsk),
		Entry("SlowBid - Works", "SlowBid", Financial_Quotes_SlowBid),
		Entry("SlowBidAsk - Works", "SlowBidAsk", Financial_Quotes_SlowBidAsk),
		Entry("SlowDueLRPBid - Works", "SlowDueLRPBid", Financial_Quotes_SlowDueLRPBid),
		Entry("SlowDueLRPAsk - Works", "SlowDueLRPAsk", Financial_Quotes_SlowDueLRPAsk),
		Entry("SlowDueNYSELRP - Works", "SlowDueNYSELRP", Financial_Quotes_SlowDueNYSELRP),
		Entry("SlowDueSetSlowListBidAsk - Works", "SlowDueSetSlowListBidAsk", Financial_Quotes_SlowDueSetSlowListBidAsk),
		Entry("ManualAskAutomatedBid - Works", "ManualAskAutomatedBid", Financial_Quotes_ManualAskAutomatedBid),
		Entry("ManualBidAutomatedAsk - Works", "ManualBidAutomatedAsk", Financial_Quotes_ManualBidAutomatedAsk),
		Entry("ManualBidAndAsk - Works", "ManualBidAndAsk", Financial_Quotes_ManualBidAndAsk),
		Entry("Opening - Works", "Opening", Financial_Quotes_Opening),
		Entry("Closing - Works", "Closing", Financial_Quotes_Closing),
		Entry("Closed - Works", "Closed", Financial_Quotes_Closed),
		Entry("Resume - Works", "Resume", Financial_Quotes_Resume),
		Entry("FastTrading - Works", "FastTrading", Financial_Quotes_FastTrading),
		Entry("TradingRangeIndicated - Works", "TradingRangeIndicated", Financial_Quotes_TradingRangeIndicated),
		Entry("MarketMakerQuotesClosed - Works", "MarketMakerQuotesClosed", Financial_Quotes_MarketMakerQuotesClosed),
		Entry("NonFirm - Works", "NonFirm", Financial_Quotes_NonFirm),
		Entry("NewsDissemination - Works", "NewsDissemination", Financial_Quotes_NewsDissemination),
		Entry("OrderInflux - Works", "OrderInflux", Financial_Quotes_OrderInflux),
		Entry("OrderImbalance - Works", "OrderImbalance", Financial_Quotes_OrderImbalance),
		Entry("DueToRelatedSecurityNewsDissemination - Works", "DueToRelatedSecurityNewsDissemination", Financial_Quotes_DueToRelatedSecurityNewsDissemination),
		Entry("DueToRelatedSecurityNewsPending - Works", "DueToRelatedSecurityNewsPending", Financial_Quotes_DueToRelatedSecurityNewsPending),
		Entry("AdditionalInformation - Works", "AdditionalInformation", Financial_Quotes_AdditionalInformation),
		Entry("NewsPending - Works", "NewsPending", Financial_Quotes_NewsPending),
		Entry("AdditionalInformationDueToRelatedSecurity - Works", "AdditionalInformationDueToRelatedSecurity", Financial_Quotes_AdditionalInformationDueToRelatedSecurity),
		Entry("DueToRelatedSecurity - Works", "DueToRelatedSecurity", Financial_Quotes_DueToRelatedSecurity),
		Entry("InViewOfCommon - Works", "InViewOfCommon", Financial_Quotes_InViewOfCommon),
		Entry("EquipmentChangeover - Works", "EquipmentChangeover", Financial_Quotes_EquipmentChangeover),
		Entry("NoOpenNoResponse - Works", "NoOpenNoResponse", Financial_Quotes_NoOpenNoResponse),
		Entry("SubPennyTrading - Works", "SubPennyTrading", Financial_Quotes_SubPennyTrading),
		Entry("AutomatedBidNoOfferNoBid - Works", "AutomatedBidNoOfferNoBid", Financial_Quotes_AutomatedBidNoOfferNoBid),
		Entry("LULDPPriceBand - Works", "LULDPPriceBand", Financial_Quotes_LULDPPriceBand),
		Entry("MarketWideCircuitBreakerLevel1 - Works", "MarketWideCircuitBreakerLevel1", Financial_Quotes_MarketWideCircuitBreakerLevel1),
		Entry("MarketWideCircuitBreakerLevel2 - Works", "MarketWideCircuitBreakerLevel2", Financial_Quotes_MarketWideCircuitBreakerLevel2),
		Entry("MarketWideCircuitBreakerLevel3 - Works", "MarketWideCircuitBreakerLevel3", Financial_Quotes_MarketWideCircuitBreakerLevel3),
		Entry("RepublishedLULDPriceBand - Works", "RepublishedLULDPriceBand", Financial_Quotes_RepublishedLULDPriceBand),
		Entry("OnDemandAuction - Works", "OnDemandAuction", Financial_Quotes_OnDemandAuction),
		Entry("CashOnlySettlement - Works", "CashOnlySettlement", Financial_Quotes_CashOnlySettlement),
		Entry("NextDaySettlement - Works", "NextDaySettlement", Financial_Quotes_NextDaySettlement),
		Entry("LULDTradingPause - Works", "LULDTradingPause", Financial_Quotes_LULDTradingPause),
		Entry("SlowDuelRPBidAsk - Works", "SlowDuelRPBidAsk", Financial_Quotes_SlowDuelRPBidAsk),
		Entry("Cancel - Works", "Cancel", Financial_Quotes_Cancel),
		Entry("CorrectedPrice - Works", "CorrectedPrice", Financial_Quotes_CorrectedPrice),
		Entry("SIPGenerated - Works", "SIPGenerated", Financial_Quotes_SIPGenerated),
		Entry("Unknown - Works", "Unknown", Financial_Quotes_Unknown),
		Entry("CrossedMarket - Works", "CrossedMarket", Financial_Quotes_CrossedMarket),
		Entry("LockedMarket - Works", "LockedMarket", Financial_Quotes_LockedMarket),
		Entry("DepthOnOfferSide - Works", "DepthOnOfferSide", Financial_Quotes_DepthOnOfferSide),
		Entry("DepthOnBidSide - Works", "DepthOnBidSide", Financial_Quotes_DepthOnBidSide),
		Entry("DepthOnBidAndOffer - Works", "DepthOnBidAndOffer", Financial_Quotes_DepthOnBidAndOffer),
		Entry("PreOpeningIndication - Works", "PreOpeningIndication", Financial_Quotes_PreOpeningIndication),
		Entry("SyndicateBid - Works", "SyndicateBid", Financial_Quotes_SyndicateBid),
		Entry("PreSyndicateBid - Works", "PreSyndicateBid", Financial_Quotes_PreSyndicateBid),
		Entry("PenaltyBid - Works", "PenaltyBid", Financial_Quotes_PenaltyBid),
		Entry("CQSGenerated - Works", "CQSGenerated", Financial_Quotes_CQSGenerated),
		Entry("Invalid - Works", "Invalid", Financial_Quotes_Invalid),
		Entry("0 - Works", 0, Financial_Quotes_Regular),
		Entry("1 - Works", 1, Financial_Quotes_RegularTwoSidedOpen),
		Entry("2 - Works", 2, Financial_Quotes_RegularOneSidedOpen),
		Entry("3 - Works", 3, Financial_Quotes_SlowAsk),
		Entry("4 - Works", 4, Financial_Quotes_SlowBid),
		Entry("5 - Works", 5, Financial_Quotes_SlowBidAsk),
		Entry("6 - Works", 6, Financial_Quotes_SlowDueLRPBid),
		Entry("7 - Works", 7, Financial_Quotes_SlowDueLRPAsk),
		Entry("8 - Works", 8, Financial_Quotes_SlowDueNYSELRP),
		Entry("9 - Works", 9, Financial_Quotes_SlowDueSetSlowListBidAsk),
		Entry("10 - Works", 10, Financial_Quotes_ManualAskAutomatedBid),
		Entry("11 - Works", 11, Financial_Quotes_ManualBidAutomatedAsk),
		Entry("12 - Works", 12, Financial_Quotes_ManualBidAndAsk),
		Entry("13 - Works", 13, Financial_Quotes_Opening),
		Entry("14 - Works", 14, Financial_Quotes_Closing),
		Entry("15 - Works", 15, Financial_Quotes_Closed),
		Entry("16 - Works", 16, Financial_Quotes_Resume),
		Entry("17 - Works", 17, Financial_Quotes_FastTrading),
		Entry("18 - Works", 18, Financial_Quotes_TradingRangeIndicated),
		Entry("19 - Works", 19, Financial_Quotes_MarketMakerQuotesClosed),
		Entry("20 - Works", 20, Financial_Quotes_NonFirm),
		Entry("21 - Works", 21, Financial_Quotes_NewsDissemination),
		Entry("22 - Works", 22, Financial_Quotes_OrderInflux),
		Entry("23 - Works", 23, Financial_Quotes_OrderImbalance),
		Entry("24 - Works", 24, Financial_Quotes_DueToRelatedSecurityNewsDissemination),
		Entry("25 - Works", 25, Financial_Quotes_DueToRelatedSecurityNewsPending),
		Entry("26 - Works", 26, Financial_Quotes_AdditionalInformation),
		Entry("27 - Works", 27, Financial_Quotes_NewsPending),
		Entry("28 - Works", 28, Financial_Quotes_AdditionalInformationDueToRelatedSecurity),
		Entry("29 - Works", 29, Financial_Quotes_DueToRelatedSecurity),
		Entry("30 - Works", 30, Financial_Quotes_InViewOfCommon),
		Entry("31 - Works", 31, Financial_Quotes_EquipmentChangeover),
		Entry("32 - Works", 32, Financial_Quotes_NoOpenNoResponse),
		Entry("33 - Works", 33, Financial_Quotes_SubPennyTrading),
		Entry("34 - Works", 34, Financial_Quotes_AutomatedBidNoOfferNoBid),
		Entry("35 - Works", 35, Financial_Quotes_LULDPPriceBand),
		Entry("36 - Works", 36, Financial_Quotes_MarketWideCircuitBreakerLevel1),
		Entry("37 - Works", 37, Financial_Quotes_MarketWideCircuitBreakerLevel2),
		Entry("38 - Works", 38, Financial_Quotes_MarketWideCircuitBreakerLevel3),
		Entry("39 - Works", 39, Financial_Quotes_RepublishedLULDPriceBand),
		Entry("40 - Works", 40, Financial_Quotes_OnDemandAuction),
		Entry("41 - Works", 41, Financial_Quotes_CashOnlySettlement),
		Entry("42 - Works", 42, Financial_Quotes_NextDaySettlement),
		Entry("43 - Works", 43, Financial_Quotes_LULDTradingPause),
		Entry("71 - Works", 71, Financial_Quotes_SlowDuelRPBidAsk),
		Entry("80 - Works", 80, Financial_Quotes_Cancel),
		Entry("81 - Works", 81, Financial_Quotes_CorrectedPrice),
		Entry("82 - Works", 82, Financial_Quotes_SIPGenerated),
		Entry("83 - Works", 83, Financial_Quotes_Unknown),
		Entry("84 - Works", 84, Financial_Quotes_CrossedMarket),
		Entry("85 - Works", 85, Financial_Quotes_LockedMarket),
		Entry("86 - Works", 86, Financial_Quotes_DepthOnOfferSide),
		Entry("87 - Works", 87, Financial_Quotes_DepthOnBidSide),
		Entry("88 - Works", 88, Financial_Quotes_DepthOnBidAndOffer),
		Entry("89 - Works", 89, Financial_Quotes_PreOpeningIndication),
		Entry("90 - Works", 90, Financial_Quotes_SyndicateBid),
		Entry("91 - Works", 91, Financial_Quotes_PreSyndicateBid),
		Entry("92 - Works", 92, Financial_Quotes_PenaltyBid),
		Entry("94 - Works", 94, Financial_Quotes_CQSGenerated),
		Entry("999 - Works", 999, Financial_Quotes_Invalid))
})

var _ = Describe("Financial.Quotes.Indicator Marshal/Unmarshal Tests", func() {

	// Test that converting the Financial.Quotes.Indicator enum to JSON works for all values
	DescribeTable("MarshalJSON Tests",
		func(enum Financial_Quotes_Indicator, value string) {
			data, err := json.Marshal(enum)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("NBBNBOExecutable - Works", Financial_Quotes_NBBNBOExecutable, "\"NBBNBOExecutable\""),
		Entry("NBBBelowLowerBand - Works", Financial_Quotes_NBBBelowLowerBand, "\"NBBBelowLowerBand\""),
		Entry("NBOAboveUpperBand - Works", Financial_Quotes_NBOAboveUpperBand, "\"NBOAboveUpperBand\""),
		Entry("NBBBelowLowerBandAndNBOAboveUpperBand - Works",
			Financial_Quotes_NBBBelowLowerBandAndNBOAboveUpperBand, "\"NBBBelowLowerBandAndNBOAboveUpperBand\""),
		Entry("NBBEqualsUpperBand - Works", Financial_Quotes_NBBEqualsUpperBand, "\"NBBEqualsUpperBand\""),
		Entry("NBOEqualsLowerBand - Works", Financial_Quotes_NBOEqualsLowerBand, "\"NBOEqualsLowerBand\""),
		Entry("NBBEqualsUpperBandAndNBOAboveUpperBand - Works",
			Financial_Quotes_NBBEqualsUpperBandAndNBOAboveUpperBand, "\"NBBEqualsUpperBandAndNBOAboveUpperBand\""),
		Entry("NBBBelowLowerBandAndNBOEqualsLowerBand - Works",
			Financial_Quotes_NBBBelowLowerBandAndNBOEqualsLowerBand, "\"NBBBelowLowerBandAndNBOEqualsLowerBand\""),
		Entry("BidPriceAboveUpperLimitPriceBand - Works",
			Financial_Quotes_BidPriceAboveUpperLimitPriceBand, "\"BidPriceAboveUpperLimitPriceBand\""),
		Entry("OfferPriceBelowLowerLimitPriceBand - Works",
			Financial_Quotes_OfferPriceBelowLowerLimitPriceBand, "\"OfferPriceBelowLowerLimitPriceBand\""),
		Entry("BidAndOfferOutsidePriceBand - Works",
			Financial_Quotes_BidAndOfferOutsidePriceBand, "\"BidAndOfferOutsidePriceBand\""),
		Entry("OpeningUpdate - Works", Financial_Quotes_OpeningUpdate, "\"OpeningUpdate\""),
		Entry("IntraDayUpdate - Works", Financial_Quotes_IntraDayUpdate, "\"IntraDayUpdate\""),
		Entry("RestatedValue - Works", Financial_Quotes_RestatedValue, "\"RestatedValue\""),
		Entry("SuspendedDuringTradingHalt - Works",
			Financial_Quotes_SuspendedDuringTradingHalt, "\"SuspendedDuringTradingHalt\""),
		Entry("ReOpeningUpdate - Works", Financial_Quotes_ReOpeningUpdate, "\"ReOpeningUpdate\""),
		Entry("OutsidePriceBandRuleHours - Works", Financial_Quotes_OutsidePriceBandRuleHours, "\"OutsidePriceBandRuleHours\""),
		Entry("AuctionExtension - Works", Financial_Quotes_AuctionExtension, "\"AuctionExtension\""),
		Entry("LULDPriceBand - Works", Financial_Quotes_LULDPriceBand, "\"LULDPriceBand\""),
		Entry("RepublishedLULDPriceBandInd - Works",
			Financial_Quotes_RepublishedLULDPriceBandInd, "\"RepublishedLULDPriceBandInd\""),
		Entry("NBBLimitStateEntered - Works", Financial_Quotes_NBBLimitStateEntered, "\"NBBLimitStateEntered\""),
		Entry("NBBLimitStateExited - Works", Financial_Quotes_NBBLimitStateExited, "\"NBBLimitStateExited\""),
		Entry("NBOLimitStateEntered - Works", Financial_Quotes_NBOLimitStateEntered, "\"NBOLimitStateEntered\""),
		Entry("NBOLimitStateExited - Works", Financial_Quotes_NBOLimitStateExited, "\"NBOLimitStateExited\""),
		Entry("NBBAndNBOLimitStateEntered - Works",
			Financial_Quotes_NBBAndNBOLimitStateEntered, "\"NBBAndNBOLimitStateEntered\""),
		Entry("NBBAndNBOLimitStateExited - Works", Financial_Quotes_NBBAndNBOLimitStateExited, "\"NBBAndNBOLimitStateExited\""),
		Entry("NBBLimitStateEnteredNBOLimitStateExited - Works",
			Financial_Quotes_NBBLimitStateEnteredNBOLimitStateExited, "\"NBBLimitStateEnteredNBOLimitStateExited\""),
		Entry("NBBLimitStateExitedNBOLimitStateEntered - Works",
			Financial_Quotes_NBBLimitStateExitedNBOLimitStateEntered, "\"NBBLimitStateExitedNBOLimitStateEntered\""),
		Entry("Normal - Works", Financial_Quotes_Normal, "\"Normal\""),
		Entry("Bankrupt - Works", Financial_Quotes_Bankrupt, "\"Bankrupt\""),
		Entry("Deficient - Works", Financial_Quotes_Deficient, "\"Deficient\""),
		Entry("Delinquent - Works", Financial_Quotes_Delinquent, "\"Delinquent\""),
		Entry("BankruptAndDeficient - Works", Financial_Quotes_BankruptAndDeficient, "\"BankruptAndDeficient\""),
		Entry("BankruptAndDelinquent - Works", Financial_Quotes_BankruptAndDelinquent, "\"BankruptAndDelinquent\""),
		Entry("DeficientAndDelinquent - Works", Financial_Quotes_DeficientAndDelinquent, "\"DeficientAndDelinquent\""),
		Entry("DeficientDeliquentBankrupt - Works",
			Financial_Quotes_DeficientDeliquentBankrupt, "\"DeficientDeliquentBankrupt\""),
		Entry("Liquidation - Works", Financial_Quotes_Liquidation, "\"Liquidation\""),
		Entry("CreationsSuspended - Works", Financial_Quotes_CreationsSuspended, "\"CreationsSuspended\""),
		Entry("RedemptionsSuspended - Works", Financial_Quotes_RedemptionsSuspended, "\"RedemptionsSuspended\""),
		Entry("CreationsRedemptionsSuspended - Works",
			Financial_Quotes_CreationsRedemptionsSuspended, "\"CreationsRedemptionsSuspended\""),
		Entry("NormalTrading - Works", Financial_Quotes_NormalTrading, "\"NormalTrading\""),
		Entry("OpeningDelay - Works", Financial_Quotes_OpeningDelay, "\"OpeningDelay\""),
		Entry("TradingHalt - Works", Financial_Quotes_TradingHalt, "\"TradingHalt\""),
		Entry("TradingResume - Works", Financial_Quotes_TradingResume, "\"TradingResume\""),
		Entry("NoOpenNoResume - Works", Financial_Quotes_NoOpenNoResume, "\"NoOpenNoResume\""),
		Entry("PriceIndication - Works", Financial_Quotes_PriceIndication, "\"PriceIndication\""),
		Entry("TradingRangeIndication - Works", Financial_Quotes_TradingRangeIndication, "\"TradingRangeIndication\""),
		Entry("MarketImbalanceBuy - Works", Financial_Quotes_MarketImbalanceBuy, "\"MarketImbalanceBuy\""),
		Entry("MarketImbalanceSell - Works", Financial_Quotes_MarketImbalanceSell, "\"MarketImbalanceSell\""),
		Entry("MarketOnCloseImbalanceBuy - Works", Financial_Quotes_MarketOnCloseImbalanceBuy, "\"MarketOnCloseImbalanceBuy\""),
		Entry("MarketOnCloseImbalanceSell - Works",
			Financial_Quotes_MarketOnCloseImbalanceSell, "\"MarketOnCloseImbalanceSell\""),
		Entry("NoMarketImbalance - Works", Financial_Quotes_NoMarketImbalance, "\"NoMarketImbalance\""),
		Entry("NoMarketOnCloseImbalance - Works", Financial_Quotes_NoMarketOnCloseImbalance, "\"NoMarketOnCloseImbalance\""),
		Entry("ShortSaleRestriction - Works", Financial_Quotes_ShortSaleRestriction, "\"ShortSaleRestriction\""),
		Entry("LimitUpLimitDown - Works", Financial_Quotes_LimitUpLimitDown, "\"LimitUpLimitDown\""),
		Entry("QuotationResumption - Works", Financial_Quotes_QuotationResumption, "\"QuotationResumption\""),
		Entry("TradingResumption - Works", Financial_Quotes_TradingResumption, "\"TradingResumption\""),
		Entry("VolatilityTradingPause - Works", Financial_Quotes_VolatilityTradingPause, "\"VolatilityTradingPause\""),
		Entry("PolygonReserved - Works", Financial_Quotes_PolygonReserved, "\"PolygonReserved\""),
		Entry("HaltNewsPending - Works", Financial_Quotes_HaltNewsPending, "\"HaltNewsPending\""),
		Entry("UpdateNewsDissemination - Works", Financial_Quotes_UpdateNewsDissemination, "\"UpdateNewsDissemination\""),
		Entry("HaltSingleStockTradingPause - Works",
			Financial_Quotes_HaltSingleStockTradingPause, "\"HaltSingleStockTradingPause\""),
		Entry("HaltRegulatoryExtraordinaryMarketActivity - Works",
			Financial_Quotes_HaltRegulatoryExtraordinaryMarketActivity, "\"HaltRegulatoryExtraordinaryMarketActivity\""),
		Entry("HaltETF - Works", Financial_Quotes_HaltETF, "\"HaltETF\""),
		Entry("HaltInformationRequested - Works", Financial_Quotes_HaltInformationRequested, "\"HaltInformationRequested\""),
		Entry("HaltExchangeNonCompliance - Works", Financial_Quotes_HaltExchangeNonCompliance, "\"HaltExchangeNonCompliance\""),
		Entry("HaltFilingsNotCurrent - Works", Financial_Quotes_HaltFilingsNotCurrent, "\"HaltFilingsNotCurrent\""),
		Entry("HaltSECTradingSuspension - Works", Financial_Quotes_HaltSECTradingSuspension, "\"HaltSECTradingSuspension\""),
		Entry("HaltRegulatoryConcern - Works", Financial_Quotes_HaltRegulatoryConcern, "\"HaltRegulatoryConcern\""),
		Entry("HaltMarketOperations - Works", Financial_Quotes_HaltMarketOperations, "\"HaltMarketOperations\""),
		Entry("IPOSecurityNotYetTrading - Works", Financial_Quotes_IPOSecurityNotYetTrading, "\"IPOSecurityNotYetTrading\""),
		Entry("HaltCorporateAction - Works", Financial_Quotes_HaltCorporateAction, "\"HaltCorporateAction\""),
		Entry("QuotationNotAvailable - Works", Financial_Quotes_QuotationNotAvailable, "\"QuotationNotAvailable\""),
		Entry("HaltVolatilityTradingPause - Works",
			Financial_Quotes_HaltVolatilityTradingPause, "\"HaltVolatilityTradingPause\""),
		Entry("HaltVolatilityTradingPauseStraddleCondition - Works",
			Financial_Quotes_HaltVolatilityTradingPauseStraddleCondition, "\"HaltVolatilityTradingPauseStraddleCondition\""),
		Entry("UpdateNewsAndResumptionTimes - Works",
			Financial_Quotes_UpdateNewsAndResumptionTimes, "\"UpdateNewsAndResumptionTimes\""),
		Entry("HaltSingleStockTradingPauseQuotesOnly - Works",
			Financial_Quotes_HaltSingleStockTradingPauseQuotesOnly, "\"HaltSingleStockTradingPauseQuotesOnly\""),
		Entry("ResumeQualificationIssuesReviewedResolved - Works",
			Financial_Quotes_ResumeQualificationIssuesReviewedResolved, "\"ResumeQualificationIssuesReviewedResolved\""),
		Entry("ResumeFilingRequirementsSatisfiedResolved - Works",
			Financial_Quotes_ResumeFilingRequirementsSatisfiedResolved, "\"ResumeFilingRequirementsSatisfiedResolved\""),
		Entry("ResumeNewsNotForthcoming - Works", Financial_Quotes_ResumeNewsNotForthcoming, "\"ResumeNewsNotForthcoming\""),
		Entry("ResumeQualificationsMaintRequirementsMet - Works",
			Financial_Quotes_ResumeQualificationsMaintRequirementsMet, "\"ResumeQualificationsMaintRequirementsMet\""),
		Entry("ResumeQualificationsFilingsMet - Works",
			Financial_Quotes_ResumeQualificationsFilingsMet, "\"ResumeQualificationsFilingsMet\""),
		Entry("ResumeRegulatoryAuth - Works", Financial_Quotes_ResumeRegulatoryAuth, "\"ResumeRegulatoryAuth\""),
		Entry("NewIssueAvailable - Works", Financial_Quotes_NewIssueAvailable, "\"NewIssueAvailable\""),
		Entry("IssueAvailable - Works", Financial_Quotes_IssueAvailable, "\"IssueAvailable\""),
		Entry("MWCBCarryFromPreviousDay - Works", Financial_Quotes_MWCBCarryFromPreviousDay, "\"MWCBCarryFromPreviousDay\""),
		Entry("MWCBResume - Works", Financial_Quotes_MWCBResume, "\"MWCBResume\""),
		Entry("IPOSecurityReleasedForQuotation - Works",
			Financial_Quotes_IPOSecurityReleasedForQuotation, "\"IPOSecurityReleasedForQuotation\""),
		Entry("IPOSecurityPositioningWindowExtension - Works",
			Financial_Quotes_IPOSecurityPositioningWindowExtension, "\"IPOSecurityPositioningWindowExtension\""),
		Entry("MWCBLevel1 - Works", Financial_Quotes_MWCBLevel1, "\"MWCBLevel1\""),
		Entry("MWCBLevel2 - Works", Financial_Quotes_MWCBLevel2, "\"MWCBLevel2\""),
		Entry("MWCBLevel3 - Works", Financial_Quotes_MWCBLevel3, "\"MWCBLevel3\""),
		Entry("HaltSubPennyTrading - Works", Financial_Quotes_HaltSubPennyTrading, "\"HaltSubPennyTrading\""),
		Entry("OrderImbalanceInd - Works", Financial_Quotes_OrderImbalanceInd, "\"OrderImbalanceInd\""),
		Entry("LULDTradingPaused - Works", Financial_Quotes_LULDTradingPaused, "\"LULDTradingPaused\""),
		Entry("NONE - Works", Financial_Quotes_NONE, "\"NONE\""),
		Entry("ShortSalesRestrictionActivated - Works",
			Financial_Quotes_ShortSalesRestrictionActivated, "\"ShortSalesRestrictionActivated\""),
		Entry("ShortSalesRestrictionContinued - Works",
			Financial_Quotes_ShortSalesRestrictionContinued, "\"ShortSalesRestrictionContinued\""),
		Entry("ShortSalesRestrictionDeactivated - Works",
			Financial_Quotes_ShortSalesRestrictionDeactivated, "\"ShortSalesRestrictionDeactivated\""),
		Entry("ShortSalesRestrictionInEffect - Works",
			Financial_Quotes_ShortSalesRestrictionInEffect, "\"ShortSalesRestrictionInEffect\""),
		Entry("ShortSalesRestrictionMax - Works", Financial_Quotes_ShortSalesRestrictionMax, "\"ShortSalesRestrictionMax\""),
		Entry("RetailInterestOnBid - Works", Financial_Quotes_RetailInterestOnBid, "\"RetailInterestOnBid\""),
		Entry("RetailInterestOnAsk - Works", Financial_Quotes_RetailInterestOnAsk, "\"RetailInterestOnAsk\""),
		Entry("RetailInterestOnBidAndAsk - Works", Financial_Quotes_RetailInterestOnBidAndAsk, "\"RetailInterestOnBidAndAsk\""),
		Entry("FinraBBONoChange - Works", Financial_Quotes_FinraBBONoChange, "\"FinraBBONoChange\""),
		Entry("FinraBBODoesNotExist - Works", Financial_Quotes_FinraBBODoesNotExist, "\"FinraBBODoesNotExist\""),
		Entry("FinraBBBOExecutable - Works", Financial_Quotes_FinraBBBOExecutable, "\"FinraBBBOExecutable\""),
		Entry("FinraBBBelowLowerBand - Works", Financial_Quotes_FinraBBBelowLowerBand, "\"FinraBBBelowLowerBand\""),
		Entry("FinraBOAboveUpperBand - Works", Financial_Quotes_FinraBOAboveUpperBand, "\"FinraBOAboveUpperBand\""),
		Entry("FinraBBBelowLowerBandBOAbboveUpperBand - Works",
			Financial_Quotes_FinraBBBelowLowerBandBOAbboveUpperBand, "\"FinraBBBelowLowerBandBOAbboveUpperBand\""),
		Entry("NBBONoChange - Works", Financial_Quotes_NBBONoChange, "\"NBBONoChange\""),
		Entry("NBBOQuoteIsNBBO - Works", Financial_Quotes_NBBOQuoteIsNBBO, "\"NBBOQuoteIsNBBO\""),
		Entry("NBBONoBBNoBO - Works", Financial_Quotes_NBBONoBBNoBO, "\"NBBONoBBNoBO\""),
		Entry("NBBOBBBOShortAppendage - Works", Financial_Quotes_NBBOBBBOShortAppendage, "\"NBBOBBBOShortAppendage\""),
		Entry("NBBOBBBOLongAppendage - Works", Financial_Quotes_NBBOBBBOLongAppendage, "\"NBBOBBBOLongAppendage\""),
		Entry("HeldTradeNotLastSaleNotConsolidated - Works",
			Financial_Quotes_HeldTradeNotLastSaleNotConsolidated, "\"HeldTradeNotLastSaleNotConsolidated\""),
		Entry("HeldTradeLastSaleButNotConsolidated - Works",
			Financial_Quotes_HeldTradeLastSaleButNotConsolidated, "\"HeldTradeLastSaleButNotConsolidated\""),
		Entry("HeldTradeLastSaleAndConsolidated - Works",
			Financial_Quotes_HeldTradeLastSaleAndConsolidated, "\"HeldTradeLastSaleAndConsolidated\""),
		Entry("CTANotDueToRelatedSecurity - Works",
			Financial_Quotes_CTANotDueToRelatedSecurity, "\"CTANotDueToRelatedSecurity\""),
		Entry("CTADueToRelatedSecurity - Works", Financial_Quotes_CTADueToRelatedSecurity, "\"CTADueToRelatedSecurity\""),
		Entry("CTANotInViewOfCommon - Works", Financial_Quotes_CTANotInViewOfCommon, "\"CTANotInViewOfCommon\""),
		Entry("CTAInViewOfCommon - Works", Financial_Quotes_CTAInViewOfCommon, "\"CTAInViewOfCommon\""),
		Entry("CTAPriceIndicator - Works", Financial_Quotes_CTAPriceIndicator, "\"CTAPriceIndicator\""),
		Entry("CTANewPriceIndicator - Works", Financial_Quotes_CTANewPriceIndicator, "\"CTANewPriceIndicator\""),
		Entry("CTACorrectedPriceIndication - Works",
			Financial_Quotes_CTACorrectedPriceIndication, "\"CTACorrectedPriceIndication\""),
		Entry("CTACancelledMarketImbalance - Works",
			Financial_Quotes_CTACancelledMarketImbalance, "\"CTACancelledMarketImbalance\""))

	// Test that converting the Financial.Quotes.Indicator enum to a CSV column works for all values
	DescribeTable("MarshalCSV Tests",
		func(enum Financial_Quotes_Indicator, value string) {
			data, err := enum.MarshalCSV()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("NBBNBOExecutable - Works", Financial_Quotes_NBBNBOExecutable, "0"),
		Entry("NBBBelowLowerBand - Works", Financial_Quotes_NBBBelowLowerBand, "1"),
		Entry("NBOAboveUpperBand - Works", Financial_Quotes_NBOAboveUpperBand, "2"),
		Entry("NBBBelowLowerBandAndNBOAboveUpperBand - Works", Financial_Quotes_NBBBelowLowerBandAndNBOAboveUpperBand, "3"),
		Entry("NBBEqualsUpperBand - Works", Financial_Quotes_NBBEqualsUpperBand, "4"),
		Entry("NBOEqualsLowerBand - Works", Financial_Quotes_NBOEqualsLowerBand, "5"),
		Entry("NBBEqualsUpperBandAndNBOAboveUpperBand - Works", Financial_Quotes_NBBEqualsUpperBandAndNBOAboveUpperBand, "6"),
		Entry("NBBBelowLowerBandAndNBOEqualsLowerBand - Works", Financial_Quotes_NBBBelowLowerBandAndNBOEqualsLowerBand, "7"),
		Entry("BidPriceAboveUpperLimitPriceBand - Works", Financial_Quotes_BidPriceAboveUpperLimitPriceBand, "8"),
		Entry("OfferPriceBelowLowerLimitPriceBand - Works", Financial_Quotes_OfferPriceBelowLowerLimitPriceBand, "9"),
		Entry("BidAndOfferOutsidePriceBand - Works", Financial_Quotes_BidAndOfferOutsidePriceBand, "10"),
		Entry("OpeningUpdate - Works", Financial_Quotes_OpeningUpdate, "11"),
		Entry("IntraDayUpdate - Works", Financial_Quotes_IntraDayUpdate, "12"),
		Entry("RestatedValue - Works", Financial_Quotes_RestatedValue, "13"),
		Entry("SuspendedDuringTradingHalt - Works", Financial_Quotes_SuspendedDuringTradingHalt, "14"),
		Entry("ReOpeningUpdate - Works", Financial_Quotes_ReOpeningUpdate, "15"),
		Entry("OutsidePriceBandRuleHours - Works", Financial_Quotes_OutsidePriceBandRuleHours, "16"),
		Entry("AuctionExtension - Works", Financial_Quotes_AuctionExtension, "17"),
		Entry("LULDPriceBand - Works", Financial_Quotes_LULDPriceBand, "18"),
		Entry("RepublishedLULDPriceBandInd - Works", Financial_Quotes_RepublishedLULDPriceBandInd, "19"),
		Entry("NBBLimitStateEntered - Works", Financial_Quotes_NBBLimitStateEntered, "20"),
		Entry("NBBLimitStateExited - Works", Financial_Quotes_NBBLimitStateExited, "21"),
		Entry("NBOLimitStateEntered - Works", Financial_Quotes_NBOLimitStateEntered, "22"),
		Entry("NBOLimitStateExited - Works", Financial_Quotes_NBOLimitStateExited, "23"),
		Entry("NBBAndNBOLimitStateEntered - Works", Financial_Quotes_NBBAndNBOLimitStateEntered, "24"),
		Entry("NBBAndNBOLimitStateExited - Works", Financial_Quotes_NBBAndNBOLimitStateExited, "25"),
		Entry("NBBLimitStateEnteredNBOLimitStateExited - Works", Financial_Quotes_NBBLimitStateEnteredNBOLimitStateExited, "26"),
		Entry("NBBLimitStateExitedNBOLimitStateEntered - Works", Financial_Quotes_NBBLimitStateExitedNBOLimitStateEntered, "27"),
		Entry("Normal - Works", Financial_Quotes_Normal, "28"),
		Entry("Bankrupt - Works", Financial_Quotes_Bankrupt, "29"),
		Entry("Deficient - Works", Financial_Quotes_Deficient, "30"),
		Entry("Delinquent - Works", Financial_Quotes_Delinquent, "31"),
		Entry("BankruptAndDeficient - Works", Financial_Quotes_BankruptAndDeficient, "32"),
		Entry("BankruptAndDelinquent - Works", Financial_Quotes_BankruptAndDelinquent, "33"),
		Entry("DeficientAndDelinquent - Works", Financial_Quotes_DeficientAndDelinquent, "34"),
		Entry("DeficientDeliquentBankrupt - Works", Financial_Quotes_DeficientDeliquentBankrupt, "35"),
		Entry("Liquidation - Works", Financial_Quotes_Liquidation, "36"),
		Entry("CreationsSuspended - Works", Financial_Quotes_CreationsSuspended, "37"),
		Entry("RedemptionsSuspended - Works", Financial_Quotes_RedemptionsSuspended, "38"),
		Entry("CreationsRedemptionsSuspended - Works", Financial_Quotes_CreationsRedemptionsSuspended, "39"),
		Entry("NormalTrading - Works", Financial_Quotes_NormalTrading, "40"),
		Entry("OpeningDelay - Works", Financial_Quotes_OpeningDelay, "41"),
		Entry("TradingHalt - Works", Financial_Quotes_TradingHalt, "42"),
		Entry("TradingResume - Works", Financial_Quotes_TradingResume, "43"),
		Entry("NoOpenNoResume - Works", Financial_Quotes_NoOpenNoResume, "44"),
		Entry("PriceIndication - Works", Financial_Quotes_PriceIndication, "45"),
		Entry("TradingRangeIndication - Works", Financial_Quotes_TradingRangeIndication, "46"),
		Entry("MarketImbalanceBuy - Works", Financial_Quotes_MarketImbalanceBuy, "47"),
		Entry("MarketImbalanceSell - Works", Financial_Quotes_MarketImbalanceSell, "48"),
		Entry("MarketOnCloseImbalanceBuy - Works", Financial_Quotes_MarketOnCloseImbalanceBuy, "49"),
		Entry("MarketOnCloseImbalanceSell - Works", Financial_Quotes_MarketOnCloseImbalanceSell, "50"),
		Entry("NoMarketImbalance - Works", Financial_Quotes_NoMarketImbalance, "51"),
		Entry("NoMarketOnCloseImbalance - Works", Financial_Quotes_NoMarketOnCloseImbalance, "52"),
		Entry("ShortSaleRestriction - Works", Financial_Quotes_ShortSaleRestriction, "53"),
		Entry("LimitUpLimitDown - Works", Financial_Quotes_LimitUpLimitDown, "54"),
		Entry("QuotationResumption - Works", Financial_Quotes_QuotationResumption, "55"),
		Entry("TradingResumption - Works", Financial_Quotes_TradingResumption, "56"),
		Entry("VolatilityTradingPause - Works", Financial_Quotes_VolatilityTradingPause, "57"),
		Entry("PolygonReserved - Works", Financial_Quotes_PolygonReserved, "58"),
		Entry("HaltNewsPending - Works", Financial_Quotes_HaltNewsPending, "59"),
		Entry("UpdateNewsDissemination - Works", Financial_Quotes_UpdateNewsDissemination, "60"),
		Entry("HaltSingleStockTradingPause - Works", Financial_Quotes_HaltSingleStockTradingPause, "61"),
		Entry("HaltRegulatoryExtraordinaryMarketActivity - Works", Financial_Quotes_HaltRegulatoryExtraordinaryMarketActivity, "62"),
		Entry("HaltETF - Works", Financial_Quotes_HaltETF, "63"),
		Entry("HaltInformationRequested - Works", Financial_Quotes_HaltInformationRequested, "64"),
		Entry("HaltExchangeNonCompliance - Works", Financial_Quotes_HaltExchangeNonCompliance, "65"),
		Entry("HaltFilingsNotCurrent - Works", Financial_Quotes_HaltFilingsNotCurrent, "66"),
		Entry("HaltSECTradingSuspension - Works", Financial_Quotes_HaltSECTradingSuspension, "67"),
		Entry("HaltRegulatoryConcern - Works", Financial_Quotes_HaltRegulatoryConcern, "68"),
		Entry("HaltMarketOperations - Works", Financial_Quotes_HaltMarketOperations, "69"),
		Entry("IPOSecurityNotYetTrading - Works", Financial_Quotes_IPOSecurityNotYetTrading, "70"),
		Entry("HaltCorporateAction - Works", Financial_Quotes_HaltCorporateAction, "71"),
		Entry("QuotationNotAvailable - Works", Financial_Quotes_QuotationNotAvailable, "72"),
		Entry("HaltVolatilityTradingPause - Works", Financial_Quotes_HaltVolatilityTradingPause, "73"),
		Entry("HaltVolatilityTradingPauseStraddleCondition - Works", Financial_Quotes_HaltVolatilityTradingPauseStraddleCondition, "74"),
		Entry("UpdateNewsAndResumptionTimes - Works", Financial_Quotes_UpdateNewsAndResumptionTimes, "75"),
		Entry("HaltSingleStockTradingPauseQuotesOnly - Works", Financial_Quotes_HaltSingleStockTradingPauseQuotesOnly, "76"),
		Entry("ResumeQualificationIssuesReviewedResolved - Works", Financial_Quotes_ResumeQualificationIssuesReviewedResolved, "77"),
		Entry("ResumeFilingRequirementsSatisfiedResolved - Works", Financial_Quotes_ResumeFilingRequirementsSatisfiedResolved, "78"),
		Entry("ResumeNewsNotForthcoming - Works", Financial_Quotes_ResumeNewsNotForthcoming, "79"),
		Entry("ResumeQualificationsMaintRequirementsMet - Works", Financial_Quotes_ResumeQualificationsMaintRequirementsMet, "80"),
		Entry("ResumeQualificationsFilingsMet - Works", Financial_Quotes_ResumeQualificationsFilingsMet, "81"),
		Entry("ResumeRegulatoryAuth - Works", Financial_Quotes_ResumeRegulatoryAuth, "82"),
		Entry("NewIssueAvailable - Works", Financial_Quotes_NewIssueAvailable, "83"),
		Entry("IssueAvailable - Works", Financial_Quotes_IssueAvailable, "84"),
		Entry("MWCBCarryFromPreviousDay - Works", Financial_Quotes_MWCBCarryFromPreviousDay, "85"),
		Entry("MWCBResume - Works", Financial_Quotes_MWCBResume, "86"),
		Entry("IPOSecurityReleasedForQuotation - Works", Financial_Quotes_IPOSecurityReleasedForQuotation, "87"),
		Entry("IPOSecurityPositioningWindowExtension - Works", Financial_Quotes_IPOSecurityPositioningWindowExtension, "88"),
		Entry("MWCBLevel1 - Works", Financial_Quotes_MWCBLevel1, "89"),
		Entry("MWCBLevel2 - Works", Financial_Quotes_MWCBLevel2, "90"),
		Entry("MWCBLevel3 - Works", Financial_Quotes_MWCBLevel3, "91"),
		Entry("HaltSubPennyTrading - Works", Financial_Quotes_HaltSubPennyTrading, "92"),
		Entry("OrderImbalanceInd - Works", Financial_Quotes_OrderImbalanceInd, "93"),
		Entry("LULDTradingPaused - Works", Financial_Quotes_LULDTradingPaused, "94"),
		Entry("NONE - Works", Financial_Quotes_NONE, "95"),
		Entry("ShortSalesRestrictionActivated - Works", Financial_Quotes_ShortSalesRestrictionActivated, "96"),
		Entry("ShortSalesRestrictionContinued - Works", Financial_Quotes_ShortSalesRestrictionContinued, "97"),
		Entry("ShortSalesRestrictionDeactivated - Works", Financial_Quotes_ShortSalesRestrictionDeactivated, "98"),
		Entry("ShortSalesRestrictionInEffect - Works", Financial_Quotes_ShortSalesRestrictionInEffect, "99"),
		Entry("ShortSalesRestrictionMax - Works", Financial_Quotes_ShortSalesRestrictionMax, "100"),
		Entry("NBBONoChange - Works", Financial_Quotes_NBBONoChange, "101"),
		Entry("NBBOQuoteIsNBBO - Works", Financial_Quotes_NBBOQuoteIsNBBO, "102"),
		Entry("NBBONoBBNoBO - Works", Financial_Quotes_NBBONoBBNoBO, "103"),
		Entry("NBBOBBBOShortAppendage - Works", Financial_Quotes_NBBOBBBOShortAppendage, "104"),
		Entry("NBBOBBBOLongAppendage - Works", Financial_Quotes_NBBOBBBOLongAppendage, "105"),
		Entry("HeldTradeNotLastSaleNotConsolidated - Works", Financial_Quotes_HeldTradeNotLastSaleNotConsolidated, "106"),
		Entry("HeldTradeLastSaleButNotConsolidated - Works", Financial_Quotes_HeldTradeLastSaleButNotConsolidated, "107"),
		Entry("HeldTradeLastSaleAndConsolidated - Works", Financial_Quotes_HeldTradeLastSaleAndConsolidated, "108"),
		Entry("RetailInterestOnBid - Works", Financial_Quotes_RetailInterestOnBid, "109"),
		Entry("RetailInterestOnAsk - Works", Financial_Quotes_RetailInterestOnAsk, "110"),
		Entry("RetailInterestOnBidAndAsk - Works", Financial_Quotes_RetailInterestOnBidAndAsk, "111"),
		Entry("FinraBBONoChange - Works", Financial_Quotes_FinraBBONoChange, "112"),
		Entry("FinraBBODoesNotExist - Works", Financial_Quotes_FinraBBODoesNotExist, "113"),
		Entry("FinraBBBOExecutable - Works", Financial_Quotes_FinraBBBOExecutable, "114"),
		Entry("FinraBBBelowLowerBand - Works", Financial_Quotes_FinraBBBelowLowerBand, "115"),
		Entry("FinraBOAboveUpperBand - Works", Financial_Quotes_FinraBOAboveUpperBand, "116"),
		Entry("FinraBBBelowLowerBandBOAbboveUpperBand - Works", Financial_Quotes_FinraBBBelowLowerBandBOAbboveUpperBand, "117"),
		Entry("CTANotDueToRelatedSecurity - Works", Financial_Quotes_CTANotDueToRelatedSecurity, "118"),
		Entry("CTADueToRelatedSecurity - Works", Financial_Quotes_CTADueToRelatedSecurity, "119"),
		Entry("CTANotInViewOfCommon - Works", Financial_Quotes_CTANotInViewOfCommon, "120"),
		Entry("CTAInViewOfCommon - Works", Financial_Quotes_CTAInViewOfCommon, "121"),
		Entry("CTAPriceIndicator - Works", Financial_Quotes_CTAPriceIndicator, "122"),
		Entry("CTANewPriceIndicator - Works", Financial_Quotes_CTANewPriceIndicator, "123"),
		Entry("CTACorrectedPriceIndication - Works", Financial_Quotes_CTACorrectedPriceIndication, "124"),
		Entry("CTACancelledMarketImbalance - Works", Financial_Quotes_CTACancelledMarketImbalance, "125"))

	// Test that converting the Financial.Quotes.Indicator enum to a AttributeValue works for all values
	DescribeTable("MarshalDynamoDBAttributeValue Tests",
		func(enum Financial_Quotes_Indicator, value string) {
			data, err := enum.MarshalDynamoDBAttributeValue()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data.(*types.AttributeValueMemberS).Value).Should(Equal(value))
		},
		Entry("NBBNBOExecutable - Works", Financial_Quotes_NBBNBOExecutable, "NBBNBOExecutable"),
		Entry("NBBBelowLowerBand - Works", Financial_Quotes_NBBBelowLowerBand, "NBBBelowLowerBand"),
		Entry("NBOAboveUpperBand - Works", Financial_Quotes_NBOAboveUpperBand, "NBOAboveUpperBand"),
		Entry("NBBBelowLowerBandAndNBOAboveUpperBand - Works",
			Financial_Quotes_NBBBelowLowerBandAndNBOAboveUpperBand, "NBBBelowLowerBandAndNBOAboveUpperBand"),
		Entry("NBBEqualsUpperBand - Works", Financial_Quotes_NBBEqualsUpperBand, "NBBEqualsUpperBand"),
		Entry("NBOEqualsLowerBand - Works", Financial_Quotes_NBOEqualsLowerBand, "NBOEqualsLowerBand"),
		Entry("NBBEqualsUpperBandAndNBOAboveUpperBand - Works",
			Financial_Quotes_NBBEqualsUpperBandAndNBOAboveUpperBand, "NBBEqualsUpperBandAndNBOAboveUpperBand"),
		Entry("NBBBelowLowerBandAndNBOEqualsLowerBand - Works",
			Financial_Quotes_NBBBelowLowerBandAndNBOEqualsLowerBand, "NBBBelowLowerBandAndNBOEqualsLowerBand"),
		Entry("BidPriceAboveUpperLimitPriceBand - Works",
			Financial_Quotes_BidPriceAboveUpperLimitPriceBand, "BidPriceAboveUpperLimitPriceBand"),
		Entry("OfferPriceBelowLowerLimitPriceBand - Works",
			Financial_Quotes_OfferPriceBelowLowerLimitPriceBand, "OfferPriceBelowLowerLimitPriceBand"),
		Entry("BidAndOfferOutsidePriceBand - Works",
			Financial_Quotes_BidAndOfferOutsidePriceBand, "BidAndOfferOutsidePriceBand"),
		Entry("OpeningUpdate - Works", Financial_Quotes_OpeningUpdate, "OpeningUpdate"),
		Entry("IntraDayUpdate - Works", Financial_Quotes_IntraDayUpdate, "IntraDayUpdate"),
		Entry("RestatedValue - Works", Financial_Quotes_RestatedValue, "RestatedValue"),
		Entry("SuspendedDuringTradingHalt - Works",
			Financial_Quotes_SuspendedDuringTradingHalt, "SuspendedDuringTradingHalt"),
		Entry("ReOpeningUpdate - Works", Financial_Quotes_ReOpeningUpdate, "ReOpeningUpdate"),
		Entry("OutsidePriceBandRuleHours - Works", Financial_Quotes_OutsidePriceBandRuleHours, "OutsidePriceBandRuleHours"),
		Entry("AuctionExtension - Works", Financial_Quotes_AuctionExtension, "AuctionExtension"),
		Entry("LULDPriceBand - Works", Financial_Quotes_LULDPriceBand, "LULDPriceBand"),
		Entry("RepublishedLULDPriceBandInd - Works",
			Financial_Quotes_RepublishedLULDPriceBandInd, "RepublishedLULDPriceBandInd"),
		Entry("NBBLimitStateEntered - Works", Financial_Quotes_NBBLimitStateEntered, "NBBLimitStateEntered"),
		Entry("NBBLimitStateExited - Works", Financial_Quotes_NBBLimitStateExited, "NBBLimitStateExited"),
		Entry("NBOLimitStateEntered - Works", Financial_Quotes_NBOLimitStateEntered, "NBOLimitStateEntered"),
		Entry("NBOLimitStateExited - Works", Financial_Quotes_NBOLimitStateExited, "NBOLimitStateExited"),
		Entry("NBBAndNBOLimitStateEntered - Works",
			Financial_Quotes_NBBAndNBOLimitStateEntered, "NBBAndNBOLimitStateEntered"),
		Entry("NBBAndNBOLimitStateExited - Works", Financial_Quotes_NBBAndNBOLimitStateExited, "NBBAndNBOLimitStateExited"),
		Entry("NBBLimitStateEnteredNBOLimitStateExited - Works",
			Financial_Quotes_NBBLimitStateEnteredNBOLimitStateExited, "NBBLimitStateEnteredNBOLimitStateExited"),
		Entry("NBBLimitStateExitedNBOLimitStateEntered - Works",
			Financial_Quotes_NBBLimitStateExitedNBOLimitStateEntered, "NBBLimitStateExitedNBOLimitStateEntered"),
		Entry("Normal - Works", Financial_Quotes_Normal, "Normal"),
		Entry("Bankrupt - Works", Financial_Quotes_Bankrupt, "Bankrupt"),
		Entry("Deficient - Works", Financial_Quotes_Deficient, "Deficient"),
		Entry("Delinquent - Works", Financial_Quotes_Delinquent, "Delinquent"),
		Entry("BankruptAndDeficient - Works", Financial_Quotes_BankruptAndDeficient, "BankruptAndDeficient"),
		Entry("BankruptAndDelinquent - Works", Financial_Quotes_BankruptAndDelinquent, "BankruptAndDelinquent"),
		Entry("DeficientAndDelinquent - Works", Financial_Quotes_DeficientAndDelinquent, "DeficientAndDelinquent"),
		Entry("DeficientDeliquentBankrupt - Works",
			Financial_Quotes_DeficientDeliquentBankrupt, "DeficientDeliquentBankrupt"),
		Entry("Liquidation - Works", Financial_Quotes_Liquidation, "Liquidation"),
		Entry("CreationsSuspended - Works", Financial_Quotes_CreationsSuspended, "CreationsSuspended"),
		Entry("RedemptionsSuspended - Works", Financial_Quotes_RedemptionsSuspended, "RedemptionsSuspended"),
		Entry("CreationsRedemptionsSuspended - Works",
			Financial_Quotes_CreationsRedemptionsSuspended, "CreationsRedemptionsSuspended"),
		Entry("NormalTrading - Works", Financial_Quotes_NormalTrading, "NormalTrading"),
		Entry("OpeningDelay - Works", Financial_Quotes_OpeningDelay, "OpeningDelay"),
		Entry("TradingHalt - Works", Financial_Quotes_TradingHalt, "TradingHalt"),
		Entry("TradingResume - Works", Financial_Quotes_TradingResume, "TradingResume"),
		Entry("NoOpenNoResume - Works", Financial_Quotes_NoOpenNoResume, "NoOpenNoResume"),
		Entry("PriceIndication - Works", Financial_Quotes_PriceIndication, "PriceIndication"),
		Entry("TradingRangeIndication - Works", Financial_Quotes_TradingRangeIndication, "TradingRangeIndication"),
		Entry("MarketImbalanceBuy - Works", Financial_Quotes_MarketImbalanceBuy, "MarketImbalanceBuy"),
		Entry("MarketImbalanceSell - Works", Financial_Quotes_MarketImbalanceSell, "MarketImbalanceSell"),
		Entry("MarketOnCloseImbalanceBuy - Works", Financial_Quotes_MarketOnCloseImbalanceBuy, "MarketOnCloseImbalanceBuy"),
		Entry("MarketOnCloseImbalanceSell - Works",
			Financial_Quotes_MarketOnCloseImbalanceSell, "MarketOnCloseImbalanceSell"),
		Entry("NoMarketImbalance - Works", Financial_Quotes_NoMarketImbalance, "NoMarketImbalance"),
		Entry("NoMarketOnCloseImbalance - Works", Financial_Quotes_NoMarketOnCloseImbalance, "NoMarketOnCloseImbalance"),
		Entry("ShortSaleRestriction - Works", Financial_Quotes_ShortSaleRestriction, "ShortSaleRestriction"),
		Entry("LimitUpLimitDown - Works", Financial_Quotes_LimitUpLimitDown, "LimitUpLimitDown"),
		Entry("QuotationResumption - Works", Financial_Quotes_QuotationResumption, "QuotationResumption"),
		Entry("TradingResumption - Works", Financial_Quotes_TradingResumption, "TradingResumption"),
		Entry("VolatilityTradingPause - Works", Financial_Quotes_VolatilityTradingPause, "VolatilityTradingPause"),
		Entry("PolygonReserved - Works", Financial_Quotes_PolygonReserved, "PolygonReserved"),
		Entry("HaltNewsPending - Works", Financial_Quotes_HaltNewsPending, "HaltNewsPending"),
		Entry("UpdateNewsDissemination - Works", Financial_Quotes_UpdateNewsDissemination, "UpdateNewsDissemination"),
		Entry("HaltSingleStockTradingPause - Works",
			Financial_Quotes_HaltSingleStockTradingPause, "HaltSingleStockTradingPause"),
		Entry("HaltRegulatoryExtraordinaryMarketActivity - Works",
			Financial_Quotes_HaltRegulatoryExtraordinaryMarketActivity, "HaltRegulatoryExtraordinaryMarketActivity"),
		Entry("HaltETF - Works", Financial_Quotes_HaltETF, "HaltETF"),
		Entry("HaltInformationRequested - Works", Financial_Quotes_HaltInformationRequested, "HaltInformationRequested"),
		Entry("HaltExchangeNonCompliance - Works", Financial_Quotes_HaltExchangeNonCompliance, "HaltExchangeNonCompliance"),
		Entry("HaltFilingsNotCurrent - Works", Financial_Quotes_HaltFilingsNotCurrent, "HaltFilingsNotCurrent"),
		Entry("HaltSECTradingSuspension - Works", Financial_Quotes_HaltSECTradingSuspension, "HaltSECTradingSuspension"),
		Entry("HaltRegulatoryConcern - Works", Financial_Quotes_HaltRegulatoryConcern, "HaltRegulatoryConcern"),
		Entry("HaltMarketOperations - Works", Financial_Quotes_HaltMarketOperations, "HaltMarketOperations"),
		Entry("IPOSecurityNotYetTrading - Works", Financial_Quotes_IPOSecurityNotYetTrading, "IPOSecurityNotYetTrading"),
		Entry("HaltCorporateAction - Works", Financial_Quotes_HaltCorporateAction, "HaltCorporateAction"),
		Entry("QuotationNotAvailable - Works", Financial_Quotes_QuotationNotAvailable, "QuotationNotAvailable"),
		Entry("HaltVolatilityTradingPause - Works",
			Financial_Quotes_HaltVolatilityTradingPause, "HaltVolatilityTradingPause"),
		Entry("HaltVolatilityTradingPauseStraddleCondition - Works",
			Financial_Quotes_HaltVolatilityTradingPauseStraddleCondition, "HaltVolatilityTradingPauseStraddleCondition"),
		Entry("UpdateNewsAndResumptionTimes - Works",
			Financial_Quotes_UpdateNewsAndResumptionTimes, "UpdateNewsAndResumptionTimes"),
		Entry("HaltSingleStockTradingPauseQuotesOnly - Works",
			Financial_Quotes_HaltSingleStockTradingPauseQuotesOnly, "HaltSingleStockTradingPauseQuotesOnly"),
		Entry("ResumeQualificationIssuesReviewedResolved - Works",
			Financial_Quotes_ResumeQualificationIssuesReviewedResolved, "ResumeQualificationIssuesReviewedResolved"),
		Entry("ResumeFilingRequirementsSatisfiedResolved - Works",
			Financial_Quotes_ResumeFilingRequirementsSatisfiedResolved, "ResumeFilingRequirementsSatisfiedResolved"),
		Entry("ResumeNewsNotForthcoming - Works", Financial_Quotes_ResumeNewsNotForthcoming, "ResumeNewsNotForthcoming"),
		Entry("ResumeQualificationsMaintRequirementsMet - Works",
			Financial_Quotes_ResumeQualificationsMaintRequirementsMet, "ResumeQualificationsMaintRequirementsMet"),
		Entry("ResumeQualificationsFilingsMet - Works",
			Financial_Quotes_ResumeQualificationsFilingsMet, "ResumeQualificationsFilingsMet"),
		Entry("ResumeRegulatoryAuth - Works", Financial_Quotes_ResumeRegulatoryAuth, "ResumeRegulatoryAuth"),
		Entry("NewIssueAvailable - Works", Financial_Quotes_NewIssueAvailable, "NewIssueAvailable"),
		Entry("IssueAvailable - Works", Financial_Quotes_IssueAvailable, "IssueAvailable"),
		Entry("MWCBCarryFromPreviousDay - Works", Financial_Quotes_MWCBCarryFromPreviousDay, "MWCBCarryFromPreviousDay"),
		Entry("MWCBResume - Works", Financial_Quotes_MWCBResume, "MWCBResume"),
		Entry("IPOSecurityReleasedForQuotation - Works",
			Financial_Quotes_IPOSecurityReleasedForQuotation, "IPOSecurityReleasedForQuotation"),
		Entry("IPOSecurityPositioningWindowExtension - Works",
			Financial_Quotes_IPOSecurityPositioningWindowExtension, "IPOSecurityPositioningWindowExtension"),
		Entry("MWCBLevel1 - Works", Financial_Quotes_MWCBLevel1, "MWCBLevel1"),
		Entry("MWCBLevel2 - Works", Financial_Quotes_MWCBLevel2, "MWCBLevel2"),
		Entry("MWCBLevel3 - Works", Financial_Quotes_MWCBLevel3, "MWCBLevel3"),
		Entry("HaltSubPennyTrading - Works", Financial_Quotes_HaltSubPennyTrading, "HaltSubPennyTrading"),
		Entry("OrderImbalanceInd - Works", Financial_Quotes_OrderImbalanceInd, "OrderImbalanceInd"),
		Entry("LULDTradingPaused - Works", Financial_Quotes_LULDTradingPaused, "LULDTradingPaused"),
		Entry("NONE - Works", Financial_Quotes_NONE, "NONE"),
		Entry("ShortSalesRestrictionActivated - Works",
			Financial_Quotes_ShortSalesRestrictionActivated, "ShortSalesRestrictionActivated"),
		Entry("ShortSalesRestrictionContinued - Works",
			Financial_Quotes_ShortSalesRestrictionContinued, "ShortSalesRestrictionContinued"),
		Entry("ShortSalesRestrictionDeactivated - Works",
			Financial_Quotes_ShortSalesRestrictionDeactivated, "ShortSalesRestrictionDeactivated"),
		Entry("ShortSalesRestrictionInEffect - Works",
			Financial_Quotes_ShortSalesRestrictionInEffect, "ShortSalesRestrictionInEffect"),
		Entry("ShortSalesRestrictionMax - Works", Financial_Quotes_ShortSalesRestrictionMax, "ShortSalesRestrictionMax"),
		Entry("RetailInterestOnBid - Works", Financial_Quotes_RetailInterestOnBid, "RetailInterestOnBid"),
		Entry("RetailInterestOnAsk - Works", Financial_Quotes_RetailInterestOnAsk, "RetailInterestOnAsk"),
		Entry("RetailInterestOnBidAndAsk - Works", Financial_Quotes_RetailInterestOnBidAndAsk, "RetailInterestOnBidAndAsk"),
		Entry("FinraBBONoChange - Works", Financial_Quotes_FinraBBONoChange, "FinraBBONoChange"),
		Entry("FinraBBODoesNotExist - Works", Financial_Quotes_FinraBBODoesNotExist, "FinraBBODoesNotExist"),
		Entry("FinraBBBOExecutable - Works", Financial_Quotes_FinraBBBOExecutable, "FinraBBBOExecutable"),
		Entry("FinraBBBelowLowerBand - Works", Financial_Quotes_FinraBBBelowLowerBand, "FinraBBBelowLowerBand"),
		Entry("FinraBOAboveUpperBand - Works", Financial_Quotes_FinraBOAboveUpperBand, "FinraBOAboveUpperBand"),
		Entry("FinraBBBelowLowerBandBOAbboveUpperBand - Works",
			Financial_Quotes_FinraBBBelowLowerBandBOAbboveUpperBand, "FinraBBBelowLowerBandBOAbboveUpperBand"),
		Entry("NBBONoChange - Works", Financial_Quotes_NBBONoChange, "NBBONoChange"),
		Entry("NBBOQuoteIsNBBO - Works", Financial_Quotes_NBBOQuoteIsNBBO, "NBBOQuoteIsNBBO"),
		Entry("NBBONoBBNoBO - Works", Financial_Quotes_NBBONoBBNoBO, "NBBONoBBNoBO"),
		Entry("NBBOBBBOShortAppendage - Works", Financial_Quotes_NBBOBBBOShortAppendage, "NBBOBBBOShortAppendage"),
		Entry("NBBOBBBOLongAppendage - Works", Financial_Quotes_NBBOBBBOLongAppendage, "NBBOBBBOLongAppendage"),
		Entry("HeldTradeNotLastSaleNotConsolidated - Works",
			Financial_Quotes_HeldTradeNotLastSaleNotConsolidated, "HeldTradeNotLastSaleNotConsolidated"),
		Entry("HeldTradeLastSaleButNotConsolidated - Works",
			Financial_Quotes_HeldTradeLastSaleButNotConsolidated, "HeldTradeLastSaleButNotConsolidated"),
		Entry("HeldTradeLastSaleAndConsolidated - Works",
			Financial_Quotes_HeldTradeLastSaleAndConsolidated, "HeldTradeLastSaleAndConsolidated"),
		Entry("CTANotDueToRelatedSecurity - Works",
			Financial_Quotes_CTANotDueToRelatedSecurity, "CTANotDueToRelatedSecurity"),
		Entry("CTADueToRelatedSecurity - Works", Financial_Quotes_CTADueToRelatedSecurity, "CTADueToRelatedSecurity"),
		Entry("CTANotInViewOfCommon - Works", Financial_Quotes_CTANotInViewOfCommon, "CTANotInViewOfCommon"),
		Entry("CTAInViewOfCommon - Works", Financial_Quotes_CTAInViewOfCommon, "CTAInViewOfCommon"),
		Entry("CTAPriceIndicator - Works", Financial_Quotes_CTAPriceIndicator, "CTAPriceIndicator"),
		Entry("CTANewPriceIndicator - Works", Financial_Quotes_CTANewPriceIndicator, "CTANewPriceIndicator"),
		Entry("CTACorrectedPriceIndication - Works",
			Financial_Quotes_CTACorrectedPriceIndication, "CTACorrectedPriceIndication"),
		Entry("CTACancelledMarketImbalance - Works",
			Financial_Quotes_CTACancelledMarketImbalance, "CTACancelledMarketImbalance"))

	// Test that attempting to deserialize a Financial.Quotes.Indicator will fail and
	// return an error if the value canno be deserialized from a JSON value to a string
	It("UnmarshalJSON fails - Error", func() {

		// Attempt to convert a non-parseable string value into a Financial.Quotes.Indicator
		// This should return an error
		enum := new(Financial_Quotes_Indicator)
		err := enum.UnmarshalJSON([]byte("derp"))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Quotes_Indicator"))
	})

	// Test that attempting to deserialize a Financial.Quotes.Indicator will fail and
	// return an error if the value cannot be converted to either the name value or integer
	// value of the enum option
	It("UnmarshalJSON - Value is invalid - Error", func() {

		// Attempt to convert a fake string value into a Financial.Quotes.Indicator
		// This should return an error
		enum := new(Financial_Quotes_Indicator)
		err := enum.UnmarshalJSON([]byte("\"derp\""))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Quotes_Indicator"))
	})

	// Test the conditions under which values should be convertible to a Financial.Quotes.Indicator
	DescribeTable("UnmarshalJSON Tests",
		func(value interface{}, shouldBe Financial_Quotes_Indicator) {

			// Attempt to convert the string value into a Financial.Quotes.Indicator
			// This should not fail
			var enum Financial_Quotes_Indicator
			err := enum.UnmarshalJSON([]byte(fmt.Sprintf("%v", value)))

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("NBBNBOExecutable - Works", "\"NBBNBOExecutable\"", Financial_Quotes_NBBNBOExecutable),
		Entry("NBBBelowLowerBand - Works", "\"NBBBelowLowerBand\"", Financial_Quotes_NBBBelowLowerBand),
		Entry("NBOAboveUpperBand - Works", "\"NBOAboveUpperBand\"", Financial_Quotes_NBOAboveUpperBand),
		Entry("NBBBelowLowerBandAndNBOAboveUpperBand - Works",
			"\"NBBBelowLowerBandAndNBOAboveUpperBand\"", Financial_Quotes_NBBBelowLowerBandAndNBOAboveUpperBand),
		Entry("NBBEqualsUpperBand - Works", "\"NBBEqualsUpperBand\"", Financial_Quotes_NBBEqualsUpperBand),
		Entry("NBOEqualsLowerBand - Works", "\"NBOEqualsLowerBand\"", Financial_Quotes_NBOEqualsLowerBand),
		Entry("NBBEqualsUpperBandAndNBOAboveUpperBand - Works",
			"\"NBBEqualsUpperBandAndNBOAboveUpperBand\"", Financial_Quotes_NBBEqualsUpperBandAndNBOAboveUpperBand),
		Entry("NBBBelowLowerBandAndNBOEqualsLowerBand - Works",
			"\"NBBBelowLowerBandAndNBOEqualsLowerBand\"", Financial_Quotes_NBBBelowLowerBandAndNBOEqualsLowerBand),
		Entry("BidPriceAboveUpperLimitPriceBand - Works",
			"\"BidPriceAboveUpperLimitPriceBand\"", Financial_Quotes_BidPriceAboveUpperLimitPriceBand),
		Entry("OfferPriceBelowLowerLimitPriceBand - Works",
			"\"OfferPriceBelowLowerLimitPriceBand\"", Financial_Quotes_OfferPriceBelowLowerLimitPriceBand),
		Entry("BidAndOfferOutsidePriceBand - Works",
			"\"BidAndOfferOutsidePriceBand\"", Financial_Quotes_BidAndOfferOutsidePriceBand),
		Entry("OpeningUpdate - Works", "\"OpeningUpdate\"", Financial_Quotes_OpeningUpdate),
		Entry("IntraDayUpdate - Works", "\"IntraDayUpdate\"", Financial_Quotes_IntraDayUpdate),
		Entry("RestatedValue - Works", "\"RestatedValue\"", Financial_Quotes_RestatedValue),
		Entry("SuspendedDuringTradingHalt - Works",
			"\"SuspendedDuringTradingHalt\"", Financial_Quotes_SuspendedDuringTradingHalt),
		Entry("ReOpeningUpdate - Works", "\"ReOpeningUpdate\"", Financial_Quotes_ReOpeningUpdate),
		Entry("OutsidePriceBandRuleHours - Works", "\"OutsidePriceBandRuleHours\"", Financial_Quotes_OutsidePriceBandRuleHours),
		Entry("AuctionExtension - Works", "\"AuctionExtension\"", Financial_Quotes_AuctionExtension),
		Entry("LULDPriceBand - Works", "\"LULDPriceBand\"", Financial_Quotes_LULDPriceBand),
		Entry("RepublishedLULDPriceBandInd - Works",
			"\"RepublishedLULDPriceBandInd\"", Financial_Quotes_RepublishedLULDPriceBandInd),
		Entry("NBBLimitStateEntered - Works", "\"NBBLimitStateEntered\"", Financial_Quotes_NBBLimitStateEntered),
		Entry("NBBLimitStateExited - Works", "\"NBBLimitStateExited\"", Financial_Quotes_NBBLimitStateExited),
		Entry("NBOLimitStateEntered - Works", "\"NBOLimitStateEntered\"", Financial_Quotes_NBOLimitStateEntered),
		Entry("NBOLimitStateExited - Works", "\"NBOLimitStateExited\"", Financial_Quotes_NBOLimitStateExited),
		Entry("NBBAndNBOLimitStateEntered - Works",
			"\"NBBAndNBOLimitStateEntered\"", Financial_Quotes_NBBAndNBOLimitStateEntered),
		Entry("NBBAndNBOLimitStateExited - Works",
			"\"NBBAndNBOLimitStateExited\"", Financial_Quotes_NBBAndNBOLimitStateExited),
		Entry("NBBLimitStateEnteredNBOLimitStateExited - Works",
			"\"NBBLimitStateEnteredNBOLimitStateExited\"", Financial_Quotes_NBBLimitStateEnteredNBOLimitStateExited),
		Entry("NBBLimitStateExitedNBOLimitStateEntered - Works",
			"\"NBBLimitStateExitedNBOLimitStateEntered\"", Financial_Quotes_NBBLimitStateExitedNBOLimitStateEntered),
		Entry("Normal - Works", "\"Normal\"", Financial_Quotes_Normal),
		Entry("Bankrupt - Works", "\"Bankrupt\"", Financial_Quotes_Bankrupt),
		Entry("Deficient - Works", "\"Deficient\"", Financial_Quotes_Deficient),
		Entry("Delinquent - Works", "\"Delinquent\"", Financial_Quotes_Delinquent),
		Entry("BankruptAndDeficient - Works", "\"BankruptAndDeficient\"", Financial_Quotes_BankruptAndDeficient),
		Entry("BankruptAndDelinquent - Works", "\"BankruptAndDelinquent\"", Financial_Quotes_BankruptAndDelinquent),
		Entry("DeficientAndDelinquent - Works", "\"DeficientAndDelinquent\"", Financial_Quotes_DeficientAndDelinquent),
		Entry("DeficientDeliquentBankrupt - Works",
			"\"DeficientDeliquentBankrupt\"", Financial_Quotes_DeficientDeliquentBankrupt),
		Entry("Liquidation - Works", "\"Liquidation\"", Financial_Quotes_Liquidation),
		Entry("CreationsSuspended - Works", "\"CreationsSuspended\"", Financial_Quotes_CreationsSuspended),
		Entry("RedemptionsSuspended - Works", "\"RedemptionsSuspended\"", Financial_Quotes_RedemptionsSuspended),
		Entry("CreationsRedemptionsSuspended - Works",
			"\"CreationsRedemptionsSuspended\"", Financial_Quotes_CreationsRedemptionsSuspended),
		Entry("NormalTrading - Works", "\"NormalTrading\"", Financial_Quotes_NormalTrading),
		Entry("OpeningDelay - Works", "\"OpeningDelay\"", Financial_Quotes_OpeningDelay),
		Entry("TradingHalt - Works", "\"TradingHalt\"", Financial_Quotes_TradingHalt),
		Entry("TradingResume - Works", "\"TradingResume\"", Financial_Quotes_TradingResume),
		Entry("NoOpenNoResume - Works", "\"NoOpenNoResume\"", Financial_Quotes_NoOpenNoResume),
		Entry("PriceIndication - Works", "\"PriceIndication\"", Financial_Quotes_PriceIndication),
		Entry("TradingRangeIndication - Works", "\"TradingRangeIndication\"", Financial_Quotes_TradingRangeIndication),
		Entry("MarketImbalanceBuy - Works", "\"MarketImbalanceBuy\"", Financial_Quotes_MarketImbalanceBuy),
		Entry("MarketImbalanceSell - Works", "\"MarketImbalanceSell\"", Financial_Quotes_MarketImbalanceSell),
		Entry("MarketOnCloseImbalanceBuy - Works", "\"MarketOnCloseImbalanceBuy\"", Financial_Quotes_MarketOnCloseImbalanceBuy),
		Entry("MarketOnCloseImbalanceSell - Works",
			"\"MarketOnCloseImbalanceSell\"", Financial_Quotes_MarketOnCloseImbalanceSell),
		Entry("NoMarketImbalance - Works", "\"NoMarketImbalance\"", Financial_Quotes_NoMarketImbalance),
		Entry("NoMarketOnCloseImbalance - Works", "\"NoMarketOnCloseImbalance\"", Financial_Quotes_NoMarketOnCloseImbalance),
		Entry("ShortSaleRestriction - Works", "\"ShortSaleRestriction\"", Financial_Quotes_ShortSaleRestriction),
		Entry("LimitUpLimitDown - Works", "\"LimitUpLimitDown\"", Financial_Quotes_LimitUpLimitDown),
		Entry("QuotationResumption - Works", "\"QuotationResumption\"", Financial_Quotes_QuotationResumption),
		Entry("TradingResumption - Works", "\"TradingResumption\"", Financial_Quotes_TradingResumption),
		Entry("VolatilityTradingPause - Works", "\"VolatilityTradingPause\"", Financial_Quotes_VolatilityTradingPause),
		Entry("PolygonReserved - Works", "\"PolygonReserved\"", Financial_Quotes_PolygonReserved),
		Entry("HaltNewsPending - Works", "\"HaltNewsPending\"", Financial_Quotes_HaltNewsPending),
		Entry("UpdateNewsDissemination - Works", "\"UpdateNewsDissemination\"", Financial_Quotes_UpdateNewsDissemination),
		Entry("HaltSingleStockTradingPause - Works",
			"\"HaltSingleStockTradingPause\"", Financial_Quotes_HaltSingleStockTradingPause),
		Entry("HaltRegulatoryExtraordinaryMarketActivity - Works",
			"\"HaltRegulatoryExtraordinaryMarketActivity\"", Financial_Quotes_HaltRegulatoryExtraordinaryMarketActivity),
		Entry("HaltETF - Works", "\"HaltETF\"", Financial_Quotes_HaltETF),
		Entry("HaltInformationRequested - Works", "\"HaltInformationRequested\"", Financial_Quotes_HaltInformationRequested),
		Entry("HaltExchangeNonCompliance - Works", "\"HaltExchangeNonCompliance\"", Financial_Quotes_HaltExchangeNonCompliance),
		Entry("HaltFilingsNotCurrent - Works", "\"HaltFilingsNotCurrent\"", Financial_Quotes_HaltFilingsNotCurrent),
		Entry("HaltSECTradingSuspension - Works", "\"HaltSECTradingSuspension\"", Financial_Quotes_HaltSECTradingSuspension),
		Entry("HaltRegulatoryConcern - Works", "\"HaltRegulatoryConcern\"", Financial_Quotes_HaltRegulatoryConcern),
		Entry("HaltMarketOperations - Works", "\"HaltMarketOperations\"", Financial_Quotes_HaltMarketOperations),
		Entry("IPOSecurityNotYetTrading - Works", "\"IPOSecurityNotYetTrading\"", Financial_Quotes_IPOSecurityNotYetTrading),
		Entry("HaltCorporateAction - Works", "\"HaltCorporateAction\"", Financial_Quotes_HaltCorporateAction),
		Entry("QuotationNotAvailable - Works", "\"QuotationNotAvailable\"", Financial_Quotes_QuotationNotAvailable),
		Entry("HaltVolatilityTradingPause - Works",
			"\"HaltVolatilityTradingPause\"", Financial_Quotes_HaltVolatilityTradingPause),
		Entry("HaltVolatilityTradingPauseStraddleCondition - Works",
			"\"HaltVolatilityTradingPauseStraddleCondition\"", Financial_Quotes_HaltVolatilityTradingPauseStraddleCondition),
		Entry("UpdateNewsAndResumptionTimes - Works",
			"\"UpdateNewsAndResumptionTimes\"", Financial_Quotes_UpdateNewsAndResumptionTimes),
		Entry("HaltSingleStockTradingPauseQuotesOnly - Works",
			"\"HaltSingleStockTradingPauseQuotesOnly\"", Financial_Quotes_HaltSingleStockTradingPauseQuotesOnly),
		Entry("ResumeQualificationIssuesReviewedResolved - Works",
			"\"ResumeQualificationIssuesReviewedResolved\"", Financial_Quotes_ResumeQualificationIssuesReviewedResolved),
		Entry("ResumeFilingRequirementsSatisfiedResolved - Works",
			"\"ResumeFilingRequirementsSatisfiedResolved\"", Financial_Quotes_ResumeFilingRequirementsSatisfiedResolved),
		Entry("ResumeNewsNotForthcoming - Works", "\"ResumeNewsNotForthcoming\"", Financial_Quotes_ResumeNewsNotForthcoming),
		Entry("ResumeQualificationsMaintRequirementsMet - Works",
			"\"ResumeQualificationsMaintRequirementsMet\"", Financial_Quotes_ResumeQualificationsMaintRequirementsMet),
		Entry("ResumeQualificationsFilingsMet - Works",
			"\"ResumeQualificationsFilingsMet\"", Financial_Quotes_ResumeQualificationsFilingsMet),
		Entry("ResumeRegulatoryAuth - Works", "\"ResumeRegulatoryAuth\"", Financial_Quotes_ResumeRegulatoryAuth),
		Entry("NewIssueAvailable - Works", "\"NewIssueAvailable\"", Financial_Quotes_NewIssueAvailable),
		Entry("IssueAvailable - Works", "\"IssueAvailable\"", Financial_Quotes_IssueAvailable),
		Entry("MWCBCarryFromPreviousDay - Works", "\"MWCBCarryFromPreviousDay\"", Financial_Quotes_MWCBCarryFromPreviousDay),
		Entry("MWCBResume - Works", "\"MWCBResume\"", Financial_Quotes_MWCBResume),
		Entry("IPOSecurityReleasedForQuotation - Works",
			"\"IPOSecurityReleasedForQuotation\"", Financial_Quotes_IPOSecurityReleasedForQuotation),
		Entry("IPOSecurityPositioningWindowExtension - Works",
			"\"IPOSecurityPositioningWindowExtension\"", Financial_Quotes_IPOSecurityPositioningWindowExtension),
		Entry("MWCBLevel1 - Works", "\"MWCBLevel1\"", Financial_Quotes_MWCBLevel1),
		Entry("MWCBLevel2 - Works", "\"MWCBLevel2\"", Financial_Quotes_MWCBLevel2),
		Entry("MWCBLevel3 - Works", "\"MWCBLevel3\"", Financial_Quotes_MWCBLevel3),
		Entry("HaltSubPennyTrading - Works", "\"HaltSubPennyTrading\"", Financial_Quotes_HaltSubPennyTrading),
		Entry("OrderImbalanceInd - Works", "\"OrderImbalanceInd\"", Financial_Quotes_OrderImbalanceInd),
		Entry("LULDTradingPaused - Works", "\"LULDTradingPaused\"", Financial_Quotes_LULDTradingPaused),
		Entry("NONE - Works", "\"NONE\"", Financial_Quotes_NONE),
		Entry("ShortSalesRestrictionActivated - Works",
			"\"ShortSalesRestrictionActivated\"", Financial_Quotes_ShortSalesRestrictionActivated),
		Entry("ShortSalesRestrictionContinued - Works",
			"\"ShortSalesRestrictionContinued\"", Financial_Quotes_ShortSalesRestrictionContinued),
		Entry("ShortSalesRestrictionDeactivated - Works",
			"\"ShortSalesRestrictionDeactivated\"", Financial_Quotes_ShortSalesRestrictionDeactivated),
		Entry("ShortSalesRestrictionInEffect - Works",
			"\"ShortSalesRestrictionInEffect\"", Financial_Quotes_ShortSalesRestrictionInEffect),
		Entry("ShortSalesRestrictionMax - Works", "\"ShortSalesRestrictionMax\"", Financial_Quotes_ShortSalesRestrictionMax),
		Entry("RetailInterestOnBid - Works", "\"RetailInterestOnBid\"", Financial_Quotes_RetailInterestOnBid),
		Entry("RetailInterestOnAsk - Works", "\"RetailInterestOnAsk\"", Financial_Quotes_RetailInterestOnAsk),
		Entry("RetailInterestOnBidAndAsk - Works", "\"RetailInterestOnBidAndAsk\"", Financial_Quotes_RetailInterestOnBidAndAsk),
		Entry("FinraBBONoChange - Works", "\"FinraBBONoChange\"", Financial_Quotes_FinraBBONoChange),
		Entry("FinraBBODoesNotExist - Works", "\"FinraBBODoesNotExist\"", Financial_Quotes_FinraBBODoesNotExist),
		Entry("FinraBBBOExecutable - Works", "\"FinraBBBOExecutable\"", Financial_Quotes_FinraBBBOExecutable),
		Entry("FinraBBBelowLowerBand - Works", "\"FinraBBBelowLowerBand\"", Financial_Quotes_FinraBBBelowLowerBand),
		Entry("FinraBOAboveUpperBand - Works", "\"FinraBOAboveUpperBand\"", Financial_Quotes_FinraBOAboveUpperBand),
		Entry("FinraBBBelowLowerBandBOAbboveUpperBand - Works",
			"\"FinraBBBelowLowerBandBOAbboveUpperBand\"", Financial_Quotes_FinraBBBelowLowerBandBOAbboveUpperBand),
		Entry("NBBONoChange - Works", "\"NBBONoChange\"", Financial_Quotes_NBBONoChange),
		Entry("NBBOQuoteIsNBBO - Works", "\"NBBOQuoteIsNBBO\"", Financial_Quotes_NBBOQuoteIsNBBO),
		Entry("NBBONoBBNoBO - Works", "\"NBBONoBBNoBO\"", Financial_Quotes_NBBONoBBNoBO),
		Entry("NBBOBBBOShortAppendage - Works", "\"NBBOBBBOShortAppendage\"", Financial_Quotes_NBBOBBBOShortAppendage),
		Entry("NBBOBBBOLongAppendage - Works", "\"NBBOBBBOLongAppendage\"", Financial_Quotes_NBBOBBBOLongAppendage),
		Entry("HeldTradeNotLastSaleNotConsolidated - Works",
			"\"HeldTradeNotLastSaleNotConsolidated\"", Financial_Quotes_HeldTradeNotLastSaleNotConsolidated),
		Entry("HeldTradeLastSaleButNotConsolidated - Works",
			"\"HeldTradeLastSaleButNotConsolidated\"", Financial_Quotes_HeldTradeLastSaleButNotConsolidated),
		Entry("HeldTradeLastSaleAndConsolidated - Works",
			"\"HeldTradeLastSaleAndConsolidated\"", Financial_Quotes_HeldTradeLastSaleAndConsolidated),
		Entry("CTANotDueToRelatedSecurity - Works",
			"\"CTANotDueToRelatedSecurity\"", Financial_Quotes_CTANotDueToRelatedSecurity),
		Entry("CTADueToRelatedSecurity - Works", "\"CTADueToRelatedSecurity\"", Financial_Quotes_CTADueToRelatedSecurity),
		Entry("CTANotInViewOfCommon - Works", "\"CTANotInViewOfCommon\"", Financial_Quotes_CTANotInViewOfCommon),
		Entry("CTAInViewOfCommon - Works", "\"CTAInViewOfCommon\"", Financial_Quotes_CTAInViewOfCommon),
		Entry("CTAPriceIndicator - Works", "\"CTAPriceIndicator\"", Financial_Quotes_CTAPriceIndicator),
		Entry("CTANewPriceIndicator - Works", "\"CTANewPriceIndicator\"", Financial_Quotes_CTANewPriceIndicator),
		Entry("CTACorrectedPriceIndication - Works",
			"\"CTACorrectedPriceIndication\"", Financial_Quotes_CTACorrectedPriceIndication),
		Entry("CTACancelledMarketImbalance - Works",
			"\"CTACancelledMarketImbalance\"", Financial_Quotes_CTACancelledMarketImbalance),
		Entry("'0' - Works", "\"0\"", Financial_Quotes_NBBNBOExecutable),
		Entry("'1' - Works", "\"1\"", Financial_Quotes_NBBBelowLowerBand),
		Entry("'2' - Works", "\"2\"", Financial_Quotes_NBOAboveUpperBand),
		Entry("'3' - Works", "\"3\"", Financial_Quotes_NBBBelowLowerBandAndNBOAboveUpperBand),
		Entry("'4' - Works", "\"4\"", Financial_Quotes_NBBEqualsUpperBand),
		Entry("'5' - Works", "\"5\"", Financial_Quotes_NBOEqualsLowerBand),
		Entry("'6' - Works", "\"6\"", Financial_Quotes_NBBEqualsUpperBandAndNBOAboveUpperBand),
		Entry("'7' - Works", "\"7\"", Financial_Quotes_NBBBelowLowerBandAndNBOEqualsLowerBand),
		Entry("'8' - Works", "\"8\"", Financial_Quotes_BidPriceAboveUpperLimitPriceBand),
		Entry("'9' - Works", "\"9\"", Financial_Quotes_OfferPriceBelowLowerLimitPriceBand),
		Entry("'10' - Works", "\"10\"", Financial_Quotes_BidAndOfferOutsidePriceBand),
		Entry("'11' - Works", "\"11\"", Financial_Quotes_OpeningUpdate),
		Entry("'12' - Works", "\"12\"", Financial_Quotes_IntraDayUpdate),
		Entry("'13' - Works", "\"13\"", Financial_Quotes_RestatedValue),
		Entry("'14' - Works", "\"14\"", Financial_Quotes_SuspendedDuringTradingHalt),
		Entry("'15' - Works", "\"15\"", Financial_Quotes_ReOpeningUpdate),
		Entry("'16' - Works", "\"16\"", Financial_Quotes_OutsidePriceBandRuleHours),
		Entry("'17' - Works", "\"17\"", Financial_Quotes_AuctionExtension),
		Entry("'18' - Works", "\"18\"", Financial_Quotes_LULDPriceBand),
		Entry("'19' - Works", "\"19\"", Financial_Quotes_RepublishedLULDPriceBandInd),
		Entry("'20' - Works", "\"20\"", Financial_Quotes_NBBLimitStateEntered),
		Entry("'21' - Works", "\"21\"", Financial_Quotes_NBBLimitStateExited),
		Entry("'22' - Works", "\"22\"", Financial_Quotes_NBOLimitStateEntered),
		Entry("'23' - Works", "\"23\"", Financial_Quotes_NBOLimitStateExited),
		Entry("'24' - Works", "\"24\"", Financial_Quotes_NBBAndNBOLimitStateEntered),
		Entry("'25' - Works", "\"25\"", Financial_Quotes_NBBAndNBOLimitStateExited),
		Entry("'26' - Works", "\"26\"", Financial_Quotes_NBBLimitStateEnteredNBOLimitStateExited),
		Entry("'27' - Works", "\"27\"", Financial_Quotes_NBBLimitStateExitedNBOLimitStateEntered),
		Entry("'28' - Works", "\"28\"", Financial_Quotes_Normal),
		Entry("'29' - Works", "\"29\"", Financial_Quotes_Bankrupt),
		Entry("'30' - Works", "\"30\"", Financial_Quotes_Deficient),
		Entry("'31' - Works", "\"31\"", Financial_Quotes_Delinquent),
		Entry("'32' - Works", "\"32\"", Financial_Quotes_BankruptAndDeficient),
		Entry("'33' - Works", "\"33\"", Financial_Quotes_BankruptAndDelinquent),
		Entry("'34' - Works", "\"34\"", Financial_Quotes_DeficientAndDelinquent),
		Entry("'35' - Works", "\"35\"", Financial_Quotes_DeficientDeliquentBankrupt),
		Entry("'36' - Works", "\"36\"", Financial_Quotes_Liquidation),
		Entry("'37' - Works", "\"37\"", Financial_Quotes_CreationsSuspended),
		Entry("'38' - Works", "\"38\"", Financial_Quotes_RedemptionsSuspended),
		Entry("'39' - Works", "\"39\"", Financial_Quotes_CreationsRedemptionsSuspended),
		Entry("'40' - Works", "\"40\"", Financial_Quotes_NormalTrading),
		Entry("'41' - Works", "\"41\"", Financial_Quotes_OpeningDelay),
		Entry("'42' - Works", "\"42\"", Financial_Quotes_TradingHalt),
		Entry("'43' - Works", "\"43\"", Financial_Quotes_TradingResume),
		Entry("'44' - Works", "\"44\"", Financial_Quotes_NoOpenNoResume),
		Entry("'45' - Works", "\"45\"", Financial_Quotes_PriceIndication),
		Entry("'46' - Works", "\"46\"", Financial_Quotes_TradingRangeIndication),
		Entry("'47' - Works", "\"47\"", Financial_Quotes_MarketImbalanceBuy),
		Entry("'48' - Works", "\"48\"", Financial_Quotes_MarketImbalanceSell),
		Entry("'49' - Works", "\"49\"", Financial_Quotes_MarketOnCloseImbalanceBuy),
		Entry("'50' - Works", "\"50\"", Financial_Quotes_MarketOnCloseImbalanceSell),
		Entry("'51' - Works", "\"51\"", Financial_Quotes_NoMarketImbalance),
		Entry("'52' - Works", "\"52\"", Financial_Quotes_NoMarketOnCloseImbalance),
		Entry("'53' - Works", "\"53\"", Financial_Quotes_ShortSaleRestriction),
		Entry("'54' - Works", "\"54\"", Financial_Quotes_LimitUpLimitDown),
		Entry("'55' - Works", "\"55\"", Financial_Quotes_QuotationResumption),
		Entry("'56' - Works", "\"56\"", Financial_Quotes_TradingResumption),
		Entry("'57' - Works", "\"57\"", Financial_Quotes_VolatilityTradingPause),
		Entry("'58' - Works", "\"58\"", Financial_Quotes_PolygonReserved),
		Entry("'59' - Works", "\"59\"", Financial_Quotes_HaltNewsPending),
		Entry("'60' - Works", "\"60\"", Financial_Quotes_UpdateNewsDissemination),
		Entry("'61' - Works", "\"61\"", Financial_Quotes_HaltSingleStockTradingPause),
		Entry("'62' - Works", "\"62\"", Financial_Quotes_HaltRegulatoryExtraordinaryMarketActivity),
		Entry("'63' - Works", "\"63\"", Financial_Quotes_HaltETF),
		Entry("'64' - Works", "\"64\"", Financial_Quotes_HaltInformationRequested),
		Entry("'65' - Works", "\"65\"", Financial_Quotes_HaltExchangeNonCompliance),
		Entry("'66' - Works", "\"66\"", Financial_Quotes_HaltFilingsNotCurrent),
		Entry("'67' - Works", "\"67\"", Financial_Quotes_HaltSECTradingSuspension),
		Entry("'68' - Works", "\"68\"", Financial_Quotes_HaltRegulatoryConcern),
		Entry("'69' - Works", "\"69\"", Financial_Quotes_HaltMarketOperations),
		Entry("'70' - Works", "\"70\"", Financial_Quotes_IPOSecurityNotYetTrading),
		Entry("'71' - Works", "\"71\"", Financial_Quotes_HaltCorporateAction),
		Entry("'72' - Works", "\"72\"", Financial_Quotes_QuotationNotAvailable),
		Entry("'73' - Works", "\"73\"", Financial_Quotes_HaltVolatilityTradingPause),
		Entry("'74' - Works", "\"74\"", Financial_Quotes_HaltVolatilityTradingPauseStraddleCondition),
		Entry("'75' - Works", "\"75\"", Financial_Quotes_UpdateNewsAndResumptionTimes),
		Entry("'76' - Works", "\"76\"", Financial_Quotes_HaltSingleStockTradingPauseQuotesOnly),
		Entry("'77' - Works", "\"77\"", Financial_Quotes_ResumeQualificationIssuesReviewedResolved),
		Entry("'78' - Works", "\"78\"", Financial_Quotes_ResumeFilingRequirementsSatisfiedResolved),
		Entry("'79' - Works", "\"79\"", Financial_Quotes_ResumeNewsNotForthcoming),
		Entry("'80' - Works", "\"80\"", Financial_Quotes_ResumeQualificationsMaintRequirementsMet),
		Entry("'81' - Works", "\"81\"", Financial_Quotes_ResumeQualificationsFilingsMet),
		Entry("'82' - Works", "\"82\"", Financial_Quotes_ResumeRegulatoryAuth),
		Entry("'83' - Works", "\"83\"", Financial_Quotes_NewIssueAvailable),
		Entry("'84' - Works", "\"84\"", Financial_Quotes_IssueAvailable),
		Entry("'85' - Works", "\"85\"", Financial_Quotes_MWCBCarryFromPreviousDay),
		Entry("'86' - Works", "\"86\"", Financial_Quotes_MWCBResume),
		Entry("'87' - Works", "\"87\"", Financial_Quotes_IPOSecurityReleasedForQuotation),
		Entry("'88' - Works", "\"88\"", Financial_Quotes_IPOSecurityPositioningWindowExtension),
		Entry("'89' - Works", "\"89\"", Financial_Quotes_MWCBLevel1),
		Entry("'90' - Works", "\"90\"", Financial_Quotes_MWCBLevel2),
		Entry("'91' - Works", "\"91\"", Financial_Quotes_MWCBLevel3),
		Entry("'92' - Works", "\"92\"", Financial_Quotes_HaltSubPennyTrading),
		Entry("'93' - Works", "\"93\"", Financial_Quotes_OrderImbalanceInd),
		Entry("'94' - Works", "\"94\"", Financial_Quotes_LULDTradingPaused),
		Entry("'95' - Works", "\"95\"", Financial_Quotes_NONE),
		Entry("'96' - Works", "\"96\"", Financial_Quotes_ShortSalesRestrictionActivated),
		Entry("'97' - Works", "\"97\"", Financial_Quotes_ShortSalesRestrictionContinued),
		Entry("'98' - Works", "\"98\"", Financial_Quotes_ShortSalesRestrictionDeactivated),
		Entry("'99' - Works", "\"99\"", Financial_Quotes_ShortSalesRestrictionInEffect),
		Entry("'100' - Works", "\"100\"", Financial_Quotes_ShortSalesRestrictionMax),
		Entry("'101' - Works", "\"101\"", Financial_Quotes_NBBONoChange),
		Entry("'102' - Works", "\"102\"", Financial_Quotes_NBBOQuoteIsNBBO),
		Entry("'103' - Works", "\"103\"", Financial_Quotes_NBBONoBBNoBO),
		Entry("'104' - Works", "\"104\"", Financial_Quotes_NBBOBBBOShortAppendage),
		Entry("'105' - Works", "\"105\"", Financial_Quotes_NBBOBBBOLongAppendage),
		Entry("'106' - Works", "\"106\"", Financial_Quotes_HeldTradeNotLastSaleNotConsolidated),
		Entry("'107' - Works", "\"107\"", Financial_Quotes_HeldTradeLastSaleButNotConsolidated),
		Entry("'108' - Works", "\"108\"", Financial_Quotes_HeldTradeLastSaleAndConsolidated),
		Entry("'109' - Works", "\"109\"", Financial_Quotes_RetailInterestOnBid),
		Entry("'110' - Works", "\"110\"", Financial_Quotes_RetailInterestOnAsk),
		Entry("'111' - Works", "\"111\"", Financial_Quotes_RetailInterestOnBidAndAsk),
		Entry("'112' - Works", "\"112\"", Financial_Quotes_FinraBBONoChange),
		Entry("'113' - Works", "\"113\"", Financial_Quotes_FinraBBODoesNotExist),
		Entry("'114' - Works", "\"114\"", Financial_Quotes_FinraBBBOExecutable),
		Entry("'115' - Works", "\"115\"", Financial_Quotes_FinraBBBelowLowerBand),
		Entry("'116' - Works", "\"116\"", Financial_Quotes_FinraBOAboveUpperBand),
		Entry("'117' - Works", "\"117\"", Financial_Quotes_FinraBBBelowLowerBandBOAbboveUpperBand),
		Entry("'118' - Works", "\"118\"", Financial_Quotes_CTANotDueToRelatedSecurity),
		Entry("'119' - Works", "\"119\"", Financial_Quotes_CTADueToRelatedSecurity),
		Entry("'120' - Works", "\"120\"", Financial_Quotes_CTANotInViewOfCommon),
		Entry("'121' - Works", "\"121\"", Financial_Quotes_CTAInViewOfCommon),
		Entry("'122' - Works", "\"122\"", Financial_Quotes_CTAPriceIndicator),
		Entry("'123' - Works", "\"123\"", Financial_Quotes_CTANewPriceIndicator),
		Entry("'124' - Works", "\"124\"", Financial_Quotes_CTACorrectedPriceIndication),
		Entry("'125' - Works", "\"125\"", Financial_Quotes_CTACancelledMarketImbalance),
		Entry("0 - Works", 0, Financial_Quotes_NBBNBOExecutable),
		Entry("1 - Works", 1, Financial_Quotes_NBBBelowLowerBand),
		Entry("2 - Works", 2, Financial_Quotes_NBOAboveUpperBand),
		Entry("3 - Works", 3, Financial_Quotes_NBBBelowLowerBandAndNBOAboveUpperBand),
		Entry("4 - Works", 4, Financial_Quotes_NBBEqualsUpperBand),
		Entry("5 - Works", 5, Financial_Quotes_NBOEqualsLowerBand),
		Entry("6 - Works", 6, Financial_Quotes_NBBEqualsUpperBandAndNBOAboveUpperBand),
		Entry("7 - Works", 7, Financial_Quotes_NBBBelowLowerBandAndNBOEqualsLowerBand),
		Entry("8 - Works", 8, Financial_Quotes_BidPriceAboveUpperLimitPriceBand),
		Entry("9 - Works", 9, Financial_Quotes_OfferPriceBelowLowerLimitPriceBand),
		Entry("10 - Works", 10, Financial_Quotes_BidAndOfferOutsidePriceBand),
		Entry("11 - Works", 11, Financial_Quotes_OpeningUpdate),
		Entry("12 - Works", 12, Financial_Quotes_IntraDayUpdate),
		Entry("13 - Works", 13, Financial_Quotes_RestatedValue),
		Entry("14 - Works", 14, Financial_Quotes_SuspendedDuringTradingHalt),
		Entry("15 - Works", 15, Financial_Quotes_ReOpeningUpdate),
		Entry("16 - Works", 16, Financial_Quotes_OutsidePriceBandRuleHours),
		Entry("17 - Works", 17, Financial_Quotes_AuctionExtension),
		Entry("18 - Works", 18, Financial_Quotes_LULDPriceBand),
		Entry("19 - Works", 19, Financial_Quotes_RepublishedLULDPriceBandInd),
		Entry("20 - Works", 20, Financial_Quotes_NBBLimitStateEntered),
		Entry("21 - Works", 21, Financial_Quotes_NBBLimitStateExited),
		Entry("22 - Works", 22, Financial_Quotes_NBOLimitStateEntered),
		Entry("23 - Works", 23, Financial_Quotes_NBOLimitStateExited),
		Entry("24 - Works", 24, Financial_Quotes_NBBAndNBOLimitStateEntered),
		Entry("25 - Works", 25, Financial_Quotes_NBBAndNBOLimitStateExited),
		Entry("26 - Works", 26, Financial_Quotes_NBBLimitStateEnteredNBOLimitStateExited),
		Entry("27 - Works", 27, Financial_Quotes_NBBLimitStateExitedNBOLimitStateEntered),
		Entry("28 - Works", 28, Financial_Quotes_Normal),
		Entry("29 - Works", 29, Financial_Quotes_Bankrupt),
		Entry("30 - Works", 30, Financial_Quotes_Deficient),
		Entry("31 - Works", 31, Financial_Quotes_Delinquent),
		Entry("32 - Works", 32, Financial_Quotes_BankruptAndDeficient),
		Entry("33 - Works", 33, Financial_Quotes_BankruptAndDelinquent),
		Entry("34 - Works", 34, Financial_Quotes_DeficientAndDelinquent),
		Entry("35 - Works", 35, Financial_Quotes_DeficientDeliquentBankrupt),
		Entry("36 - Works", 36, Financial_Quotes_Liquidation),
		Entry("37 - Works", 37, Financial_Quotes_CreationsSuspended),
		Entry("38 - Works", 38, Financial_Quotes_RedemptionsSuspended),
		Entry("39 - Works", 39, Financial_Quotes_CreationsRedemptionsSuspended),
		Entry("40 - Works", 40, Financial_Quotes_NormalTrading),
		Entry("41 - Works", 41, Financial_Quotes_OpeningDelay),
		Entry("42 - Works", 42, Financial_Quotes_TradingHalt),
		Entry("43 - Works", 43, Financial_Quotes_TradingResume),
		Entry("44 - Works", 44, Financial_Quotes_NoOpenNoResume),
		Entry("45 - Works", 45, Financial_Quotes_PriceIndication),
		Entry("46 - Works", 46, Financial_Quotes_TradingRangeIndication),
		Entry("47 - Works", 47, Financial_Quotes_MarketImbalanceBuy),
		Entry("48 - Works", 48, Financial_Quotes_MarketImbalanceSell),
		Entry("49 - Works", 49, Financial_Quotes_MarketOnCloseImbalanceBuy),
		Entry("50 - Works", 50, Financial_Quotes_MarketOnCloseImbalanceSell),
		Entry("51 - Works", 51, Financial_Quotes_NoMarketImbalance),
		Entry("52 - Works", 52, Financial_Quotes_NoMarketOnCloseImbalance),
		Entry("53 - Works", 53, Financial_Quotes_ShortSaleRestriction),
		Entry("54 - Works", 54, Financial_Quotes_LimitUpLimitDown),
		Entry("55 - Works", 55, Financial_Quotes_QuotationResumption),
		Entry("56 - Works", 56, Financial_Quotes_TradingResumption),
		Entry("57 - Works", 57, Financial_Quotes_VolatilityTradingPause),
		Entry("58 - Works", 58, Financial_Quotes_PolygonReserved),
		Entry("59 - Works", 59, Financial_Quotes_HaltNewsPending),
		Entry("60 - Works", 60, Financial_Quotes_UpdateNewsDissemination),
		Entry("61 - Works", 61, Financial_Quotes_HaltSingleStockTradingPause),
		Entry("62 - Works", 62, Financial_Quotes_HaltRegulatoryExtraordinaryMarketActivity),
		Entry("63 - Works", 63, Financial_Quotes_HaltETF),
		Entry("64 - Works", 64, Financial_Quotes_HaltInformationRequested),
		Entry("65 - Works", 65, Financial_Quotes_HaltExchangeNonCompliance),
		Entry("66 - Works", 66, Financial_Quotes_HaltFilingsNotCurrent),
		Entry("67 - Works", 67, Financial_Quotes_HaltSECTradingSuspension),
		Entry("68 - Works", 68, Financial_Quotes_HaltRegulatoryConcern),
		Entry("69 - Works", 69, Financial_Quotes_HaltMarketOperations),
		Entry("70 - Works", 70, Financial_Quotes_IPOSecurityNotYetTrading),
		Entry("71 - Works", 71, Financial_Quotes_HaltCorporateAction),
		Entry("72 - Works", 72, Financial_Quotes_QuotationNotAvailable),
		Entry("73 - Works", 73, Financial_Quotes_HaltVolatilityTradingPause),
		Entry("74 - Works", 74, Financial_Quotes_HaltVolatilityTradingPauseStraddleCondition),
		Entry("75 - Works", 75, Financial_Quotes_UpdateNewsAndResumptionTimes),
		Entry("76 - Works", 76, Financial_Quotes_HaltSingleStockTradingPauseQuotesOnly),
		Entry("77 - Works", 77, Financial_Quotes_ResumeQualificationIssuesReviewedResolved),
		Entry("78 - Works", 78, Financial_Quotes_ResumeFilingRequirementsSatisfiedResolved),
		Entry("79 - Works", 79, Financial_Quotes_ResumeNewsNotForthcoming),
		Entry("80 - Works", 80, Financial_Quotes_ResumeQualificationsMaintRequirementsMet),
		Entry("81 - Works", 81, Financial_Quotes_ResumeQualificationsFilingsMet),
		Entry("82 - Works", 82, Financial_Quotes_ResumeRegulatoryAuth),
		Entry("83 - Works", 83, Financial_Quotes_NewIssueAvailable),
		Entry("84 - Works", 84, Financial_Quotes_IssueAvailable),
		Entry("85 - Works", 85, Financial_Quotes_MWCBCarryFromPreviousDay),
		Entry("86 - Works", 86, Financial_Quotes_MWCBResume),
		Entry("87 - Works", 87, Financial_Quotes_IPOSecurityReleasedForQuotation),
		Entry("88 - Works", 88, Financial_Quotes_IPOSecurityPositioningWindowExtension),
		Entry("89 - Works", 89, Financial_Quotes_MWCBLevel1),
		Entry("90 - Works", 90, Financial_Quotes_MWCBLevel2),
		Entry("91 - Works", 91, Financial_Quotes_MWCBLevel3),
		Entry("92 - Works", 92, Financial_Quotes_HaltSubPennyTrading),
		Entry("93 - Works", 93, Financial_Quotes_OrderImbalanceInd),
		Entry("94 - Works", 94, Financial_Quotes_LULDTradingPaused),
		Entry("95 - Works", 95, Financial_Quotes_NONE),
		Entry("96 - Works", 96, Financial_Quotes_ShortSalesRestrictionActivated),
		Entry("97 - Works", 97, Financial_Quotes_ShortSalesRestrictionContinued),
		Entry("98 - Works", 98, Financial_Quotes_ShortSalesRestrictionDeactivated),
		Entry("99 - Works", 99, Financial_Quotes_ShortSalesRestrictionInEffect),
		Entry("100 - Works", 100, Financial_Quotes_ShortSalesRestrictionMax),
		Entry("101 - Works", 101, Financial_Quotes_NBBONoChange),
		Entry("102 - Works", 102, Financial_Quotes_NBBOQuoteIsNBBO),
		Entry("103 - Works", 103, Financial_Quotes_NBBONoBBNoBO),
		Entry("104 - Works", 104, Financial_Quotes_NBBOBBBOShortAppendage),
		Entry("105 - Works", 105, Financial_Quotes_NBBOBBBOLongAppendage),
		Entry("106 - Works", 106, Financial_Quotes_HeldTradeNotLastSaleNotConsolidated),
		Entry("107 - Works", 107, Financial_Quotes_HeldTradeLastSaleButNotConsolidated),
		Entry("108 - Works", 108, Financial_Quotes_HeldTradeLastSaleAndConsolidated),
		Entry("109 - Works", 109, Financial_Quotes_RetailInterestOnBid),
		Entry("110 - Works", 110, Financial_Quotes_RetailInterestOnAsk),
		Entry("111 - Works", 111, Financial_Quotes_RetailInterestOnBidAndAsk),
		Entry("112 - Works", 112, Financial_Quotes_FinraBBONoChange),
		Entry("113 - Works", 113, Financial_Quotes_FinraBBODoesNotExist),
		Entry("114 - Works", 114, Financial_Quotes_FinraBBBOExecutable),
		Entry("115 - Works", 115, Financial_Quotes_FinraBBBelowLowerBand),
		Entry("116 - Works", 116, Financial_Quotes_FinraBOAboveUpperBand),
		Entry("117 - Works", 117, Financial_Quotes_FinraBBBelowLowerBandBOAbboveUpperBand),
		Entry("118 - Works", 118, Financial_Quotes_CTANotDueToRelatedSecurity),
		Entry("119 - Works", 119, Financial_Quotes_CTADueToRelatedSecurity),
		Entry("120 - Works", 120, Financial_Quotes_CTANotInViewOfCommon),
		Entry("121 - Works", 121, Financial_Quotes_CTAInViewOfCommon),
		Entry("122 - Works", 122, Financial_Quotes_CTAPriceIndicator),
		Entry("123 - Works", 123, Financial_Quotes_CTANewPriceIndicator),
		Entry("124 - Works", 124, Financial_Quotes_CTACorrectedPriceIndication),
		Entry("125 - Works", 125, Financial_Quotes_CTACancelledMarketImbalance))

	// Test that attempting to deserialize a Financial.Quotes.Indicator will fial and return an
	// error if the value cannot be converted to either the name value or integer value
	// of the enum option
	It("UnmarshalCSV - Value is empty - Error", func() {

		// Attempt to convert a fake string value into a Financial.Quotes.Indicator
		// This should return an error
		enum := new(Financial_Quotes_Indicator)
		err := enum.UnmarshalCSV("")

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"\" cannot be mapped to a gopb.Financial_Quotes_Indicator"))
	})

	// Test the conditions under which values should be convertible to a Financial.Quotes.Indicator
	DescribeTable("UnmarshalCSV Tests",
		func(value string, shouldBe Financial_Quotes_Indicator) {

			// Attempt to convert the value into a Financial.Quotes.Indicator
			// This should not fail
			var enum Financial_Quotes_Indicator
			err := enum.UnmarshalCSV(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("NBBNBOExecutable - Works", "NBBNBOExecutable", Financial_Quotes_NBBNBOExecutable),
		Entry("NBBBelowLowerBand - Works", "NBBBelowLowerBand", Financial_Quotes_NBBBelowLowerBand),
		Entry("NBOAboveUpperBand - Works", "NBOAboveUpperBand", Financial_Quotes_NBOAboveUpperBand),
		Entry("NBBBelowLowerBandAndNBOAboveUpperBand - Works",
			"NBBBelowLowerBandAndNBOAboveUpperBand", Financial_Quotes_NBBBelowLowerBandAndNBOAboveUpperBand),
		Entry("NBBEqualsUpperBand - Works", "NBBEqualsUpperBand", Financial_Quotes_NBBEqualsUpperBand),
		Entry("NBOEqualsLowerBand - Works", "NBOEqualsLowerBand", Financial_Quotes_NBOEqualsLowerBand),
		Entry("NBBEqualsUpperBandAndNBOAboveUpperBand - Works",
			"NBBEqualsUpperBandAndNBOAboveUpperBand", Financial_Quotes_NBBEqualsUpperBandAndNBOAboveUpperBand),
		Entry("NBBBelowLowerBandAndNBOEqualsLowerBand - Works",
			"NBBBelowLowerBandAndNBOEqualsLowerBand", Financial_Quotes_NBBBelowLowerBandAndNBOEqualsLowerBand),
		Entry("BidPriceAboveUpperLimitPriceBand - Works",
			"BidPriceAboveUpperLimitPriceBand", Financial_Quotes_BidPriceAboveUpperLimitPriceBand),
		Entry("OfferPriceBelowLowerLimitPriceBand - Works",
			"OfferPriceBelowLowerLimitPriceBand", Financial_Quotes_OfferPriceBelowLowerLimitPriceBand),
		Entry("BidAndOfferOutsidePriceBand - Works",
			"BidAndOfferOutsidePriceBand", Financial_Quotes_BidAndOfferOutsidePriceBand),
		Entry("OpeningUpdate - Works", "OpeningUpdate", Financial_Quotes_OpeningUpdate),
		Entry("IntraDayUpdate - Works", "IntraDayUpdate", Financial_Quotes_IntraDayUpdate),
		Entry("RestatedValue - Works", "RestatedValue", Financial_Quotes_RestatedValue),
		Entry("SuspendedDuringTradingHalt - Works", "SuspendedDuringTradingHalt", Financial_Quotes_SuspendedDuringTradingHalt),
		Entry("ReOpeningUpdate - Works", "ReOpeningUpdate", Financial_Quotes_ReOpeningUpdate),
		Entry("OutsidePriceBandRuleHours - Works", "OutsidePriceBandRuleHours", Financial_Quotes_OutsidePriceBandRuleHours),
		Entry("AuctionExtension - Works", "AuctionExtension", Financial_Quotes_AuctionExtension),
		Entry("LULDPriceBand - Works", "LULDPriceBand", Financial_Quotes_LULDPriceBand),
		Entry("RepublishedLULDPriceBandInd - Works",
			"RepublishedLULDPriceBandInd", Financial_Quotes_RepublishedLULDPriceBandInd),
		Entry("NBBLimitStateEntered - Works", "NBBLimitStateEntered", Financial_Quotes_NBBLimitStateEntered),
		Entry("NBBLimitStateExited - Works", "NBBLimitStateExited", Financial_Quotes_NBBLimitStateExited),
		Entry("NBOLimitStateEntered - Works", "NBOLimitStateEntered", Financial_Quotes_NBOLimitStateEntered),
		Entry("NBOLimitStateExited - Works", "NBOLimitStateExited", Financial_Quotes_NBOLimitStateExited),
		Entry("NBBAndNBOLimitStateEntered - Works", "NBBAndNBOLimitStateEntered", Financial_Quotes_NBBAndNBOLimitStateEntered),
		Entry("NBBAndNBOLimitStateExited - Works", "NBBAndNBOLimitStateExited", Financial_Quotes_NBBAndNBOLimitStateExited),
		Entry("NBBLimitStateEnteredNBOLimitStateExited - Works",
			"NBBLimitStateEnteredNBOLimitStateExited", Financial_Quotes_NBBLimitStateEnteredNBOLimitStateExited),
		Entry("NBBLimitStateExitedNBOLimitStateEntered - Works",
			"NBBLimitStateExitedNBOLimitStateEntered", Financial_Quotes_NBBLimitStateExitedNBOLimitStateEntered),
		Entry("Normal - Works", "Normal", Financial_Quotes_Normal),
		Entry("Bankrupt - Works", "Bankrupt", Financial_Quotes_Bankrupt),
		Entry("Deficient - Works", "Deficient", Financial_Quotes_Deficient),
		Entry("Delinquent - Works", "Delinquent", Financial_Quotes_Delinquent),
		Entry("BankruptAndDeficient - Works", "BankruptAndDeficient", Financial_Quotes_BankruptAndDeficient),
		Entry("BankruptAndDelinquent - Works", "BankruptAndDelinquent", Financial_Quotes_BankruptAndDelinquent),
		Entry("DeficientAndDelinquent - Works", "DeficientAndDelinquent", Financial_Quotes_DeficientAndDelinquent),
		Entry("DeficientDeliquentBankrupt - Works", "DeficientDeliquentBankrupt", Financial_Quotes_DeficientDeliquentBankrupt),
		Entry("Liquidation - Works", "Liquidation", Financial_Quotes_Liquidation),
		Entry("CreationsSuspended - Works", "CreationsSuspended", Financial_Quotes_CreationsSuspended),
		Entry("RedemptionsSuspended - Works", "RedemptionsSuspended", Financial_Quotes_RedemptionsSuspended),
		Entry("CreationsRedemptionsSuspended - Works",
			"CreationsRedemptionsSuspended", Financial_Quotes_CreationsRedemptionsSuspended),
		Entry("NormalTrading - Works", "NormalTrading", Financial_Quotes_NormalTrading),
		Entry("OpeningDelay - Works", "OpeningDelay", Financial_Quotes_OpeningDelay),
		Entry("TradingHalt - Works", "TradingHalt", Financial_Quotes_TradingHalt),
		Entry("TradingResume - Works", "TradingResume", Financial_Quotes_TradingResume),
		Entry("NoOpenNoResume - Works", "NoOpenNoResume", Financial_Quotes_NoOpenNoResume),
		Entry("PriceIndication - Works", "PriceIndication", Financial_Quotes_PriceIndication),
		Entry("TradingRangeIndication - Works", "TradingRangeIndication", Financial_Quotes_TradingRangeIndication),
		Entry("MarketImbalanceBuy - Works", "MarketImbalanceBuy", Financial_Quotes_MarketImbalanceBuy),
		Entry("MarketImbalanceSell - Works", "MarketImbalanceSell", Financial_Quotes_MarketImbalanceSell),
		Entry("MarketOnCloseImbalanceBuy - Works", "MarketOnCloseImbalanceBuy", Financial_Quotes_MarketOnCloseImbalanceBuy),
		Entry("MarketOnCloseImbalanceSell - Works", "MarketOnCloseImbalanceSell", Financial_Quotes_MarketOnCloseImbalanceSell),
		Entry("NoMarketImbalance - Works", "NoMarketImbalance", Financial_Quotes_NoMarketImbalance),
		Entry("NoMarketOnCloseImbalance - Works", "NoMarketOnCloseImbalance", Financial_Quotes_NoMarketOnCloseImbalance),
		Entry("ShortSaleRestriction - Works", "ShortSaleRestriction", Financial_Quotes_ShortSaleRestriction),
		Entry("LimitUpLimitDown - Works", "LimitUpLimitDown", Financial_Quotes_LimitUpLimitDown),
		Entry("QuotationResumption - Works", "QuotationResumption", Financial_Quotes_QuotationResumption),
		Entry("TradingResumption - Works", "TradingResumption", Financial_Quotes_TradingResumption),
		Entry("VolatilityTradingPause - Works", "VolatilityTradingPause", Financial_Quotes_VolatilityTradingPause),
		Entry("PolygonReserved - Works", "PolygonReserved", Financial_Quotes_PolygonReserved),
		Entry("HaltNewsPending - Works", "HaltNewsPending", Financial_Quotes_HaltNewsPending),
		Entry("UpdateNewsDissemination - Works", "UpdateNewsDissemination", Financial_Quotes_UpdateNewsDissemination),
		Entry("HaltSingleStockTradingPause - Works",
			"HaltSingleStockTradingPause", Financial_Quotes_HaltSingleStockTradingPause),
		Entry("HaltRegulatoryExtraordinaryMarketActivity - Works",
			"HaltRegulatoryExtraordinaryMarketActivity", Financial_Quotes_HaltRegulatoryExtraordinaryMarketActivity),
		Entry("HaltETF - Works", "HaltETF", Financial_Quotes_HaltETF),
		Entry("HaltInformationRequested - Works", "HaltInformationRequested", Financial_Quotes_HaltInformationRequested),
		Entry("HaltExchangeNonCompliance - Works", "HaltExchangeNonCompliance", Financial_Quotes_HaltExchangeNonCompliance),
		Entry("HaltFilingsNotCurrent - Works", "HaltFilingsNotCurrent", Financial_Quotes_HaltFilingsNotCurrent),
		Entry("HaltSECTradingSuspension - Works", "HaltSECTradingSuspension", Financial_Quotes_HaltSECTradingSuspension),
		Entry("HaltRegulatoryConcern - Works", "HaltRegulatoryConcern", Financial_Quotes_HaltRegulatoryConcern),
		Entry("HaltMarketOperations - Works", "HaltMarketOperations", Financial_Quotes_HaltMarketOperations),
		Entry("IPOSecurityNotYetTrading - Works", "IPOSecurityNotYetTrading", Financial_Quotes_IPOSecurityNotYetTrading),
		Entry("HaltCorporateAction - Works", "HaltCorporateAction", Financial_Quotes_HaltCorporateAction),
		Entry("QuotationNotAvailable - Works", "QuotationNotAvailable", Financial_Quotes_QuotationNotAvailable),
		Entry("HaltVolatilityTradingPause - Works", "HaltVolatilityTradingPause", Financial_Quotes_HaltVolatilityTradingPause),
		Entry("HaltVolatilityTradingPauseStraddleCondition - Works",
			"HaltVolatilityTradingPauseStraddleCondition", Financial_Quotes_HaltVolatilityTradingPauseStraddleCondition),
		Entry("UpdateNewsAndResumptionTimes - Works",
			"UpdateNewsAndResumptionTimes", Financial_Quotes_UpdateNewsAndResumptionTimes),
		Entry("HaltSingleStockTradingPauseQuotesOnly - Works",
			"HaltSingleStockTradingPauseQuotesOnly", Financial_Quotes_HaltSingleStockTradingPauseQuotesOnly),
		Entry("ResumeQualificationIssuesReviewedResolved - Works",
			"ResumeQualificationIssuesReviewedResolved", Financial_Quotes_ResumeQualificationIssuesReviewedResolved),
		Entry("ResumeFilingRequirementsSatisfiedResolved - Works",
			"ResumeFilingRequirementsSatisfiedResolved", Financial_Quotes_ResumeFilingRequirementsSatisfiedResolved),
		Entry("ResumeNewsNotForthcoming - Works", "ResumeNewsNotForthcoming", Financial_Quotes_ResumeNewsNotForthcoming),
		Entry("ResumeQualificationsMaintRequirementsMet - Works",
			"ResumeQualificationsMaintRequirementsMet", Financial_Quotes_ResumeQualificationsMaintRequirementsMet),
		Entry("ResumeQualificationsFilingsMet - Works",
			"ResumeQualificationsFilingsMet", Financial_Quotes_ResumeQualificationsFilingsMet),
		Entry("ResumeRegulatoryAuth - Works", "ResumeRegulatoryAuth", Financial_Quotes_ResumeRegulatoryAuth),
		Entry("NewIssueAvailable - Works", "NewIssueAvailable", Financial_Quotes_NewIssueAvailable),
		Entry("IssueAvailable - Works", "IssueAvailable", Financial_Quotes_IssueAvailable),
		Entry("MWCBCarryFromPreviousDay - Works", "MWCBCarryFromPreviousDay", Financial_Quotes_MWCBCarryFromPreviousDay),
		Entry("MWCBResume - Works", "MWCBResume", Financial_Quotes_MWCBResume),
		Entry("IPOSecurityReleasedForQuotation - Works",
			"IPOSecurityReleasedForQuotation", Financial_Quotes_IPOSecurityReleasedForQuotation),
		Entry("IPOSecurityPositioningWindowExtension - Works",
			"IPOSecurityPositioningWindowExtension", Financial_Quotes_IPOSecurityPositioningWindowExtension),
		Entry("MWCBLevel1 - Works", "MWCBLevel1", Financial_Quotes_MWCBLevel1),
		Entry("MWCBLevel2 - Works", "MWCBLevel2", Financial_Quotes_MWCBLevel2),
		Entry("MWCBLevel3 - Works", "MWCBLevel3", Financial_Quotes_MWCBLevel3),
		Entry("HaltSubPennyTrading - Works", "HaltSubPennyTrading", Financial_Quotes_HaltSubPennyTrading),
		Entry("OrderImbalanceInd - Works", "OrderImbalanceInd", Financial_Quotes_OrderImbalanceInd),
		Entry("LULDTradingPaused - Works", "LULDTradingPaused", Financial_Quotes_LULDTradingPaused),
		Entry("NONE - Works", "NONE", Financial_Quotes_NONE),
		Entry("ShortSalesRestrictionActivated - Works",
			"ShortSalesRestrictionActivated", Financial_Quotes_ShortSalesRestrictionActivated),
		Entry("ShortSalesRestrictionContinued - Works",
			"ShortSalesRestrictionContinued", Financial_Quotes_ShortSalesRestrictionContinued),
		Entry("ShortSalesRestrictionDeactivated - Works",
			"ShortSalesRestrictionDeactivated", Financial_Quotes_ShortSalesRestrictionDeactivated),
		Entry("ShortSalesRestrictionInEffect - Works",
			"ShortSalesRestrictionInEffect", Financial_Quotes_ShortSalesRestrictionInEffect),
		Entry("ShortSalesRestrictionMax - Works", "ShortSalesRestrictionMax", Financial_Quotes_ShortSalesRestrictionMax),
		Entry("RetailInterestOnBid - Works", "RetailInterestOnBid", Financial_Quotes_RetailInterestOnBid),
		Entry("RetailInterestOnAsk - Works", "RetailInterestOnAsk", Financial_Quotes_RetailInterestOnAsk),
		Entry("RetailInterestOnBidAndAsk - Works",
			"RetailInterestOnBidAndAsk", Financial_Quotes_RetailInterestOnBidAndAsk),
		Entry("FinraBBONoChange - Works", "FinraBBONoChange", Financial_Quotes_FinraBBONoChange),
		Entry("FinraBBODoesNotExist - Works", "FinraBBODoesNotExist", Financial_Quotes_FinraBBODoesNotExist),
		Entry("FinraBBBOExecutable - Works", "FinraBBBOExecutable", Financial_Quotes_FinraBBBOExecutable),
		Entry("FinraBBBelowLowerBand - Works", "FinraBBBelowLowerBand", Financial_Quotes_FinraBBBelowLowerBand),
		Entry("FinraBOAboveUpperBand - Works", "FinraBOAboveUpperBand", Financial_Quotes_FinraBOAboveUpperBand),
		Entry("FinraBBBelowLowerBandBOAbboveUpperBand - Works",
			"FinraBBBelowLowerBandBOAbboveUpperBand", Financial_Quotes_FinraBBBelowLowerBandBOAbboveUpperBand),
		Entry("NBBONoChange - Works", "NBBONoChange", Financial_Quotes_NBBONoChange),
		Entry("NBBOQuoteIsNBBO - Works", "NBBOQuoteIsNBBO", Financial_Quotes_NBBOQuoteIsNBBO),
		Entry("NBBONoBBNoBO - Works", "NBBONoBBNoBO", Financial_Quotes_NBBONoBBNoBO),
		Entry("NBBOBBBOShortAppendage - Works", "NBBOBBBOShortAppendage", Financial_Quotes_NBBOBBBOShortAppendage),
		Entry("NBBOBBBOLongAppendage - Works", "NBBOBBBOLongAppendage", Financial_Quotes_NBBOBBBOLongAppendage),
		Entry("HeldTradeNotLastSaleNotConsolidated - Works",
			"HeldTradeNotLastSaleNotConsolidated", Financial_Quotes_HeldTradeNotLastSaleNotConsolidated),
		Entry("HeldTradeLastSaleButNotConsolidated - Works",
			"HeldTradeLastSaleButNotConsolidated", Financial_Quotes_HeldTradeLastSaleButNotConsolidated),
		Entry("HeldTradeLastSaleAndConsolidated - Works",
			"HeldTradeLastSaleAndConsolidated", Financial_Quotes_HeldTradeLastSaleAndConsolidated),
		Entry("CTANotDueToRelatedSecurity - Works", "CTANotDueToRelatedSecurity", Financial_Quotes_CTANotDueToRelatedSecurity),
		Entry("CTADueToRelatedSecurity - Works", "CTADueToRelatedSecurity", Financial_Quotes_CTADueToRelatedSecurity),
		Entry("CTANotInViewOfCommon - Works", "CTANotInViewOfCommon", Financial_Quotes_CTANotInViewOfCommon),
		Entry("CTAInViewOfCommon - Works", "CTAInViewOfCommon", Financial_Quotes_CTAInViewOfCommon),
		Entry("CTAPriceIndicator - Works", "CTAPriceIndicator", Financial_Quotes_CTAPriceIndicator),
		Entry("CTANewPriceIndicator - Works", "CTANewPriceIndicator", Financial_Quotes_CTANewPriceIndicator),
		Entry("CTACorrectedPriceIndication - Works",
			"CTACorrectedPriceIndication", Financial_Quotes_CTACorrectedPriceIndication),
		Entry("CTACancelledMarketImbalance - Works",
			"CTACancelledMarketImbalance", Financial_Quotes_CTACancelledMarketImbalance),
		Entry("0 - Works", "0", Financial_Quotes_NBBNBOExecutable),
		Entry("1 - Works", "1", Financial_Quotes_NBBBelowLowerBand),
		Entry("2 - Works", "2", Financial_Quotes_NBOAboveUpperBand),
		Entry("3 - Works", "3", Financial_Quotes_NBBBelowLowerBandAndNBOAboveUpperBand),
		Entry("4 - Works", "4", Financial_Quotes_NBBEqualsUpperBand),
		Entry("5 - Works", "5", Financial_Quotes_NBOEqualsLowerBand),
		Entry("6 - Works", "6", Financial_Quotes_NBBEqualsUpperBandAndNBOAboveUpperBand),
		Entry("7 - Works", "7", Financial_Quotes_NBBBelowLowerBandAndNBOEqualsLowerBand),
		Entry("8 - Works", "8", Financial_Quotes_BidPriceAboveUpperLimitPriceBand),
		Entry("9 - Works", "9", Financial_Quotes_OfferPriceBelowLowerLimitPriceBand),
		Entry("10 - Works", "10", Financial_Quotes_BidAndOfferOutsidePriceBand),
		Entry("11 - Works", "11", Financial_Quotes_OpeningUpdate),
		Entry("12 - Works", "12", Financial_Quotes_IntraDayUpdate),
		Entry("13 - Works", "13", Financial_Quotes_RestatedValue),
		Entry("14 - Works", "14", Financial_Quotes_SuspendedDuringTradingHalt),
		Entry("15 - Works", "15", Financial_Quotes_ReOpeningUpdate),
		Entry("16 - Works", "16", Financial_Quotes_OutsidePriceBandRuleHours),
		Entry("17 - Works", "17", Financial_Quotes_AuctionExtension),
		Entry("18 - Works", "18", Financial_Quotes_LULDPriceBand),
		Entry("19 - Works", "19", Financial_Quotes_RepublishedLULDPriceBandInd),
		Entry("20 - Works", "20", Financial_Quotes_NBBLimitStateEntered),
		Entry("21 - Works", "21", Financial_Quotes_NBBLimitStateExited),
		Entry("22 - Works", "22", Financial_Quotes_NBOLimitStateEntered),
		Entry("23 - Works", "23", Financial_Quotes_NBOLimitStateExited),
		Entry("24 - Works", "24", Financial_Quotes_NBBAndNBOLimitStateEntered),
		Entry("25 - Works", "25", Financial_Quotes_NBBAndNBOLimitStateExited),
		Entry("26 - Works", "26", Financial_Quotes_NBBLimitStateEnteredNBOLimitStateExited),
		Entry("27 - Works", "27", Financial_Quotes_NBBLimitStateExitedNBOLimitStateEntered),
		Entry("28 - Works", "28", Financial_Quotes_Normal),
		Entry("29 - Works", "29", Financial_Quotes_Bankrupt),
		Entry("30 - Works", "30", Financial_Quotes_Deficient),
		Entry("31 - Works", "31", Financial_Quotes_Delinquent),
		Entry("32 - Works", "32", Financial_Quotes_BankruptAndDeficient),
		Entry("33 - Works", "33", Financial_Quotes_BankruptAndDelinquent),
		Entry("34 - Works", "34", Financial_Quotes_DeficientAndDelinquent),
		Entry("35 - Works", "35", Financial_Quotes_DeficientDeliquentBankrupt),
		Entry("36 - Works", "36", Financial_Quotes_Liquidation),
		Entry("37 - Works", "37", Financial_Quotes_CreationsSuspended),
		Entry("38 - Works", "38", Financial_Quotes_RedemptionsSuspended),
		Entry("39 - Works", "39", Financial_Quotes_CreationsRedemptionsSuspended),
		Entry("40 - Works", "40", Financial_Quotes_NormalTrading),
		Entry("41 - Works", "41", Financial_Quotes_OpeningDelay),
		Entry("42 - Works", "42", Financial_Quotes_TradingHalt),
		Entry("43 - Works", "43", Financial_Quotes_TradingResume),
		Entry("44 - Works", "44", Financial_Quotes_NoOpenNoResume),
		Entry("45 - Works", "45", Financial_Quotes_PriceIndication),
		Entry("46 - Works", "46", Financial_Quotes_TradingRangeIndication),
		Entry("47 - Works", "47", Financial_Quotes_MarketImbalanceBuy),
		Entry("48 - Works", "48", Financial_Quotes_MarketImbalanceSell),
		Entry("49 - Works", "49", Financial_Quotes_MarketOnCloseImbalanceBuy),
		Entry("50 - Works", "50", Financial_Quotes_MarketOnCloseImbalanceSell),
		Entry("51 - Works", "51", Financial_Quotes_NoMarketImbalance),
		Entry("52 - Works", "52", Financial_Quotes_NoMarketOnCloseImbalance),
		Entry("53 - Works", "53", Financial_Quotes_ShortSaleRestriction),
		Entry("54 - Works", "54", Financial_Quotes_LimitUpLimitDown),
		Entry("55 - Works", "55", Financial_Quotes_QuotationResumption),
		Entry("56 - Works", "56", Financial_Quotes_TradingResumption),
		Entry("57 - Works", "57", Financial_Quotes_VolatilityTradingPause),
		Entry("58 - Works", "58", Financial_Quotes_PolygonReserved),
		Entry("59 - Works", "59", Financial_Quotes_HaltNewsPending),
		Entry("60 - Works", "60", Financial_Quotes_UpdateNewsDissemination),
		Entry("61 - Works", "61", Financial_Quotes_HaltSingleStockTradingPause),
		Entry("62 - Works", "62", Financial_Quotes_HaltRegulatoryExtraordinaryMarketActivity),
		Entry("63 - Works", "63", Financial_Quotes_HaltETF),
		Entry("64 - Works", "64", Financial_Quotes_HaltInformationRequested),
		Entry("65 - Works", "65", Financial_Quotes_HaltExchangeNonCompliance),
		Entry("66 - Works", "66", Financial_Quotes_HaltFilingsNotCurrent),
		Entry("67 - Works", "67", Financial_Quotes_HaltSECTradingSuspension),
		Entry("68 - Works", "68", Financial_Quotes_HaltRegulatoryConcern),
		Entry("69 - Works", "69", Financial_Quotes_HaltMarketOperations),
		Entry("70 - Works", "70", Financial_Quotes_IPOSecurityNotYetTrading),
		Entry("71 - Works", "71", Financial_Quotes_HaltCorporateAction),
		Entry("72 - Works", "72", Financial_Quotes_QuotationNotAvailable),
		Entry("73 - Works", "73", Financial_Quotes_HaltVolatilityTradingPause),
		Entry("74 - Works", "74", Financial_Quotes_HaltVolatilityTradingPauseStraddleCondition),
		Entry("75 - Works", "75", Financial_Quotes_UpdateNewsAndResumptionTimes),
		Entry("76 - Works", "76", Financial_Quotes_HaltSingleStockTradingPauseQuotesOnly),
		Entry("77 - Works", "77", Financial_Quotes_ResumeQualificationIssuesReviewedResolved),
		Entry("78 - Works", "78", Financial_Quotes_ResumeFilingRequirementsSatisfiedResolved),
		Entry("79 - Works", "79", Financial_Quotes_ResumeNewsNotForthcoming),
		Entry("80 - Works", "80", Financial_Quotes_ResumeQualificationsMaintRequirementsMet),
		Entry("81 - Works", "81", Financial_Quotes_ResumeQualificationsFilingsMet),
		Entry("82 - Works", "82", Financial_Quotes_ResumeRegulatoryAuth),
		Entry("83 - Works", "83", Financial_Quotes_NewIssueAvailable),
		Entry("84 - Works", "84", Financial_Quotes_IssueAvailable),
		Entry("85 - Works", "85", Financial_Quotes_MWCBCarryFromPreviousDay),
		Entry("86 - Works", "86", Financial_Quotes_MWCBResume),
		Entry("87 - Works", "87", Financial_Quotes_IPOSecurityReleasedForQuotation),
		Entry("88 - Works", "88", Financial_Quotes_IPOSecurityPositioningWindowExtension),
		Entry("89 - Works", "89", Financial_Quotes_MWCBLevel1),
		Entry("90 - Works", "90", Financial_Quotes_MWCBLevel2),
		Entry("91 - Works", "91", Financial_Quotes_MWCBLevel3),
		Entry("92 - Works", "92", Financial_Quotes_HaltSubPennyTrading),
		Entry("93 - Works", "93", Financial_Quotes_OrderImbalanceInd),
		Entry("94 - Works", "94", Financial_Quotes_LULDTradingPaused),
		Entry("95 - Works", "95", Financial_Quotes_NONE),
		Entry("96 - Works", "96", Financial_Quotes_ShortSalesRestrictionActivated),
		Entry("97 - Works", "97", Financial_Quotes_ShortSalesRestrictionContinued),
		Entry("98 - Works", "98", Financial_Quotes_ShortSalesRestrictionDeactivated),
		Entry("99 - Works", "99", Financial_Quotes_ShortSalesRestrictionInEffect),
		Entry("100 - Works", "100", Financial_Quotes_ShortSalesRestrictionMax),
		Entry("101 - Works", "101", Financial_Quotes_NBBONoChange),
		Entry("102 - Works", "102", Financial_Quotes_NBBOQuoteIsNBBO),
		Entry("103 - Works", "103", Financial_Quotes_NBBONoBBNoBO),
		Entry("104 - Works", "104", Financial_Quotes_NBBOBBBOShortAppendage),
		Entry("105 - Works", "105", Financial_Quotes_NBBOBBBOLongAppendage),
		Entry("106 - Works", "106", Financial_Quotes_HeldTradeNotLastSaleNotConsolidated),
		Entry("107 - Works", "107", Financial_Quotes_HeldTradeLastSaleButNotConsolidated),
		Entry("108 - Works", "108", Financial_Quotes_HeldTradeLastSaleAndConsolidated),
		Entry("109 - Works", "109", Financial_Quotes_RetailInterestOnBid),
		Entry("110 - Works", "110", Financial_Quotes_RetailInterestOnAsk),
		Entry("111 - Works", "111", Financial_Quotes_RetailInterestOnBidAndAsk),
		Entry("112 - Works", "112", Financial_Quotes_FinraBBONoChange),
		Entry("113 - Works", "113", Financial_Quotes_FinraBBODoesNotExist),
		Entry("114 - Works", "114", Financial_Quotes_FinraBBBOExecutable),
		Entry("115 - Works", "115", Financial_Quotes_FinraBBBelowLowerBand),
		Entry("116 - Works", "116", Financial_Quotes_FinraBOAboveUpperBand),
		Entry("117 - Works", "117", Financial_Quotes_FinraBBBelowLowerBandBOAbboveUpperBand),
		Entry("118 - Works", "118", Financial_Quotes_CTANotDueToRelatedSecurity),
		Entry("119 - Works", "119", Financial_Quotes_CTADueToRelatedSecurity),
		Entry("120 - Works", "120", Financial_Quotes_CTANotInViewOfCommon),
		Entry("121 - Works", "121", Financial_Quotes_CTAInViewOfCommon),
		Entry("122 - Works", "122", Financial_Quotes_CTAPriceIndicator),
		Entry("123 - Works", "123", Financial_Quotes_CTANewPriceIndicator),
		Entry("124 - Works", "124", Financial_Quotes_CTACorrectedPriceIndication),
		Entry("125 - Works", "125", Financial_Quotes_CTACancelledMarketImbalance))

	// Tests that, if the attribute type submitted to UnmarshalDynamoDBAttributeValue is not one we
	// recognize, then the function will return an error
	It("UnmarshalDynamoDBAttributeValue - AttributeValue type invalid - Error", func() {
		value := new(Financial_Quotes_Indicator)
		err := attributevalue.Unmarshal(&types.AttributeValueMemberBOOL{Value: true}, &value)
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("Attribute value of *types.AttributeValueMemberBOOL could not be converted to a Financial.Quotes.Indicator"))
	})

	// Tests the conditions under which UnmarshalDynamoDBAttributeValue is called and no error is generated
	DescribeTable("UnmarshalDynamoDBAttributeValue - AttributeValue Conditions",
		func(raw types.AttributeValue, expected Financial_Quotes_Indicator) {
			var value Financial_Quotes_Indicator
			err := attributevalue.Unmarshal(raw, &value)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(expected))
		},
		Entry("Value is []bytes, NBBNBOExecutable - Works",
			&types.AttributeValueMemberB{Value: []byte("NBBNBOExecutable")}, Financial_Quotes_NBBNBOExecutable),
		Entry("Value is []bytes, NBBBelowLowerBand - Works",
			&types.AttributeValueMemberB{Value: []byte("NBBBelowLowerBand")}, Financial_Quotes_NBBBelowLowerBand),
		Entry("Value is []bytes, NBOAboveUpperBand - Works",
			&types.AttributeValueMemberB{Value: []byte("NBOAboveUpperBand")}, Financial_Quotes_NBOAboveUpperBand),
		Entry("Value is []bytes, NBBBelowLowerBandAndNBOAboveUpperBand - Works",
			&types.AttributeValueMemberB{Value: []byte("NBBBelowLowerBandAndNBOAboveUpperBand")}, Financial_Quotes_NBBBelowLowerBandAndNBOAboveUpperBand),
		Entry("Value is []bytes, NBBEqualsUpperBand - Works",
			&types.AttributeValueMemberB{Value: []byte("NBBEqualsUpperBand")}, Financial_Quotes_NBBEqualsUpperBand),
		Entry("Value is []bytes, NBOEqualsLowerBand - Works",
			&types.AttributeValueMemberB{Value: []byte("NBOEqualsLowerBand")}, Financial_Quotes_NBOEqualsLowerBand),
		Entry("Value is []bytes, NBBEqualsUpperBandAndNBOAboveUpperBand - Works",
			&types.AttributeValueMemberB{Value: []byte("NBBEqualsUpperBandAndNBOAboveUpperBand")}, Financial_Quotes_NBBEqualsUpperBandAndNBOAboveUpperBand),
		Entry("Value is []bytes, NBBBelowLowerBandAndNBOEqualsLowerBand - Works",
			&types.AttributeValueMemberB{Value: []byte("NBBBelowLowerBandAndNBOEqualsLowerBand")}, Financial_Quotes_NBBBelowLowerBandAndNBOEqualsLowerBand),
		Entry("Value is []bytes, BidPriceAboveUpperLimitPriceBand - Works",
			&types.AttributeValueMemberB{Value: []byte("BidPriceAboveUpperLimitPriceBand")}, Financial_Quotes_BidPriceAboveUpperLimitPriceBand),
		Entry("Value is []bytes, OfferPriceBelowLowerLimitPriceBand - Works",
			&types.AttributeValueMemberB{Value: []byte("OfferPriceBelowLowerLimitPriceBand")}, Financial_Quotes_OfferPriceBelowLowerLimitPriceBand),
		Entry("Value is []bytes, BidAndOfferOutsidePriceBand - Works",
			&types.AttributeValueMemberB{Value: []byte("BidAndOfferOutsidePriceBand")}, Financial_Quotes_BidAndOfferOutsidePriceBand),
		Entry("Value is []bytes, OpeningUpdate - Works",
			&types.AttributeValueMemberB{Value: []byte("OpeningUpdate")}, Financial_Quotes_OpeningUpdate),
		Entry("Value is []bytes, IntraDayUpdate - Works",
			&types.AttributeValueMemberB{Value: []byte("IntraDayUpdate")}, Financial_Quotes_IntraDayUpdate),
		Entry("Value is []bytes, RestatedValue - Works",
			&types.AttributeValueMemberB{Value: []byte("RestatedValue")}, Financial_Quotes_RestatedValue),
		Entry("Value is []bytes, SuspendedDuringTradingHalt - Works",
			&types.AttributeValueMemberB{Value: []byte("SuspendedDuringTradingHalt")}, Financial_Quotes_SuspendedDuringTradingHalt),
		Entry("Value is []bytes, ReOpeningUpdate - Works",
			&types.AttributeValueMemberB{Value: []byte("ReOpeningUpdate")}, Financial_Quotes_ReOpeningUpdate),
		Entry("Value is []bytes, OutsidePriceBandRuleHours - Works",
			&types.AttributeValueMemberB{Value: []byte("OutsidePriceBandRuleHours")}, Financial_Quotes_OutsidePriceBandRuleHours),
		Entry("Value is []bytes, AuctionExtension - Works",
			&types.AttributeValueMemberB{Value: []byte("AuctionExtension")}, Financial_Quotes_AuctionExtension),
		Entry("Value is []bytes, LULDPriceBand - Works",
			&types.AttributeValueMemberB{Value: []byte("LULDPriceBand")}, Financial_Quotes_LULDPriceBand),
		Entry("Value is []bytes, RepublishedLULDPriceBandInd - Works",
			&types.AttributeValueMemberB{Value: []byte("RepublishedLULDPriceBandInd")}, Financial_Quotes_RepublishedLULDPriceBandInd),
		Entry("Value is []bytes, NBBLimitStateEntered - Works",
			&types.AttributeValueMemberB{Value: []byte("NBBLimitStateEntered")}, Financial_Quotes_NBBLimitStateEntered),
		Entry("Value is []bytes, NBBLimitStateExited - Works",
			&types.AttributeValueMemberB{Value: []byte("NBBLimitStateExited")}, Financial_Quotes_NBBLimitStateExited),
		Entry("Value is []bytes, NBOLimitStateEntered - Works",
			&types.AttributeValueMemberB{Value: []byte("NBOLimitStateEntered")}, Financial_Quotes_NBOLimitStateEntered),
		Entry("Value is []bytes, NBOLimitStateExited - Works",
			&types.AttributeValueMemberB{Value: []byte("NBOLimitStateExited")}, Financial_Quotes_NBOLimitStateExited),
		Entry("Value is []bytes, NBBAndNBOLimitStateEntered - Works",
			&types.AttributeValueMemberB{Value: []byte("NBBAndNBOLimitStateEntered")}, Financial_Quotes_NBBAndNBOLimitStateEntered),
		Entry("Value is []bytes, NBBAndNBOLimitStateExited - Works",
			&types.AttributeValueMemberB{Value: []byte("NBBAndNBOLimitStateExited")}, Financial_Quotes_NBBAndNBOLimitStateExited),
		Entry("Value is []bytes, NBBLimitStateEnteredNBOLimitStateExited - Works",
			&types.AttributeValueMemberB{Value: []byte("NBBLimitStateEnteredNBOLimitStateExited")}, Financial_Quotes_NBBLimitStateEnteredNBOLimitStateExited),
		Entry("Value is []bytes, NBBLimitStateExitedNBOLimitStateEntered - Works",
			&types.AttributeValueMemberB{Value: []byte("NBBLimitStateExitedNBOLimitStateEntered")}, Financial_Quotes_NBBLimitStateExitedNBOLimitStateEntered),
		Entry("Value is []bytes, Normal - Works",
			&types.AttributeValueMemberB{Value: []byte("Normal")}, Financial_Quotes_Normal),
		Entry("Value is []bytes, Bankrupt - Works",
			&types.AttributeValueMemberB{Value: []byte("Bankrupt")}, Financial_Quotes_Bankrupt),
		Entry("Value is []bytes, Deficient - Works",
			&types.AttributeValueMemberB{Value: []byte("Deficient")}, Financial_Quotes_Deficient),
		Entry("Value is []bytes, Delinquent - Works",
			&types.AttributeValueMemberB{Value: []byte("Delinquent")}, Financial_Quotes_Delinquent),
		Entry("Value is []bytes, BankruptAndDeficient - Works",
			&types.AttributeValueMemberB{Value: []byte("BankruptAndDeficient")}, Financial_Quotes_BankruptAndDeficient),
		Entry("Value is []bytes, BankruptAndDelinquent - Works",
			&types.AttributeValueMemberB{Value: []byte("BankruptAndDelinquent")}, Financial_Quotes_BankruptAndDelinquent),
		Entry("Value is []bytes, DeficientAndDelinquent - Works",
			&types.AttributeValueMemberB{Value: []byte("DeficientAndDelinquent")}, Financial_Quotes_DeficientAndDelinquent),
		Entry("Value is []bytes, DeficientDeliquentBankrupt - Works",
			&types.AttributeValueMemberB{Value: []byte("DeficientDeliquentBankrupt")}, Financial_Quotes_DeficientDeliquentBankrupt),
		Entry("Value is []bytes, Liquidation - Works",
			&types.AttributeValueMemberB{Value: []byte("Liquidation")}, Financial_Quotes_Liquidation),
		Entry("Value is []bytes, CreationsSuspended - Works",
			&types.AttributeValueMemberB{Value: []byte("CreationsSuspended")}, Financial_Quotes_CreationsSuspended),
		Entry("Value is []bytes, RedemptionsSuspended - Works",
			&types.AttributeValueMemberB{Value: []byte("RedemptionsSuspended")}, Financial_Quotes_RedemptionsSuspended),
		Entry("Value is []bytes, CreationsRedemptionsSuspended - Works",
			&types.AttributeValueMemberB{Value: []byte("CreationsRedemptionsSuspended")}, Financial_Quotes_CreationsRedemptionsSuspended),
		Entry("Value is []bytes, NormalTrading - Works",
			&types.AttributeValueMemberB{Value: []byte("NormalTrading")}, Financial_Quotes_NormalTrading),
		Entry("Value is []bytes, OpeningDelay - Works",
			&types.AttributeValueMemberB{Value: []byte("OpeningDelay")}, Financial_Quotes_OpeningDelay),
		Entry("Value is []bytes, TradingHalt - Works",
			&types.AttributeValueMemberB{Value: []byte("TradingHalt")}, Financial_Quotes_TradingHalt),
		Entry("Value is []bytes, TradingResume - Works",
			&types.AttributeValueMemberB{Value: []byte("TradingResume")}, Financial_Quotes_TradingResume),
		Entry("Value is []bytes, NoOpenNoResume - Works",
			&types.AttributeValueMemberB{Value: []byte("NoOpenNoResume")}, Financial_Quotes_NoOpenNoResume),
		Entry("Value is []bytes, PriceIndication - Works",
			&types.AttributeValueMemberB{Value: []byte("PriceIndication")}, Financial_Quotes_PriceIndication),
		Entry("Value is []bytes, TradingRangeIndication - Works",
			&types.AttributeValueMemberB{Value: []byte("TradingRangeIndication")}, Financial_Quotes_TradingRangeIndication),
		Entry("Value is []bytes, MarketImbalanceBuy - Works",
			&types.AttributeValueMemberB{Value: []byte("MarketImbalanceBuy")}, Financial_Quotes_MarketImbalanceBuy),
		Entry("Value is []bytes, MarketImbalanceSell - Works",
			&types.AttributeValueMemberB{Value: []byte("MarketImbalanceSell")}, Financial_Quotes_MarketImbalanceSell),
		Entry("Value is []bytes, MarketOnCloseImbalanceBuy - Works",
			&types.AttributeValueMemberB{Value: []byte("MarketOnCloseImbalanceBuy")}, Financial_Quotes_MarketOnCloseImbalanceBuy),
		Entry("Value is []bytes, MarketOnCloseImbalanceSell - Works",
			&types.AttributeValueMemberB{Value: []byte("MarketOnCloseImbalanceSell")}, Financial_Quotes_MarketOnCloseImbalanceSell),
		Entry("Value is []bytes, NoMarketImbalance - Works",
			&types.AttributeValueMemberB{Value: []byte("NoMarketImbalance")}, Financial_Quotes_NoMarketImbalance),
		Entry("Value is []bytes, NoMarketOnCloseImbalance - Works",
			&types.AttributeValueMemberB{Value: []byte("NoMarketOnCloseImbalance")}, Financial_Quotes_NoMarketOnCloseImbalance),
		Entry("Value is []bytes, ShortSaleRestriction - Works",
			&types.AttributeValueMemberB{Value: []byte("ShortSaleRestriction")}, Financial_Quotes_ShortSaleRestriction),
		Entry("Value is []bytes, LimitUpLimitDown - Works",
			&types.AttributeValueMemberB{Value: []byte("LimitUpLimitDown")}, Financial_Quotes_LimitUpLimitDown),
		Entry("Value is []bytes, QuotationResumption - Works",
			&types.AttributeValueMemberB{Value: []byte("QuotationResumption")}, Financial_Quotes_QuotationResumption),
		Entry("Value is []bytes, TradingResumption - Works",
			&types.AttributeValueMemberB{Value: []byte("TradingResumption")}, Financial_Quotes_TradingResumption),
		Entry("Value is []bytes, VolatilityTradingPause - Works",
			&types.AttributeValueMemberB{Value: []byte("VolatilityTradingPause")}, Financial_Quotes_VolatilityTradingPause),
		Entry("Value is []bytes, PolygonReserved - Works",
			&types.AttributeValueMemberB{Value: []byte("PolygonReserved")}, Financial_Quotes_PolygonReserved),
		Entry("Value is []bytes, HaltNewsPending - Works",
			&types.AttributeValueMemberB{Value: []byte("HaltNewsPending")}, Financial_Quotes_HaltNewsPending),
		Entry("Value is []bytes, UpdateNewsDissemination - Works",
			&types.AttributeValueMemberB{Value: []byte("UpdateNewsDissemination")}, Financial_Quotes_UpdateNewsDissemination),
		Entry("Value is []bytes, HaltSingleStockTradingPause - Works",
			&types.AttributeValueMemberB{Value: []byte("HaltSingleStockTradingPause")}, Financial_Quotes_HaltSingleStockTradingPause),
		Entry("Value is []bytes, HaltRegulatoryExtraordinaryMarketActivity - Works",
			&types.AttributeValueMemberB{Value: []byte("HaltRegulatoryExtraordinaryMarketActivity")}, Financial_Quotes_HaltRegulatoryExtraordinaryMarketActivity),
		Entry("Value is []bytes, HaltETF - Works",
			&types.AttributeValueMemberB{Value: []byte("HaltETF")}, Financial_Quotes_HaltETF),
		Entry("Value is []bytes, HaltInformationRequested - Works",
			&types.AttributeValueMemberB{Value: []byte("HaltInformationRequested")}, Financial_Quotes_HaltInformationRequested),
		Entry("Value is []bytes, HaltExchangeNonCompliance - Works",
			&types.AttributeValueMemberB{Value: []byte("HaltExchangeNonCompliance")}, Financial_Quotes_HaltExchangeNonCompliance),
		Entry("Value is []bytes, HaltFilingsNotCurrent - Works",
			&types.AttributeValueMemberB{Value: []byte("HaltFilingsNotCurrent")}, Financial_Quotes_HaltFilingsNotCurrent),
		Entry("Value is []bytes, HaltSECTradingSuspension - Works",
			&types.AttributeValueMemberB{Value: []byte("HaltSECTradingSuspension")}, Financial_Quotes_HaltSECTradingSuspension),
		Entry("Value is []bytes, HaltRegulatoryConcern - Works",
			&types.AttributeValueMemberB{Value: []byte("HaltRegulatoryConcern")}, Financial_Quotes_HaltRegulatoryConcern),
		Entry("Value is []bytes, HaltMarketOperations - Works",
			&types.AttributeValueMemberB{Value: []byte("HaltMarketOperations")}, Financial_Quotes_HaltMarketOperations),
		Entry("Value is []bytes, IPOSecurityNotYetTrading - Works",
			&types.AttributeValueMemberB{Value: []byte("IPOSecurityNotYetTrading")}, Financial_Quotes_IPOSecurityNotYetTrading),
		Entry("Value is []bytes, HaltCorporateAction - Works",
			&types.AttributeValueMemberB{Value: []byte("HaltCorporateAction")}, Financial_Quotes_HaltCorporateAction),
		Entry("Value is []bytes, QuotationNotAvailable - Works",
			&types.AttributeValueMemberB{Value: []byte("QuotationNotAvailable")}, Financial_Quotes_QuotationNotAvailable),
		Entry("Value is []bytes, HaltVolatilityTradingPause - Works",
			&types.AttributeValueMemberB{Value: []byte("HaltVolatilityTradingPause")}, Financial_Quotes_HaltVolatilityTradingPause),
		Entry("Value is []bytes, HaltVolatilityTradingPauseStraddleCondition - Works",
			&types.AttributeValueMemberB{Value: []byte("HaltVolatilityTradingPauseStraddleCondition")}, Financial_Quotes_HaltVolatilityTradingPauseStraddleCondition),
		Entry("Value is []bytes, UpdateNewsAndResumptionTimes - Works",
			&types.AttributeValueMemberB{Value: []byte("UpdateNewsAndResumptionTimes")}, Financial_Quotes_UpdateNewsAndResumptionTimes),
		Entry("Value is []bytes, HaltSingleStockTradingPauseQuotesOnly - Works",
			&types.AttributeValueMemberB{Value: []byte("HaltSingleStockTradingPauseQuotesOnly")}, Financial_Quotes_HaltSingleStockTradingPauseQuotesOnly),
		Entry("Value is []bytes, ResumeQualificationIssuesReviewedResolved - Works",
			&types.AttributeValueMemberB{Value: []byte("ResumeQualificationIssuesReviewedResolved")}, Financial_Quotes_ResumeQualificationIssuesReviewedResolved),
		Entry("Value is []bytes, ResumeFilingRequirementsSatisfiedResolved - Works",
			&types.AttributeValueMemberB{Value: []byte("ResumeFilingRequirementsSatisfiedResolved")}, Financial_Quotes_ResumeFilingRequirementsSatisfiedResolved),
		Entry("Value is []bytes, ResumeNewsNotForthcoming - Works",
			&types.AttributeValueMemberB{Value: []byte("ResumeNewsNotForthcoming")}, Financial_Quotes_ResumeNewsNotForthcoming),
		Entry("Value is []bytes, ResumeQualificationsMaintRequirementsMet - Works",
			&types.AttributeValueMemberB{Value: []byte("ResumeQualificationsMaintRequirementsMet")}, Financial_Quotes_ResumeQualificationsMaintRequirementsMet),
		Entry("Value is []bytes, ResumeQualificationsFilingsMet - Works",
			&types.AttributeValueMemberB{Value: []byte("ResumeQualificationsFilingsMet")}, Financial_Quotes_ResumeQualificationsFilingsMet),
		Entry("Value is []bytes, ResumeRegulatoryAuth - Works",
			&types.AttributeValueMemberB{Value: []byte("ResumeRegulatoryAuth")}, Financial_Quotes_ResumeRegulatoryAuth),
		Entry("Value is []bytes, NewIssueAvailable - Works",
			&types.AttributeValueMemberB{Value: []byte("NewIssueAvailable")}, Financial_Quotes_NewIssueAvailable),
		Entry("Value is []bytes, IssueAvailable - Works",
			&types.AttributeValueMemberB{Value: []byte("IssueAvailable")}, Financial_Quotes_IssueAvailable),
		Entry("Value is []bytes, MWCBCarryFromPreviousDay - Works",
			&types.AttributeValueMemberB{Value: []byte("MWCBCarryFromPreviousDay")}, Financial_Quotes_MWCBCarryFromPreviousDay),
		Entry("Value is []bytes, MWCBResume - Works",
			&types.AttributeValueMemberB{Value: []byte("MWCBResume")}, Financial_Quotes_MWCBResume),
		Entry("Value is []bytes, IPOSecurityReleasedForQuotation - Works",
			&types.AttributeValueMemberB{Value: []byte("IPOSecurityReleasedForQuotation")}, Financial_Quotes_IPOSecurityReleasedForQuotation),
		Entry("Value is []bytes, IPOSecurityPositioningWindowExtension - Works",
			&types.AttributeValueMemberB{Value: []byte("IPOSecurityPositioningWindowExtension")}, Financial_Quotes_IPOSecurityPositioningWindowExtension),
		Entry("Value is []bytes, MWCBLevel1 - Works",
			&types.AttributeValueMemberB{Value: []byte("MWCBLevel1")}, Financial_Quotes_MWCBLevel1),
		Entry("Value is []bytes, MWCBLevel2 - Works",
			&types.AttributeValueMemberB{Value: []byte("MWCBLevel2")}, Financial_Quotes_MWCBLevel2),
		Entry("Value is []bytes, MWCBLevel3 - Works",
			&types.AttributeValueMemberB{Value: []byte("MWCBLevel3")}, Financial_Quotes_MWCBLevel3),
		Entry("Value is []bytes, HaltSubPennyTrading - Works",
			&types.AttributeValueMemberB{Value: []byte("HaltSubPennyTrading")}, Financial_Quotes_HaltSubPennyTrading),
		Entry("Value is []bytes, OrderImbalanceInd - Works",
			&types.AttributeValueMemberB{Value: []byte("OrderImbalanceInd")}, Financial_Quotes_OrderImbalanceInd),
		Entry("Value is []bytes, LULDTradingPaused - Works",
			&types.AttributeValueMemberB{Value: []byte("LULDTradingPaused")}, Financial_Quotes_LULDTradingPaused),
		Entry("Value is []bytes, NONE - Works",
			&types.AttributeValueMemberB{Value: []byte("NONE")}, Financial_Quotes_NONE),
		Entry("Value is []bytes, ShortSalesRestrictionActivated - Works",
			&types.AttributeValueMemberB{Value: []byte("ShortSalesRestrictionActivated")}, Financial_Quotes_ShortSalesRestrictionActivated),
		Entry("Value is []bytes, ShortSalesRestrictionContinued - Works",
			&types.AttributeValueMemberB{Value: []byte("ShortSalesRestrictionContinued")}, Financial_Quotes_ShortSalesRestrictionContinued),
		Entry("Value is []bytes, ShortSalesRestrictionDeactivated - Works",
			&types.AttributeValueMemberB{Value: []byte("ShortSalesRestrictionDeactivated")}, Financial_Quotes_ShortSalesRestrictionDeactivated),
		Entry("Value is []bytes, ShortSalesRestrictionInEffect - Works",
			&types.AttributeValueMemberB{Value: []byte("ShortSalesRestrictionInEffect")}, Financial_Quotes_ShortSalesRestrictionInEffect),
		Entry("Value is []bytes, ShortSalesRestrictionMax - Works",
			&types.AttributeValueMemberB{Value: []byte("ShortSalesRestrictionMax")}, Financial_Quotes_ShortSalesRestrictionMax),
		Entry("Value is []bytes, NBBONoChange - Works",
			&types.AttributeValueMemberB{Value: []byte("NBBONoChange")}, Financial_Quotes_NBBONoChange),
		Entry("Value is []bytes, NBBOQuoteIsNBBO - Works",
			&types.AttributeValueMemberB{Value: []byte("NBBOQuoteIsNBBO")}, Financial_Quotes_NBBOQuoteIsNBBO),
		Entry("Value is []bytes, NBBONoBBNoBO - Works",
			&types.AttributeValueMemberB{Value: []byte("NBBONoBBNoBO")}, Financial_Quotes_NBBONoBBNoBO),
		Entry("Value is []bytes, NBBOBBBOShortAppendage - Works",
			&types.AttributeValueMemberB{Value: []byte("NBBOBBBOShortAppendage")}, Financial_Quotes_NBBOBBBOShortAppendage),
		Entry("Value is []bytes, NBBOBBBOLongAppendage - Works",
			&types.AttributeValueMemberB{Value: []byte("NBBOBBBOLongAppendage")}, Financial_Quotes_NBBOBBBOLongAppendage),
		Entry("Value is []bytes, HeldTradeNotLastSaleNotConsolidated - Works",
			&types.AttributeValueMemberB{Value: []byte("HeldTradeNotLastSaleNotConsolidated")}, Financial_Quotes_HeldTradeNotLastSaleNotConsolidated),
		Entry("Value is []bytes, HeldTradeLastSaleButNotConsolidated - Works",
			&types.AttributeValueMemberB{Value: []byte("HeldTradeLastSaleButNotConsolidated")}, Financial_Quotes_HeldTradeLastSaleButNotConsolidated),
		Entry("Value is []bytes, HeldTradeLastSaleAndConsolidated - Works",
			&types.AttributeValueMemberB{Value: []byte("HeldTradeLastSaleAndConsolidated")}, Financial_Quotes_HeldTradeLastSaleAndConsolidated),
		Entry("Value is []bytes, RetailInterestOnBid - Works",
			&types.AttributeValueMemberB{Value: []byte("RetailInterestOnBid")}, Financial_Quotes_RetailInterestOnBid),
		Entry("Value is []bytes, RetailInterestOnAsk - Works",
			&types.AttributeValueMemberB{Value: []byte("RetailInterestOnAsk")}, Financial_Quotes_RetailInterestOnAsk),
		Entry("Value is []bytes, RetailInterestOnBidAndAsk - Works",
			&types.AttributeValueMemberB{Value: []byte("RetailInterestOnBidAndAsk")}, Financial_Quotes_RetailInterestOnBidAndAsk),
		Entry("Value is []bytes, FinraBBONoChange - Works",
			&types.AttributeValueMemberB{Value: []byte("FinraBBONoChange")}, Financial_Quotes_FinraBBONoChange),
		Entry("Value is []bytes, FinraBBODoesNotExist - Works",
			&types.AttributeValueMemberB{Value: []byte("FinraBBODoesNotExist")}, Financial_Quotes_FinraBBODoesNotExist),
		Entry("Value is []bytes, FinraBBBOExecutable - Works",
			&types.AttributeValueMemberB{Value: []byte("FinraBBBOExecutable")}, Financial_Quotes_FinraBBBOExecutable),
		Entry("Value is []bytes, FinraBBBelowLowerBand - Works",
			&types.AttributeValueMemberB{Value: []byte("FinraBBBelowLowerBand")}, Financial_Quotes_FinraBBBelowLowerBand),
		Entry("Value is []bytes, FinraBOAboveUpperBand - Works",
			&types.AttributeValueMemberB{Value: []byte("FinraBOAboveUpperBand")}, Financial_Quotes_FinraBOAboveUpperBand),
		Entry("Value is []bytes, FinraBBBelowLowerBandBOAbboveUpperBand - Works",
			&types.AttributeValueMemberB{Value: []byte("FinraBBBelowLowerBandBOAbboveUpperBand")}, Financial_Quotes_FinraBBBelowLowerBandBOAbboveUpperBand),
		Entry("Value is []bytes, CTANotDueToRelatedSecurity - Works",
			&types.AttributeValueMemberB{Value: []byte("CTANotDueToRelatedSecurity")}, Financial_Quotes_CTANotDueToRelatedSecurity),
		Entry("Value is []bytes, CTADueToRelatedSecurity - Works",
			&types.AttributeValueMemberB{Value: []byte("CTADueToRelatedSecurity")}, Financial_Quotes_CTADueToRelatedSecurity),
		Entry("Value is []bytes, CTANotInViewOfCommon - Works",
			&types.AttributeValueMemberB{Value: []byte("CTANotInViewOfCommon")}, Financial_Quotes_CTANotInViewOfCommon),
		Entry("Value is []bytes, CTAInViewOfCommon - Works",
			&types.AttributeValueMemberB{Value: []byte("CTAInViewOfCommon")}, Financial_Quotes_CTAInViewOfCommon),
		Entry("Value is []bytes, CTAPriceIndicator - Works",
			&types.AttributeValueMemberB{Value: []byte("CTAPriceIndicator")}, Financial_Quotes_CTAPriceIndicator),
		Entry("Value is []bytes, CTANewPriceIndicator - Works",
			&types.AttributeValueMemberB{Value: []byte("CTANewPriceIndicator")}, Financial_Quotes_CTANewPriceIndicator),
		Entry("Value is []bytes, CTACorrectedPriceIndication - Works",
			&types.AttributeValueMemberB{Value: []byte("CTACorrectedPriceIndication")}, Financial_Quotes_CTACorrectedPriceIndication),
		Entry("Value is []bytes, CTACancelledMarketImbalance - Works",
			&types.AttributeValueMemberB{Value: []byte("CTACancelledMarketImbalance")}, Financial_Quotes_CTACancelledMarketImbalance),
		Entry("Value is numeric, 0 - Works",
			&types.AttributeValueMemberN{Value: "0"}, Financial_Quotes_NBBNBOExecutable),
		Entry("Value is numeric, 1 - Works",
			&types.AttributeValueMemberN{Value: "1"}, Financial_Quotes_NBBBelowLowerBand),
		Entry("Value is numeric, 2 - Works",
			&types.AttributeValueMemberN{Value: "2"}, Financial_Quotes_NBOAboveUpperBand),
		Entry("Value is numeric, 3 - Works",
			&types.AttributeValueMemberN{Value: "3"}, Financial_Quotes_NBBBelowLowerBandAndNBOAboveUpperBand),
		Entry("Value is numeric, 4 - Works",
			&types.AttributeValueMemberN{Value: "4"}, Financial_Quotes_NBBEqualsUpperBand),
		Entry("Value is numeric, 5 - Works",
			&types.AttributeValueMemberN{Value: "5"}, Financial_Quotes_NBOEqualsLowerBand),
		Entry("Value is numeric, 6 - Works",
			&types.AttributeValueMemberN{Value: "6"}, Financial_Quotes_NBBEqualsUpperBandAndNBOAboveUpperBand),
		Entry("Value is numeric, 7 - Works",
			&types.AttributeValueMemberN{Value: "7"}, Financial_Quotes_NBBBelowLowerBandAndNBOEqualsLowerBand),
		Entry("Value is numeric, 8 - Works",
			&types.AttributeValueMemberN{Value: "8"}, Financial_Quotes_BidPriceAboveUpperLimitPriceBand),
		Entry("Value is numeric, 9 - Works",
			&types.AttributeValueMemberN{Value: "9"}, Financial_Quotes_OfferPriceBelowLowerLimitPriceBand),
		Entry("Value is numeric, 10 - Works",
			&types.AttributeValueMemberN{Value: "10"}, Financial_Quotes_BidAndOfferOutsidePriceBand),
		Entry("Value is numeric, 11 - Works",
			&types.AttributeValueMemberN{Value: "11"}, Financial_Quotes_OpeningUpdate),
		Entry("Value is numeric, 12 - Works",
			&types.AttributeValueMemberN{Value: "12"}, Financial_Quotes_IntraDayUpdate),
		Entry("Value is numeric, 13 - Works",
			&types.AttributeValueMemberN{Value: "13"}, Financial_Quotes_RestatedValue),
		Entry("Value is numeric, 14 - Works",
			&types.AttributeValueMemberN{Value: "14"}, Financial_Quotes_SuspendedDuringTradingHalt),
		Entry("Value is numeric, 15 - Works",
			&types.AttributeValueMemberN{Value: "15"}, Financial_Quotes_ReOpeningUpdate),
		Entry("Value is numeric, 16 - Works",
			&types.AttributeValueMemberN{Value: "16"}, Financial_Quotes_OutsidePriceBandRuleHours),
		Entry("Value is numeric, 17 - Works",
			&types.AttributeValueMemberN{Value: "17"}, Financial_Quotes_AuctionExtension),
		Entry("Value is numeric, 18 - Works",
			&types.AttributeValueMemberN{Value: "18"}, Financial_Quotes_LULDPriceBand),
		Entry("Value is numeric, 19 - Works",
			&types.AttributeValueMemberN{Value: "19"}, Financial_Quotes_RepublishedLULDPriceBandInd),
		Entry("Value is numeric, 20 - Works",
			&types.AttributeValueMemberN{Value: "20"}, Financial_Quotes_NBBLimitStateEntered),
		Entry("Value is numeric, 21 - Works",
			&types.AttributeValueMemberN{Value: "21"}, Financial_Quotes_NBBLimitStateExited),
		Entry("Value is numeric, 22 - Works",
			&types.AttributeValueMemberN{Value: "22"}, Financial_Quotes_NBOLimitStateEntered),
		Entry("Value is numeric, 23 - Works",
			&types.AttributeValueMemberN{Value: "23"}, Financial_Quotes_NBOLimitStateExited),
		Entry("Value is numeric, 24 - Works",
			&types.AttributeValueMemberN{Value: "24"}, Financial_Quotes_NBBAndNBOLimitStateEntered),
		Entry("Value is numeric, 25 - Works",
			&types.AttributeValueMemberN{Value: "25"}, Financial_Quotes_NBBAndNBOLimitStateExited),
		Entry("Value is numeric, 26 - Works",
			&types.AttributeValueMemberN{Value: "26"}, Financial_Quotes_NBBLimitStateEnteredNBOLimitStateExited),
		Entry("Value is numeric, 27 - Works",
			&types.AttributeValueMemberN{Value: "27"}, Financial_Quotes_NBBLimitStateExitedNBOLimitStateEntered),
		Entry("Value is numeric, 28 - Works",
			&types.AttributeValueMemberN{Value: "28"}, Financial_Quotes_Normal),
		Entry("Value is numeric, 29 - Works",
			&types.AttributeValueMemberN{Value: "29"}, Financial_Quotes_Bankrupt),
		Entry("Value is numeric, 30 - Works",
			&types.AttributeValueMemberN{Value: "30"}, Financial_Quotes_Deficient),
		Entry("Value is numeric, 31 - Works",
			&types.AttributeValueMemberN{Value: "31"}, Financial_Quotes_Delinquent),
		Entry("Value is numeric, 32 - Works",
			&types.AttributeValueMemberN{Value: "32"}, Financial_Quotes_BankruptAndDeficient),
		Entry("Value is numeric, 33 - Works",
			&types.AttributeValueMemberN{Value: "33"}, Financial_Quotes_BankruptAndDelinquent),
		Entry("Value is numeric, 34 - Works",
			&types.AttributeValueMemberN{Value: "34"}, Financial_Quotes_DeficientAndDelinquent),
		Entry("Value is numeric, 35 - Works",
			&types.AttributeValueMemberN{Value: "35"}, Financial_Quotes_DeficientDeliquentBankrupt),
		Entry("Value is numeric, 36 - Works",
			&types.AttributeValueMemberN{Value: "36"}, Financial_Quotes_Liquidation),
		Entry("Value is numeric, 37 - Works",
			&types.AttributeValueMemberN{Value: "37"}, Financial_Quotes_CreationsSuspended),
		Entry("Value is numeric, 38 - Works",
			&types.AttributeValueMemberN{Value: "38"}, Financial_Quotes_RedemptionsSuspended),
		Entry("Value is numeric, 39 - Works",
			&types.AttributeValueMemberN{Value: "39"}, Financial_Quotes_CreationsRedemptionsSuspended),
		Entry("Value is numeric, 40 - Works",
			&types.AttributeValueMemberN{Value: "40"}, Financial_Quotes_NormalTrading),
		Entry("Value is numeric, 41 - Works",
			&types.AttributeValueMemberN{Value: "41"}, Financial_Quotes_OpeningDelay),
		Entry("Value is numeric, 42 - Works",
			&types.AttributeValueMemberN{Value: "42"}, Financial_Quotes_TradingHalt),
		Entry("Value is numeric, 43 - Works",
			&types.AttributeValueMemberN{Value: "43"}, Financial_Quotes_TradingResume),
		Entry("Value is numeric, 44 - Works",
			&types.AttributeValueMemberN{Value: "44"}, Financial_Quotes_NoOpenNoResume),
		Entry("Value is numeric, 45 - Works",
			&types.AttributeValueMemberN{Value: "45"}, Financial_Quotes_PriceIndication),
		Entry("Value is numeric, 46 - Works",
			&types.AttributeValueMemberN{Value: "46"}, Financial_Quotes_TradingRangeIndication),
		Entry("Value is numeric, 47 - Works",
			&types.AttributeValueMemberN{Value: "47"}, Financial_Quotes_MarketImbalanceBuy),
		Entry("Value is numeric, 48 - Works",
			&types.AttributeValueMemberN{Value: "48"}, Financial_Quotes_MarketImbalanceSell),
		Entry("Value is numeric, 49 - Works",
			&types.AttributeValueMemberN{Value: "49"}, Financial_Quotes_MarketOnCloseImbalanceBuy),
		Entry("Value is numeric, 50 - Works",
			&types.AttributeValueMemberN{Value: "50"}, Financial_Quotes_MarketOnCloseImbalanceSell),
		Entry("Value is numeric, 51 - Works",
			&types.AttributeValueMemberN{Value: "51"}, Financial_Quotes_NoMarketImbalance),
		Entry("Value is numeric, 52 - Works",
			&types.AttributeValueMemberN{Value: "52"}, Financial_Quotes_NoMarketOnCloseImbalance),
		Entry("Value is numeric, 53 - Works",
			&types.AttributeValueMemberN{Value: "53"}, Financial_Quotes_ShortSaleRestriction),
		Entry("Value is numeric, 54 - Works",
			&types.AttributeValueMemberN{Value: "54"}, Financial_Quotes_LimitUpLimitDown),
		Entry("Value is numeric, 55 - Works",
			&types.AttributeValueMemberN{Value: "55"}, Financial_Quotes_QuotationResumption),
		Entry("Value is numeric, 56 - Works",
			&types.AttributeValueMemberN{Value: "56"}, Financial_Quotes_TradingResumption),
		Entry("Value is numeric, 57 - Works",
			&types.AttributeValueMemberN{Value: "57"}, Financial_Quotes_VolatilityTradingPause),
		Entry("Value is numeric, 58 - Works",
			&types.AttributeValueMemberN{Value: "58"}, Financial_Quotes_PolygonReserved),
		Entry("Value is numeric, 59 - Works",
			&types.AttributeValueMemberN{Value: "59"}, Financial_Quotes_HaltNewsPending),
		Entry("Value is numeric, 60 - Works",
			&types.AttributeValueMemberN{Value: "60"}, Financial_Quotes_UpdateNewsDissemination),
		Entry("Value is numeric, 61 - Works",
			&types.AttributeValueMemberN{Value: "61"}, Financial_Quotes_HaltSingleStockTradingPause),
		Entry("Value is numeric, 62 - Works",
			&types.AttributeValueMemberN{Value: "62"}, Financial_Quotes_HaltRegulatoryExtraordinaryMarketActivity),
		Entry("Value is numeric, 63 - Works",
			&types.AttributeValueMemberN{Value: "63"}, Financial_Quotes_HaltETF),
		Entry("Value is numeric, 64 - Works",
			&types.AttributeValueMemberN{Value: "64"}, Financial_Quotes_HaltInformationRequested),
		Entry("Value is numeric, 65 - Works",
			&types.AttributeValueMemberN{Value: "65"}, Financial_Quotes_HaltExchangeNonCompliance),
		Entry("Value is numeric, 66 - Works",
			&types.AttributeValueMemberN{Value: "66"}, Financial_Quotes_HaltFilingsNotCurrent),
		Entry("Value is numeric, 67 - Works",
			&types.AttributeValueMemberN{Value: "67"}, Financial_Quotes_HaltSECTradingSuspension),
		Entry("Value is numeric, 68 - Works",
			&types.AttributeValueMemberN{Value: "68"}, Financial_Quotes_HaltRegulatoryConcern),
		Entry("Value is numeric, 69 - Works",
			&types.AttributeValueMemberN{Value: "69"}, Financial_Quotes_HaltMarketOperations),
		Entry("Value is numeric, 70 - Works",
			&types.AttributeValueMemberN{Value: "70"}, Financial_Quotes_IPOSecurityNotYetTrading),
		Entry("Value is numeric, 71 - Works",
			&types.AttributeValueMemberN{Value: "71"}, Financial_Quotes_HaltCorporateAction),
		Entry("Value is numeric, 72 - Works",
			&types.AttributeValueMemberN{Value: "72"}, Financial_Quotes_QuotationNotAvailable),
		Entry("Value is numeric, 73 - Works",
			&types.AttributeValueMemberN{Value: "73"}, Financial_Quotes_HaltVolatilityTradingPause),
		Entry("Value is numeric, 74 - Works",
			&types.AttributeValueMemberN{Value: "74"}, Financial_Quotes_HaltVolatilityTradingPauseStraddleCondition),
		Entry("Value is numeric, 75 - Works",
			&types.AttributeValueMemberN{Value: "75"}, Financial_Quotes_UpdateNewsAndResumptionTimes),
		Entry("Value is numeric, 76 - Works",
			&types.AttributeValueMemberN{Value: "76"}, Financial_Quotes_HaltSingleStockTradingPauseQuotesOnly),
		Entry("Value is numeric, 77 - Works",
			&types.AttributeValueMemberN{Value: "77"}, Financial_Quotes_ResumeQualificationIssuesReviewedResolved),
		Entry("Value is numeric, 78 - Works",
			&types.AttributeValueMemberN{Value: "78"}, Financial_Quotes_ResumeFilingRequirementsSatisfiedResolved),
		Entry("Value is numeric, 79 - Works",
			&types.AttributeValueMemberN{Value: "79"}, Financial_Quotes_ResumeNewsNotForthcoming),
		Entry("Value is numeric, 80 - Works",
			&types.AttributeValueMemberN{Value: "80"}, Financial_Quotes_ResumeQualificationsMaintRequirementsMet),
		Entry("Value is numeric, 81 - Works",
			&types.AttributeValueMemberN{Value: "81"}, Financial_Quotes_ResumeQualificationsFilingsMet),
		Entry("Value is numeric, 82 - Works",
			&types.AttributeValueMemberN{Value: "82"}, Financial_Quotes_ResumeRegulatoryAuth),
		Entry("Value is numeric, 83 - Works",
			&types.AttributeValueMemberN{Value: "83"}, Financial_Quotes_NewIssueAvailable),
		Entry("Value is numeric, 84 - Works",
			&types.AttributeValueMemberN{Value: "84"}, Financial_Quotes_IssueAvailable),
		Entry("Value is numeric, 85 - Works",
			&types.AttributeValueMemberN{Value: "85"}, Financial_Quotes_MWCBCarryFromPreviousDay),
		Entry("Value is numeric, 86 - Works",
			&types.AttributeValueMemberN{Value: "86"}, Financial_Quotes_MWCBResume),
		Entry("Value is numeric, 87 - Works",
			&types.AttributeValueMemberN{Value: "87"}, Financial_Quotes_IPOSecurityReleasedForQuotation),
		Entry("Value is numeric, 88 - Works",
			&types.AttributeValueMemberN{Value: "88"}, Financial_Quotes_IPOSecurityPositioningWindowExtension),
		Entry("Value is numeric, 89 - Works",
			&types.AttributeValueMemberN{Value: "89"}, Financial_Quotes_MWCBLevel1),
		Entry("Value is numeric, 90 - Works",
			&types.AttributeValueMemberN{Value: "90"}, Financial_Quotes_MWCBLevel2),
		Entry("Value is numeric, 91 - Works",
			&types.AttributeValueMemberN{Value: "91"}, Financial_Quotes_MWCBLevel3),
		Entry("Value is numeric, 92 - Works",
			&types.AttributeValueMemberN{Value: "92"}, Financial_Quotes_HaltSubPennyTrading),
		Entry("Value is numeric, 93 - Works",
			&types.AttributeValueMemberN{Value: "93"}, Financial_Quotes_OrderImbalanceInd),
		Entry("Value is numeric, 94 - Works",
			&types.AttributeValueMemberN{Value: "94"}, Financial_Quotes_LULDTradingPaused),
		Entry("Value is numeric, 95 - Works",
			&types.AttributeValueMemberN{Value: "95"}, Financial_Quotes_NONE),
		Entry("Value is numeric, 96 - Works",
			&types.AttributeValueMemberN{Value: "96"}, Financial_Quotes_ShortSalesRestrictionActivated),
		Entry("Value is numeric, 97 - Works",
			&types.AttributeValueMemberN{Value: "97"}, Financial_Quotes_ShortSalesRestrictionContinued),
		Entry("Value is numeric, 98 - Works",
			&types.AttributeValueMemberN{Value: "98"}, Financial_Quotes_ShortSalesRestrictionDeactivated),
		Entry("Value is numeric, 99 - Works",
			&types.AttributeValueMemberN{Value: "99"}, Financial_Quotes_ShortSalesRestrictionInEffect),
		Entry("Value is numeric, 100 - Works",
			&types.AttributeValueMemberN{Value: "100"}, Financial_Quotes_ShortSalesRestrictionMax),
		Entry("Value is numeric, 101 - Works",
			&types.AttributeValueMemberN{Value: "101"}, Financial_Quotes_NBBONoChange),
		Entry("Value is numeric, 102 - Works",
			&types.AttributeValueMemberN{Value: "102"}, Financial_Quotes_NBBOQuoteIsNBBO),
		Entry("Value is numeric, 103 - Works",
			&types.AttributeValueMemberN{Value: "103"}, Financial_Quotes_NBBONoBBNoBO),
		Entry("Value is numeric, 104 - Works",
			&types.AttributeValueMemberN{Value: "104"}, Financial_Quotes_NBBOBBBOShortAppendage),
		Entry("Value is numeric, 105 - Works",
			&types.AttributeValueMemberN{Value: "105"}, Financial_Quotes_NBBOBBBOLongAppendage),
		Entry("Value is numeric, 106 - Works",
			&types.AttributeValueMemberN{Value: "106"}, Financial_Quotes_HeldTradeNotLastSaleNotConsolidated),
		Entry("Value is numeric, 107 - Works",
			&types.AttributeValueMemberN{Value: "107"}, Financial_Quotes_HeldTradeLastSaleButNotConsolidated),
		Entry("Value is numeric, 108 - Works",
			&types.AttributeValueMemberN{Value: "108"}, Financial_Quotes_HeldTradeLastSaleAndConsolidated),
		Entry("Value is numeric, 109 - Works",
			&types.AttributeValueMemberN{Value: "109"}, Financial_Quotes_RetailInterestOnBid),
		Entry("Value is numeric, 110 - Works",
			&types.AttributeValueMemberN{Value: "110"}, Financial_Quotes_RetailInterestOnAsk),
		Entry("Value is numeric, 111 - Works",
			&types.AttributeValueMemberN{Value: "111"}, Financial_Quotes_RetailInterestOnBidAndAsk),
		Entry("Value is numeric, 112 - Works",
			&types.AttributeValueMemberN{Value: "112"}, Financial_Quotes_FinraBBONoChange),
		Entry("Value is numeric, 113 - Works",
			&types.AttributeValueMemberN{Value: "113"}, Financial_Quotes_FinraBBODoesNotExist),
		Entry("Value is numeric, 114 - Works",
			&types.AttributeValueMemberN{Value: "114"}, Financial_Quotes_FinraBBBOExecutable),
		Entry("Value is numeric, 115 - Works",
			&types.AttributeValueMemberN{Value: "115"}, Financial_Quotes_FinraBBBelowLowerBand),
		Entry("Value is numeric, 116 - Works",
			&types.AttributeValueMemberN{Value: "116"}, Financial_Quotes_FinraBOAboveUpperBand),
		Entry("Value is numeric, 117 - Works",
			&types.AttributeValueMemberN{Value: "117"}, Financial_Quotes_FinraBBBelowLowerBandBOAbboveUpperBand),
		Entry("Value is numeric, 118 - Works",
			&types.AttributeValueMemberN{Value: "118"}, Financial_Quotes_CTANotDueToRelatedSecurity),
		Entry("Value is numeric, 119 - Works",
			&types.AttributeValueMemberN{Value: "119"}, Financial_Quotes_CTADueToRelatedSecurity),
		Entry("Value is numeric, 120 - Works",
			&types.AttributeValueMemberN{Value: "120"}, Financial_Quotes_CTANotInViewOfCommon),
		Entry("Value is numeric, 121 - Works",
			&types.AttributeValueMemberN{Value: "121"}, Financial_Quotes_CTAInViewOfCommon),
		Entry("Value is numeric, 122 - Works",
			&types.AttributeValueMemberN{Value: "122"}, Financial_Quotes_CTAPriceIndicator),
		Entry("Value is numeric, 123 - Works",
			&types.AttributeValueMemberN{Value: "123"}, Financial_Quotes_CTANewPriceIndicator),
		Entry("Value is numeric, 124 - Works",
			&types.AttributeValueMemberN{Value: "124"}, Financial_Quotes_CTACorrectedPriceIndication),
		Entry("Value is numeric, 125 - Works",
			&types.AttributeValueMemberN{Value: "125"}, Financial_Quotes_CTACancelledMarketImbalance),
		Entry("Value is NULL - Works", new(types.AttributeValueMemberNULL), Financial_Quotes_Indicator(0)),
		Entry("Value is string, NBBNBOExecutable - Works",
			&types.AttributeValueMemberS{Value: "NBBNBOExecutable"}, Financial_Quotes_NBBNBOExecutable),
		Entry("Value is string, NBBBelowLowerBand - Works",
			&types.AttributeValueMemberS{Value: "NBBBelowLowerBand"}, Financial_Quotes_NBBBelowLowerBand),
		Entry("Value is string, NBOAboveUpperBand - Works",
			&types.AttributeValueMemberS{Value: "NBOAboveUpperBand"}, Financial_Quotes_NBOAboveUpperBand),
		Entry("Value is string, NBBBelowLowerBandAndNBOAboveUpperBand - Works",
			&types.AttributeValueMemberS{Value: "NBBBelowLowerBandAndNBOAboveUpperBand"}, Financial_Quotes_NBBBelowLowerBandAndNBOAboveUpperBand),
		Entry("Value is string, NBBEqualsUpperBand - Works",
			&types.AttributeValueMemberS{Value: "NBBEqualsUpperBand"}, Financial_Quotes_NBBEqualsUpperBand),
		Entry("Value is string, NBOEqualsLowerBand - Works",
			&types.AttributeValueMemberS{Value: "NBOEqualsLowerBand"}, Financial_Quotes_NBOEqualsLowerBand),
		Entry("Value is string, NBBEqualsUpperBandAndNBOAboveUpperBand - Works",
			&types.AttributeValueMemberS{Value: "NBBEqualsUpperBandAndNBOAboveUpperBand"}, Financial_Quotes_NBBEqualsUpperBandAndNBOAboveUpperBand),
		Entry("Value is string, NBBBelowLowerBandAndNBOEqualsLowerBand - Works",
			&types.AttributeValueMemberS{Value: "NBBBelowLowerBandAndNBOEqualsLowerBand"}, Financial_Quotes_NBBBelowLowerBandAndNBOEqualsLowerBand),
		Entry("Value is string, BidPriceAboveUpperLimitPriceBand - Works",
			&types.AttributeValueMemberS{Value: "BidPriceAboveUpperLimitPriceBand"}, Financial_Quotes_BidPriceAboveUpperLimitPriceBand),
		Entry("Value is string, OfferPriceBelowLowerLimitPriceBand - Works",
			&types.AttributeValueMemberS{Value: "OfferPriceBelowLowerLimitPriceBand"}, Financial_Quotes_OfferPriceBelowLowerLimitPriceBand),
		Entry("Value is string, BidAndOfferOutsidePriceBand - Works",
			&types.AttributeValueMemberS{Value: "BidAndOfferOutsidePriceBand"}, Financial_Quotes_BidAndOfferOutsidePriceBand),
		Entry("Value is string, OpeningUpdate - Works",
			&types.AttributeValueMemberS{Value: "OpeningUpdate"}, Financial_Quotes_OpeningUpdate),
		Entry("Value is string, IntraDayUpdate - Works",
			&types.AttributeValueMemberS{Value: "IntraDayUpdate"}, Financial_Quotes_IntraDayUpdate),
		Entry("Value is string, RestatedValue - Works",
			&types.AttributeValueMemberS{Value: "RestatedValue"}, Financial_Quotes_RestatedValue),
		Entry("Value is string, SuspendedDuringTradingHalt - Works",
			&types.AttributeValueMemberS{Value: "SuspendedDuringTradingHalt"}, Financial_Quotes_SuspendedDuringTradingHalt),
		Entry("Value is string, ReOpeningUpdate - Works",
			&types.AttributeValueMemberS{Value: "ReOpeningUpdate"}, Financial_Quotes_ReOpeningUpdate),
		Entry("Value is string, OutsidePriceBandRuleHours - Works",
			&types.AttributeValueMemberS{Value: "OutsidePriceBandRuleHours"}, Financial_Quotes_OutsidePriceBandRuleHours),
		Entry("Value is string, AuctionExtension - Works",
			&types.AttributeValueMemberS{Value: "AuctionExtension"}, Financial_Quotes_AuctionExtension),
		Entry("Value is string, LULDPriceBand - Works",
			&types.AttributeValueMemberS{Value: "LULDPriceBand"}, Financial_Quotes_LULDPriceBand),
		Entry("Value is string, RepublishedLULDPriceBandInd - Works",
			&types.AttributeValueMemberS{Value: "RepublishedLULDPriceBandInd"}, Financial_Quotes_RepublishedLULDPriceBandInd),
		Entry("Value is string, NBBLimitStateEntered - Works",
			&types.AttributeValueMemberS{Value: "NBBLimitStateEntered"}, Financial_Quotes_NBBLimitStateEntered),
		Entry("Value is string, NBBLimitStateExited - Works",
			&types.AttributeValueMemberS{Value: "NBBLimitStateExited"}, Financial_Quotes_NBBLimitStateExited),
		Entry("Value is string, NBOLimitStateEntered - Works",
			&types.AttributeValueMemberS{Value: "NBOLimitStateEntered"}, Financial_Quotes_NBOLimitStateEntered),
		Entry("Value is string, NBOLimitStateExited - Works",
			&types.AttributeValueMemberS{Value: "NBOLimitStateExited"}, Financial_Quotes_NBOLimitStateExited),
		Entry("Value is string, NBBAndNBOLimitStateEntered - Works",
			&types.AttributeValueMemberS{Value: "NBBAndNBOLimitStateEntered"}, Financial_Quotes_NBBAndNBOLimitStateEntered),
		Entry("Value is string, NBBAndNBOLimitStateExited - Works",
			&types.AttributeValueMemberS{Value: "NBBAndNBOLimitStateExited"}, Financial_Quotes_NBBAndNBOLimitStateExited),
		Entry("Value is string, NBBLimitStateEnteredNBOLimitStateExited - Works",
			&types.AttributeValueMemberS{Value: "NBBLimitStateEnteredNBOLimitStateExited"}, Financial_Quotes_NBBLimitStateEnteredNBOLimitStateExited),
		Entry("Value is string, NBBLimitStateExitedNBOLimitStateEntered - Works",
			&types.AttributeValueMemberS{Value: "NBBLimitStateExitedNBOLimitStateEntered"}, Financial_Quotes_NBBLimitStateExitedNBOLimitStateEntered),
		Entry("Value is string, Normal - Works",
			&types.AttributeValueMemberS{Value: "Normal"}, Financial_Quotes_Normal),
		Entry("Value is string, Bankrupt - Works",
			&types.AttributeValueMemberS{Value: "Bankrupt"}, Financial_Quotes_Bankrupt),
		Entry("Value is string, Deficient - Works",
			&types.AttributeValueMemberS{Value: "Deficient"}, Financial_Quotes_Deficient),
		Entry("Value is string, Delinquent - Works",
			&types.AttributeValueMemberS{Value: "Delinquent"}, Financial_Quotes_Delinquent),
		Entry("Value is string, BankruptAndDeficient - Works",
			&types.AttributeValueMemberS{Value: "BankruptAndDeficient"}, Financial_Quotes_BankruptAndDeficient),
		Entry("Value is string, BankruptAndDelinquent - Works",
			&types.AttributeValueMemberS{Value: "BankruptAndDelinquent"}, Financial_Quotes_BankruptAndDelinquent),
		Entry("Value is string, DeficientAndDelinquent - Works",
			&types.AttributeValueMemberS{Value: "DeficientAndDelinquent"}, Financial_Quotes_DeficientAndDelinquent),
		Entry("Value is string, DeficientDeliquentBankrupt - Works",
			&types.AttributeValueMemberS{Value: "DeficientDeliquentBankrupt"}, Financial_Quotes_DeficientDeliquentBankrupt),
		Entry("Value is string, Liquidation - Works",
			&types.AttributeValueMemberS{Value: "Liquidation"}, Financial_Quotes_Liquidation),
		Entry("Value is string, CreationsSuspended - Works",
			&types.AttributeValueMemberS{Value: "CreationsSuspended"}, Financial_Quotes_CreationsSuspended),
		Entry("Value is string, RedemptionsSuspended - Works",
			&types.AttributeValueMemberS{Value: "RedemptionsSuspended"}, Financial_Quotes_RedemptionsSuspended),
		Entry("Value is string, CreationsRedemptionsSuspended - Works",
			&types.AttributeValueMemberS{Value: "CreationsRedemptionsSuspended"}, Financial_Quotes_CreationsRedemptionsSuspended),
		Entry("Value is string, NormalTrading - Works",
			&types.AttributeValueMemberS{Value: "NormalTrading"}, Financial_Quotes_NormalTrading),
		Entry("Value is string, OpeningDelay - Works",
			&types.AttributeValueMemberS{Value: "OpeningDelay"}, Financial_Quotes_OpeningDelay),
		Entry("Value is string, TradingHalt - Works",
			&types.AttributeValueMemberS{Value: "TradingHalt"}, Financial_Quotes_TradingHalt),
		Entry("Value is string, TradingResume - Works",
			&types.AttributeValueMemberS{Value: "TradingResume"}, Financial_Quotes_TradingResume),
		Entry("Value is string, NoOpenNoResume - Works",
			&types.AttributeValueMemberS{Value: "NoOpenNoResume"}, Financial_Quotes_NoOpenNoResume),
		Entry("Value is string, PriceIndication - Works",
			&types.AttributeValueMemberS{Value: "PriceIndication"}, Financial_Quotes_PriceIndication),
		Entry("Value is string, TradingRangeIndication - Works",
			&types.AttributeValueMemberS{Value: "TradingRangeIndication"}, Financial_Quotes_TradingRangeIndication),
		Entry("Value is string, MarketImbalanceBuy - Works",
			&types.AttributeValueMemberS{Value: "MarketImbalanceBuy"}, Financial_Quotes_MarketImbalanceBuy),
		Entry("Value is string, MarketImbalanceSell - Works",
			&types.AttributeValueMemberS{Value: "MarketImbalanceSell"}, Financial_Quotes_MarketImbalanceSell),
		Entry("Value is string, MarketOnCloseImbalanceBuy - Works",
			&types.AttributeValueMemberS{Value: "MarketOnCloseImbalanceBuy"}, Financial_Quotes_MarketOnCloseImbalanceBuy),
		Entry("Value is string, MarketOnCloseImbalanceSell - Works",
			&types.AttributeValueMemberS{Value: "MarketOnCloseImbalanceSell"}, Financial_Quotes_MarketOnCloseImbalanceSell),
		Entry("Value is string, NoMarketImbalance - Works",
			&types.AttributeValueMemberS{Value: "NoMarketImbalance"}, Financial_Quotes_NoMarketImbalance),
		Entry("Value is string, NoMarketOnCloseImbalance - Works",
			&types.AttributeValueMemberS{Value: "NoMarketOnCloseImbalance"}, Financial_Quotes_NoMarketOnCloseImbalance),
		Entry("Value is string, ShortSaleRestriction - Works",
			&types.AttributeValueMemberS{Value: "ShortSaleRestriction"}, Financial_Quotes_ShortSaleRestriction),
		Entry("Value is string, LimitUpLimitDown - Works",
			&types.AttributeValueMemberS{Value: "LimitUpLimitDown"}, Financial_Quotes_LimitUpLimitDown),
		Entry("Value is string, QuotationResumption - Works",
			&types.AttributeValueMemberS{Value: "QuotationResumption"}, Financial_Quotes_QuotationResumption),
		Entry("Value is string, TradingResumption - Works",
			&types.AttributeValueMemberS{Value: "TradingResumption"}, Financial_Quotes_TradingResumption),
		Entry("Value is string, VolatilityTradingPause - Works",
			&types.AttributeValueMemberS{Value: "VolatilityTradingPause"}, Financial_Quotes_VolatilityTradingPause),
		Entry("Value is string, PolygonReserved - Works",
			&types.AttributeValueMemberS{Value: "PolygonReserved"}, Financial_Quotes_PolygonReserved),
		Entry("Value is string, HaltNewsPending - Works",
			&types.AttributeValueMemberS{Value: "HaltNewsPending"}, Financial_Quotes_HaltNewsPending),
		Entry("Value is string, UpdateNewsDissemination - Works",
			&types.AttributeValueMemberS{Value: "UpdateNewsDissemination"}, Financial_Quotes_UpdateNewsDissemination),
		Entry("Value is string, HaltSingleStockTradingPause - Works",
			&types.AttributeValueMemberS{Value: "HaltSingleStockTradingPause"}, Financial_Quotes_HaltSingleStockTradingPause),
		Entry("Value is string, HaltRegulatoryExtraordinaryMarketActivity - Works",
			&types.AttributeValueMemberS{Value: "HaltRegulatoryExtraordinaryMarketActivity"}, Financial_Quotes_HaltRegulatoryExtraordinaryMarketActivity),
		Entry("Value is string, HaltETF - Works",
			&types.AttributeValueMemberS{Value: "HaltETF"}, Financial_Quotes_HaltETF),
		Entry("Value is string, HaltInformationRequested - Works",
			&types.AttributeValueMemberS{Value: "HaltInformationRequested"}, Financial_Quotes_HaltInformationRequested),
		Entry("Value is string, HaltExchangeNonCompliance - Works",
			&types.AttributeValueMemberS{Value: "HaltExchangeNonCompliance"}, Financial_Quotes_HaltExchangeNonCompliance),
		Entry("Value is string, HaltFilingsNotCurrent - Works",
			&types.AttributeValueMemberS{Value: "HaltFilingsNotCurrent"}, Financial_Quotes_HaltFilingsNotCurrent),
		Entry("Value is string, HaltSECTradingSuspension - Works",
			&types.AttributeValueMemberS{Value: "HaltSECTradingSuspension"}, Financial_Quotes_HaltSECTradingSuspension),
		Entry("Value is string, HaltRegulatoryConcern - Works",
			&types.AttributeValueMemberS{Value: "HaltRegulatoryConcern"}, Financial_Quotes_HaltRegulatoryConcern),
		Entry("Value is string, HaltMarketOperations - Works",
			&types.AttributeValueMemberS{Value: "HaltMarketOperations"}, Financial_Quotes_HaltMarketOperations),
		Entry("Value is string, IPOSecurityNotYetTrading - Works",
			&types.AttributeValueMemberS{Value: "IPOSecurityNotYetTrading"}, Financial_Quotes_IPOSecurityNotYetTrading),
		Entry("Value is string, HaltCorporateAction - Works",
			&types.AttributeValueMemberS{Value: "HaltCorporateAction"}, Financial_Quotes_HaltCorporateAction),
		Entry("Value is string, QuotationNotAvailable - Works",
			&types.AttributeValueMemberS{Value: "QuotationNotAvailable"}, Financial_Quotes_QuotationNotAvailable),
		Entry("Value is string, HaltVolatilityTradingPause - Works",
			&types.AttributeValueMemberS{Value: "HaltVolatilityTradingPause"}, Financial_Quotes_HaltVolatilityTradingPause),
		Entry("Value is string, HaltVolatilityTradingPauseStraddleCondition - Works",
			&types.AttributeValueMemberS{Value: "HaltVolatilityTradingPauseStraddleCondition"}, Financial_Quotes_HaltVolatilityTradingPauseStraddleCondition),
		Entry("Value is string, UpdateNewsAndResumptionTimes - Works",
			&types.AttributeValueMemberS{Value: "UpdateNewsAndResumptionTimes"}, Financial_Quotes_UpdateNewsAndResumptionTimes),
		Entry("Value is string, HaltSingleStockTradingPauseQuotesOnly - Works",
			&types.AttributeValueMemberS{Value: "HaltSingleStockTradingPauseQuotesOnly"}, Financial_Quotes_HaltSingleStockTradingPauseQuotesOnly),
		Entry("Value is string, ResumeQualificationIssuesReviewedResolved - Works",
			&types.AttributeValueMemberS{Value: "ResumeQualificationIssuesReviewedResolved"}, Financial_Quotes_ResumeQualificationIssuesReviewedResolved),
		Entry("Value is string, ResumeFilingRequirementsSatisfiedResolved - Works",
			&types.AttributeValueMemberS{Value: "ResumeFilingRequirementsSatisfiedResolved"}, Financial_Quotes_ResumeFilingRequirementsSatisfiedResolved),
		Entry("Value is string, ResumeNewsNotForthcoming - Works",
			&types.AttributeValueMemberS{Value: "ResumeNewsNotForthcoming"}, Financial_Quotes_ResumeNewsNotForthcoming),
		Entry("Value is string, ResumeQualificationsMaintRequirementsMet - Works",
			&types.AttributeValueMemberS{Value: "ResumeQualificationsMaintRequirementsMet"}, Financial_Quotes_ResumeQualificationsMaintRequirementsMet),
		Entry("Value is string, ResumeQualificationsFilingsMet - Works",
			&types.AttributeValueMemberS{Value: "ResumeQualificationsFilingsMet"}, Financial_Quotes_ResumeQualificationsFilingsMet),
		Entry("Value is string, ResumeRegulatoryAuth - Works",
			&types.AttributeValueMemberS{Value: "ResumeRegulatoryAuth"}, Financial_Quotes_ResumeRegulatoryAuth),
		Entry("Value is string, NewIssueAvailable - Works",
			&types.AttributeValueMemberS{Value: "NewIssueAvailable"}, Financial_Quotes_NewIssueAvailable),
		Entry("Value is string, IssueAvailable - Works",
			&types.AttributeValueMemberS{Value: "IssueAvailable"}, Financial_Quotes_IssueAvailable),
		Entry("Value is string, MWCBCarryFromPreviousDay - Works",
			&types.AttributeValueMemberS{Value: "MWCBCarryFromPreviousDay"}, Financial_Quotes_MWCBCarryFromPreviousDay),
		Entry("Value is string, MWCBResume - Works",
			&types.AttributeValueMemberS{Value: "MWCBResume"}, Financial_Quotes_MWCBResume),
		Entry("Value is string, IPOSecurityReleasedForQuotation - Works",
			&types.AttributeValueMemberS{Value: "IPOSecurityReleasedForQuotation"}, Financial_Quotes_IPOSecurityReleasedForQuotation),
		Entry("Value is string, IPOSecurityPositioningWindowExtension - Works",
			&types.AttributeValueMemberS{Value: "IPOSecurityPositioningWindowExtension"}, Financial_Quotes_IPOSecurityPositioningWindowExtension),
		Entry("Value is string, MWCBLevel1 - Works",
			&types.AttributeValueMemberS{Value: "MWCBLevel1"}, Financial_Quotes_MWCBLevel1),
		Entry("Value is string, MWCBLevel2 - Works",
			&types.AttributeValueMemberS{Value: "MWCBLevel2"}, Financial_Quotes_MWCBLevel2),
		Entry("Value is string, MWCBLevel3 - Works",
			&types.AttributeValueMemberS{Value: "MWCBLevel3"}, Financial_Quotes_MWCBLevel3),
		Entry("Value is string, HaltSubPennyTrading - Works",
			&types.AttributeValueMemberS{Value: "HaltSubPennyTrading"}, Financial_Quotes_HaltSubPennyTrading),
		Entry("Value is string, OrderImbalanceInd - Works",
			&types.AttributeValueMemberS{Value: "OrderImbalanceInd"}, Financial_Quotes_OrderImbalanceInd),
		Entry("Value is string, LULDTradingPaused - Works",
			&types.AttributeValueMemberS{Value: "LULDTradingPaused"}, Financial_Quotes_LULDTradingPaused),
		Entry("Value is string, NONE - Works",
			&types.AttributeValueMemberS{Value: "NONE"}, Financial_Quotes_NONE),
		Entry("Value is string, ShortSalesRestrictionActivated - Works",
			&types.AttributeValueMemberS{Value: "ShortSalesRestrictionActivated"}, Financial_Quotes_ShortSalesRestrictionActivated),
		Entry("Value is string, ShortSalesRestrictionContinued - Works",
			&types.AttributeValueMemberS{Value: "ShortSalesRestrictionContinued"}, Financial_Quotes_ShortSalesRestrictionContinued),
		Entry("Value is string, ShortSalesRestrictionDeactivated - Works",
			&types.AttributeValueMemberS{Value: "ShortSalesRestrictionDeactivated"}, Financial_Quotes_ShortSalesRestrictionDeactivated),
		Entry("Value is string, ShortSalesRestrictionInEffect - Works",
			&types.AttributeValueMemberS{Value: "ShortSalesRestrictionInEffect"}, Financial_Quotes_ShortSalesRestrictionInEffect),
		Entry("Value is string, ShortSalesRestrictionMax - Works",
			&types.AttributeValueMemberS{Value: "ShortSalesRestrictionMax"}, Financial_Quotes_ShortSalesRestrictionMax),
		Entry("Value is string, NBBONoChange - Works",
			&types.AttributeValueMemberS{Value: "NBBONoChange"}, Financial_Quotes_NBBONoChange),
		Entry("Value is string, NBBOQuoteIsNBBO - Works",
			&types.AttributeValueMemberS{Value: "NBBOQuoteIsNBBO"}, Financial_Quotes_NBBOQuoteIsNBBO),
		Entry("Value is string, NBBONoBBNoBO - Works",
			&types.AttributeValueMemberS{Value: "NBBONoBBNoBO"}, Financial_Quotes_NBBONoBBNoBO),
		Entry("Value is string, NBBOBBBOShortAppendage - Works",
			&types.AttributeValueMemberS{Value: "NBBOBBBOShortAppendage"}, Financial_Quotes_NBBOBBBOShortAppendage),
		Entry("Value is string, NBBOBBBOLongAppendage - Works",
			&types.AttributeValueMemberS{Value: "NBBOBBBOLongAppendage"}, Financial_Quotes_NBBOBBBOLongAppendage),
		Entry("Value is string, HeldTradeNotLastSaleNotConsolidated - Works",
			&types.AttributeValueMemberS{Value: "HeldTradeNotLastSaleNotConsolidated"}, Financial_Quotes_HeldTradeNotLastSaleNotConsolidated),
		Entry("Value is string, HeldTradeLastSaleButNotConsolidated - Works",
			&types.AttributeValueMemberS{Value: "HeldTradeLastSaleButNotConsolidated"}, Financial_Quotes_HeldTradeLastSaleButNotConsolidated),
		Entry("Value is string, HeldTradeLastSaleAndConsolidated - Works",
			&types.AttributeValueMemberS{Value: "HeldTradeLastSaleAndConsolidated"}, Financial_Quotes_HeldTradeLastSaleAndConsolidated),
		Entry("Value is string, RetailInterestOnBid - Works",
			&types.AttributeValueMemberS{Value: "RetailInterestOnBid"}, Financial_Quotes_RetailInterestOnBid),
		Entry("Value is string, RetailInterestOnAsk - Works",
			&types.AttributeValueMemberS{Value: "RetailInterestOnAsk"}, Financial_Quotes_RetailInterestOnAsk),
		Entry("Value is string, RetailInterestOnBidAndAsk - Works",
			&types.AttributeValueMemberS{Value: "RetailInterestOnBidAndAsk"}, Financial_Quotes_RetailInterestOnBidAndAsk),
		Entry("Value is string, FinraBBONoChange - Works",
			&types.AttributeValueMemberS{Value: "FinraBBONoChange"}, Financial_Quotes_FinraBBONoChange),
		Entry("Value is string, FinraBBODoesNotExist - Works",
			&types.AttributeValueMemberS{Value: "FinraBBODoesNotExist"}, Financial_Quotes_FinraBBODoesNotExist),
		Entry("Value is string, FinraBBBOExecutable - Works",
			&types.AttributeValueMemberS{Value: "FinraBBBOExecutable"}, Financial_Quotes_FinraBBBOExecutable),
		Entry("Value is string, FinraBBBelowLowerBand - Works",
			&types.AttributeValueMemberS{Value: "FinraBBBelowLowerBand"}, Financial_Quotes_FinraBBBelowLowerBand),
		Entry("Value is string, FinraBOAboveUpperBand - Works",
			&types.AttributeValueMemberS{Value: "FinraBOAboveUpperBand"}, Financial_Quotes_FinraBOAboveUpperBand),
		Entry("Value is string, FinraBBBelowLowerBandBOAbboveUpperBand - Works",
			&types.AttributeValueMemberS{Value: "FinraBBBelowLowerBandBOAbboveUpperBand"}, Financial_Quotes_FinraBBBelowLowerBandBOAbboveUpperBand),
		Entry("Value is string, CTANotDueToRelatedSecurity - Works",
			&types.AttributeValueMemberS{Value: "CTANotDueToRelatedSecurity"}, Financial_Quotes_CTANotDueToRelatedSecurity),
		Entry("Value is string, CTADueToRelatedSecurity - Works",
			&types.AttributeValueMemberS{Value: "CTADueToRelatedSecurity"}, Financial_Quotes_CTADueToRelatedSecurity),
		Entry("Value is string, CTANotInViewOfCommon - Works",
			&types.AttributeValueMemberS{Value: "CTANotInViewOfCommon"}, Financial_Quotes_CTANotInViewOfCommon),
		Entry("Value is string, CTAInViewOfCommon - Works",
			&types.AttributeValueMemberS{Value: "CTAInViewOfCommon"}, Financial_Quotes_CTAInViewOfCommon),
		Entry("Value is string, CTAPriceIndicator - Works",
			&types.AttributeValueMemberS{Value: "CTAPriceIndicator"}, Financial_Quotes_CTAPriceIndicator),
		Entry("Value is string, CTANewPriceIndicator - Works",
			&types.AttributeValueMemberS{Value: "CTANewPriceIndicator"}, Financial_Quotes_CTANewPriceIndicator),
		Entry("Value is string, CTACorrectedPriceIndication - Works",
			&types.AttributeValueMemberS{Value: "CTACorrectedPriceIndication"}, Financial_Quotes_CTACorrectedPriceIndication),
		Entry("Value is string, CTACancelledMarketImbalance - Works",
			&types.AttributeValueMemberS{Value: "CTACancelledMarketImbalance"}, Financial_Quotes_CTACancelledMarketImbalance))

	// Test that attempting to deserialize a Financial.Quotes.Indicator will fial and return an
	// error if the value cannot be converted to either the name value or integer value of the enum option
	It("Scan - Value is nil - Error", func() {

		// Attempt to convert a fake string value into a Financial.Quotes.Indicator
		// This should return an error
		var enum *Financial_Quotes_Indicator
		err := enum.Scan(nil)

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of %!q(<nil>) had an invalid type of <nil>"))
		Expect(enum).Should(BeNil())
	})

	// Test the conditions under which values should be convertible to a Financial.Quotes.Indicator
	DescribeTable("Scan Tests",
		func(value interface{}, shouldBe Financial_Quotes_Indicator) {

			// Attempt to convert the value into a Financial.Quotes.Indicator
			// This should not fail
			var enum Financial_Quotes_Indicator
			err := enum.Scan(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("NBBNBOExecutable - Works", "NBBNBOExecutable", Financial_Quotes_NBBNBOExecutable),
		Entry("NBBBelowLowerBand - Works", "NBBBelowLowerBand", Financial_Quotes_NBBBelowLowerBand),
		Entry("NBOAboveUpperBand - Works", "NBOAboveUpperBand", Financial_Quotes_NBOAboveUpperBand),
		Entry("NBBBelowLowerBandAndNBOAboveUpperBand - Works",
			"NBBBelowLowerBandAndNBOAboveUpperBand", Financial_Quotes_NBBBelowLowerBandAndNBOAboveUpperBand),
		Entry("NBBEqualsUpperBand - Works", "NBBEqualsUpperBand", Financial_Quotes_NBBEqualsUpperBand),
		Entry("NBOEqualsLowerBand - Works", "NBOEqualsLowerBand", Financial_Quotes_NBOEqualsLowerBand),
		Entry("NBBEqualsUpperBandAndNBOAboveUpperBand - Works",
			"NBBEqualsUpperBandAndNBOAboveUpperBand", Financial_Quotes_NBBEqualsUpperBandAndNBOAboveUpperBand),
		Entry("NBBBelowLowerBandAndNBOEqualsLowerBand - Works",
			"NBBBelowLowerBandAndNBOEqualsLowerBand", Financial_Quotes_NBBBelowLowerBandAndNBOEqualsLowerBand),
		Entry("BidPriceAboveUpperLimitPriceBand - Works",
			"BidPriceAboveUpperLimitPriceBand", Financial_Quotes_BidPriceAboveUpperLimitPriceBand),
		Entry("OfferPriceBelowLowerLimitPriceBand - Works",
			"OfferPriceBelowLowerLimitPriceBand", Financial_Quotes_OfferPriceBelowLowerLimitPriceBand),
		Entry("BidAndOfferOutsidePriceBand - Works",
			"BidAndOfferOutsidePriceBand", Financial_Quotes_BidAndOfferOutsidePriceBand),
		Entry("OpeningUpdate - Works", "OpeningUpdate", Financial_Quotes_OpeningUpdate),
		Entry("IntraDayUpdate - Works", "IntraDayUpdate", Financial_Quotes_IntraDayUpdate),
		Entry("RestatedValue - Works", "RestatedValue", Financial_Quotes_RestatedValue),
		Entry("SuspendedDuringTradingHalt - Works", "SuspendedDuringTradingHalt", Financial_Quotes_SuspendedDuringTradingHalt),
		Entry("ReOpeningUpdate - Works", "ReOpeningUpdate", Financial_Quotes_ReOpeningUpdate),
		Entry("OutsidePriceBandRuleHours - Works", "OutsidePriceBandRuleHours", Financial_Quotes_OutsidePriceBandRuleHours),
		Entry("AuctionExtension - Works", "AuctionExtension", Financial_Quotes_AuctionExtension),
		Entry("LULDPriceBand - Works", "LULDPriceBand", Financial_Quotes_LULDPriceBand),
		Entry("RepublishedLULDPriceBandInd - Works",
			"RepublishedLULDPriceBandInd", Financial_Quotes_RepublishedLULDPriceBandInd),
		Entry("NBBLimitStateEntered - Works", "NBBLimitStateEntered", Financial_Quotes_NBBLimitStateEntered),
		Entry("NBBLimitStateExited - Works", "NBBLimitStateExited", Financial_Quotes_NBBLimitStateExited),
		Entry("NBOLimitStateEntered - Works", "NBOLimitStateEntered", Financial_Quotes_NBOLimitStateEntered),
		Entry("NBOLimitStateExited - Works", "NBOLimitStateExited", Financial_Quotes_NBOLimitStateExited),
		Entry("NBBAndNBOLimitStateEntered - Works", "NBBAndNBOLimitStateEntered", Financial_Quotes_NBBAndNBOLimitStateEntered),
		Entry("NBBAndNBOLimitStateExited - Works", "NBBAndNBOLimitStateExited", Financial_Quotes_NBBAndNBOLimitStateExited),
		Entry("NBBLimitStateEnteredNBOLimitStateExited - Works",
			"NBBLimitStateEnteredNBOLimitStateExited", Financial_Quotes_NBBLimitStateEnteredNBOLimitStateExited),
		Entry("NBBLimitStateExitedNBOLimitStateEntered - Works",
			"NBBLimitStateExitedNBOLimitStateEntered", Financial_Quotes_NBBLimitStateExitedNBOLimitStateEntered),
		Entry("Normal - Works", "Normal", Financial_Quotes_Normal),
		Entry("Bankrupt - Works", "Bankrupt", Financial_Quotes_Bankrupt),
		Entry("Deficient - Works", "Deficient", Financial_Quotes_Deficient),
		Entry("Delinquent - Works", "Delinquent", Financial_Quotes_Delinquent),
		Entry("BankruptAndDeficient - Works", "BankruptAndDeficient", Financial_Quotes_BankruptAndDeficient),
		Entry("BankruptAndDelinquent - Works", "BankruptAndDelinquent", Financial_Quotes_BankruptAndDelinquent),
		Entry("DeficientAndDelinquent - Works", "DeficientAndDelinquent", Financial_Quotes_DeficientAndDelinquent),
		Entry("DeficientDeliquentBankrupt - Works", "DeficientDeliquentBankrupt", Financial_Quotes_DeficientDeliquentBankrupt),
		Entry("Liquidation - Works", "Liquidation", Financial_Quotes_Liquidation),
		Entry("CreationsSuspended - Works", "CreationsSuspended", Financial_Quotes_CreationsSuspended),
		Entry("RedemptionsSuspended - Works", "RedemptionsSuspended", Financial_Quotes_RedemptionsSuspended),
		Entry("CreationsRedemptionsSuspended - Works",
			"CreationsRedemptionsSuspended", Financial_Quotes_CreationsRedemptionsSuspended),
		Entry("NormalTrading - Works", "NormalTrading", Financial_Quotes_NormalTrading),
		Entry("OpeningDelay - Works", "OpeningDelay", Financial_Quotes_OpeningDelay),
		Entry("TradingHalt - Works", "TradingHalt", Financial_Quotes_TradingHalt),
		Entry("TradingResume - Works", "TradingResume", Financial_Quotes_TradingResume),
		Entry("NoOpenNoResume - Works", "NoOpenNoResume", Financial_Quotes_NoOpenNoResume),
		Entry("PriceIndication - Works", "PriceIndication", Financial_Quotes_PriceIndication),
		Entry("TradingRangeIndication - Works", "TradingRangeIndication", Financial_Quotes_TradingRangeIndication),
		Entry("MarketImbalanceBuy - Works", "MarketImbalanceBuy", Financial_Quotes_MarketImbalanceBuy),
		Entry("MarketImbalanceSell - Works", "MarketImbalanceSell", Financial_Quotes_MarketImbalanceSell),
		Entry("MarketOnCloseImbalanceBuy - Works", "MarketOnCloseImbalanceBuy", Financial_Quotes_MarketOnCloseImbalanceBuy),
		Entry("MarketOnCloseImbalanceSell - Works", "MarketOnCloseImbalanceSell", Financial_Quotes_MarketOnCloseImbalanceSell),
		Entry("NoMarketImbalance - Works", "NoMarketImbalance", Financial_Quotes_NoMarketImbalance),
		Entry("NoMarketOnCloseImbalance - Works", "NoMarketOnCloseImbalance", Financial_Quotes_NoMarketOnCloseImbalance),
		Entry("ShortSaleRestriction - Works", "ShortSaleRestriction", Financial_Quotes_ShortSaleRestriction),
		Entry("LimitUpLimitDown - Works", "LimitUpLimitDown", Financial_Quotes_LimitUpLimitDown),
		Entry("QuotationResumption - Works", "QuotationResumption", Financial_Quotes_QuotationResumption),
		Entry("TradingResumption - Works", "TradingResumption", Financial_Quotes_TradingResumption),
		Entry("VolatilityTradingPause - Works", "VolatilityTradingPause", Financial_Quotes_VolatilityTradingPause),
		Entry("PolygonReserved - Works", "PolygonReserved", Financial_Quotes_PolygonReserved),
		Entry("HaltNewsPending - Works", "HaltNewsPending", Financial_Quotes_HaltNewsPending),
		Entry("UpdateNewsDissemination - Works", "UpdateNewsDissemination", Financial_Quotes_UpdateNewsDissemination),
		Entry("HaltSingleStockTradingPause - Works",
			"HaltSingleStockTradingPause", Financial_Quotes_HaltSingleStockTradingPause),
		Entry("HaltRegulatoryExtraordinaryMarketActivity - Works",
			"HaltRegulatoryExtraordinaryMarketActivity", Financial_Quotes_HaltRegulatoryExtraordinaryMarketActivity),
		Entry("HaltETF - Works", "HaltETF", Financial_Quotes_HaltETF),
		Entry("HaltInformationRequested - Works", "HaltInformationRequested", Financial_Quotes_HaltInformationRequested),
		Entry("HaltExchangeNonCompliance - Works", "HaltExchangeNonCompliance", Financial_Quotes_HaltExchangeNonCompliance),
		Entry("HaltFilingsNotCurrent - Works", "HaltFilingsNotCurrent", Financial_Quotes_HaltFilingsNotCurrent),
		Entry("HaltSECTradingSuspension - Works", "HaltSECTradingSuspension", Financial_Quotes_HaltSECTradingSuspension),
		Entry("HaltRegulatoryConcern - Works", "HaltRegulatoryConcern", Financial_Quotes_HaltRegulatoryConcern),
		Entry("HaltMarketOperations - Works", "HaltMarketOperations", Financial_Quotes_HaltMarketOperations),
		Entry("IPOSecurityNotYetTrading - Works", "IPOSecurityNotYetTrading", Financial_Quotes_IPOSecurityNotYetTrading),
		Entry("HaltCorporateAction - Works", "HaltCorporateAction", Financial_Quotes_HaltCorporateAction),
		Entry("QuotationNotAvailable - Works", "QuotationNotAvailable", Financial_Quotes_QuotationNotAvailable),
		Entry("HaltVolatilityTradingPause - Works", "HaltVolatilityTradingPause", Financial_Quotes_HaltVolatilityTradingPause),
		Entry("HaltVolatilityTradingPauseStraddleCondition - Works",
			"HaltVolatilityTradingPauseStraddleCondition", Financial_Quotes_HaltVolatilityTradingPauseStraddleCondition),
		Entry("UpdateNewsAndResumptionTimes - Works",
			"UpdateNewsAndResumptionTimes", Financial_Quotes_UpdateNewsAndResumptionTimes),
		Entry("HaltSingleStockTradingPauseQuotesOnly - Works",
			"HaltSingleStockTradingPauseQuotesOnly", Financial_Quotes_HaltSingleStockTradingPauseQuotesOnly),
		Entry("ResumeQualificationIssuesReviewedResolved - Works",
			"ResumeQualificationIssuesReviewedResolved", Financial_Quotes_ResumeQualificationIssuesReviewedResolved),
		Entry("ResumeFilingRequirementsSatisfiedResolved - Works",
			"ResumeFilingRequirementsSatisfiedResolved", Financial_Quotes_ResumeFilingRequirementsSatisfiedResolved),
		Entry("ResumeNewsNotForthcoming - Works", "ResumeNewsNotForthcoming", Financial_Quotes_ResumeNewsNotForthcoming),
		Entry("ResumeQualificationsMaintRequirementsMet - Works",
			"ResumeQualificationsMaintRequirementsMet", Financial_Quotes_ResumeQualificationsMaintRequirementsMet),
		Entry("ResumeQualificationsFilingsMet - Works",
			"ResumeQualificationsFilingsMet", Financial_Quotes_ResumeQualificationsFilingsMet),
		Entry("ResumeRegulatoryAuth - Works", "ResumeRegulatoryAuth", Financial_Quotes_ResumeRegulatoryAuth),
		Entry("NewIssueAvailable - Works", "NewIssueAvailable", Financial_Quotes_NewIssueAvailable),
		Entry("IssueAvailable - Works", "IssueAvailable", Financial_Quotes_IssueAvailable),
		Entry("MWCBCarryFromPreviousDay - Works", "MWCBCarryFromPreviousDay", Financial_Quotes_MWCBCarryFromPreviousDay),
		Entry("MWCBResume - Works", "MWCBResume", Financial_Quotes_MWCBResume),
		Entry("IPOSecurityReleasedForQuotation - Works",
			"IPOSecurityReleasedForQuotation", Financial_Quotes_IPOSecurityReleasedForQuotation),
		Entry("IPOSecurityPositioningWindowExtension - Works",
			"IPOSecurityPositioningWindowExtension", Financial_Quotes_IPOSecurityPositioningWindowExtension),
		Entry("MWCBLevel1 - Works", "MWCBLevel1", Financial_Quotes_MWCBLevel1),
		Entry("MWCBLevel2 - Works", "MWCBLevel2", Financial_Quotes_MWCBLevel2),
		Entry("MWCBLevel3 - Works", "MWCBLevel3", Financial_Quotes_MWCBLevel3),
		Entry("HaltSubPennyTrading - Works", "HaltSubPennyTrading", Financial_Quotes_HaltSubPennyTrading),
		Entry("OrderImbalanceInd - Works", "OrderImbalanceInd", Financial_Quotes_OrderImbalanceInd),
		Entry("LULDTradingPaused - Works", "LULDTradingPaused", Financial_Quotes_LULDTradingPaused),
		Entry("NONE - Works", "NONE", Financial_Quotes_NONE),
		Entry("ShortSalesRestrictionActivated - Works",
			"ShortSalesRestrictionActivated", Financial_Quotes_ShortSalesRestrictionActivated),
		Entry("ShortSalesRestrictionContinued - Works",
			"ShortSalesRestrictionContinued", Financial_Quotes_ShortSalesRestrictionContinued),
		Entry("ShortSalesRestrictionDeactivated - Works",
			"ShortSalesRestrictionDeactivated", Financial_Quotes_ShortSalesRestrictionDeactivated),
		Entry("ShortSalesRestrictionInEffect - Works",
			"ShortSalesRestrictionInEffect", Financial_Quotes_ShortSalesRestrictionInEffect),
		Entry("ShortSalesRestrictionMax - Works", "ShortSalesRestrictionMax", Financial_Quotes_ShortSalesRestrictionMax),
		Entry("RetailInterestOnBid - Works", "RetailInterestOnBid", Financial_Quotes_RetailInterestOnBid),
		Entry("RetailInterestOnAsk - Works", "RetailInterestOnAsk", Financial_Quotes_RetailInterestOnAsk),
		Entry("RetailInterestOnBidAndAsk - Works",
			"RetailInterestOnBidAndAsk", Financial_Quotes_RetailInterestOnBidAndAsk),
		Entry("FinraBBONoChange - Works", "FinraBBONoChange", Financial_Quotes_FinraBBONoChange),
		Entry("FinraBBODoesNotExist - Works", "FinraBBODoesNotExist", Financial_Quotes_FinraBBODoesNotExist),
		Entry("FinraBBBOExecutable - Works", "FinraBBBOExecutable", Financial_Quotes_FinraBBBOExecutable),
		Entry("FinraBBBelowLowerBand - Works", "FinraBBBelowLowerBand", Financial_Quotes_FinraBBBelowLowerBand),
		Entry("FinraBOAboveUpperBand - Works", "FinraBOAboveUpperBand", Financial_Quotes_FinraBOAboveUpperBand),
		Entry("FinraBBBelowLowerBandBOAbboveUpperBand - Works",
			"FinraBBBelowLowerBandBOAbboveUpperBand", Financial_Quotes_FinraBBBelowLowerBandBOAbboveUpperBand),
		Entry("NBBONoChange - Works", "NBBONoChange", Financial_Quotes_NBBONoChange),
		Entry("NBBOQuoteIsNBBO - Works", "NBBOQuoteIsNBBO", Financial_Quotes_NBBOQuoteIsNBBO),
		Entry("NBBONoBBNoBO - Works", "NBBONoBBNoBO", Financial_Quotes_NBBONoBBNoBO),
		Entry("NBBOBBBOShortAppendage - Works", "NBBOBBBOShortAppendage", Financial_Quotes_NBBOBBBOShortAppendage),
		Entry("NBBOBBBOLongAppendage - Works", "NBBOBBBOLongAppendage", Financial_Quotes_NBBOBBBOLongAppendage),
		Entry("HeldTradeNotLastSaleNotConsolidated - Works",
			"HeldTradeNotLastSaleNotConsolidated", Financial_Quotes_HeldTradeNotLastSaleNotConsolidated),
		Entry("HeldTradeLastSaleButNotConsolidated - Works",
			"HeldTradeLastSaleButNotConsolidated", Financial_Quotes_HeldTradeLastSaleButNotConsolidated),
		Entry("HeldTradeLastSaleAndConsolidated - Works",
			"HeldTradeLastSaleAndConsolidated", Financial_Quotes_HeldTradeLastSaleAndConsolidated),
		Entry("CTANotDueToRelatedSecurity - Works", "CTANotDueToRelatedSecurity", Financial_Quotes_CTANotDueToRelatedSecurity),
		Entry("CTADueToRelatedSecurity - Works", "CTADueToRelatedSecurity", Financial_Quotes_CTADueToRelatedSecurity),
		Entry("CTANotInViewOfCommon - Works", "CTANotInViewOfCommon", Financial_Quotes_CTANotInViewOfCommon),
		Entry("CTAInViewOfCommon - Works", "CTAInViewOfCommon", Financial_Quotes_CTAInViewOfCommon),
		Entry("CTAPriceIndicator - Works", "CTAPriceIndicator", Financial_Quotes_CTAPriceIndicator),
		Entry("CTANewPriceIndicator - Works", "CTANewPriceIndicator", Financial_Quotes_CTANewPriceIndicator),
		Entry("CTACorrectedPriceIndication - Works",
			"CTACorrectedPriceIndication", Financial_Quotes_CTACorrectedPriceIndication),
		Entry("CTACancelledMarketImbalance - Works",
			"CTACancelledMarketImbalance", Financial_Quotes_CTACancelledMarketImbalance),
		Entry("0 - Works", 0, Financial_Quotes_NBBNBOExecutable),
		Entry("1 - Works", 1, Financial_Quotes_NBBBelowLowerBand),
		Entry("2 - Works", 2, Financial_Quotes_NBOAboveUpperBand),
		Entry("3 - Works", 3, Financial_Quotes_NBBBelowLowerBandAndNBOAboveUpperBand),
		Entry("4 - Works", 4, Financial_Quotes_NBBEqualsUpperBand),
		Entry("5 - Works", 5, Financial_Quotes_NBOEqualsLowerBand),
		Entry("6 - Works", 6, Financial_Quotes_NBBEqualsUpperBandAndNBOAboveUpperBand),
		Entry("7 - Works", 7, Financial_Quotes_NBBBelowLowerBandAndNBOEqualsLowerBand),
		Entry("8 - Works", 8, Financial_Quotes_BidPriceAboveUpperLimitPriceBand),
		Entry("9 - Works", 9, Financial_Quotes_OfferPriceBelowLowerLimitPriceBand),
		Entry("10 - Works", 10, Financial_Quotes_BidAndOfferOutsidePriceBand),
		Entry("11 - Works", 11, Financial_Quotes_OpeningUpdate),
		Entry("12 - Works", 12, Financial_Quotes_IntraDayUpdate),
		Entry("13 - Works", 13, Financial_Quotes_RestatedValue),
		Entry("14 - Works", 14, Financial_Quotes_SuspendedDuringTradingHalt),
		Entry("15 - Works", 15, Financial_Quotes_ReOpeningUpdate),
		Entry("16 - Works", 16, Financial_Quotes_OutsidePriceBandRuleHours),
		Entry("17 - Works", 17, Financial_Quotes_AuctionExtension),
		Entry("18 - Works", 18, Financial_Quotes_LULDPriceBand),
		Entry("19 - Works", 19, Financial_Quotes_RepublishedLULDPriceBandInd),
		Entry("20 - Works", 20, Financial_Quotes_NBBLimitStateEntered),
		Entry("21 - Works", 21, Financial_Quotes_NBBLimitStateExited),
		Entry("22 - Works", 22, Financial_Quotes_NBOLimitStateEntered),
		Entry("23 - Works", 23, Financial_Quotes_NBOLimitStateExited),
		Entry("24 - Works", 24, Financial_Quotes_NBBAndNBOLimitStateEntered),
		Entry("25 - Works", 25, Financial_Quotes_NBBAndNBOLimitStateExited),
		Entry("26 - Works", 26, Financial_Quotes_NBBLimitStateEnteredNBOLimitStateExited),
		Entry("27 - Works", 27, Financial_Quotes_NBBLimitStateExitedNBOLimitStateEntered),
		Entry("28 - Works", 28, Financial_Quotes_Normal),
		Entry("29 - Works", 29, Financial_Quotes_Bankrupt),
		Entry("30 - Works", 30, Financial_Quotes_Deficient),
		Entry("31 - Works", 31, Financial_Quotes_Delinquent),
		Entry("32 - Works", 32, Financial_Quotes_BankruptAndDeficient),
		Entry("33 - Works", 33, Financial_Quotes_BankruptAndDelinquent),
		Entry("34 - Works", 34, Financial_Quotes_DeficientAndDelinquent),
		Entry("35 - Works", 35, Financial_Quotes_DeficientDeliquentBankrupt),
		Entry("36 - Works", 36, Financial_Quotes_Liquidation),
		Entry("37 - Works", 37, Financial_Quotes_CreationsSuspended),
		Entry("38 - Works", 38, Financial_Quotes_RedemptionsSuspended),
		Entry("39 - Works", 39, Financial_Quotes_CreationsRedemptionsSuspended),
		Entry("40 - Works", 40, Financial_Quotes_NormalTrading),
		Entry("41 - Works", 41, Financial_Quotes_OpeningDelay),
		Entry("42 - Works", 42, Financial_Quotes_TradingHalt),
		Entry("43 - Works", 43, Financial_Quotes_TradingResume),
		Entry("44 - Works", 44, Financial_Quotes_NoOpenNoResume),
		Entry("45 - Works", 45, Financial_Quotes_PriceIndication),
		Entry("46 - Works", 46, Financial_Quotes_TradingRangeIndication),
		Entry("47 - Works", 47, Financial_Quotes_MarketImbalanceBuy),
		Entry("48 - Works", 48, Financial_Quotes_MarketImbalanceSell),
		Entry("49 - Works", 49, Financial_Quotes_MarketOnCloseImbalanceBuy),
		Entry("50 - Works", 50, Financial_Quotes_MarketOnCloseImbalanceSell),
		Entry("51 - Works", 51, Financial_Quotes_NoMarketImbalance),
		Entry("52 - Works", 52, Financial_Quotes_NoMarketOnCloseImbalance),
		Entry("53 - Works", 53, Financial_Quotes_ShortSaleRestriction),
		Entry("54 - Works", 54, Financial_Quotes_LimitUpLimitDown),
		Entry("55 - Works", 55, Financial_Quotes_QuotationResumption),
		Entry("56 - Works", 56, Financial_Quotes_TradingResumption),
		Entry("57 - Works", 57, Financial_Quotes_VolatilityTradingPause),
		Entry("58 - Works", 58, Financial_Quotes_PolygonReserved),
		Entry("59 - Works", 59, Financial_Quotes_HaltNewsPending),
		Entry("60 - Works", 60, Financial_Quotes_UpdateNewsDissemination),
		Entry("61 - Works", 61, Financial_Quotes_HaltSingleStockTradingPause),
		Entry("62 - Works", 62, Financial_Quotes_HaltRegulatoryExtraordinaryMarketActivity),
		Entry("63 - Works", 63, Financial_Quotes_HaltETF),
		Entry("64 - Works", 64, Financial_Quotes_HaltInformationRequested),
		Entry("65 - Works", 65, Financial_Quotes_HaltExchangeNonCompliance),
		Entry("66 - Works", 66, Financial_Quotes_HaltFilingsNotCurrent),
		Entry("67 - Works", 67, Financial_Quotes_HaltSECTradingSuspension),
		Entry("68 - Works", 68, Financial_Quotes_HaltRegulatoryConcern),
		Entry("69 - Works", 69, Financial_Quotes_HaltMarketOperations),
		Entry("70 - Works", 70, Financial_Quotes_IPOSecurityNotYetTrading),
		Entry("71 - Works", 71, Financial_Quotes_HaltCorporateAction),
		Entry("72 - Works", 72, Financial_Quotes_QuotationNotAvailable),
		Entry("73 - Works", 73, Financial_Quotes_HaltVolatilityTradingPause),
		Entry("74 - Works", 74, Financial_Quotes_HaltVolatilityTradingPauseStraddleCondition),
		Entry("75 - Works", 75, Financial_Quotes_UpdateNewsAndResumptionTimes),
		Entry("76 - Works", 76, Financial_Quotes_HaltSingleStockTradingPauseQuotesOnly),
		Entry("77 - Works", 77, Financial_Quotes_ResumeQualificationIssuesReviewedResolved),
		Entry("78 - Works", 78, Financial_Quotes_ResumeFilingRequirementsSatisfiedResolved),
		Entry("79 - Works", 79, Financial_Quotes_ResumeNewsNotForthcoming),
		Entry("80 - Works", 80, Financial_Quotes_ResumeQualificationsMaintRequirementsMet),
		Entry("81 - Works", 81, Financial_Quotes_ResumeQualificationsFilingsMet),
		Entry("82 - Works", 82, Financial_Quotes_ResumeRegulatoryAuth),
		Entry("83 - Works", 83, Financial_Quotes_NewIssueAvailable),
		Entry("84 - Works", 84, Financial_Quotes_IssueAvailable),
		Entry("85 - Works", 85, Financial_Quotes_MWCBCarryFromPreviousDay),
		Entry("86 - Works", 86, Financial_Quotes_MWCBResume),
		Entry("87 - Works", 87, Financial_Quotes_IPOSecurityReleasedForQuotation),
		Entry("88 - Works", 88, Financial_Quotes_IPOSecurityPositioningWindowExtension),
		Entry("89 - Works", 89, Financial_Quotes_MWCBLevel1),
		Entry("90 - Works", 90, Financial_Quotes_MWCBLevel2),
		Entry("91 - Works", 91, Financial_Quotes_MWCBLevel3),
		Entry("92 - Works", 92, Financial_Quotes_HaltSubPennyTrading),
		Entry("93 - Works", 93, Financial_Quotes_OrderImbalanceInd),
		Entry("94 - Works", 94, Financial_Quotes_LULDTradingPaused),
		Entry("95 - Works", 95, Financial_Quotes_NONE),
		Entry("96 - Works", 96, Financial_Quotes_ShortSalesRestrictionActivated),
		Entry("97 - Works", 97, Financial_Quotes_ShortSalesRestrictionContinued),
		Entry("98 - Works", 98, Financial_Quotes_ShortSalesRestrictionDeactivated),
		Entry("99 - Works", 99, Financial_Quotes_ShortSalesRestrictionInEffect),
		Entry("100 - Works", 100, Financial_Quotes_ShortSalesRestrictionMax),
		Entry("101 - Works", 101, Financial_Quotes_NBBONoChange),
		Entry("102 - Works", 102, Financial_Quotes_NBBOQuoteIsNBBO),
		Entry("103 - Works", 103, Financial_Quotes_NBBONoBBNoBO),
		Entry("104 - Works", 104, Financial_Quotes_NBBOBBBOShortAppendage),
		Entry("105 - Works", 105, Financial_Quotes_NBBOBBBOLongAppendage),
		Entry("106 - Works", 106, Financial_Quotes_HeldTradeNotLastSaleNotConsolidated),
		Entry("107 - Works", 107, Financial_Quotes_HeldTradeLastSaleButNotConsolidated),
		Entry("108 - Works", 108, Financial_Quotes_HeldTradeLastSaleAndConsolidated),
		Entry("109 - Works", 109, Financial_Quotes_RetailInterestOnBid),
		Entry("110 - Works", 110, Financial_Quotes_RetailInterestOnAsk),
		Entry("111 - Works", 111, Financial_Quotes_RetailInterestOnBidAndAsk),
		Entry("112 - Works", 112, Financial_Quotes_FinraBBONoChange),
		Entry("113 - Works", 113, Financial_Quotes_FinraBBODoesNotExist),
		Entry("114 - Works", 114, Financial_Quotes_FinraBBBOExecutable),
		Entry("115 - Works", 115, Financial_Quotes_FinraBBBelowLowerBand),
		Entry("116 - Works", 116, Financial_Quotes_FinraBOAboveUpperBand),
		Entry("117 - Works", 117, Financial_Quotes_FinraBBBelowLowerBandBOAbboveUpperBand),
		Entry("118 - Works", 118, Financial_Quotes_CTANotDueToRelatedSecurity),
		Entry("119 - Works", 119, Financial_Quotes_CTADueToRelatedSecurity),
		Entry("120 - Works", 120, Financial_Quotes_CTANotInViewOfCommon),
		Entry("121 - Works", 121, Financial_Quotes_CTAInViewOfCommon),
		Entry("122 - Works", 122, Financial_Quotes_CTAPriceIndicator),
		Entry("123 - Works", 123, Financial_Quotes_CTANewPriceIndicator),
		Entry("124 - Works", 124, Financial_Quotes_CTACorrectedPriceIndication),
		Entry("125 - Works", 125, Financial_Quotes_CTACancelledMarketImbalance))
})

var _ = Describe("Financial.Trades.Condition Marshal/Unmarshal Tests", func() {

	// Test that converting the Financial.Trades.Condition enum to JSON works for all values
	DescribeTable("MarshalJSON Tests",
		func(enum Financial_Trades_Condition, value string) {
			data, err := json.Marshal(enum)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("RegularSale - Works", Financial_Trades_RegularSale, "\"RegularSale\""),
		Entry("Acquisition - Works", Financial_Trades_Acquisition, "\"Acquisition\""),
		Entry("AveragePriceTrade - Works", Financial_Trades_AveragePriceTrade, "\"AveragePriceTrade\""),
		Entry("AutomaticExecution - Works", Financial_Trades_AutomaticExecution, "\"AutomaticExecution\""),
		Entry("BunchedTrade - Works", Financial_Trades_BunchedTrade, "\"BunchedTrade\""),
		Entry("BunchedSoldTrade - Works", Financial_Trades_BunchedSoldTrade, "\"BunchedSoldTrade\""),
		Entry("CAPElection - Works", Financial_Trades_CAPElection, "\"CAPElection\""),
		Entry("CashSale - Works", Financial_Trades_CashSale, "\"CashSale\""),
		Entry("ClosingPrints - Works", Financial_Trades_ClosingPrints, "\"ClosingPrints\""),
		Entry("CrossTrade - Works", Financial_Trades_CrossTrade, "\"CrossTrade\""),
		Entry("DerivativelyPriced - Works", Financial_Trades_DerivativelyPriced, "\"DerivativelyPriced\""),
		Entry("Distribution - Works", Financial_Trades_Distribution, "\"Distribution\""),
		Entry("FormT - Works", Financial_Trades_FormT, "\"FormT\""),
		Entry("ExtendedTradingHours - Works", Financial_Trades_ExtendedTradingHours, "\"ExtendedTradingHours\""),
		Entry("IntermarketSweep - Works", Financial_Trades_IntermarketSweep, "\"IntermarketSweep\""),
		Entry("MarketCenterOfficialClose - Works", Financial_Trades_MarketCenterOfficialClose, "\"MarketCenterOfficialClose\""),
		Entry("MarketCenterOfficialOpen - Works", Financial_Trades_MarketCenterOfficialOpen, "\"MarketCenterOfficialOpen\""),
		Entry("MarketCenterOpeningTrade - Works", Financial_Trades_MarketCenterOpeningTrade, "\"MarketCenterOpeningTrade\""),
		Entry("MarketCenterReopeningTrade - Works",
			Financial_Trades_MarketCenterReopeningTrade, "\"MarketCenterReopeningTrade\""),
		Entry("MarketCenterClosingTrade - Works", Financial_Trades_MarketCenterClosingTrade, "\"MarketCenterClosingTrade\""),
		Entry("NextDay - Works", Financial_Trades_NextDay, "\"NextDay\""),
		Entry("PriceVariationTrade - Works", Financial_Trades_PriceVariationTrade, "\"PriceVariationTrade\""),
		Entry("PriorReferencePrice - Works", Financial_Trades_PriorReferencePrice, "\"PriorReferencePrice\""),
		Entry("Rule155Trade - Works", Financial_Trades_Rule155Trade, "\"Rule155Trade\""),
		Entry("Rule127NYSE - Works", Financial_Trades_Rule127NYSE, "\"Rule127NYSE\""),
		Entry("OpeningPrints - Works", Financial_Trades_OpeningPrints, "\"OpeningPrints\""),
		Entry("Opened - Works", Financial_Trades_Opened, "\"Opened\""),
		Entry("StoppedStock - Works", Financial_Trades_StoppedStock, "\"StoppedStock\""),
		Entry("ReOpeningPrints - Works", Financial_Trades_ReOpeningPrints, "\"ReOpeningPrints\""),
		Entry("Seller - Works", Financial_Trades_Seller, "\"Seller\""),
		Entry("SoldLast - Works", Financial_Trades_SoldLast, "\"SoldLast\""),
		Entry("SoldLastAndStoppedStock - Works", Financial_Trades_SoldLastAndStoppedStock, "\"SoldLastAndStoppedStock\""),
		Entry("SoldOut - Works", Financial_Trades_SoldOut, "\"SoldOut\""),
		Entry("SoldOutOfSequence - Works", Financial_Trades_SoldOutOfSequence, "\"SoldOutOfSequence\""),
		Entry("SplitTrade - Works", Financial_Trades_SplitTrade, "\"SplitTrade\""),
		Entry("StockOption - Works", Financial_Trades_StockOption, "\"StockOption\""),
		Entry("YellowFlagRegularTrade - Works", Financial_Trades_YellowFlagRegularTrade, "\"YellowFlagRegularTrade\""),
		Entry("OddLotTrade - Works", Financial_Trades_OddLotTrade, "\"OddLotTrade\""),
		Entry("CorrectedConsolidatedClose - Works",
			Financial_Trades_CorrectedConsolidatedClose, "\"CorrectedConsolidatedClose\""),
		Entry("Unknown - Works", Financial_Trades_Unknown, "\"Unknown\""),
		Entry("Held - Works", Financial_Trades_Held, "\"Held\""),
		Entry("TradeThruExempt - Works", Financial_Trades_TradeThruExempt, "\"TradeThruExempt\""),
		Entry("NonEligible - Works", Financial_Trades_NonEligible, "\"NonEligible\""),
		Entry("NonEligibleExtended - Works", Financial_Trades_NonEligibleExtended, "\"NonEligibleExtended\""),
		Entry("Cancelled - Works", Financial_Trades_Cancelled, "\"Cancelled\""),
		Entry("Recovery - Works", Financial_Trades_Recovery, "\"Recovery\""),
		Entry("Correction - Works", Financial_Trades_Correction, "\"Correction\""),
		Entry("AsOf - Works", Financial_Trades_AsOf, "\"AsOf\""),
		Entry("AsOfCorrection - Works", Financial_Trades_AsOfCorrection, "\"AsOfCorrection\""),
		Entry("AsOfCancel - Works", Financial_Trades_AsOfCancel, "\"AsOfCancel\""),
		Entry("OOB - Works", Financial_Trades_OOB, "\"OOB\""),
		Entry("Summary - Works", Financial_Trades_Summary, "\"Summary\""),
		Entry("ContingentTrade - Works", Financial_Trades_ContingentTrade, "\"ContingentTrade\""),
		Entry("QualifiedContingentTrade - Works", Financial_Trades_QualifiedContingentTrade, "\"QualifiedContingentTrade\""),
		Entry("Errored - Works", Financial_Trades_Errored, "\"Errored\""),
		Entry("OpeningReopeningTradeDetail - Works",
			Financial_Trades_OpeningReopeningTradeDetail, "\"OpeningReopeningTradeDetail\""),
		Entry("Placeholder - Works", Financial_Trades_Placeholder, "\"Placeholder\""),
		Entry("ShortSaleRestrictionActivated - Works",
			Financial_Trades_ShortSaleRestrictionActivated, "\"ShortSaleRestrictionActivated\""),
		Entry("ShortSaleRestrictionContinued - Works",
			Financial_Trades_ShortSaleRestrictionContinued, "\"ShortSaleRestrictionContinued\""),
		Entry("ShortSaleRestrictionDeactivated - Works",
			Financial_Trades_ShortSaleRestrictionDeactivated, "\"ShortSaleRestrictionDeactivated\""),
		Entry("ShortSaleRestrictionInEffect - Works",
			Financial_Trades_ShortSaleRestrictionInEffect, "\"ShortSaleRestrictionInEffect\""),
		Entry("FinancialStatusBankrupt - Works", Financial_Trades_FinancialStatusBankrupt, "\"FinancialStatusBankrupt\""),
		Entry("FinancialStatusDeficient - Works", Financial_Trades_FinancialStatusDeficient, "\"FinancialStatusDeficient\""),
		Entry("FinancialStatusDelinquent - Works", Financial_Trades_FinancialStatusDelinquent, "\"FinancialStatusDelinquent\""),
		Entry("FinancialStatusBankruptAndDeficient - Works",
			Financial_Trades_FinancialStatusBankruptAndDeficient, "\"FinancialStatusBankruptAndDeficient\""),
		Entry("FinancialStatusBankruptAndDelinquent - Works",
			Financial_Trades_FinancialStatusBankruptAndDelinquent, "\"FinancialStatusBankruptAndDelinquent\""),
		Entry("FinancialStatusDeficientAndDelinquent - Works",
			Financial_Trades_FinancialStatusDeficientAndDelinquent, "\"FinancialStatusDeficientAndDelinquent\""),
		Entry("FinancialStatusDeficientDelinquentBankrupt - Works",
			Financial_Trades_FinancialStatusDeficientDelinquentBankrupt, "\"FinancialStatusDeficientDelinquentBankrupt\""),
		Entry("FinancialStatusLiquidation - Works",
			Financial_Trades_FinancialStatusLiquidation, "\"FinancialStatusLiquidation\""),
		Entry("FinancialStatusCreationsSuspended - Works",
			Financial_Trades_FinancialStatusCreationsSuspended, "\"FinancialStatusCreationsSuspended\""),
		Entry("FinancialStatusRedemptionsSuspended - Works",
			Financial_Trades_FinancialStatusRedemptionsSuspended, "\"FinancialStatusRedemptionsSuspended\""),
		Entry("Canceled - Works", Financial_Trades_Canceled, "\"Canceled\""),
		Entry("LateAndOutOfSequence - Works", Financial_Trades_LateAndOutOfSequence, "\"LateAndOutOfSequence\""),
		Entry("LastAndCanceled - Works", Financial_Trades_LastAndCanceled, "\"LastAndCanceled\""),
		Entry("Late - Works", Financial_Trades_Late, "\"Late\""),
		Entry("OpeningTradeAndCanceled - Works", Financial_Trades_OpeningTradeAndCanceled, "\"OpeningTradeAndCanceled\""),
		Entry("OpeningTradeLateAndOutOfSequence - Works",
			Financial_Trades_OpeningTradeLateAndOutOfSequence, "\"OpeningTradeLateAndOutOfSequence\""),
		Entry("OnlyTradeAndCanceled - Works", Financial_Trades_OnlyTradeAndCanceled, "\"OnlyTradeAndCanceled\""),
		Entry("OpeningTradeAndLate - Works", Financial_Trades_OpeningTradeAndLate, "\"OpeningTradeAndLate\""),
		Entry("AutomaticExecutionOption - Works", Financial_Trades_AutomaticExecutionOption, "\"AutomaticExecutionOption\""),
		Entry("ReopeningTrade - Works", Financial_Trades_ReopeningTrade, "\"ReopeningTrade\""),
		Entry("IntermarketSweepOrder - Works", Financial_Trades_IntermarketSweepOrder, "\"IntermarketSweepOrder\""),
		Entry("SingleLegAuctionNonISO - Works", Financial_Trades_SingleLegAuctionNonISO, "\"SingleLegAuctionNonISO\""),
		Entry("SingleLegAuctionISO - Works", Financial_Trades_SingleLegAuctionISO, "\"SingleLegAuctionISO\""),
		Entry("SingleLegCrossNonISO - Works", Financial_Trades_SingleLegCrossNonISO, "\"SingleLegCrossNonISO\""),
		Entry("SingleLegCrossISO - Works", Financial_Trades_SingleLegCrossISO, "\"SingleLegCrossISO\""),
		Entry("SingleLegFloorTrade - Works", Financial_Trades_SingleLegFloorTrade, "\"SingleLegFloorTrade\""),
		Entry("MultiLegAutoElectronicTrade - Works",
			Financial_Trades_MultiLegAutoElectronicTrade, "\"MultiLegAutoElectronicTrade\""),
		Entry("MultiLegAuction - Works", Financial_Trades_MultiLegAuction, "\"MultiLegAuction\""),
		Entry("MultiLegCross - Works", Financial_Trades_MultiLegCross, "\"MultiLegCross\""),
		Entry("MultiLegFloorTrade - Works", Financial_Trades_MultiLegFloorTrade, "\"MultiLegFloorTrade\""),
		Entry("MultiLegAutoElectronicTradeAgainstSingleLeg - Works",
			Financial_Trades_MultiLegAutoElectronicTradeAgainstSingleLeg, "\"MultiLegAutoElectronicTradeAgainstSingleLeg\""),
		Entry("StockOptionsAuction - Works", Financial_Trades_StockOptionsAuction, "\"StockOptionsAuction\""),
		Entry("MultiLegAuctionAgainstSingleLeg - Works",
			Financial_Trades_MultiLegAuctionAgainstSingleLeg, "\"MultiLegAuctionAgainstSingleLeg\""),
		Entry("MultiLegFloorTradeAgainstSingleLeg - Works",
			Financial_Trades_MultiLegFloorTradeAgainstSingleLeg, "\"MultiLegFloorTradeAgainstSingleLeg\""),
		Entry("StockOptionsAutoElectronicTrade - Works",
			Financial_Trades_StockOptionsAutoElectronicTrade, "\"StockOptionsAutoElectronicTrade\""),
		Entry("StockOptionsCross - Works", Financial_Trades_StockOptionsCross, "\"StockOptionsCross\""),
		Entry("StockOptionsFloorTrade - Works", Financial_Trades_StockOptionsFloorTrade, "\"StockOptionsFloorTrade\""),
		Entry("StockOptionsAutoElectronicTradeAgainstSingleLeg - Works",
			Financial_Trades_StockOptionsAutoElectronicTradeAgainstSingleLeg, "\"StockOptionsAutoElectronicTradeAgainstSingleLeg\""),
		Entry("StockOptionsAuctionAgainstSingleLeg - Works",
			Financial_Trades_StockOptionsAuctionAgainstSingleLeg, "\"StockOptionsAuctionAgainstSingleLeg\""),
		Entry("StockOptionsFloorTradeAgainstSingleLeg - Works",
			Financial_Trades_StockOptionsFloorTradeAgainstSingleLeg, "\"StockOptionsFloorTradeAgainstSingleLeg\""),
		Entry("MultiLegFloorTradeOfProprietaryProducts - Works",
			Financial_Trades_MultiLegFloorTradeOfProprietaryProducts, "\"MultiLegFloorTradeOfProprietaryProducts\""),
		Entry("MultilateralCompressionTradeOfProprietaryProducts - Works",
			Financial_Trades_MultilateralCompressionTradeOfProprietaryProducts, "\"MultilateralCompressionTradeOfProprietaryProducts\""),
		Entry("ExtendedHoursTrade - Works", Financial_Trades_ExtendedHoursTrade, "\"ExtendedHoursTrade\""))

	// Test that converting the Financial.Trades.Condition enum to a CSV column works for all values
	DescribeTable("MarshalCSV Tests",
		func(enum Financial_Trades_Condition, value string) {
			data, err := enum.MarshalCSV()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("RegularSale - Works", Financial_Trades_RegularSale, "0"),
		Entry("Acquisition - Works", Financial_Trades_Acquisition, "1"),
		Entry("AveragePriceTrade - Works", Financial_Trades_AveragePriceTrade, "2"),
		Entry("AutomaticExecution - Works", Financial_Trades_AutomaticExecution, "3"),
		Entry("BunchedTrade - Works", Financial_Trades_BunchedTrade, "4"),
		Entry("BunchedSoldTrade - Works", Financial_Trades_BunchedSoldTrade, "5"),
		Entry("CAPElection - Works", Financial_Trades_CAPElection, "6"),
		Entry("CashSale - Works", Financial_Trades_CashSale, "7"),
		Entry("ClosingPrints - Works", Financial_Trades_ClosingPrints, "8"),
		Entry("CrossTrade - Works", Financial_Trades_CrossTrade, "9"),
		Entry("DerivativelyPriced - Works", Financial_Trades_DerivativelyPriced, "10"),
		Entry("Distribution - Works", Financial_Trades_Distribution, "11"),
		Entry("FormT - Works", Financial_Trades_FormT, "12"),
		Entry("ExtendedTradingHours - Works", Financial_Trades_ExtendedTradingHours, "13"),
		Entry("IntermarketSweep - Works", Financial_Trades_IntermarketSweep, "14"),
		Entry("MarketCenterOfficialClose - Works", Financial_Trades_MarketCenterOfficialClose, "15"),
		Entry("MarketCenterOfficialOpen - Works", Financial_Trades_MarketCenterOfficialOpen, "16"),
		Entry("MarketCenterOpeningTrade - Works", Financial_Trades_MarketCenterOpeningTrade, "17"),
		Entry("MarketCenterReopeningTrade - Works", Financial_Trades_MarketCenterReopeningTrade, "18"),
		Entry("MarketCenterClosingTrade - Works", Financial_Trades_MarketCenterClosingTrade, "19"),
		Entry("NextDay - Works", Financial_Trades_NextDay, "20"),
		Entry("PriceVariationTrade - Works", Financial_Trades_PriceVariationTrade, "21"),
		Entry("PriorReferencePrice - Works", Financial_Trades_PriorReferencePrice, "22"),
		Entry("Rule155Trade - Works", Financial_Trades_Rule155Trade, "23"),
		Entry("Rule127NYSE - Works", Financial_Trades_Rule127NYSE, "24"),
		Entry("OpeningPrints - Works", Financial_Trades_OpeningPrints, "25"),
		Entry("Opened - Works", Financial_Trades_Opened, "26"),
		Entry("StoppedStock - Works", Financial_Trades_StoppedStock, "27"),
		Entry("ReOpeningPrints - Works", Financial_Trades_ReOpeningPrints, "28"),
		Entry("Seller - Works", Financial_Trades_Seller, "29"),
		Entry("SoldLast - Works", Financial_Trades_SoldLast, "30"),
		Entry("SoldLastAndStoppedStock - Works", Financial_Trades_SoldLastAndStoppedStock, "31"),
		Entry("SoldOut - Works", Financial_Trades_SoldOut, "32"),
		Entry("SoldOutOfSequence - Works", Financial_Trades_SoldOutOfSequence, "33"),
		Entry("SplitTrade - Works", Financial_Trades_SplitTrade, "34"),
		Entry("StockOption - Works", Financial_Trades_StockOption, "35"),
		Entry("YellowFlagRegularTrade - Works", Financial_Trades_YellowFlagRegularTrade, "36"),
		Entry("OddLotTrade - Works", Financial_Trades_OddLotTrade, "37"),
		Entry("CorrectedConsolidatedClose - Works", Financial_Trades_CorrectedConsolidatedClose, "38"),
		Entry("Unknown - Works", Financial_Trades_Unknown, "39"),
		Entry("Held - Works", Financial_Trades_Held, "40"),
		Entry("TradeThruExempt - Works", Financial_Trades_TradeThruExempt, "41"),
		Entry("NonEligible - Works", Financial_Trades_NonEligible, "42"),
		Entry("NonEligibleExtended - Works", Financial_Trades_NonEligibleExtended, "43"),
		Entry("Cancelled - Works", Financial_Trades_Cancelled, "44"),
		Entry("Recovery - Works", Financial_Trades_Recovery, "45"),
		Entry("Correction - Works", Financial_Trades_Correction, "46"),
		Entry("AsOf - Works", Financial_Trades_AsOf, "47"),
		Entry("AsOfCorrection - Works", Financial_Trades_AsOfCorrection, "48"),
		Entry("AsOfCancel - Works", Financial_Trades_AsOfCancel, "49"),
		Entry("OOB - Works", Financial_Trades_OOB, "50"),
		Entry("Summary - Works", Financial_Trades_Summary, "51"),
		Entry("ContingentTrade - Works", Financial_Trades_ContingentTrade, "52"),
		Entry("QualifiedContingentTrade - Works", Financial_Trades_QualifiedContingentTrade, "53"),
		Entry("Errored - Works", Financial_Trades_Errored, "54"),
		Entry("OpeningReopeningTradeDetail - Works",
			Financial_Trades_OpeningReopeningTradeDetail, "55"),
		Entry("Placeholder - Works", Financial_Trades_Placeholder, "56"),
		Entry("ShortSaleRestrictionActivated - Works",
			Financial_Trades_ShortSaleRestrictionActivated, "57"),
		Entry("ShortSaleRestrictionContinued - Works",
			Financial_Trades_ShortSaleRestrictionContinued, "58"),
		Entry("ShortSaleRestrictionDeactivated - Works",
			Financial_Trades_ShortSaleRestrictionDeactivated, "59"),
		Entry("ShortSaleRestrictionInEffect - Works",
			Financial_Trades_ShortSaleRestrictionInEffect, "60"),
		Entry("FinancialStatusBankrupt - Works", Financial_Trades_FinancialStatusBankrupt, "62"),
		Entry("FinancialStatusDeficient - Works", Financial_Trades_FinancialStatusDeficient, "63"),
		Entry("FinancialStatusDelinquent - Works", Financial_Trades_FinancialStatusDelinquent, "64"),
		Entry("FinancialStatusBankruptAndDeficient - Works",
			Financial_Trades_FinancialStatusBankruptAndDeficient, "65"),
		Entry("FinancialStatusBankruptAndDelinquent - Works",
			Financial_Trades_FinancialStatusBankruptAndDelinquent, "66"),
		Entry("FinancialStatusDeficientAndDelinquent - Works",
			Financial_Trades_FinancialStatusDeficientAndDelinquent, "67"),
		Entry("FinancialStatusDeficientDelinquentBankrupt - Works",
			Financial_Trades_FinancialStatusDeficientDelinquentBankrupt, "68"),
		Entry("FinancialStatusLiquidation - Works", Financial_Trades_FinancialStatusLiquidation, "69"),
		Entry("FinancialStatusCreationsSuspended - Works",
			Financial_Trades_FinancialStatusCreationsSuspended, "70"),
		Entry("FinancialStatusRedemptionsSuspended - Works",
			Financial_Trades_FinancialStatusRedemptionsSuspended, "71"),
		Entry("Canceled - Works", Financial_Trades_Canceled, "201"),
		Entry("LateAndOutOfSequence - Works", Financial_Trades_LateAndOutOfSequence, "202"),
		Entry("LastAndCanceled - Works", Financial_Trades_LastAndCanceled, "203"),
		Entry("Late - Works", Financial_Trades_Late, "204"),
		Entry("OpeningTradeAndCanceled - Works", Financial_Trades_OpeningTradeAndCanceled, "205"),
		Entry("OpeningTradeLateAndOutOfSequence - Works",
			Financial_Trades_OpeningTradeLateAndOutOfSequence, "206"),
		Entry("OnlyTradeAndCanceled - Works", Financial_Trades_OnlyTradeAndCanceled, "207"),
		Entry("OpeningTradeAndLate - Works", Financial_Trades_OpeningTradeAndLate, "208"),
		Entry("AutomaticExecutionOption - Works", Financial_Trades_AutomaticExecutionOption, "209"),
		Entry("ReopeningTrade - Works", Financial_Trades_ReopeningTrade, "210"),
		Entry("IntermarketSweepOrder - Works", Financial_Trades_IntermarketSweepOrder, "219"),
		Entry("SingleLegAuctionNonISO - Works", Financial_Trades_SingleLegAuctionNonISO, "227"),
		Entry("SingleLegAuctionISO - Works", Financial_Trades_SingleLegAuctionISO, "228"),
		Entry("SingleLegCrossNonISO - Works", Financial_Trades_SingleLegCrossNonISO, "229"),
		Entry("SingleLegCrossISO - Works", Financial_Trades_SingleLegCrossISO, "230"),
		Entry("SingleLegFloorTrade - Works", Financial_Trades_SingleLegFloorTrade, "231"),
		Entry("MultiLegAutoElectronicTrade - Works",
			Financial_Trades_MultiLegAutoElectronicTrade, "232"),
		Entry("MultiLegAuction - Works", Financial_Trades_MultiLegAuction, "233"),
		Entry("MultiLegCross - Works", Financial_Trades_MultiLegCross, "234"),
		Entry("MultiLegFloorTrade - Works", Financial_Trades_MultiLegFloorTrade, "235"),
		Entry("MultiLegAutoElectronicTradeAgainstSingleLeg - Works",
			Financial_Trades_MultiLegAutoElectronicTradeAgainstSingleLeg, "236"),
		Entry("StockOptionsAuction - Works", Financial_Trades_StockOptionsAuction, "237"),
		Entry("MultiLegAuctionAgainstSingleLeg - Works",
			Financial_Trades_MultiLegAuctionAgainstSingleLeg, "238"),
		Entry("MultiLegFloorTradeAgainstSingleLeg - Works",
			Financial_Trades_MultiLegFloorTradeAgainstSingleLeg, "239"),
		Entry("StockOptionsAutoElectronicTrade - Works",
			Financial_Trades_StockOptionsAutoElectronicTrade, "240"),
		Entry("StockOptionsCross - Works", Financial_Trades_StockOptionsCross, "241"),
		Entry("StockOptionsFloorTrade - Works", Financial_Trades_StockOptionsFloorTrade, "242"),
		Entry("StockOptionsAutoElectronicTradeAgainstSingleLeg - Works",
			Financial_Trades_StockOptionsAutoElectronicTradeAgainstSingleLeg, "243"),
		Entry("StockOptionsAuctionAgainstSingleLeg - Works",
			Financial_Trades_StockOptionsAuctionAgainstSingleLeg, "244"),
		Entry("StockOptionsFloorTradeAgainstSingleLeg - Works",
			Financial_Trades_StockOptionsFloorTradeAgainstSingleLeg, "245"),
		Entry("MultiLegFloorTradeOfProprietaryProducts - Works",
			Financial_Trades_MultiLegFloorTradeOfProprietaryProducts, "246"),
		Entry("MultilateralCompressionTradeOfProprietaryProducts - Works",
			Financial_Trades_MultilateralCompressionTradeOfProprietaryProducts, "247"),
		Entry("ExtendedHoursTrade - Works", Financial_Trades_ExtendedHoursTrade, "248"))

	// Test that converting the Financial.Trades.Condition enum to a AttributeValue works for all values
	DescribeTable("MarshalDynamoDBAttributeValue Tests",
		func(enum Financial_Trades_Condition, value string) {
			data, err := enum.MarshalDynamoDBAttributeValue()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data.(*types.AttributeValueMemberS).Value).Should(Equal(value))
		},
		Entry("RegularSale - Works", Financial_Trades_RegularSale, "RegularSale"),
		Entry("Acquisition - Works", Financial_Trades_Acquisition, "Acquisition"),
		Entry("AveragePriceTrade - Works", Financial_Trades_AveragePriceTrade, "AveragePriceTrade"),
		Entry("AutomaticExecution - Works", Financial_Trades_AutomaticExecution, "AutomaticExecution"),
		Entry("BunchedTrade - Works", Financial_Trades_BunchedTrade, "BunchedTrade"),
		Entry("BunchedSoldTrade - Works", Financial_Trades_BunchedSoldTrade, "BunchedSoldTrade"),
		Entry("CAPElection - Works", Financial_Trades_CAPElection, "CAPElection"),
		Entry("CashSale - Works", Financial_Trades_CashSale, "CashSale"),
		Entry("ClosingPrints - Works", Financial_Trades_ClosingPrints, "ClosingPrints"),
		Entry("CrossTrade - Works", Financial_Trades_CrossTrade, "CrossTrade"),
		Entry("DerivativelyPriced - Works", Financial_Trades_DerivativelyPriced, "DerivativelyPriced"),
		Entry("Distribution - Works", Financial_Trades_Distribution, "Distribution"),
		Entry("FormT - Works", Financial_Trades_FormT, "FormT"),
		Entry("ExtendedTradingHours - Works", Financial_Trades_ExtendedTradingHours, "ExtendedTradingHours"),
		Entry("IntermarketSweep - Works", Financial_Trades_IntermarketSweep, "IntermarketSweep"),
		Entry("MarketCenterOfficialClose - Works", Financial_Trades_MarketCenterOfficialClose, "MarketCenterOfficialClose"),
		Entry("MarketCenterOfficialOpen - Works", Financial_Trades_MarketCenterOfficialOpen, "MarketCenterOfficialOpen"),
		Entry("MarketCenterOpeningTrade - Works", Financial_Trades_MarketCenterOpeningTrade, "MarketCenterOpeningTrade"),
		Entry("MarketCenterReopeningTrade - Works",
			Financial_Trades_MarketCenterReopeningTrade, "MarketCenterReopeningTrade"),
		Entry("MarketCenterClosingTrade - Works", Financial_Trades_MarketCenterClosingTrade, "MarketCenterClosingTrade"),
		Entry("NextDay - Works", Financial_Trades_NextDay, "NextDay"),
		Entry("PriceVariationTrade - Works", Financial_Trades_PriceVariationTrade, "PriceVariationTrade"),
		Entry("PriorReferencePrice - Works", Financial_Trades_PriorReferencePrice, "PriorReferencePrice"),
		Entry("Rule155Trade - Works", Financial_Trades_Rule155Trade, "Rule155Trade"),
		Entry("Rule127NYSE - Works", Financial_Trades_Rule127NYSE, "Rule127NYSE"),
		Entry("OpeningPrints - Works", Financial_Trades_OpeningPrints, "OpeningPrints"),
		Entry("Opened - Works", Financial_Trades_Opened, "Opened"),
		Entry("StoppedStock - Works", Financial_Trades_StoppedStock, "StoppedStock"),
		Entry("ReOpeningPrints - Works", Financial_Trades_ReOpeningPrints, "ReOpeningPrints"),
		Entry("Seller - Works", Financial_Trades_Seller, "Seller"),
		Entry("SoldLast - Works", Financial_Trades_SoldLast, "SoldLast"),
		Entry("SoldLastAndStoppedStock - Works", Financial_Trades_SoldLastAndStoppedStock, "SoldLastAndStoppedStock"),
		Entry("SoldOut - Works", Financial_Trades_SoldOut, "SoldOut"),
		Entry("SoldOutOfSequence - Works", Financial_Trades_SoldOutOfSequence, "SoldOutOfSequence"),
		Entry("SplitTrade - Works", Financial_Trades_SplitTrade, "SplitTrade"),
		Entry("StockOption - Works", Financial_Trades_StockOption, "StockOption"),
		Entry("YellowFlagRegularTrade - Works", Financial_Trades_YellowFlagRegularTrade, "YellowFlagRegularTrade"),
		Entry("OddLotTrade - Works", Financial_Trades_OddLotTrade, "OddLotTrade"),
		Entry("CorrectedConsolidatedClose - Works",
			Financial_Trades_CorrectedConsolidatedClose, "CorrectedConsolidatedClose"),
		Entry("Unknown - Works", Financial_Trades_Unknown, "Unknown"),
		Entry("Held - Works", Financial_Trades_Held, "Held"),
		Entry("TradeThruExempt - Works", Financial_Trades_TradeThruExempt, "TradeThruExempt"),
		Entry("NonEligible - Works", Financial_Trades_NonEligible, "NonEligible"),
		Entry("NonEligibleExtended - Works", Financial_Trades_NonEligibleExtended, "NonEligibleExtended"),
		Entry("Cancelled - Works", Financial_Trades_Cancelled, "Cancelled"),
		Entry("Recovery - Works", Financial_Trades_Recovery, "Recovery"),
		Entry("Correction - Works", Financial_Trades_Correction, "Correction"),
		Entry("AsOf - Works", Financial_Trades_AsOf, "AsOf"),
		Entry("AsOfCorrection - Works", Financial_Trades_AsOfCorrection, "AsOfCorrection"),
		Entry("AsOfCancel - Works", Financial_Trades_AsOfCancel, "AsOfCancel"),
		Entry("OOB - Works", Financial_Trades_OOB, "OOB"),
		Entry("Summary - Works", Financial_Trades_Summary, "Summary"),
		Entry("ContingentTrade - Works", Financial_Trades_ContingentTrade, "ContingentTrade"),
		Entry("QualifiedContingentTrade - Works", Financial_Trades_QualifiedContingentTrade, "QualifiedContingentTrade"),
		Entry("Errored - Works", Financial_Trades_Errored, "Errored"),
		Entry("OpeningReopeningTradeDetail - Works",
			Financial_Trades_OpeningReopeningTradeDetail, "OpeningReopeningTradeDetail"),
		Entry("Placeholder - Works", Financial_Trades_Placeholder, "Placeholder"),
		Entry("ShortSaleRestrictionActivated - Works",
			Financial_Trades_ShortSaleRestrictionActivated, "ShortSaleRestrictionActivated"),
		Entry("ShortSaleRestrictionContinued - Works",
			Financial_Trades_ShortSaleRestrictionContinued, "ShortSaleRestrictionContinued"),
		Entry("ShortSaleRestrictionDeactivated - Works",
			Financial_Trades_ShortSaleRestrictionDeactivated, "ShortSaleRestrictionDeactivated"),
		Entry("ShortSaleRestrictionInEffect - Works",
			Financial_Trades_ShortSaleRestrictionInEffect, "ShortSaleRestrictionInEffect"),
		Entry("FinancialStatusBankrupt - Works", Financial_Trades_FinancialStatusBankrupt, "FinancialStatusBankrupt"),
		Entry("FinancialStatusDeficient - Works", Financial_Trades_FinancialStatusDeficient, "FinancialStatusDeficient"),
		Entry("FinancialStatusDelinquent - Works", Financial_Trades_FinancialStatusDelinquent, "FinancialStatusDelinquent"),
		Entry("FinancialStatusBankruptAndDeficient - Works",
			Financial_Trades_FinancialStatusBankruptAndDeficient, "FinancialStatusBankruptAndDeficient"),
		Entry("FinancialStatusBankruptAndDelinquent - Works",
			Financial_Trades_FinancialStatusBankruptAndDelinquent, "FinancialStatusBankruptAndDelinquent"),
		Entry("FinancialStatusDeficientAndDelinquent - Works",
			Financial_Trades_FinancialStatusDeficientAndDelinquent, "FinancialStatusDeficientAndDelinquent"),
		Entry("FinancialStatusDeficientDelinquentBankrupt - Works",
			Financial_Trades_FinancialStatusDeficientDelinquentBankrupt, "FinancialStatusDeficientDelinquentBankrupt"),
		Entry("FinancialStatusLiquidation - Works",
			Financial_Trades_FinancialStatusLiquidation, "FinancialStatusLiquidation"),
		Entry("FinancialStatusCreationsSuspended - Works",
			Financial_Trades_FinancialStatusCreationsSuspended, "FinancialStatusCreationsSuspended"),
		Entry("FinancialStatusRedemptionsSuspended - Works",
			Financial_Trades_FinancialStatusRedemptionsSuspended, "FinancialStatusRedemptionsSuspended"),
		Entry("Canceled - Works", Financial_Trades_Canceled, "Canceled"),
		Entry("LateAndOutOfSequence - Works", Financial_Trades_LateAndOutOfSequence, "LateAndOutOfSequence"),
		Entry("LastAndCanceled - Works", Financial_Trades_LastAndCanceled, "LastAndCanceled"),
		Entry("Late - Works", Financial_Trades_Late, "Late"),
		Entry("OpeningTradeAndCanceled - Works", Financial_Trades_OpeningTradeAndCanceled, "OpeningTradeAndCanceled"),
		Entry("OpeningTradeLateAndOutOfSequence - Works",
			Financial_Trades_OpeningTradeLateAndOutOfSequence, "OpeningTradeLateAndOutOfSequence"),
		Entry("OnlyTradeAndCanceled - Works", Financial_Trades_OnlyTradeAndCanceled, "OnlyTradeAndCanceled"),
		Entry("OpeningTradeAndLate - Works", Financial_Trades_OpeningTradeAndLate, "OpeningTradeAndLate"),
		Entry("AutomaticExecutionOption - Works", Financial_Trades_AutomaticExecutionOption, "AutomaticExecutionOption"),
		Entry("ReopeningTrade - Works", Financial_Trades_ReopeningTrade, "ReopeningTrade"),
		Entry("IntermarketSweepOrder - Works", Financial_Trades_IntermarketSweepOrder, "IntermarketSweepOrder"),
		Entry("SingleLegAuctionNonISO - Works", Financial_Trades_SingleLegAuctionNonISO, "SingleLegAuctionNonISO"),
		Entry("SingleLegAuctionISO - Works", Financial_Trades_SingleLegAuctionISO, "SingleLegAuctionISO"),
		Entry("SingleLegCrossNonISO - Works", Financial_Trades_SingleLegCrossNonISO, "SingleLegCrossNonISO"),
		Entry("SingleLegCrossISO - Works", Financial_Trades_SingleLegCrossISO, "SingleLegCrossISO"),
		Entry("SingleLegFloorTrade - Works", Financial_Trades_SingleLegFloorTrade, "SingleLegFloorTrade"),
		Entry("MultiLegAutoElectronicTrade - Works",
			Financial_Trades_MultiLegAutoElectronicTrade, "MultiLegAutoElectronicTrade"),
		Entry("MultiLegAuction - Works", Financial_Trades_MultiLegAuction, "MultiLegAuction"),
		Entry("MultiLegCross - Works", Financial_Trades_MultiLegCross, "MultiLegCross"),
		Entry("MultiLegFloorTrade - Works", Financial_Trades_MultiLegFloorTrade, "MultiLegFloorTrade"),
		Entry("MultiLegAutoElectronicTradeAgainstSingleLeg - Works",
			Financial_Trades_MultiLegAutoElectronicTradeAgainstSingleLeg, "MultiLegAutoElectronicTradeAgainstSingleLeg"),
		Entry("StockOptionsAuction - Works", Financial_Trades_StockOptionsAuction, "StockOptionsAuction"),
		Entry("MultiLegAuctionAgainstSingleLeg - Works",
			Financial_Trades_MultiLegAuctionAgainstSingleLeg, "MultiLegAuctionAgainstSingleLeg"),
		Entry("MultiLegFloorTradeAgainstSingleLeg - Works",
			Financial_Trades_MultiLegFloorTradeAgainstSingleLeg, "MultiLegFloorTradeAgainstSingleLeg"),
		Entry("StockOptionsAutoElectronicTrade - Works",
			Financial_Trades_StockOptionsAutoElectronicTrade, "StockOptionsAutoElectronicTrade"),
		Entry("StockOptionsCross - Works", Financial_Trades_StockOptionsCross, "StockOptionsCross"),
		Entry("StockOptionsFloorTrade - Works", Financial_Trades_StockOptionsFloorTrade, "StockOptionsFloorTrade"),
		Entry("StockOptionsAutoElectronicTradeAgainstSingleLeg - Works",
			Financial_Trades_StockOptionsAutoElectronicTradeAgainstSingleLeg, "StockOptionsAutoElectronicTradeAgainstSingleLeg"),
		Entry("StockOptionsAuctionAgainstSingleLeg - Works",
			Financial_Trades_StockOptionsAuctionAgainstSingleLeg, "StockOptionsAuctionAgainstSingleLeg"),
		Entry("StockOptionsFloorTradeAgainstSingleLeg - Works",
			Financial_Trades_StockOptionsFloorTradeAgainstSingleLeg, "StockOptionsFloorTradeAgainstSingleLeg"),
		Entry("MultiLegFloorTradeOfProprietaryProducts - Works",
			Financial_Trades_MultiLegFloorTradeOfProprietaryProducts, "MultiLegFloorTradeOfProprietaryProducts"),
		Entry("MultilateralCompressionTradeOfProprietaryProducts - Works",
			Financial_Trades_MultilateralCompressionTradeOfProprietaryProducts, "MultilateralCompressionTradeOfProprietaryProducts"),
		Entry("ExtendedHoursTrade - Works", Financial_Trades_ExtendedHoursTrade, "ExtendedHoursTrade"))

	// Test that attempting to deserialize a Financial.Trades.Condition will fail and
	// return an error if the value canno be deserialized from a JSON value to a string
	It("UnmarshalJSON fails - Error", func() {

		// Attempt to convert a non-parseable string value into a Financial.Trades.Condition
		// This should return an error
		enum := new(Financial_Trades_Condition)
		err := enum.UnmarshalJSON([]byte("derp"))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Trades_Condition"))
	})

	// Test that attempting to deserialize a Financial.Trades.Condition will fail and
	// return an error if the value cannot be converted to either the name value or integer
	// value of the enum option
	It("UnmarshalJSON - Value is invalid - Error", func() {

		// Attempt to convert a fake string value into a Financial.Trades.Condition
		// This should return an error
		enum := new(Financial_Trades_Condition)
		err := enum.UnmarshalJSON([]byte("\"derp\""))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Trades_Condition"))
	})

	// Test the conditions under which values should be convertible to a Financial.Trades.Condition
	DescribeTable("UnmarshalJSON Tests",
		func(value interface{}, shouldBe Financial_Trades_Condition) {

			// Attempt to convert the string value into a Financial.Trades.Condition
			// This should not fail
			var enum Financial_Trades_Condition
			err := enum.UnmarshalJSON([]byte(fmt.Sprintf("%v", value)))

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("RegularSale - Works", "\"RegularSale\"", Financial_Trades_RegularSale),
		Entry("Acquisition - Works", "\"Acquisition\"", Financial_Trades_Acquisition),
		Entry("AveragePriceTrade - Works", "\"AveragePriceTrade\"", Financial_Trades_AveragePriceTrade),
		Entry("AutomaticExecution - Works", "\"AutomaticExecution\"", Financial_Trades_AutomaticExecution),
		Entry("BunchedTrade - Works", "\"BunchedTrade\"", Financial_Trades_BunchedTrade),
		Entry("BunchedSoldTrade - Works", "\"BunchedSoldTrade\"", Financial_Trades_BunchedSoldTrade),
		Entry("CAPElection - Works", "\"CAPElection\"", Financial_Trades_CAPElection),
		Entry("CashSale - Works", "\"CashSale\"", Financial_Trades_CashSale),
		Entry("ClosingPrints - Works", "\"ClosingPrints\"", Financial_Trades_ClosingPrints),
		Entry("CrossTrade - Works", "\"CrossTrade\"", Financial_Trades_CrossTrade),
		Entry("DerivativelyPriced - Works", "\"DerivativelyPriced\"", Financial_Trades_DerivativelyPriced),
		Entry("Distribution - Works", "\"Distribution\"", Financial_Trades_Distribution),
		Entry("FormT - Works", "\"FormT\"", Financial_Trades_FormT),
		Entry("ExtendedTradingHours - Works", "\"ExtendedTradingHours\"", Financial_Trades_ExtendedTradingHours),
		Entry("IntermarketSweep - Works", "\"IntermarketSweep\"", Financial_Trades_IntermarketSweep),
		Entry("MarketCenterOfficialClose - Works", "\"MarketCenterOfficialClose\"", Financial_Trades_MarketCenterOfficialClose),
		Entry("MarketCenterOfficialOpen - Works", "\"MarketCenterOfficialOpen\"", Financial_Trades_MarketCenterOfficialOpen),
		Entry("MarketCenterOpeningTrade - Works", "\"MarketCenterOpeningTrade\"", Financial_Trades_MarketCenterOpeningTrade),
		Entry("MarketCenterReopeningTrade - Works", "\"MarketCenterReopeningTrade\"", Financial_Trades_MarketCenterReopeningTrade),
		Entry("MarketCenterClosingTrade - Works", "\"MarketCenterClosingTrade\"", Financial_Trades_MarketCenterClosingTrade),
		Entry("NextDay - Works", "\"NextDay\"", Financial_Trades_NextDay),
		Entry("PriceVariationTrade - Works", "\"PriceVariationTrade\"", Financial_Trades_PriceVariationTrade),
		Entry("PriorReferencePrice - Works", "\"PriorReferencePrice\"", Financial_Trades_PriorReferencePrice),
		Entry("Rule155Trade - Works", "\"Rule155Trade\"", Financial_Trades_Rule155Trade),
		Entry("Rule127NYSE - Works", "\"Rule127NYSE\"", Financial_Trades_Rule127NYSE),
		Entry("OpeningPrints - Works", "\"OpeningPrints\"", Financial_Trades_OpeningPrints),
		Entry("Opened - Works", "\"Opened\"", Financial_Trades_Opened),
		Entry("StoppedStock - Works", "\"StoppedStock\"", Financial_Trades_StoppedStock),
		Entry("ReOpeningPrints - Works", "\"ReOpeningPrints\"", Financial_Trades_ReOpeningPrints),
		Entry("Seller - Works", "\"Seller\"", Financial_Trades_Seller),
		Entry("SoldLast - Works", "\"SoldLast\"", Financial_Trades_SoldLast),
		Entry("SoldLastAndStoppedStock - Works", "\"SoldLastAndStoppedStock\"", Financial_Trades_SoldLastAndStoppedStock),
		Entry("SoldOut - Works", "\"SoldOut\"", Financial_Trades_SoldOut),
		Entry("SoldOutOfSequence - Works", "\"SoldOutOfSequence\"", Financial_Trades_SoldOutOfSequence),
		Entry("SplitTrade - Works", "\"SplitTrade\"", Financial_Trades_SplitTrade),
		Entry("StockOption - Works", "\"StockOption\"", Financial_Trades_StockOption),
		Entry("YellowFlagRegularTrade - Works", "\"YellowFlagRegularTrade\"", Financial_Trades_YellowFlagRegularTrade),
		Entry("OddLotTrade - Works", "\"OddLotTrade\"", Financial_Trades_OddLotTrade),
		Entry("CorrectedConsolidatedClose - Works", "\"CorrectedConsolidatedClose\"", Financial_Trades_CorrectedConsolidatedClose),
		Entry("Unknown - Works", "\"Unknown\"", Financial_Trades_Unknown),
		Entry("Held - Works", "\"Held\"", Financial_Trades_Held),
		Entry("TradeThruExempt - Works", "\"TradeThruExempt\"", Financial_Trades_TradeThruExempt),
		Entry("NonEligible - Works", "\"NonEligible\"", Financial_Trades_NonEligible),
		Entry("NonEligibleExtended - Works", "\"NonEligibleExtended\"", Financial_Trades_NonEligibleExtended),
		Entry("Cancelled - Works", "\"Cancelled\"", Financial_Trades_Cancelled),
		Entry("Recovery - Works", "\"Recovery\"", Financial_Trades_Recovery),
		Entry("Correction - Works", "\"Correction\"", Financial_Trades_Correction),
		Entry("AsOf - Works", "\"AsOf\"", Financial_Trades_AsOf),
		Entry("AsOfCorrection - Works", "\"AsOfCorrection\"", Financial_Trades_AsOfCorrection),
		Entry("AsOfCancel - Works", "\"AsOfCancel\"", Financial_Trades_AsOfCancel),
		Entry("OOB - Works", "\"OOB\"", Financial_Trades_OOB),
		Entry("Summary - Works", "\"Summary\"", Financial_Trades_Summary),
		Entry("ContingentTrade - Works", "\"ContingentTrade\"", Financial_Trades_ContingentTrade),
		Entry("QualifiedContingentTrade - Works", "\"QualifiedContingentTrade\"", Financial_Trades_QualifiedContingentTrade),
		Entry("Errored - Works", "\"Errored\"", Financial_Trades_Errored),
		Entry("OpeningReopeningTradeDetail - Works",
			"\"OpeningReopeningTradeDetail\"", Financial_Trades_OpeningReopeningTradeDetail),
		Entry("Placeholder - Works", "\"Placeholder\"", Financial_Trades_Placeholder),
		Entry("ShortSaleRestrictionActivated - Works",
			"\"ShortSaleRestrictionActivated\"", Financial_Trades_ShortSaleRestrictionActivated),
		Entry("ShortSaleRestrictionContinued - Works",
			"\"ShortSaleRestrictionContinued\"", Financial_Trades_ShortSaleRestrictionContinued),
		Entry("ShortSaleRestrictionDeactivated - Works",
			"\"ShortSaleRestrictionDeactivated\"", Financial_Trades_ShortSaleRestrictionDeactivated),
		Entry("ShortSaleRestrictionInEffect - Works",
			"\"ShortSaleRestrictionInEffect\"", Financial_Trades_ShortSaleRestrictionInEffect),
		Entry("FinancialStatusBankrupt - Works", "\"FinancialStatusBankrupt\"", Financial_Trades_FinancialStatusBankrupt),
		Entry("FinancialStatusDeficient - Works", "\"FinancialStatusDeficient\"", Financial_Trades_FinancialStatusDeficient),
		Entry("FinancialStatusDelinquent - Works", "\"FinancialStatusDelinquent\"", Financial_Trades_FinancialStatusDelinquent),
		Entry("FinancialStatusBankruptAndDeficient - Works",
			"\"FinancialStatusBankruptAndDeficient\"", Financial_Trades_FinancialStatusBankruptAndDeficient),
		Entry("FinancialStatusBankruptAndDelinquent - Works",
			"\"FinancialStatusBankruptAndDelinquent\"", Financial_Trades_FinancialStatusBankruptAndDelinquent),
		Entry("FinancialStatusDeficientAndDelinquent - Works",
			"\"FinancialStatusDeficientAndDelinquent\"", Financial_Trades_FinancialStatusDeficientAndDelinquent),
		Entry("FinancialStatusDeficientDelinquentBankrupt - Works",
			"\"FinancialStatusDeficientDelinquentBankrupt\"", Financial_Trades_FinancialStatusDeficientDelinquentBankrupt),
		Entry("FinancialStatusLiquidation - Works", "\"FinancialStatusLiquidation\"", Financial_Trades_FinancialStatusLiquidation),
		Entry("FinancialStatusCreationsSuspended - Works",
			"\"FinancialStatusCreationsSuspended\"", Financial_Trades_FinancialStatusCreationsSuspended),
		Entry("FinancialStatusRedemptionsSuspended - Works",
			"\"FinancialStatusRedemptionsSuspended\"", Financial_Trades_FinancialStatusRedemptionsSuspended),
		Entry("Canceled - Works", "\"Canceled\"", Financial_Trades_Canceled),
		Entry("LateAndOutOfSequence - Works", "\"LateAndOutOfSequence\"", Financial_Trades_LateAndOutOfSequence),
		Entry("LastAndCanceled - Works", "\"LastAndCanceled\"", Financial_Trades_LastAndCanceled),
		Entry("Late - Works", "\"Late\"", Financial_Trades_Late),
		Entry("OpeningTradeAndCanceled - Works", "\"OpeningTradeAndCanceled\"", Financial_Trades_OpeningTradeAndCanceled),
		Entry("OpeningTradeLateAndOutOfSequence - Works",
			"\"OpeningTradeLateAndOutOfSequence\"", Financial_Trades_OpeningTradeLateAndOutOfSequence),
		Entry("OnlyTradeAndCanceled - Works", "\"OnlyTradeAndCanceled\"", Financial_Trades_OnlyTradeAndCanceled),
		Entry("OpeningTradeAndLate - Works", "\"OpeningTradeAndLate\"", Financial_Trades_OpeningTradeAndLate),
		Entry("AutomaticExecutionOption - Works", "\"AutomaticExecutionOption\"", Financial_Trades_AutomaticExecutionOption),
		Entry("ReopeningTrade - Works", "\"ReopeningTrade\"", Financial_Trades_ReopeningTrade),
		Entry("IntermarketSweepOrder - Works", "\"IntermarketSweepOrder\"", Financial_Trades_IntermarketSweepOrder),
		Entry("SingleLegAuctionNonISO - Works", "\"SingleLegAuctionNonISO\"", Financial_Trades_SingleLegAuctionNonISO),
		Entry("SingleLegAuctionISO - Works", "\"SingleLegAuctionISO\"", Financial_Trades_SingleLegAuctionISO),
		Entry("SingleLegCrossNonISO - Works", "\"SingleLegCrossNonISO\"", Financial_Trades_SingleLegCrossNonISO),
		Entry("SingleLegCrossISO - Works", "\"SingleLegCrossISO\"", Financial_Trades_SingleLegCrossISO),
		Entry("SingleLegFloorTrade - Works", "\"SingleLegFloorTrade\"", Financial_Trades_SingleLegFloorTrade),
		Entry("MultiLegAutoElectronicTrade - Works",
			"\"MultiLegAutoElectronicTrade\"", Financial_Trades_MultiLegAutoElectronicTrade),
		Entry("MultiLegAuction - Works", "\"MultiLegAuction\"", Financial_Trades_MultiLegAuction),
		Entry("MultiLegCross - Works", "\"MultiLegCross\"", Financial_Trades_MultiLegCross),
		Entry("MultiLegFloorTrade - Works", "\"MultiLegFloorTrade\"", Financial_Trades_MultiLegFloorTrade),
		Entry("MultiLegAutoElectronicTradeAgainstSingleLeg - Works",
			"\"MultiLegAutoElectronicTradeAgainstSingleLeg\"", Financial_Trades_MultiLegAutoElectronicTradeAgainstSingleLeg),
		Entry("StockOptionsAuction - Works", "\"StockOptionsAuction\"", Financial_Trades_StockOptionsAuction),
		Entry("MultiLegAuctionAgainstSingleLeg - Works",
			"\"MultiLegAuctionAgainstSingleLeg\"", Financial_Trades_MultiLegAuctionAgainstSingleLeg),
		Entry("MultiLegFloorTradeAgainstSingleLeg - Works",
			"\"MultiLegFloorTradeAgainstSingleLeg\"", Financial_Trades_MultiLegFloorTradeAgainstSingleLeg),
		Entry("StockOptionsAutoElectronicTrade - Works",
			"\"StockOptionsAutoElectronicTrade\"", Financial_Trades_StockOptionsAutoElectronicTrade),
		Entry("StockOptionsCross - Works", "\"StockOptionsCross\"", Financial_Trades_StockOptionsCross),
		Entry("StockOptionsFloorTrade - Works", "\"StockOptionsFloorTrade\"", Financial_Trades_StockOptionsFloorTrade),
		Entry("StockOptionsAutoElectronicTradeAgainstSingleLeg - Works",
			"\"StockOptionsAutoElectronicTradeAgainstSingleLeg\"", Financial_Trades_StockOptionsAutoElectronicTradeAgainstSingleLeg),
		Entry("StockOptionsAuctionAgainstSingleLeg - Works",
			"\"StockOptionsAuctionAgainstSingleLeg\"", Financial_Trades_StockOptionsAuctionAgainstSingleLeg),
		Entry("StockOptionsFloorTradeAgainstSingleLeg - Works",
			"\"StockOptionsFloorTradeAgainstSingleLeg\"", Financial_Trades_StockOptionsFloorTradeAgainstSingleLeg),
		Entry("MultiLegFloorTradeOfProprietaryProducts - Works",
			"\"MultiLegFloorTradeOfProprietaryProducts\"", Financial_Trades_MultiLegFloorTradeOfProprietaryProducts),
		Entry("MultilateralCompressionTradeOfProprietaryProducts - Works",
			"\"MultilateralCompressionTradeOfProprietaryProducts\"", Financial_Trades_MultilateralCompressionTradeOfProprietaryProducts),
		Entry("ExtendedHoursTrade - Works", "\"ExtendedHoursTrade\"", Financial_Trades_ExtendedHoursTrade),
		Entry("'0' - Works", "\"0\"", Financial_Trades_RegularSale),
		Entry("'1' - Works", "\"1\"", Financial_Trades_Acquisition),
		Entry("'2' - Works", "\"2\"", Financial_Trades_AveragePriceTrade),
		Entry("'3' - Works", "\"3\"", Financial_Trades_AutomaticExecution),
		Entry("'4' - Works", "\"4\"", Financial_Trades_BunchedTrade),
		Entry("'5' - Works", "\"5\"", Financial_Trades_BunchedSoldTrade),
		Entry("'6' - Works", "\"6\"", Financial_Trades_CAPElection),
		Entry("'7' - Works", "\"7\"", Financial_Trades_CashSale),
		Entry("'8' - Works", "\"8\"", Financial_Trades_ClosingPrints),
		Entry("'9' - Works", "\"9\"", Financial_Trades_CrossTrade),
		Entry("'10' - Works", "\"10\"", Financial_Trades_DerivativelyPriced),
		Entry("'11' - Works", "\"11\"", Financial_Trades_Distribution),
		Entry("'12' - Works", "\"12\"", Financial_Trades_FormT),
		Entry("'13' - Works", "\"13\"", Financial_Trades_ExtendedTradingHours),
		Entry("'14' - Works", "\"14\"", Financial_Trades_IntermarketSweep),
		Entry("'15' - Works", "\"15\"", Financial_Trades_MarketCenterOfficialClose),
		Entry("'16' - Works", "\"16\"", Financial_Trades_MarketCenterOfficialOpen),
		Entry("'17' - Works", "\"17\"", Financial_Trades_MarketCenterOpeningTrade),
		Entry("'18' - Works", "\"18\"", Financial_Trades_MarketCenterReopeningTrade),
		Entry("'19' - Works", "\"19\"", Financial_Trades_MarketCenterClosingTrade),
		Entry("'20' - Works", "\"20\"", Financial_Trades_NextDay),
		Entry("'21' - Works", "\"21\"", Financial_Trades_PriceVariationTrade),
		Entry("'22' - Works", "\"22\"", Financial_Trades_PriorReferencePrice),
		Entry("'23' - Works", "\"23\"", Financial_Trades_Rule155Trade),
		Entry("'24' - Works", "\"24\"", Financial_Trades_Rule127NYSE),
		Entry("'25' - Works", "\"25\"", Financial_Trades_OpeningPrints),
		Entry("'26' - Works", "\"26\"", Financial_Trades_Opened),
		Entry("'27' - Works", "\"27\"", Financial_Trades_StoppedStock),
		Entry("'28' - Works", "\"28\"", Financial_Trades_ReOpeningPrints),
		Entry("'29' - Works", "\"29\"", Financial_Trades_Seller),
		Entry("'30' - Works", "\"30\"", Financial_Trades_SoldLast),
		Entry("'31' - Works", "\"31\"", Financial_Trades_SoldLastAndStoppedStock),
		Entry("'32' - Works", "\"32\"", Financial_Trades_SoldOut),
		Entry("'33' - Works", "\"33\"", Financial_Trades_SoldOutOfSequence),
		Entry("'34' - Works", "\"34\"", Financial_Trades_SplitTrade),
		Entry("'35' - Works", "\"35\"", Financial_Trades_StockOption),
		Entry("'36' - Works", "\"36\"", Financial_Trades_YellowFlagRegularTrade),
		Entry("'37' - Works", "\"37\"", Financial_Trades_OddLotTrade),
		Entry("'38' - Works", "\"38\"", Financial_Trades_CorrectedConsolidatedClose),
		Entry("'39' - Works", "\"39\"", Financial_Trades_Unknown),
		Entry("'40' - Works", "\"40\"", Financial_Trades_Held),
		Entry("'41' - Works", "\"41\"", Financial_Trades_TradeThruExempt),
		Entry("'42' - Works", "\"42\"", Financial_Trades_NonEligible),
		Entry("'43' - Works", "\"43\"", Financial_Trades_NonEligibleExtended),
		Entry("'44' - Works", "\"44\"", Financial_Trades_Cancelled),
		Entry("'45' - Works", "\"45\"", Financial_Trades_Recovery),
		Entry("'46' - Works", "\"46\"", Financial_Trades_Correction),
		Entry("'47' - Works", "\"47\"", Financial_Trades_AsOf),
		Entry("'48' - Works", "\"48\"", Financial_Trades_AsOfCorrection),
		Entry("'49' - Works", "\"49\"", Financial_Trades_AsOfCancel),
		Entry("'50' - Works", "\"50\"", Financial_Trades_OOB),
		Entry("'51' - Works", "\"51\"", Financial_Trades_Summary),
		Entry("'52' - Works", "\"52\"", Financial_Trades_ContingentTrade),
		Entry("'53' - Works", "\"53\"", Financial_Trades_QualifiedContingentTrade),
		Entry("'54' - Works", "\"54\"", Financial_Trades_Errored),
		Entry("'55' - Works", "\"55\"", Financial_Trades_OpeningReopeningTradeDetail),
		Entry("'56' - Works", "\"56\"", Financial_Trades_Placeholder),
		Entry("'57' - Works", "\"57\"", Financial_Trades_ShortSaleRestrictionActivated),
		Entry("'58' - Works", "\"58\"", Financial_Trades_ShortSaleRestrictionContinued),
		Entry("'59' - Works", "\"59\"", Financial_Trades_ShortSaleRestrictionDeactivated),
		Entry("'60' - Works", "\"60\"", Financial_Trades_ShortSaleRestrictionInEffect),
		Entry("'62' - Works", "\"62\"", Financial_Trades_FinancialStatusBankrupt),
		Entry("'63' - Works", "\"63\"", Financial_Trades_FinancialStatusDeficient),
		Entry("'64' - Works", "\"64\"", Financial_Trades_FinancialStatusDelinquent),
		Entry("'65' - Works", "\"65\"", Financial_Trades_FinancialStatusBankruptAndDeficient),
		Entry("'66' - Works", "\"66\"", Financial_Trades_FinancialStatusBankruptAndDelinquent),
		Entry("'67' - Works", "\"67\"", Financial_Trades_FinancialStatusDeficientAndDelinquent),
		Entry("'68' - Works", "\"68\"", Financial_Trades_FinancialStatusDeficientDelinquentBankrupt),
		Entry("'69' - Works", "\"69\"", Financial_Trades_FinancialStatusLiquidation),
		Entry("'70' - Works", "\"70\"", Financial_Trades_FinancialStatusCreationsSuspended),
		Entry("'71' - Works", "\"71\"", Financial_Trades_FinancialStatusRedemptionsSuspended),
		Entry("'201' - Works", "\"201\"", Financial_Trades_Canceled),
		Entry("'202' - Works", "\"202\"", Financial_Trades_LateAndOutOfSequence),
		Entry("'203' - Works", "\"203\"", Financial_Trades_LastAndCanceled),
		Entry("'204' - Works", "\"204\"", Financial_Trades_Late),
		Entry("'205' - Works", "\"205\"", Financial_Trades_OpeningTradeAndCanceled),
		Entry("'206' - Works", "\"206\"", Financial_Trades_OpeningTradeLateAndOutOfSequence),
		Entry("'207' - Works", "\"207\"", Financial_Trades_OnlyTradeAndCanceled),
		Entry("'208' - Works", "\"208\"", Financial_Trades_OpeningTradeAndLate),
		Entry("'209' - Works", "\"209\"", Financial_Trades_AutomaticExecutionOption),
		Entry("'210' - Works", "\"210\"", Financial_Trades_ReopeningTrade),
		Entry("'219' - Works", "\"219\"", Financial_Trades_IntermarketSweepOrder),
		Entry("'227' - Works", "\"227\"", Financial_Trades_SingleLegAuctionNonISO),
		Entry("'228' - Works", "\"228\"", Financial_Trades_SingleLegAuctionISO),
		Entry("'229' - Works", "\"229\"", Financial_Trades_SingleLegCrossNonISO),
		Entry("'230' - Works", "\"230\"", Financial_Trades_SingleLegCrossISO),
		Entry("'231' - Works", "\"231\"", Financial_Trades_SingleLegFloorTrade),
		Entry("'232' - Works", "\"232\"", Financial_Trades_MultiLegAutoElectronicTrade),
		Entry("'233' - Works", "\"233\"", Financial_Trades_MultiLegAuction),
		Entry("'234' - Works", "\"234\"", Financial_Trades_MultiLegCross),
		Entry("'235' - Works", "\"235\"", Financial_Trades_MultiLegFloorTrade),
		Entry("'236' - Works", "\"236\"", Financial_Trades_MultiLegAutoElectronicTradeAgainstSingleLeg),
		Entry("'237' - Works", "\"237\"", Financial_Trades_StockOptionsAuction),
		Entry("'238' - Works", "\"238\"", Financial_Trades_MultiLegAuctionAgainstSingleLeg),
		Entry("'239' - Works", "\"239\"", Financial_Trades_MultiLegFloorTradeAgainstSingleLeg),
		Entry("'240' - Works", "\"240\"", Financial_Trades_StockOptionsAutoElectronicTrade),
		Entry("'241' - Works", "\"241\"", Financial_Trades_StockOptionsCross),
		Entry("'242' - Works", "\"242\"", Financial_Trades_StockOptionsFloorTrade),
		Entry("'243' - Works", "\"243\"", Financial_Trades_StockOptionsAutoElectronicTradeAgainstSingleLeg),
		Entry("'244' - Works", "\"244\"", Financial_Trades_StockOptionsAuctionAgainstSingleLeg),
		Entry("'245' - Works", "\"245\"", Financial_Trades_StockOptionsFloorTradeAgainstSingleLeg),
		Entry("'246' - Works", "\"246\"", Financial_Trades_MultiLegFloorTradeOfProprietaryProducts),
		Entry("'247' - Works", "\"247\"", Financial_Trades_MultilateralCompressionTradeOfProprietaryProducts),
		Entry("'248' - Works", "\"248\"", Financial_Trades_ExtendedHoursTrade),
		Entry("0 - Works", 0, Financial_Trades_RegularSale),
		Entry("1 - Works", 1, Financial_Trades_Acquisition),
		Entry("2 - Works", 2, Financial_Trades_AveragePriceTrade),
		Entry("3 - Works", 3, Financial_Trades_AutomaticExecution),
		Entry("4 - Works", 4, Financial_Trades_BunchedTrade),
		Entry("5 - Works", 5, Financial_Trades_BunchedSoldTrade),
		Entry("6 - Works", 6, Financial_Trades_CAPElection),
		Entry("7 - Works", 7, Financial_Trades_CashSale),
		Entry("8 - Works", 8, Financial_Trades_ClosingPrints),
		Entry("9 - Works", 9, Financial_Trades_CrossTrade),
		Entry("10 - Works", 10, Financial_Trades_DerivativelyPriced),
		Entry("11 - Works", 11, Financial_Trades_Distribution),
		Entry("12 - Works", 12, Financial_Trades_FormT),
		Entry("13 - Works", 13, Financial_Trades_ExtendedTradingHours),
		Entry("14 - Works", 14, Financial_Trades_IntermarketSweep),
		Entry("15 - Works", 15, Financial_Trades_MarketCenterOfficialClose),
		Entry("16 - Works", 16, Financial_Trades_MarketCenterOfficialOpen),
		Entry("17 - Works", 17, Financial_Trades_MarketCenterOpeningTrade),
		Entry("18 - Works", 18, Financial_Trades_MarketCenterReopeningTrade),
		Entry("19 - Works", 19, Financial_Trades_MarketCenterClosingTrade),
		Entry("20 - Works", 20, Financial_Trades_NextDay),
		Entry("21 - Works", 21, Financial_Trades_PriceVariationTrade),
		Entry("22 - Works", 22, Financial_Trades_PriorReferencePrice),
		Entry("23 - Works", 23, Financial_Trades_Rule155Trade),
		Entry("24 - Works", 24, Financial_Trades_Rule127NYSE),
		Entry("25 - Works", 25, Financial_Trades_OpeningPrints),
		Entry("26 - Works", 26, Financial_Trades_Opened),
		Entry("27 - Works", 27, Financial_Trades_StoppedStock),
		Entry("28 - Works", 28, Financial_Trades_ReOpeningPrints),
		Entry("29 - Works", 29, Financial_Trades_Seller),
		Entry("30 - Works", 30, Financial_Trades_SoldLast),
		Entry("31 - Works", 31, Financial_Trades_SoldLastAndStoppedStock),
		Entry("32 - Works", 32, Financial_Trades_SoldOut),
		Entry("33 - Works", 33, Financial_Trades_SoldOutOfSequence),
		Entry("34 - Works", 34, Financial_Trades_SplitTrade),
		Entry("35 - Works", 35, Financial_Trades_StockOption),
		Entry("36 - Works", 36, Financial_Trades_YellowFlagRegularTrade),
		Entry("37 - Works", 37, Financial_Trades_OddLotTrade),
		Entry("38 - Works", 38, Financial_Trades_CorrectedConsolidatedClose),
		Entry("39 - Works", 39, Financial_Trades_Unknown),
		Entry("40 - Works", 40, Financial_Trades_Held),
		Entry("41 - Works", 41, Financial_Trades_TradeThruExempt),
		Entry("42 - Works", 42, Financial_Trades_NonEligible),
		Entry("43 - Works", 43, Financial_Trades_NonEligibleExtended),
		Entry("44 - Works", 44, Financial_Trades_Cancelled),
		Entry("45 - Works", 45, Financial_Trades_Recovery),
		Entry("46 - Works", 46, Financial_Trades_Correction),
		Entry("47 - Works", 47, Financial_Trades_AsOf),
		Entry("48 - Works", 48, Financial_Trades_AsOfCorrection),
		Entry("49 - Works", 49, Financial_Trades_AsOfCancel),
		Entry("50 - Works", 50, Financial_Trades_OOB),
		Entry("51 - Works", 51, Financial_Trades_Summary),
		Entry("52 - Works", 52, Financial_Trades_ContingentTrade),
		Entry("53 - Works", 53, Financial_Trades_QualifiedContingentTrade),
		Entry("54 - Works", 54, Financial_Trades_Errored),
		Entry("55 - Works", 55, Financial_Trades_OpeningReopeningTradeDetail),
		Entry("56 - Works", 56, Financial_Trades_Placeholder),
		Entry("57 - Works", 57, Financial_Trades_ShortSaleRestrictionActivated),
		Entry("58 - Works", 58, Financial_Trades_ShortSaleRestrictionContinued),
		Entry("59 - Works", 59, Financial_Trades_ShortSaleRestrictionDeactivated),
		Entry("60 - Works", 60, Financial_Trades_ShortSaleRestrictionInEffect),
		Entry("62 - Works", 62, Financial_Trades_FinancialStatusBankrupt),
		Entry("63 - Works", 63, Financial_Trades_FinancialStatusDeficient),
		Entry("64 - Works", 64, Financial_Trades_FinancialStatusDelinquent),
		Entry("65 - Works", 65, Financial_Trades_FinancialStatusBankruptAndDeficient),
		Entry("66 - Works", 66, Financial_Trades_FinancialStatusBankruptAndDelinquent),
		Entry("67 - Works", 67, Financial_Trades_FinancialStatusDeficientAndDelinquent),
		Entry("68 - Works", 68, Financial_Trades_FinancialStatusDeficientDelinquentBankrupt),
		Entry("69 - Works", 69, Financial_Trades_FinancialStatusLiquidation),
		Entry("70 - Works", 70, Financial_Trades_FinancialStatusCreationsSuspended),
		Entry("71 - Works", 71, Financial_Trades_FinancialStatusRedemptionsSuspended),
		Entry("201 - Works", 201, Financial_Trades_Canceled),
		Entry("202 - Works", 202, Financial_Trades_LateAndOutOfSequence),
		Entry("203 - Works", 203, Financial_Trades_LastAndCanceled),
		Entry("204 - Works", 204, Financial_Trades_Late),
		Entry("205 - Works", 205, Financial_Trades_OpeningTradeAndCanceled),
		Entry("206 - Works", 206, Financial_Trades_OpeningTradeLateAndOutOfSequence),
		Entry("207 - Works", 207, Financial_Trades_OnlyTradeAndCanceled),
		Entry("208 - Works", 208, Financial_Trades_OpeningTradeAndLate),
		Entry("209 - Works", 209, Financial_Trades_AutomaticExecutionOption),
		Entry("210 - Works", 210, Financial_Trades_ReopeningTrade),
		Entry("219 - Works", 219, Financial_Trades_IntermarketSweepOrder),
		Entry("227 - Works", 227, Financial_Trades_SingleLegAuctionNonISO),
		Entry("228 - Works", 228, Financial_Trades_SingleLegAuctionISO),
		Entry("229 - Works", 229, Financial_Trades_SingleLegCrossNonISO),
		Entry("230 - Works", 230, Financial_Trades_SingleLegCrossISO),
		Entry("231 - Works", 231, Financial_Trades_SingleLegFloorTrade),
		Entry("232 - Works", 232, Financial_Trades_MultiLegAutoElectronicTrade),
		Entry("233 - Works", 233, Financial_Trades_MultiLegAuction),
		Entry("234 - Works", 234, Financial_Trades_MultiLegCross),
		Entry("235 - Works", 235, Financial_Trades_MultiLegFloorTrade),
		Entry("236 - Works", 236, Financial_Trades_MultiLegAutoElectronicTradeAgainstSingleLeg),
		Entry("237 - Works", 237, Financial_Trades_StockOptionsAuction),
		Entry("238 - Works", 238, Financial_Trades_MultiLegAuctionAgainstSingleLeg),
		Entry("239 - Works", 239, Financial_Trades_MultiLegFloorTradeAgainstSingleLeg),
		Entry("240 - Works", 240, Financial_Trades_StockOptionsAutoElectronicTrade),
		Entry("241 - Works", 241, Financial_Trades_StockOptionsCross),
		Entry("242 - Works", 242, Financial_Trades_StockOptionsFloorTrade),
		Entry("243 - Works", 243, Financial_Trades_StockOptionsAutoElectronicTradeAgainstSingleLeg),
		Entry("244 - Works", 244, Financial_Trades_StockOptionsAuctionAgainstSingleLeg),
		Entry("245 - Works", 245, Financial_Trades_StockOptionsFloorTradeAgainstSingleLeg),
		Entry("246 - Works", 246, Financial_Trades_MultiLegFloorTradeOfProprietaryProducts),
		Entry("247 - Works", 247, Financial_Trades_MultilateralCompressionTradeOfProprietaryProducts),
		Entry("248 - Works", 248, Financial_Trades_ExtendedHoursTrade))

	// Test that attempting to deserialize a Financial.Trades.Condition will fial and return an
	// error if the value cannot be converted to either the name value or integer value
	// of the enum option
	It("UnmarshalCSV - Value is empty - Error", func() {

		// Attempt to convert a fake string value into a Financial.Trades.Condition
		// This should return an error
		enum := new(Financial_Trades_Condition)
		err := enum.UnmarshalCSV("")

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"\" cannot be mapped to a gopb.Financial_Trades_Condition"))
	})

	// Test the conditions under which values should be convertible to a Financial.Trades.Condition
	DescribeTable("UnmarshalCSV Tests",
		func(value string, shouldBe Financial_Trades_Condition) {

			// Attempt to convert the value into a Financial.Trades.Condition
			// This should not fail
			var enum Financial_Trades_Condition
			err := enum.UnmarshalCSV(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("RegularSale - Works", "RegularSale", Financial_Trades_RegularSale),
		Entry("Acquisition - Works", "Acquisition", Financial_Trades_Acquisition),
		Entry("AveragePriceTrade - Works", "AveragePriceTrade", Financial_Trades_AveragePriceTrade),
		Entry("AutomaticExecution - Works", "AutomaticExecution", Financial_Trades_AutomaticExecution),
		Entry("BunchedTrade - Works", "BunchedTrade", Financial_Trades_BunchedTrade),
		Entry("BunchedSoldTrade - Works", "BunchedSoldTrade", Financial_Trades_BunchedSoldTrade),
		Entry("CAPElection - Works", "CAPElection", Financial_Trades_CAPElection),
		Entry("CashSale - Works", "CashSale", Financial_Trades_CashSale),
		Entry("ClosingPrints - Works", "ClosingPrints", Financial_Trades_ClosingPrints),
		Entry("CrossTrade - Works", "CrossTrade", Financial_Trades_CrossTrade),
		Entry("DerivativelyPriced - Works", "DerivativelyPriced", Financial_Trades_DerivativelyPriced),
		Entry("Distribution - Works", "Distribution", Financial_Trades_Distribution),
		Entry("FormT - Works", "FormT", Financial_Trades_FormT),
		Entry("ExtendedTradingHours - Works", "ExtendedTradingHours", Financial_Trades_ExtendedTradingHours),
		Entry("IntermarketSweep - Works", "IntermarketSweep", Financial_Trades_IntermarketSweep),
		Entry("MarketCenterOfficialClose - Works", "MarketCenterOfficialClose", Financial_Trades_MarketCenterOfficialClose),
		Entry("MarketCenterOfficialOpen - Works", "MarketCenterOfficialOpen", Financial_Trades_MarketCenterOfficialOpen),
		Entry("MarketCenterOpeningTrade - Works", "MarketCenterOpeningTrade", Financial_Trades_MarketCenterOpeningTrade),
		Entry("MarketCenterReopeningTrade - Works", "MarketCenterReopeningTrade", Financial_Trades_MarketCenterReopeningTrade),
		Entry("MarketCenterClosingTrade - Works", "MarketCenterClosingTrade", Financial_Trades_MarketCenterClosingTrade),
		Entry("NextDay - Works", "NextDay", Financial_Trades_NextDay),
		Entry("PriceVariationTrade - Works", "PriceVariationTrade", Financial_Trades_PriceVariationTrade),
		Entry("PriorReferencePrice - Works", "PriorReferencePrice", Financial_Trades_PriorReferencePrice),
		Entry("Rule155Trade - Works", "Rule155Trade", Financial_Trades_Rule155Trade),
		Entry("Rule127NYSE - Works", "Rule127NYSE", Financial_Trades_Rule127NYSE),
		Entry("OpeningPrints - Works", "OpeningPrints", Financial_Trades_OpeningPrints),
		Entry("Opened - Works", "Opened", Financial_Trades_Opened),
		Entry("StoppedStock - Works", "StoppedStock", Financial_Trades_StoppedStock),
		Entry("ReOpeningPrints - Works", "ReOpeningPrints", Financial_Trades_ReOpeningPrints),
		Entry("Seller - Works", "Seller", Financial_Trades_Seller),
		Entry("SoldLast - Works", "SoldLast", Financial_Trades_SoldLast),
		Entry("SoldLastAndStoppedStock - Works", "SoldLastAndStoppedStock", Financial_Trades_SoldLastAndStoppedStock),
		Entry("SoldOut - Works", "SoldOut", Financial_Trades_SoldOut),
		Entry("SoldOutOfSequence - Works", "SoldOutOfSequence", Financial_Trades_SoldOutOfSequence),
		Entry("SplitTrade - Works", "SplitTrade", Financial_Trades_SplitTrade),
		Entry("StockOption - Works", "StockOption", Financial_Trades_StockOption),
		Entry("YellowFlagRegularTrade - Works", "YellowFlagRegularTrade", Financial_Trades_YellowFlagRegularTrade),
		Entry("OddLotTrade - Works", "OddLotTrade", Financial_Trades_OddLotTrade),
		Entry("CorrectedConsolidatedClose - Works", "CorrectedConsolidatedClose", Financial_Trades_CorrectedConsolidatedClose),
		Entry("Unknown - Works", "Unknown", Financial_Trades_Unknown),
		Entry("Held - Works", "Held", Financial_Trades_Held),
		Entry("TradeThruExempt - Works", "TradeThruExempt", Financial_Trades_TradeThruExempt),
		Entry("NonEligible - Works", "NonEligible", Financial_Trades_NonEligible),
		Entry("NonEligibleExtended - Works", "NonEligibleExtended", Financial_Trades_NonEligibleExtended),
		Entry("Cancelled - Works", "Cancelled", Financial_Trades_Cancelled),
		Entry("Recovery - Works", "Recovery", Financial_Trades_Recovery),
		Entry("Correction - Works", "Correction", Financial_Trades_Correction),
		Entry("AsOf - Works", "AsOf", Financial_Trades_AsOf),
		Entry("AsOfCorrection - Works", "AsOfCorrection", Financial_Trades_AsOfCorrection),
		Entry("AsOfCancel - Works", "AsOfCancel", Financial_Trades_AsOfCancel),
		Entry("OOB - Works", "OOB", Financial_Trades_OOB),
		Entry("Summary - Works", "Summary", Financial_Trades_Summary),
		Entry("ContingentTrade - Works", "ContingentTrade", Financial_Trades_ContingentTrade),
		Entry("QualifiedContingentTrade - Works", "QualifiedContingentTrade", Financial_Trades_QualifiedContingentTrade),
		Entry("Errored - Works", "Errored", Financial_Trades_Errored),
		Entry("OpeningReopeningTradeDetail - Works", "OpeningReopeningTradeDetail", Financial_Trades_OpeningReopeningTradeDetail),
		Entry("Placeholder - Works", "Placeholder", Financial_Trades_Placeholder),
		Entry("ShortSaleRestrictionActivated - Works",
			"ShortSaleRestrictionActivated", Financial_Trades_ShortSaleRestrictionActivated),
		Entry("ShortSaleRestrictionContinued - Works",
			"ShortSaleRestrictionContinued", Financial_Trades_ShortSaleRestrictionContinued),
		Entry("ShortSaleRestrictionDeactivated - Works",
			"ShortSaleRestrictionDeactivated", Financial_Trades_ShortSaleRestrictionDeactivated),
		Entry("ShortSaleRestrictionInEffect - Works",
			"ShortSaleRestrictionInEffect", Financial_Trades_ShortSaleRestrictionInEffect),
		Entry("FinancialStatusBankrupt - Works", "FinancialStatusBankrupt", Financial_Trades_FinancialStatusBankrupt),
		Entry("FinancialStatusDeficient - Works", "FinancialStatusDeficient", Financial_Trades_FinancialStatusDeficient),
		Entry("FinancialStatusDelinquent - Works", "FinancialStatusDelinquent", Financial_Trades_FinancialStatusDelinquent),
		Entry("FinancialStatusBankruptAndDeficient - Works",
			"FinancialStatusBankruptAndDeficient", Financial_Trades_FinancialStatusBankruptAndDeficient),
		Entry("FinancialStatusBankruptAndDelinquent - Works",
			"FinancialStatusBankruptAndDelinquent", Financial_Trades_FinancialStatusBankruptAndDelinquent),
		Entry("FinancialStatusDeficientAndDelinquent - Works",
			"FinancialStatusDeficientAndDelinquent", Financial_Trades_FinancialStatusDeficientAndDelinquent),
		Entry("FinancialStatusDeficientDelinquentBankrupt - Works",
			"FinancialStatusDeficientDelinquentBankrupt", Financial_Trades_FinancialStatusDeficientDelinquentBankrupt),
		Entry("FinancialStatusLiquidation - Works", "FinancialStatusLiquidation", Financial_Trades_FinancialStatusLiquidation),
		Entry("FinancialStatusCreationsSuspended - Works",
			"FinancialStatusCreationsSuspended", Financial_Trades_FinancialStatusCreationsSuspended),
		Entry("FinancialStatusRedemptionsSuspended - Works",
			"FinancialStatusRedemptionsSuspended", Financial_Trades_FinancialStatusRedemptionsSuspended),
		Entry("Canceled - Works", "Canceled", Financial_Trades_Canceled),
		Entry("LateAndOutOfSequence - Works", "LateAndOutOfSequence", Financial_Trades_LateAndOutOfSequence),
		Entry("LastAndCanceled - Works", "LastAndCanceled", Financial_Trades_LastAndCanceled),
		Entry("Late - Works", "Late", Financial_Trades_Late),
		Entry("OpeningTradeAndCanceled - Works", "OpeningTradeAndCanceled", Financial_Trades_OpeningTradeAndCanceled),
		Entry("OpeningTradeLateAndOutOfSequence - Works",
			"OpeningTradeLateAndOutOfSequence", Financial_Trades_OpeningTradeLateAndOutOfSequence),
		Entry("OnlyTradeAndCanceled - Works", "OnlyTradeAndCanceled", Financial_Trades_OnlyTradeAndCanceled),
		Entry("OpeningTradeAndLate - Works", "OpeningTradeAndLate", Financial_Trades_OpeningTradeAndLate),
		Entry("AutomaticExecutionOption - Works", "AutomaticExecutionOption", Financial_Trades_AutomaticExecutionOption),
		Entry("ReopeningTrade - Works", "ReopeningTrade", Financial_Trades_ReopeningTrade),
		Entry("IntermarketSweepOrder - Works", "IntermarketSweepOrder", Financial_Trades_IntermarketSweepOrder),
		Entry("SingleLegAuctionNonISO - Works", "SingleLegAuctionNonISO", Financial_Trades_SingleLegAuctionNonISO),
		Entry("SingleLegAuctionISO - Works", "SingleLegAuctionISO", Financial_Trades_SingleLegAuctionISO),
		Entry("SingleLegCrossNonISO - Works", "SingleLegCrossNonISO", Financial_Trades_SingleLegCrossNonISO),
		Entry("SingleLegCrossISO - Works", "SingleLegCrossISO", Financial_Trades_SingleLegCrossISO),
		Entry("SingleLegFloorTrade - Works", "SingleLegFloorTrade", Financial_Trades_SingleLegFloorTrade),
		Entry("MultiLegAutoElectronicTrade - Works",
			"MultiLegAutoElectronicTrade", Financial_Trades_MultiLegAutoElectronicTrade),
		Entry("MultiLegAuction - Works", "MultiLegAuction", Financial_Trades_MultiLegAuction),
		Entry("MultiLegCross - Works", "MultiLegCross", Financial_Trades_MultiLegCross),
		Entry("MultiLegFloorTrade - Works", "MultiLegFloorTrade", Financial_Trades_MultiLegFloorTrade),
		Entry("MultiLegAutoElectronicTradeAgainstSingleLeg - Works",
			"MultiLegAutoElectronicTradeAgainstSingleLeg", Financial_Trades_MultiLegAutoElectronicTradeAgainstSingleLeg),
		Entry("StockOptionsAuction - Works", "StockOptionsAuction", Financial_Trades_StockOptionsAuction),
		Entry("MultiLegAuctionAgainstSingleLeg - Works",
			"MultiLegAuctionAgainstSingleLeg", Financial_Trades_MultiLegAuctionAgainstSingleLeg),
		Entry("MultiLegFloorTradeAgainstSingleLeg - Works",
			"MultiLegFloorTradeAgainstSingleLeg", Financial_Trades_MultiLegFloorTradeAgainstSingleLeg),
		Entry("StockOptionsAutoElectronicTrade - Works",
			"StockOptionsAutoElectronicTrade", Financial_Trades_StockOptionsAutoElectronicTrade),
		Entry("StockOptionsCross - Works", "StockOptionsCross", Financial_Trades_StockOptionsCross),
		Entry("StockOptionsFloorTrade - Works", "StockOptionsFloorTrade", Financial_Trades_StockOptionsFloorTrade),
		Entry("StockOptionsAutoElectronicTradeAgainstSingleLeg - Works",
			"StockOptionsAutoElectronicTradeAgainstSingleLeg", Financial_Trades_StockOptionsAutoElectronicTradeAgainstSingleLeg),
		Entry("StockOptionsAuctionAgainstSingleLeg - Works",
			"StockOptionsAuctionAgainstSingleLeg", Financial_Trades_StockOptionsAuctionAgainstSingleLeg),
		Entry("StockOptionsFloorTradeAgainstSingleLeg - Works",
			"StockOptionsFloorTradeAgainstSingleLeg", Financial_Trades_StockOptionsFloorTradeAgainstSingleLeg),
		Entry("MultiLegFloorTradeOfProprietaryProducts - Works",
			"MultiLegFloorTradeOfProprietaryProducts", Financial_Trades_MultiLegFloorTradeOfProprietaryProducts),
		Entry("MultilateralCompressionTradeOfProprietaryProducts - Works",
			"MultilateralCompressionTradeOfProprietaryProducts", Financial_Trades_MultilateralCompressionTradeOfProprietaryProducts),
		Entry("ExtendedHoursTrade - Works", "ExtendedHoursTrade", Financial_Trades_ExtendedHoursTrade),
		Entry("0 - Works", "0", Financial_Trades_RegularSale),
		Entry("1 - Works", "1", Financial_Trades_Acquisition),
		Entry("2 - Works", "2", Financial_Trades_AveragePriceTrade),
		Entry("3 - Works", "3", Financial_Trades_AutomaticExecution),
		Entry("4 - Works", "4", Financial_Trades_BunchedTrade),
		Entry("5 - Works", "5", Financial_Trades_BunchedSoldTrade),
		Entry("6 - Works", "6", Financial_Trades_CAPElection),
		Entry("7 - Works", "7", Financial_Trades_CashSale),
		Entry("8 - Works", "8", Financial_Trades_ClosingPrints),
		Entry("9 - Works", "9", Financial_Trades_CrossTrade),
		Entry("10 - Works", "10", Financial_Trades_DerivativelyPriced),
		Entry("11 - Works", "11", Financial_Trades_Distribution),
		Entry("12 - Works", "12", Financial_Trades_FormT),
		Entry("13 - Works", "13", Financial_Trades_ExtendedTradingHours),
		Entry("14 - Works", "14", Financial_Trades_IntermarketSweep),
		Entry("15 - Works", "15", Financial_Trades_MarketCenterOfficialClose),
		Entry("16 - Works", "16", Financial_Trades_MarketCenterOfficialOpen),
		Entry("17 - Works", "17", Financial_Trades_MarketCenterOpeningTrade),
		Entry("18 - Works", "18", Financial_Trades_MarketCenterReopeningTrade),
		Entry("19 - Works", "19", Financial_Trades_MarketCenterClosingTrade),
		Entry("20 - Works", "20", Financial_Trades_NextDay),
		Entry("21 - Works", "21", Financial_Trades_PriceVariationTrade),
		Entry("22 - Works", "22", Financial_Trades_PriorReferencePrice),
		Entry("23 - Works", "23", Financial_Trades_Rule155Trade),
		Entry("24 - Works", "24", Financial_Trades_Rule127NYSE),
		Entry("25 - Works", "25", Financial_Trades_OpeningPrints),
		Entry("26 - Works", "26", Financial_Trades_Opened),
		Entry("27 - Works", "27", Financial_Trades_StoppedStock),
		Entry("28 - Works", "28", Financial_Trades_ReOpeningPrints),
		Entry("29 - Works", "29", Financial_Trades_Seller),
		Entry("30 - Works", "30", Financial_Trades_SoldLast),
		Entry("31 - Works", "31", Financial_Trades_SoldLastAndStoppedStock),
		Entry("32 - Works", "32", Financial_Trades_SoldOut),
		Entry("33 - Works", "33", Financial_Trades_SoldOutOfSequence),
		Entry("34 - Works", "34", Financial_Trades_SplitTrade),
		Entry("35 - Works", "35", Financial_Trades_StockOption),
		Entry("36 - Works", "36", Financial_Trades_YellowFlagRegularTrade),
		Entry("37 - Works", "37", Financial_Trades_OddLotTrade),
		Entry("38 - Works", "38", Financial_Trades_CorrectedConsolidatedClose),
		Entry("39 - Works", "39", Financial_Trades_Unknown),
		Entry("40 - Works", "40", Financial_Trades_Held),
		Entry("41 - Works", "41", Financial_Trades_TradeThruExempt),
		Entry("42 - Works", "42", Financial_Trades_NonEligible),
		Entry("43 - Works", "43", Financial_Trades_NonEligibleExtended),
		Entry("44 - Works", "44", Financial_Trades_Cancelled),
		Entry("45 - Works", "45", Financial_Trades_Recovery),
		Entry("46 - Works", "46", Financial_Trades_Correction),
		Entry("47 - Works", "47", Financial_Trades_AsOf),
		Entry("48 - Works", "48", Financial_Trades_AsOfCorrection),
		Entry("49 - Works", "49", Financial_Trades_AsOfCancel),
		Entry("50 - Works", "50", Financial_Trades_OOB),
		Entry("51 - Works", "51", Financial_Trades_Summary),
		Entry("52 - Works", "52", Financial_Trades_ContingentTrade),
		Entry("53 - Works", "53", Financial_Trades_QualifiedContingentTrade),
		Entry("54 - Works", "54", Financial_Trades_Errored),
		Entry("55 - Works", "55", Financial_Trades_OpeningReopeningTradeDetail),
		Entry("56 - Works", "56", Financial_Trades_Placeholder),
		Entry("57 - Works", "57", Financial_Trades_ShortSaleRestrictionActivated),
		Entry("58 - Works", "58", Financial_Trades_ShortSaleRestrictionContinued),
		Entry("59 - Works", "59", Financial_Trades_ShortSaleRestrictionDeactivated),
		Entry("60 - Works", "60", Financial_Trades_ShortSaleRestrictionInEffect),
		Entry("62 - Works", "62", Financial_Trades_FinancialStatusBankrupt),
		Entry("63 - Works", "63", Financial_Trades_FinancialStatusDeficient),
		Entry("64 - Works", "64", Financial_Trades_FinancialStatusDelinquent),
		Entry("65 - Works", "65", Financial_Trades_FinancialStatusBankruptAndDeficient),
		Entry("66 - Works", "66", Financial_Trades_FinancialStatusBankruptAndDelinquent),
		Entry("67 - Works", "67", Financial_Trades_FinancialStatusDeficientAndDelinquent),
		Entry("68 - Works", "68", Financial_Trades_FinancialStatusDeficientDelinquentBankrupt),
		Entry("69 - Works", "69", Financial_Trades_FinancialStatusLiquidation),
		Entry("70 - Works", "70", Financial_Trades_FinancialStatusCreationsSuspended),
		Entry("71 - Works", "71", Financial_Trades_FinancialStatusRedemptionsSuspended),
		Entry("201 - Works", "201", Financial_Trades_Canceled),
		Entry("202 - Works", "202", Financial_Trades_LateAndOutOfSequence),
		Entry("203 - Works", "203", Financial_Trades_LastAndCanceled),
		Entry("204 - Works", "204", Financial_Trades_Late),
		Entry("205 - Works", "205", Financial_Trades_OpeningTradeAndCanceled),
		Entry("206 - Works", "206", Financial_Trades_OpeningTradeLateAndOutOfSequence),
		Entry("207 - Works", "207", Financial_Trades_OnlyTradeAndCanceled),
		Entry("208 - Works", "208", Financial_Trades_OpeningTradeAndLate),
		Entry("209 - Works", "209", Financial_Trades_AutomaticExecutionOption),
		Entry("210 - Works", "210", Financial_Trades_ReopeningTrade),
		Entry("219 - Works", "219", Financial_Trades_IntermarketSweepOrder),
		Entry("227 - Works", "227", Financial_Trades_SingleLegAuctionNonISO),
		Entry("228 - Works", "228", Financial_Trades_SingleLegAuctionISO),
		Entry("229 - Works", "229", Financial_Trades_SingleLegCrossNonISO),
		Entry("230 - Works", "230", Financial_Trades_SingleLegCrossISO),
		Entry("231 - Works", "231", Financial_Trades_SingleLegFloorTrade),
		Entry("232 - Works", "232", Financial_Trades_MultiLegAutoElectronicTrade),
		Entry("233 - Works", "233", Financial_Trades_MultiLegAuction),
		Entry("234 - Works", "234", Financial_Trades_MultiLegCross),
		Entry("235 - Works", "235", Financial_Trades_MultiLegFloorTrade),
		Entry("236 - Works", "236", Financial_Trades_MultiLegAutoElectronicTradeAgainstSingleLeg),
		Entry("237 - Works", "237", Financial_Trades_StockOptionsAuction),
		Entry("238 - Works", "238", Financial_Trades_MultiLegAuctionAgainstSingleLeg),
		Entry("239 - Works", "239", Financial_Trades_MultiLegFloorTradeAgainstSingleLeg),
		Entry("240 - Works", "240", Financial_Trades_StockOptionsAutoElectronicTrade),
		Entry("241 - Works", "241", Financial_Trades_StockOptionsCross),
		Entry("242 - Works", "242", Financial_Trades_StockOptionsFloorTrade),
		Entry("243 - Works", "243", Financial_Trades_StockOptionsAutoElectronicTradeAgainstSingleLeg),
		Entry("244 - Works", "244", Financial_Trades_StockOptionsAuctionAgainstSingleLeg),
		Entry("245 - Works", "245", Financial_Trades_StockOptionsFloorTradeAgainstSingleLeg),
		Entry("246 - Works", "246", Financial_Trades_MultiLegFloorTradeOfProprietaryProducts),
		Entry("247 - Works", "247", Financial_Trades_MultilateralCompressionTradeOfProprietaryProducts),
		Entry("248 - Works", "248", Financial_Trades_ExtendedHoursTrade))

	// Tests that, if the attribute type submitted to UnmarshalDynamoDBAttributeValue is not one we
	// recognize, then the function will return an error
	It("UnmarshalDynamoDBAttributeValue - AttributeValue type invalid - Error", func() {
		value := new(Financial_Trades_Condition)
		err := attributevalue.Unmarshal(&types.AttributeValueMemberBOOL{Value: true}, &value)
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("Attribute value of *types.AttributeValueMemberBOOL could not be converted to a Financial.Trades.Condition"))
	})

	// Tests the conditions under which UnmarshalDynamoDBAttributeValue is called and no error is generated
	DescribeTable("UnmarshalDynamoDBAttributeValue - AttributeValue Conditions",
		func(raw types.AttributeValue, expected Financial_Trades_Condition) {
			var value Financial_Trades_Condition
			err := attributevalue.Unmarshal(raw, &value)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(expected))
		},
		Entry("Value is []bytes, RegularSale - Works",
			&types.AttributeValueMemberB{Value: []byte("RegularSale")}, Financial_Trades_RegularSale),
		Entry("Value is []bytes, Acquisition - Works",
			&types.AttributeValueMemberB{Value: []byte("Acquisition")}, Financial_Trades_Acquisition),
		Entry("Value is []bytes, AveragePriceTrade - Works",
			&types.AttributeValueMemberB{Value: []byte("AveragePriceTrade")}, Financial_Trades_AveragePriceTrade),
		Entry("Value is []bytes, AutomaticExecution - Works",
			&types.AttributeValueMemberB{Value: []byte("AutomaticExecution")}, Financial_Trades_AutomaticExecution),
		Entry("Value is []bytes, BunchedTrade - Works",
			&types.AttributeValueMemberB{Value: []byte("BunchedTrade")}, Financial_Trades_BunchedTrade),
		Entry("Value is []bytes, BunchedSoldTrade - Works",
			&types.AttributeValueMemberB{Value: []byte("BunchedSoldTrade")}, Financial_Trades_BunchedSoldTrade),
		Entry("Value is []bytes, CAPElection - Works",
			&types.AttributeValueMemberB{Value: []byte("CAPElection")}, Financial_Trades_CAPElection),
		Entry("Value is []bytes, CashSale - Works",
			&types.AttributeValueMemberB{Value: []byte("CashSale")}, Financial_Trades_CashSale),
		Entry("Value is []bytes, ClosingPrints - Works",
			&types.AttributeValueMemberB{Value: []byte("ClosingPrints")}, Financial_Trades_ClosingPrints),
		Entry("Value is []bytes, CrossTrade - Works",
			&types.AttributeValueMemberB{Value: []byte("CrossTrade")}, Financial_Trades_CrossTrade),
		Entry("Value is []bytes, DerivativelyPriced - Works",
			&types.AttributeValueMemberB{Value: []byte("DerivativelyPriced")}, Financial_Trades_DerivativelyPriced),
		Entry("Value is []bytes, Distribution - Works",
			&types.AttributeValueMemberB{Value: []byte("Distribution")}, Financial_Trades_Distribution),
		Entry("Value is []bytes, FormT - Works",
			&types.AttributeValueMemberB{Value: []byte("FormT")}, Financial_Trades_FormT),
		Entry("Value is []bytes, ExtendedTradingHours - Works",
			&types.AttributeValueMemberB{Value: []byte("ExtendedTradingHours")}, Financial_Trades_ExtendedTradingHours),
		Entry("Value is []bytes, IntermarketSweep - Works",
			&types.AttributeValueMemberB{Value: []byte("IntermarketSweep")}, Financial_Trades_IntermarketSweep),
		Entry("Value is []bytes, MarketCenterOfficialClose - Works",
			&types.AttributeValueMemberB{Value: []byte("MarketCenterOfficialClose")}, Financial_Trades_MarketCenterOfficialClose),
		Entry("Value is []bytes, MarketCenterOfficialOpen - Works",
			&types.AttributeValueMemberB{Value: []byte("MarketCenterOfficialOpen")}, Financial_Trades_MarketCenterOfficialOpen),
		Entry("Value is []bytes, MarketCenterOpeningTrade - Works",
			&types.AttributeValueMemberB{Value: []byte("MarketCenterOpeningTrade")}, Financial_Trades_MarketCenterOpeningTrade),
		Entry("Value is []bytes, MarketCenterReopeningTrade - Works",
			&types.AttributeValueMemberB{Value: []byte("MarketCenterReopeningTrade")}, Financial_Trades_MarketCenterReopeningTrade),
		Entry("Value is []bytes, MarketCenterClosingTrade - Works",
			&types.AttributeValueMemberB{Value: []byte("MarketCenterClosingTrade")}, Financial_Trades_MarketCenterClosingTrade),
		Entry("Value is []bytes, NextDay - Works",
			&types.AttributeValueMemberB{Value: []byte("NextDay")}, Financial_Trades_NextDay),
		Entry("Value is []bytes, PriceVariationTrade - Works",
			&types.AttributeValueMemberB{Value: []byte("PriceVariationTrade")}, Financial_Trades_PriceVariationTrade),
		Entry("Value is []bytes, PriorReferencePrice - Works",
			&types.AttributeValueMemberB{Value: []byte("PriorReferencePrice")}, Financial_Trades_PriorReferencePrice),
		Entry("Value is []bytes, Rule155Trade - Works",
			&types.AttributeValueMemberB{Value: []byte("Rule155Trade")}, Financial_Trades_Rule155Trade),
		Entry("Value is []bytes, Rule127NYSE - Works",
			&types.AttributeValueMemberB{Value: []byte("Rule127NYSE")}, Financial_Trades_Rule127NYSE),
		Entry("Value is []bytes, OpeningPrints - Works",
			&types.AttributeValueMemberB{Value: []byte("OpeningPrints")}, Financial_Trades_OpeningPrints),
		Entry("Value is []bytes, Opened - Works",
			&types.AttributeValueMemberB{Value: []byte("Opened")}, Financial_Trades_Opened),
		Entry("Value is []bytes, StoppedStock - Works",
			&types.AttributeValueMemberB{Value: []byte("StoppedStock")}, Financial_Trades_StoppedStock),
		Entry("Value is []bytes, ReOpeningPrints - Works",
			&types.AttributeValueMemberB{Value: []byte("ReOpeningPrints")}, Financial_Trades_ReOpeningPrints),
		Entry("Value is []bytes, Seller - Works",
			&types.AttributeValueMemberB{Value: []byte("Seller")}, Financial_Trades_Seller),
		Entry("Value is []bytes, SoldLast - Works",
			&types.AttributeValueMemberB{Value: []byte("SoldLast")}, Financial_Trades_SoldLast),
		Entry("Value is []bytes, SoldLastAndStoppedStock - Works",
			&types.AttributeValueMemberB{Value: []byte("SoldLastAndStoppedStock")}, Financial_Trades_SoldLastAndStoppedStock),
		Entry("Value is []bytes, SoldOut - Works",
			&types.AttributeValueMemberB{Value: []byte("SoldOut")}, Financial_Trades_SoldOut),
		Entry("Value is []bytes, SoldOutOfSequence - Works",
			&types.AttributeValueMemberB{Value: []byte("SoldOutOfSequence")}, Financial_Trades_SoldOutOfSequence),
		Entry("Value is []bytes, SplitTrade - Works",
			&types.AttributeValueMemberB{Value: []byte("SplitTrade")}, Financial_Trades_SplitTrade),
		Entry("Value is []bytes, StockOption - Works",
			&types.AttributeValueMemberB{Value: []byte("StockOption")}, Financial_Trades_StockOption),
		Entry("Value is []bytes, YellowFlagRegularTrade - Works",
			&types.AttributeValueMemberB{Value: []byte("YellowFlagRegularTrade")}, Financial_Trades_YellowFlagRegularTrade),
		Entry("Value is []bytes, OddLotTrade - Works",
			&types.AttributeValueMemberB{Value: []byte("OddLotTrade")}, Financial_Trades_OddLotTrade),
		Entry("Value is []bytes, CorrectedConsolidatedClose - Works",
			&types.AttributeValueMemberB{Value: []byte("CorrectedConsolidatedClose")}, Financial_Trades_CorrectedConsolidatedClose),
		Entry("Value is []bytes, Unknown - Works",
			&types.AttributeValueMemberB{Value: []byte("Unknown")}, Financial_Trades_Unknown),
		Entry("Value is []bytes, Held - Works",
			&types.AttributeValueMemberB{Value: []byte("Held")}, Financial_Trades_Held),
		Entry("Value is []bytes, TradeThruExempt - Works",
			&types.AttributeValueMemberB{Value: []byte("TradeThruExempt")}, Financial_Trades_TradeThruExempt),
		Entry("Value is []bytes, NonEligible - Works",
			&types.AttributeValueMemberB{Value: []byte("NonEligible")}, Financial_Trades_NonEligible),
		Entry("Value is []bytes, NonEligibleExtended - Works",
			&types.AttributeValueMemberB{Value: []byte("NonEligibleExtended")}, Financial_Trades_NonEligibleExtended),
		Entry("Value is []bytes, Cancelled - Works",
			&types.AttributeValueMemberB{Value: []byte("Cancelled")}, Financial_Trades_Cancelled),
		Entry("Value is []bytes, Recovery - Works",
			&types.AttributeValueMemberB{Value: []byte("Recovery")}, Financial_Trades_Recovery),
		Entry("Value is []bytes, Correction - Works",
			&types.AttributeValueMemberB{Value: []byte("Correction")}, Financial_Trades_Correction),
		Entry("Value is []bytes, AsOf - Works",
			&types.AttributeValueMemberB{Value: []byte("AsOf")}, Financial_Trades_AsOf),
		Entry("Value is []bytes, AsOfCorrection - Works",
			&types.AttributeValueMemberB{Value: []byte("AsOfCorrection")}, Financial_Trades_AsOfCorrection),
		Entry("Value is []bytes, AsOfCancel - Works",
			&types.AttributeValueMemberB{Value: []byte("AsOfCancel")}, Financial_Trades_AsOfCancel),
		Entry("Value is []bytes, OOB - Works",
			&types.AttributeValueMemberB{Value: []byte("OOB")}, Financial_Trades_OOB),
		Entry("Value is []bytes, Summary - Works",
			&types.AttributeValueMemberB{Value: []byte("Summary")}, Financial_Trades_Summary),
		Entry("Value is []bytes, ContingentTrade - Works",
			&types.AttributeValueMemberB{Value: []byte("ContingentTrade")}, Financial_Trades_ContingentTrade),
		Entry("Value is []bytes, QualifiedContingentTrade - Works",
			&types.AttributeValueMemberB{Value: []byte("QualifiedContingentTrade")}, Financial_Trades_QualifiedContingentTrade),
		Entry("Value is []bytes, Errored - Works",
			&types.AttributeValueMemberB{Value: []byte("Errored")}, Financial_Trades_Errored),
		Entry("Value is []bytes, OpeningReopeningTradeDetail - Works",
			&types.AttributeValueMemberB{Value: []byte("OpeningReopeningTradeDetail")}, Financial_Trades_OpeningReopeningTradeDetail),
		Entry("Value is []bytes, Placeholder - Works",
			&types.AttributeValueMemberB{Value: []byte("Placeholder")}, Financial_Trades_Placeholder),
		Entry("Value is []bytes, ShortSaleRestrictionActivated - Works",
			&types.AttributeValueMemberB{Value: []byte("ShortSaleRestrictionActivated")}, Financial_Trades_ShortSaleRestrictionActivated),
		Entry("Value is []bytes, ShortSaleRestrictionContinued - Works",
			&types.AttributeValueMemberB{Value: []byte("ShortSaleRestrictionContinued")}, Financial_Trades_ShortSaleRestrictionContinued),
		Entry("Value is []bytes, ShortSaleRestrictionDeactivated - Works",
			&types.AttributeValueMemberB{Value: []byte("ShortSaleRestrictionDeactivated")}, Financial_Trades_ShortSaleRestrictionDeactivated),
		Entry("Value is []bytes, ShortSaleRestrictionInEffect - Works",
			&types.AttributeValueMemberB{Value: []byte("ShortSaleRestrictionInEffect")}, Financial_Trades_ShortSaleRestrictionInEffect),
		Entry("Value is []bytes, FinancialStatusBankrupt - Works",
			&types.AttributeValueMemberB{Value: []byte("FinancialStatusBankrupt")}, Financial_Trades_FinancialStatusBankrupt),
		Entry("Value is []bytes, FinancialStatusDeficient - Works",
			&types.AttributeValueMemberB{Value: []byte("FinancialStatusDeficient")}, Financial_Trades_FinancialStatusDeficient),
		Entry("Value is []bytes, FinancialStatusDelinquent - Works",
			&types.AttributeValueMemberB{Value: []byte("FinancialStatusDelinquent")}, Financial_Trades_FinancialStatusDelinquent),
		Entry("Value is []bytes, FinancialStatusBankruptAndDeficient - Works",
			&types.AttributeValueMemberB{Value: []byte("FinancialStatusBankruptAndDeficient")}, Financial_Trades_FinancialStatusBankruptAndDeficient),
		Entry("Value is []bytes, FinancialStatusBankruptAndDelinquent - Works",
			&types.AttributeValueMemberB{Value: []byte("FinancialStatusBankruptAndDelinquent")}, Financial_Trades_FinancialStatusBankruptAndDelinquent),
		Entry("Value is []bytes, FinancialStatusDeficientAndDelinquent - Works",
			&types.AttributeValueMemberB{Value: []byte("FinancialStatusDeficientAndDelinquent")}, Financial_Trades_FinancialStatusDeficientAndDelinquent),
		Entry("Value is []bytes, FinancialStatusDeficientDelinquentBankrupt - Works",
			&types.AttributeValueMemberB{Value: []byte("FinancialStatusDeficientDelinquentBankrupt")}, Financial_Trades_FinancialStatusDeficientDelinquentBankrupt),
		Entry("Value is []bytes, FinancialStatusLiquidation - Works",
			&types.AttributeValueMemberB{Value: []byte("FinancialStatusLiquidation")}, Financial_Trades_FinancialStatusLiquidation),
		Entry("Value is []bytes, FinancialStatusCreationsSuspended - Works",
			&types.AttributeValueMemberB{Value: []byte("FinancialStatusCreationsSuspended")}, Financial_Trades_FinancialStatusCreationsSuspended),
		Entry("Value is []bytes, FinancialStatusRedemptionsSuspended - Works",
			&types.AttributeValueMemberB{Value: []byte("FinancialStatusRedemptionsSuspended")}, Financial_Trades_FinancialStatusRedemptionsSuspended),
		Entry("Value is []bytes, Canceled - Works",
			&types.AttributeValueMemberB{Value: []byte("Canceled")}, Financial_Trades_Canceled),
		Entry("Value is []bytes, LateAndOutOfSequence - Works",
			&types.AttributeValueMemberB{Value: []byte("LateAndOutOfSequence")}, Financial_Trades_LateAndOutOfSequence),
		Entry("Value is []bytes, LastAndCanceled - Works",
			&types.AttributeValueMemberB{Value: []byte("LastAndCanceled")}, Financial_Trades_LastAndCanceled),
		Entry("Value is []bytes, Late - Works",
			&types.AttributeValueMemberB{Value: []byte("Late")}, Financial_Trades_Late),
		Entry("Value is []bytes, OpeningTradeAndCanceled - Works",
			&types.AttributeValueMemberB{Value: []byte("OpeningTradeAndCanceled")}, Financial_Trades_OpeningTradeAndCanceled),
		Entry("Value is []bytes, OpeningTradeLateAndOutOfSequence - Works",
			&types.AttributeValueMemberB{Value: []byte("OpeningTradeLateAndOutOfSequence")}, Financial_Trades_OpeningTradeLateAndOutOfSequence),
		Entry("Value is []bytes, OnlyTradeAndCanceled - Works",
			&types.AttributeValueMemberB{Value: []byte("OnlyTradeAndCanceled")}, Financial_Trades_OnlyTradeAndCanceled),
		Entry("Value is []bytes, OpeningTradeAndLate - Works",
			&types.AttributeValueMemberB{Value: []byte("OpeningTradeAndLate")}, Financial_Trades_OpeningTradeAndLate),
		Entry("Value is []bytes, AutomaticExecutionOption - Works",
			&types.AttributeValueMemberB{Value: []byte("AutomaticExecutionOption")}, Financial_Trades_AutomaticExecutionOption),
		Entry("Value is []bytes, ReopeningTrade - Works",
			&types.AttributeValueMemberB{Value: []byte("ReopeningTrade")}, Financial_Trades_ReopeningTrade),
		Entry("Value is []bytes, IntermarketSweepOrder - Works",
			&types.AttributeValueMemberB{Value: []byte("IntermarketSweepOrder")}, Financial_Trades_IntermarketSweepOrder),
		Entry("Value is []bytes, SingleLegAuctionNonISO - Works",
			&types.AttributeValueMemberB{Value: []byte("SingleLegAuctionNonISO")}, Financial_Trades_SingleLegAuctionNonISO),
		Entry("Value is []bytes, SingleLegAuctionISO - Works",
			&types.AttributeValueMemberB{Value: []byte("SingleLegAuctionISO")}, Financial_Trades_SingleLegAuctionISO),
		Entry("Value is []bytes, SingleLegCrossNonISO - Works",
			&types.AttributeValueMemberB{Value: []byte("SingleLegCrossNonISO")}, Financial_Trades_SingleLegCrossNonISO),
		Entry("Value is []bytes, SingleLegCrossISO - Works",
			&types.AttributeValueMemberB{Value: []byte("SingleLegCrossISO")}, Financial_Trades_SingleLegCrossISO),
		Entry("Value is []bytes, SingleLegFloorTrade - Works",
			&types.AttributeValueMemberB{Value: []byte("SingleLegFloorTrade")}, Financial_Trades_SingleLegFloorTrade),
		Entry("Value is []bytes, MultiLegAutoElectronicTrade - Works",
			&types.AttributeValueMemberB{Value: []byte("MultiLegAutoElectronicTrade")}, Financial_Trades_MultiLegAutoElectronicTrade),
		Entry("Value is []bytes, MultiLegAuction - Works",
			&types.AttributeValueMemberB{Value: []byte("MultiLegAuction")}, Financial_Trades_MultiLegAuction),
		Entry("Value is []bytes, MultiLegCross - Works",
			&types.AttributeValueMemberB{Value: []byte("MultiLegCross")}, Financial_Trades_MultiLegCross),
		Entry("Value is []bytes, MultiLegFloorTrade - Works",
			&types.AttributeValueMemberB{Value: []byte("MultiLegFloorTrade")}, Financial_Trades_MultiLegFloorTrade),
		Entry("Value is []bytes, MultiLegAutoElectronicTradeAgainstSingleLeg - Works",
			&types.AttributeValueMemberB{Value: []byte("MultiLegAutoElectronicTradeAgainstSingleLeg")}, Financial_Trades_MultiLegAutoElectronicTradeAgainstSingleLeg),
		Entry("Value is []bytes, StockOptionsAuction - Works",
			&types.AttributeValueMemberB{Value: []byte("StockOptionsAuction")}, Financial_Trades_StockOptionsAuction),
		Entry("Value is []bytes, MultiLegAuctionAgainstSingleLeg - Works",
			&types.AttributeValueMemberB{Value: []byte("MultiLegAuctionAgainstSingleLeg")}, Financial_Trades_MultiLegAuctionAgainstSingleLeg),
		Entry("Value is []bytes, MultiLegFloorTradeAgainstSingleLeg - Works",
			&types.AttributeValueMemberB{Value: []byte("MultiLegFloorTradeAgainstSingleLeg")}, Financial_Trades_MultiLegFloorTradeAgainstSingleLeg),
		Entry("Value is []bytes, StockOptionsAutoElectronicTrade - Works",
			&types.AttributeValueMemberB{Value: []byte("StockOptionsAutoElectronicTrade")}, Financial_Trades_StockOptionsAutoElectronicTrade),
		Entry("Value is []bytes, StockOptionsCross - Works",
			&types.AttributeValueMemberB{Value: []byte("StockOptionsCross")}, Financial_Trades_StockOptionsCross),
		Entry("Value is []bytes, StockOptionsFloorTrade - Works",
			&types.AttributeValueMemberB{Value: []byte("StockOptionsFloorTrade")}, Financial_Trades_StockOptionsFloorTrade),
		Entry("Value is []bytes, StockOptionsAutoElectronicTradeAgainstSingleLeg - Works",
			&types.AttributeValueMemberB{Value: []byte("StockOptionsAutoElectronicTradeAgainstSingleLeg")}, Financial_Trades_StockOptionsAutoElectronicTradeAgainstSingleLeg),
		Entry("Value is []bytes, StockOptionsAuctionAgainstSingleLeg - Works",
			&types.AttributeValueMemberB{Value: []byte("StockOptionsAuctionAgainstSingleLeg")}, Financial_Trades_StockOptionsAuctionAgainstSingleLeg),
		Entry("Value is []bytes, StockOptionsFloorTradeAgainstSingleLeg - Works",
			&types.AttributeValueMemberB{Value: []byte("StockOptionsFloorTradeAgainstSingleLeg")}, Financial_Trades_StockOptionsFloorTradeAgainstSingleLeg),
		Entry("Value is []bytes, MultiLegFloorTradeOfProprietaryProducts - Works",
			&types.AttributeValueMemberB{Value: []byte("MultiLegFloorTradeOfProprietaryProducts")}, Financial_Trades_MultiLegFloorTradeOfProprietaryProducts),
		Entry("Value is []bytes, MultilateralCompressionTradeOfProprietaryProducts - Works",
			&types.AttributeValueMemberB{Value: []byte("MultilateralCompressionTradeOfProprietaryProducts")}, Financial_Trades_MultilateralCompressionTradeOfProprietaryProducts),
		Entry("Value is []bytes, ExtendedHoursTrade - Works",
			&types.AttributeValueMemberB{Value: []byte("ExtendedHoursTrade")}, Financial_Trades_ExtendedHoursTrade),
		Entry("Value is numeric, 0 - Works",
			&types.AttributeValueMemberN{Value: "0"}, Financial_Trades_RegularSale),
		Entry("Value is numeric, 1 - Works",
			&types.AttributeValueMemberN{Value: "1"}, Financial_Trades_Acquisition),
		Entry("Value is numeric, 2 - Works",
			&types.AttributeValueMemberN{Value: "2"}, Financial_Trades_AveragePriceTrade),
		Entry("Value is numeric, 3 - Works",
			&types.AttributeValueMemberN{Value: "3"}, Financial_Trades_AutomaticExecution),
		Entry("Value is numeric, 4 - Works",
			&types.AttributeValueMemberN{Value: "4"}, Financial_Trades_BunchedTrade),
		Entry("Value is numeric, 5 - Works",
			&types.AttributeValueMemberN{Value: "5"}, Financial_Trades_BunchedSoldTrade),
		Entry("Value is numeric, 6 - Works",
			&types.AttributeValueMemberN{Value: "6"}, Financial_Trades_CAPElection),
		Entry("Value is numeric, 7 - Works",
			&types.AttributeValueMemberN{Value: "7"}, Financial_Trades_CashSale),
		Entry("Value is numeric, 8 - Works",
			&types.AttributeValueMemberN{Value: "8"}, Financial_Trades_ClosingPrints),
		Entry("Value is numeric, 9 - Works",
			&types.AttributeValueMemberN{Value: "9"}, Financial_Trades_CrossTrade),
		Entry("Value is numeric, 10 - Works",
			&types.AttributeValueMemberN{Value: "10"}, Financial_Trades_DerivativelyPriced),
		Entry("Value is numeric, 11 - Works",
			&types.AttributeValueMemberN{Value: "11"}, Financial_Trades_Distribution),
		Entry("Value is numeric, 12 - Works",
			&types.AttributeValueMemberN{Value: "12"}, Financial_Trades_FormT),
		Entry("Value is numeric, 13 - Works",
			&types.AttributeValueMemberN{Value: "13"}, Financial_Trades_ExtendedTradingHours),
		Entry("Value is numeric, 14 - Works",
			&types.AttributeValueMemberN{Value: "14"}, Financial_Trades_IntermarketSweep),
		Entry("Value is numeric, 15 - Works",
			&types.AttributeValueMemberN{Value: "15"}, Financial_Trades_MarketCenterOfficialClose),
		Entry("Value is numeric, 16 - Works",
			&types.AttributeValueMemberN{Value: "16"}, Financial_Trades_MarketCenterOfficialOpen),
		Entry("Value is numeric, 17 - Works",
			&types.AttributeValueMemberN{Value: "17"}, Financial_Trades_MarketCenterOpeningTrade),
		Entry("Value is numeric, 18 - Works",
			&types.AttributeValueMemberN{Value: "18"}, Financial_Trades_MarketCenterReopeningTrade),
		Entry("Value is numeric, 19 - Works",
			&types.AttributeValueMemberN{Value: "19"}, Financial_Trades_MarketCenterClosingTrade),
		Entry("Value is numeric, 20 - Works",
			&types.AttributeValueMemberN{Value: "20"}, Financial_Trades_NextDay),
		Entry("Value is numeric, 21 - Works",
			&types.AttributeValueMemberN{Value: "21"}, Financial_Trades_PriceVariationTrade),
		Entry("Value is numeric, 22 - Works",
			&types.AttributeValueMemberN{Value: "22"}, Financial_Trades_PriorReferencePrice),
		Entry("Value is numeric, 23 - Works",
			&types.AttributeValueMemberN{Value: "23"}, Financial_Trades_Rule155Trade),
		Entry("Value is numeric, 24 - Works",
			&types.AttributeValueMemberN{Value: "24"}, Financial_Trades_Rule127NYSE),
		Entry("Value is numeric, 25 - Works",
			&types.AttributeValueMemberN{Value: "25"}, Financial_Trades_OpeningPrints),
		Entry("Value is numeric, 26 - Works",
			&types.AttributeValueMemberN{Value: "26"}, Financial_Trades_Opened),
		Entry("Value is numeric, 27 - Works",
			&types.AttributeValueMemberN{Value: "27"}, Financial_Trades_StoppedStock),
		Entry("Value is numeric, 28 - Works",
			&types.AttributeValueMemberN{Value: "28"}, Financial_Trades_ReOpeningPrints),
		Entry("Value is numeric, 29 - Works",
			&types.AttributeValueMemberN{Value: "29"}, Financial_Trades_Seller),
		Entry("Value is numeric, 30 - Works",
			&types.AttributeValueMemberN{Value: "30"}, Financial_Trades_SoldLast),
		Entry("Value is numeric, 31 - Works",
			&types.AttributeValueMemberN{Value: "31"}, Financial_Trades_SoldLastAndStoppedStock),
		Entry("Value is numeric, 32 - Works",
			&types.AttributeValueMemberN{Value: "32"}, Financial_Trades_SoldOut),
		Entry("Value is numeric, 33 - Works",
			&types.AttributeValueMemberN{Value: "33"}, Financial_Trades_SoldOutOfSequence),
		Entry("Value is numeric, 34 - Works",
			&types.AttributeValueMemberN{Value: "34"}, Financial_Trades_SplitTrade),
		Entry("Value is numeric, 35 - Works",
			&types.AttributeValueMemberN{Value: "35"}, Financial_Trades_StockOption),
		Entry("Value is numeric, 36 - Works",
			&types.AttributeValueMemberN{Value: "36"}, Financial_Trades_YellowFlagRegularTrade),
		Entry("Value is numeric, 37 - Works",
			&types.AttributeValueMemberN{Value: "37"}, Financial_Trades_OddLotTrade),
		Entry("Value is numeric, 38 - Works",
			&types.AttributeValueMemberN{Value: "38"}, Financial_Trades_CorrectedConsolidatedClose),
		Entry("Value is numeric, 39 - Works",
			&types.AttributeValueMemberN{Value: "39"}, Financial_Trades_Unknown),
		Entry("Value is numeric, 40 - Works",
			&types.AttributeValueMemberN{Value: "40"}, Financial_Trades_Held),
		Entry("Value is numeric, 41 - Works",
			&types.AttributeValueMemberN{Value: "41"}, Financial_Trades_TradeThruExempt),
		Entry("Value is numeric, 42 - Works",
			&types.AttributeValueMemberN{Value: "42"}, Financial_Trades_NonEligible),
		Entry("Value is numeric, 43 - Works",
			&types.AttributeValueMemberN{Value: "43"}, Financial_Trades_NonEligibleExtended),
		Entry("Value is numeric, 44 - Works",
			&types.AttributeValueMemberN{Value: "44"}, Financial_Trades_Cancelled),
		Entry("Value is numeric, 45 - Works",
			&types.AttributeValueMemberN{Value: "45"}, Financial_Trades_Recovery),
		Entry("Value is numeric, 46 - Works",
			&types.AttributeValueMemberN{Value: "46"}, Financial_Trades_Correction),
		Entry("Value is numeric, 47 - Works",
			&types.AttributeValueMemberN{Value: "47"}, Financial_Trades_AsOf),
		Entry("Value is numeric, 48 - Works",
			&types.AttributeValueMemberN{Value: "48"}, Financial_Trades_AsOfCorrection),
		Entry("Value is numeric, 49 - Works",
			&types.AttributeValueMemberN{Value: "49"}, Financial_Trades_AsOfCancel),
		Entry("Value is numeric, 50 - Works",
			&types.AttributeValueMemberN{Value: "50"}, Financial_Trades_OOB),
		Entry("Value is numeric, 51 - Works",
			&types.AttributeValueMemberN{Value: "51"}, Financial_Trades_Summary),
		Entry("Value is numeric, 52 - Works",
			&types.AttributeValueMemberN{Value: "52"}, Financial_Trades_ContingentTrade),
		Entry("Value is numeric, 53 - Works",
			&types.AttributeValueMemberN{Value: "53"}, Financial_Trades_QualifiedContingentTrade),
		Entry("Value is numeric, 54 - Works",
			&types.AttributeValueMemberN{Value: "54"}, Financial_Trades_Errored),
		Entry("Value is numeric, 55 - Works",
			&types.AttributeValueMemberN{Value: "55"}, Financial_Trades_OpeningReopeningTradeDetail),
		Entry("Value is numeric, 56 - Works",
			&types.AttributeValueMemberN{Value: "56"}, Financial_Trades_Placeholder),
		Entry("Value is numeric, 57 - Works",
			&types.AttributeValueMemberN{Value: "57"}, Financial_Trades_ShortSaleRestrictionActivated),
		Entry("Value is numeric, 58 - Works",
			&types.AttributeValueMemberN{Value: "58"}, Financial_Trades_ShortSaleRestrictionContinued),
		Entry("Value is numeric, 59 - Works",
			&types.AttributeValueMemberN{Value: "59"}, Financial_Trades_ShortSaleRestrictionDeactivated),
		Entry("Value is numeric, 60 - Works",
			&types.AttributeValueMemberN{Value: "60"}, Financial_Trades_ShortSaleRestrictionInEffect),
		Entry("Value is numeric, 62 - Works",
			&types.AttributeValueMemberN{Value: "62"}, Financial_Trades_FinancialStatusBankrupt),
		Entry("Value is numeric, 63 - Works",
			&types.AttributeValueMemberN{Value: "63"}, Financial_Trades_FinancialStatusDeficient),
		Entry("Value is numeric, 64 - Works",
			&types.AttributeValueMemberN{Value: "64"}, Financial_Trades_FinancialStatusDelinquent),
		Entry("Value is numeric, 65 - Works",
			&types.AttributeValueMemberN{Value: "65"}, Financial_Trades_FinancialStatusBankruptAndDeficient),
		Entry("Value is numeric, 66 - Works",
			&types.AttributeValueMemberN{Value: "66"}, Financial_Trades_FinancialStatusBankruptAndDelinquent),
		Entry("Value is numeric, 67 - Works",
			&types.AttributeValueMemberN{Value: "67"}, Financial_Trades_FinancialStatusDeficientAndDelinquent),
		Entry("Value is numeric, 68 - Works",
			&types.AttributeValueMemberN{Value: "68"}, Financial_Trades_FinancialStatusDeficientDelinquentBankrupt),
		Entry("Value is numeric, 69 - Works",
			&types.AttributeValueMemberN{Value: "69"}, Financial_Trades_FinancialStatusLiquidation),
		Entry("Value is numeric, 70 - Works",
			&types.AttributeValueMemberN{Value: "70"}, Financial_Trades_FinancialStatusCreationsSuspended),
		Entry("Value is numeric, 71 - Works",
			&types.AttributeValueMemberN{Value: "71"}, Financial_Trades_FinancialStatusRedemptionsSuspended),
		Entry("Value is numeric, 201 - Works",
			&types.AttributeValueMemberN{Value: "201"}, Financial_Trades_Canceled),
		Entry("Value is numeric, 202 - Works",
			&types.AttributeValueMemberN{Value: "202"}, Financial_Trades_LateAndOutOfSequence),
		Entry("Value is numeric, 203 - Works",
			&types.AttributeValueMemberN{Value: "203"}, Financial_Trades_LastAndCanceled),
		Entry("Value is numeric, 204 - Works",
			&types.AttributeValueMemberN{Value: "204"}, Financial_Trades_Late),
		Entry("Value is numeric, 205 - Works",
			&types.AttributeValueMemberN{Value: "205"}, Financial_Trades_OpeningTradeAndCanceled),
		Entry("Value is numeric, 206 - Works",
			&types.AttributeValueMemberN{Value: "206"}, Financial_Trades_OpeningTradeLateAndOutOfSequence),
		Entry("Value is numeric, 207 - Works",
			&types.AttributeValueMemberN{Value: "207"}, Financial_Trades_OnlyTradeAndCanceled),
		Entry("Value is numeric, 208 - Works",
			&types.AttributeValueMemberN{Value: "208"}, Financial_Trades_OpeningTradeAndLate),
		Entry("Value is numeric, 209 - Works",
			&types.AttributeValueMemberN{Value: "209"}, Financial_Trades_AutomaticExecutionOption),
		Entry("Value is numeric, 210 - Works",
			&types.AttributeValueMemberN{Value: "210"}, Financial_Trades_ReopeningTrade),
		Entry("Value is numeric, 219 - Works",
			&types.AttributeValueMemberN{Value: "219"}, Financial_Trades_IntermarketSweepOrder),
		Entry("Value is numeric, 227 - Works",
			&types.AttributeValueMemberN{Value: "227"}, Financial_Trades_SingleLegAuctionNonISO),
		Entry("Value is numeric, 228 - Works",
			&types.AttributeValueMemberN{Value: "228"}, Financial_Trades_SingleLegAuctionISO),
		Entry("Value is numeric, 229 - Works",
			&types.AttributeValueMemberN{Value: "229"}, Financial_Trades_SingleLegCrossNonISO),
		Entry("Value is numeric, 230 - Works",
			&types.AttributeValueMemberN{Value: "230"}, Financial_Trades_SingleLegCrossISO),
		Entry("Value is numeric, 231 - Works",
			&types.AttributeValueMemberN{Value: "231"}, Financial_Trades_SingleLegFloorTrade),
		Entry("Value is numeric, 232 - Works",
			&types.AttributeValueMemberN{Value: "232"}, Financial_Trades_MultiLegAutoElectronicTrade),
		Entry("Value is numeric, 233 - Works",
			&types.AttributeValueMemberN{Value: "233"}, Financial_Trades_MultiLegAuction),
		Entry("Value is numeric, 234 - Works",
			&types.AttributeValueMemberN{Value: "234"}, Financial_Trades_MultiLegCross),
		Entry("Value is numeric, 235 - Works",
			&types.AttributeValueMemberN{Value: "235"}, Financial_Trades_MultiLegFloorTrade),
		Entry("Value is numeric, 236 - Works",
			&types.AttributeValueMemberN{Value: "236"}, Financial_Trades_MultiLegAutoElectronicTradeAgainstSingleLeg),
		Entry("Value is numeric, 237 - Works",
			&types.AttributeValueMemberN{Value: "237"}, Financial_Trades_StockOptionsAuction),
		Entry("Value is numeric, 238 - Works",
			&types.AttributeValueMemberN{Value: "238"}, Financial_Trades_MultiLegAuctionAgainstSingleLeg),
		Entry("Value is numeric, 239 - Works",
			&types.AttributeValueMemberN{Value: "239"}, Financial_Trades_MultiLegFloorTradeAgainstSingleLeg),
		Entry("Value is numeric, 240 - Works",
			&types.AttributeValueMemberN{Value: "240"}, Financial_Trades_StockOptionsAutoElectronicTrade),
		Entry("Value is numeric, 241 - Works",
			&types.AttributeValueMemberN{Value: "241"}, Financial_Trades_StockOptionsCross),
		Entry("Value is numeric, 242 - Works",
			&types.AttributeValueMemberN{Value: "242"}, Financial_Trades_StockOptionsFloorTrade),
		Entry("Value is numeric, 243 - Works",
			&types.AttributeValueMemberN{Value: "243"}, Financial_Trades_StockOptionsAutoElectronicTradeAgainstSingleLeg),
		Entry("Value is numeric, 244 - Works",
			&types.AttributeValueMemberN{Value: "244"}, Financial_Trades_StockOptionsAuctionAgainstSingleLeg),
		Entry("Value is numeric, 245 - Works",
			&types.AttributeValueMemberN{Value: "245"}, Financial_Trades_StockOptionsFloorTradeAgainstSingleLeg),
		Entry("Value is numeric, 246 - Works",
			&types.AttributeValueMemberN{Value: "246"}, Financial_Trades_MultiLegFloorTradeOfProprietaryProducts),
		Entry("Value is numeric, 247 - Works",
			&types.AttributeValueMemberN{Value: "247"}, Financial_Trades_MultilateralCompressionTradeOfProprietaryProducts),
		Entry("Value is numeric, 248 - Works",
			&types.AttributeValueMemberN{Value: "248"}, Financial_Trades_ExtendedHoursTrade),
		Entry("Value is NULL - Works", new(types.AttributeValueMemberNULL), Financial_Trades_Condition(0)),
		Entry("Value is string, RegularSale - Works",
			&types.AttributeValueMemberS{Value: "RegularSale"}, Financial_Trades_RegularSale),
		Entry("Value is string, Acquisition - Works",
			&types.AttributeValueMemberS{Value: "Acquisition"}, Financial_Trades_Acquisition),
		Entry("Value is string, AveragePriceTrade - Works",
			&types.AttributeValueMemberS{Value: "AveragePriceTrade"}, Financial_Trades_AveragePriceTrade),
		Entry("Value is string, AutomaticExecution - Works",
			&types.AttributeValueMemberS{Value: "AutomaticExecution"}, Financial_Trades_AutomaticExecution),
		Entry("Value is string, BunchedTrade - Works",
			&types.AttributeValueMemberS{Value: "BunchedTrade"}, Financial_Trades_BunchedTrade),
		Entry("Value is string, BunchedSoldTrade - Works",
			&types.AttributeValueMemberS{Value: "BunchedSoldTrade"}, Financial_Trades_BunchedSoldTrade),
		Entry("Value is string, CAPElection - Works",
			&types.AttributeValueMemberS{Value: "CAPElection"}, Financial_Trades_CAPElection),
		Entry("Value is string, CashSale - Works",
			&types.AttributeValueMemberS{Value: "CashSale"}, Financial_Trades_CashSale),
		Entry("Value is string, ClosingPrints - Works",
			&types.AttributeValueMemberS{Value: "ClosingPrints"}, Financial_Trades_ClosingPrints),
		Entry("Value is string, CrossTrade - Works",
			&types.AttributeValueMemberS{Value: "CrossTrade"}, Financial_Trades_CrossTrade),
		Entry("Value is string, DerivativelyPriced - Works",
			&types.AttributeValueMemberS{Value: "DerivativelyPriced"}, Financial_Trades_DerivativelyPriced),
		Entry("Value is string, Distribution - Works",
			&types.AttributeValueMemberS{Value: "Distribution"}, Financial_Trades_Distribution),
		Entry("Value is string, FormT - Works",
			&types.AttributeValueMemberS{Value: "FormT"}, Financial_Trades_FormT),
		Entry("Value is string, ExtendedTradingHours - Works",
			&types.AttributeValueMemberS{Value: "ExtendedTradingHours"}, Financial_Trades_ExtendedTradingHours),
		Entry("Value is string, IntermarketSweep - Works",
			&types.AttributeValueMemberS{Value: "IntermarketSweep"}, Financial_Trades_IntermarketSweep),
		Entry("Value is string, MarketCenterOfficialClose - Works",
			&types.AttributeValueMemberS{Value: "MarketCenterOfficialClose"}, Financial_Trades_MarketCenterOfficialClose),
		Entry("Value is string, MarketCenterOfficialOpen - Works",
			&types.AttributeValueMemberS{Value: "MarketCenterOfficialOpen"}, Financial_Trades_MarketCenterOfficialOpen),
		Entry("Value is string, MarketCenterOpeningTrade - Works",
			&types.AttributeValueMemberS{Value: "MarketCenterOpeningTrade"}, Financial_Trades_MarketCenterOpeningTrade),
		Entry("Value is string, MarketCenterReopeningTrade - Works",
			&types.AttributeValueMemberS{Value: "MarketCenterReopeningTrade"}, Financial_Trades_MarketCenterReopeningTrade),
		Entry("Value is string, MarketCenterClosingTrade - Works",
			&types.AttributeValueMemberS{Value: "MarketCenterClosingTrade"}, Financial_Trades_MarketCenterClosingTrade),
		Entry("Value is string, NextDay - Works",
			&types.AttributeValueMemberS{Value: "NextDay"}, Financial_Trades_NextDay),
		Entry("Value is string, PriceVariationTrade - Works",
			&types.AttributeValueMemberS{Value: "PriceVariationTrade"}, Financial_Trades_PriceVariationTrade),
		Entry("Value is string, PriorReferencePrice - Works",
			&types.AttributeValueMemberS{Value: "PriorReferencePrice"}, Financial_Trades_PriorReferencePrice),
		Entry("Value is string, Rule155Trade - Works",
			&types.AttributeValueMemberS{Value: "Rule155Trade"}, Financial_Trades_Rule155Trade),
		Entry("Value is string, Rule127NYSE - Works",
			&types.AttributeValueMemberS{Value: "Rule127NYSE"}, Financial_Trades_Rule127NYSE),
		Entry("Value is string, OpeningPrints - Works",
			&types.AttributeValueMemberS{Value: "OpeningPrints"}, Financial_Trades_OpeningPrints),
		Entry("Value is string, Opened - Works",
			&types.AttributeValueMemberS{Value: "Opened"}, Financial_Trades_Opened),
		Entry("Value is string, StoppedStock - Works",
			&types.AttributeValueMemberS{Value: "StoppedStock"}, Financial_Trades_StoppedStock),
		Entry("Value is string, ReOpeningPrints - Works",
			&types.AttributeValueMemberS{Value: "ReOpeningPrints"}, Financial_Trades_ReOpeningPrints),
		Entry("Value is string, Seller - Works",
			&types.AttributeValueMemberS{Value: "Seller"}, Financial_Trades_Seller),
		Entry("Value is string, SoldLast - Works",
			&types.AttributeValueMemberS{Value: "SoldLast"}, Financial_Trades_SoldLast),
		Entry("Value is string, SoldLastAndStoppedStock - Works",
			&types.AttributeValueMemberS{Value: "SoldLastAndStoppedStock"}, Financial_Trades_SoldLastAndStoppedStock),
		Entry("Value is string, SoldOut - Works",
			&types.AttributeValueMemberS{Value: "SoldOut"}, Financial_Trades_SoldOut),
		Entry("Value is string, SoldOutOfSequence - Works",
			&types.AttributeValueMemberS{Value: "SoldOutOfSequence"}, Financial_Trades_SoldOutOfSequence),
		Entry("Value is string, SplitTrade - Works",
			&types.AttributeValueMemberS{Value: "SplitTrade"}, Financial_Trades_SplitTrade),
		Entry("Value is string, StockOption - Works",
			&types.AttributeValueMemberS{Value: "StockOption"}, Financial_Trades_StockOption),
		Entry("Value is string, YellowFlagRegularTrade - Works",
			&types.AttributeValueMemberS{Value: "YellowFlagRegularTrade"}, Financial_Trades_YellowFlagRegularTrade),
		Entry("Value is string, OddLotTrade - Works",
			&types.AttributeValueMemberS{Value: "OddLotTrade"}, Financial_Trades_OddLotTrade),
		Entry("Value is string, CorrectedConsolidatedClose - Works",
			&types.AttributeValueMemberS{Value: "CorrectedConsolidatedClose"}, Financial_Trades_CorrectedConsolidatedClose),
		Entry("Value is string, Unknown - Works",
			&types.AttributeValueMemberS{Value: "Unknown"}, Financial_Trades_Unknown),
		Entry("Value is string, Held - Works",
			&types.AttributeValueMemberS{Value: "Held"}, Financial_Trades_Held),
		Entry("Value is string, TradeThruExempt - Works",
			&types.AttributeValueMemberS{Value: "TradeThruExempt"}, Financial_Trades_TradeThruExempt),
		Entry("Value is string, NonEligible - Works",
			&types.AttributeValueMemberS{Value: "NonEligible"}, Financial_Trades_NonEligible),
		Entry("Value is string, NonEligibleExtended - Works",
			&types.AttributeValueMemberS{Value: "NonEligibleExtended"}, Financial_Trades_NonEligibleExtended),
		Entry("Value is string, Cancelled - Works",
			&types.AttributeValueMemberS{Value: "Cancelled"}, Financial_Trades_Cancelled),
		Entry("Value is string, Recovery - Works",
			&types.AttributeValueMemberS{Value: "Recovery"}, Financial_Trades_Recovery),
		Entry("Value is string, Correction - Works",
			&types.AttributeValueMemberS{Value: "Correction"}, Financial_Trades_Correction),
		Entry("Value is string, AsOf - Works",
			&types.AttributeValueMemberS{Value: "AsOf"}, Financial_Trades_AsOf),
		Entry("Value is string, AsOfCorrection - Works",
			&types.AttributeValueMemberS{Value: "AsOfCorrection"}, Financial_Trades_AsOfCorrection),
		Entry("Value is string, AsOfCancel - Works",
			&types.AttributeValueMemberS{Value: "AsOfCancel"}, Financial_Trades_AsOfCancel),
		Entry("Value is string, OOB - Works",
			&types.AttributeValueMemberS{Value: "OOB"}, Financial_Trades_OOB),
		Entry("Value is string, Summary - Works",
			&types.AttributeValueMemberS{Value: "Summary"}, Financial_Trades_Summary),
		Entry("Value is string, ContingentTrade - Works",
			&types.AttributeValueMemberS{Value: "ContingentTrade"}, Financial_Trades_ContingentTrade),
		Entry("Value is string, QualifiedContingentTrade - Works",
			&types.AttributeValueMemberS{Value: "QualifiedContingentTrade"}, Financial_Trades_QualifiedContingentTrade),
		Entry("Value is string, Errored - Works",
			&types.AttributeValueMemberS{Value: "Errored"}, Financial_Trades_Errored),
		Entry("Value is string, OpeningReopeningTradeDetail - Works",
			&types.AttributeValueMemberS{Value: "OpeningReopeningTradeDetail"}, Financial_Trades_OpeningReopeningTradeDetail),
		Entry("Value is string, Placeholder - Works",
			&types.AttributeValueMemberS{Value: "Placeholder"}, Financial_Trades_Placeholder),
		Entry("Value is string, ShortSaleRestrictionActivated - Works",
			&types.AttributeValueMemberS{Value: "ShortSaleRestrictionActivated"}, Financial_Trades_ShortSaleRestrictionActivated),
		Entry("Value is string, ShortSaleRestrictionContinued - Works",
			&types.AttributeValueMemberS{Value: "ShortSaleRestrictionContinued"}, Financial_Trades_ShortSaleRestrictionContinued),
		Entry("Value is string, ShortSaleRestrictionDeactivated - Works",
			&types.AttributeValueMemberS{Value: "ShortSaleRestrictionDeactivated"}, Financial_Trades_ShortSaleRestrictionDeactivated),
		Entry("Value is string, ShortSaleRestrictionInEffect - Works",
			&types.AttributeValueMemberS{Value: "ShortSaleRestrictionInEffect"}, Financial_Trades_ShortSaleRestrictionInEffect),
		Entry("Value is string, FinancialStatusBankrupt - Works",
			&types.AttributeValueMemberS{Value: "FinancialStatusBankrupt"}, Financial_Trades_FinancialStatusBankrupt),
		Entry("Value is string, FinancialStatusDeficient - Works",
			&types.AttributeValueMemberS{Value: "FinancialStatusDeficient"}, Financial_Trades_FinancialStatusDeficient),
		Entry("Value is string, FinancialStatusDelinquent - Works",
			&types.AttributeValueMemberS{Value: "FinancialStatusDelinquent"}, Financial_Trades_FinancialStatusDelinquent),
		Entry("Value is string, FinancialStatusBankruptAndDeficient - Works",
			&types.AttributeValueMemberS{Value: "FinancialStatusBankruptAndDeficient"}, Financial_Trades_FinancialStatusBankruptAndDeficient),
		Entry("Value is string, FinancialStatusBankruptAndDelinquent - Works",
			&types.AttributeValueMemberS{Value: "FinancialStatusBankruptAndDelinquent"}, Financial_Trades_FinancialStatusBankruptAndDelinquent),
		Entry("Value is string, FinancialStatusDeficientAndDelinquent - Works",
			&types.AttributeValueMemberS{Value: "FinancialStatusDeficientAndDelinquent"}, Financial_Trades_FinancialStatusDeficientAndDelinquent),
		Entry("Value is string, FinancialStatusDeficientDelinquentBankrupt - Works",
			&types.AttributeValueMemberS{Value: "FinancialStatusDeficientDelinquentBankrupt"}, Financial_Trades_FinancialStatusDeficientDelinquentBankrupt),
		Entry("Value is string, FinancialStatusLiquidation - Works",
			&types.AttributeValueMemberS{Value: "FinancialStatusLiquidation"}, Financial_Trades_FinancialStatusLiquidation),
		Entry("Value is string, FinancialStatusCreationsSuspended - Works",
			&types.AttributeValueMemberS{Value: "FinancialStatusCreationsSuspended"}, Financial_Trades_FinancialStatusCreationsSuspended),
		Entry("Value is string, FinancialStatusRedemptionsSuspended - Works",
			&types.AttributeValueMemberS{Value: "FinancialStatusRedemptionsSuspended"}, Financial_Trades_FinancialStatusRedemptionsSuspended),
		Entry("Value is string, Canceled - Works",
			&types.AttributeValueMemberS{Value: "Canceled"}, Financial_Trades_Canceled),
		Entry("Value is string, LateAndOutOfSequence - Works",
			&types.AttributeValueMemberS{Value: "LateAndOutOfSequence"}, Financial_Trades_LateAndOutOfSequence),
		Entry("Value is string, LastAndCanceled - Works",
			&types.AttributeValueMemberS{Value: "LastAndCanceled"}, Financial_Trades_LastAndCanceled),
		Entry("Value is string, Late - Works",
			&types.AttributeValueMemberS{Value: "Late"}, Financial_Trades_Late),
		Entry("Value is string, OpeningTradeAndCanceled - Works",
			&types.AttributeValueMemberS{Value: "OpeningTradeAndCanceled"}, Financial_Trades_OpeningTradeAndCanceled),
		Entry("Value is string, OpeningTradeLateAndOutOfSequence - Works",
			&types.AttributeValueMemberS{Value: "OpeningTradeLateAndOutOfSequence"}, Financial_Trades_OpeningTradeLateAndOutOfSequence),
		Entry("Value is string, OnlyTradeAndCanceled - Works",
			&types.AttributeValueMemberS{Value: "OnlyTradeAndCanceled"}, Financial_Trades_OnlyTradeAndCanceled),
		Entry("Value is string, OpeningTradeAndLate - Works",
			&types.AttributeValueMemberS{Value: "OpeningTradeAndLate"}, Financial_Trades_OpeningTradeAndLate),
		Entry("Value is string, AutomaticExecutionOption - Works",
			&types.AttributeValueMemberS{Value: "AutomaticExecutionOption"}, Financial_Trades_AutomaticExecutionOption),
		Entry("Value is string, ReopeningTrade - Works",
			&types.AttributeValueMemberS{Value: "ReopeningTrade"}, Financial_Trades_ReopeningTrade),
		Entry("Value is string, IntermarketSweepOrder - Works",
			&types.AttributeValueMemberS{Value: "IntermarketSweepOrder"}, Financial_Trades_IntermarketSweepOrder),
		Entry("Value is string, SingleLegAuctionNonISO - Works",
			&types.AttributeValueMemberS{Value: "SingleLegAuctionNonISO"}, Financial_Trades_SingleLegAuctionNonISO),
		Entry("Value is string, SingleLegAuctionISO - Works",
			&types.AttributeValueMemberS{Value: "SingleLegAuctionISO"}, Financial_Trades_SingleLegAuctionISO),
		Entry("Value is string, SingleLegCrossNonISO - Works",
			&types.AttributeValueMemberS{Value: "SingleLegCrossNonISO"}, Financial_Trades_SingleLegCrossNonISO),
		Entry("Value is string, SingleLegCrossISO - Works",
			&types.AttributeValueMemberS{Value: "SingleLegCrossISO"}, Financial_Trades_SingleLegCrossISO),
		Entry("Value is string, SingleLegFloorTrade - Works",
			&types.AttributeValueMemberS{Value: "SingleLegFloorTrade"}, Financial_Trades_SingleLegFloorTrade),
		Entry("Value is string, MultiLegAutoElectronicTrade - Works",
			&types.AttributeValueMemberS{Value: "MultiLegAutoElectronicTrade"}, Financial_Trades_MultiLegAutoElectronicTrade),
		Entry("Value is string, MultiLegAuction - Works",
			&types.AttributeValueMemberS{Value: "MultiLegAuction"}, Financial_Trades_MultiLegAuction),
		Entry("Value is string, MultiLegCross - Works",
			&types.AttributeValueMemberS{Value: "MultiLegCross"}, Financial_Trades_MultiLegCross),
		Entry("Value is string, MultiLegFloorTrade - Works",
			&types.AttributeValueMemberS{Value: "MultiLegFloorTrade"}, Financial_Trades_MultiLegFloorTrade),
		Entry("Value is string, MultiLegAutoElectronicTradeAgainstSingleLeg - Works",
			&types.AttributeValueMemberS{Value: "MultiLegAutoElectronicTradeAgainstSingleLeg"}, Financial_Trades_MultiLegAutoElectronicTradeAgainstSingleLeg),
		Entry("Value is string, StockOptionsAuction - Works",
			&types.AttributeValueMemberS{Value: "StockOptionsAuction"}, Financial_Trades_StockOptionsAuction),
		Entry("Value is string, MultiLegAuctionAgainstSingleLeg - Works",
			&types.AttributeValueMemberS{Value: "MultiLegAuctionAgainstSingleLeg"}, Financial_Trades_MultiLegAuctionAgainstSingleLeg),
		Entry("Value is string, MultiLegFloorTradeAgainstSingleLeg - Works",
			&types.AttributeValueMemberS{Value: "MultiLegFloorTradeAgainstSingleLeg"}, Financial_Trades_MultiLegFloorTradeAgainstSingleLeg),
		Entry("Value is string, StockOptionsAutoElectronicTrade - Works",
			&types.AttributeValueMemberS{Value: "StockOptionsAutoElectronicTrade"}, Financial_Trades_StockOptionsAutoElectronicTrade),
		Entry("Value is string, StockOptionsCross - Works",
			&types.AttributeValueMemberS{Value: "StockOptionsCross"}, Financial_Trades_StockOptionsCross),
		Entry("Value is string, StockOptionsFloorTrade - Works",
			&types.AttributeValueMemberS{Value: "StockOptionsFloorTrade"}, Financial_Trades_StockOptionsFloorTrade),
		Entry("Value is string, StockOptionsAutoElectronicTradeAgainstSingleLeg - Works",
			&types.AttributeValueMemberS{Value: "StockOptionsAutoElectronicTradeAgainstSingleLeg"}, Financial_Trades_StockOptionsAutoElectronicTradeAgainstSingleLeg),
		Entry("Value is string, StockOptionsAuctionAgainstSingleLeg - Works",
			&types.AttributeValueMemberS{Value: "StockOptionsAuctionAgainstSingleLeg"}, Financial_Trades_StockOptionsAuctionAgainstSingleLeg),
		Entry("Value is string, StockOptionsFloorTradeAgainstSingleLeg - Works",
			&types.AttributeValueMemberS{Value: "StockOptionsFloorTradeAgainstSingleLeg"}, Financial_Trades_StockOptionsFloorTradeAgainstSingleLeg),
		Entry("Value is string, MultiLegFloorTradeOfProprietaryProducts - Works",
			&types.AttributeValueMemberS{Value: "MultiLegFloorTradeOfProprietaryProducts"}, Financial_Trades_MultiLegFloorTradeOfProprietaryProducts),
		Entry("Value is string, MultilateralCompressionTradeOfProprietaryProducts - Works",
			&types.AttributeValueMemberS{Value: "MultilateralCompressionTradeOfProprietaryProducts"}, Financial_Trades_MultilateralCompressionTradeOfProprietaryProducts),
		Entry("Value is string, ExtendedHoursTrade - Works",
			&types.AttributeValueMemberS{Value: "ExtendedHoursTrade"}, Financial_Trades_ExtendedHoursTrade))

	// Test that attempting to deserialize a Financial.Trades.Condition will fial and return an
	// error if the value cannot be converted to either the name value or integer value of the enum option
	It("Scan - Value is nil - Error", func() {

		// Attempt to convert a fake string value into a Financial.Trades.Condition
		// This should return an error
		var enum *Financial_Trades_Condition
		err := enum.Scan(nil)

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of %!q(<nil>) had an invalid type of <nil>"))
		Expect(enum).Should(BeNil())
	})

	// Test the conditions under which values should be convertible to a Financial.Trades.Condition
	DescribeTable("Scan Tests",
		func(value interface{}, shouldBe Financial_Trades_Condition) {

			// Attempt to convert the value into a Financial.Trades.Condition
			// This should not fail
			var enum Financial_Trades_Condition
			err := enum.Scan(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("RegularSale - Works", "RegularSale", Financial_Trades_RegularSale),
		Entry("Acquisition - Works", "Acquisition", Financial_Trades_Acquisition),
		Entry("AveragePriceTrade - Works", "AveragePriceTrade", Financial_Trades_AveragePriceTrade),
		Entry("AutomaticExecution - Works", "AutomaticExecution", Financial_Trades_AutomaticExecution),
		Entry("BunchedTrade - Works", "BunchedTrade", Financial_Trades_BunchedTrade),
		Entry("BunchedSoldTrade - Works", "BunchedSoldTrade", Financial_Trades_BunchedSoldTrade),
		Entry("CAPElection - Works", "CAPElection", Financial_Trades_CAPElection),
		Entry("CashSale - Works", "CashSale", Financial_Trades_CashSale),
		Entry("ClosingPrints - Works", "ClosingPrints", Financial_Trades_ClosingPrints),
		Entry("CrossTrade - Works", "CrossTrade", Financial_Trades_CrossTrade),
		Entry("DerivativelyPriced - Works", "DerivativelyPriced", Financial_Trades_DerivativelyPriced),
		Entry("Distribution - Works", "Distribution", Financial_Trades_Distribution),
		Entry("FormT - Works", "FormT", Financial_Trades_FormT),
		Entry("ExtendedTradingHours - Works", "ExtendedTradingHours", Financial_Trades_ExtendedTradingHours),
		Entry("IntermarketSweep - Works", "IntermarketSweep", Financial_Trades_IntermarketSweep),
		Entry("MarketCenterOfficialClose - Works", "MarketCenterOfficialClose", Financial_Trades_MarketCenterOfficialClose),
		Entry("MarketCenterOfficialOpen - Works", "MarketCenterOfficialOpen", Financial_Trades_MarketCenterOfficialOpen),
		Entry("MarketCenterOpeningTrade - Works", "MarketCenterOpeningTrade", Financial_Trades_MarketCenterOpeningTrade),
		Entry("MarketCenterReopeningTrade - Works", "MarketCenterReopeningTrade", Financial_Trades_MarketCenterReopeningTrade),
		Entry("MarketCenterClosingTrade - Works", "MarketCenterClosingTrade", Financial_Trades_MarketCenterClosingTrade),
		Entry("NextDay - Works", "NextDay", Financial_Trades_NextDay),
		Entry("PriceVariationTrade - Works", "PriceVariationTrade", Financial_Trades_PriceVariationTrade),
		Entry("PriorReferencePrice - Works", "PriorReferencePrice", Financial_Trades_PriorReferencePrice),
		Entry("Rule155Trade - Works", "Rule155Trade", Financial_Trades_Rule155Trade),
		Entry("Rule127NYSE - Works", "Rule127NYSE", Financial_Trades_Rule127NYSE),
		Entry("OpeningPrints - Works", "OpeningPrints", Financial_Trades_OpeningPrints),
		Entry("Opened - Works", "Opened", Financial_Trades_Opened),
		Entry("StoppedStock - Works", "StoppedStock", Financial_Trades_StoppedStock),
		Entry("ReOpeningPrints - Works", "ReOpeningPrints", Financial_Trades_ReOpeningPrints),
		Entry("Seller - Works", "Seller", Financial_Trades_Seller),
		Entry("SoldLast - Works", "SoldLast", Financial_Trades_SoldLast),
		Entry("SoldLastAndStoppedStock - Works", "SoldLastAndStoppedStock", Financial_Trades_SoldLastAndStoppedStock),
		Entry("SoldOut - Works", "SoldOut", Financial_Trades_SoldOut),
		Entry("SoldOutOfSequence - Works", "SoldOutOfSequence", Financial_Trades_SoldOutOfSequence),
		Entry("SplitTrade - Works", "SplitTrade", Financial_Trades_SplitTrade),
		Entry("StockOption - Works", "StockOption", Financial_Trades_StockOption),
		Entry("YellowFlagRegularTrade - Works", "YellowFlagRegularTrade", Financial_Trades_YellowFlagRegularTrade),
		Entry("OddLotTrade - Works", "OddLotTrade", Financial_Trades_OddLotTrade),
		Entry("CorrectedConsolidatedClose - Works", "CorrectedConsolidatedClose", Financial_Trades_CorrectedConsolidatedClose),
		Entry("Unknown - Works", "Unknown", Financial_Trades_Unknown),
		Entry("Held - Works", "Held", Financial_Trades_Held),
		Entry("TradeThruExempt - Works", "TradeThruExempt", Financial_Trades_TradeThruExempt),
		Entry("NonEligible - Works", "NonEligible", Financial_Trades_NonEligible),
		Entry("NonEligibleExtended - Works", "NonEligibleExtended", Financial_Trades_NonEligibleExtended),
		Entry("Cancelled - Works", "Cancelled", Financial_Trades_Cancelled),
		Entry("Recovery - Works", "Recovery", Financial_Trades_Recovery),
		Entry("Correction - Works", "Correction", Financial_Trades_Correction),
		Entry("AsOf - Works", "AsOf", Financial_Trades_AsOf),
		Entry("AsOfCorrection - Works", "AsOfCorrection", Financial_Trades_AsOfCorrection),
		Entry("AsOfCancel - Works", "AsOfCancel", Financial_Trades_AsOfCancel),
		Entry("OOB - Works", "OOB", Financial_Trades_OOB),
		Entry("Summary - Works", "Summary", Financial_Trades_Summary),
		Entry("ContingentTrade - Works", "ContingentTrade", Financial_Trades_ContingentTrade),
		Entry("QualifiedContingentTrade - Works", "QualifiedContingentTrade", Financial_Trades_QualifiedContingentTrade),
		Entry("Errored - Works", "Errored", Financial_Trades_Errored),
		Entry("OpeningReopeningTradeDetail - Works", "OpeningReopeningTradeDetail", Financial_Trades_OpeningReopeningTradeDetail),
		Entry("Placeholder - Works", "Placeholder", Financial_Trades_Placeholder),
		Entry("ShortSaleRestrictionActivated - Works",
			"ShortSaleRestrictionActivated", Financial_Trades_ShortSaleRestrictionActivated),
		Entry("ShortSaleRestrictionContinued - Works",
			"ShortSaleRestrictionContinued", Financial_Trades_ShortSaleRestrictionContinued),
		Entry("ShortSaleRestrictionDeactivated - Works",
			"ShortSaleRestrictionDeactivated", Financial_Trades_ShortSaleRestrictionDeactivated),
		Entry("ShortSaleRestrictionInEffect - Works",
			"ShortSaleRestrictionInEffect", Financial_Trades_ShortSaleRestrictionInEffect),
		Entry("FinancialStatusBankrupt - Works", "FinancialStatusBankrupt", Financial_Trades_FinancialStatusBankrupt),
		Entry("FinancialStatusDeficient - Works", "FinancialStatusDeficient", Financial_Trades_FinancialStatusDeficient),
		Entry("FinancialStatusDelinquent - Works", "FinancialStatusDelinquent", Financial_Trades_FinancialStatusDelinquent),
		Entry("FinancialStatusBankruptAndDeficient - Works",
			"FinancialStatusBankruptAndDeficient", Financial_Trades_FinancialStatusBankruptAndDeficient),
		Entry("FinancialStatusBankruptAndDelinquent - Works",
			"FinancialStatusBankruptAndDelinquent", Financial_Trades_FinancialStatusBankruptAndDelinquent),
		Entry("FinancialStatusDeficientAndDelinquent - Works",
			"FinancialStatusDeficientAndDelinquent", Financial_Trades_FinancialStatusDeficientAndDelinquent),
		Entry("FinancialStatusDeficientDelinquentBankrupt - Works",
			"FinancialStatusDeficientDelinquentBankrupt", Financial_Trades_FinancialStatusDeficientDelinquentBankrupt),
		Entry("FinancialStatusLiquidation - Works", "FinancialStatusLiquidation", Financial_Trades_FinancialStatusLiquidation),
		Entry("FinancialStatusCreationsSuspended - Works",
			"FinancialStatusCreationsSuspended", Financial_Trades_FinancialStatusCreationsSuspended),
		Entry("FinancialStatusRedemptionsSuspended - Works",
			"FinancialStatusRedemptionsSuspended", Financial_Trades_FinancialStatusRedemptionsSuspended),
		Entry("Canceled - Works", "Canceled", Financial_Trades_Canceled),
		Entry("LateAndOutOfSequence - Works", "LateAndOutOfSequence", Financial_Trades_LateAndOutOfSequence),
		Entry("LastAndCanceled - Works", "LastAndCanceled", Financial_Trades_LastAndCanceled),
		Entry("Late - Works", "Late", Financial_Trades_Late),
		Entry("OpeningTradeAndCanceled - Works", "OpeningTradeAndCanceled", Financial_Trades_OpeningTradeAndCanceled),
		Entry("OpeningTradeLateAndOutOfSequence - Works",
			"OpeningTradeLateAndOutOfSequence", Financial_Trades_OpeningTradeLateAndOutOfSequence),
		Entry("OnlyTradeAndCanceled - Works", "OnlyTradeAndCanceled", Financial_Trades_OnlyTradeAndCanceled),
		Entry("OpeningTradeAndLate - Works", "OpeningTradeAndLate", Financial_Trades_OpeningTradeAndLate),
		Entry("AutomaticExecutionOption - Works", "AutomaticExecutionOption", Financial_Trades_AutomaticExecutionOption),
		Entry("ReopeningTrade - Works", "ReopeningTrade", Financial_Trades_ReopeningTrade),
		Entry("IntermarketSweepOrder - Works", "IntermarketSweepOrder", Financial_Trades_IntermarketSweepOrder),
		Entry("SingleLegAuctionNonISO - Works", "SingleLegAuctionNonISO", Financial_Trades_SingleLegAuctionNonISO),
		Entry("SingleLegAuctionISO - Works", "SingleLegAuctionISO", Financial_Trades_SingleLegAuctionISO),
		Entry("SingleLegCrossNonISO - Works", "SingleLegCrossNonISO", Financial_Trades_SingleLegCrossNonISO),
		Entry("SingleLegCrossISO - Works", "SingleLegCrossISO", Financial_Trades_SingleLegCrossISO),
		Entry("SingleLegFloorTrade - Works", "SingleLegFloorTrade", Financial_Trades_SingleLegFloorTrade),
		Entry("MultiLegAutoElectronicTrade - Works",
			"MultiLegAutoElectronicTrade", Financial_Trades_MultiLegAutoElectronicTrade),
		Entry("MultiLegAuction - Works", "MultiLegAuction", Financial_Trades_MultiLegAuction),
		Entry("MultiLegCross - Works", "MultiLegCross", Financial_Trades_MultiLegCross),
		Entry("MultiLegFloorTrade - Works", "MultiLegFloorTrade", Financial_Trades_MultiLegFloorTrade),
		Entry("MultiLegAutoElectronicTradeAgainstSingleLeg - Works",
			"MultiLegAutoElectronicTradeAgainstSingleLeg", Financial_Trades_MultiLegAutoElectronicTradeAgainstSingleLeg),
		Entry("StockOptionsAuction - Works", "StockOptionsAuction", Financial_Trades_StockOptionsAuction),
		Entry("MultiLegAuctionAgainstSingleLeg - Works",
			"MultiLegAuctionAgainstSingleLeg", Financial_Trades_MultiLegAuctionAgainstSingleLeg),
		Entry("MultiLegFloorTradeAgainstSingleLeg - Works",
			"MultiLegFloorTradeAgainstSingleLeg", Financial_Trades_MultiLegFloorTradeAgainstSingleLeg),
		Entry("StockOptionsAutoElectronicTrade - Works",
			"StockOptionsAutoElectronicTrade", Financial_Trades_StockOptionsAutoElectronicTrade),
		Entry("StockOptionsCross - Works", "StockOptionsCross", Financial_Trades_StockOptionsCross),
		Entry("StockOptionsFloorTrade - Works", "StockOptionsFloorTrade", Financial_Trades_StockOptionsFloorTrade),
		Entry("StockOptionsAutoElectronicTradeAgainstSingleLeg - Works",
			"StockOptionsAutoElectronicTradeAgainstSingleLeg", Financial_Trades_StockOptionsAutoElectronicTradeAgainstSingleLeg),
		Entry("StockOptionsAuctionAgainstSingleLeg - Works",
			"StockOptionsAuctionAgainstSingleLeg", Financial_Trades_StockOptionsAuctionAgainstSingleLeg),
		Entry("StockOptionsFloorTradeAgainstSingleLeg - Works",
			"StockOptionsFloorTradeAgainstSingleLeg", Financial_Trades_StockOptionsFloorTradeAgainstSingleLeg),
		Entry("MultiLegFloorTradeOfProprietaryProducts - Works",
			"MultiLegFloorTradeOfProprietaryProducts", Financial_Trades_MultiLegFloorTradeOfProprietaryProducts),
		Entry("MultilateralCompressionTradeOfProprietaryProducts - Works",
			"MultilateralCompressionTradeOfProprietaryProducts", Financial_Trades_MultilateralCompressionTradeOfProprietaryProducts),
		Entry("ExtendedHoursTrade - Works", "ExtendedHoursTrade", Financial_Trades_ExtendedHoursTrade),
		Entry("0 - Works", 0, Financial_Trades_RegularSale),
		Entry("1 - Works", 1, Financial_Trades_Acquisition),
		Entry("2 - Works", 2, Financial_Trades_AveragePriceTrade),
		Entry("3 - Works", 3, Financial_Trades_AutomaticExecution),
		Entry("4 - Works", 4, Financial_Trades_BunchedTrade),
		Entry("5 - Works", 5, Financial_Trades_BunchedSoldTrade),
		Entry("6 - Works", 6, Financial_Trades_CAPElection),
		Entry("7 - Works", 7, Financial_Trades_CashSale),
		Entry("8 - Works", 8, Financial_Trades_ClosingPrints),
		Entry("9 - Works", 9, Financial_Trades_CrossTrade),
		Entry("10 - Works", 10, Financial_Trades_DerivativelyPriced),
		Entry("11 - Works", 11, Financial_Trades_Distribution),
		Entry("12 - Works", 12, Financial_Trades_FormT),
		Entry("13 - Works", 13, Financial_Trades_ExtendedTradingHours),
		Entry("14 - Works", 14, Financial_Trades_IntermarketSweep),
		Entry("15 - Works", 15, Financial_Trades_MarketCenterOfficialClose),
		Entry("16 - Works", 16, Financial_Trades_MarketCenterOfficialOpen),
		Entry("17 - Works", 17, Financial_Trades_MarketCenterOpeningTrade),
		Entry("18 - Works", 18, Financial_Trades_MarketCenterReopeningTrade),
		Entry("19 - Works", 19, Financial_Trades_MarketCenterClosingTrade),
		Entry("20 - Works", 20, Financial_Trades_NextDay),
		Entry("21 - Works", 21, Financial_Trades_PriceVariationTrade),
		Entry("22 - Works", 22, Financial_Trades_PriorReferencePrice),
		Entry("23 - Works", 23, Financial_Trades_Rule155Trade),
		Entry("24 - Works", 24, Financial_Trades_Rule127NYSE),
		Entry("25 - Works", 25, Financial_Trades_OpeningPrints),
		Entry("26 - Works", 26, Financial_Trades_Opened),
		Entry("27 - Works", 27, Financial_Trades_StoppedStock),
		Entry("28 - Works", 28, Financial_Trades_ReOpeningPrints),
		Entry("29 - Works", 29, Financial_Trades_Seller),
		Entry("30 - Works", 30, Financial_Trades_SoldLast),
		Entry("31 - Works", 31, Financial_Trades_SoldLastAndStoppedStock),
		Entry("32 - Works", 32, Financial_Trades_SoldOut),
		Entry("33 - Works", 33, Financial_Trades_SoldOutOfSequence),
		Entry("34 - Works", 34, Financial_Trades_SplitTrade),
		Entry("35 - Works", 35, Financial_Trades_StockOption),
		Entry("36 - Works", 36, Financial_Trades_YellowFlagRegularTrade),
		Entry("37 - Works", 37, Financial_Trades_OddLotTrade),
		Entry("38 - Works", 38, Financial_Trades_CorrectedConsolidatedClose),
		Entry("39 - Works", 39, Financial_Trades_Unknown),
		Entry("40 - Works", 40, Financial_Trades_Held),
		Entry("41 - Works", 41, Financial_Trades_TradeThruExempt),
		Entry("42 - Works", 42, Financial_Trades_NonEligible),
		Entry("43 - Works", 43, Financial_Trades_NonEligibleExtended),
		Entry("44 - Works", 44, Financial_Trades_Cancelled),
		Entry("45 - Works", 45, Financial_Trades_Recovery),
		Entry("46 - Works", 46, Financial_Trades_Correction),
		Entry("47 - Works", 47, Financial_Trades_AsOf),
		Entry("48 - Works", 48, Financial_Trades_AsOfCorrection),
		Entry("49 - Works", 49, Financial_Trades_AsOfCancel),
		Entry("50 - Works", 50, Financial_Trades_OOB),
		Entry("51 - Works", 51, Financial_Trades_Summary),
		Entry("52 - Works", 52, Financial_Trades_ContingentTrade),
		Entry("53 - Works", 53, Financial_Trades_QualifiedContingentTrade),
		Entry("54 - Works", 54, Financial_Trades_Errored),
		Entry("55 - Works", 55, Financial_Trades_OpeningReopeningTradeDetail),
		Entry("56 - Works", 56, Financial_Trades_Placeholder),
		Entry("57 - Works", 57, Financial_Trades_ShortSaleRestrictionActivated),
		Entry("58 - Works", 58, Financial_Trades_ShortSaleRestrictionContinued),
		Entry("59 - Works", 59, Financial_Trades_ShortSaleRestrictionDeactivated),
		Entry("60 - Works", 60, Financial_Trades_ShortSaleRestrictionInEffect),
		Entry("62 - Works", 62, Financial_Trades_FinancialStatusBankrupt),
		Entry("63 - Works", 63, Financial_Trades_FinancialStatusDeficient),
		Entry("64 - Works", 64, Financial_Trades_FinancialStatusDelinquent),
		Entry("65 - Works", 65, Financial_Trades_FinancialStatusBankruptAndDeficient),
		Entry("66 - Works", 66, Financial_Trades_FinancialStatusBankruptAndDelinquent),
		Entry("67 - Works", 67, Financial_Trades_FinancialStatusDeficientAndDelinquent),
		Entry("68 - Works", 68, Financial_Trades_FinancialStatusDeficientDelinquentBankrupt),
		Entry("69 - Works", 69, Financial_Trades_FinancialStatusLiquidation),
		Entry("70 - Works", 70, Financial_Trades_FinancialStatusCreationsSuspended),
		Entry("71 - Works", 71, Financial_Trades_FinancialStatusRedemptionsSuspended),
		Entry("201 - Works", 201, Financial_Trades_Canceled),
		Entry("202 - Works", 202, Financial_Trades_LateAndOutOfSequence),
		Entry("203 - Works", 203, Financial_Trades_LastAndCanceled),
		Entry("204 - Works", 204, Financial_Trades_Late),
		Entry("205 - Works", 205, Financial_Trades_OpeningTradeAndCanceled),
		Entry("206 - Works", 206, Financial_Trades_OpeningTradeLateAndOutOfSequence),
		Entry("207 - Works", 207, Financial_Trades_OnlyTradeAndCanceled),
		Entry("208 - Works", 208, Financial_Trades_OpeningTradeAndLate),
		Entry("209 - Works", 209, Financial_Trades_AutomaticExecutionOption),
		Entry("210 - Works", 210, Financial_Trades_ReopeningTrade),
		Entry("219 - Works", 219, Financial_Trades_IntermarketSweepOrder),
		Entry("227 - Works", 227, Financial_Trades_SingleLegAuctionNonISO),
		Entry("228 - Works", 228, Financial_Trades_SingleLegAuctionISO),
		Entry("229 - Works", 229, Financial_Trades_SingleLegCrossNonISO),
		Entry("230 - Works", 230, Financial_Trades_SingleLegCrossISO),
		Entry("231 - Works", 231, Financial_Trades_SingleLegFloorTrade),
		Entry("232 - Works", 232, Financial_Trades_MultiLegAutoElectronicTrade),
		Entry("233 - Works", 233, Financial_Trades_MultiLegAuction),
		Entry("234 - Works", 234, Financial_Trades_MultiLegCross),
		Entry("235 - Works", 235, Financial_Trades_MultiLegFloorTrade),
		Entry("236 - Works", 236, Financial_Trades_MultiLegAutoElectronicTradeAgainstSingleLeg),
		Entry("237 - Works", 237, Financial_Trades_StockOptionsAuction),
		Entry("238 - Works", 238, Financial_Trades_MultiLegAuctionAgainstSingleLeg),
		Entry("239 - Works", 239, Financial_Trades_MultiLegFloorTradeAgainstSingleLeg),
		Entry("240 - Works", 240, Financial_Trades_StockOptionsAutoElectronicTrade),
		Entry("241 - Works", 241, Financial_Trades_StockOptionsCross),
		Entry("242 - Works", 242, Financial_Trades_StockOptionsFloorTrade),
		Entry("243 - Works", 243, Financial_Trades_StockOptionsAutoElectronicTradeAgainstSingleLeg),
		Entry("244 - Works", 244, Financial_Trades_StockOptionsAuctionAgainstSingleLeg),
		Entry("245 - Works", 245, Financial_Trades_StockOptionsFloorTradeAgainstSingleLeg),
		Entry("246 - Works", 246, Financial_Trades_MultiLegFloorTradeOfProprietaryProducts),
		Entry("247 - Works", 247, Financial_Trades_MultilateralCompressionTradeOfProprietaryProducts),
		Entry("248 - Works", 248, Financial_Trades_ExtendedHoursTrade))
})

var _ = Describe("Financial.Trades.CorrectionCode Marshal/Unmarshal Tests", func() {

	// Test that converting the Financial.Trades.CorrectionCode enum to JSON works for all values
	DescribeTable("MarshalJSON Tests",
		func(enum Financial_Trades_CorrectionCode, value string) {
			data, err := json.Marshal(enum)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("NotCorrected - Works", Financial_Trades_NotCorrected, "\"NotCorrected\""),
		Entry("LateCorrected - Works", Financial_Trades_LateCorrected, "\"LateCorrected\""),
		Entry("Erroneous - Works", Financial_Trades_Erroneous, "\"Erroneous\""),
		Entry("Cancel - Works", Financial_Trades_Cancel, "\"Cancel\""),
		Entry("CancelRecord - Works", Financial_Trades_CancelRecord, "\"CancelRecord\""),
		Entry("ErrorRecord - Works", Financial_Trades_ErrorRecord, "\"ErrorRecord\""),
		Entry("CorrectionRecord - Works", Financial_Trades_CorrectionRecord, "\"CorrectionRecord\""))

	// Test that converting the Financial.Trades.CorrectionCode enum to a CSV column works for all values
	DescribeTable("MarshalCSV Tests",
		func(enum Financial_Trades_CorrectionCode, value string) {
			data, err := enum.MarshalCSV()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("NotCorrected - Works", Financial_Trades_NotCorrected, "00"),
		Entry("LateCorrected - Works", Financial_Trades_LateCorrected, "01"),
		Entry("Erroneous - Works", Financial_Trades_Erroneous, "07"),
		Entry("Cancel - Works", Financial_Trades_Cancel, "08"),
		Entry("CancelRecord - Works", Financial_Trades_CancelRecord, "10"),
		Entry("ErrorRecord - Works", Financial_Trades_ErrorRecord, "11"),
		Entry("CorrectionRecord - Works", Financial_Trades_CorrectionRecord, "12"))

	// Test that converting the Financial.Trades.CorrectionCode enum to a AttributeValue works for all values
	DescribeTable("MarshalDynamoDBAttributeValue Tests",
		func(enum Financial_Trades_CorrectionCode, value string) {
			data, err := enum.MarshalDynamoDBAttributeValue()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data.(*types.AttributeValueMemberS).Value).Should(Equal(value))
		},
		Entry("NotCorrected - Works", Financial_Trades_NotCorrected, "NotCorrected"),
		Entry("LateCorrected - Works", Financial_Trades_LateCorrected, "LateCorrected"),
		Entry("Erroneous - Works", Financial_Trades_Erroneous, "Erroneous"),
		Entry("Cancel - Works", Financial_Trades_Cancel, "Cancel"),
		Entry("CancelRecord - Works", Financial_Trades_CancelRecord, "CancelRecord"),
		Entry("ErrorRecord - Works", Financial_Trades_ErrorRecord, "ErrorRecord"),
		Entry("CorrectionRecord - Works", Financial_Trades_CorrectionRecord, "CorrectionRecord"))

	// Test that attempting to deserialize a Financial.Trades.CorrectionCode will fail and
	// return an error if the value canno be deserialized from a JSON value to a string
	It("UnmarshalJSON fails - Error", func() {

		// Attempt to convert a non-parseable string value into a Financial.Trades.CorrectionCode
		// This should return an error
		enum := new(Financial_Trades_CorrectionCode)
		err := enum.UnmarshalJSON([]byte("derp"))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Trades_CorrectionCode"))
	})

	// Test that attempting to deserialize a Financial.Trades.CorrectionCode will fail and
	// return an error if the value cannot be converted to either the name value or integer
	// value of the enum option
	It("UnmarshalJSON - Value is invalid - Error", func() {

		// Attempt to convert a fake string value into a Financial.Trades.CorrectionCode
		// This should return an error
		enum := new(Financial_Trades_CorrectionCode)
		err := enum.UnmarshalJSON([]byte("\"derp\""))

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"derp\" cannot be mapped to a gopb.Financial_Trades_CorrectionCode"))
	})

	// Test the conditions under which values should be convertible to a Financial.Trades.CorrectionCode
	DescribeTable("UnmarshalJSON Tests",
		func(value interface{}, shouldBe Financial_Trades_CorrectionCode) {

			// Attempt to convert the string value into a Financial.Trades.CorrectionCode
			// This should not fail
			var enum Financial_Trades_CorrectionCode
			err := enum.UnmarshalJSON([]byte(fmt.Sprintf("%v", value)))

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("NotCorrected - Works", "\"NotCorrected\"", Financial_Trades_NotCorrected),
		Entry("LateCorrected - Works", "\"LateCorrected\"", Financial_Trades_LateCorrected),
		Entry("Erroneous - Works", "\"Erroneous\"", Financial_Trades_Erroneous),
		Entry("Cancel - Works", "\"Cancel\"", Financial_Trades_Cancel),
		Entry("CancelRecord - Works", "\"CancelRecord\"", Financial_Trades_CancelRecord),
		Entry("ErrorRecord - Works", "\"ErrorRecord\"", Financial_Trades_ErrorRecord),
		Entry("CorrectionRecord - Works", "\"CorrectionRecord\"", Financial_Trades_CorrectionRecord),
		Entry("00 - Works", "\"00\"", Financial_Trades_NotCorrected),
		Entry("01 - Works", "\"01\"", Financial_Trades_LateCorrected),
		Entry("07 - Works", "\"07\"", Financial_Trades_Erroneous),
		Entry("08 - Works", "\"08\"", Financial_Trades_Cancel),
		Entry("'0' - Works", "\"0\"", Financial_Trades_NotCorrected),
		Entry("'1' - Works", "\"1\"", Financial_Trades_LateCorrected),
		Entry("'7' - Works", "\"7\"", Financial_Trades_Erroneous),
		Entry("'8' - Works", "\"8\"", Financial_Trades_Cancel),
		Entry("'10' - Works", "\"10\"", Financial_Trades_CancelRecord),
		Entry("'11' - Works", "\"11\"", Financial_Trades_ErrorRecord),
		Entry("'12' - Works", "\"12\"", Financial_Trades_CorrectionRecord),
		Entry("0 - Works", 0, Financial_Trades_NotCorrected),
		Entry("1 - Works", 1, Financial_Trades_LateCorrected),
		Entry("7 - Works", 7, Financial_Trades_Erroneous),
		Entry("8 - Works", 8, Financial_Trades_Cancel),
		Entry("10 - Works", 10, Financial_Trades_CancelRecord),
		Entry("11 - Works", 11, Financial_Trades_ErrorRecord),
		Entry("12 - Works", 12, Financial_Trades_CorrectionRecord))

	// Test that attempting to deserialize a Financial.Trades.CorrectionCode will fial and return an
	// error if the value cannot be converted to either the name value or integer value
	// of the enum option
	It("UnmarshalCSV - Value is empty - Error", func() {

		// Attempt to convert a fake string value into a Financial.Trades.CorrectionCode
		// This should return an error
		enum := new(Financial_Trades_CorrectionCode)
		err := enum.UnmarshalCSV("")

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of \"\" cannot be mapped to a gopb.Financial_Trades_CorrectionCode"))
	})

	// Test the conditions under which values should be convertible to a Financial.Trades.CorrectionCode
	DescribeTable("UnmarshalCSV Tests",
		func(value string, shouldBe Financial_Trades_CorrectionCode) {

			// Attempt to convert the value into a Financial.Trades.CorrectionCode
			// This should not fail
			var enum Financial_Trades_CorrectionCode
			err := enum.UnmarshalCSV(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("NotCorrected - Works", "NotCorrected", Financial_Trades_NotCorrected),
		Entry("LateCorrected - Works", "LateCorrected", Financial_Trades_LateCorrected),
		Entry("Erroneous - Works", "Erroneous", Financial_Trades_Erroneous),
		Entry("Cancel - Works", "Cancel", Financial_Trades_Cancel),
		Entry("CancelRecord - Works", "CancelRecord", Financial_Trades_CancelRecord),
		Entry("ErrorRecord - Works", "ErrorRecord", Financial_Trades_ErrorRecord),
		Entry("CorrectionRecord - Works", "CorrectionRecord", Financial_Trades_CorrectionRecord),
		Entry("00 - Works", "00", Financial_Trades_NotCorrected),
		Entry("01 - Works", "01", Financial_Trades_LateCorrected),
		Entry("07 - Works", "07", Financial_Trades_Erroneous),
		Entry("08 - Works", "08", Financial_Trades_Cancel),
		Entry("0 - Works", "0", Financial_Trades_NotCorrected),
		Entry("1 - Works", "1", Financial_Trades_LateCorrected),
		Entry("7 - Works", "7", Financial_Trades_Erroneous),
		Entry("8 - Works", "8", Financial_Trades_Cancel),
		Entry("10 - Works", "10", Financial_Trades_CancelRecord),
		Entry("11 - Works", "11", Financial_Trades_ErrorRecord),
		Entry("12 - Works", "12", Financial_Trades_CorrectionRecord))

	// Tests that, if the attribute type submitted to UnmarshalDynamoDBAttributeValue is not one we
	// recognize, then the function will return an error
	It("UnmarshalDynamoDBAttributeValue - AttributeValue type invalid - Error", func() {
		value := new(Financial_Trades_CorrectionCode)
		err := attributevalue.Unmarshal(&types.AttributeValueMemberBOOL{Value: true}, &value)
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("Attribute value of *types.AttributeValueMemberBOOL could not be converted to a Financial.Trades.CorrectionCode"))
	})

	// Tests the conditions under which UnmarshalDynamoDBAttributeValue is called and no error is generated
	DescribeTable("UnmarshalDynamoDBAttributeValue - AttributeValue Conditions",
		func(raw types.AttributeValue, expected Financial_Trades_CorrectionCode) {
			var value Financial_Trades_CorrectionCode
			err := attributevalue.Unmarshal(raw, &value)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(value).Should(Equal(expected))
		},
		Entry("Value is []bytes, NotCorrected - Works",
			&types.AttributeValueMemberB{Value: []byte("NotCorrected")}, Financial_Trades_NotCorrected),
		Entry("Value is []bytes, LateCorrected - Works",
			&types.AttributeValueMemberB{Value: []byte("LateCorrected")}, Financial_Trades_LateCorrected),
		Entry("Value is []bytes, Erroneous - Works",
			&types.AttributeValueMemberB{Value: []byte("Erroneous")}, Financial_Trades_Erroneous),
		Entry("Value is []bytes, Cancel - Works",
			&types.AttributeValueMemberB{Value: []byte("Cancel")}, Financial_Trades_Cancel),
		Entry("Value is []bytes, CancelRecord - Works",
			&types.AttributeValueMemberB{Value: []byte("CancelRecord")}, Financial_Trades_CancelRecord),
		Entry("Value is []bytes, ErrorRecord - Works",
			&types.AttributeValueMemberB{Value: []byte("ErrorRecord")}, Financial_Trades_ErrorRecord),
		Entry("Value is []bytes, CorrectionRecord - Works",
			&types.AttributeValueMemberB{Value: []byte("CorrectionRecord")}, Financial_Trades_CorrectionRecord),
		Entry("Value is numeric, 0 - Works",
			&types.AttributeValueMemberN{Value: "0"}, Financial_Trades_NotCorrected),
		Entry("Value is numeric, 1 - Works",
			&types.AttributeValueMemberN{Value: "1"}, Financial_Trades_LateCorrected),
		Entry("Value is numeric, 7 - Works",
			&types.AttributeValueMemberN{Value: "7"}, Financial_Trades_Erroneous),
		Entry("Value is numeric, 8 - Works",
			&types.AttributeValueMemberN{Value: "8"}, Financial_Trades_Cancel),
		Entry("Value is numeric, 10 - Works",
			&types.AttributeValueMemberN{Value: "10"}, Financial_Trades_CancelRecord),
		Entry("Value is numeric, 11 - Works",
			&types.AttributeValueMemberN{Value: "11"}, Financial_Trades_ErrorRecord),
		Entry("Value is numeric, 12 - Works",
			&types.AttributeValueMemberN{Value: "12"}, Financial_Trades_CorrectionRecord),
		Entry("Value is NULL - Works", new(types.AttributeValueMemberNULL), Financial_Trades_CorrectionCode(0)),
		Entry("Value is string, NotCorrected - Works",
			&types.AttributeValueMemberS{Value: "NotCorrected"}, Financial_Trades_NotCorrected),
		Entry("Value is string, LateCorrected - Works",
			&types.AttributeValueMemberS{Value: "LateCorrected"}, Financial_Trades_LateCorrected),
		Entry("Value is string, Erroneous - Works",
			&types.AttributeValueMemberS{Value: "Erroneous"}, Financial_Trades_Erroneous),
		Entry("Value is string, Cancel - Works",
			&types.AttributeValueMemberS{Value: "Cancel"}, Financial_Trades_Cancel),
		Entry("Value is string, CancelRecord - Works",
			&types.AttributeValueMemberS{Value: "CancelRecord"}, Financial_Trades_CancelRecord),
		Entry("Value is string, ErrorRecord - Works",
			&types.AttributeValueMemberS{Value: "ErrorRecord"}, Financial_Trades_ErrorRecord),
		Entry("Value is string, CorrectionRecord - Works",
			&types.AttributeValueMemberS{Value: "CorrectionRecord"}, Financial_Trades_CorrectionRecord))

	// Test that attempting to deserialize a Financial.Trades.CorrectionCode will fial and return an
	// error if the value cannot be converted to either the name value or integer value of the enum option
	It("Scan - Value is nil - Error", func() {

		// Attempt to convert a fake string value into a Financial.Trades.CorrectionCode
		// This should return an error
		var enum *Financial_Trades_CorrectionCode
		err := enum.Scan(nil)

		// Verify the error
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(Equal("value of %!q(<nil>) had an invalid type of <nil>"))
		Expect(enum).Should(BeNil())
	})

	// Test the conditions under which values should be convertible to a Financial.Trades.CorrectionCode
	DescribeTable("Scan Tests",
		func(value interface{}, shouldBe Financial_Trades_CorrectionCode) {

			// Attempt to convert the value into a Financial.Trades.CorrectionCode
			// This should not fail
			var enum Financial_Trades_CorrectionCode
			err := enum.Scan(value)

			// Verify that the deserialization was successful
			Expect(err).ShouldNot(HaveOccurred())
			Expect(enum).Should(Equal(shouldBe))
		},
		Entry("NotCorrected - Works", "NotCorrected", Financial_Trades_NotCorrected),
		Entry("LateCorrected - Works", "LateCorrected", Financial_Trades_LateCorrected),
		Entry("Erroneous - Works", "Erroneous", Financial_Trades_Erroneous),
		Entry("Cancel - Works", "Cancel", Financial_Trades_Cancel),
		Entry("CancelRecord - Works", "CancelRecord", Financial_Trades_CancelRecord),
		Entry("ErrorRecord - Works", "ErrorRecord", Financial_Trades_ErrorRecord),
		Entry("CorrectionRecord - Works", "CorrectionRecord", Financial_Trades_CorrectionRecord),
		Entry("00 - Works", "00", Financial_Trades_NotCorrected),
		Entry("01 - Works", "01", Financial_Trades_LateCorrected),
		Entry("07 - Works", "07", Financial_Trades_Erroneous),
		Entry("08 - Works", "08", Financial_Trades_Cancel),
		Entry("0 - Works", 0, Financial_Trades_NotCorrected),
		Entry("1 - Works", 1, Financial_Trades_LateCorrected),
		Entry("7 - Works", 7, Financial_Trades_Erroneous),
		Entry("8 - Works", 8, Financial_Trades_Cancel),
		Entry("10 - Works", 10, Financial_Trades_CancelRecord),
		Entry("11 - Works", 11, Financial_Trades_ErrorRecord),
		Entry("12 - Works", 12, Financial_Trades_CorrectionRecord))
})
