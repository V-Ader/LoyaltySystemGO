package main

import (
	"github.com/V-Ader/Loyality_GO/api/router"
)

func main() {

	user_router := router.New()
	user_router.Run(":8080")
}
