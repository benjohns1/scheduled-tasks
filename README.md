# Scheduled Tasks
[![Go Report Card](https://goreportcard.com/badge/github.com/benjohns1/scheduled-tasks)](https://goreportcard.com/report/github.com/benjohns1/scheduled-tasks)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](LICENSE)
## Task app with scheduled recurrences
To test and run this you'll first need to:
1. Install [Docker](https://www.docker.com/products/docker-desktop)
2. Copy `.env.default` to `.env` (these environment variables are injected into containers and used in the app)

### Run tests
#### Run unit tests locally
1. `go test ./...`

#### Run integration tests locally
Connects to DB containers for integration testing
1. `docker-compose up`
2. `go test ./... -tags="integration"`
3. `docker-compose down`

### Build & run

#### Dev/test environment
Run & build the app locally, run a transient DB in Docker container
1. `docker-compose up`
2. `set GOOS=<your-local-OS>`
3. Build and run the API server:
   1. `cd cmd/srv`
   2. `go build && ./srv`
4. Build and run the scheduler process:
   1. `cd cmd/sched`
   2. `go build && ./sched`
5. API Server: `localhost:8080`
6. DB Adminer: `localhost:8081`
7. Tear it down: `docker-compose down`

#### Staging environment
Build the app locally, run it and a transient DB in Docker containers
1. Build the server and image:
   1. `set GOOS=linux`
   2. `cd cmd/srv`
   4. `go build`
   5. `docker build --no-cache -t scheduled-tasks_api .`
   6. `cd ../sched`
   7. `go build`
   8. `docker build --no-cache -t scheduled-tasks_sched .`
   9. `cd ../..`
2. `docker-compose -f docker-compose.stage.yml up`
3. Server: `localhost:8080`
4. DB Adminer: `localhost:8081`
5. Tear it down: `docker-compose -f docker-compose.stage.yml down`