package client

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/V-Ader/Loyality_GO/api/resource/common"
	"github.com/V-Ader/Loyality_GO/database"
	"github.com/gin-gonic/gin"
)

type ClientService struct {
	transactionMutex sync.Mutex
}

func (s *ClientService) TransactionLock() {
	s.transactionMutex.Lock()
}

func (s *ClientService) TransactionUnLock() {
	s.transactionMutex.Unlock()
}

func extractPagination(context *gin.Context) (int, int) {
	page, _ := strconv.Atoi(context.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(context.DefaultQuery("pageSize", "10"))
	return page, pageSize
}

func (s *ClientService) ExecutGet(dbConnection *sql.DB, context *gin.Context) ([]common.Entity, error) {
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

	clients := []common.Entity{}
	for results.Next() {
		var client Client
		err = results.Scan(&client.Id, &client.Name, &client.Email)
		if err != nil {
			return nil, err
		}
		clients = append(clients, &client)
	}
	return clients, nil
}

func (s *ClientService) ExecutGetById(dbConnection *sql.DB, context *gin.Context) (common.Entity, error) {
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

func (s *ClientService) ExecutePost(dbConnection *sql.DB, context *gin.Context) error {
	var clientData ClientDataRequest

	if err := context.BindJSON(&clientData); err != nil {
		return err
	}

	query := "INSERT INTO clients (id, name, email) VALUES (nextval('client_seq'), $1, $2)"
	_, err := dbConnection.Exec(query, clientData.Name, clientData.Email)
	return err
}

func (s *ClientService) ExecutePut(dbConnection *sql.DB, context *gin.Context) error {
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

func (s *ClientService) ExecutePatch(dbConnection *sql.DB, context *gin.Context) error {
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

func (s *ClientService) ExecuteDelete(dbConnection *sql.DB, context *gin.Context) error {
	query := "DELETE FROM clients where id = $1"
	_, err := dbConnection.Exec(query, context.Param("id"))
	return err
}
