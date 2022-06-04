package controllers

import (
	"net/http"
	"strconv"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/sections"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
)

type SectionController struct {
	service sections.Service
}

type reqSections struct {
	SectionNumber      int `json:"section_number"`
	CurrentTemperature int `json:"current_temperature"`
	MinimumTemperature int `json:"minimum_temperature"`
	CurrentCapacity    int `json:"current_capacity"`
	MininumCapacity    int `json:"minimum_capacity"`
	MaximumCapacity    int `json:"maximum_capacity"`
	WarehouseId        int `json:"warehouse_id"`
	ProductTypeId      int `json:"product_type_id"`
}

func NewSection(s sections.Service) *SectionController {
	return &SectionController{
		service: s,
	}
}

func (s *SectionController) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData reqSections

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("invalid request input"))
			return
		}

		section, resp := s.service.Create(
			requestData.SectionNumber,
			requestData.CurrentTemperature,
			requestData.MinimumTemperature,
			requestData.CurrentCapacity,
			requestData.MininumCapacity,
			requestData.MaximumCapacity,
			requestData.WarehouseId,
			requestData.ProductTypeId,
		)

		if resp.Err != nil {
			c.JSON(resp.Code, gin.H{
				"error": resp.Err.Error(),
			})
			return
		}

		c.JSON(
			resp.Code,
			web.NewResponse(section),
		)
		return
	}
}

func (s *SectionController) GetOne() gin.HandlerFunc {
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

		section, resp := s.service.GetOne(parsedId)

		if resp.Err != nil {
			c.JSON(
				http.StatusNotFound,
				web.DecodeError(resp.Err.Error()),
			)
			return
		}

		c.JSON(
			http.StatusOK,
			web.NewResponse(section),
		)
		return
	}
}

func (s *SectionController) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		sectionsList, resp := s.service.GetAll()

		if resp.Err != nil {
			c.JSON(
				resp.Code,
				web.DecodeError(resp.Err.Error()),
			)
			return
		}

		c.JSON(
			http.StatusOK,
			web.NewResponse(sectionsList),
		)
		return
	}
}

func (s *SectionController) Delete() gin.HandlerFunc {
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

		c.JSON(resp.Code, web.NewResponse("section with id "+id+" was deleted"))
	}
}

func (s *SectionController) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData reqSections

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

		section, resp := s.service.Update(
			parsedId,
			requestData.SectionNumber,
			requestData.CurrentTemperature,
			requestData.MinimumTemperature,
			requestData.CurrentCapacity,
			requestData.MininumCapacity,
			requestData.MaximumCapacity,
			requestData.WarehouseId,
			requestData.ProductTypeId,
		)

		if resp.Err != nil {
			c.JSON(resp.Code, web.DecodeError(resp.Err.Error()))
			return
		}

		c.JSON(resp.Code, web.NewResponse(section))
		return
	}
}

// func (s *SectionController) UpdateAddress() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		id, err := strconv.Atoi(c.Param("id"))
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, web.DecodeError("invalid id"))
// 			return
// 		}

// 		var requestData request
// 		err = c.ShouldBindJSON(&requestData)

// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusBadRequest, web.DecodeError("invalid request data"))
// 			return
// 		}

// 		section, resp := s.service.UpdateAddress(id, requestData.Address)

// 		if resp.Err != nil {
// 			c.JSON(resp.Code, resp.Err.Error())
// 		}

// 		c.JSON(resp.Code, web.NewResponse(section))
// 	}
// }
