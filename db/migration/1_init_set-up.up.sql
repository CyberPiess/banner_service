create table banners
(
    id          serial
        primary key,
    content     jsonb,
    is_active   boolean,
    create_time timestamp with time zone,
    update_time timestamp with time zone
);

alter table banners
    owner to "test_user";

insert into banners (content, is_active, create_time, update_time)
values ('{"title": "some_title", "text": "some_text", "url": "some_url"}', true, current_timestamp, null);

create table features
(
    feature_id   bigint not null,
    banner_id    bigint not null,
    primary key (feature_id, banner_id)
);

alter table features
    owner to "test_user";

insert into features (feature_id, banner_id)
values (1,1);

create table tags
(
    tag_id    bigint not null,
    banner_id bigint not null,
    primary key (tag_id, banner_id)
);

alter table tags
    owner to "test_user";

insert into tags (tag_id, banner_id)
values (1,1), (2, 1), (3, 1);

create table valid_tokens
(
    id                serial
        primary key,
    token             text,
    permission_level text
);

alter table valid_tokens
    owner to "test_user";

insert into valid_tokens(token, permission_level)
values ('user_token', 'user'), ('admin_token', 'admin');