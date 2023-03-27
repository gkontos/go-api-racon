# Go API RACON

This project is an end to end example api written in go.   The example contains routing, route security, social authentication, jwt token handling, database connectivity, and tests which can be used to bootstrap your next go project.  The project is intended to provide a more complete context for building an application in go than the examples for specific functionality can provide.

This project can be run in google app engine.  This project can also be run locally or within a container using go's ListenAndServe().  

This project is not intended to securely handle secrets.  If you are running a publicly hosted project, don't rely on the environment variables in app.yaml or start.sh for project setup.

RACON: A radar beacon used to aid nautical navigation.

## App Engine 
To run this project in app engine, you will need to setup
- a google cloud project
- ouath 2 client 
- a cloud sql instance
- app engine 
- a public/private key pair ( use ssh-keygen on mac )

Once these are setup, you can add the appropriate env variables to app.yaml and 
```
./deploy.sh
```
will push the project to app engine.

## Local development 
To run this project locally, you will need to setup
- a public/private keypair ( use ssh-keygen on mac )
- ouath 2 client
- a postgres instance

Once these are setup, you can add the appropriate env variables to start.sh and 
```
./start.sh
```
will start a local server for the project.

### Start local db container
The example uses postgres for the backend database which can be deployed locally with docker:
```
sudo mkdir /var/run/postgres
docker run -p 5432:5432 -v /var/run/postgresql/:/var/run/postgresql -v $PWD/local-app.sql:/docker-entrypoint-initdb.d/init-app-db.sql --name postgres -e POSTGRES_PASSWORD=pgpass -d postgres:latest
```
- starts docker container locally w/ socket /var/run/postgresql, user postgres, db postgres
- loads db from local-app.sql file

## container deployment
Since the go http server is run from the main.go file, you can use this project to create a docker image. A dockerfile has not been included.

## Key Components

### router
- This project uses chi because it's just a router.  Gin also looked nice and has more recent github activity, but it handles a much wider range of functionality than needed for an API.  

### logging
- This project uses zerologger because it has a decent API and it's reportedly fast.  Being able to switch between a structured logger for deployment and a console logger for local development was also important.


## Todo
- Implement a connection to secrets manager as a better method of secrets management than env variables
- More comphrensive testing
