package main

import (
	"bufio"
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
	const iter int = 10000
	const batch int = 10000
	cltCalc := false
	if cltCalc {
		const binWidth float64 = 0.00000005
		const left float64 = 0.0
		const right float64 = 1.0
		numBins := ((right - left) / binWidth)
		frequency := make([]int, int(numBins))
		fmt.Printf("Number of bins = %v, bin width = %v", numBins, binWidth)

		const samples int = 10000
		for i := 0; i < samples; i++ {
			total, _ := flipCustomRand(iter, flips, batch)
			result := float64(total) / float64(iter*flips)
			binIdx := int(result / binWidth)
			frequency[binIdx]++
		}

		f, _ := os.Create("frequency1.txt")
		defer f.Close()

		writer := bufio.NewWriter(f)
		for i := 0; i < int(numBins); i++ {
			fmt.Fprintf(writer, "%0.6f, %d\n", float64(i)*binWidth, frequency[i])
		}
		writer.Flush()
	} else {
		total, timeTaken := flipCustomRand(iter, flips, batch)
		fmt.Printf("Final outcome: %d out of %d\n", total, iter*flips)
		fmt.Printf("Time taken: %v", timeTaken)
	}
}

func flip(iter, flips, batch int) (int, time.Duration) {
	c := make(chan int)
	total := 0
	start := time.Now()
	for it := 0; it < iter; it++ {
		rand.Seed(start.UnixNano() + int64(it))
		for i := 0; i < flips; i++ {
			if i%batch == 0 {
				go func() {
					result := 0
					for j := 0; j < batch; j++ {
						p := rand.Float64()
						if p >= 0.5 {
							result++
						}
					}
					c <- result
				}()
			}
		}
		count := 0
		for i := 0; i < flips/batch; i++ {
			result := <-c
			count += result
		}
		total += count
	}
	end := time.Now()
	diff := end.Sub(start)
	return total, diff
}

func flipCustomRand(iter, flips, batch int) (int, time.Duration) {
	c := make(chan int)
	total := 0
	start := time.Now()
	for it := 0; it < iter; it++ {
		for i := 0; i < flips; i++ {
			if i%batch == 0 {
				go func() {
					r := rand.New(rand.NewSource(start.UnixNano() + int64(it+i)))
					result := 0
					for j := 0; j < batch; j++ {
						p := r.Float64()
						if p >= 0.5 {
							result++
						}
					}
					c <- result
				}()
			}
		}
		count := 0
		for i := 0; i < flips/batch; i++ {
			result := <-c
			count += result
		}
		total += count
	}
	end := time.Now()
	diff := end.Sub(start)
	return total, diff
}
