package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Trishank-Omniful/Onboarding-Task/models"
	"github.com/Trishank-Omniful/Onboarding-Task/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HubController struct {
	Repo *repository.HubRepository
}

func NewHubController(repo *repository.HubRepository) *HubController {
	return &HubController{Repo: repo}
}

func (ctrl *HubController) GetAllHubs(c *gin.Context) {
	hubs, err := ctrl.Repo.GetAllHubs()
	if err != nil {
		log.Print("Failed to get all hubs: ", err)
		c.JSON(500, gin.H{"error": err})
		return
	}
	c.JSON(200, hubs)
}

func (ctrl *HubController) GetHubById(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Print("Invalid ID Parameter: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	hub, err := ctrl.Repo.GetHubById(uint(id))

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hub Not Found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server Error"})
	}

	c.JSON(http.StatusOK, hub)
}

func (ctrl *HubController) CreateHub(c *gin.Context) {
	var hub models.Hub
	err := c.ShouldBindJSON(&hub)
	if err != nil {
		log.Print("Issue While Parsing JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	err = ctrl.Repo.CreateHub(&hub)
	if err != nil {
		log.Print("Record Already Exists: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Record Already Exists",
		})
		return
	}
	c.JSON(200, hub)
}

func (ctrl *HubController) UpdateHub(c *gin.Context) {
	var updatedData models.Hub
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Print("Issue While Converting Id to uint: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	err = c.ShouldBindJSON(&updatedData)
	if err != nil {
		log.Print("Issue While Parsing JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	updatedData.ID = uint(id)
	err = ctrl.Repo.UpdateHub(&updatedData)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hub not found"})
		return
	} else if err != nil {
		c.JSON(500, gin.H{"error": "Server Error"})
		return
	}

	c.JSON(200, updatedData)
}

func (ctrl *HubController) DeleteHub(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Print("Issue While Converting Id to uint: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	err = ctrl.Repo.DeleteHub(uint(id))

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hub not found"})
		return
	} else if err != nil {
		c.JSON(500, gin.H{"error": "Server Error"})
		return
	}

	c.JSON(200, gin.H{"Success": true})
}
