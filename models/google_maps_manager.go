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
	c, err := maps.NewClient(maps.WithAPIKey(k))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	print("Google Maps Manager created")
	return &GoogleMapsManager{
		ApiKey:     "your_api_key",
		MapsClient: c,
	}

}

func (gm *GoogleMapsManager) GetDistanceMatrix() {}
func (gm *GoogleMapsManager) GetRoute()          {}
func (gm *GoogleMapsManager) GetProximity()      {}
