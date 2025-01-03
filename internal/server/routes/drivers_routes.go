package Server

import (
	"encoding/json"
	"log"
	"malawi-ride-share-backend/models"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type LocationUpdate struct {
	driverId  string
	latitude  float64
	longitude float64
}

func DriversEndpoint(router *http.ServeMux, dm *models.DriverManager) {
	router.HandleFunc("/ws/drivers", func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")
		driverId := r.Header.Get("DriverId")

		println("----")
		println(token)
		println(driverId)

		if token == "" || driverId == "" {
			//log.Println("Failed")
			//http.Error(w, "Missing token or driverId", http.StatusBadRequest)
			//return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Failed to upgrade connection", err)
		}

		d := models.NewDriver(driverId, conn, dm)
		dm.AddDriver(d)

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v\n", err)
				break
			}

			for d, avaliable := range dm.Drivers {
				if d.DriverId == driverId && avaliable {
					var receivedMessage models.Location
					if err := json.Unmarshal(message, &receivedMessage); err != nil {
						log.Println("Error unmarshaling JSON:", err)
						continue
					}

					d.Location = &receivedMessage
				}
			}

			log.Printf("Received: %s\n", message)

			// Echo the message back to the client
			if err := conn.WriteMessage(messageType, message); err != nil {
				log.Printf("Error writing message: %v\n", err)
				break
			}
		}
	})
}
