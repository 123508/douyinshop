create database if not exists tiktok;
use tiktok;

-- 用户组表
create table if not exists user_group(

     id binary(16) not null default (UUID_TO_BIN(UUID(), 1)) comment '用户组ID',
     name varchar(255) not null comment '用户组名称',
     code varchar(100) not null default CONCAT('group_', REPLACE(UUID(), '-', '')) comment '用户组唯一标识符',
     status tinyint(1) not null default 0 comment '权限是否启用:0不启用 1启用',

    -- 父级优化：允许NULL表示根节点
     parent_id binary(16) null comment '父级id,没有就置空',
     path varchar(1000) generated always as (
         case
             when parent_id is null then BIN_TO_UUID(id)
             else CONCAT_WS('/', BIN_TO_UUID(parent_id), BIN_TO_UUID(id)) -- 仅拼接父ID
             end
         ) virtual comment '用户组路径',

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

     primary key (id),

     index idx_name_deleted (name,is_deleted),
     unique index udx_code_deleted (code,deleted_at_fixed) comment '确保用户组标识符唯一存在',

     constraint chk_status check ( status in ( 0,1 ) )
)engine=InnoDB default charset=utf8mb4 collate=utf8mb4_0900_ai_ci
    COMMENT='用户组表';