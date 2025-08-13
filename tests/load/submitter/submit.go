package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type GameResult struct {
	GameID    string `json:"game_id" binding:"required"`
	ServerID  string `json:"server_id" binding:"required"`
	Player1ID string `json:"player1_id" binding:"required"`
	Player2ID string `json:"player2_id" binding:"required"`
	Result    int    `json:"result" binding:"gte=0,lte=2"`
}

var (
	ip         = flag.String("ip", "127.0.0.1", "Server IP")
	port       = flag.Int("port", 8080, "Server Port")
	goroutines = flag.Int("goroutines", 10, "Number of concurrent submitters")
)

func main() {
	flag.Parse()
	url := fmt.Sprintf("http://%s:%d/submit-game", *ip, *port)

	for i := 0; i < *goroutines; i++ {
		go func(id int) {
			for {
				payload := GameResult{
					GameID:    fmt.Sprintf("game-%d-%d", id, rand.Intn(1000)),
					ServerID:  fmt.Sprintf("server-%d", rand.Intn(5)),
					Player1ID: fmt.Sprintf("p%d", rand.Intn(10)),
					Player2ID: fmt.Sprintf("p%d", rand.Intn(10)),
					Result:    rand.Intn(3), // 0, 1, or 2
				}

				data, err := json.Marshal(payload)
				if err != nil {
					log.Printf("[goroutine %d] marshal error: %v", id, err)
					continue
				}

				resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
				if err != nil {
					log.Printf("[goroutine %d] post error: %v", id, err)
					time.Sleep(1 * time.Second)
					continue
				}
				log.Printf("[goroutine %d] sent game %s vs %s (result=%d)",
					id, payload.Player1ID, payload.Player2ID, payload.Result)

				resp.Body.Close()
				time.Sleep(500 * time.Millisecond)
			}
		}(i)
	}
	select {}
}
