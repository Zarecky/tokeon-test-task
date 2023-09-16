-- name: Create :one
INSERT INTO users (id, username, created_at)
	VALUES ($1, $2, now())
	RETURNING *;

-- name: Get :one
SELECT * FROM users WHERE id=$1 LIMIT 1;

-- name: Update :one
UPDATE users
	SET username=$1, updated_at=now(), deleted_at=$2
	WHERE id=$3
	RETURNING *;

