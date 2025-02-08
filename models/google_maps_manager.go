package models

import (
	"log"
	"os"

	"googlemaps.github.io/maps"
)

type GoogleMapsManager struct {
	ApiKey     string
	MapsClient *maps.Client
}

func NewGoogleMapsManager() *GoogleMapsManager {
	k := os.Getenv("GAPI_KEY")
	print(k)
	c, err := maps.NewClient(maps.WithAPIKey(k))

	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	print("Google Maps Manager created")

	return &GoogleMapsManager{
		ApiKey:     k,
		MapsClient: c,
	}

}

func (gm *GoogleMapsManager) GetDistanceMatrix() {}
func (gm *GoogleMapsManager) GetRoute()          {}
func (gm *GoogleMapsManager) GetProximity()      {}
