package models

import (
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

func (dm *DriverManager) serveWS(w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection", err )
	}

	d := NewDriver("driverId", conn, dm)
	dm.addDriver(d)
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