package plustable

import (
	"strings"
)

const (
	rowSep = "|+|"
	colSep = "^"
)

// Parse a "plus table" string into its table rows.
func Parse(str string) [][]string {
	if str == "" {
		return nil
	}

	rows := strings.Split(str, rowSep)
	parsed := make([][]string, len(rows))

	for i, row := range rows {
		parsed[i] = strings.Split(row, colSep)
	}

	return parsed
}
