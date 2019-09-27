package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	// ./worker -workerId 1 -port 3001
	// go func() {
	start := time.Now()
	cmd := exec.Command("bin/worker.linux", "-workerId", "1", "-port", "3001")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd.Start() failed with %s: %s\n", err, stderr.String())
	}
	// }()

	// Have to wait?
	time.Sleep(1 * time.Second)
	// pr, pw := io.Pipe()
	//`http://localhost:3001/rnd?n=100`
	res, err := http.Get("http://localhost:3001/rnd?n=40")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	io.Copy(os.Stdout, res.Body)
	elapsed := time.Since(start)
	fmt.Println(elapsed)
}
