package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/employees"
	"github.com/emidioreb/mercado-fresco-lerigophers/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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

func NewEmployeeHandler(r *gin.Engine, es employees.Service) {
	controllerEmployee := NewEmployee(es)
	employeeGroup := r.Group("/api/v1/employees")
	{
		employeeGroup.GET("/:id", controllerEmployee.GetOne())
		employeeGroup.GET("/", controllerEmployee.GetAll())
		employeeGroup.POST("/", controllerEmployee.Create())
		employeeGroup.DELETE("/:id", controllerEmployee.Delete())
		employeeGroup.PATCH("/:id", controllerEmployee.Update())
	}
}

func (s *EmployeeController) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData reqEmployee

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("invalid request input"))
			return
		}

		if strings.ReplaceAll(requestData.CardNumberId, " ", "") == "" {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("empty card_number_id not allowed"))
			return
		}

		if len(requestData.CardNumberId) > 45 {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("card_number_id too long: max 45 characters"))
			return
		}

		if len(requestData.FirstName) > 45 {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("first_name too long: max 45 characters"))
			return
		}

		if len(requestData.LastName) > 45 {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("last_name too long: max 45 characters"))
			return
		}

		if requestData.WarehouseId == 0 {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("warehouse_id must be informed and greather than 0"))
			return
		}

		employee, resp := s.service.Create(
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
			web.NewResponse(employee),
		)
	}
}

func (s *EmployeeController) GetOne() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		parsedId, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, web.DecodeError("id must be a number"))
			return
		}

		employee, resp := s.service.GetOne(parsedId)

		if resp.Err != nil {
			c.JSON(
				resp.Code,
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
		var requestValidatorType reqEmployee
		requestData := make(map[string]interface{})

		id := c.Param("id")

		parsedId, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, web.DecodeError("id must be a number"))
			return
		}

		if err := c.ShouldBindBodyWith(&requestData, binding.JSON); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("invalid request data"))
			return
		}

		if len(requestData) == 0 {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("invalid request data - body needed"))
			return
		}

		if err := c.ShouldBindBodyWith(&requestValidatorType, binding.JSON); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("invalid type of data"))
			return
		}

		if requestData["card_number_id"] != nil {
			if strings.ReplaceAll(requestData["card_number_id"].(string), " ", "") == "" {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("empty card_number_id not allowed"))
				return
			}
			if len(requestData["card_number_id"].(string)) > 45 {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("card_number_id too long: max 45 characters"))
				return
			}
		}

		if requestData["first_name"] != nil {
			if len(requestData["first_name"].(string)) > 45 {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("first_name too long: max 45 characters"))
				return
			}
		}

		if requestData["last_name"] != nil {
			if len(requestData["last_name"].(string)) > 45 {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("last_name too long: max 45 characters"))
				return
			}
		}

		if requestData["warehouse_id"] != nil {
			if int(requestData["warehouse_id"].(float64)) == 0 {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.DecodeError("warehouse_id must be greather than 0"))
				return
			}
		}

		employee, resp := s.service.Update(parsedId, requestData)

		if resp.Err != nil {
			c.JSON(resp.Code, web.DecodeError(resp.Err.Error()))
			return
		}

		c.JSON(resp.Code, web.NewResponse(employee))
	}
}
