package main

import (
	"fmt"
	"math/rand"
	"text/template"
	"time"

	"github.com/segmentio/ksuid"
)

var seed = rand.NewSource(time.Now().UnixNano())

var funcMap = template.FuncMap{
	"email": email,
	"enum": enum,
	"intn": intn,
	"ip": ip,
	"ksuid": ksuid.New,
	"string": ksuid.New,
	"url": randURL,
	"time": time.Now,
}

func email() string {
	return fmt.Sprintf("%s@%s.com", ksuid.New().String(), ksuid.New().String())
}

func enum (args ...interface{}) interface{} {
	i := rand.Intn(len(args))
	return args[i]
}

func intn (max int) int {
	return rand.Intn(max)
}

func ip() string {
	return fmt.Sprintf("%d.%d.%d.%d", intn(255), intn(255), intn(255), intn(255))
}

func randURL() string {
	return fmt.Sprintf("https://%s.com/%s", ksuid.New().String(), ksuid.New().String())
}

