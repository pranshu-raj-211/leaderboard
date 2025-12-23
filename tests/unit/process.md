# Process of adding tests, improving code, and learning from mistakes

## Validating GameResult (models.go)
Smallest part of the code, independent from other parts (has downstream dependencies), perfect starting point for tests.

Decided on a basic table driven test (similar to the one I wrote in dbos), to check for an error and report when it's different to expectations (got error when none expected, or got no error when one should have been returned).

This basic test uncovered one issue - Validate did not check for `ServerID`, so not passing that in worked, Go initialized that to an empty string. This is undesirable behavior and was fixed.

This simple error check when expected fell apart soon when testcases with multiple errors were tried out (based on previous case - kept ServerID empty when result was too high). I realized that the tests would need some way of identifying which type of thing had failed (different from just adding a message to the error).

To fix this, I created a `ValidationError` which had a `Field` and `Reason` keys. The `Field` showed what part of the GameResult object did not have something as expected. The Reason was a message based on the type of check applied during validation. This provides more clarity, and is more deliberate about what kind of errors are being found on a more granular level, which I believe is good design for this code.

### Learnings
- `errors.As` takes in an error, a target, and if the error is of the type target it assigns the value of the error to target. This can later be used as I did for the last test check involving the Field.
- Typed errors - The `ValidationError` that I created is a pattern in Go named a type error. It's elegant and reminds me of Python. Lots of people do not like the implicit implementation of interfaces in Go, but I believe it has a lot of uses, and particularly like the uses it has in testing (much better to not implement interfaces and inherit classes explicitly in this case - for mocking or dependency injection as an example).
- zap.NewNop is a no op logger (doesn't actually log anything or write internal errors, perfect for testing code that has logging dependencies).

## Game submissions (backend/submit_game.go)
Work in progress

This is related to the `/submit-game` endpoint, required for game servers to submit game results, which are used to update scores in the leaderboard. This has quite a few dependencies, but significantly less than other endpoints, and the code is easy to test as well (few binding, validation and Redis update).

Designing a test for this made me realize how tightly everything is coupled in my code. I did make a basic version of dependency injection for the broadcaster in the SSE streaming code, but Redis and other dependencies are still as tightly coupled as ever. This is a problem, since I'd have to use those as is for testing, while I agree things can be done without mocking, mocks will make things simpler, and if used correctly, make the code easier to reason about (provided I do not go overboard with things).

DI is something I really looked forward to using in FastAPI, it made my life a whole lot easier, and I look forward to understanding how it's done in production systems running Go.