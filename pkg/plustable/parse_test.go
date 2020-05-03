package plustable

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const exampleString = `1^Locked^QAM256^33^663.0^-9.6^38.9^32447^8490^|+|2^Locked^QAM256^5^483.0^-10.2^34.3^837031^84292^|+|3^Locked^QAM256^6^489.0^-10.7^34.5^975251^97720^|+|4^Locked^QAM256^7^495.0^-10.5^38.0^791410^45091^|+|5^Locked^QAM256^8^`

func TestParse(t *testing.T) {
	testcases := []struct {
		input    string
		expected [][]string
	}{
		{input: "",
			expected: [][]string{}},
		{input: "r0c0^r0c1",
			expected: [][]string{}},
	}

	for _, tc := range testcases {
		t.Run(tc.input, func(t *testing.T) {
			tbl := Parse(tc.input)
			require.NotNil(t, tbl)
			assert.ElementsMatch(t, tbl, tc.expected)
		})
	}

}
