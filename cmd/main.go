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
	Server.AuthEndpoint(db,r)
	Server.LocationsEnpoint(db,r)
	Server.DriversEndpoint(db,r, dm)
	Server.UserEndpoint(db, r)

	server := http.Server {
		Addr: ":8081",
		Handler: Server.RecoveryMiddleware(r),
	}

	server.ListenAndServe()

}


