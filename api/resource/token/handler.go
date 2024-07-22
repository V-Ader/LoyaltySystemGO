package token

import (
	"database/sql"
	"net/http"

	"github.com/V-Ader/Loyality_GO/api/resource/response"

	"github.com/gin-gonic/gin"
)

func Post(dbConnection *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		val, err := ExecutePost(dbConnection, context)
		if err != nil {
			context.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		} else {
			context.JSON(http.StatusOK, TokenResponse{Token: val})
		}
	}
}

func Delete(dbConnection *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		err := ExecuteDelete(dbConnection, context)
		if err != nil {
			context.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		} else {
			context.Status(http.StatusOK)
		}
	}
}
