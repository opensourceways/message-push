create table cloud_event_message
(
    id                integer default nextval('message_center.message_cloud_event_id_seq'::regclass) not null
        constraint message_cloud_event_pk
            primary key,
    source            text,
    type              text,
    data_content_type text,
    data_schema       text,
    spec_version      text                                                                           not null,
    data              json,
    time              timestamp,
    created_time      timestamp,
    updated_time      timestamp
);

comment on table cloud_event_message is '清洗后的标准消息';

comment on column cloud_event_message.source is '消息源';

comment on column cloud_event_message.type is '消息类型';

comment on column cloud_event_message.data_content_type is '数据内容类型';

comment on column cloud_event_message.data_schema is '数据分类';

comment on column cloud_event_message.spec_version is '版本';

comment on column cloud_event_message.data is '数据详情';

comment on column cloud_event_message.time is '事件时间';

alter table cloud_event_message
    owner to postgres;

