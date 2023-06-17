-- +migrate Up
CREATE TABLE IF NOT EXISTS discord_webhooks (
    id TEXT NOT NULL PRIMARY KEY,
    webhook_id TEXT NOT NULL,
    webhook_token TEXT NOT NULL
);

-- +migrate Down
DROP TABLE discord_webhooks;