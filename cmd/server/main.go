package main

import (
	"github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/buyer"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/employees"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/products"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/sections"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/sellers"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/warehouses"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	repoBuyer := buyers.NewRepository()
	serviceBuyer := buyers.NewService(repoBuyer)
	controllerBuyer := controllers.NewBuyer(serviceBuyer)

	buyerGroup := server.Group("/api/v1/buyers")
	{
		buyerGroup.GET("/:id", controllerBuyer.GetOne())
		buyerGroup.GET("/", controllerBuyer.GetAll())
		buyerGroup.POST("/", controllerBuyer.Create())
		buyerGroup.DELETE("/:id", controllerBuyer.Delete())
		buyerGroup.PUT("/:id", controllerBuyer.Update())
		buyerGroup.PATCH("/:id", controllerBuyer.UpdateLastName())
	}

	repository := sellers.NewRepository()
	service := sellers.NewService(repository)
	controller := controllers.NewSeller(service)

	sellerGroup := server.Group("/api/v1/sellers")
	{
		sellerGroup.GET("/:id", controller.GetOne())
		sellerGroup.GET("/", controller.GetAll())
		sellerGroup.POST("/", controller.Create())
		sellerGroup.DELETE("/:id", controller.Delete())
		sellerGroup.PATCH("/:id", controller.Update())
		//sellerGroup.PUT("/:id", controller.Update())
	}

	repoWarehouse := warehouses.NewRepository()
	serviceWarehouse := warehouses.NewService(repoWarehouse)
	controllerWarehouse := controllers.NewWarehouse(serviceWarehouse)

	warehouseGroup := server.Group("/api/v1/warehouses")

	{
		warehouseGroup.GET("/", controllerWarehouse.GetAll())
		warehouseGroup.GET("/:id", controllerWarehouse.GetOne())
		warehouseGroup.POST("/", controllerWarehouse.Create())
		warehouseGroup.DELETE("/:id", controllerWarehouse.Delete())
		warehouseGroup.PUT("/:id", controllerWarehouse.Update())
		warehouseGroup.PATCH("/:id", controllerWarehouse.UpdateTelephone())
	}

	repoSection := sections.NewRepository()
	serviceSection := sections.NewService(repoSection)
	controllerSection := controllers.NewSection(serviceSection)

	sectionGroup := server.Group("/api/v1/sections")
	{
		sectionGroup.GET("/:id", controllerSection.GetOne())
		sectionGroup.GET("/", controllerSection.GetAll())
		sectionGroup.POST("/", controllerSection.Create())
		sectionGroup.DELETE("/:id", controllerSection.Delete())
		sectionGroup.PATCH("/:id", controllerSection.Update())
	}

	repoProduct := products.NewRepository()
	serviceProduct := products.NewService(repoProduct)
	controllerProduct := controllers.NewProduct(serviceProduct)

	productGroup := server.Group("/api/v1/products")
	{
		productGroup.GET("/:id", controllerProduct.GetOne())
		productGroup.GET("/", controllerProduct.GetAll())
		productGroup.POST("/", controllerProduct.Create())
		productGroup.DELETE("/:id", controllerProduct.Delete())
		productGroup.PUT("/:id", controllerProduct.Update())
		productGroup.PATCH("/:id", controllerProduct.UpdateExpirationRate())
	}

	repoEmployee := employees.NewRepository()
	serviceEmployee := employees.NewService(repoEmployee)
	controllerEmployee := controllers.NewEmployee(serviceEmployee)

	employeeGroup := server.Group("/api/v1/employees")
	{
		employeeGroup.GET("/:id", controllerEmployee.GetOne())
		employeeGroup.GET("/", controllerEmployee.GetAll())
		employeeGroup.POST("/", controllerEmployee.Create())
		employeeGroup.DELETE("/:id", controllerEmployee.Delete())
		employeeGroup.PUT("/:id", controllerEmployee.Update())
		employeeGroup.PATCH("/:id", controllerEmployee.UpdateFirstName())
	}

	server.Run(":4000")
}
