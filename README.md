# Scheduled Tasks
[![Go Report Card](https://goreportcard.com/badge/github.com/benjohns1/scheduled-tasks/services)](https://goreportcard.com/report/github.com/benjohns1/scheduled-tasks/services)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](LICENSE)
## Setup
To test and run this locally you'll first need to:
1. Install [Go](https://golang.org/) :-)
2. Install [Node.js](https://nodejs.org/) and [npm](https://www.npmjs.com/)
3. Install [Docker](https://www.docker.com/products/docker-desktop)
4. Copy `.env.default` to `.env` (these environment variables are injected into containers and used in the app)

## API and Scheduling Services
### Run tests for the API and scheduling services
#### Run unit tests locally
1. `cd services`
2. `go test ./...`

#### Run integration tests locally
Connects to DB containers for integration testing
1. `docker-compose up`
2. In a new terminal: `cd services`
3. `go test ./... -tags="integration"`
4. `cd ../`
5. `docker-compose down`

### Build & run the API and scheduling services

#### Dev/test environment
Run & build the services locally, run a transient DB in Docker container
1. `docker-compose up`
2. Build and run the API server & scheduler processes (in a new terminal):
   1. `set GOOS=<your-local-OS>`
   2. `cd services/cmd/srv`
   3. `go build && ./srv`
3. API Server: `localhost:8080`
4. DB Adminer: `localhost:8081`
5. Tear it down: `docker-compose down`

#### Staging environment
Build the services locally, run them and a transient DB in Docker containers
1. Build the server and image:
   1. `set GOOS=linux`
   2. `cd services/cmd/srv`
   4. `go build`
   5. `docker build --no-cache -t scheduled-tasks .`
   9. `cd ../../..`
2. `docker-compose -f docker-compose.stage.yml up`
3. Server: `localhost:8080`
4. DB Adminer: `localhost:8081`
5. Tear it down: `docker-compose -f docker-compose.stage.yml down`