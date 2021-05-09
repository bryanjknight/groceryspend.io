-- name: GetorCreateUserByAuthProviderId :one
INSERT INTO "user" (
  auth_provider_id, auth_specific_id
) VALUES (
  $1, $2
)
ON CONFLICT(auth_provider_id, auth_specific_id)
-- TODO: perf impact of doing this
DO UPDATE SET auth_provider_id=EXCLUDED.auth_provider_id
RETURNING *;



-- name: DeleteAuthor :exec
DELETE FROM "user"
WHERE id = $1;