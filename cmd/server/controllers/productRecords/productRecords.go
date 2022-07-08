package controllers

import (
	"net/http"

	product_records "github.com/emidioreb/mercado-fresco-lerigophers/internal/productRecords"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
)

type ProductRecordController struct {
	service product_records.Service
}

type reqProductRecord struct {
	LastUpdateDate string  `json:"last_update_date" binding:"required"`
	PurchasePrice  float64 `json:"purchase_price" binding:"required"`
	SalePrice      float64 `json:"sale_price" binding:"required"`
	ProductId      int     `json:"product_id" binding:"required"`
}

func NewProductRecord(s product_records.Service) *ProductRecordController {
	return &ProductRecordController{
		service: s,
	}
}

func (s *ProductRecordController) CreateProductRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData reqProductRecord

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("invalid request input"))
		}

		productRecord, resp := s.service.CreateProductRecord(
			requestData.LastUpdateDate,
			requestData.PurchasePrice,
			requestData.SalePrice,
			requestData.ProductId,
		)

		if resp.Err != nil {
			c.JSON(resp.Code, gin.H{
				"error": resp.Err.Error(),
			})
			return
		}

		c.JSON(
			resp.Code,
			web.NewResponse(productRecord),
		)

	}
}
