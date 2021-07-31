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

	t.Run("test top score", func(t *testing.T) {
		scores := map[string]float64{"item1": 0.11, "item2": 0.2, "item3": 0.2, "item4": 0.16, "item5": 0.14}

		topScoreItem := bandit.GetTopScore(scores)

		require.Equal(t, 0.2, topScoreItem)
	})

	t.Run("test items with top score", func(t *testing.T) {
		scores := map[string]float64{"item1": 0.11, "item2": 0.2, "item3": 0.2, "item4": 0.16, "item5": 0.14}

		itemsWithTopScore := bandit.GetItemsWithTopScore(scores, 0.2)

		require.Len(t, itemsWithTopScore, 2)
		require.Contains(t, itemsWithTopScore, "item2")
		require.Contains(t, itemsWithTopScore, "item3")
	})

	t.Run("test nonclicked results rate", func(t *testing.T) {
		items := []string{"item1", "item2", "item3", "item4", "item5"}
		clicks := map[string]int{}
		views := map[string]int{"item1": 5000, "item2": 5000, "item3": 5000, "item4": 5000, "item5": 5000}

		results := map[string]int{}

		for i := 0; i < 10000; i++ {
			item, _ := bandit.Use(items, clicks, views)

			views[item]++
			results[item]++
		}

		for _, value := range results {
			require.Greater(t, value, 1900, "total count should be greater than 1900")
			require.Less(t, value, 2100, "total count should be less than 2100")
		}
	})

	t.Run("test clicked results rate", func(t *testing.T) {
		items := []string{"item1", "item2", "item3", "item4", "item5"}
		clicks := map[string]int{"item3": 1}
		views := map[string]int{"item1": 5000, "item2": 5000, "item3": 5000, "item4": 5000, "item5": 5000}

		results := map[string]int{}

		for i := 0; i < 10000; i++ {
			item, _ := bandit.Use(items, clicks, views)

			views[item]++
			results[item]++
		}

		for key, value := range results {
			if key == "item3" {
				continue
			}

			require.Less(t, value, results["item3"], "total count should be less than item3 value")
		}
	})

	t.Run("test empty slice", func(t *testing.T) {
		item, err := bandit.Use([]string{}, map[string]int{}, map[string]int{})

		require.ErrorIs(t, err, ErrEmptySlice)
		require.Empty(t, item)
	})

	t.Run("test views checker", func(t *testing.T) {
		err := bandit.CheckOneView([]string{"item1", "item2"}, map[string]int{})

		require.ErrorIs(t, err, ErrNoViewsForItem)
	})

	t.Run("test empty views", func(t *testing.T) {
		item, err := bandit.Use([]string{"item1", "item2"}, map[string]int{}, map[string]int{})

		require.ErrorIs(t, err, ErrNoViewsForItem)
		require.Empty(t, item)
	})
}
