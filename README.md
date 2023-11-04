# Airport Lookup - Go

The goal of this project is to serve as an example for building Go applications with the 
[Tarmac application framework](https://github.com/tarmac-project/tarmac). As Tarmac 
continues to evolve, so does this project; contributions, of course, are welcomed!

Tarmac is a new approach to building services and leverages WebAssembly to offer a language-agnostic application 
framework.

The goal of Tarmac is to abstract the non-functionals and let you focus purely on the business logic of your 
applications and less on the infrastructure and boilerplate code. However, building services with Tarmac is different 
from building a traditional application.

Rather than considering an application as a single entity, consider it a collection of serverless functions. A key 
difference is that this collection of serverless functions runs together and can communicate with each other through 
internal calls.

## Airport Lookup API

This service is a simple API that provides details of Airports. On boot, the service will fetch a CSV containing
detailed airport information and store the data within a MySQL database.

The lookup API is accessible via a POST request to `/`.

```console
$ curl -X POST http://localhost/ -d '{"local_code": "PHX"}' 
{"local_code": "PHX", "name": "Phoenix Sky Harbor International Airport", "country": "US", "emoji": "🇺🇸", "type": "large_airport", "type_emoji": "✈️", "status": "open"}
```

## Running the Service

To run the entire service, use the following make commands.

```console
$ make run
```

This project leverages Docker Compose to run the service and its dependencies. Once the service runs, you can access 
the API by calling <http://localhost/>.

### Non-Functionals

This service is meant to showcase the out of the box non-functional capabilities of Tarmac. This section provides a
high-level overview of the non-functionals implemented within this service.

#### Health and Readiness Checks

This service offers out-of-the-box health and readiness checks accessible via the following end-points.

`/health`

The health end-point reflects the liveness of this service. The logic behind this end-point is to return a `200 OK` once 
the HTTP server is up and running.

`/ready`

The ready end-point reflects the readiness of this service. The logic behind this end-point is to validate that all 
service dependencies are up and accessible. If everything is up and ready, a `200 OK` is returned.

#### Observability

This service leverages metrics & logging for observability. Users can enable Debug and Trace logging via the 
environment variable configuration in the `docker-compose.yml` file.

Applications metrics are exposed via the `/metrics` end-point.

#### Caching

This service leverages Redis as a caching layer; as requests to the API arrive, the handler function will store 
results within Redis and use these results on subsequent requests.

## Project Structure

There are three critical directories within this repository: config, functions, & pkg.

### Config

The configuration directory holds both the Tarmac configuration and any ancillary configuration files.

### Functions

The functions directory is home to the source code for the various serverless functions that comprise this application. 

### Pkg

The pkg directory is a traditional Go packages directory with packages used and imported throughout this application.
