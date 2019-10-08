package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func spawnWorker(workerID string, port string) {
	cmd := exec.Command("bin/worker.linux", "-workerId", workerID, "-port", port)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd.Start() failed with %s: %s\n", err, stderr.String())
	}
	fmt.Printf("Spawned %s on :%s\n", workerID, port)
}

func callEndpoint(endpoint string, ch chan<- int) {
	var nums []int
	fmt.Println("Consumer started:", endpoint)
	res, err := http.Get(endpoint)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	p := make([]byte, 7) // rnd=99\n is 7 bytes
	for {
		n, err := res.Body.Read(p)
		if err != nil || err == io.EOF { // response ended or worker crashed
			break
		}
		rnd := string(p[:n])
		i := strings.Index(rnd, "\n")
		num, _ := strconv.Atoi(rnd[4:i])
		ch <- num
		nums = append(nums, num)
	}

	fmt.Printf("Consumer finished: %s, length:%d\n", endpoint, len(nums))
}

func worker(jobs <-chan int, results chan<- int) {
	for i := range jobs {
		workerID := strconv.Itoa(i)
		port := strconv.Itoa(3000 + i)
		spawnWorker(workerID, port)
		// Have to wait?
		time.Sleep(5 * time.Millisecond)
		endpoint := "http://localhost:" + port + "/rnd?n=15"
		go callEndpoint(endpoint, results)
	}
}

func main() {
	/*
		TODOs:
		What is Data sample? 150 data samples
	*/
	var nums []int
	num := 0
	start := time.Now()
	jobs := make(chan int, 30)
	results := make(chan int, 200)

	go worker(jobs, results)
	// go worker(jobs, results)

	for i := 1; i < 20; i++ {
		jobs <- i
	}
	close(jobs)

	for {
		num = <-results
		// fmt.Println(num, nums)
		nums = append(nums, num)
		if len(nums) >= 150 {
			break
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("Time elapsed: %v Total number: %d\n", elapsed, len(nums))
}
