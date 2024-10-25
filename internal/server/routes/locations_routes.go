package Server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LocationRequestBody struct {
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name string `json:"name"`
	Address string `json:"address"`
	City string `json:"city"`
	Country string `json:"country"`

}

func LocationsEnpoint(db *sql.DB, r *http.ServeMux) {
	r.HandleFunc("POST /locations", func(w http.ResponseWriter,  r *http.Request){
		if r.Method != http.MethodPost {
			http.Error(w,"Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		
		var requestBody LocationRequestBody
		fmt.Print(body)
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			fmt.Print(requestBody)
			http.Error(w,"Failed to parse JSON", http.StatusBadRequest)
		}

		createLocationStatement := `INSERT INTO Locations(latitude,longitude,address,city,country) VALUES($1,$2,$3,$4,$5)`
		_, err = db.Exec(createLocationStatement, requestBody.Latitude, requestBody.Longitude,requestBody.Address,requestBody.City,requestBody.Country )

		if err != nil {
			http.Error(w, "Failed to write to database", http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": Location Recieved"}`))
	})
}
