create table refresh_tokens
(
    refresh_token varchar(200),
    user_id varchar(100) REFERENCES users(user_id),
    expires_at timestamptz,
    created_at timestamptz
);