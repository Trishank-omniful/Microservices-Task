package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Trishank-Omniful/Onboarding-Task/db"
	"github.com/Trishank-Omniful/Onboarding-Task/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		log.Print("Invalid ID: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var hub models.Hub
	client := db.GetDB()
	result := client.First(&hub, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hub not found"})
		return
	}

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, hub)
}

func CreateHub(c *gin.Context) {
	var hub models.Hub
	err := c.BindJSON(&hub)
	if err != nil {
		log.Print("Issue While Parsing JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"Success": false,
		})
		return
	}
	client := db.GetDB()
	client.Create(&hub)
	c.JSON(200, hub)
}

func UpdateHub(c *gin.Context) {
	var hub models.Hub
	var updatedData models.Hub
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Print("Issue While Converting Id to uint: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"Success": false,
		})
		return
	}
	err = c.BindJSON(&updatedData)
	if err != nil {
		log.Print("Issue While Parsing JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"Success": false,
		})
		return
	}

	client := db.GetDB()

	result := client.First(&hub, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hub not found"})
		return
	}

	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Server Error"})
		return
	}

	update := client.Model(&hub).Updates(updatedData)

	if update.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to update hub", "details": update.Error.Error()})
		return
	}

	c.JSON(200, hub)
}

func DeleteHub(c *gin.Context) {
	var hub models.Hub
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Print("Issue While Converting Id to uint: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"Success": false,
		})
		return
	}
	client := db.GetDB()
	delete := client.Where("id = ?", id).Delete(&hub)

	if delete.Error != nil {
		c.JSON(500, gin.H{"error": "Server Error"})
		return
	}

	c.JSON(200, hub)
}
