# xm


## Prerequisites
```azure
NOTE: make commands where only tested in linux
```
 - Docker
 - Docker Compose

## Installation
```bash
make init
# adjunst .env file content if needed - defults should work
```

## Usage
```bash
# start the server
make up

# run examples - See cmd/example
make example

#stop the server
make down
```

## Run tests
```bash
make test # runs in dev container
```

## Run linter
```bash
make lint
```

## Get help
```bash
make help
```

# Features
- [x] Production ready (needs more testing, standardize log messages, simplify docker file, fix few TODOs)
- [x] Dockerized
- [x] JWT authentication
- [x] Unit and functional tests
- [x] Linter
- [x] Makefile
- [x] README
- [x] Example usage (this is equivalent to an integration test - see cmd/example)



