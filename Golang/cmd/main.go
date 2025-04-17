package main

import (
	"Golang/internal/writer"

	"fmt"
	"os"
)

func main() {
	// Define file path and permissions
	filePath := "./example.tsfile"
	flags := os.O_RDWR | os.O_CREATE
	mode := os.FileMode(0644)

	// Initialize the TSFile writer
	tsFileWriter := &writer.TsFileWriter{}

	// Open (Create) the file
	err := tsFileWriter.Open(filePath, flags, mode)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}

	// File created successfully
	fmt.Println("TSFile created successfully.")
}
