-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedUrl :one
SELECT * FROM feeds WHERE url = $1;

-- name: GetFeedName :one
SELECT name FROM feeds WHERE id = $1;

-- name: MarkFeedFetched :one
Update feeds set updated_at = $1, last_fetched_at = $1 WHERE id = $2
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds order by last_fetched_at asc nulls first limit 1;