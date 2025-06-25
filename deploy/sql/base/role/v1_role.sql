create database if not exists tiktok;
use tiktok;

-- 角色表
create table if not exists roles(
      id  binary(16) not null  comment '角色表',
      code varchar(100) not null comment '角色唯一标识符',
      name varchar(100) not null  default '' comment '角色名称',
      description varchar(255) not null default '' comment '角色详细描述',

      status tinyint(1) not null default 0 comment '角色是否启用:0不启用 1启用',

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
      created_by binary(16) null comment '创建人id',   -- 允许null
      updated_by binary(16) null comment '修改人id',   -- 允许null

      primary key (id) comment '主键',

      unique index roles_udx_code (code,deleted_at_fixed) comment '保证角色标识符唯一存在',

      index roles_idx_code_deleted_status(code,is_deleted,status),

      constraint roles_chk_status check ( status in (0,1) )
)engine=InnoDB default charset=utf8mb4 collate=utf8mb4_0900_ai_ci
    COMMENT='角色表';