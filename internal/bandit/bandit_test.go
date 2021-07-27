package bandit

import (
	"fmt"
	"testing"
)

func TestBandit(t *testing.T) {
	bandit := New()

	fmt.Println(bandit)
}
