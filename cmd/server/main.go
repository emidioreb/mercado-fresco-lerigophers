package main

import (
	"database/sql"
	"log"

	buyersController "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/buyers"
	controllers "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/carriers"
	employeesController "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/employees"
	inboundOrdersController "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/inboundOrders"
	localitiesController "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/localities"
	productBatchesController "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/productBatches"
	productRecordsController "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/productRecords"
	productsController "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/products"
	purchaseOrdersController "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/purchaseOrders"
	sectionsController "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/sections"
	sellersController "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/sellers"
	warehousesController "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/warehouses"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/buyers"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/carriers"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/employees"
	order_status "github.com/emidioreb/mercado-fresco-lerigophers/internal/orderStatus"
	product_records "github.com/emidioreb/mercado-fresco-lerigophers/internal/productRecords"
	product_types "github.com/emidioreb/mercado-fresco-lerigophers/internal/productTypes"
	purchase_orders "github.com/emidioreb/mercado-fresco-lerigophers/internal/purchaseOrders"

	inboundorders "github.com/emidioreb/mercado-fresco-lerigophers/internal/inboundOrders"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/localities"
	product_batches "github.com/emidioreb/mercado-fresco-lerigophers/internal/productBatches"

	"github.com/emidioreb/mercado-fresco-lerigophers/internal/products"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/sections"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/sellers"
	"github.com/emidioreb/mercado-fresco-lerigophers/internal/warehouses"
	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	server := gin.Default()
	dataSource := "root:root@tcp(localhost:4000)/mercado_fresco?parseTime=true"

	conn, _ := sql.Open("mysql", dataSource)
	_, err := conn.Query("USE mercado_fresco")
	if err != nil {
		log.Fatal("Couldn't connect to database: mercado_fresco do not exists")
	}

	conn, err = sql.Open("mysql", dataSource)
	if err != nil {
		log.Fatal("failed to connect to mariadb")
	}

	repoProductBatches := product_batches.NewMariaDbRepository(conn)
	serviceProductBatches := product_batches.NewService(repoProductBatches)
	controllerProductBatches := productBatchesController.NewProductBatch(serviceProductBatches)
	ProductBatchesGroup := server.Group("/api/v1/productBatches")
	{
		ProductBatchesGroup.POST("/", controllerProductBatches.CreateProductBatch())
		ProductBatchesGroup.GET("/reportProducts", controllerProductBatches.GetReportSection())
	}

	repoLocalities := localities.NewMariaDbRepository(conn)
	serviceLocality := localities.NewService(repoLocalities)
	controllerLocality := localitiesController.NewLocality(serviceLocality)
	localityGroup := server.Group("/api/v1/localities")
	{
		localityGroup.POST("/", controllerLocality.CreateLocality())
		localityGroup.GET("/reportSellers", controllerLocality.GetReportSellers())
		localityGroup.GET("/reportCarries", controllerLocality.GetReportCarriers())
	}

	repoBuyer := buyers.NewMariaDbRepository(conn)
	serviceBuyer := buyers.NewService(repoBuyer)
	controllerBuyer := buyersController.NewBuyer(serviceBuyer)
	buyerGroup := server.Group("/api/v1/buyers")
	{
		buyerGroup.GET("/:id", controllerBuyer.GetOne())
		buyerGroup.GET("/", controllerBuyer.GetAll())
		buyerGroup.POST("/", controllerBuyer.Create())
		buyerGroup.DELETE("/:id", controllerBuyer.Delete())
		buyerGroup.PATCH("/:id", controllerBuyer.Update())
		buyerGroup.GET("/reportPurchaseOrders", controllerBuyer.GetReportPurchaseOrders())
	}

	repoSellers := sellers.NewMariaDbRepository(conn)
	service := sellers.NewService(repoSellers, repoLocalities)
	sellersController.NewSellerHandler(server, service)

	repoWarehouse := warehouses.NewMariaDbRepository(conn)
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

	repoProductType := product_types.NewMariaDbRepository(conn)

	repoSection := sections.NewMariaDbRepository(conn)
	serviceSection := sections.NewService(repoSection, repoWarehouse, repoProductType)
	controllerSection := sectionsController.NewSection(serviceSection)

	sectionGroup := server.Group("/api/v1/sections")
	{
		sectionGroup.GET("/:id", controllerSection.GetOne())
		sectionGroup.GET("/", controllerSection.GetAll())
		sectionGroup.POST("/", controllerSection.Create())
		sectionGroup.DELETE("/:id", controllerSection.Delete())
		sectionGroup.PATCH("/:id", controllerSection.Update())
	}

	repoProduct := products.NewMariaDbRepository(conn)
	serviceProduct := products.NewService(repoProduct, repoSellers)
	controllerProduct := productsController.NewProduct(serviceProduct)

	productGroup := server.Group("/api/v1/products")
	{
		productGroup.GET("/:id", controllerProduct.GetOne())
		productGroup.GET("/", controllerProduct.GetAll())
		productGroup.POST("/", controllerProduct.Create())
		productGroup.DELETE("/:id", controllerProduct.Delete())
		productGroup.PATCH("/:id", controllerProduct.Update())
		productGroup.GET("/reportRecords", controllerProduct.GetReportRecords())
	}

	repoProductRecords := product_records.NewMariaDbRepository(conn)
	serviceProductRecords := product_records.NewService(repoProductRecords, repoProduct)
	controllerProductRecords := productRecordsController.NewProductRecord(serviceProductRecords)
	ProductRecordsGroup := server.Group("/api/v1/productRecords")
	{
		ProductRecordsGroup.POST("/", controllerProductRecords.CreateProductRecord())
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

	repoCarriers := carriers.NewMariaDbRepository(conn)
	serviceCarriers := carriers.NewService(repoCarriers)
	controllerCarriers := controllers.NewCarry(serviceCarriers)

	carriersGroup := server.Group("/api/v1/carries")
	{
		carriersGroup.POST("/", controllerCarriers.Create())
	}

	repoInbound := inboundorders.NewMariaDbRepository(conn)
	serviceInbound := inboundorders.NewService(repoInbound, repoWarehouse, repoEmployee, repoProductBatches)
	controllerInbound := inboundOrdersController.NewInboud(serviceInbound)

	inboundGroup := server.Group("/api/v1")
	{
		inboundGroup.GET("employees/reportInboundOrders", controllerInbound.GetReportInboundOrders())
		inboundGroup.POST("/inboundOrders", controllerInbound.CreateInboundOrders())
	}

	repoOrderStatus := order_status.NewMariaDbRepository(conn)
	repoPurchaseOrders := purchase_orders.NewMariaDbRepository(conn)
	servicePurchaseOrders := purchase_orders.NewService(repoPurchaseOrders, repoBuyer, repoProductRecords, repoOrderStatus)
	controllerPurchaseOrders := purchaseOrdersController.NewPurchaseOrder(servicePurchaseOrders)
	PurchaseOrdersGroup := server.Group("/api/v1/purchaseOrders")
	{
		PurchaseOrdersGroup.POST("/", controllerPurchaseOrders.CreatePurchaseOrder())
	}
	server.Run(":4400")
}
