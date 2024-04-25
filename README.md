
# Go REST API with Redis Caching
This project consists of a REST API in Go that utilizes Redis for caching and PostgreSQL for data storage. The API retrieves product information by ID, leveraging Redis to minimize database queries.

## Setup and Operation
* Database Configuration: Initialized a PostgreSQL database, created a products table, and populated it with sample data.
* Redis Integration: Configured a Redis instance in a Docker container, available for caching operations.
* API Implementation: Developed a Go application with the Gin framework to serve HTTP requests, interfacing with Redis and PostgreSQL.
* Caching Mechanism: Incorporated logic to check the Redis cache before querying the database, caching new queries with a TTL of 15 seconds.
* Testing and Verification: Manual testing via tools like curl or Postman confirms the API's functionality and effective caching.
* Execute go run main.go to start the server and access the API at localhost:8080/products/:id.

![image](https://github.com/yessetassan/GoAss3/assets/139701904/66640719-19ca-41fc-a7b2-946177073e65)

