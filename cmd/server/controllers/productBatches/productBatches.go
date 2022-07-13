package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	product_batches "github.com/emidioreb/mercado-fresco-lerigophers/internal/productBatches"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
)

type ProductBatchController struct {
	service product_batches.Service
}

type ReqProductBatch struct {
	BatchNumber        int    `json:"batch_number" binding:"required"`
	CurrentQuantity    int    `json:"current_quantity" binding:"required"`
	CurrentTemperature int    `json:"current_temperature" binding:"required"`
	InitialQuantity    int    `json:"initial_quantity" binding:"required"`
	ManufacturingHour  int    `json:"manufacturing_hour" binding:"required"`
	MinimumTemperature int    `json:"minumum_temperature" binding:"required"`
	ProductId          int    `json:"product_id" binding:"required"`
	SectionId          int    `json:"section_id" binding:"required"`
	DueDate            string `json:"due_date" binding:"required"`
	ManufacturingDate  string `json:"manufacturing_date" binding:"required"`
}

func NewProductBatch(s product_batches.Service) *ProductBatchController {
	return &ProductBatchController{
		service: s,
	}
}

func NewProductBatchHandler(r *gin.Engine, pb product_batches.Service) {
	controllerProductBatches := NewProductBatch(pb)
	ProductBatchesGroup := r.Group("/productBatches")
	{
		ProductBatchesGroup.POST("/", controllerProductBatches.CreateProductBatch())
		ProductBatchesGroup.GET("/reportProducts", controllerProductBatches.GetReportSection())
	}
}

func (s *ProductBatchController) CreateProductBatch() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData ReqProductBatch

		if err := c.ShouldBindJSON(&requestData); err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("invalid request input"))
			return
		}

		if requestData.BatchNumber < 0 {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("batch_number must be greather than 0"))
			return
		}

		if requestData.CurrentQuantity < 0 {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("current_quantity must be greather than 0"))
			return
		}

		if requestData.InitialQuantity < 0 {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("initial_quantity must be greather than 0"))
			return
		}

		if requestData.ProductId < 0 {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("product_id must be greather than 0"))
			return
		}

		if requestData.SectionId < 0 {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("section_id must be greather than 0"))
			return
		}

		const layout = "2006-01-02"

		duedate, errDate := time.Parse(layout, requestData.DueDate)

		if errDate != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("due_date format incorrect, model: YYYY-MM-DD"))
			return
		}

		manufacturingdate, errDate := time.Parse(layout, requestData.ManufacturingDate)
		if errDate != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("manufacturing_date format incorrect, model: YYYY-MM-DD"))
			return
		}

		productBatch, resp := s.service.CreateProductBatch(
			requestData.BatchNumber,
			requestData.CurrentQuantity,
			requestData.CurrentTemperature,
			requestData.InitialQuantity,
			requestData.ManufacturingHour,
			requestData.MinimumTemperature,
			requestData.ProductId,
			requestData.SectionId,
			duedate,
			manufacturingdate,
		)

		if resp.Err != nil {
			c.JSON(resp.Code, gin.H{
				"error": resp.Err.Error(),
			})
			return
		}

		c.JSON(
			resp.Code,
			web.NewResponse(productBatch),
		)
	}
}

func (s *ProductBatchController) GetReportSection() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		if id != "" {
			parsedId, err := strconv.Atoi(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, web.DecodeError("id must be a number"))
				return
			}

			reportSections, resp := s.service.GetReportSection(parsedId)
			if resp.Err != nil {
				c.JSON(
					http.StatusNotFound,
					web.DecodeError(resp.Err.Error()),
				)
				return
			}

			c.JSON(
				http.StatusOK,
				web.NewResponse(reportSections),
			)
		} else {
			reportSections, resp := s.service.GetReportSection(0)
			if resp.Err != nil {
				c.JSON(
					http.StatusNotFound,
					web.DecodeError(resp.Err.Error()),
				)
				return
			}

			c.JSON(
				http.StatusOK,
				web.NewResponse(reportSections),
			)
		}

	}

}
