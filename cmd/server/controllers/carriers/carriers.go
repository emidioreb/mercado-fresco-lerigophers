package controllers

import (
	"net/http"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/carriers"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
)

type CarryController struct {
	service carriers.Service
}

type reqCarries struct {
	Cid         string `json:"cid" binding:"required"`
	CompanyName string `json:"company_name" binding:"required"`
	Address     string `json:"address" binding:"required"`
	Telephone   string `json:"telephone" binding:"required"`
	LocalityId  string `json:"locality_id" binding:"required"`
}

func NewCarry(s carriers.Service) *CarryController {
	return &CarryController{
		service: s,
	}
}

func NewCarryHandler(r *gin.Engine, cs carriers.Service) {
	carryController := NewCarry(cs)
	r.POST("/api/v1/carries/", carryController.Create())
}

func (s *CarryController) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData reqCarries

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("invalid request input"))
			return
		}

		if len(requestData.Cid) > 255 {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("CID too long: max 255 characters"))
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

		carry, resp := s.service.Create(
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
			web.NewResponse(carry),
		)
	}
}
