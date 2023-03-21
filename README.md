# User management service

## Instructions

Assuming you are on the project directory and Docker is installed.

### Docker

For Apple M1 users :
`sudo docker buildx build -t user-management-service:latest --platform=linux/arm64 . --load`
For **amd64** we build the container image for our microservice
`DOCKER_BUILDKIT=0 docker build -t user-management-service:latest .`.
After docker building process is completed succesfully,
`docker compose up`
optionally (if you want to run it as background daemon).
`docker compose up -d` .
Then **user-management-service** should be ready to make CRUD operations on it.

### Local

If `export MONGO_URL=<mongo_url:mongo_port>`set, program will use it. Otherwise `mongodb://127.0.0.1:27017` is the default url.
Additionally **MUST** :
`export KAFKA_URL=<kafka_url>:<kafka_port>`
To run the application locally first make sure that mongo instance is running on the system.
Use `go build cmd/main.go -o user-management-service`.
Finally application is ready to be used with `./user-management-service`.

### Why Hexagonal Architecture?

By implementing hexagonal architecture basic API functionality could easily be divided into _adapter_,_application_,_core_ layers. By doing so we could abstract each layer from another using interfaces like contracts. Since protobuf messages structures should not be sent directly to the database _model_ approach has been used to transform and manipulate data and vice a versa.

### Why MongoDB?

First of all great documentation. Since this is the first time I used MongoDB, documentation has huge effect on database selection. It is also NoSQL which makes it easy to generate data and play with it.

### Why Kafka?

Kafka is used to produce events and let other interested servers notified about the changes. Since kafka is well documented, easily built up I preferred Kafka to uses broker.

## Improvements

### Scaling

From vertical scaling perspective increasing the machine specs would help which application is running. From the horizontal perspective a task queue which filled by client request and workers listening on that queue would help us on making concurrent jobs in the application.

### Deployment

An automated CI/CD pipeline to run tests and deployment could be added. After deployment monitoring tools could be added to collect, analyze and debug services.
