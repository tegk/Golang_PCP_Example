package main

import (
	"fmt"
	"sync"
	"time"
	"math/rand"
)

//creating buffered channels with the channel "done"
//used to block terminating main()

var wg sync.WaitGroup

type Job struct {
	Work   string
}

//sending produced to channel "jobs" and after closes the channel.
func producer(jobs chan<- *Job, ){
	for c := 'a'; c <= 'z'; c++ {
		jobs <- &Job{Work: fmt.Sprintf("%c", c)}
	}
	close(jobs)
}

func writer(results <-chan *Job, done chan<- bool){
	for {
		j, more :=  <-results
		if more {
			fmt.Println(j.Work)
		} else {
			done <- true
			return
		}
	}
}

func worker(jobs <-chan *Job, results chan<- *Job) {
	for {
		j, more := <-jobs
		if more {
			time.Sleep(time.Duration(rand.Float32() * float32(time.Second)))
			results <- j
		} else {
			wg.Done()
			return
		}
	}
}

func main() {
	var jobs = make(chan *Job)
	var results = make(chan *Job, 100)
	var done = make(chan bool, 1)

	for w := 1; w <= 4; w++ {
		wg.Add(1)
		go worker(jobs, results)
	}

	go producer(jobs)
	go writer(results, done)
	wg.Wait()
	close(results)
	<- done
}
