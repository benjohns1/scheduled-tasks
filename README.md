# Scheduled Tasks
[![Go Report Card](https://goreportcard.com/badge/github.com/benjohns1/scheduled-tasks/services)](https://goreportcard.com/report/github.com/benjohns1/scheduled-tasks/services)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](LICENSE)
## Setup
### Local
To test and run this locally you'll first need to:
1. Install [Go](https://golang.org/) :-)
2. Install [Node.js](https://nodejs.org/) and [npm](https://www.npmjs.com/)
3. Install [Docker](https://www.docker.com/products/docker-desktop)
4. Install client web app node modules `cd app` and `npm install`
5. Install end-to-end test runner `cd app-test` and `npm install`
6. Copy `default.env` to `env/local-dev/.env` (these environment variables are injected into containers and used in the apps when running locally)
7. Create an Auth0 tenant for the `local-dev` environment (See Auth0 Authentication section below)
8. Create another Auth0 tenant for the `local-test` environment if you wish

### AWS Staging
App will be public, but NOT protected with an SSL cert (which will break auth requests in chromium browsers)
1. Install [Terraform](https://www.terraform.io)
2. Setup AWS credentials on your local machine
3. Create a hosted zone for your domain name with a staging subdomain in AWS Route53
5. Create a key pair to connect to your EC2 instance in AWS EC2 under Network & Security -> Key Pairs
6. Copy `default.secret.auto.tfvars` to `env/local-stage/.secret.auto.tfvars` (these are used when spinning up cloud infrastructure with terraform)
7. Create an Auth0 tenant for the staging environment and set the .secret.auto.tfvars appropriately
8. Build/push the containers and spin up AWS infrastructure, will incur costs (`terraform destroy` to tear down)
```sh
cd ./services
docker build -t benjohns1/scheduled-tasks-services .
docker push benjohns1/scheduled-tasks-services

cd ../app
docker build -t benjohns1/scheduled-tasks-webapp .
docker push benjohns1/scheduled-tasks-webapp

cd ../env/aws-stage
terraform taint module.ecs.aws_ecs_task_definition.tasks
terraform apply
```

### Auth0 Tenant Setup
Sign-up and create an [Auth0](https://auth0.com) tenant for each environment you want with the following:
1. An API (e.g. 'Dev API') with the HS256 signing alg and a single permission of 'type:anon', then set the AUTH0_API_* env vars appropriately
2. A machine to machine application (e.g. 'Dev Anon User') with 'type:anon' scoped access to your API, then set the AUTH0_ANON_CLIENT_* env vars
3. A machine to machine application (e.g. 'Dev E2E Test User') with no scoped access to your API, then set the AUTH0_E2E_DEV_CLIENT_* env vars
4. A single page application (e.g. 'Dev Web App') with Allowed Callback URLs, Allowed Web Origins, Allowed Logout URLs, and Allowed Origins (CORS) set to http://localhost:3000, then set the AUTH0_DOMAIN and AUTH0_WEBAPP_CLIENT_ID env vars

## Quick Run Scripts
Setup local environments in one command (Windows only, for now)
* Development (hot reloading for the Sapper app and Cypress interactive test runner):
  * `cd env/local-dev` and `run`
* Test build (builds and creates container images for the web app and API):
  * `cd env/local-test` and `run`

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
1. `cd ./env/local-test`
2. Point to the test version of the env vars:
   * on nix: `export ENV_FILEPATH=$PWD/.env`
   * -or- on windows: `set ENV_FILEPATH=%CD%\.env`
3. Build the app and container images: `docker-compose build`
4. Start the containers: `docker-compose up`
5. Tear it down: `docker-compose down`

## Development Environment
Run the app and services locally with a transient DB container
1. `cd ./env/local-dev`
2. Point to the dev version of the env vars:
   * on nix: `export ENV_FILEPATH=$PWD/.env`
   * -or- on windows: `set ENV_FILEPATH=%CD%\.env`
3. Start DBs: `docker-compose up`
4. Modify `./services` code, rebuild as-needed
5. Rebuild and run the services in `./services/cmd/srv`:
   * on nix: `go build && ./srv`
   * -or- on windows: `go build && srv`
6. Start the web app with hot reloading in `./app` with: `npm run dev`
7. Open cypress for live testing in `./app-test` with: `npm run cy:open`
8. Modify `./app` code
9. Tear down the DBs in `./env/local-dev` with: `docker-compose down`

## Local Development Testing
### Run services unit tests
1. Run tests in `./services` with: `go test ./...`

### Run services integration tests
Connects to DB containers for integration testing
1. Start DBs in `./env/local-dev` with: `docker-compose up`
1. Point to the test version of the env vars:
   * on nix: `export ENV_FILEPATH=$PWD/.env`
   * -or- on windows: `set ENV_FILEPATH=%CD%\.env`
1. Run tests in `./services` with: `go test ./... -tags='integration'`
1. Tear down the DBs in `./env/local-dev` with: `docker-compose down`

### Run client web app cypress tests
1. Start DBs in `./env/local-dev` with: `docker-compose up`
1. Build and run the API server & scheduler processes (in a new terminal) in `./services/cmd/srv` with: `go build && ./srv`
1. Run the tests in `./app-test` with: `npm test`
1. Tear down the DBs in `./env/local-dev` with: `docker-compose down`
