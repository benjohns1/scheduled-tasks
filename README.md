# Scheduled Tasks
## Task app with scheduled recurrences
### Run all tests locally
`go test ./...`
### Run local test environment
1. Install [Docker](https://www.docker.com/products/docker-desktop)
2. Copy `.env.default` to `.env` (these environment variables are injected into containers and used in the app)
3. Build the server and image:
   1. `cd cmd/server`
   2. `set GOOS=linux`
   3. `go build`
   4. `docker build --no-cache -t scheduled-tasks_api .`
   5. `cd ../..`
4. `docker-compose up`

Go to [localhost:8080]() for the server api and [localhost:8081]()

To tear it down: `docker-compose down`