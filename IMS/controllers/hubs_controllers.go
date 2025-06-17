package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Trishank-Omniful/Onboarding-Task/db"
	"github.com/Trishank-Omniful/Onboarding-Task/models"
	"github.com/gin-gonic/gin"
)

func GetAllHubs(c *gin.Context) {
	var hubs []models.Hub
	client := db.GetDB()
	result := client.Find(&hubs)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error})
		return
	}
	c.JSON(200, hubs)
}

func GetHubById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Print("Issue While Converting Id to uint: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"Success": false,
		})
		return
	}
	var hub models.Hub
	client := db.GetDB()
	result := client.Where("id = ?", id).First(&hub)

	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error})
		return
	}
	c.JSON(200, hub)
}

func CreateHub(c *gin.Context) {
	var hub models.Hub
	err := c.BindJSON(&hub)
	if err != nil {
		log.Print("Issue While Converting Id to uint: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"Success": false,
		})
		return
	}
	client := db.GetDB()
	client.Create(&hub)
	c.JSON(200, hub)
}
