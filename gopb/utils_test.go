package gopb

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Create a new test runner we'll use to test all the
// modules in the utils package
func TestUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils Suite")
}

var _ = Describe("UnixTimestamp Tests", func() {

	// Test that the New function creates a valid timestamp from a time
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
		timestamp := UnixTimestamp{
			Seconds:     1654127993,
			Nanoseconds: 983651350,
		}

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
			generateTimestamp(1655510000, 900838091), generateTimestamp(3311020001, 801676182)))

	// Tests that the AddDuration function works under various conditions
	DescribeTable("AddDuration - Works",
		func(duration time.Duration, result *UnixTimestamp) {
			timestamp := generateTimestamp(1655510000, 900838091)
			Expect(timestamp.AddDuration(duration)).Should(Equal(result))
		},
		Entry("Add nanoseconds - Works", 15*time.Nanosecond, generateTimestamp(1655510000, 900838106)),
		Entry("Add microseconds - Works", 15*time.Microsecond, generateTimestamp(1655510000, 900853091)),
		Entry("Add milliseconds - Works", 15*time.Millisecond, generateTimestamp(1655510000, 915838091)),
		Entry("Add seconds - Works", 15*time.Second, generateTimestamp(1655510015, 900838091)),
		Entry("Add minutes - Works", 15*time.Minute, generateTimestamp(1655510900, 900838091)),
		Entry("Add hours - Works", 15*time.Hour, generateTimestamp(1655564000, 900838091)))

	// Tests the conditions determining whether IsValid will return true or false
	DescribeTable("IsValid - Conditions",
		func(timestamp *UnixTimestamp, result bool) {
			Expect(timestamp.IsValid()).Should(Equal(result))
		},
		Entry("Timestamp is nil - False", nil, false),
		Entry("Seconds < Minimum Timestamp - False",
			&UnixTimestamp{Seconds: -62135596801, Nanoseconds: 983651350}, false),
		Entry("Seconds > Maximum Timestamp - False",
			&UnixTimestamp{Seconds: 253402300800, Nanoseconds: 983651350}, false),
		Entry("Nanoseconds > 1 second - False",
			&UnixTimestamp{Seconds: 1654127993, Nanoseconds: 1000000000}, false),
		Entry("Nanoseconds negative - False",
			&UnixTimestamp{Seconds: 1654127993, Nanoseconds: -1}, false),
		Entry("Valid - True", &UnixTimestamp{Seconds: 1654127993, Nanoseconds: 983651350}, true))

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
		Entry("Seconds < Minimum Timestamp - False",
			&UnixTimestamp{Seconds: -62135596801, Nanoseconds: 983651350}, true,
			"timestamp (-62135596801, 983651350) before 0001-01-01"),
		Entry("Seconds > Maximum Timestamp - False",
			&UnixTimestamp{Seconds: 253402300800, Nanoseconds: 983651350}, true,
			"timestamp (253402300800, 983651350) after 9999-12-31"),
		Entry("Nanoseconds > 1 second - False",
			&UnixTimestamp{Seconds: 1654127993, Nanoseconds: 1000000000}, true,
			"timestamp (1654127993, 1000000000) has out-of-range nanos"),
		Entry("Nanoseconds negative - False",
			&UnixTimestamp{Seconds: 1654127993, Nanoseconds: -1}, true,
			"timestamp (1654127993, -1) has out-of-range nanos"),
		Entry("Valid - True", &UnixTimestamp{Seconds: 1654127993, Nanoseconds: 983651350}, false, ""))

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
		func(rawValue string, message string) {

			// Attempt to convert a fake string value into a Timestamp
			// This should return an error
			timestamp := new(UnixTimestamp)
			err := timestamp.Scan(rawValue)

			// Verify the error
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(Equal(message))
		},
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
})

// Helper function that generates a UnixTimestamp from a value for seconds and nanoseconds
func generateTimestamp(seconds int64, nanos int32) *UnixTimestamp {
	return &UnixTimestamp{
		Seconds:     seconds,
		Nanoseconds: nanos,
	}
}
