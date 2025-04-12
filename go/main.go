package main

import (
	"fmt"
	"go/gowrapper"
)

func main() {
	// Open TsFile
	reader, err := gowrapper.NewReader("data.ts")
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	// Read data (example)
	fmt.Println("TsFile opened successfully!")
	// batch, err := reader.ReadBatch()
	// ... (implement based on your query logic)
}
