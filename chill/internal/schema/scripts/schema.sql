-- Initialize schema

CREATE DATABASE IF NOT EXISTS chill;
use chill;

CREATE TABLE `command_log` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `build_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'build no',
  `command` text COMMENT 'command',
  `remark` varchar(1000) NOT NULL DEFAULT '' COMMENT 'remark',
  `status` varchar(10) NOT NULL DEFAULT '' COMMENT 'status: SUCCESSFUL, FAILED',
  `ctime` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created at',
  PRIMARY KEY (`id`),
  KEY `build_no_idx` (`build_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Command Log';

CREATE TABLE `build_info` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT 'build name',
  `status` varchar(10) NOT NULL DEFAULT '' COMMENT 'status: SUCCESSFUL, FAILED',
  `ctime` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created at',
  `utime` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated at',
  `commit_id` varchar(45) NOT NULL DEFAULT '' COMMENT 'commit id',
  `tag` varchar(45) NOT NULL DEFAULT '' COMMENT 'tag',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_uk` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Build Info';

CREATE TABLE `build_log` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `build_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'build no',
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT 'build name',
  `status` varchar(10) NOT NULL DEFAULT '' COMMENT 'status: SUCCESSFUL, FAILED',
  `remark` varchar(1000) NOT NULL DEFAULT '' COMMENT 'remark',
  `ctime` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created at',
  `utime` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated at',
  `build_start_time` timestamp NULL DEFAULT NULL COMMENT 'build start time',
  `build_end_time` timestamp NULL DEFAULT NULL COMMENT 'build end time',
  `commit_id` varchar(45) NOT NULL DEFAULT '' COMMENT 'commit id',
  `tag` varchar(45) NOT NULL DEFAULT '' COMMENT 'tag',
  PRIMARY KEY (`id`),
  UNIQUE KEY `build_no_uk` (`build_no`),
  KEY `name_idx` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Build Log';
