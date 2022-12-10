CREATE DATABASE `main`;

USE `main`;

CREATE TABLE `paths` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `person_id` varchar(64)  NOT NULL,
  `path` varchar(256)  NOT NULL,
  `departure_time` bigint,
  `insert_time` bigint NOT NULL COMMENT 'record insertion time',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;


#FLUSH PRIVILEGES;
#CREATE USER 'root'@'localhost' IDENTIFIED BY 'local';
#GRANT ALL PRIVILEGES ON *.* TO 'root'@'%';
#CREATE USER 'me'@'localhost' IDENTIFIED BY '<password>';
