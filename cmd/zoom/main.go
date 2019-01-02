package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"sync"
)

var method = flag.String("method", "", "the http method to use when making the request")
var n = flag.Int("n", 0, "the total number of requests you want to make")
var concurrency = flag.Int("concurrency", 1, "the number of requests you want to have in flight at any given time")
var url = flag.String("url", "", "the url that you want to send the request to")
var template = flag.String("template", "", "the template that will be used to produce the body of a request")

func main() {
	l := log.New(os.Stdout, "zoom: ", log.LstdFlags)
	flag.Parse()

	// We need to send the all of the requests into the channel so that our doers can execute them
	reqs := make(chan *http.Request)
	go func() {
		for i := 0; i < *n; i++ {
			req, err := http.NewRequest(*method, *url, nil)
			if err != nil {
				panic(err)
			}
			reqs <- req
		}
		close(reqs)
	}()

	// Create the specified amount of doers for the load test
	var wg sync.WaitGroup
	wg.Add(*concurrency)
	s := start()
	for i :=0; i < *concurrency; i++ {
		go func() {
			defer wg.Done()
			do(s, reqs)
		}()
	}

	// Wait for all of the doers to finish processing the requests
	wg.Wait()
	s.stop()

	// Display the results
	l.Println("---------- finished load test ----------")
	s.results(l)
}
