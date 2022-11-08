package gopb

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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

var _ = Describe("UnixDuration Marshal/Unmarshal Tests", func() {

	// Test that converting a Duration to JSON works for all values
	DescribeTable("MarshalJSON Tests",
		func(duration *UnixDuration, value string) {
			data, err := duration.MarshalJSON()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("Duration is nil - Works", nil, ""),
		Entry("Duration has value - Works",
			&UnixDuration{Seconds: 1654127993, Nanoseconds: 983651350}, "1654127993983651350"))

	// Test that converting a Duration to a CSV column works for all values
	DescribeTable("MarshalCSV Tests",
		func(duration *UnixDuration, value string) {
			data, err := duration.MarshalCSV()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(value))
		},
		Entry("Duration is nil - Works", nil, ""),
		Entry("Duration has value - Works",
			&UnixDuration{Seconds: 1654127993, Nanoseconds: 983651350}, "1654127993983651350"))

	// Test that converting a Duration to a AttributeValue works for all values
	DescribeTable("MarshalDynamoDBAttributeValue Tests",
		func(duration *UnixDuration, value string) {
			data, err := duration.MarshalDynamoDBAttributeValue()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data.(*types.AttributeValueMemberS).Value).Should(Equal(value))
		},
		Entry("Duration is nil - Works", nil, ""),
		Entry("Duration has value - Works",
			&UnixDuration{Seconds: 1654127993, Nanoseconds: 983651350}, "1654127993983651350"))

	// Test that converting a Duration to an SQL value for all values
	DescribeTable("Value Tests",
		func(duration *UnixDuration, value string) {
			data, err := duration.Value()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data).Should(Equal(value))
		},
		Entry("Duration is nil - Works", nil, ""),
		Entry("Duration has value - Works",
			&UnixDuration{Seconds: 1654127993, Nanoseconds: 983651350}, "1654127993983651350"))

	// Test that attempting to deserialize a Duration will fail and return an error if the
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
		Entry("Seconds < Minimum Duration - Error", "-62135596801983651350", false,
			"duration (-62135596801, 983651350) before 0001-01-01"),
		Entry("Seconds > Maximum Duration - Error", "253402300800983651350", false,
			"duration (253402300800, 983651350) after 9999-12-31"))

	// Test that, if UnmarshalJSON is called with a value of nil then the duration will be nil
	It("UnmarshalJSON - Nil - Nil", func() {

		// Attempt to convert a non-parseable string value into a duration
		// This should not return an error
		var duration *UnixDuration
		err := duration.UnmarshalJSON(nil)
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the duration
		Expect(duration).Should(BeNil())
	})

	// Test that, if UnmarshalJSON is called with an empty string then the duration will be nil
	It("UnmarshalJSON - Empty string - Nil", func() {

		// Attempt to convert a non-parseable string value into a duration
		// This should not return an error
		var duration *UnixDuration
		err := duration.UnmarshalJSON([]byte(""))
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the duration
		Expect(duration).Should(BeNil())
	})

	// Test that, if the UnmarshalJSON function is called with a valid UNIX duration, then it
	// will be parsed into a Duration object
	It("UnmarshalJSON - Non-empty string - Works", func() {

		// Attempt to convert a non-parseable string value into a duration
		// This should not return an error
		var duration *UnixDuration
		err := json.Unmarshal([]byte("1654127993983651350"), &duration)
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the duration
		Expect(duration).ShouldNot(BeNil())
		Expect(duration.Seconds).Should(Equal(int64(1654127993)))
		Expect(duration.Nanoseconds).Should(Equal(int32(983651350)))
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
		Entry("Seconds < Minimum Duration - Error", "-62135596801983651350",
			"duration (-62135596801, 983651350) before 0001-01-01"),
		Entry("Seconds > Maximum Duration - Error", "253402300800983651350",
			"duration (253402300800, 983651350) after 9999-12-31"))

	// Test that, if UnmarshalCSV is called with an empty string then the duration will be nil
	It("UnmarshalCSV - Empty string - Nil", func() {

		// Attempt to convert a non-parseable string value into a duration
		// This should not return an error
		var duration *UnixDuration
		err := duration.UnmarshalCSV("")
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the duration
		Expect(duration).Should(BeNil())
	})

	// Test that, if the UnmarshalCSV function is called with a valid UNIX duration, then it
	// will be parsed into a Duration object
	It("UnmarshalCSV - Non-empty string - Works", func() {

		// Attempt to convert a non-parseable string value into a duration
		// This should not return an error
		duration := new(UnixDuration)
		err := duration.UnmarshalCSV("1654127993983651350")
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the duration
		Expect(duration).ShouldNot(BeNil())
		Expect(duration.Seconds).Should(Equal(int64(1654127993)))
		Expect(duration.Nanoseconds).Should(Equal(int32(983651350)))
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

	// Test that attempting to deserialize a Duration will fail and return an error if the
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
		Entry("Seconds < Minimum Duration - Error", "-62135596801983651350",
			"duration (-62135596801, 983651350) before 0001-01-01"),
		Entry("Seconds > Maximum Duration - Error", "253402300800983651350",
			"duration (253402300800, 983651350) after 9999-12-31"),
		Entry("Nanoseconds > 1 second - Error", "1654127993-10000000",
			"duration (1654127993, -10000000) has out-of-range nanos"))

	// Test that, if Scan is called with a value of nil then the duration will be nil
	It("Scan - Nil - Nil", func() {

		// Attempt to convert nil string value into a duration
		// This should not return an error
		var duration *UnixDuration
		err := duration.Scan(nil)
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the duration
		Expect(duration).Should(BeNil())
	})

	// Test that, if Scan is called with an empty string then the duration will be nil
	It("Scan - Empty string - Nil", func() {

		// Attempt to convert an empty string value into a duration
		// This should not return an error
		var duration *UnixDuration
		err := duration.Scan("")
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the duration
		Expect(duration).Should(BeNil())
	})

	// Test that, if the Scan function is called with a valid UNIX duration, then it
	// will be parsed into a Duration object
	It("Scan - Non-empty string - Works", func() {

		// Attempt to convert a UNIX duration string value into a duration
		// This should not return an error
		duration := new(UnixDuration)
		err := duration.Scan("1654127993983651350")
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the duration
		Expect(duration).ShouldNot(BeNil())
		Expect(duration.Seconds).Should(Equal(int64(1654127993)))
		Expect(duration.Nanoseconds).Should(Equal(int32(983651350)))
	})
})
