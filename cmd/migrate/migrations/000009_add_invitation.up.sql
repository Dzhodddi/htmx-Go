CREATE TABLE IF NOT EXISTS user_invetetions (
    token bytea PRIMARY KEY,
    user_id bigint NOT NULL
);