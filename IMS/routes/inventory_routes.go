package routes

import (
	"github.com/Trishank-Omniful/Onboarding-Task/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterInventoryRoutes(router *gin.RouterGroup, ctrl *controllers.InventoryController) {
	router.POST("/inventory", ctrl.UpsertInventory)
	router.GET("/inventory", ctrl.GetInventory)
	router.POST("/inventory/batch", ctrl.UpsertInventoryBatch)
	router.POST("/inventory/batch/hub-skus", ctrl.GetInventoriesByHubAndSKUs)
	router.POST("/inventory/atomic/reduce", ctrl.AtomicReduceInventory)
	router.POST("/inventory/check-availability", ctrl.CheckInventoryAvailability)
}
