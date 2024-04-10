create table valid_tokens
(
    id                serial
        primary key,
    token             text,
    permission_level text
);