-- MySQL dump 10.13  Distrib 5.6.23, for osx10.8 (x86_64)
--
-- Host: localhost    Database: fabsmith
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
INSERT INTO `auth` VALUES (19,'FFEF8A8E','f7c19341b9c14c27136b4653514f1b7d7ad16b1c2306181481956fb93b749c74c0337dcb2622d86644d83406e98d45b782c4588a3f94d25ce79547d26f7a11ae','53d2ab2f6759bf41bff8a4bbb93975fb31cd4a914a6750f3d021ba3be5ea8fd4','','2016-12-07 15:00:10');
/*!40000 ALTER TABLE `auth` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `categories`
--

DROP TABLE IF EXISTS `categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `categories` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `location_id` int(11) DEFAULT NULL,
  `short_name` varchar(20) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `archived` tinyint(1) DEFAULT '0',
  `old_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=137 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `categories`
--

LOCK TABLES `categories` WRITE;
/*!40000 ALTER TABLE `categories` DISABLE KEYS */;
INSERT INTO `categories` VALUES (9,1,'3dprinter','3D Printer',0,1),(10,1,'cnc','CNC Mill',0,2),(11,1,'heatpress','Heatpress',0,3),(12,1,'knitting','Knitting Machine',0,4),(13,1,'lasercutter','Lasercutter',0,5),(14,1,'vinylcutter','Vinylcutter',0,6),(15,1,'lasercutter-advanced','Lasercutter (Advanced)',0,7),(16,1,'cnc-pcb','CNC Mill (PCB)',0,8),(136,1,'T-Shirt-Printer','T-Shirt-Printer',0,NULL);
/*!40000 ALTER TABLE `categories` ENABLE KEYS */;
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
) ENGINE=InnoDB AUTO_INCREMENT=3333 DEFAULT CHARSET=latin1;
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
) ENGINE=InnoDB AUTO_INCREMENT=10084 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `invoices`
--

LOCK TABLES `invoices` WRITE;
/*!40000 ALTER TABLE `invoices` DISABLE KEYS */;
INSERT INTO `invoices` VALUES (10083,1,0,'',0,'',2,2017,0,0,19,'draft',0,0,0,0,0,NULL,NULL,NULL,1);
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
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `locations`
--

LOCK TABLES `locations` WRITE;
/*!40000 ALTER TABLE `locations` DISABLE KEYS */;
INSERT INTO `locations` VALUES (1,'Test Lab','','','','Berlin','','','',1,'','',1,1,1,1,0,'Europe/Berlin','location-logo-1.svg',NULL);
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
) ENGINE=InnoDB AUTO_INCREMENT=43 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `machine_maintenances`
--

LOCK TABLES `machine_maintenances` WRITE;
/*!40000 ALTER TABLE `machine_maintenances` DISABLE KEYS */;
/*!40000 ALTER TABLE `machine_maintenances` ENABLE KEYS */;
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
  `safety_guidelines` text,
  `links` text,
  `materials` text,
  `required_software` text,
  `image` varchar(255) DEFAULT NULL,
  `image_small` varchar(255) DEFAULT NULL,
  `available` tinyint(1) DEFAULT NULL,
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
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `machines`
--

LOCK TABLES `machines` WRITE;
/*!40000 ALTER TABLE `machines` DISABLE KEYS */;
INSERT INTO `machines` VALUES (2,1,'01 Vincent','1-Mara','MakerBot Replicator 2 3D printer\n\nTechnical Specifications\nPrint Technology: 	\nFused Filament Fabrication \nBuild Volume: 	9.7”W x 6.4”L x 6.1”H\n[24.6 cm x 16.3 cm x 15.5 cm] \nLayer Height Settings: 	High 100 microns [0.0039 in] \nMedium 200 microns [0.0079 in] \nLow 300 microns [0.0118 in] \nPositioning Precision: 	XY: 11 microns [0.0004 in]; Z: 2.5 microns [0.0001 in] \nFilament Diameter: 	1.75 mm [0.069 in] 0.4 mm [0.015 in] \nNozzle Diameter: 	\nSoftware Bundle: 	MakerBot MakerWareTM\nFile Types: 	STL, OBJ, Thing\nSupports:	 Windows [7+], Ubuntu [11.10+], \nMac OS X [10.6+] \nDIMENSIONS Without Spools:	 19.1 x 12.8 x 14.7 in [49 x 32 x 38 cm] \nDIMENSIONS With Spools:	19.1 x 16.5 x 14.7 in [49 x 42 x 38 cm] \nShipping Box: \n22.75 x 22.75 x 16.75 in [57.8 x 57.8 x 42.5 cm] \nWeight: 	27.8 lbs [12.6 kg] \nShipping Weight: 	39.0 lbs [79.7 kg] [All packages] \nTEMPERATURE Ambient Operation: 	15°–32° C [60°–90° F] \nStorage Temperature: 	 0°–32° C [32°–90° F]\nELECTRICAL AC Input: 	100–240V, ~4 amps, 50–60 Hz \nPower Requirements: 	24V DC @ 9.2 amps\nConnectivity: 	SD card [FAT16, max 2 GB] \nChassis:	Powder-coated steel \nBody:	PVC Panels \nBuild Platform:	356 aluminum \n XYZ Bearings: 	Wear-resistant, oil-infused bronze \nStepper Motors: 	1.8° step angle with 1/16 micro-stepping.','Watch out','https://www.google.de Google Search\n\n\nwww.makerbot.com','PLA','Makerbot Desktop','machine-2.jpg','machine-2-small.jpg',0,'',NULL,0.1,'minute','',1,'[]',0,1,NULL,NULL,0,9,'Replicator 2','','28.5 cm x 15.3 cm x 15.5 cm','','','',1,'mfi',1,0,'2017-02-08 11:50:17'),(3,1,'Epilog 6030','ZLC','Tech Specification\n\nEngraving Area	610 x 305 mm (24\" x 12\")\nMaximum Material Thickness	 197 mm (7.75\")\nLaser Wattage	40 watts\nLaser Source	State-of-the-art, digitally controlled, air-cooled CO2 laser tubes are fully modular, permanently aligned and field replaceable.\nIntelligent Memory Capacity	Multiple file storage up to 64 MB. Rolling buffer allows files of any size to be engraved.\nAir Assist	Attached air compressor to remove heat and combustible gases from the cutting surface by directing a constant stream of compressed air across the cutting surface. \nLaser Dashboard	The Laser Dashboard™ controls your Epilog Laser\'s settings from a wide range of software packages - from design programs to spreadsheet applications to CAD drawing packages. \nRed Dot Pointer	Since the laser beam is invisible, the Red Dot Pointer on Epilog\'s Zing Laser allows you to have a visual reference for locating where the laser will fire. \nRelocatable Home	When engraving items that are not easily placed at the top corner of the laser, you can set a new home position by hand with the convenient Movable Home Position feature on the Zing Laser. \nOperating Modes	Optimized raster, vector or combined modes. \nMotion Control System	High-speed micro stepper motors. \nX-Axis Bearings	Shielded Roller Bearing Assembly on a Ceramic Coated Aluminum Guide Rail. \nBelts	Advanced B-style Kevlar Belts. \nResolution	User controlled from 100 to 1000 dpi. \nSpeed and Power Control	Computer or manually control speed and power in 1% increments to 100%. Vector color mapping links speed, power and focus to any RGB color. \nPrint Interface	10 Base-T Ethernet or USB Connection. Compatible with Windows® XP/Vista/7/8. \nSize (W x D x H)	965 x 692 x 381 mm (38\" x 27.25\" x 15\")\nWeight	64 kg (140 lbs)\nElectrical Requirements	Auto-switching power supply accommodates 110 to 240 volts, 50 or 60 Hz, single phase, 15 amp AC.\nMaximum Table Weight	The Zing 16 and 24 have a static table weight of 22.7 kg (50 lbs) and a lifting table weight of 11.5 kg (25 lbs). \nVentilation System	350 - 400 CFM (595-680 m3/hr) external exhaust to the outside or internal filtration system is required. There is one output port, 4\" in diameter.','',NULL,NULL,NULL,'machine-3.jpg','machine-3-small.jpg',1,'',NULL,0.8,'minute','asd',1,'',0,0,NULL,5,0,13,'Laser cutter','','600mm x 300mm x 114mm','','','',1,'mfi',1,0,'2017-02-08 11:50:15'),(4,1,'02 Yoda','2-Messi','Craft Bot Plus printer for your medium-sized prints\n\nUse Craftware software you must.','',NULL,NULL,NULL,'machine-4.jpg','machine-4-small.jpg',0,'',NULL,0.1,'minute','Use Craftware software you must',1,'[]',0,0,NULL,NULL,0,9,'Craft Bot Plus','','252 mm x 199 mm x 150 mm','','','',1,'mfi',1,0,'2017-02-08 11:50:12'),(6,1,'07 Fabienne','I3B2','i3Berlin 3D Printer.\n\nTech Specs:\n\nTechnology	 FDM (Fused Deposition Modeling)\nCategory 	Assembled\nPrint volume	 200x200x200mm [LxWxH]\nPrinter dimensions 	400x440x380mm [LxWxH]\nPrintable materials	 Thermoplastics (PLA, ABS, Nylon, PC, PP, PET, Wood Filaments, ...)\nHotend E3D	\nFilament diameter Default 	1.75mm.\nNozzle diameter Default	0.4mm ( 0,25mm/0,6mm/1mm available )\nLayer height	 0.02mm – 0.3mm (Default 0.4 nozzle)\nMax. power consumption 	110/220V 350W\nSoftware 	Cura/Kisslicer (opensource) Simplify3D (proprietary)\nFirmware	 Arduino/Marlin\nLicence 	GPL\nSource files	 https://github.com/open3dengineering/i3_Berlin',NULL,NULL,NULL,NULL,'machine-6.jpg','machine-6-small.jpg',0,'',NULL,0.1,'minute','',1,'',0,1,NULL,NULL,0,9,'i3 Berlin','','','','','',1,'mfi',1,0,'2017-02-08 11:50:16'),(7,1,'04 Mia','MBR','MakerBot Replicator generation 5. For your medium-sized prints.\n\nTechnical Specsifications\n\nPrint Technology: Fused Filament Fabrication\nBuild Volume: 252 X x 199 Y x 150 Z mm [9.9 L x 7.8 W x 5.9 H in]\nLayer Height Settings: \nHigh 100 microns [0.0039 in]\nMedium 200 microns [0.0079 in]\nLow 300 microns [0.0118 in]\nPositioning Precision: XY: 11 microns [0.0004 in]; Z: 2.5 microns [0.0001 in]\nFilament Diameter: 1.75 mm [0.069 in]\nNozzle Diameter: 0.4 mm [0.015 in]\nSoftware Bundle: MakerBot Desktop\nFile Types: STL, OBJ, Thing\nSupports: Windows [7+], Ubuntu [11.10+], Mac OS X [10.6+]\nConnectivity: USB Stick',NULL,NULL,NULL,NULL,'machine-7.jpg','machine-7-small.jpg',0,'',NULL,0.1,'minute','',1,'',0,0,NULL,NULL,0,9,'Replicator 5 gen','','252 mm x 199 mm x 150 mm','','','',1,'mfi',1,0,'2017-02-08 11:50:19'),(8,1,'05 Pumpkin','I3B1','i3Berlin 3D Printer with dual extruder\n\nTech specs:\nTechnology	 FDM (Fused Deposition Modeling)\nCategory 	Assembled\nPrint volume	 200x200x200mm [LxWxH]\nPrinter dimensions 	400x440x380mm [LxWxH]\nPrintable materials	 Thermoplastics (PLA, ABS, Nylon, PC, PP, PET, Wood Filaments, ...)\nHotend E3D	\nFilament diameter Default 	1.75mm.\nNozzle diameter Default	0.4mm ( 0,25mm/0,6mm/1mm available )\nLayer height	 0.02mm – 0.3mm (Default 0.4 nozzle)\nMax. power consumption 	110/220V 350W\nSoftware 	Cura/Kisslicer (opensource) Simplify3D (proprietary)\nFirmware	 Arduino/Marlin\nLicence 	GPL\nSource files	 https://github.com/open3dengineering/i3_Berlin','',NULL,NULL,NULL,'machine-8.png','machine-8-small.png',0,'',NULL,0.1,'minute','',1,'',1,0,NULL,NULL,0,9,'i3 Berlin','','','','','',1,'mfi',1,0,'2017-02-08 11:50:18'),(9,1,'Electronic Station 1','ES1','Electronic Station 1 - all the soldering irons, the main multi plug basically.','',NULL,NULL,NULL,'machine-9.JPG','machine-9-small.JPG',1,'',NULL,0.1,'minute','',1,'',0,0,NULL,NULL,0,0,'','','','','','',1,'mfi',0,0,'2017-02-08 11:50:12'),(10,1,'06 Honey Bunny','I3B2','i3Berlin 3D Printer.\n\nTech Specs:\n\nTechnology	 FDM (Fused Deposition Modeling)\nCategory 	Assembled\nPrint volume	 200x200x200mm [LxWxH]\nPrinter dimensions 	400x440x380mm [LxWxH]\nPrintable materials	 Thermoplastics (PLA, ABS, Nylon, PC, PP, PET, Wood Filaments, ...)\nHotend E3D	\nFilament diameter Default 	1.75mm.\nNozzle diameter Default	0.4mm ( 0,25mm/0,6mm/1mm available )\nLayer height	 0.02mm – 0.3mm (Default 0.4 nozzle)\nMax. power consumption 	110/220V 350W\nSoftware 	Cura/Kisslicer (opensource) Simplify3D (proprietary)\nFirmware	 Arduino/Marlin\nLicence 	GPL\nSource files	 https://github.com/open3dengineering/i3_Berlin','',NULL,NULL,NULL,'machine-10.jpg','machine-10-small.jpg',0,'',NULL,0.1,'minute','',1,'',1,0,NULL,NULL,0,9,'i3 Berlin','','','','','',1,'mfi',1,0,'2017-02-08 11:50:03');
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
) ENGINE=InnoDB AUTO_INCREMENT=10058 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `membership`
--

LOCK TABLES `membership` WRITE;
/*!40000 ALTER TABLE `membership` DISABLE KEYS */;
INSERT INTO `membership` VALUES (3,1,'Staff Membership','FLS',12,0,100,'[2,97,4,12,7,8,10,6,13,40,74,98,29,11,77,43,9,48,23,3,28,15,20,34,44,18,31,17,14,37,42,35,38,32,33,16,39,24,21,26,19,25,27,41]',1,1,0,'[15,9,10,0,13,11,16,14,12]');
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
) ENGINE=InnoDB AUTO_INCREMENT=96 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `migrations`
--

LOCK TABLES `migrations` WRITE;
/*!40000 ALTER TABLE `migrations` DISABLE KEYS */;
INSERT INTO `migrations` VALUES (94,'MachineMaterialsSoftware_20170210_123506','2017-02-10 11:37:07','ALTER TABLE machines ADD COLUMN materials text AFTER links; ALTER TABLE machines ADD COLUMN required_software text AFTER materials','ALTER TABLE machines DROP COLUMN materials; ALTER TABLE machines DROP COLUMN required_software','rollback'),(95,'MachineMaterialsSoftware_20170210_123506','2017-02-10 11:37:38','ALTER TABLE machines ADD COLUMN materials text AFTER links; ALTER TABLE machines ADD COLUMN required_software text AFTER materials',NULL,'update');
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
) ENGINE=InnoDB AUTO_INCREMENT=174 DEFAULT CHARSET=utf8;
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
) ENGINE=InnoDB AUTO_INCREMENT=32322 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `permission`
--

LOCK TABLES `permission` WRITE;
/*!40000 ALTER TABLE `permission` DISABLE KEYS */;
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
) ENGINE=InnoDB AUTO_INCREMENT=41 DEFAULT CHARSET=utf8;
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
) ENGINE=InnoDB AUTO_INCREMENT=39243 DEFAULT CHARSET=utf8;
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
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8;
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
) ENGINE=InnoDB AUTO_INCREMENT=116 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `settings`
--

LOCK TABLES `settings` WRITE;
/*!40000 ALTER TABLE `settings` DISABLE KEYS */;
INSERT INTO `settings` VALUES (109,1,'TermsUrl',NULL,'https://fablab.berlin/de/content/18-agb-fab-lab',NULL),(110,1,'Currency',NULL,'€',NULL),(111,1,'FastbillTemplateId',946226,NULL,NULL),(112,1,'VAT',NULL,NULL,19),(113,1,'ReservationNotificationEmail',NULL,'',NULL),(114,1,'MailchimpApiKey',NULL,'8ce59590d15a8adbcb437ed46369e4f9-us6',NULL),(115,1,'MailchimpListId',NULL,'7a8cb19a0d',NULL);
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
) ENGINE=InnoDB AUTO_INCREMENT=11252 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES (19,'Test','User','testuser','testuser@example.com','There is one','0',696,0,'','',0,1,'2015-06-04 06:34:51','Fastbill Kd-Nr. 696','123456','','','DE',NULL,1,0,0,NULL,NULL);
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
) ENGINE=InnoDB AUTO_INCREMENT=1393 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_locations`
--

LOCK TABLES `user_locations` WRITE;
/*!40000 ALTER TABLE `user_locations` DISABLE KEYS */;
INSERT INTO `user_locations` VALUES (8,1,19,'admin',0);
/*!40000 ALTER TABLE `user_locations` ENABLE KEYS */;
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
) ENGINE=InnoDB AUTO_INCREMENT=471 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_memberships`
--

LOCK TABLES `user_memberships` WRITE;
/*!40000 ALTER TABLE `user_memberships` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_memberships` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2017-02-16 14:45:57
