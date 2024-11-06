package Server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type LocationUpdate struct {
	driverId string
	latitude float64
	longitude float64
}

var upgrader = websocket.Upgrader{ 
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

type DriverManager struct {}

func NewDriverManager() *DriverManager {
	return &DriverManager{}
}

func (dm *DriverManager) serveWS(w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection", err )
	}

	defer conn.Close()
}

func DriversEndpoint(db *sql.DB, router *http.ServeMux){
	router.HandleFunc("/ws/drivers", func(w http.ResponseWriter, r *http.Request){

	})
}