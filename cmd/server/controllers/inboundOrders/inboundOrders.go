package inboundorders

import (
	"net/http"
	"time"

	inboundorders "github.com/emidioreb/mercado-fresco-lerigophers/internal/inboundOrders"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
)

type InboundOrdersController struct {
	service inboundorders.Service
}

type reqInboundOrder struct {
	OrderNumber    string `json:"order_number" binding:"required"`
	OrderDate      string `json:"order_date" binding:"required"`
	EmployeeId     int    `json:"employee_id" binding:"required"`
	ProductBatchId int    `json:"product_batch_id" binding:"required"`
	WarehouseId    int    `json:"warehouse_id" binding:"required"`
}

func NewInboud(s inboundorders.Service) *InboundOrdersController {
	return &InboundOrdersController{
		service: s,
	}
}

func (s *InboundOrdersController) CreateInboundOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData reqInboundOrder
		const layout = "2006-01-02"

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("invalid request input"))
			return
		}

		if len(requestData.OrderNumber) > 255 {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("order_number too long: max 255 characters"))
			return
		}

		_, errDate := time.Parse(layout, requestData.OrderDate)

		if errDate != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("order_date format incorrect, model: YYYY-MM-DD"))
			return
		}

		inboundOrders, resp := s.service.CreateInboundOrders(
			requestData.OrderNumber,
			requestData.OrderDate,
			requestData.EmployeeId,
			requestData.ProductBatchId,
			requestData.WarehouseId,
		)

		if resp.Err != nil {
			c.JSON(resp.Code, gin.H{
				"error": resp.Err.Error(),
			})
			return
		}

		c.JSON(
			resp.Code,
			web.NewResponse(inboundOrders),
		)
	}
}

func (s *InboundOrdersController) GetReportInboundOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		reportInbounds, resp := s.service.GetReportInboundOrders(id)
		if resp.Err != nil {
			c.JSON(
				http.StatusInternalServerError,
				web.DecodeError(resp.Err.Error()),
			)
			return
		}

		c.JSON(
			http.StatusOK,
			web.NewResponse(reportInbounds),
		)
	}

}
