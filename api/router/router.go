package router

import (
	"fmt"

	"github.com/V-Ader/Loyality_GO/api/resource/client"

	db "github.com/V-Ader/Loyality_GO/database"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	router := gin.Default()
	if router == nil {
		fmt.Println("Failed to create router")
		return nil
	}

	dbConnection, err := db.Connect()
	if err != nil {
		fmt.Println("Failed to connect with database", err)
		return nil
	}

	clientGroup := router.Group("clients")
	{
		clientGroup.GET("", client.GetAll(dbConnection))
		clientGroup.GET("/:id", client.Get(dbConnection))
		clientGroup.POST("/", client.Post(dbConnection))
		// clientGroup.PUT("/:id", client.Put(dbConnection))
		// clientGroup.PATCH("/:id", client.Patch(dbConnection))
		clientGroup.DELETE("/:id", client.Delete(dbConnection))
	}

	return router
}
