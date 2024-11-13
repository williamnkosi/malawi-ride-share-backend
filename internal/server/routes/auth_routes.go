package Server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	ServerUtils "malawi-ride-share-backend/internal/server/utils"
	"malawi-ride-share-backend/models"
	"net/http"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)


func AuthEndpoint(db *sql.DB,  router *http.ServeMux) {
	router.HandleFunc("GET /auth", func(w http.ResponseWriter, r *http.Request){
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request methog", http.StatusMethodNotAllowed)
			return 
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body",http.StatusBadRequest)
			return 
		}

		defer r.Body.Close()


		var requestBody models.Credentials
		err = json.Unmarshal(body,&requestBody)
		if err != nil {
			http.Error(w, "Could not parse JSON", http.StatusInternalServerError)
			return
		}

		

		u := models.User{}
		err = db.QueryRow("SELECT id, first_name, last_name, phone_number, email, password_hash, role  FROM users WHERE phone_number=$1", requestBody.PhoneNumber).Scan(&u.Id, &u.FirstName, &u.LastName, &u.PhoneNumber, &u.Email, &u.Password, &u.Role)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w,"Incorrect Phone Number or Password", http.StatusBadRequest)
			} else {
				http.Error(w, "Interal server error", http.StatusInternalServerError)
			}
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(requestBody.Password)); err != nil {
			http.Error(w, "Incorrect phone number or password", http.StatusBadGateway)
		
			return
		}

		tokenString, err := ServerUtils.GenerateToken(u.Id, u.FirstName, u.LastName, u.PhoneNumber)
		if err != nil {
			http.Error(w, "Couldn't create token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"token":` + tokenString + `}`))
	})
		
	



}


