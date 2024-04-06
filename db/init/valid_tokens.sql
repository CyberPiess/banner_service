create table valid_tokens
(
    id                serial
        primary key,
    token             text,
    perlmission_level text
);