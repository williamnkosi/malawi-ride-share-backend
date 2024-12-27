package models

import (
	"math/rand"
	"time"
)

type RideStatus int

const (
	Requested = iota
	Accepted
	DriverOnRoute
	DriverArrived
	RideInProgress
	Completed
	Cancelled
	Failed
)

type Rider struct {
	RiderId        string
	RiderFirstName string
	RiderLastName  string
}
type Ride struct {
	RideId              string
	RideStatus          RideStatus
	RiderStartLocation  Location
	DriverStartLocation Location
	Destination         Location
	Driver              Driver
	Rider               Rider
	dm                  *DriverManager
}

type RideShareManager struct {
	Rides map[string]*Ride
}

func NewRide(r *Ride) *Ride {

	rand.Seed(time.Now().UnixNano())
	return &Ride{
		RideId: "testing",
	}
}

func (rm *RideShareManager) RequestDriverLocation(dm *DriverManager) {

}

func (rm *RideShareManager) RequestRide() {}
func (rm *RideShareManager) AcceptRide()  {}
func (rm *RideShareManager) StartRide()   {}
func (rm *RideShareManager) CancelRide()  {}
