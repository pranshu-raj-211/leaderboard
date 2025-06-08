The config package handles application configuration and logging setup.

## Features

- YAML-based configuration
- Structured logging with Zap
- Environment-aware settings
- Centralized error handling

## Configuration Parameters

- Redis settings (address, password, retries)
- Server settings (host, port)
- Leaderboard settings (update intervals, limits)

## Logger

Uses Uber's Zap logger with:
- JSON formatting
- File and stdout output
- Error level control