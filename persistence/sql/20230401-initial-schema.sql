-- +migrate Up
CREATE TABLE IF NOT EXISTS discord_messages (
    id TEXT NOT NULL PRIMARY KEY,
    message_id TEXT NOT NULL
);

-- +migrate Down
DROP TABLE discord_messages;