# Refactoring code for dependency injection and decoupling
At the time of writing tests, the code was a bit too coupled. There's use of globals everywhere, functions would perform operations in a very restrictive manner using these, which did not allow for testing. 

I refactored the code with the goal of allowing dependency injection and adding a bit more modularity. This started with the Redis client (which is the main dependency) and the changes propagated further as needed.

The changes are added through this [pull request](https://github.com/pranshu-raj-211/leaderboard/pull/18).

## How were the changes designed
The main issue was the global redisclient usage, which was used by all the handlers and the broadcaster. Injecting the dependency instead would allow for mocking, which would ease testing (as you don't have to spin up redis, and keep test code separate).

Changes actually began with moving the interfaces defined in `submit_game.go` and `main.go` into a separate interfaces package. A struct was added to replace the redis.Z return type by the `GetTopNPlayers` function.

All redis functions were changed to methods that used a `LeaderboardStore` pointer, which would have contained a redis client. The  `InitRedis` function, which previously operated on the global redis client instance was replaced by `NewRedisClient`, which instead took in the config and logger(to warn in case of redis connection failures during retries) and returned a Redis client based on the config provided. This is much better design than using the defaults, not hardcoding anything. The `return config.Error` in case of error during a retry was replaced with a logger.warn, and returns now contain an error field as well.

A constructor was created to assign the client to the struct instance (which can be bypassed, but there were some errors while trying to do that so I kept this).

## Handler updates
Two function handlers in the `player.go` and `leaderboard.go` files were changed to now return closures and accept the store as a dependency, a pattern that I started using in `submit_game.go`. This enables testing - inject the store at time of assigning routes to the router. Instead of taking in struct instances as inputs any struct which implements the `LeaderboardStore` interface is accepted.

There were a few minor changes - using explicit http status signals instead of harcoded numbers and the use of `ctx.JSON` to marshal responses.

Respective changes were made to main. The redis client functions they call do not have a key param anymore, so that was also fixed. (I wonder if redis functions should be moved to a separate file, just keep constructors in that file)?

---

Apart from these, `sse.go` and `main.go` got some updates to make things compile (and some extra checks in main). Test updates were done to reflect these changes in the logic, and further tests for the leaderboard and player files have yet to be written (in part due to design of the function mocking the actual redis function calls not being something I can come up with at the moment - which I need to work on).

Future changes (after the test additions), will be on:
1. Refactoring `sse.go`
2. Refactoring `config.go`
3. Removing `return config.Log` wherever seen, replace with returning of actual errors
4. Adding documentation and docstrings to the code
5. Evaluating what metrics and code are unused and can be removed

But first, this needs to be in a stable state, with no bugs before other experiments (profiling, TCP tuning) can be done.