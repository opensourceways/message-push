create table message_push_record
(
    recipient_id            text,
    time_uuid               timeuuid,
    created_at              timestamp,
    event_data              map<text, text>,
    event_data_content_type text,
    event_data_schema       text,
    event_id                text,
    event_source            text,
    event_source_url        text,
    event_spec_version      text,
    event_time              timestamp,
    event_type              text,
    event_user              text,
    push_address            text,
    push_state              text,
    push_time               timestamp,
    push_type               text,
    remark                  text,
    title                   text,
    summary                 text,
    primary key (recipient_id, time_uuid)
)
            with clustering order by (time_uuid desc);

create index message_push_record_source_index
    on message_push_record (event_source);

