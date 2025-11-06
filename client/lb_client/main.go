package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	ip          = flag.String("ip", "127.0.0.1", "Server IP")
	connections = flag.Int("conn", 10000, "Number of SSE connections")
)

func main() {
	flag.Usage = func() {
		io.WriteString(os.Stderr, `SSE client generator
Example usage: ./client -ip=127.0.0.1 -conn=10000
`)
		flag.PrintDefaults()
	}
	flag.Parse()

	url := fmt.Sprintf("http://%s:8080/stream-leaderboard", *ip)
	var conns []*http.Response

	for i := 0; i < *connections; i++ {
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("conn %d failed: %v", i, err)
			break
		}
		defer resp.Body.Close()
		conns = append(conns, resp)
		log.Printf("Connection %d established", i)
	}

	// keep alive to check persistence
	for {
		time.Sleep(30 * time.Second)
	}
}
