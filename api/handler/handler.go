package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/V-Ader/Loyality_GO/api/resource/cache"
	"github.com/V-Ader/Loyality_GO/api/resource/common"
	"github.com/V-Ader/Loyality_GO/api/resource/response"
	"github.com/V-Ader/Loyality_GO/api/service"

	"github.com/gin-gonic/gin"
)

var (
	tokenCache *cache.TokenCache
)

func init() {
	tokenCache = cache.NewTokenCache(5*time.Minute, 10*time.Minute)
}

func Token(dbConnection *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		token, err := tokenCache.CreateToken()
		if err != nil {
			context.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		} else {
			context.JSON(http.StatusOK, response.Response{Data: token})
		}
	}
}

type dbFunctionProcessing func(*sql.DB, *gin.Context) ([]common.Entity, *common.RequestError)

func Execute(executeFunction dbFunctionProcessing, dbConnection *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		users, err := executeFunction(dbConnection, context)
		if err != nil {
			context.JSON(err.StatusCode, response.ErrorResponse{Message: err.Err.Error()})
		} else {
			context.JSON(http.StatusOK, response.Response{Data: users})
		}
	}
}

func GetAll(service service.Service, dbConnection *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		users, err := service.ExecutGet(dbConnection, context)
		if err != nil {
			context.JSON(err.StatusCode, response.ErrorResponse{Message: err.Err.Error()})
		} else {
			context.JSON(http.StatusOK, response.Response{Data: users})
		}
	}
}

func Get(service service.Service, dbConnection *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		result, err := service.ExecutGetById(dbConnection, context)
		if err != nil {
			context.JSON(err.StatusCode, response.ErrorResponse{Message: err.Err.Error()})
		} else {
			context.Header("ETag", result.GetHash())
			context.JSON(http.StatusOK, response.Response{Data: &result})
		}
	}
}

func Post(service service.Service, dbConnection *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		err := tokenCache.RemoveToken(context.Query("token"))
		if err != nil {
			context.JSON(err.StatusCode, response.ErrorResponse{Message: err.Err.Error()})
			return
		}
		err = service.ExecutePost(dbConnection, context)
		if err != nil {
			context.JSON(err.StatusCode, response.ErrorResponse{Message: err.Err.Error()})
		} else {
			context.Status(http.StatusCreated)
		}
	}
}

func CheckIfMatch(context *gin.Context, entity common.Entity) bool {
	ifMatchCondition := context.GetHeader("If-Match")
	return ifMatchCondition == entity.GetHash()
}

func Put(service service.Service, dbConnection *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		service.TransactionLock()
		defer service.TransactionUnLock()

		entity, err := service.ExecutGetById(dbConnection, context)
		if err != nil {
			context.JSON(err.StatusCode, response.ErrorResponse{Message: err.Err.Error()})
			return
		}

		if !CheckIfMatch(context, entity) {
			context.JSON(http.StatusPreconditionFailed, response.ErrorResponse{Message: "If-Match does not match ETag"})
			return
		}

		err = service.ExecutePut(dbConnection, context)
		if err != nil {
			context.JSON(err.StatusCode, response.ErrorResponse{Message: err.Err.Error()})
			return
		}

		context.Status(http.StatusNoContent)
	}
}

func Patch(service service.Service, dbConnection *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		service.TransactionLock()
		defer service.TransactionUnLock()

		err := service.ExecutePatch(dbConnection, context)
		if err != nil {
			context.IndentedJSON(err.StatusCode, response.ErrorResponse{Message: err.Err.Error()})
		} else {
			context.Status(http.StatusNoContent)
		}
	}
}

func Delete(service service.Service, dbConnection *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		err := service.ExecuteDelete(dbConnection, context)
		if err != nil {
			context.IndentedJSON(err.StatusCode, response.ErrorResponse{Message: err.Err.Error()})
		} else {
			context.Status(http.StatusOK)
		}
	}
}
