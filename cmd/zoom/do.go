package main

import (
	"net/http"
	"time"
)

func do(s *stats, reqs <-chan *http.Request, errs chan error) {
	c := http.Client{}
	for req := range reqs {
		start := time.Now()
		resp, err := c.Do(req)
		if err != nil {
			errs <- err
			continue
		}
		resp.Body.Close()
		s.incrementRequests(1)
		s.recordStatus(resp.StatusCode)
		s.recordLatency(time.Since(start))
	}
}
