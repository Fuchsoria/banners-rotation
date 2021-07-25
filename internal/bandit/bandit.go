package bandit

import (
	"math"
)

type Bandit struct{}

func (b *Bandit) GetScore(viewsCount float64, clicksCount float64, totalUses float64) float64 {
	clickViewRate := clicksCount / viewsCount
	banditRate := math.Sqrt((math.Log2(totalUses) / viewsCount))

	return clickViewRate + banditRate
}

func (b *Bandit) GetTopScore(scores map[string]float64) string {
	maxValue := 0.0
	maxValueKey := ""

	for key, score := range scores {
		if score > maxValue {
			maxValue = score
			maxValueKey = key
		}
	}

	return maxValueKey
}

func (b *Bandit) Use(items []string, clicks map[string]int, views map[string]int) string {
	itemsScore := make(map[string]float64)

	for _, item := range items {
		itemsScore[item] = b.GetScore(float64(views[item]), float64(clicks[item]), float64(len(views)))
	}

	itemID := b.GetTopScore(itemsScore)

	return itemID
}

func New() *Bandit {
	return &Bandit{}
}
