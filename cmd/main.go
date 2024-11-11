package main

import (
	"malawi-ride-share-backend/internal/database"
	Server "malawi-ride-share-backend/internal/server/routes"
	models "malawi-ride-share-backend/models"
	"net/http"

	_ "github.com/lib/pq"
)





func main() {
	r := http.NewServeMux()
	db := database.InitializeDataBase()
	dm := models.NewDriverManager()
	Server.LocationsEnpoint(db,r)
	Server.DriversEndpoint(db,r, dm)
	Server.UserEndpoint(db, r)

	server := http.Server {
		Addr: ":8081",
		Handler: r,
	}

	server.ListenAndServe()
	// r.Use(ServerMiddleware.CustomRecovery())
	// Server.AuthEndpoint(db, r)
	// Server.UserEndpoint(db, r)
	// err := r.Run()
	// if err != nil {
	// 	return
	// } // listen and serve on 0.0.0.0:8080

}


