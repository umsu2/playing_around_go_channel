package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// make a simple server that serves on 8000, take some request , delay for a 100 ms and response but for every 10 ms will do some amount of work
// make a cancelable req on a client that calls that server  and cancels either with timeout or with actual cancellation
func main() {
	//fmt.Print(time.Now())
	//c := time.After(time.Second)
	//t := <-c
	//fmt.Print(t)
	http.HandleFunc("/test", handleSearch)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	c := make(chan error, 1)
	go func() {
		fmt.Println("blah")
		c <- work(w, r)
		fmt.Println("blahx2")
	}()
	select {
	case <-r.Context().Done():
		<-c
		fmt.Println(r.Context().Err())
		panic("context done?")
	case err := <-c:
		if err != nil {
			fmt.Println(err)
			panic("context not done returning from work")
		}
	}

}
func work(w http.ResponseWriter, r *http.Request) error {
	t0 := time.Now()
	c := 0
	var done bool
	go func() {
		<-r.Context().Done()
		fmt.Println("worker detected context done returning")
		done = true
		return
	}()
	for {
		c++
		fmt.Println("doing work...", c)
		time.Sleep(10 * time.Millisecond)
		if done {
			return nil
		}
		if time.Since(t0) >= time.Second {
			fmt.Println("finishing work")
			_ = json.NewEncoder(w).Encode(map[string]string{"success": "true"})
			return nil
		}
	}
}

func someGoRoutine(inputChannel chan int, errChan chan error) {

}
