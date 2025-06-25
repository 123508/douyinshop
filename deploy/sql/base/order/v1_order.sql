create database if not exists tiktok;
use tiktok;

-- 订单表  
create table if not exists orders(
        -- 订单基本信息
        id  binary(16) not null default UUID_TO_BIN(uuid(),1) comment '订单id',
        user_id binary(16) not null comment '用户id',
        shop_id binary(16) not null comment '店铺id',
        coupon_id binary(16) null comment '优惠券id',
        promotion_id binary(16) null comment '活动id',
        order_no varchar(32) not null comment '订单号',
        pay_method tinyint not null comment '支付方式',  -- 1 支付宝  2 微信  3 银行卡  4 其他
        amount decimal(18,2) not null comment '支付金额,精确到后两位',
        amount_unit tinyint not null default 0 comment '支付价格单位', -- 0分 1角 2元 3美元 4欧元 5日元 6港币 7新台币 8韩元 9英镑
        remark varchar(4000) default null comment '用户下单备注',

        -- 0待付款 1待接单 2已接单 3运输中 4待收货 5已完成 6已取消 7退款中 8已退款 9商家拒单
        final_status int unsigned not null default 0 comment '用户订单的最终状态',

        -- 时间戳
        created_at datetime(3) not null default current_timestamp(3) comment '创建时间',
        updated_at datetime(3) not null default current_timestamp(3) on update CURRENT_TIMESTAMP(3) comment '更新时间',

        -- 删除三元组
        deleted_at datetime(3) default null comment '删除时间(软删除)',
        is_deleted tinyint(1) generated always as (IF(deleted_at is null, 0, 1)) virtual ,
        deleted_date date generated always as (COALESCE(DATE(deleted_at), '9999-12-31')) virtual ,

        -- 审计字段
        version  int not null default 0 comment '版本号',

        -- 唯一约束的虚拟列
        deleted_at_fixed datetime(3) generated always as (COALESCE(deleted_at, '1970-01-01 00:00:00.000')) virtual,

        primary key (id),
        index idx_status_deleted (final_status,is_deleted),
        index idx_user_deleted(user_id,is_deleted),
        index idx_shop_deleted(shop_id,is_deleted),
        -- 保证订单号唯一存在
        unique index udx_order_no(order_no,deleted_at_fixed),
        created_by binary(16) null comment '创建人id',   -- 允许null
        updated_by binary(16) null comment '修改人id',   -- 允许null

        constraint chk_status check ( final_status between 0 and 9),
        constraint chk_pay  check ( pay_method between 1 and 4)
)engine=InnoDB default charset=utf8mb4 collate=utf8mb4_0900_ai_ci
    COMMENT='订单表';


-- 冷数据表,负责存储指定时间后的数据(配合脚本实现自动迁移)
# create table order_archive like orders;