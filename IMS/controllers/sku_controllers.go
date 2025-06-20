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

type SkuController struct {
	Repo *repository.SkuRepository
}

func NewSkuController(repo *repository.SkuRepository) *SkuController {
	return &SkuController{Repo: repo}
}

func (ctrl *SkuController) GetAllSkus(c *gin.Context) {
	skus, err := ctrl.Repo.GetAllSkus()
	if err != nil {
		log.Print("Failed to get all SKUs: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrGetAllSKUs})
		return
	}
	c.JSON(http.StatusOK, skus)
}

func (ctrl *SkuController) GetSkuById(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Print("Invalid ID Parameter: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidID})
		return
	}

	sku, err := ctrl.Repo.GetSkuById(uint(id))

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrSKUNotFound})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrServerError})
		return
	}

	c.JSON(http.StatusOK, sku)
}

func (ctrl *SkuController) CreateSku(c *gin.Context) {
	var sku models.SKU
	err := c.ShouldBindJSON(&sku)
	if err != nil {
		log.Print("Issue While Parsing JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.ErrParsingJSON,
		})
		return
	}

	if err := validators.ValidateSKU(&sku); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = ctrl.Repo.CreateSku(&sku)
	if err != nil {
		log.Print("Failed to create SKU: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": constants.ErrSKUCreate,
		})
		return
	}
	c.JSON(http.StatusCreated, sku)
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
		c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrSKUNotFound})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrServerError})
		return
	}

	c.JSON(http.StatusOK, updatedData)
}

func (ctrl *SkuController) DeleteSku(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Print("Invalid ID Parameter: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidID})
		return
	}

	err = ctrl.Repo.DeleteSku(uint(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrSKUNotFound})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SKU deleted successfully"})
}

func (ctrl *SkuController) GetSkusByTenantAndSeller(c *gin.Context) {
	var body struct {
		TenantID string   `json:"tenant_id"`
		SellerID string   `json:"seller_id"`
		SkuCodes []string `json:"sku_codes"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	skus, err := ctrl.Repo.GetSkusByTenantAndSeller(body.TenantID, body.SellerID, body.SkuCodes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get SKUs"})
		return
	}

	c.JSON(http.StatusOK, skus)
}

func (ctrl *SkuController) CreateSKUsBatch(c *gin.Context) {
	var skus []models.SKU
	err := c.ShouldBindJSON(&skus)
	if err != nil {
		log.Print("Issue While Parsing JSON: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.ErrParsingJSON,
		})
		return
	}

	if err := validators.ValidateBatchSize(len(skus)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	for i, sku := range skus {
		if err := validators.ValidateSKU(&sku); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     err.Error(),
				"sku_index": i,
			})
			return
		}
	}

	err = ctrl.Repo.CreateSKUsBatch(skus)
	if err != nil {
		log.Print("Batch SKU creation failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": constants.ErrBatchOperation,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "SKUs created successfully",
		"count":   len(skus),
	})
}

func (ctrl *SkuController) GetSKUsByIDs(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrParsingJSON})
		return
	}

	skus, err := ctrl.Repo.GetSKUsByIDs(request.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrServerError})
		return
	}

	c.JSON(http.StatusOK, skus)
}

func (ctrl *SkuController) GetSKUsByCodes(c *gin.Context) {
	var request struct {
		Codes []string `json:"codes"`
	}

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrParsingJSON})
		return
	}

	skus, err := ctrl.Repo.GetSKUsByCodes(request.Codes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrServerError})
		return
	}

	c.JSON(http.StatusOK, skus)
}
