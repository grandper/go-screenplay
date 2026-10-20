package sentence_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay/v2/internal/sentence"
)

func TestRank(t *testing.T) {
	t.Parallel()

	t.Run("writes a position in English", func(t *testing.T) {
		t.Parallel()

		for position, rank := range map[int]string{
			1: "1st", 2: "2nd", 3: "3rd", 4: "4th", 11: "11th", 12: "12th", 13: "13th", 21: "21st", 22: "22nd", 112: "112th",
		} {
			assert.Equal(t, rank, sentence.Rank(position))
		}
	})
}
