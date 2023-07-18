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
		cells := strings.Split(row, colSep)

		parsed[i] = make([]string, len(cells))
		for idx, cell := range cells {
			parsed[i][idx] = strings.TrimSpace(cell)
		}
	}

	return parsed
}
