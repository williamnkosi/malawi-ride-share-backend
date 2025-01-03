package models

import "github.com/gorilla/websocket"

type DriversList map[*Driver]bool

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}
type Driver struct {
	Connection *websocket.Conn
	Manager    *DriverManager
	DriverId   string
	Location   *Location
}

type ResponseDriverData struct {
	DriverId string   `json:"driverId"`
	Location Location `json:"location"`
}

func NewDriver(driverId string, connection *websocket.Conn, manager *DriverManager) *Driver {
	return &Driver{
		Connection: connection,
		Manager:    manager,
		DriverId:   driverId,
	}
}

func (d *Driver) TrimData() ResponseDriverData {
	var td = &ResponseDriverData{}
	td.DriverId = d.DriverId
	td.Location = *d.Location
	return *td

}
