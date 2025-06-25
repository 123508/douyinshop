create database if not exists tiktok;
use tiktok;

-- 物流表
create table if not exists order_delivery(
    id              binary(16) not null default UUID_TO_BIN(uuid(),1) comment '订单物流id',
    order_id        binary(16) not null comment '订单ID',
    delivery_no     varchar(64) not null comment '运单号/快递单号',

    -- 物流公司信息
    express_code    varchar(32) not null comment '物流公司编码（如SF, YTO, JD等）',
    express_name    varchar(64) not null comment '物流公司名称',

    -- 发货信息
    sender_name     varchar(50) not null comment '发件人姓名',
    sender_phone    varchar(32) not null comment '发件人电话',
    sender_address  varchar(512) not null comment '发件人地址',

    -- 收货信息
    receiver_name   varchar(50) not null comment '收件人姓名',
    receiver_phone  varchar(32) not null comment '收件人电话',
    receiver_address varchar(512) not null comment '收件人地址',

    status          TINYINT(1) not null default 0 comment '物流状态: 0待发货 1运输中 2派送中 3已签收 4异常 5退回',
    shipped_at      datetime(3) default null comment '发货时间',
    delivered_at    datetime(3) default null comment '签收时间',
    remark          varchar(255) default null comment '备注',

    -- 时间戳
    created_at      datetime(3) not null default CURRENT_TIMESTAMP(3) comment '创建时间',
    updated_at      datetime(3) not null default CURRENT_TIMESTAMP(3) on update CURRENT_TIMESTAMP(3) comment '更新时间',


    primary key (id),
    index idx_order(order_id),
    index idx_delivery_no(delivery_no),
    index idx_express_code(express_code),
    constraint chk_status check ( status between 0 and 5)
)  engine =InnoDB charset=utf8mb4 collate=utf8mb4_0900_ai_ci
    comment='订单物流表';