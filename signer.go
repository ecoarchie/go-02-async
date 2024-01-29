package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

func Crc32ToChan(data string, out chan<- string) {
		out <- DataSignerCrc32(data)
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
			out <- res
			wg.Done()
		}()
	}
	wg.Wait()
}

func Crc32Th(inData interface{}, ch chan<- string, i int) {
	s := strconv.Itoa(i)
	ch <- DataSignerCrc32(s + inData.(string))
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
				close(ch)
			}

			out <- step
			wg.Done()

		}(inData)
	}
	wg.Wait()
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
}

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	in := make(chan interface{})
	for _, j := range jobs {
		wg.Add(1)
		out := make(chan interface{})
		go func(job func(in, out chan interface{}), wait *sync.WaitGroup, in, out chan interface{}) {
			defer wait.Done()

			job(in, out)
			close(out)
		}(j, wg, in, out)
		in = out
	}
	wg.Wait()
}
