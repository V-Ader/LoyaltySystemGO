package card

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/V-Ader/Loyality_GO/api/resource/common"
	"github.com/V-Ader/Loyality_GO/database"
	"github.com/gin-gonic/gin"
)

type CardService struct{}

func extractPagination(context *gin.Context) (int, int) {
	page, _ := strconv.Atoi(context.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(context.DefaultQuery("pageSize", "10"))
	return page, pageSize
}

func (s *CardService) ExecutGet(dbConnection *sql.DB, context *gin.Context) ([]common.Entity, error) {
	var query string
	var args []interface{}

	page, pageSize := extractPagination(context)

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = "SELECT * FROM cards ORDER BY id LIMIT $1 OFFSET $2"
		args = []interface{}{pageSize, offset}
	} else {
		query = "SELECT * FROM cards ORDER BY id"
	}

	results, err := dbConnection.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	cards := []common.Entity{}
	for results.Next() {
		var card Card
		err = results.Scan(&card.Id, &card.Issuer_id, &card.Owner_id, &card.Active, &card.Usages, &card.Capacity)
		if err != nil {
			return nil, err
		}
		cards = append(cards, &card)
	}
	return cards, nil
}

func (s *CardService) ExecutGetById(dbConnection *sql.DB, context *gin.Context) (common.Entity, error) {
	id := context.Param("id")
	query := "SELECT * FROM cards WHERE id = $1"
	row := dbConnection.QueryRow(query, id)

	var card Card
	err := row.Scan(&card.Id, &card.Issuer_id, &card.Owner_id, &card.Active, &card.Usages, &card.Capacity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("card not found")
		}
		return nil, err
	}

	return &card, nil
}

func (s *CardService) ExecutePost(dbConnection *sql.DB, context *gin.Context) error {
	var cardData CardDataRequest

	if err := context.BindJSON(&cardData); err != nil {
		return err
	}

	query := "INSERT INTO cards (id, issuer_id, owner_id, active, usages, capacity) VALUES (nextval('card_seq'), $1, $2, $3, $4, $5)"
	_, err := dbConnection.Exec(query, cardData.Issuer_id, cardData.Owner_id, cardData.Active, cardData.Usages, cardData.Capacity)
	return err
}

func (s *CardService) ExecutePut(dbConnection *sql.DB, context *gin.Context) error {
	id := context.Param("id")
	var cardUpdate CardDataRequest

	if err := context.BindJSON(&cardUpdate); err != nil {
		return err
	}

	updates := map[string]interface{}{
		"issuer_id": cardUpdate.Issuer_id,
		"owner_id":  cardUpdate.Owner_id,
		"active":    cardUpdate.Active,
		"usages":    cardUpdate.Usages,
		"capacity":  cardUpdate.Capacity,
	}

	query, args := database.BuildUpsertQuery("cards", updates, id)

	_, err := dbConnection.Exec(query, args...)
	return err
}

func (s *CardService) ExecutePatch(dbConnection *sql.DB, context *gin.Context) error {
	id := context.Param("id")
	var cardPatch CardPatchRequest

	if err := context.BindJSON(&cardPatch); err != nil {
		return err
	}

	updates := map[string]interface{}{}
	if cardPatch.Issuer_id != nil {
		updates["issuer_id"] = *cardPatch.Issuer_id
	}
	if cardPatch.Owner_id != nil {
		updates["owner_id"] = *cardPatch.Owner_id
	}
	if cardPatch.Active != nil {
		updates["active"] = *cardPatch.Active
	}
	if cardPatch.Usages != nil {
		updates["usages"] = *cardPatch.Usages
	}
	if cardPatch.Capacity != nil {
		updates["capacity"] = *cardPatch.Capacity
	}

	if len(updates) == 0 {
		return errors.New("no fields provided for update")
	}

	query, args := database.BuildUpdateQuery("cards", updates, id)

	_, err := dbConnection.Exec(query, args...)
	return err
}

func (s *CardService) ExecuteDelete(dbConnection *sql.DB, context *gin.Context) error {
	query := "DELETE FROM cards where id = $1"
	_, err := dbConnection.Exec(query, context.Param("id"))
	return err
}
