package bot

import (
	"math"
	"testing"
)

func TestSplit(t *testing.T) {
	text := ""
	LIMIT := 3
	num := int(math.Ceil(float64(len(text))/float64(LIMIT)))
	for i := 0; i < num; i++ {
		start := i * LIMIT
		end := func() int {
			if (i+1)*LIMIT > len(text) {
				return len(text)
			} else {
				return (i + 1) * LIMIT
			}
		}()
		t.Log(text[start:end])
	}
}
