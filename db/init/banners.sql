create table banners
(
    id          serial
        primary key,
    content     jsonb,
    is_active   boolean,
    create_time timestamp with time zone,
    update_time timestamp with time zone
);
