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
SELECT f.name, f.url, u.name as username
FROM feeds f
JOIN users u ON f.user_id = u.id
WHERE f.user_id = u.id;

-- name: GetFeedByUrl :one
SELECT *
FROM feeds
WHERE url = $1
LIMIT 1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $2, updated_at = $3
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT id, name, url
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
