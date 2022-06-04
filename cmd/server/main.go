package main

import (
	"github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/sections"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/sellers"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/warehouses"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	repository := sellers.NewRepository()
	service := sellers.NewService(repository)
	controller := controllers.NewSeller(service)

	sellerGroup := server.Group("/api/v1/sellers")
	{
		sellerGroup.GET("/:id", controller.GetOne())
		sellerGroup.GET("/", controller.GetAll())
		sellerGroup.POST("/", controller.Create())
		sellerGroup.DELETE("/:id", controller.Delete())
		sellerGroup.PUT("/:id", controller.Update())
		sellerGroup.PATCH("/:id", controller.UpdateAddress())
	}

	repoWarehouse := warehouses.NewRepository()
	serviceWarehouse := warehouses.NewService(repoWarehouse)
	controllerWarehouse := controllers.NewWarehouse(serviceWarehouse)

	warehouseGroup := server.Group("/api/v1/warehouses")
	{
		warehouseGroup.GET("/:id", controllerWarehouse.GetOne())
		warehouseGroup.GET("/", controllerWarehouse.GetAll())
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
		sectionGroup.PUT("/:id", controllerSection.Update())
		sectionGroup.PATCH("/:id", controllerSection.UpdateCurrCapacity())
	}

	server.Run(":4000")
}
