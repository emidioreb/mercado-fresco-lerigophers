package controllers

import (
	"net/http"
	"time"

	purchase_orders "github.com/emidioreb/mercado-fresco-lerigophers/internal/purchaseOrders"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
)

type PurchaseOrdersController struct {
	service purchase_orders.Service
}

type reqPurchaseOrders struct {
	OrderNumber     string `json:"order_number" binding:"required"`
	OrderDate       string `json:"order_date" binding:"required"`
	TrackingCode    string `json:"tracking_code"`
	BuyerId         int    `json:"buyer_id" binding:"required"`
	ProductRecordId int    `json:"product_record_id" binding:"required"`
	OrderStatusId   int    `json:"order_status_id" binding:"required"`
}

func NewPurchaseOrder(s purchase_orders.Service) *PurchaseOrdersController {
	return &PurchaseOrdersController{
		service: s,
	}
}

func (s *PurchaseOrdersController) CreatePurchaseOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData reqPurchaseOrders

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("invalid request input"))
			return
		}

		if requestData.BuyerId < 0 {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("buyer_id must be greather than 0"))
			return
		}

		if requestData.ProductRecordId < 0 {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("product_record_id must be greather than 0"))
			return
		}

		if requestData.OrderStatusId < 0 {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("order_status_id must be greather than 0"))
			return
		}

		const layout = "2006-01-02"

		orderDate, errDate := time.Parse(layout, requestData.OrderDate)

		if errDate != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("order_date format incorrect, model: YYYY-MM-DD"))
			return
		}

		purchaseOrder, resp := s.service.CreatePurchaseOrders(
			requestData.OrderNumber,
			orderDate,
			requestData.TrackingCode,
			requestData.BuyerId,
			requestData.ProductRecordId,
			requestData.OrderStatusId,
		)

		if resp.Err != nil {
			c.JSON(resp.Code, gin.H{
				"error": resp.Err.Error(),
			})
			return
		}

		c.JSON(
			resp.Code,
			web.NewResponse(purchaseOrder),
		)
	}
}