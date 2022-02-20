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
) ENGINE = Innodb AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8 COMMENT '忽略中间件 IP 配置';

