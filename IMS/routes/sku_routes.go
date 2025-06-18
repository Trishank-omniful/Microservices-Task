package routes

import (
	"github.com/Trishank-Omniful/Onboarding-Task/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterSkuRoutes(router *gin.RouterGroup, controller *controllers.SkuController) {
	hubGroup := router.Group("/sku")
	{
		hubGroup.GET("", controller.GetAllSkus)
		hubGroup.GET("/:id", controller.GetSkuById)
		hubGroup.POST("", controller.CreateSku)
		hubGroup.PUT("/:id", controller.UpdateSku)
		hubGroup.DELETE("/:id", controller.DeleteSku)
	}
}
