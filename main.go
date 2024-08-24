package main

import (
	"database/sql"
	"fmt"
	"log"
	"malawi-ride-share-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)


func main() {
	r := gin.Default()

	initializeDataBase()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
 	r.POST("/user", registerUser)
 	r.Run() // listen and serve on 0.0.0.0:8080


}

func initializeDataBase() {
	connStr := "postgres://postgres:password@localhost:5432"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")
}

func registerUser(c *gin.Context) {
	user := models.User{}
	if err := c.BindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
}