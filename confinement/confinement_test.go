package confinement

import (
	"slices"
	"testing"
	"time"
)

func Test_confinement(t *testing.T) {
	expected := []int{2, 4, 6, 8, 10}

	start := time.Now()
	actual := Run([]int{1, 2, 3, 4, 5})
	last := time.Since(start)

	if !slices.Equal(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}

	maxDuration := time.Second + 10*time.Millisecond
	if last < time.Second || last > maxDuration {
		t.Errorf("expected: %v, actual: %v", maxDuration, last)
	}
}
