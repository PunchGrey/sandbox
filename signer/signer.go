package main

import (
	"fmt"
)

func main() {

	folowJobs := []job{
		job(func(in, out chan interface{}) {
			out <- int(3)
			close(out)
		}),
		job(func(in, out chan interface{}) {
			val := <-in
			out <- 2 * val.(int)
		}),
		job(func(in, out chan interface{}) {
			fmt.Println(<-in)
		}),
	}

	ExecutePipeline(folowJobs...)
	fmt.Scanln()

}

//ExecutePipeline функция обрабатывающая последлвательно массив функций, типа job
func ExecutePipeline(jobs ...job) {
	in := make(chan interface{}, 1)
	out := make(chan interface{}, 1)
	for i, itemJob := range jobs {
		if i%2 != 0 {
			fmt.Println("not", i%2)
			itemJob(in, out)
		} else {
			fmt.Println("eq", i%2)
			itemJob(out, in)
		}
	}

}
