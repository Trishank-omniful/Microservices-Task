package controllers

import (
	"fmt"
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

	if hubID != "" && skuID != "" {
		var hubIDUint, skuIDUint uint
		if _, err := fmt.Sscanf(hubID, "%d", &hubIDUint); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hub_id format"})
			return
		}
		if _, err := fmt.Sscanf(skuID, "%d", &skuIDUint); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sku_id format"})
			return
		}

		inventory, err := ctrl.Repo.GetInventoryByHubAndSKU(hubIDUint, skuIDUint)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get inventory"})
			return
		}
		c.JSON(http.StatusOK, inventory)
		return
	}

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

	if request.HubID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "hub_id is required"})
		return
	}

	inventories, err := ctrl.Repo.GetInventoryWithZeroDefaults(request.HubID, request.SKUIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrServerError})
		return
	}

	c.JSON(http.StatusOK, inventories)
}

func (ctrl *InventoryController) AtomicReduceInventory(c *gin.Context) {
	var request struct {
		HubID            uint `json:"hub_id"`
		SKUID            uint `json:"sku_id"`
		QuantityToReduce int  `json:"quantity_to_reduce"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrParsingJSON})
		return
	}

	if request.HubID == 0 || request.SKUID == 0 || request.QuantityToReduce <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "hub_id, sku_id, and positive quantity_to_reduce are required"})
		return
	}

	updatedInventory, err := ctrl.Repo.AtomicReduceInventory(request.HubID, request.SKUID, request.QuantityToReduce)
	if err != nil {
		if err.Error() == "insufficient inventory" {
			c.JSON(http.StatusConflict, gin.H{"error": constants.ErrInsufficientInventory})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrAtomicOperation})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "Inventory reduced successfully",
		"updated_inventory": updatedInventory,
	})
}

func (ctrl *InventoryController) CheckInventoryAvailability(c *gin.Context) {
	var request struct {
		HubID            uint `json:"hub_id"`
		SKUID            uint `json:"sku_id"`
		RequiredQuantity int  `json:"required_quantity"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrParsingJSON})
		return
	}

	if request.HubID == 0 || request.SKUID == 0 || request.RequiredQuantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "hub_id, sku_id, and positive required_quantity are required"})
		return
	}

	available, err := ctrl.Repo.CheckInventoryAvailability(request.HubID, request.SKUID, request.RequiredQuantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"hub_id":            request.HubID,
		"sku_id":            request.SKUID,
		"required_quantity": request.RequiredQuantity,
		"available":         available,
	})
}
