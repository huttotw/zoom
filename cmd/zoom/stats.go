package main

import (
	"log"
	"sort"
	"sync"
	"time"
)


type stats struct {
	mutex sync.Mutex
	statuses map[int]int
	latencies latencies
	end time.Time
	begin time.Time
	requests int
}

func start() *stats {
	s := stats{
		statuses: make(map[int]int),
		latencies: make([]time.Duration, 0),
		requests: 0,
		begin: time.Now(),
	}
	return &s
}

func (s *stats) incrementRequests(n int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.requests++
}

func (s *stats) stop() {
	s.end = time.Now()
}

func (s *stats) recordLatency(dur time.Duration) {
	s.latencies = append(s.latencies, dur)
}

func (s *stats) recordStatus(status int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.statuses[status]++
}

type latencies []time.Duration

func (l latencies) Len() int {
	return len(l)
}

func (l latencies) Less(i, j int) bool {
	return l[i] < l[j]
}

func (l latencies) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
	}



func (s *stats) results(l *log.Logger) {
	if s.requests == 0 {
		l.Println("no results")
		return
	}

	end := s.end
	if s.end.IsZero() {
		end = time.Now()
	}

	sort.Sort(s.latencies)
	dur := end.Sub(s.begin).String()
	reqPerS := float64(s.requests) / end.Sub(s.begin).Seconds()
	max := s.latencies[len(s.latencies) - 1].String()
	min := s.latencies[0].String()
	medIndex := len(s.latencies) / 2
	med := s.latencies[medIndex]
	p99Index := 0.99 * float64(len(s.latencies))
	p99 := s.latencies[int(p99Index - 1)]
	p95Index := 0.95 * float64(len(s.latencies))
	p95 := s.latencies[int(p95Index - 1)]
	p90Index := 0.90 * float64(len(s.latencies))
	p90 := s.latencies[int(p90Index - 1)]
	errorRate := float64(s.statuses[400] + s.statuses[500]) / float64(s.requests)

	l.Printf("total requests:\t%d", s.requests)
	l.Printf("duration:\t\t%s", dur)
	l.Printf("reqs/s:\t\t%f", reqPerS)
	l.Printf("max latency:\t\t%s", max)
	l.Printf("min latency:\t\t%s", min)
	l.Printf("median latency:\t%s", med)
	l.Printf("p99 latency:\t\t%s", p99)
	l.Printf("p95 latency:\t\t%s", p95)
	l.Printf("p90 latency:\t\t%s", p90)
	for k, v := range s.statuses {
		l.Printf("status %d:\t\t%d", k, v)
	}
	l.Printf("error rate:\t\t%f", errorRate)
}