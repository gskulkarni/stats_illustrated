package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	const iter int = 1000000000
	total := 0
	start := time.Now()
	var result int
	r := rand.New(rand.NewSource(start.UnixNano()))
	for it := 0; it < iter; it++ {
		p := r.Float64()
		if p >= 0.5 {
			result = 1
		} else {
			result = 0
		}
		total += result
	}
	end := time.Now()
	diff := end.Sub(start)
	fmt.Printf("Final outcome: %d out of %d\n", total, iter)
	fmt.Printf("Time taken: %v", diff)
}
