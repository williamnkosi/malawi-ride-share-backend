package models

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
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

func (dm *DriverManager) ServeWS(w http.ResponseWriter, r *http.Request) {

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

	d := NewDriver(driverId, conn, dm)
	dm.AddDriver(d)

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

func (dm *DriverManager) AddDriver(d *Driver) {
	dm.Lock()
	defer dm.Unlock()
	dm.drivers[d] = true
}

func (dm *DriverManager) RemoveDriver(d *Driver) {
	dm.Lock()
	defer dm.Unlock()

	if _, ok := dm.drivers[d]; !ok {
		d.connection.Close()
		delete(dm.drivers, d)
	}

}

func (dm *DriverManager) GetAllDrivers() []ResponseDriverData {
	dm.RLock()
	l := []ResponseDriverData{}
	for d, avaliable := range dm.drivers {
		if avaliable && d.location != nil {
			trimmedData := d.TrimData()
			l = append(l, trimmedData)
		}
	}
	defer dm.RUnlock()

	return l
}

func (dm *DriverManager) GetDriversByProximity() []ResponseDriverData {
	dm.RLock()
	l := []ResponseDriverData{}
	defer dm.RUnlock()

	return l
}
