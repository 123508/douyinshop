create database if not exists tiktok;
use tiktok;

-- 地址表
create table if not exists address_books (
    -- 基本信息
    id  binary(16) not null default UUID_TO_BIN(uuid(),1) comment '地址表的唯一标识(默认使用uuid7,失败的时候兜底为有序uuid1)',
    user_id binary(16) not null comment '持有该地址的用户唯一标识',
    consignee varchar(60) not null default '' comment '收货人',
    gender  tinyint unsigned default 0 comment '性别:0=未知 1=男 2=女',
    phone    varchar(15) not null default '' comment '电话号码(支持国际电话号码)',

    -- 地址信息
    stress_address varchar(255) not null default '' comment '详细地址',
    country_code char(2) not null default 'CN' comment 'ISO国家代码',
    city   varchar(100) not null default '' comment '城市',
    state varchar(100) not null default '' comment '省份/州',
    post_code varchar(12) not null default '' comment '邮政编码',

    -- 自定义标签部分
    label VARCHAR(50) not null default '' comment '用户自定义标签',

    -- 时间戳
    created_at datetime(3) not null default current_timestamp(3) comment '创建时间',
    updated_at datetime(3) not null default current_timestamp(3) on update CURRENT_TIMESTAMP(3) comment '更新时间',

    -- 删除三元组
    deleted_at datetime(3) default null comment '删除时间(软删除)',
    is_deleted tinyint(1) generated always as (IF(deleted_at is null, 0, 1)) virtual ,
    deleted_date date generated always as (COALESCE(DATE(deleted_at), '9999-12-31')) virtual ,

    -- 唯一约束的虚拟列
    deleted_at_fixed datetime(3) generated always as (COALESCE(deleted_at, '1970-01-01 00:00:00.000')) virtual,

    -- 默认地址
    is_default tinyint(1) not null default 0 comment '是否为默认地址:0不是 1是',

    -- 验证字段
    verification_status ENUM('pending', 'verified', 'invalid') not null default 'pending' comment '验证状态',
    last_verified_at datetime(3) null comment '最后验证时间',

    -- 审计字段
    version  int not null default 0 comment '版本号',

    primary key (id) comment '主键',
    -- 索引优化                               ,
    index index_user_deleted_phone (user_id,is_deleted,phone) comment '建立联合索引方便查询',
    index idx_location (country_code, state, city),
    index idx_default (user_id, is_default) ,
    index idx_verification (verification_status),

    -- 唯一约束确保单个默认地址
    unique index udx_user_default (user_id, is_default,deleted_at_fixed)
        comment '确保每个用户只有一个默认地址'
        invisible , -- MySQL 8.0+ 隐藏索引
    -- 约束
    constraint chk_gender check (gender in (0,1,2)) ,
    constraint chk_default check (is_default in (0,1)) -- 字段检查
) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_0900_ai_ci
    COMMENT='地址表';
