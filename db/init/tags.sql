create table tags
(
    tag_id    bigint not null,
    banner_id bigint not null,
    primary key (tag_id, banner_id)
);