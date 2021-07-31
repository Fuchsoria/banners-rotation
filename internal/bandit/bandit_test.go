package bandit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBandit(t *testing.T) {
	bandit := New()

	t.Run("test score value", func(t *testing.T) {
		var views float64 = 1000
		var clicks float64 = 10
		var totalUses float64 = 10000

		score := bandit.GetScore(views, clicks, totalUses)

		// calc result: https://user-images.githubusercontent.com/43413472/127745648-e7d3e3d0-6e50-4119-9da7-0b909f2ba6a2.png
		require.Greater(t, score, 0.1457)
		require.Less(t, score, 0.1458)
	})
}
