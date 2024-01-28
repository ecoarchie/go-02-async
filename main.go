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
	close(out)

}

func reader(in, out chan interface{}) {
	dataRaw := <-in
	data, ok := dataRaw.(string)
	if !ok {
		fmt.Errorf("cant convert result data to string")
	}
	fmt.Println("data = ", data)
	// testResult = data
	// for t := range in {
	// 	fmt.Println(t)
	// }

}

func main() {

	// chan0 := make(chan interface{})
	// chan1 := make(chan interface{})
	// chan2 := make(chan interface{})
	// chan3 := make(chan interface{})
	// chan4 := make(chan interface{})
	// chan5 := make(chan interface{})
	// go generator(chan0, chan1)

	// go SingleHash(chan1, chan2)

	// go MultiHash(chan2, chan3)

	// go CombineResults(chan3, chan4)

	// reader(chan4, chan5)
	start := time.Now()
	ExecutePipeline(generator, SingleHash, MultiHash, CombineResults, reader)
	fmt.Printf("Passed: %v\n", time.Since(start))
}
