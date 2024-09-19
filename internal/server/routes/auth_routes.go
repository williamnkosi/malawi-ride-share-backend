package Server

import (
	"database/sql"
	"errors"
	ServerUtils "malawi-ride-share-backend/internal/server/utils"
	"malawi-ride-share-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)


func AuthEndpoint(db *sql.DB, r *gin.Engine) {
	const AUTH_ENPOINT = "/auth"
	r.GET(AUTH_ENPOINT, func(c *gin.Context) {
		u := models.User{}
		l := models.Credentials{}
		if err := c.BindJSON(&l); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := db.QueryRow("SELECT id, first_name, last_name, phone_number, email, password_hash, role  FROM users WHERE phone_number=$1", l.PhoneNumber).Scan(&u.Id, &u.FirstName, &u.LastName, &u.PhoneNumber, &u.Email, &u.Password, &u.Role)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect Phone Number or Password"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
			}
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(l.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect Phone Number or Password"})
			return
		}

		tokenString, err := ServerUtils.GenerateToken(u.Id, u.FirstName, u.LastName, u.PhoneNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": tokenString})
	})

}


