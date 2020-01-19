package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Millisecond)
	go func() {
		time.Sleep(time.Millisecond * 25)
		cancel()
	}()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8000/test", nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))
}
