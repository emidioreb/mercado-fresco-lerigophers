-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema mydb
-- -----------------------------------------------------
-- -----------------------------------------------------
-- Schema mercado_fresco
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema mercado_fresco
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `mercado_fresco` DEFAULT CHARACTER SET utf8mb3 ;
USE `mercado_fresco` ;

-- -----------------------------------------------------
-- Table `mercado_fresco`.`buyers`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mercado_fresco`.`buyers` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `card_number_id` VARCHAR(45) NOT NULL,
  `first_name` VARCHAR(45) NULL DEFAULT NULL,
  `last_name` VARCHAR(45) NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
  UNIQUE INDEX `card_number_id_UNIQUE` (`card_number_id` ASC) VISIBLE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `mercado_fresco`.`warehouses`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mercado_fresco`.`warehouses` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `warehouse_code` VARCHAR(255) NOT NULL,
  `address` VARCHAR(255) NULL DEFAULT NULL,
  `telephone` VARCHAR(20) NULL DEFAULT NULL,
  `minimum_temperature` INT NULL DEFAULT NULL,
  `minimum_capacity` INT UNSIGNED NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
  UNIQUE INDEX `cid_UNIQUE` (`warehouse_code` ASC) VISIBLE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `mercado_fresco`.`employees`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mercado_fresco`.`employees` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `card_number_id` VARCHAR(45) NOT NULL,
  `first_name` VARCHAR(45) NULL DEFAULT NULL,
  `last_name` VARCHAR(45) NULL DEFAULT NULL,
  `warehouse_id` INT UNSIGNED NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `card_number_id_UNIQUE` (`card_number_id` ASC) VISIBLE,
  INDEX `fk_employees_warehouses_idx` (`warehouse_id` ASC) VISIBLE,
  CONSTRAINT `fk_employees_warehouses`
    FOREIGN KEY (`warehouse_id`)
    REFERENCES `mercado_fresco`.`warehouses` (`id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `mercado_fresco`.`product_type`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mercado_fresco`.`product_type` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NULL DEFAULT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `mercado_fresco`.`localities`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mercado_fresco`.`localities` (
  `id` VARCHAR(255) NOT NULL,
  `locality_name` VARCHAR(255) NOT NULL,
  `province_name` VARCHAR(255) NOT NULL,
  `country_name` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `mercado_fresco`.`sellers`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mercado_fresco`.`sellers` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `cid` BIGINT(8) UNSIGNED NOT NULL,
  `company_name` VARCHAR(255) NULL DEFAULT NULL,
  `address` VARCHAR(255) NULL DEFAULT NULL,
  `telephone` VARCHAR(20) NULL DEFAULT NULL,
  `locality_id` VARCHAR(255) NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
  UNIQUE INDEX `cid_UNIQUE` (`cid` ASC) VISIBLE,
  INDEX `fk_sellers_localities_idx` (`locality_id` ASC) VISIBLE,
  CONSTRAINT `fk_sellers_localities`
    FOREIGN KEY (`locality_id`)
    REFERENCES `mercado_fresco`.`localities` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `mercado_fresco`.`products`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mercado_fresco`.`products` (
  `id` INT UNSIGNED AUTO_INCREMENT NOT NULL,
  `product_code` VARCHAR(255) NOT NULL,
  `description` VARCHAR(255) NULL DEFAULT NULL,
  `width` FLOAT NULL DEFAULT NULL,
  `height` FLOAT NULL DEFAULT NULL,
  `length` FLOAT NULL DEFAULT NULL,
  `net_weight` FLOAT NULL DEFAULT NULL,
  `expiration_rate` FLOAT NULL DEFAULT NULL,
  `recommended_freezing_temperature` FLOAT NULL DEFAULT NULL,
  `freezing_rate` FLOAT NULL DEFAULT NULL,
  `product_type_id` INT UNSIGNED NULL DEFAULT NULL,
  `seller_id` INT UNSIGNED NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `product_code_UNIQUE` (`product_code` ASC) VISIBLE,
  INDEX `fk_products_product_type_idx` (`product_type_id` ASC) VISIBLE,
  INDEX `fk_products_sellers_idx` (`seller_id` ASC) VISIBLE,
  CONSTRAINT `fk_products_product_type`
    FOREIGN KEY (`product_type_id`)
    REFERENCES `mercado_fresco`.`product_type` (`id`),
  CONSTRAINT `fk_products_sellers`
    FOREIGN KEY (`seller_id`)
    REFERENCES `mercado_fresco`.`sellers` (`id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `mercado_fresco`.`sections`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mercado_fresco`.`sections` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `section_number` INT NOT NULL,
  `current_temperature` INT NULL DEFAULT NULL,
  `minimum_temperature` INT NULL DEFAULT NULL,
  `current_capacity` INT UNSIGNED NULL DEFAULT NULL,
  `minimum_capacity` INT UNSIGNED NULL DEFAULT NULL,
  `maximum_capacity` INT UNSIGNED NULL DEFAULT NULL,
  `warehouse_id` INT UNSIGNED NULL DEFAULT NULL,
  `product_type_id` INT UNSIGNED NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
  UNIQUE INDEX `section_number_UNIQUE` (`section_number` ASC) VISIBLE,
  INDEX `fk_sections_warehouses_idx` (`warehouse_id` ASC) VISIBLE,
  INDEX `fk_sections_product_type_idx` (`product_type_id` ASC) VISIBLE,
  CONSTRAINT `fk_sections_product_type`
    FOREIGN KEY (`product_type_id`)
    REFERENCES `mercado_fresco`.`product_type` (`id`),
  CONSTRAINT `fk_sections_warehouses`
    FOREIGN KEY (`warehouse_id`)
    REFERENCES `mercado_fresco`.`warehouses` (`id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `mercado_fresco`.`carriers`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mercado_fresco`.`carriers` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `cid` VARCHAR(255) NOT NULL,
  `company_name` VARCHAR(255) NULL,
  `address` VARCHAR(255) NULL,
  `telephone` VARCHAR(255) NULL,
  `locality_id` VARCHAR(255) NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
  UNIQUE INDEX `cid_UNIQUE` (`cid` ASC) VISIBLE,
  INDEX `fk_carriers_localities_idx` (`locality_id` ASC) VISIBLE,
  CONSTRAINT `fk_carriers_localities`
    FOREIGN KEY (`locality_id`)
    REFERENCES `mercado_fresco`.`localities` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `mercado_fresco`.`order_status`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mercado_fresco`.`order_status` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `description` VARCHAR(255) NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `mercado_fresco`.`product_records`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mercado_fresco`.`product_records` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `last_update_date` DATE NULL,
  `purchase_price` INT NULL DEFAULT NULL,
  `sale_price` INT NULL,
  `product_id` INT UNSIGNED NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
  INDEX `fk_product_records_products_idx` (`product_id` ASC) VISIBLE,
  CONSTRAINT `fk_product_records_products`
    FOREIGN KEY (`product_id`)
    REFERENCES `mercado_fresco`.`products` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `mercado_fresco`.`purchase_orders`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mercado_fresco`.`purchase_orders` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `order_number` VARCHAR(255) NULL,
  `order_date` DATE NULL DEFAULT NULL,
  `tracking_code` VARCHAR(255) NULL DEFAULT NULL,
  `buyer_id` INT UNSIGNED NULL,
  `order_status_id` INT UNSIGNED NULL,
  `product_record_id` INT UNSIGNED NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
  INDEX `fk_purchase_orders_order_status_idx` (`order_status_id` ASC) VISIBLE,
  INDEX `fk_purchase_orders_buyer_idx` (`buyer_id` ASC) VISIBLE,
  INDEX `fk_purchase_orders_product_records_idx` (`product_record_id` ASC) VISIBLE,
  CONSTRAINT `fk_purchase_orders_buyer`
    FOREIGN KEY (`buyer_id`)
    REFERENCES `mercado_fresco`.`buyers` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_purchase_orders_order_status`
    FOREIGN KEY (`order_status_id`)
    REFERENCES `mercado_fresco`.`order_status` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_purchase_orders_product_records`
    FOREIGN KEY (`product_record_id`)
    REFERENCES `mercado_fresco`.`product_records` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `mercado_fresco`.`product_batches`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mercado_fresco`.`product_batches` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `batch_number` INT UNSIGNED NOT NULL,
  `current_quatity` INT NULL DEFAULT NULL,
  `current_temperature` INT NULL,
  `due_date` DATE NULL,
  `initial_quantity` INT NULL,
  `manufacturing_date` DATE NULL,
  `manufacturing_hour` INT NULL,
  `minimum_temperature` INT NULL,
  `product_id` INT UNSIGNED NULL,
  `section_id` INT UNSIGNED NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
  INDEX `fk_product_batches_sections_idx` (`section_id` ASC) VISIBLE,
  INDEX `fk_product_batches_products_idx` (`product_id` ASC) VISIBLE,
  UNIQUE INDEX `batch_number_UNIQUE` (`batch_number` ASC) VISIBLE,
  CONSTRAINT `fk_product_batches_sections`
    FOREIGN KEY (`section_id`)
    REFERENCES `mercado_fresco`.`sections` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_product_batches_products`
    FOREIGN KEY (`product_id`)
    REFERENCES `mercado_fresco`.`products` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb3;


-- -----------------------------------------------------
-- Table `mercado_fresco`.`inbound_orders`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mercado_fresco`.`inbound_orders` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `order_number` VARCHAR(255) NULL,
  `order_date` DATE NULL DEFAULT NULL,
  `employee_id` INT UNSIGNED NULL,
  `product_batch_id` INT UNSIGNED NULL,
  `warehouse_id` INT UNSIGNED NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
  INDEX `fk_inbound_orders_employee_idx` (`employee_id` ASC) VISIBLE,
  INDEX `fk_inbound_orders_products_batches_idx` (`product_batch_id` ASC) VISIBLE,
  INDEX `fk_inbound_orders_products_wareHouses_idx` (`warehouse_id` ASC) VISIBLE,
  CONSTRAINT `fk_inbound_orders_employee`
    FOREIGN KEY (`employee_id`)
    REFERENCES `mercado_fresco`.`employees` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_inbound_orders_products_batches`
    FOREIGN KEY (`product_batch_id`)
    REFERENCES `mercado_fresco`.`product_batches` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_inbound_orders_products_wareHouses`
    FOREIGN KEY (`warehouse_id`)
    REFERENCES `mercado_fresco`.`warehouses` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb3;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;

--  Order status
INSERT INTO `mercado_fresco`.`order_status` (description) VALUES ("ok");
INSERT INTO `mercado_fresco`.`order_status` (description) VALUES ("in progress");
INSERT INTO `mercado_fresco`.`order_status` (description) VALUES ("canceled");

-- Product types
INSERT INTO `mercado_fresco`.`product_type` (name) VALUES ("electronic");
INSERT INTO `mercado_fresco`.`product_type` (name) VALUES ("freezed");
INSERT INTO `mercado_fresco`.`product_type` (name) VALUES ("food");
INSERT INTO `mercado_fresco`.`product_type` (name) VALUES ("data storage");
INSERT INTO `mercado_fresco`.`product_type` (name) VALUES ("test");

-- Localities
INSERT INTO `mercado_fresco`.`localities` (`id`, `locality_name`, `province_name`, `country_name`) VALUES ("1", "Presidente Dutra", "MA", "BR");
INSERT INTO `mercado_fresco`.`localities` (`id`, `locality_name`, `province_name`, `country_name`) VALUES ("2", "Osasco", "SP", "BR");
INSERT INTO `mercado_fresco`.`localities` (`id`, `locality_name`, `province_name`, `country_name`) VALUES ("3", "Aparecida de Goi??nia", "GO", "BR");
INSERT INTO `mercado_fresco`.`localities` (`id`, `locality_name`, `province_name`, `country_name`) VALUES ("4", "Tuntum", "MA", "BR");
INSERT INTO `mercado_fresco`.`localities` (`id`, `locality_name`, `province_name`, `country_name`) VALUES ("5", "Barra do Corda", "MA", "BR");
INSERT INTO `mercado_fresco`.`localities` (`id`, `locality_name`, `province_name`, `country_name`) VALUES ("6", "Florian??polis", "SC", "BR");

-- Sellers
INSERT INTO `mercado_fresco`.`sellers` (`id`, `cid`, `company_name`, `address`, `telephone`, `locality_id`) VALUES ('1', '1', 'Mercado Livre', 'Av. Tancredo Neves', '123', '1');
INSERT INTO `mercado_fresco`.`sellers` (`id`, `cid`, `company_name`, `address`, `telephone`, `locality_id`) VALUES ('2', '2', 'Mercado Livre', 'Av. Das Na????es Unidas', '123', '2');
INSERT INTO `mercado_fresco`.`sellers` (`id`, `cid`, `company_name`, `address`, `telephone`, `locality_id`) VALUES ('3', '3', 'Mercado Livre', 'Av. Olavo Sampaio', '123', '3');
INSERT INTO `mercado_fresco`.`sellers` (`id`, `cid`, `company_name`, `address`, `telephone`, `locality_id`) VALUES ('4', '4', 'Mercado Livre', 'Av. Paulista', '123', '4');
INSERT INTO `mercado_fresco`.`sellers` (`id`, `cid`, `company_name`, `address`, `telephone`, `locality_id`) VALUES ('5', '5', 'Mercado Livre', 'Av. Das Flores', '456', '5');
INSERT INTO `mercado_fresco`.`sellers` (`id`, `cid`, `company_name`, `address`, `telephone`, `locality_id`) VALUES ('6', '6', 'Mercado Livre', 'Av. Campo Dantas', '456', '6');

-- Products
INSERT INTO `mercado_fresco`.`products` (`id`, `product_code`, `description`, `width`, `height`, `length`, `net_weight`, `expiration_rate`, `recommended_freezing_temperature`, `freezing_rate`, `product_type_id`, `seller_id`) VALUES ('1', '1', 'Tomato', '250', '250', '250', '10', '50', '34', '50', '1', '1');
INSERT INTO `mercado_fresco`.`products` (`id`, `product_code`, `description`, `width`, `height`, `length`, `net_weight`, `expiration_rate`, `recommended_freezing_temperature`, `freezing_rate`, `product_type_id`, `seller_id`) VALUES ('2', '2', 'Purple Onion', '250', '250', '250', '10', '60', '34', '50', '2', '2');
INSERT INTO `mercado_fresco`.`products` (`id`, `product_code`, `description`, `width`, `height`, `length`, `net_weight`, `expiration_rate`, `recommended_freezing_temperature`, `freezing_rate`, `product_type_id`, `seller_id`) VALUES ('3', '3', 'RAM 16GB DDR5', '260', '250', '250', '10', '70', '34', '50', '3', '3');
INSERT INTO `mercado_fresco`.`products` (`id`, `product_code`, `description`, `width`, `height`, `length`, `net_weight`, `expiration_rate`, `recommended_freezing_temperature`, `freezing_rate`, `product_type_id`, `seller_id`) VALUES ('4', '4', 'RTX 3080 16GB', '250', '250', '250', '10', '90', '34', '50', '4', '4');
INSERT INTO `mercado_fresco`.`products` (`id`, `product_code`, `description`, `width`, `height`, `length`, `net_weight`, `expiration_rate`, `recommended_freezing_temperature`, `freezing_rate`, `product_type_id`, `seller_id`) VALUES ('5', '5', 'Intel Core i9', '250', '250', '250', '10', '80', '34', '50', '5', '5');

-- Warehouses
INSERT INTO `mercado_fresco`.`warehouses`(`warehouse_code`,`address`,`telephone`,`minimum_temperature`,`minimum_capacity`)VALUES("Cod#1","Rua Herc??lio Luz, Florian??polis","48999001122",10,10);
INSERT INTO `mercado_fresco`.`warehouses`(`warehouse_code`,`address`,`telephone`,`minimum_temperature`,`minimum_capacity`)VALUES("Cod#2","Avenida Paulista, SP","11999001133",5,20);

-- Buyers
INSERT INTO `mercado_fresco`.`buyers`(`card_number_id`,`first_name`,`last_name`)VALUES("Card#1","Vitor","Souza");
INSERT INTO `mercado_fresco`.`buyers`(`card_number_id`,`first_name`,`last_name`)VALUES("Card#2","Lucas","Bulh??es");

-- Product Records
INSERT INTO `mercado_fresco`.`product_records`(`last_update_date`,`purchase_price`,`sale_price`,`product_id`)VALUES("2022-01-01",5,25,2);
INSERT INTO `mercado_fresco`.`product_records`(`last_update_date`,`purchase_price`,`sale_price`,`product_id`)VALUES("2022-01-05",6,23,2);

-- Sections
INSERT INTO `mercado_fresco`.`sections`(`section_number`,`current_temperature`,`minimum_temperature`,`current_capacity`,`minimum_capacity`,`maximum_capacity`,`warehouse_id`,`product_type_id`)VALUES(1,25,5,200,5,800,1,1);
INSERT INTO `mercado_fresco`.`sections`(`section_number`,`current_temperature`,`minimum_temperature`,`current_capacity`,`minimum_capacity`,`maximum_capacity`,`warehouse_id`,`product_type_id`)VALUES(13,25,5,200,5,800,2,2);
