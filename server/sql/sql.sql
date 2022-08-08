-- kv表
drop table if exists `account`;
create table `account` (
    id bigint,
    `money` bigint not null default 0,
    primary key (`id`)
) engine=innodb default charset=utf8;

-- 订单表
drop table if exists `order`;
create table `order` (
    id bigint auto_increment,
    `from` bigint not null,
    `from_money` bigint not null,
    `to` bigint not null,
    `to_money` bigint not null,
    `ext` varchar(128),
    primary key (`id`)
) engine=innodb default charset=utf8;

-- 唯一key
drop table if exists `unique_key`;
create table `unique_key`(
    `key` varchar(128) not null,
    primary key (`key`)
) engine=innodb default charset=utf8;