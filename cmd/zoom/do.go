package main

import (
	"net/http"
	"time"
)

func do(s *stats, reqs <-chan *http.Request) error {
	for req := range reqs {
		start := time.Now()
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		s.incrementRequests(1)
		s.recordStatus(resp.StatusCode)
		s.recordLatency(time.Since(start))
	}

	return nil
}
