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
) ENGINE=InnoDB AUTO_INCREMENT=197 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `activations`
--

LOCK TABLES `activations` WRITE;
/*!40000 ALTER TABLE `activations` DISABLE KEYS */;
INSERT INTO `activations` VALUES (100,NULL,7,2,0,'2015-06-25 17:40:13','2015-06-25 17:40:16',3,0,0,0,0,'',0,0,16,'€','hour'),(101,NULL,7,3,0,'2015-06-25 17:40:22','2015-06-25 17:40:27',6,0,0,0,0,'',0,0,1,'€','minute'),(102,NULL,7,3,0,'2015-06-25 17:40:30','2015-06-25 17:40:32',3,0,0,0,0,'',0,0,1,'€','minute'),(103,NULL,7,3,0,'2015-07-15 15:50:40','2015-07-15 15:50:44',4,0,0,0,0,'',0,0,1,'€','minute'),(104,NULL,7,2,0,'2015-07-15 15:51:04','2015-07-15 15:51:10',6,0,0,0,0,'',0,0,16,'€','hour'),(105,NULL,7,2,0,'2015-07-15 15:51:18','2015-07-15 15:51:24',7,0,0,0,0,'',0,0,16,'€','hour'),(106,NULL,7,2,0,'2015-07-15 15:51:32','2015-07-15 15:51:37',5,0,0,0,0,'',0,0,16,'€','hour'),(107,NULL,7,2,0,'2015-07-15 18:04:57','2015-07-15 18:05:17',20,0,0,0,0,'',0,0,16,'€','hour'),(108,NULL,7,1,0,'2015-07-15 18:07:34','2015-07-15 18:07:45',11,0,0,0,0,'',0,0,0.1,'€','minute'),(109,NULL,7,1,0,'2015-07-15 18:07:48','2015-07-15 18:07:56',8,0,0,0,0,'',0,0,0.1,'€','minute'),(110,NULL,7,1,0,'2015-07-15 18:08:05','2015-07-15 18:08:21',17,0,0,0,0,'',0,0,0.1,'€','minute'),(111,NULL,7,2,0,'2015-07-28 11:42:17','2015-07-28 11:42:28',12,0,0,0,0,'',0,0,16,'€','hour'),(112,NULL,7,1,0,'2015-07-28 11:42:33','2015-07-28 11:42:39',6,0,0,0,0,'',0,0,0.1,'€','minute'),(113,NULL,7,1,0,'2015-07-28 11:42:42','2015-07-28 11:42:46',4,0,0,0,0,'',0,0,0.1,'€','minute'),(114,NULL,7,1,0,'2015-07-28 11:42:48','2015-07-28 11:42:50',3,0,0,0,0,'',0,0,0.1,'€','minute'),(115,NULL,7,1,0,'2015-07-28 11:42:54','2015-07-28 11:43:04',10,0,0,0,0,'',0,0,0.1,'€','minute'),(116,NULL,7,1,0,'2015-07-28 11:47:19','2015-07-28 11:47:23',4,0,0,0,0,'',0,0,0.1,'€','minute'),(117,NULL,7,1,0,'2015-07-28 11:47:24','2015-07-28 11:47:26',3,0,0,0,0,'',0,0,0.1,'€','minute'),(118,NULL,7,2,0,'2015-07-28 11:47:35','2015-07-28 11:47:38',4,0,0,0,0,'',0,0,16,'€','hour'),(119,NULL,1,2,0,'2015-07-29 12:06:18','2015-07-29 12:06:32',15,0,0,0,0,'',0,0,16,'€','hour'),(120,NULL,1,1,0,'2015-07-29 12:06:34','2015-07-29 12:06:42',8,0,0,0,0,'',0,0,0.1,'€','minute'),(121,NULL,2,1,0,'2015-07-29 12:06:58','2015-07-29 12:07:07',10,0,0,0,0,'',0,0,0.1,'€','minute'),(122,NULL,1,1,0,'2015-08-18 16:40:45','2015-08-18 16:40:53',8,0,0,0,0,'',0,0,0.1,'€','minute'),(123,NULL,1,2,0,'2015-08-18 16:40:58','2015-08-18 16:42:53',115,0,0,0,0,'',0,0,16,'€','hour'),(124,NULL,1,2,0,'2015-08-18 16:42:54','2015-08-18 16:43:10',16,0,0,0,0,'',0,0,16,'€','hour'),(125,NULL,1,2,0,'2015-08-18 16:43:11','2015-08-18 16:44:56',105,0,0,0,0,'',0,0,16,'€','hour'),(126,NULL,1,1,0,'2015-08-18 17:15:59','2015-08-18 17:30:24',865,0,0,0,0,'',0,0,0.1,'€','minute'),(127,NULL,1,1,0,'2015-08-18 18:24:38','2015-08-18 18:25:00',22,0,0,0,0,'',0,0,0.1,'€','minute'),(128,NULL,2,1,0,'2015-09-17 11:37:20','2015-09-17 11:37:32',13,0,0,0,0,'',0,0,0.1,'€','minute'),(129,NULL,2,1,0,'2015-09-17 11:37:34','2015-09-17 11:37:38',4,0,0,0,0,'',0,0,0.1,'€','minute'),(130,NULL,2,1,0,'2015-09-17 11:49:02','2015-09-17 11:49:06',4,0,0,0,0,'',0,0,0.1,'€','minute'),(131,NULL,2,1,0,'2015-09-17 11:52:32','2015-09-17 11:54:18',106,0,0,0,0,'',0,0,0.1,'€','minute'),(132,NULL,2,1,0,'2015-09-17 11:54:19','2015-09-17 11:54:23',4,0,0,0,0,'',0,0,0.1,'€','minute'),(133,NULL,2,1,0,'2015-09-17 11:54:24','2015-09-17 11:54:28',5,0,0,0,0,'',0,0,0.1,'€','minute'),(134,NULL,2,1,0,'2015-09-17 11:54:41','2015-09-17 11:54:43',2,0,0,0,0,'',0,0,0.1,'€','minute'),(135,NULL,2,1,0,'2015-09-17 11:54:45','2015-09-17 11:54:49',4,0,0,0,0,'',0,0,0.1,'€','minute'),(136,NULL,2,1,0,'2015-09-17 11:55:17','2015-09-17 11:59:24',248,0,0,0,0,'',0,0,0.1,'€','minute'),(137,NULL,2,1,0,'2015-09-17 11:59:25','2015-09-17 11:59:39',14,0,0,0,0,'',0,0,0.1,'€','minute'),(138,NULL,2,1,0,'2015-09-17 11:59:43','2015-09-17 11:59:46',3,0,0,0,0,'',0,0,0.1,'€','minute'),(139,NULL,2,1,0,'2015-09-17 12:03:41','2015-09-17 12:03:43',3,0,0,0,0,'',0,0,0.1,'€','minute'),(140,NULL,2,1,0,'2015-09-17 12:03:44','2015-09-17 12:03:46',3,0,0,0,0,'',0,0,0.1,'€','minute'),(141,NULL,2,1,0,'2015-09-17 12:06:00','2015-09-17 12:06:02',2,0,0,0,0,'',0,0,0.1,'€','minute'),(142,NULL,2,1,0,'2015-09-17 12:06:13','2015-09-17 17:23:02',19010,0,0,0,0,'',0,0,0.1,'€','minute'),(143,NULL,1,2,0,'2015-09-17 12:10:31','2015-09-17 12:14:07',217,0,0,0,0,'',0,0,16,'€','hour'),(144,NULL,1,2,0,'2015-09-17 12:15:00','2015-09-17 12:15:03',3,0,0,0,0,'',0,0,16,'€','hour'),(145,NULL,1,2,0,'2015-09-17 12:15:07','2015-09-17 12:15:09',3,0,0,0,0,'',0,0,16,'€','hour'),(146,NULL,1,2,0,'2015-09-17 12:15:28','2015-09-17 12:16:31',63,0,0,0,0,'',0,0,16,'€','hour'),(147,NULL,1,2,0,'2015-09-17 12:16:41','2015-09-17 13:24:09',4049,0,0,0,0,'',0,0,16,'€','hour'),(148,NULL,2,3,0,'2015-09-17 12:37:10','2015-09-17 12:37:12',2,0,0,0,0,'',0,0,1,'€','minute'),(149,NULL,2,2,0,'2015-09-17 13:24:18','2015-09-17 13:24:21',3,0,0,0,0,'',0,0,16,'€','hour'),(150,NULL,1,2,0,'2015-09-17 13:24:28','2015-09-17 18:11:29',17221,0,0,0,0,'',0,0,16,'€','hour'),(151,NULL,2,1,0,'2015-09-17 17:23:06','2015-09-17 18:13:24',3018,0,0,0,0,'',0,0,0.1,'€','minute'),(152,NULL,1,2,0,'2015-09-17 18:13:16','2015-09-17 18:13:18',3,0,0,0,0,'',0,0,16,'€','hour'),(153,NULL,1,2,0,'2015-09-17 18:39:25','2015-09-17 18:40:02',37,0,0,0,0,'',0,0,16,'€','hour'),(154,NULL,2,1,0,'2015-09-17 18:39:45','2015-09-17 18:40:15',30,0,0,0,0,'',0,0,0.1,'€','minute'),(155,NULL,2,1,0,'2015-09-17 18:40:17','2015-09-17 19:03:52',1416,0,0,0,0,'',0,0,0.1,'€','minute'),(156,NULL,1,2,0,'2015-09-17 18:40:22','2015-09-17 19:04:18',1437,0,0,0,0,'',0,0,16,'€','hour'),(157,NULL,2,1,0,'2015-09-17 19:03:56','2015-09-17 19:04:00',4,0,0,0,0,'',0,0,0.1,'€','minute'),(158,NULL,1,1,0,'2015-09-17 19:04:04','2015-09-17 19:04:10',7,0,0,0,0,'',0,0,0.1,'€','minute'),(159,NULL,1,2,0,'2015-09-17 19:09:40','2015-09-17 19:23:26',827,0,0,0,0,'',0,0,16,'€','hour'),(160,NULL,1,1,0,'2015-09-17 19:11:23','2015-09-17 19:11:27',5,0,0,0,0,'',0,0,0.1,'€','minute'),(161,NULL,1,1,0,'2015-09-17 19:23:01','2015-09-17 19:23:03',2,0,0,0,0,'',0,0,0.1,'€','minute'),(162,0,2,1,0,'2015-09-17 19:23:14','2015-10-13 12:43:10',2222396,0,0,0,0,'',0,0,0.1,'€','minute'),(163,NULL,1,2,0,'2015-09-17 19:23:31','2015-09-18 14:10:43',67633,0,0,0,0,'',0,0,16,'€','hour'),(164,NULL,2,2,0,'2015-09-18 14:11:28','2015-09-18 14:11:35',8,0,0,0,0,'',0,0,16,'€','hour'),(165,NULL,2,2,0,'2015-09-18 14:11:36','2015-09-18 14:11:41',5,0,0,0,0,'',0,0,16,'€','hour'),(166,0,1,2,0,'2015-10-11 11:09:41','2015-10-11 12:09:53',3611,0,0,0,0,'',0,0,16,'€','hour'),(167,0,1,2,0,'2015-10-11 12:09:54','2015-10-13 09:55:26',164731,0,0,0,0,'',0,0,16,'€','hour'),(168,0,1,2,0,'2015-10-14 13:22:37','2015-10-14 13:22:42',5,0,0,0,0,'',0,0,16,'€','hour'),(169,0,2,1,0,'2015-10-14 13:42:46','2015-10-14 13:42:49',2,0,0,0,0,'',0,0,0.1,'€','minute'),(170,0,1,2,0,'2015-10-14 13:49:04','2015-10-14 13:49:16',11,0,0,0,0,'',0,0,16,'€','hour'),(171,0,1,2,0,'2015-10-14 13:50:45','2015-10-14 13:50:56',10,0,0,0,0,'',0,0,16,'€','hour'),(172,0,1,2,0,'2015-10-14 13:51:01','2015-10-14 13:51:08',6,0,0,0,0,'',0,0,16,'€','hour'),(173,0,1,2,0,'2015-10-14 13:51:21','2015-10-14 13:51:31',10,0,0,0,0,'',0,0,16,'€','hour'),(174,0,2,1,0,'2015-10-14 13:52:48','2015-10-14 13:52:56',8,0,0,0,0,'',0,0,0.1,'€','minute'),(175,0,1,2,0,'2015-10-14 13:53:33','2015-10-14 13:53:34',1,0,0,0,0,'',0,0,16,'€','hour'),(176,0,2,2,0,'2015-10-14 13:56:19','2015-10-14 13:56:21',1,0,0,0,0,'',0,0,16,'€','hour'),(177,0,1,2,0,'2015-10-14 14:13:03','2015-10-14 14:13:05',1,0,0,0,0,'',0,0,16,'€','hour'),(178,0,1,1,0,'2015-10-15 08:10:56','2015-10-15 08:11:27',30,0,0,0,0,'',0,0,0.1,'€','minute'),(179,0,1,2,0,'2015-10-15 08:12:45','2015-10-15 08:15:55',189,0,0,0,0,'',0,0,16,'€','hour'),(180,0,2,1,0,'2015-10-15 08:15:46','2015-10-15 08:15:50',3,0,0,0,0,'',0,0,0.1,'€','minute'),(181,0,2,2,0,'2015-10-15 08:15:57','2015-10-15 08:16:02',4,0,0,0,0,'',0,0,16,'€','hour'),(182,0,2,3,0,'2015-10-15 12:21:06','2015-10-15 12:21:40',33,0,0,0,0,'',0,0,1,'€','minute'),(183,0,1,2,0,'2015-10-15 14:53:37','2015-10-16 08:07:43',62045,0,0,0,0,'',0,0,16,'€','hour'),(184,0,2,2,0,'2015-10-16 08:07:49','2015-10-16 08:07:52',3,0,0,0,0,'',0,0,16,'€','hour'),(185,0,1,1,0,'2015-10-16 08:11:09','2015-10-16 08:11:12',3,0,0,0,0,'',0,0,0.1,'€','minute'),(186,0,1,1,0,'2015-10-18 12:05:46','2015-10-18 12:05:51',5,0,0,0,0,'',0,0,0.1,'€','minute'),(187,0,1,1,0,'2015-10-22 15:24:05','2015-10-22 15:24:10',4,0,0,0,0,'',0,0,0.1,'€','minute'),(188,0,1,1,0,'2015-10-23 07:24:15','2015-10-23 07:26:52',157,0,0,0,0,'',0,0,0.1,'€','minute'),(189,0,1,1,0,'2015-10-23 07:26:54','2015-10-23 07:27:01',7,0,0,0,0,'',0,0,0.1,'€','minute'),(190,0,1,1,0,'2015-10-23 07:33:54','2015-10-23 07:34:15',20,0,0,0,0,'',0,0,0.1,'€','minute'),(191,0,1,1,0,'2015-10-23 07:52:01','2015-10-23 07:52:03',2,0,0,0,0,'',0,0,0.1,'€','minute'),(192,0,1,1,0,'2015-10-23 08:55:54','2015-10-23 08:55:59',5,0,0,0,0,'',0,0,0.1,'€','minute'),(193,0,1,2,0,'2015-10-23 08:56:20','2015-10-23 08:56:23',3,0,0,0,0,'',0,0,16,'€','hour'),(194,0,2,1,0,'2015-10-27 09:32:11','2015-10-27 09:32:20',9,0,0,0,0,'',0,0,0.1,'€','minute'),(195,0,1,3,0,'2015-10-28 09:43:28','2015-10-28 09:43:32',3,0,0,0,0,'',0,0,1,'€','minute'),(196,0,1,1,0,'2015-11-19 15:07:32','2015-11-20 16:30:09',91357,0,0,0,0,'',0,0,0.1,'€','minute');
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
INSERT INTO `auth` VALUES (1,'1234','adfda4a8476a98568a5542d239dee848ccd497c9d425674ed8a59e4d1eda2871065dd89982c1c75204a8451015ec6f68a46cc73c5908c22c24a88a7595d0f5b3','5b55c9a96b1bdac9cac54ab52b5bf32278dc56dd1a25ca17332f7a24ded714e3'),(2,'04126DBA233580','f2ae94748430030d13b43359d021411e97c5b49aa7fcad70134abc27ccdd4c24de096143810bffaf824d9b5dc028bb3de8f144ee48ea06c4e0895d8a6df6cdac','f755192ae0fcb52ef329b0007ccd8f1e68d040d7801a32e96db24b047896c27d'),(4,'','827300af9fa28ba6f0f64c7a421991448eef8cf786bfce1e0fa48561840a33db3952e3a1b3d9dd126c2efcdddc7db809d13eb8dcc924fa63408b6a12bb1d86b7','a4b853b80e157787be41777ad41cb5c8d5db81e73b2300e620b8390b130a79c1'),(5,'','7a29f5e8ca373e158a148a61ea63187027c876418653606486baa5107167fb2fb209cbb723d45005966c52e35f24da3ee1d8d53ebd020ed7c9a8e31310ca1780','59dce8439c5c4394fe8ff7cb5978303c0bd9cb40eb72f019e2482b3fbc70d34a'),(6,'','0a23a7c66c19ebea6eb276a72692f219df5d1fdcfd6fae5204b97b9ee2ebe7d078e1dd08376abcc72f2c1544d77e43688ca4af94ae95deaca6d19ec24d52fdce','89adaf9a8ea22b704531214869259af9415941206d2a688b52c7859c2487d4fe'),(7,'','93faef84e77ca4f3d505dbb80986646b31f000846f0f5036d140f182fd2458d07af57f4d367ea6c3c094ed5433b825f128daf157acfbfe5a3c724aa4383cdac1','40809e87d685d2d6142e6f430871ad12070b26d23cab2d4f0c80791e6917e8bd'),(8,'','1cfd1395f1601521f5b5ccd2b035f27370ce183e168ef180b83d99ab55fa2e51dd6e3282d44103c2f839e5c984058b5bd510ab6d3bc43e89223b0da518563097','a0b4ca2a8f9bb0dd5186b0c8e838b7dc79eae961b8b0a5929994db7d6693f31a'),(9,'','2ec67587989772fb9f9b27ff4fdc93b1a25be4c9f9690f71920562958c4efc68116fd2e31f815d155d27f03246d98cd9d4c17dafe051182d2fcc834c9f239af7','3e20ed519b850e8df932357114c433e86a4bd1bf6526f9993568912e9e654c4f');
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
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `invoices`
--

LOCK TABLES `invoices` WRITE;
/*!40000 ALTER TABLE `invoices` DISABLE KEYS */;
INSERT INTO `invoices` VALUES (1,'[110,115,108,113,116,109,114,112,117,100,105,118,106,111,104,107,103,101,102]','files/invoice-20150401-20150729-XzCvxyJyFh.xlsx','2015-07-29 11:20:01','2015-04-01 00:00:00','2015-07-29 23:59:59'),(2,'[110,115,108,113,116,109,114,112,117,100,105,118,106,111,104,107,103,101,102]','files/invoice-20150301-20150729-FbtXhnhLcV.xlsx','2015-07-29 11:23:45','2015-03-01 00:00:00','2015-07-29 23:59:59'),(3,'[110,115,108,113,116,109,114,112,117,100,105,118,106,111,104,107,103,101,102]','files/invoice-20150401-20150729-FhvyAPvcag.xlsx','2015-07-29 11:26:32','2015-04-01 00:00:00','2015-07-29 23:59:59'),(4,'[108,113,116,109,114,112,117,110,115,118,106,111,104,107,105,103,121,120,119]','files/invoice-20150701-20150805-jdtXYYawew.xlsx','2015-08-05 11:39:35','2015-07-01 00:00:00','2015-08-05 23:59:59'),(5,'[110,115,108,113,116,109,114,112,117,100,105,118,106,111,104,107,103,101,102,121,120,119]','files/invoice-20150601-20150820-SqTcWObyGz.xlsx','2015-08-08 23:56:09','2015-06-01 00:00:00','2015-08-20 23:59:59'),(6,'[110,115,108,113,116,109,114,112,117,100,105,118,106,111,104,107,103,101,102,121,120,126,122,127,125,123,119,124]','files/invoice-20150101-20150914-KfYmEHRJAU.xlsx','2015-09-14 18:04:15','2015-01-01 00:00:00','2015-09-14 23:59:59'),(7,'[110,115,108,113,116,109,114,112,117,100,105,118,106,111,104,107,103,101,102,121,120,126,122,127,125,123,119,124]','files/invoice-20150101-20150914-PrvCMiGkta.xlsx','2015-09-14 18:07:10','2015-01-01 00:00:00','2015-09-14 23:59:59'),(8,'[128,133,138,131,136,141,151,129,134,139,154,132,137,142,157,130,135,140,155,164,149,165,148,160,158,161,143,153,146,156,144,159,147,152,163,145,150]','files/invoice-20150901-20150930-TETaYessnH.xlsx','2015-09-21 14:30:58','2015-09-01 00:00:00','2015-09-30 23:59:59'),(9,'[122,123,124,125,126,127]','files/invoice-20150801-20150831-IXEfuxswGo.xlsx','2015-10-28 10:33:24','2015-08-01 00:00:00','2015-08-31 23:59:59'),(10,'[]','files/invoice-20151103-20151127-EWhTlgWKAD.xlsx','2015-11-27 17:03:34','2015-11-03 00:00:00','2015-11-27 23:59:59');
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
  `connected_machines` varchar(255) DEFAULT NULL,
  `switch_ref_count` int(11) DEFAULT NULL,
  `under_maintenance` tinyint(1) DEFAULT NULL,
  `reservation_price_start` double unsigned DEFAULT NULL,
  `reservation_price_hourly` double unsigned DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `machines`
--

LOCK TABLES `machines` WRITE;
/*!40000 ALTER TABLE `machines` DISABLE KEYS */;
INSERT INTO `machines` VALUES (1,'Laydrop 3D Printer','MB3DP','NYC 3D printer 4 real and 4 life.','',0,'',NULL,0.1,'minute','',1,'[]',1,0,NULL,5),(2,'MakerBot 3D Printer','MB3DP','NYC 3D printer 4 real and 4 life.','machine-2.svg',0,'',NULL,16,'hour','',1,'',1,0,NULL,5),(3,'Zing Laser Cutter','ZLC','Cuts wood, plastic, paper. Fast.','',1,'',NULL,1,'minute','Machine related comments',1,'[]',2,0,NULL,5);
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
  `duration_months` int(11) DEFAULT NULL,
  `monthly_price` double unsigned NOT NULL,
  `machine_price_deduction` int(11) NOT NULL,
  `affected_machines` text,
  `auto_extend` tinyint(1) DEFAULT NULL,
  `auto_extend_duration_months` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `membership`
--

LOCK TABLES `membership` WRITE;
/*!40000 ALTER TABLE `membership` DISABLE KEYS */;
INSERT INTO `membership` VALUES (1,'6 Months Basic','6MB',7,75,50,'[1,2,3]',1,1),(2,'1 Month Basic','SMB',1,100,50,'[1,2,3]',1,1),(5,'Kuul Membership','',0,0,0,'',1,1),(6,'My Membership','',0,0,0,'',1,1),(7,'Test Membership','',0,0,0,'[]',1,1);
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
) ENGINE=InnoDB AUTO_INCREMENT=58 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `migrations`
--

LOCK TABLES `migrations` WRITE;
/*!40000 ALTER TABLE `migrations` DISABLE KEYS */;
INSERT INTO `migrations` VALUES (1,'Billingingaddress_20150728_120222','2015-07-28 10:16:39','ALTER TABLE user MODIFY invoice_addr TEXT','ALTER TABLE user MODIFY invoice_addr INT(11); ALTER TABLE user MODIFY ship_addr INT(11)','rollback'),(2,'Billingingaddress_20150728_120222','2015-07-28 10:16:39','ALTER TABLE user MODIFY invoice_addr TEXT; ALTER TABLE user MODIFY ship_addr TEXT','ALTER TABLE user MODIFY invoice_addr INT(11); ALTER TABLE user MODIFY ship_addr INT(11)','rollback'),(3,'Billingingaddress_20150728_120222','2015-07-28 10:16:51','ALTER TABLE user MODIFY invoice_addr TEXT; ALTER TABLE user MODIFY ship_addr TEXT',NULL,'update'),(4,'Userphone_20150728_152059','2015-07-28 13:28:03','ALTER TABLE user ADD COLUMN phone VARCHAR(50)','ALTER TABLE user DROP COLUMN phone','rollback'),(5,'Userphone_20150728_152059','2015-07-28 13:34:50','ALTER TABLE user ADD COLUMN phone VARCHAR(50)',NULL,'update'),(6,'Fastbilluseradd_20150805_170318','2015-08-05 15:04:40','ALTER TABLE user ADD COLUMN zip_code VARCHAR(100); ALTER TABLE user ADD COLUMN city VARCHAR(100); ALTER TABLE user ADD COLUMN country_code VARCHAR(2)','ALTER TABLE user DROP COLUMN zip_code; ALTER TABLE user DROP COLUMN city; ALTER TABLE user DROP COLUMN country_code','rollback'),(7,'Fastbilluseradd_20150805_170318','2015-08-05 15:16:09','ALTER TABLE user ADD COLUMN zip_code VARCHAR(100); ALTER TABLE user ADD COLUMN city VARCHAR(100); ALTER TABLE user ADD COLUMN country_code VARCHAR(2)',NULL,'update'),(8,'PricePerMonth_20150901_110810','2015-09-07 12:30:37','ALTER TABLE membership ADD monthly_price double unsigned NOT NULL AFTER price; UPDATE membership SET monthly_price = price WHERE duration = 30 AND unit = \'days\'; UPDATE membership SET monthly_price = price / 3 WHERE duration = 90 AND unit = \'days\'; UPDATE membership SET monthly_price = price / 12 WHERE duration = 365 AND unit = \'days\'; UPDATE membership SET monthly_price = price / duration * 30 WHERE duration <> 30 AND duration <> 90 AND duration <> 365 AND unit = \'days\'; UPDATE membership SET monthly_price = price WHERE unit <> \'days\'; ALTER TABLE membership DROP COLUMN price',NULL,'update'),(9,'Activationfeedback_20150908_145935','2015-09-08 14:48:47','\n		CREATE TABLE activation_feedback (\n			id bigint(20) unsigned NOT NULL AUTO_INCREMENT,\n			activation_id int(11) NOT NULL,\n			satisfaction varchar(100) DEFAULT NULL,\n			PRIMARY KEY (id)\n	)',NULL,'update'),(10,'Autoextendmemberships_20150908_172301','2015-09-08 15:42:10','ALTER TABLE membership ADD COLUMN auto_extend TINYINT(1); ALTER TABLE membership ADD COLUMN auto_extend_duration INT(11)','ALTER TABLE membership DROP COLUMN auto_extend; ALTER TABLE membership DROP COLUMN auto_extend_duration','rollback'),(11,'Autoextendmemberships_20150908_172301','2015-09-08 15:42:10','ALTER TABLE membership ADD COLUMN auto_extend TINYINT(1); ALTER TABLE membership ADD COLUMN auto_extend_duration INT(11)','ALTER TABLE membership DROP COLUMN auto_extend; ALTER TABLE membership DROP COLUMN auto_extend_duration','rollback'),(12,'Autoextendmemberships_20150908_172301','2015-09-08 15:42:10','ALTER TABLE membership ADD COLUMN auto_extend TINYINT(1) DEFAULT 1; ALTER TABLE membership ADD COLUMN auto_extend_duration INT(11) DEFAULT 30','ALTER TABLE membership DROP COLUMN auto_extend; ALTER TABLE membership DROP COLUMN auto_extend_duration','rollback'),(13,'Autoextendmemberships_20150908_172301','2015-09-08 15:42:16','ALTER TABLE membership ADD COLUMN auto_extend TINYINT(1); ALTER TABLE membership ADD COLUMN auto_extend_duration INT(11)',NULL,'update'),(14,'Usermembershipautoextend_20150909_110015','2015-09-09 09:01:54','ALTER TABLE user_membership ADD COLUMN auto_extend TINYINT(1)','ALTER TABLE user_membership DROP COLUMN auto_extend','rollback'),(15,'Usermembershipautoextend_20150909_110015','2015-09-09 09:02:02','ALTER TABLE user_membership ADD COLUMN auto_extend TINYINT(1)',NULL,'update'),(16,'Undermaintenance_20150909_113436','2015-09-21 09:51:36','ALTER TABLE machines ADD under_maintenance tinyint(1)','ALTER TABLE machines DROP COLUMN under_maintenance','rollback'),(17,'Undermaintenance_20150909_113436','2015-09-21 09:52:35','ALTER TABLE machines ADD under_maintenance tinyint(1)',NULL,'update'),(19,'Usermembershipautoextend_20150910_110015','2015-09-21 09:54:03','ALTER TABLE user_membership ADD COLUMN auto_extend TINYINT(1)',NULL,'update'),(20,'Autoextendmemberships_20150910_172301','2015-09-21 09:54:03','ALTER TABLE membership ADD COLUMN auto_extend TINYINT(1); ALTER TABLE membership ADD COLUMN auto_extend_duration INT(11)',NULL,'update'),(21,'Usermembershipterminate_20150918_115716','2015-09-21 09:59:22','ALTER TABLE user_membership ADD COLUMN is_terminated TINYINT(1)','ALTER TABLE user_membership DROP COLUMN is_terminated','rollback'),(22,'Autoextendmembershipmonths_20150921_115503','2015-09-21 09:59:11','ALTER TABLE membership CHANGE COLUMN auto_extend_duration auto_extend_duration_months INT(11)','ALTER TABLE membership CHANGE COLUMN auto_extend_duration_months auto_extend_duration INT(11)','rollback'),(23,'Usermembershipterminate_20150918_115716','2015-09-21 09:59:27','ALTER TABLE user_membership ADD COLUMN is_terminated TINYINT(1)',NULL,'update'),(24,'Autoextendmembershipmonths_20150921_115503','2015-09-21 09:59:27','ALTER TABLE membership CHANGE COLUMN auto_extend_duration auto_extend_duration_months INT(11)',NULL,'update'),(25,'Rmmembershipunitcol_20150921_123926','2015-09-21 10:54:49','ALTER TABLE membership CHANGE COLUMN duration duration_months INT(11); UPDATE membership SET duration_months=ROUND(duration_months / 30) WHERE unit=\'days\'; ALTER TABLE membership DROP COLUMN unit',NULL,'update'),(26,'Usermembershipenddate_20150924_114628','2015-10-13 09:50:59','\nUPDATE user_membership\nSET end_date = DATE_ADD(start_date, INTERVAL\n                          (SELECT duration_months\n                           FROM membership\n                           WHERE membership.id = user_membership.membership_id) MONTH)\n	',NULL,'update'),(27,'Autoextenddefaulttrue_20150924_120425','2015-10-13 09:50:59','UPDATE membership SET auto_extend = TRUE, auto_extend_duration_months = 1; UPDATE user_membership SET auto_extend = TRUE',NULL,'update'),(28,'Reservations_20150928_181415','2015-10-13 09:50:59','CREATE TABLE reservations (\n		id int(11) unsigned NOT NULL AUTO_INCREMENT,\n		machine_id int(11) NOT NULL,\n		user_id int(11) NOT NULL,\n		time_start datetime NOT NULL,\n		time_end datetime NOT NULL,\n		created datetime NOT NULL,\n		PRIMARY KEY (id)\n	)',NULL,'update'),(29,'Reservationrules_20151001_173125','2015-10-13 09:50:59','CREATE TABLE reservation_rules (\n		id int(11) unsigned NOT NULL AUTO_INCREMENT,\n		name varchar(100),\n		machine_id int(11),\n		available tinyint(1),\n		unavailable tinyint(1),\n		date_start char(10),\n		date_end char(10),\n		time_start char(5),\n		time_end char(5),\n		time_zone varchar(100),\n		monday tinyint(1),\n		tuesday tinyint(1),\n		wednesday tinyint(1),\n		thursday tinyint(1),\n		friday tinyint(1),\n		saturday tinyint(1),\n		sunday tinyint(1),\n		created datetime NOT NULL,\n		PRIMARY KEY (id)\n	)',NULL,'update'),(30,'Reservationprices_20151008_142149','2015-10-13 09:50:59','ALTER TABLE machines ADD COLUMN reservation_price_start double unsigned; ALTER TABLE machines ADD COLUMN reservation_price_hourly double unsigned',NULL,'update'),(31,'Activationcurrentprice_20151023_094534','2015-10-26 13:51:32','ALTER TABLE activations ADD COLUMN current_machine_price int(11); ALTER TABLE activations ADD COLUMN current_machine_price_currency varchar(10)','ALTER TABLE activations DROP COLUMN current_machine_price; ALTER TABLE activations DROP COLUMN current_machine_price_currency; ALTER TABLE activations DROP COLUMN current_machine_price_unit','rollback'),(32,'Activationcurrentprice_20151023_094534','2015-10-26 13:51:32','ALTER TABLE activations ADD COLUMN current_machine_price int(11); ALTER TABLE activations ADD COLUMN current_machine_price_currency varchar(10); ALTER TABLE activations ADD COLUMN current_machine_price_unit varchar(100)','ALTER TABLE activations DROP COLUMN current_machine_price; ALTER TABLE activations DROP COLUMN current_machine_price_currency; ALTER TABLE activations DROP COLUMN current_machine_price_unit','rollback'),(33,'Activationcurrentprice_20151023_094534','2015-10-26 13:51:32','ALTER TABLE activations ADD COLUMN current_machine_price double unsigned; ALTER TABLE activations ADD COLUMN current_machine_price_currency varchar(10); ALTER TABLE activations ADD COLUMN current_machine_price_unit varchar(100)','ALTER TABLE activations DROP COLUMN current_machine_price; ALTER TABLE activations DROP COLUMN current_machine_price_currency; ALTER TABLE activations DROP COLUMN current_machine_price_unit','rollback'),(34,'Reservationcurrentprice_20151026_120719','2015-10-26 13:53:13','ALTER TABLE reservations ADD COLUMN current_price double unsigned; ALTER TABLE reservations ADD COLUMN current_price_currency varchar(10); ALTER TABLE reservations ADD COLUMN current_price_unit varchar(100)','ALTER TABLE reservations DROP COLUMN current_price; ALTER TABLE reservations DROP COLUMN current_price_currency; ALTER TABLE reservations DROP COLUMN current_price_unit','rollback'),(35,'Reservationcurrentprice_20151026_120719','2015-10-26 13:53:13','ALTER TABLE reservations ADD COLUMN current_price double unsigned; ALTER TABLE reservations ADD COLUMN current_price_currency varchar(10); ALTER TABLE reservations ADD COLUMN current_price_unit varchar(100)','ALTER TABLE reservations DROP COLUMN current_price; ALTER TABLE reservations DROP COLUMN current_price_currency; ALTER TABLE reservations DROP COLUMN current_price_unit','rollback'),(36,'Activationcurrentprice_20151023_094534','2015-10-26 13:51:32','ALTER TABLE activations ADD COLUMN current_machine_price double unsigned; ALTER TABLE activations ADD COLUMN current_machine_price_currency varchar(10); ALTER TABLE activations ADD COLUMN current_machine_price_unit varchar(100); UPDATE activations a JOIN machines m ON a.machine_id = m.id SET a.current_machine_price=m.price, a.current_machine_price_currency=\'€\', a.current_machine_price_unit=m.price_unit','ALTER TABLE activations DROP COLUMN current_machine_price; ALTER TABLE activations DROP COLUMN current_machine_price_currency; ALTER TABLE activations DROP COLUMN current_machine_price_unit','rollback'),(37,'Reservationcurrentprice_20151026_120719','2015-10-26 13:53:13','ALTER TABLE reservations ADD COLUMN current_price double unsigned; ALTER TABLE reservations ADD COLUMN current_price_currency varchar(10); ALTER TABLE reservations ADD COLUMN current_price_unit varchar(100)','ALTER TABLE reservations DROP COLUMN current_price; ALTER TABLE reservations DROP COLUMN current_price_currency; ALTER TABLE reservations DROP COLUMN current_price_unit','rollback'),(38,'Activationcurrentprice_20151023_094534','2015-10-26 13:52:06','ALTER TABLE activations ADD COLUMN current_machine_price double unsigned; ALTER TABLE activations ADD COLUMN current_machine_price_currency varchar(10); ALTER TABLE activations ADD COLUMN current_machine_price_unit varchar(100); UPDATE activations a JOIN machines m ON a.machine_id = m.id SET a.current_machine_price=m.price, a.current_machine_price_currency=\'€\', a.current_machine_price_unit=m.price_unit',NULL,'update'),(39,'Reservationcurrentprice_20151026_120719','2015-10-26 13:53:13','ALTER TABLE reservations ADD COLUMN current_price double unsigned; ALTER TABLE reservations ADD COLUMN current_price_currency varchar(10); ALTER TABLE reservations ADD COLUMN current_price_unit varchar(100); UPDATE reservations r JOIN machines m ON r.machine_id = m.id SET r.current_price=m.reservation_price_start, r.current_price_currency=\'€\', r.current_price_unit=\'30 minutes\'','ALTER TABLE reservations DROP COLUMN current_price; ALTER TABLE reservations DROP COLUMN current_price_currency; ALTER TABLE reservations DROP COLUMN current_price_unit','rollback'),(40,'Reservationcurrentprice_20151026_120719','2015-10-26 13:53:41','ALTER TABLE reservations ADD COLUMN current_price double unsigned; ALTER TABLE reservations ADD COLUMN current_price_currency varchar(10); ALTER TABLE reservations ADD COLUMN current_price_unit varchar(100); UPDATE reservations r JOIN machines m ON r.machine_id = m.id SET r.current_price=m.reservation_price_hourly, r.current_price_currency=\'€\', r.current_price_unit=\'30 minutes\'',NULL,'update'),(41,'Globalconfig_20151117_095700','2015-11-18 10:55:43','CREATE TABLE settings (\n		id int(11) unsigned NOT NULL AUTO_INCREMENT,\n		name varchar(100) NOT NULL,\n		value_int int(11),\n		value_string text,\n		value_float double,\n		PRIMARY KEY (id)\n	)',NULL,'update'),(42,'Reservationdisabled_20151118_133537','2015-11-18 13:47:25','ALTER TABLE reservations ADD COLUMN disabled tinyint(1)',NULL,'update'),(43,'Products_20151118_141721','2015-11-18 13:47:25','CREATE TABLE products (\n		id int(11) unsigned NOT NULL AUTO_INCREMENT,\n		type varchar(100),\n		name varchar(100),\n		price double unsigned,\n		price_unit varchar(100),\n		PRIMARY KEY (id)\n	)',NULL,'update'),(44,'Purchases_20151119_115310','2015-11-23 09:42:37','CREATE TABLE purchases (\n		id int(11) unsigned NOT NULL AUTO_INCREMENT,\n		type varchar(100) NOT NULL,\n		product_id int(11) unsigned,\n		created datetime,\n		user_id int(11) unsigned NOT NULL,\n		time_start datetime DEFAULT NULL,\n		time_end datetime DEFAULT NULL,\n		quantity double NOT NULL,\n		price_per_unit double,\n		price_unit varchar(100),\n		vat double,\n		activation_running tinyint(1),\n		reservation_disabled tinyint(1),\n		machine_id int(11) unsigned,\n		PRIMARY KEY (id)\n	); \n		INSERT INTO purchases ( TYPE, product_id, created, user_id, time_start, time_end, quantity, price_per_unit, price_unit, vat, activation_running, reservation_disabled, machine_id )\n		SELECT \'activation\',\n		       NULL,\n		       time_start,\n		       user_id,\n		       time_start,\n		       time_end,\n		       time_total / 60,\n		       current_machine_price,\n		       current_machine_price_unit,\n		       vat_rate,\n		       active,\n		       NULL,\n		       machine_id\n		FROM activations\n		WHERE current_machine_price_unit = \"minute\"\n		UNION\n		SELECT \'activation\',\n		       NULL,\n		       time_start,\n		       user_id,\n		       time_start,\n		       time_end,\n		       time_total / 3600,\n		       current_machine_price,\n		       current_machine_price_unit,\n		       vat_rate,\n		       active,\n		       NULL,\n		       machine_id\n		FROM activations\n		WHERE current_machine_price_unit = \"hour\"\n	; \n		INSERT INTO purchases ( TYPE, product_id, created, user_id, time_start, time_end, quantity, price_per_unit, price_unit, vat, activation_running, reservation_disabled, machine_id )\n		SELECT \'reservation\',\n		       NULL,\n		       created,\n		       user_id,\n		       time_start,\n		       time_end,\n		       TIME_TO_SEC(TIMEDIFF(time_end, time_start)) / 1800,\n		       current_price,\n		       current_price_unit,\n		       NULL,\n		       NULL,\n		       disabled,\n		       machine_id\n		FROM reservations\n	',NULL,'update'),(45,'PurchaseCancelledFlag_20151124_153407','2015-11-25 09:33:54','ALTER TABLE purchases ADD COLUMN cancelled tinyint(1)',NULL,'update'),(46,'Tutorproduct_20151126_120015','2015-11-26 11:03:44','ALTER TABLE products ADD COLUMN user_id int(11); ALTER TABLE products ADD COLUMN machine_skills varchar(255)',NULL,'update'),(47,'Productcomments_20151126_170242','2015-11-26 16:05:16','ALTER TABLE products ADD COLUMN comments TEXT',NULL,'update'),(48,'PurchaseTimeEndActual_20151202_102231','2015-12-02 15:32:55','ALTER TABLE purchases ADD COLUMN time_end_actual datetime AFTER time_end; ALTER TABLE purchases CHANGE COLUMN activation_running running tinyint(1)',NULL,'update'),(49,'Productarchived_20151202_172646','2015-12-02 16:31:24','ALTER TABLE products ADD COLUMN archived TINYINT(1) DEFAULT 0','ALTER TABLE products DROP COLUMN archived','rollback'),(50,'Productarchived_20151202_172646','2015-12-02 16:31:33','ALTER TABLE products ADD COLUMN archived TINYINT(1) DEFAULT 0',NULL,'update'),(51,'PurchaseTimeEndPlanned_20151203_134204','2015-12-03 13:23:39','ALTER TABLE purchases CHANGE COLUMN time_end_actual time_end_planned datetime',NULL,'update'),(52,'Purchasearchived_20151203_142543','2015-12-03 13:34:34','ALTER TABLE purchases ADD COLUMN archived TINYINT(1) DEFAULT 0','ALTER TABLE purchases DROP COLUMN archived','rollback'),(53,'Purchasearchived_20151203_142543','2015-12-03 13:35:02','ALTER TABLE purchases ADD COLUMN archived TINYINT(1) DEFAULT 0',NULL,'update'),(54,'Purchasecomments_20151203_160933','2015-12-03 15:11:46','ALTER TABLE purchases ADD COLUMN comments TEXT DEFAULT \'\'','ALTER TABLE purchases DROP COLUMN comments','rollback'),(55,'Purchasecomments_20151203_160933','2015-12-03 15:12:38','ALTER TABLE purchases ADD COLUMN comments TEXT DEFAULT \'\'',NULL,'update'),(56,'Tutorpurchasetimer_20151208_143302','2015-12-08 13:38:01','ALTER TABLE purchases ADD COLUMN timer_start_time DATETIME','ALTER TABLE purchases DROP COLUMN timer_start_time','rollback'),(57,'Tutorpurchasetimer_20151208_143302','2015-12-08 13:38:43','ALTER TABLE purchases ADD COLUMN timer_time_start DATETIME',NULL,'update');
/*!40000 ALTER TABLE `migrations` ENABLE KEYS */;
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
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
) ENGINE=InnoDB AUTO_INCREMENT=330 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `permission`
--

LOCK TABLES `permission` WRITE;
/*!40000 ALTER TABLE `permission` DISABLE KEYS */;
INSERT INTO `permission` VALUES (183,2,2),(184,2,3),(188,2,13),(327,1,1),(328,1,2),(329,1,3);
/*!40000 ALTER TABLE `permission` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `products`
--

DROP TABLE IF EXISTS `products`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `products` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `type` varchar(100) DEFAULT NULL,
  `name` varchar(100) DEFAULT NULL,
  `price` double unsigned DEFAULT NULL,
  `price_unit` varchar(100) DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  `machine_skills` varchar(255) DEFAULT NULL,
  `comments` text,
  `archived` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `products`
--

LOCK TABLES `products` WRITE;
/*!40000 ALTER TABLE `products` DISABLE KEYS */;
INSERT INTO `products` VALUES (1,'space','Electronics Lab',50,'hour',NULL,NULL,NULL,0),(2,'','',12,'hour',2,'[2,3]',NULL,0),(3,'tutor','',13,'hour',2,'[2,3]','',1),(4,'tutor','Hermann Stufe',15,'hour',2,'[3,2,1]','Very good tutor.',0),(5,'tutor','Krisjanis Rijnieks',15,'hour',3,'[1,2]','',1),(6,'tutor','Kalvis Maitāns',16,'hour',6,'[1,2]','',1),(7,'tutor','Hermann Stufe',17,'hour',2,'[2,1]','',1),(8,'tutor','Hermann Stufe',16,'hour',2,'[2,3]','',1),(9,'tutor','Mike Smike',15,'hour',4,'[2,1]','',1),(10,'tutor','Karlik Memenjev',60,'hour',8,'[2]','',0),(11,'co-working','Exclusive space',60,'month',0,'','',0),(12,'tutor','',0,'',0,'','',1),(13,'tutor','Karlik Memenjev',50,'hour',8,'[1]','Yes. Let\'s do it!',0),(14,'tutor','',0,'',0,'','',1),(15,'space','Vinyl Lab',20,'hour',0,'','',0),(16,'space','Space Shit',50,'hour',0,'','',0),(17,'co-working','Special Cowok',1000,'month',0,'','',0),(18,'tutor','',0,'',0,'','',1),(19,'co-working','Sukatable',30,'month',0,'','',0),(20,'co-working','Sylwesters Table',100,'month',0,'','',0),(21,'co-working','My Table',30,'month',0,'','',0),(22,'co-working','No Table at All',50,'month',0,'','',0),(23,'space','Total Space',30,'hour',0,'','',0),(24,'space','Mother Space',50,'hour',0,'','',0),(25,'space','Zone Space',50,'hour',0,'','',0),(26,'co-working','My Product',34,'month',0,'','',0),(27,'tutor','',0,'',0,'','',1),(28,'tutor','',0,'',0,'','',1);
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
  `type` varchar(100) NOT NULL,
  `product_id` int(11) unsigned DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `user_id` int(11) unsigned NOT NULL,
  `time_start` datetime DEFAULT NULL,
  `time_end` datetime DEFAULT NULL,
  `time_end_planned` datetime DEFAULT NULL,
  `quantity` double NOT NULL,
  `price_per_unit` double DEFAULT NULL,
  `price_unit` varchar(100) DEFAULT NULL,
  `vat` double DEFAULT NULL,
  `running` tinyint(1) DEFAULT NULL,
  `reservation_disabled` tinyint(1) DEFAULT NULL,
  `machine_id` int(11) unsigned DEFAULT NULL,
  `cancelled` tinyint(1) DEFAULT NULL,
  `archived` tinyint(1) DEFAULT '0',
  `comments` text,
  `timer_time_start` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=174 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `purchases`
--

LOCK TABLES `purchases` WRITE;
/*!40000 ALTER TABLE `purchases` DISABLE KEYS */;
INSERT INTO `purchases` VALUES (1,'activation',NULL,'2015-06-25 17:40:22',7,'2015-06-25 17:40:22','2015-06-25 17:40:27',NULL,0.1,1,'minute',0,0,NULL,3,NULL,0,NULL,NULL),(2,'activation',NULL,'2015-06-25 17:40:30',7,'2015-06-25 17:40:30','2015-06-25 17:40:32',NULL,0.05,1,'minute',0,0,NULL,3,NULL,0,NULL,NULL),(3,'activation',NULL,'2015-07-15 15:50:40',7,'2015-07-15 15:50:40','2015-07-15 15:50:44',NULL,0.0667,1,'minute',0,0,NULL,3,NULL,0,NULL,NULL),(4,'activation',NULL,'2015-07-15 18:07:34',7,'2015-07-15 18:07:34','2015-07-15 18:07:45',NULL,0.1833,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(5,'activation',NULL,'2015-07-15 18:07:48',7,'2015-07-15 18:07:48','2015-07-15 18:07:56',NULL,0.1333,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(6,'activation',NULL,'2015-07-15 18:08:05',7,'2015-07-15 18:08:05','2015-07-15 18:08:21',NULL,0.2833,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(7,'activation',NULL,'2015-07-28 11:42:33',7,'2015-07-28 11:42:33','2015-07-28 11:42:39',NULL,0.1,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(8,'activation',NULL,'2015-07-28 11:42:42',7,'2015-07-28 11:42:42','2015-07-28 11:42:46',NULL,0.0667,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(9,'activation',NULL,'2015-07-28 11:42:48',7,'2015-07-28 11:42:48','2015-07-28 11:42:50',NULL,0.05,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(10,'activation',NULL,'2015-07-28 11:42:54',7,'2015-07-28 11:42:54','2015-07-28 11:43:04',NULL,0.1667,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(11,'activation',NULL,'2015-07-28 11:47:19',7,'2015-07-28 11:47:19','2015-07-28 11:47:23',NULL,0.0667,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(12,'activation',NULL,'2015-07-28 11:47:24',7,'2015-07-28 11:47:24','2015-07-28 11:47:26',NULL,0.05,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(13,'activation',NULL,'2015-07-29 12:06:34',1,'2015-07-29 12:06:34','2015-07-29 12:06:42',NULL,0.1333,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(14,'activation',NULL,'2015-07-29 12:06:58',2,'2015-07-29 12:06:58','2015-07-29 12:07:07',NULL,0.1667,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(15,'activation',NULL,'2015-08-18 16:40:45',1,'2015-08-18 16:40:45','2015-08-18 16:40:53',NULL,0.1333,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(16,'activation',NULL,'2015-08-18 17:15:59',1,'2015-08-18 17:15:59','2015-08-18 17:30:24',NULL,14.4167,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(17,'activation',NULL,'2015-08-18 18:24:38',1,'2015-08-18 18:24:38','2015-08-18 18:25:00',NULL,0.3667,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(18,'activation',NULL,'2015-09-17 11:37:20',2,'2015-09-17 11:37:20','2015-09-17 11:37:32',NULL,0.2167,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(19,'activation',NULL,'2015-09-17 11:37:34',2,'2015-09-17 11:37:34','2015-09-17 11:37:38',NULL,0.0667,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(20,'activation',NULL,'2015-09-17 11:49:02',2,'2015-09-17 11:49:02','2015-09-17 11:49:06',NULL,0.0667,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(21,'activation',NULL,'2015-09-17 11:52:32',2,'2015-09-17 11:52:32','2015-09-17 11:54:18',NULL,1.7667,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(22,'activation',NULL,'2015-09-17 11:54:19',2,'2015-09-17 11:54:19','2015-09-17 11:54:23',NULL,0.0667,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(23,'activation',NULL,'2015-09-17 11:54:24',2,'2015-09-17 11:54:24','2015-09-17 11:54:28',NULL,0.0833,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(24,'activation',NULL,'2015-09-17 11:54:41',2,'2015-09-17 11:54:41','2015-09-17 11:54:43',NULL,0.0333,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(25,'activation',NULL,'2015-09-17 11:54:45',2,'2015-09-17 11:54:45','2015-09-17 11:54:49',NULL,0.0667,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(26,'activation',NULL,'2015-09-17 11:55:17',2,'2015-09-17 11:55:17','2015-09-17 11:59:24',NULL,4.1333,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(27,'activation',NULL,'2015-09-17 11:59:25',2,'2015-09-17 11:59:25','2015-09-17 11:59:39',NULL,0.2333,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(28,'activation',NULL,'2015-09-17 11:59:43',2,'2015-09-17 11:59:43','2015-09-17 11:59:46',NULL,0.05,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(29,'activation',NULL,'2015-09-17 12:03:41',2,'2015-09-17 12:03:41','2015-09-17 12:03:43',NULL,0.05,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(30,'activation',NULL,'2015-09-17 12:03:44',2,'2015-09-17 12:03:44','2015-09-17 12:03:46',NULL,0.05,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(31,'activation',NULL,'2015-09-17 12:06:00',2,'2015-09-17 12:06:00','2015-09-17 12:06:02',NULL,0.0333,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(32,'activation',NULL,'2015-09-17 12:06:13',2,'2015-09-17 12:06:13','2015-09-17 17:23:02',NULL,316.8333,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(33,'activation',NULL,'2015-09-17 12:37:10',2,'2015-09-17 12:37:10','2015-09-17 12:37:12',NULL,0.0333,1,'minute',0,0,NULL,3,NULL,0,NULL,NULL),(34,'activation',NULL,'2015-09-17 17:23:06',2,'2015-09-17 17:23:06','2015-09-17 18:13:24',NULL,50.3,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(35,'activation',NULL,'2015-09-17 18:39:45',2,'2015-09-17 18:39:45','2015-09-17 18:40:15',NULL,0.5,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(36,'activation',NULL,'2015-09-17 18:40:17',2,'2015-09-17 18:40:17','2015-09-17 19:03:52',NULL,23.6,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(37,'activation',NULL,'2015-09-17 19:03:56',2,'2015-09-17 19:03:56','2015-09-17 19:04:00',NULL,0.0667,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(38,'activation',NULL,'2015-09-17 19:04:04',1,'2015-09-17 19:04:04','2015-09-17 19:04:10',NULL,0.1167,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(39,'activation',NULL,'2015-09-17 19:11:23',1,'2015-09-17 19:11:23','2015-09-17 19:11:27',NULL,0.0833,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(40,'activation',NULL,'2015-09-17 19:23:01',1,'2015-09-17 19:23:01','2015-09-17 19:23:03',NULL,0.0333,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(41,'activation',NULL,'2015-09-17 19:23:14',2,'2015-09-17 19:23:14','2015-10-13 12:43:10',NULL,37039.9333,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(42,'activation',NULL,'2015-10-14 13:42:46',2,'2015-10-14 13:42:46','2015-10-14 13:42:49',NULL,0.0333,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(43,'activation',NULL,'2015-10-14 13:52:48',2,'2015-10-14 13:52:48','2015-10-14 13:52:56',NULL,0.1333,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(44,'activation',NULL,'2015-10-15 08:10:56',1,'2015-10-15 08:10:56','2015-10-15 08:11:27',NULL,0.5,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(45,'activation',NULL,'2015-10-15 08:15:46',2,'2015-10-15 08:15:46','2015-10-15 08:15:50',NULL,0.05,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(46,'activation',NULL,'2015-10-15 12:21:06',2,'2015-10-15 12:21:06','2015-10-15 12:21:40',NULL,0.55,1,'minute',0,0,NULL,3,NULL,0,NULL,NULL),(47,'activation',NULL,'2015-10-16 08:11:09',1,'2015-10-16 08:11:09','2015-10-16 08:11:12',NULL,0.05,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(48,'activation',NULL,'2015-10-18 12:05:46',1,'2015-10-18 12:05:46','2015-10-18 12:05:51',NULL,0.0833,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(49,'activation',NULL,'2015-10-22 15:24:05',1,'2015-10-22 15:24:05','2015-10-22 15:24:10',NULL,0.0667,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(50,'activation',NULL,'2015-10-23 07:24:15',1,'2015-10-23 07:24:15','2015-10-23 07:26:52',NULL,2.6167,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(51,'activation',NULL,'2015-10-23 07:26:54',1,'2015-10-23 07:26:54','2015-10-23 07:27:01',NULL,0.1167,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(52,'activation',NULL,'2015-10-23 07:33:54',1,'2015-10-23 07:33:54','2015-10-23 07:34:15',NULL,0.3333,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(53,'activation',NULL,'2015-10-23 07:52:01',1,'2015-10-23 07:52:01','2015-10-23 07:52:03',NULL,0.0333,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(54,'activation',NULL,'2015-10-23 08:55:54',1,'2015-10-23 08:55:54','2015-10-23 08:55:59',NULL,0.0833,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(55,'activation',NULL,'2015-10-27 09:32:11',2,'2015-10-27 09:32:11','2015-10-27 09:32:20',NULL,0.15,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(56,'activation',NULL,'2015-10-28 09:43:28',1,'2015-10-28 09:43:28','2015-10-28 09:43:32',NULL,0.05,1,'minute',0,0,NULL,3,NULL,0,NULL,NULL),(57,'activation',NULL,'2015-11-19 15:07:32',1,'2015-11-19 15:07:32','2015-11-20 16:30:09',NULL,1522.6167,0.1,'minute',0,0,NULL,1,NULL,0,NULL,NULL),(58,'activation',NULL,'2015-06-25 17:40:13',7,'2015-06-25 17:40:13','2015-06-25 17:40:16',NULL,0.0008,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(59,'activation',NULL,'2015-07-15 15:51:04',7,'2015-07-15 15:51:04','2015-07-15 15:51:10',NULL,0.0017,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(60,'activation',NULL,'2015-07-15 15:51:18',7,'2015-07-15 15:51:18','2015-07-15 15:51:24',NULL,0.0019,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(61,'activation',NULL,'2015-07-15 15:51:32',7,'2015-07-15 15:51:32','2015-07-15 15:51:37',NULL,0.0014,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(62,'activation',NULL,'2015-07-15 18:04:57',7,'2015-07-15 18:04:57','2015-07-15 18:05:17',NULL,0.0056,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(63,'activation',NULL,'2015-07-28 11:42:17',7,'2015-07-28 11:42:17','2015-07-28 11:42:28',NULL,0.0033,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(64,'activation',NULL,'2015-07-28 11:47:35',7,'2015-07-28 11:47:35','2015-07-28 11:47:38',NULL,0.0011,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(65,'activation',NULL,'2015-07-29 12:06:18',1,'2015-07-29 12:06:18','2015-07-29 12:06:32',NULL,0.0042,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(66,'activation',NULL,'2015-08-18 16:40:58',1,'2015-08-18 16:40:58','2015-08-18 16:42:53',NULL,0.0319,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(67,'activation',NULL,'2015-08-18 16:42:54',1,'2015-08-18 16:42:54','2015-08-18 16:43:10',NULL,0.0044,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(68,'activation',NULL,'2015-08-18 16:43:11',1,'2015-08-18 16:43:11','2015-08-18 16:44:56',NULL,0.0292,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(69,'activation',NULL,'2015-09-17 12:10:31',1,'2015-09-17 12:10:31','2015-09-17 12:14:07',NULL,0.0603,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(70,'activation',NULL,'2015-09-17 12:15:00',1,'2015-09-17 12:15:00','2015-09-17 12:15:03',NULL,0.0008,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(71,'activation',NULL,'2015-09-17 12:15:07',1,'2015-09-17 12:15:07','2015-09-17 12:15:09',NULL,0.0008,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(72,'activation',NULL,'2015-09-17 12:15:28',1,'2015-09-17 12:15:28','2015-09-17 12:16:31',NULL,0.0175,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(73,'activation',NULL,'2015-09-17 12:16:41',1,'2015-09-17 12:16:41','2015-09-17 13:24:09',NULL,1.1247,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(74,'activation',NULL,'2015-09-17 13:24:18',2,'2015-09-17 13:24:18','2015-09-17 13:24:21',NULL,0.0008,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(75,'activation',NULL,'2015-09-17 13:24:28',1,'2015-09-17 13:24:28','2015-09-17 18:11:29',NULL,4.7836,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(76,'activation',NULL,'2015-09-17 18:13:16',1,'2015-09-17 18:13:16','2015-09-17 18:13:18',NULL,0.0008,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(77,'activation',NULL,'2015-09-17 18:39:25',1,'2015-09-17 18:39:25','2015-09-17 18:40:02',NULL,0.0103,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(78,'activation',NULL,'2015-09-17 18:40:22',1,'2015-09-17 18:40:22','2015-09-17 19:04:18',NULL,0.3992,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(79,'activation',NULL,'2015-09-17 19:09:40',1,'2015-09-17 19:09:40','2015-09-17 19:23:26',NULL,0.2297,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(80,'activation',NULL,'2015-09-17 19:23:31',1,'2015-09-17 19:23:31','2015-09-18 14:10:43',NULL,18.7869,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(81,'activation',NULL,'2015-09-18 14:11:28',2,'2015-09-18 14:11:28','2015-09-18 14:11:35',NULL,0.0022,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(82,'activation',NULL,'2015-09-18 14:11:36',2,'2015-09-18 14:11:36','2015-09-18 14:11:41',NULL,0.0014,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(83,'activation',NULL,'2015-10-11 11:09:41',1,'2015-10-11 11:09:41','2015-10-11 12:09:53',NULL,1.0031,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(84,'activation',NULL,'2015-10-11 12:09:54',1,'2015-10-11 12:09:54','2015-10-13 09:55:26',NULL,45.7586,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(85,'activation',NULL,'2015-10-14 13:22:37',1,'2015-10-14 13:22:37','2015-10-14 13:22:42',NULL,0.0014,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(86,'activation',NULL,'2015-10-14 13:49:04',1,'2015-10-14 13:49:04','2015-10-14 13:49:16',NULL,0.0031,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(87,'activation',NULL,'2015-10-14 13:50:45',1,'2015-10-14 13:50:45','2015-10-14 13:50:56',NULL,0.0028,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(88,'activation',NULL,'2015-10-14 13:51:01',1,'2015-10-14 13:51:01','2015-10-14 13:51:08',NULL,0.0017,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(89,'activation',NULL,'2015-10-14 13:51:21',1,'2015-10-14 13:51:21','2015-10-14 13:51:31',NULL,0.0028,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(90,'activation',NULL,'2015-10-14 13:53:33',1,'2015-10-14 13:53:33','2015-10-14 13:53:34',NULL,0.0003,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(91,'activation',NULL,'2015-10-14 13:56:19',2,'2015-10-14 13:56:19','2015-10-14 13:56:21',NULL,0.0003,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(92,'activation',NULL,'2015-10-14 14:13:03',1,'2015-10-14 14:13:03','2015-10-14 14:13:05',NULL,0.0003,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(93,'activation',NULL,'2015-10-15 08:12:45',1,'2015-10-15 08:12:45','2015-10-15 08:15:55',NULL,0.0525,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(94,'activation',NULL,'2015-10-15 08:15:57',2,'2015-10-15 08:15:57','2015-10-15 08:16:02',NULL,0.0011,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(95,'activation',NULL,'2015-10-15 14:53:37',1,'2015-10-15 14:53:37','2015-10-16 08:07:43',NULL,17.2347,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(96,'activation',NULL,'2015-10-16 08:07:49',2,'2015-10-16 08:07:49','2015-10-16 08:07:52',NULL,0.0008,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(97,'activation',NULL,'2015-10-23 08:56:20',1,'2015-10-23 08:56:20','2015-10-23 08:56:23',NULL,0.0008,16,'hour',0,0,NULL,2,NULL,0,NULL,NULL),(128,'reservation',0,'2015-10-13 09:55:53',2,'2015-10-14 08:00:00','2015-10-14 10:00:00',NULL,4,5,'30 minutes',0,0,0,2,0,0,'',NULL),(129,'reservation',NULL,'2015-10-13 09:00:00',2,'2015-10-13 08:00:00','2015-10-13 14:00:00',NULL,12,5,'30 minutes',NULL,NULL,NULL,2,NULL,0,NULL,NULL),(130,'reservation',NULL,'2015-10-14 12:53:21',2,'2015-10-16 09:00:00','2015-10-16 09:30:00',NULL,1,5,'30 minutes',NULL,NULL,NULL,2,NULL,0,NULL,NULL),(131,'reservation',NULL,'2015-10-14 13:11:01',2,'2015-10-15 08:00:00','2015-10-15 08:30:00',NULL,1,2.4,'30 minutes',NULL,NULL,NULL,1,NULL,0,NULL,NULL),(132,'reservation',NULL,'2015-10-14 13:12:51',1,'2015-10-22 14:30:00','2015-10-22 16:30:00',NULL,4,2.4,'30 minutes',NULL,NULL,NULL,1,NULL,0,NULL,NULL),(133,'reservation',NULL,'2015-10-14 13:13:28',2,'2015-10-16 11:00:00','2015-10-16 13:00:00',NULL,4,5,'30 minutes',NULL,NULL,NULL,2,NULL,0,NULL,NULL),(134,'reservation',NULL,'2015-10-14 13:15:50',1,'2015-10-23 14:00:00','2015-10-23 14:30:00',NULL,1,2.4,'30 minutes',NULL,NULL,NULL,1,NULL,0,NULL,NULL),(135,'reservation',NULL,'2015-10-14 13:16:52',2,'2015-10-14 10:30:00','2015-10-14 19:00:00',NULL,17,2.4,'30 minutes',NULL,NULL,NULL,1,NULL,0,NULL,NULL),(136,'reservation',NULL,'2015-10-15 08:18:43',2,'2015-10-23 09:30:00','2015-10-23 11:30:00',NULL,4,2.4,'30 minutes',NULL,NULL,NULL,1,NULL,0,NULL,NULL),(137,'reservation',NULL,'2015-10-27 16:01:15',2,'2015-10-28 10:00:00','2015-10-28 15:00:00',NULL,10,5,'30 minutes',NULL,NULL,NULL,3,NULL,0,NULL,NULL),(138,'reservation',NULL,'2015-10-24 16:01:15',2,'2015-10-27 10:00:00','2015-10-27 13:00:00',NULL,6,5,'30 minutes',NULL,NULL,NULL,3,NULL,0,NULL,NULL),(139,'reservation',NULL,'2015-10-27 16:01:15',2,'2015-10-28 16:00:00','2015-10-28 17:00:00',NULL,2,5,'30 minutes',NULL,NULL,NULL,3,NULL,0,NULL,NULL),(140,'reservation',NULL,'2015-10-27 16:01:15',1,'2015-10-29 16:00:00','2015-10-29 17:00:00',NULL,2,5,'30 minutes',NULL,NULL,NULL,3,NULL,0,NULL,NULL),(141,'reservation',NULL,'2015-10-27 16:01:15',1,'2015-10-30 16:00:00','2015-10-30 17:00:00',NULL,2,5,'30 minutes',NULL,NULL,NULL,3,NULL,0,NULL,NULL),(145,'reservation',0,'2015-11-25 09:38:08',1,'2015-11-26 09:30:00','2015-11-26 11:30:00',NULL,0,5,'30 minutes',0,0,0,3,0,0,NULL,NULL),(146,'reservation',0,'2015-11-24 09:38:08',1,'2015-11-25 09:30:00','2015-11-25 11:30:00',NULL,0,5,'30 minutes',0,0,0,3,0,0,NULL,NULL),(147,'space',15,'2015-11-25 10:14:10',7,'2015-11-25 10:14:00','2015-11-25 10:14:00',NULL,0,20,'hour',0,0,0,0,0,0,'',NULL),(148,'space',15,'2015-11-25 13:46:54',5,'2015-11-25 13:46:00','2015-11-25 13:46:00',NULL,0,20,'hour',0,0,0,0,0,0,'',NULL),(149,'space',1,'2015-11-26 09:08:50',2,'2015-11-26 09:08:00','2015-11-26 09:08:00',NULL,0,50,'hour',0,0,0,0,0,0,'',NULL),(150,'activation',0,NULL,1,'2015-11-27 17:02:44','2015-11-27 17:02:50',NULL,0.0016893587469444444,16,'hour',0,0,0,2,0,0,NULL,NULL),(151,'co-working',21,'2015-11-27 17:04:00',9,'2015-11-27 17:04:00','2015-11-27 17:04:00',NULL,0,30,'month',0,0,0,0,0,0,'',NULL),(152,'tutor',0,'2015-12-02 16:11:55',0,'2015-12-02 16:11:55','2015-12-02 16:11:55',NULL,0,0,'hour',0,0,0,0,0,1,NULL,NULL),(153,'tutor',0,'2015-12-03 11:15:27',0,'2015-12-03 11:15:27','2015-12-03 11:15:27',NULL,0,0,'hour',0,0,0,0,0,1,NULL,NULL),(154,'tutor',0,'2015-12-03 13:05:00',0,'2015-12-03 13:05:00','2015-12-03 13:05:00',NULL,0,0,'hour',0,0,0,0,0,1,NULL,NULL),(155,'tutor',0,'2015-12-03 13:15:22',0,'2015-12-03 13:15:22','2015-12-03 13:15:22',NULL,0,0,'hour',0,0,0,0,0,1,NULL,NULL),(156,'tutor',4,'2015-12-03 13:24:11',6,'2015-12-11 08:57:00','2015-12-12 09:12:00','2015-12-03 14:24:00',0.001388888888888889,15,'hour',0,0,0,0,0,0,'Some comments. Sure.','2015-12-08 14:19:36'),(157,'tutor',10,'2015-12-03 14:22:58',2,'2015-12-03 14:22:00',NULL,'2015-12-03 13:22:00',0,0,'',0,0,0,0,0,1,NULL,NULL),(158,'tutor',13,'2015-12-03 14:26:19',8,'2015-12-03 14:26:00',NULL,'2015-12-03 15:26:00',0,0,'',0,0,0,0,0,1,NULL,NULL),(159,'tutor',4,'2015-12-03 14:58:57',8,'2015-12-03 14:58:00',NULL,'2015-12-03 15:58:00',1,15,'hour',0,0,0,0,0,1,NULL,NULL),(160,'tutor',4,'2015-12-03 14:59:14',4,'2015-12-08 13:00:00',NULL,'2015-12-08 14:00:00',1,15,'hour',0,0,0,0,0,1,'Arbeit machen. Freizeit geniessen.',NULL),(161,'space',1,'2015-12-03 15:37:16',3,'2015-12-03 15:37:00','2015-12-03 16:37:00',NULL,1,50,'hour',0,0,0,0,0,0,'',NULL),(162,'co-working',17,'2015-12-03 16:17:17',3,'2015-12-03 16:17:00','2015-12-03 16:17:00',NULL,0,1000,'month',0,0,0,0,0,0,'',NULL),(163,'co-working',17,'2015-12-03 16:18:46',2,'2015-12-03 16:18:00','2015-12-03 16:18:00',NULL,0,1000,'month',0,0,0,0,0,0,'',NULL),(164,'co-working',17,'2015-12-04 18:29:18',4,'2015-12-04 18:29:00','2015-12-04 18:29:00',NULL,0,1000,'month',0,0,0,0,0,0,'',NULL),(165,'co-working',11,'2015-12-04 19:33:32',6,'2015-12-04 19:33:00','2016-01-14 19:33:00',NULL,3,60,'month',0,0,0,0,0,0,'',NULL),(166,'activation',0,NULL,1,'2015-12-07 08:56:58',NULL,NULL,0,0.1,'minute',0,1,0,1,0,0,'',NULL),(167,'tutor',4,'2015-12-07 08:57:21',6,'2015-12-07 08:57:21','2015-12-07 08:57:28','2015-12-07 08:57:21',0.0018868835122222223,15,'hour',0,0,0,0,0,0,'Some comments. Sure.',NULL),(168,'activation',0,NULL,2,'2015-12-07 08:58:56',NULL,NULL,0,16,'hour',0,1,0,2,0,0,'',NULL),(169,'tutor',0,'2015-12-08 13:03:28',0,'2015-12-08 13:03:28',NULL,'2015-12-08 13:03:28',0,0,'hour',0,0,0,0,0,1,'',NULL),(170,'tutor',4,'2015-12-08 13:26:00',6,'2015-12-08 13:26:00','2015-12-08 13:26:08','2015-12-07 08:57:21',0.011099258611666666,15,'hour',0,0,0,0,0,0,'Some comments. Sure.','2015-12-08 15:30:56'),(171,'tutor',4,'2015-12-08 13:26:09',6,'2015-12-08 13:26:09','2015-12-08 14:26:09','2015-12-07 08:57:21',0,15,'hour',0,1,0,0,0,1,'Some comments. Sure.',NULL),(172,'tutor',4,'2015-12-10 14:21:49',3,'2015-12-10 14:00:00','2015-12-10 15:21:00','2015-12-10 15:00:00',1.3856587187727778,15,'hour',0,0,0,0,0,0,'','2015-12-10 17:40:45'),(173,'tutor',4,'2015-12-10 14:29:23',7,'2015-12-10 14:20:00','2015-12-10 16:30:00',NULL,0.016647024204444444,15,'hour',0,0,0,0,0,0,'Let\'s have some fun!','2015-12-10 17:39:42');
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
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `reservation_rules`
--

LOCK TABLES `reservation_rules` WRITE;
/*!40000 ALTER TABLE `reservation_rules` DISABLE KEYS */;
INSERT INTO `reservation_rules` VALUES (5,'Laydrop Printer Rule',1,1,0,'2015-10-21','2015-10-23','08:00','20:00','',0,0,1,1,1,0,0,'2015-10-22 13:36:57'),(7,'Laydrop Extended Rule',1,1,0,'2015-10-24','2015-11-25','10:00','18:00','',1,1,1,1,1,1,1,'2015-10-22 13:47:27'),(8,'Laydrop Unavail',1,0,1,'2015-11-02','2015-11-07','12:00','14:00','',1,1,1,1,1,1,0,'2015-10-22 13:49:43'),(9,'Zing Avail',3,1,0,'2015-10-20','2015-12-30','09:00','19:00','',1,1,1,1,1,1,1,'2015-10-27 15:59:13');
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
) ENGINE=InnoDB AUTO_INCREMENT=43 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `reservations`
--

LOCK TABLES `reservations` WRITE;
/*!40000 ALTER TABLE `reservations` DISABLE KEYS */;
INSERT INTO `reservations` VALUES (2,2,1,'2015-10-14 08:00:00','2015-10-14 10:00:00','2015-10-13 09:55:53',5,'€','30 minutes',NULL),(3,2,2,'2015-10-13 08:00:00','2015-10-13 14:00:00','2015-10-13 09:00:00',5,'€','30 minutes',NULL),(23,2,2,'2015-10-16 09:00:00','2015-10-16 09:30:00','2015-10-14 12:53:21',5,'€','30 minutes',NULL),(26,1,2,'2015-10-15 08:00:00','2015-10-15 08:30:00','2015-10-14 13:11:01',2.4,'€','30 minutes',NULL),(28,1,1,'2015-10-22 14:30:00','2015-10-22 16:30:00','2015-10-14 13:12:51',2.4,'€','30 minutes',NULL),(30,2,2,'2015-10-16 11:00:00','2015-10-16 13:00:00','2015-10-14 13:13:28',5,'€','30 minutes',NULL),(31,1,1,'2015-10-23 14:00:00','2015-10-23 14:30:00','2015-10-14 13:15:50',2.4,'€','30 minutes',NULL),(33,1,2,'2015-10-14 10:30:00','2015-10-14 19:00:00','2015-10-14 13:16:52',2.4,'€','30 minutes',NULL),(34,1,2,'2015-10-23 09:30:00','2015-10-23 11:30:00','2015-10-15 08:18:43',2.4,'€','30 minutes',NULL),(37,3,2,'2015-10-28 10:00:00','2015-10-28 15:00:00','2015-10-27 16:01:15',5,'€','30 minutes',NULL),(38,3,2,'2015-10-27 10:00:00','2015-10-27 13:00:00','2015-10-24 16:01:15',5,'€','30 minutes',NULL),(39,3,2,'2015-10-28 16:00:00','2015-10-28 17:00:00','2015-10-27 16:01:15',5,'€','30 minutes',NULL),(40,3,1,'2015-10-29 16:00:00','2015-10-29 17:00:00','2015-10-27 16:01:15',5,'€','30 minutes',NULL),(42,3,1,'2015-10-30 16:00:00','2015-10-30 17:00:00','2015-10-27 16:01:15',5,'€','30 minutes',NULL);
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
  `name` varchar(100) NOT NULL,
  `value_int` int(11) DEFAULT NULL,
  `value_string` text,
  `value_float` double DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `settings`
--

LOCK TABLES `settings` WRITE;
/*!40000 ALTER TABLE `settings` DISABLE KEYS */;
INSERT INTO `settings` VALUES (1,'Currency',NULL,'€',NULL),(2,'VAT',NULL,NULL,19);
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
  `user_role` varchar(100) NOT NULL DEFAULT 'member',
  `created` datetime DEFAULT NULL,
  `comments` text,
  `phone` varchar(50) DEFAULT NULL,
  `zip_code` varchar(100) DEFAULT NULL,
  `city` varchar(100) DEFAULT NULL,
  `country_code` varchar(2) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES (1,'Michael','Hackenberg','user','user@example.com','Mein Strit 34','0',101,0,'','',0,'',NULL,'','+898 898989','LV-1010','Belingo','AF'),(2,'Hermann','Stufe','admin','admin@example.com','Keine str. 45','0',100,0,'','',0,'admin',NULL,'','+371 34523452','11223s','Ghanaz','CY'),(3,'Krisjanis','Rijnieks','kris','kris@makea.org','My Billing Address','0',102,0,'','',0,'','2015-06-26 12:25:33','','','11002','Berlin','DE'),(4,'Mike','Smike','striker','kas@example.com','My Street 123','',95,0,'My Company','',0,'','2015-07-28 06:25:21','','','','',''),(5,'Maik','Schmertz','schmertz','ktek@example.com','Super Street 123 asd','',96,0,'MyComp','',0,'','2015-07-28 08:43:03','','','','',''),(6,'Kalvis','Maitāns','kalvis','my@oh.god','My Street 999','',0,0,'GodChamber','',0,'','2015-07-28 13:38:05','','+371 234234777',NULL,NULL,NULL),(7,'Eric','Kurtenberf','anotheruser','kalvis@kudos.lv','Mein Address','',93,0,'','',0,'','2015-07-28 08:09:13','FasbillID007','+78 9991235','11122','Berlin','GH'),(8,'Karlik','Memenjev','karlik','karlik@makaroni.lv','Akas iela 3','',103,0,'','',0,'','2015-08-06 15:47:45','','+371 34577456','LV-1010','Riga','LV'),(9,'Mykael','Mustambika','musta','mail@medium.com','My Address 123','',104,0,'','',0,'','2015-08-08 21:49:12','','00224445','002334','Riga','DZ');
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
  `auto_extend` tinyint(1) DEFAULT NULL,
  `is_terminated` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=45 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_membership`
--

LOCK TABLES `user_membership` WRITE;
/*!40000 ALTER TABLE `user_membership` DISABLE KEYS */;
INSERT INTO `user_membership` VALUES (3,2,2,'2015-08-07 00:00:00','2016-01-07 00:00:00',1,NULL),(11,1,1,'2015-06-10 02:00:00','2016-01-10 02:00:00',1,NULL),(41,9,2,'2015-09-09 00:00:00','2016-01-09 00:00:00',1,NULL),(42,9,1,'2015-10-10 00:00:00','2016-05-10 00:00:00',1,NULL),(44,5,2,'2015-09-21 02:00:00','2016-01-21 02:00:00',1,NULL);
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

-- Dump completed on 2015-12-22 22:45:53
