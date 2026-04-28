package counting

import (
	"runtime"
	"sync"
	"sync/atomic"
)

func GenerateNumbers(totalNumbers int) []int {
	numbers := make([]int, totalNumbers)
	for i := 0; i < totalNumbers; i++ {
		numbers[i] = i
	}
	return numbers
}

func Add(numbers []int) int64 {
	var sum int64
	for _, n := range numbers {
		sum += int64(n)
	}
	return sum
}

func AddConcurrent(numbers []int) int64 {
	numOfCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numOfCPUs)

	var sum int64
	chunkSize := len(numbers) / numOfCPUs

	var wg sync.WaitGroup
	for i := 0; i < numOfCPUs; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == numOfCPUs-1 {
			end = len(numbers)
		}
		wg.Add(1)
		go func(s, e int) {
			defer wg.Done()
			var partialSum int64
			for _, n := range numbers[s:e] {
				partialSum += int64(n)
			}
			atomic.AddInt64(&sum, partialSum)
		}(start, end)
	}
	wg.Wait()
	return sum
}
