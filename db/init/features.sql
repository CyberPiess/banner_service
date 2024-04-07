create table features
(
    feature_id   bigint not null,
    banner_id    bigint not null,
    primary key (feature_id, banner_id)
);
