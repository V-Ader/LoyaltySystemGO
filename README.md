## Loyality_GO

# About

# Usage

# CRUD Routes

## Clients Group

| HTTP Method | Endpoint     | Handler Function          | Description               |
|-------------|--------------|---------------------------|---------------------------|
| GET         | /clients     | client.GetAll(dbConnection) | Get all clients           |
| GET         | /clients/:id | client.Get(dbConnection)     | Get a specific client by ID |
| POST        | /clients/    | client.Post(dbConnection)    | Create a new client       |
| PUT         | /clients/:id | client.Put(dbConnection)     | Update a client by ID     |
| PATCH       | /clients/:id | client.Patch(dbConnection)   | Partially update a client by ID |
| DELETE      | /clients/:id | client.Delete(dbConnection)  | Delete a client by ID     |

## Issuers Group

| HTTP Method | Endpoint     | Handler Function          | Description               |
|-------------|--------------|---------------------------|---------------------------|
| GET         | /issuers     | client.GetAll(dbConnection) | Get all issuers           |
| GET         | /issuers/:id | client.Get(dbConnection)     | Get a specific issuer by ID |
| POST        | /issuers/    | client.Post(dbConnection)    | Create a new issuer       |
| PUT         | /issuers/:id | client.Put(dbConnection)     | Update an issuer by ID     |
| PATCH       | /issuers/:id | client.Patch(dbConnection)   | Partially update an issuer by ID |
| DELETE      | /issuers/:id | client.Delete(dbConnection)  | Delete an issuer by ID     |


