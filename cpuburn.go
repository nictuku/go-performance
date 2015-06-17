// cpuburn is an HTTP server that runs a busy loop that sends requests to
// itself and records statistics. It's used to analyze the performance of the
// Go scheduler, given different parallelization of workload and system
// threads.
//
// The main effect I'm trying to understand is how the scheduler penalizes
// workloads with less parallelization.
package main

import (
	"expvar"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/nictuku/cpustat"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Got %v", r.URL.Path)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: cpuburn <number of concurrent GET goroutines>")
	}
	numRoutines, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Expected 'number of concurrent GET goroutines', got %q: %v", os.Args[1], err)
	}
	fmt.Println("Go version:", runtime.Version())

	http.HandleFunc("/", handler)
	go func() {
		err := http.ListenAndServe("localhost:8080", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	count := expvar.NewInt("total-requests")
	cpu := cpustat.New()
	var maxPerformance float64

	// Print the latest QPS, process CPU usage and QPS/core.
	go func() {
		t := time.Now()
		for {
			count.Set(0)
			time.Sleep(time.Second)
			requests, _ := strconv.ParseFloat(count.String(), 32)
			cpuUsage, err := cpu.ProcCPU()
			if err != nil {
				log.Fatal("cpu ProcCPU: %v", err)
			}
			qps := requests / time.Since(t).Seconds()
			qpsCore := qps / cpuUsage
			log.Println(numRoutines, runtime.GOMAXPROCS(0), qps, cpuUsage, qpsCore)

			if qpsCore > maxPerformance {
				maxPerformance = qpsCore
			}
			t = time.Now()
		}

	}()

	// Send GET requests to the local server as fast as possible.
	// Exit if any error occurs.
	for i := 0; i < numRoutines; i++ {
		go func() {
			for {
				resp, err := http.Get("http://localhost:8080")
				if err != nil {
					log.Fatal(err)
				}
				defer resp.Body.Close()
				// Read the body so this connection can be reused if necessary.
				_, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err)
				}
				count.Add(1)
			}

		}()
	}
	time.Sleep(10 * time.Second)
	fmt.Printf("%d %d %f\n", numRoutines, runtime.GOMAXPROCS(0), maxPerformance)
	os.Exit(0)
}
