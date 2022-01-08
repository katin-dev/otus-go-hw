-- +goose Up
CREATE TABLE events (
    id VARCHAR(255) NOT NULL,
    title VARCHAR(500) NOT NULL,
    "date" timestamp NOT NULL,
    duration INT NOT NULL,
    description VARCHAR(1000),
    user_id VARCHAR(255) NOT NULL DEFAULT '',
    notify_before INT NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE events;