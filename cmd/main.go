package main

import (
	"malawi-ride-share-backend/internal/database"
	ServerMiddleware "malawi-ride-share-backend/internal/server"
	Server "malawi-ride-share-backend/internal/server/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)




func main() {
	r := gin.Default()

	db := database.InitializeDataBase()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Use(ServerMiddleware.CustomRecovery())
	Server.AuthEndpoint(db, r)
	Server.UserEndpoint(db, r)
	err := r.Run()
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080

}


