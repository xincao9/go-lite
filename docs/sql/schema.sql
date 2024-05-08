CREATE
DATABASE `golite` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE
`golite`;

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
    `id`         bigint(11) unsigned NOT NULL AUTO_INCREMENT,
    `username`   char(32)  NOT NULL DEFAULT '',
    `password`   char(32)  NOT NULL DEFAULT '',
    `expire`     timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `token`      char(128) NOT NULL DEFAULT '',
    `role`       int(11) NOT NULL DEFAULT '0',
    `created_at` timestamp NULL DEFAULT NULL,
    `deleted_at` timestamp NULL DEFAULT NULL,
    `updated_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `username_idx` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
