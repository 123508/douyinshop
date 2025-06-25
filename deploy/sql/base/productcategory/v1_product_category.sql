create database if not exists tiktok;
use tiktok;


-- 商品-分类表
create table if not exists product_category(

    product_id binary(16) not null comment '商品id',
    category_id binary(16) not null comment '分类id',
    status   tinyint(1) not null default 1 comment '启用状态: 0禁用 1启用',

    -- 时间戳
    created_at datetime(3) not null default current_timestamp(3) comment '创建时间',
    updated_at datetime(3) not null default current_timestamp(3) on update CURRENT_TIMESTAMP(3) comment '更新时间',

    operator_id binary(16) null  comment '操作人(null=系统)',

    primary key (product_id,category_id),
    index idx_category_active(category_id,status),
    constraint chk_status check ( status in (0,1) )
)engine=InnoDB default charset=utf8mb4 collate=utf8mb4_0900_ai_ci
    COMMENT='商品-分类表';
