package client

import (
	"database/sql"
	"net/http"

	"github.com/V-Ader/Loyality_GO/api/resource/response"

	"github.com/gin-gonic/gin"
)

func GetAll(dbConnection *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		users, err := GetUsers(dbConnection)
		if err != nil {
			context.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		} else {
			context.JSON(http.StatusOK, response.ClientResponse{Data: users})
		}
	}
}

func Get(dbConnection *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		result, err := GetClientById(dbConnection, context)
		if err != nil {
			if err.Error() == "record not found" {
				context.JSON(http.StatusNotFound, response.ErrorResponse{Message: err.Error()})
			} else {
				context.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
			}
		} else {
			context.JSON(http.StatusAccepted, response.ClientResponse{Data: &result})
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

// func Put(dbConnection *sql.DB) gin.HandlerFunc {
// 	return func(context *gin.Context) {
// 		query := "INSERT INTO users (id, name) VALUES (nextval('user_seq'), $1)"
// 		_, err := dbConnection.Exec(query, context.Query("name"))
// 		if err != nil {
// 			context.IndentedJSON(http.StatusInternalServerError, response.ClientResponse{Error: err.Error()})
// 		} else {
// 			context.IndentedJSON(http.StatusOK, response.ClientResponse{Data: "User added correctly"})
// 		}
// 	}
// }

// func Patch(dbConnection *sql.DB) gin.HandlerFunc {
// 	return func(context *gin.Context) {
// 		query := "INSERT INTO users (id, name) VALUES (nextval('user_seq'), $1)"
// 		_, err := dbConnection.Exec(query, context.Query("name"))
// 		if err != nil {
// 			context.IndentedJSON(http.StatusInternalServerError, response.ClientResponse{Error: err.Error()})
// 		} else {
// 			context.IndentedJSON(http.StatusOK, response.ClientResponse{Data: "User added correctly"})
// 		}
// 	}
// }

func Delete(dbConnection *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		err := ExecuteDelte(dbConnection, context)
		if err != nil {
			context.IndentedJSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		} else {
			context.Status(http.StatusOK)
		}
	}
}
