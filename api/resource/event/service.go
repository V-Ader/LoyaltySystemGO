package event

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/V-Ader/Loyality_GO/api/resource/card"
	"github.com/V-Ader/Loyality_GO/api/resource/common"
	"github.com/gin-gonic/gin"
)

type EventService struct {
	transactionMutex    sync.Mutex
	CardServiceinstance *card.CardService
}

func (s *EventService) TransactionLock() {
	s.transactionMutex.Lock()
}

func (s *EventService) TransactionUnLock() {
	s.transactionMutex.Unlock()
}

func extractPagination(context *gin.Context) (int, int) {
	page, _ := strconv.Atoi(context.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(context.DefaultQuery("pageSize", "10"))
	return page, pageSize
}

func (s *EventService) ExecutGet(dbConnection *sql.DB, context *gin.Context) ([]common.Entity, *common.RequestError) {
	var query string
	var args []interface{}

	page, pageSize := extractPagination(context)

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = "SELECT * FROM events ORDER BY id LIMIT $1 OFFSET $2"
		args = []interface{}{pageSize, offset}
	} else {
		query = "SELECT * FROM events ORDER BY id"
	}

	results, err := dbConnection.Query(query, args...)
	if err != nil {
		return nil, &common.RequestError{StatusCode: http.StatusNotFound, Err: err}
	}
	defer results.Close()

	entities := []common.Entity{}
	for results.Next() {
		var event Event
		err = results.Scan(&event.Id, &event.Card_id, &event.Timestamp, &event.Quantity)
		if err != nil {
			return nil, &common.RequestError{StatusCode: http.StatusInternalServerError, Err: err}
		}
		entities = append(entities, &event)
	}
	return entities, nil
}

func (s *EventService) ExecutGetById(dbConnection *sql.DB, context *gin.Context) (common.Entity, *common.RequestError) {
	id := context.Param("id")
	query := "SELECT * FROM events WHERE id = $1"
	row := dbConnection.QueryRow(query, id)

	var event Event
	err := row.Scan(&event.Id, &event.Card_id, &event.Timestamp, &event.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &common.RequestError{StatusCode: http.StatusBadRequest, Err: fmt.Errorf("event not found")}
		}
		return nil, &common.RequestError{StatusCode: http.StatusInternalServerError, Err: err}
	}

	return &event, nil
}

func CreateTestContext(jsonBody string, paramID string) *gin.Context {
	req, _ := http.NewRequest(http.MethodPatch, "/", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(nil)
	c.Request = req

	c.Params = append(c.Params, gin.Param{Key: "id", Value: paramID})
	return c
}

func (s *EventService) ExecutePost(dbConnection *sql.DB, context *gin.Context) *common.RequestError {
	var eventData EventDataRequest

	if err := context.BindJSON(&eventData); err != nil {
		return &common.RequestError{StatusCode: http.StatusBadRequest, Err: err}
	}

	ctx := CreateTestContext("", fmt.Sprintf("%d", eventData.Card_id))
	s.CardServiceinstance.TransactionLock()
	defer s.CardServiceinstance.TransactionUnLock()

	entity, err := s.CardServiceinstance.ExecutGetById(dbConnection, ctx)
	if err != nil {
		return err
	}
	card, ok := entity.(*card.Card)
	if !ok {
		return &common.RequestError{StatusCode: http.StatusUnprocessableEntity, Err: fmt.Errorf("could not process the card of id %d", eventData.Card_id)}
	}

	if card.Tokens < eventData.Quantity {
		return &common.RequestError{StatusCode: http.StatusUnprocessableEntity, Err: fmt.Errorf("not enough tokens on card %d", eventData.Card_id)}
	}

	//recude quantity
	body := fmt.Sprintf("{\n\"tokens\": %d\n}", card.Tokens-eventData.Quantity)
	ctx = CreateTestContext(body, fmt.Sprintf("%d", eventData.Card_id))

	patchErr := s.CardServiceinstance.ExecutePatch(dbConnection, ctx)
	if patchErr != nil {
		return patchErr
	}

	//add event
	query := "INSERT INTO events (id, card_id, timestamp, quantity) VALUES (nextval('event_seq'), $1, $2, $3)"
	_, execErr := dbConnection.Exec(query, eventData.Card_id, time.Now(), eventData.Quantity)
	if execErr != nil {
		return &common.RequestError{StatusCode: http.StatusInternalServerError, Err: execErr}
	}
	return nil
}

func (s *EventService) ExecutePut(dbConnection *sql.DB, context *gin.Context) *common.RequestError {
	return &common.RequestError{StatusCode: http.StatusForbidden, Err: fmt.Errorf("PUT method is forbidden")}
}

func (s *EventService) ExecutePatch(dbConnection *sql.DB, context *gin.Context) *common.RequestError {
	return &common.RequestError{StatusCode: http.StatusForbidden, Err: fmt.Errorf("PUT method is forbidden")}
}

func (s *EventService) ExecuteDelete(dbConnection *sql.DB, context *gin.Context) *common.RequestError {
	query := "DELETE FROM events where id = $1"
	_, err := dbConnection.Exec(query, context.Param("id"))
	if err != nil {
		return &common.RequestError{StatusCode: http.StatusInternalServerError, Err: err}
	}
	return nil
}
