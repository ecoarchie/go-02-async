package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	for i := range in {
		crs32chan := make(chan string)
		crsmd5chan := make(chan string)
		s := strconv.Itoa(i.(int))

		wg.Add(1)
		go func(s string, out chan<- string) {
			mu.Lock()
			md5 := DataSignerMd5(s)
			mu.Unlock()
			crsmd := DataSignerCrc32(md5)
			out <- crsmd
			close(out)
		}(s, crsmd5chan)

		go func(s string, out chan<- string) {
			crs32 := DataSignerCrc32(s)
			out <- crs32
			close(out)
		}(s, crs32chan)

		go func() {
			crsmd := <-crsmd5chan
			crs32 := <-crs32chan
			res := crs32 + "~" + crsmd
			fmt.Printf("single Hash for %s - %v\n", s, res)
			out <- res
			wg.Done()
		}()
	}
	wg.Wait()
	close(out)
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for inData := range in {
		wg.Add(1)
		go func(inData interface{}) {
			chans := []chan string{
				make(chan string, 1),
				make(chan string, 1),
				make(chan string, 1),
				make(chan string, 1),
				make(chan string, 1),
				make(chan string, 1),
			}

			for i := 0; i < len(chans); i++ {
				go func(inData interface{}, ch chan string, i int) {
					s := strconv.Itoa(i)
					ch <- DataSignerCrc32(s + inData.(string))
				}(inData, chans[i], i)
			}
			var step string
			for _, ch := range chans {
				step += <-ch
			}

			out <- step
			wg.Done()

		}(inData)
	}
	wg.Wait()
	close(out)
}

// CombineResults получает все результаты, сортирует (https://golang.org/pkg/sort/),
// объединяет отсортированный результат через _ (символ подчеркивания) в одну строку
func CombineResults(in, out chan interface{}) {
	res := make([]string, 0)
	for i := range in {
		res = append(res, i.(string))
	}
	sort.Strings(res)

	o := strings.Join(res, "_")
	out <- o
	close(out)
}

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})
	for i := 0; i < len(jobs); i++ {
		out := make(chan interface{})
		if i != len(jobs)-1 {
			go jobs[i](in, out)
		} else {
			jobs[i](in, out)
		}
		in = out
	}
}
// MultiHash считает значение crc32(th+data)) (конкатенация цифры, приведённой к строке и строки),
// где th=0..5 ( т.е. 6 хешей на каждое входящее значение ), потом берёт конкатенацию результатов в порядке расчета (0..5), где data - то что пришло на вход (и ушло на выход из SingleHash)
// func MultiHash(in, out chan interface{}) {
// 	fmt.Println("here in multilehash")
// 	wg := &sync.WaitGroup{}
// 	// var res string
// 	for inData := range in {
// 		wg.Add(1)
// 		go func(inData interface{}) {
// 			var step string
// 			for i := 0; i < 6; i++ {
// 				s := strconv.Itoa(i)
// 				step += DataSignerCrc32(s + inData.(string))
// 				fmt.Printf("multihash for step %d = %s\n", i, step)
// 			}
// 			fmt.Printf("res of Multihash = %s\n\n", step)
// 			out <- step
// 			wg.Done()
// 		}(inData)
// 	}
// 	wg.Wait()
// 	close(out)
// }

// func MultiHash(in, out chan interface{}) {
// 	fmt.Println("here in multilehash")
// 	wg := &sync.WaitGroup{}
// 	// var res string
// 	for inData := range in {
// 		wg.Add(1)
// 		go func(inData interface{}) {
// 			var step string
// 			output := make(chan string, 6)
// 			wg2 := &sync.WaitGroup{}
// 			wg2.Add(6)
// 			for i := 0; i < 6; i++ {
// 				go func(i int, output chan string) {
// 					var step string
// 					s := strconv.Itoa(i)
// 					step += DataSignerCrc32(s + inData.(string))
// 					fmt.Printf("multihash for step %d = %s\n", i, step)
// 					output <- step
// 					wg2.Done()
// 				}(i, output)
// 			}
// 			wg2.Wait()
// 			close(output)
// 			for o := range output {
// 				step += o
// 			}
// 			fmt.Printf("res of Multihash = %s\n\n", step)
// 			out <- step
// 			wg.Done()
// 		}(inData)
// 	}
// 	wg.Wait()
// 	close(out)
// }