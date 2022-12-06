package gopb

import (
	"encoding/json"

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
