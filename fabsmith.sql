-- MySQL dump 10.13  Distrib 5.6.14, for osx10.7 (x86_64)
--
-- Host: localhost    Database: fabsmith
-- ------------------------------------------------------
-- Server version	5.6.14

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
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=83 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `activations`
--

LOCK TABLES `activations` WRITE;
/*!40000 ALTER TABLE `activations` DISABLE KEYS */;
INSERT INTO `activations` VALUES (1,NULL,2,1,0,'2015-06-05 10:54:17','2015-06-05 10:56:07',111,0,0,0,0,'',0,0),(2,NULL,2,1,0,'2015-06-05 10:56:09','2015-06-05 11:05:33',565,0,0,0,0,'',0,0),(3,NULL,2,1,0,'2015-06-05 11:05:35','2015-06-05 11:07:48',134,0,0,0,0,'',0,0),(8,NULL,2,1,0,'2015-06-05 11:17:26','2015-06-05 11:24:54',448,0,0,0,0,'',0,0),(9,NULL,2,1,0,'2015-06-05 11:29:12','2015-06-05 11:30:29',78,0,0,0,0,'',0,0),(10,NULL,2,1,0,'2015-06-05 11:31:06','2015-06-05 11:38:50',464,0,0,0,0,'',0,0),(17,NULL,2,1,0,'2015-06-05 11:57:28','2015-06-05 11:58:18',50,0,0,0,0,'',0,0),(18,NULL,2,1,0,'2015-06-05 11:58:32','2015-06-05 11:58:36',5,0,0,0,0,'',0,0),(20,NULL,2,1,0,'2015-06-05 12:04:05','2015-06-05 12:04:12',7,0,0,0,0,'',0,0),(21,NULL,2,1,0,'2015-06-05 12:12:15','2015-06-05 12:12:19',5,0,0,0,0,'',0,0),(22,NULL,2,1,0,'2015-06-05 12:12:22','2015-06-05 12:12:27',5,0,0,0,0,'',0,0),(23,NULL,2,1,0,'2015-06-05 12:12:31','2015-06-05 12:13:11',41,0,0,0,0,'',0,0),(24,NULL,1,1,0,'2015-06-05 12:14:30','2015-06-05 12:15:59',89,0,0,0,0,'',0,0),(25,NULL,2,1,0,'2015-06-05 12:19:21','2015-06-05 12:19:24',3,0,0,0,0,'',0,0),(26,NULL,1,1,0,'2015-06-05 12:19:33','2015-06-05 12:19:35',3,0,0,0,0,'',0,0),(27,NULL,1,1,0,'2015-06-05 12:36:47','2015-06-05 12:36:51',5,0,0,0,0,'',0,0),(28,NULL,1,1,0,'2015-06-05 12:37:39','2015-06-05 12:38:11',32,0,0,0,0,'',0,0),(29,NULL,2,1,0,'2015-06-05 12:38:43','2015-06-05 12:38:46',3,0,0,0,0,'',0,0),(30,NULL,1,2,0,'2015-06-05 12:53:19','2015-06-15 17:18:22',879904,0,0,0,0,'',0,0),(31,NULL,2,3,0,'2015-06-15 11:50:17','2015-06-15 11:50:32',15,0,0,0,0,'',0,0),(32,NULL,2,3,0,'2015-06-15 11:50:34','2015-06-15 11:50:40',6,0,0,0,0,'',0,0),(33,NULL,2,3,0,'2015-06-15 11:51:17','2015-06-15 11:52:33',77,0,0,0,0,'',0,0),(34,NULL,2,3,0,'2015-06-15 11:52:36','2015-06-15 11:52:41',6,0,0,0,0,'',0,0),(35,NULL,2,3,0,'2015-06-15 11:52:46','2015-06-15 11:52:55',10,0,0,0,0,'',0,0),(36,NULL,2,3,0,'2015-06-15 12:31:41','2015-06-15 12:31:53',13,0,0,0,0,'',0,0),(37,NULL,2,3,0,'2015-06-15 12:31:57','2015-06-15 12:32:09',12,0,0,0,0,'',0,0),(38,NULL,2,3,0,'2015-06-15 12:33:37','2015-06-15 12:33:51',14,0,0,0,0,'',0,0),(39,NULL,2,3,0,'2015-06-15 12:33:55','2015-06-15 12:34:10',16,0,0,0,0,'',0,0),(40,NULL,2,3,0,'2015-06-15 12:34:43','2015-06-15 12:34:55',12,0,0,0,0,'',0,0),(41,NULL,2,3,0,'2015-06-15 12:34:56','2015-06-15 12:35:02',6,0,0,0,0,'',0,0),(42,NULL,2,3,0,'2015-06-15 12:35:05','2015-06-15 12:35:13',8,0,0,0,0,'',0,0),(43,NULL,2,3,0,'2015-06-15 12:35:28','2015-06-15 12:35:32',5,0,0,0,0,'',0,0),(44,NULL,2,3,0,'2015-06-15 12:36:32','2015-06-15 12:38:57',145,0,0,0,0,'',0,0),(45,NULL,2,3,0,'2015-06-15 12:38:59','2015-06-15 12:39:22',24,0,0,0,0,'',0,0),(46,NULL,2,3,0,'2015-06-15 12:39:39','2015-06-15 12:42:05',147,0,0,0,0,'',0,0),(47,NULL,2,3,0,'2015-06-15 12:42:07','2015-06-15 12:42:16',9,0,0,0,0,'',0,0),(48,NULL,2,3,0,'2015-06-15 12:42:19','2015-06-15 12:42:28',9,0,0,0,0,'',0,0),(49,NULL,2,3,0,'2015-06-15 17:14:22','2015-06-15 17:14:29',7,0,0,0,0,'',0,0),(50,NULL,2,3,0,'2015-06-15 17:16:48','2015-06-15 17:17:05',17,0,0,0,0,'',0,0),(51,NULL,2,3,0,'2015-06-15 17:17:06','2015-06-15 17:18:14',69,0,0,0,0,'',0,0),(52,NULL,2,3,0,'2015-06-15 17:18:31','2015-06-15 17:18:43',12,0,0,0,0,'',0,0),(53,NULL,2,2,0,'2015-06-15 17:18:43','2015-06-17 12:02:51',153849,0,0,0,0,'',0,0),(54,NULL,2,3,0,'2015-06-15 17:40:19','2015-06-15 17:40:23',5,0,0,0,0,'',0,0),(55,NULL,2,3,0,'2015-06-17 11:52:28','2015-06-17 11:52:36',9,0,0,0,0,'',0,0),(56,NULL,1,2,0,'2015-06-17 15:33:48','2015-06-17 15:34:28',40,0,0,0,0,'',0,0),(57,NULL,1,2,0,'2015-06-17 15:34:37','2015-06-17 15:42:04',447,0,0,0,0,'',0,0),(58,NULL,1,2,0,'2015-06-17 15:42:43','2015-06-17 15:43:02',19,0,0,0,0,'',0,0),(59,NULL,1,2,0,'2015-06-17 15:43:10','2015-06-17 15:43:32',23,0,0,0,0,'',0,0),(60,NULL,2,2,0,'2015-06-18 13:35:23','2015-06-18 14:08:48',2006,0,0,0,0,'',0,0),(61,NULL,1,1,0,'2015-06-18 14:08:40','2015-06-18 14:13:14',275,0,0,0,0,'',0,0),(62,NULL,1,2,0,'2015-06-18 14:08:56','2015-06-18 14:09:06',10,0,0,0,0,'',0,0),(63,NULL,1,2,0,'2015-06-18 14:09:08','2015-06-18 14:09:22',15,0,0,0,0,'',0,0),(64,NULL,1,2,0,'2015-06-18 14:09:29','2015-06-18 14:09:38',9,0,0,0,0,'',0,0),(65,NULL,1,2,0,'2015-06-18 14:10:34','2015-06-18 14:10:47',14,0,0,0,0,'',0,0),(66,NULL,1,2,0,'2015-06-18 14:12:17','2015-06-18 14:12:25',8,0,0,0,0,'',0,0),(67,NULL,1,2,0,'2015-06-18 14:12:31','2015-06-18 14:12:46',15,0,0,0,0,'',0,0),(68,NULL,2,2,0,'2015-06-18 14:12:52','2015-06-18 14:12:59',7,0,0,0,0,'',0,0),(69,NULL,1,2,0,'2015-06-18 14:13:05','2015-06-18 14:13:47',42,0,0,0,0,'',0,0),(70,NULL,2,1,0,'2015-06-18 14:13:22','2015-06-18 14:13:28',6,0,0,0,0,'',0,0),(71,NULL,1,1,0,'2015-06-18 14:13:31','2015-06-18 14:13:43',12,0,0,0,0,'',0,0),(72,NULL,2,1,0,'2015-06-18 16:22:38','2015-06-18 16:23:32',55,0,0,0,0,'',0,0),(73,NULL,2,1,0,'2015-06-18 16:23:36','2015-06-18 16:24:32',56,0,0,0,0,'',0,0),(74,NULL,2,2,0,'2015-06-18 16:24:37','2015-06-18 16:24:43',6,0,0,0,0,'',0,0),(75,NULL,2,2,0,'2015-06-18 16:24:48','2015-06-18 16:24:52',5,0,0,0,0,'',0,0),(76,NULL,1,2,0,'2015-06-18 16:24:56','2015-06-18 16:25:36',41,0,0,0,0,'',0,0),(77,NULL,2,1,0,'2015-06-18 16:27:16','2015-06-18 16:27:23',8,0,0,0,0,'',0,0),(78,NULL,1,1,0,'2015-06-18 16:27:26','2015-06-18 16:27:38',12,0,0,0,0,'',0,0),(79,NULL,2,1,0,'2015-06-18 16:27:40','2015-06-18 16:27:52',13,0,0,0,0,'',0,0),(80,NULL,2,2,0,'2015-06-18 19:10:21','2015-06-22 16:43:30',336790,0,0,0,0,'',0,0),(81,NULL,2,1,0,'2015-06-22 16:43:37','2015-06-22 16:43:39',3,0,0,0,0,'',0,0),(82,NULL,2,3,0,'2015-06-22 16:43:40','2015-06-22 16:43:42',3,0,0,0,0,'',0,0);
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
  `salt` varchar(100) NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `auth`
--

LOCK TABLES `auth` WRITE;
/*!40000 ALTER TABLE `auth` DISABLE KEYS */;
INSERT INTO `auth` VALUES (1,'1234','6cb03f6918dc261f58f5a22486e7d47dd8f28f69b55f0f60615f6aec403910b383adfa1d49a961a688f88bf2248b8256083ee6727e3751b5c58c0e95ae186ab0','9a18f2959dd3dee005873bce01da036a0751154eac08ecd7da9a802cd2abda18'),(2,'04126DBA233580','f2ae94748430030d13b43359d021411e97c5b49aa7fcad70134abc27ccdd4c24de096143810bffaf824d9b5dc028bb3de8f144ee48ea06c4e0895d8a6df6cdac','f755192ae0fcb52ef329b0007ccd8f1e68d040d7801a32e96db24b047896c27d');
/*!40000 ALTER TABLE `auth` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `hexaswitch`
--

DROP TABLE IF EXISTS `hexaswitch`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `hexaswitch` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `machine_id` int(11) NOT NULL,
  `switch_ip` varchar(100) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `hexaswitch`
--

LOCK TABLES `hexaswitch` WRITE;
/*!40000 ALTER TABLE `hexaswitch` DISABLE KEYS */;
/*!40000 ALTER TABLE `hexaswitch` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `invoices`
--

DROP TABLE IF EXISTS `invoices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `invoices` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `activations` text NOT NULL,
  `file_path` varchar(255) NOT NULL DEFAULT '',
  `created` datetime DEFAULT NULL,
  `period_from` datetime DEFAULT NULL,
  `period_to` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `invoices`
--

LOCK TABLES `invoices` WRITE;
/*!40000 ALTER TABLE `invoices` DISABLE KEYS */;
/*!40000 ALTER TABLE `invoices` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `machines`
--

DROP TABLE IF EXISTS `machines`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `machines` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT '',
  `shortname` varchar(100) DEFAULT NULL,
  `description` text NOT NULL,
  `image` varchar(255) DEFAULT NULL,
  `available` tinyint(1) NOT NULL,
  `unavail_msg` text,
  `unavail_till` datetime DEFAULT NULL,
  `price` double unsigned NOT NULL,
  `price_unit` varchar(100) DEFAULT NULL,
  `comments` text,
  `visible` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `machines`
--

LOCK TABLES `machines` WRITE;
/*!40000 ALTER TABLE `machines` DISABLE KEYS */;
INSERT INTO `machines` VALUES (1,'Laydrop 3D Printer','MB3DP','NYC 3D printer 4 real and 4 life.','',1,'',NULL,16,'hour','',1),(2,'MakerBot 3D Printer','MB3DP','NYC 3D printer 4 real and 4 life.','machine-2.svg',1,'',NULL,16,'hour','',1),(3,'Zing Laser Cutter','ZLC','Cuts wood, plastic, paper. Fast.','',1,'',NULL,1,'minute','asd',1);
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
  `title` varchar(100) NOT NULL DEFAULT '',
  `short_name` varchar(100) NOT NULL DEFAULT '',
  `duration` int(11) NOT NULL,
  `unit` varchar(100) NOT NULL,
  `price` double NOT NULL,
  `machine_price_deduction` int(11) NOT NULL,
  `affected_machines` text,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `membership`
--

LOCK TABLES `membership` WRITE;
/*!40000 ALTER TABLE `membership` DISABLE KEYS */;
INSERT INTO `membership` VALUES (1,'6 Months Basic','6MB',160,'days',400,50,'[1,2,3]'),(2,'1 Month Basic','SMB',30,'days',100,50,'[1,2,3]');
/*!40000 ALTER TABLE `membership` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `netswitch`
--

DROP TABLE IF EXISTS `netswitch`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `netswitch` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `machine_id` int(11) unsigned NOT NULL,
  `url_on` varchar(255) NOT NULL DEFAULT '',
  `url_off` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `netswitch`
--

LOCK TABLES `netswitch` WRITE;
/*!40000 ALTER TABLE `netswitch` DISABLE KEYS */;
/*!40000 ALTER TABLE `netswitch` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `permission`
--

DROP TABLE IF EXISTS `permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `permission` (
  `id` bigint(11) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(11) NOT NULL,
  `machine_id` bigint(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=310 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `permission`
--

LOCK TABLES `permission` WRITE;
/*!40000 ALTER TABLE `permission` DISABLE KEYS */;
INSERT INTO `permission` VALUES (183,2,2),(184,2,3),(188,2,13),(308,1,1),(309,1,2);
/*!40000 ALTER TABLE `permission` ENABLE KEYS */;
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
  `invoice_addr` int(11) NOT NULL,
  `ship_addr` int(11) NOT NULL,
  `client_id` int(11) NOT NULL,
  `b2b` tinyint(1) NOT NULL,
  `company` varchar(100) NOT NULL DEFAULT '',
  `vat_user_id` varchar(100) NOT NULL DEFAULT '',
  `vat_rate` int(11) NOT NULL,
  `user_role` varchar(100) NOT NULL DEFAULT 'member',
  `created` datetime DEFAULT NULL,
  `comments` text,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES (1,'Regular','User','user','user@example.com',0,0,0,0,'','',0,'',NULL,''),(2,'Regular','Admin','admin','admin@example.com',0,0,0,0,'','',0,'admin',NULL,NULL);
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
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
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_membership`
--

LOCK TABLES `user_membership` WRITE;
/*!40000 ALTER TABLE `user_membership` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_membership` ENABLE KEYS */;
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
INSERT INTO `user_roles` VALUES (1,0,0,1),(2,1,0,0);
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

-- Dump completed on 2015-06-22 17:00:22
