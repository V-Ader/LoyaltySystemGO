package router

import (
	"fmt"

	"github.com/V-Ader/Loyality_GO/api/handler"
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

	clientService := &client.ClientService{}
	clientGroup := router.Group("clients")
	{
		clientGroup.GET("", handler.GetAll(clientService, dbConnection))
		clientGroup.GET("/:id", handler.Get(clientService, dbConnection))
		clientGroup.POST("/", handler.Post(clientService, dbConnection))
		clientGroup.PUT("/:id", handler.Put(clientService, dbConnection))
		clientGroup.PATCH("/:id", handler.Patch(clientService, dbConnection))
		clientGroup.DELETE("/:id", handler.Delete(clientService, dbConnection))
	}

	IssuerService := &issuer.IssuerService{}
	issuerGroup := router.Group("issuers")
	{
		issuerGroup.GET("", handler.GetAll(IssuerService, dbConnection))
		issuerGroup.GET("/:id", handler.Get(IssuerService, dbConnection))
		issuerGroup.POST("/", handler.Post(IssuerService, dbConnection))
		issuerGroup.PUT("/:id", handler.Put(IssuerService, dbConnection))
		issuerGroup.PATCH("/:id", handler.Patch(IssuerService, dbConnection))
		issuerGroup.DELETE("/:id", handler.Delete(IssuerService, dbConnection))
	}

	cardService := &card.CardService{}
	cardGroup := router.Group("cards")
	{
		cardGroup.GET("", handler.GetAll(cardService, dbConnection))
		cardGroup.GET("/:id", handler.Get(cardService, dbConnection))
		cardGroup.POST("/", handler.Post(cardService, dbConnection))
		cardGroup.PUT("/:id", handler.Put(cardService, dbConnection))
		cardGroup.PATCH("/:id", handler.Patch(cardService, dbConnection))
		cardGroup.DELETE("/:id", handler.Delete(cardService, dbConnection))
	}

	tokenGroup := router.Group("tokens")
	{
		tokenGroup.POST("/", handler.Token(dbConnection))
	}

	return router
}
