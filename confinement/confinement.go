// Package confinement shows how to use a confinement to limit the go routine access to a specific data.
// It avoids race condition by accessing the data concurrently.
package confinement

import (
	"log/slog"
	"sync"
	"time"
)

func process(data int) int {
	// simulate time consuming operation
	time.Sleep(time.Second)

	return data * 2
}

func processData(wg *sync.WaitGroup, resultDest *int, data int) {
	defer wg.Done()

	processedData := process(data)

	*resultDest = processedData
}

// Run apply operation using go routine.
func Run(input []int) []int {
	start := time.Now()

	var wg sync.WaitGroup

	result := make([]int, len(input))

	for i, data := range input {
		wg.Add(1)

		go processData(&wg, &result[i], data)
	}

	wg.Wait()

	slog.Info("Run finished successfully",
		slog.String("duration", time.Since(start).String()),
		slog.Any("result", result),
	)

	return result
}
