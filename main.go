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
	"sync"
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

func callEndpoint(wg *sync.WaitGroup, endpoint string, ch chan int) {
	defer wg.Done()
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
		nums = append(nums, num)
	}
	ch <- len(nums)
	fmt.Printf("Consumer finished: %s, numbers:%v length:%d\n", endpoint, nums, len(nums))
}

func main() {
	/*
		TODOs:
		What is Data sample? 150 data samples
	*/
	start := time.Now()
	var wg sync.WaitGroup
	totalNums := 0
	numsChan := make(chan int, 17)

	for i := 1; i < 18; i++ {
		workerID := strconv.Itoa(i)
		port := strconv.Itoa(3000 + i)
		spawnWorker(workerID, port)
		// Have to wait?
		time.Sleep(5 * time.Millisecond)
		endpoint := "http://localhost:" + port + "/rnd?n=13"
		wg.Add(1)
		go callEndpoint(&wg, endpoint, numsChan)
	}
	wg.Wait()
	close(numsChan)
	for s := range numsChan {
		totalNums += s
	}
	elapsed := time.Since(start)
	fmt.Printf("Time elapsed: %v Total number: %d\n", elapsed, totalNums)
}
