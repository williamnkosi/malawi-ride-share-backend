package models

import (
	"context"
	"fmt"
	"sync"

	"googlemaps.github.io/maps"
)

type DriverManager struct {
	Drivers DriversList
	sync.RWMutex
}

func NewDriverManager() *DriverManager {
	return &DriverManager{
		Drivers: make(DriversList),
	}
}

func (dm *DriverManager) AddDriver(d *Driver) {
	dm.Lock()
	defer dm.Unlock()
	dm.Drivers[d] = true
}

func (dm *DriverManager) RemoveDriver(d *Driver) {
	dm.Lock()
	defer dm.Unlock()

	if _, ok := dm.Drivers[d]; !ok {
		d.Connection.Close()
		delete(dm.Drivers, d)
	}

}

func (dm *DriverManager) GetAllDrivers() []ResponseDriverData {
	dm.RLock()
	l := []ResponseDriverData{}
	for d, avaliable := range dm.Drivers {
		if avaliable && d.Location != nil {
			trimmedData := d.TrimData()
			l = append(l, trimmedData)
		}
	}
	defer dm.RUnlock()

	return l
}

func (dm *DriverManager) GetDriversByProximity(userLocation Location, gm *GoogleMapsManager) []Driver {
	dm.RLock()
	l := []Driver{}
	for d, avaliable := range dm.Drivers {
		if avaliable && d.Location != nil {
			req := &maps.DistanceMatrixRequest{
				Origins:      []string{fmt.Sprintf("%f,%f", userLocation.Latitude, userLocation.Longitude)},
				Destinations: []string{fmt.Sprintf("%f,%f", d.Location.Latitude, d.Location.Longitude)},
				Mode:         maps.TravelModeDriving,
				Units:        maps.UnitsImperial,
			}

			resp, err := gm.MapsClient.DistanceMatrix(context.Background(), req)
			if err != nil {
				fmt.Println(err)
			}

			if len(resp.Rows) > 0 && len(resp.Rows[0].Elements) > 0 {

				if len(l) < 2 {
					l = append(l, *d)
				}

			} else {
				fmt.Println("No distance data available")
			}
		}
	}
	defer dm.RUnlock()

	return l
}
