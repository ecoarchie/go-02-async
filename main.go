package main

import (
	"fmt"
	"time"
)

func generator(in, out chan interface{}) {
	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	for _, n := range inputData {
		out <- n
	}

}

func reader(in, out chan interface{}) {
	dataRaw := <-in
	data, ok := dataRaw.(string)
	if !ok {
		fmt.Printf("cant convert result data to string")
	}
	fmt.Println("data = ", data)
}

func main() {

	start := time.Now()
	ExecutePipeline(generator, SingleHash, MultiHash, CombineResults, reader)
	fmt.Printf("Passed: %v\n", time.Since(start))
}
