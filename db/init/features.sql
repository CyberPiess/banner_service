create table features
(
    feature_id   bigint not null,
    feature_name text,
    banner_id    bigint not null,
    primary key (feature_id, banner_id)
);
