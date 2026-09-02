package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func main() {
	file, err := os.Create("data.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"id", "value"})

	const totalRows = 500_000
	for i := 0; i < totalRows; i++ {
		writer.Write([]string{strconv.Itoa(i), strconv.Itoa(i * 2)})
	}

	fmt.Println("generated", totalRows, "rows")
}
