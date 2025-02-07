-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *,
(SELECT name as creator FROM users WHERE id = user_id),
(SELECT name as feed FROM feeds WHERE id = feed_id);

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*, feeds.name as title, (SELECT name AS user FROM users WHERE users.id = $1)
FROM feed_follows
INNER JOIN feeds
ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1;

-- name: Unfollow :exec
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2;