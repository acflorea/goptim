package main

import (
	"fmt"
	"github.com/acflorea/goptim/rand"
)

func main() {
	fmt.Println("Test")

	for i := 0; i < 10; i++ {
		fmt.Println(rand.Float64(10, 15))
	}
}
