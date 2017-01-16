-- MySQL dump 10.13  Distrib 5.6.23, for osx10.8 (x86_64)
--
-- Host: localhost    Database: fabsmith_test
-- ------------------------------------------------------
-- Server version	5.6.23

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `activation_feedback`
--

DROP TABLE IF EXISTS `activation_feedback`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `activation_feedback` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `activation_id` int(11) NOT NULL,
  `satisfaction` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `activation_feedback`
--

LOCK TABLES `activation_feedback` WRITE;
/*!40000 ALTER TABLE `activation_feedback` DISABLE KEYS */;
/*!40000 ALTER TABLE `activation_feedback` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `activations`
--

DROP TABLE IF EXISTS `activations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `activations` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `invoice_id` int(11) DEFAULT NULL,
  `user_id` int(11) NOT NULL,
  `machine_id` int(11) NOT NULL,
  `active` tinyint(1) NOT NULL,
  `time_start` datetime NOT NULL,
  `time_end` datetime DEFAULT NULL,
  `time_total` int(11) NOT NULL,
  `used_kwh` float NOT NULL,
  `discount_percents` float NOT NULL,
  `discount_fixed` float NOT NULL,
  `vat_rate` float NOT NULL,
  `comment_ref` varchar(255) NOT NULL DEFAULT '',
  `invoiced` tinyint(1) NOT NULL,
  `changed` tinyint(1) NOT NULL,
  `current_machine_price` double unsigned DEFAULT NULL,
  `current_machine_price_currency` varchar(10) DEFAULT NULL,
  `current_machine_price_unit` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5481 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `activations`
--

LOCK TABLES `activations` WRITE;
/*!40000 ALTER TABLE `activations` DISABLE KEYS */;
/*!40000 ALTER TABLE `activations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `auth`
--

DROP TABLE IF EXISTS `auth`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `auth` (
  `user_id` int(11) NOT NULL,
  `nfc_key` varchar(100) NOT NULL DEFAULT '',
  `hash` varchar(300) NOT NULL DEFAULT '',
  `salt` varchar(100) NOT NULL DEFAULT '',
  `pw_reset_key` varchar(255) DEFAULT NULL,
  `pw_reset_time` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `auth`
--

LOCK TABLES `auth` WRITE;
/*!40000 ALTER TABLE `auth` DISABLE KEYS */;
/*!40000 ALTER TABLE `auth` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `coupon_usages`
--

DROP TABLE IF EXISTS `coupon_usages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `coupon_usages` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `coupon_id` int(11) unsigned DEFAULT NULL,
  `value` double DEFAULT NULL,
  `month` tinyint(4) DEFAULT NULL,
  `year` smallint(6) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `coupon_usages`
--

LOCK TABLES `coupon_usages` WRITE;
/*!40000 ALTER TABLE `coupon_usages` DISABLE KEYS */;
/*!40000 ALTER TABLE `coupon_usages` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `coupons`
--

DROP TABLE IF EXISTS `coupons`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `coupons` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `location_id` int(11) unsigned DEFAULT NULL,
  `code` varchar(100) DEFAULT NULL,
  `user_id` int(11) unsigned DEFAULT NULL,
  `value` double DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `coupons`
--

LOCK TABLES `coupons` WRITE;
/*!40000 ALTER TABLE `coupons` DISABLE KEYS */;
/*!40000 ALTER TABLE `coupons` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `invoice_user_memberships`
--

DROP TABLE IF EXISTS `invoice_user_memberships`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `invoice_user_memberships` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `location_id` int(11) unsigned DEFAULT NULL,
  `user_id` int(11) unsigned DEFAULT NULL,
  `membership_id` int(11) unsigned DEFAULT NULL,
  `user_membership_id` int(11) unsigned DEFAULT NULL,
  `start_date` varchar(255) DEFAULT NULL,
  `termination_date` varchar(255) DEFAULT NULL,
  `initial_duration_months` int(11) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  `invoice_id` int(11) DEFAULT NULL,
  `invoice_status` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4104 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `invoice_user_memberships`
--

LOCK TABLES `invoice_user_memberships` WRITE;
/*!40000 ALTER TABLE `invoice_user_memberships` DISABLE KEYS */;
/*!40000 ALTER TABLE `invoice_user_memberships` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `invoices`
--

DROP TABLE IF EXISTS `invoices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `invoices` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `location_id` int(11) unsigned DEFAULT NULL,
  `fastbill_id` int(11) unsigned DEFAULT NULL,
  `fastbill_no` varchar(100) DEFAULT NULL,
  `canceled_fastbill_id` int(11) unsigned DEFAULT NULL,
  `canceled_fastbill_no` varchar(100) DEFAULT NULL,
  `month` tinyint(3) unsigned DEFAULT NULL,
  `year` smallint(5) unsigned DEFAULT NULL,
  `customer_id` int(11) unsigned DEFAULT NULL,
  `customer_no` int(11) unsigned DEFAULT NULL,
  `user_id` int(11) unsigned DEFAULT NULL,
  `status` varchar(20) DEFAULT NULL,
  `canceled` tinyint(1) DEFAULT '0',
  `sent` tinyint(1) DEFAULT '0',
  `canceled_sent` tinyint(1) DEFAULT '0',
  `total` double DEFAULT NULL,
  `vat_percent` double DEFAULT NULL,
  `invoice_date` datetime DEFAULT NULL,
  `paid_date` datetime DEFAULT NULL,
  `due_date` datetime DEFAULT NULL,
  `current` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10030 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `invoices`
--

LOCK TABLES `invoices` WRITE;
/*!40000 ALTER TABLE `invoices` DISABLE KEYS */;
/*!40000 ALTER TABLE `invoices` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `locations`
--

DROP TABLE IF EXISTS `locations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `locations` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(100) DEFAULT NULL,
  `first_name` varchar(100) DEFAULT NULL,
  `last_name` varchar(100) DEFAULT NULL,
  `email` varchar(100) DEFAULT NULL,
  `city` varchar(100) DEFAULT NULL,
  `organization` varchar(100) DEFAULT NULL,
  `phone` varchar(100) DEFAULT NULL,
  `comments` varchar(100) DEFAULT NULL,
  `approved` tinyint(1) DEFAULT NULL,
  `xmpp_id` varchar(255) DEFAULT NULL,
  `local_ip` varchar(255) DEFAULT NULL,
  `feature_coworking` tinyint(1) DEFAULT NULL,
  `feature_setup_time` tinyint(1) DEFAULT NULL,
  `feature_spaces` tinyint(1) DEFAULT NULL,
  `feature_tutoring` tinyint(1) DEFAULT NULL,
  `feature_coupons` tinyint(1) DEFAULT NULL,
  `timezone` varchar(100) DEFAULT NULL,
  `logo` varchar(255) DEFAULT NULL,
  `university` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `locations`
--

LOCK TABLES `locations` WRITE;
/*!40000 ALTER TABLE `locations` DISABLE KEYS */;
/*!40000 ALTER TABLE `locations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `machine_maintenances`
--

DROP TABLE IF EXISTS `machine_maintenances`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `machine_maintenances` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `machine_id` int(11) unsigned DEFAULT NULL,
  `start` datetime DEFAULT NULL,
  `end` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `machine_maintenances`
--

LOCK TABLES `machine_maintenances` WRITE;
/*!40000 ALTER TABLE `machine_maintenances` DISABLE KEYS */;
/*!40000 ALTER TABLE `machine_maintenances` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `machine_types`
--

DROP TABLE IF EXISTS `machine_types`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `machine_types` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `location_id` int(11) DEFAULT NULL,
  `short_name` varchar(20) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `archived` tinyint(1) DEFAULT '0',
  `old_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=136 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `machine_types`
--

LOCK TABLES `machine_types` WRITE;
/*!40000 ALTER TABLE `machine_types` DISABLE KEYS */;
/*!40000 ALTER TABLE `machine_types` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `machines`
--

DROP TABLE IF EXISTS `machines`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `machines` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `location_id` int(11) unsigned DEFAULT NULL,
  `name` varchar(255) NOT NULL DEFAULT '',
  `shortname` varchar(100) DEFAULT NULL,
  `description` text NOT NULL,
  `image` varchar(255) DEFAULT NULL,
  `image_small` varchar(255) DEFAULT NULL,
  `available` tinyint(1) NOT NULL,
  `unavail_msg` text,
  `unavail_till` datetime DEFAULT NULL,
  `price` double unsigned NOT NULL,
  `price_unit` varchar(100) DEFAULT NULL,
  `comments` text,
  `visible` tinyint(1) DEFAULT NULL,
  `connected_machines` varchar(255) DEFAULT NULL,
  `switch_ref_count` int(11) DEFAULT NULL,
  `under_maintenance` tinyint(1) DEFAULT NULL,
  `reservation_price_start` double unsigned DEFAULT NULL,
  `reservation_price_hourly` double unsigned DEFAULT NULL,
  `grace_period` int(11) DEFAULT NULL,
  `type_id` int(11) unsigned DEFAULT NULL,
  `brand` varchar(255) DEFAULT NULL,
  `dimensions` varchar(255) DEFAULT NULL,
  `workspace_dimensions` varchar(255) DEFAULT NULL,
  `netswitch_url_on` varchar(255) DEFAULT NULL,
  `netswitch_url_off` varchar(255) DEFAULT NULL,
  `netswitch_host` varchar(255) DEFAULT NULL,
  `netswitch_sensor_port` int(5) DEFAULT NULL,
  `netswitch_type` varchar(255) DEFAULT NULL,
  `netswitch_xmpp` tinyint(1) DEFAULT NULL,
  `archived` tinyint(1) DEFAULT NULL,
  `netswitch_last_ping` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=149 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `machines`
--

LOCK TABLES `machines` WRITE;
/*!40000 ALTER TABLE `machines` DISABLE KEYS */;
/*!40000 ALTER TABLE `machines` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `membership`
--

DROP TABLE IF EXISTS `membership`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `membership` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `location_id` int(11) DEFAULT NULL,
  `title` varchar(100) NOT NULL DEFAULT '',
  `short_name` varchar(100) NOT NULL DEFAULT '',
  `duration_months` int(11) DEFAULT NULL,
  `monthly_price` double unsigned NOT NULL,
  `machine_price_deduction` int(11) NOT NULL,
  `affected_machines` text,
  `auto_extend` tinyint(1) DEFAULT NULL,
  `auto_extend_duration_months` int(11) DEFAULT NULL,
  `archived` tinyint(1) DEFAULT NULL,
  `affected_categories` text,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10608 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `membership`
--

LOCK TABLES `membership` WRITE;
/*!40000 ALTER TABLE `membership` DISABLE KEYS */;
/*!40000 ALTER TABLE `membership` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `migrations`
--

DROP TABLE IF EXISTS `migrations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `migrations` (
  `id_migration` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'surrogate key',
  `name` varchar(255) DEFAULT NULL COMMENT 'migration name, unique',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'date migrated or rolled back',
  `statements` longtext COMMENT 'SQL statements for this migration',
  `rollback_statements` longtext COMMENT 'SQL statment for rolling back migration',
  `status` enum('update','rollback') DEFAULT NULL COMMENT 'update indicates it is a normal migration while rollback means this migration is rolled back',
  PRIMARY KEY (`id_migration`)
) ENGINE=InnoDB AUTO_INCREMENT=88 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `migrations`
--

LOCK TABLES `migrations` WRITE;
/*!40000 ALTER TABLE `migrations` DISABLE KEYS */;
INSERT INTO `migrations` VALUES (1,'Billingingaddress_20150728_120222','2015-08-05 09:59:57','ALTER TABLE user MODIFY invoice_addr TEXT; ALTER TABLE user MODIFY ship_addr TEXT',NULL,'update'),(2,'Userphone_20150728_152059','2015-08-05 09:59:58','ALTER TABLE user ADD COLUMN phone VARCHAR(50)',NULL,'update'),(3,'Fastbilluseradd_20150805_170318','2015-08-07 13:12:43','ALTER TABLE user ADD COLUMN zip_code VARCHAR(100); ALTER TABLE user ADD COLUMN city VARCHAR(100); ALTER TABLE user ADD COLUMN country_code VARCHAR(2)',NULL,'update'),(4,'PricePerMonth_20150901_110810','2015-09-02 18:02:01','ALTER TABLE membership ADD monthly_price double unsigned NOT NULL AFTER price; UPDATE membership SET monthly_price = price WHERE duration = 30 AND unit = \'days\'; UPDATE membership SET monthly_price = price / 3 WHERE duration = 90 AND unit = \'days\'; UPDATE membership SET monthly_price = price / 12 WHERE duration = 365 AND unit = \'days\'; UPDATE membership SET monthly_price = price / duration * 30 WHERE duration <> 30 AND duration <> 90 AND duration <> 365 AND unit = \'days\'; UPDATE membership SET monthly_price = price WHERE unit <> \'days\'; ALTER TABLE membership DROP COLUMN price',NULL,'update'),(5,'Activationfeedback_20150908_145935','2015-09-18 08:12:25','\n		CREATE TABLE activation_feedback (\n			id bigint(20) unsigned NOT NULL AUTO_INCREMENT,\n			activation_id int(11) NOT NULL,\n			satisfaction varchar(100) DEFAULT NULL,\n			PRIMARY KEY (id)\n	)',NULL,'update'),(6,'Undermaintenance_20150909_113436','2015-09-18 08:12:26','ALTER TABLE machines ADD under_maintenance tinyint(1)',NULL,'update'),(7,'Usermembershipautoextend_20150910_110015','2015-09-25 07:44:06','ALTER TABLE user_membership ADD COLUMN auto_extend TINYINT(1)',NULL,'update'),(8,'Autoextendmemberships_20150910_172301','2015-09-25 07:44:07','ALTER TABLE membership ADD COLUMN auto_extend TINYINT(1); ALTER TABLE membership ADD COLUMN auto_extend_duration INT(11)',NULL,'update'),(9,'Usermembershipterminate_20150918_115716','2015-09-25 07:44:07','ALTER TABLE user_membership ADD COLUMN is_terminated TINYINT(1)',NULL,'update'),(10,'Autoextendmembershipmonths_20150921_115503','2015-09-25 07:44:08','ALTER TABLE membership CHANGE COLUMN auto_extend_duration auto_extend_duration_months INT(11)',NULL,'update'),(11,'Rmmembershipunitcol_20150921_123926','2015-09-25 07:44:09','ALTER TABLE membership CHANGE COLUMN duration duration_months INT(11); UPDATE membership SET duration_months=ROUND(duration_months / 30) WHERE unit=\'days\'; ALTER TABLE membership DROP COLUMN unit',NULL,'update'),(12,'Usermembershipenddate_20150924_114628','2015-09-25 07:44:09','\nUPDATE user_membership\nSET end_date = DATE_ADD(start_date, INTERVAL\n                          (SELECT duration_months\n                           FROM membership\n                           WHERE membership.id = user_membership.membership_id) MONTH)\n	',NULL,'update'),(13,'Autoextenddefaulttrue_20150924_120425','2015-09-25 07:44:09','UPDATE membership SET auto_extend = TRUE, auto_extend_duration_months = 1; UPDATE user_membership SET auto_extend = TRUE',NULL,'update'),(14,'Reservations_20150928_181415','2015-10-16 06:27:24','CREATE TABLE reservations (\n		id int(11) unsigned NOT NULL AUTO_INCREMENT,\n		machine_id int(11) NOT NULL,\n		user_id int(11) NOT NULL,\n		time_start datetime NOT NULL,\n		time_end datetime NOT NULL,\n		created datetime NOT NULL,\n		PRIMARY KEY (id)\n	)',NULL,'update'),(15,'Reservationrules_20151001_173125','2015-10-16 06:27:24','CREATE TABLE reservation_rules (\n		id int(11) unsigned NOT NULL AUTO_INCREMENT,\n		name varchar(100),\n		machine_id int(11),\n		available tinyint(1),\n		unavailable tinyint(1),\n		date_start char(10),\n		date_end char(10),\n		time_start char(5),\n		time_end char(5),\n		time_zone varchar(100),\n		monday tinyint(1),\n		tuesday tinyint(1),\n		wednesday tinyint(1),\n		thursday tinyint(1),\n		friday tinyint(1),\n		saturday tinyint(1),\n		sunday tinyint(1),\n		created datetime NOT NULL,\n		PRIMARY KEY (id)\n	)',NULL,'update'),(16,'Reservationprices_20151008_142149','2015-10-16 06:27:25','ALTER TABLE machines ADD COLUMN reservation_price_start double unsigned; ALTER TABLE machines ADD COLUMN reservation_price_hourly double unsigned',NULL,'update'),(17,'Activationcurrentprice_20151023_094534','2015-10-28 15:13:18','ALTER TABLE activations ADD COLUMN current_machine_price double unsigned; ALTER TABLE activations ADD COLUMN current_machine_price_currency varchar(10); ALTER TABLE activations ADD COLUMN current_machine_price_unit varchar(100); UPDATE activations a JOIN machines m ON a.machine_id = m.id SET a.current_machine_price=m.price, a.current_machine_price_currency=\'€\', a.current_machine_price_unit=m.price_unit',NULL,'update'),(18,'Reservationcurrentprice_20151026_120719','2015-10-28 15:13:19','ALTER TABLE reservations ADD COLUMN current_price double unsigned; ALTER TABLE reservations ADD COLUMN current_price_currency varchar(10); ALTER TABLE reservations ADD COLUMN current_price_unit varchar(100); UPDATE reservations r JOIN machines m ON r.machine_id = m.id SET r.current_price=m.reservation_price_hourly, r.current_price_currency=\'€\', r.current_price_unit=\'30 minutes\'',NULL,'update'),(19,'Globalconfig_20151117_095700','2015-12-06 18:33:16','CREATE TABLE settings (\n		id int(11) unsigned NOT NULL AUTO_INCREMENT,\n		name varchar(100) NOT NULL,\n		value_int int(11),\n		value_string text,\n		value_float double,\n		PRIMARY KEY (id)\n	)',NULL,'update'),(20,'Reservationdisabled_20151118_133537','2015-12-06 18:33:16','ALTER TABLE reservations ADD COLUMN disabled tinyint(1)',NULL,'update'),(21,'Products_20151118_141721','2015-12-06 18:33:16','CREATE TABLE products (\n		id int(11) unsigned NOT NULL AUTO_INCREMENT,\n		type varchar(100),\n		name varchar(100),\n		price double unsigned,\n		price_unit varchar(100),\n		PRIMARY KEY (id)\n	)',NULL,'update'),(22,'Purchases_20151119_115310','2015-12-06 18:33:17','CREATE TABLE purchases (\n		id int(11) unsigned NOT NULL AUTO_INCREMENT,\n		type varchar(100) NOT NULL,\n		product_id int(11) unsigned,\n		created datetime,\n		user_id int(11) unsigned NOT NULL,\n		time_start datetime DEFAULT NULL,\n		time_end datetime DEFAULT NULL,\n		quantity double NOT NULL,\n		price_per_unit double,\n		price_unit varchar(100),\n		vat double,\n		activation_running tinyint(1),\n		reservation_disabled tinyint(1),\n		machine_id int(11) unsigned,\n		PRIMARY KEY (id)\n	); \n		INSERT INTO purchases ( TYPE, product_id, created, user_id, time_start, time_end, quantity, price_per_unit, price_unit, vat, activation_running, reservation_disabled, machine_id )\n		SELECT \'activation\',\n		       NULL,\n		       time_start,\n		       user_id,\n		       time_start,\n		       time_end,\n		       time_total / 60,\n		       current_machine_price,\n		       current_machine_price_unit,\n		       vat_rate,\n		       active,\n		       NULL,\n		       machine_id\n		FROM activations\n		WHERE current_machine_price_unit = \"minute\"\n		UNION\n		SELECT \'activation\',\n		       NULL,\n		       time_start,\n		       user_id,\n		       time_start,\n		       time_end,\n		       time_total / 3600,\n		       current_machine_price,\n		       current_machine_price_unit,\n		       vat_rate,\n		       active,\n		       NULL,\n		       machine_id\n		FROM activations\n		WHERE current_machine_price_unit = \"hour\"\n	; \n		INSERT INTO purchases ( TYPE, product_id, created, user_id, time_start, time_end, quantity, price_per_unit, price_unit, vat, activation_running, reservation_disabled, machine_id )\n		SELECT \'reservation\',\n		       NULL,\n		       created,\n		       user_id,\n		       time_start,\n		       time_end,\n		       TIME_TO_SEC(TIMEDIFF(time_end, time_start)) / 1800,\n		       current_price / 2,\n		       current_price_unit,\n		       NULL,\n		       NULL,\n		       disabled,\n		       machine_id\n		FROM reservations\n	',NULL,'update'),(23,'PurchaseCancelledFlag_20151124_153407','2015-12-06 18:33:17','ALTER TABLE purchases ADD COLUMN cancelled tinyint(1)',NULL,'update'),(24,'Tutorproduct_20151126_120015','2015-12-06 18:33:18','ALTER TABLE products ADD COLUMN user_id int(11); ALTER TABLE products ADD COLUMN machine_skills varchar(255)',NULL,'update'),(25,'Productcomments_20151126_170242','2015-12-06 18:33:18','ALTER TABLE products ADD COLUMN comments TEXT',NULL,'update'),(26,'PurchaseTimeEndActual_20151202_102231','2015-12-06 18:33:19','ALTER TABLE purchases ADD COLUMN time_end_actual datetime AFTER time_end; ALTER TABLE purchases CHANGE COLUMN activation_running running tinyint(1)',NULL,'update'),(27,'Productarchived_20151202_172646','2015-12-06 18:33:20','ALTER TABLE products ADD COLUMN archived TINYINT(1) DEFAULT 0',NULL,'update'),(28,'PurchaseTimeEndPlanned_20151203_134204','2015-12-06 18:33:20','ALTER TABLE purchases CHANGE COLUMN time_end_actual time_end_planned datetime',NULL,'update'),(29,'Purchasearchived_20151203_142543','2015-12-06 18:33:21','ALTER TABLE purchases ADD COLUMN archived TINYINT(1) DEFAULT 0',NULL,'update'),(30,'Purchasecomments_20151203_160933','2015-12-06 18:33:21','ALTER TABLE purchases ADD COLUMN comments TEXT',NULL,'update'),(31,'Tutorpurchasetimer_20151208_143302','2015-12-17 19:35:58','ALTER TABLE purchases ADD COLUMN timer_time_start DATETIME',NULL,'update'),(32,'MachineXmpp_20160108_113817','2016-01-11 18:48:01','ALTER TABLE netswitch ADD COLUMN xmpp tinyint(1)',NULL,'update'),(33,'Locations_20160119_152047','2016-01-26 18:14:29','CREATE TABLE locations (\n		id int(11) unsigned NOT NULL AUTO_INCREMENT,\n		name varchar(100),\n		PRIMARY KEY (id)\n	); INSERT INTO locations VALUES (1, \"Fab Lab Berlin\"); ALTER TABLE machines ADD COLUMN location_id int(11) unsigned AFTER id; ALTER TABLE netswitch ADD COLUMN location_id int(11) unsigned AFTER id; UPDATE machines SET location_id = 1',NULL,'update'),(34,'Hostyourmachines_20160120_131417','2016-01-26 18:14:29','\n        CREATE TABLE hosts (\n            id bigint(20) unsigned NOT NULL AUTO_INCREMENT,\n            first_name varchar(100) NOT NULL,\n            last_name varchar(100) NOT NULL,\n            email varchar(100) NOT NULL,\n            location varchar(100) NOT NULL,\n            organization varchar(100) NOT NULL,\n            phone varchar(100) NOT NULL,\n            comments text NOT NULL,\n            PRIMARY KEY (id)\n    )',NULL,'update'),(35,'MergeHostsAndLocations_20160120_163019','2016-01-26 18:14:33','DROP TABLE hosts; ALTER TABLE locations CHANGE name title varchar(100); ALTER TABLE locations ADD COLUMN first_name varchar(100); ALTER TABLE locations ADD COLUMN last_name varchar(100); ALTER TABLE locations ADD COLUMN email varchar(100); ALTER TABLE locations ADD COLUMN city varchar(100); ALTER TABLE locations ADD COLUMN organization varchar(100); ALTER TABLE locations ADD COLUMN phone varchar(100); ALTER TABLE locations ADD COLUMN comments varchar(100); ALTER TABLE locations ADD COLUMN approved tinyint(1); UPDATE locations SET approved = 1 WHERE id = 1',NULL,'update'),(36,'NetswitchMfi_20160120_174425','2016-01-26 18:14:34','ALTER TABLE netswitch ADD COLUMN host varchar(255) AFTER url_off; ALTER TABLE netswitch ADD COLUMN sensor_port int(5) AFTER url_off; UPDATE netswitch SET sensor_port = 1',NULL,'update'),(37,'MachineType_20160126_151613','2016-01-26 18:14:35','ALTER TABLE machines ADD COLUMN type varchar(255); ALTER TABLE machines ADD COLUMN brand varchar(255); ALTER TABLE machines ADD COLUMN dimensions varchar(255); ALTER TABLE machines ADD COLUMN workspace_dimensions varchar(255)',NULL,'update'),(38,'MachineTypes_20160127_155843','2016-01-28 18:43:40','\n		CREATE TABLE machine_types (\n			id int(11) unsigned NOT NULL AUTO_INCREMENT,\n			short_name varchar(20),\n			name varchar(255),\n			PRIMARY KEY (id)\n	); INSERT INTO machine_types VALUES (1, \"3dprinter\", \"3D Printer\"); INSERT INTO machine_types VALUES (2, \"cnc\", \"CNC Mill\"); INSERT INTO machine_types VALUES (3, \"heatpress\", \"Heatpress\"); INSERT INTO machine_types VALUES (4, \"knitting\", \"Knitting Machine\"); INSERT INTO machine_types VALUES (5, \"lasercutter\", \"Lasercutter\"); INSERT INTO machine_types VALUES (6, \"vinylcutter\", \"Vinylcutter\"); ALTER TABLE machines ADD COLUMN type_id int(11) unsigned AFTER type; ALTER TABLE machines DROP COLUMN type',NULL,'update'),(39,'UserLocations_20160201_155547','2016-02-16 18:48:38','\n		CREATE TABLE user_locations (\n			id int(11) unsigned NOT NULL AUTO_INCREMENT,\n			location_id int(11) unsigned,\n			user_id int(11) unsigned,\n			user_role varchar(100),\n			archived tinyint(1) DEFAULT 0,\n			PRIMARY KEY (id)\n	)',NULL,'update'),(40,'UserLocationsAdminRoles_20160210_121712','2016-02-16 18:48:38','\n		INSERT INTO user_locations\n		            (location_id,\n		             user_id,\n		             user_role,\n		             archived)\n		SELECT 1,\n		       id,\n		       \"admin\",\n		       0\n		FROM   user\n		WHERE  user_role = \"admin\"\n	',NULL,'update'),(41,'ProductLocationId_20160210_145725','2016-02-16 18:48:38','\n	    ALTER TABLE products \n	      ADD COLUMN location_id INT(11) after id\n	; UPDATE products SET location_id = 1',NULL,'update'),(42,'PurchaseLocationId_20160210_182928','2016-02-16 18:48:39','\n	    ALTER TABLE purchases \n	      ADD COLUMN location_id INT(11) after id\n	; UPDATE purchases SET location_id = 1',NULL,'update'),(43,'MembershipLocationId_20160215_135752','2016-02-16 18:48:40','\n	    ALTER TABLE membership\n	      ADD COLUMN location_id int(11) AFTER id\n	; UPDATE membership SET location_id = 1',NULL,'update'),(44,'ReservationRulesLocationId_20160215_191406','2016-02-16 18:48:40','\n	    ALTER TABLE reservation_rules\n	      ADD COLUMN location_id int(11) AFTER id\n	; UPDATE reservation_rules SET location_id = 1',NULL,'update'),(45,'UserLocationsMemberRoles_20160216_131116','2016-02-16 18:48:40','UPDATE user SET user_role = \"member\" WHERE user_role <> \"admin\"; \n		INSERT INTO user_locations\n		            (location_id,\n		             user_id,\n		             user_role,\n		             archived)\n		SELECT 1,\n		       id,\n		       \"member\",\n		       0\n		FROM   user\n		WHERE  user_role <> \"admin\"\n	',NULL,'update'),(46,'LocationsXmppId_20160219_124354','2016-02-20 12:06:12','ALTER TABLE locations ADD COLUMN xmpp_id VARCHAR(255)',NULL,'update'),(47,'MergeMachinesAndNetswitches_20160219_133411','2016-02-20 12:06:13','\n		ALTER TABLE machines \n		  ADD COLUMN netswitch_url_on VARCHAR(255),\n		  ADD COLUMN netswitch_url_off VARCHAR(255),\n		  ADD COLUMN netswitch_host VARCHAR(255),\n		  ADD COLUMN netswitch_sensor_port INT(5),\n		  ADD COLUMN netswitch_xmpp TINYINT(1)\n	; \n		UPDATE machines\n		       JOIN netswitch\n		         ON machines.id = netswitch.machine_id\n		SET    netswitch_url_on = netswitch.url_on,\n		       netswitch_url_off = netswitch.url_off,\n		       netswitch_host = netswitch.host,\n		       netswitch_sensor_port = netswitch.sensor_port,\n		       netswitch_xmpp = netswitch.xmpp\n	; DROP TABLE netswitch',NULL,'update'),(48,'MachineGracePeriod_20160225_141536','2016-02-25 19:27:04','ALTER TABLE machines ADD COLUMN grace_period int(11) AFTER reservation_price_hourly',NULL,'update'),(49,'SettingsLocationId_20160225_163048','2016-02-25 19:27:04','\n	    ALTER TABLE settings \n	      ADD COLUMN location_id INT(11) after id\n	; UPDATE settings SET location_id = 1',NULL,'update'),(50,'MachineNetswitchType_20160226_165758','2016-03-02 08:14:25','ALTER TABLE machines\n	  ADD COLUMN netswitch_type VARCHAR(255) AFTER netswitch_sensor_port',NULL,'update'),(51,'LocationsLocalIp_20160303_195458','2016-03-08 18:33:33','ALTER TABLE locations ADD COLUMN local_ip VARCHAR(255); UPDATE locations SET local_ip = \'37.44.7.170\' WHERE id = 1',NULL,'update'),(52,'MonthlyEarnings_20160304_113146','2016-03-08 18:33:35','RENAME TABLE invoices TO monthly_earnings; ALTER TABLE monthly_earnings ADD COLUMN month_from TINYINT UNSIGNED AFTER id; ALTER TABLE monthly_earnings ADD COLUMN year_from SMALLINT UNSIGNED AFTER month_from; ALTER TABLE monthly_earnings ADD COLUMN month_to TINYINT UNSIGNED AFTER year_from; ALTER TABLE monthly_earnings ADD COLUMN year_to SMALLINT UNSIGNED AFTER month_to; UPDATE monthly_earnings SET month_from = MONTH(period_from); UPDATE monthly_earnings SET year_from = YEAR(period_from); UPDATE monthly_earnings SET month_to = MONTH(period_to); UPDATE monthly_earnings SET year_to = YEAR(period_to)',NULL,'update'),(53,'NetswitchTypeDefault_20160307_110139','2016-03-08 18:33:35','UPDATE machines SET netswitch_type = \'mfi\' WHERE netswitch_xmpp = 1',NULL,'update'),(54,'AuthPwReset_20160307_134938','2016-03-08 18:33:36','ALTER TABLE auth ADD COLUMN pw_reset_key VARCHAR(255); ALTER TABLE auth ADD COLUMN pw_reset_time DATETIME',NULL,'update'),(55,'LocationsFeatureToggles_20160314_204905','2016-03-16 13:08:48','ALTER TABLE locations ADD COLUMN feature_coworking TINYINT(1); ALTER TABLE locations ADD COLUMN feature_setup_time TINYINT(1); ALTER TABLE locations ADD COLUMN feature_spaces TINYINT(1); ALTER TABLE locations ADD COLUMN feature_tutoring TINYINT(1); UPDATE locations SET feature_coworking = 1 WHERE id = 1; UPDATE locations SET feature_setup_time = 1 WHERE id = 1; UPDATE locations SET feature_spaces = 1 WHERE id = 1; UPDATE locations SET feature_tutoring = 1 WHERE id = 1',NULL,'update'),(56,'MonthlyEarningLocationId_20160318_170010','2016-03-20 18:11:15','\n	    ALTER TABLE monthly_earnings \n	      ADD COLUMN location_id INT(11) after id\n	; UPDATE monthly_earnings SET location_id = 1',NULL,'update'),(57,'UserTestuser_20160321_142446','2016-03-21 13:27:47','ALTER TABLE user ADD COLUMN test_user TINYINT(1)',NULL,'update'),(58,'UserNoAutoInvoicing_20160413_141041','2016-04-13 20:32:28','ALTER TABLE user ADD COLUMN no_auto_invoicing TINYINT(1)',NULL,'update'),(59,'MachineArchived_20160419_155227','2016-04-20 17:48:50','ALTER TABLE machines ADD COLUMN archived TINYINT(1)',NULL,'update'),(60,'MembershipsArchived_20160419_180659','2016-04-20 17:48:50','ALTER TABLE membership ADD COLUMN archived TINYINT(1)',NULL,'update'),(61,'Coupons_20160421_144423','2016-04-26 09:43:42','\n		CREATE TABLE coupons (\n			id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,\n			location_id INT(11) UNSIGNED,\n			code VARCHAR(100),\n			user_id INT(11) UNSIGNED,\n			value DOUBLE,\n			PRIMARY KEY (id)\n	); \n		CREATE TABLE coupon_usages (\n			id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,\n			coupon_id INT(11) UNSIGNED,\n			value DOUBLE,\n			month TINYINT,\n			year SMALLINT,\n			PRIMARY KEY (id)\n	)',NULL,'update'),(62,'LocationsFeatureToggleCoupons_20160422_150028','2016-04-26 09:43:42','ALTER TABLE locations ADD COLUMN feature_coupons TINYINT(1); UPDATE locations SET feature_coupons = 1 WHERE id = 1',NULL,'update'),(63,'PurchasesInvoiceIdAndStatus_20160527_170456','2016-05-30 12:36:12','ALTER TABLE purchases ADD COLUMN invoice_id int(11) unsigned; ALTER TABLE purchases ADD COLUMN invoice_status varchar(20); ALTER TABLE user_membership ADD COLUMN invoice_id int(11) unsigned; ALTER TABLE user_membership ADD COLUMN invoice_status varchar(20)',NULL,'update'),(64,'InvoicesTable_20160531_154603','2016-08-26 07:20:53','\n		CREATE TABLE invoices (\n			id int(11) unsigned NOT NULL AUTO_INCREMENT,\n			location_id int(11) unsigned,\n			fastbill_id int(11) unsigned,\n			fastbill_no varchar(100),\n			canceled_fastbill_id int(11) unsigned,\n			canceled_fastbill_no varchar(100),\n			month tinyint unsigned,\n			year smallint unsigned,\n			customer_id int(11) unsigned,\n			customer_no int(11) unsigned,\n			user_id int(11) unsigned,\n			status varchar(20),\n			canceled tinyint(1) DEFAULT 0,\n			sent tinyint(1) DEFAULT 0,\n			canceled_sent tinyint(1) DEFAULT 0,\n			total real,\n			vat_percent real,\n			invoice_date DATETIME,\n			paid_date DATETIME,\n			due_date DATETIME,\n			current TINYINT(1),\n			PRIMARY KEY (id)\n	)',NULL,'update'),(65,'PurchasesInvoiceStatusNotNull_20160601_153246','2016-08-26 07:20:54','ALTER TABLE purchases DROP COLUMN invoice_status; ALTER TABLE user_membership DROP COLUMN invoice_status; ALTER TABLE purchases ADD COLUMN invoice_status varchar(20) NOT NULL DEFAULT \'\'; ALTER TABLE user_membership ADD COLUMN invoice_status varchar(20) NOT NULL DEFAULT \'\'',NULL,'update'),(66,'PurchasesInvoiceIdNotNull_20160602_140525','2016-08-26 07:23:05','ALTER TABLE purchases CHANGE invoice_id invoice_id int(11) unsigned NOT NULL; ALTER TABLE user_membership CHANGE invoice_id invoice_id int(11) unsigned NOT NULL',NULL,'update'),(67,'LocationsTimezone_20160826_131359','2016-08-26 11:42:05','ALTER TABLE locations ADD COLUMN timezone varchar(100)',NULL,'update'),(68,'UserVatFbTemplateId_20160826_183558','2016-08-29 11:55:42','ALTER TABLE user ADD COLUMN fastbill_template_id int(11) unsigned; ALTER TABLE user ADD COLUMN eu_delivery tinyint(1)',NULL,'update'),(69,'MachineLastPing_20160909_163523','2016-09-09 16:17:16','ALTER TABLE machines ADD COLUMN netswitch_last_ping DATETIME',NULL,'update'),(70,'LocationsLogo_20160926_111453','2016-09-26 12:07:53','ALTER TABLE locations ADD COLUMN logo varchar(255)',NULL,'update'),(71,'UserLocationsUnique_20160926_150303','2016-09-26 13:20:07','ALTER TABLE user_locations ADD UNIQUE unique_user_locations (user_id, location_id)',NULL,'update'),(72,'MachineImageSmall_20160927_152202','2016-09-27 13:59:21','ALTER TABLE machines ADD COLUMN image_small varchar(255) AFTER image; UPDATE machines SET image_small = image',NULL,'update'),(73,'PurchaseCustomName_20161017_135651','2016-10-19 12:55:41','ALTER TABLE purchases ADD COLUMN custom_name varchar(255) AFTER machine_id',NULL,'update'),(74,'PurchaseRemoveTimeEnd_20161018_174831','2016-10-21 11:27:10','ALTER TABLE purchases DROP COLUMN time_end',NULL,'update'),(75,'InvoiceUserMemberships_20161025_190752','2016-11-15 17:47:01','\nCREATE TABLE user_memberships (\n	id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,\n	location_id INT(11) UNSIGNED,\n	user_id INT(11) UNSIGNED,\n	membership_id INT(11) UNSIGNED,\n	start_date DATETIME,\n	termination_date DATETIME,\n	initial_duration_months INT(11),\n	auto_extend TINYINT(1),\n	created DATETIME,\n	updated DATETIME,\n	PRIMARY KEY (id)\n)\n	; \nCREATE TABLE invoice_user_memberships (\n	id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,\n	location_id INT(11) UNSIGNED,\n	user_id INT(11) UNSIGNED,\n	membership_id INT(11) UNSIGNED,\n	user_membership_id INT(11) UNSIGNED,\n	start_date DATETIME,\n	termination_date DATETIME,\n	initial_duration_months INT(11),\n	created DATETIME,\n	updated DATETIME,\n	invoice_id INT(11),\n	invoice_status VARCHAR(100),\n	PRIMARY KEY (id)\n)\n	',NULL,'update'),(76,'PopulateInvoiceUserMemberships_20161028_134609','2016-11-15 17:47:06','',NULL,'update'),(77,'UserMembershipsRemoveAutoExtend_20161104_175233','2016-11-15 17:47:06','ALTER TABLE user_memberships DROP COLUMN auto_extend',NULL,'update'),(78,'UserMembershipsStringDates_20161115_193212','2016-11-22 15:49:09','ALTER TABLE user_memberships CHANGE start_date start_date VARCHAR(255); UPDATE user_memberships SET start_date = SUBSTR(start_date, 1, 10); ALTER TABLE user_memberships CHANGE termination_date termination_date VARCHAR(255); UPDATE user_memberships SET termination_date = SUBSTR(termination_date, 1, 10); ALTER TABLE invoice_user_memberships CHANGE start_date start_date VARCHAR(255); UPDATE invoice_user_memberships SET start_date = SUBSTR(start_date, 1, 10); ALTER TABLE invoice_user_memberships CHANGE termination_date termination_date VARCHAR(255); UPDATE invoice_user_memberships SET termination_date = SUBSTR(termination_date, 1, 10)',NULL,'update'),(79,'MachineMaintenances_20161130_115623','2016-12-05 12:32:32','\nCREATE TABLE machine_maintenances (\n	id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,\n	machine_id INT(11) UNSIGNED,\n	start DATETIME,\n	end DATETIME,\n	PRIMARY KEY (id)\n)\n	',NULL,'update'),(80,'LocationsUniversity_20161207_154447','2016-12-07 15:52:44','ALTER TABLE locations ADD COLUMN university TINYINT(1)',NULL,'update'),(81,'UserUniversity_20161208_132309','2016-12-08 13:47:09','ALTER TABLE user ADD COLUMN student_id varchar(50); ALTER TABLE user ADD COLUMN security_briefing varchar(255)',NULL,'update'),(82,'UserRoleRefactoring_20161212_112724','2016-12-13 16:46:43','ALTER TABLE user ADD COLUMN super_admin TINYINT(1) AFTER user_role; ALTER TABLE user DROP COLUMN user_role',NULL,'update'),(83,'MachineTypeArchived_20170106_145406','2017-01-06 15:11:24','ALTER TABLE machine_types ADD COLUMN archived TINYINT(1) DEFAULT 0',NULL,'update'),(84,'PermissionCategoryId_20170106_165735','2017-01-16 13:45:01','SET sql_mode = \'\'; \n		CREATE TABLE permission_new (\n			id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,\n			location_id INT(11) UNSIGNED,\n			user_id INT(11) UNSIGNED,\n			category_id INT(11) UNSIGNED,\n			PRIMARY KEY (id)\n	); \nINSERT INTO permission_new\nSELECT permission.id,\n       location_id,\n       permission.user_id,\n       type_id\nFROM   permission\n       JOIN machines\n         ON machines.id = permission.machine_id\nGROUP  BY permission.user_id,\n          type_id \n	; RENAME TABLE permission TO permission_old; RENAME TABLE permission_new TO permission',NULL,'update'),(85,'MachineTypeLocationId_20170109_102614','2017-01-16 13:45:02','ALTER TABLE machine_types ADD COLUMN location_id int(11) AFTER id; ALTER TABLE machine_types ADD COLUMN old_id int(11); \nINSERT INTO machine_types (location_id, short_name, name, archived, old_id)\nSELECT l.id, t.short_name, t.name, t.archived, t.id\nFROM locations AS l,\n     machine_types AS t\n; DELETE FROM machine_types WHERE location_id IS NULL; UPDATE machines JOIN machine_types AS t ON machines.type_id = old_id SET type_id = t.id; UPDATE permission JOIN machine_types AS t ON permission.category_id = old_id SET category_id = t.id',NULL,'update'),(86,'MembershipsPopulateCategoryIds_20170112_140000','2017-01-16 13:45:02','ALTER TABLE membership ADD COLUMN affected_categories text',NULL,'update'),(87,'MembershipsPopulateCategoryIds_20170112_143836','2017-01-16 13:45:02','',NULL,'update');
/*!40000 ALTER TABLE `migrations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `monthly_earnings`
--

DROP TABLE IF EXISTS `monthly_earnings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `monthly_earnings` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `location_id` int(11) DEFAULT NULL,
  `month_from` tinyint(3) unsigned DEFAULT NULL,
  `year_from` smallint(5) unsigned DEFAULT NULL,
  `month_to` tinyint(3) unsigned DEFAULT NULL,
  `year_to` smallint(5) unsigned DEFAULT NULL,
  `activations` text NOT NULL,
  `file_path` varchar(255) NOT NULL DEFAULT '',
  `created` datetime DEFAULT NULL,
  `period_from` datetime DEFAULT NULL,
  `period_to` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=179 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `monthly_earnings`
--

LOCK TABLES `monthly_earnings` WRITE;
/*!40000 ALTER TABLE `monthly_earnings` DISABLE KEYS */;
/*!40000 ALTER TABLE `monthly_earnings` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `permission`
--

DROP TABLE IF EXISTS `permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `permission` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `location_id` int(11) unsigned DEFAULT NULL,
  `user_id` int(11) unsigned DEFAULT NULL,
  `category_id` int(11) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=31569 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `permission`
--

LOCK TABLES `permission` WRITE;
/*!40000 ALTER TABLE `permission` DISABLE KEYS */;
/*!40000 ALTER TABLE `permission` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `permission_old`
--

DROP TABLE IF EXISTS `permission_old`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `permission_old` (
  `id` bigint(11) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(11) NOT NULL,
  `machine_id` bigint(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=31547 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `permission_old`
--

LOCK TABLES `permission_old` WRITE;
/*!40000 ALTER TABLE `permission_old` DISABLE KEYS */;
/*!40000 ALTER TABLE `permission_old` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `products`
--

DROP TABLE IF EXISTS `products`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `products` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `location_id` int(11) DEFAULT NULL,
  `type` varchar(100) DEFAULT NULL,
  `name` varchar(100) DEFAULT NULL,
  `price` double unsigned DEFAULT NULL,
  `price_unit` varchar(100) DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  `machine_skills` varchar(255) DEFAULT NULL,
  `comments` text,
  `archived` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=50 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `products`
--

LOCK TABLES `products` WRITE;
/*!40000 ALTER TABLE `products` DISABLE KEYS */;
/*!40000 ALTER TABLE `products` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `purchases`
--

DROP TABLE IF EXISTS `purchases`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `purchases` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `location_id` int(11) DEFAULT NULL,
  `type` varchar(100) NOT NULL,
  `product_id` int(11) unsigned DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `user_id` int(11) unsigned NOT NULL,
  `time_start` datetime DEFAULT NULL,
  `time_end_planned` datetime DEFAULT NULL,
  `quantity` double NOT NULL,
  `price_per_unit` double DEFAULT NULL,
  `price_unit` varchar(100) DEFAULT NULL,
  `vat` double DEFAULT NULL,
  `running` tinyint(1) DEFAULT NULL,
  `reservation_disabled` tinyint(1) DEFAULT NULL,
  `machine_id` int(11) unsigned DEFAULT NULL,
  `custom_name` varchar(255) DEFAULT NULL,
  `cancelled` tinyint(1) DEFAULT NULL,
  `archived` tinyint(1) DEFAULT '0',
  `comments` text,
  `timer_time_start` datetime DEFAULT NULL,
  `invoice_id` int(11) unsigned NOT NULL,
  `invoice_status` varchar(20) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=40610 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `purchases`
--

LOCK TABLES `purchases` WRITE;
/*!40000 ALTER TABLE `purchases` DISABLE KEYS */;
/*!40000 ALTER TABLE `purchases` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `reservation_rules`
--

DROP TABLE IF EXISTS `reservation_rules`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `reservation_rules` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `location_id` int(11) DEFAULT NULL,
  `name` varchar(100) DEFAULT NULL,
  `machine_id` int(11) DEFAULT NULL,
  `available` tinyint(1) DEFAULT NULL,
  `unavailable` tinyint(1) DEFAULT NULL,
  `date_start` char(10) DEFAULT NULL,
  `date_end` char(10) DEFAULT NULL,
  `time_start` char(5) DEFAULT NULL,
  `time_end` char(5) DEFAULT NULL,
  `time_zone` varchar(100) DEFAULT NULL,
  `monday` tinyint(1) DEFAULT NULL,
  `tuesday` tinyint(1) DEFAULT NULL,
  `wednesday` tinyint(1) DEFAULT NULL,
  `thursday` tinyint(1) DEFAULT NULL,
  `friday` tinyint(1) DEFAULT NULL,
  `saturday` tinyint(1) DEFAULT NULL,
  `sunday` tinyint(1) DEFAULT NULL,
  `created` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=36 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `reservation_rules`
--

LOCK TABLES `reservation_rules` WRITE;
/*!40000 ALTER TABLE `reservation_rules` DISABLE KEYS */;
/*!40000 ALTER TABLE `reservation_rules` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `reservations`
--

DROP TABLE IF EXISTS `reservations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `reservations` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `machine_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `time_start` datetime NOT NULL,
  `time_end` datetime NOT NULL,
  `created` datetime NOT NULL,
  `current_price` double unsigned DEFAULT NULL,
  `current_price_currency` varchar(10) DEFAULT NULL,
  `current_price_unit` varchar(100) DEFAULT NULL,
  `disabled` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=189 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `reservations`
--

LOCK TABLES `reservations` WRITE;
/*!40000 ALTER TABLE `reservations` DISABLE KEYS */;
/*!40000 ALTER TABLE `reservations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `settings`
--

DROP TABLE IF EXISTS `settings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `settings` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `location_id` int(11) DEFAULT NULL,
  `name` varchar(100) NOT NULL,
  `value_int` int(11) DEFAULT NULL,
  `value_string` text,
  `value_float` double DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=118 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `settings`
--

LOCK TABLES `settings` WRITE;
/*!40000 ALTER TABLE `settings` DISABLE KEYS */;
/*!40000 ALTER TABLE `settings` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `first_name` varchar(100) NOT NULL DEFAULT '',
  `last_name` varchar(100) NOT NULL DEFAULT '',
  `username` varchar(100) NOT NULL DEFAULT '',
  `email` varchar(100) NOT NULL DEFAULT '',
  `invoice_addr` text,
  `ship_addr` text,
  `client_id` int(11) NOT NULL,
  `b2b` tinyint(1) NOT NULL,
  `company` varchar(100) NOT NULL DEFAULT '',
  `vat_user_id` varchar(100) NOT NULL DEFAULT '',
  `vat_rate` int(11) NOT NULL,
  `super_admin` tinyint(1) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `comments` text,
  `phone` varchar(50) DEFAULT NULL,
  `zip_code` varchar(100) DEFAULT NULL,
  `city` varchar(100) DEFAULT NULL,
  `country_code` varchar(2) DEFAULT NULL,
  `test_user` tinyint(1) DEFAULT NULL,
  `no_auto_invoicing` tinyint(1) DEFAULT NULL,
  `fastbill_template_id` int(11) unsigned DEFAULT NULL,
  `eu_delivery` tinyint(1) DEFAULT NULL,
  `student_id` varchar(50) DEFAULT NULL,
  `security_briefing` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=14006 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_locations`
--

DROP TABLE IF EXISTS `user_locations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_locations` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `location_id` int(11) unsigned DEFAULT NULL,
  `user_id` int(11) unsigned DEFAULT NULL,
  `user_role` varchar(100) DEFAULT NULL,
  `archived` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_user_locations` (`user_id`,`location_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4116 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_locations`
--

LOCK TABLES `user_locations` WRITE;
/*!40000 ALTER TABLE `user_locations` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_locations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_membership`
--

DROP TABLE IF EXISTS `user_membership`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_membership` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `membership_id` int(11) NOT NULL,
  `start_date` datetime NOT NULL,
  `end_date` datetime DEFAULT NULL,
  `auto_extend` tinyint(1) DEFAULT NULL,
  `is_terminated` tinyint(1) DEFAULT NULL,
  `invoice_id` int(11) unsigned NOT NULL,
  `invoice_status` varchar(20) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2975 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_membership`
--

LOCK TABLES `user_membership` WRITE;
/*!40000 ALTER TABLE `user_membership` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_membership` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_memberships`
--

DROP TABLE IF EXISTS `user_memberships`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_memberships` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `location_id` int(11) unsigned DEFAULT NULL,
  `user_id` int(11) unsigned DEFAULT NULL,
  `membership_id` int(11) unsigned DEFAULT NULL,
  `start_date` varchar(255) DEFAULT NULL,
  `termination_date` varchar(255) DEFAULT NULL,
  `initial_duration_months` int(11) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=971 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_memberships`
--

LOCK TABLES `user_memberships` WRITE;
/*!40000 ALTER TABLE `user_memberships` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_memberships` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_roles`
--

DROP TABLE IF EXISTS `user_roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_roles` (
  `user_id` int(11) NOT NULL,
  `admin` tinyint(1) NOT NULL,
  `staff` tinyint(1) NOT NULL,
  `member` tinyint(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_roles`
--

LOCK TABLES `user_roles` WRITE;
/*!40000 ALTER TABLE `user_roles` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_roles` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2017-01-16 15:16:09
