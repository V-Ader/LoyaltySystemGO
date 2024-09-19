package client

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/V-Ader/Loyality_GO/api/resource/card"
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

func (s *ClientService) ExecutGet(dbConnection *sql.DB, context *gin.Context) ([]common.Entity, *common.RequestError) {
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
		return nil, &common.RequestError{StatusCode: http.StatusNotFound, Err: err}
	}
	defer results.Close()

	clients := []common.Entity{}
	for results.Next() {
		var client Client
		err = results.Scan(&client.Id, &client.Name, &client.Email)
		if err != nil {
			return nil, &common.RequestError{StatusCode: http.StatusInternalServerError, Err: err}
		}
		clients = append(clients, &client)
	}
	return clients, nil
}

func (s *ClientService) ExecutGetCards(dbConnection *sql.DB, context *gin.Context) ([]common.Entity, *common.RequestError) {
	var query string
	var args []interface{}

	page, pageSize := extractPagination(context)
	id := context.Param("id")

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = "SELECT * FROM cards WHERE owner_id = $1 ORDER BY id LIMIT $2 OFFSET $3"
		args = []interface{}{id, pageSize, offset}
	} else {
		query = "SELECT * FROM cards WHERE owner_id = $1 ORDER BY id"
		args = []interface{}{id}
	}

	results, err := dbConnection.Query(query, args...)
	if err != nil {
		return nil, &common.RequestError{StatusCode: http.StatusNotFound, Err: err}
	}
	defer results.Close()

	cards := []common.Entity{}
	for results.Next() {
		var card card.Card
		err = results.Scan(&card.Id, &card.Issuer_id, &card.Owner_id, &card.Active, &card.Tokens, &card.Capacity)
		if err != nil {
			return nil, &common.RequestError{StatusCode: http.StatusInternalServerError, Err: err}
		}
		cards = append(cards, &card)
	}
	return cards, nil
}

func (s *ClientService) ExecutGetById(dbConnection *sql.DB, context *gin.Context) (common.Entity, *common.RequestError) {
	id := context.Param("id")
	query := "SELECT id, name, email FROM clients WHERE id = $1"
	row := dbConnection.QueryRow(query, id)

	var client Client
	err := row.Scan(&client.Id, &client.Name, &client.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &common.RequestError{StatusCode: http.StatusBadRequest, Err: fmt.Errorf("client not found")}
		}
		return nil, &common.RequestError{StatusCode: http.StatusInternalServerError, Err: err}
	}

	return &client, nil
}

func (s *ClientService) ExecutePost(dbConnection *sql.DB, context *gin.Context) *common.RequestError {
	var clientData ClientDataRequest

	if err := context.BindJSON(&clientData); err != nil {
		return &common.RequestError{StatusCode: http.StatusBadRequest, Err: err}
	}

	query := "INSERT INTO clients (id, name, email) VALUES (nextval('client_seq'), $1, $2)"
	_, err := dbConnection.Exec(query, clientData.Name, clientData.Email)
	if err != nil {
		return &common.RequestError{StatusCode: http.StatusBadRequest, Err: err}
	}
	return nil
}

func (s *ClientService) ExecutePut(dbConnection *sql.DB, context *gin.Context) *common.RequestError {
	id := context.Param("id")
	var clientUpdate ClientDataRequest

	if err := context.BindJSON(&clientUpdate); err != nil {
		return &common.RequestError{StatusCode: http.StatusBadRequest, Err: err}
	}

	updates := map[string]interface{}{
		"name":  clientUpdate.Name,
		"email": clientUpdate.Email,
	}

	query, args := database.BuildUpsertQuery("clients", updates, id)

	_, err := dbConnection.Exec(query, args...)
	if err != nil {
		return &common.RequestError{StatusCode: http.StatusInternalServerError, Err: err}
	}
	return nil
}

func (s *ClientService) ExecutePatch(dbConnection *sql.DB, context *gin.Context) *common.RequestError {
	id := context.Param("id")
	var clientPatch ClientPatchRequest

	if err := context.BindJSON(&clientPatch); err != nil {
		return &common.RequestError{StatusCode: http.StatusBadRequest, Err: err}
	}

	updates := map[string]interface{}{}
	if clientPatch.Name != nil {
		updates["name"] = *clientPatch.Name
	}
	if clientPatch.Email != nil {
		updates["email"] = *clientPatch.Email
	}

	if len(updates) == 0 {
		return &common.RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("no fields provided for update")}
	}

	query, args := database.BuildUpdateQuery("clients", updates, id)

	_, err := dbConnection.Exec(query, args...)
	if err != nil {
		return &common.RequestError{StatusCode: http.StatusInternalServerError, Err: err}
	}
	return nil
}

func (s *ClientService) ExecuteDelete(dbConnection *sql.DB, context *gin.Context) *common.RequestError {
	query := "DELETE FROM clients where id = $1"
	_, err := dbConnection.Exec(query, context.Param("id"))
	if err != nil {
		return &common.RequestError{StatusCode: http.StatusInternalServerError, Err: err}
	}
	return nil
}
