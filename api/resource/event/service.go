package event

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/V-Ader/Loyality_GO/api/resource/common"
	"github.com/gin-gonic/gin"
)

type EventService struct{}

func extractPagination(context *gin.Context) (int, int) {
	page, _ := strconv.Atoi(context.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(context.DefaultQuery("pageSize", "10"))
	return page, pageSize
}

func (s *EventService) ExecutGet(dbConnection *sql.DB, context *gin.Context) ([]common.Entity, error) {
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
		return nil, err
	}
	defer results.Close()

	entities := []common.Entity{}
	for results.Next() {
		var event Event
		err = results.Scan(&event.Id, &event.Card_id, &event.Timestamp, &event.Quantity)
		if err != nil {
			return nil, err
		}
		entities = append(entities, &event)
	}
	return entities, nil
}

func (s *EventService) ExecutGetById(dbConnection *sql.DB, context *gin.Context) (common.Entity, error) {
	id := context.Param("id")
	query := "SELECT * FROM events WHERE id = $1"
	row := dbConnection.QueryRow(query, id)

	var event Event
	err := row.Scan(&event.Id, &event.Card_id, &event.Timestamp, &event.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("event not found")
		}
		return nil, err
	}

	return &event, nil
}

func (s *EventService) ExecutePost(dbConnection *sql.DB, context *gin.Context) error {
	var eventData EventDataRequest

	if err := context.BindJSON(&eventData); err != nil {
		return err
	}

	query := "INSERT INTO events (id, card_id, timestamp, quantity) VALUES (nextval('event_seq'), $1, $2, $3)"
	_, err := dbConnection.Exec(query, eventData.Card_id, time.Now(), eventData.Quantity)
	return err
}

func (s *EventService) ExecutePut(dbConnection *sql.DB, context *gin.Context) error {
	return fmt.Errorf("PUT method is forbidden")
}

func (s *EventService) ExecutePatch(dbConnection *sql.DB, context *gin.Context) error {
	return fmt.Errorf("PUT method is forbidden")
}

func (s *EventService) ExecuteDelete(dbConnection *sql.DB, context *gin.Context) error {
	query := "DELETE FROM events where id = $1"
	_, err := dbConnection.Exec(query, context.Param("id"))
	return err
}
