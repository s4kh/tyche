package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strconv"
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
	fmt.Println("Anything from cmd:", out.String())
	fmt.Printf("Spawned %s on :%s\n", workerID, port)
}

func callEndpoint(wg *sync.WaitGroup, endpoint string) {
	defer wg.Done()

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
		fmt.Printf("%s", string(p[:n]))
	}
	fmt.Printf("Consumer finished: %s\n", endpoint)
}

func main() {
	/*
		TODOs:
		1. Consume numbers count them using new line
		2. Spawn multiple workers with different id and port
		 a. How to know which workers are ready to accept request
		3. Send requests to workers - multiple clients also
	*/
	start := time.Now()
	var wg sync.WaitGroup

	for i := 1; i < 17; i++ {
		workerID := strconv.Itoa(i)
		port := strconv.Itoa(3000 + i)
		spawnWorker(workerID, port)
		// Have to wait?
		time.Sleep(5 * time.Millisecond)
		endpoint := "http://localhost:" + port + "/rnd?n=2"
		wg.Add(1)
		go callEndpoint(&wg, endpoint)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Println("Time elapsed:", elapsed)
}
