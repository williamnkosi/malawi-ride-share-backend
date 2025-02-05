package Server

import (
	"encoding/json"
	"io"
	"malawi-ride-share-backend/models"
	"net/http"
)

func RideShareManagerEndpoint(r *http.ServeMux, dm *models.DriverManager, rm *models.RideShareManager, gm *models.GoogleMapsManager, mm *models.MessagingManager) {

	r.HandleFunc("GET /rideshare/drivers", func(w http.ResponseWriter, r *http.Request) {
		d := dm.GetAllDrivers()
		data := map[string]interface{}{
			"list": d,
		}
		res := models.Response{
			Status:  http.StatusOK,
			Message: "List of all available drivers",
			Data:    data,
		}

		res.SuccessfulResponse(w, "Drivers", d)

	})

	r.HandleFunc("GET /rideshare/request", func(w http.ResponseWriter, r *http.Request) {

		type Body struct {
			Id       string          `json:"id"`
			Location models.Location `json:"location"`
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
		}
		defer r.Body.Close()

		var requestBody Body
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		}

		d := dm.GetDriversByProximity(requestBody.Location, gm)

		for _, d := range d {
			mm.SendNotification(d.FcmToken)
			mm.SendDataMessage(d.FcmToken)
		}

	})

	//r.HandleFunc("/ws/riders/start", rm.StartRide)
	//r.HandleFunc("/ws/riders/cancel", rm.CancelRide)
}
