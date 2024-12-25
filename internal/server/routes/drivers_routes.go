package Server

import (
	"malawi-ride-share-backend/models"
	"net/http"
)

type LocationUpdate struct {
	driverId string
	latitude float64
	longitude float64
}

func DriversEndpoint( router *http.ServeMux, dm *models.DriverManager){
	router.HandleFunc("/ws/drivers", dm.ServeWS)
}