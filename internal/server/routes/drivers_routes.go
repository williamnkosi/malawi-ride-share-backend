package Server

import (
	"encoding/json"
	"log"
	Middleware "malawi-ride-share-backend/internal/middleware"
	"malawi-ride-share-backend/models"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func DriversEndpoint(router *http.ServeMux, dm *models.DriverManager) {
	driversHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")
		fcm := r.Header.Get("FcmToken")
		driverId := r.Header.Get("DriverId")

		// body, err := io.ReadAll(r.Body)
		// if err != nil {
		// 	http.Error(w, "Failed to read request body", http.StatusBadRequest)
		// 	return
		// }

		// defer r.Body.Close()

		// var requestBody models.Location
		// err = json.Unmarshal(body, &requestBody)
		// {
		// 	if err != nil {
		// 		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		// 		return
		// 	}
		// }

		if token == "" || driverId == "" || fcm == "" {
			log.Println("Failed")
			http.Error(w, "Missing token or driverId, token, Location", http.StatusBadRequest)
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Failed to upgrade connection", err)
		}

		d := models.NewDriver(driverId, fcm, conn, dm)
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
					print(d.Location)
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
	router.HandleFunc("/ws/drivers", func(w http.ResponseWriter, r *http.Request) {
		Middleware.FirebaseAuthMiddleware(driversHandler).ServeHTTP(w, r)
	})

}
