create table users
(
    user_id    varchar(100) PRIMARY KEY,
    user_name  varchar(200),
    password   varchar(256),
    created_at timestamptz,
    updated_at timestamptz
);