package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

func processRow(row []string) int {
	id, _ := strconv.Atoi(row[0])
	value, _ := strconv.Atoi(row[1])
	return id + value
}

func sequentialProcess(filename string) (int, time.Duration) {
	start := time.Now()

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Read() // pula o cabeçalho

	sum := 0
	for {
		row, err := reader.Read()
		if err != nil {
			break
		}
		sum += processRow(row)
	}

	return sum, time.Since(start)
}

func concurrentProcess(filename string, numWorkers int) (int, time.Duration) {
	start := time.Now()

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Read() // pula o cabeçalho

	rows := make(chan []string)
	results := make(chan int)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for row := range rows {
				results <- processRow(row)
			}
		}()
	}

	go func() {
		for {
			row, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			rows <- row
		}
		close(rows)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	sum := 0
	for partial := range results {
		sum += partial
	}

	return sum, time.Since(start)
}

func main() {
	//sequencial
	// sum, duration := sequentialProcess("./generator/data.csv")
	// fmt.Printf("sequential: sum=%d, duration=%s\n", sum, duration)

	//pararelo
	sum, duration := concurrentProcess("./generator/data.csv", 4)
	fmt.Printf("concurrent (4 workers): sum=%d, duration=%s\n", sum, duration)
}
