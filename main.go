package main

import (
	"malawi-ride-share-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)


func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/user", registerUser)
	r.Run() // listen and serve on 0.0.0.0:8080

}

func registerUser(c *gin.Context) {
	user := models.User{}
	if err := c.BindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
}