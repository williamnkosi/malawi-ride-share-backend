package main

import (
	"database/sql"
	"fmt"
	"log"
	"malawi-ride-share-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey  = []byte("testing-the-new-key")

func main() {
	r := gin.Default()

	db := initializeDataBase()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	authEndpoint(db,r)
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



func authEndpoint(db *sql.DB, r *gin.Engine){
	r.GET("/auth", func(c * gin.Context){
		u := models.User{}
		l := models.Credentials{}
		if err := c.BindJSON(&l); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		 err := db.QueryRow("SELECT user_id, firstName, lastName, phoneNumber, email, age, password FROM test.users WHERE phoneNumber=$1", l.PhoneNumber).Scan(&u.Id,&u.FirstName, &u.LastName,&u.PhoneNumber, &u.Email, &u.Age,  &u.Password)
		if err != nil {
			if err == sql.ErrNoRows{
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect Phone Number or Password"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
			}
			return
		}

		if err := bcrypt.CompareHashAndPassword( []byte(u.Password),[]byte(l.Password)); err != nil{
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect Phone Number or Password"})
			return 
		}

		expirationTime := time.Now().Add(2* time.Hour)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"phoneNumber": u.PhoneNumber,
			"firstName": u.FirstName,
			"lastName": u.LastName,
			"expiration-time": expirationTime,
		})
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
			return 
		}

		c.JSON(http.StatusOK, gin.H{"data": tokenString})
		
	})
}

func UserEndpoint(db *sql.DB, r *gin.Engine){
	r.POST("/user", func( c *gin.Context) {
		u := models.CreateUser{}

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

