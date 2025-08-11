package backend

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"leaderboard/src/config"
	"leaderboard/src/metrics"
	"leaderboard/src/redisclient"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type LeaderboardUpdate struct {
	Data []byte
	Hash [32]byte
}

type Client struct {
	ID      int64
	channel chan LeaderboardUpdate
	ctx     context.Context
	cancel  context.CancelFunc
}

type LeaderboardBroadcaster struct {
	clients       map[int64]*Client
	clientsMutex  sync.RWMutex
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	clientCounter int64
}

func CreateLeaderboardBroadcaster() *LeaderboardBroadcaster {
	ctx, cancel := context.WithCancel(context.Background())

	lb := &LeaderboardBroadcaster{
		// make the channel buffered - clients may be slow, messages can pile up
		clients: make(map[int64]*Client),
		ctx:     ctx,
		cancel:  cancel,
	}

	lb.wg.Add(1)
	go lb.detectLeaderboardChanges()
	return lb
}

// should have an endpoint, with proper auth - admin only
func (lb *LeaderboardBroadcaster) StopBroadcast() {
	lb.cancel()
	lb.wg.Wait()

	lb.clientsMutex.Lock()
	defer lb.clientsMutex.Unlock()
	// remove all clients
	for _, client := range lb.clients {
		client.cancel()
		close(client.channel)
	}
}

// Create new channel for client, add to map
func (lb *LeaderboardBroadcaster) AddClient() (*Client, <-chan LeaderboardUpdate) {
	lb.clientsMutex.Lock()
	lb.clientCounter++
	ctx, cancel := context.WithCancel(lb.ctx)

	client := &Client{
		ID:      lb.clientCounter,
		ctx:     ctx,
		cancel:  cancel,
		channel: make(chan LeaderboardUpdate, config.AppConfig.Server.BroadcastBufferSize),
	}
	lb.clients[lb.clientCounter] = client
	lb.clientsMutex.Unlock()

	return client, client.channel
}

// remove specific client channel - closed connection
func (lb *LeaderboardBroadcaster) RemoveClient(client *Client) {
	lb.clientsMutex.Lock()
	defer lb.clientsMutex.Unlock()

	if _, exists := lb.clients[client.ID]; exists {
		delete(lb.clients, client.ID)
		client.cancel()
		close(client.channel)
	}
}

// broadcastToAllClients sends an update to all connected clients
func (lb *LeaderboardBroadcaster) broadcastToAllClients(update LeaderboardUpdate) {
	lb.clientsMutex.RLock()

	var clientsToRemove []*Client

	// what to do in case Client channel is full, skip this client (to be changed later - add channel clearing mechanism + alerting)
	for _, client := range lb.clients {
		select {
		case client.channel <- update:
			// sent
		case <-client.ctx.Done():
			// clean
			clientsToRemove = append(clientsToRemove, client)
		default:
			metrics.FilledSSEChannels.Inc()
		}
	}
	lb.clientsMutex.RUnlock()

	if len(clientsToRemove) > 0 {
		lb.clientsMutex.Lock()
		for _, client := range clientsToRemove {
			if _, exists := lb.clients[client.ID]; exists {
				delete(lb.clients, client.ID)
				client.cancel()
				close(client.channel)
			}
		}
		lb.clientsMutex.Unlock()
	}
}

// package level var
var broadcaster *LeaderboardBroadcaster

func SetBroadcaster(b *LeaderboardBroadcaster) {
	broadcaster = b
}

// poll redis, dedup leaderboard values, push to broadcast to all clients
func (lb *LeaderboardBroadcaster) detectLeaderboardChanges() {
	defer lb.wg.Done()
	ticker := time.NewTicker(time.Duration(config.AppConfig.Server.PollingIntervalSeconds) * time.Second)
	defer ticker.Stop()

	var lastHash [32]byte

	for {
		select {
		case <-ticker.C:
			results, err := redisclient.GetTopNPlayers(lb.ctx, "leaderboard", int64(config.AppConfig.Leaderboard.TopPlayersLimit))
			if err != nil {
				metrics.RedisOperationErrors.WithLabelValues("get_top_players").Inc()
				config.Error("Failed to fetch leaderboard from Redis.", map[string]any{"Error": err, "source": "/stream-leaderboard"})
				continue
			}

			resultString := fmt.Sprintf("%+v", results)
			currentHash := sha256.Sum256([]byte(resultString))

			if currentHash != lastHash {
				lastHash = currentHash

				jsonStart := time.Now()
				jsonData, err := json.Marshal(results)
				if err != nil {
					config.Error("JSON marshaling error",
						map[string]any{"Error": err, "source": "/stream-leaderboard", "results": results})
					metrics.JSONErrors.WithLabelValues("marshal").Inc()
					continue
				}
				metrics.JSONMarshalDuration.Observe(time.Since(jsonStart).Seconds())

				update := LeaderboardUpdate{
					Data: jsonData,
					Hash: currentHash,
				}

				lb.broadcastToAllClients(update)
			}
		case <-lb.ctx.Done():
			return
		}
	}
}

func StreamLeaderboard(c *gin.Context) {
	metrics.ConcurrentClients.Inc()
	defer metrics.ConcurrentClients.Dec()

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Flush()

	metrics.ActiveSSEConnections.Inc()
	config.Info("New SSE conn", map[string]any{"Num active clients": metrics.ActiveSSEConnections})
	defer metrics.ActiveSSEConnections.Dec()

	client, channel := broadcaster.AddClient()
	defer broadcaster.RemoveClient(client)

	for {
		select {
		case update, ok := <-channel:
			if !ok {
				// channel closed
				return
			}
			fmt.Fprintf(c.Writer, "data: %s\n\n", update.Data)
			c.Writer.Flush()
			metrics.SSEMessagesSent.Inc()
		case <-c.Request.Context().Done():
			metrics.DroppedSSEConnections.Inc()
			config.Info("Closed SSE conn", map[string]any{"open": metrics.ActiveSSEConnections})
			return
		}
	}
}
