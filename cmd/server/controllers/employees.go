package controllers

import (
	"net/http"
	"strconv"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/employees"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
)

type EmployeeController struct {
	service employees.Service
}
type reqEmployee struct {
	CardNumberId string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseId  int    `json:"warehouse_id"`
}

func NewEmployee(s employees.Service) *EmployeeController {
	return &EmployeeController{
		service: s,
	}
}

func (s *EmployeeController) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData reqEmployee

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("invalid request input"))
			return
		}

		seller, resp := s.service.Create(
			requestData.CardNumberId,
			requestData.FirstName,
			requestData.LastName, requestData.WarehouseId,
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

func (s *EmployeeController) GetOne() gin.HandlerFunc {
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

		employee, resp := s.service.GetOne(parsedId)

		if resp.Err != nil {
			c.JSON(
				http.StatusNotFound,
				web.DecodeError(resp.Err.Error()),
			)
			return
		}

		c.JSON(
			http.StatusOK,
			web.NewResponse(employee),
		)
	}
}

func (s *EmployeeController) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		employeesList, resp := s.service.GetAll()

		if resp.Err != nil {
			c.JSON(
				resp.Code,
				web.DecodeError(resp.Err.Error()),
			)
			return
		}

		c.JSON(
			http.StatusOK,
			web.NewResponse(employeesList),
		)
	}
}

func (s *EmployeeController) Delete() gin.HandlerFunc {
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

		c.JSON(resp.Code, web.NewResponse("seller with id "+id+" was deleted"))
	}
}

func (s *EmployeeController) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData reqEmployee

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

		employee, resp := s.service.Update(
			parsedId,
			requestData.CardNumberId,
			requestData.FirstName,
			requestData.LastName,
			requestData.WarehouseId,
		)

		if resp.Err != nil {
			c.JSON(resp.Code, web.DecodeError(resp.Err.Error()))
			return
		}

		c.JSON(resp.Code, web.NewResponse(employee))
	}
}

func (s *EmployeeController) UpdateFirstName() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, web.DecodeError("invalid id"))
			return
		}

		var requestData reqEmployee
		err = c.ShouldBindJSON(&requestData)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, web.DecodeError("invalid request data"))
			return
		}

		employee, resp := s.service.UpdateFirstName(id, requestData.FirstName)

		if resp.Err != nil {
			c.JSON(resp.Code, resp.Err.Error())
		}

		c.JSON(resp.Code, web.NewResponse(employee))
	}
}
