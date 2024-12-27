package Server

import (
	"malawi-ride-share-backend/models"
	"net/http"
)

func RideShareManagerEndpoint(r *http.ServeMux, dm *models.DriverManager, rm *models.RideShareManager) {

	r.HandleFunc("GET /drivers", func(w http.ResponseWriter, r *http.Request) {
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
	r.HandleFunc("/ws/riders", func(w http.ResponseWriter, r *http.Request) {

	})
	//r.HandleFunc("/ws/riders/accept", rm.AcceptRide)
	//r.HandleFunc("/ws/riders/start", rm.StartRide)
	//r.HandleFunc("/ws/riders/cancel", rm.CancelRide)
}
