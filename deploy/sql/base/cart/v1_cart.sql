create database if not exists tiktok;
use tiktok;

-- 购物车表
create table if not exists cart(
    id  binary(16) not null default UUID_TO_BIN(uuid(),1) comment '购物车表的唯一标识(默认使用uuid7,失败的时候兜底为有序uuid1)',
    user_id binary(16) not null comment '用户id标识',
    product_id binary(16) not null comment '产品id标识',
    quantity  int unsigned not null default 1 comment '商品数量',
    selected  tinyint(1) not null default 1 comment '商品是否被选中: 1选中 0没选中',

    -- 时间戳
    created_at datetime(3) not null default current_timestamp(3) comment '创建时间',
    updated_at datetime(3) not null default current_timestamp(3) on update CURRENT_TIMESTAMP(3) comment '更新时间',

    -- 删除三元组
    deleted_at datetime(3) default null comment '删除时间(软删除)',
    is_deleted tinyint(1) generated always as (IF(deleted_at is null, 0, 1)) virtual ,
    deleted_date date generated always as (COALESCE(DATE(deleted_at), '9999-12-31')) virtual ,

    -- 唯一约束的虚拟列
    deleted_at_fixed datetime(3) generated always as (COALESCE(deleted_at, '1970-01-01 00:00:00.000')) virtual,

    -- 分片键（抖音方案）
#     shard_key tinyint unsigned generated always as (CRC32(user_id) % 128) virtual ,

    primary key (id) comment '主键',
    unique key udx_user_product_fixed (user_id,product_id,deleted_at_fixed),
    index idx_user_status (user_id,is_deleted,selected,created_at) comment '联合索引方便查询',
    constraint check_selected check ( selected in (0,1) ),
    constraint check_quantity check ( quantity > 0 )
) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_0900_ai_ci
    COMMENT='购物车表';