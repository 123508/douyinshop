create database if not exists tiktok;
use tiktok;

-- 订单日志表
create table if not exists order_status_log(
     id  binary(16) not null comment '订单日志id',
     order_id binary(16) not null comment '订单id',
     event_time datetime(3) not null default current_timestamp(3) comment '事件发生时间',
     description varchar(2000) not null default '' comment '日志附带消息',

     -- 0待付款 1待接单 2已接单 3运输中 4待收货 5已完成 6已取消 7退款中 8已退款 9商家拒单
     old_status int unsigned not null comment '变更前状态',
     new_status int unsigned not null comment '变更后状态',

     operator_id binary(16) not null default 0 comment '操作人id(0=系统操作)',
     operation_source enum('system','user','admin') not null default 'system'comment '操作来源: system/user/admin',

     -- 仅创建时间 (日志不可变)
     created_at datetime(3) not null default current_timestamp(3) comment '创建时间',

     primary key (id),
     -- 联合索引加速查询
     index idx_order_event(order_id,event_time),    -- 订单历史查询
     index idx_old_event(old_status,event_time),    -- 状态监控(高频)
     index idx_new_event(new_status,event_time),    -- 审计分析(低频)

    -- 防止重复状态变更 (相同订单+新状态+时间戳去重)
     unique index order_status_log_udx_order_status (order_id, new_status, (UNIX_TIMESTAMP(event_time) div 10)),

     constraint order_status_log_chk_status check ((old_status between 0 and 9) and (new_status between 0 and 9))
)engine=InnoDB default charset=utf8mb4 collate=utf8mb4_0900_ai_ci
    COMMENT='订单日志表';