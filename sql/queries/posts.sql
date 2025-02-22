-- name: CreatePost :exec
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
);

-- name: GetPostsForUser :many
SELECT p.title, p.url, p.description, p.published_at, p.created_at, p.updated_at, f.name as feed_name
FROM posts p
JOIN feeds f ON p.feed_id = f.id
WHERE feed_id IN (
    SELECT feed_id
    FROM feed_follows
    WHERE feed_follows.user_id = $1
)
ORDER BY p.published_at DESC, p.created_at DESC
LIMIT $2;
