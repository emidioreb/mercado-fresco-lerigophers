package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/products"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
	SellerId                       int     `json:"seller_id"`
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

		if strings.ReplaceAll(requestData.ProductCode, " ", "") == "" {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("empty product_code not allowed"))
			return
		}

		product, resp := s.service.Create(requestData.ProductCode, requestData.Description,
			requestData.Width, requestData.Height, requestData.Length, requestData.NetWeight, requestData.ExpirationRate,
			requestData.RecommendedFreezingTemperature, requestData.FreezingRate, requestData.ProductTypeId, requestData.SellerId)

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
		var requestValidorType reqProducts
		requestData := make(map[string]interface{})
		id := c.Param("id")

		parsedId, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, web.DecodeError("id must be a number"))
			return
		}

		if err = c.ShouldBindBodyWith(&requestData, binding.JSON); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, web.DecodeError("invalid request data"))
			return
		}

		if len(requestData) == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, web.DecodeError("invalid request data - body needed"))
			return
		}

		if err := c.ShouldBindBodyWith(&requestValidorType, binding.JSON); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, web.DecodeError("invalid type of data"))
			return
		}

		if requestData["product_code"] != nil {
			if strings.ReplaceAll(requestData["product_code"].(string), " ", "") == "" {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("empty product_code not allowed"))
				return
			}
		}

		product, resp := s.service.Update(parsedId, requestData)

		if resp.Err != nil {
			c.JSON(resp.Code, web.DecodeError(resp.Err.Error()))
			return
		}

		c.JSON(resp.Code, web.NewResponse(product))
	}
}
