package bandit

import (
	"errors"
	"math"
	"math/rand"
)

type Bandit struct{}

var (
	ErrEmptySlice     = errors.New("slice cannot be empty")
	ErrNoViewsForItem = errors.New("some items does not have views")
)

func (b *Bandit) GetScore(viewsCount float64, clicksCount float64, totalUses float64) float64 {
	clickViewRate := clicksCount / viewsCount
	banditRate := math.Sqrt(((2 * math.Log(totalUses)) / viewsCount))

	return clickViewRate + banditRate
}

func (b *Bandit) GetTopScore(scores map[string]float64) float64 {
	maxValue := 0.0

	for _, score := range scores {
		if score > maxValue {
			maxValue = score
		}
	}

	return maxValue
}

func (b *Bandit) GetItemsWithTopScore(scores map[string]float64, topScore float64) []string {
	sameScoreItems := []string{}

	for key, score := range scores {
		if score == topScore {
			sameScoreItems = append(sameScoreItems, key)
		}
	}

	return sameScoreItems
}

func (b *Bandit) GetRandomItemFromTop(topItems []string) string {
	if len(topItems) == 0 {
		return ""
	}

	rand.Shuffle(len(topItems), func(i, j int) { topItems[i], topItems[j] = topItems[j], topItems[i] })

	return topItems[0]
}

func (b *Bandit) CheckOneView(items []string, views map[string]int) error {
	for _, item := range items {
		if views[item] == 0 {
			return ErrNoViewsForItem
		}
	}

	return nil
}

func (b *Bandit) Use(items []string, clicks map[string]int, views map[string]int) (string, error) {
	if len(items) == 0 {
		return "", ErrEmptySlice
	}

	if err := b.CheckOneView(items, views); err != nil {
		return "", err
	}

	itemsScore := make(map[string]float64)

	for _, item := range items {
		itemsScore[item] = b.GetScore(float64(views[item]), float64(clicks[item]), float64(len(views)))
	}

	topScore := b.GetTopScore(itemsScore)
	itemsWithTopScore := b.GetItemsWithTopScore(itemsScore, topScore)
	itemID := b.GetRandomItemFromTop(itemsWithTopScore)

	return itemID, nil
}

func New() *Bandit {
	return &Bandit{}
}
