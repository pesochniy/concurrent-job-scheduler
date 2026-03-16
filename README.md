# Concurrent Job Scheduler (Go)

A background job processing service written in **Go** with **PostgreSQL persistence**.  
Jobs are submitted via HTTP API, stored in the database, and executed concurrently by a worker pool.

## Features

- HTTP API for job submission and retrieval
- Concurrent job processing using goroutines and channels
- Worker pool architecture
- Job state machine (pending → running → completed / failed)
- Retry policy with retry limits
- Job execution timeout
- PostgreSQL persistence (pgx)
- Scheduler polling for pending jobs
- Graceful shutdown with context
- Dockerized deployment

## Architecture

```
Client → HTTP API → PostgreSQL
                     ↓
                 Scheduler
                     ↓
               Worker Pool
```

Jobs are stored in the database, picked up by the scheduler, and executed by workers.

## API

Create job:

```bash
curl -X POST localhost:8080/jobs \
-H "Content-Type: application/json" \
-d '{"type":"email","payload":{"to":"user@example.com","subject":"Hello","body":"Test"}}'
```

List jobs:

```bash
curl localhost:8080/jobs
```

Health check:

```
GET /health
```

## Running the Project

Requirements:

- Docker
- Docker Compose

Start services:

```bash
docker compose up --build
```

API will be available at:

```
http://localhost:8080
```

Stop services:

```bash
docker compose down
```

## Project Structure

```
internal/
  domain      job payload models
  job         job model and storage
  scheduler   scheduler and worker pool
  workers     job execution logic

handlers      HTTP handlers
schema.sql    database schema
```

## Tech Stack

- Go
- PostgreSQL
- pgx
- Docker