package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin/binding"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/sellers"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
)

type SellerController struct {
	service sellers.Service
}

type reqSellersCreate struct {
	Cid         int    `json:"cid" binding:"required"`
	CompanyName string `json:"company_name" binding:"required"`
	Address     string `json:"address" binding:"required"`
	Telephone   string `json:"telephone" binding:"required"`
	LocalityId  string `json:"locality_id" binding:"required"`
}

type reqSellersUpdate struct {
	Cid         int    `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
	LocalityId  string `json:"locality_id"`
}

func NewSeller(s sellers.Service) *SellerController {
	return &SellerController{
		service: s,
	}
}

func NewSellerHandler(r *gin.RouterGroup, ss sellers.Service) {
	sellerController := NewSeller(ss)
	sellerGroup := r.Group("/sellers")
	{
		sellerGroup.GET("/:id", sellerController.GetOne())
		sellerGroup.GET("/", sellerController.GetAll())
		sellerGroup.POST("/", sellerController.Create())
		sellerGroup.DELETE("/:id", sellerController.Delete())
		sellerGroup.PATCH("/:id", sellerController.Update())
	}
}

func (s *SellerController) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData reqSellersCreate

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("invalid request input"))
			return
		}

		if len(requestData.CompanyName) > 255 {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("company_name too long: max 255 characters"))
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

		if len(requestData.LocalityId) > 255 {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("locality_id too long: max 255 characters"))
			return
		}

		seller, resp := s.service.Create(
			requestData.Cid,
			requestData.CompanyName,
			requestData.Address,
			requestData.Telephone,
			requestData.LocalityId,
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
	}
}

func (s *SellerController) GetOne() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		parsedId, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, web.DecodeError("id must be a number"))
			return
		}

		seller, resp := s.service.GetOne(parsedId)

		if resp.Err != nil {
			c.JSON(
				http.StatusNotFound,
				web.DecodeError(resp.Err.Error()),
			)
			return
		}

		c.JSON(
			http.StatusOK,
			web.NewResponse(seller),
		)
	}
}

func (s *SellerController) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		sellersList, resp := s.service.GetAll()

		if resp.Err != nil {
			c.JSON(
				resp.Code,
				web.DecodeError(resp.Err.Error()),
			)
			return
		}

		c.JSON(
			http.StatusOK,
			web.NewResponse(sellersList),
		)
	}
}

func (s *SellerController) Delete() gin.HandlerFunc {
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

		c.JSON(resp.Code, web.NewResponse("seller with id "+id+" was deleted"))
	}
}

func (s *SellerController) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestValidatorType reqSellersUpdate
		var requestData map[string]interface{}

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

		if value, ok := requestData["cid"].(float64); ok {
			if value < 1 {
				c.AbortWithStatusJSON(
					http.StatusUnprocessableEntity,
					web.DecodeError("cid must be greather than 0"),
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

		if value, ok := requestData["company_name"].(string); ok {
			if len(value) > 255 {
				c.AbortWithStatusJSON(
					http.StatusUnprocessableEntity,
					web.DecodeError("company_name too long: max 255 characters"),
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

		if value, ok := requestData["locality_id"].(string); ok {
			if len(value) > 255 {
				c.AbortWithStatusJSON(
					http.StatusUnprocessableEntity,
					web.DecodeError("locality_id too long: max 255 characters"),
				)
				return
			}
		}

		seller, resp := s.service.Update(parsedId, requestData)

		if resp.Err != nil {
			c.JSON(resp.Code, web.DecodeError(resp.Err.Error()))
			return
		}

		c.JSON(resp.Code, web.NewResponse(seller))
	}
}
