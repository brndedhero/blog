# Blog
  
This is an implementation of the backend of a blog service written in Go 1.17 created with a microservice architecture in mind. Currently implemented is a simple model for blog posts using MySQL as a DB backend with Redis caching.  
  
This has been developed with containers in mind and is configured using environment variables. To use deploy this in a Go 1.17 container and pass the following environment variables;  
  
|Name|Value|
|-|-|
|DB_HOST|Hostname for MySQL Database (string)|
|DB_NAME|Database name for MySQL Database (string)|
|DB_PASSWORD|Password for MySQL Database (string)|
|DB_USER|Username for MySQL Database (string)|
|HTTP_HOST|Hostname to listen on (string)|
|HTTP_PORT|Port to listen on (int)|
|REDIS_DB|Redis database to use (int)|
|REDIS_HOST|Hostname for Redis cluster (string)|
|REDIS_PASSWORD|Password for Redis cluster (string)|

## To do

### Tags
- [ ] Implement Tags model  

### Search functionality
- [ ] Elasticsearch integration  

### Authentication
- [ ] Implement authentication 

### Logging
- [ ] Loki integration

## Done
- [x] Implement blog post model  
- [x] Have Redis caching  
- [x] Logging as JSON format for future Loki integration  
- [x] Prometheus metrics exporter  
