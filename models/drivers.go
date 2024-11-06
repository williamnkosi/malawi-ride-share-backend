package models

import "github.com/gorilla/websocket"

type DriversList map[*Driver]bool

type Driver struct {
	connection *websocket.Conn
	manager *DriverManager
	driverId string
	longitude float64
	latitude float64
}

func NewDriver(driverId string, connection *websocket.Conn, manager *DriverManager) *Driver {
	return &Driver{
		connection: connection,
		manager: manager,
		driverId: driverId,
	}
}
