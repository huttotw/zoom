package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func do(s *stats, reqs <-chan *http.Request, errs chan error) {
	c := http.Client{
		Transport: &http.Transport{
			MaxIdleConns: 0,
			MaxIdleConnsPerHost: 100,
			MaxConnsPerHost: 0,
		},
	}
	for req := range reqs {
		start := time.Now()
		resp, err := c.Do(req)
		if err != nil {
			errs <- err
			continue
		}
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
		s.incrementRequests(1)
		s.recordStatus(resp.StatusCode)
		s.recordLatency(time.Since(start))
	}
}
