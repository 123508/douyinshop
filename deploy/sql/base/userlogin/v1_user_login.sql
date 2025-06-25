create database if not exists tiktok;
use tiktok;

-- 用户密码表
create table if not exists user_login(
    user_id binary(16) not null comment '用户id',

    password varchar(200) not null  comment '用户加密后密码',

    -- 时间戳
    created_at datetime(3) not null default current_timestamp(3) comment '创建时间',
    updated_at datetime(3) not null default current_timestamp(3) on update CURRENT_TIMESTAMP(3) comment '更新时间',

    -- 删除三元组
    deleted_at datetime(3) default null comment '删除时间(软删除)',
    is_deleted tinyint(1) generated always as (IF(deleted_at is null, 0, 1)) virtual ,
    deleted_date date generated always as (COALESCE(DATE(deleted_at), '9999-12-31')) virtual ,

    -- 审计字段
    version  int not null default 0 comment '版本号',

    primary key (user_id),
    index user_login_idx_user_pwd(user_id,password)
)engine=InnoDB default charset=utf8mb4 collate=utf8mb4_0900_ai_ci
    COMMENT='用户密码表';