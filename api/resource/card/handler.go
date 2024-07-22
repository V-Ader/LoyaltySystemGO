package card

import (
	"database/sql"
	"net/http"

	"github.com/V-Ader/Loyality_GO/api/resource/response"

	"github.com/gin-gonic/gin"
)

func GetAll(dbConnection *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		users, err := ExecutGet(dbConnection, context)
		if err != nil {
			context.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		} else {
			context.JSON(http.StatusOK, CardResponse{Data: users})
		}
	}
}

func Get(dbConnection *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		result, err := ExecutGetById(dbConnection, context)
		if err != nil {
			if err.Error() == "record not found" {
				context.JSON(http.StatusNotFound, response.ErrorResponse{Message: err.Error()})
			} else {
				context.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
			}
		} else {
			context.Header("ETag", result.getHash())
			context.JSON(http.StatusOK, CardResponse{Data: &result})
		}
	}
}

func Post(dbConnection *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		err := ExecutePost(dbConnection, context)
		if err != nil {
			context.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		} else {
			context.Status(http.StatusOK)
		}
	}
}

func CheckIfMatch(context *gin.Context, entity *Card) bool {
	ifMatchCondition := context.GetHeader("If-Match")
	return ifMatchCondition == entity.getHash()
}

func Put(dbConnection *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		entity, err := ExecutGetById(dbConnection, context)
		if err != nil {
			context.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
			return
		}

		if !CheckIfMatch(context, entity) {
			context.JSON(http.StatusPreconditionFailed, response.ErrorResponse{Message: "If-Match does not match ETag"})
			return
		}

		err = ExecutePut(dbConnection, context)
		if err != nil {
			context.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
			return
		}

		context.Status(http.StatusOK)
	}
}

func Patch(dbConnection *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		err := ExecutePatch(dbConnection, context)
		if err != nil {
			context.IndentedJSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		} else {
			context.Status(http.StatusOK)
		}
	}
}

func Delete(dbConnection *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		err := ExecuteDelete(dbConnection, context)
		if err != nil {
			context.IndentedJSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		} else {
			context.Status(http.StatusOK)
		}
	}
}
