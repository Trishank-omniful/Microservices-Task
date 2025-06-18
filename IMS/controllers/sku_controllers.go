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

type SkuController struct {
	Repo *repository.SkuRepository
}

func NewSkuController(repo *repository.SkuRepository) *SkuController {
	return &SkuController{Repo: repo}
}

func (ctrl *SkuController) GetAllSkus(c *gin.Context) {
	hubs, err := ctrl.Repo.GetAllSkus()
	if err != nil {
		log.Print("Failed to get all hubs: ", err)
		c.JSON(500, gin.H{"error": err})
		return
	}
	c.JSON(200, hubs)
}

func (ctrl *SkuController) GetSkuById(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Print("Invalid ID Parameter: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	hub, err := ctrl.Repo.GetSkuById(uint(id))

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "SKU Not Found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server Error"})
	}

	c.JSON(http.StatusOK, hub)
}

func (ctrl *SkuController) CreateSku(c *gin.Context) {
	var hub models.SKU
	err := c.ShouldBindJSON(&hub)
	if err != nil {
		log.Print("Issue While Parsing JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	err = ctrl.Repo.CreateSku(&hub)
	if err != nil {
		log.Print("Record Already Exists: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Record Already Exists",
		})
		return
	}
	c.JSON(200, hub)
}

func (ctrl *SkuController) UpdateSku(c *gin.Context) {
	var updatedData models.SKU
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
	err = ctrl.Repo.UpdateSku(&updatedData)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "SKU not found"})
		return
	} else if err != nil {
		c.JSON(500, gin.H{"error": "Server Error"})
		return
	}

	c.JSON(200, updatedData)
}

func (ctrl *SkuController) DeleteSku(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Print("Issue While Converting Id to uint: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	err = ctrl.Repo.DeleteSku(uint(id))

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "SKU not found"})
		return
	} else if err != nil {
		c.JSON(500, gin.H{"error": "Server Error"})
		return
	}

	c.JSON(200, gin.H{"Success": true})
}
