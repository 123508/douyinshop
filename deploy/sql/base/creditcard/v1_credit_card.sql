create database if not exists tiktok;
use tiktok;

-- 信用卡表
create table if not exists credit_card(
    id binary(16) not null default (UUID_TO_BIN(UUID(), 1)) comment '支付方式ID',
    user_id binary(16) not null comment '用户ID',

    -- 支付平台信息
    platform tinyint not null comment '支付平台',-- 1支付宝 2 微信  3 银行卡  4 其他
    platform_user_id varchar(128) comment '支付平台用户ID(支付宝user_id/微信openid)',
    auth_token varchar(256) not null comment '需应用层AES-256加密存储', -- 安全注释
    bank_token varchar(128) comment '支付平台返回的银行卡唯一标识',

    -- 银行卡信息
    card_last_four char(4) comment '卡号后四位',
    card_brand varchar(15) not null comment '卡品牌',
    card_type enum('CREDIT','DEBIT') not null comment '信用卡/借记卡',
    bank_name varchar(64) not null comment '银行名称',

    -- 默认设置
    is_default_card tinyint(1) not null default 0 comment '是否为默认使用卡',
    is_default_pay tinyint(1) not null default 0 comment '是否默认支付方式',

    -- 令牌管理
    token_expiry datetime(3) comment '令牌过期时间',

    -- 时间戳
    created_at datetime(3) not null default CURRENT_TIMESTAMP(3),
    updated_at datetime(3) not null default CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),

    -- 虚拟列用于唯一约束
    platform_uid varchar(128) generated always as (
         case
             when platform = 3 then bank_token
             else platform_user_id
         end
    ) virtual ,

    -- 默认支付方式虚拟列（新增）
    is_default_pay_flag tinyint generated always as(
    IF(is_default_pay, 1, null)
    ) virtual ,

    -- 删除三元组
    deleted_at datetime(3) default null comment '删除时间(软删除)',
    is_deleted tinyint(1) generated always as (IF(deleted_at is null, 0, 1)) virtual ,
    deleted_date date generated always as (COALESCE(DATE(deleted_at), '9999-12-31')) virtual ,

    -- 唯一约束的虚拟列
    deleted_at_fixed datetime(3) generated always as (COALESCE(deleted_at, '1970-01-01 00:00:00.000')) virtual,

    -- 审计字段
    version  int not null default 0 comment '版本号',
    created_by binary(16) null comment '创建人id',   -- 允许null
    updated_by binary(16) null comment '修改人id',   -- 允许null

    PRIMARY KEY (id),
    index idx_user (user_id),
    unique index uniq_user_platform (user_id, platform, platform_uid,deleted_at_fixed),
    -- 条件唯一索引（MySQL 8.0+）
    unique index uniq_bank_token (bank_token,deleted_at_fixed),
    unique index udx_user_default (user_id, is_default_card) invisible ,
    -- 默认支付方式索引（新增）
    unique index udx_user_default_pay (user_id, is_default_pay_flag) invisible ,

    -- 数据完整性约束（MySQL 8.0.16+）
    constraint chk_platform_data check (
        (platform = 3 and bank_token is not null and platform_user_id is null)
            or
        (platform in (1,2,4) and platform_user_id is not null)
    ),

    constraint chk_default check ( is_default_card in (0,1) and is_default_pay in (0,1) )

) engine =InnoDB default charset=utf8mb4 collate=utf8mb4_0900_ai_ci
    comment '用户支付方式表';