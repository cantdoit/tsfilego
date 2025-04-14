package tsfilego

import (
	"fmt"
)

func main() {
	reader, err := tsfile.NewReader("data.tsfile")
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	result, err := reader.Query("sensors", []string{"temp", "pressure"}, 0, 1000000)
	if err != nil {
		panic(err)
	}

	for result.Next() {
		if val, ok := result.GetFloat32("temp"); ok {
			fmt.Printf("Temperature: %.2f\n", val)
		}
	}
}
