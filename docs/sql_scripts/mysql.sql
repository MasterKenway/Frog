-- UA 黑名单配置
CREATE TABLE `t_ua_config`
(
    `id`          int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
    `ua`          varchar(512) NOT NULL COMMENT 'UA 黑名单配置',
    `insert_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '插入时间',
    `update_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `is_delete`   tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否删除 0 - 未删除 1 - 删除',
    PRIMARY KEY (`id`)
) ENGINE = Innodb AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8 COMMENT 'UA 黑名单配置';

CREATE TABLE `t_escape_middleware_ips`
(
    `id`          int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
    `ip`          varchar(15) NOT NULL COMMENT 'IP',
    `insert_time` timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '插入时间',
    `update_time` timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `is_delete`   tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否删除 0 - 未删除 1 - 删除',
    PRIMARY KEY (`id`),
    KEY           `ip` (`ip`)
) ENGINE = Innodb AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8mb4 COMMENT '忽略中间件 IP 配置';


CREATE TABLE `t_user`
(
    `id`          int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
    `uid`         varchar(128) NOT NULL COMMENT '用户 id',
    `username`    varchar(128) NOT NULL COMMENT '用户名',
    `password`    varchar(32)  NOT NULL COMMENT '密码',
    `email`       varchar(128) NOT NULL COMMENT '邮箱',
    `login_ips`   varchar(256) NULL COMMENT '登录 ip （json 数组）',
    `is_valid`    tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否允许登录',
    `insert_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '插入时间',
    `update_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY `id`,
    UNIQUE KEY `username` (`username`),
    KEY           `uid` (`uid`)
) ENGINE = Innodb AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8mb4 COMMENT '用户表';