package Server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"

	ServerUtils "malawi-ride-share-backend/internal/server/utils"

	"malawi-ride-share-backend/models"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)


type TestType struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber string `json:"PhoneNumber"`
}


func UserEndpoint(db *sql.DB, router *http.ServeMux) {
	CreateDriverHandleFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { 
	
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var requestBody models.CreateDriver
		err =json.Unmarshal(body, &requestBody)
		if err != nil { 
			http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		}

		registerDriverSqlStatement := `INSERT INTO drivers(user_id,license_number,vehicle_id, rating,status) VALUES($1,$2,$3,$4,$5)`

		_, err = db.Exec(registerDriverSqlStatement, requestBody.UserId, requestBody.LicenseNumber, requestBody.VehicleID, 5, "active")
		if err != nil {
			http.Error(w, "Failed to write to database", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "Driver created successfully"}`))
	})
	
	router.HandleFunc("POST /user", func(w http.ResponseWriter, r *http.Request) { 
		

		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var requestBody models.CreateUser
		err =json.Unmarshal(body, &requestBody)
		if err != nil { 
			http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		}

		if requestBody.Password1 != requestBody.Password2 {
			http.Error(w, "Passwords do not match", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password1), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Could not hash Password", http.StatusInternalServerError)
			return
		}

		var id string

		err = db.QueryRow("INSERT INTO users(first_name, last_name, phone_number,  email, password_hash, role) VALUES($1,$2,$3,$4,$5,$6) RETURNING id", requestBody.FirstName, requestBody.LastName, requestBody.PhoneNumber, requestBody.Email, hashedPassword, requestBody.Role).Scan(&id)
		if err != nil {
			http.Error(w, "Failed to write to database", http.StatusInternalServerError)
			return
		}

		tokenString, err := ServerUtils.GenerateToken(id, requestBody.PhoneNumber, requestBody.FirstName, requestBody.LastName)

		if err != nil {
			http.Error(w, "Could not create token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "User Created", "token": "` + tokenString + `"}`))

	})

	router.Handle("POST /user/driver", AuthMiddleware(CreateDriverHandleFunc ))
    
	router.HandleFunc("GET /user",  func( w http.ResponseWriter,r *http.Request){
		if r.Method != http.MethodGet {
			http.Error(w, "Invalide method", http.StatusBadRequest)
		}

		body,err := io.ReadAll(r.Body)
		if err != nil{
			http.Error(w, "Failed to read request body",  http.StatusBadRequest)
		}

		defer r.Body.Close()

		var requestBody models.Credentials
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		}

		ur := &TestType{}
		
		err = db.QueryRow(`SELECT first_name, last_name, phone_number FROM users WHRERE phone_number=$1`, requestBody.PhoneNumber).Scan(&ur.FirstName, &ur.LastName, &ur.PhoneNumber)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows){
				http.Error(w, "Incorrent number or password", http.StatusBadRequest)
			} else {
				http.Error(w, "Database query error", http.StatusInternalServerError)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		
		if err := json.NewEncoder(w).Encode(ur); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}
	})

	
}