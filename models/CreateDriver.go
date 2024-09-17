package models

type CreateDriver struct {
	UserId        string `json:"userId"`
	LicenseNumber string `json:"licenseNumber"`
	VehicleID     string `json:"vehicleId"`
}
