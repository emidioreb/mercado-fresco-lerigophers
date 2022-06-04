package controllers

import (
	"net/http"
	"strconv"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/warehouses"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
)

type WarehouseController struct {
	service warehouses.Service
}

type reqWarehouses struct {
	WarehouseCode      string `json:"warehouse_code"`
	Address            string `json:"adress"`
	Telephone          string `json:"telephone"`
	MinimumCapacity    int    `json:"minimum_capacity"`
	MaximumTemperature int    `json:"maximum_temperature"`
}

func NewWarehouse(s warehouses.Service) *WarehouseController {
	return &WarehouseController{
		service: s,
	}
}

func (s *WarehouseController) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData reqWarehouses

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("invalid request input"))
			return
		}

		warehouse, resp := s.service.Create(requestData.WarehouseCode, requestData.Address, requestData.Telephone, requestData.MinimumCapacity, requestData.MaximumTemperature)

		if resp.Err != nil {
			c.JSON(resp.Code, gin.H{
				"error": resp.Err.Error(),
			})
			return
		}

		c.JSON(
			resp.Code,
			web.NewResponse(warehouse),
		)
		return
	}
}

func (s *WarehouseController) GetOne() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if id == "" {
			c.JSON(http.StatusBadRequest, web.DecodeError("id must be informed"))
			return
		}

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

		c.JSON(
			http.StatusOK,
			web.NewResponse(warehouse),
		)
		return
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
		return
	}
}

func (s *WarehouseController) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		if id == "" {
			c.JSON(http.StatusBadRequest, web.DecodeError("id must be informed"))
			return
		}

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
		var requestData reqWarehouses

		id := c.Param("id")

		if id == "" {
			c.JSON(http.StatusBadRequest, web.DecodeError("id must be informed"))
			return
		}

		parsedId, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, web.DecodeError("id must be a number"))
			return
		}

		err = c.ShouldBindJSON(&requestData)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, web.DecodeError("invalid request data"))
			return
		}

		warehouse, resp := s.service.Update(
			parsedId,
			requestData.WarehouseCode,
			requestData.Address,
			requestData.Telephone,
			requestData.MinimumCapacity,
			requestData.MaximumTemperature,
		)

		if resp.Err != nil {
			c.JSON(resp.Code, web.DecodeError(resp.Err.Error()))
			return
		}

		c.JSON(resp.Code, web.NewResponse(warehouse))
		return
	}
}

func (s *WarehouseController) UpdateTelephone() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, web.DecodeError("invalid id"))
			return
		}

		var requestData request
		err = c.ShouldBindJSON(&requestData)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, web.DecodeError("invalid request data"))
			return
		}

		warehouse, resp := s.service.UpdateTelephone(id, requestData.Telephone)

		if resp.Err != nil {
			c.JSON(resp.Code, resp.Err.Error())
		}

		c.JSON(resp.Code, web.NewResponse(warehouse))
	}
}
