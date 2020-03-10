/*
This snippet is an example of backpressure implementation in Go.
It doesn't run in Go Playground, because it starts an HTTP Server.
The example starts an HTTP server and sends multiple requests to it. The server starts denying
requests by replying an "X" (i.e. a 502) when its buffered channel reaches capacity.
This is not the same as rate-limiting; you might be interested in https://github.com/juju/ratelimit
or https://godoc.org/golang.org/x/time/rate.
Note that asking the question: `len(ch) < cap(ch) ?` is a racey operation; the channel might
actually be full a microsecond later. In the context of short running requests this is not an
issue.
Outputs:
```
$ go run test.go
√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√
√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√X√X
√√X√√X√√X√√X√√X√√X√√X√√X√X√√√X√√X√X√√√X√X√√X√√X√√X√√X√√X√√X√√X√X√√X√√X√√X√√X√X√√X√√X√√X√√X√√X
√X√√√X√X√√X√√X√√X√√X√X√√X√X√√√X√X√√X√√X√√X√√X√√X√X√√X√√X√√X√√X√X√√X√√X√X√√X√X√√X√√X√√X√√X√√X√
√X√√X√X√√X√√X√√X√√X√√X√X√√X√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√
√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√√
```
*/
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func main() {
	requests := make(chan request, 100)

	go startServer(requests)
	go process(requests)

	makeRequests(500, 6*time.Millisecond)
}

func startServer(rq chan request) {
	http.HandleFunc("/requests", handle(rq))
	http.ListenAndServe(":9000", nil)
}

func process(rq chan request) {
	for r := range rq {
		r.process()
	}
}

func makeRequests(count int, cooldown time.Duration) {
	wg := sync.WaitGroup{}
	for i := 0; i < count; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			response, _ := http.Get("http://localhost:9000/requests")
			defer response.Body.Close()
			b, _ := ioutil.ReadAll(response.Body)
			fmt.Print(string(b))
		}()
		time.Sleep(cooldown)
	}
	wg.Wait()
}

func handle(rq chan request) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(rq) < cap(rq) {
			r := newRequest(r)
			rq <- r
			w.Write(<-r.response)
		} else {
			w.Write([]byte("X"))
		}
	}
}

type request struct {
	r        *http.Request
	response chan []byte
}

func newRequest(r *http.Request) request { return request{r, make(chan []byte)} }

func (r request) process() {
	time.Sleep(10 * time.Millisecond)
	r.response <- []byte("√")
}
