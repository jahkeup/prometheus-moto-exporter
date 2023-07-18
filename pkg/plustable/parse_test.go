package plustable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const exampleString = `1^Locked^QAM256^33^663.0^-9.6^38.9^32447^8490^|+|2^Locked^QAM256^5^483.0^-10.2^34.3^837031^84292^|+|3^Locked^QAM256^6^489.0^-10.7^34.5^975251^97720^|+|4^Locked^QAM256^7^495.0^-10.5^38.0^791410^45091^|+|5^Locked^QAM256^8^`

// The Motorola DS8611 can return strings like this for the upstream channels; notably these include whitespace within the cells.
const ds8611ExampleString = `1^Locked^SC-QAM^1^5120^24.0^42.3^|+|2^Locked^SC-QAM^2^5120^30.4^43.5^|+|3^Locked^SC-QAM^3^5120^36.8^43.5^|+|4^Locked^OFDMA^4^12000^ 4.6^35.2^`

func TestParse(t *testing.T) {
	testcases := []struct {
		input    string
		expected [][]string
	}{
		{input: "",
			expected: nil},
		{input: "r0c0^r0c1",
			expected: [][]string{[]string{"r0c0", "r0c1"}}},
		{input: "r0c0|+|r1c0",
			expected: [][]string{[]string{"r0c0"}, []string{"r1c0"}}},
		{input: "r0c0|+|r1c0^r1c1",
			expected: [][]string{[]string{"r0c0"}, []string{"r1c0", "r1c1"}}},
	}

	for _, tc := range testcases {
		t.Run(tc.input, func(t *testing.T) {
			tbl := Parse(tc.input)
			assert.ElementsMatch(t, tbl, tc.expected)
		})
	}

}

func TestParseExample(t *testing.T) {
	actual := Parse(exampleString)
	expected := [][]string{
		[]string{"1", "Locked", "QAM256", "33", "663.0", "-9.6", "38.9", "32447", "8490", ""},
		[]string{"2", "Locked", "QAM256", "5", "483.0", "-10.2", "34.3", "837031", "84292", ""},
		[]string{"3", "Locked", "QAM256", "6", "489.0", "-10.7", "34.5", "975251", "97720", ""},
		[]string{"4", "Locked", "QAM256", "7", "495.0", "-10.5", "38.0", "791410", "45091", ""},
		[]string{"5", "Locked", "QAM256", "8", ""}}

	t.Logf("actual: %#v", actual)
	assert.ElementsMatch(t, actual, expected)
}

func TestParseExample_DS8611(t *testing.T) {
	actual := Parse(ds8611ExampleString)
	expected := [][]string{
		[]string{"1", "Locked", "SC-QAM", "1", "5120", "24.0", "42.3", ""},
		[]string{"2", "Locked", "SC-QAM", "2", "5120", "30.4", "43.5", ""},
		[]string{"3", "Locked", "SC-QAM", "3", "5120", "36.8", "43.5", ""},
		[]string{"4", "Locked", "OFDMA", "4", "12000", "4.6", "35.2", ""}}

	t.Logf("actual: %#v", actual)
	assert.ElementsMatch(t, actual, expected)
}
