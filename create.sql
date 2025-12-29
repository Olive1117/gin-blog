-- create database blog;
-- use blog;
-- source create.sql;
DROP TABLE IF EXISTS blog_tag, blog_category, blog_article, blog_login, blog_article_tag;

CREATE TABLE `blog_tag` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `created_by` bigint(20) unsigned DEFAULT '0' COMMENT '创建者ID',
  `updated_by` bigint(20) unsigned DEFAULT '0' COMMENT '修改者ID',
  `deleted_by` bigint(20) unsigned DEFAULT '0' COMMENT '删除者ID',

  `name` varchar(100) DEFAULT '' COMMENT '标签名称',
  `state` tinyint(3) unsigned DEFAULT '1' COMMENT '状态 0为禁用、1为启用',
  INDEX idx_blog_tag_deleted_at (deleted_at),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章标签管理';

CREATE TABLE `blog_category` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `created_by` bigint(20) unsigned DEFAULT '0' COMMENT '创建者ID',
  `updated_by` bigint(20) unsigned DEFAULT '0' COMMENT '修改者ID',
  `deleted_by` bigint(20) unsigned DEFAULT '0' COMMENT '删除者ID',

  `name` varchar(100) DEFAULT '' COMMENT '分类名称',
  `state` tinyint(3) unsigned DEFAULT '1' COMMENT '状态 0为禁用、1为启用',
  INDEX idx_blog_catrgory_deleted_at (deleted_at),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章分类管理';

CREATE TABLE `blog_article` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `created_by` bigint(20) unsigned DEFAULT '0' COMMENT '创建者ID',
  `updated_by` bigint(20) unsigned DEFAULT '0' COMMENT '修改者ID',
  `deleted_by` bigint(20) unsigned DEFAULT '0' COMMENT '删除者ID',

  `category_id` bigint(20) unsigned DEFAULT '0' COMMENT '分类ID',
  `title` varchar(100) DEFAULT '' COMMENT '文章标题',
  `desc` varchar(255) DEFAULT '' COMMENT '简述',
  `content` text,
  `state` tinyint(3) unsigned DEFAULT '1' COMMENT '状态 0为禁用1为启用',
  INDEX idx_blog_article_deleted_at (deleted_at),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章管理';

CREATE TABLE `blog_article_tag` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `created_by` bigint(20) unsigned DEFAULT '0' COMMENT '创建者ID',
  `updated_by` bigint(20) unsigned DEFAULT '0' COMMENT '修改者ID',
  `deleted_by` bigint(20) unsigned DEFAULT '0' COMMENT '删除者ID',

  `article_id` bigint(20) unsigned DEFAULT '0' COMMENT '文章ID',
  `tag_id` bigint(20) unsigned DEFAULT '0' COMMENT '标签ID',
  CONSTRAINT `fk_article` FOREIGN KEY (`article_id`) REFERENCES `blog_article`(`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_tag` FOREIGN KEY (`tag_id`) REFERENCES `blog_tag`(`id`) ON DELETE CASCADE,
  UNIQUE KEY `uk_article_tag` (`article_id`, `tag_id`),
  INDEX idx_blog_article_tab_deleted_at (deleted_at),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章分类管理';

CREATE TABLE `blog_login` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `created_by` bigint(20) unsigned DEFAULT '0' COMMENT '创建者ID',
  `updated_by` bigint(20) unsigned DEFAULT '0' COMMENT '修改者ID',
  `deleted_by` bigint(20) unsigned DEFAULT '0' COMMENT '删除者ID',

  `username` varchar(50) DEFAULT '' COMMENT '账号',
  `password` varchar(100) DEFAULT '' COMMENT '密码',
  INDEX idx_blog_login_deleted_at (deleted_at),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `blog`.`blog_login` (`id`, `username`, `password`) VALUES (null, 'admin', '123456');