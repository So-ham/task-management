# Task-management

## Problem Breakdown & Design Decisions

Build a Go microservice for managing tasks with full CRUD (Create, Read, Update, Delete) operations.  
Key features include pagination and filtering by task status on the `GET /tasks` endpoint.  
The design follows microservice principles with clear separation of concerns between API handlers, business logic, and data layers, adhering to the Single Responsibility Principle.  

The service is stateless and designed to scale horizontally, making it suitable for deployment behind load balancers or container orchestrators.

## How to Run the Service

### Prerequisites  
- Go 1.20+ installed

### Steps  
1. Clone the repository:
2. 2. Download dependencies and run the app:  
3. The service will be accessible at:   `http://localhost:8080`

### Documentation
- API Documenation is under `docs/swagger.yaml`


### Microservices Concepts Demonstrated

- **Single Responsibility Principle:** Clear separation between API handlers, business logic (service layer), and data access (model).  
- **Stateless Design & Scalability:** Service can be scaled horizontally by running multiple instances behind a load balancer.  
- **Inter-Service Communication:** Potential communication with other microservices (e.g., User Service) through REST, gRPC, or asynchronous message queues, enabling loosely coupled services and event-driven workflows based on the feature requirement.


