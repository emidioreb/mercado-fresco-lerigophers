package main

import (
	"github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/sellers"
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

	server.Run(":4000")
}
