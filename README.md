# Pizza Order App – Monorepo

This repository contains the **Pizza Order App** system, structured using a microservices architecture. It includes a main **API Service** for handling pizza orders, an **Email Service** for sending email notifications, and a **Search Service** for searching restaurants and pizzas based on location. The services communicate via asynchronous messaging.


## Services Overview

### `api-service` – Pizza Ordering API – Message Producer

- Handles customer orders via REST APIs.
- Publishes order events (e.g., `user.registered`, `order.placed`) to a message broker.
- Implements Domain-Driven Design architecture.

### `email-service` – Email Sending Service – Message Consumer

- Listens to order events from the broker.
- Sends confirmation emails to customers.
- A lightweight background service.

### `search-service` – Search API – Message Consumer

- Handles events from the broker (e.g., RabbitMQ) and indexes them.
- Exposes search API via Gin and Elasticsearch.
- Supports location-based and text-based search.


## Tech Stack

- **Language:** Go (Golang)
- **Database:** PostgreSQL
- **Search:** Elasticsearch
- **Architecture:** Domain-Driven Design Architecture
- **Messaging:** RabbitMQ


## Project Structure

```bash
pizza-order-app/
│── api-service/               # Main API Service (Producer)
│   ├── cmd/                   # Entry point
│   │   ├── main.go            # Starts HTTP API, publishes events
│   ├── internal/
│   │   ├── application/       # Application logic/use cases
│   │   ├── domain/            # Domain models and interfaces
│   │   ├── infrastructure/    # DB, messaging
│   │   ├── interfaces/        # API controllers
│
│── email-service/             # Email Sending Service (Consumer)
│   ├── cmd/                   # Entry point
│   │   ├── main.go            # Starts email consumer
│   ├── internal/
│   │   ├── application/       # Email handling logic
│   │   ├── domain/            # Email domain models and interfaces
│   │   ├── infrastructure/    # Email transport, message broker
│
│── search-service/            # Uses events to sync ES data (Consumer)
│   ├── cmd/                   # Entry point
│   │   ├── main.go            # Starts search API and consumer
│   ├── internal/
│   │   ├── application/       # Search handling logic
│   │   ├── domain/            # Search domain models and interfaces
│   │   ├── infrastructure/    # Elasticsearch adapter, event consumer
│
├── web-user/                  # Frontend UI for users (React)
│
│── .gitignore                 # .gitignore file
│── README.md                  # You are here
```


## Project Status

> **Note:** This project is **under active development**.  
> Some features may be incomplete or subject to change.  
> You're welcome to explore or provide feedback!
