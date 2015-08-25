package main

import (
	"runtime"

	"github.com/sparck/golang-poll/api"
)

func main() {
	// Max cpu
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Bootstrap
	api.Bootstrap()
}
