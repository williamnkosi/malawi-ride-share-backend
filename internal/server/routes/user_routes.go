package Server

import (
	"database/sql"
	"errors"
	ServerMiddleware "malawi-ride-share-backend/internal/server"
	ServerUtils "malawi-ride-share-backend/internal/server/utils"
	"malawi-ride-share-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)


type TestType struct {
	FirstName   string
	LastName    string
	PhoneNumber string
}


func UserEndpoint(db *sql.DB, r *gin.Engine) {

	const USER_ENPOINT = "/user"
	r.POST(USER_ENPOINT, func(c *gin.Context) {
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

		//registerUserSqlStatement := `INSERT INTO users(first_name, last_name, phone_number,  email, password_hash, role) VALUES($1,$2,$3,$4,$5,$6)`
		//_, err = db.Exec(registerUserSqlStatement, u.FirstName, u.LastName, u.PhoneNumber, u.Email, hashedPassword, u.Role)

		//if err != nil {
		//	c.JSON(http.StatusBadRequest, gin.H{"error with database": err.Error()})
		//	return
		//}

		var id string

		err = db.QueryRow("INSERT INTO users(first_name, last_name, phone_number,  email, password_hash, role) VALUES($1,$2,$3,$4,$5,$6) RETURNING id", u.FirstName, u.LastName, u.PhoneNumber, u.Email, hashedPassword, u.Role).Scan(&id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error with database": err.Error()})
			return
		}

		tokenString, err := ServerUtils.GenerateToken(id, u.PhoneNumber, u.FirstName, u.LastName)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"status": "User created", "token": tokenString, "id": id})
	})

	r.POST(USER_ENPOINT+"/driver", ServerMiddleware.AuthMiddleware(), func(c *gin.Context) {
		var d = models.CreateDriver{}
		if err := c.BindJSON(&d); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		registerDriverSqlStatement := `INSERT INTO drivers(user_id,license_number,vehicle_id, rating,status) VALUES($1,$2,$3,$4,$5)`

		_, err := db.Exec(registerDriverSqlStatement, d.UserId, d.LicenseNumber, d.VehicleID, 5, "active")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error with database": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"status": "Driver created"})

	})

	r.GET("/user", ServerMiddleware.AuthMiddleware(), func(c *gin.Context) {

		cr := models.Credentials{}
		ur := &TestType{}
		if err := c.BindJSON(&cr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		phoneNumber := c.MustGet("phoneNumber").(string)

		if phoneNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no phoneNumber"})
			return
		}
		err := db.QueryRow(`SELECT first_name,last_name, phone_number FROM users WHERE phone_number=$1 `, phoneNumber).Scan(&ur.FirstName, &ur.LastName, &ur.PhoneNumber)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect Phone Number or Password"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": ur})

	})
}