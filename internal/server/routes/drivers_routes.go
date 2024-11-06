package Server

import (
	"database/sql"
	"malawi-ride-share-backend/models"
	"net/http"
)

type LocationUpdate struct {
	driverId string
	latitude float64
	longitude float64
}

func DriversEndpoint(db *sql.DB, router *http.ServeMux, dm *models.DriverManager){
	router.HandleFunc("/ws/drivers", dm.ServeWS)
}