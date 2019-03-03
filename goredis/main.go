package main

import (
	"fmt"
	"log"

	"github.com/jbert/goredis"
)

func main() {
	defaultPort := 6379
	gr := goredis.New()
	err := gr.ListenAndServe(fmt.Sprintf(":%d", defaultPort))
	log.Fatalf("Server exited: %s", err)
}
