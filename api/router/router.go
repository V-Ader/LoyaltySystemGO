package router

import (
	"fmt"

	"github.com/V-Ader/Loyality_GO/api/resource/card"

	"github.com/V-Ader/Loyality_GO/api/resource/client"
	"github.com/V-Ader/Loyality_GO/api/resource/issuer"

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
		clientGroup.PUT("/:id", client.Put(dbConnection))
		clientGroup.PATCH("/:id", client.Patch(dbConnection))
		clientGroup.DELETE("/:id", client.Delete(dbConnection))
	}

	issuerGroup := router.Group("issuers")
	{
		issuerGroup.GET("", issuer.GetAll(dbConnection))
		issuerGroup.GET("/:id", issuer.Get(dbConnection))
		issuerGroup.POST("/", issuer.Post(dbConnection))
		issuerGroup.PUT("/:id", issuer.Put(dbConnection))
		issuerGroup.PATCH("/:id", issuer.Patch(dbConnection))
		issuerGroup.DELETE("/:id", issuer.Delete(dbConnection))
	}

	cardGroup := router.Group("cards")
	{
		cardGroup.GET("", card.GetAll(dbConnection))
		cardGroup.GET("/:id", card.Get(dbConnection))
		cardGroup.POST("/", card.Post(dbConnection))
		cardGroup.PUT("/:id", card.Put(dbConnection))
		cardGroup.PATCH("/:id", card.Patch(dbConnection))
		cardGroup.DELETE("/:id", card.Delete(dbConnection))
	}

	tokenGroup := router.Group("tokens")
	{
		tokenGroup.POST("/clients", client.Token(dbConnection))
		tokenGroup.POST("/cards", card.Token(dbConnection))
		tokenGroup.POST("/issuers", issuer.Token(dbConnection))

	}

	return router
}
