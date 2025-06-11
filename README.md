# 📦 Collections API System

A modular API platform built with **Golang** that replicates the behavior of Postman collections.  
It consists of two main components:

- **Collections API Gateway** – GraphQL-based entrypoint to interact with collection data.
- **Collection Service** – gRPC backend that performs all business logic and data persistence.


## 🧱 Architecture Overview

- **Language**: Go (Golang)
- **Gateway**: GraphQL + gRPC Client
- **Service**: gRPC Server
- **Data Layer**: PostgreSQL via GORM
- **Logging**: Zerolog (structured and performant)
- **Communication**: Protocol Buffers over gRPC


## 🎯 Key Features

- Create and manage API request collections (like Postman)
- Organize requests into collections
- Update or delete collections and individual requests
- Query all stored collections and their nested requests


## 🔐 gRPC Service Methods

| Method                           | Description                                |
|----------------------------------|--------------------------------------------|
| `CreateCollection`              | Creates a new API collection               |
| `AddRequestToCollection`        | Adds a request to a specific collection    |
| `ListCollections`    | Lists all collections with their requests  |
| `UpdateCollection`              | Updates collection metadata (e.g., name)   |
| `UpdateRequestInCollection`     | Updates a specific request inside a collection |
| `DeleteCollection`              | Deletes a full collection                  |
| `DeleteRequestFromCollection`   | Deletes a single request from a collection |


## 🚀 Running the System

```bash

###  Clone the repository

git clone https://github.com/your-username/collections-api-system.git
cd collections-api-system

###  Run the servers
go run .\cmd\server\main.go
