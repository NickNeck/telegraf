package json

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	validJSON              = "{\"fields\":{\"a\":2.3,\"b\":1.1},\"name\":\"foo\",\"tags\":{\"x\":\"u\",\"y\":4},\"timestamp\":1510756542}"
	validJSONNewline       = "\n{\"fields\":{\"a\":2.3,\"b\":1.1},\"name\":\"foo\",\"tags\":{\"x\":\"u\",\"y\":4},\"timestamp\":1510756542}\n\n"
	validJSONArray         = "[{\"fields\":{\"a\":2.3,\"b\":1.1},\"name\":\"foo\",\"tags\":{\"x\":\"u\",\"y\":4},\"timestamp\":1510756542}]"
	validJSONArrayMultiple = "[{\"fields\":{\"a\":2.3,\"b\":1.1},\"name\":\"foo\",\"tags\":{\"x\":\"u\",\"y\":4},\"timestamp\":1510756542},{\"fields\":{\"c\":3.3,\"d\":2.1},\"name\":\"bar\",\"tags\":{\"x2\":\"u\",\"z\":4},\"timestamp\":1510756540}]"
	invalidJSON            = "I don't think this is JSON"
	invalidJSON2           = "{\"a\": 5, \"b\": \"c\": 6}}"
)

const validJSONTags = `
{
    "a": 5,
    "b": {
        "c": 6
    },
    "mytag": "foobar",
    "othertag": "baz"
}
`

const validJSONArrayTags = `
[
{
    "a": 5,
    "b": {
        "c": 6
    },
    "mytag": "foo",
    "othertag": "baz"
},
{
    "a": 7,
    "b": {
        "c": 8
    },
    "mytag": "bar",
    "othertag": "baz"
}
]
`

func TestParseValidJSON(t *testing.T) {
	parser := JSONParser{}

	// Most basic vanilla test
	metrics, err := parser.Parse([]byte(validJSON))
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, "foo", metrics[0].Name())
	assert.Equal(t, map[string]interface{}{
		"a": float64(2.3),
		"b": float64(1.1),
	}, metrics[0].Fields())
	assert.Equal(t, map[string]string{
		"x": "u",
		"y": "4",
	}, metrics[0].Tags())

	// Test that newlines are fine
	metrics, err = parser.Parse([]byte(validJSONNewline))
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, "foo", metrics[0].Name())
	assert.Equal(t, map[string]interface{}{
		"a": float64(2.3),
		"b": float64(1.1),
	}, metrics[0].Fields())
	assert.Equal(t, map[string]string{
		"x": "u",
		"y": "4",
	}, metrics[0].Tags())

	// Test that whitespace only will parse as an empty list of metrics
	metrics, err = parser.Parse([]byte("\n\t"))
	assert.NoError(t, err)
	assert.Len(t, metrics, 0)

	// Test that an empty string will parse as an empty list of metrics
	metrics, err = parser.Parse([]byte(""))
	assert.NoError(t, err)
	assert.Len(t, metrics, 0)
}

func TestParseLineValidJSON(t *testing.T) {
	parser := JSONParser{
		MetricName: "json_test",
	}

	// Most basic vanilla test
	metric, err := parser.ParseLine(validJSON)
	assert.NoError(t, err)
	assert.Equal(t, "foo", metric.Name())
	assert.Equal(t, map[string]interface{}{
		"a": float64(2.3),
		"b": float64(1.1),
	}, metric.Fields())
	assert.Equal(t, map[string]string{
		"x": "u",
		"y": "4",
	}, metric.Tags())

	// Test that newlines are fine
	metric, err = parser.ParseLine(validJSONNewline)

	assert.NoError(t, err)
	assert.Equal(t, "foo", metric.Name())
	assert.Equal(t, map[string]interface{}{
		"a": float64(2.3),
		"b": float64(1.1),
	}, metric.Fields())
	assert.Equal(t, map[string]string{
		"x": "u",
		"y": "4",
	}, metric.Tags())
}

func TestParseInvalidJSON(t *testing.T) {
	parser := JSONParser{}

	_, err := parser.Parse([]byte(invalidJSON))
	assert.Error(t, err)
	_, err = parser.Parse([]byte(invalidJSON2))
	assert.Error(t, err)
	_, err = parser.ParseLine(invalidJSON)
	assert.Error(t, err)
}

// Test that json arrays can be parsed
func TestParseValidJSONArray(t *testing.T) {
	parser := JSONParser{
		MetricName: "json_array_test",
	}

	// Most basic vanilla test
	metrics, err := parser.Parse([]byte(validJSONArray))
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, "foo", metrics[0].Name())
	assert.Equal(t, map[string]interface{}{
		"a": float64(2.3),
		"b": float64(1.1),
	}, metrics[0].Fields())
	assert.Equal(t, map[string]string{
		"x": "u",
		"y": "4",
	}, metrics[0].Tags())

	// Basic multiple datapoints
	metrics, err = parser.Parse([]byte(validJSONArrayMultiple))
	assert.NoError(t, err)
	assert.Len(t, metrics, 2)
	assert.Equal(t, "foo", metrics[0].Name())
	assert.Equal(t, map[string]interface{}{
		"a": float64(2.3),
		"b": float64(1.1),
	}, metrics[0].Fields())
	assert.Equal(t, map[string]string{
		"x": "u",
		"y": "4",
	}, metrics[0].Tags())
	assert.Equal(t, "bar", metrics[1].Name())
	assert.Equal(t, map[string]interface{}{
		"c": float64(3.3),
		"d": float64(2.1),
	}, metrics[1].Fields())
	assert.Equal(t, map[string]string{
		"x2": "u",
		"z":  "4",
	}, metrics[1].Tags())
}
