package main

import (
	"fmt"
	"log"

	"github.com/jbert/goredis"
)

type opt

func main() {
	//	defaultPort := 6379
	defaultPort := 6378
	gr := goredis.New()
	err := gr.ListenAndServe(fmt.Sprintf(":%d", defaultPort))
	log.Fatalf("Server exited: %s", err)
}
