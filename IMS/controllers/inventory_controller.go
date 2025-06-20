package controllers

import (
	"net/http"

	"github.com/Trishank-Omniful/Onboarding-Task/constants"
	"github.com/Trishank-Omniful/Onboarding-Task/models"
	"github.com/Trishank-Omniful/Onboarding-Task/repository"
	"github.com/Trishank-Omniful/Onboarding-Task/validators"
	"github.com/gin-gonic/gin"
)

type InventoryController struct {
	Repo *repository.InventoryRepository
}

func NewInventoryController(repo *repository.InventoryRepository) *InventoryController {
	return &InventoryController{Repo: repo}
}

func (ctrl *InventoryController) UpsertInventory(c *gin.Context) {
	var inventory models.Inventory
	if err := c.ShouldBindJSON(&inventory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrParsingJSON})
		return
	}

	if err := validators.ValidateInventory(&inventory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.Repo.UpsertInventory(&inventory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upsert inventory"})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

func (ctrl *InventoryController) GetInventory(c *gin.Context) {
	hubID := c.Query("hub_id")
	skuID := c.Query("sku_id")

	inventory, err := ctrl.Repo.GetInventory(hubID, skuID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get inventory"})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

func (ctrl *InventoryController) UpsertInventoryBatch(c *gin.Context) {
	var inventories []models.Inventory
	if err := c.ShouldBindJSON(&inventories); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrParsingJSON})
		return
	}

	if err := validators.ValidateBatchSize(len(inventories)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, inventory := range inventories {
		if err := validators.ValidateInventory(&inventory); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":           err.Error(),
				"inventory_index": i,
			})
			return
		}
	}

	if err := ctrl.Repo.UpsertInventoryBatch(inventories); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrBatchOperation})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Inventory batch operation completed successfully",
		"count":   len(inventories),
	})
}

func (ctrl *InventoryController) GetInventoriesByHubAndSKUs(c *gin.Context) {
	var request struct {
		HubID  uint   `json:"hub_id"`
		SKUIDs []uint `json:"sku_ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrParsingJSON})
		return
	}

	inventories, err := ctrl.Repo.GetInventoriesByHubAndSKUs(request.HubID, request.SKUIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrServerError})
		return
	}

	c.JSON(http.StatusOK, inventories)
}
