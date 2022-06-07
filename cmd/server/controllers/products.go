package controllers

import (
	"net/http"
	"strconv"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/products"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
)

type ProductController struct {
	service products.Service
}

type reqProducts struct {
	Id                             int     `json:"id"`
	ProductCode                    string  `json:"product_code"`
	Description                    string  `json:"description"`
	Width                          float64 `json:"width"`
	Height                         float64 `json:"height"`
	Length                         float64 `json:"length"`
	NetWeight                      float64 `json:"net_weight"`
	ExpirationRate                 float64 `json:"expiration_rate"`
	RecommendedFreezingTemperature float64 `json:"recommended_freezing_temperature"`
	FreezingRate                   float64 `json:"freezing_rate"`
	ProductTypeId                  int     `json:"product_type_id"`
}

func NewProduct(s products.Service) *ProductController {
	return &ProductController{
		service: s,
	}
}

func (s *ProductController) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData reqProducts

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("invalid request input"))
			return
		}

		product, resp := s.service.Create(requestData.ProductCode, requestData.Description,
			requestData.Width, requestData.Height, requestData.Length, requestData.NetWeight, requestData.ExpirationRate,
			requestData.RecommendedFreezingTemperature, requestData.FreezingRate, requestData.ProductTypeId)

		if resp.Err != nil {
			c.JSON(resp.Code, gin.H{
				"error": resp.Err.Error(),
			})
			return
		}

		c.JSON(
			resp.Code,
			web.NewResponse(product),
		)
	}
}

func (s *ProductController) GetOne() gin.HandlerFunc {
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

		product, resp := s.service.GetOne(parsedId)

		if resp.Err != nil {
			c.JSON(
				http.StatusNotFound,
				web.DecodeError(resp.Err.Error()),
			)
			return
		}

		c.JSON(
			http.StatusOK,
			web.NewResponse(product),
		)
	}
}

func (s *ProductController) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		ProductsList, resp := s.service.GetAll()

		if resp.Err != nil {
			c.JSON(
				resp.Code,
				web.DecodeError(resp.Err.Error()),
			)
			return
		}

		c.JSON(
			http.StatusOK,
			web.NewResponse(ProductsList),
		)
	}
}

func (s *ProductController) Delete() gin.HandlerFunc {
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

		c.JSON(resp.Code, web.NewResponse("Product with id "+id+" was deleted"))
	}
}

func (s *ProductController) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData reqProducts

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

		product, resp := s.service.Update(
			parsedId,
			requestData.ProductCode,
			requestData.Description,
			requestData.Width,
			requestData.Height,
			requestData.Length,
			requestData.NetWeight,
			requestData.ExpirationRate,
			requestData.RecommendedFreezingTemperature,
			requestData.FreezingRate,
			requestData.ProductTypeId,
		)

		if resp.Err != nil {
			c.JSON(resp.Code, web.DecodeError(resp.Err.Error()))
			return
		}

		c.JSON(resp.Code, web.NewResponse(product))
	}
}

func (s *ProductController) UpdateExpirationRate() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, web.DecodeError("invalid id"))
			return
		}

		var requestData reqProducts
		err = c.ShouldBindJSON(&requestData)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, web.DecodeError("invalid request data"))
			return
		}

		product, resp := s.service.UpdateExpirationRate(id, requestData.ExpirationRate)

		if resp.Err != nil {
			c.JSON(resp.Code, resp.Err.Error())
		}

		c.JSON(resp.Code, web.NewResponse(product))
	}
}
