# Pizza Order App – Monorepo

This repository contains the **Pizza Order App** system, structured as a microservices architecture. It includes a main **API Service** for handling pizza orders, an **Email Service** for sending email notifications, and **Search Service** for searching restaurants and pizzas based on location. The services communicate via asynchronous messaging.


## Services Overview

### `api-service` - Pizza Ordering API – Message Producer
- Handles customer orders via REST API.
- Publishes order events (e.g., `user.registered`, `order.placed`) to a message broker (e.g. RabbitMQ).
- Implements Clean/Domain-Driven Design Architecture.

### `email-service` - Email Sending Service – Message Consumer
- Listens to order events from the broker.
- Sends confirmation emails to customers.
- Lightweight background service.

### `search-service` - Search API – Message Consumer
- Handles events from the broker (RabbitMQ), indexes them.
- Exposes search API via Gin and Elasticsearch.
- Supports location-based and text search.

## Tech Stack

- **Language:** Go (Golang)
- **Database:** PostgreSQL
- **Search:** Elasticsearch
- **Architecture:** Clean Architecture
- **Messaging:** RabbitMQ


## Project Structure

```bash
pizza-order-app/
│── api-service/               # Main API Service (Producer)
│   ├── cmd/                   # Entrypoint
│   │   ├── main.go            # Starts HTTP API, Publishes Events
│   ├── internal/
│   │   ├── application/       # Application logic/use-cases
│   │   ├── domain/            # Domain models and interfaces
│   │   ├── infrastructure/    # DB, Messaging
│   │   ├── interfaces/        # API controllers
│
│── email-service/             # Email Sending Service (Consumer)
│   ├── cmd/                   # Entrypoint
│   │   ├── main.go            # Starts Email Consumer
│   ├── internal/
│   │   ├── application/       # Email handling logic
│   │   ├── domain/            # Email domain models and interfaces
│   │   ├── infrastructure/    # Email transport, message broker
│
│── search-service/            # Use events to sync ES data (Search API and Consumer)
│   ├── cmd/                   # Entrypoint
│   │   ├── main.go            # Starts Search API and Consumer
│   ├── internal/
│   │   ├── application/       # Search handling logic
│   │   ├── domain/            # Search domain models and interfaces
│   │   ├── infrastructure/    # Elasticsearch adapter, Event consumer
│
├── web-user/                  # Frontend UI for Users (React)
│
│── .gitignore                 # .gitignore file
│── README.md                  # You are here ✨
```


## Project Status

> **Note:** This project is **under active development**.
> Some features may be incomplete or subject to change.
> You're welcome to explore or give feedback!
