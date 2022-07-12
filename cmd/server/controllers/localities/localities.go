package controllers

import (
	"net/http"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/localities"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
)

type LocalityController struct {
	service localities.Service
}

type reqLocality struct {
	Id           string `json:"id" binding:"required"`
	LocalityName string `json:"locality_name" binding:"required"`
	ProvinceName string `json:"province_name" binding:"required"`
	CountryName  string `json:"country_name" binding:"required"`
}

func NewLocality(s localities.Service) *LocalityController {
	return &LocalityController{
		service: s,
	}
}

func (s *LocalityController) CreateLocality() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData reqLocality

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("invalid request input"))
			return
		}

		if len(requestData.Id) > 255 {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("id too long: max 255 characters"))
			return
		}

		if len(requestData.CountryName) > 255 {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("country_name too long: max 255 characters"))
			return
		}

		if len(requestData.LocalityName) > 255 {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("locality_name too long: max 255 characters"))
			return
		}

		if len(requestData.ProvinceName) > 255 {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("province_name too long: max 255 characters"))
			return
		}

		seller, resp := s.service.CreateLocality(
			requestData.Id,
			requestData.LocalityName,
			requestData.ProvinceName,
			requestData.CountryName,
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

func (s *LocalityController) GetReportSellers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			reportSellers []localities.ReportSellers
			resp          web.ResponseCode
		)
		id := c.Query("id")
		if id != "" {
			reportSellers, resp = s.service.GetReportOneSeller(id)
		} else {
			reportSellers, resp = s.service.GetAllReportSellers()
		}

		if resp.Err != nil {
			c.JSON(
				resp.Code,
				web.DecodeError(resp.Err.Error()),
			)
			return
		}

		c.JSON(
			http.StatusOK,
			web.NewResponse(reportSellers),
		)
	}
}

func (s *LocalityController) GetReportCarriers() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		reportCarries, resp := s.service.GetReportCarriers(id)
		if resp.Err != nil {
			c.JSON(
				http.StatusNotFound,
				web.DecodeError(resp.Err.Error()),
			)
			return
		}

		c.JSON(
			http.StatusOK,
			web.NewResponse(reportCarries),
		)
	}

}
