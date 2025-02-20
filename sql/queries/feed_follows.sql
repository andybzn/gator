-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT inserted_feed_follow.*,
    feeds.name AS feedname,
    users.name AS username
FROM inserted_feed_follow
JOIN users ON users.id = inserted_feed_follow.user_id
JOIN feeds on feeds.id = inserted_feed_follow.feed_id;

-- name: GetFeedFollowsForUser :many
SELECT *
FROM feed_follows
JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;

-- name: RemoveFeedFollow :exec
DELETE FROM feed_follows
WHERE feed_id = $1
    AND user_id = $2;
