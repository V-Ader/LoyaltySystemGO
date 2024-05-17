package client

import (
	"database/sql"
	"fmt"

	"github.com/V-Ader/Loyality_GO/database"
	"github.com/gin-gonic/gin"
)

func GetUsers(dbConnection *sql.DB) ([]Client, error) {
	results, err := dbConnection.Query("SELECT * FROM clients")
	if err != nil {
		return nil, err
	}
	defer results.Close()

	clients := []Client{}
	for results.Next() {
		var client Client
		err = results.Scan(&client.Id, &client.Name, &client.Email)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return clients, nil
}

func GetClientById(dbConnection *sql.DB, context *gin.Context) (*Client, error) {
	id := context.Param("id")
	query := "SELECT id, name, email FROM clients WHERE id = $1"
	row := dbConnection.QueryRow(query, id)

	var client Client
	err := row.Scan(&client.Id, &client.Name, &client.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("client not found")
		}
		return nil, err
	}

	return &client, nil
}

func ExecutePost(dbConnection *sql.DB, context *gin.Context) error {
	query := "INSERT INTO clients (id, name, email) VALUES (nextval('client_seq'), $1, $2)"
	_, err := dbConnection.Exec(query, context.Query("name"), context.Query("email"))
	return err
}

func ExecutePut(dbConnection *sql.DB, context *gin.Context) error {
	id := context.Param("id")
	var clientUpdate ClientUpdateRequest

	if err := context.BindJSON(&clientUpdate); err != nil {
		return err
	}

	updates := map[string]interface{}{
		"name":  clientUpdate.Name,
		"email": clientUpdate.Email,
	}

	query, args := database.BuildUpsertQuery("clients", updates, id)

	_, err := dbConnection.Exec(query, args...)
	return err
}

func ExecutePatch(dbConnection *sql.DB, context *gin.Context) error {
	id := context.Param("id")
	var clientPatch ClientPatchRequest

	if err := context.BindJSON(&clientPatch); err != nil {
		return err
	}

	updates := map[string]interface{}{}
	if clientPatch.Name != nil {
		updates["name"] = *clientPatch.Name
	}
	if clientPatch.Email != nil {
		updates["email"] = *clientPatch.Email
	}

	if len(updates) == 0 {
		return fmt.Errorf("no fields provided for update")
	}

	query, args := database.BuildUpdateQuery("clients", updates, id)

	_, err := dbConnection.Exec(query, args...)
	return err
}

func ExecuteDelte(dbConnection *sql.DB, context *gin.Context) error {
	query := "DELETE FROM clients where id = $1"
	_, err := dbConnection.Exec(query, context.Param("id"))
	return err
}
