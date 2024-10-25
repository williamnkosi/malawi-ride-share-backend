package Server

import (
	"database/sql"
	"encoding/json"
	"errors"
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



func LocationsEnpoint(db *sql.DB, router *http.ServeMux) {
	router.HandleFunc("POST /locations", func(w http.ResponseWriter,  r *http.Request){
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
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
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

	router.HandleFunc("GET /locations/{id}", func(w http.ResponseWriter,  r *http.Request){
		type Location struct {
			Id string
			Latitude float64
			Longitude float64
			Address string
			City string
			Country string
		}
		l := Location{}
		if r.Method != http.MethodGet {
			http.Error(w,"Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		id := r.PathValue("id")

		if id == "" {
			http.Error(w,"Didnt provide a valid value", http.StatusBadRequest)
		}

		LocationQueryStatement := `SELECT id, latitude, longitude, address, city, country FROM Locations WHERE id=` + id

		err := db.QueryRow(LocationQueryStatement).Scan(&l.Id, &l.Latitude,&l.Longitude, &l.Address, &l.City, &l.Country)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w,"No query result", http.StatusBadRequest)
				
			} else {
				http.Error(w,"Couldn't complete process",http.StatusInternalServerError)
		
			}
			return 
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(l); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	router.HandleFunc("PUT /locations", func(w http.ResponseWriter,  r *http.Request){
		type UpdateLocationBody struct {
			Id string `json:"id"`
			Latitude float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Name string `json:"name"`
			Address string `json:"address"`
			City string `json:"city"`
			Country string `json:"country"`
		}
		if r.Method != http.MethodPut {
			http.Error(w,"Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var requestBody UpdateLocationBody
		err = json.Unmarshal(body,&requestBody); if err != nil {
			http.Error(w,"Failed to parse JSON", http.StatusBadRequest)
			return
		}

		locationQuery := "UPDATE locations SET latitude=$1, longitude=$2, address=$3, city=$4, country=$5 WHERE id=$6"
		result, err := db.Exec(locationQuery, requestBody.Latitude, requestBody.Longitude, requestBody.Address, requestBody.City, requestBody.Country, requestBody.Id)
		if err != nil {
			http.Error(w, "Failed to update location", http.StatusInternalServerError)
		}

		rowsAffecteded, err := result.RowsAffected()
		if err != nil || rowsAffecteded == 0 {http.Error(w,"Locaton was not found", http.StatusBadRequest) }

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(requestBody); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
	})

	router.HandleFunc("DELETE /locations/{id}", func(w http.ResponseWriter,  r *http.Request){
		if r.Method != http.MethodDelete {
			http.Error(w,"Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		id := r.PathValue("id")

		if id == "" {
			http.Error(w,"Didnt provide a valid value", http.StatusBadRequest)
		}

		LocationsDeleteStatement := `DELETE FROM locations WHERE id=$1`
		_, err := db.Query(LocationsDeleteStatement, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w,"No query result", http.StatusBadRequest)
				
				} else {
					http.Error(w,"Couldn't complete process",http.StatusInternalServerError)
			
				}
				return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Location Deleted"}`))
	})
	
}
