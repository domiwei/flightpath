## How to Run

Run the following command to build and run the api server with log in stdout
```
make run
```
Run this if you don't like the output log 
```
make run-detach
```
Here is how to stop background docker containers
```
make stop
```

## Code Structue
- `app/main.go` setup all the resources and register apis, then launch `gin` api server.
- `api/api.go` registers endpoints by corresponding handler, which is defined in `flight.go`
- `api/flight.go` implements how I store the query the flights data, and algorithm to check if a route is valid or not.
- `service/db/db.go` is a simple wrapper to initialize a mysql engine.
- `Dockerfile`, `docker-compose.yml`, `Makefile` are scripts related to infrastucture helping us build the entire project.
- `init.sql` defines the schema to alter our table.

## Design and Anything to Share
- `Database`: As a simple micro service aiming to store and be queried people flight path, here I use MySQL as underlying storage to keep our data permanent on disk. The main reason why I chose SQL instead of NoSQL is because our flight data seems really structured with fixed fields while we design our db schema. And MySQL is also a good choice if we might need to join this table with other tables in order to get some mixed result data in the future.
Besides, according to the spec description in assignment, I believe this service doesn't suffer from high WRITE operations in any circumstance. In terms of READ operation, it's a good idea to cache our flights data with key "personID" (maybe passportID) so we don't need to worry about this case. In short, I believe traditional SQL DB is good enough to handle what we are going to do, and also, SQL is robust enough to protect our data.
- `departureTime`: In my implementation, departureTime is optional parameter while insert a new record. In the beginning I was thinking about the unique key of the table, and a bit having trouble on how to deal with the case that insert multiple records from identifcal person. Finally I found that what I need is another field indicating when does the tour start. Hence in my api server, user is able to post flights with departure time so that our server can store many flights records in DB.
- `Things can be improved`: A typical api server must need to consider the throughput like QPS during peak time, or average QPS for post/get api, so this server is just a really simple server only finishing the requirements. Here I list out something can be improved:
    ```
    1. Cache the result after we get flights data from underlying database. Just like local cache or redis.
    2. Check the index of database. Make sure my query meets the index.
    3. Apply ratelimit in api middleware to protect our micro service.
    4. Depends on how much data we may recive, need to re-design how we store data in DB.
    5. It is not well-defined whether a flight path is valid or not. Need to clarify it. 
    ```
- Because watched FIFA worldcup, I had limited time to complete this assignment. It would be helpful if having one more day to do the job. Just a suggestion :).



## API Endpoint
##### POST `/v1/calculate/{personID}/flights`
- Body
    ```
    {
        "path": <tour path in json array format, see example below>,
        "departureTime": uint64 // optional, departure time of first flight
    }
    ```
- Response
    ```
    200: success
    400: invalid or unreasonable path  
    500: internal server errror
    ```
- Example Body
    ```
    {
        "path": [["city2", "city3"], ["city1", "city2"], ["city3", "city4"]],
        "departureTime": 1670696785
    }
    ```
- Example Request
    ```
    curl -X POST localhost:8080/v1/calculate/kewei/flights \
    -d '{"path": [["city2", "city3"], ["city1", "city2"], ["city3", "city4"]], "departureTime": 1670696785}'
    ```
- Example Response
    ```
    {}
    ```

##### GET `/v1/calculate/{personID}/flights`
- Response
    ```
    200: success
    500: internal server errror
    ```
- Example Request
    ```
    curl -X GET "localhost:8080/v1/calculate/kewei/flights"
    ```
- Example Response
    ```
    [
        {
            "path":["city1","city2","city3","city4"],
            "source":"city1",
            "destination":"city4",
            "departureTime":1670696785
        }
    ]
    ```

