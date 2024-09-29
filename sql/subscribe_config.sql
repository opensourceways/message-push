create table subscribe_config
(
    id           integer default nextval('message_center.subscrbe_config_id_seq'::regclass) not null
        constraint subscrbe_config_pk
            primary key,
    source       text                                                                       not null,
    event_type   text,
    spec_version text,
    mode_filter  jsonb,
    created_at   timestamp,
    updated_at   timestamp,
    is_deleted   boolean,
    recipient_id text
);

comment on column subscribe_config.source is '消息源';

comment on column subscribe_config.event_type is '事件类型';

comment on column subscribe_config.spec_version is '版本';

comment on column subscribe_config.mode_filter is '模式过滤';

comment on column subscribe_config.recipient_id is '接收人id';


create index subscribe_config_event_index
    on subscribe_config (source, event_type, spec_version);

