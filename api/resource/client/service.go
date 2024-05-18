package client

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/V-Ader/Loyality_GO/api/resource/cache"
	"github.com/V-Ader/Loyality_GO/database"
	"github.com/gin-gonic/gin"
)

var tokenCache *cache.TokenCache

func init() {
	tokenCache = cache.NewTokenCache(5*time.Minute, 10*time.Minute)
}

func extractPagination(context *gin.Context) (int, int) {
	pageStr := context.Query("page")
	pageSizeStr := context.Query("pageSize")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 0
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 0 {
		pageSize = 0
	}

	return page, pageSize
}

func GetUsers(dbConnection *sql.DB, context *gin.Context) ([]Client, error) {
	var query string
	var args []interface{}

	page, pageSize := extractPagination(context)

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = "SELECT * FROM clients ORDER BY id LIMIT $1 OFFSET $2"
		args = []interface{}{pageSize, offset}
	} else {
		query = "SELECT * FROM clients ORDER BY id"
	}

	results, err := dbConnection.Query(query, args...)
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
	deduplicationToken := context.Query("deduplicationToken")

	if err := tokenCache.ProcessToken(deduplicationToken); err != nil {
		return err
	}

	var clientData ClientDataRequest

	if err := context.BindJSON(&clientData); err != nil {
		return err
	}

	query := "INSERT INTO clients (id, name, email) VALUES (nextval('client_seq'), $1, $2)"
	_, err := dbConnection.Exec(query, clientData.Name, clientData.Email)
	return err
}

func ExecutePut(dbConnection *sql.DB, context *gin.Context) error {
	id := context.Param("id")
	var clientUpdate ClientDataRequest

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
		return errors.New("no fields provided for update")
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
