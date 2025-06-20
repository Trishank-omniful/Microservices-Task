package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Trishank-Omniful/Onboarding-Task/constants"
	"github.com/Trishank-Omniful/Onboarding-Task/models"
	"github.com/Trishank-Omniful/Onboarding-Task/repository"
	"github.com/Trishank-Omniful/Onboarding-Task/validators"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrGetAllHubs})
		return
	}
	c.JSON(http.StatusOK, hubs)
}

func (ctrl *HubController) GetHubById(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Print("Invalid ID Parameter: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidID})
		return
	}

	hub, err := ctrl.Repo.GetHubById(uint(id))

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrHubNotFound})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrServerError})
		return
	}

	c.JSON(http.StatusOK, hub)
}

func (ctrl *HubController) CreateHub(c *gin.Context) {
	var hub models.Hub
	err := c.ShouldBindJSON(&hub)
	if err != nil {
		log.Print("Issue While Parsing JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.ErrParsingJSON,
		})
		return
	}

	if err := validators.ValidateHub(&hub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = ctrl.Repo.CreateHub(&hub)
	if err != nil {
		log.Print("Failed to create hub: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": constants.ErrHubCreate,
		})
		return
	}
	c.JSON(http.StatusCreated, hub)
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
		c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrHubNotFound})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrServerError})
		return
	}

	c.JSON(http.StatusOK, updatedData)
}

func (ctrl *HubController) DeleteHub(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Print("Invalid ID Parameter: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidID})
		return
	}

	err = ctrl.Repo.DeleteHub(uint(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrHubNotFound})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hub deleted successfully"})
}

func (ctrl *HubController) CreateHubsBatch(c *gin.Context) {
	var hubs []models.Hub
	err := c.ShouldBindJSON(&hubs)
	if err != nil {
		log.Print("Issue While Parsing JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.ErrParsingJSON,
		})
		return
	}

	if err := validators.ValidateBatchSize(len(hubs)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	for i, hub := range hubs {
		if err := validators.ValidateHub(&hub); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     err.Error(),
				"hub_index": i,
			})
			return
		}
	}

	err = ctrl.Repo.CreateHubsBatch(hubs)
	if err != nil {
		log.Print("Batch hub creation failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": constants.ErrBatchOperation,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Hubs created successfully",
		"count":   len(hubs),
	})
}

func (ctrl *HubController) GetHubsByIDs(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrParsingJSON})
		return
	}

	if len(request.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No IDs provided"})
		return
	}

	hubs, err := ctrl.Repo.GetHubsByIDs(request.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrServerError})
		return
	}

	c.JSON(http.StatusOK, hubs)
}
