2025-12-21
- Build table driven test for game result validation
- Refactor GameResult.Validate() to return error values, add a new validation condition
- Refactor `/submit-game` endpoint to use dependency injection (Redis), turn into Gin HandlerFunc, add GameResult validation check before Redis update
- 