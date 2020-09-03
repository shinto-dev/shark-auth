create table refresh_tokens
(
    refresh_token varchar(200) PRIMARY KEY,
    user_id       varchar(100) REFERENCES users (user_id),
    session_id    varchar(100) NOT NULL,
    expires_at    timestamptz  NOT NULL,
    created_at    timestamptz  NOT NULL,
    deleted       boolean DEFAULT false
);