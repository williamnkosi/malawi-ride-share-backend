package Server

import (
	"database/sql"
	"net/http"
)

type LocationUpdate struct {
	driverId string
	latitude float64
	longitude float64
}

func DriversEndpoint(db *sql.DB, router *http.ServeMux){
	router.HandleFunc("/ws/drivers", func(w http.ResponseWriter, r *http.Request){

	})
}