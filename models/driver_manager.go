package models

import (
	"sync"
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

func (dm *DriverManager) GetDriversByProximity() []ResponseDriverData {
	dm.RLock()
	l := []ResponseDriverData{}
	defer dm.RUnlock()

	return l
}
