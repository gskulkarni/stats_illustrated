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

	const flips int = 100000
	const iter int = 100000
	c := make(chan int)
	total := 0
	start := time.Now()
	for it := 0; it < iter; it++ {
		rand.Seed(start.UnixNano() + int64(it))
		for i := 0; i < flips; i++ {
			go func() {
				p := rand.Float64()
				if p >= 0.5 {
					c <- 1
				} else {
					c <- 0
				}
			}()
		}
		count := 0
		streak := 0
		for i := 0; i < flips; i++ {
			result := <-c
			if result > 0 {
				streak++
			} else {
				if streak >= 30 {
					fmt.Printf("Found streak with length %d\n", streak)
				}
				streak = 0
			}
			count += result
		}
		total += count
	}
	end := time.Now()
	diff := end.Sub(start)
	fmt.Printf("Final outcome: %d out of %d\n", total, iter*flips)
	fmt.Printf("Time taken: %v", diff)
}
