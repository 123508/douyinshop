create database if not exists tiktok;
use tiktok;

-- 商品图片表
create table if not exists product_image(

      id  binary(16) not null default UUID_TO_BIN(uuid(),1) comment '商品图片id',
      product_id binary(16) not null comment '商品id',
      image_url varchar(1000) not null comment '图片url',  -- 缩短长度并NOT NULL
      sort int default 0 comment '图片顺序',

    -- 删除三元组
      deleted_at datetime(3) default null comment '删除时间(软删除)',
      is_deleted tinyint(1) generated always as (IF(deleted_at is null, 0, 1)) virtual ,
      deleted_date date generated always as (COALESCE(DATE(deleted_at), '9999-12-31')) virtual ,

      primary key (id),
      index idx_product_deleted (product_id,is_deleted)
)engine=InnoDB default charset=utf8mb4 collate=utf8mb4_0900_ai_ci
    COMMENT='商品图片表';