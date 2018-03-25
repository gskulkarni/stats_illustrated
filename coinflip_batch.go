package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
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

	const flips int = 10000000
	const batch int = 100000
	cltCalc := true
	fmt.Printf("Cores used: %d\n", runtime.NumCPU())
	if cltCalc {
		const binWidth float64 = 0.00000005
		const left float64 = 0.0
		const right float64 = 1.0
		numBins := ((right - left) / binWidth)
		frequency := make([]int, int(numBins))
		fmt.Printf("Number of bins = %v, bin width = %v", numBins, binWidth)

		const samples int = 10000
		for i := 0; i < samples; i++ {
			total, _ := flipCustomRand(flips, batch)
			result := float64(total) / float64(flips)
			binIdx := int(result / binWidth)
			frequency[binIdx]++
		}

		f, _ := os.Create("frequency10000000.txt")
		defer f.Close()

		writer := bufio.NewWriter(f)
		for i := 0; i < int(numBins); i++ {
			fmt.Fprintf(writer, "%0.6f, %d\n", float64(i)*binWidth, frequency[i])
		}
		writer.Flush()
	} else {
		total, timeTaken := flipCustomRandWg(flips, batch)
		fmt.Printf("Final outcome: %d out of %d\n", total, flips)
		fmt.Printf("Time taken: %v", timeTaken)
	}
}

func flipCustomRand(flips, batch int) (int, time.Duration) {
	c := make(chan int)
	total := 0
	start := time.Now()
	for i := 0; i < flips; i += batch {
		go func(k int) {
			r := rand.New(rand.NewSource(start.UnixNano() + int64(k)))
			result := 0
			for j := 0; j < batch; j++ {
				p := r.Float64()
				if p >= 0.5 {
					result++
				}
			}
			c <- result
		}(i)
	}
	count := 0
	for i := 0; i < flips/batch; i++ {
		result := <-c
		count += result
	}
	total += count
	end := time.Now()
	diff := end.Sub(start)
	return total, diff
}

func flipCustomRandWg(flips, batch int) (int, time.Duration) {
	total := 0
	result := make([]int, flips/batch)
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < flips; i += batch {
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			r := rand.New(rand.NewSource(start.UnixNano() + int64(k)))
			for j := 0; j < batch; j++ {
				p := r.Float64()
				if p >= 0.5 {
					result[k/batch]++
				}
			}
		}(i)
	}
	wg.Wait()
	for i := 0; i < flips/batch; i++ {
		total += result[i/batch]
	}
	end := time.Now()
	diff := end.Sub(start)
	return total, diff
}
