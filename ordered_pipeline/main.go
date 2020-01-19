package main

import (
	"fmt"
	"sync"
)

// purpose of ordered pipeline is to process things in order, even if the work takes a long period of time
// it will try to parallelize all the work given number of current workers
// if a task fails at one point, it should notify where it failed.

// imagine task 1,2,3,4,5,6,7,8,9 . and 2 workers. 1 and 2 would start, and then depends on which finishes, worker would

func main() {
	c := gen(1000)
	numOfWorkers := 10
	m := &sync.Map{}
	var wg sync.WaitGroup
	wg.Add(numOfWorkers)
	for i := 0; i < numOfWorkers; i++ {
		go doWork(c, m, i, &wg)
	}
	wg.Wait()
	m.Range(func(k interface{}, v interface{}) bool {
		fmt.Println(k)
		fmt.Println(v)
		return true
	})
}

func doWork(in <-chan int, m *sync.Map, id int, wg *sync.WaitGroup) {
	go func() {
		for n := range in {
			fmt.Println("worker", id, "doing", n)
			m.Store(n, n*n)
		}
		fmt.Println("worker", id, "done")
		wg.Done()
		return
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
