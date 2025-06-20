package routes

import (
	"github.com/Trishank-Omniful/Onboarding-Task/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterSkuRoutes(router *gin.RouterGroup, controller *controllers.SkuController) {
	skuGroup := router.Group("/sku")
	{
		skuGroup.GET("", controller.GetAllSkus)
		skuGroup.GET("/:id", controller.GetSkuById)
		skuGroup.POST("", controller.CreateSku)
		skuGroup.PUT("/:id", controller.UpdateSku)
		skuGroup.DELETE("/:id", controller.DeleteSku)
		skuGroup.POST("/filter", controller.GetSkusByTenantAndSeller)
		skuGroup.POST("/batch", controller.CreateSKUsBatch)
		skuGroup.POST("/batch/ids", controller.GetSKUsByIDs)
		skuGroup.POST("/batch/codes", controller.GetSKUsByCodes)
	}
}
