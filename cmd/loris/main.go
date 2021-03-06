package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/jbert/loris"
	"github.com/jbert/loris/store"
)

type opt struct {
	port      uint
	storeName string
	debug     bool

	cpuProfile string
	memProfile string
}

func (o *opt) validate() error {
	err := store.ValidateName(o.storeName)
	if err != nil {
		return fmt.Errorf("Invalid store name: %s", err)
	}
	return nil
}

func getOpts() *opt {
	o := opt{}
	flag.UintVar(&o.port, "port", 6378, "Listening port")
	flag.StringVar(&o.storeName, "store", "mutexmap", "Store name")
	flag.BoolVar(&o.debug, "debug", false, "Enable debug output")
	flag.StringVar(&o.cpuProfile, "cpuprofile", "", "Write CPU profile")
	flag.StringVar(&o.memProfile, "memprofile", "", "Write Memory profile")
	flag.Parse()

	err := o.validate()
	if err != nil {
		log.Printf("Error: %s", err)
		flag.PrintDefaults()
		os.Exit(-1)
	}
	return &o
}

func main() {
	o := getOpts()
	store, err := store.NewFromName(o.storeName)
	if err != nil {
		log.Fatalf("Error constructing store: %s", err)
	}
	log.Printf("Using store %s", o.storeName)
	gr := loris.NewWithStore(store, o.debug)
	if o.cpuProfile != "" {
		f, err := os.Create(o.cpuProfile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
	err = gr.ListenAndServe(fmt.Sprintf(":%d", o.port))
	log.Printf("Server exited: %s", err)
	if o.memProfile != "" {
		f, err := os.Create(o.memProfile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}
}
