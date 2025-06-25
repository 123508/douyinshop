create database if not exists tiktok;
use tiktok;

-- 商品表
create table if not exists product(

    id  binary(16) not null default UUID_TO_BIN(uuid(),1) comment '商品id',
    name varchar(512) not null comment '商品名称',
    shop_id binary(16) not null comment '店家id',
    price  decimal(18,2) not null default 0 comment '商品价格',
    status  tinyint(1) not null  default 0 comment '商品状态: 0下架  1审核中  2上架',
    sales   int unsigned not null default 0  comment '商品销量',

    -- 时间戳
    created_at datetime(3) not null default current_timestamp(3) comment '创建时间',
    updated_at datetime(3) not null default current_timestamp(3) on update CURRENT_TIMESTAMP(3) comment '更新时间',

    -- 删除三元组
    deleted_at datetime(3) default null comment '删除时间(软删除)',
    is_deleted tinyint(1) generated always as (IF(deleted_at is null, 0, 1)) virtual ,
    deleted_date date generated always as (COALESCE(DATE(deleted_at), '9999-12-31')) virtual ,

    -- 审计字段
    version  int not null default 0 comment '版本号',

    primary key (id),

    index idx_shop_deleted_status(shop_id,is_deleted,status),
    index idx_sales_deleted(sales,is_deleted),
    INDEX idx_price_shop (price,shop_id),
    INDEX idx_shop_created (shop_id, created_at),

    constraint chk_status check ( status between 0 and 2)
)engine=InnoDB default charset=utf8mb4 collate=utf8mb4_0900_ai_ci
    COMMENT='商品表';