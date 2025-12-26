# Process of adding tests, improving code, and learning from mistakes

## Testing models

### Validating GameResult (`models.go`)
Smallest part of the code, independent from other parts (has downstream dependencies), perfect starting point for tests.

Decided on a basic table driven test (similar to the one I wrote in dbos), to check for an error and report when it's different to expectations (got error when none expected, or got no error when one should have been returned).

This basic test uncovered one issue - Validate did not check for `ServerID`, so not passing that in worked, Go initialized that to an empty string. This is undesirable behavior and was fixed.

This simple error check when expected fell apart soon when testcases with multiple errors were tried out (based on previous case - kept ServerID empty when result was too high). I realized that the tests would need some way of identifying which type of thing had failed (different from just adding a message to the error).

To fix this, I created a `ValidationError` which had a `Field` and `Reason` keys. The `Field` showed what part of the GameResult object did not have something as expected. The Reason was a message based on the type of check applied during validation. This provides more clarity, and is more deliberate about what kind of errors are being found on a more granular level, which I believe is good design for this code.

#### Learnings
- `errors.As` takes in an error, a target, and if the error is of the type target it assigns the value of the error to target. This can later be used as I did for the last test check involving the Field.
- Typed errors - The `ValidationError` that I created is a pattern in Go named a type error. It's elegant and reminds me of Python. Lots of people do not like the implicit implementation of interfaces in Go, but I believe it has a lot of uses, and particularly like the uses it has in testing (much better to not implement interfaces and inherit classes explicitly in this case - for mocking or dependency injection as an example).
- zap.NewNop is a no op logger (doesn't actually log anything or write internal errors, perfect for testing code that has logging dependencies).

## Testing HTTP handlers

### Game submissions (`submit_game.go`)
This is related to the `/submit-game` endpoint, required for game servers to submit game results, which are used to update scores in the leaderboard. This has quite a few dependencies, but significantly less than other endpoints, and the code is easy to test as well (few binding, validation and Redis update).

Designing a test for this made me realize how tightly everything is coupled in my code. I did make a basic version of dependency injection for the broadcaster in the SSE streaming code, but Redis and other dependencies are still as tightly coupled as ever. This is a problem, since I'd have to use those as is for testing, and while I agree things can be done without mocking, mocks will make things simpler, and if used correctly, make the code easier to reason about (provided I do not go overboard with things).

DI is something I really looked forward to using in FastAPI, it made my life a whole lot easier, and I look forward to understanding how it's done in production systems running Go.

Here's how the implementation went:
- Added an interface that implements `UpdateLeaderboard`, which allows mocking the Redis client. Updated methods using the Redisclient (so far only the game submissions handler) to accept instances of this interface instead of an actual Redis client, so that it can be swapped out for test specific mocks instead.
- Add a game result validation step to the `SubmitGameResults` func. This serves as an intermediate check between the JSON binding and the actual leaderboard update, ensuring that an invalid result is not passed as an update.
- Modified the `SubmitGameResults` to return a HandlerFunc instead of being a handler in and of itself.
- Created 2 unit tests to test 1. A correct input is successful 2. Incorrect inputs should always fail. The 1st one is pretty trivial, so I don't think it requires a lot of test cases (kept it to 1 nominal testcase). Used a table driven test for failure checking, as there are lots of failure scenarios.
- Test setup code was changed - actually use the `UpdateLeaderboard` to set flags, instead of just returning a nil, and a function that creates a Gin router, sets the routes and returns the router.

IMO I think the testing code I wrote is hideous, particularly in terms of the way I handled defining the test case bodies for failure cases. There's also a bit more abstraction, that I don't totally understand well (not too used to Go) so I'll have to read up on that.

There's also the design issues - this code doesn't really follow most of Go's design conventions (which exist for an important reason), which is secondary but would become a problem at some point.

#### Learnings
- Returning handlerfuncs instead of making functions into handlers to allow for dependency injection using closures. This allows me to decouple (not use globals) and simplifies testing as a bonus.
- Testing HTTP endpoints uses all of the functionality that'd be in a normal HTTP call, just internally - marshaling, unmarshaling, sending data in bytes. JSON binding is taken care of by Gin by default.
- Decoupling dependencies for testing code, never call main while testing, tests instead build their own router.
- Design tests before building them, stricter than code in fact, make the invariants clear before building tests.

### Get Leaderboard Top N (`leaderboard.go`)


### Get Player Score (`player.go`)
Weird issue I ran into - `encoding/json` can only work on exported fields, not unexported ones. This made my test `TestPlayer_Success` fail over and over again with the error expected rank 1, got 0, which baffled me since there was no indication as to why this error was occuring like this. The fix was to change the field names in the resp struct to CamelCase (making them exported fields) from `rank` to `Rank` and from `score` to `Score`. Why does json ignore the unexported fields?

### Test SSE access pattern (`sse.go`)
SSE is unarguably the most important access pattern (and the most complicated one too). The code for it is highly coupled, and I can't seem to reduce the coupling by a lot, it works though so it's a good idea to leave it as it is for now.

There's a few things that need to be tested
1. Client lifecycle (add, remove clients, ensure no leaks happen)
2. Broadcasting (and handling of clients that are removed while in the middle of a broadcast - concurrency)
3. Change detection (something changes in store, update is pushed for broadcast)
4. Shutdown

In the near future graceful shutdown is supposed to be added in order to make the app remove all connections and degrade gracefully. This will augment the shutdown process, though it's not a critical need.

#### Testing client lifecycle
1. Add client

