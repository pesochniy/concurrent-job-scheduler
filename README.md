# Concurrent Job Scheduler — API Boilerplate

Run the server locally (requires Go >=1.20):

```bash
go run .
```

Server listens on `localhost:8080`.

Endpoints:

- `GET /health` — returns `{ "status": "ok" }`
- `GET /api/hello?name=you` — returns `{ "message": "hello you" }`

Example:

```bash
curl http://localhost:8080/health
curl http://localhost:8080/api/hello?name=alice
```

Graceful shutdown: the server listens for SIGINT/SIGTERM and will shut down within 15s.
