package models

import "github.com/gorilla/websocket"

type DriversList map[*Driver]bool

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}
type Driver struct {
	connection *websocket.Conn
	manager    *DriverManager
	driverId   string
	location   *Location
}

type ResponseDriverData struct {
	DriverId string   `json:"driverId"`
	Location Location `json:"location"`
}

func NewDriver(driverId string, connection *websocket.Conn, manager *DriverManager) *Driver {
	return &Driver{
		connection: connection,
		manager:    manager,
		driverId:   driverId,
	}
}

func (d *Driver) TrimData() ResponseDriverData {
	var td = &ResponseDriverData{}
	td.DriverId = d.driverId
	td.Location = *d.location
	return *td

}
