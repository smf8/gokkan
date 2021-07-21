# Gokkan
Gokkan is an online shop backend service written in Go.

## Setup
### Requirements
Gokkan requires `postgresql` as database.
You can use `docker-compose` to get an instance up and running.

**Note** the default settings in `internal/app/gokkan/config/default.go`
when connecting to the database.

A `pgAdmin` instance is also inside docker-compsoe. you can use `localhost:8000` to check database with `pgAdmin`

### Install
Clone the repository inside a folder **outside** of GOPATH
```shell
git clone https://github.com/smf8/gokkan
cd gokkan
make build

# important before running the application
# please make sure that postgres is up and running
# before running migrates
make migrate-up
```
aAfter that you can use `gokkan` binary to run the server.
Use `gokkan -h` to see available commands. currently available commands are:
```shell
  ./gokkan server      # start the server
```


To clear database after tests run 
```shell
make migrate-reset
```

## Usage
Echo server will start listening on port `8080` by default.

A (Postman Collection file)[gokkan_api.json] is provided to describe API behaviour