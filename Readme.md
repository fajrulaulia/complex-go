# ModulorGO

## Overview
 ModulorGO is a Backend module written by golang, easy to Write Code, Testing and Build.

## Requirement
- OS Linux or Mac.
- Make sure you system installed `docker-engine`,`docker-compose` and `make`.
> you don't need to install MySQL server, Relax.

## Commands
- `make service-up` : Create Service like MySQL Server for Local Development.
- `make service-down` : Stop and Remove  Service like MySQL Server for Local Development. 
- `make local-run`  : Run local server golang.
- `make local-test` : Run Test locally.
- `make build`: Build image Application.
- `make deploy`: Deploy image to registry.

## Development
- make sure you have access right from me :) .
- clone from `git@gitlab.com:fajrulaulia/modulor-go.git`.
- cd to `modulor-go` folder.
- run `make service-up` to start service like MySQL server.
- run `make local-run` to run your code.

## Testing
- Testing can you write in `$ROOT_PROJECT/app/test/`.
- run `make local-test`.
- automatically, a service for testing will be created.

## Migrations
- for migrations, this module have a two folder
    - `migrations/test` : for testing when you create a new service and you need DDL and DML.
    - `migrations/build`: for production when you create a new service and you need DDL and DML.

## Contributor
Contributor: Fajrul Aulia
