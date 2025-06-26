create database if not exists tiktok;
use tiktok;

-- 用户表
create table if not exists users(
    id  binary(16) not null comment '用户id',
    username varchar(60) not null default '' comment '用户名称',
    nickname varchar(60) not null default '' comment '用户昵称',
    email varchar(255) not null default '' comment '用户邮箱',
    phone varchar(20) not null  comment '用户手机',
    gender tinyint(1) not null default 0 comment '用户性别:0未知  1男  2女',
    avatar varchar(1000) comment '用户头像',
    status tinyint(1) not null default 0 comment '用户状态:0正常 1冻结',

    -- 时间戳
    created_at datetime(3) not null default current_timestamp(3) comment '创建时间',
    updated_at datetime(3) not null default current_timestamp(3) on update CURRENT_TIMESTAMP(3) comment '更新时间',

    -- 删除三元组
    deleted_at datetime(3) default null comment '删除时间(软删除)',
    is_deleted tinyint(1) generated always as (IF(deleted_at is null, 0, 1)) virtual ,
    deleted_date date generated always as (COALESCE(DATE(deleted_at), '9999-12-31')) virtual ,

    -- 唯一约束的虚拟列
    deleted_at_fixed datetime(3) generated always as (COALESCE(deleted_at, '1970-01-01 00:00:00.000')) virtual,

    -- 审计字段
    version  int not null default 0 comment '版本号',
    last_updated_by binary(16) comment '最后修改人',

    primary key (id),

    unique index users_udx_email(email,deleted_at_fixed),
    unique index users_udx_phone(phone,deleted_at_fixed),
    unique index users_udx_username(username,deleted_at_fixed),

    constraint users_chk_status check ( status in(0,1) ),
    constraint users_chk_gender check ( gender in(0,1,2) )
)engine=InnoDB default charset=utf8mb4 collate=utf8mb4_0900_ai_ci
    COMMENT='用户表';