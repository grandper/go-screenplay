package sentence

import (
	"strconv"
	"strings"
)

// Rank writes a position in English with digits: 1st, 2nd, 3rd, 4th, 11th,
// 21st. It is the rank of a pass of a loop, or of a branch of a choice.
func Rank(position int) string {
	rank := strconv.Itoa(position)

	switch {
	case strings.HasSuffix(rank, "11"), strings.HasSuffix(rank, "12"), strings.HasSuffix(rank, "13"):
		return rank + "th"
	case strings.HasSuffix(rank, "1"):
		return rank + "st"
	case strings.HasSuffix(rank, "2"):
		return rank + "nd"
	case strings.HasSuffix(rank, "3"):
		return rank + "rd"
	default:
		return rank + "th"
	}
}
