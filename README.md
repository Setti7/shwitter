# Shwitter

###### Shitpost like there is no tomorrow

---

Shwitter is like twitter but where you have fun instead of being pissed of by other people's stupidity.

# Getting Started

The easiest way of setting up your local development environment is by using docker-compose. After installing docker and
docker-compose, just run `make up` to start the dependencies for the local server.

While waiting for the docker images to download, run `make build` to build the go backend. When completed, execute
`./shwitter setup` to finalize setting up the local environment.   
After it succeeds, run `./shwitter start` to start the server with the default configuration.

To checkout more commands or help with each command, just run `./shwitter help`.

---

### Database migrations

We are using [migrate](https://github.com/golang-migrate/migrate/) to manage the database migrations. If you don't have
it installed already, please go to
[here](https://github.com/golang-migrate/migrate/blob/master/GETTING_STARTED.md) and follow the installation steps.

