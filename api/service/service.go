package service

import (
	"database/sql"

	"github.com/V-Ader/Loyality_GO/api/resource/common"
	"github.com/gin-gonic/gin"
)

type Service interface {
	ExecutGet(dbConnection *sql.DB, context *gin.Context) ([]common.Entity, error)
	ExecutGetById(dbConnection *sql.DB, context *gin.Context) (common.Entity, error)
	ExecutePost(dbConnection *sql.DB, context *gin.Context) error
	ExecutePut(dbConnection *sql.DB, context *gin.Context) error
	ExecutePatch(dbConnection *sql.DB, context *gin.Context) error
	ExecuteDelete(dbConnection *sql.DB, context *gin.Context) error
}
