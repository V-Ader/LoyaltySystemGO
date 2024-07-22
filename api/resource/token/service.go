package token

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
)

func ExecutePost(dbConnection *sql.DB, context *gin.Context) (string, error) {
	query := "INSERT INTO tokens (id) VALUES (nextval('token_seq'))"
	result, err := dbConnection.Exec(query)

	if err != nil {
		return "", err
	}
	id, err := result.LastInsertId()
	return fmt.Sprint(id), err
}

func ExecuteDelete(dbConnection *sql.DB, context *gin.Context) error {
	query := "DELETE FROM tokens where id = $1"
	_, err := dbConnection.Exec(query, context.Param("id"))
	return err
}
