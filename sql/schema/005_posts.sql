-- +goose Up
CREATE TABLE posts (
    id uuid PRIMARY KEY,
    created_at TIMESTAMP not null,
    updated_at TIMESTAMP not null,
    title text NOT NULL,
    url text unique not null,
    description text,
    published_at TIMESTAMP,
    feed_id uuid NOT NULL REFERENCES feeds(id) ON DELETE CASCADE 
);

-- +goose Down
DROP TABLE posts;
