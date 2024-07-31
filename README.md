# Loyality_GO API

## About
Loyalty_GO is a RESTful API built with Go and the Gin framework. The API provides CRUD operations for managing clients, issuers, cards, events, and tokens in a loyalty program system. The project is structured to facilitate easy management and extension.

## Getting started
### Prerequisites
- Go 1.18+
- A running instance of your preferred database (configured in database package)
- Git
- docker
- docker-compose

### Instalation
1. Clone repository
```bash
git clone https://github.com/V-Ader/LoyaltySystemGO.git
cd LoyaltySystemGO
```
2. Adjust parameters in __docker-compose.yaml__

3. Use docker compose file
```bash
docker compose up --build
```

## Usage
You can import file __Loyality.postman_collection.json__ into Postman to use predefined test requests. 

By default service will be avaible on port __:8080__

# CRUD Table

## Clients

| HTTP Method | Endpoint     | Description               |
|-------------|--------------|---------------------------|
| GET         | /clients     | Get all clients           |
| GET         | /clients/:id | Get a specific client by ID |
| POST        | /clients/    | Create a new client       |
| PUT         | /clients/:id | Update a client by ID     |
| PATCH       | /clients/:id | Partially update a client by ID |
| DELETE      | /clients/:id | Delete a client by ID     |

## Issuers

| HTTP Method | Endpoint     | Description               |
|-------------|--------------|---------------------------|
| GET         | /issuers     | Get all issuers           |
| GET         | /issuers/:id | Get a specific Issuer by ID |
| POST        | /issuers/    | Create a new issuer       |
| PUT         | /issuers/:id | Update a issuer by ID     |
| PATCH       | /issuers/:id | Partially update a issuer by ID |
| DELETE      | /issuers/:id | Delete a issuer by ID     |

## Cards

| HTTP Method | Endpoint     | Description               |
|-------------|--------------|---------------------------|
| GET         | /cards     | Get all cards           |
| GET         | /cards/:id | Get a specific card by ID |
| POST        | /cards/    | Create a new card       |
| PUT         | /cards/:id | Update a card by ID     |
| PATCH       | /cards/:id | Partially update a card by ID |
| DELETE      | /cards/:id | Delete a card by ID     |


## Events

| HTTP Method | Endpoint     | Description               |
|-------------|--------------|---------------------------|
| GET         | /events     | Get all events           |
| GET         | /events/:id | Get a specific event by ID |
| POST        | /events/    | Create a new event       |
| DELETE      | /events/:id | Delete a event by ID     |


## Token
| HTTP Method | Endpoint     | Description               |
|-------------|--------------|---------------------------|
| POST         | /tokens     | create and get new token           |