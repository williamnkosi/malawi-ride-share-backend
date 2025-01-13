package models

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

func NewRideShareManager() *RideShareManager {
	type RidesList map[string]*Ride
	return &RideShareManager{Rides: make(RidesList)}
}

func (rm *RideShareManager) RequestDriverLocation(dm *DriverManager) {

}

func (rm *RideShareManager) RequestRide() {}
func (rm *RideShareManager) AcceptRide()  {}
func (rm *RideShareManager) StartRide()   {}
func (rm *RideShareManager) CancelRide()  {}
