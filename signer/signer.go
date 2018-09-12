package main

import (
	"fmt"
	"log"
	"sync"
)

func main() {

	folowJobs := []job{
		job(func(in, out chan interface{}) {
			out <- int(3)
		}),
		job(func(in, out chan interface{}) {
			val := <-in
			out <- 2 * val.(int)
			out <- 3 * val.(int)
			out <- 2 * val.(int)
			out <- 3 * val.(int)
		}),
		job(SingleHash),
		(MultiHash),
		job(func(in, out chan interface{}) {
			fmt.Println("result: ", <-in)
			fmt.Println("result: ", <-in)
			fmt.Println("result: ", <-in)
			fmt.Println("result: ", <-in)
		}),
	}

	ExecutePipeline(folowJobs...)
	fmt.Scanln()

}

//ExecutePipeline функция обрабатывающая последлвательно массив функций, типа job
func ExecutePipeline(jobs ...job) {
	in := make(chan interface{}, 200)
	out := make(chan interface{}, 200)
	for i, itemJob := range jobs {
		if i%2 != 0 {
			fmt.Println("not", i%2)
			itemJob(out, in)
		} else {
			fmt.Println("eq", i%2)
			itemJob(in, out)
		}
	}
}

//SingleHash считает значение crc32(data)+"~"+crc32(md5(data))
func SingleHash(in, out chan interface{}) {
	wgsh := &sync.WaitGroup{}

	for {
		select {
		case dataRaw := <-in:
			wgsh.Add(1)
			fmt.Println(dataRaw)
			go func(dataRaw interface{}, out chan interface{}) {
				defer wgsh.Done()
				intData, ok := (dataRaw.(int))
				if !ok {
					log.Fatal("con't convert data to string")
				}
				data := string(intData)
				result := DataSignerCrc32(data) + "~" + DataSignerCrc32(DataSignerMd5(data))
				out <- result
			}(dataRaw, out)
		default:
			wgsh.Wait()
			return
		}
	}
}

//MultiHash считает значение crc32(th+data)
func MultiHash(in, out chan interface{}) {
	wgsh := &sync.WaitGroup{}
	for {
		select {
		case dataRaw := <-in:
			wgsh.Add(1)
			go func(dataRaw interface{}, out chan interface{}) {
				defer wgsh.Done()
				data, ok := (dataRaw.(string))
				if !ok {
					log.Fatal("con't convert data to string")
				}
				result := multiHashOne(data)
				out <- result
			}(dataRaw, out)
		default:
			wgsh.Wait()
			return
		}
	}
}

func multiHashOne(data string) string {
	type numAndData struct {
		data string
		num  int
	}
	var arrResult [6]string
	var result string
	mh := make(chan numAndData, 6)
	end := make(chan struct{})
	wg := &sync.WaitGroup{}

	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func(th int) {
			mh <- numAndData{data: DataSignerCrc32(string(th) + data), num: th}
		}(i)
	}

	go func(end chan struct{}) {
		wg.Wait()
		end <- struct{}{}
	}(end)

LOOP:
	for {
		select {
		case item := <-mh:
			arrResult[item.num] = item.data
			wg.Done()
		case <-end:
			break LOOP
		}
	}

	for _, item := range arrResult {
		result = result + item
	}
	return result
}
