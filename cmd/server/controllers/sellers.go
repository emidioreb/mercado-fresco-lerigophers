package controllers

import (
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/sellers"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type SellerController struct {
	service sellers.Service
}

type request struct {
	Cid         int    `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
}

func NewSeller(s sellers.Service) *SellerController {
	return &SellerController{
		service: s,
	}
}

func (s *SellerController) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData request

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("invalid request input"))
			return
		}

		seller, resp := s.service.Create(
			requestData.Cid,
			requestData.CompanyName,
			requestData.Address, requestData.Telephone,
		)

		if resp.Err != nil {
			c.JSON(resp.Code, gin.H{
				"error": resp.Err.Error(),
			})
			return
		}

		c.JSON(
			resp.Code,
			web.NewResponse(seller),
		)
		return
	}
}

func (s *SellerController) GetOne() gin.HandlerFunc {
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

		seller, err := s.service.GetOne(parsedId)

		if err != nil {
			c.JSON(
				http.StatusNotFound,
				web.DecodeError(err.Error()),
			)
			return
		}

		c.JSON(
			http.StatusOK,
			web.NewResponse(seller),
		)
		return
	}
}

func (s *SellerController) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		sellers, err := s.service.GetAll()

		if err != nil {
			c.JSON(
				http.StatusNotFound,
				web.DecodeError(err.Error()),
			)
			return
		}

		c.JSON(
			http.StatusOK,
			web.NewResponse(sellers),
		)
		return
	}
}

func (s *SellerController) Delete() gin.HandlerFunc {
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

		err = s.service.Delete(parsedId)
		if err != nil {
			c.JSON(http.StatusNotFound, web.DecodeError(err.Error()))
			return
		}

		c.JSON(http.StatusNoContent, web.NewResponse("seller with id "+id+" was deleted"))
	}
}

func (s *SellerController) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData request

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

		seller, resp := s.service.Update(
			parsedId,
			requestData.Cid,
			requestData.CompanyName,
			requestData.Address,
			requestData.Telephone,
		)

		if resp.Err != nil {
			c.JSON(resp.Code, web.DecodeError(resp.Err.Error()))
			return
		}

		c.JSON(resp.Code, web.NewResponse(seller))
		return
	}
}

func (s *SellerController) UpdateAddress() gin.HandlerFunc {
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

		seller, resp := s.service.UpdateAddress(id, requestData.Address)

		if resp.Err != nil {
			c.JSON(resp.Code, resp.Err.Error())
		}

		c.JSON(resp.Code, web.NewResponse(seller))
	}
}
