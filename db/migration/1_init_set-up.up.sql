create table banners
(
    id          serial
        primary key,
    content     jsonb,
    is_active   boolean,
    create_time timestamp with time zone,
    update_time timestamp with time zone
);


insert into banners (content, is_active, create_time, update_time)
values ('{"title": "some_title", "text": "some_text", "url": "some_url"}', true, '0001-01-01T00:00:00Z', null);
insert into banners (content, is_active, create_time, update_time)
values('{"title": "some_title", "text": "some_text", "url": "some_url"}', false, '0001-01-01T00:00:00Z', null);
insert into banners (content, is_active, create_time, update_time)
values('{"title": "some_title", "text": "some_text", "url": "some_url"}', false, '0001-01-01T00:00:00Z', null);

create table features
(
    feature_id   bigint not null,
    banner_id    bigint not null,
    primary key (feature_id, banner_id)
);

insert into features (feature_id, banner_id)
values (1,1), (1,2), (2,3);

create table tags
(
    tag_id    bigint not null,
    banner_id bigint not null,
    primary key (tag_id, banner_id)
);


insert into tags (tag_id, banner_id)
values (1,1), (2, 1), (3, 1), (1,3), (2, 3), (3, 3), (4,2), (5, 2), (6, 2);

create table valid_tokens
(
    id                serial
        primary key,
    token             text,
    permission_level text
);


insert into valid_tokens(token, permission_level)
values ('user_token', 'user'), ('admin_token', 'admin');