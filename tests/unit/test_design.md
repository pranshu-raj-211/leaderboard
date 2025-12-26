## Things to test
### Handler behavior
Don't go for the redis operation or JSON binding. Just test what the inner code in the func is supposed to do. Mock any redis functions if used.

### sse.go
#### Client lifecycle
1. Adding clients - once added, counters go up, channels and other stuff are created
2. Removing clients - remove actually occurs, is idempotent, does not panic
3. Closed channels are never written to

#### Broadcaster
1. Test broadcast - same update received by all active clients
2. No broadcast to happen if lb state unchanged
3. Broadcasts do not block indefinitely
4. Backpressure (draining channels when full, broadcast of state happens eventually)

#### Handler
1. Calling handler creates clients
2. At end of handler (client side closes conn) the client is removed
3. Heartbeats sent
4. Closing connection closes client, removes channel

#### Concurrency
1. Channels are closed only once (test concurrency bug fixed using sync.Once)