package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
	"time"
)

var method = flag.String("method", "", "the http method to use when making the request")
var n = flag.Int("n", 0, "the total number of requests you want to make")
var concurrency = flag.Int("concurrency", 1, "the number of requests you want to have in flight at any given time")
var url = flag.String("url", "", "the url that you want to send the request to")
var temp = flag.String("template", "", "the template that will be used to produce the body of a request")
var tempFile = flag.String("template-file", "", "the template file that will be used to produce the body of a request")

func main() {
	l := log.New(os.Stdout, "zoom: ", log.LstdFlags)
	flag.Parse()

	// Parse the incoming template so that we can execute before each request
	var tmpl *template.Template
	var err error
	if *temp != "" {
		tmpl, err = template.New("main").Funcs(funcMap).Parse(*temp)
		if err != nil {
			panic(err)
		}
	} else {
		_, file := filepath.Split(*tempFile)
		tmpl, err = template.New(file).Funcs(funcMap).ParseFiles(*tempFile)
		if err != nil {
			panic(err)
		}
	}


	// We need to send the all of the requests into the channel so that our doers can execute them
	reqs := make(chan *http.Request)
	go func() {
		i := 0
		for i < *n || *n == 0 {
			buf := bytes.NewBuffer(nil)
			err := tmpl.Execute(buf, nil)
			if err != nil {
				panic(err)
			}

			req, err := http.NewRequest(*method, *url, buf)
			if err != nil {
				panic(err)
			}
			req.Header.Set("User-Agent", "zoom")
			reqs <- req
			i++
		}
		close(reqs)
	}()

	s := start()
	l.Println("---------- starting load test ----------")
	s.results(l)

	// Show the results every 15 seconds
	tick := time.Tick(time.Second * 15)
	go func() {
		for {
			select {
			case <-tick:
				l.Println("---------- running load test ----------")
				s.results(l)
			}
		}
	}()

	// Create the specified amount of doers for the load test
	var wg sync.WaitGroup
	wg.Add(*concurrency)
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
