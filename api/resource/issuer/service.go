package issuer

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

func ExecutGet(dbConnection *sql.DB, context *gin.Context) ([]Issuer, error) {
	var query string
	var args []interface{}

	page, pageSize := extractPagination(context)

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = "SELECT * FROM issuers ORDER BY id LIMIT $1 OFFSET $2"
		args = []interface{}{pageSize, offset}
	} else {
		query = "SELECT * FROM issuers ORDER BY id"
	}

	results, err := dbConnection.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	issuers := []Issuer{}
	for results.Next() {
		var issuer Issuer
		err = results.Scan(&issuer.Id, &issuer.Name)
		if err != nil {
			return nil, err
		}
		issuers = append(issuers, issuer)
	}
	return issuers, nil
}

func ExecutGetById(dbConnection *sql.DB, context *gin.Context) (*Issuer, error) {
	id := context.Param("id")
	query := "SELECT id, name FROM issuers WHERE id = $1"
	row := dbConnection.QueryRow(query, id)

	var issuer Issuer
	err := row.Scan(&issuer.Id, &issuer.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("issuer not found")
		}
		return nil, err
	}

	return &issuer, nil
}

func ExecutePost(dbConnection *sql.DB, context *gin.Context) error {
	deduplicationToken := context.Query("deduplicationToken")

	if err := tokenCache.ProcessToken(deduplicationToken); err != nil {
		return err
	}

	var issuerData IssuerDataRequest

	if err := context.BindJSON(&issuerData); err != nil {
		return err
	}

	query := "INSERT INTO issuers (id, name) VALUES (nextval('issuer_seq'), $1)"
	_, err := dbConnection.Exec(query, issuerData.Name)
	return err
}

func ExecutePut(dbConnection *sql.DB, context *gin.Context) error {
	id := context.Param("id")
	var issuerData IssuerDataRequest

	if err := context.BindJSON(&issuerData); err != nil {
		return err
	}

	updates := map[string]interface{}{
		"name": issuerData.Name,
	}

	query, args := database.BuildUpsertQuery("issuers", updates, id)

	_, err := dbConnection.Exec(query, args...)
	return err
}

func ExecutePatch(dbConnection *sql.DB, context *gin.Context) error {
	id := context.Param("id")
	var issuerData IssuerPatchRequest

	if err := context.BindJSON(&issuerData); err != nil {
		return err
	}

	updates := map[string]interface{}{}
	if issuerData.Name != nil {
		updates["name"] = *issuerData.Name
	}

	if len(updates) == 0 {
		return errors.New("no fields provided for update")
	}

	query, args := database.BuildUpdateQuery("issuers", updates, id)

	_, err := dbConnection.Exec(query, args...)
	return err
}

func ExecuteDelete(dbConnection *sql.DB, context *gin.Context) error {
	query := "DELETE FROM issuers where id = $1"
	_, err := dbConnection.Exec(query, context.Param("id"))
	return err
}
