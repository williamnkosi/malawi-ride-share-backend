package main

import (
	"database/sql"
	"fmt"
	"log"
	"malawi-ride-share-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)


func main() {
	r := gin.Default()

	db := initializeDataBase()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	UserEndpoint(db,r)
 	r.Run() // listen and serve on 0.0.0.0:8080


}

func initializeDataBase() *sql.DB {
	connStr := "postgres://postgres:password@localhost:5432/postgres"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")
	return db;
}

func UserEndpoint(db *sql.DB, r *gin.Engine){
	r.POST("/user", func( c *gin.Context) {
		u := models.User{}

		if err := c.BindJSON(&u); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if u.Password1 != u.Password2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
			return  
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password1), bcrypt.DefaultCost)
		if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
            return
        }
		
		registerUserSqlStatement := `INSERT INTO test.users(firstName, lastName, phoneNumber, email, age, password) VALUES($1,$2,$3,$4,$5,$6)`
		_ , err = db.Exec(registerUserSqlStatement, u.FirstName, u.LastName, u.PhoneNumber, u.Email,u.Age, hashedPassword)
		if(err != nil){
			c.JSON(http.StatusBadRequest, gin.H{"error with database": err.Error()})
			return
		} 
		c.JSON(http.StatusCreated, gin.H{"status": "User created"})
	})
}

