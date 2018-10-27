package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Tuple struct {
	ind int
	str string
}

func ExtractData(dataRaw interface{}) string {
	switch dataRaw.(type) {
	case int:
		data, _ := dataRaw.(int)
		return strconv.Itoa(data)
	case string:
		data, _ := dataRaw.(string)
		return data
	case Tuple:
		tuple, _ := dataRaw.(Tuple)
		return strconv.Itoa(tuple.ind) + tuple.str
	}
	return ""
}

func calculateSingleHash(wg *sync.WaitGroup, input string, md5 string, out chan interface{}) {
	defer wg.Done()
	firstComponent := make(chan interface{})
	secondComponent := make(chan interface{})

	go func() {
		firstComponent <- DataSignerCrc32(input)
	}()

	go func() {
		secondComponent <- DataSignerCrc32(md5)
	}()

	out <- ExtractData(<-firstComponent) + "~" + ExtractData(<-secondComponent)
}

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	for input := range in {
		wg.Add(1)
		str := ExtractData(input)
		md5 := DataSignerMd5(str)
		go calculateSingleHash(wg, str, md5, out)
	}
}

func calculateMultiHashSingle(wg *sync.WaitGroup, data Tuple, out chan Tuple) {
	defer wg.Done()
	out <- Tuple{str: DataSignerCrc32(ExtractData(data)), ind: data.ind}
}

func calculateMultiHash(wg *sync.WaitGroup, input string, out chan interface{}) {
	defer wg.Done()
	wgSingle := &sync.WaitGroup{}
	outSingle := make(chan Tuple, 6)
	wgSingle.Add(6)
	for i := 0; i <= 5; i++ {
		go calculateMultiHashSingle(wgSingle, Tuple{ind: i, str: input}, outSingle)
	}
	wgSingle.Wait()
	close(outSingle)

	var list []Tuple
	for hash := range outSingle {
		list = append(list, hash)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].ind < list[j].ind
	})

	builder := strings.Builder{}
	for _, elem := range list {
		builder.WriteString(elem.str)
	}
	out <- builder.String()
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	for input := range in {
		wg.Add(1)
		go calculateMultiHash(wg, ExtractData(input), out)
	}
}

func CombineResults(in, out chan interface{}) {
	inputs := make([]string, 0)
	for data := range in {
		inputs = append(inputs, ExtractData(data))
	}
	sort.Strings(inputs)
	out <- strings.Join(inputs, "_")
}

func ExecuteJob(wg *sync.WaitGroup, currentJob job, in, out chan interface{}){
	defer wg.Done()
	currentJob(in, out)
	close(out)
}

func ExecutePipeline(jobs []job) {
	wg := &sync.WaitGroup{}
	currentIn := make(chan interface{})
	currentOut := make(chan interface{})

	wg.Add(len(jobs))
	for _, currentJob := range jobs {
		go ExecuteJob(wg, currentJob, currentIn, currentOut)
		currentIn = currentOut
		currentOut = make(chan interface{})
	}
	wg.Wait()
	close(currentOut)
}
