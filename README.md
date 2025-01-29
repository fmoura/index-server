# Index Server

This project is a home task given by Gosolve. You can find the actual task description [here](recruitment_task.md)

## Introduction

The project was based on the [GoFr](gofr.dev) framework. It is a good quick start for APIs written in Go as it gives out of the box:
- Configuration file
- Configurable Logging
- Swagger rendering
- Easy Routing

Note: Other libraries like Echo ( for request routing and handling), Slog for logging and Viper for configuration could have a better end result. The way GoFr structs logging and configurations are not the best.

## Features

As asked by the [task description](recruitment_task.md), this project is a service that exposes just an HTTP endpoint.

This endpoint searches for the `index` of a given `value` at the `input data`. The `input data` is a collection of sorted numbers ranging from 0 til 1000000.

The service can respond an exact `index` if the given `value` is actually present at the `input data` or a close `index` and `actual value` if `value` and `actual value` are within a `conformation level`.

Internally the search is done using a *Binary Search* algorithm

## Configuration

The server configuration is handled by [GoFr](gofr.dev). To configure the application it uses `.env` files placed at the `configuration` folder, as well as the system environment variables.

### Conformation Level

The `index` endpoint can respond to index searches with a approximate answer if the given `value` is not part of the `input data`, if the actual value is within a `conformation level`. To set the `conformation level` use the `CONFORMATION_LEVEL` config. The confirmation level must be a integer at the range 0 til 100 and it represents the conformation level in percents. The default value is `10`

```conf
CONFORMATION_LEVEL=10
```

### HTTP Port

To configure which HTTP Port the server will accept HTTP requests is configured by the `HTTP_PORT` config. The default [GoFr](gofr.dev) port is`8000`

Example:
```conf
HTTP_PORT=8000
```

### Log Level

Logging is also handled by [GoFr](gofr.dev), and the log level is configured by the `LOG_LEVEL` config. GoFr has an API to change the log level while the service is running, more information about it can be found [here](https://gofr.dev/docs/advanced-guide/remote-log-level-change). The default value is `INFO`

Example:
```conf
LOG_LEVEL=INFO
```
## Run

To run the server just open a shell and run:

```shell
make run
```

## Usage

The only exposed API endpoint is the `index/{value}`. It searches for the given `value` index at the `input data` where `value` is a integer in the range from 0 til 1000000. Any other value will raise a `invalid parameter error` with `HTTP 400 status`

To experiment the API, once the server is running (see [Run](#run) section above) open a browser and access the Swagger page at `{ADDRESS}/.well-known/swagger`, where `SERVER:PORT` is the server address (ex: http://localhost:8000) or open a
system shell and run:

```shell
curl http://localhost:8000/index/200
```
As the given `value` is `200`, and it is present at the `input data`, you should receive a response with `HTTP 200` status with payload like:
```json
{
  "data": {
    "index": 2,
    "value": 200
  }
}
```

Now try the following:

```shell
curl http://localhost:8000/index/110
```
As the given value is `110`, and it is not present at the `input data`, but is within the conformation criteria, you should receive a `HTTP 200` status response with the following content:
```json
{
  "data": {
    "index": 1,
    "value": 100
  }
}
```

Now let's search for a `value` that do not fit the data or conformation criteria:
```shell
curl http://localhost:8000/index/150
```
The value is `150` that is not present at the `input_data` and also falls out of the conformation criteria, so the response comes with `HTTP 200` status  but a message that informs the value is not found

```json
{
  "data": {
    "index": -1,
    "value": 150,
    "error message": "Value not found"
  }
}
```

As explained at the beginning of this section, try a invalid value:
```shell
curl http://localhost:8000/index/value
```

You should receive a response with `HTTP 400` status and payload like this:
```
{
  "error": {
    "message": "'1' invalid parameter(s): value"
  }
}
```

## Test

To run all the tests, simply run:

```shell
make test
```

## Architecture

The project folder structure follows Go best practices and GoFr conventions

The Server is structured on the *three tier* architecture for better modularity and separation of concerns:

### Presentation Tier

The presentation tier is represented by [`IndexHandler`](https://github.com/fmoura/index-server/blob/39af9bcf1e0898a9ae09edc9ea1e36d5e252bb58/internal/handlers/index.go) at `handlers` package. It is responsible for handling the HTTP requests and calling the `IndexService` to fulfill the request, and give well formatted *json* responses

### Business Logic Tier

The business logic tier is represented by the [`IndexService`](https://github.com/fmoura/index-server/blob/39af9bcf1e0898a9ae09edc9ea1e36d5e252bb58/internal/service/index.go) at the `service` package. It do the actual index search


### Data Tier

The data tier is represented by the [`TextDataProvider`](https://github.com/fmoura/index-server/blob/39af9bcf1e0898a9ae09edc9ea1e36d5e252bb58/internal/data/provider.go) at the `data` package. It loads the `input data`