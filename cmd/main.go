package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	Server "malawi-ride-share-backend/internal/server/routes"
	models "malawi-ride-share-backend/models"
	"net/http"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"

	_ "github.com/lib/pq"
)

func main() {
	log.Println("/Users/williamnkosi/repo/malawi-ride-share-backend/cmd/serviceAccountKey.json")
	opt := option.WithCredentialsFile("/Users/williamnkosi/repo/malawi-ride-share-backend/cmd/serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}

	r := http.NewServeMux()
	//db := database.InitializeDataBase()
	dm := models.NewDriverManager()
	rm := models.NewRideShareManager()
	gm := models.NewGoogleMapsManager()
	mm := models.NewMessagingManager(app)
	// Server.AuthEndpoint(db,r)
	// Server.LocationsEnpoint(db,r)
	Server.DriversEndpoint(r, dm)
	Server.RideShareManagerEndpoint(r, dm, rm, gm, mm)
	// Server.UserEndpoint(db, r)

	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		{
			if err != nil {
				http.Error(w, "Failed to read request body", http.StatusBadRequest)
				return
			}
		}

		defer r.Body.Close()
		type TestType struct {
			Token string `json:"token"`
		}

		var requestBody TestType
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		}

		var token = requestBody.Token
		mm.SendDataMessage(token)
		//w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "pong")
	})

	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: r,
	}

	server.ListenAndServe()

}
