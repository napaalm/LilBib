/*
 * lilbib_example.sql
 *
 * Struttura del database di sviluppo per lilbib.
 *
 * Copyright (c) 2020 Antonio Napolitano <nap@antonionapolitano.eu>
 * Copyright (c) 2020 Davide Vendramin <davidevendramin5@gmail.com>
 *
 * This file is part of LilBib.
 *
 * LilBib is free software; you can redistribute it and/or modify it
 * under the terms of the Affero GNU General Public License as
 * published by the Free Software Foundation; either version 3, or (at
 * your option) any later version.
 *
 * LilBib is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
 * or FITNESS FOR A PARTICULAR PURPOSE.  See the Affero GNU General
 * Public License for more details.
 *
 * You should have received a copy of the Affero GNU General Public
 * License along with LilBib; see the file LICENSE. If not see
 * <http://www.gnu.org/licenses/>.
 */

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

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
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

TRUNCATE `Autore`;
INSERT INTO `Autore` (`Codice`, `Nome`, `Cognome`) VALUES
(1,	'Alice',	'Adelfi'),
(2,	'Barbara',	'Bonato'),
(3,	'Carlo',	'Carraretto'),
(4,	'Dario',	'Donati'),
(5,	'Enrico',	'Emilietti'),
(6,	'Fulvio',	'Fabbian'),
(7,	'Gerardo',	'Gianfigliazzi');

DROP TABLE IF EXISTS `Genere`;
CREATE TABLE `Genere` (
  `Codice` int(11) NOT NULL AUTO_INCREMENT,
  `Nome` varchar(20) CHARACTER SET utf32 NOT NULL,
  PRIMARY KEY (`Codice`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

TRUNCATE `Genere`;
INSERT INTO `Genere` (`Codice`, `Nome`) VALUES
(1,	'Avventura e azione'),
(2,	'Biografia'),
(4,	'Distopia'),
(5,	'Erotico'),
(6,	'Fantascienza'),
(7,	'Giallo');

DROP TABLE IF EXISTS `Libro`;
CREATE TABLE `Libro` (
  `Codice` int(11) NOT NULL AUTO_INCREMENT,
  `Titolo` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Autore` int(11) NOT NULL,
  `Genere` int(11) NOT NULL,
  `Prenotato` tinyint(4) NOT NULL,
  `Hashz` char(44) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`Codice`),
  KEY `Genere` (`Genere`),
  KEY `Autore` (`Autore`),
  CONSTRAINT `Libro_ibfk_4` FOREIGN KEY (`Genere`) REFERENCES `Genere` (`Codice`) ON UPDATE CASCADE,
  CONSTRAINT `Libro_ibfk_6` FOREIGN KEY (`Autore`) REFERENCES `Autore` (`Codice`) ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

TRUNCATE `Libro`;
INSERT INTO `Libro` (`Codice`, `Titolo`, `Autore`, `Genere`, `Prenotato`, `Hashz`) VALUES
(1,	'Alla Ricerca Di Adele',	1,	1,	1,	''),
(2,	'Beatrice Bortolin - La Vita',	2,	2,	0,	''),
(4,	'Dieci E Lode',	4,	4,	0,	'');

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
  CONSTRAINT `Prestito_ibfk_3` FOREIGN KEY (`Libro`) REFERENCES `Libro` (`Codice`) ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

TRUNCATE `Prestito`;
INSERT INTO `Prestito` (`Codice`, `Libro`, `Utente`, `Data_prenotazione`, `Durata`, `Data_restituzione`) VALUES
(1,	1,	'Alfio.Ammannati',	1013126400,	1209600,	NULL),
(2,	2,	'Bartolomeo.Bianchi',	1013126400,	1209600,	1014336000);

CREATE USER 'lilbib'@'%' IDENTIFIED BY 'secret';
GRANT ALL PRIVILEGES ON lilbib.* TO 'lilbib'@'%';
