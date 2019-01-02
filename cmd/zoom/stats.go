package main

import (
	"log"
	"net/http"
	"sync"
	"time"
	"sort"
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

	switch status {
	case
		http.StatusContinue,
		http.StatusSwitchingProtocols,
		http.StatusProcessing:
			s.statuses[100]++
	case
		http.StatusOK,
		http.StatusCreated,
		http.StatusAccepted,
		http.StatusNonAuthoritativeInfo,
		http.StatusNoContent,
		http.StatusResetContent,
		http.StatusPartialContent,
		http.StatusMultiStatus,
		http.StatusAlreadyReported,
		http.StatusIMUsed:
			s.statuses[200]++
	case
		http.StatusMultipleChoices,
		http.StatusMovedPermanently,
		http.StatusFound,
		http.StatusSeeOther,
		http.StatusNotModified,
		http.StatusUseProxy ,
		http.StatusTemporaryRedirect,
		http.StatusPermanentRedirect:
			s.statuses[300]++
	case
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusPaymentRequired,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusMethodNotAllowed,
		http.StatusNotAcceptable,
		http.StatusProxyAuthRequired,
		http.StatusRequestTimeout,
		http.StatusConflict,
		http.StatusGone,
		http.StatusLengthRequired,
		http.StatusPreconditionFailed,
		http.StatusRequestEntityTooLarge,
		http.StatusRequestURITooLong,
		http.StatusUnsupportedMediaType,
		http.StatusRequestedRangeNotSatisfiable,
		http.StatusExpectationFailed,
		http.StatusTeapot,
		http.StatusMisdirectedRequest,
		http.StatusUnprocessableEntity,
		http.StatusLocked,
		http.StatusFailedDependency,
		http.StatusUpgradeRequired,
		http.StatusPreconditionRequired,
		http.StatusTooManyRequests,
		http.StatusRequestHeaderFieldsTooLarge,
		http.StatusUnavailableForLegalReasons:
			s.statuses[400]++
	case
		http.StatusInternalServerError,
		http.StatusNotImplemented,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
		http.StatusHTTPVersionNotSupported,
		http.StatusVariantAlsoNegotiates,
		http.StatusInsufficientStorage,
		http.StatusLoopDetected,
		http.StatusNotExtended,
		http.StatusNetworkAuthenticationRequired:
			s.statuses[500]++
	}
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

	sort.Sort(s.latencies)
	dur := s.end.Sub(s.begin).String()
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
	l.Printf("max latency:\t\t%s", max)
	l.Printf("min latency:\t\t%s", min)
	l.Printf("median latency:\t%s", med)
	l.Printf("p99 latency:\t\t%s", p99)
	l.Printf("p95 latency:\t\t%s", p95)
	l.Printf("p90 latency:\t\t%s", p90)
	l.Printf("1xx:\t\t\t%d", s.statuses[100])
	l.Printf("2xx:\t\t\t%d", s.statuses[200])
	l.Printf("3xx:\t\t\t%d", s.statuses[300])
	l.Printf("4xx:\t\t\t%d", s.statuses[400])
	l.Printf("5xx:\t\t\t%d", s.statuses[500])
	l.Printf("error rate:\t\t%f", errorRate)
}