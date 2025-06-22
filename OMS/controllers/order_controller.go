package controllers

import (
	"github.com/Trishank-omniful/Onboarding-Task/clients"
	"github.com/gin-gonic/gin"
)

type OrderController struct {
	s3Client clients.S3ClientInterface
}

func NewOrderController(s3Client clients.S3ClientInterface) *OrderController {
	return &OrderController{
		s3Client: s3Client,
	}
}

func (c *OrderController) BulkUploadCSV(g *gin.Context) {

}
