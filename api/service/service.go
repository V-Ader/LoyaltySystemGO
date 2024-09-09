package service

import (
	"database/sql"

	"github.com/V-Ader/Loyality_GO/api/resource/common"
	"github.com/gin-gonic/gin"
)

type Service interface {
	ExecutGet(dbConnection *sql.DB, context *gin.Context) ([]common.Entity, *common.RequestError)
	ExecutGetById(dbConnection *sql.DB, context *gin.Context) (common.Entity, *common.RequestError)
	ExecutePost(dbConnection *sql.DB, context *gin.Context) *common.RequestError
	ExecutePut(dbConnection *sql.DB, context *gin.Context) *common.RequestError
	ExecutePatch(dbConnection *sql.DB, context *gin.Context) *common.RequestError
	ExecuteDelete(dbConnection *sql.DB, context *gin.Context) *common.RequestError
	TransactionLock()
	TransactionUnLock()
}
