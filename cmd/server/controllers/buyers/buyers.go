package buyers_controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/buyers"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type BuyerController struct {
	service buyers.Service
}

type reqBuyers struct {
	Id           int    `json:"id"`
	CardNumberId string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}

func NewBuyer(s buyers.Service) *BuyerController {
	return &BuyerController{
		service: s,
	}
}

func NewBuyerHandler(r *gin.Engine, bs buyers.Service) {
	buyerController := NewBuyer(bs)
	buyerGroup := r.Group("/api/v1/buyers")
	{
		buyerGroup.GET("/:id", buyerController.GetOne())
		buyerGroup.GET("/", buyerController.GetAll())
		buyerGroup.POST("/", buyerController.Create())
		buyerGroup.DELETE("/:id", buyerController.Delete())
		buyerGroup.PATCH("/:id", buyerController.Update())
		buyerGroup.GET("/reportPurchaseOrders", buyerController.GetReportPurchaseOrders())
	}
}

func (s *BuyerController) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData reqBuyers

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("invalid request input"))
			return
		}

		if strings.ReplaceAll(requestData.CardNumberId, " ", "") == "" {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("empty card_number_id not allowed"))
			return
		}

		buyer, resp := s.service.Create(requestData.CardNumberId, requestData.FirstName, requestData.LastName)

		if resp.Err != nil {
			c.JSON(resp.Code, gin.H{
				"error": resp.Err.Error(),
			})
			return
		}

		c.JSON(resp.Code,
			web.NewResponse(buyer),
		)
	}
}

func (s *BuyerController) GetOne() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		parsedId, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, web.DecodeError("id must be a number"))
			return
		}

		buyer, resp := s.service.GetOne(parsedId)

		if resp.Err != nil {
			c.JSON(
				http.StatusNotFound,
				web.DecodeError(resp.Err.Error()),
			)
			return
		}

		c.JSON(
			http.StatusOK,
			web.NewResponse(buyer),
		)
	}
}

func (s *BuyerController) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		BuyersList, resp := s.service.GetAll()

		if resp.Err != nil {
			c.JSON(
				resp.Code,
				web.DecodeError(resp.Err.Error()),
			)
			return
		}

		c.JSON(
			http.StatusOK,
			web.NewResponse(BuyersList),
		)
	}
}

func (s *BuyerController) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		parsedId, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, web.DecodeError("id must be a number"))
			return
		}

		resp := s.service.Delete(parsedId)
		if resp.Err != nil {
			c.JSON(resp.Code, web.DecodeError(resp.Err.Error()))
			return
		}

		c.JSON(resp.Code, web.NewResponse("Buyer with id "+id+" was deleted"))
	}
}

func (s *BuyerController) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestValidatorType reqBuyers
		requestData := make(map[string]interface{})

		id := c.Param("id")

		parsedId, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, web.DecodeError("id must be a number"))
			return
		}

		if err := c.ShouldBindBodyWith(&requestData, binding.JSON); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, web.DecodeError("invalid request data"))
			return
		}

		if len(requestData) == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, web.DecodeError("invalid request data - body needed"))
			return
		}

		if err := c.ShouldBindBodyWith(&requestValidatorType, binding.JSON); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, web.DecodeError("invalid type of data"))
			return
		}

		if requestData["card_number_id"] != nil {
			if strings.ReplaceAll(requestData["card_number_id"].(string), " ", "") == "" {
				c.AbortWithStatusJSON(http.StatusBadRequest, web.DecodeError("empty card_number_id not allowed"))
				return
			}
		}

		buyer, resp := s.service.Update(parsedId, requestData)

		if resp.Err != nil {
			c.JSON(resp.Code, web.DecodeError(resp.Err.Error()))
			return
		}

		c.JSON(resp.Code, web.NewResponse(buyer))
	}
}

func (s *BuyerController) GetReportPurchaseOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			reportPurchaseOrders []buyers.ReportPurchaseOrders
			resp                 web.ResponseCode
		)

		id := c.Query("id")
		if id != "" {
			parsedId, err := strconv.Atoi(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, web.DecodeError("id must be a number"))
				return
			}
			reportPurchaseOrders, resp = s.service.GetReportPurchaseOrders(parsedId)

		} else {
			reportPurchaseOrders, resp = s.service.GetReportPurchaseOrders(0)
		}

		if resp.Err != nil {
			c.JSON(
				resp.Code,
				web.DecodeError(resp.Err.Error()),
			)
			return
		}

		c.JSON(
			resp.Code,
			web.NewResponse(reportPurchaseOrders),
		)

	}

}
