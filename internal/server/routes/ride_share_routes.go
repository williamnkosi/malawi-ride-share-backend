package Server

import (
	"malawi-ride-share-backend/models"
	"net/http"
)

func RideShareManagerEndpoint(r *http.ServeMux, dm *models.DriverManager, rm *models.RideShareManager) {

	//r.HandleFunc("/ws/drivers", dm.ServeWS)
	r.HandleFunc("/ws/riders", func(w http.ResponseWriter, r *http.Request) {
		dm.addDriver(models.NewDriver("1", nil, dm))
	})
	//r.HandleFunc("/ws/riders/accept", rm.AcceptRide)
	//r.HandleFunc("/ws/riders/start", rm.StartRide)
	//r.HandleFunc("/ws/riders/cancel", rm.CancelRide)
}
