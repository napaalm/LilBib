-- Adminer 4.7.7 MySQL dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;

SET NAMES utf8mb4;

DROP DATABASE IF EXISTS `lilbib`;
CREATE DATABASE `lilbib` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */;
USE `lilbib`;

DROP TABLE IF EXISTS `Autore`;
CREATE TABLE `Autore` (
  `Codice` int(11) NOT NULL AUTO_INCREMENT,
  `Nome` varchar(20) CHARACTER SET utf8 NOT NULL,
  `Cognome` varchar(20) CHARACTER SET utf8 NOT NULL,
  PRIMARY KEY (`Codice`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


DROP TABLE IF EXISTS `Genere`;
CREATE TABLE `Genere` (
  `Codice` int(11) NOT NULL AUTO_INCREMENT,
  `Nome` varchar(20) CHARACTER SET utf32 NOT NULL,
  PRIMARY KEY (`Codice`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


DROP TABLE IF EXISTS `Libro`;
CREATE TABLE `Libro` (
  `Codice` int(11) NOT NULL AUTO_INCREMENT,
  `Titolo` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Autore` int(11) NOT NULL,
  `Genere` int(11) NOT NULL,
  `Prenotato` tinyint(4) NOT NULL,
  `Hashz` char(16) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`Codice`),
  KEY `Autore` (`Autore`),
  KEY `Genere` (`Genere`),
  CONSTRAINT `Libro_ibfk_3` FOREIGN KEY (`Autore`) REFERENCES `Autore` (`Codice`) ON UPDATE CASCADE,
  CONSTRAINT `Libro_ibfk_4` FOREIGN KEY (`Genere`) REFERENCES `Genere` (`Codice`) ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


DROP TABLE IF EXISTS `Prestito`;
CREATE TABLE `Prestito` (
  `Codice` int(11) NOT NULL AUTO_INCREMENT,
  `Libro` int(11) NOT NULL,
  `Utente` varchar(20) CHARACTER SET utf8mb4 NOT NULL,
  `Data_prenotazione` int(11) NOT NULL,
  `Durata` int(11) NOT NULL,
  `Data_restituzione` int(11) DEFAULT NULL,
  PRIMARY KEY (`Codice`),
  KEY `Libro` (`Libro`),
  CONSTRAINT `Prestito_ibfk_2` FOREIGN KEY (`Libro`) REFERENCES `Libro` (`Codice`) ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- 2020-06-02 16:36:59
