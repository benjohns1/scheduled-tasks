# Scheduled Tasks
## Task app with scheduled recurrences

### Run unit tests locally
`go test ./...`

### Run integration tests locally
`go test ./... -tags="integration"`

### Build & run
1. Install [Docker](https://www.docker.com/products/docker-desktop)
2. Copy `.env.default` to `.env` (these environment variables are injected into containers and used in the app)

#### Dev environment
Run & build the app locally, run a transient DB in Docker container
1. `docker-compose up`
2. Build and run the server:
   1. `cd cmd/server`
   2. `set GOOS=<your-local-OS>`
   3. `go build && ./server`
3. Server: `localhost:8080`
4. DB Adminer: `localhost:8081`
5. Tear it down: `docker-compose down`

#### Test environment
Build the app locally, run it and a transient DB in Docker containers
1. Build the server and image:
   1. `cd cmd/server`
   2. `set GOOS=linux`
   3. `go build`
   4. `docker build --no-cache -t scheduled-tasks_api .`
   5. `cd ../..`
2. `docker-compose -f docker-compose.test.yml up`
3. Server: `localhost:8080`
4. DB Adminer: `localhost:8081`
4. Tear it down: `docker-compose -f docker-compose.test.yml down`