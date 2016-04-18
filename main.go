package main

import (
	"runtime"

	"github.com/scristofari/golang-poll/api"
)

func main() {
	// Max goroutine
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Bootstrap
	api.Bootstrap()
}
