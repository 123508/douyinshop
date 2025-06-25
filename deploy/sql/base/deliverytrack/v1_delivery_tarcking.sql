create database if not exists tiktok;
use tiktok;

-- 物流轨迹表
create table if not exists delivery_tracking (

    id binary(16)not null  comment '主键',
    delivery_id binary(16) not null comment '物流记录ID',
    status tinyint(1) not null comment '变更后状态:0待发货 1运输中 2派送中 3已签收 4异常 5退回',
    location varchar(255) comment '当前位置',
    description varchar(500) comment '状态描述',

    created_at datetime(3) not null default CURRENT_TIMESTAMP(3) comment '创建时间',
    operator_id binary(16) null comment '操作人id(null=系统)',

    primary key (id),

    index delivery_tracking_idx_delivery (delivery_id),

    constraint delivery_tracking_chk_status check ( status between 0 and 5)

) engine =InnoDB charset=utf8mb4 collate=utf8mb4_0900_ai_ci
    comment= '物流轨迹表';