package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/warehouses"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type WarehouseController struct {
	service warehouses.Service
}

type reqWarehouses struct {
	WarehouseCode      string `json:"warehouse_code"`
	Address            string `json:"address"`
	Telephone          string `json:"telephone"`
	MinimumCapacity    int    `json:"minimum_capacity"`
	MinimumTemperature int    `json:"minimum_temperature"`
}

func NewWarehouse(s warehouses.Service) *WarehouseController {
	return &WarehouseController{
		service: s,
	}
}

func NewSellerHandler(r *gin.Engine, wh warehouses.Service) {
	controllerWarehouse := NewWarehouse(wh)
	sellerGroup := r.Group("/api/v1/sellers")
	{
		sellerGroup.GET("/:id", controllerWarehouse.GetOne())
		sellerGroup.GET("/", controllerWarehouse.GetAll())
		sellerGroup.POST("/", controllerWarehouse.Create())
		sellerGroup.DELETE("/:id", controllerWarehouse.Delete())
		sellerGroup.PATCH("/:id", controllerWarehouse.Update())
	}
}

func (s *WarehouseController) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData reqWarehouses

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("invalid request input"))
			return
		}

		if strings.ReplaceAll(requestData.WarehouseCode, " ", "") == "" {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("empty warehouse_code not allowed"))
			return
		}

		if len(requestData.WarehouseCode) > 255 {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("warehouse_code too long: max 255 characters"))
			return
		}

		if len(requestData.Address) > 255 {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("address too long: max 255 characters"))
			return
		}

		if len(requestData.Telephone) > 20 {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("telephone too long: max 20 characters"))
			return
		}

		warehouse, resp := s.service.Create(requestData.WarehouseCode, requestData.Address, requestData.Telephone, requestData.MinimumCapacity, requestData.MinimumTemperature)

		if resp.Err != nil {
			c.JSON(resp.Code, gin.H{
				"error": resp.Err.Error(),
			})
			return
		}

		c.JSON(resp.Code, web.NewResponse(warehouse))
	}
}

func (s *WarehouseController) GetOne() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		parsedId, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, web.DecodeError("id must be a number"))
			return
		}

		warehouse, resp := s.service.GetOne(parsedId)

		if resp.Err != nil {
			c.JSON(
				http.StatusNotFound,
				web.DecodeError(resp.Err.Error()),
			)
			return
		}

		c.JSON(http.StatusOK, web.NewResponse(warehouse))
	}
}

func (s *WarehouseController) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		warehousesList, resp := s.service.GetAll()

		if resp.Err != nil {
			c.JSON(
				resp.Code,
				web.DecodeError(resp.Err.Error()),
			)
			return
		}

		c.JSON(
			http.StatusOK,
			web.NewResponse(warehousesList),
		)
	}
}

func (s *WarehouseController) Delete() gin.HandlerFunc {
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

		c.JSON(resp.Code, web.NewResponse("warehouse with id "+id+" was deleted"))
	}
}

func (s *WarehouseController) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestValidatorType reqWarehouses
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

		if requestData["warehouse_code"] != nil {
			if strings.ReplaceAll(requestData["warehouse_code"].(string), " ", "") == "" {
				c.AbortWithStatusJSON(http.StatusBadRequest, web.DecodeError("empty warehouse_code not allowed"))
				return
			}
		}

		if value, ok := requestData["warehouse_code"].(string); ok {
			if len(value) > 255 {
				c.AbortWithStatusJSON(
					http.StatusUnprocessableEntity,
					web.DecodeError("warehouse_code too long: max 255 characters"),
				)
				return
			}
		}

		if value, ok := requestData["address"].(string); ok {
			if len(value) > 255 {
				c.AbortWithStatusJSON(
					http.StatusUnprocessableEntity,
					web.DecodeError("address too long: max 255 characters"),
				)
				return
			}
		}

		if value, ok := requestData["telephone"].(string); ok {
			if len(value) > 20 {
				c.AbortWithStatusJSON(
					http.StatusUnprocessableEntity,
					web.DecodeError("telephone too long: max 20 characters"),
				)
				return
			}
		}

		if value, ok := requestData["minimum_capacity"].(int); ok {
			if value < 0 {
				c.AbortWithStatusJSON(
					http.StatusUnprocessableEntity,
					web.DecodeError("minimum_capacity must be greather than 0"),
				)
				return
			}
		}

		warehouse, resp := s.service.Update(parsedId, requestData)
		if resp.Err != nil {
			c.JSON(resp.Code, web.DecodeError(resp.Err.Error()))
			return
		}

		c.JSON(resp.Code, web.NewResponse(warehouse))
	}
}
