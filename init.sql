CREATE DATABASE `main`;

USE `main`;

CREATE TABLE `paths` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `person_id` varchar(64)  NOT NULL,
  `path` varchar(256)  NOT NULL,
  `departure_time` bigint,
  `insert_time` bigint NOT NULL COMMENT 'record insertion time',
  CONSTRAINT person_record_constraint UNIQUE (`person_id`, `departure_time`),
  INDEX person_record (`person_id`, `departure_time`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
