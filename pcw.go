package main

import (
	"fmt"
	"sync"
)

var jobs = make(chan int, 100)
var results = make(chan int, 100)
var wg sync.WaitGroup
var done = make(chan bool, 1)

func producer(){
	for j := 1; j <= 20; j++ {
		fmt.Println("Produced job", j)
		jobs <- j
	}
	close(jobs)
}

func writer(){
	for {
		j, more :=  <-results
		if more {
			fmt.Println("Writer received job", j)
		} else {
			fmt.Println("Writer received all jobs")
			done <- true
			return
		}
	}
}

func worker(id int, jobs <-chan int, results chan<- int) {
	for {
		j, more := <-jobs
		if more {
			fmt.Println("Worker", id, "received job", j)
			results <- j
		} else {
			fmt.Println("Worker", id, "received all jobs")
			wg.Done()
			return
		}
	}
}

func main() {
	wg.Add(5)

	for w := 1; w <= 5; w++ {
		go worker(w, jobs, results)
	}

	go producer()
	go writer()
	wg.Wait()
	close(results)
	<- done
}
