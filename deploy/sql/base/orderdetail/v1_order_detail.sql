create database if not exists tiktok;
use tiktok;

-- 订单明细表
create table if not exists order_details(
    id  binary(16) not null  comment '订单明细id',
    order_id binary(16) not null comment '订单id',
    coupon_id binary(16) null comment '优惠券id',
    promotion_id binary(16) null comment '活动id',
    order_no varchar(32) not null comment '订单号',
    product_id binary(16) not null comment '商品id',

    sku_name varchar(512) not null default ''  comment '订单商品名',
    sku_image varchar(512) not null  default '' comment '商品图片',
    sku_quantity int unsigned not null  default 0 comment '订单商品数量',
    sku_price decimal(18,2) not null comment '商品价格',
    sku_spec varchar(255) not null comment '商品规格',
    sku_unit  tinyint not null default 0 comment '商品价格单位', -- 0分 1角 2元 3美元 4欧元 5日元 6港币 7新台币 8韩元 9英镑
    sku_barcode varchar(127) null comment '商品条码',


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
    index order_details_idx_order_deleted(order_id,is_deleted),
    index order_details_idx_created(created_at),
    constraint order_details_chk_quantity check ( sku_quantity >= 0 )

)engine=InnoDB default charset=utf8mb4 collate=utf8mb4_0900_ai_ci
    COMMENT='订单明细表';