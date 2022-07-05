package main

import (
	"database/sql"
	"log"

	buyersController "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/buyers"
	employeesController "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/employees"
	productsController "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/products"
	sectionsController "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/sections"
	sellersController "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/sellers"
	warehousesController "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/warehouses"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/buyers"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/employees"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/products"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/sections"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/sellers"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/warehouses"
	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	server := gin.Default()

	dataSource := "root:123456@tcp(localhost:4400)/mercado_fresco?parseTime=true"

	conn, _ := sql.Open("mysql", dataSource)
	_, err := conn.Query("USE mercado_fresco")
	if err != nil {
		log.Fatal("Couldn't connect to database: mercado_fresco do not exists")
	}

	conn, err = sql.Open("mysql", dataSource)
	if err != nil {
		log.Fatal("failed to connect to mariadb")
	}

	repoBuyer := buyers.NewRepository()
	serviceBuyer := buyers.NewService(repoBuyer)
	controllerBuyer := buyersController.NewBuyer(serviceBuyer)
	buyerGroup := server.Group("/api/v1/buyers")
	{
		buyerGroup.GET("/:id", controllerBuyer.GetOne())
		buyerGroup.GET("/", controllerBuyer.GetAll())
		buyerGroup.POST("/", controllerBuyer.Create())
		buyerGroup.DELETE("/:id", controllerBuyer.Delete())
		buyerGroup.PATCH("/:id", controllerBuyer.Update())
	}

	repoSellers := sellers.NewMariaDbRepository(conn)
	service := sellers.NewService(repoSellers)
	controller := sellersController.NewSeller(service)
	sellerGroup := server.Group("/api/v1/sellers")
	{
		sellerGroup.GET("/:id", controller.GetOne())
		sellerGroup.GET("/", controller.GetAll())
		sellerGroup.POST("/", controller.Create())
		sellerGroup.DELETE("/:id", controller.Delete())
		sellerGroup.PATCH("/:id", controller.Update())
	}

	repoWarehouse := warehouses.NewRepository()
	serviceWarehouse := warehouses.NewService(repoWarehouse)
	controllerWarehouse := warehousesController.NewWarehouse(serviceWarehouse)

	warehouseGroup := server.Group("/api/v1/warehouses")

	{
		warehouseGroup.GET("/", controllerWarehouse.GetAll())
		warehouseGroup.GET("/:id", controllerWarehouse.GetOne())
		warehouseGroup.POST("/", controllerWarehouse.Create())
		warehouseGroup.DELETE("/:id", controllerWarehouse.Delete())
		warehouseGroup.PATCH("/:id", controllerWarehouse.Update())
	}

	repoSection := sections.NewRepository()
	serviceSection := sections.NewService(repoSection)
	controllerSection := sectionsController.NewSection(serviceSection)

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
	controllerProduct := productsController.NewProduct(serviceProduct)

	productGroup := server.Group("/api/v1/products")
	{
		productGroup.GET("/:id", controllerProduct.GetOne())
		productGroup.GET("/", controllerProduct.GetAll())
		productGroup.POST("/", controllerProduct.Create())
		productGroup.DELETE("/:id", controllerProduct.Delete())
		productGroup.PATCH("/:id", controllerProduct.Update())
	}

	repoEmployee := employees.NewMariaDbRepository(conn)
	serviceEmployee := employees.NewService(repoEmployee, repoWarehouse)
	controllerEmployee := employeesController.NewEmployee(serviceEmployee)

	employeeGroup := server.Group("/api/v1/employees")
	{
		employeeGroup.GET("/:id", controllerEmployee.GetOne())
		employeeGroup.GET("/", controllerEmployee.GetAll())
		employeeGroup.POST("/", controllerEmployee.Create())
		employeeGroup.DELETE("/:id", controllerEmployee.Delete())
		employeeGroup.PATCH("/:id", controllerEmployee.Update())
	}

	server.Run(":4000")
}
