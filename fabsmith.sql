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
-- Table structure for table `activation`
--

DROP TABLE IF EXISTS `activation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `activation` (
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
) ENGINE=InnoDB AUTO_INCREMENT=111 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `activation`
--

LOCK TABLES `activation` WRITE;
/*!40000 ALTER TABLE `activation` DISABLE KEYS */;
INSERT INTO `activation` VALUES (50,NULL,1,2,0,'2014-10-24 20:13:10','2014-10-24 20:14:29',79,0,0,0,0,'',0,0),(51,NULL,1,2,0,'2014-10-24 20:41:46','2014-10-24 20:42:20',35,0,0,0,0,'',0,0),(52,NULL,1,2,0,'2014-10-24 20:42:38','2014-10-24 20:42:53',16,0,0,0,0,'',0,0),(53,NULL,1,2,0,'2014-10-27 17:06:52','2014-10-27 17:07:10',18,0,0,0,0,'',0,0),(54,NULL,1,2,0,'2014-10-27 18:32:32','2014-10-27 18:32:46',15,0,0,0,0,'',0,0),(55,NULL,1,2,0,'2014-10-27 18:39:04','2014-10-27 18:39:16',12,0,0,0,0,'',0,0),(56,NULL,1,2,0,'2014-10-27 18:48:17','2014-10-27 18:48:29',13,0,0,0,0,'',0,0),(57,NULL,1,2,0,'2014-10-27 18:58:27','2014-10-27 18:58:39',13,0,0,0,0,'',0,0),(58,NULL,2,2,0,'2014-10-30 14:35:38','2014-10-30 18:25:58',13821,0,0,0,0,'',0,0),(71,NULL,1,1,0,'2014-10-30 16:15:52','2014-10-30 17:23:37',4065,0,0,0,0,'',0,0),(72,NULL,1,1,0,'2014-10-30 17:23:49','2014-10-30 17:24:09',20,0,0,0,0,'',0,0),(73,NULL,1,1,0,'2014-10-30 17:25:12','2014-10-30 17:28:42',211,0,0,0,0,'',0,0),(74,NULL,1,1,0,'2014-10-30 17:28:56','2014-10-30 17:29:04',8,0,0,0,0,'',0,0),(75,NULL,1,1,0,'2014-10-30 17:31:53','2014-10-30 17:31:58',5,0,0,0,0,'',0,0),(76,NULL,1,1,0,'2014-10-30 17:33:55','2014-10-30 17:34:00',6,0,0,0,0,'',0,0),(77,NULL,1,3,0,'2014-10-30 17:34:29','2014-10-30 17:34:44',16,0,0,0,0,'',0,0),(78,NULL,1,1,0,'2014-10-30 17:35:07','2014-10-30 17:35:18',11,0,0,0,0,'',0,0),(79,NULL,1,1,0,'2014-10-30 17:37:24','2014-10-30 17:37:39',15,0,0,0,0,'',0,0),(80,NULL,1,1,0,'2014-10-30 17:38:21','2014-10-30 17:38:33',12,0,0,0,0,'',0,0),(81,NULL,1,1,0,'2014-10-30 17:38:47','2014-10-30 17:38:52',6,0,0,0,0,'',0,0),(82,NULL,1,1,0,'2014-10-30 17:38:55','2014-10-30 17:38:58',4,0,0,0,0,'',0,0),(83,NULL,1,1,0,'2014-10-30 17:41:47','2014-10-30 17:43:00',73,0,0,0,0,'',0,0),(84,NULL,1,1,0,'2014-10-30 17:44:11','2014-10-30 17:44:57',46,0,0,0,0,'',0,0),(85,NULL,1,1,0,'2014-10-30 17:50:54','2014-10-30 18:25:23',2070,0,0,0,0,'',0,0),(86,NULL,1,3,0,'2014-10-30 18:25:15','2014-10-30 18:25:40',25,0,0,0,0,'',0,0),(87,NULL,2,2,0,'2014-10-30 18:26:09','2014-10-30 18:26:13',5,0,0,0,0,'',0,0),(88,NULL,1,2,0,'2014-10-30 18:26:39','2014-10-30 18:32:55',377,0,0,0,0,'',0,0),(89,NULL,1,1,0,'2014-10-30 18:32:47','2014-10-30 18:32:58',11,0,0,0,0,'',0,0),(90,NULL,1,1,0,'2014-10-30 18:33:40','2014-10-30 18:34:30',50,0,0,0,0,'',0,0),(91,NULL,1,2,0,'2014-10-30 18:34:21','2014-10-30 18:34:32',12,0,0,0,0,'',0,0),(92,NULL,1,1,0,'2014-10-30 18:34:46','2014-10-30 18:52:02',1036,0,0,0,0,'',0,0),(93,NULL,1,2,0,'2014-10-30 18:51:22','2014-10-30 18:52:22',60,0,0,0,0,'',0,0),(94,NULL,2,2,0,'2014-10-30 18:56:49','2014-10-31 13:14:21',65852,0,0,0,0,'',0,0),(95,NULL,2,5,0,'2014-10-31 13:13:32','2014-10-31 14:29:06',4535,0,0,0,0,'',0,0),(96,NULL,1,1,0,'2014-10-31 13:53:42','2014-10-31 14:07:18',816,0,0,0,0,'',0,0),(97,NULL,1,2,0,'2014-10-31 14:04:54','2014-10-31 14:05:08',14,0,0,0,0,'',0,0),(98,NULL,1,2,0,'2014-10-31 14:05:13','2014-10-31 14:28:39',1407,0,0,0,0,'',0,0),(99,NULL,2,3,0,'2014-10-31 14:06:21','2014-10-31 14:33:25',1624,0,0,0,0,'',0,0),(100,NULL,1,1,0,'2014-10-31 14:32:35','2014-10-31 19:29:54',17840,0,0,0,0,'',0,0),(101,NULL,2,2,0,'2014-10-31 14:34:08','2014-10-31 19:27:49',17622,0,0,0,0,'',0,0),(102,NULL,1,3,0,'2014-10-31 17:59:30','2014-10-31 18:00:45',76,0,0,0,0,'',0,0),(103,NULL,1,2,0,'2014-10-31 19:28:39','2014-10-31 19:30:28',110,0,0,0,0,'',0,0),(104,NULL,1,1,0,'2014-11-17 20:12:47','2014-12-01 12:24:05',1181478,0,0,0,0,'',0,0),(105,NULL,2,3,0,'2014-11-17 20:13:24','2014-12-01 12:24:24',1181461,0,0,0,0,'',0,0),(106,NULL,1,2,0,'2014-11-30 17:13:54','2014-12-01 12:24:10',69016,0,0,0,0,'',0,0),(107,NULL,1,3,0,'2014-12-01 12:24:37','2014-12-01 12:59:53',2116,0,0,0,0,'',0,0),(108,NULL,1,1,0,'2014-12-01 12:52:26','2014-12-01 12:52:30',4,0,0,0,0,'',0,0),(109,NULL,1,2,0,'2014-12-01 14:37:15','2014-12-09 14:19:11',690117,0,0,0,0,'',0,0),(110,NULL,1,2,1,'2014-12-09 14:19:23',NULL,0,0,0,0,0,'',0,0);
/*!40000 ALTER TABLE `activation` ENABLE KEYS */;
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
  `password` varchar(100) NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `auth`
--

LOCK TABLES `auth` WRITE;
/*!40000 ALTER TABLE `auth` DISABLE KEYS */;
INSERT INTO `auth` VALUES (1,'0','bfd59291e825b5f2bbf1eb76569f8fe7'),(2,'0','bfd59291e825b5f2bbf1eb76569f8fe7');
/*!40000 ALTER TABLE `auth` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `machine`
--

DROP TABLE IF EXISTS `machine`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `machine` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT '',
  `description` text NOT NULL,
  `available` tinyint(1) unsigned NOT NULL,
  `unavail_msg` text NOT NULL,
  `unavail_till` datetime NOT NULL,
  `calc_by_energy` tinyint(1) unsigned NOT NULL,
  `calc_by_time` tinyint(1) unsigned NOT NULL,
  `costs_per_kwh` float unsigned NOT NULL,
  `costs_per_min` float unsigned NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `machine`
--

LOCK TABLES `machine` WRITE;
/*!40000 ALTER TABLE `machine` DISABLE KEYS */;
INSERT INTO `machine` VALUES (1,'i3berlin 3D Printer','The tools you make. Your tools, your make.',1,'','0000-00-00 00:00:00',0,1,0,0.2),(2,'MakerBot 3D printer','NYC 3D printer 4 real and 4 life.',0,'','0000-00-00 00:00:00',0,1,0,0.2),(3,'Zing Laser Cutter','Cuts wood, plastic, paper. Fast.',1,'','0000-00-00 00:00:00',0,1,0,1),(4,'CNC Router','Cuts steel, plutanium, uranium. Drill on steroids.',0,'','0000-00-00 00:00:00',0,1,0,1),(5,'Hand Drill','A man is a man if he does not know how to handle one.',1,'','0000-00-00 00:00:00',0,1,0,0.2);
/*!40000 ALTER TABLE `machine` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `permission`
--

DROP TABLE IF EXISTS `permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `permission` (
  `user_id` int(11) unsigned NOT NULL,
  `machine_id` int(11) unsigned NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `permission`
--

LOCK TABLES `permission` WRITE;
/*!40000 ALTER TABLE `permission` DISABLE KEYS */;
INSERT INTO `permission` VALUES (1,1),(1,2),(1,3),(2,2),(2,3),(2,4),(2,5);
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
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES (1,'Krisjanis','Rijnieks','kris','krisjanis.rijnieks@gmail.com',0,0,0,0,'','',0),(2,'Kruger','Ultimus','kruxy','kru@xy.io',0,0,0,0,'','',0);
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
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
INSERT INTO `user_roles` VALUES (1,0,0,0);
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

-- Dump completed on 2014-12-09 16:35:53
