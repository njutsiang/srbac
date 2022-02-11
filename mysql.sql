-- noinspection SqlDialectInspectionForFile

/*
SQLyog Community v13.1.9 (64 bit)
MySQL - 5.7.27 : Database - srbac
*********************************************************************
*/

/*!40101 SET NAMES utf8 */;
/*!40101 SET SQL_MODE=''*/;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

/*Table structure for table `api_item` */

CREATE TABLE `api_item` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `service_id` int(10) unsigned NOT NULL COMMENT '服务 id',
  `method` char(8) NOT NULL COMMENT '请求方式',
  `uri` varchar(128) NOT NULL COMMENT '接口路由',
  `name` varchar(32) NOT NULL COMMENT '接口名称',
  `is_anonymous_access` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否允许匿名文档：1=允许，0=不允许',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `created_at` int(10) unsigned NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `service_id_method_uri` (`service_id`,`method`,`uri`),
  CONSTRAINT `api_item_service_id` FOREIGN KEY (`service_id`) REFERENCES `service` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='接口节点';

/*Table structure for table `data_item` */

CREATE TABLE `data_item` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `service_id` int(10) unsigned NOT NULL COMMENT '服务 id',
  `key` varchar(32) NOT NULL COMMENT '权限标识',
  `name` varchar(32) NOT NULL COMMENT '权限名称',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `created_at` int(10) unsigned NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `service_id_key` (`service_id`,`key`),
  CONSTRAINT `data_item_service_id` FOREIGN KEY (`service_id`) REFERENCES `service` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据节点';

/*Table structure for table `menu_item` */

CREATE TABLE `menu_item` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `service_id` int(10) unsigned NOT NULL COMMENT '服务 id',
  `key` varchar(32) NOT NULL COMMENT '菜单标识',
  `name` varchar(32) NOT NULL COMMENT '菜单名称',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `created_at` int(11) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `service_id_key` (`service_id`,`key`),
  CONSTRAINT `menu_item_service_id` FOREIGN KEY (`service_id`) REFERENCES `service` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='菜单节点';

/*Table structure for table `role` */

CREATE TABLE `role` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `key` varchar(32) NOT NULL COMMENT '角色标识',
  `name` varchar(32) NOT NULL COMMENT '角色名称',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `created_at` int(10) unsigned NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `key` (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色';

/*Table structure for table `role_api_item` */

CREATE TABLE `role_api_item` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `role_id` int(10) unsigned NOT NULL COMMENT '角色 id',
  `service_id` int(10) unsigned NOT NULL COMMENT '服务 id',
  `api_item_id` int(10) unsigned NOT NULL COMMENT '接口节点 id',
  `created_at` int(10) unsigned NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `role_id_api_item_id` (`role_id`,`api_item_id`),
  KEY `role_api_item_service_id` (`service_id`),
  KEY `role_api_item_api_item_id` (`api_item_id`),
  CONSTRAINT `role_api_item_api_item_id` FOREIGN KEY (`api_item_id`) REFERENCES `api_item` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `role_api_item_role_id` FOREIGN KEY (`role_id`) REFERENCES `role` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `role_api_item_service_id` FOREIGN KEY (`service_id`) REFERENCES `service` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色接口节点分配';

/*Table structure for table `role_data_item` */

CREATE TABLE `role_data_item` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `role_id` int(10) unsigned NOT NULL COMMENT '角色 id',
  `service_id` int(10) unsigned NOT NULL COMMENT '服务 id',
  `data_item_id` int(10) unsigned NOT NULL COMMENT '数据节点 id',
  `created_at` int(10) unsigned NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `role_id_data_item_id` (`role_id`,`data_item_id`),
  KEY `role_data_item_service_id` (`service_id`),
  KEY `role_data_item_data_item_id` (`data_item_id`),
  CONSTRAINT `role_data_item_data_item_id` FOREIGN KEY (`data_item_id`) REFERENCES `data_item` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `role_data_item_role_id` FOREIGN KEY (`role_id`) REFERENCES `role` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `role_data_item_service_id` FOREIGN KEY (`service_id`) REFERENCES `service` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色数据节点分配';

/*Table structure for table `role_menu_item` */

CREATE TABLE `role_menu_item` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `role_id` int(10) unsigned NOT NULL COMMENT '角色 id',
  `service_id` int(10) unsigned NOT NULL COMMENT '服务 id',
  `menu_item_id` int(10) unsigned NOT NULL COMMENT '菜单节点 id',
  `created_at` int(10) unsigned NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `role_id_menu_item_id` (`role_id`,`menu_item_id`),
  KEY `role_menu_item_service_id` (`service_id`),
  KEY `role_menu_item_menu_item_id` (`menu_item_id`),
  CONSTRAINT `role_menu_item_menu_item_id` FOREIGN KEY (`menu_item_id`) REFERENCES `menu_item` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `role_menu_item_role_id` FOREIGN KEY (`role_id`) REFERENCES `role` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `role_menu_item_service_id` FOREIGN KEY (`service_id`) REFERENCES `service` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色菜单节点分配';

/*Table structure for table `role_service` */

CREATE TABLE `role_service` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `role_id` int(10) unsigned NOT NULL COMMENT '角色 id',
  `service_id` int(10) unsigned NOT NULL COMMENT '服务 id',
  `created_at` int(10) unsigned NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `role_id_service_id` (`role_id`,`service_id`),
  KEY `role_service_service_id` (`service_id`),
  CONSTRAINT `role_service_role_id` FOREIGN KEY (`role_id`) REFERENCES `role` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `role_service_service_id` FOREIGN KEY (`service_id`) REFERENCES `service` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色服务分配';

/*Table structure for table `service` */

CREATE TABLE `service` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `key` varchar(32) NOT NULL COMMENT '服务标识',
  `name` varchar(32) NOT NULL DEFAULT '' COMMENT '服务名称',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `created_at` int(10) unsigned NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `key` (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='服务';

/*Table structure for table `user` */

CREATE TABLE `user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(32) NOT NULL DEFAULT '' COMMENT '姓名',
  `username` varchar(32) NOT NULL COMMENT '用户名',
  `password` varchar(128) NOT NULL DEFAULT '' COMMENT '密码',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '状态：1=启用，0=禁用',
  `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `created_at` int(10) unsigned NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户';

/*Table structure for table `user_api_item` */

CREATE TABLE `user_api_item` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` int(10) unsigned NOT NULL COMMENT '用户 id',
  `service_id` int(10) unsigned NOT NULL COMMENT '服务 id',
  `api_item_id` int(10) unsigned NOT NULL COMMENT '接口节点 id',
  `created_at` int(10) unsigned NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id_api_item_id` (`user_id`,`api_item_id`),
  KEY `user_api_item_service_id` (`service_id`),
  KEY `user_api_item_api_item_id` (`api_item_id`),
  CONSTRAINT `user_api_item_api_item_id` FOREIGN KEY (`api_item_id`) REFERENCES `api_item` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `user_api_item_service_id` FOREIGN KEY (`service_id`) REFERENCES `service` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `user_api_item_user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户接口节点分配';

/*Table structure for table `user_data_item` */

CREATE TABLE `user_data_item` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` int(10) unsigned NOT NULL COMMENT '用户 id',
  `service_id` int(10) unsigned NOT NULL COMMENT '服务 id',
  `data_item_id` int(10) unsigned NOT NULL COMMENT '数据节点 id',
  `created_at` int(10) unsigned NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id_data_item_id` (`user_id`,`data_item_id`),
  KEY `user_data_item_service_id` (`service_id`),
  KEY `user_data_item_data_item_id` (`data_item_id`),
  CONSTRAINT `user_data_item_data_item_id` FOREIGN KEY (`data_item_id`) REFERENCES `data_item` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `user_data_item_service_id` FOREIGN KEY (`service_id`) REFERENCES `service` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `user_data_item_user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户数据节点分配';

/*Table structure for table `user_menu_item` */

CREATE TABLE `user_menu_item` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` int(10) unsigned NOT NULL COMMENT '用户 id',
  `service_id` int(10) unsigned NOT NULL COMMENT '服务 id',
  `menu_item_id` int(10) unsigned NOT NULL COMMENT '菜单节点 id',
  `created_at` int(10) unsigned NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id_menu_item_id` (`user_id`,`menu_item_id`),
  KEY `user_menu_item_service_id` (`service_id`),
  KEY `user_menu_item_menu_item_id` (`menu_item_id`),
  CONSTRAINT `user_menu_item_menu_item_id` FOREIGN KEY (`menu_item_id`) REFERENCES `menu_item` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `user_menu_item_service_id` FOREIGN KEY (`service_id`) REFERENCES `service` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `user_menu_item_user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户菜单节点分配';

/*Table structure for table `user_role` */

CREATE TABLE `user_role` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` int(10) unsigned NOT NULL COMMENT '用户 id',
  `role_id` int(10) unsigned NOT NULL COMMENT '角色 id',
  `created_at` int(10) unsigned NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id_role_id` (`user_id`,`role_id`),
  KEY `user_role_role_id` (`role_id`),
  CONSTRAINT `user_role_role_id` FOREIGN KEY (`role_id`) REFERENCES `role` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `user_role_user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色分配';

/*Table structure for table `user_service` */

CREATE TABLE `user_service` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` int(10) unsigned NOT NULL COMMENT '用户 id',
  `service_id` int(10) unsigned NOT NULL COMMENT '服务 id',
  `created_at` int(10) unsigned NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id_service_id` (`user_id`,`service_id`),
  KEY `user_service_service_id` (`service_id`),
  CONSTRAINT `user_service_service_id` FOREIGN KEY (`service_id`) REFERENCES `service` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `user_service_user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户服务分配';

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
