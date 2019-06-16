# Scheduled Tasks
[![Go Report Card](https://goreportcard.com/badge/github.com/benjohns1/scheduled-tasks/services)](https://goreportcard.com/report/github.com/benjohns1/scheduled-tasks/services)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](LICENSE)
## Setup
To test and run this locally you'll first need to:
1. Install [Go](https://golang.org/) :-)
2. Install [Node.js](https://nodejs.org/) and [npm](https://www.npmjs.com/)
3. Install [Docker](https://www.docker.com/products/docker-desktop)
4. Copy `.env.default` to `.env` (these environment variables are injected into containers and used in the app)
5. Install client web app node modules `cd app` and `npm install`

## Staging Environment
1. Build the services:
   1. if on windows: `set GOOS=linux`
   2. `cd services/cmd/srv`
   4. `env GOOS=linux GOARCH=386 go build`  
   or, if on windows: `go build`
   5. if on windows: `set GOOS=windows`
   6. `cd ../../..`
2. Rebuild app & service images: `docker-compose -f docker-compose.stage.yml build`
3. Start the containers: `docker-compose -f docker-compose.stage.yml up`
4. Web app: `localhost:3000`
5. API server: `localhost:8080`
6. DB adminer: `localhost:8081`
7. Tear it down: `docker-compose -f docker-compose.stage.yml down`

## Testing
### Run services unit tests
1. `cd services`
2. Run tests: `go test ./...`

### Run services integration tests
Connects to DB containers for integration testing
1. `docker-compose up`
2. In a new terminal: `cd services`
3. Run tests: `go test ./... -tags="integration"`
4. `cd ../`
5. `docker-compose down`

### Run client web app cypress tests
1. `docker-compose up`
2. Build and run the API server & scheduler processes (in a new terminal):
   1. `cd services/cmd/srv`
   1. `go build && ./srv`
3. In a new terminal: `cd app`
4. `npm test`
5. `cd ../`
6. `docker-compose down`

## Development Environment

Run the app and services locally with a transient DB container
1. `docker-compose up`
2. Modify `./services` code, rebuild as-needed
3. Rebuild and run the services (in a new terminal):
   1. `cd services/cmd/srv`
   2. `go build && ./srv`
5. API server: `localhost:8080`
6. DB adminer: `localhost:8081`
7. `cd app`
8. Start the web app with hot reloading: `npm run dev`
9. Open cypress for live testing: `npm run cy:open`
10. Modify `./app` code
11. Web app: `localhost:3000`
12. Tear it down: `docker-compose down`