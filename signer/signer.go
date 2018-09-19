package main

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"log"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

func main() {

	/*	folowJobs := []job{
		job(func(in, out chan interface{}) {
			out <- int(3)
			//			close(out)
		}),
		job(func(in, out chan interface{}) {
			val := <-in
			out <- 2 * val.(int)
			out <- 3 * val.(int)
			out <- 2 * val.(int)
			out <- 3 * val.(int)
			out <- 5 * val.(int)
			//			close(out)
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			fmt.Println("result: ", <-in)
		}),
	}*/

	//	testExpected := "1173136728138862632818075107442090076184424490584241521304_1696913515191343735512658979631549563179965036907783101867_27225454331033649287118297354036464389062965355426795162684_29568666068035183841425683795340791879727309630931025356555_3994492081516972096677631278379039212655368881548151736_4958044192186797981418233587017209679042592862002427381542_4958044192186797981418233587017209679042592862002427381542"
	testResult := "NOT_SET"

	// это небольшая защита от попыток не вызывать мои функции расчета
	// я преопределяю фукции на свои которые инкрементят локальный счетчик
	// переопределение возможо потому что я объявил функцию как переменную, в которой лежит функция
	var (
		DataSignerSalt         string = "" // на сервере будет другое значение
		OverheatLockCounter    uint32
		OverheatUnlockCounter  uint32
		DataSignerMd5Counter   uint32
		DataSignerCrc32Counter uint32
	)
	OverheatLock = func() {
		atomic.AddUint32(&OverheatLockCounter, 1)
		for {
			if swapped := atomic.CompareAndSwapUint32(&dataSignerOverheat, 0, 1); !swapped {
				fmt.Println("OverheatLock happend")
				time.Sleep(time.Second)
			} else {
				break
			}
		}
	}
	OverheatUnlock = func() {
		atomic.AddUint32(&OverheatUnlockCounter, 1)
		for {
			if swapped := atomic.CompareAndSwapUint32(&dataSignerOverheat, 1, 0); !swapped {
				fmt.Println("OverheatUnlock happend")
				time.Sleep(time.Second)
			} else {
				break
			}
		}
	}
	DataSignerMd5 = func(data string) string {
		atomic.AddUint32(&DataSignerMd5Counter, 1)
		OverheatLock()
		defer OverheatUnlock()
		data += DataSignerSalt
		dataHash := fmt.Sprintf("%x", md5.Sum([]byte(data)))
		time.Sleep(10 * time.Millisecond)
		return dataHash
	}
	DataSignerCrc32 = func(data string) string {
		atomic.AddUint32(&DataSignerCrc32Counter, 1)
		data += DataSignerSalt
		crcH := crc32.ChecksumIEEE([]byte(data))
		dataHash := strconv.FormatUint(uint64(crcH), 10)
		time.Sleep(time.Second)
		return dataHash
	}

	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	// inputData := []int{0,1}

	folowJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			dataRaw := <-in
			data, ok := dataRaw.(string)
			if !ok {
				fmt.Println("cant convert result data to string")
			}
			testResult = data
			fmt.Println(testResult)
		}),
	}

	ExecutePipeline(folowJobs...)
	//	fmt.Println(forTestResult(0))

	fmt.Scanln()

}

//ExecutePipeline функция обрабатывающая последлвательно массив функций, типа job
func ExecutePipeline(jobs ...job) {
	wgGlob := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	l := len(jobs)
	if l == 0 {
		return
	}

	tube := make([]chan interface{}, 0, l)
	for i := 0; i < l; i++ {
		tube = append(tube, make(chan interface{}, 50))
	}

	wgGlob.Add(1)
	go func() {
		defer wgGlob.Done()
		defer close(tube[0])
		jobs[0](make(chan interface{}, 50), tube[0])
	}()

	for i, itemJob := range jobs[1:] {
		wgGlob.Add(1)

		go func(i int, itemJob job) {
			defer wgGlob.Done()
			defer close(tube[i])
			mu.Lock()
			itemJob(tube[i-1], tube[i])

			mu.Unlock()
		}(i+1, itemJob)
		time.Sleep(time.Millisecond)
	}

	wgGlob.Wait()
}

//SingleHash считает значение crc32(data)+"~"+crc32(md5(data))
func SingleHash(in, out chan interface{}) {
	wgsh := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	for {
		select {
		case dataRaw, ok := <-in:
			if !ok {
				wgsh.Wait()
				return
			}
			wgsh.Add(1)
			go func(dataRaw interface{}, out chan interface{}) {
				defer wgsh.Done()
				oneCh := make(chan string, 1)
				twoCh := make(chan string, 1)

				var dataMd5 string

				intData, ok := (dataRaw.(int))
				if !ok {
					log.Fatal("con't convert data to string")
				}
				data := strconv.Itoa(intData)
				mu.Lock()
				for {
					if dataSignerOverheat == 0 {
						dataMd5 = DataSignerMd5(data)
						break
					}
				}
				mu.Unlock()

				go func() {
					oneCh <- DataSignerCrc32(data)
					close(oneCh)
				}()
				go func() {
					twoCh <- DataSignerCrc32(dataMd5)
					close(twoCh)
				}()

				//result := DataSignerCrc32(data) + "~" + DataSignerCrc32(dataMd5)
				result := <-oneCh + "~" + <-twoCh
				out <- result
			}(dataRaw, out)
		}
	}
}

//MultiHash считает значение crc32(th+data)
func MultiHash(in, out chan interface{}) {
	wgsh := &sync.WaitGroup{}
	for {
		select {
		case dataRaw, ok := <-in:
			if !ok {
				wgsh.Wait()
				return
			}
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
			mh <- numAndData{data: DataSignerCrc32(strconv.Itoa(th) + data), num: th}
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

//CombineResults получает все результаты, сортирует, объединяет отсортированный результат через _
func CombineResults(in, out chan interface{}) {
	var sliceResult []string
	var result string
LOOPC:
	for {
		select {
		case dataRaw, ok := <-in:
			if !ok {
				break LOOPC
			}
			sliceResult = append(sliceResult, dataRaw.(string))
		}
	}
	sort.Strings(sliceResult)
	for i, item := range sliceResult {
		if i == 0 {
			result = item
		} else {
			result = result + "_" + item
		}
	}
	out <- result
}
