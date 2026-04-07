2026-04-07
- Add no metrics mode using a config feature flag (don't register metrics, don't expose endpoint, no containers for prom and grafana)

2025-12-27
- Update README to latest state
- Add a github actions workflow to run tests at each pull request to main
- Check for nil configs at the broadcaster constructor, return error if nil

2025-12-26
- Add functions to enable testing in `sse.go`
- Add basic unit tests for `AddClient` and `RemoveClient`
- Minor improvements to other parts of `sse.go` and tests

2025-12-25
- Remove global config from inside handler, use injected values
- Add tests for `leaderboard.go` and `player.go`
- Remove use of global config in `sse.go`

2025-12-24
- Create `interfaces` package for easier dep injection, major changes to `redisclient`, make a constructor returning a client instance instead of using global.
- Subsequent changes for interfaces and redis in handler and test code

2025-12-23
- Start dependency injection changes for redis
- Add test for `/submit-game` endpoint

2025-12-21
- Build table driven test for game result validation
- Refactor GameResult.Validate() to return error values, add a new validation condition
- Refactor `/submit-game` endpoint to use dependency injection (Redis), turn into Gin HandlerFunc, add GameResult validation check before Redis update
- Start documentation of testing process