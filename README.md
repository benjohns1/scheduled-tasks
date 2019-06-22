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

## Quick Run Scripts (windows only for now)
* Staging: `run.stage`
* Development: `run.dev.webapp` (hot reloading for Svelte app and Cypress tests)

Web app: http://localhost:3000  
API server: http://localhost:8080  
DB adminer: http://localhost:8081  

## Staging Environment
1. Build the services:
   * on nix: `(cd services/cmd/srv && env GOOS=linux GOARCH=386 go build)`  
   * -or- on windows: `cmd /C "cd services/cmd/srv&&set GOOS=linux&&set GOARCH=386&&go build"`
2. Build the container images: `docker-compose -f docker-compose.stage.yml build`
3. Start the containers: `docker-compose -f docker-compose.stage.yml up`
   * Web app: `localhost:3000`
   * API server: `localhost:8080`
   * DB adminer: `localhost:8081`
4. Tear it down: `docker-compose -f docker-compose.stage.yml down`

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
   * on nix: `(cd services/cmd/srv && go build && ./srv)`
   * -or- on windows: `cmd /C "cd services/cmd/srv&&go build&&srv"`
   * API server: `localhost:8080`
   * DB adminer: `localhost:8081`
6. `cd app`
7. Start the web app with hot reloading: `npm run dev`
8. Open cypress for live testing: `npm run cy:open`
9.  Modify `./app` code
    * Web app: `localhost:3000`
    * Tear it down: `docker-compose down`