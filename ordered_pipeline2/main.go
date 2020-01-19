package main

import (
	"errors"
	"fmt"
	"time"
)

// purpose of ordered pipeline is to process things in order, even if the work takes a long period of time
// it will try to parallelize all the work given number of current workers
// if a task fails at one point, it should notify where it failed.

// imagine task 1,2,3,4,5,6,7,8,9 . and 2 workers. 1 and 2 would start, and then depends on which finishes, worker would

type Task struct {
	ID       int
	resultCh chan int
	errorCh  chan error
	canceled <-chan struct{}
}

func main() {
	c := gen(1000)
	// for each task create a new chanel for returned results
	finalResult := []int{}
	_ = finalResult
	tasks := make(chan Task, 1000)
	cancelCh := make(chan struct{})
	taskArray := []*Task{}
	for i := range c {
		t := Task{
			ID:       i,
			resultCh: make(chan int),
			errorCh:  make(chan error),
			canceled: cancelCh,
		}
		tasks <- t
		taskArray = append(taskArray, &t)
	}
	numOfWorkers := 100
	for i := 0; i < numOfWorkers; i++ {
		go doWork(i, tasks)
	}
	go func() {
		//time.Sleep(1100 * time.Millisecond)
		// propagates all cancelled signal to all go routines
		//close(cancelCh)
	}()
	defer func() {
		fmt.Println(finalResult)
	}()
	for _, i := range taskArray {
		select {
		case r := <-i.resultCh:
			finalResult = append(finalResult, r)
		case err := <-i.errorCh:
			fmt.Println(err)
			return
		case <-i.canceled:
			fmt.Println("task cancelled")
			return
		}
	}
}

func doWork(id int, t <-chan Task) {

	go func() {
		for n := range t {
			select {
			case <-n.canceled:
				fmt.Println("cancel got detected")
				return
			default:
				time.Sleep(1 * time.Second)
				if n.ID == 800 {
					n.errorCh <- errors.New("my 800 error")
					return
				}
				r := n.ID * n.ID
				n.resultCh <- r
				fmt.Println("worker", id, "finished on task", n.ID)
			}
		}
	}()
}

func gen(work int) <-chan int {
	out := make(chan int)
	go func() {
		for i := 0; i < work; i++ {
			out <- i
		}
		close(out)
	}()
	return out
}
