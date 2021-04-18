-- MariaDB dump 10.17  Distrib 10.4.10-MariaDB, for osx10.15 (x86_64)
--
-- Host: localhost    Database: ecp
-- ------------------------------------------------------
-- Server version	5.7.24

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Current Database: `ecp`
--

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `ecp` /*!40100 DEFAULT CHARACTER SET latin1 */;

USE `ecp`;

--
-- Table structure for table `account`
--

DROP TABLE IF EXISTS `account`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `account` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `address` varchar(50) NOT NULL,
  `code` varchar(8) CHARACTER SET utf8 DEFAULT NULL,
  `referrals` json DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `account_address_uindex` (`address`),
  UNIQUE KEY `account_code_uindex` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `referral_earn`
--

DROP TABLE IF EXISTS `referral_earn`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `referral_earn` (
  `address` varchar(50) NOT NULL,
  `block_number` varchar(50) NOT NULL,
  `amount` decimal(20,8) NOT NULL,
  `time_stamp` varchar(50) NOT NULL,
  PRIMARY KEY (`address`,`block_number`),
  KEY `referral_earn_address_index` (`address`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `transaction`
--

DROP TABLE IF EXISTS `transaction`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `transaction` (
  `block_number` varchar(50) NOT NULL,
  `time_stamp` varchar(50) DEFAULT NULL,
  `hash` varchar(255) DEFAULT NULL,
  `nonce` varchar(50) DEFAULT NULL,
  `block_hash` varchar(255) DEFAULT NULL,
  `from` varchar(50) DEFAULT NULL,
  `contract_address` varchar(50) DEFAULT NULL,
  `to` varchar(50) DEFAULT NULL,
  `value` varchar(50) DEFAULT NULL,
  `token_name` varchar(50) DEFAULT NULL,
  `token_decimal` varchar(50) DEFAULT NULL,
  `token_symbol` varchar(50) DEFAULT NULL,
  `transaction_index` varchar(50) DEFAULT NULL,
  `gas` varchar(50) DEFAULT NULL,
  `gas_price` varchar(50) DEFAULT NULL,
  `gas_used` varchar(50) DEFAULT NULL,
  `cumulative_gas_used` varchar(50) DEFAULT NULL,
  `input` varchar(50) DEFAULT NULL,
  `confirmations` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`block_number`),
  KEY `transaction_from_index` (`from`),
  KEY `transaction_to_index` (`to`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-04-18 14:00:06
