package segment

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:generate go run gen_test_cases.go --prefix wordBreak --dataPath data/WordBreakTest.txt --outputPath word_break_test_cases.go

func TestWordBreaker(t *testing.T) {
	for i, tc := range wordBreakTestCases() {
		t.Run(strconv.FormatInt(int64(i), 10), func(t *testing.T) {
			var wb WordBreaker
			var segments [][]rune
			var seg []rune
			for _, r := range tc.inputString {
				canBreakBefore := wb.ProcessRune(r)
				if len(seg) > 0 && canBreakBefore {
					segments = append(segments, seg)
					seg = nil
				}
				seg = append(seg, r)
			}
			if len(seg) > 0 {
				segments = append(segments, seg)
			}
			assert.Equal(t, tc.segments, segments, tc.description)
		})
	}
}
