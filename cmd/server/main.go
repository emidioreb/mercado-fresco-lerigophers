package main

import (
	"database/sql"
	"log"
	"os"

	buyersController "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/buyers"
	carriersController "github.com/emidioreb/mercado-fresco-lerigophers/cmd/server/controllers/carriers"
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
	"github.com/joho/godotenv"

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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dataSource := os.Getenv("SERVER_URI")
	PORT := ":" + os.Getenv("PORT")

	conn, _ := sql.Open("mysql", dataSource)
	_, err = conn.Query("USE mercado_fresco")
	if err != nil {
		log.Fatal("Couldn't connect to database: mercado_fresco do not exists")
	}

	conn, err = sql.Open("mysql", dataSource)
	if err != nil {
		log.Fatal("failed to connect to mariadb")
	}

	repoProductBatches := product_batches.NewMariaDbRepository(conn)
	serviceProductBatches := product_batches.NewService(repoProductBatches)
	productBatchesController.NewProductBatchHandler(server, serviceProductBatches)

	repoLocalities := localities.NewMariaDbRepository(conn)
	serviceLocality := localities.NewService(repoLocalities)
	localitiesController.NewLocalityHandle(server, serviceLocality)

	repoBuyer := buyers.NewMariaDbRepository(conn)
	serviceBuyer := buyers.NewService(repoBuyer)
	buyersController.NewBuyerHandler(server, serviceBuyer)

	repoSellers := sellers.NewMariaDbRepository(conn)
	service := sellers.NewService(repoSellers, repoLocalities)
	sellersController.NewSellerHandler(server, service)

	repoWarehouse := warehouses.NewMariaDbRepository(conn)
	serviceWarehouse := warehouses.NewService(repoWarehouse)
	warehousesController.NewWarehouseHandler(server, serviceWarehouse)

	repoProductType := product_types.NewMariaDbRepository(conn)

	repoSection := sections.NewMariaDbRepository(conn)
	serviceSection := sections.NewService(repoSection, repoWarehouse, repoProductType)
	sectionsController.NewSectionHandler(server, serviceSection)

	repoProduct := products.NewMariaDbRepository(conn)
	serviceProduct := products.NewService(repoProduct, repoSellers)
	productsController.NewProductHandler(server, serviceProduct)

	repoProductRecords := product_records.NewMariaDbRepository(conn)
	serviceProductRecords := product_records.NewService(repoProductRecords, repoProduct)

	productRecordsController.NewProductRecordHandler(server, serviceProductRecords)

	repoEmployee := employees.NewMariaDbRepository(conn)
	serviceEmployee := employees.NewService(repoEmployee, repoWarehouse)
	employeesController.NewEmployeeHandler(server, serviceEmployee)

	repoCarriers := carriers.NewMariaDbRepository(conn)
	serviceCarriers := carriers.NewService(repoCarriers)
	carriersController.NewCarryHandler(server, serviceCarriers)

	repoInbound := inboundorders.NewMariaDbRepository(conn)
	serviceInbound := inboundorders.NewService(repoInbound, repoWarehouse, repoEmployee, repoProductBatches)
	inboundOrdersController.NewInboundHandler(server, serviceInbound)

	repoOrderStatus := order_status.NewMariaDbRepository(conn)

	repoPurchaseOrders := purchase_orders.NewMariaDbRepository(conn)
	servicePurchaseOrders := purchase_orders.NewService(repoPurchaseOrders, repoBuyer, repoProductRecords, repoOrderStatus)
	purchaseOrdersController.NewPurchaseOrderHandler(server, servicePurchaseOrders)

	server.Run(PORT)
}
