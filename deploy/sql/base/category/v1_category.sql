create database if not exists tiktok;
use tiktok;

-- 分类表
create table if not exists category(
      id  binary(16) not null comment '分类表id',
      name varchar(100) not null  default '' comment '分类名称',
      status tinyint(1) not null default 0 comment '分类是否启用:0不启用 1启用',

      -- 时间戳
      created_at datetime(3) not null default current_timestamp(3) comment '创建时间',
      updated_at datetime(3) not null default current_timestamp(3)  on update CURRENT_TIMESTAMP(3) comment '更新时间',

      -- 删除三元组
      deleted_at datetime(3) default null comment '删除时间(软删除)',
      is_deleted tinyint(1) generated always as (IF(deleted_at is null, 0, 1)) virtual ,
      deleted_date date generated always as (COALESCE(DATE(deleted_at), '9999-12-31')) virtual ,

      -- 唯一约束的虚拟列
      deleted_at_fixed datetime(3) generated always as (COALESCE(deleted_at, '1970-01-01 00:00:00.000')) virtual,

      -- 审计字段
      version  int not null default 0 comment '版本号',

      primary key (id) comment '主键',

      unique index category_udx_name_deleted(name,is_deleted,deleted_at_fixed),

      -- 约束
      constraint category_chk_status check ( status in (0,1) )
)engine=InnoDB default charset=utf8mb4 collate=utf8mb4_0900_ai_ci
    comment = '分类表'
