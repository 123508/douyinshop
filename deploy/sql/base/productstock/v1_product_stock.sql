create database if not exists tiktok;
use tiktok;

-- 库存表
create table if not exists product_stock(

     product_id binary(16) not null comment '商品id',
     stock int unsigned not null default 0 comment '可用库存',

     -- 时间戳
     created_at datetime(3) not null default current_timestamp(3) comment '创建时间',
     updated_at datetime(3) not null default current_timestamp(3) on update CURRENT_TIMESTAMP(3) comment '更新时间',

     -- 审计字段
     version  int not null default 0 comment '版本号',

     primary key (product_id)
)engine=InnoDB default charset=utf8mb4 collate=utf8mb4_0900_ai_ci
    COMMENT='商品库存表';