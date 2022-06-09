package controllers

import (
	"net/http"
	"strconv"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/sections"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
		var requestValidatorType reqSections
		requestData := make(map[string]interface{})
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

		section, resp := s.service.Update(parsedId, requestData)

		if resp.Err != nil {
			c.JSON(resp.Code, web.DecodeError(resp.Err.Error()))
			return
		}

		c.JSON(resp.Code, web.NewResponse(section))
	}
}
