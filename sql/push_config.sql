create table push_config
(
    id           serial
        constraint push_config_pk
            primary key,
    subscribe_id integer,
    push_type    text,
    push_address text,
    created_at   timestamp,
    updated_at   timestamp,
    is_deleted   boolean
);

comment on table push_config is '推送配置';

comment on column push_config.subscribe_id is '订阅id';

comment on column push_config.push_type is '推送类型';

comment on column push_config.push_address is '推送地址';


create index push_config_subcribe_index
    on push_config (subscribe_id);

