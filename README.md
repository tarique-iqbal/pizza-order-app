# Pizza Order App - Monorepo

This repository contains the **Pizza Order App** system, structured as a microservices architecture. It includes a main **API Service** for handling pizza orders and an **Email Service** for sending email notifications. The services communicate via asynchronous messaging.


## Services Overview

### `api-service` - Pizza Ordering API (Producer)
- Handles customer orders via REST API.
- Publishes order events (e.g., `user.registered`, `order.placed`) to a message broker (e.g. RabbitMQ).
- Implements Clean/DDD Architecture.

### `email-service` - Email Notification Service (Consumer)
- Listens to order events from the broker.
- Sends confirmation emails to customers.
- Lightweight background service.


## Tech Stack

- **Language:** Go (Golang)
- **Database:** PostgreSQL
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
│   │   ├── domain/            # Email domain models/interfaces
│   │   ├── infrastructure/    # Email transport, message broker
│
├── user-client/               # Frontend UI for Users (React)
│
│── .gitignore                 # .gitignore file
│── README.md                  # You are here ✨
```


## Project Status

> **Note:** This project is **under active development**.  
> Some features may be incomplete or subject to change.  
> You're welcome to explore or give feedback!
