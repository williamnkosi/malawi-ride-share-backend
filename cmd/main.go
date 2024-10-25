package main

import (
	"malawi-ride-share-backend/internal/database"
	Server "malawi-ride-share-backend/internal/server/routes"
	"net/http"

	_ "github.com/lib/pq"
)





func main() {
	//r := gin.Default()

	r := http.NewServeMux()
	db := database.InitializeDataBase()
	Server.LocationsEnpoint(db,r)

	server := http.Server {
		Addr: ":8081",
		Handler: r,
	}

	server.ListenAndServe()

	//db := database.InitializeDataBase()
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })
	// r.Use(ServerMiddleware.CustomRecovery())
	// Server.AuthEndpoint(db, r)
	// Server.UserEndpoint(db, r)
	// err := r.Run()
	// if err != nil {
	// 	return
	// } // listen and serve on 0.0.0.0:8080

}


