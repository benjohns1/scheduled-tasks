# Scheduled Tasks
[![Go Report Card](https://goreportcard.com/badge/github.com/benjohns1/scheduled-tasks/services)](https://goreportcard.com/report/github.com/benjohns1/scheduled-tasks/services)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](LICENSE)
## Setup
To test and run this locally you'll first need to:
1. Install [Go](https://golang.org/) :-)
2. Install [Node.js](https://nodejs.org/) and [npm](https://www.npmjs.com/)
3. Install [Docker](https://www.docker.com/products/docker-desktop)
4. Install client web app node modules `cd app` and `npm install`
5. Copy `default.env` to `env/local-dev/.env` and `env/local-test/.env` (these environment variables are injected into containers and used in the app)
6. Sign-up and create an [Auth0](https://auth0.com) tenant with the following:
   1. An API (e.g. 'Dev API') with the HS256 signing alg and a single permission of 'type:anon', then set the AUTH0_API_* env vars appropriately
   2. A machine to machine application (e.g. 'Dev Anon User') with 'type:anon' scoped access to your API, then set the AUTH0_ANON_CLIENT_* env vars
   3. A machine to machine application (e.g. 'Dev E2E Test User') with no scoped access to your API, then set the AUTH0_E2E_DEV_CLIENT_* env vars
   4. A single page application (e.g. 'Dev Web App') with Allowed Callback URLs, Allowed Web Origins, Allowed Logout URLs, and Allowed Origins (CORS) set to http://localhost:3000, then set the AUTH0_DOMAIN and AUTH0_WEBAPP_CLIENT_ID env vars
7. Repeat step 6 for the local-test environment if you wish

## Quick Run Scripts (Windows only, for now)
* Development (hot reloading for Svelte app and Cypress tests):
  * `cd env/local-dev`
  * `run`
* Test build (builds and creates container images for the web app and API):
  * `cd env/local-test`
  * `run`

## Default URLs
  * local-dev
    * Web app: http://localhost:3000  
    * DB adminer: http://localhost:3001  
    * API server: http://localhost:3002
  * local-test
    * Web app: http://localhost:3100  
    * DB adminer: http://localhost:3101  
    * API server: http://localhost:3102


## Test Build Environment
1. Build the services:
   * on nix: `(cd services/cmd/srv && env GOOS=linux GOARCH=386 go build)`  
   * -or- on windows: `cmd /C "cd services/cmd/srv&&set GOOS=linux&&set GOARCH=386&&go build"`
2. `cd env/local-test`
3. Build the container images: `docker-compose build`
4. Start the containers: `docker-compose up`
5. Tear it down: `docker-compose down`
6. `cd ../..`

## Development Environment

Run the app and services locally with a transient DB container
1. Start DB: `cd env/local-dev` and `docker-compose up`
2. Modify `./services` code, rebuild as-needed
3. Rebuild and run the services (in a new terminal):
   * on nix: `(cd services/cmd/srv && go build && ./srv)`
   * -or- on windows: `cmd /C "cd services/cmd/srv&&go build&&srv"`
4. `cd ./app`
5. Start the web app with hot reloading: `npm run dev`
6. Open cypress for live testing: `npm run cy:open`
7. Modify `./app` code
8. Tear down the DB: `cd env/local-dev` and `docker-compose down`

## Local Development Testing
### Run services unit tests
1. `cd services`
2. Run tests: `go test ./...`

### Run services integration tests
Connects to DB containers for integration testing
1. Start DBs: `cd env/local-dev` and `docker-compose up`
2. In a new terminal: `cd services`
3. Run tests: `go test ./... -tags="integration"`
4. Tear down DBs: `cd ../env/local-dev` and `docker-compose down`

### Run client web app cypress tests
1. Start DBs: `cd env/local-dev` and `docker-compose up`
2. Build and run the API server & scheduler processes (in a new terminal):
   1. `cd services/cmd/srv`
   1. `go build && ./srv`
3. In a new terminal: `cd app`
4. `npm test`
4. Tear down DBs: `cd ../env/local-dev` and `docker-compose down`
