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
    `uid`         varchar(256)  NOT NULL COMMENT '用户 id',
    `username`    varchar(128)  NOT NULL COMMENT '用户名',
    `password`    varchar(32)   NOT NULL COMMENT '密码',
    `email`       varchar(128)  NOT NULL COMMENT '邮箱',
    `avatar`      varchar(1024) NOT NULL DEFAULT '' COMMENT '头像',
    `login_ips`   text NULL COMMENT '登录 ip （json 数组）',
    `is_valid`    tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否允许登录',
    `insert_time` timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '插入时间',
    `update_time` timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `is_delete`   tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否删除 0 - 未删除 1 - 删除',
    PRIMARY KEY (`id`),
    UNIQUE KEY `username` (`username`),
    KEY           `uid` (`uid`)
) ENGINE = Innodb AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8mb4 COMMENT '用户表';

CREATE TABLE `t_rental_info`
(
    `id`              int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
    `title`           varchar(128)  NOT NULL COMMENT '标题',
    `cover`           varchar(1024) NOT NULL COMMENT '封面 (COS 链接)',
    `pics`            text          NOT NULL COMMENT '房子照片 (COS 链接)',
    `area`            float         NOT NULL COMMENT '面积',
    `price`           float         NOT NULL COMMENT '价格',
    `rent_avail_time` varchar(64)   NOT NULL COMMENT '入住时间',
    `rent_term`       varchar(64)   NOT NULL COMMENT '租房周期',
    `province`        varchar(256)  NOT NULL COMMENT '省份',
    `city`            varchar(256)  NOT NULL COMMENT '城市',
    `location`        varchar(256) COMMENT ' 位置 ',
    `desc`            text          NOT NULL COMMENT ' 描述 ',
    `tags`            text          NOT NULL COMMENT ' 标签 (json 数组) ',
    `house_type`      varchar(256)  NOT NULL COMMENT ' 户型 ',
    `room_type`       varchar(1024) NOT NULL COMMENT ' 房型 (图片链接) ',
    `furniture`       text          NOT NULL COMMENT ' 家具 (json 数组) ',
    `type`            int           NOT NULL COMMENT ' 0 - 整租 1 - 合租 ',
    `rooms`           text          NOT NULL COMMENT ' 如果为合租，出租的房间 ',
    `insert_time`     timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '插入时间',
    `update_time`     timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `is_delete`       tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否删除 0 - 未删除 1 - 删除',
    PRIMARY KEY `id` (`id`),
    KEY               `price` (`price`),
    KEY               `area` (`area`),
) ENGINE = Innodb AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8mb4 COMMENT ' 租房信息表 ';


CREATE TABLE `t_comment`
(
    `id`          int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
    `rental_id`   int unsigned NOT NULL COMMENT '住房信息 ID',
    `uid`         varchar(256)  NOT NULL COMMENT '用户 uid',
    `username`    varchar(128)  NOT NULL COMMENT '用户名称',
    `avatar`      varchar(1024) NOT NULL COMMENT '用户头像',
    `content`     text          NOT NULL COMMENT '评论内容',
    `insert_time` timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '插入时间',
    `update_time` timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `is_delete`   tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否删除 0 - 未删除 1 - 删除',
    PRIMARY KEY `id` (`id`),
    KEY           `rental_id_time` (`rental_id`, `insert_time`)
)ENGINE = Innodb AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8mb4 COMMENT '租房评论';

CREATE TABLE `t_subscribe`
(
    `id`          int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
    `uid`         varchar(256) NOT NULL COMMENT '用户 uid',
    `type`        int          NOT NULL COMMENT '订阅类型',
    `rental_id`   int unsigned NULL COMMENT '住房信息 ID',
    `tag`         varchar(64) NULL COMMENT '订阅 Tag',
    `city`        varchar(64) NULL COMMENT '订阅城市',
    `location`    varchar(64) NULL COMMENT '订阅位置',
    `is_valid`    tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否取消',
    `insert_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '插入时间',
    `update_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `is_delete`   tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否删除 0 - 未删除 1 - 删除',
    PRIMARY KEY `id` (`id`),
    KEY           `uid_time` (`uid`, `insert_time`),
)ENGINE = Innodb AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8mb4 COMMENT '用户订阅租房信息';

CREATE TABLE `t_notification`
(
    `id`          int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
    `type`        int          NOT NULL COMMENT '通知类型',
    `comment_id`  int unsigned NULL `评论`,
    `rental_id`   int unsigned NULL `租房信息 ID`,
    `content`     varchar(128) NOT NULL COMMENT '通知内容',
    `is_read`     tinyint(1) NOT NULL DEFAULT 0 COMMENT '已读',
    `insert_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '插入时间',
    `update_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `is_delete`   tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否删除 0 - 未删除 1 - 删除',
    PRIMARY KEY `id` (`id`),
    KEY           `id` (`id`, `insert_time`),
)ENGINE = Innodb AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8mb4 COMMENT '用户通知表';

CREATE TABLE `l_log`
(
    `id`          int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
    `time`        timestamp    NOT NULL COMMENT '',
    `level`       varchar(16)  NOT NULL COMMENT '',
    `caller`      varchar(256) NOT NULL COMMENT '',
    `request_id`  varchar(256) NOT NULL COMMENT '',
    `message`     text         NOT NULL COMMENT '',
    `insert_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '插入时间',
    `update_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `is_delete`   tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否删除 0 - 未删除 1 - 删除',
    PRIMARY KEY `id` (`id`),
    KEY           `request_id` (`request_id`)
)ENGINE = Innodb AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8mb4 COMMENT '日志表';