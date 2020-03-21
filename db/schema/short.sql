DROP TABLE IF EXISTS `short`;
CREATE TABLE `short` (
`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
`long_url` varchar(1024) NOT NULL,
`short_url` varchar(11) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
`create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
PRIMARY KEY (`id`),
UNIQUE KEY `idx_uniq_short_long_url` (`long_url`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
