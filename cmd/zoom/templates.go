package main

import (
	"text/template"
	"math/rand"
	"time"

	"github.com/segmentio/ksuid"
)

var seed = rand.NewSource(time.Now().UnixNano())

var funcMap = template.FuncMap{
	"ksuid": ksuid.New,
	"intn": intn,
}

func intn (max int) int {
	return rand.Intn(max)
}