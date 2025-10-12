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
SELECT feeds.id, feeds.created_at, feeds.updated_at, feeds.name, feeds.url, feeds.user_id, users.name as user_name
FROM feeds
JOIN users 
ON feeds.user_id = users.id;

-- name: GetFeedByUrl :one
SELECT * FROM feeds 
WHERE url = $1 LIMIT 1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at = $1, last_fetched_at = $2
WHERE id = $3;
