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
(2,	'Abdul',	'Abdullah'),
(3,	'Barbara',	'Bonato'),
(4,	'Baldassarre',	'Baccani'),
(5,	'Carlo',	'Carraretto'),
(6,	'Cinzia',	'Camilleri'),
(7,	'Dario',	'Donati'),
(8,	'Duccio',	'Doritos'),
(9,	'Dalila',	'Danzanti'),
(10,	'Enrico',	'Emilietti'),
(11,	'Elio',		'E le storie tese'),
(12,	'Fulvio',	'Fabbian'),
(13,	'Gerardo',	'Gianfigliazzi'),
(14,	'Gaia',		'Gibboni'),
(15,	'Gabriele',	'Gastrite'),
(16,	'Harald',	'Hardrada'),
(17,	'Ivan',		'Iacopini'),
(18,	'Irene',	'Iddau');

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
(3,	'Carme'),
(4,	'Distopia'),
(5,	'Erotico'),
(6,	'Fantascienza'),
(7,	'Giallo'),
(8,	'Horror'),
(9,	'Imeneo (lirica)');

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
(2,	'Alla Salute!',	2,	1,	0,	''),
(3,	'Argo Is Better Than Nuvola. Change My Mind.',	2,	1,	0,	''),
(4,	'Beatrice Bortolin - La Vita',	4,	2,	0,	''),
(5,	'Bear Grylls',	3,	2,	0,	''),
(6,	'Complotto Terrapiattista',	5,	3,	0,	''),
(7,	'Castagne Autunnali',	5,	3,	0,	''),
(8,	'Cigarini All''Atalanta',	5,	3,	0,	''),
(9,	'C++ E Java A Colazione',	6,	3,	0,	''),
(10,	'Dieci E Lode',	7,	4,	0,	''),
(11,	'Duffy Duck E Il Calice Di Fuoco',	8,	4,	0,	''),
(12,	'Dragon Ball Z - Goku Si Converte Al GodSucca',	8,	4,	0,	''),
(13,	'Dente Di Leone',	9,	4,	0,	''),
(14,	'Durlindana Si Spezza',	7,	4,	0,	''),
(15,	'Elton John: Vite Parallele',	10,	5,	0,	''),
(16,	'Edimburgo Ad Aprile',	10,	5,	0,	''),
(17,	'Edin Dzeko Gelataio',	11,	5,	0,	''),
(18,	'Ettore, Andromaca E Priamo',	11,	5,	0,	''),
(19,	'Funghi Nucleari',	12,	6,	0,	''),
(20,	'Fallimenti Aziendali',	12,	6,	0,	''),
(21,	'Gas A Martello, PRIMA DENTROOO',	13,	7,	0,	''),
(22,	'Gigi D''Alessio, Il Killer Di Apparati Uditivi',	13,	7,	0,	''),
(23,	'Gallio',	14,	7,	0,	''),
(24,	'Gelsomino Turlupinato',	15,	7,	0,	''),
(25,	'Ho La Faringe Infiammata',	16,	8,	0,	''),
(26,	'Highway To Hell',	16,	8,	0,	''),
(27,	'Iolanda Vive In Islanda',	17,	9,	0,	''),
(28,	'Istrici A Colazione',	17,	9,	0,	''),
(29,	'Io Sono Dio',	17,	9,	0,	''),
(30,	'Isabella Di Castiglia',	18,	9,	0,	''),
(31,	'Iennarientu: Meraviglie Sarde',	18,	9,	0,	'');

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
