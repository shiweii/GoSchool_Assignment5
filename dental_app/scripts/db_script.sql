CREATE DATABASE  IF NOT EXISTS `dental_db` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `dental_db`;
-- MySQL dump 10.13  Distrib 8.0.29, for macos12 (x86_64)
--
-- Host: 127.0.0.1    Database: dental_db
-- ------------------------------------------------------
-- Server version	8.0.29

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `APIAccessKey`
--

DROP TABLE IF EXISTS `APIAccessKey`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `APIAccessKey` (
  `Username` varchar(20) NOT NULL,
  `AccessKey` varchar(36) NOT NULL,
  PRIMARY KEY (`Username`),
  CONSTRAINT `fk_APIAccessKey_User_Username` FOREIGN KEY (`Username`) REFERENCES `User` (`Username`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `APIAccessKey`
--

LOCK TABLES `APIAccessKey` WRITE;
/*!40000 ALTER TABLE `APIAccessKey` DISABLE KEYS */;
INSERT INTO `APIAccessKey` VALUES ('admin','834f4e94-dcfc-11ec-a171-2aa4b7a723b7'),('amosburt','a83c53f0-dcfc-11ec-a171-2aa4b7a723b7'),('Apollo_A','3eb1c852-e22b-11ec-ba97-0242ac150002'),('aPond','ac670196-dcfc-11ec-a171-2aa4b7a723b7'),('bellaC','c1b80dc2-e180-11ec-b4e4-0242ac150002'),('BPotts','b1e2c8a8-dcfc-11ec-a171-2aa4b7a723b7'),('camina88','b8514840-dcfc-11ec-a171-2aa4b7a723b7'),('cOswald','c08ac4f0-dcfc-11ec-a171-2aa4b7a723b7'),('jHolden','c45779ca-dcfc-11ec-a171-2aa4b7a723b7'),('kara.Thrace','c8fb77ec-dcfc-11ec-a171-2aa4b7a723b7'),('maJones','ccb71b48-dcfc-11ec-a171-2aa4b7a723b7'),('riverS','d3652a66-dcfc-11ec-a171-2aa4b7a723b7'),('roster','d7d5308c-dcfc-11ec-a171-2aa4b7a723b7'),('we2332-8','d359c531-e228-11ec-ba97-0242ac150002');
/*!40000 ALTER TABLE `APIAccessKey` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Appointment`
--

DROP TABLE IF EXISTS `Appointment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Appointment` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `Patient` varchar(20) NOT NULL,
  `Dentist` varchar(20) NOT NULL,
  `Date` date NOT NULL,
  `Session` int NOT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `AppointmentSession` (`Dentist`,`Date`,`Session`),
  KEY `Patient_idx` (`Patient`),
  KEY `Dentist_idx` (`Dentist`),
  KEY `Session_idx` (`Session`),
  CONSTRAINT `fk_Appointment_AppointmentSession_ID` FOREIGN KEY (`Session`) REFERENCES `AppointmentSession` (`ID`),
  CONSTRAINT `fk_Appointment_User_Dentist_Username` FOREIGN KEY (`Dentist`) REFERENCES `User` (`Username`),
  CONSTRAINT `fk_Appointment_User_Patient_Username` FOREIGN KEY (`Patient`) REFERENCES `User` (`Username`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Appointment`
--

LOCK TABLES `Appointment` WRITE;
/*!40000 ALTER TABLE `Appointment` DISABLE KEYS */;
INSERT INTO `Appointment` VALUES (2,'cOswald','jHolden','2022-07-29',2),(3,'cOswald','camina88','2022-07-29',2),(4,'cOswald','amosburt','2022-08-17',6),(5,'riverS','camina88','2022-06-29',5),(14,'riverS','camina88','2022-08-08',4),(16,'we2332-8','amosburt','2022-11-01',7),(20,'kara.Thrace','camina88','2022-07-29',4),(21,'we2332-8','amosburt','2022-09-25',4),(22,'riverS','camina88','2022-10-18',5),(25,'riverS','amosburt','2022-06-28',4);
/*!40000 ALTER TABLE `Appointment` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `AppointmentSession`
--

DROP TABLE IF EXISTS `AppointmentSession`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `AppointmentSession` (
  `ID` int NOT NULL,
  `StartTime` varchar(5) NOT NULL,
  `EndTime` varchar(5) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `AppointmentSession`
--

LOCK TABLES `AppointmentSession` WRITE;
/*!40000 ALTER TABLE `AppointmentSession` DISABLE KEYS */;
INSERT INTO `AppointmentSession` VALUES (1,'09:00','10:00'),(2,'10:00','11:00'),(3,'11:00','12:00'),(4,'13:00','14:00'),(5,'14:00','15:00'),(6,'15:00','16:00'),(7,'16:00','18:00');
/*!40000 ALTER TABLE `AppointmentSession` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Authentication`
--

DROP TABLE IF EXISTS `Authentication`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Authentication` (
  `Username` varchar(20) NOT NULL,
  `Password` varchar(60) NOT NULL,
  PRIMARY KEY (`Username`),
  UNIQUE KEY `Username_UNIQUE` (`Username`),
  CONSTRAINT `fk_Authentication_User_Username` FOREIGN KEY (`Username`) REFERENCES `User` (`Username`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Authentication`
--

LOCK TABLES `Authentication` WRITE;
/*!40000 ALTER TABLE `Authentication` DISABLE KEYS */;
INSERT INTO `Authentication` VALUES ('admin','$2a$04$LDptspeQILSMS8XGrTRETetmfUekz22v2msichhvmYuYQJ8tzxAH.'),('amosburt','$2a$04$LDptspeQILSMS8XGrTRETetmfUekz22v2msichhvmYuYQJ8tzxAH.'),('aPond','$2a$04$LDptspeQILSMS8XGrTRETetmfUekz22v2msichhvmYuYQJ8tzxAH.'),('bellaC','$2a$04$c5bNY7tYOq/lBw1zioekGeikoWC5bMn4OHMXoX2Zkt3hNOp7DXS2W'),('BPotts','$2a$04$LDptspeQILSMS8XGrTRETetmfUekz22v2msichhvmYuYQJ8tzxAH.'),('camina88','$2a$04$LDptspeQILSMS8XGrTRETetmfUekz22v2msichhvmYuYQJ8tzxAH.'),('cOswald','$2a$04$LDptspeQILSMS8XGrTRETetmfUekz22v2msichhvmYuYQJ8tzxAH.'),('jHolden','$2a$04$LDptspeQILSMS8XGrTRETetmfUekz22v2msichhvmYuYQJ8tzxAH.'),('kara.Thrace','$2a$04$LDptspeQILSMS8XGrTRETetmfUekz22v2msichhvmYuYQJ8tzxAH.'),('maJones','$2a$04$LDptspeQILSMS8XGrTRETetmfUekz22v2msichhvmYuYQJ8tzxAH.'),('riverS','$2a$04$639OYoAGau8cE4xrdsOObelaxRn0IKhNDF7GdpSx7IhFXItapPSkG'),('roster','$2a$04$LDptspeQILSMS8XGrTRETetmfUekz22v2msichhvmYuYQJ8tzxAH.'),('we2332-8','$2a$04$cIMT2o/nN/ZguUznGoE8feo2YDB1pe0MVQCRSos6aeRHM2GvNW/0u');
/*!40000 ALTER TABLE `Authentication` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `User`
--

DROP TABLE IF EXISTS `User`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `User` (
  `Username` varchar(20) NOT NULL,
  `FirstName` varchar(30) NOT NULL,
  `LastName` varchar(30) NOT NULL,
  `MobileNumber` int DEFAULT NULL,
  `IsDeleted` tinyint(1) NOT NULL,
  `Role` varchar(10) NOT NULL,
  PRIMARY KEY (`Username`),
  UNIQUE KEY `Username_UNIQUE` (`Username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `User`
--

LOCK TABLES `User` WRITE;
/*!40000 ALTER TABLE `User` DISABLE KEYS */;
INSERT INTO `User` VALUES ('admin','admin','admin',NULL,0,'admin'),('amosburt','Amos','Burton',NULL,0,'dentist'),('Apollo_A','Lee','Adama',92345678,0,'patient'),('aPond','Amy','Pond',96402533,1,'patient'),('bellaC','fsssf','rrwer',0,0,'patient'),('BPotts','Bille','Potts',96345216,0,'patient'),('camina88','Camina','Drummer',NULL,0,'dentist'),('cOswald','Clara','Oswald',83724523,0,'patient'),('jHolden','James','Holden',NULL,0,'dentist'),('kara.Thrace','Kara','Thrace',84623847,0,'patient'),('maJones','Martha','Jones',87166096,0,'patient'),('riverS','River','Song',96418403,0,'patient'),('roster','Rosy','Tyler',85551064,1,'patient'),('we2332-8','test','testere',98374623,0,'patient');
/*!40000 ALTER TABLE `User` ENABLE KEYS */;
UNLOCK TABLES;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
/*!50003 CREATE*/ /*!50017 DEFINER=`root`@`%`*/ /*!50003 TRIGGER `trgDeleteSessionIfIsDeleted` AFTER UPDATE ON `User` FOR EACH ROW BEGIN
	IF !(NEW.IsDeleted <=> OLD.IsDeleted) THEN
      IF NEW.IsDeleted = true THEN 
		IF (SELECT count(*) FROM UserSession WHERE Username = OLD.Username) > 0 THEN
			DELETE FROM UserSession WHERE Username = OLD.Username;
        END IF;
      END IF;
   END IF;
END */;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;

--
-- Table structure for table `UserSession`
--

DROP TABLE IF EXISTS `UserSession`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `UserSession` (
  `Username` varchar(20) NOT NULL,
  `SessionID` varchar(36) NOT NULL,
  `ExpireDT` datetime NOT NULL,
  `LoginDT` datetime NOT NULL,
  PRIMARY KEY (`Username`),
  CONSTRAINT `fk_UserSession_User_Username` FOREIGN KEY (`Username`) REFERENCES `User` (`Username`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `UserSession`
--

LOCK TABLES `UserSession` WRITE;
/*!40000 ALTER TABLE `UserSession` DISABLE KEYS */;
INSERT INTO `UserSession` VALUES ('admin','1482c8dc-e25c-11ec-a7d9-0242ac160002','2022-06-03 10:09:25','2022-06-02 10:09:25'),('riverS','d9b0310e-e25a-11ec-a7d9-0242ac160002','2022-06-03 10:00:37','2022-06-02 10:00:37');
/*!40000 ALTER TABLE `UserSession` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping events for database 'dental_db'
--
/*!50106 SET @save_time_zone= @@TIME_ZONE */ ;
/*!50106 DROP EVENT IF EXISTS `eventDeleteUserSession` */;
DELIMITER ;;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;;
/*!50003 SET character_set_client  = utf8mb4 */ ;;
/*!50003 SET character_set_results = utf8mb4 */ ;;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;;
/*!50003 SET @saved_time_zone      = @@time_zone */ ;;
/*!50003 SET time_zone             = 'SYSTEM' */ ;;
/*!50106 CREATE*/ /*!50117 DEFINER=`root`@`%`*/ /*!50106 EVENT `eventDeleteUserSession` ON SCHEDULE EVERY 1 MINUTE STARTS '2022-06-02 07:14:22' ENDS '2032-06-02 07:14:22' ON COMPLETION NOT PRESERVE ENABLE DO DELETE FROM UserSession WHERE ExpireDT < NOW() */ ;;
/*!50003 SET time_zone             = @saved_time_zone */ ;;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;;
/*!50003 SET character_set_client  = @saved_cs_client */ ;;
/*!50003 SET character_set_results = @saved_cs_results */ ;;
/*!50003 SET collation_connection  = @saved_col_connection */ ;;
DELIMITER ;
/*!50106 SET TIME_ZONE= @save_time_zone */ ;

--
-- Dumping routines for database 'dental_db'
--
/*!50003 DROP PROCEDURE IF EXISTS `spAccessKeyExist` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `spAccessKeyExist`(pAccessKey VARCHAR(36))
BEGIN
	SELECT EXISTS (
		SELECT * 
		FROM APIAccessKey api
		INNER JOIN User u
		ON api.Username = u.Username
		WHERE AccessKey = pAccessKey
		AND isDeleted = false
	);
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spAppointmentCreate` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `spAppointmentCreate`(
	pPatient VARCHAR(20),
    pDentist VARcHAR(20),
    pDate DATE,
    pSession INT
)
BEGIN
	INSERT INTO Appointment (Patient, Dentist, Date, Session)
    VALUES (pPatient, pDentist, pDate, pSession);
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spAppointmentUpdate` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `spAppointmentUpdate`(
	pID INT,
    pPatient VARCHAR(20),
    pDentist VARcHAR(20),
    pDate DATE,
    pSession INT
)
BEGIN
	UPDATE Appointment SET
		Patient = COALESCE(pPatient, Patient),
        Dentist = COALESCE(pDentist, Dentist),
        Date = COALESCE(pDate, Date),
        Session = COALESCE(pSession, Session)
	WHERE ID = pID;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spAuthenticationCreate` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `spAuthenticationCreate`(pUsername VARCHAR(20), pPassword VARCHAR(60))
BEGIN
	INSERT INTO Authentication (Username, Password)
    VALUES (pUsername, pPassword);
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spAuthenticationGet` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `spAuthenticationGet`(pUsername VARCHAR(20))
BEGIN
	SELECT auth.*, akey.AccessKey 
    FROM Authentication auth
    INNER JOIN APIAccessKey akey
    ON auth.Username = akey.Username
    INNER JOIN User u
    ON auth.Username = u.Username
    WHERE auth.Username = pUsername
    AND u.IsDeleted = false;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spAuthenticationUpdate` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `spAuthenticationUpdate`(pUsername VARCHAR(20), pPassword VARCHAR(60))
BEGIN
	UPDATE Authentication 
    SET Password = pPassword
    WHERE Username = pUsername;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spUserCreate` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `spUserCreate`(
	pUsername VARCHAR(20),
    pFirstName VARCHAR(30),
    pLastName VARCHAR(30),
    pMobileNumber INT,
    pRole VARCHAR(10)
)
BEGIN

INSERT INTO User (Username, FirstName, LastName, MobileNumber, isDeleted, Role)
VALUES (pUsername, pFirstName, pLastName, pMobileNumber, false, pRole);

INSERT INTO APIAccessKey(Username, AccessKey)
VALUES (pUsername, UUID());

END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spUserExistByUsername` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `spUserExistByUsername`(pUsername varchar(20))
BEGIN
	SELECT EXISTS (SELECT * FROM User WHERE Username = pUsername AND isDeleted = false);
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spUserGet` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `spUserGet`(pRole varchar(10), pUsername varchar(20))
BEGIN
	IF pRole = 'admin' THEN
		SELECT * FROM User WHERE Username = pUsername;
	ELSE
		SELECT * FROM User WHERE Username = pUsername AND isDeleted = false;
END IF;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spUserGetByAccessKey` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `spUserGetByAccessKey`(pAccessKey VARCHAR(36))
BEGIN
	SELECT u.* FROM User u
	INNER JOIN APIAccessKey api
	ON u.Username = api.Username
	WHERE api.AccessKey = pAccessKey;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spUserGetDentistAvailability` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `spUserGetDentistAvailability`(pDentist VARCHAR(20), pDate DATE)
BEGIN
	SELECT ap.Dentist, ap.Date, aps.ID, aps.StartTime, aps.EndTime 
	FROM (
		SELECT * 
		FROM Appointment ap
		WHERE ap.Dentist = pDentist
		AND ap.Date = pDate
	) ap 
	RIGHT JOIN AppointmentSession aps
	ON ap.Session = aps.ID;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spUserGetRoleByAccessKey` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `spUserGetRoleByAccessKey`(pAccessKey VARCHAR(36))
BEGIN
	SELECT u.Role FROM User u
	INNER JOIN APIAccessKey api
	ON u.Username = api.Username
	WHERE api.AccessKey = pAccessKey;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spUserSessionCreate` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `spUserSessionCreate`(pUsername VARCHAR(20), pExpireDT DATETIME, pLoginDT DATETIME)
BEGIN
	INSERT INTO UserSession(Username, SessionID, ExpireDT, LoginDT)
    VALUES (pUsername, UUID(), pExpireDT, pLoginDT);
    
    SELECT SessionID FROM UserSession WHERE USername = pUsername;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spUserSessionDelete` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `spUserSessionDelete`(pUsername VARCHAR(60))
BEGIN
	IF (SELECT count(*) FROM UserSession WHERE username = pUsername) > 0 THEN
		DELETE FROM UserSession WHERE Username = pUsername;
	END IF;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spUserSessionDeleteByUsername` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `spUserSessionDeleteByUsername`(pUsername VARCHAR(60))
BEGIN
	IF (SELECT count(*) FROM UserSession WHERE username = pUsername) > 0 THEN
		DELETE FROM UserSession WHERE Username = pUsername;
	END IF;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spUserSessionExistByUsername` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `spUserSessionExistByUsername`(pUsername VARCHAR(20))
BEGIN
	SELECT EXISTS (SELECT * FROM UserSession WHERE Username = pUsername);
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spUserSessionGet` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `spUserSessionGet`(pSessionID VARcHAR(36))
BEGIN
	SELECT us.Username, u.FirstName, u.LastName, u.Role, us.ExpireDT, api.AccessKey
	FROM UserSession us
	INNER JOIN APIAccessKey api
	ON us.Username = api.Username
	INNER JOIN User u
	ON us.Username = u.Username
	WHERE us.SessionID = pSessionID;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spUserSessionGetAll` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `spUserSessionGetAll`()
BEGIN
	SELECT us.Username, u.Role, us.SessionID, us.LoginDT
	FROM UserSession us
	INNER JOIN User u
	ON us.Username = u.Username;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spUserUpdate` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `spUserUpdate`(
	pUsername VARCHAR(20),
    pFirstName VARCHAR(30),
    pLastName VARCHAR(30),
    pMobileNumber INT,
    pIsDeleted BOOL,
    pRole VARCHAR(10)
)
BEGIN
	UPDATE User SET
		FirstName = COALESCE(pFirstName, FirstName),
        LastName = COALESCE(pLastName, LastName),
        MobileNumber = COALESCE(pMobileNumber, MobileNumber),
        IsDeleted = COALESCE(pIsDeleted, IsDeleted),
        Role = COALESCE(pRole, Role)
	WHERE Username = pUsername;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-06-02 18:34:39
