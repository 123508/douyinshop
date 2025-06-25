create database if not exists tiktok;
use tiktok;

-- 用户日志表
create table if not exists user_login_log(
     id  binary(16) not null comment '日志id',
     user_id binary(16) not null comment '用户id',
     login_type tinyint unsigned not null  default 0 comment '用户登录方式: 0密码 1短信 2扫码 3第三方',
     login_status tinyint(1) unsigned not null default 0 comment '登录结果: 0成功  1失败',
     fail_reason varchar(255) null comment '登录失败原因',
     login_ip varchar(50) default null comment '用户登录ip',
     user_agent varchar(255) default null comment '浏览器UA',
     device varchar(64) default null comment '设备标识',
     created_at datetime(3) not null default current_timestamp(3) comment '登录时间',

     primary key (id),
     index user_login_log_idx_user(user_id)
)engine=InnoDB default charset=utf8mb4 collate=utf8mb4_0900_ai_ci
    COMMENT='用户日志表';