create database if not exists tiktok;
use tiktok;

-- 权限表
create table if not exists permission (
    id  binary(16) not null  comment '权限表的唯一标识',
    code varchar(100) not null  comment '权限唯一标识符',
    name varchar(100) not null  default '' comment '权限名称',
    description varchar(255) not null default '' comment '权限详细描述',

    -- 父级优化：允许NULL表示根节点
    parent_id binary(16) null comment '父级id,没有就置空',
    path varchar(1000) generated always as (
        case
            when parent_id is null then BIN_TO_UUID(id)
            else CONCAT_WS('/', BIN_TO_UUID(parent_id), BIN_TO_UUID(id)) -- 仅拼接父ID
        end
        ) virtual comment '权限路径',

    type enum(
        'API',      -- 接口权限
        'MENU',     -- 菜单权限
        'BUTTON',   -- 按钮权限
        'DATA',     -- 数据权限
        'FIELD',    -- 字段权限
        'MODULE',   -- 模块权限
        'FILE',     -- 文件权限
        'TASK'      -- 任务权限
        ) not null default 'API' comment '权限类型',
    resource varchar(200) not null default '' comment '权限对应资源',
    method varchar(15) not null default '' comment '权限对应方法类型',
    status tinyint(1) not null default 0 comment '权限是否启用:0不启用 1启用',

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

    -- 唯一约束优化
    unique index permission_udx_code (code,deleted_at_fixed) comment '保证权限标识符唯一存在',

    index permission_idx_code_deleted_status(code,is_deleted,status),

    -- API权限唯一约束
    unique index permission_udx_resource_method (resource, method, deleted_at_fixed)
        comment '同一资源+方法只能有一个权限',

    -- 树形结构索引
    index permission_idx_parent_id (parent_id),

    -- 类型查询优化
    index permission_idx_type (type),

    -- 启用类型约束
    constraint permission_chk_status check ( status in (0,1) ),

    -- 类型枚举约束
    constraint permission_chk_type check (type in ('API','MENU','BUTTON','DATA','FIELD','MODULE','FILE','TASK'))
)engine=InnoDB default charset=utf8mb4 collate=utf8mb4_0900_ai_ci
    COMMENT='权限表';