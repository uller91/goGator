-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;


-- name: GetPostForUser :many
SELECT posts.*, feeds.name as feed_name from posts
join feed_follows on feed_follows.feed_id = posts.feed_id
join feeds on posts.feed_id = feeds.id
WHERE feed_follows.user_id = $1
order by published_at desc 
limit $2;