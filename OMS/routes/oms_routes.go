package routes

import "github.com/gin-gonic/gin"

func RegisterOMSRoutes(router *gin.RouterGroup) {
	omsGroup := router.Group("oms")
	{
		omsGroup.POST("/orders/buld-upload", orderCtrl.BulkUploadCSV)
	}
}
