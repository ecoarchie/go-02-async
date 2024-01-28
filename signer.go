package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func Crc32ToChan(data string, out chan<- string) {
		out <- DataSignerCrc32(data)
		close(out)
}

func Crc32FromMd5(s string, out chan<- string, mu *sync.Mutex) {
		mu.Lock()
		md5 := DataSignerMd5(s)
		mu.Unlock()
		Crc32ToChan(md5, out)
}

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	for i := range in {
		crs32chan := make(chan string)
		crsmd5chan := make(chan string)
		s := strconv.Itoa(i.(int))

		wg.Add(1)
		go Crc32FromMd5(s, crsmd5chan, mu)
		go Crc32ToChan(s, crs32chan)

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

func Crc32Th(inData interface{}, ch chan<- string, i int) {
	s := strconv.Itoa(i)
	ch <- DataSignerCrc32(s + inData.(string))
	close(ch)
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
				go Crc32Th(inData, chans[i], i)
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