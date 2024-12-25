package models

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{ 
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

type DriverManager struct {
	drivers DriversList
	sync.RWMutex
}

func NewDriverManager() *DriverManager {
	return &DriverManager{
		drivers: make(DriversList),
	}
}

func (dm *DriverManager) ServeWS(w http.ResponseWriter, r *http.Request){
	print("working")
	token := r.Header.Get("Authorization")
	driverId := r.Header.Get("DriverId")

	log.Println(token)
	log.Println(driverId)

	if token == "" || driverId == "" {
		http.Error(w, "Missing token or driverId", http.StatusBadRequest)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection", err )
	}

	d := NewDriver(driverId, conn, dm)
	dm.addDriver(d)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v\n", err)
			break
		}

		for d, avaliable := range dm.drivers {
			if d.driverId == driverId && avaliable {
				var receivedMessage Location
				if err := json.Unmarshal(message, &receivedMessage); err != nil {
					log.Println("Error unmarshaling JSON:", err)
					continue
				}

				d.location = &receivedMessage
			}
		}

		log.Printf("Received: %s\n", message)

		// Echo the message back to the client
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Printf("Error writing message: %v\n", err)
			break
		}
	}
}

func (dm * DriverManager) addDriver(d *Driver){
	dm.Lock()
	defer dm.Unlock()
	dm.drivers[d] = true
}

func (dm *DriverManager) removeDriver(d *Driver){
	dm.Lock()
	defer dm.Unlock()

	if _, ok := dm.drivers[d]; !ok {
		d.connection.Close()
		delete(dm.drivers, d)
	}
	
}