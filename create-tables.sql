-- create database blog;
-- use blog;
-- source create-tables.sql;
DROP TABLE IF EXISTS blog_tag, blog_category, blog_article, blog_auth, blog_article_tag;

CREATE TABLE `blog_tag` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) DEFAULT '' COMMENT '标签名称',
  -- `created_on` int(10) unsigned DEFAULT '0' COMMENT '创建时间',
  `created_at` TIMESTAMP NULL DEFAULT NULL COMMENT '创建时间',
  `created_by` varchar(100) DEFAULT '' COMMENT '创建人',
  -- `modified_on` int(10) unsigned DEFAULT '0' COMMENT '修改时间',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT '修改时间',
  `updated_by` varchar(100) DEFAULT '' COMMENT '修改人',
  `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',
  `deleted_by` varchar(100) DEFAULT '' COMMENT '删除人',
  `state` tinyint(3) unsigned DEFAULT '1' COMMENT '状态 0为禁用、1为启用',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章标签管理';

CREATE TABLE `blog_category` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) DEFAULT '' COMMENT '分类名称',
  `created_at` TIMESTAMP NULL DEFAULT NULL COMMENT '创建时间',
  `created_by` varchar(100) DEFAULT '' COMMENT '创建人',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT '修改时间',
  `updated_by` varchar(100) DEFAULT '' COMMENT '修改人',
  `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',
  `deleted_by` varchar(100) DEFAULT '' COMMENT '删除人',
  `state` tinyint(3) unsigned DEFAULT '1' COMMENT '状态 0为禁用、1为启用',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章分类管理';

CREATE TABLE `blog_article` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `category_id` int(10) unsigned DEFAULT '0' COMMENT '分类ID',
  `title` varchar(100) DEFAULT '' COMMENT '文章标题',
  `desc` varchar(255) DEFAULT '' COMMENT '简述',
  `content` text,
  `created_at` TIMESTAMP NULL DEFAULT NULL COMMENT '创建时间',
  `created_by` varchar(100) DEFAULT '' COMMENT '创建人',
  `updated_at` TIMESTAMP NULL DEFAULT NULL COMMENT '修改时间',
  `updated_by` varchar(255) DEFAULT '' COMMENT '修改人',
  `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',
  `deleted_by` varchar(100) DEFAULT '' COMMENT '删除人',
  `state` tinyint(3) unsigned DEFAULT '1' COMMENT '状态 0为禁用1为启用',
  INDEX idx_deleted_at (deleted_at),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章管理';

CREATE TABLE `blog_article_tag` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `article_id` int(10) unsigned DEFAULT '0' COMMENT '文章ID',
  `tag_id` int(10) unsigned DEFAULT '0' COMMENT '标签ID',
  `created_at` TIMESTAMP NULL DEFAULT NULL COMMENT '创建时间',
  `created_by` varchar(100) DEFAULT '' COMMENT '创建人',
  CONSTRAINT `fk_article` FOREIGN KEY (`article_id`) REFERENCES `blog_article`(`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_tag` FOREIGN KEY (`tag_id`) REFERENCES `blog_tag`(`id`) ON DELETE CASCADE,
  UNIQUE KEY `uk_article_tag` (`article_id`, `tag_id`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章分类管理';

CREATE TABLE `blog_auth` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) DEFAULT '' COMMENT '账号',
  `password` varchar(100) DEFAULT '' COMMENT '密码',
  `deleted_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `blog`.`blog_auth` (`id`, `username`, `password`) VALUES (null, 'admin', '123456');