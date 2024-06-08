
# Airport Lookup

This project is a practical demonstration of how to build applications with [Tarmac](http://github.com/tarmac-project/tarmac). It's a hands-on guide that shows you how to leverage Tarmac's unique features to create efficient and scalable distributed systems.

Tarmac is a unique approach to building distributed systems. It is an application framework that lets you create your application as a set of serverless-like functions.

Unlike traditional serverless functions, where you would deploy each function into a serverless platform, Tarmac offers a more efficient approach. With Tarmac, you can compose your functions into a single monolithic application instance or a set of microservices, enhancing your system's scalability and efficiency.

Tarmac lets you create small, testable logic, handles non-functional requirements, and enables you to run them where and how you see fit.

## Airport Lookup API

This project is a collection of functions that provide a simple API to look up details of Airports worldwide. These functions are WebAssembly modules that the Tarmac application server will run.

The table below describes each function and its purpose.

| Function | Description |
|----------|-------------|
| handlers/lookup | This function is the main API handler. It receives a POST request with a local_code parameter and returns the airport details. It first looks up the airport data in its Key:Value cache and, if not found, searches the SQL database and updates the cache. |
| data/init | This function initializes the SQL database and triggers the load function. |
| data/load | This function calls the fetch function parses the CSV data and stores it in the SQL database. |
| data/fetch | This function downloads the airport data CSV file. |

Currently, this project showcases two deployment models: a set of microservices where the API traffic is split from the data management tasks or as a single monolith.

### Monolith

In the monolith model, all functions run within a single application instance called lookup. This single lookup service will perform all tasks.

```
               +----------------------------+                
               | Lookup                     |                
               | +---------+                |                
               | |data/init|                |  +------------+
               | +----+----+                +-->SQL Database|
               |      |                     |  +------------+
               | +----v----+                |                
               | |data/load|                |                
               | +----+----+                |                
               |      |        +----------+ |                
               |      +-------->data/fetch| |                
               |               +----------+ |                
               |                            |                
+----------+   | +---------------+          |  +---------+   
|API Client+---+->handlers/lookup|          +-->K/V Cache|   
+----------+   | +---------------+          |  +---------+   
               |                            |                
               +----------------------------+                
```
 
### Microservices

In the microservices model, API traffic is handled by a lookup service. This service is the frontend API layer; it serves client requests by first looking up the Airport data in its Key:Value cache and, if not found, searching the SQL database.

A second data-manager service is a backend management process. Its job is to manage the Airport data by periodically downloading and storing it within the SQL database.

```
               +------------------------------+              
               |Data Manager                  |              
               |   +---------+                |              
               |   |data/init|                |              
               |   +----+----+                |              
               |        |                     |              
               |   +----v----+                |              
               |   |data/load|                |              
               |   +----+----+                |              
               |        |        +----------+ |              
               |        +-------->data/fetch| |              
               |                 +----------+ |              
               |                              |              
               +-----+-----------------+------+              
                     |                 |                     
                +----v----+     +------v-----+               
                |K/V Cache|     |SQL Database|               
                +----^----+     +------^-----+               
                     |                 |                     
               +-----+-----------------+------+              
               |Lookup                        |              
+----------+   | +--------------+             |              
|API Client+---+->handler/lookup|             |              
+----------+   | +--------------+             |              
               |                              |              
               +------------------------------+              
```

The microservices model provides faster and more consistent API response times, as the lookup service no longer has the overhead of downloading the airport data and loading that data into the database.

## Running the service

To run the microservices deployment, use the following command:

```console
$ make run
```

This project leverages Docker Compose to run the services and their dependencies. Once the lookup service starts and data is loaded, you can access the API by calling <http://localhost/>.

```console
$ curl -X POST http://localhost/ -d '{"local_code": "PHX"}'
{"local_code": "PHX", "name": "Phoenix Sky Harbor International Airport", "country": "US", "emoji": "🇺🇸", "type": "large_airport", "type_emoji": "✈️", "status": "open"}
```

### Non-Functionals

One of the best things about Tarmac is its out-of-the-box, non-functional capabilities. With Tarmac, you can think less about making systems reliable and more about your application logic.

The below section provides a high-level overview of the non-functionals implemented within this service.

#### Health and Readiness Checks

This service offers out-of-the-box health and readiness checks accessible via the following end-points.

##### `/health`

The health end-point reflects the liveness of this service. The logic behind this end-point is to return a 200 OK once the HTTP server runs.

##### `/ready`

The ready end-point reflects the readiness of this service. The logic behind this end-point is to validate that all service dependencies are up and accessible. If everything is up and ready, a 200 OK is returned.

#### Observability

This service leverages metrics and logging for observability.

##### Logging

Users can enable Debug and Trace logging via the environment variable configuration.

##### Metrics

Application metrics are exposed via the `/metrics` end-point. A dashboard is available at <http://localhost:3000> (username: admin, password: example).

#### Caching

This service leverages caching. Tarmac supports multiple caching backends via its Key:Value store interface. The default deployment uses Redis, with an in-memory cache used for the load-testing deployment.

## Project Structure

There are three critical directories within this repository: config, functions, & pkg.

### Config

The configuration directory holds both the Tarmac configuration and any ancillary configuration files.

### Functions

The functions directory is home to the source code for the various serverless functions that comprise this application. 

### Pkg

The pkg directory is a traditional Go packages directory with packages used and imported throughout this application.

## Contributions

Contributions to this project are welcome! If you have a simple suggestion, please open a pull request. Feel free to open an issue if you have a more complex suggestion.
