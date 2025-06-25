create database if not exists tiktok;
use tiktok;

-- 店铺表
create table if not exists shops(
    id  binary(16) not null default UUID_TO_BIN(uuid(),1) comment '店铺id',
    name varchar(255) not null comment '店铺名称',
    user_id binary(16) not null comment '店主id',
    type tinyint(1) not null default 0 comment '店铺类型: 0普通 1旗舰 2品牌',
    status tinyint(1) not null default 1 comment '状态: 0禁用 1启用',
    description varchar(2000) not null default '' comment '店铺描述',
    logo_url varchar(1000) not null comment '店铺logo的url',
    license_url varchar(255) not null default '' comment '经营许可证',
    score	decimal(3,2) comment '店铺评分',
    contact_phone varchar(32) comment '店铺电话',

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

    index idx_name_deleted_status(name,is_deleted,status),

    constraint chk_type check ( type in (0,1,2) ),
    constraint chk_status check ( status in (0,1) )
)engine=InnoDB default charset=utf8mb4 collate=utf8mb4_0900_ai_ci
    COMMENT='店铺表';