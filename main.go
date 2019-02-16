package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

func main() {

	started := time.Now()

	fileName := os.Args[1]
	word := os.Args[2]

	processorNumbers := runtime.NumCPU()

	jobs := make(chan string, 100000)
	results := make(chan *string, 100000)

	go func() {
		// start reader
		err := readFile(fileName, jobs)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		// start workers
		wg := &sync.WaitGroup{}
		for i := 0; i <= processorNumbers; i++ {
			wg.Add(1)
			go worker(wg, word, jobs, results)
		}

		wg.Wait()
		close(results)
	}()

	list := calculator(results)

	duration := time.Since(started)

	fmt.Printf("%.2f, %s \n", duration.Seconds()*1000, strings.Join(list, ", "))
}

func readFile(fileName string, jobs chan<- string) error {
	defer close(jobs)

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	s := bufio.NewScanner(reader)
	s.Split(bufio.ScanLines)

	for s.Scan() {
		jobs <- s.Text()
	}

	return nil

}

func worker(wg *sync.WaitGroup, word string, job <-chan string, results chan<- *string) {
	defer wg.Done()

	freq := getFreq(word)

	for {
		select {
		case testable, ok := <-job:
			if !ok {
				return
			}

			results <- isAnagram(freq, testable)
		}
	}
}

func isAnagram(freq1 map[int32]int8, word2 string) *string {
	freq2 := make(map[int32]int8)
	for _, l := range word2 {
		n1, found1 := freq1[l]
		if !found1 {
			return nil
		}

		n2 := freq2[l]
		n2++

		if n2 > n1 {
			return nil
		}

		freq2[l] = n2
	}

	for l1, count := range freq1 {
		count2, found2 := freq2[l1]
		if !found2 {
			return nil
		}

		if count != count2 {
			return nil
		}
	}

	return &word2
}

func getFreq(word string) map[int32]int8 {
	freq := make(map[int32]int8)
	for _, l := range word {
		n := freq[l]
		n++

		freq[l] = n
	}

	return freq
}

func calculator(results <-chan *string) []string {
	list := make([]string, 0)

	for {
		select {
		case word, ok := <-results:
			if !ok {
				return list
			}

			if word != nil {
				list = append(list, *word)
			}
		}
	}
}
